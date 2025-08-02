# Validation Issues Fix

## Issues Found

1. **WebSocket Connection Error**: The WebSocket tries to close before connection is established (timing issue)
2. **Missing Associations Data**: The associations are not being fetched separately for the Mapping page

## Solutions

### 1. Fix WebSocket Connection Timing

**File: `renderer/src/hooks/useWebSocket.js`**

```javascript
// Add connection state tracking to prevent premature close
import { useEffect, useRef, useCallback, useState } from 'react';
import { useAuth } from './useAuth';

const WS_URL = 'ws://localhost:42423/api/ws';
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
        console.log('WebSocket connected');
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
          console.error('Failed to parse WebSocket message:', error);
        }
      };

      wsRef.current.onerror = (error) => {
        console.error('WebSocket error:', error);
        setConnectionState('error');
      };

      wsRef.current.onclose = () => {
        console.log('WebSocket disconnected');
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
      console.error('Failed to create WebSocket connection:', error);
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
    if (isAuthenticated) {
      connect();
    } else {
      disconnect();
    }

    return () => {
      disconnect();
    };
  }, [isAuthenticated]); // Remove connect and disconnect from deps to prevent loops

  const sendMessage = useCallback((message) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify(message));
    } else {
      console.warn('WebSocket is not connected');
    }
  }, []);

  return {
    sendMessage,
    isConnected: wsRef.current?.readyState === WebSocket.OPEN,
    connectionState
  };
}
```

### 2. Add Associations to App Data Store

**File: `renderer/src/stores/appDataStore.js`**

Add associations to the store:

```javascript
const useAppDataStore = create(
  devtools(
    (set, get) => ({
      accounts: [],
      associations: [], // Add this
      config: null,
      selectedAccountId: null,
      selectedCharacterId: null,
      refreshKey: Date.now(),
      hasInitialFetch: false,
      loading: {
        accounts: false,
        associations: false, // Add this
        config: false
      },
      error: {
        accounts: null,
        associations: null, // Add this
        config: null
      },

      fetchAppData: async (forceRefresh = false) => {
        const state = get();
        
        if (!forceRefresh && state.hasInitialFetch && state.config) {
          return {
            accounts: state.accounts,
            associations: state.associations,
            config: state.config
          };
        }

        set({ 
          loading: { accounts: true, associations: true, config: true },
          error: { accounts: null, associations: null, config: null }
        });

        try {
          const [accountsRes, associationsRes, configRes] = await Promise.all([
            apiService.getAccounts(),
            apiService.getAssociations(), // Add this
            apiService.getConfig()
          ]);

          const accounts = accountsRes?.data || [];
          const associations = associationsRes?.data || []; // Add this
          const config = configRes?.data || {};

          set({
            accounts,
            associations, // Add this
            config,
            hasInitialFetch: true,
            loading: { accounts: false, associations: false, config: false }
          });

          return { accounts, associations, config };
        } catch (error) {
          console.error('Failed to fetch app data:', error);
          set({
            hasInitialFetch: true,
            loading: { accounts: false, associations: false, config: false },
            error: { 
              accounts: error.message,
              associations: error.message,
              config: error.message
            }
          });
          throw error;
        }
      },

      // Add fetchAssociations method
      fetchAssociations: async () => {
        set({ loading: { ...get().loading, associations: true } });
        try {
          const response = await apiService.getAssociations();
          const associations = response?.data || [];
          set({ 
            associations,
            loading: { ...get().loading, associations: false },
            error: { ...get().error, associations: null }
          });
          return associations;
        } catch (error) {
          console.error('Failed to fetch associations:', error);
          set({ 
            loading: { ...get().loading, associations: false },
            error: { ...get().error, associations: error.message }
          });
          throw error;
        }
      },

      // Update refreshData to include associations
      refreshData: async () => {
        set({ refreshKey: Date.now() });
        const [accounts, associations, config] = await Promise.all([
          get().fetchAccounts(),
          get().fetchAssociations(), // Add this
          get().fetchConfig()
        ]);
        return { accounts, associations, config };
      },

      // ... rest of the store
    }),
    {
      name: 'app-data-store'
    }
  )
);
```

### 3. Update useAppData Hook

**File: `renderer/src/hooks/useAppData.js`**

Add associations to the return value:

```javascript
export function useAppData() {
  const {
    accounts,
    associations, // Add this
    config,
    selectedAccountId,
    selectedCharacterId,
    refreshKey,
    hasInitialFetch,
    loading,
    error,
    fetchAppData,
    fetchAccounts,
    fetchAssociations, // Add this
    fetchConfig,
    updateAccount,
    updateConfig,
    deleteAccount,
    selectAccount,
    selectCharacter,
    refreshData,
    reset
  } = useAppDataStore();

  // ... existing code ...

  return {
    accounts,
    associations, // Add this
    config,
    selectedAccountId,
    selectedCharacterId,
    refreshKey,
    loading,
    error,
    isLoading: loading.accounts || loading.associations || loading.config, // Update this
    fetchAppData: (forceRefresh) => execute(() => fetchAppData(forceRefresh), { showToast: false }),
    fetchAccounts: () => execute(fetchAccounts, { showToast: false }),
    fetchAssociations: () => execute(fetchAssociations, { showToast: false }), // Add this
    fetchConfig: () => execute(fetchConfig, { showToast: false }),
    updateAccount: handleUpdateAccount,
    updateConfig: handleUpdateConfig,
    deleteAccount: handleDeleteAccount,
    selectAccount,
    selectCharacter,
    refreshData: handleRefreshData,
    getSelectedAccount,
    getAccountCharacters,
    getSelectedCharacter,
    reset
  };
}
```

### 4. Update Routes.jsx

**File: `renderer/src/Routes.jsx`**

Update to use associations from the store:

```javascript
function AppRoutes({ characters }) {
    const { isAuthenticated } = useAuth();
    const { accounts, associations, config, isLoading } = useAppData(); // Add associations
    const { skillPlans, eveProfiles, eveConversions, loading: eveLoading } = useEveData();
    
    if (!isAuthenticated) {
        return <Landing />;
    }
    
    if (isLoading || eveLoading.skillPlans || eveLoading.eveProfiles || eveLoading.eveConversions) {
        return (
            <div className="flex items-center justify-center min-h-screen bg-gray-900 text-teal-200">
                <p>Loading...</p>
            </div>
        );
    }
    
    // Extract data from stores
    const roles = config?.Roles || [];
    // Remove this line: const associations = accounts.flatMap(acc => acc.associations || []);
    const userSelections = config?.DropDownSelections || {};
    const currentSettingsDir = config?.SettingsDir || '';
    const lastBackupDir = config?.LastBackupDir || [];

    return (
        <Routes>
            {/* ... existing routes ... */}
            <Route
                path="/mapping"
                element={
                    <Mapping
                        associations={associations} // Now from store
                        subDirs={eveProfiles}
                    />
                }
            />
            {/* ... rest of routes ... */}
        </Routes>
    );
}
```

## Testing

1. **WebSocket Testing**:
   - Check console for WebSocket connection logs
   - Verify no premature close errors
   - Confirm reconnection works on disconnect

2. **Associations Testing**:
   - Navigate to Mapping page
   - Verify associations are loaded
   - Check that character-account mappings display correctly

## Next Steps

After implementing these fixes:
1. Test all functionality
2. Verify data persistence
3. Check that all CRUD operations work correctly