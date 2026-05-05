import PropTypes from 'prop-types';

const SwatchBridge = ({ containerRef, pairs, width, height, className = '' }) => {
    if (!pairs || pairs.length === 0) {
        return <svg aria-hidden width={width} height={height} className={className} />;
    }
    const origin = containerRef.current?.getBoundingClientRect?.() || { left: 0, top: 0 };

    return (
        <svg
            aria-hidden
            width={width}
            height={height}
            className={`pointer-events-none ${className}`}
            style={{ position: 'absolute', inset: 0 }}
        >
            {pairs.map((pair) => {
                const fromX = pair.fromRect.right - origin.left;
                const fromY = pair.fromRect.top + pair.fromRect.height / 2 - origin.top;
                const toX = pair.toRect.left - origin.left;
                const toY = pair.toRect.top + pair.toRect.height / 2 - origin.top;
                const dx = toX - fromX;
                const c1x = fromX + dx * 0.5;
                const c2x = fromX + dx * 0.5;
                const d = `M ${fromX} ${fromY} C ${c1x} ${fromY}, ${c2x} ${toY}, ${toX} ${toY}`;
                return (
                    <path
                        key={pair.key}
                        d={d}
                        fill="none"
                        stroke={pair.color || 'var(--rule-1)'}
                        strokeWidth={1}
                        strokeLinecap="round"
                        opacity={0.85}
                    />
                );
            })}
        </svg>
    );
};

SwatchBridge.propTypes = {
    containerRef: PropTypes.shape({ current: PropTypes.any }).isRequired,
    pairs: PropTypes.arrayOf(PropTypes.shape({
        key: PropTypes.string.isRequired,
        color: PropTypes.string,
        fromRect: PropTypes.object.isRequired,
        toRect: PropTypes.object.isRequired,
    })).isRequired,
    width: PropTypes.number.isRequired,
    height: PropTypes.number.isRequired,
    className: PropTypes.string,
};

export default SwatchBridge;
