// src/hooks/useFuzzworksStatus.js
//
// Subscribes to the application WebSocket hub and exposes the current
// Fuzzworks update state — `idle` until the backend emits a
// `fuzzworks:status` message, then `updating` / `ready` / `error` based
// on what the backend reports. If Plan A's backend isn't merged yet, no
// messages arrive and the chip stays in `idle`.
import { useState, useCallback } from 'react';
import { useWebSocket } from './useWebSocket';

const TERMINAL_STATES = new Set(['ready', 'error']);

export function useFuzzworksStatus() {
    const [status, setStatus] = useState({ state: 'idle' });

    const onMessage = useCallback((message) => {
        if (!message || message.type !== 'fuzzworks:status') return;
        const { state, error } = message;
        if (state === 'updating') {
            setStatus({ state: 'updating' });
        } else if (TERMINAL_STATES.has(state)) {
            setStatus({ state, error });
        }
    }, []);

    useWebSocket(onMessage);

    return status;
}
