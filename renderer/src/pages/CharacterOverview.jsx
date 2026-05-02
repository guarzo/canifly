// src/pages/CharacterOverview.jsx
//
// Table-first character overview, redesigned per DESIGN.md and PRODUCT.md.
// Dense, keyboard-navigable, persistent. The single page where a multi-account
// EVE pilot answers "what is the state of every character?" in seconds.

import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import PropTypes from 'prop-types';
import {
    AccountBalanceOutlined,
    AccountCircleOutlined,
    PlaceOutlined,
    Refresh as RefreshIcon,
    VisibilityOutlined,
    VisibilityOffOutlined,
    DeleteOutline,
    OpenInNewOutlined,
    Search as SearchIcon,
    ExpandMoreOutlined,
    ChevronRightOutlined,
    MoreVert as MoreVertIcon,
} from '@mui/icons-material';
import {
    IconButton,
    Tooltip,
    Menu,
    MenuItem,
    Select,
    TextField,
} from '@mui/material';

import Subheader from '../components/ui/Subheader.jsx';
import StatusDot from '../components/ui/StatusDot.jsx';
import SegmentedControl from '../components/ui/SegmentedControl.jsx';
import Kbd from '../components/ui/Kbd.jsx';
import { overviewInstructions } from '../utils/instructions.jsx';
import { useAppData } from '../hooks/useAppData';
import { useAsyncOperation } from '../hooks/useAsyncOperation';
import { updateCharacter, deleteCharacter, refreshCharacter } from '../api/accountsApi';
import useAppDataStore from '../stores/appDataStore';
import { logger } from '../utils/logger';

// ── Storage keys ──────────────────────────────────────────────────────────
const LS = {
    view: 'co.view',
    sort: 'co.sort',
    showHidden: 'co.showHidden',
    filter: 'co.filter',
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

// ── Formatters ────────────────────────────────────────────────────────────
const formatSP = (sp) => {
    if (!sp || sp <= 0) return '—';
    if (sp >= 1_000_000_000) return `${(sp / 1_000_000_000).toFixed(2)}B`;
    if (sp >= 1_000_000) return `${(sp / 1_000_000).toFixed(1)}M`;
    if (sp >= 1_000) return `${(sp / 1_000).toFixed(0)}K`;
    return String(sp);
};

const formatDuration = (ms) => {
    if (ms == null || ms <= 0) return '—';
    const totalMin = Math.floor(ms / 60_000);
    const days = Math.floor(totalMin / (60 * 24));
    const hours = Math.floor((totalMin % (60 * 24)) / 60);
    const mins = totalMin % 60;
    if (days > 0) return `${days}d ${String(hours).padStart(2, '0')}h`;
    if (hours > 0) return `${hours}h ${String(mins).padStart(2, '0')}m`;
    return `${mins}m`;
};

// Derive queue ETA from the skill queue's last finish_date.
const deriveQueueEta = (character) => {
    const queue = character?.Character?.SkillQueue;
    if (!Array.isArray(queue) || queue.length === 0) return null;
    const last = queue[queue.length - 1];
    if (!last?.finish_date) return null;
    const ms = new Date(last.finish_date).getTime() - Date.now();
    return ms > 0 ? ms : null;
};

const deriveStatus = (character) => {
    const queue = character?.Character?.SkillQueue;
    const hasQueue = Array.isArray(queue) && queue.length > 0;
    if (!hasQueue) return 'idle';
    if (character.MCT) return 'training';
    return 'queued';
};

// ── The page ──────────────────────────────────────────────────────────────
const VIEW_OPTIONS = [
    { value: 'account', label: 'Account', icon: <AccountBalanceOutlined fontSize="small" /> },
    { value: 'role', label: 'Role', icon: <AccountCircleOutlined fontSize="small" /> },
    { value: 'location', label: 'Location', icon: <PlaceOutlined fontSize="small" /> },
];

const CharacterOverview = ({ roles = [], skillConversions = {} }) => {
    const { accounts = [], updateAccount, deleteAccount, fetchAccounts } = useAppData();
    const { execute } = useAsyncOperation();

    // Persisted UI state.
    const [view, setView] = useState(() => readLS(LS.view, 'account'));
    const [sortOrder, setSortOrder] = useState(() => readLS(LS.sort, 'asc'));
    const [showHidden, setShowHidden] = useState(() => readLS(LS.showHidden, false));
    const [filter, setFilter] = useState(() => readLS(LS.filter, ''));
    const [isRefreshing, setIsRefreshing] = useState(false);
    const [expanded, setExpanded] = useState(() => new Set());
    const [focusedRowId, setFocusedRowId] = useState(null);

    useEffect(() => writeLS(LS.view, view), [view]);
    useEffect(() => writeLS(LS.sort, sortOrder), [sortOrder]);
    useEffect(() => writeLS(LS.showHidden, showHidden), [showHidden]);
    useEffect(() => writeLS(LS.filter, filter), [filter]);

    const filterRef = useRef(null);
    const refreshRef = useRef(null);

    // Visible accounts (after the show-hidden toggle).
    const visibleAccounts = useMemo(() => {
        if (!accounts) return [];
        return showHidden ? accounts : accounts.filter((a) => a.Visible !== false);
    }, [accounts, showHidden]);

    // Flatten all characters (with their account context).
    const allCharacters = useMemo(() => {
        const out = [];
        for (const acc of visibleAccounts) {
            const accountName = acc.Name || 'Unknown Account';
            for (const ch of acc.Characters || []) {
                out.push({
                    ...ch,
                    accountId: acc.ID,
                    accountName,
                    accountVisible: acc.Visible !== false,
                });
            }
        }
        return out;
    }, [visibleAccounts]);

    // Filter pass.
    const filteredCharacters = useMemo(() => {
        const q = filter.trim().toLowerCase();
        if (!q) return allCharacters;
        return allCharacters.filter((ch) => {
            const c = ch.Character || {};
            const name = (c.CharacterName || '').toLowerCase();
            const loc = (c.LocationName || '').toLowerCase();
            const role = (ch.Role || '').toLowerCase();
            const ship = (c.ShipTypeName || c.CurrentShip || '').toLowerCase();
            return name.includes(q) || loc.includes(q) || role.includes(q) || ship.includes(q);
        });
    }, [allCharacters, filter]);

    // Group characters into ordered groups by current view.
    const groups = useMemo(() => {
        const map = new Map();
        if (view === 'account') {
            for (const ch of filteredCharacters) {
                const key = ch.accountName;
                if (!map.has(key)) map.set(key, []);
                map.get(key).push(ch);
            }
        } else if (view === 'role') {
            for (const r of roles) map.set(r, []);
            map.set('Unassigned', []);
            for (const ch of filteredCharacters) {
                const key = ch.Role || 'Unassigned';
                if (!map.has(key)) map.set(key, []);
                map.get(key).push(ch);
            }
        } else {
            for (const ch of filteredCharacters) {
                const key = ch.Character?.LocationName || 'Unknown Location';
                if (!map.has(key)) map.set(key, []);
                map.get(key).push(ch);
            }
        }
        // Sort group keys; drop empty role groups.
        let keys = [...map.keys()];
        if (view === 'role') keys = keys.filter((k) => (map.get(k) || []).length > 0);
        keys.sort((a, b) => sortOrder === 'asc' ? a.localeCompare(b) : b.localeCompare(a));
        return keys.map((k) => ({
            key: k,
            characters: [...(map.get(k) || [])].sort((a, b) =>
                (a.Character?.CharacterName || '').localeCompare(b.Character?.CharacterName || ''),
            ),
        }));
    }, [filteredCharacters, view, sortOrder, roles]);

    // Header summary.
    const summary = useMemo(() => {
        const charCount = allCharacters.length;
        const accountCount = visibleAccounts.length;
        const totalSp = allCharacters.reduce(
            (sum, ch) => sum + (ch.Character?.CharacterSkillsResponse?.total_sp || 0), 0,
        );
        return { charCount, accountCount, totalSp };
    }, [allCharacters, visibleAccounts]);

    // ── Handlers ──────────────────────────────────────────────────────────
    const toggleExpanded = useCallback((id) => {
        setExpanded((prev) => {
            const next = new Set(prev);
            if (next.has(id)) next.delete(id); else next.add(id);
            return next;
        });
    }, []);

    const handleRefreshAll = useCallback(async () => {
        setIsRefreshing(true);
        try {
            for (const acc of accounts || []) {
                for (const ch of acc.Characters || []) {
                    try {
                        await refreshCharacter(ch.Character.CharacterID);
                    } catch (err) {
                        logger.error(`Refresh failed for ${ch.Character.CharacterName}:`, err);
                    }
                }
            }
            await fetchAccounts();
        } finally {
            setIsRefreshing(false);
        }
    }, [accounts, fetchAccounts]);

    const handleUpdateCharacter = useCallback(async (characterId, updates) => {
        await execute(
            () => updateCharacter(characterId, updates),
            { successMessage: 'Character updated' },
        );
    }, [execute]);

    const handleRemoveCharacter = useCallback(async (characterId) => {
        await execute(
            () => deleteCharacter(characterId),
            { successMessage: 'Character removed' },
        );
    }, [execute]);

    // ── Keyboard shortcuts ────────────────────────────────────────────────
    const flatRowIds = useMemo(
        () => groups.flatMap((g) => g.characters.map((c) => c.Character.CharacterID)),
        [groups],
    );

    useEffect(() => {
        const onKey = (e) => {
            const tag = (e.target?.tagName || '').toLowerCase();
            const isTyping = tag === 'input' || tag === 'textarea' || e.target?.isContentEditable;
            // Slash always focuses filter, even in inputs.
            if (e.key === '/' && !isTyping) {
                e.preventDefault();
                filterRef.current?.focus();
                return;
            }
            if (isTyping) return;
            if (e.key === 'j' || e.key === 'ArrowDown') {
                e.preventDefault();
                setFocusedRowId((cur) => {
                    const idx = flatRowIds.indexOf(cur);
                    return flatRowIds[Math.min(flatRowIds.length - 1, idx < 0 ? 0 : idx + 1)] ?? cur;
                });
            } else if (e.key === 'k' || e.key === 'ArrowUp') {
                e.preventDefault();
                setFocusedRowId((cur) => {
                    const idx = flatRowIds.indexOf(cur);
                    return flatRowIds[Math.max(0, idx <= 0 ? 0 : idx - 1)] ?? cur;
                });
            } else if (e.key === 'Enter' && focusedRowId != null) {
                e.preventDefault();
                toggleExpanded(focusedRowId);
            } else if (e.key === 'Escape') {
                setExpanded(new Set());
            } else if (e.key === 'r') {
                e.preventDefault();
                refreshRef.current?.click();
            } else if (e.key === 'g') {
                // chord: g a / g r / g l
                const handler = (e2) => {
                    if (e2.key === 'a') setView('account');
                    else if (e2.key === 'r') setView('role');
                    else if (e2.key === 'l') setView('location');
                    window.removeEventListener('keydown', handler, true);
                };
                window.addEventListener('keydown', handler, true);
                setTimeout(() => window.removeEventListener('keydown', handler, true), 1200);
            }
        };
        window.addEventListener('keydown', onKey);
        return () => window.removeEventListener('keydown', onKey);
    }, [flatRowIds, focusedRowId, toggleExpanded]);

    // ── Render ────────────────────────────────────────────────────────────
    const meta = `${summary.charCount} characters across ${summary.accountCount} ${
        summary.accountCount === 1 ? 'account' : 'accounts'
    } · ${formatSP(summary.totalSp)} SP total`;

    const headerActions = (
        <>
            <SegmentedControl
                value={view}
                onChange={setView}
                ariaLabel="Group by"
                options={VIEW_OPTIONS}
            />
            <Tooltip title={`Sort ${sortOrder === 'asc' ? 'descending' : 'ascending'}`}>
                <IconButton
                    aria-label="Sort"
                    onClick={() => setSortOrder((s) => s === 'asc' ? 'desc' : 'asc')}
                    size="small"
                    className="!h-8 !w-8 !rounded-md !text-ink-2 hover:!bg-surface-2 hover:!text-ink-1"
                >
                    <span className="font-mono text-meta tabular">{sortOrder === 'asc' ? 'A→Z' : 'Z→A'}</span>
                </IconButton>
            </Tooltip>
            <Tooltip title={showHidden ? 'Hide hidden accounts' : 'Show hidden accounts'}>
                <IconButton
                    aria-label={showHidden ? 'Hide hidden accounts' : 'Show hidden accounts'}
                    onClick={() => setShowHidden((v) => !v)}
                    size="small"
                    className="!h-8 !w-8 !rounded-md hover:!bg-surface-2"
                    sx={{ color: showHidden ? 'var(--accent)' : 'var(--ink-3)' }}
                >
                    {showHidden ? <VisibilityOutlined fontSize="small" /> : <VisibilityOffOutlined fontSize="small" />}
                </IconButton>
            </Tooltip>
            <Tooltip title="Refresh from ESI · R">
                <span>
                    <IconButton
                        ref={refreshRef}
                        aria-label="Refresh from ESI"
                        onClick={handleRefreshAll}
                        disabled={isRefreshing}
                        size="small"
                        className="!h-8 !w-8 !rounded-md !text-ink-2 hover:!bg-surface-2 hover:!text-ink-1"
                    >
                        <RefreshIcon
                            fontSize="small"
                            sx={isRefreshing ? {
                                animation: 'co-spin 1s linear infinite',
                                '@keyframes co-spin': { to: { transform: 'rotate(360deg)' } },
                            } : undefined}
                        />
                    </IconButton>
                </span>
            </Tooltip>
        </>
    );

    return (
        <div className="px-6 pb-12 pt-8 max-w-[1280px] mx-auto">
            <Subheader
                title="Character Overview"
                meta={meta}
                actions={headerActions}
                tip={overviewInstructions}
            />

            {/* Filter bar */}
            <div className="mb-4 flex items-center gap-3">
                <div className="flex-1 max-w-md flex items-center gap-2 rounded-md border border-rule-1 bg-surface-1 px-2.5 py-1.5 focus-within:border-accent">
                    <SearchIcon fontSize="small" sx={{ color: 'var(--ink-3)' }} />
                    <input
                        ref={filterRef}
                        value={filter}
                        onChange={(e) => setFilter(e.target.value)}
                        placeholder="Filter by name, location, role, or ship"
                        aria-label="Filter characters"
                        className="bg-transparent flex-1 outline-none text-body text-ink-1 placeholder:text-ink-3"
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
                <span className="text-meta text-ink-3 tabular">
                    {filter ? `${filteredCharacters.length} match${filteredCharacters.length === 1 ? '' : 'es'}` : null}
                </span>
            </div>

            {/* Empty (zero accounts) */}
            {accounts.length === 0 ? (
                <div className="mt-16 text-center">
                    <p className="text-body text-ink-2">No accounts found.</p>
                    <p className="mt-1 text-meta text-ink-3">Add an EVE account from the header to get started.</p>
                </div>
            ) : filteredCharacters.length === 0 ? (
                <div className="mt-16 text-center">
                    <p className="text-body text-ink-2">
                        No characters match <span className="font-mono text-ink-1">{filter}</span>.
                    </p>
                    <button
                        type="button"
                        onClick={() => setFilter('')}
                        className="mt-2 text-meta text-accent hover:text-accent-strong"
                    >
                        Clear filter (Esc)
                    </button>
                </div>
            ) : (
                <div role="table" aria-label="Characters" className="rounded-lg border border-rule-1 bg-surface-1 overflow-hidden">
                    {/* Column header */}
                    <div
                        role="row"
                        className="hidden sm:grid grid-cols-[28px_minmax(180px,1.4fr)_1fr_84px_92px_minmax(140px,1fr)_minmax(120px,0.8fr)_28px] gap-3 px-4 py-2 text-meta text-ink-3 border-b border-rule-1 bg-surface-2"
                    >
                        <div role="columnheader" aria-label="Status" />
                        <div role="columnheader">Character</div>
                        <div role="columnheader">Account</div>
                        <div role="columnheader" className="text-right">SP</div>
                        <div role="columnheader" className="text-right">Queue</div>
                        <div role="columnheader">Location</div>
                        <div role="columnheader">Role</div>
                        <div role="columnheader" aria-label="Actions" />
                    </div>

                    {groups.map((group) => (
                        <GroupBlock
                            key={group.key}
                            groupKey={group.key}
                            view={view}
                            characters={group.characters}
                            roles={roles}
                            skillConversions={skillConversions}
                            expanded={expanded}
                            onToggleExpand={toggleExpanded}
                            focusedRowId={focusedRowId}
                            setFocusedRowId={setFocusedRowId}
                            onUpdateCharacter={handleUpdateCharacter}
                            onRemoveCharacter={handleRemoveCharacter}
                            onUpdateAccount={(accountId, updates) =>
                                execute(() => updateAccount(accountId, updates), { successMessage: 'Account updated' })
                            }
                            onRemoveAccount={(accountId) =>
                                execute(() => deleteAccount(accountId), { successMessage: 'Account removed' })
                            }
                            allAccounts={accounts}
                        />
                    ))}
                </div>
            )}

            {/* Shortcut legend — text-micro, easy to ignore, available when wanted. */}
            <div className="mt-4 flex flex-wrap items-center gap-x-4 gap-y-1 text-micro text-ink-3 font-mono">
                <span><Kbd>/</Kbd> filter</span>
                <span><Kbd>j</Kbd>/<Kbd>k</Kbd> move</span>
                <span><Kbd>Enter</Kbd> expand</span>
                <span><Kbd>Esc</Kbd> collapse</span>
                <span><Kbd>g</Kbd> then <Kbd>a</Kbd>/<Kbd>r</Kbd>/<Kbd>l</Kbd> group</span>
                <span><Kbd>r</Kbd> refresh</span>
            </div>
        </div>
    );
};

// ── Group block ──────────────────────────────────────────────────────────
function GroupBlock({
    groupKey,
    view,
    characters,
    roles,
    skillConversions,
    expanded,
    onToggleExpand,
    focusedRowId,
    setFocusedRowId,
    onUpdateCharacter,
    onRemoveCharacter,
    onUpdateAccount,
    onRemoveAccount,
    allAccounts,
}) {
    // Group meta on the right side (chars + SP).
    const groupSp = characters.reduce(
        (sum, ch) => sum + (ch.Character?.CharacterSkillsResponse?.total_sp || 0), 0,
    );
    const account = view === 'account'
        ? allAccounts.find((a) => (a.Name || 'Unknown Account') === groupKey)
        : null;

    return (
        <section role="rowgroup">
            <header
                role="row"
                className="sticky top-0 z-10 flex items-center justify-between gap-3 px-4 py-2 border-b border-rule-1 bg-surface-2"
            >
                <div className="flex items-center gap-2 min-w-0">
                    <h2 className="text-h3 text-ink-1 truncate uppercase tracking-wide">
                        {groupKey}
                    </h2>
                    {view === 'account' && account && account.Visible === false ? (
                        <span className="px-1.5 py-0.5 rounded-sm border border-rule-1 text-micro text-ink-3 uppercase tracking-wide">
                            Hidden
                        </span>
                    ) : null}
                </div>
                <div className="flex items-center gap-3 shrink-0">
                    <span className="text-meta text-ink-3 tabular">
                        {characters.length} {characters.length === 1 ? 'char' : 'chars'} · {formatSP(groupSp)} SP
                    </span>
                    {view === 'account' && account ? (
                        <AccountMenu
                            account={account}
                            onUpdate={onUpdateAccount}
                            onRemove={onRemoveAccount}
                        />
                    ) : null}
                </div>
            </header>

            {characters.map((ch, idx) => (
                <CharacterRow
                    key={ch.Character.CharacterID}
                    character={ch}
                    isLast={idx === characters.length - 1}
                    isExpanded={expanded.has(ch.Character.CharacterID)}
                    onToggleExpand={() => onToggleExpand(ch.Character.CharacterID)}
                    isFocused={focusedRowId === ch.Character.CharacterID}
                    onFocus={() => setFocusedRowId(ch.Character.CharacterID)}
                    roles={roles}
                    skillConversions={skillConversions}
                    onUpdateCharacter={onUpdateCharacter}
                    onRemoveCharacter={onRemoveCharacter}
                />
            ))}
        </section>
    );
}

// ── Account menu (kebab) on group header ────────────────────────────────
function AccountMenu({ account, onUpdate, onRemove }) {
    const [anchor, setAnchor] = useState(null);
    const open = Boolean(anchor);
    const close = () => setAnchor(null);
    return (
        <>
            <Tooltip title="Account actions">
                <IconButton
                    aria-label="Account actions"
                    size="small"
                    onClick={(e) => setAnchor(e.currentTarget)}
                    className="!h-7 !w-7 !rounded-md !text-ink-3 hover:!bg-surface-3 hover:!text-ink-1"
                >
                    <MoreVertIcon fontSize="small" />
                </IconButton>
            </Tooltip>
            <Menu anchorEl={anchor} open={open} onClose={close}>
                <MenuItem onClick={() => { onUpdate(account.ID, { Visible: !(account.Visible !== false) }); close(); }}>
                    {account.Visible === false ? 'Show account' : 'Hide account'}
                </MenuItem>
                <MenuItem onClick={() => { onUpdate(account.ID, { Status: account.Status === 'Alpha' ? 'Omega' : 'Alpha' }); close(); }}>
                    Toggle Alpha / Omega ({account.Status === 'Alpha' ? 'α' : 'Ω'})
                </MenuItem>
                <MenuItem
                    onClick={() => { onRemove(account.Name); close(); }}
                    sx={{ color: 'var(--status-error)' }}
                >
                    Remove account
                </MenuItem>
            </Menu>
        </>
    );
}

// ── Row ──────────────────────────────────────────────────────────────────
function CharacterRow({
    character,
    isLast,
    isExpanded,
    onToggleExpand,
    isFocused,
    onFocus,
    roles,
    skillConversions,
    onUpdateCharacter,
    onRemoveCharacter,
}) {
    const c = character.Character || {};
    const id = c.CharacterID;
    const portrait = `https://images.evetech.net/characters/${id}/portrait?size=64`;
    const sp = c.CharacterSkillsResponse?.total_sp || 0;
    const eta = deriveQueueEta(character);
    const status = deriveStatus(character);
    const role = character.Role || '';

    const openZkill = (e) => {
        e.stopPropagation();
        const url = `https://zkillboard.com/character/${id}/`;
        if (window.electronAPI?.openExternal) window.electronAPI.openExternal(url);
        else window.open(url, '_blank', 'noopener,noreferrer');
    };

    return (
        <>
            <div
                role="row"
                tabIndex={0}
                aria-expanded={isExpanded}
                onClick={onToggleExpand}
                onFocus={onFocus}
                onKeyDown={(e) => {
                    if (e.key === 'Enter' || e.key === ' ') {
                        e.preventDefault();
                        onToggleExpand();
                    }
                }}
                className={[
                    'grid grid-cols-[28px_minmax(180px,1.4fr)_1fr_84px_92px_minmax(140px,1fr)_minmax(120px,0.8fr)_28px]',
                    'gap-3 px-4 items-center',
                    'h-10 cursor-pointer select-none',
                    'transition-colors duration-fast ease-out-quart',
                    !isLast ? 'border-b border-rule-1' : '',
                    isFocused ? 'bg-surface-2 shadow-rail-accent' : 'hover:bg-surface-2',
                ].join(' ')}
            >
                {/* Status + expansion glyph (status when collapsed, chevron when expanded) */}
                <div className="flex items-center justify-center">
                    {isExpanded ? (
                        <ExpandMoreOutlined fontSize="small" sx={{ color: 'var(--accent)' }} />
                    ) : (
                        <StatusDot state={status} />
                    )}
                </div>

                {/* Character: portrait + name */}
                <div className="flex items-center gap-2.5 min-w-0">
                    <img
                        src={portrait}
                        alt=""
                        aria-hidden
                        loading="lazy"
                        className="h-6 w-6 rounded-sm border border-rule-1 shrink-0"
                    />
                    <span className="text-body text-ink-1 truncate">{c.CharacterName}</span>
                    <Tooltip title="Open zKillboard">
                        <IconButton
                            size="small"
                            onClick={openZkill}
                            aria-label={`Open zKillboard for ${c.CharacterName}`}
                            className="!h-6 !w-6 !text-ink-3 hover:!bg-surface-3 hover:!text-ink-1"
                        >
                            <OpenInNewOutlined sx={{ fontSize: 14 }} />
                        </IconButton>
                    </Tooltip>
                </div>

                {/* Account */}
                <span className="text-meta text-ink-3 truncate">
                    {character.accountName}
                </span>

                {/* SP */}
                <span className="text-body text-ink-1 font-mono tabular text-right" data-numeric="true">
                    {formatSP(sp)}
                </span>

                {/* Queue ETA */}
                <span
                    className={`text-body font-mono tabular text-right ${
                        status === 'training' ? 'text-status-training' : 'text-ink-2'
                    }`}
                    data-numeric="true"
                    title={status === 'idle' ? 'No skill queue' : 'Time until queue completes'}
                >
                    {status === 'idle' ? '—' : (eta != null ? formatDuration(eta) : 'queued')}
                </span>

                {/* Location */}
                <span className="text-meta text-ink-2 truncate">
                    {c.LocationName || '—'}
                </span>

                {/* Role */}
                <span className="text-meta text-ink-2 truncate">
                    {role || <span className="text-ink-3">—</span>}
                </span>

                {/* Kebab */}
                <div className="flex items-center justify-center" onClick={(e) => e.stopPropagation()}>
                    <RowMenu character={character} onRemove={onRemoveCharacter} />
                </div>
            </div>

            {isExpanded ? (
                <ExpandedRow
                    character={character}
                    roles={roles}
                    skillConversions={skillConversions}
                    onUpdateCharacter={onUpdateCharacter}
                />
            ) : null}
        </>
    );
}

function RowMenu({ character, onRemove }) {
    const [anchor, setAnchor] = useState(null);
    const open = Boolean(anchor);
    const close = () => setAnchor(null);
    return (
        <>
            <IconButton
                size="small"
                aria-label="Character actions"
                onClick={(e) => { e.stopPropagation(); setAnchor(e.currentTarget); }}
                className="!h-7 !w-7 !rounded-md !text-ink-3 hover:!bg-surface-3 hover:!text-ink-1"
            >
                <MoreVertIcon fontSize="small" />
            </IconButton>
            <Menu anchorEl={anchor} open={open} onClose={close}>
                <MenuItem
                    onClick={() => { onRemove(character.Character.CharacterID); close(); }}
                    sx={{ color: 'var(--status-error)' }}
                >
                    <DeleteOutline fontSize="small" sx={{ mr: 1 }} />
                    Remove character
                </MenuItem>
            </Menu>
        </>
    );
}

// ── Expanded row ─────────────────────────────────────────────────────────
function ExpandedRow({ character, roles, skillConversions, onUpdateCharacter }) {
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
            {/* Left: character meta + role editor */}
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

            {/* Right: skill queue */}
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
}

export default CharacterOverview;

// ── PropTypes ────────────────────────────────────────────────────────────
CharacterOverview.propTypes = {
    roles: PropTypes.array,
    skillConversions: PropTypes.object,
};

GroupBlock.propTypes = {
    groupKey: PropTypes.string.isRequired,
    view: PropTypes.string.isRequired,
    characters: PropTypes.array.isRequired,
    roles: PropTypes.array.isRequired,
    skillConversions: PropTypes.object.isRequired,
    expanded: PropTypes.instanceOf(Set).isRequired,
    onToggleExpand: PropTypes.func.isRequired,
    focusedRowId: PropTypes.number,
    setFocusedRowId: PropTypes.func.isRequired,
    onUpdateCharacter: PropTypes.func.isRequired,
    onRemoveCharacter: PropTypes.func.isRequired,
    onUpdateAccount: PropTypes.func.isRequired,
    onRemoveAccount: PropTypes.func.isRequired,
    allAccounts: PropTypes.array.isRequired,
};

AccountMenu.propTypes = {
    account: PropTypes.object.isRequired,
    onUpdate: PropTypes.func.isRequired,
    onRemove: PropTypes.func.isRequired,
};

CharacterRow.propTypes = {
    character: PropTypes.object.isRequired,
    isLast: PropTypes.bool.isRequired,
    isExpanded: PropTypes.bool.isRequired,
    onToggleExpand: PropTypes.func.isRequired,
    isFocused: PropTypes.bool.isRequired,
    onFocus: PropTypes.func.isRequired,
    roles: PropTypes.array.isRequired,
    skillConversions: PropTypes.object.isRequired,
    onUpdateCharacter: PropTypes.func.isRequired,
    onRemoveCharacter: PropTypes.func.isRequired,
};

RowMenu.propTypes = {
    character: PropTypes.object.isRequired,
    onRemove: PropTypes.func.isRequired,
};

ExpandedRow.propTypes = {
    character: PropTypes.object.isRequired,
    roles: PropTypes.array.isRequired,
    skillConversions: PropTypes.object.isRequired,
    onUpdateCharacter: PropTypes.func.isRequired,
};
