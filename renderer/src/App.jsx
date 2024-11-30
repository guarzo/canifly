import React, { useState, useEffect } from 'react';
import Header from './components/Header';
import Footer from './components/Footer';
import MainContent from './components/MainContent';
import Landing from './components/Landing';
import ErrorBoundary from './components/ErrorBoundary';
import { ToastContainer, toast } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';
import '@fortawesome/fontawesome-free/css/all.min.css';

const App = () => {
    const [homeData, setHomeData] = useState(null);
    const [unassignedCharacters, setUnassignedCharacters] = useState([]); // State for unassigned characters
    const [isAuthenticated, setIsAuthenticated] = useState(false);
    const [isLoading, setIsLoading] = useState(true);

    // Fetch home data
    const fetchHomeData = async () => {
        try {
            const response = await fetch('/api/home-data', {
                credentials: 'include', // Ensure cookies are sent
            });

            if (response.status === 401) {
                setIsAuthenticated(false);
                setHomeData(null);
                setUnassignedCharacters([]);
                toast.warning("Please log in to access your data.");
            } else if (response.ok) {
                const data = await response.json();
                setHomeData(data);
                setIsAuthenticated(true);
                fetchUnassignedCharacters();  // Fetch unassigned characters after successful login
            } else {
                const errorData = await response.json();
                toast.error(errorData.error || "An unexpected error occurred.");
                setIsAuthenticated(false);
            }
        } catch (error) {
            console.error("Error fetching home data:", error);
            toast.error("Failed to load data. Please try again later.");
            setIsAuthenticated(false);
        }
        setIsLoading(false);
    };

    // Fetch unassigned characters
    const fetchUnassignedCharacters = async () => {
        try {
            const response = await fetch('/api/unassigned-characters', {
                credentials: 'include',
            });

            if (!response.ok) {
                throw new Error("Failed to fetch unassigned characters");
            }

            const data = await response.json();
            setUnassignedCharacters(data);
        } catch (error) {
            console.error("Error fetching unassigned characters:", error);
            toast.error("Failed to load unassigned characters.");
            setUnassignedCharacters([]); // Empty array if fetch fails
        }
    };

    // Handle character assignment
    const handleAssignCharacter = async (characterID, accountID) => {
        try {
            const response = await fetch('/api/assign-character', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ characterID }),
                credentials: 'include',
            });

            if (!response.ok) {
                throw new Error("Failed to assign character.");
            }

            toast.success("Character assigned successfully!");

            // After assignment, re-fetch the unassigned characters
            fetchUnassignedCharacters();
        } catch (error) {
            console.error("Error assigning character:", error);
            toast.error("Failed to assign character.");
        }
    };

    useEffect(() => {
        fetchHomeData();  // Fetch home data on initial load
    }, []);

    if (isLoading) {
        return <div className="bg-gray-800 text-white p-4 rounded shadow-lg">Loading...</div>;
    }

    if (!isAuthenticated) {
        return <Landing />;
    }

    return (
        <ErrorBoundary>
            <div className="min-h-screen flex flex-col bg-gray-900 text-teal-200">
                <Header loggedIn={isAuthenticated} handleLogout={() => setIsAuthenticated(false)} />
                <main className="flex-grow bg-gray-800 p-4">
                    <MainContent
                        accounts={homeData?.Accounts || []}
                        unassignedCharacters={unassignedCharacters} // Pass updated unassigned characters
                        onAssignCharacter={handleAssignCharacter}  // Pass the assign handler
                    />
                </main>
                <Footer />
            </div>
            <ToastContainer />
        </ErrorBoundary>
    );
};

export default App;
