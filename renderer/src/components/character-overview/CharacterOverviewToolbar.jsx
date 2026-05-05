// src/components/character-overview/CharacterOverviewToolbar.jsx
//
// Subheader actions cluster: group-by segmented control, sort toggle,
// show-hidden toggle, and refresh button. Stateless — receives the
// current values + setters from the page.
import { IconButton, Tooltip } from '@mui/material';
import {
    AccountBalanceOutlined,
    AccountCircleOutlined,
    PlaceOutlined,
    Refresh as RefreshIcon,
    VisibilityOutlined,
    VisibilityOffOutlined,
} from '@mui/icons-material';
import SegmentedControl from '../ui/SegmentedControl.jsx';

const VIEW_OPTIONS = [
    { value: 'account',  label: 'Account',  icon: <AccountBalanceOutlined fontSize="small" /> },
    { value: 'role',     label: 'Role',     icon: <AccountCircleOutlined fontSize="small" /> },
    { value: 'location', label: 'Location', icon: <PlaceOutlined fontSize="small" /> },
];

const CharacterOverviewToolbar = ({
    view, onViewChange,
    sortOrder, onToggleSortOrder,
    showHidden, onToggleShowHidden,
    isRefreshing, onRefreshAll,
    refreshRef,
}) => (
    <>
        <SegmentedControl
            value={view}
            onChange={onViewChange}
            ariaLabel="Group by"
            options={VIEW_OPTIONS}
        />
        <Tooltip title={`Sort ${sortOrder === 'asc' ? 'descending' : 'ascending'}`}>
            <IconButton
                aria-label="Sort"
                onClick={onToggleSortOrder}
                size="small"
                className="h-8! w-8! rounded-md! text-ink-2! hover:bg-surface-2! hover:text-ink-1!"
            >
                <span className="font-mono text-meta tabular">{sortOrder === 'asc' ? 'A→Z' : 'Z→A'}</span>
            </IconButton>
        </Tooltip>
        <Tooltip title={showHidden ? 'Hide hidden accounts' : 'Show hidden accounts'}>
            <IconButton
                aria-label={showHidden ? 'Hide hidden accounts' : 'Show hidden accounts'}
                onClick={onToggleShowHidden}
                size="small"
                className="h-8! w-8! rounded-md! hover:bg-surface-2!"
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
                    aria-busy={isRefreshing}
                    onClick={onRefreshAll}
                    disabled={isRefreshing}
                    size="small"
                    className="h-8! w-8! rounded-md! text-ink-2! hover:bg-surface-2! hover:text-ink-1!"
                >
                    <RefreshIcon
                        fontSize="small"
                        className={isRefreshing ? 'animate-skeleton-pulse' : undefined}
                    />
                </IconButton>
            </span>
        </Tooltip>
    </>
);

export default CharacterOverviewToolbar;
