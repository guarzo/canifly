// src/api/configApi.js
//
// App configuration and EVE-developer credential endpoints.
import { apiRequest } from './apiClient';

export async function getConfig() {
    return apiRequest(`/api/config`, {
        method: 'GET',
        credentials: 'include'
    }, {
        errorMessage: 'Failed to fetch configuration.'
    });
}

export async function updateConfig(updates) {
    return apiRequest(`/api/config`, {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify(updates)
    }, {
        errorMessage: 'Failed to update configuration.'
    });
}

export async function checkEVEConfiguration() {
    return apiRequest(`/api/config/eve/status`, {
        method: 'GET',
        credentials: 'include'
    }, {
        errorMessage: 'Failed to check EVE configuration.'
    });
}

export async function saveEVECredentials(clientId, clientSecret) {
    return apiRequest(`/api/config/eve/credentials`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ clientId, clientSecret })
    }, {
        successMessage: 'EVE credentials saved successfully!',
        errorMessage: 'Failed to save EVE credentials.'
    });
}
