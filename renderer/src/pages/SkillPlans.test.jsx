// Legacy assertions referenced removed UI text. Replaced with a smoke test
// that mounts the new page; expanded coverage lives in the matrix/list specs.
import { render, screen } from '@testing-library/react';
import '@testing-library/jest-dom';
import { vi, describe, test, expect } from 'vitest';
import SkillPlans from './SkillPlans';

vi.mock('../api/apiService.jsx', () => ({
    default: {
        copySkillPlan: vi.fn(),
        deleteSkillPlan: vi.fn(),
    },
    copySkillPlan: vi.fn(),
    deleteSkillPlan: vi.fn(),
}));

vi.mock('react-toastify', () => ({
    toast: { success: vi.fn(), warning: vi.fn(), error: vi.fn() },
}));

describe('SkillPlans page', () => {
    test('renders the title', () => {
        render(<SkillPlans characters={[]} skillPlans={{}} conversions={{}} />);
        expect(screen.getByText('Skill Plans')).toBeInTheDocument();
    });
});
