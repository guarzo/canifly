// src/hooks/useLoginCallback.js
export function useLoginCallback() {
    /**
     * Custom hook to handle login callback logic.
     *
     * This hook sets an interval to repeatedly attempt login finalization until the user is authenticated
     * or until it gives up after a certain number of attempts.
     *
     * @param {string} state - The state parameter from the OAuth flow (or similar).
     * @param {boolean} isAuthenticated - Whether the user is currently authenticated.
     * @param {boolean} loggedOut - Whether the user is currently logged out.
     * @param {Function} loginRefresh - Function to try fetching user data after finalization.
     * @param {Function} setLoggedOut - Setter for loggedOut.
     * @param {Function} setIsAuthenticated - Setter for isAuthenticated.
     */
    return function logInCallBack(state, { isAuthenticated, loggedOut, loginRefresh, setLoggedOut, setIsAuthenticated, backEndURL }) {
        console.log("logInCallBack called with state:", state);
        setLoggedOut(false);

        let attempts = 0;
        const maxAttempts = 6;
        let finalized = false;

        const interval = setInterval(async () => {
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
    };
}
