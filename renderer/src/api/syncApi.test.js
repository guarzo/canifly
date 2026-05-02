// src/api/syncApi.test.js
import { vi } from 'vitest';
import {
    saveUserSelections,
    syncSubdirectory,
    syncAllSubdirectories,
    chooseSettingsDir,
    backupDirectory,
    resetToDefaultDirectory,
} from './syncApi';
import { apiRequest } from './apiClient';

vi.mock('./apiClient', () => ({
    apiRequest: vi.fn()
}));

describe('syncApi', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    test('saveUserSelections calls apiRequest correctly', async () => {
        apiRequest.mockResolvedValue('selections saved');
        const newSelections = { theme: 'dark' };
        const result = await saveUserSelections(newSelections);
        expect(apiRequest).toHaveBeenCalledWith(`/api/config`, {
            method: 'PATCH',
            headers: { 'Content-Type': 'application/json' },
            credentials: 'include',
            body: JSON.stringify({ userSelections: newSelections }),
        }, {
            errorMessage: 'Failed to save user selections.',
        });
        expect(result).toBe('selections saved');
    });

    test('syncSubdirectory calls apiRequest correctly', async () => {
        apiRequest.mockResolvedValue('synced');
        const result = await syncSubdirectory('profile1', 'user123', 'char456');
        expect(apiRequest).toHaveBeenCalledWith(`/api/sync-subdirectory`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            credentials: 'include',
            body: JSON.stringify({ subDir: 'profile1', userId: 'user123', charId: 'char456' })
        }, {
            errorMessage: 'Sync operation failed.'
        });
        expect(result).toBe('synced');
    });

    test('syncAllSubdirectories calls apiRequest correctly', async () => {
        apiRequest.mockResolvedValue('all synced');
        const result = await syncAllSubdirectories('profile1', 'user123', 'char456');
        expect(apiRequest).toHaveBeenCalledWith(`/api/sync-all-subdirectories`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            credentials: 'include',
            body: JSON.stringify({ subDir: 'profile1', userId: 'user123', charId: 'char456' })
        }, {
            errorMessage: 'Sync-All operation failed.'
        });
        expect(result).toBe('all synced');
    });

    test('chooseSettingsDir calls apiRequest correctly', async () => {
        apiRequest.mockResolvedValue('chosen');
        const result = await chooseSettingsDir('/path/to/dir');
        expect(apiRequest).toHaveBeenCalledWith(`/api/config`, {
            method: 'PATCH',
            headers: { 'Content-Type': 'application/json' },
            credentials: 'include',
            body: JSON.stringify({ settingsDir: '/path/to/dir' }),
        }, {
            errorMessage: 'Failed to choose settings directory.'
        });
        expect(result).toBe('chosen');
    });

    test('backupDirectory calls apiRequest correctly', async () => {
        apiRequest.mockResolvedValue('backed up');
        const result = await backupDirectory('/target/dir', '/backup/dir');
        expect(apiRequest).toHaveBeenCalledWith(`/api/backup-directory`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            credentials: 'include',
            body: JSON.stringify({ targetDir: '/target/dir', backupDir: '/backup/dir' }),
        }, {
            errorMessage: 'Backup operation failed.'
        });
        expect(result).toBe('backed up');
    });

    test('resetToDefaultDirectory calls apiRequest correctly', async () => {
        apiRequest.mockResolvedValue('reset');
        const result = await resetToDefaultDirectory();
        expect(apiRequest).toHaveBeenCalledWith(`/api/config`, {
            method: 'PATCH',
            headers: { 'Content-Type': 'application/json' },
            credentials: 'include',
            body: JSON.stringify({ settingsDir: '' }),
        }, {
            errorMessage: 'Failed to reset directory.'
        });
        expect(result).toBe('reset');
    });
});
