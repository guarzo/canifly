// src/components/dashboard/GroupCard.test.jsx
import React from 'react';
import { render, screen } from '@testing-library/react';
import '@testing-library/jest-dom';
import { vi } from 'vitest'; // Import vi from vitest
import GroupCard from './GroupCard';

describe('GroupCard', () => {
    const mockCharacters = [
        {
            MCT: false,
            Role: 'Pvp',
            Character: {
                CharacterName: 'CharacterOne',
                CharacterID: 1111,
                CharacterSkillsResponse: { total_sp: 2000000 },
                LocationName: 'Jita'
            }
        },
        {
            MCT: true,
            Role: 'Logistics',
            Character: {
                CharacterName: 'CharacterTwo',
                CharacterID: 2222,
                CharacterSkillsResponse: { total_sp: 4000000 },
                LocationName: 'Amarr'
            }
        }
    ];

    const roles = ['Pvp', 'Logistics'];
    const mockOnUpdateCharacter = vi.fn(); // Use vi.fn() instead of jest.fn()
    const mockConversions = {}

    test('renders the group name and its characters', () => {
        render(
            <GroupCard
                groupName="My Group"
                characters={mockCharacters}
                onUpdateCharacter={mockOnUpdateCharacter}
                roles={roles}
                skillConversions={mockConversions}
            />
        );

        // Check that the group name is displayed
        expect(screen.getByText('My Group')).toBeInTheDocument();

        // Check that each character name is displayed
        expect(screen.getByText('CharacterOne')).toBeInTheDocument();
        expect(screen.getByText('CharacterTwo')).toBeInTheDocument();

        // Confirm no remove buttons due to hideRemoveIcon = true
        expect(screen.queryByLabelText('Remove Character')).not.toBeInTheDocument();
    });
});
