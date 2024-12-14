import React, { useState, useEffect } from 'react';
import PropTypes from 'prop-types';
import { IconButton, Select, MenuItem, TextField, Tooltip } from '@mui/material';
import { Delete, Check as CheckIcon } from '@mui/icons-material';
import OpenInNewIcon from '@mui/icons-material/OpenInNew';
import CharacterDetailModal from "../common/CharacterDetailModal.jsx";
import { formatSP } from "../../utils/formatter.jsx";


const CharacterItem = ({
                           character,
                           onUpdateCharacter,
                           onRemoveCharacter = () => {},
                           roles,
                           hideRemoveIcon = false,
                           skillConversions,
                       }) => {
    const [role, setRole] = useState(character.Role || '');
    const [isAddingRole, setIsAddingRole] = useState(false);
    const [newRole, setNewRole] = useState('');

    const totalSp = character?.Character?.CharacterSkillsResponse?.total_sp;
    const formattedSP = totalSp ? formatSP(totalSp): '0';

    const [detailOpen, setDetailOpen] = useState(false);

    useEffect(() => {
        setRole(character.Role || '');
    }, [character.Role]);

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
            if (onUpdateCharacter && character.Character && character.Character.CharacterID) {
                onUpdateCharacter(character.Character.CharacterID, { Role: selectedRole });
            }
        }
    };

    const handleAddRole = () => {
        if (newRole.trim() !== '') {
            const trimmedRole = newRole.trim();
            setRole(trimmedRole);
            if (onUpdateCharacter && character.Character && character.Character.CharacterID) {
                onUpdateCharacter(character.Character.CharacterID, { Role: trimmedRole });
            }
            setIsAddingRole(false);
            setNewRole('');
        }
    };

    const zKillUrl = `https://zkillboard.com/character/${character.Character.CharacterID}/`;

    // MCT tooltip text
    const mctTooltip = character.MCT
        ? `Training: ${character?.Training || 'Unknown'}`
        : 'Skill queue paused';

    return (
        <div className="p-2 rounded-md shadow-sm bg-gray-700">
            <CharacterDetailModal
                open={detailOpen}
                onClose={() => setDetailOpen(false)}
                character={character}
                skillConversions={skillConversions}
            />
            <div className="flex justify-between items-center">
                <div className="flex items-center space-x-2">
                    <span className="font-semibold text-sm text-teal-200 cursor-pointer underline"
                        onClick={() => setDetailOpen(true)} >
                        {character.Character.CharacterName}
                    </span>
                    <Tooltip title="Total Skillpoints">
                        <span className="text-xs text-teal-400">
                            {formattedSP}
                        </span>
                    </Tooltip>
                    <Tooltip title="Open zKillboard">
                        <IconButton
                            aria-label="Open zKillboard"
                            size="small"
                            onClick={() => {
                                if (window.electronAPI && window.electronAPI.openExternal) {
                                    window.electronAPI.openExternal(zKillUrl);
                                } else {
                                    window.open(zKillUrl, '_blank', 'noopener,noreferrer');
                                }
                            }}
                            sx={{ color: '#99f6e4', '&:hover': { color: '#ffffff' } }}
                        >
                            <OpenInNewIcon fontSize="inherit" />
                        </IconButton>
                    </Tooltip>
                </div>

                <Tooltip title={mctTooltip}>
                    <div
                        data-testid="mct-indicator"
                        className={`w-3 h-3 rounded-full ${character.MCT ? 'bg-green-400' : 'bg-gray-400'}`}
                    ></div>
                </Tooltip>
            </div>

            {/* Second row */}
            <div className="mt-1 flex items-center justify-between">
                <div className="flex items-center space-x-2">
                    {/* Role */}
                    <div className="flex items-center">
                        <span className="text-xs text-teal-400 mr-1">Role:</span>
                        {isAddingRole ? (
                            <div className="flex items-center space-x-1">
                                <TextField
                                    size="small"
                                    value={newRole}
                                    onChange={(e) => setNewRole(e.target.value)}
                                    placeholder="Enter new role"
                                    sx={{ maxWidth: '100px' }}
                                    inputProps={{ style: { fontSize: '0.75rem' } }}
                                />
                                <IconButton aria-label="Confirm new role" onClick={handleAddRole} size="small">
                                    <CheckIcon fontSize="small" />
                                </IconButton>
                            </div>
                        ) : (
                            <Select
                                value={role}
                                onChange={handleRoleChange}
                                displayEmpty
                                size="small"
                                className="text-xs"
                                inputProps={{ 'aria-label': 'Role selection' }}
                                sx={{ fontSize: '0.75rem', maxWidth: '100px' }}
                            >
                                <MenuItem value="" disabled>
                                    Select Role
                                </MenuItem>
                                {rolesOptions.map((r) => (
                                    <MenuItem key={r} value={r} sx={{ fontSize: '0.75rem' }}>
                                        {r}
                                    </MenuItem>
                                ))}
                                <MenuItem value="add_new_role" sx={{ fontSize: '0.75rem' }}>Add New Role</MenuItem>
                            </Select>
                        )}
                    </div>
                    {/* Location */}
                    <div className="text-xs text-teal-400">
                        {character.Character.LocationName || 'Unknown'}
                    </div>
                </div>
                {/* Trash Can Icon */}
                {!hideRemoveIcon && (
                    <Tooltip title="Remove Character">
                        <IconButton
                            aria-label="Remove Character"
                            size="small"
                            onClick={() => onRemoveCharacter(character.Character.CharacterID)}
                            className="text-red-500"
                        >
                            <Delete fontSize="small" />
                        </IconButton>
                    </Tooltip>
                )}
            </div>
        </div>
    );
};

CharacterItem.propTypes = {
    character: PropTypes.object.isRequired,
    onUpdateCharacter: PropTypes.func.isRequired,
    onRemoveCharacter: PropTypes.func,
    roles: PropTypes.array.isRequired,
    hideRemoveIcon: PropTypes.bool,
    skillConversions: PropTypes.object.isRequired,
};

export default CharacterItem;
