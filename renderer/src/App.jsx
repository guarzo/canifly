import { useState, useEffect } from 'react';
import { HashRouter as Router, Routes, Route } from 'react-router-dom';
import Header from './components/Header';
import Footer from './components/Footer';
import Dashboard from './components/Dashboard';
import SkillPlans from './components/SkillPlans';
import Landing from './components/Landing';
import AddSkillPlanModal from './components/AddSkillPlanModal';
import { ThemeProvider } from '@mui/material/styles';
import ErrorBoundary from './components/ErrorBoundary';
import { ToastContainer, toast } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';
import theme from './components/theme';
import CharacterSort from './components/CharacterSort.jsx';
import Sync from './components/Sync';
import Mapping from './components/Mapping';
import helloImg from './assets/images/hello.png';

const App = () => {
    const [appData, setAppData] = useState(null);
    const [isAuthenticated, setIsAuthenticated] = useState(false);
    const [isLoading, setIsLoading] = useState(true);
    const [isSkillPlanModalOpen, setIsSkillPlanModalOpen] = useState(false);
    const [isRefreshing, setIsRefreshing] = useState(false);
    const [loggedOut, setLoggedOut] = useState(false);

    const isDev = import.meta.env.DEV;
    const backEndURL = isDev ? '' : 'http://localhost:8713';

    // Log whenever isAuthenticated changes
    useEffect(() => {
        console.log("isAuthenticated changed:", isAuthenticated);
    }, [isAuthenticated]);

    /**
     * A generic fetch function to handle various endpoints and update state.
     * @param {string} endpoint - endpoint to fetch (e.g. "/api/app-data")
     * @param {object} options - { setLoading: boolean, setRefreshing: boolean, returnSuccess: boolean }
     * @returns {boolean|undefined} If returnSuccess is true, returns boolean indicating success/failure, else undefined.
     */
    async function fetchAppEndpoint(endpoint, { setLoading = false, setRefreshing = false, returnSuccess = false } = {}) {
        console.log(`fetchAppEndpoint called with endpoint=${endpoint}, setLoading=${setLoading}, setRefreshing=${setRefreshing}, returnSuccess=${returnSuccess}`);
        console.log("State before fetch:", {
            isAuthenticated,
            isLoading,
            loggedOut,
            appDataExists: !!appData
        });

        if (loggedOut) {
            console.log("fetchAppEndpoint: loggedOut is true, returning early");
            return returnSuccess ? false : undefined;
        }

        if (setLoading) {
            console.log("Setting isLoading = true");
            setIsLoading(true);
        }
        if (setRefreshing) {
            console.log("Setting isRefreshing = true");
            setIsRefreshing(true);
        }

        try {
            const response = await fetch(`${backEndURL}${endpoint}`, { credentials: 'include' });

            if (response.status === 401) {
                console.log("fetchAppEndpoint: got 401 Unauthorized, setting isAuthenticated=false, clearing appData");
                setIsAuthenticated(false);
                setAppData(null);
                toast.warning('Please log in to access your data.');
                return returnSuccess ? false : undefined;
            }

            if (response.ok) {
                const data = await response.json();
                console.log("fetchAppEndpoint: response OK, setting appData and isAuthenticated=true");
                setAppData(data);
                setIsAuthenticated(true);
                return returnSuccess ? true : undefined;
            } else {
                const errorData = await response.json();
                console.log("fetchAppEndpoint: response not OK, error:", errorData.error);
                toast.error(errorData.error || 'An unexpected error occurred.');
                return returnSuccess ? false : undefined;
            }
        } catch (error) {
            console.error('Error fetching data:', error);
            if (isAuthenticated) {
                toast.error('Failed to load data. Please try again later.');
            }
            return returnSuccess ? false : undefined;
        } finally {
            if (setLoading) {
                console.log("Setting isLoading = false in finally");
                setIsLoading(false);
            }
            if (setRefreshing) {
                console.log("Setting isRefreshing = false in finally");
                setIsRefreshing(false);
            }
        }
    }

    // Initial fetch of app data
    async function fetchData() {
        console.log("fetchData called");
        await fetchAppEndpoint('/api/app-data', { setLoading: true });
    }

    // Used after login finalize to confirm we have data and session
    async function loginRefresh() {
        console.log("loginRefresh called");
        return fetchAppEndpoint('/api/app-data-no-cache', { setLoading: true, returnSuccess: true });
    }

    // Silent refresh data (used periodically after authenticated)
    async function silentRefreshData() {
        console.log("silentRefreshData called");
        if (!isAuthenticated || loggedOut) {
            console.log(`silentRefreshData not running because isAuthenticated=${isAuthenticated}, loggedOut=${loggedOut}`);
            return;
        }
        await fetchAppEndpoint('/api/app-data-no-cache', { setRefreshing: true });
    }

    useEffect(() => {
        console.log("loggedOut changed to:", loggedOut);
        console.trace(); // Shows the stack trace of what caused this effect
    }, [loggedOut]);


    // Callback after initiating login
    const logInCallBack = (state) => {
        console.log("logInCallBack called with state:", state);
        setLoggedOut(false);
        let attempts = 0;
        const maxAttempts = 6;
        let finalized = false; // to ensure we only finalize once

        const interval = setInterval(async () => {
            console.log(`Interval tick: attempts=${attempts}, isAuthenticated=${isAuthenticated}, finalized=${finalized}, loggedOut=${loggedOut}`);

            if (isAuthenticated) {
                console.log("logged in callback stopped - user logged in, clearing interval");
                clearInterval(interval);
                return;
            }

            console.log(`attempt ${attempts} trying to detect login`);
            attempts++;
            if (attempts > maxAttempts) {
                console.warn('Failed to detect login after multiple attempts, clearing interval.');
                clearInterval(interval);
                return;
            }

            if (!finalized) {
                console.log("Calling finalize-login endpoint...");
                const finalizeResp = await fetch(`${backEndURL}/api/finalize-login?state=${state}`, {
                    credentials: 'include'
                });

                if (finalizeResp.ok) {
                    finalized = true;
                    console.log("Finalization succeeded, now trying to fetch data via loginRefresh...");
                    const success = await loginRefresh();
                    console.log("loginRefresh returned:", success);
                    if (success) {
                        console.log("Login finalized and data fetched! Setting isAuthenticated=true and clearing interval");
                        setIsAuthenticated(true);
                        clearInterval(interval);
                    } else {
                        console.log("Session set but data fetch failed, will retry data fetch on next interval...");
                    }
                } else {
                    console.log("Not ready yet, retrying finalize-login...");
                }
            } else {
                console.log("Already finalized, just trying loginRefresh again...");
                const success = await loginRefresh();
                console.log("loginRefresh returned:", success);
                if (success) {
                    console.log("Data fetched after finalization! Setting isAuthenticated=true and clearing interval");
                    setIsAuthenticated(true);
                    clearInterval(interval);
                } else {
                    console.log("Still no data after finalization, retrying data fetch on next interval...");
                }
            }
        }, 5000);
        console.log("logInCallBack: interval created");
    };

    const handleLogout = async () => {
        console.log("handleLogout called");
        console.trace(); // Trace what caused handleLogout
        try {
            const response = await fetch(`${backEndURL}/api/logout`, {
                method: 'POST',
                credentials: 'include',
            });
            if (response.ok) {
                setIsAuthenticated(false);
                setAppData(null);
                console.log("About to set loggedOut=true in handleLogout");
                console.trace(); // Shows what triggered this state change
                setLoggedOut(true);
                toast.success('Logged out successfully!');
                console.log("Logout succeeded");
            } else {
                const errorData = await response.json();
                toast.error(errorData.error || 'Failed to log out.');
                console.log("Logout failed:", errorData);
            }
        } catch (error) {
            console.error('Error logging out:', error);
            toast.error('An error occurred during logout.');
        }
    };

    const handleToggleAccountStatus = async (accountID) => {
        console.log("handleToggleAccountStatus called with accountID:", accountID);
        try {
            const response = await fetch(`${backEndURL}/api/toggle-account-status`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ accountID }),
                credentials: 'include',
            });

            if (response.ok) {
                toast.success('Account status toggled successfully!');
                setAppData((prevAppData) => {
                    const updatedAccounts = prevAppData.Accounts.map((account) => {
                        if (account.ID === accountID) {
                            const newStatus = account.Status === 'Alpha' ? 'Omega' : 'Alpha';
                            return { ...account, Status: newStatus };
                        }
                        return account;
                    });
                    return { ...prevAppData, Accounts: updatedAccounts };
                });
            } else {
                throw new Error('Failed to toggle account status.');
            }
        } catch (error) {
            console.error('Error toggling account status:', error);
            toast.error('Failed to toggle account status.');
        }
    };

    const handleUpdateCharacter = async (characterID, updates) => {
        console.log("handleUpdateCharacter called with characterID:", characterID, "updates:", updates);
        try {
            const response = await fetch(`${backEndURL}/api/update-character`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ characterID, updates }),
                credentials: 'include',
            });

            if (response.ok) {
                toast.success('Character updated successfully!');
                setAppData((prevAppData) => {
                    const updatedRoles = prevAppData?.Roles ? [...prevAppData.Roles] : [];
                    if (updates.Role && !updatedRoles.includes(updates.Role)) {
                        updatedRoles.push(updates.Role);
                    }

                    const updatedAccounts = prevAppData.Accounts.map((account) => {
                        const updatedCharacters = account.Characters.map((character) => {
                            if (character.Character.CharacterID === characterID) {
                                return { ...character, ...updates };
                            }
                            return character;
                        });
                        return { ...account, Characters: updatedCharacters };
                    });

                    return {
                        ...prevAppData,
                        Accounts: updatedAccounts,
                        Roles: updatedRoles,
                    };
                });
            } else {
                throw new Error('Failed to update character.');
            }
        } catch (error) {
            console.error('Error updating character:', error);
            toast.error('Failed to update character.');
        }
    };

    const handleRemoveCharacter = async (characterID) => {
        console.log("handleRemoveCharacter called with characterID:", characterID);
        try {
            const response = await fetch(`${backEndURL}/api/remove-character`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ characterID }),
                credentials: 'include',
            });

            const result = await response.json();
            if (response.ok && result.success) {
                toast.success('Character removed successfully!');
                setAppData((prevAppData) => {
                    const updatedAccounts = prevAppData.Accounts.map((account) => {
                        const filteredCharacters = account.Characters.filter((character) => {
                            return character.Character.CharacterID !== characterID;
                        });
                        return { ...account, Characters: filteredCharacters };
                    });
                    return { ...prevAppData, Accounts: updatedAccounts };
                });
            } else {
                toast.error(result.error || 'Failed to remove character.');
            }
        } catch (error) {
            console.error('Error removing character:', error);
            toast.error('Error removing character.');
        }
    };

    const handleUpdateAccountName = async (accountID, newName) => {
        console.log("handleUpdateAccountName called with accountID:", accountID, "newName:", newName);
        try {
            const response = await fetch(`${backEndURL}/api/update-account-name`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ accountID, accountName: newName }),
                credentials: 'include',
            });

            const result = await response.json();
            if (response.ok && result.success) {
                toast.success('Account name updated successfully!');
                setAppData((prevAppData) => {
                    const updatedAccounts = prevAppData.Accounts.map((account) =>
                        account.ID === accountID ? { ...account, Name: newName } : account
                    );
                    return { ...prevAppData, Accounts: updatedAccounts };
                });
            } else {
                toast.error(result.error || 'Failed to update account name.');
            }
        } catch (error) {
            console.error('Error updating account name:', error);
            toast.error('Error updating account name.');
        }
    };

    const handleRemoveAccount = async (accountName) => {
        console.log("handleRemoveAccount called with accountName:", accountName);
        try {
            const response = await fetch(`${backEndURL}/api/remove-account`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ accountName }),
                credentials: 'include',
            });

            const result = await response.json();
            if (response.ok && result.success) {
                toast.success('Account removed successfully!');
                setAppData((prevAppData) => {
                    const updatedAccounts = prevAppData.Accounts.filter(
                        (account) => account.Name !== accountName
                    );
                    return { ...prevAppData, Accounts: updatedAccounts };
                });
            } else {
                toast.error(result.error || 'Failed to remove account.');
            }
        } catch (error) {
            console.error('Error removing account:', error);
            toast.error('Error removing account.');
        }
    };

    const handleAddCharacter = async (account) => {
        console.log("handleAddCharacter called with account:", account);
        try {
            const response = await fetch(`${backEndURL}/api/add-character`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ account }),
                credentials: 'include',
            });
            if (response.ok) {
                const data = await response.json();
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
            } else {
                const errorText = await response.text();
                toast.error(`Failed to initiate login: ${errorText}`);
            }
        } catch (error) {
            console.error('Error initiating add character:', error);
            toast.error('An error occurred while adding character.');
        }
    };

    const openSkillPlanModal = () => {
        console.log("openSkillPlanModal called");
        setIsSkillPlanModalOpen(true);
    };
    const closeSkillPlanModal = () => {
        console.log("closeSkillPlanModal called");
        setIsSkillPlanModalOpen(false);
    };

    const handleSaveSkillPlan = async (planName, planContents) => {
        console.log("handleSaveSkillPlan called with planName:", planName);
        try {
            const response = await fetch(`${backEndURL}/api/save-skill-plan`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ name: planName, contents: planContents }),
                credentials: 'include',
            });

            if (response.ok) {
                toast.success('Skill Plan Saved!');
                closeSkillPlanModal();
                fetchData();
            } else {
                const errorText = await response.text();
                toast.error(`Error saving skill plan: ${errorText}`);
            }
        } catch (error) {
            console.error('Error saving skill plan:', error);
            toast.error('Failed to save skill plan.');
        }
    };

    useEffect(() => {
        console.log("App mounted, calling fetchData");
        fetchData();
    }, []);

    useEffect(() => {
        console.log(`useEffect [isLoading, isAuthenticated]: isLoading=${isLoading}, isAuthenticated=${isAuthenticated}`);
        if (!isLoading && isAuthenticated) {
            console.log("Condition met to call silentRefreshData");
            silentRefreshData();
        }
    }, [isLoading, isAuthenticated]);

    console.log("Render check before return:", {
        isAuthenticated,
        loggedOut,
        isLoading,
        appDataExists: !!appData,
        appData
    });

    if (isLoading) {
        console.log("Rendering loading screen because isLoading =", isLoading);
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
    const existingAccounts = appData?.Accounts?.flatMap((account) => account.Name) || [];

    if (!isAuthenticated || loggedOut) {
        console.log("Rendering Landing (not authenticated or loggedOut)");
    } else if (!appData) {
        console.log("Rendering secondary loading screen (no appData)");
    } else {
        console.log("Rendering main Routes (authenticated and have appData)");
    }

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
