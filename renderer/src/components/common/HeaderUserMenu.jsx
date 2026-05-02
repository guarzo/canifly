// src/components/common/HeaderUserMenu.jsx
//
// Right-hand cluster of the Header — Refresh, Logout, and Close Window.
import PropTypes from 'prop-types';
import { Box, IconButton, Tooltip, CircularProgress } from '@mui/material';
import {
    LogoutOutlined,
    Close,
    RefreshOutlined as RefreshIcon,
} from '@mui/icons-material';
import { headerIconSx } from './HeaderToolbarActions.jsx';

const HeaderUserMenu = ({ isAuthenticated, isRefreshing, onRefresh, onLogout, onClose }) => (
    <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5, WebkitAppRegion: 'no-drag' }}>
        {isAuthenticated && (
            <>
                <Tooltip title="Refresh data">
                    <IconButton onClick={onRefresh} aria-label="Refresh data" sx={headerIconSx}>
                        {isRefreshing ? (
                            <CircularProgress size={16} sx={{ color: 'var(--ink-2)' }} />
                        ) : (
                            <RefreshIcon fontSize="small" />
                        )}
                    </IconButton>
                </Tooltip>
                <Tooltip title="Log out">
                    <IconButton onClick={onLogout} aria-label="Log out" sx={headerIconSx}>
                        <LogoutOutlined fontSize="small" />
                    </IconButton>
                </Tooltip>
            </>
        )}
        <Tooltip title="Close">
            <IconButton onClick={onClose} aria-label="Close window" sx={headerIconSx}>
                <Close fontSize="small" />
            </IconButton>
        </Tooltip>
    </Box>
);

HeaderUserMenu.propTypes = {
    isAuthenticated: PropTypes.bool.isRequired,
    isRefreshing: PropTypes.bool.isRequired,
    onRefresh: PropTypes.func.isRequired,
    onLogout: PropTypes.func.isRequired,
    onClose: PropTypes.func.isRequired,
};

export default HeaderUserMenu;
