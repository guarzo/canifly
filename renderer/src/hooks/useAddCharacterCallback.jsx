// src/hooks/useAddCharacterCallback.js

import { useCallback } from 'react';
import { log } from "../utils/logger.jsx";
import { finalizelogin } from "../api/apiService.jsx";

/**
 * Custom hook to finalize the "Add Character" flow via OAuth.
 * This doesn't rely on isAuthenticated, because the user might
 * already be logged in. Instead, it keeps polling finalizelogin()
 * until that call succeeds and we can refresh data successfully.
 *
 * @param {Function} loginRefresh - Function to fetch updated user data (including the new character).
 * @returns {Function} addCharacterCallback - A function that starts the polling when invoked with a state.
 */
export function useAddCharacterCallback(loginRefresh) {
    return useCallback((state) => {
        log("useAddCharacterCallback invoked with state:", state);

        let attempts = 0;
        const maxAttempts = 5;
        let finalized = false;

        const interval = setInterval(async () => {
            attempts++;
            log(`AddChar interval #${attempts}, finalized=${finalized}`);

            if (attempts > maxAttempts) {
                console.warn("Gave up waiting after multiple attempts. Clearing interval.");
                loginRefresh()
                clearInterval(interval);
                return;
            }

            if (!finalized) {
                log("Calling finalizelogin for add-character...state: ", state);
                const resp = await finalizelogin(state);
                if (resp.ok) {
                    finalized = true;
                    log("Finalization on server done, now let's try loginRefresh to load new data...");
                    const success = await loginRefresh();
                    if (success) {
                        log("Data fetch success! Clearing interval; new character should be available now.");
                        clearInterval(interval);
                    } else {
                        log("Data fetch failed or not complete yet, continuing to poll...");
                    }
                } else {
                    log("Not ready yet, continuing to poll finalizelogin...");
                }
            } else {
                // Already finalized, so we just keep trying to refresh data
                // until the new character is definitely loaded into the client.
                log("Already finalized, retrying loginRefresh...");
                const success = await loginRefresh();
                if (success) {
                    log("Successfully fetched new data, clearing interval.");
                    clearInterval(interval);
                } else {
                    log("Still no data, continuing to poll...");
                }
            }
        }, 5000);

        log("useAddCharacterCallback: interval created");
    }, [loginRefresh]);
}
