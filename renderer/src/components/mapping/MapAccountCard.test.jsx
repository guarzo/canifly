import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { vi } from 'vitest';
import MapAccountCard from './MapAccountCard';
import '@testing-library/jest-dom';

vi.mock('../../utils/formatter.jsx', () => ({
    formatDate: (isoString) => 'Oct 5, 14:30', // Mocked fixed date formatting
}));

describe('MapAccountCard', () => {
    const mockHandleUnassociate = vi.fn();
    const mockHandleDrop = vi.fn();

    const mockMtimeToColor = {
        '2023-10-05T14:30:00Z': '#4caf50',
    };

    const mapping = {
        userId: 'user123',
        name: 'TestUser',
        mtime: '2023-10-05T14:30:00Z'
    };

    const associations = [
        { userId: 'user123', charId: 'char1', charName: 'Char One', mtime: '2023-10-05T12:00:00Z' },
        { userId: 'user123', charId: 'char2', charName: 'Char Two', mtime: '2023-10-05T13:00:00Z' },
    ];

    beforeEach(() => {
        mockHandleUnassociate.mockClear();
        mockHandleDrop.mockClear();
    });

    it('renders account name and formatted date', () => {
        render(
            <MapAccountCard
                mapping={mapping}
                associations={[]}
                handleUnassociate={mockHandleUnassociate}
                handleDrop={mockHandleDrop}
                mtimeToColor={mockMtimeToColor}
            />
        );

        // Check name
        expect(screen.getByText('TestUser')).toBeInTheDocument();
        // Since we mocked formatDate, it should show "Oct 5, 14:30"
        expect(screen.getByText('Oct 5, 14:30')).toBeInTheDocument();
    });

    it('displays associated characters and calls handleUnassociate on delete', () => {
        render(
            <MapAccountCard
                mapping={mapping}
                associations={associations}
                handleUnassociate={mockHandleUnassociate}
                handleDrop={mockHandleDrop}
                mtimeToColor={mockMtimeToColor}
            />
        );

        // Both characters should be listed
        expect(screen.getByText('Char One')).toBeInTheDocument();
        expect(screen.getByText('Char Two')).toBeInTheDocument();

        // Click the unassociate button for "Char One"
        const unassociateBtn = screen.getByRole('button', { name: 'Unassociate Char One' });
        fireEvent.click(unassociateBtn);

        expect(mockHandleUnassociate).toHaveBeenCalledWith('user123', 'char1', 'Char One', 'TestUser');
    });

});
