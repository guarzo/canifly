// App.jsx
import { useState, useEffect, useCallback } from 'react';
import { HashRouter as Router, Routes, Route } from 'react-router-dom';
import { ThemeProvider } from '@mui/material/styles';
import { ToastContainer, toast } from 'react-toastify';

import { isDev, backEndURL } from './components/config';
import { fetchAppEndpoint } from './utils/api';
import { apiRequest } from './utils/apiRequest'; // we already imported this in a previous step
import { useLoginCallback } from './hooks/useLoginCallback';

import Header from './components/Header';
import Footer from './components/Footer';
import Dashboard from './components/Dashboard';
import SkillPlans from './components/SkillPlans';
import Landing from './components/Landing';
import AddSkillPlanModal from './components/AddSkillPlanModal';
import ErrorBoundary from './components/ErrorBoundary';
import CharacterSort from './components/CharacterSort';
import Sync from './components/Sync';
import Mapping from './components/Mapping';

import theme from './components/theme';
import helloImg from './assets/images/hello.png';

import 'react-toastify/dist/ReactToastify.css';

const App = () => {
    const [appData, setAppData] = useState(null);
    const [isAuthenticated, setIsAuthenticated] = useState(false);
    const [isLoading, setIsLoading] = useState(true);
    const [isSkillPlanModalOpen, setIsSkillPlanModalOpen] = useState(false);
    const [isRefreshing, setIsRefreshing] = useState(false);
    const [loggedOut, setLoggedOut] = useState(false);

    useEffect(() => {
        console.log("isAuthenticated changed:", isAuthenticated);
    }, [isAuthenticated]);

    const wrappedFetchAppEndpoint = useCallback(
        (endpoint, options) => {
            return fetchAppEndpoint(
                {
                    backEndURL,
                    endpoint,
                    loggedOut,
                    isAuthenticated,
                    setIsLoading,
                    setIsRefreshing,
                    setIsAuthenticated,
                    setAppData
                },
                options
            );
        },
        [backEndURL, loggedOut, isAuthenticated]
    );

    const fetchData = useCallback(async () => {
        console.log("fetchData called");
        await wrappedFetchAppEndpoint('/api/app-data', { setLoading: true });
    }, [wrappedFetchAppEndpoint]);

    const loginRefresh = useCallback(() => {
        console.log("loginRefresh called");
        return wrappedFetchAppEndpoint('/api/app-data-no-cache', { setLoading: true, returnSuccess: true });
    }, [wrappedFetchAppEndpoint]);

    const silentRefreshData = useCallback(async () => {
        console.log("silentRefreshData called");
        if (!isAuthenticated || loggedOut) return;
        await wrappedFetchAppEndpoint('/api/app-data-no-cache', { setRefreshing: true });
    }, [wrappedFetchAppEndpoint, isAuthenticated, loggedOut]);

    useEffect(() => {
        console.log("loggedOut changed to:", loggedOut);
        console.trace();
    }, [loggedOut]);

    const loginCallbackFn = useLoginCallback();
    const logInCallBack = (state) => {
        loginCallbackFn(state, {
            isAuthenticated,
            loggedOut,
            loginRefresh,
            setLoggedOut,
            setIsAuthenticated,
            backEndURL
        });
    };

    const handleLogout = useCallback(async () => {
        console.log("handleLogout called");
        console.trace();
        await apiRequest(`${backEndURL}/api/logout`, {
            method: 'POST',
            credentials: 'include'
        }, {
            successMessage: 'Logged out successfully!',
            errorMessage: 'Failed to log out.',
            onSuccess: () => {
                setIsAuthenticated(false);
                setAppData(null);
                setLoggedOut(true);
            }
        });
    }, [backEndURL]);

    const handleToggleAccountStatus = useCallback(async (accountID) => {
        console.log("handleToggleAccountStatus called:", accountID);
        await apiRequest(`${backEndURL}/api/toggle-account-status`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ accountID }),
            credentials: 'include'
        }, {
            successMessage: 'Account status toggled successfully!',
            errorMessage: 'Failed to toggle account status.',
            onSuccess: () => {
                setAppData((prev) => {
                    if (!prev) return prev;
                    const updatedAccounts = prev.Accounts.map((account) =>
                        account.ID === accountID
                            ? { ...account, Status: account.Status === 'Alpha' ? 'Omega' : 'Alpha' }
                            : account
                    );
                    return { ...prev, Accounts: updatedAccounts };
                });
            }
        });
    }, [backEndURL]);

    const handleUpdateCharacter = useCallback(async (characterID, updates) => {
        console.log("handleUpdateCharacter called with characterID:", characterID, "updates:", updates);
        await apiRequest(`${backEndURL}/api/update-character`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ characterID, updates }),
            credentials: 'include'
        }, {
            successMessage: 'Character updated successfully!',
            errorMessage: 'Failed to update character.',
            onSuccess: () => {
                setAppData((prev) => {
                    if (!prev) return prev;
                    const updatedRoles = prev.Roles ? [...prev.Roles] : [];
                    if (updates.Role && !updatedRoles.includes(updates.Role)) {
                        updatedRoles.push(updates.Role);
                    }

                    const updatedAccounts = prev.Accounts.map((account) => {
                        const updatedCharacters = account.Characters.map((character) =>
                            character.Character.CharacterID === characterID
                                ? { ...character, ...updates }
                                : character
                        );
                        return { ...account, Characters: updatedCharacters };
                    });

                    return { ...prev, Accounts: updatedAccounts, Roles: updatedRoles };
                });
            }
        });
    }, [backEndURL]);

    // Refactor handleRemoveCharacter
    const handleRemoveCharacter = useCallback(async (characterID) => {
        console.log("handleRemoveCharacter called with characterID:", characterID);
        await apiRequest(`${backEndURL}/api/remove-character`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ characterID }),
            credentials: 'include'
        }, {
            successMessage: 'Character removed successfully!',
            errorMessage: 'Failed to remove character.',
            onSuccess: () => {
                setAppData((prev) => {
                    if (!prev) return prev;
                    const updatedAccounts = prev.Accounts.map((account) => {
                        const filteredCharacters = account.Characters.filter(
                            (c) => c.Character.CharacterID !== characterID
                        );
                        return { ...account, Characters: filteredCharacters };
                    });
                    return { ...prev, Accounts: updatedAccounts };
                });
            }
        });
    }, [backEndURL]);

    // Refactor handleUpdateAccountName
    const handleUpdateAccountName = useCallback(async (accountID, newName) => {
        console.log("handleUpdateAccountName:", { accountID, newName });
        await apiRequest(`${backEndURL}/api/update-account-name`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ accountID, accountName: newName }),
            credentials: 'include'
        }, {
            successMessage: 'Account name updated successfully!',
            errorMessage: 'Failed to update account name.',
            onSuccess: () => {
                setAppData((prev) => {
                    if (!prev) return prev;
                    const updatedAccounts = prev.Accounts.map((account) =>
                        account.ID === accountID ? { ...account, Name: newName } : account
                    );
                    return { ...prev, Accounts: updatedAccounts };
                });
            }
        });
    }, [backEndURL]);

    // Refactor handleRemoveAccount
    const handleRemoveAccount = useCallback(async (accountName) => {
        console.log("handleRemoveAccount called with accountName:", accountName);
        await apiRequest(`${backEndURL}/api/remove-account`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ accountName }),
            credentials: 'include'
        }, {
            successMessage: 'Account removed successfully!',
            errorMessage: 'Failed to remove account.',
            onSuccess: () => {
                setAppData((prev) => {
                    if (!prev) return prev;
                    const updatedAccounts = prev.Accounts.filter((account) => account.Name !== accountName);
                    return { ...prev, Accounts: updatedAccounts };
                });
            }
        });
    }, [backEndURL]);

    // Refactor handleAddCharacter
    const handleAddCharacter = useCallback(async (account) => {
        console.log("handleAddCharacter called with account:", account);
        await apiRequest(`${backEndURL}/api/add-character`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ account }),
            credentials: 'include'
        }, {
            // No successMessage because we need to handle differently
            errorMessage: 'An error occurred while adding character.',
            onSuccess: (data) => {
                // data should have data.redirectURL if successful
                if (data.redirectURL) {
                    if (isDev) {
                        console.log("Dev mode: redirecting internally to:", data.redirectURL);
                        window.location.href = data.redirectURL;
                    } else {
                        console.log("Production mode: opening external browser to:", data.redirectURL);
                        window.electronAPI.openExternal(data.redirectURL);
                        toast.info("Please authenticate in your browser");
                    }
                } else {
                    toast.error("No redirect URL received from server.");
                }
            }
        });
    }, [backEndURL, isDev]);

    // Refactor handleSaveSkillPlan
    const handleSaveSkillPlan = useCallback(async (planName, planContents) => {
        console.log("handleSaveSkillPlan called with planName:", planName);
        await apiRequest(`${backEndURL}/api/save-skill-plan`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ name: planName, contents: planContents }),
            credentials: 'include'
        }, {
            successMessage: 'Skill Plan Saved!',
            errorMessage: 'Failed to save skill plan.',
            onSuccess: () => {
                setIsSkillPlanModalOpen(false);
                fetchData();
            }
        });
    }, [backEndURL, fetchData]);

    const openSkillPlanModal = () => {
        setIsSkillPlanModalOpen(true);
    };
    const closeSkillPlanModal = () => {
        setIsSkillPlanModalOpen(false);
    };

    useEffect(() => {
        console.log("App mounted, calling fetchData");
        fetchData();
    }, [fetchData]);

    useEffect(() => {
        console.log(`useEffect [isLoading, isAuthenticated]: isLoading=${isLoading}, isAuthenticated=${isAuthenticated}`);
        if (!isLoading && isAuthenticated) {
            silentRefreshData();
        }
    }, [isLoading, isAuthenticated, silentRefreshData]);

    if (isLoading) {
        return (
            <div className="flex flex-col items-center justify-center min-h-screen bg-gray-900 text-teal-200 space-y-4">
                <img
                    src={helloImg}
                    alt="Loading"
                    className="w-32 h-auto object-contain"
                />
                <p className="animate-pulse text-lg">Loading...</p>
            </div>
        );
    }

    const identities = appData?.Accounts?.flatMap((account) => account.Characters) || [];
    const existingAccounts = appData?.Accounts?.map((account) => account.Name) || [];

    return (
        <ErrorBoundary>
            <ThemeProvider theme={theme}>
                <Router>
                    <div className="flex flex-col min-h-screen bg-gray-900 text-teal-200">
                        <Header
                            loggedIn={isAuthenticated}
                            handleLogout={handleLogout}
                            openSkillPlanModal={openSkillPlanModal}
                            existingAccounts={existingAccounts}
                            onSilentRefresh={silentRefreshData}
                            onAddCharacter={handleAddCharacter}
                            isRefreshing={isRefreshing}
                        />
                        <main className="flex-grow container mx-auto px-4 py-8 pb-16">
                            {!isAuthenticated || loggedOut ? (
                                <Landing
                                    backEndURL={backEndURL}
                                    logInCallBack={logInCallBack}
                                />
                            ) : !appData ? (
                                <div className="flex items-center justify-center min-h-screen bg-gray-900 text-teal-200">
                                    <p>Loading...</p>
                                </div>
                            ) : (
                                <Routes>
                                    <Route
                                        path="/"
                                        element={
                                            <Dashboard
                                                accounts={appData?.Accounts || []}
                                                onToggleAccountStatus={handleToggleAccountStatus}
                                                onUpdateCharacter={handleUpdateCharacter}
                                                onUpdateAccountName={handleUpdateAccountName}
                                                onRemoveCharacter={handleRemoveCharacter}
                                                onRemoveAccount={handleRemoveAccount}
                                                roles={appData?.Roles || []}
                                            />
                                        }
                                    />
                                    <Route
                                        path="/skill-plans"
                                        element={
                                            <SkillPlans
                                                identities={identities}
                                                skillPlans={appData?.SkillPlans || {}}
                                                setAppData={setAppData}
                                                backEndURL={backEndURL}
                                            />
                                        }
                                    />
                                    <Route
                                        path="/character-sort"
                                        element={
                                            <CharacterSort
                                                accounts={appData?.Accounts || []}
                                                roles={appData?.Roles || []}
                                                onUpdateCharacter={handleUpdateCharacter}
                                            />
                                        }
                                    />
                                    <Route
                                        path="/sync"
                                        element={
                                            <Sync
                                                settingsData={appData?.SubDirs || []}
                                                associations={appData?.associations || []}
                                                currentSettingsDir={appData?.SettingsDir || ''}
                                                isDefaultDir={appData?.IsDefaultDir ?? false}
                                                userSelections={appData?.UserSelections || {}}
                                                lastBackupDir={appData?.LastBackupDir || ''}
                                                backEndURL={backEndURL}
                                            />
                                        }
                                    />
                                    <Route
                                        path="/mapping"
                                        element={
                                            <Mapping
                                                associations={appData?.associations || []}
                                                subDirs={appData?.SubDirs || []}
                                                onRefreshData={silentRefreshData}
                                                backEndURL={backEndURL}
                                            />
                                        }
                                    />
                                    <Route path="*" element={<div>Route Not Found</div>} />
                                </Routes>
                            )}
                        </main>
                        <Footer />
                        {isSkillPlanModalOpen && (
                            <AddSkillPlanModal
                                onClose={closeSkillPlanModal}
                                onSave={handleSaveSkillPlan}
                            />
                        )}
                        <ToastContainer />
                    </div>
                </Router>
            </ThemeProvider>
        </ErrorBoundary>
    );
};

export default App;
