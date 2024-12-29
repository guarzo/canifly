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
            json: vi.fn().mockResolvedValue(mockData),
        });

        const result = await apiRequest('http://test-url', {}, {
            onSuccess,
            successMessage: 'Success!',
        });

        expect(fetch).toHaveBeenCalledWith('http://test-url', {});
        expect(toast.success).toHaveBeenCalledWith('Success!');
        expect(onSuccess).toHaveBeenCalledWith(mockData);
        expect(result).toEqual(mockData);

        // Should NOT call error toast or onError
        expect(toast.error).not.toHaveBeenCalled();
        expect(onError).not.toHaveBeenCalled();
    });

    test('handles successful text response', async () => {
        const mockText = 'Some plain text';
        global.fetch.mockResolvedValue({
            ok: true,
            headers: { get: () => 'text/plain' },
            text: vi.fn().mockResolvedValue(mockText),
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

    // This test specifically for 401 errors (since your code only calls onError if status === 401)
    test('handles 401 JSON error response with given errorMessage and onError', async () => {
        const mockError = { error: 'Unauthorized' };
        global.fetch.mockResolvedValue({
            ok: false,
            status: 401, // <--- important!
            headers: { get: () => 'application/json' },
            json: vi.fn().mockResolvedValue(mockError),
        });

        await apiRequest('http://test-url', {}, {
            onError,
            errorMessage: 'Custom error message for 401'
        });

        // Check that the code calls toast.error with "Unauthorized"
        expect(toast.error).toHaveBeenCalledWith('Unauthorized');
        // Check that onError is called with the same message
        expect(onError).toHaveBeenCalledWith('Unauthorized');

        // Should NOT call toast.success or onSuccess
        expect(toast.success).not.toHaveBeenCalled();
    });

    // New test to confirm non-401 errors do NOT trigger onError or toast
    // because your code returns early if (status !== 401).
    test('ignores non-401 error response (e.g. 403)', async () => {
        const mockError = { error: 'Forbidden' };
        global.fetch.mockResolvedValue({
            ok: false,
            status: 403, // <--- non-401
            headers: { get: () => 'application/json' },
            json: vi.fn().mockResolvedValue(mockError),
        });

        await apiRequest('http://test-url', {}, {
            onError,
            errorMessage: 'Should not appear for 403'
        });

        // Expect no toasts or onError since code returns early
        expect(toast.error).not.toHaveBeenCalled();
        expect(onError).not.toHaveBeenCalled();
    });

    test('handles error non-JSON response with fallback errorMessage on 401', async () => {
        global.fetch.mockResolvedValue({
            ok: false,
            status: 401, // 401 triggers onError
            headers: { get: () => 'text/html' },
            text: vi.fn().mockResolvedValue('Some non-JSON error'),
        });

        await apiRequest('http://test-url', {}, {
            onError,
            errorMessage: 'A fallback error for 401'
        });

        // onError is called with 'A fallback error for 401'
        // because "result?.error" is undefined
        expect(onError).toHaveBeenCalledWith('A fallback error for 401');
        expect(toast.error).toHaveBeenCalledWith('A fallback error for 401');
    });

    test('handles network error', async () => {
        global.fetch.mockRejectedValue(new Error('Network failed'));

        await apiRequest('http://test-url', {}, {
            onError,
            errorMessage: 'Request failed'
        });

        // The catch block calls onError with the thrown error message
        expect(onError).toHaveBeenCalledWith('Network failed');
        // No success toast
        expect(toast.success).not.toHaveBeenCalled();
    });

    // Test ensures no error toast is shown when disableErrorToast = true (401 scenario)
    test('does not call toast.error when disableErrorToast=true', async () => {
        global.fetch.mockResolvedValue({
            ok: false,
            status: 401,
            headers: { get: () => 'application/json' },
            json: vi.fn().mockResolvedValue({ error: 'Unauthorized' }),
        });

        await apiRequest('http://test-url', {}, {
            onError,
            errorMessage: 'Some error',
            disableErrorToast: true, // <--- This is the key
        });

        // Verify that toast.error was NOT called
        expect(toast.error).not.toHaveBeenCalled();

        // We DO still call onError callback
        expect(onError).toHaveBeenCalledWith('Unauthorized');
    });
});
