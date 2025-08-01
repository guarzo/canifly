import { create } from 'zustand';
import { devtools } from 'zustand/middleware';
import apiService from '../api/apiService';

const useAppDataStore = create(
  devtools(
    (set, get) => ({
      accounts: [],
      config: null,
      dashboards: [],
      selectedAccountId: null,
      selectedCharacterId: null,
      refreshKey: Date.now(),
      loading: {
        accounts: false,
        config: false,
        dashboards: false
      },
      error: {
        accounts: null,
        config: null,
        dashboards: null
      },

      fetchAppData: async (forceRefresh = false) => {
        const state = get();
        
        if (!forceRefresh && state.accounts.length > 0 && state.config) {
          return {
            accounts: state.accounts,
            config: state.config,
            dashboards: state.dashboards
          };
        }

        set({ 
          loading: { accounts: true, config: true, dashboards: true },
          error: { accounts: null, config: null, dashboards: null }
        });

        try {
          const [accountsRes, configRes, dashboardsRes] = await Promise.all([
            apiService.getAccounts(),
            apiService.getConfig(),
            apiService.getDashboards()
          ]);

          const accounts = accountsRes?.data || [];
          const config = configRes?.data || {};
          const dashboards = dashboardsRes?.data || [];

          set({
            accounts,
            config,
            dashboards,
            loading: { accounts: false, config: false, dashboards: false }
          });

          return { accounts, config, dashboards };
        } catch (error) {
          console.error('Failed to fetch app data:', error);
          set({
            loading: { accounts: false, config: false, dashboards: false },
            error: { 
              accounts: error.message,
              config: error.message,
              dashboards: error.message
            }
          });
          throw error;
        }
      },

      fetchAccounts: async () => {
        set({ loading: { ...get().loading, accounts: true } });
        try {
          const response = await apiService.getAccounts();
          const accounts = response?.data || [];
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
        return get().fetchAppData(true);
      },

      reset: () => {
        set({
          accounts: [],
          config: null,
          dashboards: [],
          selectedAccountId: null,
          selectedCharacterId: null,
          refreshKey: Date.now(),
          loading: {
            accounts: false,
            config: false,
            dashboards: false
          },
          error: {
            accounts: null,
            config: null,
            dashboards: null
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