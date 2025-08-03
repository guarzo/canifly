// src/components/dashboard/GroupCard.jsx
import React from 'react';
import PropTypes from 'prop-types';
import CharacterItem from './CharacterItem.jsx';
import GlassCard from '../ui/GlassCard.jsx';
import { Badge } from '@mui/material';

const GroupCard = ({ groupName, characters, onUpdateCharacter, roles, skillConversions }) => {
    const getGroupIcon = (name) => {
        if (name === 'Unassigned') return '❓';
        if (name.toLowerCase().includes('pvp')) return '⚔️';
        if (name.toLowerCase().includes('mining')) return '⛏️';
        if (name.toLowerCase().includes('industry')) return '🏭';
        if (name.toLowerCase().includes('trade')) return '💰';
        if (name.toLowerCase().includes('exploration')) return '🔭';
        return '📍';
    };

    return (
        <GlassCard className="p-6 max-w-sm hover:border-blue-500/30 transition-all duration-300">
            {/* Group Header */}
            <div className="flex justify-between items-center mb-4">
                <div className="flex items-center gap-2">
                    <span className="text-2xl">{getGroupIcon(groupName)}</span>
                    <span className="text-lg font-display font-semibold">
                        {groupName}
                    </span>
                </div>
                <Badge 
                    badgeContent={characters.length} 
                    color="primary"
                    sx={{ 
                        '& .MuiBadge-badge': { 
                            backgroundColor: '#14b8a6',
                            fontWeight: 'bold' 
                        } 
                    }}
                />
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
        </GlassCard>
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
