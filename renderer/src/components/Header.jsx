// Header.jsx

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
} from '@mui/icons-material';
import { styled } from '@mui/material/styles';
import { toast } from 'react-toastify';
import AccountPromptModal from './AccountPromptModal';

const StyledAppBar = styled(AppBar)(() => ({
    backgroundImage: 'linear-gradient(to right, #1f2937, #1f2937)',
    color: '#14b8a6',
    boxShadow: 'inset 0 -4px 0 0 #14b8a6',
    borderBottom: '4px solid #14b8a6',
}));

const Header = ({ loggedIn, handleLogout, openSkillPlanModal }) => {
    const location = useLocation();
    const [drawerOpen, setDrawerOpen] = useState(false);
    const [modalOpen, setModalOpen] = useState(false);

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
        // Instead of direct redirect, open the modal to get account name
        setModalOpen(true);
    };

    const handleCloseModal = () => {
        setModalOpen(false);
    };

    const handleAddCharacterSubmit = async (account) => {
        try {
            const response = await fetch('/api/add-character', {
                method: 'POST',
                headers: {'Content-Type': 'application/json'},
                body: JSON.stringify({account}),
                credentials: 'include',
            })
                if (response.ok) {
                const data = await response.json();
                if (data.redirectURL) {
                    window.location.href = data.redirectURL;
                } else {
                    toast.error("No redirect URL received from server.");
                }
            } else {
                const errorText = await response.text();
                toast.error(`Failed to initiate login: ${errorText}`);
            }
        } catch (error) {
            console.error('Error initiating add character:', error);
            toast.error('An error occurred while adding character.');
        } finally {
            setModalOpen(false);
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
                                {/* Add Character Icon */}
                                <IconButton onClick={handleAddCharacterClick} title="Add Character">
                                    <AddCircleOutline sx={{ color: '#22c55e' }} />
                                </IconButton>
                                {/* Add Skill Plan Icon */}
                                <IconButton onClick={openSkillPlanModal} title="Add Skill Plan">
                                    <SkillPlansIcon sx={{ color: '#f59e0b' }} />
                                </IconButton>
                                {/* Logout Icon */}
                                <IconButton onClick={handleLogout} title="Logout">
                                    <ExitToApp sx={{ color: '#ef4444' }} />
                                </IconButton>
                            </>
                        )}
                        <IconButton onClick={handleCloseWindow} title="Close">
                            <Close sx={{ color: '#9ca3af' }} />
                        </IconButton>
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

            {/* Account Prompt Modal for Add Character */}
            <AccountPromptModal
                isOpen={modalOpen}
                onClose={handleCloseModal}
                onSubmit={handleAddCharacterSubmit}
                title="Add Character - Enter Account Name"
            />
        </>
    );
};

Header.propTypes = {
    loggedIn: PropTypes.bool.isRequired,
    handleLogout: PropTypes.func.isRequired,
    openSkillPlanModal: PropTypes.func.isRequired,
};

export default Header;
