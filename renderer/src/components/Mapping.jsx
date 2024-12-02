// Mapping.jsx

import React, { useState, useEffect } from 'react';
import { toast, ToastContainer } from 'react-toastify';
import Swal from 'sweetalert2';
import withReactContent from 'sweetalert2-react-content';
import 'react-toastify/dist/ReactToastify.css';
import 'sweetalert2/dist/sweetalert2.min.css';
import { Button, Grid, Typography, Card, CardContent, List, ListItem, ListItemText, IconButton, CircularProgress } from '@mui/material';
import HomeIcon from '@mui/icons-material/Home';
import DeleteIcon from '@mui/icons-material/Delete';
import PeopleIcon from '@mui/icons-material/People';

const MySwal = withReactContent(Swal);

const Mapping = () => {
    const [mappings, setMappings] = useState([]);
    const [associations, setAssociations] = useState([]);
    const [availableCharacters, setAvailableCharacters] = useState([]);
    const [isLoading, setIsLoading] = useState(true);
    const [lastModificationTimes, setLastModificationTimes] = useState({});

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
                processLastModificationTimes(data.mappings);
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

        setAvailableCharacters(availableChars);
    };

    // Process last modification times for color coding
    const processLastModificationTimes = (mappingsData) => {
        const times = {};
        mappingsData.forEach(mapping => {
            mapping.availableCharFiles.forEach(char => {
                times[char.charId] = char.mtime;
            });
        });
        setLastModificationTimes(times);
    };

    // Determine border color based on modification time
    const getBorderColor = (mtime) => {
        const colors = [
            '#4caf50', // Green
            '#2196f3', // Blue
            '#ff9800', // Orange
            '#f44336', // Red
            '#9c27b0', // Purple
            '#00bcd4', // Cyan
        ];
        const uniqueTimes = Object.keys(lastModificationTimes).sort(); // Sort times for consistency
        const index = uniqueTimes.indexOf(mtime);
        return colors[index % colors.length] || '#000'; // Default to black if out of range
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
            <Grid container spacing={2} justifyContent="center" alignItems="center">
                <Grid item xs={12}>
                    <Typography variant="h4" align="center" gutterBottom>
                        Map Characters to Accounts
                    </Typography>
                </Grid>

                {/* Navigation and Action Buttons */}
                <Grid item xs={12} container spacing={2} justifyContent="center">
                    {/* Removed navigation button as navigation is handled via the header */}
                </Grid>

                {/* Loading Spinner */}
                {isLoading && (
                    <Grid item xs={12} container justifyContent="center">
                        <CircularProgress />
                    </Grid>
                )}

                {/* Users and Characters Sections */}
                <Grid item xs={12}>
                    <Grid container spacing={3}>
                        {/* Accounts Section */}
                        <Grid item xs={12} md={6}>
                            <Typography variant="h5" gutterBottom>
                                Accounts
                            </Typography>
                            {mappings.map(mapping => (
                                <Card
                                    key={`${mapping.subDir}-account`}
                                    onDragOver={(e) => e.preventDefault()}
                                    onDrop={(e) => handleDrop(e, mapping.subDir)}
                                    aria-label={`Drop character to associate with account ${mapping.subDir}`}
                                    style={{
                                        marginBottom: '16px',
                                        border: `2px solid ${getBorderColor(associations.find(assoc => assoc.userId === mapping.subDir)?.mtime || '#000')}`
                                    }}
                                >
                                    <CardContent>
                                        <Typography variant="h6">
                                            {mapping.subDir.replace('settings_', '')}
                                        </Typography>
                                        <List>
                                            {associations
                                                .filter(assoc => assoc.userId === mapping.subDir)
                                                .map(assoc => (
                                                    <ListItem
                                                        key={assoc.charId} // Assuming charId is unique after filtering
                                                        secondaryAction={
                                                            <IconButton
                                                                edge="end"
                                                                aria-label="unassociate"
                                                                onClick={() => handleUnassociate(mapping.subDir, assoc.charId, assoc.charName)}
                                                            >
                                                                <DeleteIcon />
                                                            </IconButton>
                                                        }
                                                    >
                                                        <ListItemText
                                                            primary={`${assoc.charName} (${assoc.charId})`}
                                                            secondary={`Last Modified: ${new Date(assoc.mtime).toLocaleString()}`}
                                                        />
                                                    </ListItem>
                                                ))}
                                        </List>
                                    </CardContent>
                                </Card>
                            ))}
                        </Grid>

                        {/* Available Characters Section */}
                        <Grid item xs={12} md={6}>
                            <Typography variant="h5" gutterBottom>
                                Available Characters
                            </Typography>
                            <Grid container spacing={2}>
                                {availableCharacters.map(char => (
                                    <Grid item xs={12} sm={6} key={char.charId}>
                                        <Card
                                            draggable
                                            onDragStart={(e) => handleDragStart(e, char.charId)}
                                            aria-label={`Draggable character ${char.name}`}
                                            title={`Drag to associate with an account`}
                                            style={{
                                                border: `2px solid ${getBorderColor(char.mtime)}`
                                            }}
                                        >
                                            <CardContent>
                                                <Typography variant="h6">
                                                    {char.name}
                                                </Typography>
                                                <Typography variant="body2" color="textSecondary">
                                                    ID: {char.charId}
                                                </Typography>
                                                <Typography variant="caption" color="textSecondary">
                                                    Last Modified: {new Date(char.mtime).toLocaleString()}
                                                </Typography>
                                            </CardContent>
                                        </Card>
                                    </Grid>
                                ))}
                                {availableCharacters.length === 0 && !isLoading && (
                                    <Grid item xs={12}>
                                        <Typography variant="body1" align="center">
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
