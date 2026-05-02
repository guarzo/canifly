import PropTypes from 'prop-types';

/**
 * SkillProgress — flat ring matched to the design system.
 * No gradient, no glow, no pulse. Hue + weight encode state; the percent
 * is set in tabular mono so columns of these line up.
 */
const SIZE = 128;
const RADIUS = 45;
const STROKE = 8;
const CENTER = SIZE / 2;
const CIRCUMFERENCE = 2 * Math.PI * RADIUS;

const SkillProgress = ({ skill, progress, isTraining }) => {
    const offset = CIRCUMFERENCE - (progress / 100) * CIRCUMFERENCE;
    const stroke = isTraining
        ? 'var(--status-training)'
        : 'var(--accent)';

    return (
        <div
            className="relative"
            style={{ width: SIZE, height: SIZE }}
            role="img"
            aria-label={`${skill.name}: ${progress}% complete`}
        >
            <svg
                className="-rotate-90"
                width={SIZE}
                height={SIZE}
            >
                <circle
                    cx={CENTER}
                    cy={CENTER}
                    r={RADIUS}
                    stroke="var(--rule-1)"
                    strokeWidth={STROKE}
                    fill="none"
                />
                <circle
                    cx={CENTER}
                    cy={CENTER}
                    r={RADIUS}
                    stroke={stroke}
                    strokeWidth={STROKE}
                    fill="none"
                    strokeLinecap="round"
                    strokeDasharray={CIRCUMFERENCE}
                    strokeDashoffset={offset}
                    style={{
                        transition:
                            'stroke-dashoffset var(--motion-duration, 180ms) cubic-bezier(0.2, 0, 0, 1)',
                    }}
                />
            </svg>

            <div className="absolute inset-0 flex flex-col items-center justify-center">
                <span className="text-h2 text-ink-1 font-mono tabular">
                    {progress}%
                </span>
                <span className="text-meta text-ink-3">{skill.name}</span>
            </div>
        </div>
    );
};

SkillProgress.propTypes = {
    skill: PropTypes.shape({ name: PropTypes.string.isRequired }).isRequired,
    progress: PropTypes.number.isRequired,
    isTraining: PropTypes.bool,
};

export default SkillProgress;
