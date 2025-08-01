// src/Routes.jsx

import { Routes, Route } from 'react-router-dom';
import { useAuth } from './hooks/useAuth';
import { useAppData } from './hooks/useAppData';

import CharacterOverview from './pages/CharacterOverview.jsx';
import SkillPlans from './pages/SkillPlans.jsx';
import Landing from './pages/Landing.jsx';
import Sync from './pages/Sync.jsx';
import Mapping from './pages/Mapping.jsx';

function AppRoutes({ characters }) {
    const { isAuthenticated } = useAuth();
    const { accounts, config, dashboards, isLoading } = useAppData();
    
    if (!isAuthenticated) {
        return <Landing />;
    }
    
    if (isLoading) {
        return (
            <div className="flex items-center justify-center min-h-screen bg-gray-900 text-teal-200">
                <p>Loading...</p>
            </div>
        );
    }
    
    // Extract data from stores
    const roles = config?.Roles || [];
    const skillPlans = config?.SkillPlans || {};
    const eveProfiles = config?.EveProfiles || [];
    const associations = accounts.flatMap(acc => acc.associations || []);
    const userSelections = config?.DropDownSelections || {};
    const currentSettingsDir = config?.SettingsDir || '';
    const lastBackupDir = config?.LastBackupDir || '';
    const eveConversions = config?.EveConversions || {};

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
            <Route path="*" element={<div>Route Not Found</div>} />
        </Routes>
    );
}

export default AppRoutes;
