// src/api/authApi.js
//
// Authentication / session endpoints.
import { apiRequest } from './apiClient';

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

export async function getSession() {
    return apiRequest(`/api/session`, {
        method: 'GET',
        credentials: 'include'
    }, {
        errorMessage: 'Failed to get session status.'
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
