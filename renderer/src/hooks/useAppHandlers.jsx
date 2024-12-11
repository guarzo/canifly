// src/hooks/useAppHandlers.jsx

import { useCallback } from 'react';
import { log } from '../utils/logger';
import { createPostHandler, createAddCharacterHandler } from '../utils/createRequestHandler';
import {
    removeCharacterFromAppData,
    updateCharacterInAppData,
    toggleAccountStatusInAppData,
    updateAccountNameInAppData,
    removeAccountFromAppData
} from '../utils/appDataTransforms';
import { apiRequest } from '../utils/apiRequest';

/**
 * Custom hook that encapsulates all the handler functions used in App.
 * @param {Function} setAppData - Setter for the appData state.
 * @param {Function} fetchData - Function to re-fetch data after certain operations.
 * @param {Function} setIsAuthenticated - Setter for isAuthenticated state.
 * @param {Function} setLoggedOut - Setter for loggedOut state.
 * @param {Function} setIsSkillPlanModalOpen - Setter for isSkillPlanModalOpen state.
 * @returns {Object} Handlers object.
 */
export function useAppHandlers({
                                   setAppData,
                                   fetchData,
                                   setIsAuthenticated,
                                   setLoggedOut,
                                   setIsSkillPlanModalOpen,
                               }) {

    const handleLogout = useCallback(async () => {
        console.trace("handleLogout called");
        console.trace("handleLogout call stack:");
        await apiRequest(`/api/logout`, {
            method: 'POST',
            credentials: 'include'
        }, {
            successMessage: 'Logged out successfully!',
            errorMessage: 'Failed to log out.',
            onSuccess: () => {
                setIsAuthenticated(false);
                setAppData(null);
                console.log("handleLogout success");
                console.trace("handleLogout call stack:");
                setLoggedOut(true);
            }
        });
    }, [setIsAuthenticated, setAppData, setLoggedOut]);

    const handleToggleAccountStatus = useCallback(async (accountID) => {
        log("handleToggleAccountStatus called:", accountID);
        const toggleAccountStatus = createPostHandler({
            endpoint: '/api/toggle-account-status',
            setAppData,
            successMessage: 'Account status toggled successfully!',
            errorMessage: 'Failed to toggle account status.',
            onSuccess: () => {
                setAppData((prev) => toggleAccountStatusInAppData(prev, accountID));
            }
        });
        await toggleAccountStatus({ accountID });
    }, [setAppData]);

    const handleUpdateCharacter = useCallback(async (characterID, updates) => {
        log("handleUpdateCharacter called with characterID:", characterID, "updates:", updates);
        const updateCharacterReq = createPostHandler({
            endpoint: '/api/update-character',
            setAppData,
            successMessage: 'Character updated successfully!',
            errorMessage: 'Failed to update character.',
            onSuccess: () => {
                setAppData((prev) => updateCharacterInAppData(prev, characterID, updates));
            }
        });
        await updateCharacterReq({ characterID, updates });
    }, [setAppData]);

    const handleRemoveCharacter = useCallback(async (characterID) => {
        log("handleRemoveCharacter called with characterID:", characterID);
        const removeCharacterReq = createPostHandler({
            endpoint: '/api/remove-character',
            setAppData,
            successMessage: 'Character removed successfully!',
            errorMessage: 'Failed to remove character.',
            onSuccess: () => {
                setAppData((prev) => removeCharacterFromAppData(prev, characterID));
            }
        });
        await removeCharacterReq({ characterID });
    }, [setAppData]);

    const handleUpdateAccountName = useCallback(async (accountID, newName) => {
        log("handleUpdateAccountName:", { accountID, newName });
        const updateAccountNameReq = createPostHandler({
            endpoint: '/api/update-account-name',
            setAppData,
            successMessage: 'Account name updated successfully!',
            errorMessage: 'Failed to update account name.',
            onSuccess: () => {
                setAppData((prev) => updateAccountNameInAppData(prev, accountID, newName));
            }
        });
        await updateAccountNameReq({ accountID, accountName: newName });
    }, [setAppData]);

    const handleRemoveAccount = useCallback(async (accountName) => {
        log("handleRemoveAccount called with accountName:", accountName);
        const removeAccountReq = createPostHandler({
            endpoint: '/api/remove-account',
            setAppData,
            successMessage: 'Account removed successfully!',
            errorMessage: 'Failed to remove account.',
            onSuccess: () => {
                setAppData((prev) => removeAccountFromAppData(prev, accountName));
            }
        });
        await removeAccountReq({ accountName });
    }, [setAppData]);

    const handleAddCharacter = useCallback(async (account) => {
        const addCharacter = createAddCharacterHandler();
        await addCharacter(account);
    }, []);

    const handleSaveSkillPlan = useCallback(async (planName, planContents) => {
        log("handleSaveSkillPlan called with planName:", planName);
        // This one is slightly different because it refreshes data after success
        await apiRequest(`/api/save-skill-plan`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ name: planName, contents: planContents }),
            credentials: 'include'
        }, {
            successMessage: 'Skill Plan Saved!',
            errorMessage: 'Failed to save skill plan.',
            onSuccess: () => {
                setIsSkillPlanModalOpen(false);
                fetchData();
            }
        });
    }, [setIsSkillPlanModalOpen, fetchData]);

    return {
        handleLogout,
        handleToggleAccountStatus,
        handleUpdateCharacter,
        handleRemoveCharacter,
        handleUpdateAccountName,
        handleRemoveAccount,
        handleAddCharacter,
        handleSaveSkillPlan
    };
}
