// renderer/src/pages/Profiles.test.jsx
import { render, screen } from '@testing-library/react';
import { describe, test, expect, vi } from 'vitest';
import '@testing-library/jest-dom';
import { ThemeProvider } from '@mui/material/styles';
import { MemoryRouter } from 'react-router-dom';
import theme from '../Theme';

vi.mock('../api/syncApi', () => ({
    saveUserSelections: vi.fn(),
    syncSubdirectory: vi.fn(),
    syncAllSubdirectories: vi.fn(),
    chooseSettingsDir: vi.fn(),
    backupDirectory: vi.fn(),
    resetToDefaultDirectory: vi.fn(),
}));
vi.mock('../hooks/useAppData', () => ({
    useAppData: () => ({ refreshData: vi.fn() }),
}));

import Profiles from './Profiles';

const wrap = (ui) => render(
    <MemoryRouter><ThemeProvider theme={theme}>{ui}</ThemeProvider></MemoryRouter>,
);

describe('Profiles smart default', () => {
    const subDirs = [
        {
            profile: 'settings_default',
            availableUserFiles: [{ userId: 'u1', name: 'User1', mtime: '2026-04-01T10:00:00.000Z' }],
            availableCharFiles: [
                { charId: 'c1', name: 'Char1', mtime: '2026-04-01T10:00:00.000Z' },
                { charId: 'c2', name: 'Char2', mtime: '2026-04-01T10:00:00.000Z' },
            ],
        },
    ];

    test('defaults to mapping view when there are unmatched characters', () => {
        wrap(<Profiles subDirs={subDirs} associations={[]} settingsData={[]} userSelections={{}} currentSettingsDir="" lastBackupDir="" />);
        // Mapping mode shows the column headings.
        expect(screen.getByText(/user files/i)).toBeInTheDocument();
    });

    test('defaults to sync view when all characters are matched', () => {
        const allMatched = [{ userId: 'u1', charId: 'c1', charName: 'Char1' }, { userId: 'u1', charId: 'c2', charName: 'Char2' }];
        wrap(<Profiles subDirs={subDirs} associations={allMatched} settingsData={[]} userSelections={{}} currentSettingsDir="" lastBackupDir="" />);
        // Sync mode (with empty settingsData) shows the empty-sync state.
        expect(screen.getByText(/no eve profiles found/i)).toBeInTheDocument();
    });

    test('renders Profiles title in the header', () => {
        wrap(<Profiles subDirs={[]} associations={[]} settingsData={[]} userSelections={{}} currentSettingsDir="" lastBackupDir="" />);
        expect(screen.getByRole('heading', { name: 'Profiles' })).toBeInTheDocument();
    });
});

describe('Profiles reset-to-default button visibility', () => {
    test('hides the reset-to-default button when isDefaultDir is true (sync mode)', () => {
        wrap(
            <Profiles
                subDirs={[]}
                associations={[]}
                settingsData={[{ profile: 'settings_default', availableCharFiles: [], availableUserFiles: [] }]}
                userSelections={{}}
                currentSettingsDir="/path/to/Tranquility"
                isDefaultDir={true}
                lastBackupDir=""
            />,
        );
        expect(screen.queryByRole('button', { name: /reset to default directory/i })).toBeNull();
    });

    test('shows the reset-to-default button when isDefaultDir is false (sync mode)', () => {
        wrap(
            <Profiles
                subDirs={[]}
                associations={[]}
                settingsData={[{ profile: 'settings_default', availableCharFiles: [], availableUserFiles: [] }]}
                userSelections={{}}
                currentSettingsDir="/some/custom/dir"
                isDefaultDir={false}
                lastBackupDir=""
            />,
        );
        expect(screen.getByRole('button', { name: /reset to default directory/i })).toBeInTheDocument();
    });
});
