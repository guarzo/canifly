import React from 'react';
import PropTypes from 'prop-types';
import { FaUserPlus } from 'react-icons/fa'; // Importing the assign icon

const MainContent = ({ accounts, unassignedCharacters, onAssignCharacter }) => {
    const characters = unassignedCharacters || [];

    return (
        <main className="container mx-auto px-4 py-8 bg-gradient-to-b from-gray-800 to-gray-700 mt-24">
            {/* Accounts Section */}
            <div className="bg-gray-800 text-gray-100 p-6 rounded-md shadow-md mb-8">
                <h2 className="text-xl font-semibold text-teal-200 mb-4 border-b-2 border-teal-500 pb-2">Accounts</h2>
                {accounts && accounts.length > 0 ? (
                    <ul className="list-disc pl-6">
                        {accounts.map((account) => (
                            <li key={account.id} className="text-md text-gray-300">
                                {account.name}
                            </li>
                        ))}
                    </ul>
                ) : (
                    <p className="text-md text-gray-500">No accounts found.</p>
                )}
            </div>

            {/* Unassigned Characters Section */}
            <div className="bg-gray-800 text-gray-100 p-6 rounded-md shadow-md">
                <h2 className="text-xl font-semibold text-teal-200 mb-4 border-b-2 border-teal-500 pb-2">Unassigned Characters</h2>
                {characters && characters.length > 0 ? (
                    <ul className="list-disc pl-6">
                        {characters.map((character) => (
                            <li
                                key={character.Character.CharacterID}
                                className="text-md text-gray-300 flex items-center justify-between mb-2"
                            >
                                <span className="flex-grow">{character.Character.CharacterName}</span>
                                <button
                                    onClick={() => onAssignCharacter(character.Character.CharacterID)}
                                    className="flex items-center p-2 bg-teal-500 text-white rounded hover:bg-teal-600 focus:outline-none"
                                    title="Assign Character"
                                >
                                    <FaUserPlus />
                                </button>
                            </li>
                        ))}
                    </ul>
                ) : (
                    <p className="text-md text-gray-500">No unassigned characters found.</p>
                )}
            </div>
        </main>
    );
};

MainContent.propTypes = {
    accounts: PropTypes.arrayOf(
        PropTypes.shape({
            id: PropTypes.number.isRequired,
            name: PropTypes.string.isRequired,
        })
    ),
    unassignedCharacters: PropTypes.arrayOf(
        PropTypes.shape({
            Character: PropTypes.shape({
                CharacterID: PropTypes.number.isRequired,
                CharacterName: PropTypes.string.isRequired,
            }).isRequired,
        })
    ),
    onAssignCharacter: PropTypes.func.isRequired,
};

export default MainContent;
