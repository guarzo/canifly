import { useCallback } from 'react';
import { toast } from 'react-toastify';
import { apiRequest } from '../utils/apiRequest';
import { log, trace } from '../utils/logger';
import { isDev, backEndURL } from '../Config';

/**
 * Custom hook that encapsulates all the handler functions used in App.
 * @param {Function} setAppData - Setter for the appData state.
 * @param {Function} fetchData - Function to re-fetch data after certain operations.
 * @param {Function} setIsAuthenticated - Setter for isAuthenticated state.
 * @param {Function} setLoggedOut - Setter for loggedOut state.
 * @param {Function} setIsSkillPlanModalOpen - Setter for isSkillPlanModalOpen state.
 * @param {Function} wrappedFetchAppEndpoint - Ref to fetch endpoints with state updates.
 * @param {boolean} isAuthenticated - Current authentication state.
 * @param {boolean} loggedOut - Current loggedOut state.
 * @returns {Object} Handlers object.
 */
export function useAppHandlers({
                                   setAppData,
                                   fetchData,
                                   setIsAuthenticated,
                                   setLoggedOut,
                                   setIsSkillPlanModalOpen,
                                   wrappedFetchAppEndpoint,
                                   isAuthenticated,
                                   loggedOut
                               }) {

    const handleLogout = useCallback(async () => {
        log("handleLogout called");
        trace();
        await apiRequest(`${backEndURL}/api/logout`, {
            method: 'POST',
            credentials: 'include'
        }, {
            successMessage: 'Logged out successfully!',
            errorMessage: 'Failed to log out.',
            onSuccess: () => {
                setIsAuthenticated(false);
                setAppData(null);
                setLoggedOut(true);
            }
        });
    }, [setIsAuthenticated, setAppData, setLoggedOut]);

    const handleToggleAccountStatus = useCallback(async (accountID) => {
        log("handleToggleAccountStatus called:", accountID);
        await apiRequest(`${backEndURL}/api/toggle-account-status`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ accountID }),
            credentials: 'include'
        }, {
            successMessage: 'Account status toggled successfully!',
            errorMessage: 'Failed to toggle account status.',
            onSuccess: () => {
                setAppData((prev) => {
                    if (!prev) return prev;
                    const updatedAccounts = prev.Accounts.map((account) =>
                        account.ID === accountID
                            ? { ...account, Status: account.Status === 'Alpha' ? 'Omega' : 'Alpha' }
                            : account
                    );
                    return { ...prev, Accounts: updatedAccounts };
                });
            }
        });
    }, [setAppData]);

    const handleUpdateCharacter = useCallback(async (characterID, updates) => {
        log("handleUpdateCharacter called with characterID:", characterID, "updates:", updates);
        await apiRequest(`${backEndURL}/api/update-character`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ characterID, updates }),
            credentials: 'include'
        }, {
            successMessage: 'Character updated successfully!',
            errorMessage: 'Failed to update character.',
            onSuccess: () => {
                setAppData((prev) => {
                    if (!prev) return prev;
                    const updatedRoles = [...prev.Roles];
                    if (updates.Role && !updatedRoles.includes(updates.Role)) {
                        updatedRoles.push(updates.Role);
                    }

                    const updatedAccounts = prev.Accounts.map((account) => {
                        const updatedCharacters = account.Characters.map((character) =>
                            character.Character.CharacterID === characterID
                                ? { ...character, ...updates }
                                : character
                        );
                        return { ...account, Characters: updatedCharacters };
                    });

                    return { ...prev, Accounts: updatedAccounts, Roles: updatedRoles };
                });
            }
        });
    }, [setAppData]);

    const handleRemoveCharacter = useCallback(async (characterID) => {
        log("handleRemoveCharacter called with characterID:", characterID);
        await apiRequest(`${backEndURL}/api/remove-character`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ characterID }),
            credentials: 'include'
        }, {
            successMessage: 'Character removed successfully!',
            errorMessage: 'Failed to remove character.',
            onSuccess: () => {
                setAppData((prev) => {
                    if (!prev) return prev;
                    const updatedAccounts = prev.Accounts.map((account) => {
                        const filteredCharacters = account.Characters.filter(
                            (c) => c.Character.CharacterID !== characterID
                        );
                        return { ...account, Characters: filteredCharacters };
                    });
                    return { ...prev, Accounts: updatedAccounts };
                });
            }
        });
    }, [setAppData]);

    const handleUpdateAccountName = useCallback(async (accountID, newName) => {
        log("handleUpdateAccountName:", { accountID, newName });
        await apiRequest(`${backEndURL}/api/update-account-name`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ accountID, accountName: newName }),
            credentials: 'include'
        }, {
            successMessage: 'Account name updated successfully!',
            errorMessage: 'Failed to update account name.',
            onSuccess: () => {
                setAppData((prev) => {
                    if (!prev) return prev;
                    const updatedAccounts = prev.Accounts.map((account) =>
                        account.ID === accountID ? { ...account, Name: newName } : account
                    );
                    return { ...prev, Accounts: updatedAccounts };
                });
            }
        });
    }, [setAppData]);

    const handleRemoveAccount = useCallback(async (accountName) => {
        log("handleRemoveAccount called with accountName:", accountName);
        await apiRequest(`${backEndURL}/api/remove-account`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ accountName }),
            credentials: 'include'
        }, {
            successMessage: 'Account removed successfully!',
            errorMessage: 'Failed to remove account.',
            onSuccess: () => {
                setAppData((prev) => {
                    if (!prev) return prev;
                    const updatedAccounts = prev.Accounts.filter(
                        (account) => account.Name !== accountName
                    );
                    return { ...prev, Accounts: updatedAccounts };
                });
            }
        });
    }, [setAppData]);

    const handleAddCharacter = useCallback(async (account) => {
        log("handleAddCharacter called with account:", account);
        await apiRequest(`${backEndURL}/api/add-character`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ account }),
            credentials: 'include'
        }, {
            errorMessage: 'An error occurred while adding character.',
            onSuccess: (data) => {
                if (data.redirectURL) {
                    if (isDev) {
                        log("Dev mode: redirecting internally to:", data.redirectURL);
                        window.location.href = data.redirectURL;
                    } else {
                        log("Production mode: opening external browser to:", data.redirectURL);
                        window.electronAPI.openExternal(data.redirectURL);
                        toast.info("Please authenticate in your browser");
                    }
                } else {
                    toast.error("No redirect URL received from server.");
                }
            }
        });
    }, []);

    const handleSaveSkillPlan = useCallback(async (planName, planContents) => {
        log("handleSaveSkillPlan called with planName:", planName);
        await apiRequest(`${backEndURL}/api/save-skill-plan`, {
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
