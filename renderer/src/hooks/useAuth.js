import useAuthStore from '../stores/authStore';
import { useAsyncOperation } from './useAsyncOperation';

export function useAuth() {
  const {
    isAuthenticated,
    authCheckComplete,
    user,
    loading,
    error,
    checkAuth,
    login,
    logout,
    reset
  } = useAuthStore();

  const { execute: executeLogin } = useAsyncOperation();
  const { execute: executeLogout } = useAsyncOperation();
  const { execute: executeCheckAuth } = useAsyncOperation();

  // Don't trigger auth check here - let App.jsx handle it once

  const handleLogin = async (account) => {
    return executeLogin(() => login(account), {
      successMessage: 'Login initiated',
      errorMessage: 'Login failed'
    });
  };

  const handleLogout = async () => {
    return executeLogout(logout, {
      successMessage: 'Successfully logged out',
      errorMessage: 'Logout failed'
    });
  };

  const refreshAuth = async () => {
    return executeCheckAuth(() => checkAuth(true), {  // Force refresh
      showToast: false
    });
  };

  return {
    isAuthenticated,
    authCheckComplete,
    user,
    loading,
    error,
    login: handleLogin,
    logout: handleLogout,
    refreshAuth,
    reset
  };
}