/**
 * Normalize appData to ensure certain fields always exist in known formats under the new model.
 *
 * New structure:
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
 *     EveConversions: {},
 *   }
 * }
 *
 * This function ensures each nested array/object exists, without adding old top-level fields.
 */
import { warn }  from "./logger.jsx"
export function normalizeAppData(data) {
    if (!data) {
        warn("normalizeAppData: Received null or undefined appData from API. Returning null.");
        return null;
    }

    // Validate LoggedIn
    let LoggedIn = data.LoggedIn;
    if (typeof LoggedIn !== 'boolean') {
        warn("normalizeAppData: appData.LoggedIn is not a boolean. Defaulting to false.");
        LoggedIn = false;
    }

    // Validate AccountData
    let AccountData = data.AccountData;
    if (typeof AccountData !== 'object' || AccountData === null) {
        warn("normalizeAppData: appData.AccountData is not an object. Defaulting to empty object.");
        AccountData = {};
    }
    const Accounts = Array.isArray(AccountData.Accounts) ? AccountData.Accounts : (warn("normalizeAppData: appData.AccountData.Accounts is not an array. Defaulting to []."), []);
    const Associations = Array.isArray(AccountData.Associations) ? AccountData.Associations : (warn("normalizeAppData: appData.AccountData.Associations is not an array. Defaulting to []."), []);

    const UserAccount = (typeof AccountData.UserAccount === 'object' && AccountData.UserAccount !== null)
        ? AccountData.UserAccount
        : (() => {
            if (AccountData.UserAccount !== undefined) {
                warn("normalizeAppData: appData.AccountData.UserAccount is not an object. Defaulting to {}.");
            }
            return {};
        })();

    // Validate ConfigData
    let ConfigData = data.ConfigData;
    if (typeof ConfigData !== 'object' || ConfigData === null) {
        warn("normalizeAppData: appData.ConfigData is not an object. Defaulting to empty object.");
        ConfigData = {};
    }
    const Roles = Array.isArray(ConfigData.Roles) ? ConfigData.Roles : (warn("normalizeAppData: appData.ConfigData.Roles is not an array. Defaulting to []."), []);
    const DropDownSelections = (typeof ConfigData.DropDownSelections === 'object' && ConfigData.DropDownSelections !== null)
        ? ConfigData.DropDownSelections
        : (warn("normalizeAppData: appData.ConfigData.DropDownSelections is not an object. Defaulting to {}."), {});

    const SettingsDir = typeof ConfigData.SettingsDir === 'string' ? ConfigData.SettingsDir : (warn("normalizeAppData: appData.ConfigData.SettingsDir is not a string. Defaulting to ''."), '');
    const LastBackupDir = typeof ConfigData.LastBackupDir === 'string' ? ConfigData.LastBackupDir : (warn("normalizeAppData: appData.ConfigData.LastBackupDir is not a string. Defaulting to ''."), '');

    // Validate EveData
    let EveData = data.EveData;
    if (typeof EveData !== 'object' || EveData === null) {
        warn("normalizeAppData: appData.EveData is not an object. Defaulting to empty object.");
        EveData = {};
    }
    const SkillPlans = (typeof EveData.SkillPlans === 'object' && EveData.SkillPlans !== null) ? EveData.SkillPlans : (warn("normalizeAppData: appData.EveData.SkillPlans is not an object. Defaulting to {}."), {});
    const EveProfiles = Array.isArray(EveData.EveProfiles) ? EveData.EveProfiles : (warn("normalizeAppData: appData.EveData.EveProfiles is not an array. Defaulting to []."), []);
    const EveConversions = (typeof EveData.EveConversions === 'object' && EveData.EveConversions !== null) ? EveData.EveConversions : (warn("normalizeAppData: appData.EveData.EveConversions is not an object. Defaulting to {}."), {});

    return {
        LoggedIn,
        AccountData: {
            Accounts,
            Associations,
            UserAccount
        },
        ConfigData: {
            Roles,
            DropDownSelections,
            SettingsDir,
            LastBackupDir
        },
        EveData: {
            SkillPlans,
            EveProfiles,
            EveConversions,
        }
    };
}
