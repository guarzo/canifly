// src/api/api.test.js
import { vi } from 'vitest';
import { fetchAppEndpoint } from './api';
import { toast } from 'react-toastify';

vi.mock('react-toastify', () => ({
    toast: {
        warning: vi.fn(),
        error: vi.fn(),
    },
}));

describe('fetchAppEndpoint', () => {
    const originalLog = console.log;
    const originalError = console.error;

    beforeAll(() => {
        // Mock console.log and console.error to silence them
        console.log = vi.fn();
        console.error = vi.fn();
    });

    afterAll(() => {
        // Restore the original console methods after all tests
        console.log = originalLog;
        console.error = originalError;
    });

    const backEndURL = 'http://test-backend';
    const endpoint = '/api/test';

    let setIsLoading, setIsRefreshing, setIsAuthenticated, setAppData;

    beforeEach(() => {
        setIsLoading = vi.fn();
        setIsRefreshing = vi.fn();
        setIsAuthenticated = vi.fn();
        setAppData = vi.fn();
        global.fetch = vi.fn();
        vi.clearAllMocks();
    });

    test('returns early if loggedOut is true', async () => {
        const result = await fetchAppEndpoint({
            backEndURL,
            endpoint,
            loggedOut: true,
            isAuthenticated: false,
            setIsLoading,
            setIsRefreshing,
            setIsAuthenticated,
            setAppData
        }, { returnSuccess: true });

        expect(global.fetch).not.toHaveBeenCalled();
        expect(setIsLoading).not.toHaveBeenCalled();
        expect(setIsRefreshing).not.toHaveBeenCalled();
        expect(setIsAuthenticated).not.toHaveBeenCalled();
        expect(setAppData).not.toHaveBeenCalled();

        expect(result).toBe(false);
    });

    test('handles 401 unauthorized', async () => {
        global.fetch.mockResolvedValue({
            status: 401,
            ok: false,
            json: vi.fn().mockResolvedValue({ error: 'Unauthorized' })
        });

        const result = await fetchAppEndpoint({
            backEndURL,
            endpoint,
            loggedOut: false,
            isAuthenticated: true,
            setIsLoading,
            setIsRefreshing,
            setIsAuthenticated,
            setAppData
        }, { setLoading: true, returnSuccess: true });

        expect(setIsLoading).toHaveBeenCalledWith(true);
        expect(setIsLoading).toHaveBeenCalledWith(false);
        expect(setIsRefreshing).not.toHaveBeenCalled();
        expect(setIsAuthenticated).toHaveBeenCalledWith(false);
        expect(setAppData).toHaveBeenCalledWith(null);
        expect(toast.warning).toHaveBeenCalledWith('Please log in to access your data.');

        expect(result).toBe(false);
    });

    test('handles successful response', async () => {
        const mockData = { foo: 'bar' };
        global.fetch.mockResolvedValue({
            status: 200,
            ok: true,
            json: vi.fn().mockResolvedValue(mockData)
        });

        const result = await fetchAppEndpoint({
            backEndURL,
            endpoint,
            loggedOut: false,
            isAuthenticated: false,
            setIsLoading,
            setIsRefreshing,
            setIsAuthenticated,
            setAppData
        }, { setLoading: true, setRefreshing: true, returnSuccess: true });

        expect(setIsLoading).toHaveBeenCalledWith(true);
        expect(setIsRefreshing).toHaveBeenCalledWith(true);
        expect(setIsLoading).toHaveBeenCalledWith(false);
        expect(setIsRefreshing).toHaveBeenCalledWith(false);

        expect(setIsAuthenticated).toHaveBeenCalledWith(true);
        expect(setAppData).toHaveBeenCalledWith(mockData);

        expect(toast.warning).not.toHaveBeenCalled();
        expect(toast.error).not.toHaveBeenCalled();

        expect(result).toBe(true);
    });

    test('handles error response (non-2xx)', async () => {
        const mockError = 'Some error occurred.';
        global.fetch.mockResolvedValue({
            status: 500,
            ok: false,
            json: vi.fn().mockResolvedValue({ error: mockError })
        });

        const result = await fetchAppEndpoint({
            backEndURL,
            endpoint,
            loggedOut: false,
            isAuthenticated: true,
            setIsLoading,
            setIsRefreshing,
            setIsAuthenticated,
            setAppData
        }, { returnSuccess: true });

        expect(setIsLoading).not.toHaveBeenCalled();
        expect(setIsRefreshing).not.toHaveBeenCalled();
        expect(setIsAuthenticated).not.toHaveBeenCalled();
        expect(setAppData).not.toHaveBeenCalled();

        expect(toast.error).toHaveBeenCalledWith(mockError);

        expect(result).toBe(false);
    });

    test('handles network error', async () => {
        global.fetch.mockRejectedValue(new Error('Network Error'));

        const result = await fetchAppEndpoint({
            backEndURL,
            endpoint,
            loggedOut: false,
            isAuthenticated: true,
            setIsLoading,
            setIsRefreshing,
            setIsAuthenticated,
            setAppData
        }, { setLoading: true, returnSuccess: true });

        expect(setIsLoading).toHaveBeenCalledWith(true);
        expect(setIsLoading).toHaveBeenCalledWith(false);

        expect(setIsRefreshing).not.toHaveBeenCalled();

        expect(toast.error).toHaveBeenCalledWith('Failed to load data. Please try again later.');

        expect(setIsAuthenticated).not.toHaveBeenCalled();
        expect(setAppData).not.toHaveBeenCalled();

        expect(result).toBe(false);
    });
});
