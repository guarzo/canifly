// src/components/Characters/CharactersByRole.jsx
import React, { useMemo, useState } from 'react';
import PropTypes from 'prop-types';
import {
    Container,
    Card,
    CardContent,
    Typography,
    List,
    ListItem,
    ListItemText,
    IconButton,
    Divider,
} from '@mui/material';
import { ArrowUpward, ArrowDownward } from '@mui/icons-material';

const CharactersByRole = ({ accounts, roles }) => {
    const [sortOrder, setSortOrder] = useState('asc');

    // Combine all characters
    const allCharacters = useMemo(() => {
        let chars = [];
        (accounts || []).forEach((account) => {
            const accountName = account.Name || 'Unknown Account';
            chars = chars.concat(
                (account.Characters || []).map((char) => ({
                    ...char,
                    accountName,
                }))
            );
        });
        return chars;
    }, [accounts]);

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

    const sortedRoles = useMemo(() => {
        const roleList = Object.keys(roleMap);
        roleList.sort((a, b) =>
            sortOrder === 'asc' ? a.localeCompare(b) : b.localeCompare(a)
        );
        return roleList;
    }, [roleMap, sortOrder]);

    const toggleSortOrder = () => {
        setSortOrder((prev) => (prev === 'asc' ? 'desc' : 'asc'));
    };

    return (
        <Container maxWidth="md" className="mt-16">
            <Card elevation={3} className="bg-gray-800 text-gray-100">
                <CardContent>
                    <div className="flex items-center justify-between">
                        <Typography variant="h4" gutterBottom>
                            Characters by Role
                        </Typography>
                        <IconButton onClick={toggleSortOrder} color="inherit">
                            {sortOrder === 'asc' ? <ArrowUpward /> : <ArrowDownward />}
                        </IconButton>
                    </div>

                    {sortedRoles.map((role, index) => (
                        <div key={role}>
                            <Typography variant="h6" color="primary" gutterBottom>
                                {role}
                            </Typography>
                            <List dense>
                                {roleMap[role].map((character) => (
                                    <ListItem key={character.Character.CharacterID}>
                                        <ListItemText
                                            primary={character.Character.CharacterName}
                                            secondary={`Account: ${character.accountName} | Location: ${
                                                character.Character.LocationName || 'Unknown'
                                            }`}
                                        />
                                    </ListItem>
                                ))}
                            </List>
                            {index < sortedRoles.length - 1 && <Divider light sx={{ marginY: 2 }} />}
                        </div>
                    ))}
                </CardContent>
            </Card>
        </Container>
    );
};

CharactersByRole.propTypes = {
    accounts: PropTypes.array.isRequired,
    roles: PropTypes.array.isRequired,
};

export default CharactersByRole;
