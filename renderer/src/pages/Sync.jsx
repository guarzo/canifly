// src/components/sync/Sync.jsx

import { useState, useEffect, useCallback } from 'react';
import PropTypes from 'prop-types';
import { toast } from 'react-toastify';
import { useConfirmDialog } from '../hooks/useConfirmDialog.jsx';
import {
    CircularProgress,
    Grid,
    Box,
    Typography
} from '@mui/material';
import SyncActionsBar from '../components/sync/SyncActionsBar.jsx';
import SubDirectoryCard from '../components/sync/SubDirectoryCard.jsx';

import {
    saveUserSelections,
    syncSubdirectory,
    syncAllSubdirectories,
    chooseSettingsDir,
    backupDirectory,
    resetToDefaultDirectory
} from '../api/apiService.jsx';

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

    const saveSelections = useCallback(async (newSelections) => {
        const result = await saveUserSelections(newSelections, backEndURL);
        if (!result || !result.success) {
            // Toast handled by apiRequest if error occurs
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

            saveSelections(updated);
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
            const result = await syncSubdirectory(profile, userId, charId, backEndURL);
            if (result && result.success) {
                toast.success(result.message);
            } else {
                // error handled by apiRequest if any
            }
        } catch (error) {
            console.error('Error syncing:', error);
            // Toast handled by apiRequest or fallback
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
            const result = await syncAllSubdirectories(profile, userId, charId, backEndURL);
            if (result && result.success) {
                toast.success(`Sync-All complete: ${result.message}`);
            } else {
                // error handled by apiRequest if any
            }
        } catch (error) {
            console.error('Error syncing-all:', error);
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

            const result = await chooseSettingsDir(chosenDir, backEndURL);
            if (result && result.success) {
                setIsDefaultDir(false);
                toast.success(`Settings directory chosen: ${chosenDir}`);
            }
        } catch (error) {
            console.error('Error choosing settings directory:', error);
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
            const result = await backupDirectory(currentSettingsDir, chosenDir, backEndURL);
            if (result && result.success) {
                toast.success(result.message);
            }
        } catch (error) {
            console.error('Error during backup:', error);
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
            const result = await resetToDefaultDirectory(backEndURL);
            if (result && result.success) {
                setIsDefaultDir(true);
                toast.success('Directory reset to default: Tranquility');
            }
        } catch (error) {
            console.error('Error resetting directory:', error);
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <div className="bg-gray-900 min-h-screen text-teal-200 px-4 pb-10 pt-16">
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
