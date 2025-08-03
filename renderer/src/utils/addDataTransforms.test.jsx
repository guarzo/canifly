import { describe, it, expect } from 'vitest';
import {
    updateAccountNameInAppData,
    removeCharacterFromAppData,
    updateCharacterInAppData,
    toggleAccountStatusInAppData,
    removeAccountFromAppData
} from './appDataTransforms';

describe('appDataTransforms', () => {
    const initialAppData = {
        LoggedIn: true,
        AccountData: {
            Accounts: [
                {
                    ID: 1,
                    Name: 'AccountOne',
                    Status: 'Alpha',
                    Characters: [
                        { Character: { CharacterID: 100, CharacterName: 'Char100' }, Role: 'Miner' },
                        { Character: { CharacterID: 101, CharacterName: 'Char101' }, Role: 'Pvp' }
                    ]
                },
                {
                    ID: 2,
                    Name: 'AccountTwo',
                    Status: 'Omega',
                    Characters: [
                        { Character: { CharacterID: 200, CharacterName: 'Char200' }, Role: 'Hauler' }
                    ]
                }
            ],
            Associations: []
        },
        ConfigData: {
            Roles: ['Miner', 'Pvp'],
            SettingsDir: '/some/dir',
            LastBackupDir: '',
            DropDownSelections: {}
        },
        EveData: {
            SkillPlans: {},
            EveProfiles: []
        }
    };

    describe('updateAccountNameInAppData', () => {
        it('updates the name of the matching account', () => {
            const result = updateAccountNameInAppData(initialAppData, 1, 'NewAccountName');
            const updated = result.AccountData.Accounts.find(a => a.ID === 1);
            expect(updated.Name).toBe('NewAccountName');
        });

        it('returns prev if no match', () => {
            const result = updateAccountNameInAppData(initialAppData, 999, 'NoChange');
            expect(result).toEqual(initialAppData); // no changes
        });

        it('returns prev if prev is null', () => {
            expect(updateAccountNameInAppData(null, 1, 'Name')).toBeNull();
        });
    });

    describe('removeCharacterFromAppData', () => {
        it('removes a character by ID', () => {
            const result = removeCharacterFromAppData(initialAppData, 101);
            const account1 = result.AccountData.Accounts.find(a => a.ID === 1);
            expect(account1.Characters.find(c => c.Character.CharacterID === 101)).toBeUndefined();
        });

        it('no changes if character ID not found', () => {
            const result = removeCharacterFromAppData(initialAppData, 999);
            expect(result).toEqual(initialAppData);
        });

        it('returns prev if prev is null', () => {
            expect(removeCharacterFromAppData(null, 100)).toBeNull();
        });
    });

    describe('updateCharacterInAppData', () => {
        it('updates character fields and adds new role if needed', () => {
            const result = updateCharacterInAppData(initialAppData, 100, { Role: 'Logistics', Extra: 'Data' });
            const account1 = result.AccountData.Accounts.find(a => a.ID === 1);
            const char = account1.Characters.find(c => c.Character.CharacterID === 100);
            expect(char.Role).toBe('Logistics');
            expect(char.Extra).toBe('Data');
            expect(result.ConfigData.Roles).toContain('Logistics');
        });

        it('does not duplicate roles if already exists', () => {
            const result = updateCharacterInAppData(initialAppData, 100, { Role: 'Pvp' });
            expect(result.ConfigData.Roles.filter(r => r === 'Pvp').length).toBe(1);
        });

        it('returns prev if no character found', () => {
            const result = updateCharacterInAppData(initialAppData, 999, { Role: 'Whatever' });
            expect(result).toEqual(initialAppData);
        });

        it('returns prev if prev is null', () => {
            expect(updateCharacterInAppData(null, 100, {})).toBeNull();
        });
    });

    describe('toggleAccountStatusInAppData', () => {
        it('toggles Alpha to Omega', () => {
            const result = toggleAccountStatusInAppData(initialAppData, 1);
            const account1 = result.AccountData.Accounts.find(a => a.ID === 1);
            expect(account1.Status).toBe('Omega');
        });

        it('toggles Omega to Alpha', () => {
            const result = toggleAccountStatusInAppData(initialAppData, 2);
            const account2 = result.AccountData.Accounts.find(a => a.ID === 2);
            expect(account2.Status).toBe('Alpha');
        });

        it('no change if account not found', () => {
            const result = toggleAccountStatusInAppData(initialAppData, 999);
            expect(result).toEqual(initialAppData);
        });

        it('returns prev if prev is null', () => {
            expect(toggleAccountStatusInAppData(null, 1)).toBeNull();
        });
    });

    describe('removeAccountFromAppData', () => {
        it('removes account by name', () => {
            const result = removeAccountFromAppData(initialAppData, 'AccountTwo');
            expect(result.AccountData.Accounts.find(a => a.Name === 'AccountTwo')).toBeUndefined();
        });

        it('no change if account name not found', () => {
            const result = removeAccountFromAppData(initialAppData, 'NonExistingAccount');
            expect(result).toEqual(initialAppData);
        });

        it('returns prev if prev is null', () => {
            expect(removeAccountFromAppData(null, 'AccountOne')).toBeNull();
        });
    });
});
