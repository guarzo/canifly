// src/Routes.jsx

import { Routes, Route } from 'react-router-dom';
import PropTypes from 'prop-types';

import CharacterOverview from './pages/CharacterOverview.jsx';
import SkillPlans from './pages/SkillPlans.jsx';
import Landing from './pages/Landing.jsx';
import Sync from './pages/Sync.jsx';
import Mapping from './pages/Mapping.jsx';

function AppRoutes({
                       isAuthenticated,
                       loggedOut,
                       appData,
                       handleToggleAccountStatus,
                       handleUpdateCharacter,
                       handleUpdateAccountName,
                       handleRemoveCharacter,
                       handleRemoveAccount,
                       handleDeleteSkillPlan,
                       handleCopySkillPlan,
                       silentRefreshData,
                       setAppData,
                       characters,
                       logInCallBack,
                       handleToggleAccountVisibility,
                   }) {
    if (!isAuthenticated || loggedOut) {
        return <Landing logInCallBack={logInCallBack} />;
    } else if (!appData) {
        return (
            <div className="flex items-center justify-center min-h-screen bg-gray-900 text-teal-200">
                <p>Loading...</p>
            </div>
        );
    } else {
        // Using new model fields:
        const accounts = appData.AccountData?.Accounts || [];
        const roles = appData.ConfigData?.Roles || [];
        const skillPlans = appData.EveData?.SkillPlans || {};
        const eveProfiles = appData.EveData?.EveProfiles || [];
        const associations = appData.AccountData?.Associations || [];
        const userSelections = appData.ConfigData?.DropDownSelections || {};
        const currentSettingsDir = appData.ConfigData?.SettingsDir || '';
        const lastBackupDir = appData.ConfigData?.LastBackupDir || '';
        const eveConversions = appData.EveData?.EveConversions || {};

        return (
            <Routes>
                <Route
                    path="/"
                    element={
                        <CharacterOverview
                            accounts={accounts}
                            onToggleAccountStatus={handleToggleAccountStatus}
                            onUpdateCharacter={handleUpdateCharacter}
                            onUpdateAccountName={handleUpdateAccountName}
                            onRemoveCharacter={handleRemoveCharacter}
                            onRemoveAccount={handleRemoveAccount}
                            roles={roles}
                            skillConversions={eveConversions}
                            onToggleAccountVisibility={handleToggleAccountVisibility}
                        />
                    }
                />
                <Route
                    path="/skill-plans"
                    element={
                        <SkillPlans
                            characters={characters}
                            skillPlans={skillPlans} // Using EveData.SkillPlans
                            setAppData={setAppData}
                            conversions={eveConversions}
                            onDeleteSkillPlan={handleDeleteSkillPlan}
                            onCopySkillPlan={handleCopySkillPlan}
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
                            onRefreshData={silentRefreshData}
                        />
                    }
                />
                <Route path="*" element={<div>Route Not Found</div>} />
            </Routes>
        );
    }
}

AppRoutes.propTypes = {
    isAuthenticated: PropTypes.bool.isRequired,
    loggedOut: PropTypes.bool.isRequired,
    appData: PropTypes.shape({
        AccountData: PropTypes.shape({
            Accounts: PropTypes.array,
            Associations: PropTypes.array
        }),
        ConfigData: PropTypes.shape({
            Roles: PropTypes.array,
            SettingsDir: PropTypes.string,
            LastBackupDir: PropTypes.string,
            DropDownSelections: PropTypes.object,
        }),
        EveData: PropTypes.shape({
            SkillPlans: PropTypes.object,
            EveProfiles: PropTypes.array,
            EveConversions: PropTypes.object,
        })
    }),
    handleToggleAccountStatus: PropTypes.func.isRequired,
    handleUpdateCharacter: PropTypes.func.isRequired,
    handleUpdateAccountName: PropTypes.func.isRequired,
    handleRemoveCharacter: PropTypes.func.isRequired,
    handleRemoveAccount: PropTypes.func.isRequired,
    silentRefreshData: PropTypes.func.isRequired,
    setAppData: PropTypes.func.isRequired,
    characters: PropTypes.array.isRequired,
    logInCallBack: PropTypes.func.isRequired,
    handleToggleAccountVisibility: PropTypes.func.isRequired,
};

export default AppRoutes;
