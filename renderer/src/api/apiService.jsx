// src/api/apiService.jsx

/**
 * API Service - Frontend API interface for CanIFly
 * 
 * This file contains both new RESTful API functions and legacy functions.
 * The backend supports both patterns during the transition period.
 * 
 * NEW RESTful Endpoints:
 * - GET    /api/accounts         - List all accounts
 * - GET    /api/accounts/:id     - Get specific account
 * - PATCH  /api/accounts/:id     - Update account (name, status, visibility)
 * - DELETE /api/accounts/:id     - Delete account
 * - GET    /api/characters/:id   - Get character details
 * - PATCH  /api/characters/:id   - Update character
 * - DELETE /api/characters/:id   - Remove character
 * - GET    /api/config           - Get configuration
 * - PATCH  /api/config           - Update configuration
 * 
 * Legacy endpoints are maintained for backward compatibility.
 */

import { apiRequest } from './apiRequest';
import { normalizeAppData } from '../utils/dataNormalizer';
import {isDev} from '../Config';

// New RESTful API functions
export async function getAccounts() {
    return apiRequest(`/api/accounts`, {
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

export async function getConfig() {
    return apiRequest(`/api/config`, {
        method: 'GET',
        credentials: 'include'
    }, {
        errorMessage: 'Failed to fetch configuration.'
    });
}

// Update configuration
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

// Get dashboards
export async function getDashboards() {
    return apiRequest(`/api/dashboards`, {
        method: 'GET',
        credentials: 'include'
    }, {
        errorMessage: 'Failed to get dashboards.'
    });
}

// Get session status
export async function getSession() {
    return apiRequest(`/api/session`, {
        method: 'GET',
        credentials: 'include'
    }, {
        errorMessage: 'Failed to get session status.'
    });
}

// Update account (general purpose)
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

// Create skill plan
export async function createSkillPlan(planName, planContents) {
    return apiRequest(`/api/skill-plans`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ planName, planContents })
    }, {
        errorMessage: 'Failed to create skill plan.'
    });
}

// Copy skill plan
export async function copySkillPlan(sourcePlanName, targetPlanName) {
    return apiRequest(`/api/skill-plans/copy`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ sourcePlanName, targetPlanName })
    }, {
        errorMessage: 'Failed to copy skill plan.'
    });
}

// Additional API functions

export async function logout() {
    return apiRequest('/api/logout', {
        method: 'POST',
        credentials: 'include'
    }, {
        errorMessage: 'Failed to log out.'
    });
}

export async function toggleAccountStatus(accountID) {
    // Use new RESTful endpoint
    return apiRequest(`/api/accounts/${accountID}`, {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ isActive: true }) // Toggle active status
    }, {
        errorMessage: 'Failed to toggle account status.'
    });
}

export async function toggleAccountVisibility(accountID) {
    // Use new RESTful endpoint
    return apiRequest(`/api/accounts/${accountID}`, {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ isVisible: true }) // Toggle visibility
    }, {
        errorMessage: 'Failed to toggle account visibility.'
    });
}


export async function updateCharacter(characterID, updates) {
    // Use new RESTful endpoint
    return apiRequest(`/api/characters/${characterID}`, {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(updates),
        credentials: 'include'
    }, {
        errorMessage: 'Failed to update character.'
    });
}

export async function removeCharacter(characterID) {
    // Use new RESTful endpoint
    return apiRequest(`/api/characters/${characterID}`, {
        method: 'DELETE',
        credentials: 'include'
    }, {
        errorMessage: 'Failed to remove character.'
    });
}

// Alias for removeCharacter to match RESTful naming
export const deleteCharacter = removeCharacter;

export async function updateAccountName(accountID, newName) {
    // Use new RESTful endpoint
    return apiRequest(`/api/accounts/${accountID}`, {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name: newName }),
        credentials: 'include'
    }, {
        errorMessage: 'Failed to update account name.'
    });
}

export async function removeAccount(accountName, accountID = null) {
    // If accountID is provided, use new RESTful endpoint
    if (accountID) {
        return deleteAccount(accountID);
    }
    
    // Otherwise, use legacy endpoint
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
    // Use new RESTful endpoint
    return apiRequest(`/api/config`, {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ userSelections: newSelections }),
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
    // Use new RESTful endpoint
    return apiRequest(`/api/config`, {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ settingsDir: directory }),
    }, {
        errorMessage: 'Failed to choose settings directory.'
    });
}

// Alias for backward compatibility
export const chooseSettingsDirectory = chooseSettingsDir;

export async function chooseDefaultDirectory(directory) {
    return apiRequest(`/api/config`, {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ defaultDir: directory }),
    }, {
        errorMessage: 'Failed to choose default directory.'
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
    // Use new RESTful endpoint to set empty string (will use default)
    return apiRequest(`/api/config`, {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ settingsDir: '' }),
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

// Default export for backward compatibility
export default {
    // New RESTful functions
    getAccounts,
    getAccount,
    updateAccount,
    deleteAccount,
    getCharacter,
    updateCharacter,
    deleteCharacter,
    getConfig,
    updateConfig,
    getDashboards,
    
    // Authentication
    getSession,
    logout,
    initiateLogin,
    finalizelogin,
    
    // Other functions
    addCharacter,
    createSkillPlan,
    copySkillPlan,
    deleteSkillPlan,
    chooseDefaultDirectory,
    chooseSettingsDirectory,
    backupDirectory,
    resetToDefaultDirectory,
    associateCharacter,
    unassociateCharacter,
    
    // Utility
    normalizeAppData
};

