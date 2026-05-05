// renderer/src/components/profiles/MappingView.test.jsx
import { render, screen } from '@testing-library/react';
import { describe, test, expect } from 'vitest';
import '@testing-library/jest-dom';
import { ThemeProvider } from '@mui/material/styles';
import theme from '../../Theme';
import MappingView from './MappingView';

const wrap = (ui) => render(<ThemeProvider theme={theme}>{ui}</ThemeProvider>);

const subDirs = [
    {
        profile: 'settings_default',
        availableUserFiles: [{ userId: 'u1', name: 'User1', mtime: '2026-04-01T10:00:00.000Z' }],
        availableCharFiles: [{ charId: 'c1', name: 'Char1', mtime: '2026-04-01T10:00:00.000Z' }],
    },
];

describe('MappingView', () => {
    test('renders the user file column and character column when there is data', () => {
        wrap(<MappingView subDirs={subDirs} associations={[]} filter="" view="all" sortOrder="mtime-desc" />);
        expect(screen.getByText(/user files/i)).toBeInTheDocument();
        expect(screen.getByText(/unassociated characters/i)).toBeInTheDocument();
        expect(screen.getByText('User1')).toBeInTheDocument();
        expect(screen.getByText('Char1')).toBeInTheDocument();
    });

    test('shows the no-accounts empty state when subDirs is empty', () => {
        wrap(<MappingView subDirs={[]} associations={[]} filter="" view="all" sortOrder="mtime-desc" />);
        expect(screen.getByText(/no accounts found/i)).toBeInTheDocument();
    });
});
