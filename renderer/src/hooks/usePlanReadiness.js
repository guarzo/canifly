//
// Single source of truth for (character, plan) → readiness. Pulled out of
// PlanMatrix's inlined cellStatus/cellTooltip and PlanList's child mapping
// so the matrix, list, and timeline views all encode the same rules.

import { useMemo } from 'react';
import { calculateDaysFromToday } from '../utils/formatter.jsx';

export function computePlanReadiness(character, planName) {
    const c = character?.Character || {};
    if (c.QualifiedPlans?.[planName]) {
        return { state: 'ready', eta: null, missingCount: 0, missingSkills: null };
    }
    if (c.PendingPlans?.[planName]) {
        const queue = c.SkillQueue;
        const training = Array.isArray(queue) && queue.length > 0 && Boolean(character.MCT);
        const eta = calculateDaysFromToday(c.PendingFinishDates?.[planName]) || null;
        return {
            state: training ? 'training' : 'queued',
            eta,
            missingCount: 0,
            missingSkills: null,
        };
    }
    const missingSkills = c.MissingSkills?.[planName] || null;
    const missingCount = missingSkills ? Object.keys(missingSkills).length : 0;
    return { state: 'idle', eta: null, missingCount, missingSkills };
}

export function usePlanReadiness(character, planName) {
    return useMemo(() => computePlanReadiness(character, planName), [character, planName]);
}
