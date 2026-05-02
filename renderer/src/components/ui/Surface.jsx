import React from 'react';
import PropTypes from 'prop-types';

const TIER_BG = {
    0: 'bg-surface-0',
    1: 'bg-surface-1',
    2: 'bg-surface-2',
    3: 'bg-surface-3',
};

/**
 * A panel. The lowest-level container in the new system.
 * - `tier` picks the surface layer (0–3).
 * - `bordered` adds a 1px hairline rule. Off by default.
 * - `radius` is one of the design-system radii. Off by default for table-style content.
 */
const Surface = React.forwardRef(function Surface(
    { tier = 1, bordered = false, radius = 'none', as: Tag = 'div', className = '', children, ...rest },
    ref,
) {
    const radiusClass =
        radius === 'sm' ? 'rounded-sm' :
        radius === 'md' ? 'rounded-md' :
        radius === 'lg' ? 'rounded-lg' :
        '';
    const borderClass = bordered ? 'border border-rule-1' : '';
    return (
        <Tag
            ref={ref}
            className={`${TIER_BG[tier] || TIER_BG[1]} ${radiusClass} ${borderClass} ${className}`.trim()}
            {...rest}
        >
            {children}
        </Tag>
    );
});

Surface.propTypes = {
    tier: PropTypes.oneOf([0, 1, 2, 3]),
    bordered: PropTypes.bool,
    radius: PropTypes.oneOf(['none', 'sm', 'md', 'lg']),
    as: PropTypes.elementType,
    className: PropTypes.string,
    children: PropTypes.node,
};

export default Surface;
