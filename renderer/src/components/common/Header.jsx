import { useState } from 'react';
import { useLocation, Link } from 'react-router-dom';
import { toast } from 'react-toastify';
import PropTypes from 'prop-types';
import {
    AppBar,
    Toolbar,
    IconButton,
    Box,
    Drawer,
    List,
    ListItem,
    ListItemIcon,
    ListItemText,
    ListItemButton,
    Tooltip,
    CircularProgress,
} from '@mui/material';
import {
    Menu as MenuIcon,
    PersonAddAlt1Outlined,
    LogoutOutlined,
    Close,
    SpaceDashboardOutlined as CharacterOverviewIcon,
    ListAltOutlined as SkillPlansIcon,
    SyncOutlined as SyncIcon,
    AccountTreeOutlined as MappingIcon,
    RefreshOutlined as RefreshIcon,
    SettingsOutlined as SettingsIcon,
    AddchartOutlined,
} from '@mui/icons-material';
import { styled } from '@mui/material/styles';
import { useAuth } from '../../hooks/useAuth';
import { useAppData } from '../../hooks/useAppData';
import { useAsyncOperation } from '../../hooks/useAsyncOperation';
import AccountPromptModal from './AccountPromptModal.jsx';

// Calm AppBar — surface-1 with a single hairline rule, no glass, no shimmer.
const StyledAppBar = styled(AppBar)(() => ({
    background: 'var(--surface-1)',
    backgroundImage: 'none',
    backdropFilter: 'none',
    borderBottom: '1px solid var(--rule-1)',
    boxShadow: 'none',
    color: 'var(--ink-1)',
}));

const StyledDrawer = styled(Drawer)(() => ({
    '& .MuiPaper-root': {
        background: 'var(--surface-1)',
        backgroundImage: 'none',
        backdropFilter: 'none',
        borderRight: '1px solid var(--rule-1)',
        boxShadow: 'none',
        color: 'var(--ink-1)',
        width: 260,
    },
}));

const Header = ({ openSkillPlanModal, existingAccounts }) => {
    const location = useLocation();
    const [drawerOpen, setDrawerOpen] = useState(false);
    const [modalOpen, setModalOpen] = useState(false);

    const { isAuthenticated, logout } = useAuth();
    const { refreshData } = useAppData();
    const { execute: executeRefresh, isLoading: isRefreshing } = useAsyncOperation();
    const { execute: executeAddCharacter } = useAsyncOperation();

    const handleCloseWindow = () => {
        if (window.electronAPI && window.electronAPI.closeWindow) {
            window.electronAPI.closeWindow();
        }
    };

    const toggleDrawer = (open) => () => setDrawerOpen(open);

    const navigationLinks = [
        { text: 'Overview', icon: <CharacterOverviewIcon fontSize="small" />, path: '/' },
        { text: 'Skill Plans', icon: <SkillPlansIcon fontSize="small" />, path: '/skill-plans' },
        { text: 'Mapping', icon: <MappingIcon fontSize="small" />, path: '/mapping' },
        { text: 'Sync', icon: <SyncIcon fontSize="small" />, path: '/sync' },
        { text: 'Settings', icon: <SettingsIcon fontSize="small" />, path: '/settings' },
    ];

    const handleAddCharacterClick = () => setModalOpen(true);
    const handleCloseModal = () => setModalOpen(false);

    const handleAddCharacterSubmit = async (account) => {
        const result = await executeAddCharacter(async () => {
            const { addCharacter } = await import('../../api/apiService');
            return addCharacter(account);
        }, { showToast: false });

        if (result && result.redirectURL) {
            if (window.electronAPI && window.electronAPI.openExternal) {
                window.electronAPI.openExternal(result.redirectURL);
            } else {
                window.open(result.redirectURL, '_blank');
            }
            toast.info('Complete authorization in your browser');

            const state = result.state;
            if (state) {
                let attempts = 0;
                const maxAttempts = 30;
                const pollAuth = setInterval(async () => {
                    attempts++;
                    if (attempts > maxAttempts) {
                        clearInterval(pollAuth);
                        toast.warning('Character add timeout. Please try again.');
                        return;
                    }
                    try {
                        const { finalizelogin } = await import('../../api/apiService');
                        const finalizeResult = await finalizelogin(state);
                        if (finalizeResult && finalizeResult.success) {
                            clearInterval(pollAuth);
                            toast.success('Character added');
                        } else if (finalizeResult && finalizeResult.error === 'state_not_found') {
                            clearInterval(pollAuth);
                            toast.error('Session expired. Please try again.');
                        }
                    } catch (error) {
                        // Polling failure is not fatal — keep retrying until maxAttempts.
                        void error;
                    }
                }, 2000);
            }
        } else {
            toast.success('Character added');
            await refreshData();
        }

        setModalOpen(false);
    };

    const handleRefreshClick = async () => {
        await executeRefresh(() => refreshData(), { showToast: false });
    };

    // Reusable styling for header icon buttons — no scale, no glow.
    const headerIconSx = {
        color: 'var(--ink-2)',
        width: 32,
        height: 32,
        borderRadius: '6px',
        transition: 'background-color var(--motion-duration-fast, 120ms) var(--motion-ease, ease), color var(--motion-duration-fast, 120ms) var(--motion-ease, ease)',
        '&:hover': {
            backgroundColor: 'var(--surface-2)',
            color: 'var(--ink-1)',
            transform: 'none',
        },
    };

    return (
        <>
            <StyledAppBar position="fixed">
                <Toolbar
                    variant="dense"
                    sx={{
                        WebkitAppRegion: 'drag',
                        minHeight: 44,
                        gap: 0.5,
                        px: 1.5,
                    }}
                >
                    {isAuthenticated && (
                        <>
                            <Tooltip title="Menu">
                                <IconButton
                                    edge="start"
                                    aria-label="Menu"
                                    onClick={toggleDrawer(true)}
                                    sx={{ ...headerIconSx, WebkitAppRegion: 'no-drag' }}
                                >
                                    <MenuIcon fontSize="small" />
                                </IconButton>
                            </Tooltip>
                            <Tooltip title="Add character">
                                <IconButton
                                    onClick={handleAddCharacterClick}
                                    aria-label="Add character"
                                    sx={{ ...headerIconSx, WebkitAppRegion: 'no-drag' }}
                                >
                                    <PersonAddAlt1Outlined fontSize="small" />
                                </IconButton>
                            </Tooltip>
                            <Tooltip title="Add skill plan">
                                <IconButton
                                    onClick={openSkillPlanModal}
                                    aria-label="Add skill plan"
                                    sx={{ ...headerIconSx, WebkitAppRegion: 'no-drag' }}
                                >
                                    <AddchartOutlined fontSize="small" />
                                </IconButton>
                            </Tooltip>
                        </>
                    )}

                    <Box sx={{ flex: 1, display: 'flex', justifyContent: 'center', minWidth: 0 }}>
                        <span
                            style={{
                                fontFamily: 'Inter, system-ui, sans-serif',
                                fontWeight: 600,
                                fontSize: 13,
                                letterSpacing: '0.06em',
                                textTransform: 'uppercase',
                                color: 'var(--ink-2)',
                            }}
                        >
                            CanIFly
                        </span>
                    </Box>

                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5, WebkitAppRegion: 'no-drag' }}>
                        {isAuthenticated && (
                            <>
                                <Tooltip title="Refresh data">
                                    <IconButton
                                        onClick={handleRefreshClick}
                                        aria-label="Refresh data"
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
                                        onClick={logout}
                                        aria-label="Log out"
                                        sx={headerIconSx}
                                    >
                                        <LogoutOutlined fontSize="small" />
                                    </IconButton>
                                </Tooltip>
                            </>
                        )}
                        <Tooltip title="Close">
                            <IconButton
                                onClick={handleCloseWindow}
                                aria-label="Close window"
                                sx={headerIconSx}
                            >
                                <Close fontSize="small" />
                            </IconButton>
                        </Tooltip>
                    </Box>
                </Toolbar>
            </StyledAppBar>

            <StyledDrawer anchor="left" open={drawerOpen} onClose={toggleDrawer(false)} disableScrollLock>
                <div
                    role="presentation"
                    onClick={toggleDrawer(false)}
                    onKeyDown={toggleDrawer(false)}
                    style={{ flex: 1, display: 'flex', flexDirection: 'column' }}
                >
                    <Box sx={{ px: 2, pt: 2, pb: 1.5 }}>
                        <span
                            style={{
                                fontSize: 11,
                                fontWeight: 500,
                                letterSpacing: '0.08em',
                                textTransform: 'uppercase',
                                color: 'var(--ink-3)',
                            }}
                        >
                            CanIFly
                        </span>
                    </Box>
                    <Box sx={{ height: '1px', bgcolor: 'var(--rule-1)', mx: 2 }} />
                    <List sx={{ flex: 1, py: 1, px: 1 }}>
                        {navigationLinks.map((item) => {
                            const selected = location.pathname === item.path;
                            return (
                                <ListItem key={item.text} disablePadding sx={{ mb: 0.25 }}>
                                    <ListItemButton
                                        component={Link}
                                        to={item.path}
                                        selected={selected}
                                        sx={{
                                            borderRadius: '6px',
                                            px: 1.5,
                                            py: 1,
                                            transition: 'background-color var(--motion-duration-fast, 120ms) var(--motion-ease, ease)',
                                            color: 'var(--ink-2)',
                                            '& .MuiListItemIcon-root': {
                                                color: 'var(--ink-3)',
                                                minWidth: 28,
                                            },
                                            '&:hover': {
                                                backgroundColor: 'var(--surface-2)',
                                                color: 'var(--ink-1)',
                                                transform: 'none',
                                                '& .MuiListItemIcon-root': {
                                                    color: 'var(--ink-1)',
                                                    transform: 'none',
                                                },
                                            },
                                            '&.Mui-selected': {
                                                backgroundColor: 'var(--accent-soft)',
                                                color: 'var(--ink-1)',
                                                '& .MuiListItemIcon-root': { color: 'var(--accent)' },
                                                '&:hover': { backgroundColor: 'var(--accent-soft)' },
                                            },
                                        }}
                                    >
                                        <ListItemIcon>{item.icon}</ListItemIcon>
                                        <ListItemText
                                            primary={item.text}
                                            primaryTypographyProps={{
                                                sx: {
                                                    fontSize: 14,
                                                    fontWeight: selected ? 500 : 400,
                                                    color: 'inherit',
                                                },
                                            }}
                                        />
                                    </ListItemButton>
                                </ListItem>
                            );
                        })}
                    </List>
                </div>
            </StyledDrawer>

            <AccountPromptModal
                isOpen={modalOpen}
                onClose={handleCloseModal}
                onSubmit={handleAddCharacterSubmit}
                existingAccounts={existingAccounts}
                title="Add Character - Enter Account Name"
            />
        </>
    );
};

Header.propTypes = {
    openSkillPlanModal: PropTypes.func.isRequired,
    existingAccounts: PropTypes.array,
};

export default Header;
