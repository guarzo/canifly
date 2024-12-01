// App.jsx
import React, { useState, useEffect } from 'react';
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

const App = () => {
    const [homeData, setHomeData] = useState(null);
    const [unassignedCharacters, setUnassignedCharacters] = useState([]);
    const [isAuthenticated, setIsAuthenticated] = useState(false);
    const [isLoading, setIsLoading] = useState(true);
    const [isSkillPlanModalOpen, setIsSkillPlanModalOpen] = useState(false);

    // Fetch home data
    const fetchHomeData = async () => {
        try {
            const response = await fetch('/api/home-data', {
                credentials: 'include', // Include cookies for authentication
            });
            if (response.status === 401) {
                setIsAuthenticated(false);
                setHomeData(null);
                toast.warning('Please log in to access your data.');
            } else if (response.ok) {
                const data = await response.json();
                setHomeData(data);
                setIsAuthenticated(true);
            } else {
                const errorData = await response.json();
                toast.error(errorData.error || 'An unexpected error occurred.');
            }
        } catch (error) {
            console.error('Error fetching home data:', error);
            toast.error('Failed to load data. Please try again later.');
        }
    };

    // Fetch unassigned characters
    const fetchUnassignedCharacters = async () => {
        try {
            const response = await fetch('/api/unassigned-characters', {
                credentials: 'include', // Include cookies for authentication
            });
            if (response.ok) {
                const data = await response.json();
                setUnassignedCharacters(data || []);
            } else {
                throw new Error('Failed to fetch unassigned characters.');
            }
        } catch (error) {
            console.error('Error fetching unassigned characters:', error);
            toast.error('Failed to load unassigned characters.');
        }
    };

    // Combined data fetching
    const fetchData = async () => {
        setIsLoading(true);
        try {
            await Promise.all([fetchHomeData(), fetchUnassignedCharacters()]);
        } catch (error) {
            console.error('Error fetching data:', error);
        } finally {
            setIsLoading(false);
        }
    };

    // Handle logout
    const handleLogout = async () => {
        try {
            const response = await fetch('/api/logout', {
                method: 'POST',
                credentials: 'include',
            });
            if (response.ok) {
                setIsAuthenticated(false);
                setHomeData(null);
                setUnassignedCharacters([]);
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

    // Assign character to account
    const handleAssignCharacter = async (characterID, accountID) => {
        try {
            const response = await fetch('/api/assign-character', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ characterID, accountID }),
                credentials: 'include',
            });
            if (response.ok) {
                toast.success('Character assigned successfully!');
                fetchData(); // Refresh data
            } else {
                throw new Error('Failed to assign character.');
            }
        } catch (error) {
            console.error('Error assigning character:', error);
            toast.error('Failed to assign character.');
        }
    };

    // Toggle account status (Alpha/Omega)
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
                fetchData(); // Refresh data
            } else {
                throw new Error('Failed to toggle account status.');
            }
        } catch (error) {
            console.error('Error toggling account status:', error);
            toast.error('Failed to toggle account status.');
        }
    };

    // Update character properties (Role, MCT, etc.)
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
                fetchData(); // Refresh data
            } else {
                throw new Error('Failed to update character.');
            }
        } catch (error) {
            console.error('Error updating character:', error);
            toast.error('Failed to update character.');
        }
    };

    // Update account name
    const handleUpdateAccountName = async (accountID, newName) => {
        try {
            const response = await fetch('/api/update-account-name', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ accountID, newName }),
                credentials: 'include',
            });
            if (response.ok) {
                toast.success('Account name updated successfully!');
                fetchData(); // Refresh data
            } else {
                throw new Error('Failed to update account name.');
            }
        } catch (error) {
            console.error('Error updating account name:', error);
            toast.error('Failed to update account name.');
        }
    };

    // Skill Plan Modal Handlers
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
                fetchData(); // Refresh data
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
        fetchData(); // Fetch data on initial load
    }, []);

    if (isLoading) {
        return (
            <div className="flex items-center justify-center min-h-screen bg-gray-900 text-teal-200">
                <p>Loading...</p>
            </div>
        );
    }

    const identities = homeData?.Accounts?.flatMap((account) => account.Characters) || [];
    console.log(homeData)
    console.log(identities)

    return (
        <ErrorBoundary>
            <ThemeProvider theme={theme}>
            <Router>
                <div className="flex flex-col min-h-screen bg-gray-900 text-teal-200">
                    <Header
                        loggedIn={isAuthenticated}
                        handleLogout={handleLogout}
                        openSkillPlanModal={openSkillPlanModal}
                    />
                    <main className="flex-grow container mx-auto px-4 py-8">
                        {!isAuthenticated ? (
                            <Landing />
                        ) : (
                            <Routes>
                                <Route
                                    path="/"
                                    element={
                                        <Dashboard
                                            accounts={homeData?.Accounts || []}
                                            unassignedCharacters={unassignedCharacters || []}
                                            onAssignCharacter={handleAssignCharacter}
                                            onToggleAccountStatus={handleToggleAccountStatus}
                                            onUpdateCharacter={handleUpdateCharacter}
                                            onUpdateAccountName={handleUpdateAccountName}
                                        />
                                    }
                                />
                                <Route
                                    path="/skill-plans"
                                    element={
                                        <SkillPlans
                                            identities={identities}
                                            skillPlans={homeData?.SkillPlans || {}}
                                        />
                                    }
                                />
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

