// src/api/apiService.test.js
import { vi } from 'vitest';
import {
    getAppData,
    getAppDataNoCache,
    logout,
    toggleAccountStatus,
    updateCharacter,
    removeCharacter,
    updateAccountName,
    removeAccount,
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
import { normalizeAppData } from '../utils/dataNormalizer';

vi.mock('./apiRequest', () => ({
    apiRequest: vi.fn()
}));
vi.mock('../utils/dataNormalizer', () => ({
    normalizeAppData: vi.fn()
}));

describe('apiService', () => {
    const backEndURL = 'http://backend.test';

    beforeEach(() => {
        vi.clearAllMocks();
    });

    describe('getAppData', () => {
        test('returns normalized data on success', async () => {
            const mockData = { foo: 'bar' };
            const normalizedData = { foo: 'normalized' };
            apiRequest.mockResolvedValue(mockData);
            normalizeAppData.mockReturnValue(normalizedData);

            const result = await getAppData();
            // Depending on isDev, errorMessage may be present or not. Let's assume isDev is true for testing.
            expect(apiRequest).toHaveBeenCalledWith(`/api/app-data`, { credentials: 'include' }, { errorMessage: 'Failed to load app data.' });
            expect(normalizeAppData).toHaveBeenCalledWith(mockData);
            expect(result).toBe(normalizedData);
        });

        test('returns null if no response', async () => {
            apiRequest.mockResolvedValue(null);
            const result = await getAppData();
            expect(result).toBeNull();
        });
    });

    describe('getAppDataNoCache', () => {
        test('returns normalized data on success', async () => {
            const mockData = { baz: 'qux' };
            const normalizedData = { baz: 'normalized' };
            apiRequest.mockResolvedValue(mockData);
            normalizeAppData.mockReturnValue(normalizedData);

            const result = await getAppDataNoCache();
            expect(apiRequest).toHaveBeenCalledWith(`/api/app-data-no-cache`, { credentials: 'include' }, { errorMessage: 'Failed to load data.' });
            expect(normalizeAppData).toHaveBeenCalledWith(mockData);
            expect(result).toBe(normalizedData);
        });

        test('returns null if no response', async () => {
            apiRequest.mockResolvedValue(null);
            const result = await getAppDataNoCache();
            expect(result).toBeNull();
        });
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

    describe('toggleAccountStatus', () => {
        test('calls apiRequest correctly', async () => {
            apiRequest.mockResolvedValue('toggled');
            const result = await toggleAccountStatus(123);
            expect(apiRequest).toHaveBeenCalledWith('/api/toggle-account-status', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                credentials: 'include',
                body: JSON.stringify({ accountID: 123 })
            }, {
                errorMessage: 'Failed to toggle account status.'
            });
            expect(result).toBe('toggled');
        });
    });

    describe('updateCharacter', () => {
        test('calls apiRequest correctly', async () => {
            apiRequest.mockResolvedValue('updated');
            const updates = { Role: 'Pvp' };
            const result = await updateCharacter(456, updates);
            expect(apiRequest).toHaveBeenCalledWith('/api/update-character', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ characterID: 456, updates }),
                credentials: 'include'
            }, {
                errorMessage: 'Failed to update character.'
            });
            expect(result).toBe('updated');
        });
    });

    describe('removeCharacter', () => {
        test('calls apiRequest correctly', async () => {
            apiRequest.mockResolvedValue('removed');
            const result = await removeCharacter(789);
            expect(apiRequest).toHaveBeenCalledWith('/api/remove-character', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ characterID: 789 }),
                credentials: 'include'
            }, {
                errorMessage: 'Failed to remove character.'
            });
            expect(result).toBe('removed');
        });
    });

    describe('updateAccountName', () => {
        test('calls apiRequest correctly', async () => {
            apiRequest.mockResolvedValue('name updated');
            const result = await updateAccountName(42, 'NewName');
            expect(apiRequest).toHaveBeenCalledWith('/api/update-account-name', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ accountID: 42, accountName: 'NewName' }),
                credentials: 'include'
            }, {
                errorMessage: 'Failed to update account name.'
            });
            expect(result).toBe('name updated');
        });
    });

    // In apiService.test.js, update the removeAccount test:
    describe('removeAccount', () => {
        test('calls apiRequest correctly', async () => {
            apiRequest.mockResolvedValue('account removed');
            const result = await removeAccount('TestAccount');
            expect(apiRequest).toHaveBeenCalledWith('/api/remove-account', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ accountName: 'TestAccount' }),
                credentials: 'include'
            }, {
                // Include successMessage since the code now has it
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
            expect(apiRequest).toHaveBeenCalledWith('/api/save-skill-plan', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ name: 'MyPlan', contents: { skill: 'Level5' } }),
                credentials: 'include'
            }, {
                // Include successMessage here as well
                successMessage: 'Skill Plan Saved!',
                errorMessage: 'Failed to save skill plan.'
            });
            expect(result).toBe('saved');
        });
    });


    describe('addCharacter', () => {
        test('calls apiRequest correctly', async () => {
            apiRequest.mockResolvedValue({ redirectURL: 'http://redirect.url' });
            window.isDev = true;

            // Mock window.location
            delete window.location;
            window.location = { href: 'http://localhost:3000/' };

            const account = { Name: 'MyAccount' };
            const result = await addCharacter(account);

            expect(apiRequest).toHaveBeenCalledWith('/api/add-character', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ account }),
                credentials: 'include'
            }, {
                errorMessage: 'An error occurred while adding character.',
                onSuccess: expect.any(Function)
            });

            const onSuccess = apiRequest.mock.calls[0][2].onSuccess;
            onSuccess({ redirectURL: 'http://redirect.url' });

            // Now we can check the mocked location
            expect(window.location.href).toBe('http://redirect.url');

            expect(result).toEqual({ redirectURL: 'http://redirect.url' });
        });
    });

    describe('saveUserSelections', () => {
        test('calls apiRequest correctly', async () => {
            apiRequest.mockResolvedValue('selections saved');
            const newSelections = { theme: 'dark' };
            const result = await saveUserSelections(newSelections);
            expect(apiRequest).toHaveBeenCalledWith(`/api/save-user-selections`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                credentials: 'include',
                body: JSON.stringify(newSelections),
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
                body: JSON.stringify({ profile: 'profile1', userId: 'user123', charId: 'char456' })
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
                body: JSON.stringify({ profile: 'profile1', userId: 'user123', charId: 'char456' })
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
            expect(apiRequest).toHaveBeenCalledWith(`/api/choose-settings-dir`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                credentials: 'include',
                body: JSON.stringify({ directory: '/path/to/dir' }),
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
            expect(apiRequest).toHaveBeenCalledWith(`/api/reset-to-default-directory`, {
                method: 'POST',
                credentials: 'include',
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
            expect(apiRequest).toHaveBeenCalledWith(`/api/associate-character`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                credentials: 'include',
                body: JSON.stringify({ userId: 'user1', charId: 'char1', userName: 'UserName', charName: 'CharName' })
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
            expect(apiRequest).toHaveBeenCalledWith(`/api/unassociate-character`, {
                method: 'POST',
                headers: {'Content-Type': 'application/json'},
                credentials: 'include',
                body: JSON.stringify({ userId: 'user1', charId: 'char1', userName: 'UserName', charName: 'CharName' })
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
            expect(apiRequest).toHaveBeenCalledWith(`/api/delete-skill-plan?planName=MyPlan`, {
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
