// src/Routes.jsx

import { Routes, Route } from 'react-router-dom';
import { useAuth } from './hooks/useAuth';
import { useAppData } from './hooks/useAppData';
import { useEveData } from './hooks/useEveData';

import CharacterOverview from './pages/CharacterOverview.jsx';
import SkillPlans from './pages/SkillPlans.jsx';
import Landing from './pages/Landing.jsx';
import Sync from './pages/Sync.jsx';
import Mapping from './pages/Mapping.jsx';
import Settings from './pages/Settings.jsx';

function AppRoutes({ characters }) {
    const { isAuthenticated } = useAuth();
    const { accounts, associations, config, isLoading } = useAppData();
    const { skillPlans, eveProfiles, eveConversions, loading: eveLoading } = useEveData();
    
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
        <Routes>
            <Route
                path="/"
                element={
                    <CharacterOverview
                        roles={roles}
                        skillConversions={eveConversions}
                    />
                }
            />
            <Route
                path="/skill-plans"
                element={
                    <SkillPlans
                        characters={characters}
                        skillPlans={skillPlans}
                        conversions={eveConversions}
                    />
                }
            />
            <Route
                path="/sync"
                element={
                    <Sync
                        settingsData={eveProfiles}
                        associations={associations}
                        currentSettingsDir={currentSettingsDir}
                        userSelections={userSelections}
                        lastBackupDir={lastBackupDir}
                    />
                }
            />
            <Route
                path="/mapping"
                element={
                    <Mapping
                        associations={associations}
                        subDirs={eveProfiles}
                    />
                }
            />
            <Route
                path="/settings"
                element={<Settings />}
            />
            <Route path="*" element={<div>Route Not Found</div>} />
        </Routes>
    );
}

export default AppRoutes;
