import { describe, test, expect } from 'vitest';
import { computeMappingDerivations } from './useMappingDerivations';

describe('computeMappingDerivations', () => {
    const subDirs = [
        {
            profile: 'settings_default',
            availableUserFiles: [
                { userId: 'u1', name: 'User1', mtime: '2026-04-01T10:00:30.000Z' },
                { userId: 'u2', name: 'User2', mtime: '2026-04-01T11:00:00.000Z' },
            ],
            availableCharFiles: [
                { charId: 'c1', name: 'Char1', mtime: '2026-04-01T10:00:30.000Z' },
                { charId: 'c2', name: 'Char2', mtime: '2026-04-01T12:00:00.000Z' },
            ],
        },
    ];

    test('rounds mtimes to the minute', () => {
        const { accounts } = computeMappingDerivations(subDirs, []);
        const u1 = accounts.find((a) => a.userId === 'u1');
        expect(u1.mtime).toBe('2026-04-01T10:00:00.000Z');
    });

    test('orders accounts newest mtime first', () => {
        const { accounts } = computeMappingDerivations(subDirs, []);
        expect(accounts.map((a) => a.userId)).toEqual(['u2', 'u1']);
    });

    test('filters out already-associated characters', () => {
        const associations = [{ userId: 'u1', charId: 'c1', charName: 'Char1' }];
        const { availableCharacters } = computeMappingDerivations(subDirs, associations);
        expect(availableCharacters.map((c) => c.charId)).toEqual(['c2']);
    });

    test('assigns a color from the palette to every distinct mtime bucket', () => {
        const { mtimeToColor } = computeMappingDerivations(subDirs, []);
        const keys = Object.keys(mtimeToColor);
        expect(keys.length).toBe(3); // u1@10:00, u2@11:00, c2@12:00
        for (const k of keys) {
            expect(typeof mtimeToColor[k]).toBe('string');
        }
    });

    test('returns empty derivations when subDirs is empty', () => {
        const result = computeMappingDerivations([], []);
        expect(result).toEqual({ accounts: [], availableCharacters: [], mtimeToColor: {} });
    });
});
