import { useState } from 'react';
import PropTypes from 'prop-types';
import { useLocation, Link } from 'react-router-dom';
import {
    AppBar,
    Toolbar,
    IconButton,
    Typography,
    Drawer,
    List,
    ListItem,
    ListItemIcon,
    ListItemText,
    Divider,
    ListItemButton,
    Tooltip,
    CircularProgress
} from '@mui/material';
import {
    Menu as MenuIcon,
    AddCircleOutline,
    ExitToApp,
    Close,
    Dashboard as DashboardIcon,
    ListAlt as SkillPlansIcon,
    People as RoleIcon,
    Sync as SyncIcon,
    AccountTree as MappingIcon,
    Cached as RefreshIcon
} from '@mui/icons-material';
import { styled } from '@mui/material/styles';
import AccountPromptModal from './AccountPromptModal';

const StyledAppBar = styled(AppBar)(() => ({
    backgroundImage: 'linear-gradient(to right, #1f2937, #1f2937)',
    color: '#14b8a6',
    boxShadow: 'inset 0 -4px 0 0 #14b8a6',
    borderBottom: '4px solid #14b8a6',
}));

const Header = ({ loggedIn, handleLogout, openSkillPlanModal, existingAccounts, onSilentRefresh, onAddCharacter }) => {
    const location = useLocation();
    const [drawerOpen, setDrawerOpen] = useState(false);
    const [modalOpen, setModalOpen] = useState(false);
    const [isRefreshing, setIsRefreshing] = useState(false);

    const handleCloseWindow = () => {
        if (window.electronAPI && window.electronAPI.closeWindow) {
            window.electronAPI.closeWindow();
        } else {
            console.error('Electron API not available');
        }
    };

    const toggleDrawer = (open) => () => {
        setDrawerOpen(open);
    };

    const navigationLinks = [
        { text: 'Dashboard', icon: <DashboardIcon />, path: '/' },
        { text: 'Skill Plans', icon: <SkillPlansIcon />, path: '/skill-plans' },
        { text: 'Character Sort', icon: <RoleIcon />, path: '/character-sort' },
        { text: 'Sync', icon: <SyncIcon />, path: '/sync' },
        { text: 'Mapping', icon: <MappingIcon />, path: '/mapping' },
    ];

    const handleAddCharacterClick = () => {
        setModalOpen(true);
    };

    const handleCloseModal = () => {
        setModalOpen(false);
    };

    // Updated to call onAddCharacter from props
    const handleAddCharacterSubmit = async (account) => {
        // Just call the prop from App
        await onAddCharacter(account);
        setModalOpen(false);
    };

    const handleRefreshClick = async () => {
        if (!onSilentRefresh) return;
        setIsRefreshing(true);
        try {
            await onSilentRefresh();
            console.log('Data refreshed!');
        } catch (err) {
            console.error('Error during silent refresh:', err);
        } finally {
            setIsRefreshing(false);
        }
    };

    return (
        <>
            <StyledAppBar position="fixed">
                <Toolbar style={{ WebkitAppRegion: 'drag' }}>
                    {loggedIn && (
                        <IconButton
                            edge="start"
                            color="inherit"
                            aria-label="menu"
                            onClick={toggleDrawer(true)}
                            style={{ WebkitAppRegion: 'no-drag' }}
                        >
                            <MenuIcon />
                        </IconButton>
                    )}

                    <Typography
                        variant="h6"
                        className="flex-grow text-center"
                        sx={{ color: '#14b8a6' }}
                    >
                        Can I Fly?
                    </Typography>

                    <div
                        className="flex items-center space-x-3"
                        style={{ WebkitAppRegion: 'no-drag' }}
                    >
                        {loggedIn && (
                            <>
                                <Tooltip title="Add Character">
                                    <IconButton onClick={handleAddCharacterClick}>
                                        <AddCircleOutline sx={{ color: '#22c55e' }} />
                                    </IconButton>
                                </Tooltip>
                                <Tooltip title="Add Skill Plan">
                                    <IconButton onClick={openSkillPlanModal}>
                                        <SkillPlansIcon sx={{ color: '#f59e0b' }} />
                                    </IconButton>
                                </Tooltip>
                                <Tooltip title="Refresh Data">
                                    <IconButton onClick={handleRefreshClick}>
                                        {isRefreshing ? (
                                            <CircularProgress size={24} sx={{ color: '#9ca3af' }} />
                                        ) : (
                                            <RefreshIcon sx={{ color: '#9ca3af' }} />
                                        )}
                                    </IconButton>
                                </Tooltip>
                                <Tooltip title="Logout">
                                    <IconButton onClick={handleLogout}>
                                        <ExitToApp sx={{ color: '#ef4444' }} />
                                    </IconButton>
                                </Tooltip>
                            </>
                        )}
                        <Tooltip title="Close">
                            <IconButton onClick={handleCloseWindow}>
                                <Close sx={{ color: '#9ca3af' }} />
                            </IconButton>
                        </Tooltip>
                    </div>
                </Toolbar>
            </StyledAppBar>

            <Drawer anchor="left" open={drawerOpen} onClose={toggleDrawer(false)}>
                <div
                    role="presentation"
                    onClick={toggleDrawer(false)}
                    onKeyDown={toggleDrawer(false)}
                    style={{ width: 250 }}
                >
                    <List>
                        {navigationLinks.map((item) => (
                            <ListItem key={item.text} disablePadding>
                                <ListItemButton
                                    component={Link}
                                    to={item.path}
                                    selected={location.pathname === item.path}
                                >
                                    <ListItemIcon sx={{ color: '#14b8a6' }}>{item.icon}</ListItemIcon>
                                    <ListItemText primary={item.text} />
                                </ListItemButton>
                            </ListItem>
                        ))}
                    </List>
                    <Divider />
                </div>
            </Drawer>

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
    loggedIn: PropTypes.bool.isRequired,
    handleLogout: PropTypes.func.isRequired,
    openSkillPlanModal: PropTypes.func.isRequired,
    existingAccounts: PropTypes.array.isRequired,
    onSilentRefresh: PropTypes.func,
    onAddCharacter: PropTypes.func.isRequired, // New prop
};

export default Header;
