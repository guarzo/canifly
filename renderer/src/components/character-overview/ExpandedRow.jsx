// src/components/character-overview/ExpandedRow.jsx
//
// The detail panel revealed when a row is expanded. Houses the role
// editor (with inline "add new role") and the skill queue list.
import { useEffect, useMemo, useState } from 'react';
import PropTypes from 'prop-types';
import { ChevronRightOutlined } from '@mui/icons-material';
import { Select, MenuItem, TextField } from '@mui/material';
import useAppDataStore from '../../stores/appDataStore';
import { formatDuration } from './utils';

const ExpandedRow = ({ character, roles, skillConversions, onUpdateCharacter }) => {
    const c = character.Character || {};
    const queue = Array.isArray(c.SkillQueue) ? c.SkillQueue : [];
    const [role, setRole] = useState(character.Role || '');
    const [addingRole, setAddingRole] = useState(false);
    const [newRole, setNewRole] = useState('');

    useEffect(() => setRole(character.Role || ''), [character.Role]);

    const rolesOptions = useMemo(() => {
        const arr = [...roles];
        if (role && !arr.includes(role) && role !== 'add_new_role') arr.push(role);
        return arr;
    }, [roles, role]);

    const handleRole = (e) => {
        const v = e.target.value;
        if (v === 'add_new_role') { setAddingRole(true); return; }
        setRole(v);
        if (c.CharacterID) onUpdateCharacter(c.CharacterID, { Role: v });
    };

    const commitNewRole = async () => {
        const trimmed = newRole.trim();
        if (!trimmed) return;
        setRole(trimmed);
        if (c.CharacterID) {
            await onUpdateCharacter(c.CharacterID, { Role: trimmed });
            await useAppDataStore.getState().fetchConfig();
        }
        setAddingRole(false);
        setNewRole('');
    };

    return (
        <div className="border-b border-rule-1 bg-surface-0 px-4 py-4 grid grid-cols-1 lg:grid-cols-[minmax(0,1fr)_minmax(0,1.6fr)] gap-6">
            <div className="space-y-3">
                <div>
                    <div className="text-meta text-ink-3 mb-1">Role</div>
                    {addingRole ? (
                        <div className="flex items-center gap-2">
                            <TextField
                                size="small"
                                autoFocus
                                value={newRole}
                                onChange={(e) => setNewRole(e.target.value)}
                                onKeyDown={(e) => { if (e.key === 'Enter') commitNewRole(); }}
                                placeholder="New role"
                            />
                            <button
                                type="button"
                                onClick={commitNewRole}
                                className="px-2.5 h-8 rounded-md bg-accent text-accent-ink text-meta hover:bg-accent-strong"
                            >
                                Save
                            </button>
                            <button
                                type="button"
                                onClick={() => { setAddingRole(false); setNewRole(''); }}
                                className="px-2.5 h-8 rounded-md border border-rule-1 text-ink-2 text-meta hover:bg-surface-2"
                            >
                                Cancel
                            </button>
                        </div>
                    ) : (
                        <Select
                            value={role}
                            onChange={handleRole}
                            displayEmpty
                            size="small"
                            sx={{ minWidth: 160 }}
                        >
                            <MenuItem value="" disabled>Select role</MenuItem>
                            {rolesOptions.map((r) => (
                                <MenuItem key={r} value={r}>{r}</MenuItem>
                            ))}
                            <MenuItem value="add_new_role">+ Add new role…</MenuItem>
                        </Select>
                    )}
                </div>
                <dl className="grid grid-cols-[max-content_1fr] gap-x-4 gap-y-1.5 text-meta">
                    <dt className="text-ink-3">Location</dt>
                    <dd className="text-ink-1">{c.LocationName || '—'}</dd>
                    {character.CorporationName ? (
                        <>
                            <dt className="text-ink-3">Corporation</dt>
                            <dd className="text-ink-1">{character.CorporationName}</dd>
                        </>
                    ) : null}
                    {character.AllianceName ? (
                        <>
                            <dt className="text-ink-3">Alliance</dt>
                            <dd className="text-ink-1">{character.AllianceName}</dd>
                        </>
                    ) : null}
                    <dt className="text-ink-3">Total SP</dt>
                    <dd className="text-ink-1 font-mono tabular">
                        {(c.CharacterSkillsResponse?.total_sp || 0).toLocaleString()}
                    </dd>
                </dl>
            </div>
            <div>
                <div className="flex items-baseline justify-between mb-2">
                    <div className="text-meta text-ink-3">Skill queue</div>
                    <div className="text-micro text-ink-3">{queue.length} item{queue.length === 1 ? '' : 's'}</div>
                </div>
                {queue.length === 0 ? (
                    <p className="text-meta text-ink-3">No skill queue.</p>
                ) : (
                    <ol className="rounded-md border border-rule-1 overflow-hidden">
                        {queue.slice(0, 10).map((item, i) => {
                            const skillName = skillConversions?.[String(item.skill_id)] || `Skill #${item.skill_id}`;
                            const ms = item.finish_date ? new Date(item.finish_date).getTime() - Date.now() : null;
                            return (
                                <li
                                    key={`${item.skill_id}-${i}`}
                                    className="grid grid-cols-[1fr_44px_92px] gap-3 items-center px-3 h-9 border-b border-rule-1 last:border-b-0 bg-surface-1"
                                >
                                    <span className="text-meta text-ink-1 truncate">{skillName}</span>
                                    <span className="text-meta text-ink-2 font-mono tabular text-right">L{item.finished_level}</span>
                                    <span className="text-meta text-ink-2 font-mono tabular text-right">
                                        {ms != null && ms > 0 ? formatDuration(ms) : <ChevronRightOutlined sx={{ fontSize: 14, color: 'var(--ink-3)' }} />}
                                    </span>
                                </li>
                            );
                        })}
                        {queue.length > 10 ? (
                            <li className="px-3 h-8 flex items-center text-micro text-ink-3 bg-surface-1">
                                + {queue.length - 10} more
                            </li>
                        ) : null}
                    </ol>
                )}
            </div>
        </div>
    );
};

ExpandedRow.propTypes = {
    character: PropTypes.object.isRequired,
    roles: PropTypes.array.isRequired,
    skillConversions: PropTypes.object.isRequired,
    onUpdateCharacter: PropTypes.func.isRequired,
};

export default ExpandedRow;
