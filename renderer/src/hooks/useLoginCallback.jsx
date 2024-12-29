// src/hooks/useLoginCallback.js
import { useCallback, useState } from 'react';
import { log } from "../utils/logger.jsx";
import { finalizelogin } from "../api/apiService.jsx";
import { useOAuthPoll } from "./useOAuthPoll";

/**
 * Custom hook to handle login callback logic via polling.
 *
 * @param {boolean} isAuthenticated    - Whether the user is currently authenticated.
 * @param {Function} loginRefresh      - Function to fetch user data after finalization.
 * @param {Function} setLoggedOut      - Setter for loggedOut state.
 * @param {Function} setIsAuthenticated - Setter for isAuthenticated state.
 * @returns {Function} startPoll
 */
export function useLoginCallback(
    isAuthenticated,
    loginRefresh,
    setLoggedOut,
    setIsAuthenticated
) {
    // We'll define the same pattern: a finalizeFn + afterFinalize

    const finalizeFn = useCallback(async (state) => {
        // If user is already authenticated, skip?
        if (isAuthenticated) {
            log("User already isAuthenticated; skip finalizing");
            return { success: true };
        }

        log("Calling finalize-login endpoint...", state);
        const resp = await finalizelogin(state);

        // Must return { success: true/false } format
        if (resp && resp.success) {
            return { success: true };
        }
        return { success: false };
    }, [isAuthenticated]);

    const afterFinalize = useCallback(async () => {
        // Once finalization is done, do loginRefresh
        log("Login finalization succeeded; calling loginRefresh...");
        await loginRefresh();
        // Then setIsAuthenticated
        setIsAuthenticated(true);
    }, [loginRefresh, setIsAuthenticated]);

    // Use the common hook
    const { startPoll, pollingActive } = useOAuthPoll(
        finalizeFn,
        afterFinalize,
        25,     // e.g. maxAttempts
        5000    // interval
    );

    // The original code also called setLoggedOut(false) when starting
     return useCallback((state) => {
        log("useLoginCallback -> startPoll called with state:", state);
        setLoggedOut(false);
        startPoll(state); // from the shared hook
    }, [setLoggedOut, startPoll]);
}
