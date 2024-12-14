import React from 'react';
import { render, screen, act, fireEvent } from '@testing-library/react';
import CharacterItem from './CharacterItem';
import '@testing-library/jest-dom';
import userEvent from '@testing-library/user-event';

beforeEach(() => {
    window.electronAPI = { openExternal: vi.fn() };
});

afterEach(() => {
    delete window.electronAPI;
});

describe('CharacterItem', () => {
    const mockCharacter = {
        MCT: true,
        Role: 'Pvp',
        Training: 'Minmatar Dreadnought',
        Character: {
            CharacterName: 'TestCharacter',
            CharacterID: 123456,
            CharacterSkillsResponse: { total_sp: 5000000 },
            LocationName: 'Jita'
        }
    };

    const mockSkillConversions = {}

    const roles = ['Pvp', 'Logistics', 'Scout'];
    let mockOnUpdateCharacter;
    let mockOnRemoveCharacter;

    beforeEach(() => {
        mockOnUpdateCharacter = vi.fn();
        mockOnRemoveCharacter = vi.fn();
    });

    test('renders character name, total sp, and location', () => {
        render(
            <CharacterItem
                character={mockCharacter}
                onUpdateCharacter={mockOnUpdateCharacter}
                onRemoveCharacter={mockOnRemoveCharacter}
                roles={roles}
                skillConversions={mockSkillConversions}
            />
        );

        expect(screen.getByText('TestCharacter')).toBeInTheDocument();
        expect(screen.getByText('5M SP')).toBeInTheDocument();
        expect(screen.getByText('Jita')).toBeInTheDocument();
    });

    test('displays MCT tooltip', async () => {
        render(
            <CharacterItem
                character={mockCharacter}
                onUpdateCharacter={mockOnUpdateCharacter}
                onRemoveCharacter={mockOnRemoveCharacter}
                roles={roles}
                skillConversions={mockSkillConversions}
            />
        );

        const mctIndicator = screen.getByTestId('mct-indicator');

        // Wrap the hover action in act
        await act(async () => {
            await userEvent.hover(mctIndicator);
        });

        const tooltipText = await screen.findByText('Training: Minmatar Dreadnought');
        expect(tooltipText).toBeInTheDocument();
    });

    test('no remove icon if hideRemoveIcon is true', () => {
        render(
            <CharacterItem
                character={mockCharacter}
                onUpdateCharacter={mockOnUpdateCharacter}
                onRemoveCharacter={mockOnRemoveCharacter}
                roles={roles}
                hideRemoveIcon={true}
                skillConversions={mockSkillConversions}
            />
        );

        const removeBtn = screen.queryByLabelText('Remove Character');
        expect(removeBtn).not.toBeInTheDocument();
    });


});
