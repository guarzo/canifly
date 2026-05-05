// renderer/src/components/skillplan/PlanMatrix.jsx
import { useMemo, useState } from 'react';
import PropTypes from 'prop-types';
import { Tooltip } from '@mui/material';
import StatusDot from '../ui/StatusDot.jsx';
import EveTypeIcon from '../ui/EveTypeIcon.jsx';
import CharacterDetailModal from '../common/CharacterDetailModal.jsx';
import MatrixShell from './MatrixShell.jsx';
import { computePlanReadiness } from '../../hooks/usePlanReadiness.js';

const formatSP = (sp) => {
    if (!sp || sp <= 0) return '—';
    if (sp >= 1_000_000_000) return `${(sp / 1_000_000_000).toFixed(2)}B`;
    if (sp >= 1_000_000) return `${(sp / 1_000_000).toFixed(1)}M`;
    if (sp >= 1_000) return `${(sp / 1_000).toFixed(0)}K`;
    return String(sp);
};

const cellTooltip = (character, planName, readiness) => {
    const c = character?.Character || {};
    const name = c.CharacterName;
    if (readiness.state === 'ready') return `Ready · ${name} · ${planName}`;
    if (readiness.state === 'training' || readiness.state === 'queued') {
        return `${readiness.state === 'training' ? 'Training' : 'Queued'} · ${name} · ${planName}${readiness.eta ? ` · ${readiness.eta}` : ''}`;
    }
    if (readiness.missingCount > 0) {
        return `${readiness.missingCount} skill${readiness.missingCount === 1 ? '' : 's'} missing · ${name} · ${planName}`;
    }
    return `Idle · ${name} · ${planName}`;
};

const FROZEN_W = 220 + 72 + 92;
const NAME_W = 220;
const SP_W = 72;
const SUMMARY_W = 92;

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
                const readinessByPlan = {};
                for (const p of planNames) {
                    const r = computePlanReadiness(ch, p);
                    readinessByPlan[p] = r;
                    if (r.state === 'ready') ready += 1;
                    else if (r.state === 'training' || r.state === 'queued') training += 1;
                }
                return { character: ch, name: c.CharacterName || '', id: c.CharacterID, sp, ready, training, readinessByPlan };
            })
            .sort((a, b) => a.name.localeCompare(b.name));
    }, [characters, planNames, filter]);

    if (planNames.length === 0) {
        return (
            <div className="mt-12 text-center">
                <p className="text-body text-ink-2">No skill plans yet.</p>
                <p className="mt-1 text-meta text-ink-3">Add a plan from the header to populate the matrix.</p>
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

    const frozenHeader = (
        <>
            <span className="text-meta text-ink-3 uppercase tracking-wide flex-1">Character</span>
            <span className="text-meta text-ink-3 uppercase tracking-wide tabular text-right" style={{ width: SP_W }}>SP</span>
            <span className="text-meta text-ink-3 uppercase tracking-wide tabular text-right" style={{ width: SUMMARY_W }}>Ready / Train</span>
        </>
    );

    const columnHeaders = planNames.map((p) => (
        <div key={p} title={p} className="flex flex-col items-center justify-end gap-1.5">
            <EveTypeIcon name={p} conversions={conversions} size={5} />
            <span
                className="text-meta text-ink-2 whitespace-nowrap"
                style={{ writingMode: 'vertical-rl', transform: 'rotate(180deg)', maxHeight: 96, overflow: 'hidden', textOverflow: 'ellipsis' }}
            >
                {p}
            </span>
        </div>
    ));

    const matrixRows = rows.map((row) => ({
        key: row.id,
        frozenCell: (
            <>
                <button
                    type="button"
                    onClick={() => setDetailFor(row.character)}
                    className="flex items-center gap-2.5 min-w-0 flex-1 text-left focus:outline-hidden focus-visible:ring-2 focus-visible:ring-accent rounded-sm"
                    style={{ width: NAME_W }}
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
                <span className="text-body font-mono tabular text-right text-ink-2" style={{ width: SP_W }}>{formatSP(row.sp)}</span>
                <span className="text-meta font-mono tabular text-right text-ink-3" style={{ width: SUMMARY_W }}>
                    <span className="text-status-ready">{row.ready}</span>
                    <span className="text-ink-4"> / </span>
                    <span className="text-status-training">{row.training}</span>
                </span>
            </>
        ),
        cells: planNames.map((p) => {
            const r = row.readinessByPlan[p];
            return (
                <Tooltip key={p} title={cellTooltip(row.character, p, r)} arrow disableInteractive>
                    <span><StatusDot state={r.state} /></span>
                </Tooltip>
            );
        }),
    }));

    return (
        <>
            <CharacterDetailModal
                open={Boolean(detailFor)}
                onClose={() => setDetailFor(null)}
                character={detailFor}
                skillConversions={conversions}
            />
            <MatrixShell
                ariaLabel="Skill plan readiness matrix"
                frozenHeader={frozenHeader}
                columnHeaders={columnHeaders}
                rows={matrixRows}
                frozenWidth={FROZEN_W}
            />
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
