import { create } from 'zustand';
import { devtools } from 'zustand/middleware';
import apiService from '../api/apiService';

const useAppDataStore = create(
  devtools(
    (set, get) => ({
      accounts: [],
      associations: [],
      config: null,
      selectedAccountId: null,
      selectedCharacterId: null,
      refreshKey: Date.now(),
      hasInitialFetch: false,
      loading: {
        accounts: false,
        associations: false,
        config: false
      },
      error: {
        accounts: null,
        associations: null,
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
            apiService.getAssociations(),
            apiService.getConfig()
          ]);

          // Handle both paginated and non-paginated responses
          const accounts = Array.isArray(accountsRes) ? accountsRes : (accountsRes?.data || []);
          const associations = associationsRes?.data || [];
          const config = configRes?.data || {};

          set({
            accounts,
            associations,
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

      fetchAccounts: async (bypassCache = false) => {
        console.log('fetchAccounts called - fetching updated account list', bypassCache ? '(bypassing cache)' : '');
        set({ loading: { ...get().loading, accounts: true } });
        try {
          const response = await apiService.getAccounts(bypassCache);
          console.log('Accounts API response:', response);
          // Handle both paginated and non-paginated responses
          // If response has a 'data' field, it's paginated; otherwise it's the raw array
          const accounts = Array.isArray(response) ? response : (response?.data || []);
          console.log('Setting accounts in store:', accounts);
          set({ 
            accounts,
            loading: { ...get().loading, accounts: false },
            error: { ...get().error, accounts: null }
          });
          return accounts;
        } catch (error) {
          console.error('Failed to fetch accounts:', error);
          set({ 
            loading: { ...get().loading, accounts: false },
            error: { ...get().error, accounts: error.message }
          });
          throw error;
        }
      },

      fetchConfig: async () => {
        set({ loading: { ...get().loading, config: true } });
        try {
          const response = await apiService.getConfig();
          const config = response?.data || {};
          set({ 
            config,
            loading: { ...get().loading, config: false },
            error: { ...get().error, config: null }
          });
          return config;
        } catch (error) {
          console.error('Failed to fetch config:', error);
          set({ 
            loading: { ...get().loading, config: false },
            error: { ...get().error, config: error.message }
          });
          throw error;
        }
      },

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

      updateAccount: async (accountId, updates) => {
        try {
          const response = await apiService.updateAccount(accountId, updates);
          if (response?.success) {
            await get().fetchAccounts();
          }
          return response;
        } catch (error) {
          console.error('Failed to update account:', error);
          throw error;
        }
      },

      updateConfig: async (updates) => {
        try {
          const response = await apiService.updateConfig(updates);
          if (response?.success) {
            await get().fetchConfig();
          }
          return response;
        } catch (error) {
          console.error('Failed to update config:', error);
          throw error;
        }
      },

      deleteAccount: async (accountId) => {
        try {
          const response = await apiService.deleteAccount(accountId);
          if (response?.success) {
            await get().fetchAccounts();
          }
          return response;
        } catch (error) {
          console.error('Failed to delete account:', error);
          throw error;
        }
      },

      selectAccount: (accountId) => {
        set({ selectedAccountId: accountId });
      },

      selectCharacter: (characterId) => {
        set({ selectedCharacterId: characterId });
      },

      refreshData: async () => {
        set({ refreshKey: Date.now() });
        const [accounts, associations, config] = await Promise.all([
          get().fetchAccounts(),
          get().fetchAssociations(),
          get().fetchConfig()
        ]);
        return { accounts, associations, config };
      },

      reset: () => {
        set({
          accounts: [],
          associations: [],
          config: null,
          selectedAccountId: null,
          selectedCharacterId: null,
          refreshKey: Date.now(),
          hasInitialFetch: false,
          loading: {
            accounts: false,
            associations: false,
            config: false
          },
          error: {
            accounts: null,
            associations: null,
            config: null
          }
        });
      }
    }),
    {
      name: 'AppDataStore'
    }
  )
);

export default useAppDataStore;