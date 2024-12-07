// Sync.jsx

import React, { useState, useEffect } from 'react';
import { toast } from 'react-toastify';
import Swal from 'sweetalert2';
import withReactContent from 'sweetalert2-react-content';
import 'sweetalert2/dist/sweetalert2.min.css';
import {
    Button,
    Select,
    MenuItem,
    FormControl,
    InputLabel,
    CircularProgress,
    Grid,
    Typography,
    Card,
    CardContent,
    Box,
    Divider,
    Tooltip,
} from '@mui/material';
import LoadingButton from '@mui/lab/LoadingButton';
import BackupIcon from '@mui/icons-material/Backup';
import DeleteIcon from '@mui/icons-material/Delete';
import FolderOpenIcon from '@mui/icons-material/FolderOpen';
import UndoIcon from '@mui/icons-material/Undo';
import SyncIcon from '@mui/icons-material/Sync';
import SyncAllIcon from '@mui/icons-material/SyncAlt';

const MySwal = withReactContent(Swal);

const Sync = () => {
    const [settingsData, setSettingsData] = useState([]);
    const [associations, setAssociations] = useState([]);
    const [currentSettingsDir, setCurrentSettingsDir] = useState('');
    const [isLoading, setIsLoading] = useState(true);
    const [isDefaultDir, setIsDefaultDir] = useState(false);
    const [selections, setSelections] = useState({});

    useEffect(() => {
        loadDataAndRender();
    }, []);

    const loadDataAndRender = async () => {
        try {
            setIsLoading(true);
            const data = await window.electronAPI.loadSettings();
            console.log('Sync.jsx - loadSettings Response:', data);
            if (data && data.settingsData) {
                setSettingsData(data.settingsData);
                setAssociations(data.associations);
                setCurrentSettingsDir(data.currentSettingsDir);
                setIsDefaultDir(data.isDefaultDir);

                const initialSelections = {};
                data.settingsData.forEach(subDir => {
                    initialSelections[subDir.subDir] = {
                        charId: '',
                        userId: '',
                    };
                });
                setSelections(initialSelections);
            } else {
                toast.error('Failed to load settings: No data received.');
            }
        } catch (error) {
            console.error('Error loading settings:', error);
            toast.error('Failed to load settings: ' + error.message);
        } finally {
            setIsLoading(false);
        }
    };

    const handleSelectionChange = (subDir, field, value) => {
        setSelections(prev => ({
            ...prev,
            [subDir]: {
                ...prev[subDir],
                [field]: value,
            }
        }));
    };

    const handleSync = async (subDir) => {
        const { userId, charId } = selections[subDir];
        if (!userId || !charId) {
            toast.error('Please select both a user and a character to sync.');
            return;
        }

        const charName = getCharacterNameById(charId);

        const confirmSync = await MySwal.fire({
            title: 'Confirm Sync',
            text: `Use "${charName}" (${charId}) on account ${userId} to overwrite all files in profile "${subDir}"?`,
            icon: 'question',
            showCancelButton: true,
            confirmButtonText: 'Yes, sync it!',
            cancelButtonText: 'Cancel',
        });

        if (!confirmSync.isConfirmed) return;

        try {
            setIsLoading(true);
            toast.info('Syncing...', { autoClose: 1500 });
            const result = await window.electronAPI.syncSubdirectory(subDir, userId, charId);
            if (result.success) {
                toast.success(result.message);
                await loadDataAndRender();
            } else {
                toast.error(`Sync failed: ${result.message}`);
            }
        } catch (error) {
            console.error('Error syncing:', error);
            toast.error('Sync operation failed.');
        } finally {
            setIsLoading(false);
        }
    };

    const handleSyncAll = async (subDir) => {
        const { userId, charId } = selections[subDir];
        if (!userId || !charId) {
            toast.error('Please select both a user and a character to sync-all.');
            return;
        }

        const charName = getCharacterNameById(charId);

        const confirmSyncAll = await MySwal.fire({
            title: 'Confirm Sync All',
            text: `Use "${charName}" with ${userId} to overwrite files across all profiles?`,
            icon: 'question',
            showCancelButton: true,
            confirmButtonText: 'Yes, sync all!',
            cancelButtonText: 'Cancel',
        });

        if (!confirmSyncAll.isConfirmed) return;

        try {
            setIsLoading(true);
            const result = await window.electronAPI.syncAllSubdirectories(subDir, userId, charId);
            if (result.success) {
                toast.success(result.message);
                await loadDataAndRender();
            } else {
                toast.error(`Sync-All failed: ${result.message}`);
            }
        } catch (error) {
            console.error('Error syncing-all:', error);
            toast.error('Sync-All operation failed.');
        } finally {
            setIsLoading(false);
        }
    };

    const handleChooseSettingsDir = async () => {
        try {
            const result = await window.electronAPI.chooseSettingsDir();
            if (result.success) {
                setCurrentSettingsDir(result.settingsDir);
                const isDefault = await window.electronAPI.checkIfDefaultDirectory(result.settingsDir);
                setIsDefaultDir(isDefault);
                await loadDataAndRender();
            } else {
                toast.error('Failed to choose settings directory.');
            }
        } catch (error) {
            console.error('Error choosing settings directory:', error);
            toast.error('Failed to choose settings directory.');
        }
    };

    const handleResetToDefault = async () => {
        const confirmReset = await MySwal.fire({
            title: 'Reset to Default',
            text: 'Are you sure you want to reset the settings directory to default (Tranquility)?',
            icon: 'warning',
            showCancelButton: true,
            confirmButtonText: 'Yes, reset it!',
            cancelButtonText: 'Cancel',
        });

        if (!confirmReset.isConfirmed) return;

        try {
            setIsLoading(true);
            const result = await window.electronAPI.resetToDefaultDirectory();
            if (result.success) {
                toast.success('Directory reset to default: Tranquility');
                setCurrentSettingsDir(result.settingsDir);
                setIsDefaultDir(true);
                await loadDataAndRender();
            } else {
                toast.error(`Failed to reset directory: ${result.message}`);
            }
        } catch (error) {
            console.error('Error resetting directory:', error);
            toast.error('Failed to reset directory. Please try again.');
        } finally {
            setIsLoading(false);
        }
    };

    const handleBackup = async () => {
        const confirmBackup = await MySwal.fire({
            title: 'Confirm Backup',
            text: 'Are you sure you want to backup your settings?',
            icon: 'question',
            showCancelButton: true,
            confirmButtonText: 'Yes, backup it!',
            cancelButtonText: 'Cancel',
        });

        if (!confirmBackup.isConfirmed) return;

        try {
            setIsLoading(true);
            const result = await window.electronAPI.backupDirectory(currentSettingsDir);
            if (result.success) {
                toast.success(result.message);
            } else {
                toast.error(`Backup failed: ${result.message}`);
            }
        } catch (error) {
            console.error('Error during backup:', error);
            toast.error('Backup operation failed.');
        } finally {
            setIsLoading(false);
        }
    };

    const handleDeleteBackups = async () => {
        const confirmDelete = await MySwal.fire({
            title: 'Confirm Delete',
            text: 'Are you sure you want to delete all backups?',
            icon: 'warning',
            showCancelButton: true,
            confirmButtonText: 'Yes, delete them!',
            cancelButtonText: 'Cancel',
        });

        if (!confirmDelete.isConfirmed) return;

        try {
            setIsLoading(true);
            const result = await window.electronAPI.deleteBackups();
            if (result.success) {
                toast.success(result.message);
            } else {
                toast.error(`Delete backups failed: ${result.message}`);
            }
        } catch (error) {
            console.error('Error deleting backups:', error);
            toast.error('Delete backups operation failed.');
        } finally {
            setIsLoading(false);
        }
    };

    const getCharacterNameById = (charId) => {
        const assoc = associations.find(a => a.charId === charId);
        return assoc ? assoc.charName : 'Unknown';
    };

    return (
        <Box sx={{ pt: 10, px: 4, pb: 10, backgroundColor: 'background.default', minHeight: '100vh' }}>
            {/* Top Buttons Section */}
            <Grid container spacing={2} justifyContent="center" alignItems="center" sx={{ marginBottom: 4 }}>
                <Grid item>
                    <Tooltip title="Backup Settings">
                        <span>
                            <Button
                                variant="contained"
                                color="primary"
                                startIcon={<BackupIcon />}
                                onClick={handleBackup}
                                disabled={isLoading}
                                fullWidth
                                sx={{
                                    transition: 'transform 0.2s, box-shadow 0.2s',
                                    '&:hover': {
                                        transform: 'scale(1.02)',
                                        boxShadow: 6,
                                    },
                                }}
                            >
                                Backup
                            </Button>
                        </span>
                    </Tooltip>
                </Grid>
                <Grid item>
                    <Tooltip title="Delete All Backups">
                        <span>
                            <Button
                                variant="contained"
                                color="secondary"
                                startIcon={<DeleteIcon />}
                                onClick={handleDeleteBackups}
                                disabled={isLoading}
                                fullWidth
                                sx={{
                                    transition: 'transform 0.2s, box-shadow 0.2s',
                                    '&:hover': {
                                        transform: 'scale(1.02)',
                                        boxShadow: 6,
                                    },
                                }}
                            >
                                Delete Backups
                            </Button>
                        </span>
                    </Tooltip>
                </Grid>
            </Grid>

            {isLoading && (
                <Box display="flex" justifyContent="center" alignItems="center" sx={{ marginBottom: 4 }}>
                    <CircularProgress color="primary" />
                </Box>
            )}

            <Grid container spacing={4}>
                {settingsData.map(subDir => (
                    <Grid item xs={12} sm={6} md={4} key={subDir.subDir}>
                        <Card
                            sx={{
                                boxShadow: 3,
                                transition: 'transform 0.2s, box-shadow 0.2s',
                                '&:hover': {
                                    transform: 'scale(1.02)',
                                    boxShadow: 6,
                                },
                                backgroundColor: 'background.paper',
                            }}
                        >
                            <CardContent>
                                <Typography
                                    variant="h6"
                                    gutterBottom
                                    sx={{ textAlign: 'center', fontWeight: 600 }}
                                >
                                    {subDir.subDir.replace('settings_', '')}
                                </Typography>
                                <Divider sx={{ marginY: 2 }} />

                                <Box
                                    sx={{
                                        padding: 2,
                                        border: '1px solid rgba(0, 0, 0, 0.12)',
                                        borderRadius: 2,
                                        bgcolor: 'rgba(0, 0, 0, 0.03)',
                                        marginBottom: 2,
                                    }}
                                >
                                    <FormControl fullWidth margin="normal">
                                        <InputLabel id={`char-select-label-${subDir.subDir}`}>
                                            Select Character
                                        </InputLabel>
                                        <Select
                                            labelId={`char-select-label-${subDir.subDir}`}
                                            id={`char-select-${subDir.subDir}`}
                                            value={selections[subDir.subDir]?.charId || ''}
                                            label="Select Character"
                                            onChange={(e) =>
                                                handleSelectionChange(
                                                    subDir.subDir,
                                                    'charId',
                                                    e.target.value
                                                )
                                            }
                                            sx={{
                                                backgroundColor: 'background.paper',
                                                borderRadius: 1,
                                                '& .MuiSelect-select': {
                                                    padding: '10px 14px',
                                                },
                                            }}
                                        >
                                            <MenuItem value="">
                                                <em>-- Select Character --</em>
                                            </MenuItem>
                                            {subDir.availableCharFiles.map(char => (
                                                <MenuItem key={char.charId} value={char.charId}>
                                                    {char.name}
                                                </MenuItem>
                                            ))}
                                        </Select>
                                    </FormControl>
                                    <FormControl fullWidth margin="normal">
                                        <InputLabel id={`user-select-label-${subDir.subDir}`}>Select User</InputLabel>
                                        <Select
                                            labelId={`user-select-label-${subDir.subDir}`}
                                            id={`user-select-${subDir.subDir}`}
                                            value={selections[subDir.subDir]?.userId || ''}
                                            label="Select User"
                                            onChange={(e) =>
                                                handleSelectionChange(
                                                    subDir.subDir,
                                                    'userId',
                                                    e.target.value
                                                )
                                            }
                                            sx={{
                                                backgroundColor: 'background.paper',
                                                borderRadius: 1,
                                                '& .MuiSelect-select': {
                                                    padding: '10px 14px',
                                                },
                                            }}
                                        >
                                            <MenuItem value="">
                                                <em>-- Select User --</em>
                                            </MenuItem>
                                            {subDir.availableUserFiles.map(user => (
                                                <MenuItem key={user.userId} value={user.userId}>
                                                    {user.userId}
                                                </MenuItem>
                                            ))}
                                        </Select>
                                    </FormControl>
                                </Box>

                                <Grid container spacing={2}>
                                    <Grid item xs={6}>
                                        <Tooltip title="Sync Profile">
                                            <span style={{ width: '100%', display: 'block' }}>
                                                <LoadingButton
                                                    variant="contained"
                                                    color="primary"
                                                    onClick={() => handleSync(subDir.subDir)}
                                                    loading={isLoading}
                                                    fullWidth
                                                    disabled={
                                                        isLoading ||
                                                        !selections[subDir.subDir]?.charId ||
                                                        !selections[subDir.subDir]?.userId
                                                    }
                                                    sx={{
                                                        transition: 'transform 0.2s, box-shadow 0.2s',
                                                        '&:hover': {
                                                            transform: 'scale(1.02)',
                                                            boxShadow: 6,
                                                        },
                                                        minWidth: 0,
                                                        padding: '6px',
                                                    }}
                                                    startIcon={<SyncIcon color="inherit" />}
                                                />
                                            </span>
                                        </Tooltip>
                                    </Grid>
                                    <Grid item xs={6}>
                                        <Tooltip title="Sync All Profiles">
                                            <span style={{ width: '100%', display: 'block' }}>
                                                <LoadingButton
                                                    variant="contained"
                                                    color="secondary"
                                                    onClick={() => handleSyncAll(subDir.subDir)}
                                                    loading={isLoading}
                                                    fullWidth
                                                    disabled={
                                                        isLoading ||
                                                        !selections[subDir.subDir]?.charId ||
                                                        !selections[subDir.subDir]?.userId
                                                    }
                                                    sx={{
                                                        transition: 'transform 0.2s, box-shadow 0.2s',
                                                        '&:hover': {
                                                            transform: 'scale(1.02)',
                                                            boxShadow: 6,
                                                        },
                                                        minWidth: 0,
                                                        padding: '6px',
                                                    }}
                                                    startIcon={<SyncAllIcon color="inherit" />}
                                                />
                                            </span>
                                        </Tooltip>
                                    </Grid>
                                </Grid>
                            </CardContent>
                        </Card>
                    </Grid>
                ))}
            </Grid>

            <Grid container spacing={2} justifyContent="center" alignItems="center" sx={{ marginTop: 4 }}>
                <Grid item>
                    <Tooltip title="Choose Settings Directory">
                        <span>
                            <Button
                                variant="contained"
                                color="info"
                                startIcon={<FolderOpenIcon />}
                                onClick={handleChooseSettingsDir}
                                disabled={isLoading}
                                fullWidth
                                sx={{
                                    transition: 'transform 0.2s, box-shadow 0.2s',
                                    '&:hover': {
                                        transform: 'scale(1.02)',
                                        boxShadow: 6,
                                    },
                                }}
                            >
                                Choose Directory
                            </Button>
                        </span>
                    </Tooltip>
                </Grid>
                {!isDefaultDir && (
                    <Grid item>
                        <Tooltip title="Reset to Default Directory">
                            <span>
                                <Button
                                    variant="contained"
                                    color="warning"
                                    startIcon={<UndoIcon />}
                                    onClick={handleResetToDefault}
                                    disabled={isLoading}
                                    fullWidth
                                    sx={{
                                        transition: 'transform 0.2s, box-shadow 0.2s',
                                        '&:hover': {
                                            transform: 'scale(1.02)',
                                            boxShadow: 6,
                                        },
                                    }}
                                >
                                    Reset to Default
                                </Button>
                            </span>
                        </Tooltip>
                    </Grid>
                )}
            </Grid>
        </Box>
    );
};

export default Sync;
