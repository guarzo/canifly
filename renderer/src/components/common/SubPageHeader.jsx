import React, { useState } from 'react';
import PropTypes from 'prop-types';
import { motion, AnimatePresence } from 'framer-motion';
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
        <motion.div 
            className="max-w-7xl mx-auto mb-6"
            initial={{ opacity: 0, y: -20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5 }}
        >
            <Box className="glass rounded-xl p-6 relative overflow-hidden">
                <div className="absolute inset-0 bg-gradient-to-r from-teal-500/5 to-blue-500/5 pointer-events-none" />
                <Box display="flex" alignItems="center" className="relative z-10">
                    <Typography 
                        variant="h4" 
                        className="text-gradient font-display flex-1"
                        sx={{ fontWeight: 700 }}
                    >
                        {title}
                    </Typography>
                    {instructions && (
                        <Tooltip title={showInstructions ? "Hide instructions" : "Show instructions"}>
                            <motion.div
                                whileHover={{ scale: 1.1 }}
                                whileTap={{ scale: 0.9 }}
                            >
                                <IconButton
                                    onClick={toggleInstructions}
                                    className="text-teal-400 hover:text-teal-300 hover:bg-teal-500/10"
                                    size="small"
                                >
                                    {showInstructions ? <HelpFilledIcon fontSize="small" /> : <HelpIcon fontSize="small" />}
                                </IconButton>
                            </motion.div>
                        </Tooltip>
                    )}
                </Box>
                <AnimatePresence>
                    {instructions && showInstructions && (
                        <motion.div
                            initial={{ opacity: 0, height: 0 }}
                            animate={{ opacity: 1, height: 'auto' }}
                            exit={{ opacity: 0, height: 0 }}
                            transition={{ duration: 0.3 }}
                        >
                            <Typography 
                                variant="body2" 
                                className="text-gray-400 mt-3"
                            >
                                {instructions}
                            </Typography>
                        </motion.div>
                    )}
                </AnimatePresence>
            </Box>
        </motion.div>
    );
};

SubPageHeader.propTypes = {
    title: PropTypes.string.isRequired,
    instructions: PropTypes.string,
    storageKey: PropTypes.string.isRequired,
};

export default SubPageHeader;
