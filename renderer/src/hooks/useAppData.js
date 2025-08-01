import { useEffect } from 'react';
import useAppDataStore from '../stores/appDataStore';
import { useAuth } from './useAuth';
import { useAsyncOperation } from './useAsyncOperation';

export function useAppData() {
  const {
    accounts,
    config,
    dashboards,
    selectedAccountId,
    selectedCharacterId,
    refreshKey,
    loading,
    error,
    fetchAppData,
    fetchAccounts,
    fetchConfig,
    updateAccount,
    updateConfig,
    deleteAccount,
    selectAccount,
    selectCharacter,
    refreshData,
    reset
  } = useAppDataStore();

  const { isAuthenticated } = useAuth();
  const { execute } = useAsyncOperation();

  useEffect(() => {
    if (isAuthenticated && accounts.length === 0) {
      execute(() => fetchAppData(false), { showToast: false });
    }
  }, [isAuthenticated, accounts.length, fetchAppData, execute]);

  useEffect(() => {
    if (!isAuthenticated) {
      reset();
    }
  }, [isAuthenticated, reset]);

  const handleUpdateAccount = async (accountId, updates) => {
    return execute(
      () => updateAccount(accountId, updates),
      {
        successMessage: 'Account updated successfully',
        errorMessage: 'Failed to update account'
      }
    );
  };

  const handleUpdateConfig = async (updates) => {
    return execute(
      () => updateConfig(updates),
      {
        successMessage: 'Configuration updated successfully',
        errorMessage: 'Failed to update configuration'
      }
    );
  };

  const handleDeleteAccount = async (accountId) => {
    return execute(
      () => deleteAccount(accountId),
      {
        successMessage: 'Account deleted successfully',
        errorMessage: 'Failed to delete account'
      }
    );
  };

  const handleRefreshData = async () => {
    return execute(
      () => refreshData(),
      {
        successMessage: 'Data refreshed successfully',
        errorMessage: 'Failed to refresh data'
      }
    );
  };

  const getSelectedAccount = () => {
    return accounts.find(acc => acc.id === selectedAccountId);
  };

  const getAccountCharacters = (accountId) => {
    const account = accounts.find(acc => acc.id === accountId);
    return account?.characters || [];
  };

  const getSelectedCharacter = () => {
    for (const account of accounts) {
      const character = account.characters?.find(char => char.id === selectedCharacterId);
      if (character) return character;
    }
    return null;
  };

  return {
    accounts,
    config,
    dashboards,
    selectedAccountId,
    selectedCharacterId,
    refreshKey,
    loading,
    error,
    isLoading: loading.accounts || loading.config || loading.dashboards,
    fetchAppData: (forceRefresh) => execute(() => fetchAppData(forceRefresh), { showToast: false }),
    fetchAccounts: () => execute(fetchAccounts, { showToast: false }),
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