import { vi, describe, test, expect, beforeEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';

const mockToastSuccess = vi.fn();
const mockToastError = vi.fn();
const mockToastInfo = vi.fn();
vi.mock('react-toastify', () => ({
    toast: {
        success: (...args) => mockToastSuccess(...args),
        error: (...args) => mockToastError(...args),
        info: (...args) => mockToastInfo(...args),
    },
}));

const mockLoggerError = vi.fn();
vi.mock('../utils/logger', () => ({
    logger: {
        debug: vi.fn(),
        error: (...args) => mockLoggerError(...args),
    },
}));

import { useSyncAction } from './useSyncAction';

describe('useSyncAction', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    test('toggles isLoading around the operation', async () => {
        const { result } = renderHook(() => useSyncAction());
        const op = vi.fn().mockResolvedValue({ success: true, message: 'done' });
        let pending;
        act(() => { pending = result.current.run(op); });
        expect(result.current.isLoading).toBe(true);
        await act(async () => { await pending; });
        expect(result.current.isLoading).toBe(false);
    });

    test('toasts result.message on success when no successMessage given', async () => {
        const { result } = renderHook(() => useSyncAction());
        await act(async () => {
            await result.current.run(() => Promise.resolve({ success: true, message: 'ok' }));
        });
        expect(mockToastSuccess).toHaveBeenCalledWith('ok');
    });

    test('uses successMessage override when provided', async () => {
        const { result } = renderHook(() => useSyncAction());
        await act(async () => {
            await result.current.run(() => Promise.resolve({ success: true }), { successMessage: 'custom' });
        });
        expect(mockToastSuccess).toHaveBeenCalledWith('custom');
    });

    test('logs and rethrows on error', async () => {
        const { result } = renderHook(() => useSyncAction());
        const err = new Error('boom');
        await act(async () => {
            await expect(result.current.run(() => Promise.reject(err), { errorContext: 'syncSubdirectory' }))
                .rejects.toThrow('boom');
        });
        expect(mockLoggerError).toHaveBeenCalled();
        expect(mockToastError).toHaveBeenCalled();
        expect(result.current.isLoading).toBe(false);
    });
});
