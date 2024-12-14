import React from 'react';
import { render, screen } from '@testing-library/react';
import AccountCard from './AccountCard';

test('renders account name and status', () => {
    const mockAccount = {
        Name: 'TestAccount',
        Status: 'Omega',
        Characters: [
            {
                MCT: false,
                Role: 'Pvp',
                Training: 'Minmatar Dreadnought',
                Character: {
                    CharacterName: 'TestCharacter',
                    CharacterID: 123456, // A unique ID
                    CharacterSkillsResponse: { total_sp: 5000000 },
                    LocationName: 'Jita'
                }
            }
        ]
    };

    const mockConversions = {}

    render(
        <AccountCard
            account={mockAccount}
            onToggleAccountStatus={() => {}}
            onUpdateAccountName={() => {}}
            onUpdateCharacter={() => {}}
            onRemoveCharacter={() => {}}
            onRemoveAccount={() => {}}
            roles={[]}
            skillConversions={mockConversions}
        />
    );

    // Check account name
    expect(screen.getByText('TestAccount')).toBeInTheDocument();

    // Check character name inside account
    expect(screen.getByText('TestCharacter')).toBeInTheDocument();
});
