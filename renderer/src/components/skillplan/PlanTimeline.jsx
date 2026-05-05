//
// New default view for Skill Plans. One ProgressRail per plan, with one
// segment per character encoding readiness state. Plan rows align with the
// FROZEN column width used by PlanMatrix so the eye reads them as the same
// data, just rotated.

import { useMemo } from 'react';
import PropTypes from 'prop-types';
import EveTypeIcon from '../ui/EveTypeIcon.jsx';
import ProgressRail from '../ui/ProgressRail.jsx';
import { computePlanReadiness } from '../../hooks/usePlanReadiness.js';

const PlanTimeline = ({ characters, skillPlans, conversions, filter }) => {
    const planNames = useMemo(() => {
        const q = filter.trim().toLowerCase();
        return Object.keys(skillPlans)
            .filter((p) => (q ? p.toLowerCase().includes(q) : true))
            .sort((a, b) => a.localeCompare(b));
    }, [skillPlans, filter]);

    const sortedChars = useMemo(
        () => [...characters].sort((a, b) =>
            (a.Character?.CharacterName || '').localeCompare(b.Character?.CharacterName || ''),
        ),
        [characters],
    );

    if (Object.keys(skillPlans).length === 0) {
        return (
            <div className="mt-12 text-center">
                <p className="text-body text-ink-2">No skill plans yet.</p>
                <p className="mt-1 text-meta text-ink-3">Add a plan from the header to populate the timeline.</p>
            </div>
        );
    }

    if (planNames.length === 0) {
        return (
            <div className="mt-12 text-center">
                <p className="text-body text-ink-2">No skill plans match.</p>
            </div>
        );
    }

    return (
        <div role="table" aria-label="Skill plan timeline" className="rounded-lg border border-rule-1 bg-surface-1 overflow-hidden">
            {planNames.map((p, idx) => {
                const segments = sortedChars.map((ch) => {
                    const r = computePlanReadiness(ch, p);
                    return {
                        key: String(ch.Character?.CharacterID),
                        weight: 1,
                        state: r.state,
                        label: `${ch.Character?.CharacterName} · ${r.state}${r.eta ? ' · ' + r.eta : ''}${r.missingCount ? ` · ${r.missingCount} missing` : ''}`,
                    };
                });
                const isLast = idx === planNames.length - 1;
                return (
                    <div
                        key={p}
                        role="row"
                        className={`grid items-center gap-3 px-4 h-14 grid-cols-[28px_minmax(220px,1.6fr)_minmax(0,3fr)] ${isLast ? '' : 'border-b border-rule-1'}`}
                    >
                        <EveTypeIcon name={p} conversions={conversions} size={6} />
                        <span className="text-body text-ink-1 truncate" title={p}>{p}</span>
                        <ProgressRail segments={segments} ariaLabel={`${p} readiness`} height={10} />
                    </div>
                );
            })}
        </div>
    );
};

PlanTimeline.propTypes = {
    characters: PropTypes.array.isRequired,
    skillPlans: PropTypes.object.isRequired,
    conversions: PropTypes.object.isRequired,
    filter: PropTypes.string.isRequired,
};

export default PlanTimeline;
