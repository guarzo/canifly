// src/api/accountsApi.js
//
// Accounts, characters, and associations endpoints.
import { apiRequest } from './apiClient';

export async function getAccounts(bypassCache = false) {
    const cacheParam = bypassCache ? '&bypass_cache=true' : '';
    const url = `/api/accounts?limit=1000${cacheParam}`;
    return apiRequest(url, {
        method: 'GET',
        credentials: 'include'
    }, {
        errorMessage: 'Failed to fetch accounts.'
    });
}

export async function getAccount(accountID) {
    return apiRequest(`/api/accounts/${accountID}`, {
        method: 'GET',
        credentials: 'include'
    }, {
        errorMessage: 'Failed to fetch account.'
    });
}

export async function updateAccount(accountID, updates) {
    return apiRequest(`/api/accounts/${accountID}`, {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify(updates)
    }, {
        errorMessage: 'Failed to update account.'
    });
}

export async function deleteAccount(accountID) {
    return apiRequest(`/api/accounts/${accountID}`, {
        method: 'DELETE',
        credentials: 'include'
    }, {
        successMessage: 'Account removed successfully!',
        errorMessage: 'Failed to remove account.'
    });
}

export async function getCharacter(characterID) {
    return apiRequest(`/api/characters/${characterID}`, {
        method: 'GET',
        credentials: 'include'
    }, {
        errorMessage: 'Failed to fetch character.'
    });
}

export async function updateCharacter(characterID, updates) {
    return apiRequest(`/api/characters/${characterID}`, {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(updates),
        credentials: 'include'
    }, {
        errorMessage: 'Failed to update character.'
    });
}

export async function deleteCharacter(characterID) {
    return apiRequest(`/api/characters/${characterID}`, {
        method: 'DELETE',
        credentials: 'include'
    }, {
        errorMessage: 'Failed to remove character.'
    });
}

export async function refreshCharacter(characterID) {
    return apiRequest(`/api/characters/${characterID}/refresh`, {
        method: 'POST',
        credentials: 'include'
    }, {
        successMessage: 'Character data refreshed successfully!',
        errorMessage: 'Failed to refresh character data.'
    });
}

export async function addCharacter(account) {
    return apiRequest('/api/add-character', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ account }),
        credentials: 'include'
    });
}

export async function getAssociations() {
    return apiRequest(`/api/associations`, {
        method: 'GET',
        credentials: 'include'
    }, {
        errorMessage: 'Failed to fetch associations.'
    });
}

export async function associateCharacter(userId, charId, _userName, charName) {
    return apiRequest(`/api/associations`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ userId, characterId: charId, charName })
    }, {
        errorMessage: 'Association operation failed.'
    });
}

export async function unassociateCharacter(userId, charId, _userName, _charName) {
    return apiRequest(`/api/associations/${encodeURIComponent(userId)}/${encodeURIComponent(charId)}`, {
        method: 'DELETE',
        credentials: 'include'
    }, {
        errorMessage: 'Unassociation operation failed.'
    });
}
