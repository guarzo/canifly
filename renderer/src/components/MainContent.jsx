import React from "react";
import PropTypes from "prop-types";
import { Select, MenuItem, Tooltip } from "@mui/material";

const MainContent = ({
                         accounts,
                         unassignedCharacters,
                         onAssignCharacter,
                         onToggleAccountStatus,
                         onUpdateCharacter,
                         onUpdateAccountName,
                     }) => {
    return (
        <main className="container mx-auto px-4 py-6 space-y-6 mt-16">
            {/* Accounts Section */}
            <div className="space-y-6">
                {accounts.map((account) => (
                    <div
                        key={account.ID}
                        className="p-4 rounded-md shadow-lg bg-gradient-to-r from-gray-900 to-gray-800 text-teal-200"
                    >
                        {/* Account Header */}
                        <div className="flex justify-between items-center mb-4">
                            <input
                                className="bg-transparent border-b border-teal-400 text-lg font-bold focus:outline-none"
                                defaultValue={account.Name}
                                onBlur={(e) => onUpdateAccountName(account.ID, e.target.value)}
                            />
                            <Tooltip title="Toggle Account Status">
                                <button
                                    onClick={() => onToggleAccountStatus(account.ID)}
                                    className="text-2xl font-bold text-white"
                                >
                                    {account.Status === "Alpha" ? "α" : "Ω"}
                                </button>
                            </Tooltip>
                        </div>

                        {/* Characters in Account */}
                        <div className="space-y-4">
                            {account.Characters.map((character, index) => (
                                <div
                                    key={`${character.Character.CharacterID}-${index}`}
                                    className="p-3 rounded-md shadow-md bg-gray-800"
                                >
                                    <div className="flex justify-between items-center">
                                        <span className="font-semibold">{character.Character.CharacterName}</span>
                                        <button
                                            onClick={() =>
                                                onUpdateCharacter(character.Character.CharacterID, {
                                                    MCT: !character.MCT,
                                                })
                                            }
                                            className={`cursor-pointer ${
                                                character.MCT ? "text-green-400" : "text-red-400"
                                            }`}
                                        >
                                            MCT: {character.MCT ? "Enabled" : "Disabled"}
                                        </button>
                                    </div>
                                    <div className="mt-2">Role: {character.Role || "None"}</div>
                                </div>
                            ))}
                        </div>
                    </div>
                ))}
            </div>

            {/* Unassigned Characters Section */}
            {unassignedCharacters && unassignedCharacters.length > 0 && (
                <div>
                    <div className="border-b border-teal-500 my-4" />
                    <div className="space-y-4">
                        {unassignedCharacters.map((character) => (
                            <div
                                key={character.Character.CharacterID}
                                className="p-4 rounded-md shadow-md bg-gray-800 text-teal-200"
                            >
                                <span className="font-semibold">{character.Character.CharacterName}</span>
                                <Select
                                    defaultValue=""
                                    onChange={(e) =>
                                        onAssignCharacter(character.Character.CharacterID, e.target.value)
                                    }
                                    displayEmpty
                                    className="mt-2 w-full"
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
                </div>
            )}
        </main>
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
                    MCT: PropTypes.bool.isRequired,
                    Role: PropTypes.string,
                })
            ).isRequired,
        })
    ).isRequired,
    unassignedCharacters: PropTypes.arrayOf(
        PropTypes.shape({
            Character: PropTypes.shape({
                CharacterID: PropTypes.number.isRequired,
                CharacterName: PropTypes.string.isRequired,
            }).isRequired,
        })
    ).isRequired,
    onAssignCharacter: PropTypes.func.isRequired,
    onToggleAccountStatus: PropTypes.func.isRequired,
    onUpdateCharacter: PropTypes.func.isRequired,
    onUpdateAccountName: PropTypes.func.isRequired,
};

export default MainContent;
