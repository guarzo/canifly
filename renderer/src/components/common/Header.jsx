import { useState } from 'react';
import PropTypes from 'prop-types';
import { useLocation, Link } from 'react-router-dom';
import {
    AppBar,
    Toolbar,
    IconButton,
    Typography,
    Box,
    Drawer,
    List,
    ListItem,
    ListItemIcon,
    ListItemText,
    ListItemButton,
    Tooltip,
    CircularProgress,
    Divider
} from '@mui/material';
import {
    Menu as MenuIcon,
    AddCircleOutline,
    ExitToApp,
    Close,
    Dashboard as CharacterOverviewIcon,
    ListAlt as SkillPlansIcon,
    Sync as SyncIcon,
    AccountTree as MappingIcon,
    Cached as RefreshIcon
} from '@mui/icons-material';
import { styled } from '@mui/material/styles';
import AccountPromptModal from './AccountPromptModal.jsx';
import nav_img1 from '../../assets/images/nav-logo.png';
import nav_img2 from '../../assets/images/nav-logo2.webp';

const StyledAppBar = styled(AppBar)(() => ({
    backgroundImage: 'linear-gradient(to right, #1f2937, #1f2937)',
    color: '#14b8a6',
    boxShadow: 'inset 0 -4px 0 0 #14b8a6',
    borderBottom: '4px solid #14b8a6',
}));

const StyledDrawer = styled(Drawer)(() => ({
    '& .MuiPaper-root': {
        background: 'linear-gradient(to bottom, #1f2937, #111827)',
        overflow: 'hidden',
        color: '#5eead4',
        display: 'flex',
        flexDirection: 'column',
        justifyContent: 'space-between',
        width: 250 // Ensure the drawer has a fixed width
    }
}));

const Header = ({ loggedIn, handleLogout, openSkillPlanModal, existingAccounts, onSilentRefresh, onAddCharacter, isRefreshing }) => {
    const location = useLocation();
    const [drawerOpen, setDrawerOpen] = useState(false);
    const [modalOpen, setModalOpen] = useState(false);
    // State to toggle between two images each time the nav is opened
    const [useAlternateImage, setUseAlternateImage] = useState(false);

    const handleCloseWindow = () => {
        if (window.electronAPI && window.electronAPI.closeWindow) {
            window.electronAPI.closeWindow();
        } else {
            console.error('Electron API not available');
        }
    };

    const toggleDrawer = (open) => () => {
        setDrawerOpen(open);
        if (open === true) {
            setUseAlternateImage((prev) => !prev);
        }
    };

    const navigationLinks = [
        { text: 'Overview', icon: <CharacterOverviewIcon />, path: '/' },
        { text: 'Skill Plans', icon: <SkillPlansIcon />, path: '/skill-plans' },
        { text: 'Mapping', icon: <MappingIcon />, path: '/mapping' },
        { text: 'Sync', icon: <SyncIcon />, path: '/sync' },
    ];

    const handleAddCharacterClick = () => {
        setModalOpen(true);
    };

    const handleCloseModal = () => {
        setModalOpen(false);
    };

    const handleAddCharacterSubmit = async (account) => {
        await onAddCharacter(account);
        setModalOpen(false);
    };

    const handleRefreshClick = async () => {
        if (!onSilentRefresh) return;
        await onSilentRefresh();
    };

    const chosenImage = useAlternateImage ? nav_img2 : nav_img1;

    return (
        <>
            <StyledAppBar position="fixed">
                <Toolbar style={{ WebkitAppRegion: 'drag', display: 'flex', alignItems: 'center' }}>
                    {loggedIn && (
                        <>
                            <IconButton
                                edge="start"
                                color="inherit"
                                aria-label="menu"
                                onClick={toggleDrawer(true)}
                                style={{ WebkitAppRegion: 'no-drag' }}
                            >
                                <MenuIcon />
                            </IconButton>
                            <Tooltip title="Add Character">
                                <IconButton onClick={handleAddCharacterClick} style={{ WebkitAppRegion: 'no-drag' }}>
                                    <AddCircleOutline sx={{ color: '#22c55e' }} />
                                </IconButton>
                            </Tooltip>
                            <Tooltip title="Add Skill Plan">
                                <IconButton onClick={openSkillPlanModal} style={{ WebkitAppRegion: 'no-drag' }}>
                                    <SkillPlansIcon sx={{ color: '#f59e0b' }} />
                                </IconButton>
                            </Tooltip>
                        </>
                    )}

                    <Box sx={{ flex: 1, display: 'flex', justifyContent: 'center' }}>
                        <Typography
                            variant="h6"
                            sx={{ color: '#14b8a6', textAlign: 'center' }}
                        >
                            Can I Fly?
                        </Typography>
                    </Box>

                    <Box sx={{ display: 'flex', alignItems: 'center', WebkitAppRegion: 'no-drag' }}>
                        {loggedIn && (
                            <>
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
                    </Box>
                </Toolbar>
            </StyledAppBar>

            <StyledDrawer anchor="left" open={drawerOpen} onClose={toggleDrawer(false)} disableScrollLock>
                <div
                    role="presentation"
                    onClick={toggleDrawer(false)}
                    onKeyDown={toggleDrawer(false)}
                    style={{ flex: 1, display: 'flex', flexDirection: 'column', justifyContent: 'space-between' }}
                >
                    <List sx={{ flex: 1 }}>
                        {navigationLinks.map((item, index) => (
                            <div key={item.text}>
                                <ListItem disablePadding>
                                    <ListItemButton
                                        component={Link}
                                        to={item.path}
                                        selected={location.pathname === item.path}
                                        sx={{
                                            '&:hover': {
                                                backgroundColor: '#0f172a',
                                                '& .MuiListItemText-primary': {
                                                    color: '#a7f3d0',
                                                },
                                                '& .MuiListItemIcon-root': {
                                                    color: '#a7f3d0',
                                                },
                                            },
                                            '&.Mui-selected': {
                                                backgroundColor: '#134e4a',
                                                '&:hover': {
                                                    backgroundColor: '#145a54',
                                                    '& .MuiListItemText-primary': {
                                                        color: '#a7f3d0',
                                                    },
                                                    '& .MuiListItemIcon-root': {
                                                        color: '#a7f3d0',
                                                    },
                                                },
                                            },
                                        }}
                                    >
                                        <ListItemIcon sx={{ color: '#5eead4' }}>{item.icon}</ListItemIcon>
                                        <ListItemText
                                            primary={item.text}
                                            primaryTypographyProps={{ sx: { color: '#5eead4' } }}
                                        />
                                    </ListItemButton>
                                </ListItem>

                                {index === 1 && (
                                    <Divider
                                        sx={{
                                            backgroundColor: '#14b8a6',
                                            marginY: '4px',
                                            opacity: 0.5
                                        }}
                                    />
                                )}
                            </div>
                        ))}
                    </List>
                    {/* Image at the bottom of the nav drawer */}
                    <Box sx={{ p: 2, textAlign: 'center' }}>
                        <img
                            src={chosenImage}
                            alt="Nav Logo"
                            style={{
                                maxWidth: '220px',
                                height: 'auto',
                                display: 'block',
                                margin: '0 auto'
                            }}
                        />
                    </Box>
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
    loggedIn: PropTypes.bool.isRequired,
    handleLogout: PropTypes.func.isRequired,
    openSkillPlanModal: PropTypes.func.isRequired,
    existingAccounts: PropTypes.array.isRequired,
    onSilentRefresh: PropTypes.func,
    onAddCharacter: PropTypes.func.isRequired,
    isRefreshing: PropTypes.bool.isRequired
};

export default Header;
