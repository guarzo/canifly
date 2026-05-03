import PropTypes from 'prop-types';

/**
 * SegmentedControl — keyboard-navigable, radio-group-semantic alternative to
 * MUI's pill ToggleButtonGroup. Each option is an icon + tooltip + accessible
 * name. Selected segment fills with --accent.
 *
 * Options shape: { value, label, icon? }.
 */
const SegmentedControl = ({ value, onChange, options, ariaLabel, className = '' }) => {
    const handleKey = (e, idx) => {
        if (e.key === 'ArrowRight' || e.key === 'ArrowDown') {
            e.preventDefault();
            const next = options[(idx + 1) % options.length];
            onChange(next.value);
        } else if (e.key === 'ArrowLeft' || e.key === 'ArrowUp') {
            e.preventDefault();
            const prev = options[(idx - 1 + options.length) % options.length];
            onChange(prev.value);
        }
    };

    return (
        <div
            role="radiogroup"
            aria-label={ariaLabel}
            className={`inline-flex items-center rounded-md border border-rule-1 bg-surface-1 p-0.5 ${className}`}
        >
            {options.map((opt, idx) => {
                const selected = opt.value === value;
                return (
                    <button
                        key={opt.value}
                        type="button"
                        role="radio"
                        aria-checked={selected}
                        aria-label={opt.label}
                        title={opt.label}
                        tabIndex={selected ? 0 : -1}
                        onClick={() => onChange(opt.value)}
                        onKeyDown={(e) => handleKey(e, idx)}
                        className={[
                            'inline-flex items-center justify-center gap-1.5',
                            'h-7 px-2.5 rounded-sm text-meta',
                            'transition-colors duration-fast ease-out-quart',
                            selected
                                ? 'bg-accent text-accent-ink'
                                : 'text-ink-2 hover:bg-surface-2 hover:text-ink-1',
                        ].join(' ')}
                    >
                        {opt.icon ? <span aria-hidden className="inline-flex">{opt.icon}</span> : null}
                        {opt.showLabel ? <span>{opt.label}</span> : null}
                    </button>
                );
            })}
        </div>
    );
};

SegmentedControl.propTypes = {
    value: PropTypes.string.isRequired,
    onChange: PropTypes.func.isRequired,
    options: PropTypes.arrayOf(PropTypes.shape({
        value: PropTypes.string.isRequired,
        label: PropTypes.string.isRequired,
        icon: PropTypes.node,
        showLabel: PropTypes.bool,
    })).isRequired,
    ariaLabel: PropTypes.string.isRequired,
    className: PropTypes.string,
};

export default SegmentedControl;
