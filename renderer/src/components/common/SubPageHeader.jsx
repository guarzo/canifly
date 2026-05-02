import { useState } from 'react';
import PropTypes from 'prop-types';

/**
 * SubPageHeader — calm typographic page header.
 * Consumed by Mapping, Sync, SkillPlans, Settings, etc. The CharacterOverview
 * pilot uses the new `Subheader` primitive directly; this shim keeps the
 * unmigrated routes on the new system without forcing every page to change.
 *
 * Public props are unchanged from the previous implementation:
 * - title: string (required)
 * - instructions: string (optional)
 * - storageKey: string (required) — persists the show/hide tip preference
 */
const SubPageHeader = ({ title, instructions, storageKey }) => {
    const [showTip, setShowTip] = useState(() => {
        try {
            const stored = localStorage.getItem(storageKey);
            return stored === null ? false : JSON.parse(stored);
        } catch {
            return false;
        }
    });

    const toggle = () => {
        const next = !showTip;
        setShowTip(next);
        try { localStorage.setItem(storageKey, JSON.stringify(next)); } catch { /* swallow */ }
    };

    return (
        <header className="max-w-[1280px] mx-auto mb-5">
            <div className="flex items-baseline justify-between gap-4">
                <h1 className="text-display text-ink-1 truncate">{title}</h1>
                {instructions ? (
                    <button
                        type="button"
                        onClick={toggle}
                        aria-label={showTip ? 'Hide help' : 'Show help'}
                        aria-expanded={showTip}
                        className="inline-flex h-8 w-8 items-center justify-center rounded-md text-ink-3 hover:bg-surface-2 hover:text-ink-1 transition-colors duration-fast ease-out-quart"
                    >
                        <span aria-hidden className="font-mono text-meta">?</span>
                    </button>
                ) : null}
            </div>
            {instructions && showTip ? (
                <p className="mt-3 max-w-[65ch] text-body text-ink-3">{instructions}</p>
            ) : null}
            <div className="mt-4 h-px w-full bg-rule-1" />
        </header>
    );
};

SubPageHeader.propTypes = {
    title: PropTypes.string.isRequired,
    instructions: PropTypes.string,
    storageKey: PropTypes.string.isRequired,
};

export default SubPageHeader;
