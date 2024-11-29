// src/components/Header.jsx

import React from 'react';
import PropTypes from 'prop-types';

/**
 * Header Component
 *
 * Displays the application title and user actions like "Add Character," "Add Skill Plan," and "Logout."
 *
 * Props:
 * - title: The title of the application.
 * - loggedIn: Boolean indicating if the user is logged in.
 * - handleLogout: Function to handle user logout.
 */
const Header = ({ title, loggedIn, handleLogout }) => {
    const handleAddCharacter = () => {
        window.location.href = 'http://localhost:8713/auth-character'; // Initiate OAuth flow
    };

    const handleAddSkillPlan = () => {
        // Open the Add Skill Plan Modal by dispatching a custom event
        window.dispatchEvent(new Event('openAddSkillPlanModal'));
    };

    return (
        <header className="fixed top-0 left-0 right-0 bg-white dark:bg-gray-800 shadow-md z-10">
            <div className="container mx-auto px-4 py-4 flex justify-between items-center">
                <h1 className="text-xl font-bold text-gray-800 dark:text-white">{title}</h1>
                <div className="flex items-center space-x-4">
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
                            <button
                                onClick={handleLogout}
                                className="flex items-center p-2 bg-red-500 text-white rounded hover:bg-red-600 focus:outline-none"
                                title="Logout"
                            >
                                <i className="fas fa-sign-out-alt"></i>
                            </button>
                        </>
                    )}
                </div>
            </div>
        </header>
    );
};

Header.propTypes = {
    title: PropTypes.string.isRequired,
    loggedIn: PropTypes.bool.isRequired,
    handleLogout: PropTypes.func.isRequired,
};

export default Header;
