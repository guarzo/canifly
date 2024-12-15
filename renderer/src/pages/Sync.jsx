// src/components/sync/Sync.jsx

import React, { useState, useEffect, useCallback } from 'react';
import PropTypes from 'prop-types';
import { toast } from 'react-toastify';
import { useConfirmDialog } from '../hooks/useConfirmDialog.jsx';
import {
    CircularProgress,
    Grid,
    Box,
    Typography,
} from '@mui/material';

import SyncActionsBar from '../components/sync/SyncActionsBar.jsx';
import SubDirectoryCard from '../components/sync/SubDirectoryCard.jsx';
import {syncInstructions} from '../utils/instructions.jsx';

import {
    saveUserSelections,
    syncSubdirectory,
    syncAllSubdirectories,
    chooseSettingsDir,
    backupDirectory,
    resetToDefaultDirectory
} from '../api/apiService.jsx';
import PageHeader from "../components/common/SubPageHeader.jsx";

const Sync = ({
                  settingsData,
                  associations,
                  currentSettingsDir,
                  userSelections,
                  lastBackupDir,
              }) => {
    const [isLoading, setIsLoading] = useState(false);
    const [selections, setSelections] = useState({});
    const [showConfirmDialog, confirmDialog] = useConfirmDialog();
    const [isDefaultDir, setIsDefaultDir] = useState(false);
    const [message, setMessage] = useState('');

    // Load instruction visibility from localStorage
    const [showInstructions, setShowInstructions] = useState(() => {
        const stored = localStorage.getItem('showSyncInstructions');
        return stored === null ? true : JSON.parse(stored);
    });

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

    const saveSelectionsCallback = useCallback(async (newSelections) => {
        const result = await saveUserSelections(newSelections);
        if (!result || !result.success) {
            // Errors handled by apiRequest/toast internally
        }
    }, []);

    const handleSelectionChange = (profile, field, value) => {
        setSelections(prev => {
            const updated = {
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

            saveSelectionsCallback(updated);
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
            const result = await syncSubdirectory(profile, userId, charId);
            if (result && result.success) {
                toast.success(result.message);
                setMessage('Synced successfully!');
            }
        } catch (error) {
            console.error('Error syncing:', error);
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
            const result = await syncAllSubdirectories(profile, userId, charId);
            if (result && result.success) {
                toast.success(`Sync-All complete: ${result.message}`);
                setMessage(`Sync-All complete: ${result.message}`);
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

            const result = await chooseSettingsDir(chosenDir);
            if (result && result.success) {
                setIsDefaultDir(false);
                toast.success(`Settings directory chosen: ${chosenDir}`);
                setMessage('Settings directory chosen!');
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

            const chosenDir = await window.electronAPI.chooseDirectory(lastBackupDir || '');
            if (!chosenDir) {
                toast.info('No backup directory chosen. Backup canceled.');
                setIsLoading(false);
                return;
            }

            toast.info('Starting backup...');
            const result = await backupDirectory(currentSettingsDir, chosenDir);
            if (result && result.success) {
                toast.success(result.message);
                setMessage('Backup complete!');
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
            const result = await resetToDefaultDirectory();
            if (result && result.success) {
                setIsDefaultDir(true);
                toast.success('Directory reset to default: Tranquility');
                setMessage('Directory reset to default: Tranquility');
            }
        } catch (error) {
            console.error('Error resetting directory:', error);
        } finally {
            setIsLoading(false);
        }
    };

    const toggleInstructions = () => {
        const newValue = !showInstructions;
        setShowInstructions(newValue);
        localStorage.setItem('showSyncInstructions', JSON.stringify(newValue));
    };

    return (
        <div className="bg-gray-900 min-h-screen text-teal-200 px-4 pb-10 pt-16">
            <PageHeader
                title="Sync Profile Settings"
                instructions={syncInstructions}
                storageKey="showSyncInstructions"
            />
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

            {message && (
                <Box className="max-w-7xl mx-auto mt-4">
                    <Typography>{message}</Typography>
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
};

export default Sync;
