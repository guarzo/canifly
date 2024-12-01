import React, { useState } from 'react';
import PropTypes from 'prop-types';

/**
 * MainContent Component
 *
 * Displays accounts and their associated characters, as well as unassigned characters.
 */
const MainContent = ({
                         accounts,
                         unassignedCharacters,
                         onAssignCharacter,
                         onUpdateCharacter,
                         onToggleAccountStatus,
                     }) => {
    const [editingAccountId, setEditingAccountId] = useState(null);
    const [accountNameEdits, setAccountNameEdits] = useState({});

    // Handle inline editing for account names
    const handleEditAccountName = (accountId, newName) => {
        setAccountNameEdits({ ...accountNameEdits, [accountId]: newName });
    };

    const handleSaveAccountName = (accountId) => {
        const newName = accountNameEdits[accountId];
        if (newName) {
            console.log('Saving account name:', newName, accountId);
            setEditingAccountId(null);
            // Trigger update logic here if needed
        }
    };

    return (
        <div className="space-y-8">
            {/* Accounts Section */}
            <div>
                <h2 className="text-2xl font-bold text-teal-300 mb-4">Accounts</h2>
                <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
                    {accounts.map((account) => (
                        <div key={account.ID} className="bg-gray-800 p-4 rounded-lg shadow-lg">
                            {/* Account Header */}
                            <div className="flex items-center justify-between">
                                {editingAccountId === account.ID ? (
                                    <input
                                        type="text"
                                        value={accountNameEdits[account.ID] || account.Name}
                                        onChange={(e) =>
                                            handleEditAccountName(account.ID, e.target.value)
                                        }
                                        onBlur={() => handleSaveAccountName(account.ID)}
                                        className="bg-gray-700 text-white px-2 py-1 rounded"
                                    />
                                ) : (
                                    <h3
                                        className="text-lg font-semibold text-teal-200 cursor-pointer"
                                        onClick={() => setEditingAccountId(account.ID)}
                                    >
                                        {account.Name}
                                    </h3>
                                )}
                                <button
                                    className="text-sm text-teal-400 hover:text-teal-600"
                                    onClick={() => onToggleAccountStatus(account.ID)}
                                >
                                    {account.Status === 'Alpha' ? 'Ω' : 'Α'}
                                </button>
                            </div>

                            {/* Account Characters */}
                            <div className="mt-4 space-y-2">
                                {account.Characters.map((character) => (
                                    <div
                                        key={character.Character.CharacterID}
                                        className="bg-gray-700 p-2 rounded flex justify-between items-center"
                                    >
                                        <div>
                                            <p className="text-teal-200 font-medium">
                                                {character.Character.CharacterName}
                                            </p>
                                            <p className="text-gray-400 text-sm">
                                                Role: {character.Role || 'Unassigned'} |{' '}
                                                {character.MCT ? 'MCT Active' : 'MCT Inactive'}
                                            </p>
                                        </div>
                                        <button
                                            className="text-sm text-green-400 hover:text-green-600"
                                            onClick={() =>
                                                onUpdateCharacter(character.Character.CharacterID, {
                                                    role: character.Role || 'New Role',
                                                    MCT: !character.MCT,
                                                })
                                            }
                                        >
                                            Update
                                        </button>
                                    </div>
                                ))}
                            </div>
                        </div>
                    ))}
                </div>
            </div>

            {/* Unassigned Characters Section */}
            {unassignedCharacters?.length > 0 && (
                <div>
                    <h2 className="text-2xl font-bold text-teal-300 mb-4">Unassigned Characters</h2>
                    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
                        {unassignedCharacters.map((character) => (
                            <div
                                key={character.Character.CharacterID}
                                className="bg-gray-800 p-4 rounded-lg shadow-lg flex justify-between items-center"
                            >
                                <p className="text-teal-200 font-medium">
                                    {character.Character.CharacterName}
                                </p>
                                <button
                                    className="text-sm text-green-400 hover:text-green-600"
                                    onClick={() =>
                                        onAssignCharacter(character.Character.CharacterID)
                                    }
                                >
                                    Assign
                                </button>
                            </div>
                        ))}
                    </div>
                </div>
            )}
        </div>
    );
};

MainContent.propTypes = {
    accounts: PropTypes.arrayOf(
        PropTypes.shape({
            ID: PropTypes.number.isRequired,
            Name: PropTypes.string.isRequired,
            Status: PropTypes.string.isRequired,
            Characters: PropTypes.arrayOf(
                PropTypes.shape({
                    Character: PropTypes.shape({
                        CharacterID: PropTypes.number.isRequired,
                        CharacterName: PropTypes.string.isRequired,
                    }).isRequired,
                    Role: PropTypes.string,
                    MCT: PropTypes.bool,
                })
            ),
        })
    ).isRequired,
    unassignedCharacters: PropTypes.arrayOf(
        PropTypes.shape({
            Character: PropTypes.shape({
                CharacterID: PropTypes.number.isRequired,
                CharacterName: PropTypes.string.isRequired,
            }).isRequired,
        })
    ),
    onAssignCharacter: PropTypes.func.isRequired,
    onUpdateCharacter: PropTypes.func.isRequired,
    onToggleAccountStatus: PropTypes.func.isRequired,
};

export default MainContent;
