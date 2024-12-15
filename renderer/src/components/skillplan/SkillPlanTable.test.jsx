// src/components/dashboard/SkillPlanTable.test.jsx
import React from 'react';
import { render, screen } from '@testing-library/react';
import '@testing-library/jest-dom';
import userEvent from '@testing-library/user-event';
import { vi } from 'vitest'; // If using Vitest
import SkillPlanTable from './SkillPlanTable';

// Mocking the utility function
vi.mock('../../utils/formatter.jsx', () => ({
    calculateDaysFromToday: vi.fn(() => 3), // Always return 3 for simplicity
}));

describe('SkillPlanTable', () => {
    let user;

    beforeAll(() => {
        user = userEvent.setup();
    });

    beforeEach(() => {
        // Mock window functions
        window.copySkillPlan = vi.fn();
        window.deleteSkillPlan = vi.fn();
    });

    afterEach(() => {
        vi.resetAllMocks();
    });

    const mockSkillPlans = {
        "Plan A": {
            Name: "Plan A",
            QualifiedCharacters: ["CharacterOne"],
            PendingCharacters: ["CharacterTwo"],
            MissingCharacters: ["CharacterThree"],
            Skills: ["ABC"]
        },
        "Plan B": {
            Name: "Plan B",
            QualifiedCharacters: [],
            PendingCharacters: [],
            MissingCharacters: []
        }
    };

    const mockCharacters = [
        {
            Character: {
                CharacterName: "CharacterOne",
                PendingFinishDates: { "Plan A": "2025-12-31" }
            }
        },
        {
            Character: {
                CharacterName: "CharacterTwo",
                PendingFinishDates: { "Plan A": "2025-12-31" }
            }
        },
        {
            Character: {
                CharacterName: "CharacterThree",
                PendingFinishDates: {}
            }
        }
    ];

    const mockConversions =  {}
    const mockOnCopy = vi.fn();

    const mockOnDelete = vi.fn()


    test('renders skill plans and allows expansion', async () => {
        render(<SkillPlanTable skillPlans={mockSkillPlans} characters={mockCharacters}
                               conversions={mockConversions} onDeleteSkillPlan={mockOnDelete} onCopySkillPlan={mockOnCopy}/>);

        // Check if both plan names are rendered
        expect(screen.getByText('Plan A')).toBeInTheDocument();
        expect(screen.getByText('Plan B')).toBeInTheDocument();

        // Initially, children are collapsed
        expect(screen.queryByText('↳ CharacterOne')).not.toBeInTheDocument();

        // Expand Plan A
        const expandButton = screen.getByRole('button', { name: /expand/i });
        await user.click(expandButton);

        // Now the children should be visible
        expect(screen.getByText('↳ CharacterOne')).toBeInTheDocument();
        expect(screen.getByText('↳ CharacterTwo')).toBeInTheDocument();
        expect(screen.getByText('↳ CharacterThree')).toBeInTheDocument();
    });

    test('copy and delete skill plan actions', async () => {
        render(<SkillPlanTable skillPlans={mockSkillPlans} characters={mockCharacters} conversions={mockConversions}
        onDeleteSkillPlan={mockOnDelete} onCopySkillPlan={mockOnCopy}/>);
        const user = userEvent.setup();

        // For Plan A
        const planACopyButton = screen.getAllByRole('button', { name: /copy skill plan/i })[0];
        const planADeleteButton = screen.getAllByRole('button', { name: /delete skill plan/i })[0];

        await user.click(planACopyButton);
        expect(mockOnCopy).toHaveBeenCalledWith('Plan A', ['ABC']);
        
        await user.click(planADeleteButton);
        expect(mockOnDelete).toHaveBeenCalledWith('Plan A');
    });
});
