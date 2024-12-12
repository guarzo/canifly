import { apiRequest } from '../utils/apiRequest';

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
