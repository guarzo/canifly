//
// Pure derivation of (accounts, availableCharacters, mtimeToColor) from the
// raw subDirs payload + current associations. Pulled out of Mapping.jsx so
// the merged Profiles page can reuse it (Sync mode also wants the user file
// → user name mapping for nicer Select labels).

import { useMemo } from 'react';

const SWATCH_PALETTE = [
    'oklch(0.72 0.13 145)',
    'oklch(0.74 0.14 25)',
    'oklch(0.80 0.13 80)',
    'oklch(0.70 0.14 305)',
    'oklch(0.75 0.12 200)',
    'oklch(0.72 0.13 350)',
    'oklch(0.78 0.10 240)',
    'oklch(0.74 0.10 120)',
];

function roundToMinute(mtime) {
    const d = new Date(mtime);
    if (Number.isNaN(d.getTime())) return null;
    d.setSeconds(0, 0);
    return d.toISOString();
}

export function computeMappingDerivations(subDirs, associations) {
    if (!Array.isArray(subDirs) || subDirs.length === 0) {
        return { accounts: [], availableCharacters: [], mtimeToColor: {} };
    }

    const userMap = {};
    subDirs.forEach((mapping) => {
        if (!mapping || !Array.isArray(mapping.availableUserFiles)) return;
        mapping.availableUserFiles.forEach((userFile) => {
            if (!userFile || !userFile.userId) return;
            const roundedMtime = roundToMinute(userFile.mtime);
            if (!roundedMtime) return;
            if (
                !userMap[userFile.userId] ||
                new Date(roundedMtime) > new Date(userMap[userFile.userId].mtime)
            ) {
                userMap[userFile.userId] = { ...userFile, mtime: roundedMtime };
            }
        });
    });
    const accounts = Object.values(userMap).sort(
        (a, b) => new Date(b.mtime) - new Date(a.mtime),
    );

    const charMap = {};
    subDirs.forEach((mapping) => {
        if (!mapping || !Array.isArray(mapping.availableCharFiles)) return;
        mapping.availableCharFiles.forEach((charFile) => {
            if (!charFile || !charFile.charId) return;
            const roundedMtime = roundToMinute(charFile.mtime);
            if (!roundedMtime) return;
            const { charId } = charFile;
            if (!charMap[charId] || new Date(roundedMtime) > new Date(charMap[charId].mtime)) {
                charMap[charId] = { ...charFile, mtime: roundedMtime, profile: mapping.profile };
            }
        });
    });

    const associatedCharIds = new Set((associations || []).map((a) => a.charId));
    const availableCharacters = Object.values(charMap)
        .filter((ch) => !associatedCharIds.has(ch.charId))
        .sort((a, b) => new Date(b.mtime) - new Date(a.mtime));

    const allMtimes = [
        ...accounts.map((a) => a.mtime),
        ...availableCharacters.map((c) => c.mtime),
    ];
    const ordered = Array.from(new Set(allMtimes)).sort(
        (a, b) => new Date(a) - new Date(b),
    );
    const mtimeToColor = ordered.reduce((acc, mtime, index) => {
        acc[mtime] = SWATCH_PALETTE[index % SWATCH_PALETTE.length];
        return acc;
    }, {});

    return { accounts, availableCharacters, mtimeToColor };
}

export function useMappingDerivations(subDirs, associations) {
    return useMemo(
        () => computeMappingDerivations(subDirs, associations),
        [subDirs, associations],
    );
}
