import { create } from 'zustand';
import { devtools, persist } from 'zustand/middleware';
import apiService from '../api/apiService';
import { clearSessionCookie } from '../utils/clearCookies';

// Track if auth check is in progress globally
let authCheckInProgress = false;

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
          // Prevent multiple concurrent auth checks
          if (authCheckInProgress || get().authCheckComplete) {
            return get().isAuthenticated;
          }
          
          authCheckInProgress = true;
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
            
            // If we get a specific error about authentication, clear the invalid cookie
            if (error.message && error.message.includes('user is not logged in')) {
              console.log('Clearing invalid session cookie');
              clearSessionCookie();
            }
            
            set({ 
              isAuthenticated: false,
              authCheckComplete: true,
              user: null,
              loading: false,
              error: error.message 
            });
            return false;
          } finally {
            authCheckInProgress = false;
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
          authCheckInProgress = false;
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