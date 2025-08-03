import { useEffect, useRef, useCallback, useState } from 'react';
import { useAuth } from './useAuth';
import { logger } from '../utils/logger';

// Use relative URL for WebSocket to work with proxy
const WS_URL = window.location.protocol === 'https:' 
  ? `wss://${window.location.host}/api/ws`
  : `ws://${window.location.host}/api/ws`;
const RECONNECT_DELAY = 5000; // 5 seconds
const PING_INTERVAL = 30000; // 30 seconds

export function useWebSocket(onMessage) {
  const { isAuthenticated } = useAuth();
  const wsRef = useRef(null);
  const reconnectTimeoutRef = useRef(null);
  const pingIntervalRef = useRef(null);
  const [connectionState, setConnectionState] = useState('disconnected');

  const connect = useCallback(() => {
    if (!isAuthenticated || connectionState === 'connecting' || wsRef.current?.readyState === WebSocket.OPEN) {
      return;
    }

    try {
      setConnectionState('connecting');
      wsRef.current = new WebSocket(WS_URL);

      wsRef.current.onopen = () => {
        logger.debug('WebSocket connected');
        setConnectionState('connected');
        
        // Start ping interval
        pingIntervalRef.current = setInterval(() => {
          if (wsRef.current?.readyState === WebSocket.OPEN) {
            wsRef.current.send(JSON.stringify({ type: 'ping' }));
          }
        }, PING_INTERVAL);
      };

      wsRef.current.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data);
          if (data.type !== 'pong' && onMessage) {
            onMessage(data);
          }
        } catch (error) {
          logger.error('Failed to parse WebSocket message:', error);
        }
      };

      wsRef.current.onerror = (error) => {
        logger.error('WebSocket error:', error);
        setConnectionState('error');
      };

      wsRef.current.onclose = () => {
        logger.debug('WebSocket disconnected');
        setConnectionState('disconnected');
        
        // Clear ping interval
        if (pingIntervalRef.current) {
          clearInterval(pingIntervalRef.current);
          pingIntervalRef.current = null;
        }

        // Attempt to reconnect after delay if authenticated
        if (isAuthenticated) {
          reconnectTimeoutRef.current = setTimeout(() => {
            connect();
          }, RECONNECT_DELAY);
        }
      };
    } catch (error) {
      logger.error('Failed to create WebSocket connection:', error);
      setConnectionState('error');
    }
  }, [isAuthenticated, onMessage, connectionState]);

  const disconnect = useCallback(() => {
    if (reconnectTimeoutRef.current) {
      clearTimeout(reconnectTimeoutRef.current);
      reconnectTimeoutRef.current = null;
    }

    if (pingIntervalRef.current) {
      clearInterval(pingIntervalRef.current);
      pingIntervalRef.current = null;
    }

    if (wsRef.current && wsRef.current.readyState !== WebSocket.CLOSED) {
      wsRef.current.close();
      wsRef.current = null;
    }
    
    setConnectionState('disconnected');
  }, []);

  useEffect(() => {
    // Add a small delay to ensure backend is ready
    const timer = setTimeout(() => {
      if (isAuthenticated) {
        connect();
      } else {
        disconnect();
      }
    }, 1000); // 1 second delay

    return () => {
      clearTimeout(timer);
      disconnect();
    };
  }, [isAuthenticated]); // Remove connect and disconnect from deps to prevent loops

  const sendMessage = useCallback((message) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify(message));
    } else {
      logger.warn('WebSocket is not connected');
    }
  }, []);

  return {
    sendMessage,
    isConnected: wsRef.current?.readyState === WebSocket.OPEN,
    connectionState
  };
}