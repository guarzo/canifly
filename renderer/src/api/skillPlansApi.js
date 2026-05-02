// src/api/skillPlansApi.js
//
// Skill plan endpoints (server-managed plans, not the EVE static data).
import { apiRequest } from './apiClient';

export async function getSkillPlans() {
    return apiRequest(`/api/skill-plans`, {
        method: 'GET',
        credentials: 'include'
    }, {
        errorMessage: 'Failed to fetch skill plans.'
    });
}

export async function getSkillPlan(planName) {
    return apiRequest(`/api/skill-plans/${encodeURIComponent(planName)}`, {
        method: 'GET',
        credentials: 'include'
    }, {
        errorMessage: 'Failed to fetch skill plan.'
    });
}

export async function createSkillPlan(planName, planContents) {
    return apiRequest(`/api/skill-plans`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ name: planName, content: planContents })
    }, {
        errorMessage: 'Failed to create skill plan.'
    });
}

export async function copySkillPlan(sourcePlanName, targetPlanName) {
    return apiRequest(`/api/skill-plans/${encodeURIComponent(sourcePlanName)}/copy`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ newName: targetPlanName })
    }, {
        errorMessage: 'Failed to copy skill plan.'
    });
}

export async function deleteSkillPlan(planName) {
    return apiRequest(`/api/skill-plans/${encodeURIComponent(planName)}`, {
        method: 'DELETE',
        credentials: 'include',
    }, {
        errorMessage: 'Failed to delete skill plan.'
    });
}

// saveSkillPlan does a probe GET to decide between PUT (update) and POST (create).
// The probe uses raw fetch (relative URL) deliberately to avoid the apiRequest
// 4xx-swallowing behavior; behaviour is preserved verbatim.
export async function saveSkillPlan(planName, planContents) {
    try {
        const response = await fetch(`/api/skill-plans/${encodeURIComponent(planName)}`, {
            method: 'GET',
            credentials: 'include'
        });

        if (response.ok) {
            return apiRequest(`/api/skill-plans/${encodeURIComponent(planName)}`, {
                method: 'PUT',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ content: planContents }),
                credentials: 'include'
            }, {
                successMessage: 'Skill Plan Updated!',
                errorMessage: 'Failed to update skill plan.'
            });
        }
        return apiRequest('/api/skill-plans', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ name: planName, content: planContents }),
            credentials: 'include'
        }, {
            successMessage: 'Skill Plan Created!',
            errorMessage: 'Failed to create skill plan.'
        });
    } catch (_error) {
        return apiRequest('/api/skill-plans', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ name: planName, content: planContents }),
            credentials: 'include'
        }, {
            successMessage: 'Skill Plan Created!',
            errorMessage: 'Failed to create skill plan.'
        });
    }
}
