import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { vi } from 'vitest';
import SubDirectoryCard from './SubDirectoryCard';
import '@testing-library/jest-dom';

describe('SubDirectoryCard', () => {
    const mockHandleSelectionChange = vi.fn();
    const mockHandleSync = vi.fn();
    const mockHandleSyncAll = vi.fn();

    const subDir = {
        profile: 'settings_testProfile',
        availableCharFiles: [
            { charId: 'char1', name: 'Character One' },
            { charId: 'char2', name: 'Character Two' },
        ],
        availableUserFiles: [
            { userId: 'userA', name: 'User A' },
            { userId: 'userB', name: 'User B' },
        ]
    };

    const selections = {
        'settings_testProfile': {
            charId: '',
            userId: ''
        }
    };

    beforeEach(() => {
        mockHandleSelectionChange.mockClear();
        mockHandleSync.mockClear();
        mockHandleSyncAll.mockClear();
    });

    it('renders subdirectory name and dropdowns', () => {
        render(
            <SubDirectoryCard
                subDir={subDir}
                selections={selections}
                handleSelectionChange={mockHandleSelectionChange}
                handleSync={mockHandleSync}
                handleSyncAll={mockHandleSyncAll}
                isLoading={false}
            />
        );

        // displaySubDir = 'testProfile' after removing 'settings_'
        expect(screen.getByText('testProfile')).toBeInTheDocument();

        // Check character dropdown
        expect(screen.getByLabelText('Select Character')).toBeInTheDocument();
        // Check user dropdown
        expect(screen.getByLabelText('Select User')).toBeInTheDocument();
    });

    it('calls handleSelectionChange when user selects a character', () => {
        render(
            <SubDirectoryCard
                subDir={subDir}
                selections={selections}
                handleSelectionChange={mockHandleSelectionChange}
                handleSync={mockHandleSync}
                handleSyncAll={mockHandleSyncAll}
                isLoading={false}
            />
        );

        const charSelect = screen.getByLabelText('Select Character');
        fireEvent.mouseDown(charSelect); // Open the menu
        const charOption = screen.getByText('Character One');
        fireEvent.click(charOption);

        expect(mockHandleSelectionChange).toHaveBeenCalledWith('settings_testProfile', 'charId', 'char1');
    });

    it('calls handleSelectionChange when user selects a user', () => {
        render(
            <SubDirectoryCard
                subDir={subDir}
                selections={selections}
                handleSelectionChange={mockHandleSelectionChange}
                handleSync={mockHandleSync}
                handleSyncAll={mockHandleSyncAll}
                isLoading={false}
            />
        );

        const userSelect = screen.getByLabelText('Select User');
        fireEvent.mouseDown(userSelect);
        const userOption = screen.getByText('userA');
        fireEvent.click(userOption);

        expect(mockHandleSelectionChange).toHaveBeenCalledWith('settings_testProfile', 'userId', 'userA');
    });

    it('enables sync buttons only when both char and user are selected', () => {
        const updatedSelections = {
            'settings_testProfile': {
                charId: 'char1',
                userId: 'userA'
            }
        };

        render(
            <SubDirectoryCard
                subDir={subDir}
                selections={updatedSelections}
                handleSelectionChange={mockHandleSelectionChange}
                handleSync={mockHandleSync}
                handleSyncAll={mockHandleSyncAll}
                isLoading={false}
            />
        );

        // Now both should be enabled
        const syncButton = screen.getByRole('button', { name: /sync this specific profile/i });
        const syncAllButton = screen.getByRole('button', { name: /sync all profiles/i });
        expect(syncButton).toBeEnabled();
        expect(syncAllButton).toBeEnabled();

        // Click sync button
        fireEvent.click(syncButton);
        expect(mockHandleSync).toHaveBeenCalledWith('settings_testProfile');

        // Click sync all button
        fireEvent.click(syncAllButton);
        expect(mockHandleSyncAll).toHaveBeenCalledWith('settings_testProfile');
    });

    it('disables sync buttons if either char or user is not selected', () => {
        const partialSelections = {
            'settings_testProfile': {
                charId: 'char1',
                userId: '' // user not selected
            }
        };

        render(
            <SubDirectoryCard
                subDir={subDir}
                selections={partialSelections}
                handleSelectionChange={mockHandleSelectionChange}
                handleSync={mockHandleSync}
                handleSyncAll={mockHandleSyncAll}
                isLoading={false}
            />
        );

        const syncButton = screen.getByRole('button', { name: /sync this specific profile/i });
        const syncAllButton = screen.getByRole('button', { name: /sync all profiles/i });

        expect(syncButton).toBeDisabled();
        expect(syncAllButton).toBeDisabled();
    });

    it('shows loading state on sync buttons when isLoading is true', () => {
        const updatedSelections = {
            'settings_testProfile': {
                charId: 'char1',
                userId: 'userA'
            }
        };

        render(
            <SubDirectoryCard
                subDir={subDir}
                selections={updatedSelections}
                handleSelectionChange={mockHandleSelectionChange}
                handleSync={mockHandleSync}
                handleSyncAll={mockHandleSyncAll}
                isLoading={true}
            />
        );

        const syncButton = screen.getByRole('button', { name: /sync this specific profile/i });
        expect(syncButton).toBeDisabled(); // loading state also disables
    });
});
