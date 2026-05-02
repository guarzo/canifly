// src/api/skillPlansApi.test.js
import { vi } from 'vitest';
import { saveSkillPlan, deleteSkillPlan } from './skillPlansApi';
import { apiRequest } from './apiClient';

vi.mock('./apiClient', () => ({
    apiRequest: vi.fn()
}));

describe('skillPlansApi', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    describe('saveSkillPlan', () => {
        test('falls through to POST create on probe error', async () => {
            global.fetch = vi.fn().mockRejectedValue(new Error('boom'));
            apiRequest.mockResolvedValue('saved');
            const result = await saveSkillPlan('MyPlan', { skill: 'Level5' });
            expect(apiRequest).toHaveBeenCalledWith('/api/skill-plans', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ name: 'MyPlan', content: { skill: 'Level5' } }),
                credentials: 'include'
            }, {
                successMessage: 'Skill Plan Created!',
                errorMessage: 'Failed to create skill plan.'
            });
            expect(result).toBe('saved');
        });
    });

    describe('deleteSkillPlan', () => {
        test('calls apiRequest correctly', async () => {
            apiRequest.mockResolvedValue('deleted');
            const result = await deleteSkillPlan('MyPlan');
            expect(apiRequest).toHaveBeenCalledWith(`/api/skill-plans/MyPlan`, {
                method: 'DELETE',
                credentials: 'include',
            }, {
                errorMessage: 'Failed to delete skill plan.'
            });
            expect(result).toBe('deleted');
        });
    });
});
