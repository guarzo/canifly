// src/pages/Mapping.jsx
//
// Map EVE character files to user files. The signal users actually rely on
// is the mtime bucket — files saved within the same minute share a swatch
// color and almost always belong together. Everything else is chrome.
//
// Redesigned per DESIGN.md: dense rows, neutrals + status, no card grid,
// no framer-motion, no glass, no gradient.

import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import PropTypes from 'prop-types';
import { toast } from 'react-toastify';
import { Search as SearchIcon } from '@mui/icons-material';

import Subheader from '../components/ui/Subheader.jsx';
import SegmentedControl from '../components/ui/SegmentedControl.jsx';
import Kbd from '../components/ui/Kbd.jsx';
import AccountCard from '../components/profiles/MapAccountCard.jsx';
import CharacterCard from '../components/profiles/MapCharacterCard.jsx';
import { useConfirmDialog } from '../hooks/useConfirmDialog.jsx';
import { associateCharacter, unassociateCharacter } from '../api/accountsApi';
import { mappingInstructions } from '../utils/instructions';
import { useAppData } from '../hooks/useAppData';
import { log } from '../utils/logger';

// ── Storage ──────────────────────────────────────────────────────────────
const LS = {
    filter: 'mp.filter',
    sortOrder: 'mp.sortOrder',
    view: 'mp.view',
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

// ── Helpers ──────────────────────────────────────────────────────────────
function roundToMinute(mtime) {
    const date = new Date(mtime);
    date.setSeconds(0);
    return date.toISOString();
}

// A small palette of distinct mtime-bucket swatches. These are intentionally
// not the accent or status hues — they encode grouping, not state.
const SWATCH_PALETTE = [
    'oklch(0.72 0.13 145)', // green
    'oklch(0.74 0.14 25)',  // red-orange
    'oklch(0.80 0.13 80)',  // amber
    'oklch(0.70 0.14 305)', // violet
    'oklch(0.75 0.12 200)', // teal
    'oklch(0.72 0.13 350)', // pink
    'oklch(0.78 0.10 240)', // blue
    'oklch(0.74 0.10 120)', // moss
];

const VIEW_OPTIONS = [
    { value: 'all', label: 'All' },
    { value: 'unmatched', label: 'Unmatched' },
    { value: 'matched', label: 'Matched' },
];

const Mapping = ({ associations: initialAssociations, subDirs }) => {
    const { refreshData } = useAppData();
    const [accounts, setAccounts] = useState([]);
    const [availableCharacters, setAvailableCharacters] = useState([]);
    const [associations, setAssociations] = useState(initialAssociations);
    const [mtimeToColor, setMtimeToColor] = useState({});
    const [showConfirmDialog, confirmDialog] = useConfirmDialog();

    const [filter, setFilter] = useState(() => readLS(LS.filter, ''));
    const [sortOrder, setSortOrder] = useState(() => readLS(LS.sortOrder, 'mtime-desc'));
    const [view, setView] = useState(() => readLS(LS.view, 'all'));

    useEffect(() => writeLS(LS.filter, filter), [filter]);
    useEffect(() => writeLS(LS.sortOrder, sortOrder), [sortOrder]);
    useEffect(() => writeLS(LS.view, view), [view]);

    const filterRef = useRef(null);

    // Build the unique account & character sets from the raw subDirs payload.
    useEffect(() => {
        if (!subDirs || subDirs.length === 0) {
            setAccounts([]);
            setAvailableCharacters([]);
            setMtimeToColor({});
            return;
        }

        const userMap = {};
        subDirs.forEach((mapping) => {
            mapping.availableUserFiles.forEach((userFile) => {
                const roundedMtime = roundToMinute(userFile.mtime);
                if (
                    !userMap[userFile.userId] ||
                    new Date(roundedMtime) > new Date(userMap[userFile.userId].mtime)
                ) {
                    userMap[userFile.userId] = { ...userFile, mtime: roundedMtime };
                }
            });
        });
        const uniqueAccounts = Object.values(userMap).sort(
            (a, b) => new Date(b.mtime) - new Date(a.mtime),
        );

        const charMap = {};
        subDirs.forEach((mapping) => {
            mapping.availableCharFiles.forEach((charFile) => {
                const roundedMtime = roundToMinute(charFile.mtime);
                const { charId } = charFile;
                if (!charMap[charId] || new Date(roundedMtime) > new Date(charMap[charId].mtime)) {
                    charMap[charId] = { ...charFile, mtime: roundedMtime, profile: mapping.profile };
                }
            });
        });

        const associatedCharIds = new Set(associations.map((a) => a.charId));
        const uniqueChars = Object.values(charMap)
            .filter((ch) => !associatedCharIds.has(ch.charId))
            .sort((a, b) => new Date(b.mtime) - new Date(a.mtime));

        // Color map keyed by rounded mtime; only assign colors to mtimes that
        // appear in *both* sides (or in any list — same logic as before).
        const allMtimes = [
            ...uniqueAccounts.map((a) => a.mtime),
            ...uniqueChars.map((c) => c.mtime),
        ];
        const ordered = Array.from(new Set(allMtimes)).sort(
            (a, b) => new Date(a) - new Date(b),
        );
        const colorMapping = ordered.reduce((acc, mtime, index) => {
            acc[mtime] = SWATCH_PALETTE[index % SWATCH_PALETTE.length];
            return acc;
        }, {});

        setAccounts(uniqueAccounts);
        setAvailableCharacters(uniqueChars);
        setMtimeToColor(colorMapping);
    }, [subDirs, associations]);

    log('Mapping mtime palette size', Object.keys(mtimeToColor).length);

    // ── Drag & drop ──────────────────────────────────────────────────────
    const handleDragStart = (event, charId) => {
        event.dataTransfer.setData('text/plain', charId);
    };

    const handleDrop = async (event, userId, userName) => {
        event.preventDefault();
        const charId = event.dataTransfer.getData('text/plain');

        if (associations.some((assoc) => assoc.charId === charId)) {
            toast.error('That character is already associated!');
            return;
        }

        const char = availableCharacters.find((c) => c.charId === charId);
        const charName = char?.name;
        if (!char) {
            toast.error('Character not found.');
            return;
        }

        const confirmAssoc = await showConfirmDialog({
            title: 'Confirm Association',
            message: `Associate "${charName}" with account "${userName}"?`,
        });
        if (!confirmAssoc.isConfirmed) return;

        const result = await associateCharacter(userId, charId, userName, charName);
        if (result && result.success) {
            toast.success(result.message);
            setAvailableCharacters((prev) => prev.filter((ch) => ch.charId !== charId));
            setAssociations((prev) => [...prev, { userId, charId, charName, mtime: char.mtime }]);
            if (refreshData) await refreshData();
        }
    };

    const handleUnassociate = async (userId, charId, charName, userName) => {
        const confirmUnassoc = await showConfirmDialog({
            title: 'Confirm Unassociation',
            message: `Unassociate "${charName}" from account "${userName}"?`,
        });
        if (!confirmUnassoc.isConfirmed) return;

        const result = await unassociateCharacter(userId, charId, userName, charName);
        if (result && result.success) {
            toast.success(result.message);
            setAssociations((prev) =>
                prev.filter((a) => a.charId !== charId || a.userId !== userId),
            );
            if (refreshData) await refreshData();
        }
    };

    // ── Derived lists ────────────────────────────────────────────────────
    const q = filter.trim().toLowerCase();

    const matchesQuery = useCallback((str) => {
        if (!q) return true;
        return (str || '').toLowerCase().includes(q);
    }, [q]);

    const filteredAccounts = useMemo(() => {
        const list = accounts.filter((a) => {
            const assocs = associations.filter((x) => x.userId === a.userId);
            const hasMatchedFilter =
                view === 'all' ||
                (view === 'matched' && assocs.length > 0) ||
                (view === 'unmatched' && assocs.length === 0);
            if (!hasMatchedFilter) return false;
            if (!q) return true;
            const nameHit = matchesQuery(a.name) || matchesQuery(a.userId);
            const assocHit = assocs.some(
                (x) => matchesQuery(x.charName) || matchesQuery(x.charId),
            );
            return nameHit || assocHit;
        });
        list.sort((a, b) => {
            switch (sortOrder) {
                case 'mtime-asc':
                    return new Date(a.mtime) - new Date(b.mtime);
                case 'name-asc':
                    return (a.name || '').localeCompare(b.name || '');
                case 'name-desc':
                    return (b.name || '').localeCompare(a.name || '');
                case 'mtime-desc':
                default:
                    return new Date(b.mtime) - new Date(a.mtime);
            }
        });
        return list;
    }, [accounts, associations, view, q, sortOrder, matchesQuery]);

    const filteredCharacters = useMemo(() => {
        if (view === 'matched') return [];
        const list = availableCharacters.filter((c) => {
            if (!q) return true;
            return matchesQuery(c.name) || matchesQuery(c.charId);
        });
        list.sort((a, b) => {
            switch (sortOrder) {
                case 'mtime-asc':
                    return new Date(a.mtime) - new Date(b.mtime);
                case 'name-asc':
                    return (a.name || '').localeCompare(b.name || '');
                case 'name-desc':
                    return (b.name || '').localeCompare(a.name || '');
                case 'mtime-desc':
                default:
                    return new Date(b.mtime) - new Date(a.mtime);
            }
        });
        return list;
    }, [availableCharacters, view, q, sortOrder, matchesQuery]);

    // ── Keyboard ─────────────────────────────────────────────────────────
    useEffect(() => {
        const onKey = (e) => {
            const tag = (e.target?.tagName || '').toLowerCase();
            const isTyping = tag === 'input' || tag === 'textarea' || e.target?.isContentEditable;
            if (e.key === '/' && !isTyping) {
                e.preventDefault();
                filterRef.current?.focus();
            } else if (e.key === 'Escape' && document.activeElement === filterRef.current) {
                setFilter('');
            }
        };
        window.addEventListener('keydown', onKey);
        return () => window.removeEventListener('keydown', onKey);
    }, []);

    // ── Header ───────────────────────────────────────────────────────────
    const associatedCount = associations.length;
    const totalChars = associatedCount + availableCharacters.length;
    const meta = `${accounts.length} user file${accounts.length === 1 ? '' : 's'} · ${associatedCount}/${totalChars} characters mapped`;

    const sortLabels = {
        'mtime-desc': 'Newest first',
        'mtime-asc': 'Oldest first',
        'name-asc': 'Name A→Z',
        'name-desc': 'Name Z→A',
    };
    const cycleSort = () => {
        const order = ['mtime-desc', 'mtime-asc', 'name-asc', 'name-desc'];
        setSortOrder(order[(order.indexOf(sortOrder) + 1) % order.length]);
    };

    const headerActions = (
        <>
            <SegmentedControl
                value={view}
                onChange={setView}
                ariaLabel="Filter mapping view"
                options={VIEW_OPTIONS.map((o) => ({ ...o, showLabel: true }))}
            />
            <button
                type="button"
                onClick={cycleSort}
                aria-label={`Sort: ${sortLabels[sortOrder]}`}
                title={`Sort: ${sortLabels[sortOrder]}`}
                className="inline-flex items-center h-8 px-2.5 rounded-md border border-rule-1 bg-surface-1 text-meta text-ink-2 hover:bg-surface-2 hover:text-ink-1 font-mono tabular transition-colors duration-fast ease-out-quart"
            >
                {sortLabels[sortOrder]}
            </button>
        </>
    );

    // ── Render ───────────────────────────────────────────────────────────
    return (
        <div className="px-6 pb-12 pt-8 max-w-[1280px] mx-auto">
            <Subheader
                title="Mapping"
                meta={meta}
                actions={headerActions}
                tip={mappingInstructions}
            />

            {/* Filter bar */}
            <div className="mb-4 flex items-center gap-3">
                <div className="flex-1 max-w-md flex items-center gap-2 rounded-md border border-rule-1 bg-surface-1 px-2.5 py-1.5 focus-within:border-accent">
                    <SearchIcon fontSize="small" sx={{ color: 'var(--ink-3)' }} />
                    <input
                        ref={filterRef}
                        value={filter}
                        onChange={(e) => setFilter(e.target.value)}
                        placeholder="Filter by character or user name / id"
                        aria-label="Filter mapping"
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
                <span className="text-meta text-ink-3 tabular">
                    {q
                        ? `${filteredAccounts.length} user · ${filteredCharacters.length} char`
                        : null}
                </span>
            </div>

            {accounts.length === 0 ? (
                <div className="mt-16 text-center">
                    <p className="text-body text-ink-2">No accounts found.</p>
                    <p className="mt-1 text-meta text-ink-3">
                        EVE settings directories will appear here once detected.
                    </p>
                </div>
            ) : (
                <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                    {/* Left column: user files */}
                    <section aria-label="User files">
                        <div className="mb-2 flex items-baseline justify-between">
                            <h2 className="text-h3 text-ink-1 uppercase tracking-wide">
                                User files
                            </h2>
                            <span className="text-meta text-ink-3 tabular">
                                {filteredAccounts.length} of {accounts.length}
                            </span>
                        </div>

                        {filteredAccounts.length === 0 ? (
                            <div className="rounded-lg border border-rule-1 bg-surface-1 px-4 py-8 text-center text-meta text-ink-3">
                                No user files match the current filter.
                            </div>
                        ) : (
                            <div className="flex flex-col gap-3">
                                {filteredAccounts.map((mapping) => (
                                    <AccountCard
                                        key={`${mapping.userId}-${mapping.mtime}`}
                                        mapping={mapping}
                                        associations={associations}
                                        handleUnassociate={handleUnassociate}
                                        handleDrop={handleDrop}
                                        mtimeToColor={mtimeToColor}
                                    />
                                ))}
                            </div>
                        )}
                    </section>

                    {/* Right column: unassociated characters */}
                    {view !== 'matched' ? (
                        <section aria-label="Unassociated characters">
                            <div className="mb-2 flex items-baseline justify-between">
                                <h2 className="text-h3 text-ink-1 uppercase tracking-wide">
                                    Unassociated characters
                                </h2>
                                <span className="text-meta text-ink-3 tabular">
                                    {filteredCharacters.length}
                                </span>
                            </div>

                            {availableCharacters.length === 0 ? (
                                <div className="rounded-lg border border-rule-1 bg-surface-1 px-4 py-8 text-center text-meta text-ink-3">
                                    All characters mapped.
                                </div>
                            ) : filteredCharacters.length === 0 ? (
                                <div className="rounded-lg border border-rule-1 bg-surface-1 px-4 py-8 text-center text-meta text-ink-3">
                                    No characters match the current filter.
                                </div>
                            ) : (
                                <div
                                    role="list"
                                    data-testid="available-characters"
                                    className="rounded-lg border border-rule-1 bg-surface-1 overflow-hidden"
                                >
                                    {filteredCharacters.map((char) => (
                                        <CharacterCard
                                            key={`${char.charId}-${char.mtime}`}
                                            char={char}
                                            handleDragStart={handleDragStart}
                                            mtimeToColor={mtimeToColor}
                                        />
                                    ))}
                                </div>
                            )}
                        </section>
                    ) : null}
                </div>
            )}

            {/* Shortcut legend */}
            <div className="mt-6 flex flex-wrap items-center gap-x-4 gap-y-1 text-micro text-ink-3 font-mono">
                <span><Kbd>/</Kbd> filter</span>
                <span><Kbd>Esc</Kbd> clear</span>
                <span>drag a character row onto a user file to associate</span>
            </div>

            {confirmDialog}
        </div>
    );
};

Mapping.propTypes = {
    associations: PropTypes.array.isRequired,
    subDirs: PropTypes.array.isRequired,
};

export default Mapping;
