import React from 'react';
import PropTypes from 'prop-types';

const MainContent = ({ accounts, unassignedCharacters, onAssignCharacter }) => {
    const characters = unassignedCharacters || [];

    return (
        <div className="container mx-auto p-6">
            <h1 className="text-3xl font-bold text-center text-gray-800 mb-6">Main Content</h1>

            {/* Display Accounts */}
            <div className="mb-12 bg-white shadow-lg rounded-lg p-6">
                <h2 className="text-2xl font-semibold text-gray-700 mb-4">Accounts</h2>
                {accounts && accounts.length > 0 ? (
                    <ul className="list-disc pl-6">
                        {accounts.map((account) => (
                            <li key={account.id} className="text-md text-gray-700">
                                {account.name}
                            </li>
                        ))}
                    </ul>
                ) : (
                    <p className="text-gray-600">No accounts found.</p>
                )}
            </div>

            {/* Display Unassigned Characters */}
            <div className="bg-white shadow-lg rounded-lg p-6">
                <h2 className="text-2xl font-semibold text-gray-700 mb-4">Unassigned Characters</h2>
                {characters.length > 0 ? (
                    <ul className="list-disc pl-6">
                        {characters.map((character) => (
                            <li key={character.Character.CharacterID} className="flex justify-between items-center py-2">
                                <span className="text-md text-gray-700">{character.Character.CharacterName}</span>
                                <button
                                    onClick={() => onAssignCharacter(character.Character.CharacterID, character.Character.AccountID)}
                                    className="bg-blue-500 text-white px-4 py-2 rounded-md hover:bg-blue-600 transition duration-200"
                                >
                                    Assign
                                </button>
                            </li>
                        ))}
                    </ul>
                ) : (
                    <p className="text-gray-600">No unassigned characters found.</p>
                )}
            </div>
        </div>
    );
};

MainContent.propTypes = {
    accounts: PropTypes.array.isRequired,
    unassignedCharacters: PropTypes.array,
    onAssignCharacter: PropTypes.func.isRequired,  // Ensure this prop is passed
};

export default MainContent;
