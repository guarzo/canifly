// src/api/apiService.jsx

import { apiRequest } from './apiRequest';
import { normalizeAppData } from '../utils/dataNormalizer';
import {isDev} from '../Config';

export async function getAppData() {
    const response = await apiRequest(`/api/app-data`, {
        credentials: 'include'
    }, {
        errorMessage: isDev ? 'Failed to load app data.' : undefined
    });
    return response ? normalizeAppData(response) : null;
}

export async function getAppDataNoCache() {
    const response = await apiRequest(`/api/app-data-no-cache`, {
        credentials: 'include'
    }, {
        errorMessage: 'Failed to load data.'
    });
    return response ? normalizeAppData(response) : null;
}

export async function logout() {
    return apiRequest('/api/logout', {
        method: 'POST',
        credentials: 'include'
    }, {
        errorMessage: 'Failed to log out.'
    });
}

export async function toggleAccountStatus(accountID) {
    return apiRequest('/api/toggle-account-status', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ accountID })
    }, {
        errorMessage: 'Failed to toggle account status.'
    });
}

export async function toggleAccountVisibility(accountID) {
    return apiRequest('/api/toggle-account-visibility', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ accountID })
    }, {
        errorMessage: 'Failed to toggle account visibility.'
    });
}


export async function updateCharacter(characterID, updates) {
    return apiRequest('/api/update-character', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ characterID, updates }),
        credentials: 'include'
    }, {
        errorMessage: 'Failed to update character.'
    });
}

export async function removeCharacter(characterID) {
    return apiRequest('/api/remove-character', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ characterID }),
        credentials: 'include'
    }, {
        errorMessage: 'Failed to remove character.'
    });
}

export async function updateAccountName(accountID, newName) {
    return apiRequest('/api/update-account-name', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ accountID, accountName: newName }),
        credentials: 'include'
    }, {
        errorMessage: 'Failed to update account name.'
    });
}

export async function removeAccount(accountName) {
    return apiRequest('/api/remove-account', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ accountName }),
        credentials: 'include'
    }, {
        successMessage: 'Account removed successfully!',
        errorMessage: 'Failed to remove account.'
    });
}

export async function addCharacter(account) {
    return await apiRequest(
        '/api/add-character',
        {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ account }),
            credentials: 'include'
        },
    );
}



export async function saveSkillPlan(planName, planContents) {
    return apiRequest('/api/save-skill-plan', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name: planName, contents: planContents }),
        credentials: 'include'
    }, {
        successMessage: 'Skill Plan Saved!',
        errorMessage: 'Failed to save skill plan.'
    });
}


export async function saveUserSelections(newSelections) {
    return apiRequest(`/api/save-user-selections`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify(newSelections),
    }, {
        errorMessage: 'Failed to save user selections.',
    });
}

export async function syncSubdirectory(profile, userId, charId) {
    return apiRequest(`/api/sync-subdirectory`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ subDir: profile, userId, charId })
    }, {
        errorMessage: 'Sync operation failed.'
    });
}

export async function syncAllSubdirectories(profile, userId, charId) {
    return apiRequest(`/api/sync-all-subdirectories`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ subDir: profile, userId, charId })
    }, {
        errorMessage: 'Sync-All operation failed.'
    });
}

export async function chooseSettingsDir(directory) {
    return apiRequest(`/api/choose-settings-dir`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ directory }),
    }, {
        errorMessage: 'Failed to choose settings directory.'
    });
}

export async function backupDirectory(targetDir, backupDir) {
    return apiRequest(`/api/backup-directory`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ targetDir, backupDir }),
    }, {
        errorMessage: 'Backup operation failed.'
    });
}

export async function resetToDefaultDirectory() {
    return apiRequest(`/api/reset-to-default-directory`, {
        method: 'POST',
        credentials: 'include',
    }, {
        errorMessage: 'Failed to reset directory.'
    });
}

export async function associateCharacter(userId, charId, userName, charName) {
    return apiRequest(`/api/associate-character`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ userId, charId, userName, charName })
    }, {
        errorMessage: 'Association operation failed.'
    });
}

export async function unassociateCharacter(userId, charId, userName, charName) {
    return apiRequest(`/api/unassociate-character`, {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        credentials: 'include',
        body: JSON.stringify({ userId, charId, userName, charName })
    }, {
        errorMessage: 'Unassociation operation failed.'
    });
}

export async function deleteSkillPlan(planName) {
    return apiRequest(`/api/delete-skill-plan?planName=${encodeURIComponent(planName)}`, {
        method: 'DELETE',
        credentials: 'include',
    }, {
        errorMessage: 'Failed to delete skill plan.'
    });
}

export async function initiateLogin(account) {
    // Removed isDev parameter; we can handle isDev in the component if needed
    return apiRequest(`/api/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ account }),
        credentials: 'include',
    }, {
        errorMessage: 'Failed to initiate login.'
    });
}


export async function finalizelogin(state) {
    return apiRequest(`/api/finalize-login?state=${state}`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
    }, {
        disableErrorToast: true,
    });
}

