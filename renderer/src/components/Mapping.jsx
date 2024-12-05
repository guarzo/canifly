// Mapping.jsx

import { useState, useEffect } from 'react';
import { toast } from 'react-toastify'; // Removed ToastContainer
import Swal from 'sweetalert2';
import withReactContent from 'sweetalert2-react-content';
import 'sweetalert2/dist/sweetalert2.min.css';
import { Grid, CircularProgress, Box } from '@mui/material';
import AccountCard from './MapAccountCard';
import CharacterCard from './MapCharacterCard';

const MySwal = withReactContent(Swal);

const Mapping = () => {
    const [mappings, setMappings] = useState([]); // Unique Accounts
    const [associations, setAssociations] = useState([]); // Associations
    const [availableCharacters, setAvailableCharacters] = useState([]); // Available Characters
    const [isLoading, setIsLoading] = useState(true);
    const [mtimeToColor, setMtimeToColor] = useState({});

    useEffect(() => {
        loadDataAndRender();
    }, []);

    // Normalization Function
    const normalizeAssociation = (assoc) => {
        return {
            ...assoc,
            charId: typeof assoc.charId === 'string' ? assoc.charId : (assoc.charId?.charId || 'Unknown'),
            charName: typeof assoc.charName === 'string' ? assoc.charName : (assoc.charId?.charName || 'Unknown'),
            mtime: assoc.mtime || (typeof assoc.charId === 'object' ? assoc.charId.mtime : new Date().toISOString()),
            userId: typeof assoc.userId === 'string' ? assoc.userId : 'Unknown',
        };
    };

    const loadDataAndRender = async () => {
        try {
            setIsLoading(true);
            const data = await window.electronAPI.loadMappings();
            if (data && data.mappings) {
                // Step 1: Map each subDir to its latest user based on mtime
                const subDirToLatestUserId = {};

                data.mappings.forEach(mapping => {
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
                console.log('Unique Accounts:', uniqueAccounts);
                setMappings(uniqueAccounts);

                // Step 3: Normalize associations
                const mappedAssociations = data.associations.map(assoc => normalizeAssociation(assoc));

                // Debugging: Verify the structure of mapped associations
                console.log('Mapped Associations:', mappedAssociations);
                mappedAssociations.forEach(assoc => {
                    console.assert(typeof assoc.charId === 'string', `charId is not a string: ${assoc.charId}`);
                    console.assert(typeof assoc.userId === 'string' && assoc.userId !== 'Unknown', `userId is invalid: ${assoc.userId}`);
                    console.assert(typeof assoc.charName === 'string', `charName is not a string: ${assoc.charName}`);
                });

                setAssociations(mappedAssociations);

                // Step 4: Process available characters using mappedAssociations
                processAvailableCharacters(data.mappings, mappedAssociations);

                // Step 5: Assign colors
                assignColors(data.mappings);
            } else {
                toast.error('Failed to load mappings: No data received.');
            }
        } catch (error) {
            console.error('Error loading mappings:', error);
            toast.error('Failed to load mappings: ' + error.message);
        } finally {
            setIsLoading(false);
        }
    };

    const processAvailableCharacters = (mappingsData, associationsData) => {
        const associatedCharIds = new Set(associationsData.map(a => a.charId));

        // Debugging: Verify associatedCharIds
        console.log('Associated Character IDs:', associatedCharIds);

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

        // Keep only the most recent `mtime` for each `charId`
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

            // Pass only charId as a string
            const result = await window.electronAPI.associateCharacter(userId, charId);

            // Ensure result.message is a string
            const successMessage = typeof result.message === 'string' ? result.message : JSON.stringify(result.message);
            const errorMessage = typeof result.message === 'string' ? result.message : JSON.stringify(result.message);

            if (result.success) {
                toast.success(successMessage);
                await loadDataAndRender();
            } else {
                toast.error(`Association failed: ${errorMessage}`);
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
            const result = await window.electronAPI.unassociateCharacter(userId, charId);

            // Ensure result.message is a string
            const successMessage = typeof result.message === 'string' ? result.message : JSON.stringify(result.message);
            const errorMessage = typeof result.message === 'string' ? result.message : JSON.stringify(result.message);

            if (result.success) {
                toast.success(successMessage);
                await loadDataAndRender();
            } else {
                toast.error(`Unassociation failed: ${errorMessage}`);
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
