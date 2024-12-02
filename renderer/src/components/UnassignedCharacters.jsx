// UnassignedCharacters.jsx
import React from 'react';
import PropTypes from 'prop-types';
import { Select, MenuItem } from '@mui/material';

const UnassignedCharacters = ({
                                  characters,
                                  accounts,
                                  onAssignCharacter,
                              }) => {
    const handleAssign = (characterID, accountID) => {
        onAssignCharacter(characterID, accountID === 'new' ? null : accountID);
    };

    return (
        <div className="space-y-4">
            {characters.map((character) => (
                <div
                    key={character.Character.CharacterID}
                    className="p-4 rounded-md shadow-md bg-gray-800 text-teal-200"
                >
          <span className="font-semibold">
            {character.Character.CharacterName}
          </span>
                    <div className="mt-2">
                        <Select
                            defaultValue=""
                            onChange={(e) =>
                                handleAssign(
                                    character.Character.CharacterID,
                                    e.target.value
                                )
                            }
                            displayEmpty
                            className="w-full bg-gray-700 text-teal-200"
                        >
                            <MenuItem value="" disabled>
                                Assign to Account
                            </MenuItem>
                            {accounts.map((account) => (
                                <MenuItem key={account.ID} value={account.ID}>
                                    {account.Name}
                                </MenuItem>
                            ))}
                            <MenuItem value="new">Create New Account</MenuItem>
                        </Select>
                    </div>
                </div>
            ))}
        </div>
    );
};

UnassignedCharacters.propTypes = {
    characters: PropTypes.array.isRequired,
    accounts: PropTypes.array.isRequired,
    onAssignCharacter: PropTypes.func.isRequired,
};

export default UnassignedCharacters;
