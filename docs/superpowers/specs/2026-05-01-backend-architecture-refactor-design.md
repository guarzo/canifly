# Backend Architecture Refactor — Design Spec

**Date:** 2026-05-01
**Status:** Approved
**Scope:** Plan A of three (Plan B: frontend module cleanup; Plan C: dependency refresh)

## Goal

Break up the 1,065-line `EVEDataServiceImpl` god service into five focused services, eliminate the construction cycle that requires `nil`-then-setter patching, classify startup errors as required vs optional, move Fuzzworks bootstrap off the critical path, and split the oversized `auth.go` handler — without changing any handler signature, HTTP route, or external behavior.

## Findings addressed

This spec resolves findings #1, #2, #3, #4, #6, #7, #8 from the 2026-05-01 codebase review (`improve_findings.md`). Findings #5, #9 (frontend) are Plan B; finding #10 (deps) is Plan C.

## Current state (problem)

`internal/services/eve/eve_data_service_impl.go` (1,065 lines, 52 methods) implements five interfaces:
- `ESIAPIService` (corp/alliance/user-info ESI passthroughs)
- `CharacterService` (10 methods: GetCharacter, ProcessIdentity, RefreshCharacterData, …)
- `SkillPlanService` (CRUD + parser + 200+ lines of plan-vs-character evaluation)
- `ProfileService` (LoadCharacterSettings, BackupDir, SyncDir, SyncAllDir)
- `CacheableService` (Get/Set/SaveCache/LoadCache — disk-backed)

`internal/server/services.go` constructs it with `nil` HTTP client and `nil` account service, then calls `SetHTTPClient` and `SetAccountManagementService` post-hoc. The cycle exists because the HTTP client uses `eveDataService` as its cache, while the EVE service needs the HTTP client to make ESI calls.

`internal/services/interfaces/eve_data_composite.go` self-documents: `// This is a temporary interface for migration purposes`. The migration was started ~8 months ago and frozen.

`UpdateCharacterFields(characterID int64, updates map[string]interface{})` is stringly-typed; callers pass magic key names; the implementation does runtime type assertions.

Startup (`services.go`) silently downgrades all init errors to warnings, including a misleading `// Continue with embedded data as fallback` comment for Fuzzworks (the `embed/` package was removed — there is no fallback).

`internal/handlers/auth.go` is 566 lines mixing OAuth flow, session, and login-state polling.

## Target state (solution)

### Package layout after refactor

```
internal/services/
  cache/
    http_cache.go            (existing, unchanged)
    persistent_cache.go      (NEW)
  character/
    character_service.go     (NEW)
  skillplan/
    skillplan_service.go     (NEW)
  profile/
    profile_service.go       (NEW)
  eve/
    esi_client.go            (NEW — corp/alliance/user-info passthroughs)
    eve_data_service_impl.go (DELETED at the end)
  interfaces/
    cacheable.go             (kept; impl now PersistentCacheService)
    character.go             (kept)
    skill_plan.go            (kept)
    profile.go               (kept)
    esi_api.go               (kept)
    eve_data_composite.go    (DELETED)
    eve_data.go              (DELETED — aggregate interface)
```

### Construction order (no cycle, no setters)

```
1. storage           = NewStorageService(...)
2. config            = NewConfigurationService(storage, ...)
3. persistentCache   = NewPersistentCacheService(storage, logger)
4. authClient        = NewDynamicAuthClient(config, ...)
5. httpClient        = NewEsiHttpClient(authClient, persistentCache, logger)
6. characterService  = NewCharacterService(httpClient, repos, persistentCache, logger)
7. skillplanService  = NewSkillPlanService(httpClient, skillRepo, characterService, logger)
8. profileService    = NewProfileService(eveProfileRepo, config, logger)
9. esiClient         = NewESIClient(httpClient, persistentCache, logger)
10. accountMgmt      = NewAccountManagementService(storage, esiClient, characterService, authClient, logger)
11. syncService      = NewSyncService(profileService, characterService, config, logger)
```

No nil arguments. No `Set*` methods. Each service receives only what it actually needs.

### `AppServices` struct (after)

```go
type AppServices struct {
    StorageService           interfaces.StorageService
    AccountManagementService interfaces.AccountManagementService
    ConfigurationService     interfaces.ConfigurationService

    // Split EVE Services — each is now a distinct concrete type
    ESIAPIService    interfaces.ESIAPIService     // *eve.ESIClient
    CharacterService interfaces.CharacterService  // *character.Service
    SkillPlanService interfaces.SkillPlanService  // *skillplan.Service
    ProfileService   interfaces.ProfileService    // *profile.Service
    CacheableService interfaces.CacheableService  // *cache.PersistentCacheService

    SyncService      interfaces.SyncService
    LoginService     interfaces.LoginService
    AuthClient       interfaces.AuthClient
    HTTPCacheService interfaces.HTTPCacheService
    WebSocketHub     *handlers.WebSocketHub
}
```

The `EVEDataService` field is removed. Verify no callers in `cmd/` or anywhere outside `server/` read it before deletion.

### CharacterUpdate type (replaces map-based API)

`internal/model/character_update.go`:

```go
package model

// CharacterUpdate is a partial update for a character. Nil pointers are no-ops.
type CharacterUpdate struct {
    Role          *string
    MCT           *bool
    Training      *bool
    SkillPlanName *string
    // …only fields used by current callers; add as needed
}
```

`CharacterService` exposes:

```go
UpdateCharacter(characterID int64, update model.CharacterUpdate) error
```

The `interfaces.CharacterService` interface is updated; the old `UpdateCharacterFields(id, map[string]interface{})` is deleted, not deprecated.

### Startup error classification (services.go)

Each initialization step in `GetServices` is documented as **required** (return error, app exits) or **optional** (warn-and-continue). Header comment block:

```go
// Initialization steps fall into two categories:
//
// REQUIRED — failure aborts startup:
//   * storage directories
//   * configuration service load
//   * skill repo (skill plans + skill types)
//   * system repo
//   * auth client
//
// OPTIONAL — failure is logged and startup continues:
//   * EVE credentials (user can set via UI)
//   * Fuzzworks initial download (runs async; cached data used until ready)
//   * Settings directory creation (best-effort; recreated on demand)
```

The misleading `// Continue with embedded data as fallback` comment is replaced with `// Async Fuzzworks update will populate data when complete; existing cached data is used in the meantime.`

### Async Fuzzworks bootstrap

`fuzzworksService.Initialize(ctx)` moves into a goroutine started during `GetServices` but does not block return:

```go
if autoUpdate {
    fuzzworksService := fuzzworks.New(logger, cfg.BasePath, false)
    go func() {
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
        defer cancel()
        webSocketHub.Broadcast("fuzzworks:status", map[string]string{"state": "updating"})
        if err := fuzzworksService.Initialize(ctx); err != nil {
            logger.Errorf("Fuzzworks update failed: %v", err)
            webSocketHub.Broadcast("fuzzworks:status", map[string]string{"state": "error", "error": err.Error()})
            return
        }
        webSocketHub.Broadcast("fuzzworks:status", map[string]string{"state": "ready"})
    }()
}
```

WebSocketHub message format follows the existing convention used elsewhere in the codebase (verify in `internal/handlers/websocket.go` during implementation).

### auth.go split

Pure file-split, no logic changes:

```
internal/handlers/auth.go         (DELETED)
internal/handlers/auth_oauth.go   (NEW — login, callback, refresh handlers)
internal/handlers/auth_session.go (NEW — logout, session check, login state polling)
```

Both files share the existing `AuthHandler` struct (defined in `auth_oauth.go` since OAuth is the primary concern). The constructor `NewAuthHandler` stays in `auth_oauth.go`.

## Migration order (PR-by-PR)

Each step ends with all tests green and is independently revertible. Behavior changes are isolated to specific tasks (#7 async Fuzzworks, #6 init hardening) — all other tasks are pure refactors.

1. **Extract `PersistentCacheService`** — moves Get/Set/SaveCache/LoadCache into `services/cache/persistent_cache.go`. EVE service keeps thin delegating shims (still satisfies `CacheableService`). HTTP client now depends on the cache directly. **The cycle is broken here.**
2. **Extract `ProfileService`** — smallest service (4 methods); validates the extraction pattern.
3. **Extract `CharacterService`** + introduce `CharacterUpdate` struct. Removes the stringly-typed map.
4. **Extract `SkillPlanService`** — includes the parser, CRUD, and the 200+ line plan-vs-character evaluation logic.
5. **Extract `ESIClient`** — corp/alliance/user-info passthroughs.
6. **Delete the corpse** — remove `EVEDataServiceImpl`, `EVEDataComposite`, the `EVEDataService` interface, the aggregate field on `AppServices`, and migration comments. After this step, `interfaces/` contains only the narrow interfaces.
7. **Startup error classification** — annotate `services.go` and convert silent warnings to fatal errors for required steps.
8. **Async Fuzzworks bootstrap** — move into goroutine + WebSocket progress.
9. **Split `auth.go`** into `auth_oauth.go` + `auth_session.go`.

## Risks and mitigations

| Risk | Mitigation |
|---|---|
| Plan-evaluation tests are tightly coupled to `EVEDataServiceImpl` | Move tests in the same commit as the extraction; assertions stay identical, only imports change |
| `cmd/` or external code reads `appServices.EVEDataService` | Grep before Task 6; if found, update those call sites in the same PR that deletes the field |
| `testutil/mock_interfaces.go` is hand-written, not generated | Update mocks alongside each interface change; verify the mock methods match the trimmed interfaces |
| Async Fuzzworks could race with handlers reading the same data | Verify the storage layer (`SkillStore`, `SystemStore`) already uses appropriate locking; if not, that becomes a Task 8 sub-step |
| `WebSocketHub` API may not support the message shape proposed | Verify in `internal/handlers/websocket.go` during Task 8; adjust to match existing convention |

## Out of scope

- Frontend changes (Plan B)
- Dependency upgrades (Plan C)
- New features
- Test framework changes
- Logging refactor
- Splitting the SkillPlan evaluation logic from CRUD (rejected during brainstorming as larger than necessary)

## Acceptance criteria

- `go test -race ./...` passes at the end of every task
- `go vet ./...` clean
- `internal/services/eve/eve_data_service_impl.go` is deleted
- `internal/services/interfaces/eve_data_composite.go` and `eve_data.go` are deleted
- `internal/server/services.go` contains no `Set*` calls and no `nil` arguments to constructors
- `internal/handlers/auth.go` does not exist; `auth_oauth.go` and `auth_session.go` together contain the same handlers
- `UpdateCharacterFields` and its `map[string]interface{}` parameter no longer appear in the codebase (grep returns zero matches)
- App starts within 2 seconds (Fuzzworks no longer blocks startup)
- All API routes in `internal/server/router.go` respond identically to before (manual smoke test of login → character list → skill plan check → settings sync)
