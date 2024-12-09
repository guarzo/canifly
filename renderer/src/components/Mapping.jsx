import { useState, useEffect } from 'react';
import PropTypes from 'prop-types';
import { toast } from 'react-toastify';
import Swal from 'sweetalert2';
import withReactContent from 'sweetalert2-react-content';
import 'sweetalert2/dist/sweetalert2.min.css';
import { Grid, CircularProgress, Box } from '@mui/material';
import AccountCard from './MapAccountCard';
import CharacterCard from './MapCharacterCard';

const MySwal = withReactContent(Swal);

const Mapping = ({ associations: initialAssociations, subDirs }) => {
    const [accounts, setAccounts] = useState([]);
    const [availableCharacters, setAvailableCharacters] = useState([]);
    const [associations, setAssociations] = useState(initialAssociations);
    const [isLoading, setIsLoading] = useState(false);
    const [mtimeToColor, setMtimeToColor] = useState({});

    useEffect(() => {
        if (!subDirs.length) return;

        // Deduplicate accounts by userId, choosing the userFile with latest mtime
        const userMap = {};
        subDirs.forEach(mapping => {
            mapping.availableUserFiles.forEach(userFile => {
                const { userId, mtime } = userFile;
                if (!userMap[userId] || new Date(mtime) > new Date(userMap[userId].mtime)) {
                    userMap[userId] = userFile;
                }
            });
        });

        // Convert userMap to array and sort by descending mtime
        const uniqueAccounts = Object.values(userMap).sort((a, b) => new Date(b.mtime) - new Date(a.mtime));
        setAccounts(uniqueAccounts);

        // Deduplicate chars by charId, choosing the charFile with latest mtime
        const charMap = {};
        subDirs.forEach(mapping => {
            mapping.availableCharFiles.forEach(charFile => {
                const { charId, mtime } = charFile;
                if (!charMap[charId] || new Date(mtime) > new Date(charMap[charId].mtime)) {
                    charMap[charId] = { ...charFile, subDir: mapping.subDir };
                }
            });
        });

        // Filter out associated chars and sort by descending mtime
        const associatedCharIds = new Set(associations.map(a => a.charId));
        const uniqueChars = Object.values(charMap)
            .filter(ch => !associatedCharIds.has(ch.charId))
            .sort((a, b) => new Date(b.mtime) - new Date(a.mtime));
        setAvailableCharacters(uniqueChars);

        assignColors(uniqueAccounts, uniqueChars);
    }, [subDirs, associations]);

    const assignColors = (uniqueAccounts, uniqueChars) => {
        const predefinedColors = ['#4caf50', '#f44336', '#ff9800', '#9c27b0', '#00bcd4', '#e91e63'];

        // Extract mtimes from both accounts and chars
        const accountMtimes = uniqueAccounts.map(a => a.mtime);
        const charMtimes = uniqueChars.map(c => c.mtime);
        const allMtimes = [...accountMtimes, ...charMtimes];

        const uniqueMtimes = Array.from(new Set(allMtimes));

        // Sort mtimes so that color assignment is stable
        uniqueMtimes.sort((a, b) => new Date(a) - new Date(b));

        const colorMapping = uniqueMtimes.reduce((acc, mtime, index) => {
            acc[mtime] = predefinedColors[index % predefinedColors.length];
            return acc;
        }, {});

        setMtimeToColor(colorMapping);
    };

    const handleDragStart = (event, charId) => {
        event.dataTransfer.setData('text/plain', charId);
    };

    const handleDrop = async (event, userId) => {
        event.preventDefault();
        const charId = event.dataTransfer.getData('text/plain');
        const char = availableCharacters.find(c => c.charId === charId);

        if (!char) {
            toast.error('Character not found.');
            return;
        }

        const confirmAssoc = await MySwal.fire({
            title: 'Confirm Association',
            text: `Associate "${char.name}" (${charId}) with account "${userId}"?`,
            icon: 'question',
            showCancelButton: true,
            confirmButtonText: 'Yes, associate it!',
            cancelButtonText: 'Cancel',
        });

        if (!confirmAssoc.isConfirmed) return;

        try {
            setIsLoading(true);
            const response = await fetch('/api/associate-character', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                credentials: 'include',
                body: JSON.stringify({ userId, charId })
            });
            const result = await response.json();
            if (response.ok && result.success) {
                toast.success(result.message);

                // Update local state:
                // Remove the character from availableCharacters
                setAvailableCharacters(prev => prev.filter(ch => ch.charId !== charId));

                // Add a new association
                const charName = char.name;
                const newAssoc = { userId, charId, charName, mtime: char.mtime };
                setAssociations(prev => [...prev, newAssoc]);
            } else {
                toast.error(`Association failed: ${result.message}`);
            }
        } catch (error) {
            console.error('Error associating character:', error);
            toast.error('Association operation failed.');
        } finally {
            setIsLoading(false);
        }
    };

    const handleUnassociate = async (userId, charId, charName) => {
        const confirmUnassoc = await MySwal.fire({
            title: 'Confirm Unassociation',
            text: `Unassociate "${charName}" (${charId}) from account "${userId}"?`,
            icon: 'warning',
            showCancelButton: true,
            confirmButtonText: 'Yes, unassociate it!',
            cancelButtonText: 'Cancel',
        });

        if (!confirmUnassoc.isConfirmed) return;

        try {
            setIsLoading(true);
            const response = await fetch('/api/unassociate-character', {
                method: 'POST',
                headers: {'Content-Type': 'application/json'},
                credentials: 'include',
                body: JSON.stringify({ userId, charId })
            });
            const result = await response.json();
            if (response.ok && result.success) {
                toast.success(result.message);

                // Update local state:
                setAssociations(prev => prev.filter(a => a.charId !== charId || a.userId !== userId));
            } else {
                toast.error(`Unassociation failed: ${result.message}`);
            }
        } catch (error) {
            console.error('Error unassociating character:', error);
            toast.error('Unassociation operation failed.');
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <div className="pt-16">
            <Box sx={{ paddingTop: '64px', padding: 2, minHeight: '100vh' }}>
                {isLoading ? (
                    <Box display="flex" justifyContent="center" alignItems="center" minHeight="50vh">
                        <CircularProgress color="primary" />
                    </Box>
                ) : (
                    <Grid container spacing={4}>
                        <Grid item xs={12} md={6}>
                            {accounts.map(mapping => (
                                <AccountCard
                                    key={`${mapping.userId}-${mapping.mtime}`}
                                    mapping={mapping}
                                    associations={associations}
                                    handleUnassociate={handleUnassociate}
                                    handleDrop={handleDrop}
                                    mtimeToColor={mtimeToColor}
                                />
                            ))}
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
                                    <Box textAlign="center" width="100%">
                                        No available characters to associate.
                                    </Box>
                                )}
                            </Grid>
                        </Grid>
                    </Grid>
                )}
            </Box>
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
            subDir: PropTypes.string.isRequired,
            availableCharFiles: PropTypes.arrayOf(
                PropTypes.shape({
                    file: PropTypes.string.isRequired,
                    charId: PropTypes.string.isRequired,
                    name: PropTypes.string.isRequired,
                    mtime: PropTypes.string.isRequired,
                })
            ).isRequired,
            availableUserFiles: PropTypes.arrayOf(
                PropTypes.shape({
                    file: PropTypes.string.isRequired,
                    userId: PropTypes.string.isRequired,
                    name: PropTypes.string.isRequired,
                    mtime: PropTypes.string.isRequired,
                })
            ).isRequired,
        })
    ).isRequired,
};

export default Mapping;
