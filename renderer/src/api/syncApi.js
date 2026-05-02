// src/api/syncApi.js
//
// Settings/profile sync, backup, and directory-selection endpoints.
import { apiRequest } from './apiClient';

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
