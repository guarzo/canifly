// src/Routes.jsx

import { Routes, Route } from 'react-router-dom';
import PropTypes from 'prop-types';

import Dashboard from './components/dashboard/Dashboard.jsx';
import SkillPlans from './components/skillplan/SkillPlans.jsx';
import Landing from './components/landing/Landing.jsx';
import Sync from './components/sync/Sync.jsx';
import Mapping from './components/mapping/Mapping.jsx';

function AppRoutes({
                       isAuthenticated,
                       loggedOut,
                       appData,
                       handleToggleAccountStatus,
                       handleUpdateCharacter,
                       handleUpdateAccountName,
                       handleRemoveCharacter,
                       handleRemoveAccount,
                       silentRefreshData,
                       setAppData,
                       identities,
                       backEndURL,
                       logInCallBack
                   }) {
    if (!isAuthenticated || loggedOut) {
        return <Landing backEndURL={backEndURL} logInCallBack={logInCallBack} />;
    } else if (!appData) {
        return (
            <div className="flex items-center justify-center min-h-screen bg-gray-900 text-teal-200">
                <p>Loading...</p>
            </div>
        );
    } else {
        // Adjusting references to new model structure:
        // Accounts & Associations
        const accounts = appData.AccountData?.Accounts || [];
        const associations = appData.AccountData?.Associations || [];

        // ConfigData fields
        const roles = appData.ConfigData?.Roles || [];
        const userSelections = appData.ConfigData?.DropDownSelections || {};
        const currentSettingsDir = appData.ConfigData?.SettingsDir || '';
        const lastBackupDir = appData.ConfigData?.LastBackupDir || '';

        // EveData fields
        const skillPlans = appData.EveData?.SkillPlans || {};
        const eveProfiles = appData.EveData?.EveProfiles || [];

        return (
            <Routes>
                <Route
                    path="/"
                    element={
                        <Dashboard
                            accounts={accounts}
                            onToggleAccountStatus={handleToggleAccountStatus}
                            onUpdateCharacter={handleUpdateCharacter}
                            onUpdateAccountName={handleUpdateAccountName}
                            onRemoveCharacter={handleRemoveCharacter}
                            onRemoveAccount={handleRemoveAccount}
                            roles={roles}
                        />
                    }
                />
                <Route
                    path="/skill-plans"
                    element={
                        <SkillPlans
                            identities={identities}
                            skillPlans={skillPlans}
                            setAppData={setAppData}
                            backEndURL={backEndURL}
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
                            backEndURL={backEndURL}
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
                            backEndURL={backEndURL}
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
        })
    }),
    handleToggleAccountStatus: PropTypes.func.isRequired,
    handleUpdateCharacter: PropTypes.func.isRequired,
    handleUpdateAccountName: PropTypes.func.isRequired,
    handleRemoveCharacter: PropTypes.func.isRequired,
    handleRemoveAccount: PropTypes.func.isRequired,
    silentRefreshData: PropTypes.func.isRequired,
    setAppData: PropTypes.func.isRequired,
    identities: PropTypes.array.isRequired,
    backEndURL: PropTypes.string.isRequired,
    logInCallBack: PropTypes.func.isRequired
};

export default AppRoutes;
