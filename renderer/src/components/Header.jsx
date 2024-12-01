import React from 'react';
import PropTypes from 'prop-types';
import AddCircleOutlineIcon from '@mui/icons-material/AddCircleOutline';
import ExitToAppIcon from '@mui/icons-material/ExitToApp';
import CloseIcon from '@mui/icons-material/Close';
const { ipcRenderer } = require('electron');

const Header = ({ loggedIn, handleLogout }) => {
    const handleAddCharacter = () => {
        window.location.href = 'http://localhost:8713/auth-character';
    };

    const handleCloseWindow = () => {
        ipcRenderer.send('close-window');
    };

    return (
        <header
            className="fixed top-0 left-0 right-0 bg-gradient-to-r from-gray-900 to-gray-800 shadow-md z-10 py-3"
            style={{ WebkitAppRegion: 'drag' }}
        >
            <div className="container mx-auto px-4 flex justify-between items-center">
                <div className="flex space-x-3" style={{ WebkitAppRegion: 'no-drag' }}>
                    {loggedIn && (
                        <button
                            onClick={handleAddCharacter}
                            className="p-2 text-green-400 hover:text-green-600"
                            title="Add Character"
                        >
                            <AddCircleOutlineIcon />
                        </button>
                    )}
                </div>

                <h1 className="text-2xl font-bold text-teal-200 text-center flex-grow">Can I Fly?</h1>

                <div className="flex space-x-3" style={{ WebkitAppRegion: 'no-drag' }}>
                    {loggedIn && (
                        <>
                            <button
                                onClick={handleLogout}
                                className="p-2 text-red-400 hover:text-red-600"
                                title="Logout"
                            >
                                <ExitToAppIcon />
                            </button>
                            <button
                                onClick={handleCloseWindow}
                                className="p-2 text-gray-400 hover:text-gray-600"
                                title="Close"
                            >
                                <CloseIcon />
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
