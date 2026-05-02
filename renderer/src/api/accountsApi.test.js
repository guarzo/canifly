// src/api/accountsApi.test.js
import { vi } from 'vitest';
import {
    updateAccount,
    deleteAccount,
    updateCharacter,
    deleteCharacter,
    associateCharacter,
    unassociateCharacter,
} from './accountsApi';
import { apiRequest } from './apiClient';

vi.mock('./apiClient', () => ({
    apiRequest: vi.fn()
}));

describe('accountsApi', () => {
    beforeEach(() => {
        vi.clearAllMocks();
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
});
