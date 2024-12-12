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
 *   }
 * }
 *
 * This function ensures each nested array/object exists, without adding old top-level fields.
 */

export function normalizeAppData(data) {
    if (!data) return null;

    return {
        ...data,
        AccountData: {
            ...data.AccountData,
            Accounts: Array.isArray(data?.AccountData?.Accounts) ? data.AccountData.Accounts : [],
            Associations: Array.isArray(data?.AccountData?.Associations) ? data.AccountData.Associations : []
        },
        ConfigData: {
            ...data.ConfigData,
            Roles: Array.isArray(data?.ConfigData?.Roles) ? data.ConfigData.Roles : [],
            DropDownSelections: data?.ConfigData?.DropDownSelections || {},
            SettingsDir: data?.ConfigData?.SettingsDir || '',
            LastBackupDir: data?.ConfigData?.LastBackupDir || ''
        },
        EveData: {
            ...data.EveData,
            SkillPlans: data?.EveData?.SkillPlans || {},
            EveProfiles: Array.isArray(data?.EveData?.EveProfiles) ? data.EveData.EveProfiles : []
        },
        // LoggedIn is top-level and should remain as is, just ensure boolean
        LoggedIn: typeof data.LoggedIn === 'boolean' ? data.LoggedIn : false
    };
}
