// src/components/Characters/CharactersByLocation.jsx
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

const CharactersByLocation = ({ accounts, unassignedCharacters }) => {
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

    // Group characters by location
    const locationMap = useMemo(() => {
        const map = {};
        allCharacters.forEach((character) => {
            const location = character.Character.LocationName || 'Unknown Location';
            if (!map[location]) {
                map[location] = [];
            }
            map[location].push(character);
        });
        return map;
    }, [allCharacters]);

    const sortedLocations = useMemo(() => {
        const locations = Object.keys(locationMap);
        locations.sort((a, b) =>
            sortOrder === 'asc' ? a.localeCompare(b) : b.localeCompare(a)
        );
        return locations;
    }, [locationMap, sortOrder]);

    const toggleSortOrder = () => {
        setSortOrder((prev) => (prev === 'asc' ? 'desc' : 'asc'));
    };

    return (
        <Container maxWidth="md" className="mt-16">
            <Card elevation={3} className="bg-gray-800 text-gray-100">
                <CardContent>
                    <div className="flex items-center justify-between">
                        <Typography variant="h4" gutterBottom>
                            Characters by Location
                        </Typography>
                        <IconButton onClick={toggleSortOrder} color="inherit">
                            {sortOrder === 'asc' ? <ArrowUpward /> : <ArrowDownward />}
                        </IconButton>
                    </div>

                    {sortedLocations.map((location, index) => (
                        <div key={location}>
                            <Typography variant="h6" color="primary" gutterBottom>
                                {location}
                            </Typography>
                            <List dense>
                                {locationMap[location].map((character) => (
                                    <ListItem key={character.Character.CharacterID}>
                                        <ListItemText
                                            primary={character.Character.CharacterName}
                                            secondary={`Account: ${character.accountName} | Role: ${
                                                character.Role || 'Unassigned'
                                            }`}
                                        />
                                    </ListItem>
                                ))}
                            </List>
                            {index < sortedLocations.length - 1 && <Divider light sx={{ marginY: 2 }} />}
                        </div>
                    ))}
                </CardContent>
            </Card>
        </Container>
    );
};

CharactersByLocation.propTypes = {
    accounts: PropTypes.array.isRequired,
    unassignedCharacters: PropTypes.array,
};

export default CharactersByLocation;
