// src/api/esiApi.js
//
// EVE static-data endpoints (skill plans, profiles, conversions).
import { apiRequest } from './apiClient';

export async function getEveSkillPlans() {
    return apiRequest(`/api/eve/skill-plans`, {
        method: 'GET',
        credentials: 'include'
    }, {
        errorMessage: 'Failed to fetch EVE skill plans.'
    });
}

export async function getEveProfiles() {
    return apiRequest(`/api/eve/profiles`, {
        method: 'GET',
        credentials: 'include'
    }, {
        errorMessage: 'Failed to fetch EVE profiles.'
    });
}

export async function getEveConversions() {
    return apiRequest(`/api/eve/conversions`, {
        method: 'GET',
        credentials: 'include'
    }, {
        errorMessage: 'Failed to fetch EVE conversions.'
    });
}
