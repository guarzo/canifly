/**
 * Normalize appData to ensure certain fields always exist in known formats.
 * @param {Object|null} data - The raw data from the server.
 * @returns {Object|null} Normalized data, or null if data was null.
 */
export function normalizeAppData(data) {
    if (!data) return null;
    return {
        ...data,
        Accounts: Array.isArray(data.Accounts) ? data.Accounts : [],
        Roles: Array.isArray(data.Roles) ? data.Roles : [],
        SkillPlans: data.SkillPlans || {},
        SubDirs: Array.isArray(data.SubDirs) ? data.SubDirs : [],
        associations: Array.isArray(data.associations) ? data.associations : [],
        UserSelections: data.UserSelections || {}
    };
}
