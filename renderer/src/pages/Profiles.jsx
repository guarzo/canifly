// renderer/src/pages/Profiles.jsx
//
// Merged page: Mapping + Sync as two views in one shell. Default mode is
// computed: if any unmatched characters exist → start in Mapping; else → Sync.
// The user can switch via SegmentedControl; the choice is persisted, and the
// query string `?view=mapping|sync` overrides on entry (so /mapping and /sync
// redirects land in the right mode).

import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { useSearchParams } from 'react-router-dom';
import { toast } from 'react-toastify';
import {
    BackupOutlined,
    FolderOpenOutlined,
    RestartAltOutlined,
    GridOnOutlined,
    SyncOutlined,
} from '@mui/icons-material';
import { IconButton, Tooltip } from '@mui/material';

import PageShell from '../components/ui/PageShell.jsx';
import FilterBar from '../components/ui/FilterBar.jsx';
import SegmentedControl from '../components/ui/SegmentedControl.jsx';
import Kbd from '../components/ui/Kbd.jsx';
import MappingView from '../components/profiles/MappingView.jsx';
import SyncView from '../components/profiles/SyncView.jsx';
import { useConfirmDialog } from '../hooks/useConfirmDialog.jsx';
import { useSyncAction } from '../hooks/useSyncAction.js';
import { useMappingDerivations } from '../hooks/useMappingDerivations.js';
import { mappingInstructions, syncInstructions } from '../utils/instructions';
import {
    chooseSettingsDir,
    backupDirectory,
    resetToDefaultDirectory,
} from '../api/syncApi';

const LS = {
    view: 'pf.view',
    filter: 'pf.filter',
    mapView: 'pf.map.view',
    mapSort: 'pf.map.sort',
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
    { value: 'mapping', label: 'Map', icon: <GridOnOutlined fontSize="small" />, showLabel: true },
    { value: 'sync', label: 'Sync', icon: <SyncOutlined fontSize="small" />, showLabel: true },
];

const MAP_VIEW_OPTIONS = [
    { value: 'all', label: 'All' },
    { value: 'unmatched', label: 'Unmatched' },
    { value: 'matched', label: 'Matched' },
];

const SORT_LABELS = {
    'mtime-desc': 'Newest first',
    'mtime-asc': 'Oldest first',
    'name-asc': 'Name A→Z',
    'name-desc': 'Name Z→A',
};

const Profiles = ({
    subDirs = [],
    associations = [],
    settingsData = [],
    userSelections = {},
    currentSettingsDir = '',
    isDefaultDir = false,
    lastBackupDir = '',
}) => {
    const [params, setParams] = useSearchParams();
    const { run, isLoading } = useSyncAction();
    const [showConfirmDialog, confirmDialog] = useConfirmDialog();

    // Smart default: unmatched characters → Mapping; else Sync. The persisted
    // value wins on second visit; the query string wins over both.
    const { availableCharacters } = useMappingDerivations(subDirs, associations);
    const initialView = useMemo(() => {
        const fromQuery = params.get('view');
        if (fromQuery === 'mapping' || fromQuery === 'sync') return fromQuery;
        const hasUnmatched = availableCharacters.length > 0;
        const stored = readLS(LS.view, null);
        // Only honour the stored preference when it still makes sense:
        // don't stay in mapping mode if there are no unmatched characters.
        if (stored === 'sync') return 'sync';
        if (stored === 'mapping' && hasUnmatched) return 'mapping';
        return hasUnmatched ? 'mapping' : 'sync';
    // We only want this on first mount; subsequent switches are user-driven.
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, []);

    const [view, setView] = useState(initialView);
    const [filter, setFilter] = useState(() => readLS(LS.filter, ''));
    const [mapView, setMapView] = useState(() => readLS(LS.mapView, 'all'));
    const [mapSort, setMapSort] = useState(() => readLS(LS.mapSort, 'mtime-desc'));
    const filterRef = useRef(null);

    useEffect(() => {
        writeLS(LS.view, view);
        // Mirror to URL so a page reload keeps the same mode. Use the
        // functional setParams form to avoid stale closures over `params`.
        setParams((prev) => {
            const next = new URLSearchParams(prev);
            next.set('view', view);
            return next;
        }, { replace: true });
    }, [view, setParams]);
    useEffect(() => writeLS(LS.filter, filter), [filter]);
    useEffect(() => writeLS(LS.mapView, mapView), [mapView]);
    useEffect(() => writeLS(LS.mapSort, mapSort), [mapSort]);

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

    const handleChooseSettingsDir = useCallback(async () => {
        const chosen = await window.electronAPI.chooseDirectory();
        if (!chosen) { toast.info('No directory chosen.'); return; }
        await run(async () => {
            const result = await chooseSettingsDir(chosen);
            return { ...result, message: result?.message || `Settings directory: ${chosen}` };
        }, { errorContext: 'chooseSettingsDir' });
    }, [run]);

    const handleBackup = useCallback(async () => {
        const chosen = await window.electronAPI.chooseDirectory(lastBackupDir || '');
        if (!chosen) { toast.info('Backup canceled.'); return; }
        toast.info('Starting backup…');
        await run(() => backupDirectory(currentSettingsDir, chosen), { errorContext: 'backupDirectory' });
    }, [run, currentSettingsDir, lastBackupDir]);

    const handleResetToDefault = useCallback(async () => {
        const ok = await showConfirmDialog({
            title: 'Reset to default',
            message: 'Reset the settings directory to the default Tranquility location?',
        });
        if (!ok.isConfirmed) return;
        await run(async () => {
            const result = await resetToDefaultDirectory();
            return { ...result, message: 'Reset to default: Tranquility' };
        }, { errorContext: 'resetToDefaultDirectory' });
    }, [run, showConfirmDialog]);

    const cycleSort = () => {
        const order = ['mtime-desc', 'mtime-asc', 'name-asc', 'name-desc'];
        setMapSort(order[(order.indexOf(mapSort) + 1) % order.length]);
    };

    // Header summary — distinct per mode.
    const meta = view === 'mapping'
        ? `${associations.length}/${associations.length + availableCharacters.length} characters mapped`
        : currentSettingsDir
            ? `${settingsData.length} profile${settingsData.length === 1 ? '' : 's'} · ${currentSettingsDir}`
            : `${settingsData.length} profile${settingsData.length === 1 ? '' : 's'}`;

    const headerActions = view === 'sync' ? (
        <>
            <SegmentedControl value={view} onChange={setView} ariaLabel="Mode" options={VIEW_OPTIONS} />
            <Tooltip title="Backup current settings directory">
                <span>
                    <IconButton
                        aria-label="Backup settings"
                        onClick={handleBackup}
                        disabled={isLoading}
                        size="small"
                        className="h-8! w-8! rounded-md! text-ink-2! hover:bg-surface-2! hover:text-ink-1!"
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
                        className="h-8! w-8! rounded-md! text-ink-2! hover:bg-surface-2! hover:text-ink-1!"
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
                            className="h-8! w-8! rounded-md! text-ink-2! hover:bg-surface-2! hover:text-ink-1!"
                        >
                            <RestartAltOutlined fontSize="small" />
                        </IconButton>
                    </span>
                </Tooltip>
            ) : null}
        </>
    ) : (
        <>
            <SegmentedControl value={view} onChange={setView} ariaLabel="Mode" options={VIEW_OPTIONS} />
            <SegmentedControl
                value={mapView}
                onChange={setMapView}
                ariaLabel="Filter mapping view"
                options={MAP_VIEW_OPTIONS.map((o) => ({ ...o, showLabel: true }))}
            />
            <button
                type="button"
                onClick={cycleSort}
                aria-label={`Sort: ${SORT_LABELS[mapSort]}`}
                title={`Sort: ${SORT_LABELS[mapSort]}`}
                className="inline-flex items-center h-8 px-2.5 rounded-md border border-rule-1 bg-surface-1 text-meta text-ink-2 hover:bg-surface-2 hover:text-ink-1 font-mono tabular transition-colors duration-fast ease-out-quart"
            >
                {SORT_LABELS[mapSort]}
            </button>
        </>
    );

    const tip = view === 'mapping' ? mappingInstructions : syncInstructions;
    const placeholder = view === 'mapping'
        ? 'Filter by character or user name / id'
        : 'Filter profiles, character ids, user ids';

    return (
        <PageShell title="Profiles" meta={meta} actions={headerActions} tip={tip}>
            <FilterBar
                inputRef={filterRef}
                value={filter}
                onChange={setFilter}
                placeholder={placeholder}
                ariaLabel={view === 'mapping' ? 'Filter mapping' : 'Filter profiles'}
            />

            {view === 'mapping' ? (
                <MappingView
                    subDirs={subDirs}
                    associations={associations}
                    filter={filter}
                    view={mapView}
                    sortOrder={mapSort}
                />
            ) : (
                <SyncView
                    settingsData={settingsData}
                    associations={associations}
                    userSelections={userSelections}
                    filter={filter}
                />
            )}

            <div className="mt-6 flex flex-wrap items-center gap-x-4 gap-y-1 text-micro text-ink-3 font-mono">
                <span><Kbd>/</Kbd> filter</span>
                <span><Kbd>Esc</Kbd> clear</span>
                {view === 'mapping' ? <span>drag a character row onto a user file to associate</span> : null}
            </div>

            {confirmDialog}
        </PageShell>
    );
};

export default Profiles;
