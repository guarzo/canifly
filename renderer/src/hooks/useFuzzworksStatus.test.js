// src/hooks/useFuzzworksStatus.test.js
import { vi, describe, test, expect, beforeEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';

let registeredOnMessage;
vi.mock('./useWebSocket', () => ({
    useWebSocket: (onMessage) => {
        registeredOnMessage = onMessage;
        return { sendMessage: vi.fn(), isConnected: false, connectionState: 'disconnected' };
    },
}));

import { useFuzzworksStatus } from './useFuzzworksStatus';

describe('useFuzzworksStatus', () => {
    beforeEach(() => { registeredOnMessage = null; });

    test('starts in idle state when no message has been received', () => {
        const { result } = renderHook(() => useFuzzworksStatus());
        expect(result.current).toEqual({ state: 'idle' });
    });

    test('ignores unrelated WebSocket message types', () => {
        const { result } = renderHook(() => useFuzzworksStatus());
        act(() => registeredOnMessage({ type: 'something:else', state: 'updating' }));
        expect(result.current).toEqual({ state: 'idle' });
    });

    test('transitions to updating then ready', () => {
        const { result } = renderHook(() => useFuzzworksStatus());
        act(() => registeredOnMessage({ type: 'fuzzworks:status', state: 'updating' }));
        expect(result.current).toEqual({ state: 'updating' });
        act(() => registeredOnMessage({ type: 'fuzzworks:status', state: 'ready' }));
        expect(result.current.state).toBe('ready');
    });

    test('captures the error message in the error state', () => {
        const { result } = renderHook(() => useFuzzworksStatus());
        act(() => registeredOnMessage({ type: 'fuzzworks:status', state: 'error', error: 'boom' }));
        expect(result.current).toEqual({ state: 'error', error: 'boom' });
    });
});
