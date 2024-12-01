import React from 'react';
import PropTypes from 'prop-types';
import { Select, MenuItem } from '@mui/material';

const UnassignedCharacters = ({
                                  characters,
                                  accounts,
                                  onAssignCharacter,
                              }) => {
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
                    <Select
                        defaultValue=""
                        onChange={(e) =>
                            onAssignCharacter(
                                character.Character.CharacterID,
                                e.target.value
                            )
                        }
                        displayEmpty
                        className="mt-2 w-full bg-gray-700 text-teal-200"
                    >
                        <MenuItem value="" disabled>
                            Assign to Account
                        </MenuItem>
                        {accounts.map((account) => (
                            <MenuItem key={account.ID} value={account.ID}>
                                {account.Name}
                            </MenuItem>
                        ))}
                    </Select>
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
