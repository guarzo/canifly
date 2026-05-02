// src/pages/Sync.jsx
//
// Sync — pair an EVE character file with a user file per profile, then push the
// settings out. Redesigned per DESIGN.md: dense table, calm accent action, no
// glass / gradients / hover-lift. The page makes the selected pair obvious
// (status dot + arrow + accent fill) and keeps the destructive bits (reset,
// directory choice, backup) in a quiet header action cluster.

import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import PropTypes from 'prop-types';
import { toast } from 'react-toastify';
import {
    BackupOutlined,
    FolderOpenOutlined,
    RestartAltOutlined,
    Search as SearchIcon,
} from '@mui/icons-material';
import { IconButton, Tooltip } from '@mui/material';

import Subheader from '../components/ui/Subheader.jsx';
import Kbd from '../components/ui/Kbd.jsx';
import StatusDot from '../components/ui/StatusDot.jsx';
import SyncProfileRow from '../components/sync/SyncProfileRow.jsx';
import { useConfirmDialog } from '../hooks/useConfirmDialog.jsx';
import { syncInstructions } from '../utils/instructions.jsx';
import {
    saveUserSelections,
    syncSubdirectory,
    syncAllSubdirectories,
    chooseSettingsDir,
    backupDirectory,
    resetToDefaultDirectory,
} from '../api/syncApi';
import { logger } from '../utils/logger';

// ── Persistence ───────────────────────────────────────────────────────────
const LS = {
    profile: 'sy.profile',   // last focused profile key
    filter: 'sy.filter',     // filter text
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

// ── Page ──────────────────────────────────────────────────────────────────
const Sync = ({
    settingsData = [],
    associations = [],
    currentSettingsDir = '',
    userSelections = {},
    lastBackupDir = '',
}) => {
    const [isLoading, setIsLoading] = useState(false);
    const [selections, setSelections] = useState({});
    const [showConfirmDialog, confirmDialog] = useConfirmDialog();
    const [isDefaultDir, setIsDefaultDir] = useState(false);
    const [filter, setFilter] = useState(() => readLS(LS.filter, ''));
    const [focusedProfile, setFocusedProfile] = useState(() => readLS(LS.profile, null));

    const filterRef = useRef(null);

    useEffect(() => writeLS(LS.filter, filter), [filter]);
    useEffect(() => writeLS(LS.profile, focusedProfile), [focusedProfile]);

    // Initialize selections whenever settingsData/userSelections change.
    useEffect(() => {
        if (!settingsData || settingsData.length === 0) return;
        const initial = { ...userSelections };
        for (const sd of settingsData) {
            if (!initial[sd.profile]) initial[sd.profile] = { charId: '', userId: '' };
        }
        setSelections(initial);
    }, [settingsData, userSelections]);

    const persistSelections = useCallback(async (next) => {
        await saveUserSelections(next);
    }, []);

    const handleSelectionChange = useCallback((profile, field, value) => {
        setSelections((prev) => {
            const next = {
                ...prev,
                [profile]: { ...prev[profile], [field]: value },
            };
            // Auto-select user when a known character is associated.
            if (field === 'charId' && value) {
                const assoc = associations.find((a) => a.charId === value);
                if (assoc) next[profile].userId = assoc.userId;
            }
            persistSelections(next);
            return next;
        });
    }, [associations, persistSelections]);

    const handleSync = useCallback(async (profile) => {
        const sel = selections[profile] || {};
        const { charId, userId } = sel;
        if (!charId || !userId) {
            toast.error('Select a character file and a user file first.');
            return;
        }
        const ok = await showConfirmDialog({
            title: 'Sync profile',
            message: `Sync ${profile.replace(/^settings_/, '')} with the chosen pair?`,
        });
        if (!ok.isConfirmed) return;
        try {
            setIsLoading(true);
            logger.debug('Sync', profile, userId, charId);
            const result = await syncSubdirectory(profile, userId, charId);
            if (result?.success) toast.success(result.message);
        } catch (err) {
            logger.error('Sync failed', err);
        } finally {
            setIsLoading(false);
        }
    }, [selections, showConfirmDialog]);

    const handleSyncAll = useCallback(async (profile) => {
        const sel = selections[profile] || {};
        const { charId, userId } = sel;
        if (!charId || !userId) {
            toast.error('Select a character file and a user file first.');
            return;
        }
        const ok = await showConfirmDialog({
            title: 'Sync all profiles',
            message: 'Apply the selected pair to every profile? This overwrites their current pair.',
        });
        if (!ok.isConfirmed) return;
        try {
            setIsLoading(true);
            const result = await syncAllSubdirectories(profile, userId, charId);
            if (result?.success) toast.success(`Sync-All complete: ${result.message}`);
        } catch (err) {
            logger.error('Sync-all failed', err);
        } finally {
            setIsLoading(false);
        }
    }, [selections, showConfirmDialog]);

    const handleChooseSettingsDir = useCallback(async () => {
        try {
            setIsLoading(true);
            const chosen = await window.electronAPI.chooseDirectory();
            if (!chosen) { toast.info('No directory chosen.'); setIsLoading(false); return; }
            const result = await chooseSettingsDir(chosen);
            if (result?.success) {
                setIsDefaultDir(false);
                toast.success(`Settings directory: ${chosen}`);
            }
        } catch (err) {
            logger.error('chooseSettingsDir failed', err);
        } finally {
            setIsLoading(false);
        }
    }, []);

    const handleBackup = useCallback(async () => {
        try {
            setIsLoading(true);
            const chosen = await window.electronAPI.chooseDirectory(lastBackupDir || '');
            if (!chosen) { toast.info('Backup canceled.'); setIsLoading(false); return; }
            toast.info('Starting backup…');
            const result = await backupDirectory(currentSettingsDir, chosen);
            if (result?.success) toast.success(result.message);
        } catch (err) {
            logger.error('backup failed', err);
        } finally {
            setIsLoading(false);
        }
    }, [currentSettingsDir, lastBackupDir]);

    const handleResetToDefault = useCallback(async () => {
        const ok = await showConfirmDialog({
            title: 'Reset to default',
            message: 'Reset the settings directory to the default Tranquility location?',
        });
        if (!ok.isConfirmed) return;
        try {
            setIsLoading(true);
            const result = await resetToDefaultDirectory();
            if (result?.success) {
                setIsDefaultDir(true);
                toast.success('Reset to default: Tranquility');
            }
        } catch (err) {
            logger.error('resetToDefault failed', err);
        } finally {
            setIsLoading(false);
        }
    }, [showConfirmDialog]);

    // Filter pass.
    const filtered = useMemo(() => {
        const q = filter.trim().toLowerCase();
        if (!q) return settingsData;
        return settingsData.filter((sd) => {
            const name = sd.profile.toLowerCase();
            if (name.includes(q)) return true;
            const sel = selections[sd.profile];
            if (sel?.charId && sel.charId.toLowerCase().includes(q)) return true;
            if (sel?.userId && sel.userId.toLowerCase().includes(q)) return true;
            return false;
        });
    }, [settingsData, selections, filter]);

    // Header summary.
    const summary = useMemo(() => {
        const total = settingsData.length;
        const paired = settingsData.reduce((n, sd) => {
            const sel = selections[sd.profile];
            return n + (sel?.charId && sel?.userId ? 1 : 0);
        }, 0);
        return { total, paired };
    }, [settingsData, selections]);

    // Keyboard shortcuts.
    useEffect(() => {
        const onKey = (e) => {
            const tag = (e.target?.tagName || '').toLowerCase();
            const isTyping = tag === 'input' || tag === 'textarea' || e.target?.isContentEditable;
            if (e.key === '/' && !isTyping) {
                e.preventDefault();
                filterRef.current?.focus();
            }
        };
        window.addEventListener('keydown', onKey);
        return () => window.removeEventListener('keydown', onKey);
    }, []);

    const meta = currentSettingsDir
        ? `${summary.paired} / ${summary.total} paired · ${currentSettingsDir}`
        : `${summary.paired} / ${summary.total} paired`;

    const headerActions = (
        <>
            <Tooltip title="Backup current settings directory">
                <span>
                    <IconButton
                        aria-label="Backup settings"
                        onClick={handleBackup}
                        disabled={isLoading}
                        size="small"
                        className="!h-8 !w-8 !rounded-md !text-ink-2 hover:!bg-surface-2 hover:!text-ink-1"
                    >
                        <BackupOutlined fontSize="small" />
                    </IconButton>
                </span>
            </Tooltip>
            <Tooltip title="Choose settings directory">
                <span>
                    <IconButton
                        aria-label="Choose settings directory"
                        onClick={handleChooseSettingsDir}
                        disabled={isLoading}
                        size="small"
                        className="!h-8 !w-8 !rounded-md !text-ink-2 hover:!bg-surface-2 hover:!text-ink-1"
                    >
                        <FolderOpenOutlined fontSize="small" />
                    </IconButton>
                </span>
            </Tooltip>
            {!isDefaultDir ? (
                <Tooltip title="Reset to default (Tranquility)">
                    <span>
                        <IconButton
                            aria-label="Reset to default directory"
                            onClick={handleResetToDefault}
                            disabled={isLoading}
                            size="small"
                            className="!h-8 !w-8 !rounded-md !text-ink-2 hover:!bg-surface-2 hover:!text-ink-1"
                        >
                            <RestartAltOutlined fontSize="small" />
                        </IconButton>
                    </span>
                </Tooltip>
            ) : null}
        </>
    );

    const empty = !settingsData || settingsData.length === 0;

    return (
        <div className="px-6 pb-12 pt-8 max-w-[1280px] mx-auto">
            <Subheader
                title="Sync"
                meta={meta}
                actions={headerActions}
                tip={syncInstructions}
            />

            {/* Filter bar. */}
            <div className="mb-4 flex items-center gap-3">
                <div className="flex-1 max-w-md flex items-center gap-2 rounded-md border border-rule-1 bg-surface-1 px-2.5 py-1.5 focus-within:border-accent">
                    <SearchIcon fontSize="small" sx={{ color: 'var(--ink-3)' }} />
                    <input
                        ref={filterRef}
                        value={filter}
                        onChange={(e) => setFilter(e.target.value)}
                        placeholder="Filter profiles, character ids, user ids"
                        aria-label="Filter profiles"
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
                    {filter ? `${filtered.length} match${filtered.length === 1 ? '' : 'es'}` : null}
                </span>
            </div>

            {empty ? (
                <div className="mt-16 text-center">
                    <p className="text-body text-ink-2">No EVE profiles found.</p>
                    <p className="mt-1 text-meta text-ink-3">
                        Choose a settings directory from the header to scan for profiles.
                    </p>
                </div>
            ) : filtered.length === 0 ? (
                <div className="mt-16 text-center">
                    <p className="text-body text-ink-2">
                        No profiles match <span className="font-mono text-ink-1">{filter}</span>.
                    </p>
                    <button
                        type="button"
                        onClick={() => setFilter('')}
                        className="mt-2 text-meta text-accent hover:text-accent-strong"
                    >
                        Clear filter
                    </button>
                </div>
            ) : (
                <div
                    role="table"
                    aria-label="Profiles"
                    aria-busy={isLoading}
                    className="rounded-lg border border-rule-1 bg-surface-1 overflow-hidden"
                >
                    {/* Column header. */}
                    <div
                        role="row"
                        className="hidden sm:grid grid-cols-[20px_minmax(140px,1fr)_minmax(180px,1.4fr)_20px_minmax(180px,1.4fr)_72px_36px] gap-3 px-4 py-2 text-meta text-ink-3 border-b border-rule-1 bg-surface-2"
                    >
                        <div role="columnheader" aria-label="Pair status" />
                        <div role="columnheader">Profile</div>
                        <div role="columnheader">Character file</div>
                        <div role="columnheader" aria-hidden="true" />
                        <div role="columnheader">User file</div>
                        <div role="columnheader" className="text-right">Action</div>
                        <div role="columnheader" aria-label="Sync all" />
                    </div>

                    {filtered.map((sd) => (
                        <SyncProfileRow
                            key={sd.profile}
                            subDir={sd}
                            selection={selections[sd.profile]}
                            onSelectionChange={handleSelectionChange}
                            onSync={handleSync}
                            onSyncAll={handleSyncAll}
                            isLoading={isLoading}
                            isFocused={focusedProfile === sd.profile}
                            onFocus={() => setFocusedProfile(sd.profile)}
                        />
                    ))}
                </div>
            )}

            {/* Legend + shortcut hints. */}
            <div className="mt-4 flex flex-wrap items-center gap-x-4 gap-y-1 text-micro text-ink-3 font-mono">
                <span className="inline-flex items-center gap-1.5">
                    <StatusDot state="ready" /> paired
                </span>
                <span className="inline-flex items-center gap-1.5">
                    <StatusDot state="queued" /> partial
                </span>
                <span className="inline-flex items-center gap-1.5">
                    <StatusDot state="idle" /> empty
                </span>
                <span className="ml-auto"><Kbd>/</Kbd> filter</span>
            </div>

            {confirmDialog}
        </div>
    );
};

Sync.propTypes = {
    settingsData: PropTypes.array,
    associations: PropTypes.array,
    currentSettingsDir: PropTypes.string,
    userSelections: PropTypes.object,
    lastBackupDir: PropTypes.string,
};

export default Sync;
