// src/hooks/useOAuthPoll.js
import { useState, useEffect, useCallback, useRef } from 'react';
import { log } from '../utils/logger.jsx';

/**
 * Generic hook to handle OAuth polling logic.
 *
 * @param {Function} finalizeFn     - Function that tries to finalize OAuth. Should return { success: true } on success.
 * @param {Function} afterFinalize  - Function to call once finalize is successful (like loginRefresh).
 * @param {number}   maxAttempts    - Maximum polling attempts before giving up.
 * @param {number}   intervalMs     - Interval in milliseconds between polls.
 * @returns {Object} {
 *   startPoll: (stateString) => void,
 *   pollingActive: boolean,
 *   pollingState: string | null,
 * }
 */
export function useOAuthPoll(
    finalizeFn,
    afterFinalize,
    maxAttempts = 5,
    intervalMs = 5000
) {
    const [pollingActive, setPollingActive] = useState(false);
    const [pollingState, setPollingState] = useState(null);

    // Track if we've already finalized
    const finalizedRef = useRef(false);

    useEffect(() => {
        if (!pollingActive || !pollingState) return;

        log(`useOAuthPoll: starting for state=${pollingState}`);
        let attempts = 0;

        const intervalId = setInterval(async () => {
            attempts++;
            log(`OAuth polling attempt #${attempts}, finalized=${finalizedRef.current}`);

            // If we exceed max attempts, give up
            if (attempts > maxAttempts) {
                console.warn('Gave up after multiple attempts. Clearing interval.');
                clearInterval(intervalId);
                setPollingActive(false);
                // Optionally call one last afterFinalize
                await afterFinalize();
                return;
            }

            if (!finalizedRef.current) {
                // Try to finalize
                const resp = await finalizeFn(pollingState);
                if (resp && resp.success) {
                    finalizedRef.current = true;
                    log('Finalization success. Running afterFinalize, then stopping polling...');
                    await afterFinalize();
                    clearInterval(intervalId);
                    setPollingActive(false);
                } else {
                    log('Not ready yet, continuing to poll finalizeFn...');
                }
            } else {
                // Already finalized => do one last afterFinalize
                log('Already finalized, calling afterFinalize, then stopping...');
                await afterFinalize();
                clearInterval(intervalId);
                setPollingActive(false);
            }
        }, intervalMs);

        return () => {
            clearInterval(intervalId);
        };
    }, [pollingActive, pollingState, finalizeFn, afterFinalize, maxAttempts, intervalMs]);

    // Function to start the polling from scratch
    const startPoll = useCallback(
        (oauthState) => {
            log('useOAuthPoll -> startPoll called with oauthState:', oauthState);
            finalizedRef.current = false; // reset
            setPollingState(oauthState);
            setPollingActive(true);
        },
        []
    );

    return { startPoll, pollingActive, pollingState };
}
