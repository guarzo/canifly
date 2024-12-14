// src/pages/CharacterOverview.test.jsx
import React from 'react';
import { render, screen } from '@testing-library/react';
import '@testing-library/jest-dom';
import userEvent from '@testing-library/user-event';
import { vi } from 'vitest';
import Overview from './CharacterOverview.jsx';

describe('Overview', () => {
    let user;
    beforeAll(() => {
        user = userEvent.setup();
    });

    const mockAccounts = [
        {
            ID: 1,
            Name: 'Account A',
            Characters: [
                {
                    Character: {
                        CharacterName: 'Alice',
                        CharacterID: 1111, // Provide a unique ID here
                        LocationName: 'Jita',
                        CharacterSkillsResponse: { total_sp: 1000000 }
                    },
                    Role: 'Pvp',
                    MCT: false,
                }
            ]
        },
        {
            ID: 2,
            Name: 'Account B',
            Characters: [
                {
                    Character: {
                        CharacterName: 'Bob',
                        CharacterID: 2222, // And a unique ID here as well
                        LocationName: 'Amarr',
                        CharacterSkillsResponse: { total_sp: 2000000 }
                    },
                    Role: 'Logistics',
                    MCT: true,
                }
            ]
        }
    ];

    const mockConversions = {}


    const roles = ['Pvp', 'Logistics'];

    const onToggleAccountStatus = vi.fn();
    const onUpdateCharacter = vi.fn();
    const onUpdateAccountName = vi.fn();
    const onRemoveCharacter = vi.fn();
    const onRemoveAccount = vi.fn();

    test('renders accounts by default', () => {
        render(
            <Overview
                accounts={mockAccounts}
                onToggleAccountStatus={onToggleAccountStatus}
                onUpdateCharacter={onUpdateCharacter}
                onUpdateAccountName={onUpdateAccountName}
                onRemoveCharacter={onRemoveCharacter}
                onRemoveAccount={onRemoveAccount}
                roles={roles}
                skillConversions={mockConversions}
            />
        );

        // By default, grouping is by account, so Account A and Account B should appear
        expect(screen.getByText('Account A')).toBeInTheDocument();
        expect(screen.getByText('Account B')).toBeInTheDocument();
        // Character names rendered by AccountCard would appear too (if AccountCard works as expected)
        expect(screen.getByText('Alice')).toBeInTheDocument();
        expect(screen.getByText('Bob')).toBeInTheDocument();
    });

    test('can toggle sorting', async () => {
        render(
            <Overview
                accounts={mockAccounts}
                onToggleAccountStatus={onToggleAccountStatus}
                onUpdateCharacter={onUpdateCharacter}
                onUpdateAccountName={onUpdateAccountName}
                onRemoveCharacter={onRemoveCharacter}
                onRemoveAccount={onRemoveAccount}
                roles={roles}
                skillConversions={mockConversions}
            />
        );

        // Initially sorted in ascending order, Account A should appear before Account B
        const accountCards = screen.getAllByText(/Account [A|B]/i);
        expect(accountCards[0]).toHaveTextContent('Account A');
        expect(accountCards[1]).toHaveTextContent('Account B');

        // Click the sort icon to switch to descending order
        const sortButton = screen.getByRole('button', { name: /sort/i });
        await user.click(sortButton);

        // After toggling, Account B should come first if it's sorted by name desc
        const reorderedAccounts = screen.getAllByText(/Account [A|B]/i);
        expect(reorderedAccounts[0]).toHaveTextContent('Account B');
        expect(reorderedAccounts[1]).toHaveTextContent('Account A');
    });

    test('displays no accounts message if empty', () => {
        render(
            <Overview
                accounts={[]}
                onToggleAccountStatus={onToggleAccountStatus}
                onUpdateCharacter={onUpdateCharacter}
                onUpdateAccountName={onUpdateAccountName}
                onRemoveCharacter={onRemoveCharacter}
                onRemoveAccount={onRemoveAccount}
                roles={roles}
                skillConversions={mockConversions}
            />
        );

        expect(screen.getByText('No accounts found.')).toBeInTheDocument();
    });
});
