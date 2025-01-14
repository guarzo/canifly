/**
 * Utility functions for transforming appData with the new model structure.
 *
 * Remember the new structure:
 * appData = {
 *   LoggedIn: boolean,
 *   AccountData: {
 *     Accounts: [],
 *     Associations: []
 *   },
 *   ConfigData: {
 *     Roles: [],
 *     SettingsDir: string,
 *     LastBackupDir: string,
 *     DropDownSelections: {}
 *   },
 *   EveData: {
 *     SkillPlans: {},
 *     EveProfiles: []
 *   }
 * }
 */

export function updateAccountNameInAppData(prev, accountID, newName) {
    if (!prev) return prev;

    const updatedAccounts = prev.AccountData.Accounts.map((account) =>
        account.ID === accountID ? { ...account, Name: newName } : account
    );

    return {
        ...prev,
        AccountData: {
            ...prev.AccountData,
            Accounts: updatedAccounts
        }
    };
}

export function removeCharacterFromAppData(prev, characterID) {
    if (!prev) return prev;
    const updatedAccounts = prev.AccountData.Accounts.map((account) => {
        const filteredCharacters = account.Characters.filter(
            (c) => c.Character.CharacterID !== characterID
        );
        return { ...account, Characters: filteredCharacters };
    });

    return {
        ...prev,
        AccountData: {
            ...prev.AccountData,
            Accounts: updatedAccounts
        }
    };
}

export function updateCharacterInAppData(prev, characterID, updates) {
    if (!prev) return prev;

    let characterFound = false;

    const updatedAccounts = prev.AccountData.Accounts.map((account) => {
        const updatedCharacters = account.Characters.map((character) => {
            if (character.Character.CharacterID === characterID) {
                characterFound = true;
                return { ...character, ...updates };
            }
            return character;
        });
        return { ...account, Characters: updatedCharacters };
    });

    // If character not found, return prev as is
    if (!characterFound) {
        return prev;
    }

    const updatedRoles = Array.isArray(prev.ConfigData.Roles) ? [...prev.ConfigData.Roles] : [];
    if (updates.Role && !updatedRoles.includes(updates.Role)) {
        updatedRoles.push(updates.Role);
    }

    return {
        ...prev,
        AccountData: {
            ...prev.AccountData,
            Accounts: updatedAccounts
        },
        ConfigData: {
            ...prev.ConfigData,
            Roles: updatedRoles
        }
    };
}


export function toggleAccountStatusInAppData(prev, accountID) {
    if (!prev) return prev;

    const updatedAccounts = prev.AccountData.Accounts.map((account) =>
        account.ID === accountID
            ? { ...account, Status: account.Status === 'Alpha' ? 'Omega' : 'Alpha' }
            : account
    );

    return {
        ...prev,
        AccountData: {
            ...prev.AccountData,
            Accounts: updatedAccounts
        }
    };
}


export function toggleAccountVisibilityInAppData(prev, accountID) {
    if (!prev) return prev;

    const updatedAccounts = prev.AccountData.Accounts.map((account) =>
        account.ID === accountID
            ? { ...account, Visible: !account.Visible }
            : account
    );

    return {
        ...prev,
        AccountData: {
            ...prev.AccountData,
            Accounts: updatedAccounts
        }
    };
}

export function removeAccountFromAppData(prev, accountName) {
    if (!prev) return prev;
    const updatedAccounts = prev.AccountData.Accounts.filter(
        (account) => account.Name !== accountName
    );

    return {
        ...prev,
        AccountData: {
            ...prev.AccountData,
            Accounts: updatedAccounts
        }
    };
}

export function removePlanFromSkillPlans(prev, planName) {
    if (!prev) return prev;

    console.log("planName: ", planName)
    console.log(prev.EveData.SkillPlans)
    const updatedSkillPlans = { ...prev.EveData.SkillPlans };
    delete updatedSkillPlans[planName];
    console.log(updatedSkillPlans)

    return {
        ...prev,
        EveData: {
            ...prev.EveData,
            SkillPlans: updatedSkillPlans
        }
    };
}