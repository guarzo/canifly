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
import {
    getSkillPlans,
    getSkillPlan,
    createSkillPlan,
    copySkillPlan,
    deleteSkillPlan,
    saveSkillPlan,
} from './skillPlansApi';
import {
    saveUserSelections,
    syncSubdirectory,
    syncAllSubdirectories,
    chooseSettingsDir,
    chooseDefaultDirectory,
    backupDirectory,
    resetToDefaultDirectory,
} from './syncApi';
import {
    getConfig,
    updateConfig,
    checkEVEConfiguration,
    saveEVECredentials,
} from './configApi';
import {
    getEveSkillPlans,
    getEveProfiles,
    getEveConversions,
} from './esiApi';
import { getFuzzworksStatus, updateFuzzworks } from './fuzzworksApi';
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
    getSkillPlans,
    getSkillPlan,
    createSkillPlan,
    copySkillPlan,
    deleteSkillPlan,
    saveSkillPlan,
    saveUserSelections,
    syncSubdirectory,
    syncAllSubdirectories,
    chooseSettingsDir,
    chooseDefaultDirectory,
    backupDirectory,
    resetToDefaultDirectory,
    getConfig,
    updateConfig,
    checkEVEConfiguration,
    saveEVECredentials,
    getEveSkillPlans,
    getEveProfiles,
    getEveConversions,
    getFuzzworksStatus,
    updateFuzzworks,
};

// ─── Config (re-exported above from configApi) ─────────────────────────

export async function getSession() {
    return apiRequest(`/api/session`, {
        method: 'GET',
        credentials: 'include'
    }, {
        errorMessage: 'Failed to get session status.'
    });
}

// ─── Skill plans (re-exported above from skillPlansApi) ─────────────────

// ─── EVE static data (re-exported above from esiApi) ──────────────────

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

// ─── Sync / directory (re-exported above from syncApi) ─────────────────

// ─── Fuzzworks (re-exported above from fuzzworksApi) ──────────────────

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
