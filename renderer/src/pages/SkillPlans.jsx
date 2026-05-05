// renderer/src/pages/SkillPlans.jsx
import { useEffect, useMemo, useRef, useState } from 'react';
import PropTypes from 'prop-types';
import {
    GridOnOutlined,
    ListAltOutlined,
    ViewTimelineOutlined,
} from '@mui/icons-material';

import PageShell from '../components/ui/PageShell.jsx';
import FilterBar from '../components/ui/FilterBar.jsx';
import StatusDot from '../components/ui/StatusDot.jsx';
import SegmentedControl from '../components/ui/SegmentedControl.jsx';
import Kbd from '../components/ui/Kbd.jsx';
import PlanMatrix from '../components/skillplan/PlanMatrix.jsx';
import PlanList from '../components/skillplan/PlanList.jsx';
import PlanTimeline from '../components/skillplan/PlanTimeline.jsx';
import { skillPlanInstructions } from '../utils/instructions.jsx';
import { useAsyncOperation } from '../hooks/useAsyncOperation';
import { copySkillPlan, deleteSkillPlan } from '../api/skillPlansApi';

const LS = {
    view: 'sp.view',
    filter: 'sp.filter',
};

const readLS = (key, fallback) => {
    try {
        const raw = localStorage.getItem(key);
        return raw === null ? fallback : JSON.parse(raw);
    } catch {
        return fallback;
    }
};
const writeLS = (key, val) => {
    try { localStorage.setItem(key, JSON.stringify(val)); } catch { /* swallow */ }
};

const VIEW_OPTIONS = [
    { value: 'timeline', label: 'Timeline', icon: <ViewTimelineOutlined fontSize="small" /> },
    { value: 'matrix', label: 'Matrix', icon: <GridOnOutlined fontSize="small" /> },
    { value: 'plans', label: 'By plan', icon: <ListAltOutlined fontSize="small" /> },
];

const SkillPlans = ({ characters = [], skillPlans = {}, conversions = {} }) => {
    const { execute } = useAsyncOperation();

    const [view, setView] = useState(() => readLS(LS.view, 'timeline'));
    const [filter, setFilter] = useState(() => readLS(LS.filter, ''));
    const filterRef = useRef(null);

    useEffect(() => writeLS(LS.view, view), [view]);
    useEffect(() => writeLS(LS.filter, filter), [filter]);

    useEffect(() => {
        const onKey = (e) => {
            const tag = (e.target?.tagName || '').toLowerCase();
            const isTyping = tag === 'input' || tag === 'textarea' || e.target?.isContentEditable;
            if (e.key === '/' && !isTyping) {
                e.preventDefault();
                filterRef.current?.focus();
            } else if (e.key === 'Escape' && document.activeElement === filterRef.current) {
                setFilter('');
                filterRef.current?.blur();
            }
        };
        window.addEventListener('keydown', onKey);
        return () => window.removeEventListener('keydown', onKey);
    }, []);

    const handleCopy = (planName) => {
        const existing = new Set(Object.keys(skillPlans));
        let target = `${planName} (copy)`;
        let n = 2;
        while (existing.has(target)) {
            target = `${planName} (copy ${n})`;
            n += 1;
        }
        return execute(
            () => copySkillPlan(planName, target),
            { successMessage: `Copied to "${target}"` },
        );
    };
    const handleDelete = (planName) => execute(
        () => deleteSkillPlan(planName),
        { successMessage: 'Skill plan deleted' },
    );

    const summary = useMemo(() => {
        const planCount = Object.keys(skillPlans).length;
        const charCount = characters.length;
        let ready = 0;
        let training = 0;
        for (const ch of characters) {
            const c = ch.Character || {};
            for (const p of Object.keys(skillPlans)) {
                if (c.QualifiedPlans?.[p]) ready += 1;
                else if (c.PendingPlans?.[p]) training += 1;
            }
        }
        const total = planCount * charCount;
        return { planCount, charCount, ready, training, total };
    }, [characters, skillPlans]);

    const meta =
        `${summary.charCount} character${summary.charCount === 1 ? '' : 's'} ` +
        `· ${summary.planCount} plan${summary.planCount === 1 ? '' : 's'} ` +
        `· ${summary.ready} ready / ${summary.training} training of ${summary.total}`;

    const headerActions = (
        <SegmentedControl value={view} onChange={setView} ariaLabel="View mode" options={VIEW_OPTIONS} />
    );

    const placeholder =
        view === 'matrix' ? 'Filter characters' :
        view === 'plans' ? 'Filter plans' :
        'Filter plans';

    return (
        <PageShell title="Skill Plans" meta={meta} actions={headerActions} tip={skillPlanInstructions}>
            <FilterBar
                inputRef={filterRef}
                value={filter}
                onChange={setFilter}
                placeholder={placeholder}
                count={
                    <span className="flex items-center gap-x-4 gap-y-1 flex-wrap text-meta text-ink-3">
                        <span className="inline-flex items-center gap-1.5"><StatusDot state="ready" /> ready</span>
                        <span className="inline-flex items-center gap-1.5"><StatusDot state="training" /> training</span>
                        <span className="inline-flex items-center gap-1.5"><StatusDot state="queued" /> queued</span>
                        <span className="inline-flex items-center gap-1.5"><StatusDot state="idle" /> idle</span>
                    </span>
                }
            />

            {characters.length === 0 ? (
                <div className="mt-16 text-center">
                    <p className="text-body text-ink-2">No characters loaded.</p>
                    <p className="mt-1 text-meta text-ink-3">Add an EVE account to populate the matrix.</p>
                </div>
            ) : view === 'timeline' ? (
                <PlanTimeline characters={characters} skillPlans={skillPlans} conversions={conversions} filter={filter} />
            ) : view === 'matrix' ? (
                <PlanMatrix characters={characters} skillPlans={skillPlans} conversions={conversions} filter={filter} />
            ) : (
                <PlanList
                    skillPlans={skillPlans}
                    characters={characters}
                    conversions={conversions}
                    filter={filter}
                    onCopy={handleCopy}
                    onDelete={handleDelete}
                />
            )}

            <div className="mt-4 flex flex-wrap items-center gap-x-4 gap-y-1 text-micro text-ink-3 font-mono">
                <span><Kbd>/</Kbd> filter</span>
                <span><Kbd>Esc</Kbd> clear filter</span>
            </div>
        </PageShell>
    );
};

SkillPlans.propTypes = {
    characters: PropTypes.array,
    skillPlans: PropTypes.object,
    conversions: PropTypes.object,
};

export default SkillPlans;
