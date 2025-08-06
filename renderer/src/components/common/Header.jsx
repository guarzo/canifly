import { useState } from 'react';
import { useLocation, Link } from 'react-router-dom';
import { toast } from 'react-toastify';
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
    Cached as RefreshIcon,
    Settings as SettingsIcon
} from '@mui/icons-material';
import { styled } from '@mui/material/styles';
import { useAuth } from '../../hooks/useAuth';
import { useAppData } from '../../hooks/useAppData';
import { useAsyncOperation } from '../../hooks/useAsyncOperation';
import AccountPromptModal from './AccountPromptModal.jsx';
import nav_img1 from '../../assets/images/nav-logo.png';
import nav_img2 from '../../assets/images/nav-logo2.webp';

const StyledAppBar = styled(AppBar)(() => ({
    background: 'rgba(17, 24, 39, 0.8)',
    backdropFilter: 'blur(20px)',
    borderBottom: '1px solid rgba(20, 184, 166, 0.2)',
    boxShadow: '0 4px 30px rgba(0, 0, 0, 0.1)',
    '&::after': {
        content: '""',
        position: 'absolute',
        bottom: 0,
        left: 0,
        right: 0,
        height: '2px',
        background: 'linear-gradient(90deg, transparent, #14b8a6, transparent)',
        animation: 'shimmer 3s linear infinite',
    }
}));

const StyledDrawer = styled(Drawer)(() => ({
    '& .MuiPaper-root': {
        background: 'rgba(17, 24, 39, 0.95)',
        backdropFilter: 'blur(20px)',
        borderRight: '1px solid rgba(148, 163, 184, 0.1)',
        boxShadow: '0 0 40px rgba(0, 0, 0, 0.5)',
        color: '#5eead4',
        display: 'flex',
        flexDirection: 'column',
        justifyContent: 'space-between',
        width: 280,
        overflow: 'hidden',
        '&::before': {
            content: '""',
            position: 'absolute',
            top: 0,
            left: 0,
            right: 0,
            bottom: 0,
            background: 'radial-gradient(circle at top right, rgba(20, 184, 166, 0.1), transparent 60%)',
            pointerEvents: 'none',
        }
    }
}));

const Header = ({ openSkillPlanModal, existingAccounts }) => {
    const location = useLocation();
    const [drawerOpen, setDrawerOpen] = useState(false);
    const [modalOpen, setModalOpen] = useState(false);
    const [useAlternateImage, setUseAlternateImage] = useState(false);
    
    const { isAuthenticated, logout } = useAuth();
    const { refreshData } = useAppData();
    const { execute: executeRefresh, isLoading: isRefreshing } = useAsyncOperation();
    const { execute: executeAddCharacter } = useAsyncOperation();

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
        { text: 'Settings', icon: <SettingsIcon />, path: '/settings' },
    ];

    const handleAddCharacterClick = () => {
        setModalOpen(true);
    };

    const handleCloseModal = () => {
        setModalOpen(false);
    };

    const handleAddCharacterSubmit = async (account) => {
        const result = await executeAddCharacter(async () => {
            const { addCharacter } = await import('../../api/apiService');
            return addCharacter(account);
        }, {
            showToast: false // We'll handle the success/error message manually
        });
        
        // If we got a redirect URL, open it in the browser
        if (result && result.redirectURL) {
            console.log('Opening SSO page:', result.redirectURL);
            window.open(result.redirectURL, '_blank');
            toast.info('Please complete the authorization in your browser');
        } else {
            toast.success('Character added successfully');
        }
        
        setModalOpen(false);
    };

    const handleRefreshClick = async () => {
        await executeRefresh(() => refreshData(), {
            showToast: false
        });
    };

    const chosenImage = useAlternateImage ? nav_img2 : nav_img1;

    return (
        <>
            <StyledAppBar position="fixed">
                <Toolbar style={{ WebkitAppRegion: 'drag', display: 'flex', alignItems: 'center' }}>
                    {isAuthenticated && (
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
                                <IconButton 
                                    onClick={handleAddCharacterClick} 
                                    className="group"
                                    style={{ WebkitAppRegion: 'no-drag' }}
                                    sx={{
                                        transition: 'all 0.3s',
                                        '&:hover': {
                                            backgroundColor: 'rgba(34, 197, 94, 0.1)',
                                            transform: 'scale(1.1)',
                                        }
                                    }}
                                >
                                    <AddCircleOutline 
                                        className="transition-all duration-300 group-hover:rotate-90"
                                        sx={{ color: '#22c55e' }} 
                                    />
                                </IconButton>
                            </Tooltip>
                            <Tooltip title="Add Skill Plan">
                                <IconButton 
                                    onClick={openSkillPlanModal} 
                                    className="group"
                                    style={{ WebkitAppRegion: 'no-drag' }}
                                    sx={{
                                        transition: 'all 0.3s',
                                        '&:hover': {
                                            backgroundColor: 'rgba(245, 158, 11, 0.1)',
                                            transform: 'scale(1.1)',
                                        }
                                    }}
                                >
                                    <SkillPlansIcon 
                                        className="transition-all duration-300 group-hover:rotate-12"
                                        sx={{ color: '#f59e0b' }} 
                                    />
                                </IconButton>
                            </Tooltip>
                        </>
                    )}

                    <Box sx={{ flex: 1, display: 'flex', justifyContent: 'center' }}>
                        <Typography
                            variant="h6"
                            className="text-gradient font-display"
                            sx={{ 
                                textAlign: 'center',
                                fontSize: '1.5rem',
                                fontWeight: 700,
                                letterSpacing: '0.05em'
                            }}
                        >
                            CanIFly
                        </Typography>
                    </Box>

                    <Box sx={{ display: 'flex', alignItems: 'center', WebkitAppRegion: 'no-drag' }}>
                        {isAuthenticated && (
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
                                    <IconButton onClick={logout}>
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
                                            borderRadius: '12px',
                                            margin: '4px 8px',
                                            transition: 'all 0.3s cubic-bezier(0.4, 0, 0.2, 1)',
                                            position: 'relative',
                                            overflow: 'hidden',
                                            '&::before': {
                                                content: '""',
                                                position: 'absolute',
                                                top: 0,
                                                left: '-100%',
                                                width: '100%',
                                                height: '100%',
                                                background: 'linear-gradient(90deg, transparent, rgba(20, 184, 166, 0.1), transparent)',
                                                transition: 'left 0.5s',
                                            },
                                            '&:hover': {
                                                backgroundColor: 'rgba(20, 184, 166, 0.1)',
                                                transform: 'translateX(4px)',
                                                '&::before': {
                                                    left: '100%',
                                                },
                                                '& .MuiListItemText-primary': {
                                                    color: '#5eead4',
                                                },
                                                '& .MuiListItemIcon-root': {
                                                    color: '#5eead4',
                                                    transform: 'scale(1.1)',
                                                },
                                            },
                                            '&.Mui-selected': {
                                                backgroundColor: 'rgba(20, 184, 166, 0.2)',
                                                borderLeft: '3px solid #14b8a6',
                                                '&:hover': {
                                                    backgroundColor: 'rgba(20, 184, 166, 0.3)',
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

export default Header;
