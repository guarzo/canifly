// src/pages/SkillPlans.jsx
//
// Character × Skill-plan readiness, redesigned per DESIGN.md.
// The densest analytical view in the app: matrix on by default, plan-list as
// the management view. Persistent UI state under the "sp." key prefix.

import { useEffect, useMemo, useRef, useState } from 'react';
import PropTypes from 'prop-types';
import {
    GridOnOutlined,
    ListAltOutlined,
    Search as SearchIcon,
} from '@mui/icons-material';

import Subheader from '../components/ui/Subheader.jsx';
import StatusDot from '../components/ui/StatusDot.jsx';
import SegmentedControl from '../components/ui/SegmentedControl.jsx';
import Kbd from '../components/ui/Kbd.jsx';
import PlanMatrix from '../components/skillplan/PlanMatrix.jsx';
import PlanList from '../components/skillplan/PlanList.jsx';
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
    { value: 'matrix', label: 'Matrix', icon: <GridOnOutlined fontSize="small" /> },
    { value: 'plans', label: 'By plan', icon: <ListAltOutlined fontSize="small" /> },
];

const SkillPlans = ({ characters = [], skillPlans = {}, conversions = {} }) => {
    const { execute } = useAsyncOperation();

    const [view, setView] = useState(() => readLS(LS.view, 'matrix'));
    const [filter, setFilter] = useState(() => readLS(LS.filter, ''));
    const filterRef = useRef(null);

    useEffect(() => writeLS(LS.view, view), [view]);
    useEffect(() => writeLS(LS.filter, filter), [filter]);

    // Slash focuses the filter, esc clears it.
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

    const handleCopy = (planName, newPlanName) => execute(
        () => copySkillPlan(planName, newPlanName),
        { successMessage: 'Skill plan copied' },
    );
    const handleDelete = (planName) => execute(
        () => deleteSkillPlan(planName),
        { successMessage: 'Skill plan deleted' },
    );

    // Header summary across the whole matrix.
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
        <SegmentedControl
            value={view}
            onChange={setView}
            ariaLabel="View mode"
            options={VIEW_OPTIONS}
        />
    );

    const placeholder = view === 'matrix'
        ? 'Filter characters'
        : 'Filter plans';

    return (
        <div className="px-6 pb-12 pt-8 max-w-[1280px] mx-auto">
            <Subheader
                title="Skill Plans"
                meta={meta}
                actions={headerActions}
                tip={skillPlanInstructions}
            />

            {/* Filter + legend */}
            <div className="mb-4 flex items-center gap-3 flex-wrap">
                <div className="flex-1 min-w-[240px] max-w-md flex items-center gap-2 rounded-md border border-rule-1 bg-surface-1 px-2.5 py-1.5 focus-within:border-accent">
                    <SearchIcon fontSize="small" sx={{ color: 'var(--ink-3)' }} />
                    <input
                        ref={filterRef}
                        value={filter}
                        onChange={(e) => setFilter(e.target.value)}
                        placeholder={placeholder}
                        aria-label={placeholder}
                        className="bg-transparent flex-1 outline-hidden text-body text-ink-1 placeholder:text-ink-3"
                    />
                    {filter ? (
                        <button
                            type="button"
                            onClick={() => setFilter('')}
                            className="text-ink-3 hover:text-ink-1 text-meta"
                            aria-label="Clear filter"
                        >
                            clear
                        </button>
                    ) : (
                        <Kbd>/</Kbd>
                    )}
                </div>

                <div className="flex items-center gap-x-4 gap-y-1 flex-wrap text-meta text-ink-3">
                    <span className="inline-flex items-center gap-1.5"><StatusDot state="ready" /> ready</span>
                    <span className="inline-flex items-center gap-1.5"><StatusDot state="training" /> training</span>
                    <span className="inline-flex items-center gap-1.5"><StatusDot state="queued" /> queued</span>
                    <span className="inline-flex items-center gap-1.5"><StatusDot state="idle" /> idle</span>
                </div>
            </div>

            {characters.length === 0 ? (
                <div className="mt-16 text-center">
                    <p className="text-body text-ink-2">No characters loaded.</p>
                    <p className="mt-1 text-meta text-ink-3">Add an EVE account to populate the matrix.</p>
                </div>
            ) : view === 'matrix' ? (
                <PlanMatrix
                    characters={characters}
                    skillPlans={skillPlans}
                    conversions={conversions}
                    filter={filter}
                />
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
        </div>
    );
};

SkillPlans.propTypes = {
    characters: PropTypes.array,
    skillPlans: PropTypes.object,
    conversions: PropTypes.object,
};

export default SkillPlans;
