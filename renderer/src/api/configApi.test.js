// src/api/configApi.test.js
import { vi } from 'vitest';
import { getConfig, updateConfig, saveEVECredentials, checkEVEConfiguration } from './configApi';
import { apiRequest } from './apiClient';

vi.mock('./apiClient', () => ({
    apiRequest: vi.fn()
}));

describe('configApi', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    test('getConfig calls apiRequest correctly', async () => {
        apiRequest.mockResolvedValue('config');
        const result = await getConfig();
        expect(apiRequest).toHaveBeenCalledWith(`/api/config`, {
            method: 'GET',
            credentials: 'include'
        }, {
            errorMessage: 'Failed to fetch configuration.'
        });
        expect(result).toBe('config');
    });

    test('updateConfig calls apiRequest correctly', async () => {
        apiRequest.mockResolvedValue('updated');
        const result = await updateConfig({ a: 1 });
        expect(apiRequest).toHaveBeenCalledWith(`/api/config`, {
            method: 'PATCH',
            headers: { 'Content-Type': 'application/json' },
            credentials: 'include',
            body: JSON.stringify({ a: 1 })
        }, {
            errorMessage: 'Failed to update configuration.'
        });
        expect(result).toBe('updated');
    });

    test('saveEVECredentials calls apiRequest correctly', async () => {
        apiRequest.mockResolvedValue('saved');
        const result = await saveEVECredentials('cid', 'csec');
        expect(apiRequest).toHaveBeenCalledWith(`/api/config/eve/credentials`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            credentials: 'include',
            body: JSON.stringify({ clientId: 'cid', clientSecret: 'csec' })
        }, {
            successMessage: 'EVE credentials saved successfully!',
            errorMessage: 'Failed to save EVE credentials.'
        });
        expect(result).toBe('saved');
    });

    describe('checkEVEConfiguration', () => {
        test('calls /api/config/eve/status and returns the apiRequest result', async () => {
            apiRequest.mockResolvedValue({ needsConfiguration: false });
            const result = await checkEVEConfiguration();
            expect(apiRequest).toHaveBeenCalledWith(`/api/config/eve/status`, {
                method: 'GET',
                credentials: 'include'
            }, {
                errorMessage: 'Failed to check EVE configuration.'
            });
            expect(result).toEqual({ needsConfiguration: false });
        });

        test('returns null on apiRequest failure (apiClient returns null on error)', async () => {
            apiRequest.mockResolvedValue(null);
            const result = await checkEVEConfiguration();
            expect(result).toBeNull();
        });
    });
});
