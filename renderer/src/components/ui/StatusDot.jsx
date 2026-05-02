import PropTypes from 'prop-types';

const STATE_LABEL = {
    ready: 'Ready',
    training: 'Training',
    queued: 'Queued',
    idle: 'Idle',
    error: 'Error',
};

// Hue is one signal; shape is another. Color-blind users distinguish on shape alone.
function Glyph({ state }) {
    const common = { width: 10, height: 10, viewBox: '0 0 10 10', 'aria-hidden': true };
    switch (state) {
        case 'ready':
            return <svg {...common}><circle cx="5" cy="5" r="4" fill="var(--status-ready)" /></svg>;
        case 'training':
            // Half-filled disc.
            return (
                <svg {...common}>
                    <circle cx="5" cy="5" r="4" fill="none" stroke="var(--status-training)" strokeWidth="1.4" />
                    <path d="M5 1 a4 4 0 0 1 0 8 z" fill="var(--status-training)" />
                </svg>
            );
        case 'queued':
            return (
                <svg {...common}>
                    <path d="M3 2 L7 5 L3 8 Z" fill="var(--status-queued)" />
                </svg>
            );
        case 'idle':
            return <svg {...common}><circle cx="5" cy="5" r="3.6" fill="none" stroke="var(--status-idle)" strokeWidth="1.4" /></svg>;
        case 'error':
            return (
                <svg {...common}>
                    <circle cx="5" cy="5" r="4" fill="var(--status-error)" />
                    <path d="M5 2.4 L5 5.8 M5 7 L5 7.6" stroke="oklch(0.18 0.01 25)" strokeWidth="1.2" strokeLinecap="round" />
                </svg>
            );
        default:
            return <svg {...common}><circle cx="5" cy="5" r="3.6" fill="none" stroke="var(--status-idle)" strokeWidth="1.4" /></svg>;
    }
}

Glyph.propTypes = { state: PropTypes.string.isRequired };

/**
 * StatusDot encodes character/skill state with hue + shape + accessible label.
 * - `state`: ready | training | queued | idle | error.
 * - `label`: optional override of the screen-reader text.
 */
const StatusDot = ({ state = 'idle', label, className = '' }) => {
    const a11y = label || STATE_LABEL[state] || 'Status';
    return (
        <span
            className={`inline-flex items-center justify-center ${className}`}
            role="img"
            aria-label={a11y}
        >
            <Glyph state={state} />
        </span>
    );
};

StatusDot.propTypes = {
    state: PropTypes.oneOf(['ready', 'training', 'queued', 'idle', 'error']),
    label: PropTypes.string,
    className: PropTypes.string,
};

export default StatusDot;
