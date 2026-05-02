// src/api/apiService.jsx
//
// Transitional façade. The per-domain API modules are the new home for
// these functions. This file re-exports them so existing callers keep
// working until they are migrated; it will be deleted at the end of the
// frontend-module-cleanup migration.

import {
    getAccounts,
    getAccount,
    updateAccount,
    deleteAccount,
    getCharacter,
    updateCharacter,
    deleteCharacter,
    refreshCharacter,
    addCharacter,
    getAssociations,
    associateCharacter,
    unassociateCharacter,
} from './accountsApi';
import { apiRequest } from './apiClient';

export {
    getAccounts,
    getAccount,
    updateAccount,
    deleteAccount,
    getCharacter,
    updateCharacter,
    deleteCharacter,
    refreshCharacter,
    addCharacter,
    getAssociations,
    associateCharacter,
    unassociateCharacter,
};

// ─── Config ─────────────────────────────────────────────────────────────
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

export async function getSession() {
    return apiRequest(`/api/session`, {
        method: 'GET',
        credentials: 'include'
    }, {
        errorMessage: 'Failed to get session status.'
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

// ─── Skill plans ────────────────────────────────────────────────────────
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
        } else {
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

// ─── EVE static data ────────────────────────────────────────────────────
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

// ─── Auth ───────────────────────────────────────────────────────────────
export async function logout() {
    return apiRequest('/api/logout', {
        method: 'POST',
        credentials: 'include'
    }, {
        errorMessage: 'Failed to log out.'
    });
}

export async function initiateLogin(account, rememberMe = false) {
    return apiRequest(`/api/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ account, rememberMe }),
        credentials: 'include',
    }, {
        errorMessage: 'Failed to initiate login.'
    });
}

export async function finalizelogin(state, rememberMe = false) {
    const params = new URLSearchParams({ state });
    if (rememberMe) {
        params.append('remember', 'true');
    }
    return apiRequest(`/api/finalize-login?${params.toString()}`, {
        method: 'GET',
        credentials: 'include',
    }, {
        disableErrorToast: true,
    });
}

export async function validateSession() {
    return apiRequest(`/api/session/validate`, {
        method: 'GET',
        credentials: 'include'
    }, {
        disableErrorToast: true
    });
}

export async function refreshSession() {
    return apiRequest(`/api/session/refresh`, {
        method: 'POST',
        credentials: 'include'
    }, {
        errorMessage: 'Failed to refresh session.'
    });
}

export async function getActiveSessions() {
    return apiRequest(`/api/session/active`, {
        method: 'GET',
        credentials: 'include'
    }, {
        errorMessage: 'Failed to get active sessions.'
    });
}

// ─── Sync / directory ───────────────────────────────────────────────────
export async function saveUserSelections(newSelections) {
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
    return apiRequest(`/api/config`, {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ settingsDir: '' }),
    }, {
        errorMessage: 'Failed to reset directory.'
    });
}

// ─── Fuzzworks ──────────────────────────────────────────────────────────
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

// Default export for legacy callers (stores etc.). Will be removed
// once the migration is complete.
export default {
    getAccounts,
    getAccount,
    updateAccount,
    deleteAccount,
    getCharacter,
    updateCharacter,
    deleteCharacter,
    refreshCharacter,
    addCharacter,
    getAssociations,
    associateCharacter,
    unassociateCharacter,
    getConfig,
    updateConfig,
    getSession,
    validateSession,
    refreshSession,
    getActiveSessions,
    checkEVEConfiguration,
    saveEVECredentials,
    getEveSkillPlans,
    getEveProfiles,
    getEveConversions,
    logout,
    initiateLogin,
    finalizelogin,
    getSkillPlans,
    getSkillPlan,
    createSkillPlan,
    copySkillPlan,
    deleteSkillPlan,
    saveSkillPlan,
    chooseDefaultDirectory,
    chooseSettingsDir,
    backupDirectory,
    resetToDefaultDirectory,
    syncSubdirectory,
    syncAllSubdirectories,
    saveUserSelections,
    getFuzzworksStatus,
    updateFuzzworks,
};
