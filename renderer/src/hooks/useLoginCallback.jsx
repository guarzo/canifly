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

    useEffect(() => {
        // If not active or we have no `state`, do nothing
        if (!pollingActive || !pollingState) return;

        log(`Starting finalize-login polling for state=${pollingState}`);
        let attempts = 0;
        let finalized = false;
        const maxAttempts = 25;

        const intervalId = setInterval(async () => {
            attempts++;
            log(`Polling attempt ${attempts} for state=${pollingState}`);

            // If user is already authenticated (defensive check), stop polling
            if (isAuthenticated) {
                log("User isAuthenticated=true, stopping polling.");
                clearInterval(intervalId);
                setPollingActive(false);
                return;
            }

            // If we exceeded max attempts, stop polling
            if (attempts > maxAttempts) {
                console.warn("Failed to detect login after multiple attempts, stopping polling.");
                clearInterval(intervalId);
                setPollingActive(false);
                return;
            }

            if (!finalized) {
                log("Calling finalize-login endpoint...");
                const finalizeResp = await finalizelogin(pollingState);
                console.log(finalizeResp)
                if (finalizeResp && finalizeResp.success) {
                    finalized = true;
                    log("Finalization succeeded, now trying to fetch data via loginRefresh...");

                    const success = await loginRefresh();
                    log("loginRefresh returned:", success);

                    if (success) {
                        log("Login finalized and data fetched! Setting isAuthenticated=true.");
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
                // Already finalized => keep trying loginRefresh
                log("Already finalized, retrying loginRefresh...");
                const success = await loginRefresh();
                log("loginRefresh returned:", success);

                if (success) {
                    log("Data fetched after finalization! Setting isAuthenticated=true.");
                    setIsAuthenticated(true);
                    clearInterval(intervalId);
                    setPollingActive(false);
                } else {
                    log("Still no data after finalization, retrying data fetch on next interval...");
                }
            }
        }, 5000);

        // Cleanup: clear the interval if the component unmounts
        return () => {
            clearInterval(intervalId);
        };
    }, [
        pollingActive,
        pollingState,
        isAuthenticated,
        loginRefresh,
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
