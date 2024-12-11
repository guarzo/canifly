// src/Routes.jsx

import { Routes, Route } from 'react-router-dom';
import PropTypes from 'prop-types';

import Dashboard from './components/dashboard/Dashboard.jsx';
import SkillPlans from './components/skillplan/SkillPlans.jsx';
import Landing from './components/landing/Landing.jsx';
import CharacterSort from './components/dashboard/CharacterSort.jsx';
import Sync from './components/sync/Sync.jsx';
import Mapping from './components/mapping/Mapping.jsx';

/**
 * Extracted route definitions for cleaner App.jsx.
 *
 * @param {Object} props
 * @param {boolean} props.isAuthenticated
 * @param {boolean} props.loggedOut
 * @param {Object|null} props.appData
 * @param {Function} props.handleToggleAccountStatus
 * @param {Function} props.handleUpdateCharacter
 * @param {Function} props.handleUpdateAccountName
 * @param {Function} props.handleRemoveCharacter
 * @param {Function} props.handleRemoveAccount
 * @param {Function} props.handleAddCharacter
 * @param {Function} props.handleSaveSkillPlan
 * @param {Function} props.silentRefreshData
 * @param {Function} props.setAppData
 * @param {Array} props.identities
 * @param {Array} props.existingAccounts
 * @param {string} props.backEndURL
 * @param {Function} props.logInCallBack
 */
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
        return (
            <Routes>
                <Route
                    path="/"
                    element={
                        <Dashboard
                            accounts={appData.Accounts}
                            onToggleAccountStatus={handleToggleAccountStatus}
                            onUpdateCharacter={handleUpdateCharacter}
                            onUpdateAccountName={handleUpdateAccountName}
                            onRemoveCharacter={handleRemoveCharacter}
                            onRemoveAccount={handleRemoveAccount}
                            roles={appData.Roles}
                        />
                    }
                />
                <Route
                    path="/skill-plans"
                    element={
                        <SkillPlans
                            identities={identities}
                            skillPlans={appData.SkillPlans}
                            setAppData={setAppData}
                            backEndURL={backEndURL}
                        />
                    }
                />
                <Route
                    path="/character-sort"
                    element={
                        <CharacterSort
                            accounts={appData.Accounts}
                            roles={appData.Roles}
                            onUpdateCharacter={handleUpdateCharacter}
                        />
                    }
                />
                <Route
                    path="/sync"
                    element={
                        <Sync
                            settingsData={appData.SubDirs}
                            associations={appData.associations}
                            currentSettingsDir={appData.SettingsDir}
                            userSelections={appData.UserSelections}
                            lastBackupDir={appData.LastBackupDir}
                            backEndURL={backEndURL}
                        />
                    }
                />
                <Route
                    path="/mapping"
                    element={
                        <Mapping
                            associations={appData.associations}
                            subDirs={appData.SubDirs}
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
        Accounts: PropTypes.arrayOf(
            PropTypes.shape({
                Name: PropTypes.string,
                ID: PropTypes.oneOfType([PropTypes.number, PropTypes.string]),
                Characters: PropTypes.arrayOf(
                    PropTypes.shape({
                        Character: PropTypes.shape({
                            CharacterID: PropTypes.number.isRequired,
                            // Add other character props here if known
                        })
                    })
                ),
                Status: PropTypes.string,
            })
        ),
        Roles: PropTypes.arrayOf(PropTypes.string),
        SkillPlans: PropTypes.object,
        SubDirs: PropTypes.array,
        associations: PropTypes.array,
        UserSelections: PropTypes.object,
        SettingsDir: PropTypes.string,
        LastBackupDir: PropTypes.string
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
