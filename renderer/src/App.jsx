import React, { useState, useEffect } from 'react';
import { ThemeProvider, createTheme } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import Header from './components/Header';
import Footer from './components/Footer';
import MainContent from './components/MainContent';
import Landing from './components/Landing';
import ErrorBoundary from './components/ErrorBoundary';
import { ToastContainer, toast } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';

const App = () => {
    const [homeData, setHomeData] = useState(null);
    const [unassignedCharacters, setUnassignedCharacters] = useState([]);
    const [isAuthenticated, setIsAuthenticated] = useState(false);
    const [isLoading, setIsLoading] = useState(true);

    const theme = createTheme({
        palette: {
            mode: 'dark',
            primary: {
                main: '#007bff',
            },
            secondary: {
                main: '#6c757d',
            },
        },
    });

    const fetchHomeData = async () => {
        try {
            const response = await fetch('/api/home-data', {
                credentials: 'include',
            });
            if (response.status === 401) {
                setIsAuthenticated(false);
                setHomeData(null);
                setUnassignedCharacters([]);
                toast.warning('Please log in to access your data.');
            } else if (response.ok) {
                const data = await response.json();
                setHomeData(data);
                setIsAuthenticated(true);
                fetchUnassignedCharacters();
            } else {
                const errorData = await response.json();
                toast.error(errorData.error || 'An unexpected error occurred.');
                setIsAuthenticated(false);
            }
        } catch (error) {
            console.error('Error fetching home data:', error);
            toast.error('Failed to load data. Please try again later.');
            setIsAuthenticated(false);
        }
        setIsLoading(false);
    };

    const fetchUnassignedCharacters = async () => {
        try {
            const response = await fetch('/api/unassigned-characters', {
                credentials: 'include',
            });
            if (!response.ok) throw new Error('Failed to fetch unassigned characters');
            const data = await response.json();
            setUnassignedCharacters(data);
        } catch (error) {
            console.error('Error fetching unassigned characters:', error);
            toast.error('Failed to load unassigned characters.');
            setUnassignedCharacters([]);
        }
    };

    // Assign character to an account
    const handleAssignCharacter = async (characterID, accountID) => {
        try {
            const response = await fetch('/api/assign-character', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ characterID }),
                credentials: 'include',
            });
            if (!response.ok) throw new Error('Failed to assign character.');
            toast.success('Character assigned successfully!');
            await fetchHomeData();
        } catch (error) {
            console.error('Error assigning character:', error);
            toast.error('Failed to assign character.');
        }
    };

    // Update character properties (role, MCT, etc.)
    const handleUpdateCharacter = async (characterID, updates) => {
        try {
            const response = await fetch('/api/update-character', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ characterID, updates }),
                credentials: 'include',
            });
            if (!response.ok) throw new Error('Failed to update character.');
            toast.success('Character updated successfully!');
            await fetchHomeData();
        } catch (error) {
            console.error('Error updating character:', error);
            toast.error('Failed to update character.');
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
            if (!response.ok) throw new Error('Failed to toggle account status.');
            toast.success('Account status toggled successfully!');
            await fetchHomeData();
        } catch (error) {
            console.error('Error toggling account status:', error);
            toast.error('Failed to toggle account status.');
        }
    };

    useEffect(() => {
        fetchHomeData();
    }, []);

    if (isLoading) return <div>Loading...</div>;
    if (!isAuthenticated) return <Landing />;

    return (
        <ErrorBoundary>
            <div className="flex flex-col min-h-screen bg-gray-900 text-teal-200">
                <Header
                    loggedIn={isAuthenticated}
                    handleLogout={() => setIsAuthenticated(false)}
                />
                <main className="flex-grow container mx-auto px-4 py-6">
                    {isAuthenticated ? (
                        <MainContent
                            accounts={homeData?.Accounts || []}
                            unassignedCharacters={unassignedCharacters}
                            onAssignCharacter={(characterID) => console.log('Assign Character:', characterID)}
                            onUpdateCharacter={(characterID, updates) =>
                                console.log('Update Character:', characterID, updates)
                            }
                            onToggleAccountStatus={(accountID) =>
                                console.log('Toggle Account Status:', accountID)
                            }
                        />
                    ) : (
                        <Landing />
                    )}
                </main>
                <Footer />
                <ToastContainer />
            </div>
        </ErrorBoundary>
    );
};

export default App;
