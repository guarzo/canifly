// src/components/common/HeaderUserMenu.jsx
//
// Right-hand cluster of the Header — Refresh, Logout, and Close Window.
import { Box, IconButton, Tooltip, CircularProgress } from '@mui/material';
import {
    LogoutOutlined,
    Close,
    RefreshOutlined as RefreshIcon,
} from '@mui/icons-material';
import { headerIconSx } from './HeaderToolbarActions.jsx';

const HeaderUserMenu = ({ isAuthenticated, isRefreshing, onRefresh, onLogout, onClose }) => (
    <Box
        style={{ WebkitAppRegion: 'no-drag' }}
        sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}
    >
        {isAuthenticated && (
            <>
                <Tooltip title="Refresh data">
                    <IconButton
                        onClick={onRefresh}
                        aria-label="Refresh data"
                        disabled={isRefreshing}
                        aria-busy={isRefreshing}
                        style={{ WebkitAppRegion: 'no-drag' }}
                        sx={headerIconSx}
                    >
                        {isRefreshing ? (
                            <CircularProgress size={16} sx={{ color: 'var(--ink-2)' }} />
                        ) : (
                            <RefreshIcon fontSize="small" />
                        )}
                    </IconButton>
                </Tooltip>
                <Tooltip title="Log out">
                    <IconButton
                        onClick={onLogout}
                        aria-label="Log out"
                        style={{ WebkitAppRegion: 'no-drag' }}
                        sx={headerIconSx}
                    >
                        <LogoutOutlined fontSize="small" />
                    </IconButton>
                </Tooltip>
            </>
        )}
        <Tooltip title="Close">
            <IconButton
                onClick={onClose}
                aria-label="Close window"
                style={{ WebkitAppRegion: 'no-drag' }}
                sx={headerIconSx}
            >
                <Close fontSize="small" />
            </IconButton>
        </Tooltip>
    </Box>
);

export default HeaderUserMenu;
