import {useCallback} from 'react';
import { log } from '../utils/logger';
import {
    removeCharacterFromAppData,
    updateCharacterInAppData,
    toggleAccountStatusInAppData,
    updateAccountNameInAppData,
    removeAccountFromAppData, removePlanFromSkillPlans, toggleAccountVisibilityInAppData
} from '../utils/appDataTransforms';
import {isDev} from '../Config';


import {
    logout,
    toggleAccountStatus as toggleAccountStatusApi,
    updateCharacter as updateCharacterApi,
    removeCharacter as removeCharacterApi,
    updateAccountName as updateAccountNameApi,
    removeAccount as removeAccountApi,
    addCharacter as addCharacterApi,
    saveSkillPlan as saveSkillPlanApi,
    deleteSkillPlan as deleteSkillPlanApi, toggleAccountVisibility,
} from '../api/apiService.jsx';
import {toast} from "react-toastify";
import {useLoginCallback} from "./useLoginCallback.jsx";
import {useAddCharacterCallback} from "./useAddCharacterCallback.jsx";

/**
 * Custom hook that encapsulates all the handler functions used in App.
 * @param {Function} setAppData - Setter for the appData state.
 * @param {Function} fetchData - Function to re-fetch data after certain operations.
 * @param {Function} setIsAuthenticated - Setter for isAuthenticated state.
 * @param {Function} setLoggedOut - Setter for loggedOut state.
 * @param {Function} setIsSkillPlanModalOpen - Setter for isSkillPlanModalOpen state.
 * @param isAuthenticated
 * @param loggedOut
 * @param loginRefresh
 * @returns {Object} Handlers object.
 */
export function useAppHandlers({
                                   setAppData,
                                   fetchData,
                                   setIsAuthenticated,
                                   setLoggedOut,
                                   setIsSkillPlanModalOpen,
                                   isAuthenticated,
                                   loggedOut,
                                   loginRefresh,
                               }) {

    const logInCallBack = useLoginCallback(
        isAuthenticated,
        loggedOut,
        loginRefresh,
        setLoggedOut,
        setIsAuthenticated
    );

    const addCharacterCallback = useAddCharacterCallback(fetchData);

    const handleLogout = useCallback(async () => {
        log("handleLogout called");
        const result = await logout();
        if (result && result.success) {
            setIsAuthenticated(false);
            setAppData(null);
            setLoggedOut(true);
        }
    }, [setIsAuthenticated, setAppData, setLoggedOut]);

    const handleToggleAccountStatus = useCallback(async (accountID) => {
        log("handleToggleAccountStatus called:", accountID);
        const result = await toggleAccountStatusApi(accountID);
        if (result && result.success) {
            setAppData((prev) => toggleAccountStatusInAppData(prev, accountID));
        }
    }, [setAppData]);

    const handleToggleAccountVisibility = useCallback(async (accountID) => {
        log("handleToggleAccountVisibility called:", accountID);
        const result = await toggleAccountVisibility(accountID);
        if (result && result.success) {
            setAppData((prev) => toggleAccountVisibilityInAppData(prev, accountID));
        }
    }, [setAppData]);

    const handleUpdateCharacter = useCallback(async (characterID, updates) => {
        log("handleUpdateCharacter called with characterID:", characterID, "updates:", updates);
        const result = await updateCharacterApi(characterID, updates);
        if (result && result.success) {
            setAppData((prev) => updateCharacterInAppData(prev, characterID, updates));
        }
    }, [setAppData]);

    const handleRemoveCharacter = useCallback(async (characterID) => {
        log("handleRemoveCharacter called with characterID:", characterID);
        const result = await removeCharacterApi(characterID);
        if (result && result.success) {
            setAppData((prev) => removeCharacterFromAppData(prev, characterID));
        }
    }, [setAppData]);

    const handleUpdateAccountName = useCallback(async (accountID, newName) => {
        log("handleUpdateAccountName:", { accountID, newName });
        const result = await updateAccountNameApi(accountID, newName);
        if (result && result.success) {
            setAppData((prev) => updateAccountNameInAppData(prev, accountID, newName));
        }
    }, [setAppData]);

    const handleRemoveAccount = useCallback(async (accountName) => {
        log("handleRemoveAccount called with accountName:", accountName);
        const result = await removeAccountApi(accountName);
        if (result && result.success) {
            setAppData((prev) => removeAccountFromAppData(prev, accountName));
        }
    }, [setAppData]);

    // 2) The handleAddCharacter that calls the API, opens external link, and calls the new callback
    const handleAddCharacter = useCallback(async (account) => {
        log("handleAddCharacter called with account:", account);

        const data = await addCharacterApi(account);

        // Open the external URL if present
        if (data?.redirectURL) {
            if (isDev) {
                window.location.href = data.redirectURL;
            } else if (window.electronAPI?.openExternal) {
                window.electronAPI.openExternal(data.redirectURL);
            }
        }

        console.log(data)

        // Start polling to finalize
        if (data?.state) {
            log("Starting addCharacterCallback with state:", data.state);
            addCharacterCallback(data.state);
        }
        // No need to rely on isAuthenticated => user is already logged in
    }, [addCharacterCallback]);


    const handleSaveSkillPlan = useCallback(async (planName, planContents) => {
        log("handleSaveSkillPlan called with planName:", planName);
        const result = await saveSkillPlanApi(planName, planContents);
        if (result && result.success) {
            setIsSkillPlanModalOpen(false);
            fetchData();
        }
    }, [setIsSkillPlanModalOpen, fetchData]);

    const handleDeleteSkillPlan = useCallback(async (planName) => {
        log("handleDeleteSkillPlan called with planName:", planName);
        const result = await deleteSkillPlanApi(planName);
        console.log(result)
        if (result && result.success) {
            toast.success(`Deleted skill plan: ${planName}`, { autoClose: 1500 });
            setAppData((prev) => removePlanFromSkillPlans(prev, planName));
        }
    }, [fetchData]);

    const handleCopySkillPlan = useCallback(async (planName, skills) => {
        log("handleCopySkillPlan called with planName:", planName);
        if (!planName) {
            console.error(`Skill plan not found: ${planName}`);
            toast.warning(`Skill plan not found: ${planName}`, { autoClose: 1500 });
            return;
        }

        if (Object.keys(skills).length === 0) {
            console.warn(`No skills available to copy in the plan: ${planName}`);
            toast.warning(`No skills available to copy in the plan: ${planName}.`, {
                autoClose: 1500,
            });
            return;
        }

        const skillText = Object.entries(skills)
            .map(([skill, detail]) => `${skill} ${detail.Level}`)
            .join('\n');

        navigator.clipboard
            .writeText(skillText)
            .then(() => {
                toast.success(`Copied ${Object.keys(skills).length} skills from ${planName}.`, {
                    autoClose: 1500,
                });
            })
            .catch((err) => {
                console.error('Copy to clipboard failed:', err);
                toast.error('Failed to copy skill plan.', { autoClose: 1500 });
            });
    },[fetchData]);

    return {
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
        handleToggleAccountVisibility,
    };
}
