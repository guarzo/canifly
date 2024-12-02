// CharactersByRole.jsx
import React, { useMemo, useState } from 'react';
import PropTypes from 'prop-types';
import {
    Container,
    Paper,
    Typography,
    List,
    ListItem,
    ListItemText,
    Select,
    MenuItem,
    InputLabel,
    FormControl,
} from '@mui/material';

const CharactersByRole = ({ accounts, unassignedCharacters, roles }) => {
    const [selectedRole, setSelectedRole] = useState('');

    // Combine all characters
    const allCharacters = useMemo(() => {
        let chars = [];
        accounts.forEach((account) => {
            chars = chars.concat(account.Characters);
        });
        if (unassignedCharacters) {
            chars = chars.concat(unassignedCharacters);
        }
        return chars;
    }, [accounts, unassignedCharacters]);

    // Group characters by role
    const roleMap = useMemo(() => {
        const map = {};
        roles.forEach((role) => {
            map[role] = [];
        });
        map['Unassigned'] = [];

        allCharacters.forEach((character) => {
            const role = character.Role || 'Unassigned';
            if (!map[role]) {
                map[role] = [];
            }
            map[role].push(character);
        });
        return map;
    }, [allCharacters, roles]);

    const displayedRoles = useMemo(() => {
        return selectedRole ? [selectedRole] : Object.keys(roleMap);
    }, [selectedRole, roleMap]);

    return (
        <Container maxWidth="md" className="mt-16">
            <Paper elevation={3} className="p-6 bg-gray-800 text-gray-100">
                <Typography variant="h4" gutterBottom>
                    Characters by Role
                </Typography>

                {/* Role Filter */}
                <FormControl variant="filled" size="small" sx={{ minWidth: 200, marginBottom: 2 }}>
                    <InputLabel id="role-filter-label">Filter by Role</InputLabel>
                    <Select
                        labelId="role-filter-label"
                        value={selectedRole}
                        onChange={(e) => setSelectedRole(e.target.value)}
                        label="Filter by Role"
                    >
                        <MenuItem value="">
                            <em>All Roles</em>
                        </MenuItem>
                        {roles.map((role) => (
                            <MenuItem key={role} value={role}>
                                {role}
                            </MenuItem>
                        ))}
                        <MenuItem value="Unassigned">Unassigned</MenuItem>
                    </Select>
                </FormControl>

                {displayedRoles.map((role) => (
                    <div key={role} className="mb-4">
                        <Typography variant="h6" color="primary">
                            {role}
                        </Typography>
                        <List dense>
                            {roleMap[role].map((character) => (
                                <ListItem key={character.Character.CharacterID}>
                                    <ListItemText primary={character.Character.CharacterName} />
                                </ListItem>
                            ))}
                        </List>
                    </div>
                ))}
            </Paper>
        </Container>
    );
};

CharactersByRole.propTypes = {
    accounts: PropTypes.array.isRequired,
    unassignedCharacters: PropTypes.array,
    roles: PropTypes.array.isRequired,
};

export default CharactersByRole;
