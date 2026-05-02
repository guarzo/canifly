// src/hooks/useHeaderData.test.js
import { vi, describe, test, expect, beforeEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';

const mockLogout = vi.fn();
const mockRefreshData = vi.fn().mockResolvedValue();
const mockExecute = vi.fn(async (fn) => fn());
const mockAddCharacter = vi.fn();

vi.mock('./useAuth', () => ({
    useAuth: () => ({ isAuthenticated: true, logout: mockLogout }),
}));
vi.mock('./useAppData', () => ({
    useAppData: () => ({ refreshData: mockRefreshData }),
}));
vi.mock('./useAsyncOperation', () => ({
    useAsyncOperation: () => ({ execute: mockExecute, isLoading: false }),
}));
vi.mock('../api/accountsApi', () => ({
    addCharacter: (...args) => mockAddCharacter(...args),
}));
vi.mock('react-toastify', () => ({
    toast: { success: vi.fn(), error: vi.fn(), info: vi.fn(), warning: vi.fn() },
}));

import { useHeaderData } from './useHeaderData';

describe('useHeaderData', () => {
    beforeEach(() => {
        vi.clearAllMocks();
        window.electronAPI = undefined;
    });

    test('handleAddCharacterSubmit calls accountsApi.addCharacter and refreshes when there is no redirect', async () => {
        mockAddCharacter.mockResolvedValue({ ok: true });
        const { result } = renderHook(() => useHeaderData());
        await act(async () => {
            await result.current.handleAddCharacterSubmit('TestAccount');
        });
        expect(mockAddCharacter).toHaveBeenCalledWith('TestAccount');
        expect(mockRefreshData).toHaveBeenCalled();
    });

    test('handleCloseWindow is a no-op when electronAPI is unavailable', () => {
        const { result } = renderHook(() => useHeaderData());
        expect(() => act(() => result.current.handleCloseWindow())).not.toThrow();
    });

    test('handleCloseWindow calls electronAPI.closeWindow when available', () => {
        const closeWindow = vi.fn();
        window.electronAPI = { closeWindow };
        const { result } = renderHook(() => useHeaderData());
        act(() => result.current.handleCloseWindow());
        expect(closeWindow).toHaveBeenCalled();
    });
});
