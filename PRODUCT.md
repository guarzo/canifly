# Product

## Register

product

## Users

Hardcore EVE Online players — multi-account, often multi-character-per-account — who open CanIFly **occasionally rather than daily**: when planning the next training arc, deciding what to inject, picking a ship to fly tonight, or syncing in-game settings across alts after an EVE patch. Sessions are intent-driven and finite ("can my Loki alt fly this doctrine yet?", "sync my UI to all 14 characters") rather than ambient. The app sits *next to* EVE on a desk monitor, not inside the game.

The job to be done, in order of frequency:

1. Answer "can I fly X?" across N characters at a glance.
2. Plan the next skill plan — see what's missing, what's queued, total training time.
3. Map character files to user files and sync EVE settings without losing per-account UI tweaks.
4. Add, group, and reorganize characters by account / role / location.

## Product Purpose

CanIFly is a desk-side instrument for managing an EVE skill and settings portfolio. It exists because EVE's own UI is bad at multi-character planning and bad at settings sync, and because spreadsheets are worse. Success is measured in *seconds-to-answer*: the user opens it with a question, gets the answer, acts on it, and closes the app.

It is **not** a companion app you live in. It is a tool you reach for, use precisely, and put down.

## Brand Personality

Three words: **precise, calm, expert.**

Voice is the voice of a flight manual or a senior FC's notes — terse, factual, no exclamation marks, no marketing energy. "Sabre: ready in 3d 14h" beats "🚀 You're almost there!". Labels are real labels, not motivational phrases.

The interface should feel like Linear, Things 3, Raycast, Soulver — tools that respect a user's time and attention. Chrome is borders and weight contrast, not glows and gradients. Information density is high where the user is comparing things (character × ship matrix, skill plan progress) and generous where they are reading or deciding.

The app acknowledges EVE's context (dark surface, occasional monospace for data) but does not cosplay as a game HUD.

## Anti-references

Hard nos, ranked by current risk:

1. **Generic SaaS dashboard cream.** No Stripe-clone soft-shadow cards, no rounded-2xl-everywhere, no identical icon-and-heading card grids, no indigo-on-white. If a screenshot could be retitled "billing dashboard" without anyone noticing, redo it.
2. **Bootstrap / MUI-default look.** No untouched MUI Cards, no default blue buttons, no zero-identity surface. Material Design's default elevation system is banned; chrome is custom.
3. **Sci-fi game HUD (the current app).** No Orbitron, no glassmorphism / `backdrop-filter: blur` as a default, no gradient buttons, no glow shadows, no hover-lift `translateY` cards, no float/shimmer ambient animation, no neumorphic shadows. The current `Theme.jsx` is the anti-reference.
4. **Crypto / trading neon-on-black.** No bright green ↔ red on pure black, no candlestick aesthetics, no scrolling tickers.
5. **EVE-themed reflex.** No starfields, no nebula gradients, no Eve Gate iconography. The fact that this is for a space game is not the design.

## Design Principles

1. **Seconds to answer.** Every screen is judged by how fast the user gets the one fact they came for. Information density beats whitespace where comparison happens; whitespace beats density where reading happens. Never both at once.
2. **Borders, not glows.** Hierarchy comes from weight, scale, and 1px rules tinted toward the surface hue. No gradients, no shadows-as-decoration, no blur.
3. **Status without color alone.** Skill state (ready / pending / missing / queued) carries shape, weight, and text in addition to color so it survives color-blindness and grayscale printing.
4. **Keyboard-first, mouse-fine.** Every primary action is reachable from the keyboard with a visible focus ring. The app is usable as a single-handed instrument while EVE has the other hand.
5. **Restraint is the brand.** When a flourish (animation, color, illustration, gradient) would make this feel "designed," cut it. The product *is* design discipline applied to an EVE multibox problem.

## Accessibility & Inclusion

- **WCAG 2.1 AA** for all text and interactive states, including focus rings (3:1 against adjacent surface) and non-text indicators.
- **`prefers-reduced-motion`** disables all decorative motion (page-load fades, hover transitions on cards). Functional motion (popover open, drawer slide) drops to <120ms.
- **Color-blind safe status.** Ready / pending / missing / queued are encoded with at least two of: hue, shape (filled dot, ring, slash, bar), and weight. Color is never the sole signal.
- **Keyboard navigation.** Tab order matches reading order; every modal trap-focuses; Esc closes; primary actions have visible shortcut hints where they exist.
- **Tabular numerals** for all training times and skill points so scanning a column lines up.
