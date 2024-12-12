// api/apiService.jsx

import { apiRequest } from './apiRequest';
import { normalizeAppData } from '../utils/dataNormalizer';

export async function getAppData(backEndURL) {
    const response = await apiRequest(`${backEndURL}/api/app-data`, {
        credentials: 'include'
    }, {
        errorMessage: 'Failed to load app data.'
    });
    return response ? normalizeAppData(response) : null;
}

export async function getAppDataNoCache(backEndURL) {
    const response = await apiRequest(`${backEndURL}/api/app-data-no-cache`, {
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
        successMessage: 'Logged out successfully!',
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
        successMessage: 'Account status toggled successfully!',
        errorMessage: 'Failed to toggle account status.'
    });
}

export async function updateCharacter(characterID, updates) {
    return apiRequest('/api/update-character', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ characterID, updates }),
        credentials: 'include'
    }, {
        successMessage: 'Character updated successfully!',
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
        successMessage: 'Character removed successfully!',
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
        successMessage: 'Account name updated successfully!',
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
    return apiRequest('/api/add-character', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ account }),
        credentials: 'include'
    }, {
        errorMessage: 'An error occurred while adding character.',
        onSuccess: (data) => {
            if (data.redirectURL) {
                // If running in dev mode, redirect internally; otherwise, open externally
                if (window.isDev) {
                    window.location.href = data.redirectURL;
                } else {
                    window.electronAPI.openExternal(data.redirectURL);
                }
            } else {
                // If no redirect URL returned, consider showing a toast or handling the error
            }
        }
    });
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


export async function saveUserSelections(newSelections, backEndURL) {
    return apiRequest(`${backEndURL}/api/save-user-selections`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify(newSelections),
    }, {
        errorMessage: 'Failed to save user selections.',
    });
}

export async function syncSubdirectory(profile, userId, charId, backEndURL) {
    return apiRequest(`${backEndURL}/api/sync-subdirectory`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ profile, userId, charId })
    }, {
        errorMessage: 'Sync operation failed.'
    });
}

export async function syncAllSubdirectories(profile, userId, charId, backEndURL) {
    return apiRequest(`${backEndURL}/api/sync-all-subdirectories`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ profile, userId, charId })
    }, {
        errorMessage: 'Sync-All operation failed.'
    });
}

export async function chooseSettingsDir(directory, backEndURL) {
    return apiRequest(`${backEndURL}/api/choose-settings-dir`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ directory }),
    }, {
        errorMessage: 'Failed to choose settings directory.'
    });
}

export async function backupDirectory(targetDir, backupDir, backEndURL) {
    return apiRequest(`${backEndURL}/api/backup-directory`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ targetDir, backupDir }),
    }, {
        errorMessage: 'Backup operation failed.'
    });
}

export async function resetToDefaultDirectory(backEndURL) {
    return apiRequest(`${backEndURL}/api/reset-to-default-directory`, {
        method: 'POST',
        credentials: 'include',
    }, {
        errorMessage: 'Failed to reset directory.'
    });
}

export async function associateCharacter(userId, charId, userName, charName, backEndURL) {
    return apiRequest(`${backEndURL}/api/associate-character`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ userId, charId, userName, charName })
    }, {
        errorMessage: 'Association operation failed.'
    });
}

export async function unassociateCharacter(userId, charId, userName, charName, backEndURL) {
    return apiRequest(`${backEndURL}/api/unassociate-character`, {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        credentials: 'include',
        body: JSON.stringify({ userId, charId, userName, charName })
    }, {
        errorMessage: 'Unassociation operation failed.'
    });
}

export async function deleteSkillPlan(planName, backEndURL) {
    return apiRequest(`${backEndURL}/api/delete-skill-plan?planName=${encodeURIComponent(planName)}`, {
        method: 'DELETE',
        credentials: 'include',
    }, {
        errorMessage: 'Failed to delete skill plan.'
    });
}

export async function initiateLogin(account, backEndURL) {
    // Removed isDev parameter; we can handle isDev in the component if needed
    return apiRequest(`${backEndURL}/api/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ account }),
        credentials: 'include',
    }, {
        errorMessage: 'Failed to initiate login.'
    });
}
