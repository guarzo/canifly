# Post-Redesign Cleanup — Design

**Date:** 2026-05-05
**Status:** Approved (pending implementation plan)
**Predecessor:** `2026-05-04-canifly-ui-redesign.md` (PR #26, merged)

## Goal

Tie up three loose ends that surfaced during PR #26 review:

1. **`Profiles.isDefaultDir`** is locally guessed via `useState(false)` and never reflects whether the user is actually on the default Tranquility settings directory. The "Reset to default" button shows on first load even when there's nothing to reset.
2. **`PlanList`** still has an inline `<img>` for the plan icon instead of using the `EveTypeIcon` primitive introduced in PR #26.
3. **PropTypes are inert in React 19** but 46 source files still declare them and `prop-types` is still a dependency. The codebase pays for runtime validation it doesn't get.

All three are small, independent, and can ship in one PR with three commits for granular revert.

## Non-goals

- Migrating the codebase to TypeScript. That's a separate, larger effort and is parked for now.
- Reworking how config flows from backend to frontend. We add one boolean field; we do not restructure the config payload.
- Touching components outside the three areas above.

---

## Section A — `Profiles.isDefaultDir` from backend

### Backend changes (Go)

**File:** `internal/services/config/configuration_service.go`

Add a public method to the configuration service:

```go
// IsDefaultSettingsDir reports whether the currently configured settings
// directory matches the OS-default Tranquility location.
func (s *ConfigurationService) IsDefaultSettingsDir() (bool, error) {
    current, err := s.GetSettingsDir()
    if err != nil {
        return false, err
    }
    if current == "" {
        return false, nil
    }
    return current == s.getDefaultSettingsDir(), nil
}
```

**File(s) carrying the config response struct** (handler / serializer for `GET /api/config` or whichever endpoint returns the config consumed by `useAppData`):

Add `IsDefaultDir bool` to the response struct with JSON tag `isDefaultDir`. Populate it from the new method.

Tests: extend `configuration_service_test.go` to cover the new method (default case, custom dir case, empty case).

### Frontend changes (React)

**`renderer/src/Routes.jsx`** — pull the new field alongside the existing fields, pass it as a prop to `Profiles`:

```diff
 const currentSettingsDir = config?.settingsDir || config?.SettingsDir || '';
+const isDefaultDir = Boolean(config?.isDefaultDir ?? config?.IsDefaultDir ?? false);
```

```diff
 <Profiles
     subDirs={eveProfiles}
     associations={associations}
     settingsData={eveProfiles}
     userSelections={userSelections}
     currentSettingsDir={currentSettingsDir}
+    isDefaultDir={isDefaultDir}
     lastBackupDir={lastBackupDir}
 />
```

**`renderer/src/pages/Profiles.jsx`**:

- Add `isDefaultDir` to the prop list (default `false`).
- Remove `const [isDefaultDir, setIsDefaultDir] = useState(false);` (no longer local).
- Remove `setIsDefaultDir(false)` from `handleChooseSettingsDir` and `setIsDefaultDir(true)` from `handleResetToDefault`. The next `refreshData()` re-fetches config and the prop drives the UI.
- Update `Profiles.propTypes` if PropTypes are still around when this commit lands; if Section C runs first this is moot.

**Test:** extend `Profiles.test.jsx` to verify the Reset-to-default button is hidden when `isDefaultDir={true}` and shown when `isDefaultDir={false}`.

### Acceptance

- API: `GET /api/config` (or equivalent) response includes `isDefaultDir`.
- Frontend: opening `/profiles?view=sync` on a fresh install where `currentSettingsDir` matches the Tranquility default → Reset button is **hidden**. Choose a custom dir → Reset button appears. Click Reset → button hides again after `refreshData`.

---

## Section B — `PlanList` inline icon → `EveTypeIcon`

**File:** `renderer/src/components/skillplan/PlanList.jsx`

Replace lines around 109-113 (the inline `<img>` with `https://images.evetech.net/types/{typeID}/icon` and the placeholder `<span>` fallback) with:

```jsx
<EveTypeIcon name={p.name} conversions={conversions} />
```

Add the import at the top:

```jsx
import EveTypeIcon from '../ui/EveTypeIcon.jsx';
```

Drop the now-unused `typeID` and `icon` local variables computed inside the row map.

**Test:** existing `PlanList.test.jsx` continues to pass. No new tests required — `EveTypeIcon` has its own coverage.

### Acceptance

- `PlanList` renders identical UI as before.
- `git grep -n 'images.evetech.net/types' renderer/src` returns zero hits inside `PlanList.jsx`.

---

## Section C — Strip PropTypes (46 files)

### Pre-strip verification

Before stripping, confirm no code reads `<Component>.propTypes` at runtime:

```bash
grep -rn '\.propTypes' renderer/src --include="*.jsx" --include="*.js" | \
    grep -v '^[^:]*:[0-9]*:.*PropTypes' | \
    grep -v '^[^:]*:[0-9]*:.*\.propTypes\s*=\s*{'
```

Any hit that isn't an assignment (`= {`) or an import line is a real runtime reader and must be reviewed before stripping. Expected: zero hits.

### Strip

Three mechanical edits:

1. **Source files** — for each `.jsx`/`.js` under `renderer/src/`:
    - Delete the line `import PropTypes from 'prop-types';`
    - Delete the multi-line block `<Component>.propTypes = { … };` (matching from the assignment through the closing `};`).

   Implementation: a Node script run once that walks the tree, parses each file, deletes matching AST nodes, writes back. Or a `sed`-based approach for the import line plus a manual pass for the multi-line block — but a script is preferred for the multi-line case.

   The plan-writing step picks the implementation. The criterion: no diff on file behavior, only deletions.

2. **Dependency removal** — `renderer/package.json`: delete `"prop-types": "^15.8.1"` from `dependencies`. Run `npm install` to regenerate `package-lock.json`.

3. **Validate**:
    - `cd renderer && npm run lint` — must pass with the same error count as pre-strip (no new violations).
    - `cd renderer && npm test` — all 117 tests must still pass.
    - `cd renderer && npm run build` — must succeed.

### Acceptance

- `grep -rn "prop-types\|PropTypes" renderer/src --include="*.jsx" --include="*.js"` returns zero hits.
- `renderer/package.json` no longer lists `prop-types`.
- Test suite green; build green.

---

## Implementation order and PR shape

One PR, three commits in this order:

1. `feat(api): expose isDefaultDir from config service` (Section A backend) + `feat(profiles): drive isDefaultDir from backend` (Section A frontend) — these can be one commit since the frontend depends on the backend field. The plan-writing step will decide whether to split.
2. `refactor(skill-plans): use EveTypeIcon in PlanList` (Section B).
3. `chore: strip inert PropTypes; drop prop-types dep` (Section C).

Order rationale: A first because it's the only section that touches backend; B is trivial filler; C goes last because its blast radius (46 files) is largest and a regression is easiest to bisect when it's the most recent commit.

## Risks

| Risk | Likelihood | Mitigation |
|------|-----------|------------|
| The `/api/config` response shape is consumed by an Electron or webview client we forgot about | Low | Adding a new optional field is additive — old clients ignore it. |
| A test or other module reads `<Component>.propTypes` at runtime | Low | Pre-strip grep verification (above). |
| IDE hover hints for component props degrade after strip | Likely | Acceptable cost. Editors can read prop destructure + JSX usage. |
| `npm install` after removing the dep introduces unrelated lockfile churn | Medium | Run `npm install` separately; review the lockfile diff before committing. |

## What this design intentionally does not do

- Replace PropTypes with a runtime validator. That's reinventing PropTypes for one component while leaving 45 untouched — strictly worse than the strip.
- Add JSDoc `@param` annotations as a soft type substitute. Out of scope; can be added later if/when an editor-tooling motivation appears.
- Touch the dynamic-import warnings the production build emits (`accountsApi`, `configApi`). Pre-existing; unrelated.
