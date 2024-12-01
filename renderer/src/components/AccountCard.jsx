import React from 'react';
import PropTypes from 'prop-types';

const AccountCard = ({
                         account,
                         onToggleAccountStatus,
                         onUpdateAccountName,
                         onUpdateCharacter,
                     }) => {
    return (
        <div className="p-4 rounded-md shadow-lg bg-gray-800 text-teal-200">
            {/* Account Header */}
            <div className="flex justify-between items-center mb-4">
                <input
                    className="bg-transparent border-b border-teal-400 text-lg font-bold focus:outline-none"
                    defaultValue={account.Name}
                    onBlur={(e) => onUpdateAccountName(account.ID, e.target.value)}
                />
                <button
                    onClick={() => onToggleAccountStatus(account.ID)}
                    className="text-2xl font-bold text-white"
                >
                    {account.Status === 'Alpha' ? 'α' : 'Ω'}
                </button>
            </div>

            {/* Characters in Account */}
            <div className="space-y-4">
                {account.Characters.map((character) => (
                    <div
                        key={character.Character.CharacterID}
                        className="p-3 rounded-md shadow-md bg-gray-700"
                    >
                        <div className="flex justify-between items-center">
              <span className="font-semibold">
                {character.Character.CharacterName}
              </span>
                            <button
                                onClick={() =>
                                    onUpdateCharacter(character.Character.CharacterID, {
                                        MCT: !character.MCT,
                                    })
                                }
                                className={`cursor-pointer ${
                                    character.MCT ? 'text-green-400' : 'text-red-400'
                                }`}
                            >
                                MCT: {character.MCT ? 'Enabled' : 'Disabled'}
                            </button>
                        </div>
                        <div className="mt-2">Role: {character.Role || 'None'}</div>
                    </div>
                ))}
            </div>
        </div>
    );
};

AccountCard.propTypes = {
    account: PropTypes.object.isRequired,
    onToggleAccountStatus: PropTypes.func.isRequired,
    onUpdateAccountName: PropTypes.func.isRequired,
    onUpdateCharacter: PropTypes.func.isRequired,
};

export default AccountCard;
