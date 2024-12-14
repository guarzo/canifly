import React, { useState } from 'react';
import PropTypes from 'prop-types';
import { Typography, Box, Tooltip, IconButton } from '@mui/material';
import { Help as HelpFilledIcon, HelpOutline as HelpIcon } from '@mui/icons-material';

/**
 * A reusable page header that displays a title, optional instructions, and a toggle button
 * to show/hide those instructions.
 *
 * @param {string} title - The title text for the page
 * @param {string} instructions - The instructions text to display beneath the title when shown
 * @param {string} storageKey - A unique key to store the instructions visibility state in localStorage
 */
const SubPageHeader = ({ title, instructions, storageKey }) => {
    const [showInstructions, setShowInstructions] = useState(() => {
        const stored = localStorage.getItem(storageKey);
        return stored === null ? true : JSON.parse(stored);
    });

    const toggleInstructions = () => {
        const newValue = !showInstructions;
        setShowInstructions(newValue);
        localStorage.setItem(storageKey, JSON.stringify(newValue));
    };

    return (
        <Box className="max-w-7xl mx-auto mb-6">
            <Box className="bg-gradient-to-r from-gray-900 to-gray-800 p-4 rounded-md shadow-md relative">
                <Box display="flex" alignItems="center">
                    <Typography variant="h4" sx={{ color: '#14b8a6', fontWeight: 'bold', marginBottom: '0.5rem', flex: 1 }}>
                        {title}
                    </Typography>
                    {instructions && (
                        <Tooltip title={showInstructions ? "Hide instructions" : "Show instructions"}>
                            <IconButton
                                onClick={toggleInstructions}
                                sx={{ color: '#99f6e4' }}
                                size="small"
                            >
                                {showInstructions ? <HelpFilledIcon fontSize="small" /> : <HelpIcon fontSize="small" />}
                            </IconButton>
                        </Tooltip>
                    )}
                </Box>
                {instructions && showInstructions && (
                    <Typography variant="body2" sx={{ color: '#99f6e4', marginTop: '0.5rem' }}>
                        {instructions}
                    </Typography>
                )}
            </Box>
        </Box>
    );
};

SubPageHeader.propTypes = {
    title: PropTypes.string.isRequired,
    instructions: PropTypes.string,
    storageKey: PropTypes.string.isRequired,
};

export default SubPageHeader;
