// src/components/CharacterSort.jsx
import { useMemo, useState } from 'react';
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
    Tabs,
    Tab,
    Box,
    Tooltip,
} from '@mui/material';
import { ArrowUpward, ArrowDownward } from '@mui/icons-material';
import { makeStyles } from '@mui/styles';

// Define styles using makeStyles
const useStyles = makeStyles((theme) => ({
    container: {
        marginTop: theme.spacing(4),
    },
    card: {
        backgroundColor: theme.palette.background.paper,
        color: theme.palette.text.primary,
    },
    title: {
        marginBottom: theme.spacing(2),
    },
    listItem: {
        paddingLeft: theme.spacing(2),
        paddingRight: theme.spacing(2),
    },
    divider: {
        marginTop: theme.spacing(2),
        marginBottom: theme.spacing(2),
    },
    groupHeader: {
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'space-between',
    },
}));

const CharacterSort = ({ accounts, roles }) => {
    const classes = useStyles();
    const [sortOrder, setSortOrder] = useState('asc');
    const [tabIndex, setTabIndex] = useState(0); // 0: By Role, 1: By Location

    // Combine all characters from accounts, excluding unassigned characters
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

        allCharacters.forEach((character) => {
            const role = character.Role || 'Unassigned';
            if (role !== 'Unassigned' && map[role]) {
                map[role].push(character);
            }
            // Characters with 'Unassigned' role are excluded
        });

        // Remove roles with no characters
        Object.keys(map).forEach((role) => {
            if (map[role].length === 0) {
                delete map[role];
            }
        });

        return map;
    }, [allCharacters, roles]);

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

        // Remove locations with no characters
        Object.keys(map).forEach((location) => {
            if (map[location].length === 0) {
                delete map[location];
            }
        });

        return map;
    }, [allCharacters]);

    // Sorting functions
    const sortedRoles = useMemo(() => {
        const roleList = Object.keys(roleMap);
        roleList.sort((a, b) =>
            sortOrder === 'asc' ? a.localeCompare(b) : b.localeCompare(a)
        );
        return roleList;
    }, [roleMap, sortOrder]);

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

    const handleTabChange = (event, newValue) => {
        setTabIndex(newValue);
    };

    // Reusable GroupedList Component
    const GroupedList = ({ groups, groupMap, groupTitle }) => (
        <>
            {groups.map((group, index) => (
                <Box key={group} mb={3}>
                    <Box className={classes.groupHeader}>
                        <Typography variant="h6" style={{ color: '#14b8a6' }}>
                            {group}
                        </Typography>
                        {/* Optional: Add more controls per group if needed */}
                    </Box>
                    <List dense>
                        {groupMap[group]?.map((character) => (
                            <ListItem key={character.Character.CharacterID} className={classes.listItem}>
                                <ListItemText
                                    primary={character.Character.CharacterName}
                                    secondary={`Account: ${character.accountName} | Location: ${character.Character.LocationName || 'Unknown'}`}
                                />
                            </ListItem>
                        ))}
                    </List>
                    {index < groups.length - 1 && <Divider light />}
                </Box>
            ))}
        </>
    );

    GroupedList.propTypes = {
        groups: PropTypes.array.isRequired,
        groupMap: PropTypes.object.isRequired,
        groupTitle: PropTypes.string.isRequired,
    };

    return (
        <Container maxWidth="lg" className={classes.container}>
            <Card className={classes.card} elevation={3}>
                <CardContent>
                    <Box className={classes.groupHeader} mb={2}>
                        <Typography variant="h4" className={classes.title}>
                            Characters
                        </Typography>
                        <Tooltip title={`Sort Order: ${sortOrder === 'asc' ? 'Ascending' : 'Descending'}`}>
                            <IconButton onClick={toggleSortOrder} color="primary">
                                {sortOrder === 'asc' ? <ArrowUpward /> : <ArrowDownward />}
                            </IconButton>
                        </Tooltip>
                    </Box>
                    <Tabs
                        value={tabIndex}
                        onChange={handleTabChange}
                        indicatorColor="primary"
                        textColor="primary"
                        variant="fullWidth"
                    >
                        <Tab label="By Role" />
                        <Tab label="By Location" />
                    </Tabs>
                    <Box mt={2}>
                        {tabIndex === 0 && (
                            <GroupedList groups={sortedRoles} groupMap={roleMap} groupTitle="Role" />
                        )}
                        {tabIndex === 1 && (
                            <GroupedList groups={sortedLocations} groupMap={locationMap} groupTitle="Location" />
                        )}
                    </Box>
                </CardContent>
            </Card>
        </Container>
    );
};

CharacterSort.propTypes = {
    accounts: PropTypes.array.isRequired,
    roles: PropTypes.array.isRequired,
};

export default CharacterSort;
