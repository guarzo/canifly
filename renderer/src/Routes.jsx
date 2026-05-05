// src/Routes.jsx
import React from 'react';
import { Routes, Route, useLocation, Navigate } from 'react-router-dom';
import { AnimatePresence } from 'framer-motion';
import { useAuth } from './hooks/useAuth';
import { useAppData } from './hooks/useAppData';
import { useEveData } from './hooks/useEveData';
import { logger } from './utils/logger';

import CharacterOverview from './pages/CharacterOverview.jsx';
import SkillPlans from './pages/SkillPlans.jsx';
import Landing from './pages/Landing.jsx';
import Profiles from './pages/Profiles.jsx';
import Settings from './pages/Settings.jsx';
import PageTransition from './components/transitions/PageTransition.jsx';
import LoadingScreen from './components/ui/LoadingScreen.jsx';

function AppRoutes({ characters }) {
    const location = useLocation();
    const { isAuthenticated } = useAuth();
    const { associations, config, isLoading } = useAppData();
    const { skillPlans, eveProfiles, eveConversions, loading: eveLoading } = useEveData();

    React.useEffect(() => {
        logger.debug('Routes - associations:', associations);
        logger.debug('Routes - eveProfiles:', eveProfiles);
        logger.debug('Routes - eveConversions:', eveConversions);
    }, [associations, eveProfiles, eveConversions]);

    if (!isAuthenticated) {
        return <Landing />;
    }
    if (isLoading || eveLoading.skillPlans || eveLoading.eveProfiles || eveLoading.eveConversions) {
        return <LoadingScreen message="Loading…" />;
    }

    const roles = config?.roles || config?.Roles || [];
    const userSelections = config?.userSelections || config?.DropDownSelections || {};
    const currentSettingsDir = config?.settingsDir || config?.SettingsDir || '';
    const lastBackupDir = config?.lastBackupDir || config?.LastBackupDir || [];

    return (
        <AnimatePresence mode="wait">
            <Routes location={location} key={location.pathname}>
                <Route
                    path="/"
                    element={
                        <PageTransition>
                            <CharacterOverview roles={roles} skillConversions={eveConversions} />
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
                    path="/profiles"
                    element={
                        <PageTransition>
                            <Profiles
                                subDirs={eveProfiles}
                                associations={associations}
                                settingsData={eveProfiles}
                                userSelections={userSelections}
                                currentSettingsDir={currentSettingsDir}
                                lastBackupDir={lastBackupDir}
                            />
                        </PageTransition>
                    }
                />
                {/* Bookmark-preserving redirects */}
                <Route path="/mapping" element={<Navigate to="/profiles?view=mapping" replace />} />
                <Route path="/sync" element={<Navigate to="/profiles?view=sync" replace />} />
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
