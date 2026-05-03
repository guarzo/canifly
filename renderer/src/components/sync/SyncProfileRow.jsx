// src/components/sync/SyncProfileRow.jsx
//
// Dense, table-style row for one EVE profile in the Sync page.
// Each row pairs a character file with a user file and exposes a
// prominent (but calm) Sync action plus a Sync-All overflow.

import { useMemo } from 'react';
import PropTypes from 'prop-types';
import {
    ArrowForwardOutlined,
    SyncOutlined,
    SyncAltOutlined,
} from '@mui/icons-material';
import { IconButton, Tooltip, Select, MenuItem } from '@mui/material';

import StatusDot from '../ui/StatusDot.jsx';

const selectSx = {
    height: 32,
    minHeight: 32,
    fontSize: 14,
    color: 'var(--ink-1)',
    backgroundColor: 'var(--surface-1)',
    borderRadius: '6px',
    '& .MuiOutlinedInput-notchedOutline': { borderColor: 'var(--rule-1)' },
    '&:hover .MuiOutlinedInput-notchedOutline': { borderColor: 'var(--rule-2)' },
    '&.Mui-focused .MuiOutlinedInput-notchedOutline': {
        borderColor: 'var(--accent)',
        borderWidth: '1px',
    },
    '& .MuiSelect-select': { paddingTop: 6, paddingBottom: 6 },
};

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

    const sortedChars = useMemo(
        () => [...(subDir.availableCharFiles || [])].sort((a, b) => a.name.localeCompare(b.name)),
        [subDir.availableCharFiles],
    );
    const sortedUsers = useMemo(
        () => [...(subDir.availableUserFiles || [])].sort((a, b) => a.userId.localeCompare(b.userId)),
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
                'grid-cols-[20px_minmax(140px,1fr)_minmax(180px,1.4fr)_20px_minmax(180px,1.4fr)_72px_36px]',
                'border-b border-rule-1 last:border-b-0',
                'transition-colors duration-fast ease-out-quart outline-hidden',
                isFocused ? 'bg-surface-2 shadow-rail-accent' : 'hover:bg-surface-2',
            ].join(' ')}
            data-profile={subDir.profile}
        >
            <div className="flex items-center justify-center">
                <StatusDot state={status} label={ready ? 'Pair selected' : 'No pair'} />
            </div>

            <div className="min-w-0">
                <div className="text-body text-ink-1 truncate font-mono tabular" title={subDir.profile}>
                    {profileName}
                </div>
                <div className="text-micro text-ink-3 tabular">
                    {sortedChars.length} char · {sortedUsers.length} user
                </div>
            </div>

            <Select
                size="small"
                value={charId}
                onChange={(e) => onSelectionChange(subDir.profile, 'charId', e.target.value)}
                displayEmpty
                inputProps={{ 'aria-label': `Character file for ${profileName}` }}
                sx={selectSx}
                fullWidth
            >
                <MenuItem value="">
                    <span className="text-ink-3">Select character file…</span>
                </MenuItem>
                {sortedChars.map((c) => (
                    <MenuItem key={c.charId} value={c.charId}>
                        <span className="font-mono text-body text-ink-1">{c.name}</span>
                        <span className="ml-2 font-mono text-meta text-ink-3 tabular">{c.charId}</span>
                    </MenuItem>
                ))}
            </Select>

            <div className="flex items-center justify-center" aria-hidden="true">
                <ArrowForwardOutlined
                    sx={{
                        fontSize: 16,
                        color: ready ? 'var(--accent)' : 'var(--ink-4)',
                    }}
                />
            </div>

            <Select
                size="small"
                value={userId}
                onChange={(e) => onSelectionChange(subDir.profile, 'userId', e.target.value)}
                displayEmpty
                inputProps={{ 'aria-label': `User file for ${profileName}` }}
                sx={selectSx}
                fullWidth
            >
                <MenuItem value="">
                    <span className="text-ink-3">Select user file…</span>
                </MenuItem>
                {sortedUsers.map((u) => (
                    <MenuItem key={u.userId} value={u.userId}>
                        <span className="font-mono text-body text-ink-1 tabular">{u.userId}</span>
                    </MenuItem>
                ))}
            </Select>

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
