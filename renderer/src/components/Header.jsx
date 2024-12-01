// Header.jsx
import React from 'react';
import PropTypes from 'prop-types';
import { Link, useLocation } from 'react-router-dom';
import { AppBar, Toolbar, IconButton, Typography } from '@mui/material';
import {
    AddCircleOutline,
    ExitToApp,
    Close,
    Dashboard as DashboardIcon,
    ListAlt as SkillPlansIcon,
    Settings as SettingsIcon, // Import SettingsIcon
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

    const handleAddCharacter = () => {
        window.location.href = 'http://localhost:8713/auth-character';
    };

    const handleCloseWindow = () => {
        ipcRenderer.send('close-window');
    };

    return (
        <StyledAppBar position="fixed">
            <Toolbar style={{ WebkitAppRegion: 'drag' }}>
                {/* Left Side Icons */}
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
                                <SettingsIcon sx={{ color: '#f59e0b' }} /> {/* Amber color */}
                            </IconButton>
                        </>
                    )}
                </div>

                {/* Title */}
                <Typography
                    variant="h6"
                    className="flex-grow text-center"
                    sx={{ color: '#14b8a6' }}
                >
                    Can I Fly?
                </Typography>

                {/* Navigation Icons */}
                {loggedIn && (
                    <nav
                        className="mr-4 flex items-center space-x-2"
                        style={{ WebkitAppRegion: 'no-drag' }}
                    >
                        <IconButton
                            component={Link}
                            to="/"
                            title="Dashboard"
                            sx={{
                                color: location.pathname === '/' ? '#ffffff' : '#14b8a6',
                                '&:hover': { color: '#ffffff' },
                            }}
                        >
                            <DashboardIcon />
                        </IconButton>
                        <IconButton
                            component={Link}
                            to="/skill-plans"
                            title="Skill Plans"
                            sx={{
                                color:
                                    location.pathname === '/skill-plans' ? '#ffffff' : '#14b8a6',
                                '&:hover': { color: '#ffffff' },
                            }}
                        >
                            <SkillPlansIcon />
                        </IconButton>
                    </nav>
                )}

                {/* Right Side Icons */}
                <div
                    className="flex items-center space-x-3"
                    style={{ WebkitAppRegion: 'no-drag' }}
                >
                    {loggedIn && (
                        <IconButton onClick={handleLogout} title="Logout">
                            <ExitToApp sx={{ color: '#ef4444' }} /> {/* Red color */}
                        </IconButton>
                    )}
                    <IconButton onClick={handleCloseWindow} title="Close">
                        <Close sx={{ color: '#9ca3af' }} /> {/* Gray color */}
                    </IconButton>
                </div>
            </Toolbar>
        </StyledAppBar>
    );
};

Header.propTypes = {
    loggedIn: PropTypes.bool.isRequired,
    handleLogout: PropTypes.func.isRequired,
    openSkillPlanModal: PropTypes.func.isRequired,
};

export default Header;
