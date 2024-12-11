// src/hooks/useLoginCallback.js
import { useCallback } from 'react';

/**
 * Custom hook to handle login callback logic.
 *
 * @param {boolean} isAuthenticated - Whether the user is currently authenticated.
 * @param {boolean} loggedOut - Whether the user is currently logged out.
 * @param {Function} loginRefresh - Function to try fetching user data after finalization.
 * @param {Function} setLoggedOut - Setter for loggedOut state.
 * @param {Function} setIsAuthenticated - Setter for isAuthenticated state.
 * @param {string} backEndURL - Base URL for the backend.
 * @returns {Function} logInCallBack - A function that, when called with a state string, starts the login finalization polling.
 */
export function useLoginCallback(isAuthenticated, loggedOut, loginRefresh, setLoggedOut, setIsAuthenticated, backEndURL) {
    // We define logInCallBack using the values currently in this closure.
    // When this hook is re-run (due to state changes), a new logInCallBack
    // will be created with updated references to isAuthenticated, loggedOut, etc.
    return useCallback((state) => {
        console.log("logInCallBack called with state:", state);
        setLoggedOut(false); // Ensure we start in a non-logged-out state

        let attempts = 0;
        const maxAttempts = 6;
        let finalized = false;

        const interval = setInterval(async () => {
            // Here we rely on the values captured in the closure.
            // These values reflect the state at the time this hook last ran.
            // If isAuthenticated or loggedOut change, useLoginCallback is re-run,
            // recreating logInCallBack with updated values.
            console.log(`Interval tick: attempts=${attempts}, isAuthenticated=${isAuthenticated}, finalized=${finalized}, loggedOut=${loggedOut}`);

            if (isAuthenticated) {
                console.log("User authenticated, clearing interval");
                clearInterval(interval);
                return;
            }

            attempts++;
            if (attempts > maxAttempts) {
                console.warn('Failed to detect login after multiple attempts, clearing interval.');
                clearInterval(interval);
                return;
            }

            if (!finalized) {
                console.log("Calling finalize-login endpoint...");
                const finalizeResp = await fetch(`${backEndURL}/api/finalize-login?state=${state}`, {
                    credentials: 'include'
                });

                if (finalizeResp.ok) {
                    finalized = true;
                    console.log("Finalization succeeded, now trying to fetch data via loginRefresh...");
                    const success = await loginRefresh();
                    console.log("loginRefresh returned:", success);
                    if (success) {
                        console.log("Login finalized and data fetched! Setting isAuthenticated=true and clearing interval");
                        setIsAuthenticated(true);
                        clearInterval(interval);
                    } else {
                        console.log("Session set but data fetch failed, will retry data fetch on next interval...");
                    }
                } else {
                    console.log("Not ready yet, retrying finalize-login...");
                }
            } else {
                // Already finalized, just try loginRefresh again
                console.log("Already finalized, retrying loginRefresh...");
                const success = await loginRefresh();
                console.log("loginRefresh returned:", success);
                if (success) {
                    console.log("Data fetched after finalization! Setting isAuthenticated=true and clearing interval");
                    setIsAuthenticated(true);
                    clearInterval(interval);
                } else {
                    console.log("Still no data after finalization, retrying data fetch on next interval...");
                }
            }
        }, 5000);

        console.log("logInCallBack: interval created");
    }, [isAuthenticated, loggedOut, loginRefresh, setLoggedOut, setIsAuthenticated, backEndURL]);
}
