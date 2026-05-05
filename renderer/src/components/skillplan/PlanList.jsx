// PlanList — manage skill plans (copy / delete) with character roll-up.
// Used as the "By plan" view; complements the matrix.

import { useMemo, useState } from 'react';
import PropTypes from 'prop-types';
import {
    ContentCopyOutlined,
    DeleteOutlined as DeleteOutline,
    ExpandMoreOutlined,
    ChevronRightOutlined,
} from '@mui/icons-material';
import { IconButton, Tooltip } from '@mui/material';
import StatusDot from '../ui/StatusDot.jsx';
import { calculateDaysFromToday } from '../../utils/formatter.jsx';
import MissingSkillsPopover from './MissingSkillsPopover.jsx';

const PlanList = ({ skillPlans, characters, conversions, filter, onCopy, onDelete }) => {
    const [expanded, setExpanded] = useState(() => new Set());
    const toggle = (name) => setExpanded((prev) => {
        const next = new Set(prev);
        if (next.has(name)) next.delete(name); else next.add(name);
        return next;
    });

    const plans = useMemo(() => {
        const q = filter.trim().toLowerCase();
        return Object.values(skillPlans)
            .filter((sp) => (q ? (sp.Name || '').toLowerCase().includes(q) : true))
            .sort((a, b) => (a.Name || '').localeCompare(b.Name || ''))
            .map((sp) => {
                const qualified = sp.QualifiedCharacters || [];
                const pending = sp.PendingCharacters || [];
                const missing = sp.MissingCharacters || [];
                const children = [
                    ...qualified.map((n) => ({ name: n, state: 'ready', label: 'Ready' })),
                    ...pending.map((n) => {
                        const ch = characters.find((c) => c.Character?.CharacterName === n);
                        const eta = calculateDaysFromToday(ch?.Character?.PendingFinishDates?.[sp.Name]);
                        return { name: n, state: 'queued', label: eta ? `Training · ${eta}` : 'Training' };
                    }),
                    ...missing.map((n) => ({ name: n, state: 'idle', label: 'Missing skills' })),
                ];
                const aggregate = {};
                for (const charName of missing) {
                    const ch = characters.find((c) => c.Character?.CharacterName === charName);
                    const charMissing = ch?.Character?.MissingSkills?.[sp.Name] || {};
                    for (const [skill, lvl] of Object.entries(charMissing)) {
                        aggregate[skill] = Math.max(aggregate[skill] ?? 0, lvl);
                    }
                }
                return {
                    name: sp.Name,
                    qualified: qualified.length,
                    pending: pending.length,
                    missing: missing.length,
                    aggregateMissing: aggregate,
                    children,
                };
            });
    }, [skillPlans, characters, filter]);

    if (plans.length === 0) {
        return (
            <div className="mt-12 text-center">
                <p className="text-body text-ink-2">No skill plans match.</p>
            </div>
        );
    }

    const colHeader = (
        <div
            role="row"
            className="hidden sm:grid grid-cols-[28px_minmax(220px,1.6fr)_72px_72px_72px_88px] gap-3 px-4 py-2 border-b border-rule-1 bg-surface-2 text-meta text-ink-3 uppercase tracking-wide"
        >
            <div role="columnheader" aria-label="Expand" />
            <div role="columnheader">Skill plan</div>
            <div role="columnheader" className="text-right">Ready</div>
            <div role="columnheader" className="text-right">Train</div>
            <div role="columnheader" className="text-right">Missing</div>
            <div role="columnheader" aria-label="Actions" className="text-right">Actions</div>
        </div>
    );

    return (
        <div role="table" aria-label="Skill plans" className="rounded-lg border border-rule-1 bg-surface-1 overflow-hidden">
            {colHeader}
            {plans.map((p, idx) => {
                const isOpen = expanded.has(p.name);
                const isLast = idx === plans.length - 1;
                const typeID = conversions?.[p.name];
                const icon = typeID ? `https://images.evetech.net/types/${typeID}/icon` : null;
                return (
                    <div key={p.name} role="rowgroup">
                        <div
                            role="row"
                            tabIndex={0}
                            aria-expanded={isOpen}
                            onClick={() => toggle(p.name)}
                            onKeyDown={(e) => {
                                if (e.key === 'Enter' || e.key === ' ') {
                                    e.preventDefault();
                                    toggle(p.name);
                                }
                            }}
                            className={[
                                'grid grid-cols-[28px_minmax(220px,1.6fr)_72px_72px_72px_88px]',
                                'gap-3 px-4 items-center h-10 cursor-pointer select-none',
                                !isLast || isOpen ? 'border-b border-rule-1' : '',
                                'hover:bg-surface-2 focus:outline-hidden focus-visible:bg-surface-2',
                            ].join(' ')}
                        >
                            <div className="flex items-center justify-center">
                                {isOpen
                                    ? <ExpandMoreOutlined fontSize="small" sx={{ color: 'var(--accent)' }} />
                                    : <ChevronRightOutlined fontSize="small" sx={{ color: 'var(--ink-3)' }} />}
                            </div>
                            <div className="flex items-center gap-2.5 min-w-0">
                                {icon ? (
                                    <img src={icon} alt="" aria-hidden loading="lazy" className="h-6 w-6 rounded-sm border border-rule-1 shrink-0" />
                                ) : (
                                    <span className="h-6 w-6 rounded-sm border border-rule-1 bg-surface-2 shrink-0" aria-hidden />
                                )}
                                <span className="text-body text-ink-1 truncate">{p.name}</span>
                            </div>
                            <span className="text-body font-mono tabular text-right text-status-ready">{p.qualified}</span>
                            <span className="text-body font-mono tabular text-right text-status-training">{p.pending}</span>
                            {p.missing > 0 ? (
                                <div className="text-right" onClick={(e) => e.stopPropagation()}>
                                    <MissingSkillsPopover
                                        planName={p.name}
                                        missingSkills={p.aggregateMissing}
                                        trigger={
                                            <button
                                                type="button"
                                                className="text-body font-mono tabular text-ink-3 hover:text-ink-1 underline-offset-2 hover:underline"
                                                aria-label={`Show ${p.missing} missing skill${p.missing === 1 ? '' : 's'} across characters for ${p.name}`}
                                            >
                                                {p.missing}
                                            </button>
                                        }
                                    />
                                </div>
                            ) : (
                                <span className="text-body font-mono tabular text-right text-ink-3">{p.missing}</span>
                            )}
                            <div className="flex items-center justify-end gap-1" onClick={(e) => e.stopPropagation()}>
                                <Tooltip title="Copy skill plan">
                                    <IconButton
                                        size="small"
                                        aria-label={`Copy ${p.name}`}
                                        onClick={() => onCopy(p.name)}
                                        className="h-7! w-7! rounded-md! text-ink-3! hover:bg-surface-3! hover:text-ink-1!"
                                    >
                                        <ContentCopyOutlined sx={{ fontSize: 16 }} />
                                    </IconButton>
                                </Tooltip>
                                <Tooltip title="Delete skill plan">
                                    <IconButton
                                        size="small"
                                        aria-label={`Delete ${p.name}`}
                                        onClick={() => onDelete(p.name)}
                                        className="h-7! w-7! rounded-md! text-ink-3! hover:bg-surface-3!"
                                        sx={{ '&:hover': { color: 'var(--status-error)' } }}
                                    >
                                        <DeleteOutline sx={{ fontSize: 16 }} />
                                    </IconButton>
                                </Tooltip>
                            </div>
                        </div>

                        {isOpen ? (
                            <div className={`bg-surface-0 ${isLast ? '' : 'border-b border-rule-1'}`}>
                                {p.children.length === 0 ? (
                                    <div className="px-12 py-3 text-meta text-ink-3">No characters tracked for this plan.</div>
                                ) : (
                                    p.children.map((child) => (
                                        <div
                                            key={`${p.name}-${child.name}`}
                                            className="grid grid-cols-[28px_minmax(220px,1.6fr)_1fr] gap-3 items-center px-4 h-9 border-t border-rule-1"
                                        >
                                            <div className="flex items-center justify-center">
                                                <StatusDot state={child.state} />
                                            </div>
                                            <span className="text-meta text-ink-1 truncate pl-6">{child.name}</span>
                                            <span className="text-meta text-ink-3 truncate">{child.label}</span>
                                        </div>
                                    ))
                                )}
                            </div>
                        ) : null}
                    </div>
                );
            })}
        </div>
    );
};

PlanList.propTypes = {
    skillPlans: PropTypes.object.isRequired,
    characters: PropTypes.array.isRequired,
    conversions: PropTypes.object.isRequired,
    filter: PropTypes.string.isRequired,
    onCopy: PropTypes.func.isRequired,
    onDelete: PropTypes.func.isRequired,
};

export default PlanList;
