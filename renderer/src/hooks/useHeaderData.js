// src/hooks/useHeaderData.js
//
// Owns the Header's data dependencies and the OAuth-flow "add character"
// handler that polls /api/finalize-login until the EVE SSO redirect
// completes (or times out).
import { useCallback, useEffect, useRef } from 'react';
import { toast } from 'react-toastify';
import { useAuth } from './useAuth';
import { useAppData } from './useAppData';
import { useAsyncOperation } from './useAsyncOperation';

const POLL_INTERVAL_MS = 2000;
const POLL_MAX_ATTEMPTS = 30; // 30 * 2s = 1 minute

const openExternal = (url) => {
    if (window.electronAPI?.openExternal) {
        window.electronAPI.openExternal(url);
    } else {
        window.open(url, '_blank');
    }
};

export function useHeaderData() {
    const { isAuthenticated, logout } = useAuth();
    const { refreshData } = useAppData();
    const { execute: executeRefresh, isLoading: isRefreshing } = useAsyncOperation();
    const { execute: executeAddCharacter } = useAsyncOperation();

    // Tracks the active finalize-login poll so we can clear it on unmount and
    // avoid the "started a poll, navigated away, callbacks still fire" leak.
    const pollIntervalRef = useRef(null);
    const clearPoll = useCallback(() => {
        if (pollIntervalRef.current != null) {
            clearInterval(pollIntervalRef.current);
            pollIntervalRef.current = null;
        }
    }, []);
    useEffect(() => clearPoll, [clearPoll]);

    const handleRefreshClick = useCallback(
        () => executeRefresh(() => refreshData(), { showToast: false }),
        [executeRefresh, refreshData],
    );

    const handleCloseWindow = useCallback(() => {
        if (window.electronAPI?.closeWindow) {
            window.electronAPI.closeWindow();
        }
    }, []);

    // Poll /finalize-login until the OAuth handshake completes, errors, or times out.
    const pollFinalizeLogin = useCallback((state) => {
        clearPoll();
        let attempts = 0;
        pollIntervalRef.current = setInterval(async () => {
            attempts++;
            if (attempts > POLL_MAX_ATTEMPTS) {
                clearPoll();
                toast.warning('Character add timeout. Please try again.');
                return;
            }
            try {
                const { finalizelogin } = await import('../api/authApi');
                const result = await finalizelogin(state);
                if (result?.success) {
                    clearPoll();
                    toast.success('Character added');
                    // The backend also pushes a WebSocket update that triggers
                    // a refresh, but we refresh explicitly here so the UI
                    // updates even when the socket isn't connected.
                    await refreshData();
                } else if (result?.error === 'state_not_found') {
                    clearPoll();
                    toast.error('Session expired. Please try again.');
                }
            } catch (error) {
                // Polling failure is not fatal — keep retrying until maxAttempts.
                void error;
            }
        }, POLL_INTERVAL_MS);
    }, [clearPoll, refreshData]);

    const handleAddCharacterSubmit = useCallback(async (account) => {
        const result = await executeAddCharacter(async () => {
            const { addCharacter } = await import('../api/accountsApi');
            return addCharacter(account);
        }, { showToast: false });

        if (!result) {
            // executeAddCharacter returns null when the API call failed; the
            // underlying apiClient already surfaced its own error toast, so
            // here we just bail out without claiming success.
            return;
        }

        if (result.redirectURL) {
            openExternal(result.redirectURL);
            toast.info('Complete authorization in your browser');
            if (result.state) pollFinalizeLogin(result.state);
        } else {
            toast.success('Character added');
            await refreshData();
        }
    }, [executeAddCharacter, pollFinalizeLogin, refreshData]);

    return {
        isAuthenticated,
        logout,
        isRefreshing,
        handleRefreshClick,
        handleCloseWindow,
        handleAddCharacterSubmit,
    };
}
