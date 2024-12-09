import { useState, useEffect, useCallback } from 'react';
import PropTypes from 'prop-types';
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
    Divider,
    Tooltip,
    Box,
} from '@mui/material';
import LoadingButton from '@mui/lab/LoadingButton';
import BackupIcon from '@mui/icons-material/Backup';
import DeleteIcon from '@mui/icons-material/Delete';
import FolderOpenIcon from '@mui/icons-material/FolderOpen';
import UndoIcon from '@mui/icons-material/Undo';
import SyncIcon from '@mui/icons-material/Sync';
import SyncAllIcon from '@mui/icons-material/SyncAlt';

const MySwal = withReactContent(Swal);

const Sync = ({
                  settingsData,
                  associations,
                  currentSettingsDir,
                  isDefaultDir,
                  userSelections,
              }) => {
    const [isLoading, setIsLoading] = useState(false);
    const [selections, setSelections] = useState({});

    useEffect(() => {
        if (settingsData && settingsData.length > 0) {
            const initialSelections = { ...userSelections };
            settingsData.forEach(subDir => {
                if (!initialSelections[subDir.subDir]) {
                    initialSelections[subDir.subDir] = { charId: '', userId: '' };
                }
            });
            setSelections(initialSelections);
        }
    }, [settingsData, userSelections]);

    const getCharacterNameById = (charId) => {
        const assoc = associations.find(a => a.charId === charId);
        return assoc ? assoc.charName : 'Unknown';
    };

    const saveUserSelections = useCallback(async (newSelections) => {
        try {
            const response = await fetch('/api/save-user-selections', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                credentials: 'include',
                body: JSON.stringify(newSelections),
            });
            const result = await response.json();
            if (!response.ok || !result.success) {
                toast.error('Failed to save user selections.');
            }
        } catch (error) {
            console.error('Error saving user selections:', error);
            toast.error('An error occurred while saving user selections.');
        }
    }, []);

    const handleSelectionChange = (subDir, field, value) => {
        setSelections(prev => {
            const updated = {
                ...prev,
                [subDir]: {
                    ...prev[subDir],
                    [field]: value,
                }
            };
            saveUserSelections(updated);
            return updated;
        });
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
            const response = await fetch('/api/sync-subdirectory', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                credentials: 'include',
                body: JSON.stringify({ subDir, userId, charId }),
            });
            const result = await response.json();
            if (response.ok && result.success) {
                toast.success(result.message);
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
            const response = await fetch('/api/sync-all-subdirectories', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                credentials: 'include',
                body: JSON.stringify({ subDir, userId, charId }),
            });
            const result = await response.json();
            if (response.ok && result.success) {
                toast.success(result.message);
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
            setIsLoading(true);
            const response = await fetch('/api/choose-settings-dir', {
                method: 'POST',
                credentials: 'include',
            });
            const result = await response.json();
            if (response.ok && result.success) {
                toast.success('Settings directory chosen');
            } else {
                toast.error('Failed to choose settings directory.');
            }
        } catch (error) {
            console.error('Error choosing settings directory:', error);
            toast.error('Failed to choose settings directory.');
        } finally {
            setIsLoading(false);
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
            const response = await fetch('/api/reset-to-default-directory', {
                method: 'POST',
                credentials: 'include',
            });
            const result = await response.json();
            if (response.ok && result.success) {
                toast.success('Directory reset to default: Tranquility');
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
            const response = await fetch('/api/backup-directory', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                credentials: 'include',
                body: JSON.stringify({ targetDir: currentSettingsDir }),
            });
            const result = await response.json();
            if (response.ok && result.success) {
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
            const response = await fetch('/api/delete-backups', {
                method: 'POST',
                credentials: 'include',
            });
            const result = await response.json();
            if (response.ok && result.success) {
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

    return (
        <div className="pt-8 px-4 pb-10 bg-gray-900 min-h-screen text-teal-200">
            <Grid container spacing={2} justifyContent="center" alignItems="center" className="mb-4">
                <Grid item>
                    <Tooltip title="Backup Settings">
                        <span>
                            <Button
                                variant="contained"
                                color="primary"
                                onClick={handleBackup}
                                disabled={isLoading}
                                className="w-10 h-10 p-0 flex items-center justify-center"
                            >
                                <BackupIcon fontSize="small" />
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
                                onClick={handleDeleteBackups}
                                disabled={isLoading}
                                className="w-10 h-10 p-0 flex items-center justify-center"
                            >
                                <DeleteIcon fontSize="small" />
                            </Button>
                        </span>
                    </Tooltip>
                </Grid>
            </Grid>

            {isLoading && (
                <Box display="flex" justifyContent="center" alignItems="center" className="mb-4">
                    <CircularProgress color="primary" />
                </Box>
            )}

            <Grid container spacing={4}>
                {settingsData.map(subDir => (
                    <Grid item xs={12} sm={6} md={4} key={subDir.subDir}>
                        <Card className="bg-gray-800 text-teal-200 p-4 rounded-md shadow-md flex flex-col justify-between h-full">
                            <div>
                                <Typography
                                    variant="h6"
                                    gutterBottom
                                    sx={{ textAlign: 'center', fontWeight: 700 }}
                                >
                                    {subDir.subDir.replace('settings_', '')}
                                </Typography>
                                <Divider sx={{ my: 2, borderColor: 'rgba(255,255,255,0.2)' }} />
                                <div className="p-2 rounded-md border border-gray-600 bg-gray-700 mb-2">
                                    <FormControl fullWidth margin="normal">
                                        <InputLabel
                                            id={`char-select-label-${subDir.subDir}`}
                                            sx={{ color: '#fff' }}
                                        >
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
                                                borderRadius: 1,
                                                color: '#fff',
                                                '& .MuiOutlinedInput-notchedOutline': {
                                                    borderColor: 'rgba(255,255,255,0.2)',
                                                },
                                                '&:hover .MuiOutlinedInput-notchedOutline': {
                                                    borderColor: '#ffffff',
                                                },
                                                '&.Mui-focused .MuiOutlinedInput-notchedOutline': {
                                                    borderColor: '#ffffff',
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
                                        <InputLabel
                                            id={`user-select-label-${subDir.subDir}`}
                                            sx={{ color: '#fff' }}
                                        >
                                            Select User
                                        </InputLabel>
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
                                                borderRadius: 1,
                                                color: '#fff',
                                                '& .MuiOutlinedInput-notchedOutline': {
                                                    borderColor: 'rgba(255,255,255,0.2)',
                                                },
                                                '&:hover .MuiOutlinedInput-notchedOutline': {
                                                    borderColor: '#ffffff',
                                                },
                                                '&.Mui-focused .MuiOutlinedInput-notchedOutline': {
                                                    borderColor: '#ffffff',
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
                                </div>
                            </div>
                            <div>
                                <Grid container spacing={2}>
                                    <Grid item xs={6}>
                                        <Tooltip title="Sync this specific profile">
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
                                                    className="p-0"
                                                >
                                                    <SyncIcon fontSize="small" />
                                                </LoadingButton>
                                            </span>
                                        </Tooltip>
                                    </Grid>
                                    <Grid item xs={6}>
                                        <Tooltip title="Sync all profiles based on this selection">
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
                                                    className="p-0"
                                                >
                                                    <SyncAllIcon fontSize="small" />
                                                </LoadingButton>
                                            </span>
                                        </Tooltip>
                                    </Grid>
                                </Grid>
                            </div>
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
                                onClick={handleChooseSettingsDir}
                                disabled={isLoading}
                                className="w-10 h-10 p-0 flex items-center justify-center"
                            >
                                <FolderOpenIcon fontSize="small" />
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
                                    onClick={handleResetToDefault}
                                    disabled={isLoading}
                                    className="w-10 h-10 p-0 flex items-center justify-center"
                                >
                                    <UndoIcon fontSize="small" />
                                </Button>
                            </span>
                        </Tooltip>
                    </Grid>
                )}
            </Grid>
        </div>
    );
};

Sync.propTypes = {
    settingsData: PropTypes.arrayOf(
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
    associations: PropTypes.arrayOf(
        PropTypes.shape({
            userId: PropTypes.string.isRequired,
            charId: PropTypes.string.isRequired,
            charName: PropTypes.string.isRequired,
            mtime: PropTypes.string,
        })
    ).isRequired,
    currentSettingsDir: PropTypes.string.isRequired,
    isDefaultDir: PropTypes.bool.isRequired,
    userSelections: PropTypes.object.isRequired,
};

export default Sync;
