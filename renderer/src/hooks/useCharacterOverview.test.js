// src/hooks/useCharacterOverview.test.js
import { vi, describe, test, expect, beforeEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';

const updateCharacter = vi.fn();
const deleteCharacter = vi.fn();
const refreshCharacter = vi.fn();

vi.mock('../api/accountsApi', () => ({
    updateCharacter: (...args) => updateCharacter(...args),
    deleteCharacter: (...args) => deleteCharacter(...args),
    refreshCharacter: (...args) => refreshCharacter(...args),
}));

const mockUpdateAccount = vi.fn();
const mockDeleteAccount = vi.fn();
const mockFetchAccounts = vi.fn();
let mockAccounts = [];

vi.mock('./useAppData', () => ({
    useAppData: () => ({
        accounts: mockAccounts,
        updateAccount: mockUpdateAccount,
        deleteAccount: mockDeleteAccount,
        fetchAccounts: mockFetchAccounts,
    }),
}));

vi.mock('./useAsyncOperation', () => ({
    useAsyncOperation: () => ({
        execute: async (fn) => fn(),
    }),
}));

import { useCharacterOverview } from './useCharacterOverview';

describe('useCharacterOverview', () => {
    beforeEach(() => {
        vi.clearAllMocks();
        // Reset persisted state by clearing the keys this hook reads.
        try {
            localStorage.clear();
        } catch { /* swallow */ }
        mockAccounts = [
            {
                ID: 1,
                Name: 'Alpha',
                Visible: true,
                Characters: [
                    { Character: { CharacterID: 1, CharacterName: 'A1', LocationName: 'Jita' }, Role: 'Pvp' },
                    { Character: { CharacterID: 2, CharacterName: 'A2', LocationName: 'Amarr' }, Role: '' },
                ],
            },
            {
                ID: 2,
                Name: 'Bravo',
                Visible: false,
                Characters: [
                    { Character: { CharacterID: 3, CharacterName: 'B1', LocationName: 'Jita' }, Role: 'Pve' },
                ],
            },
        ];
    });

    test('hides hidden accounts in the default view', () => {
        const { result } = renderHook(() => useCharacterOverview({ roles: ['Pvp', 'Pve'] }));
        // Default view is 'account', and hidden accounts are filtered out of the visible set.
        expect(result.current.summary.accountCount).toBe(1);
    });

    test('groups by role and drops empty role buckets', () => {
        const { result } = renderHook(() => useCharacterOverview({ roles: ['Pvp', 'Pve', 'Logi'] }));
        act(() => result.current.setView('role'));
        const keys = result.current.groups.map((g) => g.key);
        // Logi has zero characters in our fixture and must be dropped.
        expect(keys).toContain('Pvp');
        expect(keys).toContain('Unassigned');
        expect(keys).not.toContain('Logi');
    });

    test('handleUpdateCharacter dispatches to accountsApi.updateCharacter', async () => {
        const { result } = renderHook(() => useCharacterOverview({ roles: [] }));
        await act(async () => {
            await result.current.handleUpdateCharacter(99, { Role: 'Logi' });
        });
        expect(updateCharacter).toHaveBeenCalledWith(99, { Role: 'Logi' });
    });

    test('toggleExpanded toggles ids in the expanded set', () => {
        const { result } = renderHook(() => useCharacterOverview({ roles: [] }));
        act(() => result.current.toggleExpanded(7));
        expect(result.current.expanded.has(7)).toBe(true);
        act(() => result.current.toggleExpanded(7));
        expect(result.current.expanded.has(7)).toBe(false);
    });
});
