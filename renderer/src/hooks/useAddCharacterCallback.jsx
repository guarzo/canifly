// src/hooks/useAddCharacterCallback.js
import { useCallback } from 'react';
import { log } from '../utils/logger.jsx';
import { finalizelogin } from '../api/apiService.jsx';
import { useOAuthPoll } from './useOAuthPoll.jsx';

/**
 * Custom hook to finalize the "Add Character" flow via OAuth polling.
 *
 * @param {Function} loginRefresh - Function to fetch updated user data (including the newly added character).
 * @returns {Function} startAddCharacterPoll
 */
export function useAddCharacterCallback(loginRefresh) {
    // 1) Provide finalizeFn
    const finalizeFn = useCallback(async (state) => {
        log('Calling finalizelogin for add-character... state:', state);
        return finalizelogin(state); // should return { success: boolean }
    }, []);

    // 2) Provide afterFinalize
    const afterFinalize = useCallback(async () => {
        log('Finalization complete; now fetching updated data...');
        await loginRefresh();
    }, [loginRefresh]);

    // 3) Use the common hook
    const { startPoll, pollingActive, pollingState } = useOAuthPoll(
        finalizeFn,
        afterFinalize,
        5,      // maxAttempts
        5000    // intervalMs
    );

    // Return just the startPoll function that we name "startAddCharacterPoll"
    return startPoll;
}
