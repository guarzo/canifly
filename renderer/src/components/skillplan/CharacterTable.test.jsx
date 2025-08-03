// src/components/dashboard/CharacterTable.test.jsx
import React from 'react';
import { render, screen } from '@testing-library/react';
import '@testing-library/jest-dom';
import userEvent from '@testing-library/user-event';
import { vi } from 'vitest';  // If you're using Vitest
import CharacterTable from './CharacterTable';

// Mock utility functions if needed
vi.mock('../../utils/formatter.jsx', () => ({
    calculateDaysFromToday: vi.fn(() => 5), // Returns a fixed number of days for simplicity
    formatNumberWithCommas: vi.fn((num) => num.toLocaleString()),
}));

describe('CharacterTable', () => {
    const mockCharacters = [
        {
            Character: {
                CharacterID: 1234,
                CharacterName: 'Test Character',
                CharacterSkillsResponse: {
                    total_sp: 1000000
                },
                QualifiedPlans: { "Plan A": true },
                PendingPlans: {},
                PendingFinishDates: {},
                MissingSkills: {}
            }
        }
    ];

    const mockSkillPlans = {
        "Plan A": {},
        "Plan B": {}
    };

    const mockConversions = {}

    test('renders character name and total skill points', () => {
        render(<CharacterTable characters={mockCharacters} skillPlans={mockSkillPlans} conversions={mockConversions} />);

        // Check if character name is present
        expect(screen.getByText('Test Character')).toBeInTheDocument();

        // Check if total skill points are formatted and displayed
        expect(screen.getByText('1,000,000')).toBeInTheDocument();

        // Both plans should create rows, but initially hidden under collapse
        // The expand/collapse icon should be visible since we have plans
        const toggleButton = screen.getByRole('button', { name: /expand/i });
        expect(toggleButton).toBeInTheDocument();
    });

    test('expands and collapses character plans', async () => {
        render(<CharacterTable characters={mockCharacters} skillPlans={mockSkillPlans} conversions={mockConversions} />);
        const user = userEvent.setup();

        // Initially, the plans are collapsed, and we only see the main row
        expect(screen.queryByText('↳ Plan A')).not.toBeInTheDocument();
        expect(screen.queryByText('↳ Plan B')).not.toBeInTheDocument();

        // Click the expand button
        const toggleButton = screen.getByRole('button', { name: /expand/i });
        await user.click(toggleButton);

        // Now we should see the plan rows
        expect(screen.getByText('↳ Plan A')).toBeInTheDocument();
        expect(screen.getByText('↳ Plan B')).toBeInTheDocument();

        // Click again to collapse
        await user.click(toggleButton);
        expect(screen.queryByText('↳ Plan A')).not.toBeInTheDocument();
        expect(screen.queryByText('↳ Plan B')).not.toBeInTheDocument();
    });
});
