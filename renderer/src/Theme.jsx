// Theme.jsx
//
// MUI theme aligned with the new restrained system. Tokens are duplicated as
// raw OKLCH values here so MUI can read them at render time without needing a
// CSS-variable bridge. Keep this in sync with `:root` in src/index.css —
// DESIGN.md is the source of truth.

import { createTheme } from '@mui/material/styles';

// ── Token mirror (keep visually in sync with :root in src/index.css) ──────
// MUI v6 cannot parse oklch() in palette decomposition (it's used by
// `lighten`/`darken`/`augmentColor`). CSS variables carry the real OKLCH
// values — these hex constants are sRGB approximations for MUI's palette
// engine and for components that read theme tokens directly. Visual parity
// is close but not exact; CSS-driven surfaces win when there's a mismatch.
const SURFACE_0 = '#1f1c19';
const SURFACE_1 = '#2a2724';
const SURFACE_2 = '#33302c';
const SURFACE_3 = '#3d3a35';
const INK_1 = '#f5f3ee';
const INK_2 = '#c4c0b8';
const INK_3 = '#878177';
const RULE_1 = '#48443f';
const RULE_2 = '#5b5750';
const ACCENT = '#3fb6c4';
const ACCENT_STRONG = '#5fcfdc';
const ACCENT_INK = '#0e1f22';
const STATUS_ERROR = '#dc7a5e';
const STATUS_READY = '#67ca8d';
const STATUS_TRAINING = '#dbc55b';
const STATUS_QUEUED = '#7caad6';

const theme = createTheme({
    palette: {
        mode: 'dark',
        primary: { main: ACCENT, light: ACCENT_STRONG, dark: ACCENT, contrastText: ACCENT_INK },
        secondary: { main: STATUS_QUEUED, contrastText: INK_1 },
        error: { main: STATUS_ERROR },
        warning: { main: STATUS_TRAINING },
        info: { main: STATUS_QUEUED },
        success: { main: STATUS_READY },
        background: { default: SURFACE_0, paper: SURFACE_1 },
        text: { primary: INK_1, secondary: INK_2, disabled: INK_3 },
        divider: RULE_1,
    },
    typography: {
        fontFamily: '"Inter", system-ui, sans-serif',
        // Headings inherit the body face — no display font.
        h1: { fontWeight: 600, letterSpacing: '-0.01em' },
        h2: { fontWeight: 600 },
        h3: { fontWeight: 600 },
        h4: { fontWeight: 600 },
        h5: { fontWeight: 600 },
        h6: { fontWeight: 600 },
        body1: { fontSize: '14px', lineHeight: 1.43 },
        body2: { fontSize: '12px', lineHeight: 1.33 },
        button: { fontWeight: 500, textTransform: 'none', letterSpacing: 0 },
    },
    shape: {
        borderRadius: 6,
    },
    components: {
        MuiCssBaseline: {
            styleOverrides: {
                body: {
                    backgroundColor: SURFACE_0,
                    color: INK_2,
                },
            },
        },
        MuiCard: {
            styleOverrides: {
                root: {
                    backgroundImage: 'none',
                    backgroundColor: SURFACE_1,
                    border: `1px solid ${RULE_1}`,
                    boxShadow: 'none',
                    transition: 'none',
                    '&:hover': {
                        // No translate, no shadow, no glow.
                        transform: 'none',
                        boxShadow: 'none',
                        borderColor: RULE_2,
                    },
                },
            },
        },
        MuiPaper: {
            styleOverrides: {
                root: {
                    backgroundImage: 'none',
                    backgroundColor: SURFACE_1,
                    backdropFilter: 'none',
                },
            },
        },
        MuiButton: {
            defaultProps: { disableElevation: true, disableRipple: false },
            styleOverrides: {
                root: {
                    textTransform: 'none',
                    fontWeight: 500,
                    borderRadius: 6,
                    padding: '6px 12px',
                    boxShadow: 'none',
                    transition: 'background-color var(--motion-duration-fast, 120ms) var(--motion-ease, ease)',
                    '&:hover': { boxShadow: 'none', transform: 'none' },
                },
                contained: {
                    background: ACCENT,
                    color: ACCENT_INK,
                    '&:hover': { background: ACCENT_STRONG },
                },
                outlined: {
                    borderWidth: 1,
                    borderColor: RULE_1,
                    color: INK_1,
                    '&:hover': { borderWidth: 1, borderColor: RULE_2, backgroundColor: SURFACE_2 },
                },
                text: {
                    color: INK_2,
                    '&:hover': { backgroundColor: SURFACE_2 },
                },
            },
        },
        MuiIconButton: {
            styleOverrides: {
                root: {
                    color: INK_2,
                    transition: 'background-color var(--motion-duration-fast, 120ms) var(--motion-ease, ease), color var(--motion-duration-fast, 120ms) var(--motion-ease, ease)',
                    '&:hover': {
                        transform: 'none',
                        backgroundColor: SURFACE_2,
                        color: INK_1,
                    },
                },
            },
        },
        MuiTextField: {
            styleOverrides: {
                root: {
                    '& .MuiOutlinedInput-root': {
                        backgroundColor: SURFACE_1,
                        backdropFilter: 'none',
                        '& fieldset': { borderColor: RULE_1 },
                        '&:hover fieldset': { borderColor: RULE_2 },
                        '&.Mui-focused fieldset': { borderColor: ACCENT, borderWidth: 1 },
                    },
                },
            },
        },
        MuiSelect: {
            styleOverrides: {
                root: {
                    backgroundColor: SURFACE_1,
                    backdropFilter: 'none',
                    '& .MuiOutlinedInput-notchedOutline': { borderColor: RULE_1 },
                    '&:hover .MuiOutlinedInput-notchedOutline': { borderColor: RULE_2 },
                    '&.Mui-focused .MuiOutlinedInput-notchedOutline': { borderColor: ACCENT, borderWidth: 1 },
                },
            },
        },
        MuiListItem: {
            styleOverrides: {
                root: {
                    borderRadius: 4,
                    transition: 'background-color var(--motion-duration-fast, 120ms) var(--motion-ease, ease)',
                    '&:hover': { backgroundColor: SURFACE_2 },
                },
            },
        },
        MuiTooltip: {
            styleOverrides: {
                tooltip: {
                    backgroundColor: SURFACE_3,
                    color: INK_1,
                    border: `1px solid ${RULE_1}`,
                    borderRadius: 4,
                    padding: '6px 8px',
                    fontSize: '12px',
                    fontWeight: 500,
                    backdropFilter: 'none',
                },
                arrow: { color: SURFACE_3 },
            },
        },
        MuiDivider: {
            styleOverrides: { root: { borderColor: RULE_1 } },
        },
        MuiToggleButton: {
            styleOverrides: {
                root: {
                    border: `1px solid ${RULE_1}`,
                    color: INK_2,
                    textTransform: 'none',
                    '&.Mui-selected': {
                        backgroundColor: ACCENT,
                        color: ACCENT_INK,
                        '&:hover': { backgroundColor: ACCENT_STRONG },
                    },
                },
            },
        },
    },
});

export default theme;
