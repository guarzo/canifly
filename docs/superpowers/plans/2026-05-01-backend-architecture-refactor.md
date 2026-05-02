# Backend Architecture Refactor Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Break up `EVEDataServiceImpl` (1,065 LOC, 52 methods, 5 interfaces) into five focused services, eliminate the construction cycle, harden startup error handling, async the Fuzzworks bootstrap, and split `auth.go`. No behavior changes except the explicit ones in Tasks 7–8.

**Architecture:** Five new packages under `internal/services/` (`cache`, `character`, `skillplan`, `profile`, plus `eve/esi_client.go`). Construction in `services.go` becomes linear with no `nil` arguments and no `Set*` calls. Handlers don't change — they already consume the narrow interfaces.

**Tech Stack:** Go 1.23.2, existing `interfaces.Logger`, existing `internal/persist` adapters, existing `WebSocketHub`.

**Spec:** `docs/superpowers/specs/2026-05-01-backend-architecture-refactor-design.md`

---

## Pre-flight: branch + baseline

### Task 0: Create the working branch and verify green baseline

**Files:** none

- [ ] **Step 1: Create branch from main**

```bash
git checkout main && git pull
git checkout -b refactor/backend-architecture
```

- [ ] **Step 2: Run baseline test suite and capture coverage**

```bash
go test -race -timeout 10m -coverprofile=baseline-coverage.out ./...
go tool cover -func=baseline-coverage.out | tail -1
```

Expected: all tests pass. Note the total coverage % — Task 6 verifies the refactor doesn't regress it.

- [ ] **Step 3: Confirm `go vet` is clean**

```bash
go vet ./...
```

Expected: no output.

---

## Task 1: Extract PersistentCacheService (breaks the cycle)

**Files:**
- Create: `internal/services/cache/persistent_cache.go`
- Create: `internal/services/cache/persistent_cache_test.go`
- Modify: `internal/services/eve/eve_data_service_impl.go` (lines around 937–1010 — Get/Set/SaveCache/LoadCache/cleanupExpired)
- Modify: `internal/server/services.go` (construction order)
- Modify: `internal/http/esi_http_client.go` (now takes cache directly, not via EVE service)

- [ ] **Step 1: Read current cache implementation**

```bash
sed -n '930,1060p' internal/services/eve/eve_data_service_impl.go
```

Note the dependencies: it uses `s.storage.LoadAPICache()`, `s.storage.SaveAPICache()`, and an in-memory map with TTL. Confirm there are no calls to `s.httpClient` or other EVE-service-only state.

- [ ] **Step 2: Write the failing test for the new service**

Create `internal/services/cache/persistent_cache_test.go`:

```go
package cache

import (
    "testing"
    "time"

    "github.com/guarzo/canifly/internal/testutil"
)

func TestPersistentCache_SetAndGet(t *testing.T) {
    storage := testutil.NewMockStorageService()
    logger := testutil.NewMockLogger()
    c := NewPersistentCacheService(storage, logger)

    c.Set("key1", []byte("value1"), 5*time.Minute)
    got, ok := c.Get("key1")
    if !ok || string(got) != "value1" {
        t.Fatalf("expected value1, got %q ok=%v", got, ok)
    }
}

func TestPersistentCache_Expires(t *testing.T) {
    storage := testutil.NewMockStorageService()
    logger := testutil.NewMockLogger()
    c := NewPersistentCacheService(storage, logger)

    c.Set("key1", []byte("value1"), 1*time.Millisecond)
    time.Sleep(5 * time.Millisecond)
    if _, ok := c.Get("key1"); ok {
        t.Fatalf("expected expired entry to miss")
    }
}

func TestPersistentCache_SaveLoadRoundTrip(t *testing.T) {
    storage := testutil.NewMockStorageService()
    logger := testutil.NewMockLogger()
    c := NewPersistentCacheService(storage, logger)

    c.Set("key1", []byte("value1"), 5*time.Minute)
    if err := c.SaveCache(); err != nil {
        t.Fatalf("SaveCache: %v", err)
    }

    c2 := NewPersistentCacheService(storage, logger)
    if err := c2.LoadCache(); err != nil {
        t.Fatalf("LoadCache: %v", err)
    }
    got, ok := c2.Get("key1")
    if !ok || string(got) != "value1" {
        t.Fatalf("round-trip lost value: got %q ok=%v", got, ok)
    }
}
```

- [ ] **Step 3: Run the test to confirm it fails**

```bash
go test ./internal/services/cache/ -run TestPersistentCache -v
```

Expected: build failure — `NewPersistentCacheService` undefined.

- [ ] **Step 4: Implement PersistentCacheService**

Create `internal/services/cache/persistent_cache.go`:

```go
package cache

import (
    "sync"
    "time"

    "github.com/guarzo/canifly/internal/services/interfaces"
)

type persistentEntry struct {
    Value      []byte    `json:"value"`
    Expiration time.Time `json:"expiration"`
}

type PersistentCacheService struct {
    storage interfaces.StorageService
    logger  interfaces.Logger
    mu      sync.RWMutex
    entries map[string]persistentEntry
}

func NewPersistentCacheService(storage interfaces.StorageService, logger interfaces.Logger) *PersistentCacheService {
    return &PersistentCacheService{
        storage: storage,
        logger:  logger,
        entries: make(map[string]persistentEntry),
    }
}

func (c *PersistentCacheService) Get(key string) ([]byte, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    e, ok := c.entries[key]
    if !ok || time.Now().After(e.Expiration) {
        return nil, false
    }
    return e.Value, true
}

func (c *PersistentCacheService) Set(key string, value []byte, ttl time.Duration) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.entries[key] = persistentEntry{Value: value, Expiration: time.Now().Add(ttl)}
}

func (c *PersistentCacheService) SaveCache() error {
    c.mu.RLock()
    snapshot := make(map[string]persistentEntry, len(c.entries))
    for k, v := range c.entries {
        snapshot[k] = v
    }
    c.mu.RUnlock()
    return c.storage.SaveAPICache(snapshot)
}

func (c *PersistentCacheService) LoadCache() error {
    loaded, err := c.storage.LoadAPICache()
    if err != nil {
        return err
    }
    c.mu.Lock()
    defer c.mu.Unlock()
    if loaded == nil {
        c.entries = make(map[string]persistentEntry)
        return nil
    }
    // Adapt the loaded shape to our internal type.
    // The storage layer currently returns map[string]persistentEntry-equivalent;
    // verify the actual signature in internal/services/storage and adjust this conversion.
    if typed, ok := loaded.(map[string]persistentEntry); ok {
        c.entries = typed
        return nil
    }
    c.logger.Warnf("LoadAPICache returned unexpected type %T", loaded)
    c.entries = make(map[string]persistentEntry)
    return nil
}

// SaveEsiCache is an alias kept for the CacheableService interface.
func (c *PersistentCacheService) SaveEsiCache() error { return c.SaveCache() }

// Compile-time check.
var _ interfaces.CacheableService = (*PersistentCacheService)(nil)
```

> **Note:** the actual `SaveAPICache` / `LoadAPICache` signature must be matched. Read `internal/services/interfaces/storage.go` and `internal/services/storage/storage_service.go` first; if the type differs from `map[string]persistentEntry`, adapt the conversion. Do not change the storage signature in this task.

- [ ] **Step 5: Run cache tests to confirm they pass**

```bash
go test ./internal/services/cache/ -run TestPersistentCache -v
```

Expected: PASS.

- [ ] **Step 6: Update HTTP client to take cache directly**

In `internal/http/esi_http_client.go`, change the constructor signature. The current call site is:
```go
httpClient := http.NewEsiHttpClient("https://esi.evetech.net", logger, authClient, eveDataService)
```

The fourth parameter is typed as `interfaces.CacheableService` (or similar). Keep that interface — `*PersistentCacheService` satisfies it. No code change inside the HTTP client.

- [ ] **Step 7: Wire PersistentCacheService into services.go**

Edit `internal/server/services.go`. Replace the existing construction sequence so the cache is built before the HTTP client and EVE service:

```go
storageService := storage.NewStorageService(cfg.BasePath, logger)
if err := storageService.EnsureDirectories(); err != nil {
    return nil, fmt.Errorf("failed to ensure directories: %w", err)
}

configurationService := configSvc.NewConfigurationService(storageService, logger, cfg.BasePath, cfg.SecretKey)

// NEW: persistent cache, owned independently of the EVE service
persistentCache := cacheSvc.NewPersistentCacheService(storageService, logger)
if err := persistentCache.LoadCache(); err != nil {
    logger.Warnf("failed to load persistent cache: %v", err)
}

// ...existing credential/autoUpdate logic unchanged...

httpClient := http.NewEsiHttpClient("https://esi.evetech.net", logger, authClient, persistentCache)

// EVE service no longer owns the cache; it delegates to persistentCache for any
// remaining cache calls (Task 6 will eliminate it entirely).
eveDataService := eveSvc.NewEVEDataServiceImpl(
    logger, httpClient, authClient, nil, configurationService,
    storageService, skillRepo, systemRepo, eveProfileRepo, persistentCache,
)
// ...rest unchanged for now...
```

- [ ] **Step 8: Update EVEDataServiceImpl constructor and delegate cache methods**

Add `persistentCache interfaces.CacheableService` to `EVEDataServiceImpl` and `NewEVEDataServiceImpl`. Replace the bodies of `Get`, `Set`, `SaveCache`, `LoadCache`, `SaveEsiCache` with delegations:

```go
func (s *EVEDataServiceImpl) Get(key string) ([]byte, bool)              { return s.persistentCache.Get(key) }
func (s *EVEDataServiceImpl) Set(key string, v []byte, ttl time.Duration){ s.persistentCache.Set(key, v, ttl) }
func (s *EVEDataServiceImpl) SaveCache() error                           { return s.persistentCache.SaveCache() }
func (s *EVEDataServiceImpl) LoadCache() error                           { return s.persistentCache.LoadCache() }
func (s *EVEDataServiceImpl) SaveEsiCache() error                        { return s.persistentCache.SaveEsiCache() }
```

Delete the `cleanupExpired`, `Shutdown`, `Clear` methods on EVE service if they only managed the cache. Move `cleanupExpired` into PersistentCacheService if it's needed (run as a goroutine in `NewPersistentCacheService`).

- [ ] **Step 9: Run full test suite**

```bash
go test -race ./...
go vet ./...
```

Expected: PASS / no warnings.

- [ ] **Step 10: Commit**

```bash
git add internal/services/cache/persistent_cache.go \
        internal/services/cache/persistent_cache_test.go \
        internal/services/eve/eve_data_service_impl.go \
        internal/server/services.go \
        internal/http/esi_http_client.go
git commit -m "refactor(backend): extract PersistentCacheService, break construction cycle

The HTTP client now depends on PersistentCacheService directly instead of
on EVEDataServiceImpl. EVE service still implements CacheableService via
delegation; Task 6 removes that delegation along with the rest of the
god service."
```

---

## Task 2: Extract ProfileService

**Files:**
- Create: `internal/services/profile/profile_service.go`
- Create: `internal/services/profile/profile_service_test.go`
- Modify: `internal/services/eve/eve_data_service_impl.go` (delete LoadCharacterSettings, BackupDir, SyncDir, SyncAllDir methods after extraction)
- Modify: `internal/server/services.go` (construct profile service, route ProfileService field to it)
- Modify: `internal/services/sync/sync_service.go` (depend on ProfileService instead of EVEDataService)

- [ ] **Step 1: Read the four methods being moved**

```bash
sed -n '816,935p' internal/services/eve/eve_data_service_impl.go
```

Note dependencies. Likely: `s.eveProfileRepo`, `s.configurationService`, `s.logger`. If it depends on cache or HTTP client, flag — that suggests behavior crossing a boundary that this extraction didn't expect.

- [ ] **Step 2: Write the new service**

Create `internal/services/profile/profile_service.go`:

```go
package profile

import (
    "github.com/guarzo/canifly/internal/model"
    "github.com/guarzo/canifly/internal/services/interfaces"
)

type Service struct {
    profileRepo interfaces.EveProfilesRepository  // verify exact interface name
    config      interfaces.ConfigurationService
    logger      interfaces.Logger
}

func NewService(profileRepo interfaces.EveProfilesRepository, config interfaces.ConfigurationService, logger interfaces.Logger) *Service {
    return &Service{profileRepo: profileRepo, config: config, logger: logger}
}

func (s *Service) LoadCharacterSettings() ([]model.EveProfile, error) {
    // Body copied verbatim from EVEDataServiceImpl.LoadCharacterSettings.
    // Replace `s.` references that targeted EVE service fields with the equivalent
    // Service fields above.
}

func (s *Service) BackupDir(targetDir, backupDir string) error {
    // Body copied verbatim.
}

func (s *Service) SyncDir(subDir, charId, userId string) (int, int, error) {
    // Body copied verbatim.
}

func (s *Service) SyncAllDir(baseSubDir, charId, userId string) (int, int, error) {
    // Body copied verbatim.
}

var _ interfaces.ProfileService = (*Service)(nil)
```

> **Implementer note:** copy each method body byte-for-byte. The only changes are the receiver type and any field references that need to be repointed. Do not refactor the bodies in this task.

- [ ] **Step 3: Move tests**

If existing profile tests live in `internal/services/eve/`, move them to `internal/services/profile/profile_service_test.go` and update imports. If profile methods have no tests, that's a finding for a follow-up — note it in the PR description but don't add tests in this task.

- [ ] **Step 4: Run tests for the new package**

```bash
go test -race ./internal/services/profile/...
```

Expected: PASS.

- [ ] **Step 5: Wire profile service in services.go**

```go
profileService := profile.NewService(eveProfileRepo, configurationService, logger)
```

Set `ProfileService: profileService` in the returned `AppServices`. Remove the EVE service from satisfying `ProfileService` (it still does until Task 6, but the field now points at the new service).

- [ ] **Step 6: Update sync service dependency**

In `internal/services/sync/sync_service.go`, change the constructor to accept `interfaces.ProfileService` where it currently accepts the EVE service for profile methods. Verify by reading the existing constructor first.

- [ ] **Step 7: Delete the four methods from EVEDataServiceImpl**

Delete `LoadCharacterSettings`, `BackupDir`, `SyncDir`, `SyncAllDir` and the `var _ interfaces.ProfileService = (*EVEDataServiceImpl)(nil)` assertion line.

- [ ] **Step 8: Run full test suite**

```bash
go test -race ./...
go vet ./...
```

Expected: PASS.

- [ ] **Step 9: Commit**

```bash
git add internal/services/profile/ internal/services/eve/eve_data_service_impl.go \
        internal/server/services.go internal/services/sync/sync_service.go
git commit -m "refactor(backend): extract ProfileService into services/profile package"
```

---

## Task 3: Introduce CharacterUpdate struct + extract CharacterService

**Files:**
- Create: `internal/model/character_update.go`
- Create: `internal/services/character/character_service.go`
- Create: `internal/services/character/character_service_test.go`
- Modify: `internal/services/interfaces/character.go` (replace UpdateCharacterFields signature)
- Modify: `internal/services/eve/eve_data_service_impl.go` (delete character methods after extraction)
- Modify: `internal/server/services.go` (wire CharacterService)
- Modify: `internal/handlers/character.go`, `internal/handlers/auth.go` (call sites of UpdateCharacterFields)
- Modify: `internal/testutil/mock_interfaces.go` (mock signature)

- [ ] **Step 1: Identify all call sites of UpdateCharacterFields**

```bash
grep -rn "UpdateCharacterFields" --include="*.go" .
```

For each call site, note which keys are passed in the map. This determines the fields of `CharacterUpdate`.

- [ ] **Step 2: Create CharacterUpdate model**

Create `internal/model/character_update.go`:

```go
package model

// CharacterUpdate is a partial update for a character.
// Nil pointer fields are ignored; non-nil fields are written.
type CharacterUpdate struct {
    Role          *string
    MCT           *bool
    Training      *bool
    SkillPlanName *string
    // Add additional fields here as call-site analysis requires.
}
```

> **Implementer note:** the fields above are the expected shape based on the spec. Replace this list with whatever keys actually appear in the map literals from Step 1. Do not include fields no caller sets.

- [ ] **Step 3: Update CharacterService interface**

Edit `internal/services/interfaces/character.go`:

```go
package interfaces

import "github.com/guarzo/canifly/internal/model"

type CharacterService interface {
    ProcessIdentity(charIdentity *model.CharacterIdentity) (*model.CharacterIdentity, error)
    DoesCharacterExist(characterID int64) (bool, *model.CharacterIdentity, error)
    UpdateCharacter(characterID int64, update model.CharacterUpdate) error
    RemoveCharacter(characterID int64) error
    RefreshCharacterData(characterID int64) (bool, error)
}
```

(Note: `UpdateCharacterFields` → `UpdateCharacter`.)

- [ ] **Step 4: Write the failing tests for character.Service**

Create `internal/services/character/character_service_test.go` with at minimum:
- `TestService_UpdateCharacter_AppliesNonNilFields` — passes an update with `Role` set, verifies the underlying repo receives that value
- `TestService_UpdateCharacter_IgnoresNilFields` — passes an update with all-nil pointers, verifies no write occurs
- `TestService_DoesCharacterExist_NotFound` — verifies the false/nil/nil path

Use `testutil` mocks for the repo and HTTP client. Each test ≤25 lines.

```go
package character_test

import (
    "testing"

    "github.com/guarzo/canifly/internal/model"
    "github.com/guarzo/canifly/internal/services/character"
    "github.com/guarzo/canifly/internal/testutil"
)

func TestService_UpdateCharacter_AppliesNonNilFields(t *testing.T) {
    repo := testutil.NewMockAccountRepo()
    svc := character.NewService(testutil.NewMockHTTPClient(), repo, testutil.NewMockAuthClient(), testutil.NewMockLogger())
    role := "leader"
    err := svc.UpdateCharacter(123, model.CharacterUpdate{Role: &role})
    if err != nil { t.Fatalf("UpdateCharacter: %v", err) }
    if got := repo.LastUpdate(); got.Role == nil || *got.Role != "leader" {
        t.Fatalf("expected role=leader, got %+v", got)
    }
}
```

> **Note:** mock helper names (`NewMockAccountRepo`, `LastUpdate`, etc.) are illustrative — use whatever the existing `testutil` package exposes; if not present, add a minimal helper in this commit.

- [ ] **Step 5: Run the tests to confirm they fail**

```bash
go test ./internal/services/character/...
```

Expected: build failure — `character` package does not exist.

- [ ] **Step 6: Implement character.Service**

Create `internal/services/character/character_service.go`. Move all 10 character methods from `EVEDataServiceImpl` verbatim. The new struct:

```go
package character

import (
    "github.com/guarzo/canifly/internal/model"
    "github.com/guarzo/canifly/internal/services/interfaces"
)

type Service struct {
    httpClient   interfaces.EsiHttpClient
    repo         interfaces.AccountRepository  // verify exact name
    authClient   interfaces.AuthClient
    logger       interfaces.Logger
}

func NewService(http interfaces.EsiHttpClient, repo interfaces.AccountRepository, auth interfaces.AuthClient, logger interfaces.Logger) *Service {
    return &Service{httpClient: http, repo: repo, authClient: auth, logger: logger}
}

// ...ProcessIdentity, DoesCharacterExist, UpdateCharacter, RemoveCharacter, RefreshCharacterData...
// (Plus the ESI passthroughs that belong here: GetCharacter, GetCharacterSkills,
//  GetCharacterSkillQueue, GetCharacterLocation, ResolveCharacterNames,
//  isCharacterTraining — all moved verbatim.)

var _ interfaces.CharacterService = (*Service)(nil)
```

Implement `UpdateCharacter` to walk the non-nil pointer fields and write each via the repo. The previous map-based implementation likely did the same with type assertions; the new version replaces those assertions with `if u.Role != nil { ... }` branches.

- [ ] **Step 7: Update the mock**

In `internal/testutil/mock_interfaces.go`, replace the `UpdateCharacterFields` method on the CharacterService mock with `UpdateCharacter` matching the new signature.

- [ ] **Step 8: Update call sites**

For each call site found in Step 1, replace:
```go
svc.UpdateCharacterFields(id, map[string]interface{}{"role": "leader"})
```
with:
```go
role := "leader"
svc.UpdateCharacter(id, model.CharacterUpdate{Role: &role})
```

- [ ] **Step 9: Wire character service in services.go**

```go
characterService := character.NewService(httpClient, accountRepo, authClient, logger)
```

Update `accountManagementService` constructor to take `characterService` if it currently took the EVE service for character operations.

- [ ] **Step 10: Delete character methods from EVEDataServiceImpl**

Delete the 10 character methods and the `var _ interfaces.CharacterService = (*EVEDataServiceImpl)(nil)` assertion. Delete `var _ interfaces.ESIAPIService = (*EVEDataServiceImpl)(nil)` lines that overlap (the remaining ESI methods stay until Task 5).

- [ ] **Step 11: Run full test suite**

```bash
go test -race ./...
go vet ./...
```

Expected: PASS.

- [ ] **Step 12: Verify no `UpdateCharacterFields` remains**

```bash
grep -rn "UpdateCharacterFields" --include="*.go" .
```

Expected: zero matches.

- [ ] **Step 13: Commit**

```bash
git add internal/model/character_update.go \
        internal/services/character/ \
        internal/services/interfaces/character.go \
        internal/services/eve/eve_data_service_impl.go \
        internal/server/services.go \
        internal/handlers/character.go internal/handlers/auth.go \
        internal/testutil/mock_interfaces.go
git commit -m "refactor(backend): extract CharacterService, replace UpdateCharacterFields(map) with typed CharacterUpdate"
```

---

## Task 4: Extract SkillPlanService

**Files:**
- Create: `internal/services/skillplan/skillplan_service.go`
- Create: `internal/services/skillplan/skillplan_service_test.go`
- Move: existing skill-plan tests from `internal/services/eve/` to `internal/services/skillplan/`
- Modify: `internal/services/eve/eve_data_service_impl.go` (delete skill-plan methods)
- Modify: `internal/server/services.go` (wire skillplan service)
- Modify: `internal/handlers/skillplan.go` (no signature change — handler depends on `interfaces.SkillPlanService` already)

- [ ] **Step 1: Identify methods to move**

Methods to move (from `eve_data_service_impl.go` line numbers in the file as of Task 3 completion):
- `GetSkillPlans`, `GetSkillName`, `GetSkillTypes`, `CheckIfDuplicatePlan`, `ParseAndSaveSkillPlan`, `parseSkillPlan`, `GetSkillPlanFile`, `DeleteSkillPlan`, `ListSkillPlans`, `RefreshRemotePlans`, `GetSkillTypeByID`
- Evaluation logic: `GetPlanAndConversionData`, `initializeUpdatedPlans`, `initializeEveConversions`, `processAccountsAndCharacters`, `mapCharacterSkills`, `mapSkillQueueLevels`, `ensureCharacterMaps`, `evaluatePlanForCharacter`, `getStatusString`, `updatePlanAndCharacterStatus`, `updateEveConversionsWithSkillNames`

Verify the list with:
```bash
grep "^func (s \*EVEDataServiceImpl)" internal/services/eve/eve_data_service_impl.go
```

- [ ] **Step 2: Create the package and move tests first**

```bash
mkdir -p internal/services/skillplan
git mv internal/services/eve/eve_data_service_impl_test.go internal/services/skillplan/skillplan_service_test.go
```

(If the test file is named differently, adjust. If skill-plan tests live in a separate file, move that one specifically.)

Edit the moved test file: change `package eve` → `package skillplan` (or `_test` variant), update imports, and update type names from `EVEDataServiceImpl` → `*Service`.

- [ ] **Step 3: Run the moved tests to confirm they fail**

```bash
go test ./internal/services/skillplan/
```

Expected: build failure — `Service` undefined.

- [ ] **Step 4: Create skillplan.Service**

Create `internal/services/skillplan/skillplan_service.go`:

```go
package skillplan

import (
    "github.com/guarzo/canifly/internal/model"
    "github.com/guarzo/canifly/internal/services/interfaces"
)

type Service struct {
    skillRepo  interfaces.SkillRepository  // verify exact name
    httpClient interfaces.EsiHttpClient
    logger     interfaces.Logger
}

func NewService(skillRepo interfaces.SkillRepository, httpClient interfaces.EsiHttpClient, logger interfaces.Logger) *Service {
    return &Service{skillRepo: skillRepo, httpClient: httpClient, logger: logger}
}

// Move all skill-plan methods here verbatim.

var _ interfaces.SkillPlanService = (*Service)(nil)
```

Move every method body listed in Step 1. Repoint `s.skillRepo`, `s.httpClient`, `s.logger` references — those are the only fields the new Service has. If a method references something else (e.g. `s.characterService`), that's an unexpected coupling — flag and decide whether to inject `interfaces.CharacterService` into Service or pass needed data as a parameter.

- [ ] **Step 5: Run skillplan tests to confirm they pass**

```bash
go test -race ./internal/services/skillplan/...
```

Expected: PASS.

- [ ] **Step 6: Wire in services.go**

```go
skillPlanService := skillplan.NewService(skillRepo, httpClient, logger)
```

Set `SkillPlanService: skillPlanService` in `AppServices`.

- [ ] **Step 7: Delete skill-plan methods from EVEDataServiceImpl**

Delete every method moved in Step 4 from `eve_data_service_impl.go`. Delete the `var _ interfaces.SkillPlanService = (*EVEDataServiceImpl)(nil)` assertion.

- [ ] **Step 8: Run full test suite**

```bash
go test -race ./...
go vet ./...
```

Expected: PASS.

- [ ] **Step 9: Commit**

```bash
git add internal/services/skillplan/ \
        internal/services/eve/eve_data_service_impl.go \
        internal/server/services.go
git commit -m "refactor(backend): extract SkillPlanService into services/skillplan package"
```

---

## Task 5: Extract ESIClient

**Files:**
- Create: `internal/services/eve/esi_client.go`
- Create: `internal/services/eve/esi_client_test.go`
- Modify: `internal/services/eve/eve_data_service_impl.go` (delete the 4 ESI methods)
- Modify: `internal/server/services.go` (construct ESIClient)

Methods to move: `GetUserInfo`, `GetCorporation`, `GetAlliance`, plus any remaining ESI passthroughs not moved with CharacterService.

- [ ] **Step 1: Confirm method list**

```bash
grep "^func (s \*EVEDataServiceImpl)" internal/services/eve/eve_data_service_impl.go
```

Expected: only 4–5 methods remain (the ESI passthroughs and possibly the cache delegations from Task 1).

- [ ] **Step 2: Create ESIClient**

`internal/services/eve/esi_client.go`:

```go
package eve

import (
    "golang.org/x/oauth2"

    "github.com/guarzo/canifly/internal/model"
    "github.com/guarzo/canifly/internal/services/interfaces"
)

type ESIClient struct {
    httpClient interfaces.EsiHttpClient
    cache      interfaces.CacheableService
    logger     interfaces.Logger
}

func NewESIClient(http interfaces.EsiHttpClient, cache interfaces.CacheableService, logger interfaces.Logger) *ESIClient {
    return &ESIClient{httpClient: http, cache: cache, logger: logger}
}

func (c *ESIClient) GetUserInfo(token *oauth2.Token) (*model.UserInfoResponse, error) {
    // Body copied verbatim from EVEDataServiceImpl.
}

func (c *ESIClient) GetCorporation(id int64, token *oauth2.Token) (*model.Corporation, error) {
    // Body copied verbatim.
}

func (c *ESIClient) GetAlliance(id int64, token *oauth2.Token) (*model.Alliance, error) {
    // Body copied verbatim.
}

var _ interfaces.ESIAPIService = (*ESIClient)(nil)
```

- [ ] **Step 3: Add a smoke test**

`internal/services/eve/esi_client_test.go`:

```go
package eve_test

import (
    "testing"

    "github.com/guarzo/canifly/internal/services/eve"
    "github.com/guarzo/canifly/internal/testutil"
)

func TestNewESIClient_SatisfiesInterface(t *testing.T) {
    c := eve.NewESIClient(testutil.NewMockHTTPClient(), testutil.NewMockCache(), testutil.NewMockLogger())
    if c == nil {
        t.Fatal("expected non-nil client")
    }
}
```

(Behavior tests for these methods, if they exist on the EVE service, move alongside. If they don't exist, the smoke test plus the compile-time interface assertion is the bar for this task.)

- [ ] **Step 4: Wire in services.go**

```go
esiClient := eve.NewESIClient(httpClient, persistentCache, logger)
```

Set `ESIAPIService: esiClient` in `AppServices`.

- [ ] **Step 5: Delete the moved methods from EVEDataServiceImpl**

Delete `GetUserInfo`, `GetCorporation`, `GetAlliance`, any other ESI passthroughs moved.

- [ ] **Step 6: Run full test suite**

```bash
go test -race ./...
go vet ./...
```

Expected: PASS.

- [ ] **Step 7: Commit**

```bash
git add internal/services/eve/esi_client.go internal/services/eve/esi_client_test.go \
        internal/services/eve/eve_data_service_impl.go \
        internal/server/services.go
git commit -m "refactor(backend): extract ESIClient for corp/alliance/user-info passthroughs"
```

---

## Task 6: Delete the corpse

**Files:**
- Delete: `internal/services/eve/eve_data_service_impl.go`
- Delete: `internal/services/eve/eve_data_service_impl_test.go` (if any tests remain — should already be moved)
- Delete: `internal/services/interfaces/eve_data_composite.go`
- Delete: `internal/services/interfaces/eve_data.go` (the aggregate `EVEDataService` interface)
- Modify: `internal/server/services.go` (remove `EVEDataService` field, remove `Set*` calls, remove `nil` constructor args)
- Modify: `internal/testutil/mock_interfaces.go` (remove EVEDataService composite mock if any)

- [ ] **Step 1: Verify nothing outside the deletion targets references EVEDataService or EVEDataServiceImpl**

```bash
grep -rn "EVEDataServiceImpl\|EVEDataComposite\|interfaces\.EVEDataService\b" --include="*.go" .
```

Expected: matches only in the files being deleted/modified above. If anything else uses the aggregate interface, refactor that caller now to use the narrow interface it actually needs.

- [ ] **Step 2: Verify cmd/ does not read appServices.EVEDataService**

```bash
grep -rn "EVEDataService" --include="*.go" internal/cmd/ cmd/ 2>/dev/null
```

If found, replace with the appropriate narrow field (`CharacterService`, `SkillPlanService`, etc.).

- [ ] **Step 3: Delete the files**

```bash
git rm internal/services/eve/eve_data_service_impl.go
git rm -f internal/services/eve/eve_data_service_impl_test.go 2>/dev/null || true
git rm internal/services/interfaces/eve_data_composite.go
git rm internal/services/interfaces/eve_data.go
```

- [ ] **Step 4: Update AppServices and GetServices**

In `internal/server/services.go`:

- Remove the `EVEDataService interfaces.EVEDataService` field.
- Remove the `// Composite service for backward compatibility (temporary)` comment block.
- Remove the `EVEDataService: eveDataService` line from the return.
- Remove `eveDataService.SetHTTPClient(httpClient)` and `eveDataService.SetAccountManagementService(...)` calls.
- Remove the `eveDataService` local variable entirely.
- Verify the construction order is now linear: storage → config → cache → http → character → skillplan → profile → esi → accountManagement → sync.

The function should no longer contain the words `nil,` or `Set*` in service construction.

- [ ] **Step 5: Update mocks**

In `internal/testutil/mock_interfaces.go`, remove any composite mock for the deleted interfaces. Keep narrow-interface mocks.

- [ ] **Step 6: Run full test suite + coverage check**

```bash
go test -race -coverprofile=post-refactor-coverage.out ./...
go tool cover -func=post-refactor-coverage.out | tail -1
go vet ./...
```

Expected: tests PASS. Coverage % should be within 1 percentage point of the baseline captured in Task 0; significant drops mean tests were lost during a move.

- [ ] **Step 7: Verify cleanliness**

```bash
grep -rn "EVEDataServiceImpl\|EVEDataComposite\|migration purposes" --include="*.go" .
```

Expected: zero matches.

- [ ] **Step 8: Commit**

```bash
git add -A
git commit -m "refactor(backend): delete EVEDataServiceImpl and aggregate interfaces

The migration started by introducing EVEDataComposite is now complete.
All five responsibilities live in their own packages:
  - services/cache (PersistentCacheService)
  - services/character (Service)
  - services/skillplan (Service)
  - services/profile (Service)
  - services/eve (ESIClient)

services.go construction is now linear with no nil arguments and no
post-hoc Set* calls."
```

---

## Task 7: Startup error classification

**Files:**
- Modify: `internal/server/services.go` (annotate each init step; convert silent warnings to fatal errors where appropriate)

- [ ] **Step 1: Add the classification header comment**

At the top of `GetServices`, before the first call, add:

```go
// Initialization steps fall into two categories:
//
// REQUIRED — failure aborts startup:
//   * storage directories
//   * configuration service load
//   * skill repo (skill plans + skill types)
//   * system repo
//   * auth client construction
//
// OPTIONAL — failure is logged and startup continues:
//   * EVE credentials (user can set via UI)
//   * Fuzzworks initial download (Task 8 makes this async; cached data used in the meantime)
//   * Settings directory creation (best-effort; recreated on demand)
//   * Persistent cache load (cold start is fine)
```

- [ ] **Step 2: Convert required steps from warn-and-continue to return**

Audit each `logger.Warnf(...)` and `logger.Errorf(...)` followed by no return. For each:
- If the step is in the REQUIRED list → change to `return nil, fmt.Errorf("...: %w", err)`.
- If OPTIONAL → keep the warn but reword the comment to make the intent clear.

Specifically required (based on the current code):

```go
// Already required — no change:
if err := storageService.EnsureDirectories(); err != nil { return nil, ... }
if err := skillRepo.LoadSkillPlans(); err != nil { return nil, ... }
if err := skillRepo.LoadSkillTypes(); err != nil { return nil, ... }
if err := systemRepo.LoadSystems(); err != nil { return nil, ... }

// Currently warns; should remain a warn (configData has a sensible default):
configData, err := configurationService.FetchConfigData()
if err != nil {
    logger.Warnf("config data unavailable; using defaults: %v", err)
    configData = &model.ConfigData{}
}

// Currently warns silently; should remain warn but with a clearer comment:
if err := configurationService.EnsureSettingsDir(); err != nil {
    // Settings dir is recreated on demand by the sync flow; failure here is non-fatal.
    logger.Warnf("settings dir not pre-created: %v", err)
}
```

- [ ] **Step 3: Replace the misleading Fuzzworks comment**

Find `// Continue with embedded data as fallback` and replace with:

```go
// Existing cached data is used; the next successful Fuzzworks update will
// refresh it. Task 8 of the backend refactor makes this async; for now it
// runs synchronously during startup.
```

- [ ] **Step 4: Run tests**

```bash
go test -race ./...
go vet ./...
```

Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/server/services.go
git commit -m "refactor(backend): classify startup errors as required vs optional"
```

---

## Task 8: Async Fuzzworks bootstrap

**Files:**
- Modify: `internal/server/services.go`
- Modify: `internal/handlers/websocket.go` (verify Broadcast API; add helper if needed)

- [ ] **Step 1: Read the current WebSocketHub API**

```bash
sed -n '1,60p' internal/handlers/websocket.go
```

Identify the broadcast method (likely `Broadcast(eventType string, payload interface{})` or similar). If it doesn't exist or has a different shape, adjust Step 3 to match.

- [ ] **Step 2: Move WebSocketHub construction earlier in services.go**

Currently `webSocketHub` is constructed near the end of `GetServices`. Move its construction up so it's available when Fuzzworks starts:

```go
webSocketHub := handlers.NewWebSocketHub(logger)
// ...later, when wiring AppServices, just reference it...
```

- [ ] **Step 3: Replace synchronous Fuzzworks initialize with a goroutine**

Replace the existing block:

```go
if autoUpdate {
    logger.Infof("Initializing Fuzzworks data service (auto-update enabled)...")
    fuzzworksService := fuzzworks.New(logger, cfg.BasePath, false)
    ctx := context.Background()
    if err := fuzzworksService.Initialize(ctx); err != nil {
        logger.Errorf("Failed to initialize Fuzzworks service: %v", err)
    }
} else {
    logger.Infof("Fuzzworks auto-update disabled in configuration")
}
```

with:

```go
if autoUpdate {
    logger.Infof("Starting Fuzzworks data update in background...")
    fuzzworksService := fuzzworks.New(logger, cfg.BasePath, false)
    go func() {
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
        defer cancel()
        webSocketHub.Broadcast("fuzzworks:status", map[string]string{"state": "updating"})
        if err := fuzzworksService.Initialize(ctx); err != nil {
            logger.Errorf("Fuzzworks update failed: %v", err)
            webSocketHub.Broadcast("fuzzworks:status", map[string]string{
                "state": "error",
                "error": err.Error(),
            })
            return
        }
        logger.Infof("Fuzzworks update complete")
        webSocketHub.Broadcast("fuzzworks:status", map[string]string{"state": "ready"})
    }()
} else {
    logger.Infof("Fuzzworks auto-update disabled in configuration")
}
```

If `WebSocketHub.Broadcast` doesn't exist, add it as a thin wrapper around the existing message-send mechanism — match the conventions already used by other broadcasts.

- [ ] **Step 4: Add startup-time race guard**

If any service reads Fuzzworks-derived data at startup (verify by re-running the storage repo loads), the goroutine MUST not invalidate state those repos have already loaded. The existing flow already loads `skillRepo` and `systemRepo` from disk before Fuzzworks runs, so the goroutine only refreshes the on-disk files for the next launch. Confirm this assumption holds; if not, add an `sync.RWMutex` to the affected store.

- [ ] **Step 5: Run tests**

```bash
go test -race ./...
go vet ./...
```

Expected: PASS.

- [ ] **Step 6: Manual smoke test**

```bash
npm start
```

Expected: app window opens within ~2 seconds (previously up to ~30s on first run when Fuzzworks downloads). Frontend network panel shows a `fuzzworks:status` websocket message with state=updating, then ready.

If the frontend does not yet handle the message, that's fine — Plan B can wire the UI badge later. The contract is one-way for now.

- [ ] **Step 7: Commit**

```bash
git add internal/server/services.go internal/handlers/websocket.go
git commit -m "refactor(backend): run Fuzzworks update asynchronously, broadcast progress"
```

---

## Task 9: Split auth.go

**Files:**
- Create: `internal/handlers/auth_oauth.go`
- Create: `internal/handlers/auth_session.go`
- Delete: `internal/handlers/auth.go`

- [ ] **Step 1: List the handlers in auth.go**

```bash
grep -n "^func " internal/handlers/auth.go
```

Categorize each:
- **OAuth**: `Login`, `Callback`, `RefreshToken`, helper auth-flow funcs, `NewAuthHandler` constructor, `AuthHandler` struct
- **Session**: `Logout`, `CheckLogin`, login-state polling endpoints, helper session funcs

- [ ] **Step 2: Create auth_oauth.go with the OAuth half**

`internal/handlers/auth_oauth.go` contains:
- `package handlers` declaration
- imports needed for OAuth handlers
- `AuthHandler` struct definition
- `NewAuthHandler` constructor
- All OAuth handlers (`Login`, `Callback`, `RefreshToken`, …)
- Helper functions used only by OAuth handlers

Move the code verbatim from `auth.go`. Do not refactor.

- [ ] **Step 3: Create auth_session.go with the session half**

`internal/handlers/auth_session.go` contains:
- `package handlers` declaration
- imports needed for session handlers
- All session handlers (`Logout`, `CheckLogin`, …)
- Helper functions used only by session handlers

- [ ] **Step 4: Delete auth.go**

```bash
git rm internal/handlers/auth.go
```

- [ ] **Step 5: Compile and test**

```bash
go build ./...
go test -race ./...
go vet ./...
```

Expected: PASS. If a helper is referenced from both files but defined in only one, leave it where it has the most callers and the other file imports it (same package — no import statement needed).

- [ ] **Step 6: File-size check**

```bash
wc -l internal/handlers/auth_oauth.go internal/handlers/auth_session.go
```

Expected: both files under 400 lines (sum ≈ 566).

- [ ] **Step 7: Manual smoke test**

```bash
npm start
```

Verify: log in completes, log out completes, page reload preserves session. Routes hit by these flows (`/login`, `/callback`, `/logout`, `/api/login-state`) must all behave identically to before.

- [ ] **Step 8: Commit**

```bash
git add internal/handlers/auth_oauth.go internal/handlers/auth_session.go
git rm internal/handlers/auth.go 2>/dev/null || true
git commit -m "refactor(handlers): split 566-line auth.go into auth_oauth.go and auth_session.go"
```

---

## Final verification

### Task 10: End-to-end smoke + push

- [ ] **Step 1: Full test suite + coverage**

```bash
go test -race -coverprofile=final-coverage.out ./...
go tool cover -func=final-coverage.out | tail -1
go vet ./...
```

Expected: PASS, coverage within 1pp of the Task 0 baseline.

- [ ] **Step 2: Frontend tests still pass (sanity check)**

```bash
npm run test:react
```

Expected: PASS. The refactor was backend-only; if anything fails, a handler signature was changed accidentally.

- [ ] **Step 3: Manual smoke test**

```bash
npm start
```

Run through:
- OAuth login
- Character list loads
- Open one skill plan and verify the qualification status shows
- Open Sync page; trigger a backup
- Log out

Each step should behave identically to pre-refactor.

- [ ] **Step 4: Verify acceptance criteria**

```bash
test ! -f internal/services/eve/eve_data_service_impl.go && echo OK || echo MISSING
test ! -f internal/services/interfaces/eve_data_composite.go && echo OK || echo MISSING
test ! -f internal/services/interfaces/eve_data.go && echo OK || echo MISSING
test ! -f internal/handlers/auth.go && echo OK || echo MISSING
grep -rn "UpdateCharacterFields\|EVEDataServiceImpl\|EVEDataComposite" --include="*.go" . | grep -v _test.go
```

Expected: all four `OK`, last grep returns nothing.

- [ ] **Step 5: Push and open PR**

```bash
git push -u origin refactor/backend-architecture
gh pr create --title "Backend architecture refactor: split EVEDataServiceImpl, async Fuzzworks, split auth.go" --body "$(cat <<'EOF'
## Summary
- Splits the 1,065-line EVEDataServiceImpl into PersistentCacheService, character.Service, skillplan.Service, profile.Service, and eve.ESIClient
- Eliminates the construction cycle in services.go (no more nil args + Set* calls)
- Replaces UpdateCharacterFields(map) with a typed CharacterUpdate struct
- Classifies startup errors as required vs optional
- Moves Fuzzworks update off the critical path (async + WebSocket progress)
- Splits 566-line auth.go into auth_oauth.go + auth_session.go

Closes findings #1, #2, #3, #4, #6, #7, #8 from the 2026-05-01 codebase review.

## Test plan
- [ ] Full Go test suite passes
- [ ] go vet clean
- [ ] Frontend tests unchanged
- [ ] Manual smoke: login → character list → skill plan → sync → logout
- [ ] App window opens in <2s (verify Fuzzworks no longer blocks startup)
EOF
)"
```

