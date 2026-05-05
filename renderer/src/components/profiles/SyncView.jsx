// renderer/src/components/profiles/SyncView.jsx
import { useCallback, useEffect, useMemo, useState } from 'react';
import PropTypes from 'prop-types';
import { toast } from 'react-toastify';

import StatusDot from '../ui/StatusDot.jsx';
import SyncProfileRow from './SyncProfileRow.jsx';
import { useConfirmDialog } from '../../hooks/useConfirmDialog.jsx';
import { useSyncAction } from '../../hooks/useSyncAction.js';
import {
    saveUserSelections,
    syncSubdirectory,
    syncAllSubdirectories,
} from '../../api/syncApi';

const SyncView = ({ settingsData, associations, userSelections, filter }) => {
    const { run, isLoading } = useSyncAction();
    const [selections, setSelections] = useState({});
    const [showConfirmDialog, confirmDialog] = useConfirmDialog();
    const [focusedProfile, setFocusedProfile] = useState(null);

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
        let computedNext = null;
        setSelections((prev) => {
            const next = { ...prev, [profile]: { ...prev[profile], [field]: value } };
            if (field === 'charId' && value) {
                const assoc = associations.find((a) => a.charId === value);
                if (assoc) next[profile].userId = assoc.userId;
            }
            computedNext = next;
            return next;
        });
        if (computedNext) {
            persistSelections(computedNext).catch((err) => {
                toast.error(err?.message || 'Failed to save selection.');
            });
        }
    }, [associations, persistSelections]);

    const handleSync = useCallback(async (profile) => {
        const sel = selections[profile] || {};
        if (!sel.charId || !sel.userId) {
            toast.error('Select a character file and a user file first.');
            return;
        }
        const ok = await showConfirmDialog({
            title: 'Sync profile',
            message: `Sync ${profile.replace(/^settings_/, '')} with the chosen pair?`,
        });
        if (!ok.isConfirmed) return;
        await run(() => syncSubdirectory(profile, sel.userId, sel.charId), { errorContext: 'syncSubdirectory' });
    }, [selections, showConfirmDialog, run]);

    const handleSyncAll = useCallback(async (profile) => {
        const sel = selections[profile] || {};
        if (!sel.charId || !sel.userId) {
            toast.error('Select a character file and a user file first.');
            return;
        }
        const ok = await showConfirmDialog({
            title: 'Sync all profiles',
            message: 'Apply the selected pair to every profile? This overwrites their current pair.',
        });
        if (!ok.isConfirmed) return;
        await run(
            () => syncAllSubdirectories(profile, sel.userId, sel.charId),
            { successMessage: 'Sync-All complete', errorContext: 'syncAllSubdirectories' },
        );
    }, [selections, showConfirmDialog, run]);

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

    const empty = !settingsData || settingsData.length === 0;

    if (empty) {
        return (
            <div className="mt-16 text-center">
                <p className="text-body text-ink-2">No EVE profiles found.</p>
                <p className="mt-1 text-meta text-ink-3">
                    Choose a settings directory from the header to scan for profiles.
                </p>
                {confirmDialog}
            </div>
        );
    }

    return (
        <>
            <div role="table" aria-label="Profiles" aria-busy={isLoading} className="rounded-lg border border-rule-1 bg-surface-1 overflow-hidden">
                <div role="row" className="hidden sm:grid grid-cols-[minmax(140px,1fr)_minmax(420px,2.4fr)_72px_36px] gap-3 px-4 py-2 text-meta text-ink-3 border-b border-rule-1 bg-surface-2">
                    <div role="columnheader">Profile</div>
                    <div role="columnheader">Character file → User file</div>
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

            <div className="mt-4 flex flex-wrap items-center gap-x-4 gap-y-1 text-micro text-ink-3 font-mono">
                <span className="inline-flex items-center gap-1.5"><StatusDot state="ready" /> paired</span>
                <span className="inline-flex items-center gap-1.5"><StatusDot state="queued" /> partial</span>
                <span className="inline-flex items-center gap-1.5"><StatusDot state="idle" /> empty</span>
            </div>
            {confirmDialog}
        </>
    );
};

SyncView.propTypes = {
    settingsData: PropTypes.array.isRequired,
    associations: PropTypes.array.isRequired,
    userSelections: PropTypes.object,
    filter: PropTypes.string.isRequired,
};

export default SyncView;
