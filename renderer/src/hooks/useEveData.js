import { useEffect } from 'react';
import useEveDataStore from '../stores/eveDataStore';
import { useAuth } from './useAuth';

export function useEveData() {
  const { isAuthenticated } = useAuth();
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
    if (isAuthenticated) {
      fetchEveData();
    } else {
      reset();
    }
  }, [isAuthenticated, fetchEveData, reset]);

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