import LoginButton from '../components/landing/LoginButton.jsx';

/**
 * Landing — pre-auth screen.
 * Calm, single-purpose: identify the app, explain it in one line, log in.
 * No particles, no gradient text, no scroll-driven animation.
 */
const Landing = () => (
    <div className="min-h-screen flex items-center justify-center bg-surface-0 px-6">
        <main className="w-full max-w-md">
            <header className="mb-8">
                <p className="text-meta text-ink-3 uppercase tracking-[0.18em] mb-3">CanIFly</p>
                <h1 className="text-display text-ink-1 mb-3">EVE pilot console.</h1>
                <p className="text-body text-ink-2 max-w-[44ch]">
                    Track every character across every account. Plan training, settle skill
                    plans, and sync EVE settings without losing per-alt UI tweaks.
                </p>
            </header>

            <div className="h-px w-12 bg-rule-2 mb-8" aria-hidden />

            <div className="mb-10">
                <LoginButton />
            </div>

            <ul className="grid grid-cols-3 gap-x-6 text-meta text-ink-3">
                <li>
                    <span className="block text-ink-1 font-mono tabular text-h3 mb-1">∞</span>
                    accounts
                </li>
                <li>
                    <span className="block text-ink-1 font-mono tabular text-h3 mb-1">≤s</span>
                    to answer
                </li>
                <li>
                    <span className="block text-ink-1 font-mono tabular text-h3 mb-1">ESI</span>
                    direct
                </li>
            </ul>
        </main>
    </div>
);

export default Landing;
