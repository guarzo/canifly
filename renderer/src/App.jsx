// src/App.jsx

import { useState, useEffect } from 'react';
import { HashRouter as Router } from 'react-router-dom';
import { ThemeProvider } from '@mui/material/styles';
import { ToastContainer } from 'react-toastify';

import useAuthStore from './stores/authStore';
import useAppDataStore from './stores/appDataStore';
import { useAppData } from './hooks/useAppData';
import { useEveData } from './hooks/useEveData';
import { useWebSocket } from './hooks/useWebSocket';
import { log } from './utils/logger';

import Header from './components/common/Header.jsx';
import Footer from './components/common/Footer.jsx';
import AddSkillPlanModal from './components/skillplan/AddSkillPlanModal.jsx';
import ErrorBoundary from './components/common/ErrorBoundary.jsx';
import FirstRunDialog from './components/setup/FirstRunDialog.jsx';
import AppRoutes from './Routes';
import theme from './Theme.jsx';
import LoadingScreen from './components/ui/LoadingScreen.jsx';
import 'react-toastify/dist/ReactToastify.css';

const App = () => {
    const [isSkillPlanModalOpen, setIsSkillPlanModalOpen] = useState(false);
    const [needsEVEConfig, setNeedsEVEConfig] = useState(false);
    const [configCheckComplete, setConfigCheckComplete] = useState(false);
    
    const { 
        isAuthenticated, 
        authCheckComplete,
        loading: authLoading,
        checkAuth
    } = useAuthStore();
    
    const { 
        accounts,
        config,
        isLoading: dataLoading,
        refreshData,
        fetchAccounts
    } = useAppData();
    
    const { fetchSkillPlans } = useEveData();
    
    // Get fetchAccounts directly from store for WebSocket updates
    const storeFetchAccounts = useAppDataStore(state => state.fetchAccounts);
    
    // Handle WebSocket messages for real-time updates
    const handleWebSocketMessage = (message) => {
        console.log('App.jsx - WebSocket message received:', message);
        log('WebSocket message received:', message);
        
        switch (message.type) {
            case 'account:updated':
                console.log('Account updated - refreshing accounts');
                // Call store's fetchAccounts directly to avoid wrapper issues
                storeFetchAccounts();
                break;
            case 'account:deleted':
                console.log('Account deleted - refreshing accounts');
                storeFetchAccounts();
                break;
            case 'skillplan:created':
            case 'skillplan:updated':
            case 'skillplan:deleted':
                fetchSkillPlans();
                break;
            default:
                console.log('Unknown WebSocket message type:', message.type);
                log('Unknown WebSocket message type:', message.type);
        }
    };
    
    // Connect to WebSocket for real-time updates
    useWebSocket(handleWebSocketMessage);

    // Run auth check only once on app startup
    useEffect(() => {
        if (!authCheckComplete) {
            checkAuth();
        }
    }, []); // Empty deps = run once on mount

    // Check EVE configuration on startup
    useEffect(() => {
        const checkConfig = async () => {
            try {
                const { checkEVEConfiguration } = await import('./api/apiService');
                const result = await checkEVEConfiguration();
                setNeedsEVEConfig(result.needsConfiguration);
            } catch (error) {
                log('Failed to check EVE configuration:', error);
                // Assume config is needed if check fails
                setNeedsEVEConfig(true);
            } finally {
                setConfigCheckComplete(true);
            }
        };
        
        checkConfig();
    }, []);

    useEffect(() => {
        log("isAuthenticated changed:", isAuthenticated);
    }, [isAuthenticated]);

    const isLoading = authLoading || !authCheckComplete || !configCheckComplete || (isAuthenticated && dataLoading);

    if (isLoading) {
        return <LoadingScreen message="Initializing..." />;
    }

    // Show first-run dialog if EVE configuration is needed
    if (needsEVEConfig) {
        return (
            <ThemeProvider theme={theme}>
                <FirstRunDialog 
                    open={true} 
                    onComplete={() => {
                        setNeedsEVEConfig(false);
                        // Dynamic auth client will pick up new credentials automatically
                    }} 
                />
            </ThemeProvider>
        );
    }

    const visibleAccounts = accounts.filter(account => account.Visible !== false);
    const characters = visibleAccounts.flatMap(account => account.Characters || []);
    const existingAccounts = accounts.map(account => account.Name);

    const openSkillPlanModal = () => setIsSkillPlanModalOpen(true);
    const closeSkillPlanModal = () => setIsSkillPlanModalOpen(false);

    return (
        <ErrorBoundary>
            <ThemeProvider theme={theme}>
                <Router>
                    <div className="flex flex-col min-h-screen bg-gray-900 text-teal-200">
                        <Header
                            openSkillPlanModal={openSkillPlanModal}
                            existingAccounts={existingAccounts}
                        />
                        <main className="flex-grow container mx-auto px-4 py-8 pb-16">
                            <AppRoutes
                                characters={characters}
                            />
                        </main>
                        <Footer />
                        {isSkillPlanModalOpen && (
                            <AddSkillPlanModal
                                onClose={closeSkillPlanModal}
                            />
                        )}
                        <ToastContainer 
                            position="top-right"
                            autoClose={3000}
                            hideProgressBar={false}
                            newestOnTop={false}
                            closeOnClick
                            rtl={false}
                            pauseOnFocusLoss
                            draggable
                            pauseOnHover
                            theme="dark"
                        />
                    </div>
                </Router>
            </ThemeProvider>
        </ErrorBoundary>
    );
};

export default App;
