const STATE_VAR = {
    ready: 'var(--status-ready)',
    training: 'var(--status-training)',
    queued: 'var(--status-queued)',
    idle: 'var(--status-idle)',
    error: 'var(--status-error)',
};

/**
 * ProgressRail — a 1-row segmented bar. Each segment width is proportional
 * to its `weight`; color is driven by `state` (ready/training/queued/idle/
 * error). The rail itself is the rule-1 hairline; segments fill it. Used for
 * the Skill Plans timeline and any future training-queue visualization.
 */
const ProgressRail = ({ segments, ariaLabel, height = 6, className = '' }) => {
    const total = segments.reduce((s, x) => s + (x.weight || 0), 0) || 1;
    return (
        <div
            role="img"
            aria-label={ariaLabel}
            className={`relative w-full overflow-hidden rounded-sm border border-rule-1 bg-surface-2 ${className}`}
            style={{ height }}
        >
            <div className="flex h-full w-full">
                {segments.map((seg) => {
                    const pct = ((seg.weight || 0) / total) * 100;
                    return (
                        <div
                            key={seg.key}
                            data-segment={seg.key}
                            title={seg.label}
                            style={{
                                width: `${pct}%`,
                                backgroundColor: STATE_VAR[seg.state] || STATE_VAR.idle,
                            }}
                        />
                    );
                })}
            </div>
        </div>
    );
};

export default ProgressRail;
