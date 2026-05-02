// src/api/fuzzworksApi.js
//
// Fuzzworks (EVE static data download / status) endpoints.
import { apiRequest } from './apiClient';

export async function getFuzzworksStatus() {
    return apiRequest(`/api/fuzzworks/status`, {
        method: 'GET',
        credentials: 'include'
    }, {
        errorMessage: 'Failed to fetch Fuzzworks status.'
    });
}

export async function updateFuzzworks() {
    return apiRequest(`/api/fuzzworks/update`, {
        method: 'POST',
        credentials: 'include'
    }, {
        errorMessage: 'Failed to update Fuzzworks data.'
    });
}
