// src/api/apiService.jsx

/**
 * API Service - Frontend API interface for CanIFly
 * 
 * RESTful Endpoints:
 * - GET    /api/accounts         - List all accounts
 * - GET    /api/accounts/:id     - Get specific account
 * - PATCH  /api/accounts/:id     - Update account (name, status, visibility)
 * - DELETE /api/accounts/:id     - Delete account
 * - GET    /api/characters/:id   - Get character details
 * - PATCH  /api/characters/:id   - Update character
 * - DELETE /api/characters/:id   - Remove character
 * - GET    /api/config           - Get configuration
 * - PATCH  /api/config           - Update configuration
 */

import { apiRequest } from './apiRequest';
import {isDev} from '../Config';

// New RESTful API functions
export async function getAccounts() {
    // Use a high limit to get all accounts - pagination handling can be added later if needed
    return apiRequest(`/api/accounts?limit=1000`, {
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

// Note: Dashboards are not a separate entity in the backend
// The dashboard page uses accounts and config data

// Get session status
export async function getSession() {
    return apiRequest(`/api/session`, {
        method: 'GET',
        credentials: 'include'
    }, {
        errorMessage: 'Failed to get session status.'
    });
}

// Check EVE configuration status
export async function checkEVEConfiguration() {
    return apiRequest(`/api/config/eve/status`, {
        method: 'GET',
        credentials: 'include'
    }, {
        errorMessage: 'Failed to check EVE configuration.'
    });
}

// Save EVE credentials
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

// Refresh character data from ESI
export async function refreshCharacter(characterID) {
    return apiRequest(`/api/characters/${characterID}/refresh`, {
        method: 'POST',
        credentials: 'include'
    }, {
        successMessage: 'Character data refreshed successfully!',
        errorMessage: 'Failed to refresh character data.'
    });
}

// Get all skill plans
export async function getSkillPlans() {
    return apiRequest(`/api/skill-plans`, {
        method: 'GET',
        credentials: 'include'
    }, {
        errorMessage: 'Failed to fetch skill plans.'
    });
}

// Get specific skill plan
export async function getSkillPlan(planName) {
    return apiRequest(`/api/skill-plans/${encodeURIComponent(planName)}`, {
        method: 'GET',
        credentials: 'include'
    }, {
        errorMessage: 'Failed to fetch skill plan.'
    });
}

// Create skill plan
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

// Copy skill plan
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

// EVE Data endpoints
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

// Additional API functions

export async function logout() {
    return apiRequest('/api/logout', {
        method: 'POST',
        credentials: 'include'
    }, {
        errorMessage: 'Failed to log out.'
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

// deleteCharacter is the primary function for removing characters
export async function deleteCharacter(characterID) {
    return apiRequest(`/api/characters/${characterID}`, {
        method: 'DELETE',
        credentials: 'include'
    }, {
        errorMessage: 'Failed to remove character.'
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
    // Check if plan exists by trying to get it first
    // If it exists, use PUT to update; if not, use POST to create
    try {
        const response = await fetch(`/api/skill-plans/${encodeURIComponent(planName)}`, {
            method: 'GET',
            credentials: 'include'
        });
        
        if (response.ok) {
            // Plan exists, use PUT to update
            return apiRequest(`/api/skill-plans/${encodeURIComponent(planName)}`, {
                method: 'PUT',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ content: planContents }),
                credentials: 'include'
            }, {
                successMessage: 'Skill Plan Updated!',
                errorMessage: 'Failed to update skill plan.'
            });
        } else {
            // Plan doesn't exist, use POST to create
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
    } catch (error) {
        // Default to create if there's an error checking
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

// Get all associations
export async function getAssociations() {
    return apiRequest(`/api/associations`, {
        method: 'GET',
        credentials: 'include'
    }, {
        errorMessage: 'Failed to fetch associations.'
    });
}

export async function associateCharacter(userId, charId, userName, charName) {
    return apiRequest(`/api/associations`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ userId, characterId: charId })
    }, {
        errorMessage: 'Association operation failed.'
    });
}

export async function unassociateCharacter(userId, charId, userName, charName) {
    return apiRequest(`/api/associations/${encodeURIComponent(userId)}/${encodeURIComponent(charId)}`, {
        method: 'DELETE',
        credentials: 'include'
    }, {
        errorMessage: 'Unassociation operation failed.'
    });
}

// Fuzzworks functions
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

export async function deleteSkillPlan(planName) {
    return apiRequest(`/api/skill-plans/${encodeURIComponent(planName)}`, {
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
        method: 'GET',
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
    
    // EVE Data functions
    getEveSkillPlans,
    getEveProfiles,
    getEveConversions,
    
    // Authentication
    getSession,
    logout,
    initiateLogin,
    finalizelogin,
    
    // Skill Plan functions
    getSkillPlans,
    getSkillPlan,
    createSkillPlan,
    copySkillPlan,
    deleteSkillPlan,
    saveSkillPlan,
    
    // Association functions
    getAssociations,
    associateCharacter,
    unassociateCharacter,
    
    // Fuzzworks functions
    getFuzzworksStatus,
    updateFuzzworks,
    
    // Other functions
    addCharacter,
    chooseDefaultDirectory,
    chooseSettingsDir,
    backupDirectory,
    resetToDefaultDirectory,
    syncSubdirectory,
    syncAllSubdirectories,
    saveUserSelections
};

