import { useEffect } from 'react';
import useEveDataStore from '../stores/eveDataStore';
import { useAuth } from './useAuth';

export function useEveData() {
  const { isLoggedIn } = useAuth();
  const {
    skillPlans,
    eveProfiles,
    eveConversions,
    loading,
    error,
    fetchEveData,
    fetchSkillPlans,
    fetchEveProfiles,
    fetchEveConversions,
    reset
  } = useEveDataStore();

  useEffect(() => {
    if (isLoggedIn) {
      fetchEveData();
    } else {
      reset();
    }
  }, [isLoggedIn]);

  return {
    skillPlans,
    eveProfiles,
    eveConversions,
    loading,
    error,
    fetchEveData,
    fetchSkillPlans,
    fetchEveProfiles,
    fetchEveConversions,
    refresh: () => fetchEveData(true)
  };
}