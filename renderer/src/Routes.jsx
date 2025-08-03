// src/Routes.jsx

import React from 'react';
import { Routes, Route, useLocation } from 'react-router-dom';
import { AnimatePresence } from 'framer-motion';
import { useAuth } from './hooks/useAuth';
import { useAppData } from './hooks/useAppData';
import { useEveData } from './hooks/useEveData';
import { logger } from './utils/logger';

import CharacterOverview from './pages/CharacterOverview.jsx';
import SkillPlans from './pages/SkillPlans.jsx';
import Landing from './pages/Landing.jsx';
import Sync from './pages/Sync.jsx';
import Mapping from './pages/Mapping.jsx';
import Settings from './pages/Settings.jsx';
import PageTransition from './components/transitions/PageTransition.jsx';

function AppRoutes({ characters }) {
    const location = useLocation();
    const { isAuthenticated } = useAuth();
    const { accounts, associations, config, isLoading } = useAppData();
    const { skillPlans, eveProfiles, eveConversions, loading: eveLoading } = useEveData();
    
    // Debug logging
    React.useEffect(() => {
        logger.debug('Routes - associations:', associations);
        logger.debug('Routes - eveProfiles:', eveProfiles);
    }, [associations, eveProfiles]);
    
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
    const userSelections = config?.DropDownSelections || {};
    const currentSettingsDir = config?.SettingsDir || '';
    const lastBackupDir = config?.LastBackupDir || [];

    return (
        <AnimatePresence mode="wait">
            <Routes location={location} key={location.pathname}>
                <Route
                    path="/"
                    element={
                        <PageTransition>
                            <CharacterOverview
                                roles={roles}
                                skillConversions={eveConversions}
                            />
                        </PageTransition>
                    }
                />
                <Route
                    path="/skill-plans"
                    element={
                        <PageTransition>
                            <SkillPlans
                                characters={characters}
                                skillPlans={skillPlans}
                                conversions={eveConversions}
                            />
                        </PageTransition>
                    }
                />
                <Route
                    path="/sync"
                    element={
                        <PageTransition>
                            <Sync
                                settingsData={eveProfiles}
                                associations={associations}
                                currentSettingsDir={currentSettingsDir}
                                userSelections={userSelections}
                                lastBackupDir={lastBackupDir}
                            />
                        </PageTransition>
                    }
                />
                <Route
                    path="/mapping"
                    element={
                        <PageTransition>
                            <Mapping
                                associations={associations}
                                subDirs={eveProfiles}
                            />
                        </PageTransition>
                    }
                />
                <Route
                    path="/settings"
                    element={
                        <PageTransition>
                            <Settings />
                        </PageTransition>
                    }
                />
                <Route path="*" element={<PageTransition><div>Route Not Found</div></PageTransition>} />
            </Routes>
        </AnimatePresence>
    );
}

export default AppRoutes;
