import { describe, test, expect } from 'vitest';
import { computePlanReadiness } from './usePlanReadiness';

describe('computePlanReadiness', () => {
    test('returns ready when QualifiedPlans contains the plan', () => {
        const ch = { Character: { QualifiedPlans: { 'Hurricane': true } } };
        expect(computePlanReadiness(ch, 'Hurricane').state).toBe('ready');
    });

    test('returns training when PendingPlans + active SkillQueue + MCT', () => {
        const ch = {
            MCT: true,
            Character: {
                PendingPlans: { 'Hurricane': true },
                SkillQueue: [{ skill_id: 1 }],
                PendingFinishDates: { 'Hurricane': '2026-06-01T00:00:00Z' },
            },
        };
        const r = computePlanReadiness(ch, 'Hurricane');
        expect(r.state).toBe('training');
        expect(r.eta).toBeTruthy();
    });

    test('returns queued when PendingPlans without active queue', () => {
        const ch = {
            MCT: false,
            Character: {
                PendingPlans: { 'Hurricane': true },
                SkillQueue: [],
            },
        };
        expect(computePlanReadiness(ch, 'Hurricane').state).toBe('queued');
    });

    test('returns idle with missingCount when MissingSkills exists', () => {
        const ch = {
            Character: {
                MissingSkills: { 'Hurricane': { 'Gunnery': 5, 'Motion Prediction': 3 } },
            },
        };
        const r = computePlanReadiness(ch, 'Hurricane');
        expect(r.state).toBe('idle');
        expect(r.missingCount).toBe(2);
        expect(r.missingSkills).toEqual({ 'Gunnery': 5, 'Motion Prediction': 3 });
    });

    test('returns idle with zero missingCount when nothing matches', () => {
        const ch = { Character: {} };
        const r = computePlanReadiness(ch, 'Hurricane');
        expect(r.state).toBe('idle');
        expect(r.missingCount).toBe(0);
    });
});
