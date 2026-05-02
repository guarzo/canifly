# Dependency Refresh — Design Spec

**Date:** 2026-05-01
**Status:** Approved
**Scope:** Plan C of three (Plan A: backend refactor; Plan B: frontend cleanup)

## Goal

Bring the toolchain and dependencies fully current after ~8 months of inactivity: Go to current stable, Electron to current major, React 18→19, Material-UI v5→v7. Each upgrade is its own PR with manual smoke test + single-OS package step so regressions are bisectable.

## Findings addressed

Finding #10 from the 2026-05-01 codebase review.

## Sequencing

**Plan C runs after Plan A merges.** Both touch `services.go` indirectly (Plan A restructures it; Plan C may add Go-version-conditional code there). Running them serially avoids merge conflicts. Plan B is independent (frontend only) and can run in parallel with either A or C.

## Current state

Last commit: 2025-08-23. Versions per `go.mod`, `package.json`, `renderer/package.json`:

- Go: 1.23.2 (current stable: 1.25.x as of May 2026)
- Electron: 32.2.0 (current: ~36)
- React: 18.x (React 19 GA)
- Material-UI: v5 (MUI v7 current)
- Vitest, ESLint, axios, electron-builder: drift unverified

No baseline `npm audit` or `govulncheck` data yet — first task captures it.

## Target state

Eight PRs in order. Each is independently revertible and includes:
- Code change + lockfile update
- `go test -race ./...` and/or `npm run test:react` passing
- `go vet ./...` clean (Go PRs)
- Manual smoke test: login → character list → skill plan check → sync backup → logout
- `npm run dist` on the maintainer's primary OS (Electron and bundling-affecting PRs only)
- Smoke result documented in PR description

### PR 1: Audit baseline (no code changes)

Capture and commit to `docs/superpowers/dependency-audit-2026-05-01.md`:

- `go list -u -m all` output
- `govulncheck ./...` output
- `npm audit --audit-level=moderate` (root)
- `cd renderer && npm audit --audit-level=moderate`
- `npm outdated` for both roots

This document drives the choices in PR 3 (which Go modules to bump beyond the security floor).

### PR 2: Go toolchain bump

- Bump `go.mod` `go` directive to current stable minor (1.25 if released, else 1.24)
- Update CI configuration and `package.json` `go:build` script if Go version is pinned there
- Run `go mod tidy`
- Address any new `go vet` warnings

### PR 3: Go module updates

For each module flagged by `govulncheck` and any module with security-only patches available:
- `go get module@latest`
- One commit per module (or grouped: `golang.org/x/*` together)

Defer non-security major upgrades (e.g., a hypothetical v2 of an existing dep) unless they unblock something. YAGNI.

### PR 4: Node + npm engines

- Bump `engines.node` in both `package.json` files to current LTS (≥22 if Electron target supports, else 20)
- Verify CI uses matching Node version
- Fresh `npm install` in each root; commit lockfile diffs

### PR 5: Electron major upgrade

Highest-risk PR. Process:
1. Read Electron breaking-changes notes between 32 and target (likely 36)
2. Bump `electron`, `electron-builder` in root `package.json`
3. Update `main.js` and `preload.js` for any deprecated/removed API
4. Run `npm run dist` on the maintainer's primary OS — must succeed
5. Manual smoke test on the packaged app: full user flow

If breaking changes are non-trivial across multiple majors (32→33→34→35→36), split into incremental PRs (32→34, 34→36, etc.) until current.

### PR 6: React 18 → 19

1. Read React 19 migration guide
2. Run the official React codemod if available
3. Bump `react`, `react-dom`, `react-router`, `@types/react` (if present) to React-19-compatible versions
4. `npm run test:react` and fix any failures
5. Manual smoke test

### PR 7: Material-UI v5 → v7

1. Run `@mui/codemod` for v5→v6, then v6→v7 (separate runs)
2. Bump `@mui/material`, `@mui/icons-material`, `@emotion/*` to compatible versions
3. Run `npm run test:react`
4. Manual visual check on every screen — Character Overview, Sync, Mapping, settings pages, all modals — and document anything that needed manual fix-up in the PR description

If a v5→v7 jump in one PR is too risky after auditing the codemod output, split into v5→v6 and v6→v7 PRs.

### PR 8: Frontend transitive cleanup

After PRs 4–7:
- `cd renderer && npm audit fix`
- `npm outdated` again — bump remaining minor/patch versions in one PR
- Same in root for Electron-side deps
- Final smoke test

## Risks and mitigations

| Risk | Mitigation |
|---|---|
| Electron major bump breaks packaging on one OS | Single-OS package per PR is the standing requirement; if a user reports a different-OS regression after merge, do a hotfix PR with that OS's package step added |
| React 19 + MUI v5/v6 incompatibility | If MUI v7 isn't compatible with React 19 by the time PR 6 runs, swap PR 6 and PR 7 order |
| Auto-updater regression in Electron | Verify auto-update flow as part of PR 5 smoke test or document explicit deferral |
| `govulncheck` flags transitive Go deps not pulled in directly | Bump the direct dep that pulls them; document if no upgrade path exists |
| Codemods produce mechanically-correct but ugly diffs | Accept the codemod output; clean-up is out of scope for this plan |
| MUI v7 visual regressions across many screens | Visual check is part of PR 7; budget time for manual touch-ups |

## Out of scope

- Switching to a different bundler (verify Vite is current; do not replace)
- TypeScript migration
- Replacing Material-UI with another UI library
- Replacing Zustand or other state libs
- Adding/changing lint rules
- Reformatting the codebase
- Refactoring code that doesn't compile under new versions beyond the minimum needed (those become follow-up findings)

## Acceptance criteria

- `govulncheck ./...` reports no findings
- `npm audit --audit-level=high` reports no findings in either root
- `go.mod` declares Go ≥ current stable minor − 1
- `engines.node` declares ≥ current LTS in both `package.json` files
- Electron at current major
- React at v19
- Material-UI at v7
- App launches, packages, and completes the full smoke test on the maintainer's primary OS at the end of every PR
- Each PR description records the smoke-test outcome
