import { useState, useEffect } from 'react';
import { toast } from 'react-toastify';
import Swal from 'sweetalert2';
import withReactContent from 'sweetalert2-react-content';
import 'sweetalert2/dist/sweetalert2.min.css';
import { Grid, CircularProgress, Box } from '@mui/material';
import AccountCard from './MapAccountCard';
import CharacterCard from './MapCharacterCard';

const MySwal = withReactContent(Swal);

const Mapping = ({ associations, subDirs, onRefreshData }) => {
    const [mappings, setMappings] = useState([]); // Unique Accounts
    const [availableCharacters, setAvailableCharacters] = useState([]);
    const [isLoading, setIsLoading] = useState(false);
    const [mtimeToColor, setMtimeToColor] = useState({});

    useEffect(() => {
        // Whenever subDirs or associations change, we recalculate mappings and available chars
        if (!subDirs.length) return;

        // Step 1: Map each subDir to its latest user based on mtime
        const subDirToLatestUserId = {};
        subDirs.forEach(mapping => {
            mapping.availableUserFiles.forEach(userFile => {
                const { userId, mtime, name } = userFile;
                if (
                    !subDirToLatestUserId[mapping.subDir] ||
                    new Date(userFile.mtime) > new Date(subDirToLatestUserId[mapping.subDir].mtime)
                ) {
                    subDirToLatestUserId[mapping.subDir] = { userId, mtime, name };
                }
            });
        });

        // Step 2: Create a list of unique accounts based on userId
        const uniqueAccountsMap = {};
        Object.values(subDirToLatestUserId).forEach(({ userId, mtime, name }) => {
            if (!uniqueAccountsMap[userId] || new Date(mtime) > new Date(uniqueAccountsMap[userId].mtime)) {
                uniqueAccountsMap[userId] = { userId, mtime, name };
            }
        });

        const uniqueAccounts = Object.values(uniqueAccountsMap);
        setMappings(uniqueAccounts);

        // Process available characters
        processAvailableCharacters(subDirs, associations);

        // Assign colors
        assignColors(subDirs);
    }, [subDirs, associations]);

    const processAvailableCharacters = (mappingsData, associationsData) => {
        const associatedCharIds = new Set(associationsData.map(a => a.charId));

        // Flatten characters and sort by `mtime`
        const availableChars = mappingsData
            .flatMap(mapping =>
                mapping.availableCharFiles.map(char => ({
                    ...char,
                    subDir: mapping.subDir,
                }))
            )
            .filter(char => !associatedCharIds.has(char.charId))
            .sort((a, b) => new Date(b.mtime) - new Date(a.mtime));

        // Keep only the most recent `mtime` per charId
        const uniqueAvailableChars = [];
        const seenCharIds = new Set();
        for (const char of availableChars) {
            if (!seenCharIds.has(char.charId)) {
                uniqueAvailableChars.push(char);
                seenCharIds.add(char.charId);
            }
        }

        setAvailableCharacters(uniqueAvailableChars);
        console.log('Unique Available Characters:', uniqueAvailableChars);
    };

    const assignColors = (mappingsData) => {
        const predefinedColors = ['#4caf50', '#f44336', '#ff9800', '#9c27b0', '#00bcd4', '#e91e63'];
        const uniqueMtims = new Set(
            mappingsData.flatMap(mapping => [
                ...mapping.availableUserFiles.map(user => user.mtime),
                ...mapping.availableCharFiles.map(char => char.mtime),
            ])
        );

        const colorMapping = Array.from(uniqueMtims).reduce((acc, mtime, index) => {
            acc[mtime] = predefinedColors[index % predefinedColors.length];
            return acc;
        }, {});

        setMtimeToColor(colorMapping);
        console.log('Mtime to Color Mapping:', colorMapping);
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
                // Refresh data
                onRefreshData();
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
                onRefreshData();
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
                            {mappings.map(mapping => (
                                <AccountCard
                                    key={mapping.userId}
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
                                        <Grid item xs={12} sm={6} key={char.charId}>
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

export default Mapping;
