import PropTypes from 'prop-types';

/**
 * Renders a keyboard shortcut hint inline with action labels.
 * Uses surface-2 + rule-1 + mono micro type. No animation, no glow.
 */
const Kbd = ({ children, className = '' }) => (
    <kbd
        className={[
            'inline-flex items-center justify-center',
            'h-[18px] min-w-[18px] px-1',
            'rounded-sm border border-rule-1 bg-surface-2',
            'font-mono text-micro text-ink-2',
            className,
        ].join(' ')}
    >
        {children}
    </kbd>
);

Kbd.propTypes = {
    children: PropTypes.node.isRequired,
    className: PropTypes.string,
};

export default Kbd;
