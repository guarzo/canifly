// src/api/apiRequest.test.jsx
import { vi } from 'vitest';
import { apiRequest } from './apiRequest';
import { toast } from 'react-toastify';

vi.mock('react-toastify', () => ({
    toast: {
        success: vi.fn(),
        error: vi.fn(),
    }
}));

describe('apiRequest', () => {
    const originalLog = console.log;
    const originalError = console.error;

    beforeAll(() => {
        // Mock console.log and console.error to suppress logs
        console.log = vi.fn();
        console.error = vi.fn();
    });

    afterAll(() => {
        // Restore original console methods after tests
        console.log = originalLog;
        console.error = originalError;
    });

    let onSuccess, onError;

    beforeEach(() => {
        onSuccess = vi.fn();
        onError = vi.fn();
        vi.clearAllMocks();
        global.fetch = vi.fn();
    });

    test('handles successful JSON response with successMessage and onSuccess', async () => {
        const mockData = { foo: 'bar' };
        global.fetch.mockResolvedValue({
            ok: true,
            headers: { get: () => 'application/json' },
            json: vi.fn().mockResolvedValue(mockData)
        });

        const result = await apiRequest('http://test-url', {}, {
            onSuccess,
            successMessage: 'Success!',
        });

        expect(fetch).toHaveBeenCalledWith('http://test-url', {});
        expect(toast.success).toHaveBeenCalledWith('Success!');
        expect(onSuccess).toHaveBeenCalledWith(mockData);
        expect(result).toEqual(mockData);
        expect(toast.error).not.toHaveBeenCalled();
        expect(onError).not.toHaveBeenCalled();
    });

    test('handles successful text response', async () => {
        const mockText = 'Some plain text';
        global.fetch.mockResolvedValue({
            ok: true,
            headers: { get: () => 'text/plain' },
            text: vi.fn().mockResolvedValue(mockText)
        });

        const result = await apiRequest('http://test-url', {}, {
            onSuccess,
            successMessage: 'It worked!'
        });

        expect(toast.success).toHaveBeenCalledWith('It worked!');
        expect(onSuccess).toHaveBeenCalledWith(mockText);
        expect(result).toEqual(mockText);
        expect(toast.error).not.toHaveBeenCalled();
        expect(onError).not.toHaveBeenCalled();
    });

    test('handles error JSON response with given errorMessage and onError', async () => {
        const mockError = { error: 'Something went wrong' };
        global.fetch.mockResolvedValue({
            ok: false,
            headers: { get: () => 'application/json' },
            json: vi.fn().mockResolvedValue(mockError)
        });

        await apiRequest('http://test-url', {}, {
            onError,
            errorMessage: 'Custom error message'
        });

        expect(onError).toHaveBeenCalledWith('Something went wrong');
        expect(toast.success).not.toHaveBeenCalled();
        expect(onSuccess).not.toHaveBeenCalled();
    });

    test('handles error non-JSON response with fallback errorMessage', async () => {
        global.fetch.mockResolvedValue({
            ok: false,
            headers: { get: () => 'text/html' },
            text: vi.fn().mockResolvedValue('Some non-JSON error')
        });

        await apiRequest('http://test-url', {}, {
            onError,
            errorMessage: 'A fallback error'
        });

        expect(onError).toHaveBeenCalledWith('A fallback error');
        expect(toast.success).not.toHaveBeenCalled();
        expect(onSuccess).not.toHaveBeenCalled();
    });

    test('handles network error', async () => {
        global.fetch.mockRejectedValue(new Error('Network failed'));

        await apiRequest('http://test-url', {}, {
            onError,
            errorMessage: 'Request failed'
        });

        expect(onError).toHaveBeenCalledWith('Network failed');
        expect(toast.success).not.toHaveBeenCalled();
        expect(onSuccess).not.toHaveBeenCalled();
    });
});
