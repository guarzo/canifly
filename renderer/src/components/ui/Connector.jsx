/**
 * Connector — a single hairline SVG path between two points. The signature
 * informational primitive of the redesign. Callers compute the `from` and
 * `to` points (relative to the connector's own coordinate space, typically
 * the parent's bounding box) and pass them in. We render a smooth cubic
 * bezier so horizontal-then-horizontal pairings get a calm S-curve.
 *
 * Use informationally only: paired rows in Mapping, source→destination in
 * Sync, segments under timeline plans. Never as a decorative divider.
 */
const Connector = ({
    from,
    to,
    color = 'var(--rule-1)',
    strokeWidth = 1,
    active = false,
    withArrow = false,
    ariaLabel,
    width,
    height,
    className = '',
}) => {
    const dx = to.x - from.x;
    const c1 = { x: from.x + dx * 0.5, y: from.y };
    const c2 = { x: from.x + dx * 0.5, y: to.y };
    const d = `M ${from.x} ${from.y} C ${c1.x} ${c1.y}, ${c2.x} ${c2.y}, ${to.x} ${to.y}`;

    const role = ariaLabel ? 'img' : 'presentation';
    const ariaProps = ariaLabel ? { 'aria-label': ariaLabel } : { 'aria-hidden': true };

    return (
        <svg
            role={role}
            {...ariaProps}
            width={width}
            height={height}
            className={className}
            style={{ overflow: 'visible' }}
        >
            <path
                d={d}
                fill="none"
                stroke={color}
                strokeWidth={active ? strokeWidth + 0.5 : strokeWidth}
                strokeLinecap="round"
                opacity={active ? 1 : 0.7}
            />
            {withArrow ? (
                <polygon
                    points={`${to.x},${to.y} ${to.x - 5},${to.y - 3} ${to.x - 5},${to.y + 3}`}
                    fill={color}
                    opacity={active ? 1 : 0.7}
                />
            ) : null}
        </svg>
    );
};

export default Connector;
