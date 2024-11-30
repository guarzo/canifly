//src/components/Header.jsx

import React from 'react';
import PropTypes from 'prop-types';
const { ipcRenderer } = require('electron'); // Import ipcRenderer for communication with main process

/**
 * Header Component
 *
 * Displays the application title and user actions like "Add Character," "Add Skill Plan," "Logout," and "Close Window."
 */
const Header = ({ loggedIn, handleLogout }) => {
    const handleAddCharacter = () => {
        window.location.href = 'http://localhost:8713/auth-character'; // Initiate OAuth flow
    };

    const handleAddSkillPlan = () => {
        window.dispatchEvent(new Event('openAddSkillPlanModal'));
    };

    const handleCloseWindow = () => {
        ipcRenderer.send('close-window'); // Send close-window event to main process
    };

    return (
        <header
            className="fixed top-0 left-0 right-0 bg-gradient-to-r from-gray-900 to-gray-800 shadow-md z-10 py-3"
            style={{ WebkitAppRegion: 'drag' }} // Enable window dragging
        >
            <div className="container mx-auto px-4 flex justify-between items-center">
                <div className="flex space-x-2" style={{ WebkitAppRegion: 'no-drag' }}>
                    {loggedIn && (
                        <>
                            <button
                                onClick={handleAddCharacter}
                                className="flex items-center p-2 bg-green-500 text-white rounded hover:bg-green-600 focus:outline-none"
                                title="Add Character"
                            >
                                <i className="fas fa-user-plus"></i>
                            </button>
                            <button
                                onClick={handleAddSkillPlan}
                                className="flex items-center p-2 bg-blue-500 text-white rounded hover:bg-blue-600 focus:outline-none"
                                title="Add Skill Plan"
                            >
                                <i className="fas fa-plus-circle"></i>
                            </button>
                        </>
                    )}
                </div>

                <h1 className="text-2xl font-bold text-teal-200 text-center flex-grow">Can I Fly?</h1>

                <div className="flex space-x-2" style={{ WebkitAppRegion: 'no-drag' }}>
                    {loggedIn && (
                        <>
                            <button
                                onClick={handleLogout}
                                className="flex items-center p-2 bg-red-500 text-white rounded hover:bg-red-600 focus:outline-none"
                                title="Logout"
                            >
                                <i className="fas fa-sign-out-alt"></i>
                            </button>
                            <button
                                onClick={handleCloseWindow}
                                className="flex items-center p-2 bg-gray-500 text-white rounded hover:bg-gray-600 focus:outline-none"
                                title="Close"
                            >
                                <i className="fas fa-times"></i>
                            </button>
                        </>
                    )}
                </div>
            </div>
        </header>
    );
};

Header.propTypes = {
    loggedIn: PropTypes.bool.isRequired,
    handleLogout: PropTypes.func.isRequired,
};

export default Header;
