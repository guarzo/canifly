# Post-Redesign Cleanup Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Tie up three loose ends from PR #26: drive `Profiles.isDefaultDir` from a new backend signal; replace the inline icon in `PlanList` with `EveTypeIcon`; strip the now-inert PropTypes declarations and the `prop-types` dependency.

**Architecture:** A → B → C in one PR with three logical commits (Section A is one commit because the frontend depends on the new backend field). Section C is mechanical — a Node script walks `renderer/src/`, deletes the `prop-types` import line and the `<Component>.propTypes = { … };` block from each file. The PR ships when lint, full vitest suite (117 tests), and `npm run build` are all green.

**Tech Stack:** Go 1.23.2 backend, React 19.2 + Tailwind v4 + Vitest renderer, MUI v6, Zustand.

**Spec:** `docs/superpowers/specs/2026-05-05-post-redesign-cleanup-design.md`

---

## Pre-flight

### Task 0: Branch + green baseline

**Files:** none

- [ ] **Step 1: Create branch from main**

```bash
git checkout main && git pull
git checkout -b feature/post-redesign-cleanup
```

- [ ] **Step 2: Verify baseline tests pass**

```bash
cd renderer && npm test
```

Expected: 30 files / 117 tests passing.

- [ ] **Step 3: Verify Go tests pass**

```bash
go test ./...
```

Expected: all green.

- [ ] **Step 4: Verify the build succeeds**

```bash
cd renderer && npm run build
```

Expected: build succeeds.

If anything is red on baseline, stop and report — do not start work on a broken baseline.

---

## Section A — Profiles.isDefaultDir from backend

### Task A.1: Add `IsDefaultSettingsDir` to the configuration service

**Files:**
- Modify: `internal/services/interfaces/configuration.go`
- Modify: `internal/services/config/configuration_service.go`
- Modify: `internal/testutil/mock_interfaces.go`

- [ ] **Step 1: Add the method to the interface**

In `internal/services/interfaces/configuration.go`, add this line inside the `ConfigurationService` interface (alongside the other Configuration Management methods, around line 17):

```go
IsDefaultSettingsDir() (bool, error)
```

The interface should now include (positions approximate, just place under the other directory methods):

```go
type ConfigurationService interface {
    UpdateSettingsDir(dir string) error
    GetSettingsDir() (string, error)
    EnsureSettingsDir() error
    IsDefaultSettingsDir() (bool, error)
    // ...existing methods...
}
```

- [ ] **Step 2: Implement the method on `ConfigurationService`**

Append this method to `internal/services/config/configuration_service.go` (place it next to `GetSettingsDir`, around line 102 after the existing method's closing brace):

```go
// IsDefaultSettingsDir reports whether the currently configured settings
// directory matches the OS-default Tranquility location. Returns false when
// no directory is configured or when auto-detection has not yet picked one.
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

- [ ] **Step 3: Update the mock to satisfy the interface**

In `internal/testutil/mock_interfaces.go`, add a method to `MockConfigService` (place it near the other Mock methods, e.g., after `EnsureSettingsDir` around line 323):

```go
func (m *MockConfigService) IsDefaultSettingsDir() (bool, error) {
    args := m.Called()
    return args.Bool(0), args.Error(1)
}
```

- [ ] **Step 4: Verify the Go build passes**

```bash
go build ./...
```

Expected: no errors. Note: there's no compile-time interface assertion (`var _ interfaces.ConfigurationService = (*MockConfigService)(nil)`) so adding the mock method is good hygiene but the build won't fail without it. If you want a stronger guarantee, add such an assertion under the mock declaration as a follow-up.

- [ ] **Step 5: Verify Go tests pass**

```bash
go test ./...
```

Expected: all green. The mock change is additive; existing tests should not be affected.

- [ ] **Step 6: Hold the commit until A.2 lands**

This change has no caller yet. Don't commit alone — Task A.2 wires it into the handler so the change ships as one logical unit.

---

### Task A.2: Wire `isDefaultDir` into the GET /api/config response

**Files:**
- Modify: `internal/handlers/config.go:31-56`

- [ ] **Step 1: Read the current handler**

Open `internal/handlers/config.go`. The `GetConfig()` method currently builds a `configData` map at lines 43-48 with four fields: `settingsDir`, `roles`, `userSelections`, `lastBackupDir`. Confirm that's the current shape before editing.

- [ ] **Step 2: Compute `isDefaultDir` and add it to the map**

Replace the inner function passed to `WithCache` (the `func() (interface{}, error) { … }` block, lines 37-54) with this version. The change: call the new service method and add an `isDefaultDir` key to the map.

```go
func() (interface{}, error) {
    config, err := h.configService.FetchConfigData()
    if err != nil {
        return nil, err
    }

    isDefault, err := h.configService.IsDefaultSettingsDir()
    if err != nil {
        // Non-fatal: log and fall back to false so the UI degrades to "not default" rather than failing the whole config request.
        h.logger.Warnf("IsDefaultSettingsDir failed: %v", err)
        isDefault = false
    }

    configData := map[string]interface{}{
        "settingsDir":    config.SettingsDir,
        "isDefaultDir":   isDefault,
        "roles":          config.Roles,
        "userSelections": config.DropDownSelections,
        "lastBackupDir":  config.LastBackupDir,
    }

    // Return in the format the frontend expects
    return map[string]interface{}{
        "data": configData,
    }, nil
},
```

- [ ] **Step 3: Verify the Go build still passes**

```bash
go build ./...
```

Expected: no errors.

- [ ] **Step 4: Verify Go tests pass**

```bash
go test ./...
```

Expected: all green.

- [ ] **Step 5: Integration-verify the field appears in the response**

Start the backend (in a second terminal: `npm run dev:go`) and curl the endpoint. From an authenticated session — easiest is to launch the full app once (`npm start`), let it auth, then curl from the same machine using the session cookie. Or simpler: bypass auth-checking by editing the test fixtures.

```bash
# After app is running and authenticated:
curl -s http://localhost:42423/api/config --cookie "$(grep -o 'session=[^;]*' ~/.config/canifly/cookies 2>/dev/null || echo '')" | jq '.data | keys'
```

Expected output includes `"isDefaultDir"`. If you can't easily authenticate, accept the unit-of-build verification and proceed — the field's presence will be confirmed end-to-end in Task A.4 (frontend test).

- [ ] **Step 6: Hold the commit until A.3 lands**

The frontend change in A.3 closes the loop. Commit the bundle in A.4 after the frontend test passes.

---

### Task A.3: Read `isDefaultDir` from config in `Routes.jsx`, pass to `Profiles`

**Files:**
- Modify: `renderer/src/Routes.jsx:37-45`

- [ ] **Step 1: Read the current routes file**

Open `renderer/src/Routes.jsx`. Confirm the lines extracting config fields look like this around lines 37-41:

```jsx
const roles = config?.roles || config?.Roles || [];
const userSelections = config?.userSelections || config?.DropDownSelections || {};
const currentSettingsDir = config?.settingsDir || config?.SettingsDir || '';
const rawLastBackup = config?.lastBackupDir || config?.LastBackupDir;
const lastBackupDir = typeof rawLastBackup === 'string' ? rawLastBackup : '';
```

- [ ] **Step 2: Add a line that pulls `isDefaultDir`**

Insert this line right after the `currentSettingsDir` line:

```jsx
const isDefaultDir = Boolean(config?.isDefaultDir ?? config?.IsDefaultDir ?? false);
```

(The dual-read mirrors the existing camelCase/PascalCase pattern in this file. The Go handler emits camelCase; the PascalCase fallback is defensive.)

- [ ] **Step 3: Pass the prop to `Profiles`**

In the `<Profiles … />` JSX block (around line 67-77), add `isDefaultDir={isDefaultDir}` between `currentSettingsDir` and `lastBackupDir`:

```jsx
<Profiles
    subDirs={eveProfiles}
    associations={associations}
    settingsData={eveProfiles}
    userSelections={userSelections}
    currentSettingsDir={currentSettingsDir}
    isDefaultDir={isDefaultDir}
    lastBackupDir={lastBackupDir}
/>
```

- [ ] **Step 4: Hold the commit until A.4 lands**

A.4 makes `Profiles.jsx` consume the prop. Commit together.

---

### Task A.4: Consume `isDefaultDir` prop in `Profiles.jsx`; drop local state

**Files:**
- Modify: `renderer/src/pages/Profiles.jsx`
- Modify: `renderer/src/pages/Profiles.test.jsx`

- [ ] **Step 1: Write a failing test in `Profiles.test.jsx`**

Append a new test to the existing `describe('Profiles smart default', ...)` block (or create a new `describe` if you prefer; either works):

```jsx
describe('Profiles reset-to-default button visibility', () => {
    test('hides the reset-to-default button when isDefaultDir is true (sync mode)', () => {
        wrap(
            <Profiles
                subDirs={[]}
                associations={[]}
                settingsData={[{ profile: 'settings_default', availableCharFiles: [], availableUserFiles: [] }]}
                userSelections={{}}
                currentSettingsDir="/path/to/Tranquility"
                isDefaultDir={true}
                lastBackupDir=""
            />,
        );
        expect(screen.queryByRole('button', { name: /reset to default directory/i })).toBeNull();
    });

    test('shows the reset-to-default button when isDefaultDir is false (sync mode)', () => {
        wrap(
            <Profiles
                subDirs={[]}
                associations={[]}
                settingsData={[{ profile: 'settings_default', availableCharFiles: [], availableUserFiles: [] }]}
                userSelections={{}}
                currentSettingsDir="/some/custom/dir"
                isDefaultDir={false}
                lastBackupDir=""
            />,
        );
        expect(screen.getByRole('button', { name: /reset to default directory/i })).toBeInTheDocument();
    });
});
```

The test relies on Sync mode rendering the header actions including the reset button. The smart-default branch in `Profiles` will pick `'sync'` because there are no unmatched characters when `subDirs=[]`. The test passes `settingsData` with one profile so the page is in Sync mode (not the empty-state branch — the existing tests proved that).

- [ ] **Step 2: Run the test to confirm it fails**

```bash
cd renderer && npm test -- Profiles
```

Expected: FAIL — current code uses `useState(false)`, so the button shows regardless. The "hides when true" test fails because the button is found.

- [ ] **Step 3: Modify `Profiles.jsx` — add `isDefaultDir` to the prop list**

Open `renderer/src/pages/Profiles.jsx`. Find the function signature for `Profiles` (around line 76-84). Add `isDefaultDir = false` to the destructured props:

```jsx
const Profiles = ({
    subDirs = [],
    associations = [],
    settingsData = [],
    userSelections = {},
    currentSettingsDir = '',
    isDefaultDir = false,
    lastBackupDir = '',
}) => {
```

- [ ] **Step 4: Remove the local `isDefaultDir` state**

Find this line (around line 87):

```jsx
const [isDefaultDir, setIsDefaultDir] = useState(false);
```

Delete it.

- [ ] **Step 5: Remove all `setIsDefaultDir` calls**

Search for `setIsDefaultDir` in `Profiles.jsx`. Two call sites exist:

In `handleChooseSettingsDir` (around line 138-145), remove the line `if (result?.success) setIsDefaultDir(false);`. The function should now look like:

```jsx
const handleChooseSettingsDir = useCallback(async () => {
    const chosen = await window.electronAPI.chooseDirectory();
    if (!chosen) { toast.info('No directory chosen.'); return; }
    await run(async () => {
        const result = await chooseSettingsDir(chosen);
        return { ...result, message: result?.message || `Settings directory: ${chosen}` };
    }, { errorContext: 'chooseSettingsDir' });
}, [run]);
```

In `handleResetToDefault` (around line 156-167), remove the line `if (result?.success) setIsDefaultDir(true);`. The function should now look like:

```jsx
const handleResetToDefault = useCallback(async () => {
    const ok = await showConfirmDialog({
        title: 'Reset to default',
        message: 'Reset the settings directory to the default Tranquility location?',
    });
    if (!ok.isConfirmed) return;
    await run(async () => {
        const result = await resetToDefaultDirectory();
        return { ...result, message: 'Reset to default: Tranquility' };
    }, { errorContext: 'resetToDefaultDirectory' });
}, [run, showConfirmDialog]);
```

- [ ] **Step 6: Update the PropTypes block**

In `Profiles.jsx` find the `Profiles.propTypes = { … };` block (around line 287-294). Add `isDefaultDir: PropTypes.bool,` between `currentSettingsDir` and `lastBackupDir`:

```jsx
Profiles.propTypes = {
    subDirs: PropTypes.array,
    associations: PropTypes.array,
    settingsData: PropTypes.array,
    userSelections: PropTypes.object,
    currentSettingsDir: PropTypes.string,
    isDefaultDir: PropTypes.bool,
    lastBackupDir: PropTypes.string,
};
```

(This block disappears entirely in Section C. That's fine — keeping it accurate in the meantime.)

- [ ] **Step 7: Run the test to confirm it passes**

```bash
cd renderer && npm test -- Profiles
```

Expected: PASS. All Profiles tests green (3 existing + 2 new = 5).

- [ ] **Step 8: Run the full vitest suite**

```bash
cd renderer && npm test
```

Expected: all 119 tests pass (117 baseline + 2 new in Profiles).

- [ ] **Step 9: Commit Section A as one unit**

```bash
git add internal/services/interfaces/configuration.go internal/services/config/configuration_service.go internal/testutil/mock_interfaces.go internal/handlers/config.go renderer/src/Routes.jsx renderer/src/pages/Profiles.jsx renderer/src/pages/Profiles.test.jsx
git commit -m "feat(profiles): drive isDefaultDir from backend config

Adds IsDefaultSettingsDir() to ConfigurationService, returns it as
isDefaultDir in the GET /api/config response, and threads it through
Routes -> Profiles as a prop. Profiles drops its local useState; the
Reset-to-default button now hides correctly on first load when the
current directory matches the OS default Tranquility location."
```

---

## Section B — PlanList inline icon → EveTypeIcon

### Task B.1: Replace inline icon with `EveTypeIcon`

**Files:**
- Modify: `renderer/src/components/skillplan/PlanList.jsx`

- [ ] **Step 1: Read the current file**

Open `renderer/src/components/skillplan/PlanList.jsx`. Locate the section that builds the row icon — currently around lines 80-115. The relevant block computes `typeID` and `icon` and renders one of two branches:

```jsx
const typeID = conversions?.[p.name];
const icon = typeID ? `https://images.evetech.net/types/${typeID}/icon` : null;
// ... later, in the row JSX:
{icon ? (
    <img src={icon} alt="" aria-hidden loading="lazy" className="h-6 w-6 rounded-sm border border-rule-1 shrink-0" />
) : (
    <span className="h-6 w-6 rounded-sm border border-rule-1 bg-surface-2 shrink-0" aria-hidden />
)}
```

- [ ] **Step 2: Add the import**

At the top of the file, after the existing imports, add:

```jsx
import EveTypeIcon from '../ui/EveTypeIcon.jsx';
```

- [ ] **Step 3: Delete the now-unused locals**

Find `const typeID = conversions?.[p.name];` and `const icon = typeID ? \`https://...\` : null;` inside the `plans.map(...)` callback. Delete both lines.

- [ ] **Step 4: Replace the JSX branch with `EveTypeIcon`**

Find the `{icon ? ( <img ... /> ) : ( <span ... /> )}` block and replace it with:

```jsx
<EveTypeIcon name={p.name} conversions={conversions} />
```

- [ ] **Step 5: Run PlanList tests**

```bash
cd renderer && npm test -- PlanList SkillPlans
```

Expected: PASS (no test changes needed — the existing copy-bug test and SkillPlans smoke test continue to pass; visual icon rendering is covered by `EveTypeIcon`'s own tests).

- [ ] **Step 6: Commit**

```bash
git add renderer/src/components/skillplan/PlanList.jsx
git commit -m "refactor(skill-plans): use EveTypeIcon in PlanList rows"
```

---

## Section C — Strip PropTypes

### Task C.0: Pre-strip verification — confirm no runtime readers

**Files:** none

- [ ] **Step 1: Grep for any non-import, non-assignment uses of `.propTypes`**

```bash
grep -rn '\.propTypes' renderer/src --include="*.jsx" --include="*.js" | \
    grep -v 'PropTypes' | \
    grep -v '\.propTypes\s*=' | \
    grep -v '^\s*//'
```

Expected: zero hits.

If there are hits, they are real runtime readers (e.g., test code asserting `Component.propTypes` shape). They must be removed or migrated **before** the strip. If you find any, stop and report.

- [ ] **Step 2: Inventory the strip target list**

```bash
grep -rln "from 'prop-types'" renderer/src --include="*.jsx" --include="*.js" | wc -l
```

Expected: 46 (rounded; should match what was reported in the spec).

Save the file list for the next task:

```bash
grep -rln "from 'prop-types'" renderer/src --include="*.jsx" --include="*.js" > /tmp/proptype-files.txt
wc -l /tmp/proptype-files.txt
```

---

### Task C.1: Strip PropTypes from renderer source

**Files:** all 46 files listed in `/tmp/proptype-files.txt`

This is a mechanical rewrite. Use the script below to do it in one pass and inspect the diff before committing.

- [ ] **Step 1: Write the strip script**

Save this to `scripts/strip-proptypes.mjs` (project root, not under `renderer/`):

```js
#!/usr/bin/env node
// Strip `import PropTypes from 'prop-types'` and `<Identifier>.propTypes = { ... };`
// blocks from React source files. Idempotent: running twice produces no further diffs.
import { readFileSync, writeFileSync } from 'node:fs';
import { argv } from 'node:process';

const files = argv.slice(2);
if (files.length === 0) {
    console.error('Usage: node scripts/strip-proptypes.mjs <file>...');
    process.exit(1);
}

let totalImportLinesRemoved = 0;
let totalBlocksRemoved = 0;

for (const file of files) {
    const before = readFileSync(file, 'utf8');
    let after = before;

    // 1) Strip the import line. Allow single or double quotes; allow trailing semicolon
    //    or none; consume up to and including the trailing newline so we don't leave a
    //    blank line behind.
    const importPattern = /^import\s+PropTypes\s+from\s+['"]prop-types['"];?\s*\r?\n/m;
    if (importPattern.test(after)) {
        after = after.replace(importPattern, '');
        totalImportLinesRemoved += 1;
    }

    // 2) Strip the `<Identifier>.propTypes = { ... };` block.
    //    The block is multi-line. We match from `<Word>.propTypes = {` through the
    //    matching closing `};` and consume any trailing blank line.
    //    JS regex doesn't have native brace-balancing, but in this codebase every
    //    propTypes block ends with `};` on its own line. Match non-greedy to the first
    //    `};` at the start of a line.
    const blockPattern = /\n[A-Za-z_$][\w$]*\.propTypes\s*=\s*\{[\s\S]*?\n\};\s*\r?\n/g;
    const blockMatches = after.match(blockPattern) || [];
    if (blockMatches.length > 0) {
        after = after.replace(blockPattern, '\n');
        totalBlocksRemoved += blockMatches.length;
    }

    if (after !== before) {
        writeFileSync(file, after, 'utf8');
        console.log(`stripped ${file}`);
    }
}

console.log(`done — removed ${totalImportLinesRemoved} import lines, ${totalBlocksRemoved} propTypes blocks`);
```

- [ ] **Step 2: Run the script over the inventory**

```bash
node scripts/strip-proptypes.mjs $(cat /tmp/proptype-files.txt)
```

Expected output: `done — removed 46 import lines, ~46 propTypes blocks` (the block count may differ slightly — a few files might have multiple components per file, or a sub-component, but each file has at least one block).

- [ ] **Step 3: Inspect a few representative diffs**

```bash
git diff renderer/src/components/ui/PageShell.jsx
git diff renderer/src/components/profiles/SyncProfileRow.jsx
git diff renderer/src/pages/Profiles.jsx
```

Each should show only deletions (the `import PropTypes…` line and the `<Component>.propTypes = { … };` block). No other changes. If you see anything else changed — function bodies, JSX, etc. — stop and revert; the script is over-matching.

- [ ] **Step 4: Verify no `PropTypes` references remain**

```bash
grep -rn "PropTypes\|prop-types" renderer/src --include="*.jsx" --include="*.js"
```

Expected: zero hits.

If any remain, they are leftovers the regex didn't catch (e.g., weird whitespace, inline `Component.propTypes = { single: pt }` on one line). Hand-edit the remaining ones, then re-run the verify.

- [ ] **Step 5: Run lint to catch syntax errors**

```bash
cd renderer && npx eslint src
```

Expected: same error count as on baseline (the codebase has 230 pre-existing errors, mostly `no-console`). Crucially: no **new** errors. If new errors appear, they're caused by malformed deletions — revert and investigate.

To compare counts:

```bash
cd renderer && npx eslint src 2>&1 | tail -3
```

Note the "X problems (Y errors, Z warnings)" line. It should be unchanged (or slightly lower if PropTypes-related lint rules dropped warnings).

- [ ] **Step 6: Run the full vitest suite**

```bash
cd renderer && npm test
```

Expected: all 119 tests pass (117 baseline + 2 new from Section A).

- [ ] **Step 7: Run the production build**

```bash
cd renderer && npm run build
```

Expected: build succeeds.

- [ ] **Step 8: Commit**

```bash
git add renderer/src scripts/strip-proptypes.mjs
git commit -m "chore: strip PropTypes (inert in React 19) from renderer source

PropTypes runtime validation was deprecated in React 19 — declarations
are no longer checked. Removed all 'import PropTypes from prop-types'
lines and all <Component>.propTypes = { … } blocks across 46 files via
scripts/strip-proptypes.mjs (kept for reference / re-runs).

The 'prop-types' npm dep is removed in the next commit."
```

---

### Task C.2: Remove `prop-types` dependency

**Files:**
- Modify: `renderer/package.json`
- Modify: `renderer/package-lock.json` (regenerated by npm)

- [ ] **Step 1: Remove the dep from `package.json`**

Open `renderer/package.json`. Find the `dependencies` block. Delete the line `"prop-types": "^15.8.1",` (the version may differ — match by package name).

- [ ] **Step 2: Regenerate the lockfile**

```bash
cd renderer && npm install
```

Expected: npm reports a small set of changes (one less direct dep, possibly some transitive deps drop). No errors.

- [ ] **Step 3: Sanity-check that no transitive code still imports it**

```bash
grep -rn "from 'prop-types'\|require('prop-types')" renderer/src renderer/electron renderer/main.js 2>&1 | head
```

Expected: zero hits in renderer source. The dep only stays installed if some other dep transitively pulls it (which is fine — it's the *direct* dependency we're removing).

- [ ] **Step 4: Run lint, tests, and build one more time**

```bash
cd renderer && npx eslint src 2>&1 | tail -3
cd renderer && npm test
cd renderer && npm run build
```

All three must pass. Lint count unchanged from after C.1.

- [ ] **Step 5: Commit**

```bash
git add renderer/package.json renderer/package-lock.json
git commit -m "chore: remove prop-types from renderer dependencies

PropTypes were stripped from source in the previous commit; this removes
the now-unused npm dependency."
```

---

## Final validation

### Task Z: PR-ready checks

**Files:** none

- [ ] **Step 1: Confirm commit history**

```bash
git log --oneline main..HEAD
```

Expected output (4 commits in this order):
- `chore: remove prop-types from renderer dependencies`
- `chore: strip PropTypes (inert in React 19) from renderer source`
- `refactor(skill-plans): use EveTypeIcon in PlanList rows`
- `feat(profiles): drive isDefaultDir from backend config`

- [ ] **Step 2: Final test + build pass**

```bash
go test ./...
cd renderer && npm test && npm run build
```

All green.

- [ ] **Step 3: Push and open PR**

```bash
git push -u origin feature/post-redesign-cleanup
gh pr create --base main --title "feat: post-redesign cleanup (isDefaultDir, PlanList icon, PropTypes strip)" --body "$(cat <<'EOF'
## Summary

Three loose ends from PR #26:

- **`Profiles.isDefaultDir`** now derives from a backend signal. Adds `IsDefaultSettingsDir()` to `ConfigurationService`, exposes it as `isDefaultDir` in `GET /api/config`, threads it through Routes → Profiles. The Reset-to-default button hides correctly on first load when the user is already on the default Tranquility directory.
- **`PlanList`** uses the `EveTypeIcon` primitive (introduced in #26) instead of the inline `<img>` it had been duplicating.
- **PropTypes stripped** from all 46 renderer source files. PropTypes runtime validation is inert in React 19; the declarations were misleading. The `prop-types` npm dep is removed.

## Spec
`docs/superpowers/specs/2026-05-05-post-redesign-cleanup-design.md`

## Test plan

- [ ] `go test ./...` passes
- [ ] `cd renderer && npm test` — 30+ files, 119+ tests passing
- [ ] `cd renderer && npm run build` succeeds
- [ ] Manual: `/profiles?view=sync` on first launch with default Tranquility dir → Reset button is hidden
- [ ] Manual: choose a custom dir → Reset button appears → click Reset → button hides again

🤖 Generated with [Claude Code](https://claude.com/claude-code)
EOF
)"
```

---

## Self-review checklist (engineer running this plan should glance over before starting)

- Section A ships as one commit; the backend method must be implemented before the handler references it (Tasks A.1 → A.2 → A.3 → A.4 → single commit at end of A.4).
- Section C must run **after** Section A and B because the `Profiles.propTypes` block added in A.4 gets deleted by C.1's script — that's fine and expected, just don't be surprised when the diff in C.1 includes the line you added in A.4.
- The `MockConfigService` mock in `internal/testutil/mock_interfaces.go` should be updated alongside the interface change in A.1. Note that the mock is currently unused by any test and the codebase has no compile-time interface assertion, so the build won't actually fail without the update — but skipping it leaves a stale mock that will silently start failing once any test uses it. Update it anyway.
