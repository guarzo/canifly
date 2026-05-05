// renderer/src/components/profiles/SyncProfileRow.jsx
import { useMemo } from 'react';
import PropTypes from 'prop-types';
import { SyncOutlined, SyncAltOutlined } from '@mui/icons-material';
import { IconButton, Tooltip } from '@mui/material';
import PairSelect from '../ui/PairSelect.jsx';

const SyncProfileRow = ({
    subDir,
    selection,
    onSelectionChange,
    onSync,
    onSyncAll,
    isLoading,
    isFocused,
    onFocus,
}) => {
    const profileName = subDir.profile.replace(/^settings_/, '');
    const charId = selection?.charId || '';
    const userId = selection?.userId || '';
    const ready = Boolean(charId && userId);

    const charOptions = useMemo(
        () => [...(subDir.availableCharFiles || [])]
            .sort((a, b) => (a.name || a.charId).localeCompare(b.name || b.charId))
            .map((c) => ({ value: c.charId, primary: c.name || c.charId, secondary: c.charId })),
        [subDir.availableCharFiles],
    );
    const userOptions = useMemo(
        () => [...(subDir.availableUserFiles || [])]
            .sort((a, b) => a.userId.localeCompare(b.userId))
            .map((u) => ({ value: u.userId, primary: u.name || u.userId, secondary: u.name ? u.userId : undefined })),
        [subDir.availableUserFiles],
    );

    const status = ready ? 'ready' : (charId || userId) ? 'queued' : 'idle';

    return (
        <div
            role="row"
            tabIndex={0}
            onFocus={onFocus}
            className={[
                'grid items-center gap-3 px-4 h-12',
                'grid-cols-[minmax(140px,1fr)_minmax(420px,2.4fr)_72px_36px]',
                'border-b border-rule-1 last:border-b-0',
                'transition-colors duration-fast ease-out-quart outline-hidden',
                isFocused ? 'bg-surface-2 shadow-rail-accent' : 'hover:bg-surface-2',
            ].join(' ')}
            data-profile={subDir.profile}
        >
            <div className="min-w-0">
                <div className="text-body text-ink-1 truncate font-mono tabular" title={subDir.profile}>
                    {profileName}
                </div>
                <div className="text-micro text-ink-3 tabular">
                    {charOptions.length} char · {userOptions.length} user
                </div>
            </div>

            <PairSelect
                state={status}
                statusLabel={ready ? 'Pair selected' : 'No pair'}
                leftValue={charId}
                leftOptions={charOptions}
                onLeftChange={(v) => onSelectionChange(subDir.profile, 'charId', v)}
                leftLabel={`Character file for ${profileName}`}
                leftPlaceholder="Select character file…"
                rightValue={userId}
                rightOptions={userOptions}
                onRightChange={(v) => onSelectionChange(subDir.profile, 'userId', v)}
                rightLabel={`User file for ${profileName}`}
                rightPlaceholder="Select user file…"
                arrowActive={ready}
            />

            <div className="flex items-center justify-end">
                <Tooltip title={ready ? 'Sync this profile' : 'Select a character and user file first'}>
                    <span>
                        <button
                            type="button"
                            aria-label={`Sync ${profileName}`}
                            onClick={() => onSync(subDir.profile)}
                            disabled={!ready || isLoading}
                            className={[
                                'inline-flex items-center justify-center gap-1.5',
                                'h-8 px-3 rounded-md text-meta',
                                'transition-colors duration-fast ease-out-quart',
                                ready
                                    ? 'bg-accent text-accent-ink hover:bg-accent-strong'
                                    : 'bg-surface-2 text-ink-4 cursor-not-allowed',
                            ].join(' ')}
                        >
                            <SyncOutlined sx={{ fontSize: 16 }} />
                            <span>Sync</span>
                        </button>
                    </span>
                </Tooltip>
            </div>

            <div className="flex items-center justify-center">
                <Tooltip title={ready ? 'Sync all profiles using this pair' : 'Select a pair to enable Sync All'}>
                    <span>
                        <IconButton
                            size="small"
                            aria-label={`Sync all profiles using ${profileName} selection`}
                            onClick={() => onSyncAll(subDir.profile)}
                            disabled={!ready || isLoading}
                            className="h-8! w-8! rounded-md!"
                            sx={{
                                color: ready ? 'var(--ink-2)' : 'var(--ink-4)',
                                '&:hover': { backgroundColor: 'var(--surface-3)', color: 'var(--ink-1)' },
                            }}
                        >
                            <SyncAltOutlined sx={{ fontSize: 18 }} />
                        </IconButton>
                    </span>
                </Tooltip>
            </div>
        </div>
    );
};

SyncProfileRow.propTypes = {
    subDir: PropTypes.shape({
        profile: PropTypes.string.isRequired,
        availableCharFiles: PropTypes.array.isRequired,
        availableUserFiles: PropTypes.array.isRequired,
    }).isRequired,
    selection: PropTypes.shape({
        charId: PropTypes.string,
        userId: PropTypes.string,
    }),
    onSelectionChange: PropTypes.func.isRequired,
    onSync: PropTypes.func.isRequired,
    onSyncAll: PropTypes.func.isRequired,
    isLoading: PropTypes.bool.isRequired,
    isFocused: PropTypes.bool,
    onFocus: PropTypes.func,
};

export default SyncProfileRow;
