import { useState, useEffect } from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
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

const App = () => {
    const [appData, setAppData] = useState(null);
    const [isAuthenticated, setIsAuthenticated] = useState(false);
    const [isLoading, setIsLoading] = useState(true);
    const [isSkillPlanModalOpen, setIsSkillPlanModalOpen] = useState(false);

    const fetchAppData = async () => {
        try {
            const response = await fetch('/api/app-data', {
                credentials: 'include',
            });
            if (response.status === 401) {
                setIsAuthenticated(false);
                setAppData(null);
                toast.warning('Please log in to access your data.');
            } else if (response.ok) {
                const data = await response.json();
                setAppData(data);
                setIsAuthenticated(true);
            } else {
                const errorData = await response.json();
                toast.error(errorData.error || 'An unexpected error occurred.');
            }
        } catch (error) {
            if (isAuthenticated) {
                console.error('Error fetching app data:', error);
                toast.error('Failed to load data. Please try again later.');
            }
        }
    };

    const fetchAppDataNoCache = async () => {
        try {
            const response = await fetch('/api/app-data-no-cache', {
                credentials: 'include',
            });

            if (response.ok) {
                const data = await response.json();
                setAppData(data);
                setIsAuthenticated(true);
            } else {
                const errorData = await response.json();
                console.log(errorData.error);
            }
        } catch (error) {
            if (isAuthenticated) {
                console.error('Error fetching no cache app data:', error);
            }
        }
    };

    const fetchData = async () => {
        setIsLoading(true);
        try {
            await fetchAppData();
        } catch (error) {
            if (isAuthenticated) {
                console.error('Error fetching data:', error);
            }
        } finally {
            setIsLoading(false);
        }
    };

    const silentRefreshData = async () => {
        try {
            await fetchAppDataNoCache();
        } catch (error) {
            if (isAuthenticated) {
                console.error('Error fetching data:', error);
            }
        }
    };

    const handleLogout = async () => {
        try {
            const response = await fetch('/api/logout', {
                method: 'POST',
                credentials: 'include',
            });
            if (response.ok) {
                setIsAuthenticated(false);
                setAppData(null);
                toast.success('Logged out successfully!');
            } else {
                const errorData = await response.json();
                toast.error(errorData.error || 'Failed to log out.');
            }
        } catch (error) {
            console.error('Error logging out:', error);
            toast.error('An error occurred during logout.');
        }
    };

    const handleToggleAccountStatus = async (accountID) => {
        try {
            const response = await fetch('/api/toggle-account-status', {
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
        try {
            const response = await fetch('/api/update-character', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ characterID, updates }),
                credentials: 'include',
            });

            if (response.ok) {
                toast.success('Character updated successfully!');
                setAppData((prevAppData) => {
                    const updatedRoles = prevAppData?.Roles
                        ? [...prevAppData.Roles]
                        : [];
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
        try {
            const response = await fetch('/api/remove-character', {
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
        try {
            const response = await fetch('/api/update-account-name', {
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
        try {
            const response = await fetch('/api/remove-account', {
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

    // NEW FUNCTION: handleAddCharacter (moved from Header)
    const handleAddCharacter = async (account) => {
        try {
            const response = await fetch('/api/add-character', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({account}),
                credentials: 'include',
            });
            if (response.ok) {
                const data = await response.json();

                // If data includes a redirectURL, we must redirect.
                // After redirect, the backend callback should update the app state.
                // If we want immediate local update, we must have character details in `data`
                // For now, assume we get some character details in `data.newCharacter`:
                if (data.redirectURL) {
                    window.location.href = data.redirectURL;
                } else if (data.newCharacter && data.accountName) {
                    // Update the local state immediately without waiting for refresh
                    setAppData((prevAppData) => {
                        const updatedAccounts = prevAppData.Accounts.map((acc) => {
                            if (acc.Name === data.accountName) {
                                return {
                                    ...acc,
                                    Characters: [...acc.Characters, data.newCharacter],
                                };
                            }
                            return acc;
                        });
                        return { ...prevAppData, Accounts: updatedAccounts };
                    });
                    toast.success('Character added successfully!');
                } else {
                    toast.error("No redirect URL or character details received from server.");
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

    const openSkillPlanModal = () => setIsSkillPlanModalOpen(true);
    const closeSkillPlanModal = () => setIsSkillPlanModalOpen(false);

    const handleSaveSkillPlan = async (planName, planContents) => {
        try {
            const response = await fetch('/api/save-skill-plan', {
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
        fetchData();
    }, []);

    useEffect(() => {
        if (!isLoading && isAuthenticated) {
            silentRefreshData();
        }
    }, [isLoading, isAuthenticated]);

    if (isLoading) {
        return (
            <div className="flex flex-col items-center justify-center min-h-screen bg-gray-900 text-teal-200 space-y-4">
                <img
                    src="/images/hello.png"
                    alt="Loading"
                    className="w-32 h-auto object-contain"
                />
                <p className="animate-pulse text-lg">Loading...</p>
            </div>
        );
    }

    const identities = appData?.Accounts?.flatMap((account) => account.Characters) || [];
    const existingAccounts = appData?.Accounts.flatMap((account) => account.Name) || [];
    console.log(appData);

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
                            onAddCharacter={handleAddCharacter} {/* PASS THE FUNCTION DOWN */}
                        />
                        <main className="flex-grow container mx-auto px-4 py-8">
                            {!isAuthenticated ? (
                                <Landing/>
                            ) : isLoading || !appData ? (
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
                                            />
                                        }
                                    />
                                    <Route
                                        path="/character-sort"
                                        element={
                                            <CharacterSort
                                                accounts={appData?.Accounts || []}
                                                roles={appData?.Roles || []}
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
                                            />
                                        }
                                    />
                                </Routes>
                            )}
                        </main>
                        <Footer/>
                        {isSkillPlanModalOpen && (
                            <AddSkillPlanModal
                                onClose={closeSkillPlanModal}
                                onSave={handleSaveSkillPlan}
                            />
                        )}
                        <ToastContainer/>
                    </div>
                </Router>
            </ThemeProvider>
        </ErrorBoundary>
    );
};

export default App;
