import { create } from 'zustand';
import { devtools, persist } from 'zustand/middleware';
import apiService from '../api/apiService';

const useAuthStore = create(
  devtools(
    persist(
      (set, get) => ({
        isAuthenticated: false,
        authCheckComplete: false,
        user: null,
        loading: false,
        error: null,

        checkAuth: async () => {
          set({ loading: true, error: null });
          try {
            const response = await apiService.getSession();
            const isAuthenticated = response?.status === 'ok' && response?.authenticated === true;
            set({ 
              isAuthenticated,
              authCheckComplete: true,
              user: isAuthenticated ? response.user : null,
              loading: false 
            });
            return isAuthenticated;
          } catch (error) {
            console.error('Auth check failed:', error);
            set({ 
              isAuthenticated: false,
              authCheckComplete: true,
              user: null,
              loading: false,
              error: error.message 
            });
            return false;
          }
        },

        login: async (account) => {
          set({ loading: true, error: null });
          try {
            const response = await apiService.initiateLogin(account);
            if (response?.redirectURL) {
              // Redirect to EVE SSO login
              window.location.href = response.redirectURL;
              return true;
            }
            return false;
          } catch (error) {
            console.error('Login failed:', error);
            set({ loading: false, error: error.message });
            return false;
          }
        },

        logout: async () => {
          set({ loading: true, error: null });
          try {
            await apiService.logout();
            set({ 
              isAuthenticated: false,
              user: null,
              loading: false,
              authCheckComplete: true 
            });
            return true;
          } catch (error) {
            console.error('Logout failed:', error);
            set({ loading: false, error: error.message });
            return false;
          }
        },

        reset: () => {
          set({
            isAuthenticated: false,
            authCheckComplete: false,
            user: null,
            loading: false,
            error: null
          });
        }
      }),
      {
        name: 'auth-storage',
        partialize: (state) => ({ 
          isAuthenticated: state.isAuthenticated,
          authCheckComplete: state.authCheckComplete,
          user: state.user 
        })
      }
    ),
    {
      name: 'AuthStore'
    }
  )
);

export default useAuthStore;