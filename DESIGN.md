# Design

The visual system for CanIFly. Read PRODUCT.md first for strategy. This file is the source of truth for tokens, type, motion, and primitives.

The system is deliberately small. Add to it only when a real surface needs a thing this file does not already give it.

## Theme

Dark only, warm slate. Not blue-black, not pure black. The base is a desaturated slate tinted slightly toward warm to avoid the cold-screen feel of the EVE client and most "developer dark" themes. Sits next to a 2nd-monitor EVE window without competing with it.

There is no light mode. If we add one later, it will be a separate token set, not a per-component override.

## Color

OKLCH everywhere. No `#000`, no `#fff`. Every neutral carries a tiny warm chroma (≈0.005) so the surface never looks clinical. Strategy is **Restrained**: neutrals carry the page; one accent below 10% surface area; status colors are functional.

Tokens live in CSS variables on `:root` and are mapped into Tailwind under `slate.*`, `ink.*`, `accent.*`, and `status.*`. Components reference token names, never raw OKLCH values.

### Surface (warm slate)

| Token | OKLCH | Use |
|---|---|---|
| `--surface-0` | `oklch(0.16 0.005 65)` | App background. The deepest layer. |
| `--surface-1` | `oklch(0.20 0.005 65)` | Default panel / row container. |
| `--surface-2` | `oklch(0.24 0.006 65)` | Sticky headers, group headers, hover row. |
| `--surface-3` | `oklch(0.28 0.007 65)` | Active row, selected segment, popover. |
| `--surface-overlay` | `oklch(0.10 0.004 65 / 0.6)` | Modal scrim. Used sparingly. |

Hue 65 is a warm beige direction. Chroma is intentionally tiny — these read as neutrals, not as "warm" in a way the user could name.

### Ink (text)

| Token | OKLCH | Use |
|---|---|---|
| `--ink-1` | `oklch(0.96 0.005 80)` | Primary text. Headings. Numbers. |
| `--ink-2` | `oklch(0.78 0.006 80)` | Secondary text. Body. |
| `--ink-3` | `oklch(0.58 0.007 80)` | Tertiary text. Meta, units, helpers. |
| `--ink-4` | `oklch(0.42 0.007 80)` | Disabled, dividers-against-text. |

`--ink-1` against `--surface-1` measures **15.4:1 contrast** — well past WCAG AA. `--ink-3` against `--surface-1` measures 4.6:1 — passes AA for body text and is used only as helper / unit copy, never as the only signal for state.

### Borders & rules

| Token | OKLCH | Use |
|---|---|---|
| `--rule-1` | `oklch(0.32 0.005 65)` | Hairline rules between rows, group dividers. |
| `--rule-2` | `oklch(0.40 0.006 65)` | Hovered border, focused container. |
| `--rule-strong` | `oklch(0.55 0.008 65)` | High-contrast border for active selection. |

All borders are **1px**. There is no thicker accent border anywhere in the system.

### Accent

A single accent. Used for the primary action, focus rings, selection state, and nothing else. Not a tertiary brand flourish.

| Token | OKLCH | Use |
|---|---|---|
| `--accent` | `oklch(0.72 0.12 200)` | Primary action fill, selected segment fill. |
| `--accent-strong` | `oklch(0.80 0.13 200)` | Hover state of `--accent`. |
| `--accent-ink` | `oklch(0.18 0.01 200)` | Text on `--accent` fill. |
| `--accent-soft` | `oklch(0.72 0.12 200 / 0.18)` | Selected row tint. Focus ring base. |

Hue 200 is a desaturated cyan-teal. Lower chroma than the previous `#14b8a6`. Reads as confident, not neon. The accent never appears as a gradient.

### Status (functional, not decorative)

Status pairs **hue + shape + label**. Color is never the only signal. The shape glyphs are part of the StatusDot primitive.

| State | Token | OKLCH | Glyph | Meaning |
|---|---|---|---|---|
| Ready | `--status-ready` | `oklch(0.78 0.13 145)` | `●` filled | Skill plan complete, no training. |
| Training | `--status-training` | `oklch(0.85 0.13 95)` | `◐` half | Queue active, time visible. |
| Queued | `--status-queued` | `oklch(0.74 0.10 240)` | `▸` chevron | Has queue, not currently training (paused/MCT off). |
| Idle | `--status-idle` | `oklch(0.58 0.005 65)` | `○` ring | No queue, no training. |
| Error | `--status-error` | `oklch(0.70 0.18 25)` | `!` glyph | ESI fail, stale data. |

Status hues are kept at moderate chroma (≤0.18) so colorblind users distinguish by glyph and weight even when hue collapses.

## Typography

Two faces. No display font.

- **Body / UI**: `Inter` 400/500/600. Already loaded.
- **Mono / Numerals**: `JetBrains Mono` 400/500. Used for SP, ETA, locations table column, and any numeric or identifier-shaped data.

Orbitron is removed.

### Scale

A 1.25 ratio scale, six steps. Hierarchy comes from weight and size, not from color shifts.

| Step | Size / line-height | Weight | Use |
|---|---|---|---|
| `text-display` | `28px / 36px` | 600 | Page title only. One per page. |
| `text-h2` | `22px / 30px` | 600 | Section heading inside a page (rare). |
| `text-h3` | `17px / 24px` | 600 | Group header label. |
| `text-body` | `14px / 20px` | 400 | Default. Row content. |
| `text-meta` | `12px / 16px` | 500 | Labels, units, sublines. |
| `text-micro` | `11px / 14px` | 500 | Shortcut hints, breadcrumbs. |

Tabular numerals (`font-feature-settings: "tnum"`) on every numeric column. Set globally on `mono` and on a `.tabular` utility for body-set numbers.

## Spacing & rhythm

A 4px base. Allowed values: `4 / 8 / 12 / 16 / 20 / 24 / 32 / 48 / 64`. **No** `6 / 10 / 14 / 18` etc. — variance is via picking different allowed steps, not inventing in-between ones.

Density tiers:

- **Dense rows** (CharacterOverview, Mapping, Sync tables): row height 40px, vertical padding 8px, horizontal padding 16px.
- **Comfortable rows** (lists with longer text): row height 48px, vertical padding 12px.
- **Reading layout** (Settings, instructions): max-width 65ch, line-height 1.55.

## Radius

| Token | Value | Use |
|---|---|---|
| `--radius-sm` | 4px | Status dot pill, badge. |
| `--radius-md` | 6px | Buttons, inputs, segmented controls. |
| `--radius-lg` | 10px | Surface panels, popovers, modals. |

No `rounded-2xl` anywhere. No `9999px` pill toggles (the current pill `ToggleButtonGroup` is removed in favor of a 6px segmented control).

## Elevation

There is no shadow system. Layering is by surface tier (`--surface-0/1/2/3`) and by hairline rules. Shadows are reserved for two cases only: popovers (`0 6px 20px oklch(0 0 0 / 0.4)`) and modals (the same). Cards do not cast shadows.

## Motion

Functional motion only. Duration cap **180ms**. Easing: `cubic-bezier(0.2, 0, 0, 1)` (ease-out-quart).

Allowed:
- Row expand/collapse — `height` transition ≤180ms.
- Popover / drawer / modal open — `opacity` + `transform: translateY(4px)` ≤180ms.
- Focus ring fade-in — ≤120ms.
- Refresh icon rotation while loading — only if `prefers-reduced-motion: no-preference`.

Banned:
- `transform: translateY` on hover for cards and rows.
- `scale` on hover.
- Float, shimmer, glow, pulse, infinite ambient animations.
- `box-shadow` transitions.
- Spring physics with bounce or overshoot.

`prefers-reduced-motion: reduce` disables all decorative motion and shortens functional motion to ≤80ms.

## Focus & keyboard

- Visible focus ring on every interactive element. `outline: 2px solid var(--accent); outline-offset: 2px;`. No `outline: none` without an explicit replacement.
- Row focus is a 1px `--rule-strong` border-left-replacement: implemented as an inset `box-shadow: inset 2px 0 0 var(--accent)` on the row, not a colored border-left (banned).
- Keyboard shortcuts visible on hover-or-focus of the action they trigger, formatted in `text-micro` mono.

## Iconography

Material icons (`@mui/icons-material`) at `20px` default, `16px` in row context. Icon-only buttons must carry an `aria-label` and a tooltip containing the same text.

No emoji in UI chrome. The current GroupCard's `⚔️ ⛏️ 🏭 💰 🔭` set is removed; group identity is text + weight.

## Primitives

The components in `renderer/src/components/ui/` are the system. Page code composes these; it does not reach for raw MUI Card / Paper / Box with custom sx.

### `Surface`

A panel. Tier prop: `0 | 1 | 2 | 3` mapping to `--surface-*`. Optional `bordered` adds a 1px `--rule-1` border. No padding by default — caller decides density.

### `Row`

A flex row sized for dense data. Props: `interactive` (adds hover, focus, role=button), `selected`, `onClick`. Renders `inset 2px 0 0 var(--accent)` when focused via keyboard. No translate, no scale, no shadow.

### `StatusDot`

A 10px glyph that encodes state via hue + shape + label. Props: `state: 'ready' | 'training' | 'queued' | 'idle' | 'error'`, `label?` (overrides default a11y name). Pairs with a `<VisuallyHidden>` label so screen readers say "Ready" not "green dot".

### `SegmentedControl`

Replaces the current pill `ToggleButtonGroup`. Radio-group semantics. Items are icon + tooltip + accessible name. Selected segment uses `--accent` fill with `--accent-ink` text. Border radius `--radius-md`. Keyboard: arrow keys move selection; space confirms.

### `IconButton`

Wraps MUI IconButton with the right defaults. No hover scale, no hover transform. Focus ring is the system's. Tooltip is required.

### `Kbd`

Renders a keyboard shortcut hint in `text-micro` mono inside a 2px-padded surface-2 container with `--rule-1` border. Used inline next to action labels.

### `Subheader`

The new page header. Title + optional meta line + actions slot. No glass, no gradient, no animated entry. Uses `text-display` for the title.

## What this system explicitly excludes

These exist in the current codebase and are removed by this redesign. If you find yourself reaching for one, you're outside the system.

- Glassmorphism (`backdrop-filter: blur`) as a default surface.
- Gradient buttons, gradient borders, `text-gradient` clip.
- Hover-lift (`translateY(-Xpx)`) on any card or row.
- Float, shimmer, glow, pulse keyframes.
- Neumorphic shadows.
- Pill-shaped (`rounded-full`) toggle groups.
- Orbitron or any other display font.
- The `eve.blue / eve.teal / eve.orange` 50-900 ramps. Replaced by the four token families above.
- Side-stripe colored borders (banned by the shared design laws).

## Implementation notes

- Tokens are declared in `:root` in `index.css`, then mapped under `theme.extend.colors` in `tailwind.config.cjs` so `bg-surface-1`, `text-ink-2`, `border-rule-1`, `text-status-training` all work.
- The MUI theme (`Theme.jsx`) reads the same tokens via `getComputedStyle(document.documentElement).getPropertyValue('--surface-1')` at boot, so MUI components inherit them. (Alternative: hard-code the OKLCH values in Theme.jsx and keep parity with CSS by review. Pick one strategy and stay there.)
- Reduced motion: a single `@media (prefers-reduced-motion: reduce)` rule in `index.css` zeroes out `--motion-duration` and disables the refresh-icon spin.
