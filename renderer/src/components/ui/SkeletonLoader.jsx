import PropTypes from 'prop-types';

/**
 * SkeletonLoader — calm placeholders. No gradient sweep, no shimmer.
 * Pulses opacity gently; reduced-motion users see a static block.
 */
const VARIANTS = {
    text: 'h-4 w-full',
    title: 'h-7 w-3/4',
    avatar: 'h-10 w-10 rounded-sm',
    card: 'h-24 w-full',
    button: 'h-8 w-20',
    row: 'h-10 w-full',
};

const SkeletonLoader = ({ variant = 'text', className = '' }) => (
    <div
        role="presentation"
        aria-hidden
        className={`bg-surface-2 rounded-sm animate-skeleton-pulse ${VARIANTS[variant] || VARIANTS.text} ${className}`}
    />
);

SkeletonLoader.propTypes = {
    variant: PropTypes.oneOf(Object.keys(VARIANTS)),
    className: PropTypes.string,
};

// Skeleton matched to the new CharacterOverview row geometry.
// Same height (40px) and column template so paint-to-data swap doesn't reflow.
export const CharacterRowSkeleton = () => (
    <div className="grid grid-cols-[28px_minmax(180px,1.4fr)_1fr_84px_92px_minmax(140px,1fr)_minmax(120px,0.8fr)_28px] gap-3 px-4 items-center h-10 border-b border-rule-1">
        <SkeletonLoader className="!w-2.5 !h-2.5 rounded-full" />
        <div className="flex items-center gap-2.5">
            <SkeletonLoader className="!h-6 !w-6 rounded-sm" />
            <SkeletonLoader className="!h-3.5 !w-32" />
        </div>
        <SkeletonLoader className="!h-3.5 !w-24" />
        <SkeletonLoader className="!h-3.5 !w-12 ml-auto" />
        <SkeletonLoader className="!h-3.5 !w-16 ml-auto" />
        <SkeletonLoader className="!h-3.5 !w-28" />
        <SkeletonLoader className="!h-3.5 !w-20" />
        <span />
    </div>
);

// Group block of 3 rows, used while data loads.
export const AccountCardSkeleton = () => (
    <div className="rounded-lg border border-rule-1 bg-surface-1 overflow-hidden">
        <div className="px-4 py-2 border-b border-rule-1 bg-surface-2 flex justify-between">
            <SkeletonLoader className="!h-4 !w-32" />
            <SkeletonLoader className="!h-3 !w-24" />
        </div>
        {[0, 1, 2].map((i) => <CharacterRowSkeleton key={i} />)}
    </div>
);

// Single dense row skeleton (used inside other lists).
export const CharacterItemSkeleton = () => <CharacterRowSkeleton />;

// Skill plan table skeleton — matches the row pattern of the new system.
export const SkillPlanTableSkeleton = () => (
    <div className="space-y-2">
        <SkeletonLoader variant="title" className="mb-3 !w-48" />
        {[0, 1, 2, 3, 4].map((i) => (
            <div key={i} className="grid grid-cols-[1fr_120px_120px_80px] gap-3 items-center h-9">
                <SkeletonLoader className="!h-3.5 !w-2/3" />
                <SkeletonLoader className="!h-3.5 !w-20" />
                <SkeletonLoader className="!h-3.5 !w-20" />
                <SkeletonLoader className="!h-3.5 !w-12" />
            </div>
        ))}
    </div>
);

export default SkeletonLoader;
