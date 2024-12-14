import { describe, it, expect, vi } from 'vitest';
import { normalizeAppData } from './dataNormalizer.jsx';
import { warn } from './logger.jsx';

// Mock the warn function
vi.mock('./logger.jsx', () => ({
    warn: vi.fn(),
}));

describe('normalizeAppData', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it('returns null and warns if data is null or undefined', () => {
        expect(normalizeAppData(null)).toBeNull();
        expect(warn).toHaveBeenCalledWith("normalizeAppData: Received null or undefined appData from API. Returning null.");

        vi.clearAllMocks();

        expect(normalizeAppData(undefined)).toBeNull();
        expect(warn).toHaveBeenCalledWith("normalizeAppData: Received null or undefined appData from API. Returning null.");
    });

    it('defaults LoggedIn to false if not boolean', () => {
        const data = { LoggedIn: 'not boolean' };
        const result = normalizeAppData(data);
        expect(result.LoggedIn).toBe(false);
        expect(warn).toHaveBeenCalledWith("normalizeAppData: appData.LoggedIn is not a boolean. Defaulting to false.");
    });

    it('defaults AccountData to empty object if not object', () => {
        const data = { LoggedIn: true, AccountData: "notAnObject" };
        const result = normalizeAppData(data);
        expect(result.AccountData).toEqual({
            Accounts: [],
            Associations: [],
            UserAccount: {}
        });
        expect(warn).toHaveBeenCalledWith("normalizeAppData: appData.AccountData is not an object. Defaulting to empty object.");
        expect(warn).toHaveBeenCalledWith("normalizeAppData: appData.AccountData.Accounts is not an array. Defaulting to [].");
        expect(warn).toHaveBeenCalledWith("normalizeAppData: appData.AccountData.Associations is not an array. Defaulting to [].");
    });

    it('handles UserAccount defaults', () => {
        const data = {
            LoggedIn: true,
            AccountData: {
                Accounts: [],
                Associations: [],
                UserAccount: "notObject"
            }
        };
        const result = normalizeAppData(data);
        expect(result.AccountData.UserAccount).toEqual({});
        expect(warn).toHaveBeenCalledWith("normalizeAppData: appData.AccountData.UserAccount is not an object. Defaulting to {}.");
    });

    it('defaults ConfigData if not object', () => {
        const data = { LoggedIn: true, AccountData: {}, ConfigData: null };
        const result = normalizeAppData(data);
        expect(result.ConfigData).toEqual({
            Roles: [],
            DropDownSelections: {},
            SettingsDir: '',
            LastBackupDir: ''
        });
        expect(warn).toHaveBeenCalledWith("normalizeAppData: appData.ConfigData is not an object. Defaulting to empty object.");
        expect(warn).toHaveBeenCalledWith("normalizeAppData: appData.ConfigData.Roles is not an array. Defaulting to [].");
        expect(warn).toHaveBeenCalledWith("normalizeAppData: appData.ConfigData.SettingsDir is not a string. Defaulting to ''.");
        expect(warn).toHaveBeenCalledWith("normalizeAppData: appData.ConfigData.LastBackupDir is not a string. Defaulting to ''.");
    });

    it('defaults EveData if not object', () => {
        const data = { LoggedIn: true, AccountData: {}, ConfigData: {}, EveData: "notObject" };
        const result = normalizeAppData(data);
        expect(result.EveData).toEqual({
            SkillPlans: {},
            EveProfiles: [],
            EveConversions: {},
        });
        expect(warn).toHaveBeenCalledWith("normalizeAppData: appData.EveData is not an object. Defaulting to empty object.");
        expect(warn).toHaveBeenCalledWith("normalizeAppData: appData.EveData.SkillPlans is not an object. Defaulting to {}.");
        expect(warn).toHaveBeenCalledWith("normalizeAppData: appData.EveData.EveProfiles is not an array. Defaulting to [].");
        expect(warn).toHaveBeenCalledWith("normalizeAppData: appData.EveData.EveConversions is not an object. Defaulting to {}.");
    });

    it('returns a fully normalized object when given proper data', () => {
        const data = {
            LoggedIn: true,
            AccountData: {
                Accounts: [{ Name: 'TestAccount' }],
                Associations: [{ userId: 'user1', charId: 'char1' }],
                UserAccount: { userId: 'user123' }
            },
            ConfigData: {
                Roles: ['Admin', 'User'],
                DropDownSelections: { key: 'value' },
                SettingsDir: '/path/to/settings',
                LastBackupDir: '/path/to/backup'
            },
            EveData: {
                SkillPlans: { PlanA: { Skills: {} } },
                EveProfiles: [{ profile: 'default' }],
                EveConversions: { Something: 'else'}
            }
        };

        const result = normalizeAppData(data);

        expect(result).toEqual({
            LoggedIn: true,
            AccountData: {
                Accounts: [{ Name: 'TestAccount' }],
                Associations: [{ userId: 'user1', charId: 'char1' }],
                UserAccount: { userId: 'user123' }
            },
            ConfigData: {
                Roles: ['Admin', 'User'],
                DropDownSelections: { key: 'value' },
                SettingsDir: '/path/to/settings',
                LastBackupDir: '/path/to/backup'
            },
            EveData: {
                SkillPlans: { PlanA: { Skills: {} } },
                EveProfiles: [{ profile: 'default' }],
                EveConversions: { Something: 'else'}
            }
        });

        // No warnings should be called in this case
        expect(warn).not.toHaveBeenCalled();
    });
});
