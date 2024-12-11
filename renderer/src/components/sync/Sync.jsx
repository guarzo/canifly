import { useState, useEffect, useCallback } from 'react';
import PropTypes from 'prop-types';
import { toast } from 'react-toastify';
import { useConfirmDialog } from '../partials/useConfirmDialog.jsx';
import {
    CircularProgress,
    Grid,
    Box,
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
    const [isDefaultDir, setIsDefaultDir] = useState(true)


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

    // Helper: Given a subDir name and userId, find the user's name from availableUserFiles
    const getUserInfo = (subDirName, userId) => {
        const subDir = settingsData.find(s => s.subDir === subDirName);
        if (!subDir) return { userName: 'Unknown', userId };
        const userFile = subDir.availableUserFiles.find(u => u.userId === userId);
        if (!userFile) return { userName: 'Unknown', userId };
        return { userName: userFile.name, userId };
    };

    // Helper: Given a subDir name and charId, find the char's name from availableCharFiles
    const getCharacterInfo = (subDirName, charId) => {
        const subDir = settingsData.find(s => s.subDir === subDirName);
        if (!subDir) return { charName: 'Unknown', charId };
        const charFile = subDir.availableCharFiles.find(c => c.charId === charId);
        if (!charFile) {
            // fallback to associations if not found
            const assoc = associations.find(a => a.charId === charId);
            return { charName: assoc ? assoc.charName : 'Unknown', charId };
        }
        return { charName: charFile.name, charId };
    };

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
    }, []);

    const handleSelectionChange = (subDir, field, value) => {
        setSelections(prev => {
            let updated = {
                ...prev,
                [subDir]: {
                    ...prev[subDir],
                    [field]: value,
                }
            };

            // If we just changed the charId, auto-select the user if we can find it in associations
            if (field === 'charId' && value) {
                const assoc = associations.find(a => a.charId === value);
                if (assoc) {
                    // If we found a matching association, update the userId as well
                    updated[subDir].userId = assoc.userId;
                }
            }

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

        const { charName } = getCharacterInfo(subDir, charId);
        const { userName } = getUserInfo(subDir, userId);
        const displaySubDir = subDir.replace('settings_', '');

        const confirmSync = await showConfirmDialog({
            title: 'Confirm Sync',
            message: `Use "${charName}" (${charId}) on account "${userName}" (${userId}) to overwrite all files in profile "${displaySubDir}"?`,
        });

        if (!confirmSync.isConfirmed) return;

        try {
            setIsLoading(true);
            toast.info('Syncing...', { autoClose: 1500 });
            const response = await fetch(`${backEndURL}/api/sync-subdirectory`, {
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

        const { charName } = getCharacterInfo(subDir, charId);
        const { userName } = getUserInfo(subDir, userId);
        const displaySubDir = subDir.replace('settings_', '');

        const confirmSyncAll = await showConfirmDialog({
            title: 'Confirm Sync All',
            message: `Use "${charName}" (${charId}) on account "${userName}" (${userId}) to overwrite files across all profiles (based on "${displaySubDir}")?`,
        });

        if (!confirmSyncAll.isConfirmed) return;

        try {
            setIsLoading(true);
            const response = await fetch(`${backEndURL}/api/sync-all-subdirectories`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                credentials: 'include',
                body: JSON.stringify({ subDir, userId, charId }),
            });
            const result = await response.json();
            if (response.ok && result.success) {
                toast.success(`Sync-All complete using "${charName}" (${charId}) on account "${userName}" (${userId}): ${result.message}`);
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
                setIsDefaultDir(false)
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
            message: 'Are you sure you want to reset the settings directory to default (Tranquility)?',
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
                setIsDefaultDir(true)
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

    return (
        <div className="bg-gray-900 min-h-screen text-teal-200 px-4 pb-10">
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

            <Grid container spacing={4}>
                {settingsData.map(subDir => (
                    <Grid item xs={12} sm={6} md={4} key={subDir.subDir}>
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