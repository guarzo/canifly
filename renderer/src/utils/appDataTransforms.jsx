// src/utils/appDataTransforms.jsx

import { fetchAppEndpoint } from './api';
import { normalizeAppData } from './dataNormalizer';

/**
 * Utility functions for transforming appData in a consistent, testable manner.
 */

export function updateAccountNameInAppData(prev, accountID, newName) {
    if (!prev) return prev;
    const updatedAccounts = prev.Accounts.map((account) =>
        account.ID === accountID ? { ...account, Name: newName } : account
    );
    return { ...prev, Accounts: updatedAccounts };
}

export function removeCharacterFromAppData(prev, characterID) {
    if (!prev) return prev;
    const updatedAccounts = prev.Accounts.map((account) => {
        const filteredCharacters = account.Characters.filter(
            (c) => c.Character.CharacterID !== characterID
        );
        return { ...account, Characters: filteredCharacters };
    });
    return { ...prev, Accounts: updatedAccounts };
}

export function updateCharacterInAppData(prev, characterID, updates) {
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
}

export function toggleAccountStatusInAppData(prev, accountID) {
    if (!prev) return prev;
    const updatedAccounts = prev.Accounts.map((account) =>
        account.ID === accountID
            ? { ...account, Status: account.Status === 'Alpha' ? 'Omega' : 'Alpha' }
            : account
    );
    return { ...prev, Accounts: updatedAccounts };
}

export function removeAccountFromAppData(prev, accountName) {
    if (!prev) return prev;
    const updatedAccounts = prev.Accounts.filter(
        (account) => account.Name !== accountName
    );
    return { ...prev, Accounts: updatedAccounts };
}


// src/utils/fetchAndNormalizeAppData.js


/**
 * A helper to fetch app data and automatically normalize it before updating state.
 *
 * @param {Object} params - Parameters for fetchAppEndpoint.
 * @param {Object} options - Options for fetchAppEndpoint.
 * @returns {Promise<boolean|undefined>} - Same return as fetchAppEndpoint.
 */
export async function fetchAndNormalizeAppData(params, options) {
    return fetchAppEndpoint(
        {
            ...params,
            // Override setAppData to include normalization step
            setAppData: (data) => {
                if (params.setAppData) {
                    params.setAppData(normalizeAppData(data));
                }
            },
        },
        options
    );
}
