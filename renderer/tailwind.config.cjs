// tailwind.config.cjs

/**
 * Tokens are declared in src/index.css as CSS variables on :root.
 * Tailwind utilities below reference those variables so the system has
 * a single source of truth.
 *
 * The legacy `eve.*` ramps are retained to keep older components compiling
 * during the pilot redesign. New code MUST NOT reach for them; each touched
 * surface should be migrated to surface/ink/accent/status tokens.
 */

module.exports = {
    darkMode: 'class',
    content: [
        './index.html',
        './src/**/*.{js,jsx,ts,tsx}',
    ],
    theme: {
        extend: {
            colors: {
                // ── New restrained system ─────────────────────────────────
                surface: {
                    0: 'var(--surface-0)',
                    1: 'var(--surface-1)',
                    2: 'var(--surface-2)',
                    3: 'var(--surface-3)',
                    overlay: 'var(--surface-overlay)',
                },
                ink: {
                    1: 'var(--ink-1)',
                    2: 'var(--ink-2)',
                    3: 'var(--ink-3)',
                    4: 'var(--ink-4)',
                },
                rule: {
                    1: 'var(--rule-1)',
                    2: 'var(--rule-2)',
                    strong: 'var(--rule-strong)',
                },
                accent: {
                    DEFAULT: 'var(--accent)',
                    strong: 'var(--accent-strong)',
                    soft: 'var(--accent-soft)',
                    ink: 'var(--accent-ink)',
                },
                status: {
                    ready: 'var(--status-ready)',
                    training: 'var(--status-training)',
                    queued: 'var(--status-queued)',
                    idle: 'var(--status-idle)',
                    error: 'var(--status-error)',
                },

                // ── Legacy ramps (kept for unmigrated components) ─────────
                eve: {
                    blue: {
                        50: '#e6f3ff', 100: '#b3daff', 200: '#80c1ff', 300: '#4da8ff',
                        400: '#1a8fff', 500: '#0076e6', 600: '#005db3', 700: '#004480',
                        800: '#002b4d', 900: '#00121a',
                    },
                    teal: {
                        50: '#e6fffb', 100: '#b3fff0', 200: '#80ffe5', 300: '#4dffda',
                        400: '#1affcf', 500: '#00e6b5', 600: '#00b38c', 700: '#008063',
                        800: '#004d3a', 900: '#001a11',
                    },
                    orange: {
                        50: '#fff5e6', 100: '#ffe0b3', 200: '#ffcc80', 300: '#ffb74d',
                        400: '#ffa21a', 500: '#e68900', 600: '#b36b00', 700: '#804d00',
                        800: '#4d2f00', 900: '#1a1000',
                    },
                },
            },
            fontFamily: {
                body: ['Inter', 'system-ui', 'sans-serif'],
                mono: ['"JetBrains Mono"', 'ui-monospace', 'monospace'],
                // `display` retained as an alias for body so legacy `font-display`
                // usages render with the new face instead of Orbitron.
                display: ['Inter', 'system-ui', 'sans-serif'],
            },
            fontSize: {
                display: ['28px', { lineHeight: '36px', letterSpacing: '-0.01em', fontWeight: '600' }],
                h2: ['22px', { lineHeight: '30px', fontWeight: '600' }],
                h3: ['17px', { lineHeight: '24px', fontWeight: '600' }],
                body: ['14px', { lineHeight: '20px' }],
                meta: ['12px', { lineHeight: '16px', fontWeight: '500' }],
                micro: ['11px', { lineHeight: '14px', fontWeight: '500' }],
            },
            borderRadius: {
                sm: '4px',
                md: '6px',
                lg: '10px',
            },
            boxShadow: {
                // Reserved for popovers and modals only.
                popover: '0 6px 20px oklch(0 0 0 / 0.4)',
                modal: '0 12px 40px oklch(0 0 0 / 0.5)',
                // Selection rail: inset accent on the leading edge of a row.
                'rail-accent': 'inset 2px 0 0 var(--accent)',
            },
            transitionTimingFunction: {
                'out-quart': 'cubic-bezier(0.2, 0, 0, 1)',
            },
            transitionDuration: {
                fast: '120ms',
                base: '180ms',
            },
        },
    },
    plugins: [],
};
