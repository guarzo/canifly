// src/components/mapping/Mapping.jsx
import React, { useState, useEffect } from 'react';
import PropTypes from 'prop-types';
import { toast } from 'react-toastify';
import { Grid, Box, Typography, IconButton, Tooltip } from '@mui/material';
import { HelpOutline as HelpIcon, Help as HelpFilledIcon } from '@mui/icons-material';
import AccountCard from '../components/mapping/MapAccountCard.jsx';
import CharacterCard from '../components/mapping/MapCharacterCard.jsx';
import { useConfirmDialog } from '../hooks/useConfirmDialog.jsx';
import { associateCharacter, unassociateCharacter } from '../api/apiService.jsx';
import { mappingInstructions } from './../utils/instructions'; // Imported from a separate file

function roundToSeconds(mtime) {
    const date = new Date(mtime);
    date.setMilliseconds(0);
    return date.toISOString();
}

const Mapping = ({ associations: initialAssociations, subDirs, onRefreshData, backEndURL }) => {
    const [accounts, setAccounts] = useState([]);
    const [availableCharacters, setAvailableCharacters] = useState([]);
    const [associations, setAssociations] = useState(initialAssociations);
    const [mtimeToColor, setMtimeToColor] = useState({});
    const [showConfirmDialog, confirmDialog] = useConfirmDialog();

    // Load instruction visibility from localStorage
    const [showInstructions, setShowInstructions] = useState(() => {
        const stored = localStorage.getItem('showMappingInstructions');
        return stored === null ? true : JSON.parse(stored);
    });

    useEffect(() => {
        if (subDirs.length === 0) return;

        const userMap = {};
        subDirs.forEach(mapping => {
            mapping.availableUserFiles.forEach(userFile => {
                const roundedMtime = roundToSeconds(userFile.mtime);
                if (!userMap[userFile.userId] || new Date(roundedMtime) > new Date(userMap[userFile.userId].mtime)) {
                    userMap[userFile.userId] = { ...userFile, mtime: roundedMtime };
                }
            });
        });

        const uniqueAccounts = Object.values(userMap)
            .sort((a, b) => new Date(b.mtime) - new Date(a.mtime));

        const charMap = {};
        subDirs.forEach(mapping => {
            mapping.availableCharFiles.forEach(charFile => {
                const roundedMtime = roundToSeconds(charFile.mtime);
                const { charId } = charFile;
                if (!charMap[charId] || new Date(roundedMtime) > new Date(charMap[charId].mtime)) {
                    charMap[charId] = { ...charFile, mtime: roundedMtime, profile: mapping.profile };
                }
            });
        });

        const associatedCharIds = new Set(associations.map(a => a.charId));
        const uniqueChars = Object.values(charMap)
            .filter(ch => !associatedCharIds.has(ch.charId))
            .sort((a, b) => new Date(b.mtime) - new Date(a.mtime));

        setAccounts(uniqueAccounts);
        setAvailableCharacters(uniqueChars);
        assignColors(uniqueAccounts, uniqueChars);
    }, [subDirs, associations]);

    const assignColors = (uniqueAccounts, uniqueChars) => {
        const predefinedColors = ['#4caf50', '#f44336', '#ff9800', '#9c27b0', '#00bcd4', '#e91e63'];

        const accountMtimes = uniqueAccounts.map(a => a.mtime);
        const charMtimes = uniqueChars.map(c => c.mtime);
        const allMtimes = [...accountMtimes, ...charMtimes];

        const uniqueMtimes = Array.from(new Set(allMtimes)).sort((a, b) => new Date(a) - new Date(b));

        const colorMapping = uniqueMtimes.reduce((acc, mtime, index) => {
            acc[mtime] = predefinedColors[index % predefinedColors.length];
            return acc;
        }, {});

        setMtimeToColor(colorMapping);
    };

    const handleDragStart = (event, charId) => {
        event.dataTransfer.setData('text/plain', charId);
    };

    const handleDrop = async (event, userId, userName) => {
        event.preventDefault();
        const charId = event.dataTransfer.getData('text/plain');
        const char = availableCharacters.find(c => c.charId === charId);
        const charName = char?.name;

        if (!char) {
            toast.error('Character not found.');
            return;
        }

        const confirmAssoc = await showConfirmDialog({
            title: 'Confirm Association',
            message: `Associate "${charName}" with account "${userName}"?`
        });

        if (!confirmAssoc.isConfirmed) return;

        const result = await associateCharacter(userId, charId, userName, charName, backEndURL);
        if (result && result.success) {
            toast.success(result.message);
            setAvailableCharacters(prev => prev.filter(ch => ch.charId !== charId));
            const newAssoc = { userId, charId, charName, mtime: char.mtime };
            setAssociations(prev => [...prev, newAssoc]);

            if (onRefreshData) {
                await onRefreshData();
            }
        }
    };

    const handleUnassociate = async (userId, charId, charName, userName) => {
        const confirmUnassoc = await showConfirmDialog({
            title: 'Confirm Unassociation',
            message: `Unassociate "${charName}" from account "${userName}"?`
        });

        if (!confirmUnassoc.isConfirmed) return;

        const result = await unassociateCharacter(userId, charId, userName, charName, backEndURL);
        if (result && result.success) {
            toast.success(result.message);
            setAssociations(prev => prev.filter(a => a.charId !== charId || a.userId !== userId));
            if (onRefreshData) {
                await onRefreshData();
            }
        }
    };

    const toggleInstructions = () => {
        const newValue = !showInstructions;
        setShowInstructions(newValue);
        localStorage.setItem('showMappingInstructions', JSON.stringify(newValue));
    };

    return (
        <div className="bg-gray-900 min-h-screen text-teal-200 px-4 pb-10 pt-16">
            <Box className="max-w-7xl mx-auto mb-6">
                <Box className="bg-gradient-to-r from-gray-900 to-gray-800 p-4 rounded-md shadow-md relative">
                    <Box display="flex" alignItems="center">
                        <Typography variant="h4" sx={{ color: '#14b8a6', fontWeight: 'bold', marginBottom: '0.5rem', flex: 1 }}>
                            Map Characters to User Files
                        </Typography>
                        <Tooltip title={showInstructions ? "Hide instructions" : "Show instructions"}>
                            <IconButton
                                onClick={toggleInstructions}
                                sx={{ color: '#99f6e4' }}
                                size="small"
                            >
                                {showInstructions ? <HelpFilledIcon fontSize="small" /> : <HelpIcon fontSize="small" />}
                            </IconButton>
                        </Tooltip>
                    </Box>
                    {showInstructions && (
                        <Typography variant="body2" sx={{ color: '#99f6e4', marginTop: '0.5rem' }}>
                            {mappingInstructions}
                        </Typography>
                    )}
                </Box>
            </Box>

            <Box className="max-w-7xl mx-auto">
                <Grid container spacing={4}>
                    <Grid item xs={12} md={6}>
                        {accounts.length === 0 ? (
                            <Box textAlign="center" className="text-gray-300">
                                No accounts found.
                            </Box>
                        ) : (
                            accounts.map(mapping => (
                                <AccountCard
                                    key={`${mapping.userId}-${mapping.mtime}`}
                                    mapping={mapping}
                                    associations={associations}
                                    handleUnassociate={handleUnassociate}
                                    handleDrop={handleDrop}
                                    mtimeToColor={mtimeToColor}
                                />
                            ))
                        )}
                    </Grid>
                    <Grid item xs={12} md={6}>
                        <Grid container spacing={2} data-testid="available-characters">
                            {availableCharacters.length ? (
                                availableCharacters.map(char => (
                                    <Grid item xs={12} sm={6} key={`${char.charId}-${char.mtime}`}>
                                        <CharacterCard
                                            char={char}
                                            handleDragStart={handleDragStart}
                                            mtimeToColor={mtimeToColor}
                                        />
                                    </Grid>
                                ))
                            ) : (
                                <Box textAlign="center" width="100%" className="text-gray-300">
                                    No available characters to associate.
                                </Box>
                            )}
                        </Grid>
                    </Grid>
                </Grid>
            </Box>

            {confirmDialog}
        </div>
    );
};

Mapping.propTypes = {
    associations: PropTypes.arrayOf(
        PropTypes.shape({
            userId: PropTypes.string.isRequired,
            charId: PropTypes.string.isRequired,
            charName: PropTypes.string.isRequired,
            mtime: PropTypes.string,
        })
    ).isRequired,
    subDirs: PropTypes.arrayOf(
        PropTypes.shape({
            profile: PropTypes.string.isRequired,
            availableCharFiles: PropTypes.array.isRequired,
            availableUserFiles: PropTypes.array.isRequired,
        })
    ).isRequired,
    onRefreshData: PropTypes.func,
    backEndURL: PropTypes.string.isRequired,
};

export default Mapping;
