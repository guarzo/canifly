// src/components/dashboard/GroupCard.jsx
import React from 'react';
import PropTypes from 'prop-types';
import CharacterItem from './CharacterItem.jsx';

const GroupCard = ({ groupName, characters, onUpdateCharacter, roles, skillConversions }) => {
    return (
        <div className="p-4 rounded-md shadow-md bg-gray-800 text-teal-200 max-w-sm">
            {/* Group Header */}
            <div className="flex justify-between items-center mb-4">
                <span className="text-sm font-bold">
                    {groupName}
                </span>
            </div>

            {/* Characters in this group */}
            <div className="space-y-2">
                {characters.map((character) => (
                    <CharacterItem
                        key={character.Character.CharacterID}
                        character={character}
                        onUpdateCharacter={onUpdateCharacter}
                        roles={roles}
                        hideRemoveIcon={true}
                        skillConversions={skillConversions}
                    />
                ))}
            </div>
        </div>
    );
};

GroupCard.propTypes = {
    groupName: PropTypes.string.isRequired,
    characters: PropTypes.array.isRequired,
    onUpdateCharacter: PropTypes.func.isRequired,
    roles: PropTypes.array.isRequired,
    skillConversions: PropTypes.object.isRequired,
};

export default GroupCard;
