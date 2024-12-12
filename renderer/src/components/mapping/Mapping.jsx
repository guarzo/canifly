import { useState, useEffect } from 'react';
import PropTypes from 'prop-types';
import { toast } from 'react-toastify';
import { Grid, Box, Typography } from '@mui/material';
import AccountCard from './MapAccountCard.jsx';
import CharacterCard from './MapCharacterCard.jsx';
import { useConfirmDialog } from '../partials/useConfirmDialog.jsx'; // Import the custom hook

function roundToSeconds(mtime) {
    const date = new Date(mtime);
    // Round down to the nearest second
    date.setMilliseconds(0);
    return date.toISOString();
}

const Mapping = ({ associations: initialAssociations, subDirs, onRefreshData, backEndURL }) => {
    const [accounts, setAccounts] = useState([]);
    const [availableCharacters, setAvailableCharacters] = useState([]);
    const [associations, setAssociations] = useState(initialAssociations);
    const [mtimeToColor, setMtimeToColor] = useState({});

    const [showConfirmDialog, confirmDialog] = useConfirmDialog();

    // Process data whenever subDirs or associations change
    useEffect(() => {
        if (subDirs.length === 0) return;

        // Deduplicate accounts by userId, choosing latest mtime
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

        // Deduplicate chars by charId with rounding
        const charMap = {};
        subDirs.forEach(mapping => {
            mapping.availableCharFiles.forEach(charFile => {
                const roundedMtime = roundToSeconds(charFile.mtime);
                const { charId } = charFile;
                if (!charMap[charId] || new Date(roundedMtime) > new Date(charMap[charId].mtime)) {
                    charMap[charId] = { ...charFile, mtime: roundedMtime, subDir: mapping.subDir };
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

        try {
            const response = await fetch(`${backEndURL}/api/associate-character`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                credentials: 'include',
                body: JSON.stringify({ userId, charId, userName, charName })
            });
            const result = await response.json();
            if (response.ok && result.success) {
                toast.success(result.message);
                // Optimistic update
                setAvailableCharacters(prev => prev.filter(ch => ch.charId !== charId));
                const newAssoc = { userId, charId, charName, mtime: char.mtime };
                setAssociations(prev => [...prev, newAssoc]);

                if (onRefreshData) {
                    await onRefreshData();
                }
            } else {
                toast.error(`Association failed: ${result.message}`);
            }
        } catch (error) {
            console.error('Error associating character:', error);
            toast.error('Association operation failed.');
        }
    };

    const handleUnassociate = async (userId, charId, charName, userName) => {
        const confirmUnassoc = await showConfirmDialog({
            title: 'Confirm Unassociation',
            message: `Unassociate "${charName}" from account "${userName}"?`
        });

        if (!confirmUnassoc.isConfirmed) return;

        try {
            const response = await fetch(`${backEndURL}/api/unassociate-character`, {
                method: 'POST',
                headers: {'Content-Type': 'application/json'},
                credentials: 'include',
                body: JSON.stringify({ userId, charId, userName, charName })
            });
            const result = await response.json();
            if (response.ok && result.success) {
                toast.success(result.message);
                setAssociations(prev => prev.filter(a => a.charId !== charId || a.userId !== userId));
                if (onRefreshData) {
                    await onRefreshData();
                }
            } else {
                toast.error(`Unassociation failed: ${result.message}`);
            }
        } catch (error) {
            console.error('Error unassociating character:', error);
            toast.error('Unassociation operation failed.');
        }
    };

    return (
        <div className="bg-gray-900 min-h-screen text-teal-200 px-4 pb-10 pt-16">
            {/* Heading and subtle instructions */}
            <Box className="max-w-7xl mx-auto mb-6">
                <Box className="bg-gradient-to-r from-gray-900 to-gray-800 p-4 rounded-md shadow-md">
                    <Typography variant="h4" sx={{ color: '#14b8a6', fontWeight: 'bold', marginBottom: '0.5rem' }}>
                        Mapping
                    </Typography>
                    <Typography variant="body2" sx={{ color: '#99f6e4' }}>
                        Drag and drop characters onto accounts to associate them. Associated characters appear under their accounts.
                    </Typography>
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
                        <Grid container spacing={2}>
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
