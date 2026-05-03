// Character × Skill-plan matrix.
// Rows: characters. Columns: skill plans. Cell: StatusDot encoding readiness.
// Sticky column headers (top) and sticky leftmost column (character) so the
// matrix scans cleanly even with many plans.

import { useMemo, useState } from 'react';
import PropTypes from 'prop-types';
import { Tooltip } from '@mui/material';
import StatusDot from '../ui/StatusDot.jsx';
import CharacterDetailModal from '../common/CharacterDetailModal.jsx';
import { calculateDaysFromToday } from '../../utils/formatter.jsx';

const formatSP = (sp) => {
    if (!sp || sp <= 0) return '—';
    if (sp >= 1_000_000_000) return `${(sp / 1_000_000_000).toFixed(2)}B`;
    if (sp >= 1_000_000) return `${(sp / 1_000_000).toFixed(1)}M`;
    if (sp >= 1_000) return `${(sp / 1_000).toFixed(0)}K`;
    return String(sp);
};

// Map a (character, plan) pair to one of the five status-dot states.
const cellStatus = (character, planName) => {
    const c = character?.Character || {};
    if (c.QualifiedPlans?.[planName]) return 'ready';
    if (c.PendingPlans?.[planName]) {
        // Currently training toward this plan if there's an active queue.
        const queue = c.SkillQueue;
        const training = Array.isArray(queue) && queue.length > 0 && character.MCT;
        return training ? 'training' : 'queued';
    }
    const missing = c.MissingSkills?.[planName];
    if (missing && Object.keys(missing).length > 0) return 'idle';
    return 'idle';
};

const cellTooltip = (character, planName) => {
    const c = character?.Character || {};
    if (c.QualifiedPlans?.[planName]) return `Ready · ${c.CharacterName} · ${planName}`;
    if (c.PendingPlans?.[planName]) {
        const eta = calculateDaysFromToday(c.PendingFinishDates?.[planName]);
        return `Training · ${c.CharacterName} · ${planName}${eta ? ` · ${eta}` : ''}`;
    }
    const missing = c.MissingSkills?.[planName];
    const n = missing ? Object.keys(missing).length : 0;
    return n > 0
        ? `${n} skill${n === 1 ? '' : 's'} missing · ${c.CharacterName} · ${planName}`
        : `Idle · ${c.CharacterName} · ${planName}`;
};

const PlanMatrix = ({ characters, skillPlans, conversions, filter }) => {
    const [detailFor, setDetailFor] = useState(null);

    const planNames = useMemo(
        () => Object.keys(skillPlans).sort((a, b) => a.localeCompare(b)),
        [skillPlans],
    );

    const rows = useMemo(() => {
        const q = filter.trim().toLowerCase();
        return characters
            .filter((ch) => {
                const name = (ch.Character?.CharacterName || '').toLowerCase();
                return q ? name.includes(q) : true;
            })
            .map((ch) => {
                const c = ch.Character || {};
                const sp = c.CharacterSkillsResponse?.total_sp || 0;
                let ready = 0;
                let training = 0;
                for (const p of planNames) {
                    const s = cellStatus(ch, p);
                    if (s === 'ready') ready += 1;
                    else if (s === 'training' || s === 'queued') training += 1;
                }
                return { character: ch, name: c.CharacterName || '', id: c.CharacterID, sp, ready, training };
            })
            .sort((a, b) => a.name.localeCompare(b.name));
    }, [characters, planNames, filter]);

    if (planNames.length === 0) {
        return (
            <div className="mt-12 text-center">
                <p className="text-body text-ink-2">No skill plans yet.</p>
                <p className="mt-1 text-meta text-ink-3">
                    Add a plan from the header to populate the matrix.
                </p>
            </div>
        );
    }

    if (rows.length === 0) {
        return (
            <div className="mt-12 text-center">
                <p className="text-body text-ink-2">No characters match this filter.</p>
            </div>
        );
    }

    // Column width is fixed so the sticky header lines up perfectly with cells.
    const colW = 36; // px
    // Leftmost frozen column width.
    const nameW = 220;
    const spW = 72;
    const summaryW = 92;

    return (
        <>
            <CharacterDetailModal
                open={Boolean(detailFor)}
                onClose={() => setDetailFor(null)}
                character={detailFor}
                skillConversions={conversions}
            />
            <div
                role="table"
                aria-label="Skill plan readiness matrix"
                className="rounded-lg border border-rule-1 bg-surface-1 overflow-auto max-h-[calc(100vh-220px)]"
            >
                {/* Column header — sticky top */}
                <div
                    role="row"
                    className="sticky top-0 z-20 flex items-end bg-surface-2 border-b border-rule-1"
                >
                    {/* Frozen header cell for the character column */}
                    <div
                        className="sticky left-0 z-30 bg-surface-2 border-r border-rule-1 flex items-end gap-3 px-4 py-2"
                        style={{ width: nameW + spW + summaryW }}
                    >
                        <span className="text-meta text-ink-3 uppercase tracking-wide flex-1">Character</span>
                        <span className="text-meta text-ink-3 uppercase tracking-wide tabular text-right" style={{ width: spW }}>SP</span>
                        <span className="text-meta text-ink-3 uppercase tracking-wide tabular text-right" style={{ width: summaryW }}>Ready / Train</span>
                    </div>
                    {/* Plan column headers — vertical orientation for density */}
                    {planNames.map((p) => {
                        const typeID = conversions?.[p];
                        const icon = typeID ? `https://images.evetech.net/types/${typeID}/icon` : null;
                        return (
                            <div
                                key={p}
                                role="columnheader"
                                title={p}
                                className="shrink-0 flex flex-col items-center justify-end gap-1.5 px-1 pb-2 pt-3 border-r border-rule-1 last:border-r-0"
                                style={{ width: colW, minHeight: 132 }}
                            >
                                {icon ? (
                                    <img
                                        src={icon}
                                        alt=""
                                        aria-hidden
                                        loading="lazy"
                                        className="h-5 w-5 rounded-sm border border-rule-1 shrink-0"
                                    />
                                ) : null}
                                <span
                                    className="text-meta text-ink-2 whitespace-nowrap"
                                    style={{
                                        writingMode: 'vertical-rl',
                                        transform: 'rotate(180deg)',
                                        maxHeight: 96,
                                        overflow: 'hidden',
                                        textOverflow: 'ellipsis',
                                    }}
                                >
                                    {p}
                                </span>
                            </div>
                        );
                    })}
                </div>

                {/* Body rows */}
                {rows.map((row, rIdx) => {
                    const isLast = rIdx === rows.length - 1;
                    return (
                        <div
                            role="row"
                            key={row.id}
                            className={`flex items-stretch ${isLast ? '' : 'border-b border-rule-1'} hover:bg-surface-2`}
                        >
                            {/* Frozen character cell */}
                            <div
                                className="sticky left-0 z-10 bg-surface-1 border-r border-rule-1 flex items-center gap-3 px-4 h-10"
                                style={{ width: nameW + spW + summaryW }}
                            >
                                <button
                                    type="button"
                                    onClick={() => setDetailFor(row.character)}
                                    className="flex items-center gap-2.5 min-w-0 flex-1 text-left focus:outline-hidden focus-visible:ring-2 focus-visible:ring-accent rounded-sm"
                                    style={{ width: nameW }}
                                >
                                    <img
                                        src={`https://images.evetech.net/characters/${row.id}/portrait?size=64`}
                                        alt=""
                                        aria-hidden
                                        loading="lazy"
                                        className="h-6 w-6 rounded-sm border border-rule-1 shrink-0"
                                    />
                                    <span className="text-body text-ink-1 truncate">{row.name}</span>
                                </button>
                                <span className="text-body font-mono tabular text-right text-ink-2" style={{ width: spW }}>
                                    {formatSP(row.sp)}
                                </span>
                                <span className="text-meta font-mono tabular text-right text-ink-3" style={{ width: summaryW }}>
                                    <span className="text-status-ready">{row.ready}</span>
                                    <span className="text-ink-4"> / </span>
                                    <span className="text-status-training">{row.training}</span>
                                </span>
                            </div>

                            {/* Status cells */}
                            {planNames.map((p) => {
                                const state = cellStatus(row.character, p);
                                return (
                                    <Tooltip key={p} title={cellTooltip(row.character, p)} arrow disableInteractive>
                                        <div
                                            role="cell"
                                            className="shrink-0 flex items-center justify-center border-r border-rule-1 last:border-r-0 h-10"
                                            style={{ width: colW }}
                                        >
                                            <StatusDot state={state} />
                                        </div>
                                    </Tooltip>
                                );
                            })}
                        </div>
                    );
                })}
            </div>
        </>
    );
};

PlanMatrix.propTypes = {
    characters: PropTypes.array.isRequired,
    skillPlans: PropTypes.object.isRequired,
    conversions: PropTypes.object.isRequired,
    filter: PropTypes.string.isRequired,
};

export default PlanMatrix;
