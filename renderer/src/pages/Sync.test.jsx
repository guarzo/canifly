import React from 'react';
import { render, screen, fireEvent, act } from '@testing-library/react';
import { vi } from 'vitest';
import Sync from './Sync';
import '@testing-library/jest-dom';

// Mock useConfirmDialog — always confirm.
vi.mock('../hooks/useConfirmDialog.jsx', () => ({
    useConfirmDialog: () => {
        const showConfirmDialog = () => Promise.resolve({ isConfirmed: true });
        return [showConfirmDialog, null];
    }
}));

// Mock react-toastify so we can assert on toast calls (messages aren't in the DOM).
vi.mock('react-toastify', () => ({
    toast: {
        success: vi.fn(),
        error: vi.fn(),
        info: vi.fn(),
    },
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

import { toast } from 'react-toastify';
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
        backEndURL: 'http://localhost:42423',
    };

    beforeEach(() => {
        vi.clearAllMocks();
    });

    it('renders with no settingsData and shows empty state', () => {
        render(<Sync {...defaultProps} />);
        // Subheader title is "Sync"
        expect(screen.getByText('Sync')).toBeInTheDocument();
        // Empty-state message
        expect(screen.getByText(/No EVE profiles found/i)).toBeInTheDocument();
    });

    it('renders SyncProfileRow when settingsData is provided', () => {
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

        // Profile name (without 'settings_')
        expect(screen.getByText('profileA')).toBeInTheDocument();
        // Per-row select inputs are labeled with the profile name.
        expect(screen.getByLabelText('Character file for profileA')).toBeInTheDocument();
        expect(screen.getByLabelText('User file for profileA')).toBeInTheDocument();
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
            userSelections: {
                'settings_profileA': { charId: '', userId: '' }
            }
        };
        render(<Sync {...props} />);

        const charSelect = screen.getByLabelText('Character file for profileA');
        fireEvent.mouseDown(charSelect);
        fireEvent.click(screen.getByText('Char One'));

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

        const syncButton = screen.getByRole('button', { name: /^sync profileA$/i });
        await act(async () => {
            fireEvent.click(syncButton);
        });

        expect(syncSubdirectory).toHaveBeenCalledWith('settings_profileA', 'userA', 'char1');
        expect(toast.success).toHaveBeenCalledWith('Synced successfully!');
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

        const syncAllButton = screen.getByRole('button', { name: /sync all profiles using profileB/i });
        await act(async () => {
            fireEvent.click(syncAllButton);
        });

        expect(syncAllSubdirectories).toHaveBeenCalledWith('settings_profileB', 'userB', 'char2');
        expect(toast.success).toHaveBeenCalledWith('Sync-All complete: Sync-All successful!');
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
        expect(toast.success).toHaveBeenCalledWith('Backup complete!');
    });

    it('reset to default directory', async () => {
        render(<Sync {...defaultProps} />);
        const resetBtn = screen.getByRole('button', { name: /reset to default directory/i });

        await act(async () => {
            fireEvent.click(resetBtn);
        });

        expect(resetToDefaultDirectory).toHaveBeenCalledWith();
        expect(toast.success).toHaveBeenCalledWith('Reset to default: Tranquility');
    });
});
