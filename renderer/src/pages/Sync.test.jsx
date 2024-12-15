import React from 'react';
import { render, screen, fireEvent, act } from '@testing-library/react';
import { vi } from 'vitest';
import Sync from './Sync';
import '@testing-library/jest-dom';

// Mock useConfirmDialog
vi.mock('../hooks/useConfirmDialog.jsx', () => ({
    useConfirmDialog: () => {
        const showConfirmDialog = () => {
            // Always resolve with isConfirmed: true
            return Promise.resolve({ isConfirmed: true });
        };
        const confirmDialog = null;
        return [showConfirmDialog, confirmDialog];
    }
}));

// Mock apiService calls
vi.mock('../api/apiService.jsx', () => ({
    saveUserSelections: vi.fn().mockResolvedValue({ success: true }),
    syncSubdirectory: vi.fn().mockResolvedValue({ success: true, message: 'Synced successfully!' }),
    syncAllSubdirectories: vi.fn().mockResolvedValue({ success: true, message: 'Sync-All successful!' }),
    chooseSettingsDir: vi.fn().mockResolvedValue({ success: true }),
    backupDirectory: vi.fn().mockResolvedValue({ success: true, message: 'Backup complete!' }),
    resetToDefaultDirectory: vi.fn().mockResolvedValue({ success: true }),
}));

import {
    saveUserSelections,
    syncSubdirectory,
    syncAllSubdirectories,
    chooseSettingsDir,
    backupDirectory,
    resetToDefaultDirectory
} from '../api/apiService.jsx';

// Mock electronAPI
window.electronAPI = {
    chooseDirectory: vi.fn().mockResolvedValue('/chosen/dir')
};

describe('Sync component', () => {
    const defaultProps = {
        settingsData: [],
        associations: [],
        currentSettingsDir: '/current/dir',
        userSelections: {},
        lastBackupDir: '',
        backEndURL: 'http://localhost:8713',
    };

    beforeEach(() => {
        vi.clearAllMocks();
    });

    it('renders with no settingsData and shows no subdirectories', () => {
        render(<Sync {...defaultProps} />);
        expect(screen.getByText('Sync Profile Settings')).toBeInTheDocument();
        // No subDir => no SubDirectoryCard
        expect(screen.queryByText('-- Select Character --')).not.toBeInTheDocument();
    });

    it('renders SubDirectoryCards when settingsData is provided', () => {
        const props = {
            ...defaultProps,
            settingsData: [{
                profile: 'settings_profileA',
                availableCharFiles: [
                    { charId: 'char1', name: 'Char One' }
                ],
                availableUserFiles: [
                    { userId: 'userA', name: 'User A' }
                ]
            }],
            userSelections: {
                'settings_profileA': { charId: '', userId: '' }
            }
        };
        render(<Sync {...props} />);

        // The subdir name displayed without 'settings_'
        expect(screen.getByText('profileA')).toBeInTheDocument();
        // Character and user selects
        expect(screen.getByLabelText('Select Character')).toBeInTheDocument();
        expect(screen.getByLabelText('Select User')).toBeInTheDocument();
    });

    it('handles character selection and auto-user selection if associated', () => {
        const associations = [{ userId: 'userA', charId: 'char1', charName: 'Char One' }];
        const props = {
            ...defaultProps,
            associations,
            settingsData: [{
                profile: 'settings_profileA',
                availableCharFiles: [
                    { charId: 'char1', name: 'Char One' }
                ],
                availableUserFiles: [
                    { userId: 'userA', name: 'User A' }
                ]
            }],
            userSelections: {}
        };
        render(<Sync {...props} />);

        const charSelect = screen.getByLabelText('Select Character');
        fireEvent.mouseDown(charSelect);
        fireEvent.click(screen.getByText('Char One'));

        // After selection, user should be auto-selected if assoc found
        expect(saveUserSelections).toHaveBeenCalled();
    });

    it('syncs a single profile', async () => {
        const props = {
            ...defaultProps,
            settingsData: [{
                profile: 'settings_profileA',
                availableCharFiles: [
                    { charId: 'char1', name: 'Char One' }
                ],
                availableUserFiles: [
                    { userId: 'userA', name: 'User A' }
                ]
            }],
            userSelections: {
                'settings_profileA': { charId: 'char1', userId: 'userA' }
            }
        };
        render(<Sync {...props} />);

        const syncButton = screen.getByRole('button', { name: /sync this specific profile/i });
        await act(async () => {
            fireEvent.click(syncButton);
        });

        expect(syncSubdirectory).toHaveBeenCalledWith('settings_profileA', 'userA', 'char1');
        expect(screen.getByText('Synced successfully!')).toBeInTheDocument();
    });

    it('syncs all profiles', async () => {
        const props = {
            ...defaultProps,
            settingsData: [{
                profile: 'settings_profileB',
                availableCharFiles: [
                    { charId: 'char2', name: 'Char Two' }
                ],
                availableUserFiles: [
                    { userId: 'userB', name: 'User B' }
                ]
            }],
            userSelections: {
                'settings_profileB': { charId: 'char2', userId: 'userB' }
            }
        };
        render(<Sync {...props} />);

        const syncAllButton = screen.getByRole('button', { name: /sync all profiles based on this selection/i });
        await act(async () => {
            fireEvent.click(syncAllButton);
        });

        expect(syncAllSubdirectories).toHaveBeenCalledWith('settings_profileB', 'userB', 'char2');
        expect(screen.getByText('Sync-All complete: Sync-All successful!')).toBeInTheDocument();
    });

    it('chooses settings directory', async () => {
        render(<Sync {...defaultProps} />);
        const chooseDirBtn = screen.getByRole('button', { name: /choose settings directory/i });

        await act(async () => {
            fireEvent.click(chooseDirBtn);
        });

        expect(window.electronAPI.chooseDirectory).toHaveBeenCalled();
        expect(chooseSettingsDir).toHaveBeenCalledWith('/chosen/dir');
    });

    it('backup directory chosen', async () => {
        render(<Sync {...defaultProps} />);
        const backupBtn = screen.getByRole('button', { name: /backup settings/i });

        await act(async () => {
            fireEvent.click(backupBtn);
        });

        expect(window.electronAPI.chooseDirectory).toHaveBeenCalledWith('');
        expect(backupDirectory).toHaveBeenCalledWith('/current/dir', '/chosen/dir');
        expect(screen.getByText('Backup complete!')).toBeInTheDocument();
    });

    it('reset to default directory', async () => {
        render(<Sync {...defaultProps} />);
        const resetBtn = screen.getByRole('button', { name: /reset to default directory/i });

        await act(async () => {
            fireEvent.click(resetBtn);
        });

        expect(resetToDefaultDirectory).toHaveBeenCalledWith();
        expect(screen.getByText('Directory reset to default: Tranquility')).toBeInTheDocument();
    });
});
