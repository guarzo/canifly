// src/api/apiService.test.js
import { vi } from 'vitest';
import {
    logout,
    updateCharacter,
    deleteCharacter,
    deleteAccount,
    updateAccount,
    addCharacter,
    saveSkillPlan,
    saveUserSelections,
    syncSubdirectory,
    syncAllSubdirectories,
    chooseSettingsDir,
    backupDirectory,
    resetToDefaultDirectory,
    associateCharacter,
    unassociateCharacter,
    deleteSkillPlan,
    initiateLogin
} from './apiService';
import { apiRequest } from './apiRequest';

vi.mock('./apiRequest', () => ({
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


    describe('updateCharacter', () => {
        test('calls apiRequest correctly', async () => {
            apiRequest.mockResolvedValue('updated');
            const updates = { Role: 'Pvp' };
            const result = await updateCharacter(456, updates);
            expect(apiRequest).toHaveBeenCalledWith('/api/characters/456', {
                method: 'PATCH',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(updates),
                credentials: 'include'
            }, {
                errorMessage: 'Failed to update character.'
            });
            expect(result).toBe('updated');
        });
    });

    describe('deleteCharacter', () => {
        test('calls apiRequest correctly', async () => {
            apiRequest.mockResolvedValue('removed');
            const result = await deleteCharacter(789);
            expect(apiRequest).toHaveBeenCalledWith('/api/characters/789', {
                method: 'DELETE',
                credentials: 'include'
            }, {
                errorMessage: 'Failed to remove character.'
            });
            expect(result).toBe('removed');
        });
    });

    describe('updateAccount', () => {
        test('calls apiRequest correctly for name update', async () => {
            apiRequest.mockResolvedValue('account updated');
            const result = await updateAccount(42, { name: 'NewName' });
            expect(apiRequest).toHaveBeenCalledWith('/api/accounts/42', {
                method: 'PATCH',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ name: 'NewName' }),
                credentials: 'include'
            }, {
                errorMessage: 'Failed to update account.'
            });
            expect(result).toBe('account updated');
        });
    });

    describe('deleteAccount', () => {
        test('calls apiRequest correctly', async () => {
            apiRequest.mockResolvedValue('account removed');
            const result = await deleteAccount(123);
            expect(apiRequest).toHaveBeenCalledWith('/api/accounts/123', {
                method: 'DELETE',
                credentials: 'include'
            }, {
                successMessage: 'Account removed successfully!',
                errorMessage: 'Failed to remove account.'
            });
            expect(result).toBe('account removed');
        });
    });

// In apiService.test.js, update the saveSkillPlan test:
    describe('saveSkillPlan', () => {
        test('calls apiRequest correctly', async () => {
            apiRequest.mockResolvedValue('saved');
            const result = await saveSkillPlan('MyPlan', { skill: 'Level5' });
            expect(apiRequest).toHaveBeenCalledWith('/api/skill-plans', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ name: 'MyPlan', content: { skill: 'Level5' } }),
                credentials: 'include'
            }, {
                // Include successMessage here as well
                successMessage: 'Skill Plan Created!',
                errorMessage: 'Failed to create skill plan.'
            });
            expect(result).toBe('saved');
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

    describe('associateCharacter', () => {
        test('calls apiRequest correctly', async () => {
            apiRequest.mockResolvedValue('associated');
            const result = await associateCharacter('user1', 'char1', 'UserName', 'CharName');
            expect(apiRequest).toHaveBeenCalledWith(`/api/associations`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                credentials: 'include',
                body: JSON.stringify({ userId: 'user1', characterId: 'char1' })
            }, {
                errorMessage: 'Association operation failed.'
            });
            expect(result).toBe('associated');
        });
    });

    describe('unassociateCharacter', () => {
        test('calls apiRequest correctly', async () => {
            apiRequest.mockResolvedValue('unassociated');
            const result = await unassociateCharacter('user1', 'char1', 'UserName', 'CharName');
            expect(apiRequest).toHaveBeenCalledWith(`/api/associations/user1/char1`, {
                method: 'DELETE',
                credentials: 'include'
            }, {
                errorMessage: 'Unassociation operation failed.'
            });
            expect(result).toBe('unassociated');
        });
    });

    describe('deleteSkillPlan', () => {
        test('calls apiRequest correctly', async () => {
            apiRequest.mockResolvedValue('deleted');
            const result = await deleteSkillPlan('MyPlan');
            expect(apiRequest).toHaveBeenCalledWith(`/api/skill-plans/MyPlan`, {
                method: 'DELETE',
                credentials: 'include',
            }, {
                errorMessage: 'Failed to delete skill plan.'
            });
            expect(result).toBe('deleted');
        });
    });

    describe('initiateLogin', () => {
        test('calls apiRequest correctly', async () => {
            apiRequest.mockResolvedValue('login started');
            const account = { Name: 'LoginAccount' };
            const result = await initiateLogin(account);
            expect(apiRequest).toHaveBeenCalledWith(`/api/login`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ account }),
                credentials: 'include',
            }, {
                errorMessage: 'Failed to initiate login.'
            });
            expect(result).toBe('login started');
        });
    });
});
