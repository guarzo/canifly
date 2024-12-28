// src/App.jsx

import { useState, useEffect, useCallback } from 'react';
import { HashRouter as Router } from 'react-router-dom';
import { ThemeProvider } from '@mui/material/styles';
import { ToastContainer } from 'react-toastify';

import { useLoginCallback } from './hooks/useLoginCallback';
import { log, trace } from './utils/logger';
import { useAppHandlers } from './hooks/useAppHandlers';
import { getAppData, getAppDataNoCache } from './api/apiService';

import Header from './components/common/Header.jsx';
import Footer from './components/common/Footer.jsx';
import AddSkillPlanModal from './components/skillplan/AddSkillPlanModal.jsx';
import ErrorBoundary from './components/common/ErrorBoundary.jsx';
import AppRoutes from './Routes';
import theme from './Theme.jsx';
import helloImg from './assets/images/hello.png';
import 'react-toastify/dist/ReactToastify.css';

const App = () => {
    const [appData, setAppData] = useState(null);
    const [isAuthenticated, setIsAuthenticated] = useState(false);
    const [isLoading, setIsLoading] = useState(true);
    const [isSkillPlanModalOpen, setIsSkillPlanModalOpen] = useState(false);
    const [isRefreshing, setIsRefreshing] = useState(false);
    const [loggedOut, setLoggedOut] = useState(false);


    useEffect(() => {
        log("isAuthenticated changed:", isAuthenticated);
    }, [isAuthenticated]);

    const fetchData = useCallback(async () => {
        log("fetchData called");
        setIsLoading(true);
        const data = await getAppData();
        if (data) {
            setIsAuthenticated(data.LoggedIn);
            setAppData(data);
        }
        setIsLoading(false);
    }, []);

    const loginRefresh = useCallback(async () => {
        log("loginRefresh called");
        setIsLoading(true);
        const data = await getAppDataNoCache();
        setIsLoading(false);
        if (!data) {
            return false;
        }
        setIsAuthenticated(data.LoggedIn);
        setAppData(data);
        return true;
    }, []);

    const silentRefreshData = useCallback(async () => {
        log("silentRefreshData called");
        if (!isAuthenticated || loggedOut) return;
        setIsRefreshing(true);
        const data = await getAppDataNoCache();
        setIsRefreshing(false);
        if (data) {
            setIsAuthenticated(data.LoggedIn);
            setAppData(data);
        }
    }, [isAuthenticated, loggedOut]);

    useEffect(() => {
        log("loggedOut changed to:", loggedOut);
        trace();
    }, [loggedOut]);

    const loginCallbackFn = useLoginCallback(isAuthenticated, loggedOut, loginRefresh, setLoggedOut, setIsAuthenticated);

    const logInCallBack = (state) => {
        loginCallbackFn(state);
    };

    const {
        handleLogout,
        handleToggleAccountStatus,
        handleUpdateCharacter,
        handleRemoveCharacter,
        handleUpdateAccountName,
        handleRemoveAccount,
        handleAddCharacter,
        handleSaveSkillPlan,
        handleDeleteSkillPlan,
        handleCopySkillPlan,
    } = useAppHandlers({
        setAppData,
        fetchData,
        setIsAuthenticated,
        setLoggedOut,
        setIsSkillPlanModalOpen,
        isAuthenticated,
        loggedOut,
        loginRefresh,
    });

    const openSkillPlanModal = () => setIsSkillPlanModalOpen(true);
    const closeSkillPlanModal = () => setIsSkillPlanModalOpen(false);

    useEffect(() => {
        log("App mounted, calling fetchData");
        fetchData();
    }, [fetchData]);

    useEffect(() => {
        log(`useEffect [isLoading, isAuthenticated]: isLoading=${isLoading}, isAuthenticated=${isAuthenticated}`);
        if (!isLoading && isAuthenticated) {
            silentRefreshData();
        }
    }, [isLoading, isAuthenticated, silentRefreshData]);

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

    const accounts = appData?.AccountData?.Accounts || [];
    const characters = accounts.flatMap((account) => account.Characters) || [];
    const existingAccounts = accounts.map((account) => account.Name) || [];

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
                            <AppRoutes
                                isAuthenticated={isAuthenticated}
                                loggedOut={loggedOut}
                                appData={appData}
                                handleToggleAccountStatus={handleToggleAccountStatus}
                                handleUpdateCharacter={handleUpdateCharacter}
                                handleUpdateAccountName={handleUpdateAccountName}
                                handleRemoveCharacter={handleRemoveCharacter}
                                handleRemoveAccount={handleRemoveAccount}
                                handleDeleteSkillPlan={handleDeleteSkillPlan}
                                handleCopySkillPlan={handleCopySkillPlan}
                                silentRefreshData={silentRefreshData}
                                setAppData={setAppData}
                                characters={characters}
                                logInCallBack={logInCallBack}
                            />
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
