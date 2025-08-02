import { create } from 'zustand';
import { devtools } from 'zustand/middleware';
import apiService from '../api/apiService';

const useEveDataStore = create(
  devtools(
    (set, get) => ({
      skillPlans: {},
      eveProfiles: [],
      eveConversions: {},
      loading: {
        skillPlans: false,
        eveProfiles: false,
        eveConversions: false
      },
      error: {
        skillPlans: null,
        eveProfiles: null,
        eveConversions: null
      },
      hasInitialFetch: false,

      fetchEveData: async (forceRefresh = false) => {
        const state = get();
        
        if (!forceRefresh && state.hasInitialFetch && state.skillPlans) {
          return {
            skillPlans: state.skillPlans,
            eveProfiles: state.eveProfiles,
            eveConversions: state.eveConversions
          };
        }

        set({ 
          loading: { skillPlans: true, eveProfiles: true, eveConversions: true },
          error: { skillPlans: null, eveProfiles: null, eveConversions: null }
        });

        try {
          const [skillPlansRes, eveProfilesRes, eveConversionsRes] = await Promise.all([
            apiService.getEveSkillPlans(),
            apiService.getEveProfiles(),
            apiService.getEveConversions()
          ]);

          const skillPlans = skillPlansRes?.data || {};
          const eveProfiles = eveProfilesRes?.data || [];
          const eveConversions = eveConversionsRes?.data || {};

          set({
            skillPlans,
            eveProfiles,
            eveConversions,
            hasInitialFetch: true,
            loading: { skillPlans: false, eveProfiles: false, eveConversions: false }
          });

          return { skillPlans, eveProfiles, eveConversions };
        } catch (error) {
          console.error('Failed to fetch EVE data:', error);
          set({
            hasInitialFetch: true,
            loading: { skillPlans: false, eveProfiles: false, eveConversions: false },
            error: { 
              skillPlans: error.message,
              eveProfiles: error.message,
              eveConversions: error.message
            }
          });
          throw error;
        }
      },

      fetchSkillPlans: async () => {
        set({ loading: { ...get().loading, skillPlans: true } });
        try {
          const response = await apiService.getEveSkillPlans();
          const skillPlans = response?.data || {};
          set({ 
            skillPlans,
            loading: { ...get().loading, skillPlans: false },
            error: { ...get().error, skillPlans: null }
          });
          return skillPlans;
        } catch (error) {
          console.error('Failed to fetch skill plans:', error);
          set({ 
            loading: { ...get().loading, skillPlans: false },
            error: { ...get().error, skillPlans: error.message }
          });
          throw error;
        }
      },

      fetchEveProfiles: async () => {
        set({ loading: { ...get().loading, eveProfiles: true } });
        try {
          const response = await apiService.getEveProfiles();
          const eveProfiles = response?.data || [];
          set({ 
            eveProfiles,
            loading: { ...get().loading, eveProfiles: false },
            error: { ...get().error, eveProfiles: null }
          });
          return eveProfiles;
        } catch (error) {
          console.error('Failed to fetch EVE profiles:', error);
          set({ 
            loading: { ...get().loading, eveProfiles: false },
            error: { ...get().error, eveProfiles: error.message }
          });
          throw error;
        }
      },

      fetchEveConversions: async () => {
        set({ loading: { ...get().loading, eveConversions: true } });
        try {
          const response = await apiService.getEveConversions();
          const eveConversions = response?.data || {};
          set({ 
            eveConversions,
            loading: { ...get().loading, eveConversions: false },
            error: { ...get().error, eveConversions: null }
          });
          return eveConversions;
        } catch (error) {
          console.error('Failed to fetch EVE conversions:', error);
          set({ 
            loading: { ...get().loading, eveConversions: false },
            error: { ...get().error, eveConversions: error.message }
          });
          throw error;
        }
      },

      reset: () => {
        set({
          skillPlans: {},
          eveProfiles: [],
          eveConversions: {},
          loading: {
            skillPlans: false,
            eveProfiles: false,
            eveConversions: false
          },
          error: {
            skillPlans: null,
            eveProfiles: null,
            eveConversions: null
          },
          hasInitialFetch: false
        });
      }
    }),
    {
      name: 'EveDataStore'
    }
  )
);

export default useEveDataStore;