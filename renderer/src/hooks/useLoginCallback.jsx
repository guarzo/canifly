// useLoginCallback.js

import { useState, useEffect, useCallback } from 'react';
import { log } from "../utils/logger.jsx";
import { finalizelogin } from "../api/apiService.jsx";

/**
 * Custom hook to handle login callback logic via polling.
 *
 * @param {boolean} isAuthenticated - Whether the user is currently authenticated.
 * @param {boolean} loggedOut - Whether the user is currently logged out.
 * @param {Function} loginRefresh - Function to fetch user data after finalization.
 * @param {Function} setLoggedOut - Setter for loggedOut state.
 * @param {Function} setIsAuthenticated - Setter for isAuthenticated state.
 * @returns {Function} startPoll - A function that, when called with a `state` string,
 *                                initiates polling to finalize login and fetch user data.
 */
export function useLoginCallback(
    isAuthenticated,
    loggedOut,
    loginRefresh,
    setLoggedOut,
    setIsAuthenticated
) {
    // Track whether we're currently polling
    const [pollingActive, setPollingActive] = useState(false);

    // Store the OAuth state we want to finalize
    const [pollingState, setPollingState] = useState(null);

    const [finalized, setFinalized] = useState(false);


    useEffect(() => {
        if (!pollingActive || !pollingState) return;

        log(`Starting finalize-login polling for state=${pollingState}`);
        let attempts = 0;
        const maxAttempts = 25;

        const intervalId = setInterval(async () => {
            attempts++;
            log(`Polling attempt ${attempts} for state=${pollingState}`);

            if (isAuthenticated) {
                log("User isAuthenticated=true, stopping polling.");
                clearInterval(intervalId);
                setPollingActive(false);
                return;
            }

            if (attempts > maxAttempts) {
                console.warn("Failed to detect login after multiple attempts, stopping polling.");
                clearInterval(intervalId);
                setPollingActive(false);
                return;
            }

            if (!finalized) {
                log("Calling finalize-login endpoint...");
                const finalizeResp = await finalizelogin(pollingState);
                if (finalizeResp && finalizeResp.success) {
                    setFinalized(true); // <-- set state to true

                    log("Finalization succeeded, now trying to fetch data via loginRefresh...");
                    const success = await loginRefresh();
                    if (success) {
                        setIsAuthenticated(true);
                        clearInterval(intervalId);
                        setPollingActive(false);
                    } else {
                        log("Session set but data fetch failed, will retry data fetch on next interval...");
                    }
                } else {
                    log("Not ready yet, retrying finalize-login on next interval...");
                }
            } else {
                // If already finalized, keep retrying to fetch data
                log("Already finalized, retrying loginRefresh...");
                const success = await loginRefresh();
                if (success) {
                    setIsAuthenticated(true);
                    clearInterval(intervalId);
                    setPollingActive(false);
                } else {
                    log("Still no data after finalization, retrying data fetch on next interval...");
                }
            }
        }, 5000);

        return () => {
            clearInterval(intervalId);
        };
    }, [
        pollingActive,
        pollingState,
        isAuthenticated,
        loginRefresh,
        finalized,
        setIsAuthenticated,
    ]);


    return useCallback(
        (state) => {
            log("startPoll called with state:", state);
            setLoggedOut(false);    // ensure we start in a non-logged-out state
            setPollingState(state);
            setPollingActive(true); // triggers the effect above
        },
        [setLoggedOut]
    );
}
