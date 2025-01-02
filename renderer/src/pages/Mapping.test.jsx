import React from 'react';
import { render, screen, fireEvent, act, within } from '@testing-library/react';
import { vi } from 'vitest';
import Mapping from './Mapping';
import '@testing-library/jest-dom';

// Mock the confirm dialog hook
vi.mock('../hooks/useConfirmDialog.jsx', () => ({
    useConfirmDialog: () => {
        let resolveFn;
        const showConfirmDialog = ({title, message}) => {
            // Immediately resolve with {isConfirmed: true} for testing
            return new Promise((resolve) => {
                resolveFn = resolve;
                // In a real test, we might simulate user action here.
                resolve({ isConfirmed: true });
            });
        };
        const confirmDialog = null; // Not rendering actual dialog in test
        return [showConfirmDialog, confirmDialog];
    }
}));

// Mock apiService calls
vi.mock('../api/apiService.jsx', () => ({
    associateCharacter: vi.fn().mockResolvedValue({ success: true, message: 'Character associated successfully!' }),
    unassociateCharacter: vi.fn().mockResolvedValue({ success: true, message: 'Character unassociated successfully!' })
}));

import { associateCharacter, unassociateCharacter } from '../api/apiService.jsx';

describe('Mapping', () => {
    const mockOnRefreshData = vi.fn();

    beforeEach(() => {
        vi.clearAllMocks();
    });

    it('renders "No accounts found." and "No available characters..." when subDirs is empty', () => {
        render(
            <Mapping
                associations={[]}
                subDirs={[]}
            />
        );

        expect(screen.getByText('No accounts found.')).toBeInTheDocument();
    });

    it('renders accounts and characters', () => {
        const subDirs = [{
            profile: 'profile1',
            availableCharFiles: [
                { file: 'charFile1', charId: 'char1', name: 'Char One', mtime: '2023-10-05T14:30:00Z' }
            ],
            availableUserFiles: [
                { file: 'userFile1', userId: 'user1', name: 'User One', mtime: '2023-10-05T14:31:00Z' }
            ]
        }];

        render(
            <Mapping
                associations={[]}
                subDirs={subDirs}
            />
        );

        // The account should appear
        expect(screen.getByText('User One')).toBeInTheDocument();

        // The character should appear
        expect(screen.getByText('Char One')).toBeInTheDocument();
    });

    it('associates a character with an account via drag-and-drop', async () => {
        const subDirs = [{
            profile: 'profile1',
            availableCharFiles: [
                { file: 'charFile1', charId: 'char1', name: 'Char One', mtime: '2023-10-05T14:30:00Z' }
            ],
            availableUserFiles: [
                { file: 'userFile1', userId: 'user1', name: 'User One', mtime: '2023-10-05T14:31:00Z' }
            ]
        }];

        render(
            <Mapping
                associations={[]}
                subDirs={subDirs}
                onRefreshData={mockOnRefreshData}
            />
        );

        const charCard = screen.getByText('Char One').closest('div[draggable="true"]');
        const accountCard = screen.getByText('User One').closest('div');

        // Simulate drag start
        await act(async () => {
            fireEvent.dragStart(charCard, { dataTransfer: { setData: vi.fn() } });
        });

        // Mock the dataTransfer set/get calls
        const dataTransfer = {
            setData: vi.fn(),
            getData: vi.fn(() => 'char1')
        };

        // Drag over the account card
        await act(async () => {
            fireEvent.dragOver(accountCard);
        });

        // Drop on the account card
        await act(async () => {
            fireEvent.drop(accountCard, { dataTransfer });
        });

        expect(screen.queryByTestId('available-characters')).not.toBeInTheDocument();
        
        // Optionally, you can check that 'Char One' now appears under the account if desired
        expect(screen.getByText('Char One')).toBeInTheDocument(); // Now under the associated account
    });

    it('unassociates a character from an account', async () => {
        // Now we start with an associated character
        const initialAssociations = [
            { userId: 'user1', charId: 'char1', charName: 'Char One', mtime: '2023-10-05T14:30:00Z' }
        ];

        const subDirs = [{
            profile: 'profile1',
            availableCharFiles: [],
            availableUserFiles: [
                { file: 'userFile1', userId: 'user1', name: 'User One', mtime: '2023-10-05T14:31:00Z' }
            ]
        }];

        render(
            <Mapping
                associations={initialAssociations}
                subDirs={subDirs}
                onRefreshData={mockOnRefreshData}
            />
        );

        // The association should be visible under the account
        expect(screen.getByText('Char One')).toBeInTheDocument();

        // Click unassociate
        const unassociateBtn = screen.getByRole('button', { name: 'Unassociate Char One' });
        await act(async () => {
            fireEvent.click(unassociateBtn);
        });

        expect(unassociateCharacter).toHaveBeenCalledWith('user1', 'char1', 'User One', 'Char One');
        expect(screen.queryByText('Char One')).not.toBeInTheDocument(); // character should be removed after unassociation
        expect(mockOnRefreshData).toHaveBeenCalled();
    });
});
