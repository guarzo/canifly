import React from 'react';
/**
 * Subheader — the page header in the new system.
 * Replaces the glass + gradient SubPageHeader on migrated routes.
 * - `title` is required.
 * - `meta` is the secondary line (e.g. "14 characters across 6 accounts · 1.2B SP").
 * - `actions` is a slot for icon buttons / segmented controls on the right.
 * - `tip` is an optional one-liner help, hidden by default; toggled by the help icon.
 */
const Subheader = ({ title, meta, actions, tip }) => {
    const [showTip, setShowTip] = React.useState(false);
    return (
        <header className="mb-5">
            <div className="flex items-baseline justify-between gap-4">
                <div className="min-w-0">
                    <h1 className="text-display text-ink-1 truncate">{title}</h1>
                    {meta ? (
                        <p className="mt-1 text-meta text-ink-3 tabular">
                            {meta}
                        </p>
                    ) : null}
                </div>
                <div className="flex items-center gap-1.5 shrink-0">
                    {actions}
                    {tip ? (
                        <button
                            type="button"
                            aria-label={showTip ? 'Hide help' : 'Show help'}
                            aria-expanded={showTip}
                            onClick={() => setShowTip((v) => !v)}
                            className="inline-flex h-8 w-8 items-center justify-center rounded-md text-ink-3 hover:bg-surface-2 hover:text-ink-1 transition-colors duration-fast ease-out-quart"
                        >
                            <span aria-hidden className="font-mono text-meta">?</span>
                        </button>
                    ) : null}
                </div>
            </div>
            {tip && showTip ? (
                <p className="mt-3 max-w-[65ch] text-body text-ink-3">{tip}</p>
            ) : null}
            <div className="mt-4 h-px w-full bg-rule-1" />
        </header>
    );
};

export default Subheader;
