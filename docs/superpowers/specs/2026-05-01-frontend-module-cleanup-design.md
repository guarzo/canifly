# Frontend Module Cleanup — Design Spec

**Date:** 2026-05-01
**Status:** Approved
**Scope:** Plan B of three (Plan A: backend refactor; Plan C: dependency refresh)

## Goal

Split the 557-line `apiService.jsx` god-module into per-backend-domain API clients, break the two oversized page/component files (`CharacterOverview.jsx` 491, `Header.jsx` 394) into a custom hook + focused subcomponents, and wire up the Fuzzworks status badge that Plan A introduced. No new features.

## Findings addressed

Findings #5 and #9 from the 2026-05-01 codebase review.

## Sequencing

Plan B can run in parallel with Plan A or Plan C — it touches only `renderer/`. Plan A's Task 8 (async Fuzzworks) emits `fuzzworks:status` WebSocket messages; this plan's Header refactor wires the receiving badge. If Plan A hasn't merged yet when Plan B reaches the Header task, the badge component is still built but its hook gracefully no-ops on absent messages.

## Current state

`renderer/src/api/apiService.jsx` (557 lines) exports every API call the frontend makes — accounts, skill plans, sync, config, EVE data, fuzzworks, auth — in one file. Mirrors the backend god-service problem. `apiService.test.jsx` (322 lines) tests the lot.

`renderer/src/pages/CharacterOverview.jsx` (491 lines) and `renderer/src/components/common/Header.jsx` (394 lines) each combine data fetching, derived state, modal state, and presentation in one file.

## Target state

### API client layout

```
renderer/src/api/
  apiClient.js              (NEW — shared axios instance, base URL, interceptors)
  apiClient.test.js         (NEW — tests for client config + interceptors)
  accountsApi.js            (NEW — /api/accounts/*)
  configApi.js              (NEW — /api/config/*)
  esiApi.js                 (NEW — /api/esi/*)
  skillPlansApi.js          (NEW — /api/skill-plans/*)
  fuzzworksApi.js           (NEW — /api/fuzzworks/*)
  syncApi.js                (NEW — /api/sync/*)
  authApi.js                (NEW — auth endpoints + login state)
  apiService.jsx            (DELETED at end of plan)
  apiService.test.jsx       (DELETED — assertions split across the new test files)
  apiRequest.test.jsx       (kept; renamed to apiClient.test.js if it tests the same thing — verify during implementation)
```

`apiClient.js` owns the configured axios instance (or the existing fetch wrapper — verify), the base URL constant `http://localhost:42423`, and any cross-cutting interceptors currently in `apiService.jsx`. Each domain file imports the client and exports named functions only. No default exports.

### CharacterOverview.jsx layout

```
renderer/src/pages/
  CharacterOverview.jsx                  (≤150 lines — composition only)

renderer/src/hooks/
  useCharacterOverview.js                (NEW, ≤100 lines — data fetching + derived state)

renderer/src/components/character-overview/
  CharacterOverviewHeader.jsx            (NEW, ≤200 lines)
  CharacterOverviewTable.jsx             (NEW, ≤200 lines)
  CharacterOverviewFilters.jsx           (NEW, ≤200 lines)
  CharacterOverviewModals.jsx            (NEW, ≤200 lines — modal state container)
```

The hook owns all `useEffect`, all data-loading `useState`, and any derivation memos. Subcomponents are stateless; they receive props and callbacks from the page. Subcomponent count is a target, not a contract — split by visual section in the current file.

### Header.jsx layout

```
renderer/src/components/common/
  Header.jsx                  (≤150 lines — composition only)
  HeaderNav.jsx               (NEW, ≤200 lines — nav links, route highlight)
  HeaderUserMenu.jsx          (NEW, ≤200 lines — user dropdown, logout)
  HeaderStatusBadges.jsx      (NEW, ≤200 lines — Fuzzworks status, ws connection state)

renderer/src/hooks/
  useHeaderData.js            (NEW, ≤100 lines — websocket subs, derived status)
  useFuzzworksStatus.js       (NEW, ≤50 lines — listens for fuzzworks:status messages)
```

`useFuzzworksStatus` subscribes to the WebSocket hub and exposes `{state: 'idle'|'updating'|'ready'|'error', error?: string}`. `HeaderStatusBadges` reads it and renders an MUI Chip in the appropriate color.

## Migration order

1. **Add `apiClient.js`** with the shared instance and interceptors. `apiService.jsx` re-exports from it temporarily so call sites compile. Tests exist for the client.
2. **Split per-domain API files** — one PR per domain: accounts → skillPlans → sync → config → esi → fuzzworks → auth. Each PR moves one domain's exports out of `apiService.jsx`, updates callers, and moves the relevant test cases.
3. **Delete `apiService.jsx`** and `apiService.test.jsx` once empty.
4. **Refactor `CharacterOverview.jsx`** — extract `useCharacterOverview.js` first (purely move state + effects), then extract subcomponents one at a time. The page file shrinks with each step.
5. **Refactor `Header.jsx`** — extract `useHeaderData.js` and `useFuzzworksStatus.js`, then subcomponents (Nav, UserMenu, StatusBadges).
6. **Wire the Fuzzworks badge** — `HeaderStatusBadges` consumes `useFuzzworksStatus` and renders the chip. Verify with Plan A's backend if available; otherwise the chip stays in `idle` state until backend lands.

## Risks and mitigations

| Risk | Mitigation |
|---|---|
| Vitest test imports break across many files | Codemod-style search-and-replace per domain; run `npm run test:react` after each split |
| Hidden coupling: shared module-level state in `apiService.jsx` | Read the file first; lift any module state into `apiClient.js` or eliminate it |
| Subcomponents need shared state | Lift shared state into the hook; subcomponents receive props/callbacks |
| Visual regression in Character Overview / Header | Manual smoke test after Tasks 4 and 5 — open Character Overview, open Header user menu, click each nav link, verify badges |
| Fuzzworks badge depends on backend that may not be merged | Hook gracefully renders `idle` when no messages arrive; works whether or not Plan A is merged |

## Out of scope

- TypeScript migration (separate, larger discussion)
- UI/visual changes beyond the new Fuzzworks chip
- Adding new tests where coverage was missing pre-refactor (existing tests preserved; new tests required only for the new hooks and the API client)
- State-library changes (Zustand store under `renderer/src/stores/` continues as-is)
- Refactoring other large files (`Sync.jsx`, `CharacterTable.jsx`, etc.) — explicitly deferred

## Acceptance criteria

- `apiService.jsx` and `apiService.test.jsx` do not exist
- Each new domain API file is ≤200 lines
- `CharacterOverview.jsx` and `Header.jsx` are each ≤150 lines (composition only)
- Each new subcomponent is ≤200 lines; each new hook is ≤100 lines
- `npm run test:react` passes
- `npm run lint` passes
- New tests exist for `apiClient.js`, `useCharacterOverview.js`, `useHeaderData.js`, `useFuzzworksStatus.js`
- Manual smoke test: load app, navigate Character Overview, open at least one modal, click each header nav item, log out — all behave identically to pre-refactor
- With Plan A merged: launching the app shows the Fuzzworks chip in `updating` then `ready` state during initial data load
