import React, { useState, useEffect } from 'react';
import Header from './components/Header';
import Footer from './components/Footer';
import MainContent from './components/MainContent';
import Landing from './components/Landing';
import ErrorBoundary from './components/ErrorBoundary';
import { ToastContainer, toast } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';
import '@fortawesome/fontawesome-free/css/all.min.css'; // Include FontAwesome icons

const App = () => {
    const [homeData, setHomeData] = useState(null);
    const [unassignedCharacters, setUnassignedCharacters] = useState([]);
    const [isAuthenticated, setIsAuthenticated] = useState(false);
    const [isLoading, setIsLoading] = useState(true);

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
        await Promise.all([fetchHomeData(), fetchUnassignedCharacters()]);
        setIsLoading(false);
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

    if (!isAuthenticated) {
        return <Landing />;
    }

    return (
        <ErrorBoundary>
            <div className="flex flex-col min-h-screen bg-gray-900 text-teal-200">
                <Header loggedIn={isAuthenticated} handleLogout={handleLogout} />
                <main className="flex-grow container mx-auto px-4 py-8">
                    <MainContent
                        accounts={homeData?.Accounts || []}
                        unassignedCharacters={unassignedCharacters || []}
                        onAssignCharacter={handleAssignCharacter}
                        onToggleAccountStatus={handleToggleAccountStatus}
                        onUpdateCharacter={handleUpdateCharacter}
                        onUpdateAccountName={handleUpdateAccountName}
                    />
                </main>
                <Footer />
                <ToastContainer />
            </div>
        </ErrorBoundary>
    );
};

export default App;
