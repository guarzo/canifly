// src/components/AccountCard.jsx
import React, { useState } from 'react';
import PropTypes from 'prop-types';
import { IconButton, Select, MenuItem, TextField, Tooltip } from '@mui/material';
import { Delete, Check as CheckIcon } from '@mui/icons-material';
import CharacterItem from './CharacterItem'; // Import the new CharacterItem component

const AccountCard = ({
                         account,
                         onToggleAccountStatus,
                         onUpdateAccountName,
                         onUpdateCharacter,
                         onRemoveCharacter,
                         roles,
                     }) => {
    const [isEditingName, setIsEditingName] = useState(false);
    const [accountName, setAccountName] = useState(account.Name);

    const handleNameChange = (e) => {
        setAccountName(e.target.value);
    };

    const handleNameBlur = () => {
        if (accountName !== account.Name) {
            onUpdateAccountName(account.ID, accountName);
        }
        setIsEditingName(false);
    };

    const handleNameKeyDown = (e) => {
        if (e.key === 'Enter') {
            e.preventDefault();
            e.target.blur();
        }
    };

    const startEditingName = () => {
        setIsEditingName(true);
    };

    return (
        <div className="p-4 rounded-md shadow-md bg-gray-800 text-teal-200 max-w-sm">
            {/* Account Header */}
            <div className="flex justify-between items-center mb-4">
                {isEditingName ? (
                    <input
                        className="bg-transparent border-b border-teal-400 text-sm font-bold focus:outline-none"
                        value={accountName}
                        onChange={handleNameChange}
                        onBlur={handleNameBlur}
                        onKeyDown={handleNameKeyDown}
                        autoFocus
                    />
                ) : (
                    <span className="text-sm font-bold cursor-pointer" onClick={startEditingName}>
                        {account.Name}
                    </span>
                )}
                <Tooltip title="Toggle Account Status">
                    <button
                        onClick={() => onToggleAccountStatus(account.ID)}
                        className="text-xl font-bold text-white"
                    >
                        {account.Status === 'Alpha' ? 'α' : 'Ω'}
                    </button>
                </Tooltip>
            </div>

            {/* Characters in Account */}
            <div className="space-y-2">
                {account.Characters.map((character) => (
                    <CharacterItem
                        key={character.Character.CharacterID}
                        character={character}
                        onUpdateCharacter={onUpdateCharacter}
                        onRemoveCharacter={onRemoveCharacter}
                        roles={roles}
                    />
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
    onRemoveCharacter: PropTypes.func.isRequired,
    roles: PropTypes.array.isRequired,
};

export default AccountCard;
