// src/api/apiService.test.js
import { vi } from 'vitest';
import {
    logout,
    saveUserSelections,
    syncSubdirectory,
    syncAllSubdirectories,
    chooseSettingsDir,
    backupDirectory,
    resetToDefaultDirectory,
    initiateLogin
} from './apiService';
import { apiRequest } from './apiClient';

vi.mock('./apiClient', () => ({
    apiRequest: vi.fn()
}));

describe('apiService', () => {
    const backEndURL = 'http://backend.test';

    beforeEach(() => {
        vi.clearAllMocks();
    });


    describe('logout', () => {
        test('calls apiRequest with correct parameters', async () => {
            apiRequest.mockResolvedValue('success');
            const result = await logout();
            expect(apiRequest).toHaveBeenCalledWith('/api/logout', {
                method: 'POST',
                credentials: 'include'
            }, {
                errorMessage: 'Failed to log out.'
            });
            expect(result).toBe('success');
        });
    });


    describe('saveUserSelections', () => {
        test('calls apiRequest correctly', async () => {
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
    });

    describe('syncSubdirectory', () => {
        test('calls apiRequest correctly', async () => {
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
    });

    describe('syncAllSubdirectories', () => {
        test('calls apiRequest correctly', async () => {
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
    });

    describe('chooseSettingsDir', () => {
        test('calls apiRequest correctly', async () => {
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
    });

    describe('backupDirectory', () => {
        test('calls apiRequest correctly', async () => {
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
    });

    describe('resetToDefaultDirectory', () => {
        test('calls apiRequest correctly', async () => {
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

    describe('initiateLogin', () => {
        test('calls apiRequest correctly', async () => {
            apiRequest.mockResolvedValue('login started');
            const account = 'LoginAccount';
            const rememberMe = false;
            const result = await initiateLogin(account, rememberMe);
            expect(apiRequest).toHaveBeenCalledWith(`/api/login`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ account, rememberMe }),
                credentials: 'include',
            }, {
                errorMessage: 'Failed to initiate login.'
            });
            expect(result).toBe('login started');
        });

        test('calls apiRequest with rememberMe true', async () => {
            apiRequest.mockResolvedValue('login started');
            const account = 'LoginAccount';
            const rememberMe = true;
            const result = await initiateLogin(account, rememberMe);
            expect(apiRequest).toHaveBeenCalledWith(`/api/login`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ account, rememberMe }),
                credentials: 'include',
            }, {
                errorMessage: 'Failed to initiate login.'
            });
            expect(result).toBe('login started');
        });

        test('defaults rememberMe to false when not provided', async () => {
            apiRequest.mockResolvedValue('login started');
            const account = 'LoginAccount';
            const result = await initiateLogin(account);
            expect(apiRequest).toHaveBeenCalledWith(`/api/login`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ account, rememberMe: false }),
                credentials: 'include',
            }, {
                errorMessage: 'Failed to initiate login.'
            });
            expect(result).toBe('login started');
        });
    });
});
