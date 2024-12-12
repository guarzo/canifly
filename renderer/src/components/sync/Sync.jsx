import { useState, useEffect, useCallback } from 'react';
import PropTypes from 'prop-types';
import { toast } from 'react-toastify';
import { useConfirmDialog } from '../partials/useConfirmDialog.jsx';
import {
    CircularProgress,
    Grid,
    Box,
    Typography
} from '@mui/material';
import SyncActionsBar from './SyncActionsBar.jsx';
import SubDirectoryCard from './SubDirectoryCard.jsx';

const Sync = ({
                  settingsData,
                  associations,
                  currentSettingsDir,
                  userSelections,
                  lastBackupDir,
                  backEndURL,
              }) => {
    const [isLoading, setIsLoading] = useState(false);
    const [selections, setSelections] = useState({});
    const [showConfirmDialog, confirmDialog] = useConfirmDialog();
    const [isDefaultDir, setIsDefaultDir] = useState(true);

    useEffect(() => {
        if (settingsData && settingsData.length > 0) {
            const initialSelections = { ...userSelections };
            settingsData.forEach(subDir => {
                if (!initialSelections[subDir.profile]) {
                    initialSelections[subDir.profile] = { charId: '', userId: '' };
                }
            });
            setSelections(initialSelections);
        }
    }, [settingsData, userSelections]);

    const saveUserSelections = useCallback(async (newSelections) => {
        try {
            const response = await fetch(`${backEndURL}/api/save-user-selections`, {
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
    }, [backEndURL]);

    const handleSelectionChange = (profile, field, value) => {
        setSelections(prev => {
            let updated = {
                ...prev,
                [profile]: {
                    ...prev[profile],
                    [field]: value,
                }
            };

            // Auto-select user if associated with charId
            if (field === 'charId' && value) {
                const assoc = associations.find(a => a.charId === value);
                if (assoc) {
                    updated[profile].userId = assoc.userId;
                }
            }

            saveUserSelections(updated);
            return updated;
        });
    };

    const handleSync = async (profile) => {
        const { userId, charId } = selections[profile];
        if (!userId || !charId) {
            toast.error('Please select both a user and a character to sync.');
            return;
        }

        const confirmSync = await showConfirmDialog({
            title: 'Confirm Sync',
            message: 'Are you sure you want to sync this profile with the chosen character and user?',
        });

        if (!confirmSync.isConfirmed) return;

        try {
            setIsLoading(true);
            toast.info('Syncing...', { autoClose: 1500 });
            const response = await fetch(`${backEndURL}/api/sync-subdirectory`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                credentials: 'include',
                body: JSON.stringify({ profile, userId, charId }),
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

    const handleSyncAll = async (profile) => {
        const { userId, charId } = selections[profile];
        if (!userId || !charId) {
            toast.error('Please select both a user and a character for Sync-All.');
            return;
        }

        const confirmSyncAll = await showConfirmDialog({
            title: 'Confirm Sync All',
            message: 'Are you sure you want to sync all profiles with these selections?',
        });

        if (!confirmSyncAll.isConfirmed) return;

        try {
            setIsLoading(true);
            const response = await fetch(`${backEndURL}/api/sync-all-subdirectories`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                credentials: 'include',
                body: JSON.stringify({ profile, userId, charId }),
            });
            const result = await response.json();
            if (response.ok && result.success) {
                toast.success(`Sync-All complete: ${result.message}`);
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
            const chosenDir = await window.electronAPI.chooseDirectory();
            if (!chosenDir) {
                toast.info('No directory chosen.');
                setIsLoading(false);
                return;
            }

            const response = await fetch(`${backEndURL}/api/choose-settings-dir`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                credentials: 'include',
                body: JSON.stringify({ directory: chosenDir }),
            });

            const result = await response.json();
            if (response.ok && result.success) {
                setIsDefaultDir(false);
                toast.success(`Settings directory chosen: ${chosenDir}`);
            } else {
                toast.error(`Failed to choose settings directory: ${result.error || 'Unknown error'}`);
            }
        } catch (error) {
            console.error('Error choosing settings directory:', error);
            toast.error('Failed to choose settings directory.');
        } finally {
            setIsLoading(false);
        }
    };

    const handleBackup = async () => {
        try {
            setIsLoading(true);

            const chosenDir = await window.electronAPI.chooseDirectory(lastBackupDir || null);
            if (!chosenDir) {
                toast.info('No backup directory chosen. Backup canceled.');
                setIsLoading(false);
                return;
            }

            toast.info('Starting backup...');

            const response = await fetch(`${backEndURL}/api/backup-directory`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                credentials: 'include',
                body: JSON.stringify({ targetDir: currentSettingsDir, backupDir: chosenDir }),
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

    const handleResetToDefault = async () => {
        const confirmReset = await showConfirmDialog({
            title: 'Reset to Default',
            message: 'Reset the settings directory to default (Tranquility)?',
        });

        if (!confirmReset.isConfirmed) return;

        try {
            setIsLoading(true);
            const response = await fetch(`${backEndURL}/api/reset-to-default-directory`, {
                method: 'POST',
                credentials: 'include',
            });
            const result = await response.json();
            if (response.ok && result.success) {
                setIsDefaultDir(true);
                toast.success('Directory reset to default: Tranquility');
            } else {
                toast.error(`Failed to reset directory: ${result.message}`);
            }
        } catch (error) {
            console.error('Error resetting directory:', error);
            toast.error('Failed to reset directory.');
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <div className="bg-gray-900 min-h-screen text-teal-200 px-4 pb-10 pt-16">
            {/* Top heading with subtle instructions */}
            <Box className="max-w-7xl mx-auto mb-6">
                <Box className="bg-gradient-to-r from-gray-900 to-gray-800 p-4 rounded-md shadow-md">
                    <Typography variant="h4" sx={{ color: '#14b8a6', fontWeight: 'bold', marginBottom: '0.5rem' }}>
                        Sync Settings
                    </Typography>
                    <Typography variant="body2" sx={{ color: '#99f6e4' }}>
                        Select a character and user for each profile, then sync or sync all. You can also choose or reset the settings directory and back up your configuration.
                    </Typography>
                </Box>
            </Box>

            <SyncActionsBar
                handleBackup={handleBackup}
                handleChooseSettingsDir={handleChooseSettingsDir}
                handleResetToDefault={handleResetToDefault}
                isDefaultDir={isDefaultDir}
                isLoading={isLoading}
            />

            {isLoading && (
                <Box display="flex" justifyContent="center" alignItems="center" className="mb-4">
                    <CircularProgress color="primary" />
                </Box>
            )}

            <Grid container spacing={4} className="max-w-7xl mx-auto">
                {settingsData.map(subDir => (
                    <Grid item xs={12} sm={6} md={4} key={subDir.profile}>
                        <SubDirectoryCard
                            subDir={subDir}
                            selections={selections}
                            handleSelectionChange={handleSelectionChange}
                            handleSync={handleSync}
                            handleSyncAll={handleSyncAll}
                            isLoading={isLoading}
                        />
                    </Grid>
                ))}
            </Grid>
            {confirmDialog}
        </div>
    );
};

Sync.propTypes = {
    settingsData: PropTypes.array.isRequired,
    associations: PropTypes.array.isRequired,
    currentSettingsDir: PropTypes.string.isRequired,
    lastBackupDir: PropTypes.string.isRequired,
    userSelections: PropTypes.object.isRequired,
    backEndURL: PropTypes.string.isRequired,
};

export default Sync;
