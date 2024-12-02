// Header.jsx
import React, { useState } from 'react';
import PropTypes from 'prop-types';
import { Link, useLocation } from 'react-router-dom';
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
} from '@mui/material';
import {
    Menu as MenuIcon,
    AddCircleOutline,
    ExitToApp,
    Close,
    Dashboard as DashboardIcon,
    ListAlt as SkillPlansIcon,
    Map as LocationIcon,
    People as RoleIcon,
} from '@mui/icons-material';
import { styled } from '@mui/material/styles';
const { ipcRenderer } = require('electron');

// Custom styled AppBar
const StyledAppBar = styled(AppBar)(({ theme }) => ({
    backgroundImage: 'linear-gradient(to right, #1f2937, #1f2937)',
    color: '#14b8a6',
    boxShadow: 'inset 0 -4px 0 0 #14b8a6',
    borderBottom: '4px solid #14b8a6',
}));

const Header = ({ loggedIn, handleLogout, openSkillPlanModal }) => {
    const location = useLocation();
    const [drawerOpen, setDrawerOpen] = useState(false);

    const handleAddCharacter = () => {
        window.location.href = 'http://localhost:8713/auth-character';
    };

    const handleCloseWindow = () => {
        ipcRenderer.send('close-window');
    };

    const toggleDrawer = (open) => () => {
        setDrawerOpen(open);
    };

    const navigationLinks = [
        { text: 'Dashboard', icon: <DashboardIcon />, path: '/' },
        { text: 'Skill Plans', icon: <SkillPlansIcon />, path: '/skill-plans' },
        { text: 'By Location', icon: <LocationIcon />, path: '/characters-by-location' },
        { text: 'By Role', icon: <RoleIcon />, path: '/characters-by-role' },
    ];

    return (
        <>
            <StyledAppBar position="fixed">
                <Toolbar style={{ WebkitAppRegion: 'drag' }}>
                    {/* Menu Icon */}
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

                    {/* Title */}
                    <Typography
                        variant="h6"
                        className="flex-grow text-center"
                        sx={{ color: '#14b8a6' }}
                    >
                        Can I Fly?
                    </Typography>

                    {/* Right Side Icons */}
                    <div
                        className="flex items-center space-x-3"
                        style={{ WebkitAppRegion: 'no-drag' }}
                    >
                        {loggedIn && (
                            <>
                                {/* Add Character Icon */}
                                <IconButton onClick={handleAddCharacter} title="Add Character">
                                    <AddCircleOutline sx={{ color: '#22c55e' }} /> {/* Green color */}
                                </IconButton>
                                {/* Add Skill Plan Icon */}
                                <IconButton onClick={openSkillPlanModal} title="Add Skill Plan">
                                    <SkillPlansIcon sx={{ color: '#f59e0b' }} /> {/* Amber color */}
                                </IconButton>
                                {/* Logout Icon */}
                                <IconButton onClick={handleLogout} title="Logout">
                                    <ExitToApp sx={{ color: '#ef4444' }} /> {/* Red color */}
                                </IconButton>
                            </>
                        )}
                        <IconButton onClick={handleCloseWindow} title="Close">
                            <Close sx={{ color: '#9ca3af' }} /> {/* Gray color */}
                        </IconButton>
                    </div>
                </Toolbar>
            </StyledAppBar>

            {/* Navigation Drawer */}
            <Drawer anchor="left" open={drawerOpen} onClose={toggleDrawer(false)}>
                <div
                    role="presentation"
                    onClick={toggleDrawer(false)}
                    onKeyDown={toggleDrawer(false)}
                    style={{ width: 250 }}
                >
                    <List>
                        {navigationLinks.map((item) => (
                            <ListItem
                                button
                                key={item.text}
                                component={Link}
                                to={item.path}
                                selected={location.pathname === item.path}
                            >
                                <ListItemIcon sx={{ color: '#14b8a6' }}>{item.icon}</ListItemIcon>
                                <ListItemText primary={item.text} />
                            </ListItem>
                        ))}
                    </List>
                    <Divider />
                </div>
            </Drawer>
        </>
    );
};

Header.propTypes = {
    loggedIn: PropTypes.bool.isRequired,
    handleLogout: PropTypes.func.isRequired,
    openSkillPlanModal: PropTypes.func.isRequired,
};

export default Header;
