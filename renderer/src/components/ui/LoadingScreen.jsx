/**
 * LoadingScreen — calm centered placeholder.
 * No glow, no spinning logo, no gradient. A small pulsing dot is enough
 * to signal liveness; everything else is restraint.
 */
const LoadingScreen = ({ message = 'Loading…' }) => (
    <div
        role="status"
        aria-live="polite"
        aria-busy="true"
        className="flex min-h-screen flex-col items-center justify-center bg-surface-0 px-6"
    >
        <div className="flex items-center gap-3">
            <span
                aria-hidden
                className="inline-block h-2 w-2 rounded-sm bg-accent animate-skeleton-pulse"
            />
            <span className="text-meta text-ink-3 tabular">{message}</span>
        </div>
    </div>
);

export default LoadingScreen;
