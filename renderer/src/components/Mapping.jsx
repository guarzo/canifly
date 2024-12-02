// Mapping.jsx

import React, { useState, useEffect } from 'react';
import { toast, ToastContainer } from 'react-toastify';
import Swal from 'sweetalert2';
import withReactContent from 'sweetalert2-react-content';
import 'react-toastify/dist/ReactToastify.css';
import 'sweetalert2/dist/sweetalert2.min.css';
import { Grid, Typography, CircularProgress, useTheme } from '@mui/material';
import AccountCard from './MapAccountCard';
import CharacterCard from './MapCharacterCard';

const MySwal = withReactContent(Swal);

const Mapping = () => {
    const [mappings, setMappings] = useState([]);
    const [associations, setAssociations] = useState([]);
    const [availableCharacters, setAvailableCharacters] = useState([]);
    const [isLoading, setIsLoading] = useState(true);
    const [mtimeToColor, setMtimeToColor] = useState({});
    const theme = useTheme();

    // Predefined list of colors to assign to unique mtimes
    const predefinedColors = [
        theme.palette.primary.main,    // Teal
        theme.palette.secondary.main,  // Red
        '#f59e0b',                     // Amber
        '#8b5cf6',                     // Violet
        '#10b981',                     // Emerald
        '#ec4899',                     // Pink
        '#f97316',                     // Orange
        '#3b82f6',                     // Blue
    ];

    useEffect(() => {
        loadDataAndRender();
    }, []);

    // Function to load mappings and associations from the main process
    const loadDataAndRender = async () => {
        try {
            setIsLoading(true);
            const data = await window.electronAPI.loadMappings();
            console.log('Mapping.jsx - loadMappings Response:', data);
            if (data && data.mappings) {
                setMappings(data.mappings);
                setAssociations(data.associations);
                processAvailableCharacters(data.mappings, data.associations);
                assignColors(data.mappings, data.associations);
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

    // Process available characters by filtering out already associated ones
    const processAvailableCharacters = (mappingsData, associationsData) => {
        let allCharacters = [];
        mappingsData.forEach(mapping => {
            mapping.availableCharFiles.forEach(char => {
                allCharacters.push({
                    ...char,
                    subDir: mapping.subDir,
                });
            });
        });

        // Filter out already associated characters
        const associatedCharIds = associationsData.map(a => a.charId);
        const availableChars = allCharacters.filter(char => !associatedCharIds.includes(char.charId));

        // Sort available characters by mtime descending to have the latest first
        availableChars.sort((a, b) => new Date(b.mtime) - new Date(a.mtime));

        // Remove duplicates by keeping only the first occurrence (latest)
        const uniqueAvailableChars = [];
        const seenCharIds = new Set();

        availableChars.forEach(char => {
            if (!seenCharIds.has(char.charId)) {
                uniqueAvailableChars.push(char);
                seenCharIds.add(char.charId);
            }
        });

        setAvailableCharacters(uniqueAvailableChars);
    };

    // Assign colors based on unique mtimes
    const assignColors = (mappingsData, associationsData) => {
        const uniqueMtims = new Set();

        // Extract mtimes from accounts
        mappingsData.forEach(mapping => {
            mapping.availableUserFiles.forEach(user => {
                uniqueMtims.add(user.mtime);
            });
        });

        // Extract mtimes from characters
        mappingsData.forEach(mapping => {
            mapping.availableCharFiles.forEach(char => {
                uniqueMtims.add(char.mtime);
            });
        });

        const mtimeArray = Array.from(uniqueMtims);
        const colorMapping = {};

        mtimeArray.forEach((mtime, index) => {
            colorMapping[mtime] = predefinedColors[index % predefinedColors.length];
        });

        setMtimeToColor(colorMapping);
    };

    // Handle drag start event for characters
    const handleDragStart = (event, charId) => {
        event.dataTransfer.setData('text/plain', charId);
    };

    // Handle drop event on accounts
    const handleDrop = async (event, userId) => {
        event.preventDefault();
        const charId = event.dataTransfer.getData('text/plain');

        if (!charId) {
            toast.error('No character selected for association.');
            return;
        }

        // Get character name
        const char = availableCharacters.find(c => c.charId === charId);
        const charName = char ? char.name : 'Unknown';

        // Show confirmation modal
        const confirmAssoc = await MySwal.fire({
            title: 'Confirm Association',
            text: `Associate "${charName}" (${charId}) with account ${userId}?`,
            icon: 'question',
            showCancelButton: true,
            confirmButtonText: 'Yes, associate it!',
            cancelButtonText: 'Cancel',
        });

        if (!confirmAssoc.isConfirmed) return;

        try {
            setIsLoading(true);
            const result = await window.electronAPI.associateCharacter(userId, charId);
            if (result.success) {
                toast.success(result.message);
                await loadDataAndRender();
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

    // Handle unassociating a character from an account
    const handleUnassociate = async (userId, charId, charName) => {
        const confirmUnassoc = await MySwal.fire({
            title: 'Confirm Unassociation',
            text: `Unassociate "${charName}" (${charId}) from account ${userId}?`,
            icon: 'warning',
            showCancelButton: true,
            confirmButtonText: 'Yes, unassociate it!',
            cancelButtonText: 'Cancel',
        });

        if (!confirmUnassoc.isConfirmed) return;

        try {
            setIsLoading(true);
            const result = await window.electronAPI.unassociateCharacter(userId, charId);
            if (result.success) {
                toast.success(result.message);
                await loadDataAndRender();
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
        <div>
            <ToastContainer />
            <Grid container spacing={2} justifyContent="center" alignItems="flex-start">
                {/* Header Spacer */}
                <Grid item xs={12}>
                    {/* Spacer to maintain layout */}
                    <div style={{ height: '60px' }}></div>
                </Grid>

                {/* Loading Spinner */}
                {isLoading && (
                    <Grid item xs={12} container justifyContent="center">
                        <CircularProgress color="primary" />
                    </Grid>
                )}

                {/* Users and Characters Sections */}
                <Grid item xs={12}>
                    <Grid container spacing={3}>
                        {/* Accounts Section */}
                        <Grid item xs={12} md={6}>
                            <Typography variant="h5" gutterBottom color="text.primary">
                                Accounts
                            </Typography>
                            {mappings.map(mapping => (
                                <AccountCard
                                    key={`account-${mapping.subDir}`}
                                    mapping={mapping}
                                    associations={associations}
                                    handleUnassociate={handleUnassociate}
                                    handleDrop={handleDrop}
                                    mtimeToColor={mtimeToColor}
                                />
                            ))}
                        </Grid>

                        {/* Available Characters Section */}
                        <Grid item xs={12} md={6}>
                            <Typography variant="h5" gutterBottom color="text.primary">
                                Available Characters
                            </Typography>
                            <Grid container spacing={2}>
                                {availableCharacters.map(char => (
                                    <Grid item xs={12} sm={6} key={`char-${char.charId}`}>
                                        <CharacterCard
                                            char={char}
                                            handleDragStart={handleDragStart}
                                            mtimeToColor={mtimeToColor}
                                        />
                                    </Grid>
                                ))}
                                {availableCharacters.length === 0 && !isLoading && (
                                    <Grid item xs={12}>
                                        <Typography variant="body1" align="center" color="text.secondary">
                                            No available characters to associate.
                                        </Typography>
                                    </Grid>
                                )}
                            </Grid>
                        </Grid>
                    </Grid>
                </Grid>
            </Grid>
        </div>
    );

};

export default Mapping;
