// src/App.jsx

import { useState, useEffect } from 'react';
import { HashRouter as Router } from 'react-router-dom';
import { ThemeProvider } from '@mui/material/styles';
import { ToastContainer } from 'react-toastify';

import { useAuth } from './hooks/useAuth';
import { useAppData } from './hooks/useAppData';
import { log } from './utils/logger';

import Header from './components/common/Header.jsx';
import Footer from './components/common/Footer.jsx';
import AddSkillPlanModal from './components/skillplan/AddSkillPlanModal.jsx';
import ErrorBoundary from './components/common/ErrorBoundary.jsx';
import AppRoutes from './Routes';
import theme from './Theme.jsx';
import helloImg from './assets/images/hello.png';
import 'react-toastify/dist/ReactToastify.css';

const App = () => {
    const [isSkillPlanModalOpen, setIsSkillPlanModalOpen] = useState(false);
    
    const { 
        isAuthenticated, 
        authCheckComplete,
        loading: authLoading 
    } = useAuth();
    
    const { 
        accounts,
        config,
        isLoading: dataLoading,
        refreshData 
    } = useAppData();

    useEffect(() => {
        log("isAuthenticated changed:", isAuthenticated);
    }, [isAuthenticated]);

    const isLoading = authLoading || !authCheckComplete || (isAuthenticated && dataLoading);

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
