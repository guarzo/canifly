// src/App.jsx

import { useState, useEffect, useCallback } from 'react';
import { HashRouter as Router } from 'react-router-dom';
import { ThemeProvider } from '@mui/material/styles';
import { ToastContainer } from 'react-toastify';

import { backEndURL } from './Config';
import { fetchAppEndpoint } from './utils/api';
import { useLoginCallback } from './hooks/useLoginCallback';
import { log, trace } from './utils/logger';
import { normalizeAppData } from './utils/dataNormalizer.jsx';
import { useAppHandlers } from './hooks/useAppHandlers';

import Header from './components/partials/Header.jsx';
import Footer from './components/partials/Footer.jsx';
import AddSkillPlanModal from './components/skillplan/AddSkillPlanModal.jsx';
import ErrorBoundary from './components/partials/ErrorBoundary.jsx';
import AppRoutes from './Routes';

import theme from './theme.jsx';
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

    const wrappedFetchAppEndpoint = useCallback(
        async (endpoint, options) => {
            const result = await fetchAppEndpoint(
                {
                    backEndURL,
                    endpoint,
                    loggedOut,
                    isAuthenticated,
                    setIsLoading,
                    setIsRefreshing,
                    setIsAuthenticated,
                    setAppData: (data) => setAppData(normalizeAppData(data))
                },
                options
            );
            return result;
        },
        [backEndURL, loggedOut, isAuthenticated]
    );

    const fetchData = useCallback(async () => {
        log("fetchData called");
        await wrappedFetchAppEndpoint('/api/app-data', { setLoading: true });
    }, [wrappedFetchAppEndpoint]);

    const loginRefresh = useCallback(() => {
        log("loginRefresh called");
        return wrappedFetchAppEndpoint('/api/app-data-no-cache', { setLoading: true, returnSuccess: true });
    }, [wrappedFetchAppEndpoint]);

    const silentRefreshData = useCallback(async () => {
        log("silentRefreshData called");
        if (!isAuthenticated || loggedOut) return;
        await wrappedFetchAppEndpoint('/api/app-data-no-cache', { setRefreshing: true });
    }, [wrappedFetchAppEndpoint, isAuthenticated, loggedOut]);

    useEffect(() => {
        log("loggedOut changed to:", loggedOut);
        trace();
    }, [loggedOut]);

    const loginCallbackFn = useLoginCallback(isAuthenticated, loggedOut, loginRefresh, setLoggedOut, setIsAuthenticated, backEndURL);

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
        handleSaveSkillPlan
    } = useAppHandlers({
        setAppData,
        fetchData,
        setIsAuthenticated,
        setLoggedOut,
        setIsSkillPlanModalOpen,
        wrappedFetchAppEndpoint,
        isAuthenticated,
        loggedOut
    });

    const openSkillPlanModal = () => {
        setIsSkillPlanModalOpen(true);
    };

    const closeSkillPlanModal = () => {
        setIsSkillPlanModalOpen(false);
    };

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

    const identities = appData?.Accounts.flatMap((account) => account.Characters) || [];
    const existingAccounts = appData?.Accounts.map((account) => account.Name) || [];

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
                                silentRefreshData={silentRefreshData}
                                setAppData={setAppData}
                                identities={identities}
                                backEndURL={backEndURL}
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
