// renderer/src/components/profiles/SyncView.test.jsx
import { render, screen } from '@testing-library/react';
import { describe, test, expect, vi } from 'vitest';
import '@testing-library/jest-dom';
import { ThemeProvider } from '@mui/material/styles';
import theme from '../../Theme';

vi.mock('../../api/syncApi', () => ({
    saveUserSelections: vi.fn().mockResolvedValue({ success: true }),
    syncSubdirectory: vi.fn().mockResolvedValue({ success: true }),
    syncAllSubdirectories: vi.fn().mockResolvedValue({ success: true }),
}));

import SyncView from './SyncView';

const wrap = (ui) => render(<ThemeProvider theme={theme}>{ui}</ThemeProvider>);

describe('SyncView', () => {
    test('renders the empty state when settingsData is empty', () => {
        wrap(<SyncView settingsData={[]} associations={[]} userSelections={{}} filter="" />);
        expect(screen.getByText(/no eve profiles found/i)).toBeInTheDocument();
    });

    test('renders one row per profile', () => {
        const settingsData = [
            { profile: 'settings_default', availableCharFiles: [], availableUserFiles: [] },
            { profile: 'settings_alt', availableCharFiles: [], availableUserFiles: [] },
        ];
        wrap(<SyncView settingsData={settingsData} associations={[]} userSelections={{}} filter="" />);
        expect(screen.getByText('default')).toBeInTheDocument();
        expect(screen.getByText('alt')).toBeInTheDocument();
    });
});
