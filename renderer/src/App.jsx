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
    const [unassignedCharacters, setUnassignedCharacters] = useState(null);  // State to hold unassigned characters
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
                setUnassignedCharacters(null);  // Clear unassigned characters if not authenticated
                toast.warning("Please log in to access your data.");
            } else if (response.ok) {
                const data = await response.json();
                console.log("Home Data:", data);
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
        console.log("Fetching unassigned characters..."); // Log to ensure the function is called
        try {
            const response = await fetch('/api/unassigned-characters', {
                credentials: 'include', // Include credentials to make the request
            });

            if (!response.ok) {
                throw new Error("Failed to fetch unassigned characters");
            }

            const data = await response.json();
            console.log("Unassigned characters data:", data); // Log the response data
            setUnassignedCharacters(data);  // Update state with the unassigned characters
        } catch (error) {
            console.error("Error fetching unassigned characters:", error);
            toast.error("Failed to load unassigned characters.");
            setUnassignedCharacters([]);  // Set to an empty array if the fetch fails
        }
    };

    // Function to handle assigning character to an account
    const handleAssignCharacter = async (characterID, accountID) => {
        try {
            const response = await fetch('/api/assign-character', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    characterID,
                    accountID,
                }),
                credentials: 'include',
            });

            if (!response.ok) {
                throw new Error("Failed to assign character.");
            }

            toast.success("Character assigned successfully!");
            fetchUnassignedCharacters();  // Refresh unassigned characters list
        } catch (error) {
            console.error("Error assigning character:", error);
            toast.error("Failed to assign character.");
        }
    };

    useEffect(() => {
        fetchHomeData();  // Fetch home data on initial load
    }, []);

    // Show loading message while data is being fetched
    if (isLoading) {
        return <div>Loading...</div>;
    }

    // Show landing page if user is not authenticated
    if (!isAuthenticated) {
        return <Landing />;
    }

    // Show main content with the unassigned characters and accounts
    return (
        <ErrorBoundary>
            <div className="min-h-screen flex flex-col">
                <Header loggedIn={isAuthenticated} handleLogout={() => setIsAuthenticated(false)} />
                <main className="flex-grow">
                    <MainContent
                        accounts={homeData?.Accounts || []}  // Pass accounts to MainContent
                        unassignedCharacters={unassignedCharacters}  // Pass unassigned characters here
                        onAssignCharacter={handleAssignCharacter}  // Pass the assign handler to MainContent
                    />
                </main>
                <Footer />
            </div>
            <ToastContainer />
        </ErrorBoundary>
    );
};

export default App;
