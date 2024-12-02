import React, { useState, useEffect } from 'react';
import PropTypes from 'prop-types';
import { IconButton, Select, MenuItem, TextField } from '@mui/material';
import { Delete, Check as CheckIcon } from '@mui/icons-material';

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
        <div className="p-2 rounded-md shadow-md bg-gray-800 text-teal-200 max-w-sm">
            {/* Account Header */}
            <div className="flex justify-between items-center mb-2">
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
                <button
                    onClick={() => onToggleAccountStatus(account.ID)}
                    className="text-xl font-bold text-white"
                >
                    {account.Status === 'Alpha' ? 'α' : 'Ω'}
                </button>
            </div>

            {/* Characters in Account */}
            <div className="space-y-2">
                {account.Characters.map((character) => {
                    const [role, setRole] = useState(character.Role || '');
                    const [isAddingRole, setIsAddingRole] = useState(false);
                    const [newRole, setNewRole] = useState('');

                    // Update role when character.Role changes
                    useEffect(() => {
                        setRole(character.Role || '');
                    }, [character.Role]);

                    // Combine roles and the current role if it's not already included
                    const rolesOptions = React.useMemo(() => {
                        const combinedRoles = [...roles];
                        if (role && !roles.includes(role) && role !== 'add_new_role') {
                            combinedRoles.push(role);
                        }
                        return combinedRoles;
                    }, [roles, role]);

                    const handleRoleChange = (event) => {
                        const selectedRole = event.target.value;
                        if (selectedRole === 'add_new_role') {
                            setIsAddingRole(true);
                        } else {
                            setRole(selectedRole);
                            onUpdateCharacter(character.Character.CharacterID, { Role: selectedRole });
                        }
                    };

                    const handleAddRole = () => {
                        if (newRole.trim() !== '') {
                            const trimmedRole = newRole.trim();
                            setRole(trimmedRole);
                            onUpdateCharacter(character.Character.CharacterID, { Role: trimmedRole });
                            setIsAddingRole(false);
                            setNewRole('');
                        }
                    };

                    return (
                        <div
                            key={character.Character.CharacterID}
                            className="p-2 rounded-md shadow-sm bg-gray-700"
                        >
                            <div className="flex justify-between items-center">
                <span className="font-semibold text-sm">
                  {character.Character.CharacterName}
                </span>
                                {/* MCT Toggle */}
                                <div
                                    onClick={() =>
                                        onUpdateCharacter(character.Character.CharacterID, {
                                            MCT: !character.MCT,
                                        })
                                    }
                                    className={`cursor-pointer w-3 h-3 rounded-full ${
                                        character.MCT ? 'bg-green-400' : 'bg-gray-400'
                                    }`}
                                ></div>
                            </div>
                            {/* Role, Location, and Trash Can on the Same Line */}
                            <div className="mt-1 flex items-center justify-between">
                                {/* Role */}
                                <div className="flex items-center">
                                    <span className="text-xs text-teal-400">Role:</span>
                                    {isAddingRole ? (
                                        <div className="flex items-center space-x-1">
                                            <TextField
                                                size="small"
                                                value={newRole}
                                                onChange={(e) => setNewRole(e.target.value)}
                                                placeholder="Enter new role"
                                            />
                                            <IconButton onClick={handleAddRole} size="small">
                                                <CheckIcon fontSize="small" />
                                            </IconButton>
                                        </div>
                                    ) : (
                                        <Select
                                            value={role}
                                            onChange={handleRoleChange}
                                            displayEmpty
                                            size="small"
                                            className="ml-1 text-xs"
                                            sx={{ fontSize: '0.75rem' }}
                                        >
                                            <MenuItem value="" disabled>
                                                Select Role
                                            </MenuItem>
                                            {rolesOptions.map((r) => (
                                                <MenuItem key={r} value={r}>
                                                    {r}
                                                </MenuItem>
                                            ))}
                                            <MenuItem value="add_new_role">Add New Role</MenuItem>
                                        </Select>
                                    )}
                                </div>
                                {/* Location */}
                                <div className="text-xs text-teal-400">
                                    {character.Character.LocationName || 'Unknown'}
                                </div>
                                {/* Trash Can Icon */}
                                <IconButton
                                    size="small"
                                    onClick={() => onRemoveCharacter(character.Character.CharacterID)}
                                    className="text-red-500"
                                >
                                    <Delete fontSize="small" />
                                </IconButton>
                            </div>
                        </div>
                    );
                })}
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
