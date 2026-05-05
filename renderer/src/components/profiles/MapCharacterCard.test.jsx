import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { vi } from 'vitest';
import CharacterCard from './MapCharacterCard';
import '@testing-library/jest-dom';

vi.mock('../../utils/formatter.jsx', () => ({
    formatDate: (isoString) => 'Oct 5, 14:30', // Mock a fixed date output
}));

describe('CharacterCard', () => {
    const mockHandleDragStart = vi.fn();
    const mockMtimeToColor = {
        '2023-10-05T14:30:00Z': '#4caf50',
    };

    const char = {
        name: 'Test Character',
        charId: 'char123',
        mtime: '2023-10-05T14:30:00Z'
    };

    beforeEach(() => {
        mockHandleDragStart.mockClear();
    });

    it('renders character name and formatted date', () => {
        render(
            <CharacterCard
                char={char}
                handleDragStart={mockHandleDragStart}
                mtimeToColor={mockMtimeToColor}
            />
        );

        expect(screen.getByText('Test Character')).toBeInTheDocument();
        expect(screen.getByText('Oct 5, 14:30')).toBeInTheDocument();
    });

    it('calls handleDragStart on drag start with charId', () => {
        render(
            <CharacterCard
                char={char}
                handleDragStart={mockHandleDragStart}
                mtimeToColor={mockMtimeToColor}
            />
        );

        const card = screen.getByText('Test Character').closest('div[draggable="true"]');
        fireEvent.dragStart(card);

        expect(mockHandleDragStart).toHaveBeenCalledWith(expect.any(Object), 'char123');
    });
});
