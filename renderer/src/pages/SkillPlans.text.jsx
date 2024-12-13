// src/components/skillplan/SkillPlans.test.jsx
import React from 'react';
import { render, screen } from '@testing-library/react';
import '@testing-library/jest-dom';
import userEvent from '@testing-library/user-event';
import { vi } from 'vitest';
import SkillPlans from './SkillPlans';

// Mock the API and toast if necessary
vi.mock('../api/apiService.jsx', () => ({
    deleteSkillPlan: vi.fn().mockResolvedValue({ success: true }),
}));

vi.mock('react-toastify', () => ({
    toast: {
        success: vi.fn(),
        warning: vi.fn(),
        error: vi.fn(),
    }
}));

describe('SkillPlans', () => {
    const user = userEvent.setup();

    const mockCharacters = [
        {
            Character: {
                CharacterName: 'TestCharacter',
                LocationName: 'Jita',
                CharacterSkillsResponse: { total_sp: 1000000 }
            }
        }
    ];

    const mockSkillPlans = {
        "Plan A": {
            Name: "Plan A",
            Skills: {
                "SkillX": { Level: 5 },
            },
            QualifiedCharacters: ["TestCharacter"],
            PendingCharacters: [],
            MissingCharacters: []
        }
    };

    const mockSetAppData = vi.fn();

    beforeEach(() => {
        mockSetAppData.mockClear();
    });

    test('renders and defaults to characters view', () => {
        render(
            <SkillPlans
                characters={mockCharacters}
                skillPlans={mockSkillPlans}
                setAppData={mockSetAppData}
                backEndURL="http://mock-backend"
            />
        );

        // Check main heading
        expect(screen.getByText('Skill Plans')).toBeInTheDocument();

        // By default view is 'characters', so "By Character" should be visible
        expect(screen.getByText('By Character')).toBeInTheDocument();
        expect(screen.queryByText('By Skill Plan')).not.toBeInTheDocument();
    });

    test('can switch to skill plans view', async () => {
        render(
            <SkillPlans
                characters={mockCharacters}
                skillPlans={mockSkillPlans}
                setAppData={mockSetAppData}
                backEndURL="http://mock-backend"
            />
        );

        // The toggle buttons are controlled, so let's find the "View Skill Plans" button
        const plansToggle = screen.getByRole('button', { name: /View Skill Plans/i });
        await user.click(plansToggle);

        // Now we should see the "By Skill Plan" section
        expect(screen.getByText('By Skill Plan')).toBeInTheDocument();
    });

    test('window.copySkillPlan and window.deleteSkillPlan functions are defined', () => {
        render(
            <SkillPlans
                characters={mockCharacters}
                skillPlans={mockSkillPlans}
                setAppData={mockSetAppData}
                backEndURL="http://mock-backend"
            />
        );

        // After rendering, these functions should be set
        expect(typeof window.copySkillPlan).toBe('function');
        expect(typeof window.deleteSkillPlan).toBe('function');
    });
});
