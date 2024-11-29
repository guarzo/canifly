import React, { useState, useEffect } from 'react';
import Header from './components/Header';
import Footer from './components/Footer';
import MainContent from './components/MainContent';
import AddSkillPlanModal from './components/AddSkillPlanModal';
import Landing from './components/Landing';
import ErrorBoundary from './components/ErrorBoundary';
import { ToastContainer, toast } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';
import '@fortawesome/fontawesome-free/css/all.min.css';

const App = () => {
    const [homeData, setHomeData] = useState(null);
    const [isAuthenticated, setIsAuthenticated] = useState(false);
    const [isLoading, setIsLoading] = useState(true);
    const [isSkillPlanModalOpen, setIsSkillPlanModalOpen] = useState(false);

    useEffect(() => {
        const fetchData = async () => {
            try {
                const response = await fetch('/api/home-data', {
                    credentials: 'include', // Ensure cookies are sent
                });
                if (response.status === 401) {
                    setIsAuthenticated(false);
                    toast.warning("Please log in to access your data.");
                } else if (response.ok) {
                    const data = await response.json();
                    console.log("Home Data:", data);
                    setHomeData(data);
                    setIsAuthenticated(true);
                } else {
                    const errorData = await response.json();
                    toast.error(errorData.error || "An unexpected error occurred.");
                    setIsAuthenticated(false);
                }
                setIsLoading(false);
            } catch (error) {
                console.error("Error fetching home data:", error);
                toast.error("Failed to load data. Please try again later.");
                setIsAuthenticated(false);
                setIsLoading(false);
            }
        };

        fetchData();
    }, []);

    const handleLogout = async () => {
        try {
            const response = await fetch('/api/logout', {
                method: 'POST',
                credentials: 'include',
            });
            if (response.ok) {
                const data = await response.json();
                if (data.success) {
                    setIsAuthenticated(false);
                    toast.success("Logged out successfully.");
                } else {
                    toast.error("Logout failed.");
                }
            } else {
                const errorText = await response.text();
                toast.error(`Logout failed: ${errorText}`);
            }
        } catch (error) {
            console.error("Error during logout:", error);
            toast.error("An error occurred during logout.");
        }
    };

    useEffect(() => {
        const handleOpenSkillPlanModal = () => {
            setIsSkillPlanModalOpen(true);
        };

        window.addEventListener('openAddSkillPlanModal', handleOpenSkillPlanModal);

        return () => {
            window.removeEventListener('openAddSkillPlanModal', handleOpenSkillPlanModal);
        };
    }, []);

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
                toast.success("Skill Plan Saved!");
                setIsSkillPlanModalOpen(false);

                const updatedHomeDataResponse = await fetch('/api/home-data', {
                    credentials: 'include',
                });
                if (updatedHomeDataResponse.ok) {
                    const updatedData = await updatedHomeDataResponse.json();
                    setHomeData(updatedData);
                    setIsAuthenticated(true);
                } else {
                    toast.error("Failed to refresh data after saving skill plan.");
                }
            } else {
                const errorText = await response.text();
                toast.error(`Error saving skill plan: ${errorText}`);
            }
        } catch (error) {
            console.error("Error saving skill plan:", error);
            toast.error("Failed to save skill plan.");
        }
    };

    console.log("Passing to MainContent:", {
        identities: homeData?.TabulatorIdentities || [],
        skillPlans: homeData?.TabulatorSkillPlans || {},
        matchingSkillPlans: homeData?.MatchingSkillPlans || {},
        matchingCharacters: homeData?.MatchingCharacters || {},
    });

    if (isLoading) {
        return (
            <div className="flex items-center justify-center min-h-screen bg-gradient-to-b from-gray-800 to-gray-700 text-gray-100">
                <div className="text-center">
                    <svg
                        className="animate-spin h-12 w-12 text-teal-400 mx-auto mb-4"
                        xmlns="http://www.w3.org/2000/svg"
                        fill="none"
                        viewBox="0 0 24 24"
                    >
                        <circle
                            className="opacity-25"
                            cx="12"
                            cy="12"
                            r="10"
                            stroke="currentColor"
                            strokeWidth="4"
                        ></circle>
                        <path
                            className="opacity-75"
                            fill="currentColor"
                            d="M4 12a8 8 0 018-8v8H4z"
                        ></path>
                    </svg>
                    <p className="text-xl text-teal-200">Loading...</p>
                </div>
            </div>
        );
    }

    if (!isAuthenticated) {
        return <Landing />;
    }

    return (
        <ErrorBoundary>
            <div className="flex flex-col min-h-screen bg-gradient-to-b from-gray-800 to-gray-700 text-gray-100">
                <Header
                    title={homeData?.Title || "Can I Fly?"}
                    loggedIn={homeData?.LoggedIn || false}
                    handleLogout={handleLogout}
                />
                <MainContent
                    identities={homeData?.TabulatorIdentities || []}
                    skillPlans={homeData?.TabulatorSkillPlans || {}}
                    matchingSkillPlans={homeData?.MatchingSkillPlans || {}}
                    matchingCharacters={homeData?.MatchingCharacters || {}}
                />
                <Footer />
                {isSkillPlanModalOpen && (
                    <AddSkillPlanModal
                        onClose={closeSkillPlanModal}
                        onSave={handleSaveSkillPlan}
                    />
                )}
                <ToastContainer />
            </div>
        </ErrorBoundary>
    );
};

export default App;
