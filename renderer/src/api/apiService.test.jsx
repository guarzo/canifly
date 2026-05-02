// src/api/apiService.test.js
//
// Transitional tests for the apiService.jsx façade. Domain-specific
// assertions live in {accounts,skillPlans,sync,esi,fuzzworks,auth,config}Api.test.*.
import { vi } from 'vitest';
import { logout, initiateLogin } from './apiService';
import { apiRequest } from './apiClient';

vi.mock('./apiClient', () => ({
    apiRequest: vi.fn()
}));

describe('apiService façade', () => {
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

    describe('initiateLogin', () => {
        test('defaults rememberMe to false when not provided', async () => {
            apiRequest.mockResolvedValue('login started');
            const result = await initiateLogin('LoginAccount');
            expect(apiRequest).toHaveBeenCalledWith(`/api/login`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ account: 'LoginAccount', rememberMe: false }),
                credentials: 'include',
            }, {
                errorMessage: 'Failed to initiate login.'
            });
            expect(result).toBe('login started');
        });
    });
});
