import { useMemo, useState } from 'react';
import PropTypes from 'prop-types';
import {
    Typography,
    IconButton,
    Divider,
    Box,
    Tooltip,
    ToggleButtonGroup,
    ToggleButton,
} from '@mui/material';
import { ArrowUpward, ArrowDownward, AccountCircle, Place } from '@mui/icons-material';
import CharacterItem from './CharacterItem'; // Import the character item component

const CharacterSort = ({ accounts, roles, onUpdateCharacter }) => {
    const [sortOrder, setSortOrder] = useState('asc');
    const [view, setView] = useState('role'); // 'role' or 'location'

    // Combine all characters from accounts
    const allCharacters = useMemo(() => {
        let chars = [];
        (accounts || []).forEach((account) => {
            const accountName = account.Name || 'Unknown Account';
            chars = chars.concat(
                (account.Characters || []).map((char) => ({
                    ...char,
                    accountName,
                    Role: char.Role || '', // Ensure Role field
                    MCT: typeof char.MCT === 'boolean' ? char.MCT : false, // Ensure MCT boolean
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

    const handleViewChange = (event, newValue) => {
        if (newValue !== null) {
            setView(newValue);
        }
    };

    // Reusable GroupedList Component using CharacterItem
    const GroupedList = ({ groups, groupMap }) => (
        <>
            {groups.length === 0 ? (
                <Box textAlign="center" mt={4}>
                    <Typography variant="body1" sx={{ color: '#99f6e4' }}>
                        No characters found.
                    </Typography>
                </Box>
            ) : (
                groups.map((group, index) => (
                    <Box key={group} mb={3}>
                        <Box
                            sx={{
                                display: 'flex',
                                alignItems: 'center',
                                justifyContent: 'space-between',
                                mb: 1,
                            }}
                        >
                            <Typography
                                variant="h6"
                                sx={{ color: '#14b8a6', fontWeight: 'bold' }}
                            >
                                {group}
                            </Typography>
                        </Box>
                        {/* Display characters using CharacterItem */}
                        <Box display="flex" flexDirection="column" gap={1}>
                            {groupMap[group]?.map((character) => (
                                <CharacterItem
                                    key={character.Character.CharacterID}
                                    character={character}
                                    onUpdateCharacter={onUpdateCharacter}
                                    roles={roles}
                                    hideRemoveIcon={true}
                                />
                            ))}
                        </Box>
                        {index < groups.length - 1 && (
                            <Divider sx={{ borderColor: 'rgba(255,255,255,0.2)', my: 2 }} />
                        )}
                    </Box>
                ))
            )}
        </>
    );

    GroupedList.propTypes = {
        groups: PropTypes.array.isRequired,
        groupMap: PropTypes.object.isRequired,
    };

    const groupsToDisplay = view === 'role' ? sortedRoles : sortedLocations;
    const mapToDisplay = view === 'role' ? roleMap : locationMap;

    // Determine icon styling based on sort order
    const sortIconColor = sortOrder === 'asc' ? '#14b8a6' : '#f59e0b';
    const sortIcon = sortOrder === 'asc' ? <ArrowUpward fontSize="small" sx={{ color: sortIconColor }} /> : <ArrowDownward fontSize="small" sx={{ color: sortIconColor }} />;

    return (
        <div className="bg-gray-900 min-h-screen text-teal-200 px-4 pb-10 pt-16">
            {/* Top bar with role/location toggle on left and asc/desc on right */}
            <Box
                sx={{
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'space-between',
                    mb: 3,
                }}
            >
                {/* Role/Location Toggle */}
                <ToggleButtonGroup
                    value={view}
                    exclusive
                    onChange={handleViewChange}
                    sx={{
                        backgroundColor: 'rgba(255,255,255,0.05)',
                        borderRadius: '9999px',
                        '.MuiToggleButton-root': {
                            textTransform: 'none',
                            color: '#99f6e4',
                            fontWeight: 'normal',
                            border: 'none',
                            borderRadius: '9999px',
                            '&.Mui-selected': {
                                backgroundColor: '#14b8a6 !important',
                                color: '#ffffff !important',
                                fontWeight: 'bold',
                            },
                            '&:hover': {
                                backgroundColor: 'rgba(255,255,255,0.1)',
                            },
                            minWidth: '40px',
                            minHeight: '40px',
                        },
                    }}
                >
                    <ToggleButton value="role">
                        <Tooltip title="Role">
                            <AccountCircle fontSize="small" />
                        </Tooltip>
                    </ToggleButton>
                    <ToggleButton value="location">
                        <Tooltip title="Location">
                            <Place fontSize="small" />
                        </Tooltip>
                    </ToggleButton>
                </ToggleButtonGroup>

                {/* Sort Order Control */}
                <Box
                    sx={{
                        display: 'flex',
                        alignItems: 'center',
                        gap: 1,
                        backgroundColor: 'rgba(255,255,255,0.05)',
                        borderRadius: '9999px',
                        paddingX: 1,
                        paddingY: 0.5,
                    }}
                >
                        <IconButton
                            onClick={toggleSortOrder}
                            sx={{
                                '&:hover': { backgroundColor: 'rgba(255,255,255,0.1)' },
                                padding: '4px',
                            }}
                            size="small"
                        >
                            {sortIcon}
                        </IconButton>
                </Box>
            </Box>

            <Box>
                <GroupedList groups={groupsToDisplay} groupMap={mapToDisplay} />
            </Box>
        </div>
    );
};

CharacterSort.propTypes = {
    accounts: PropTypes.array.isRequired,
    roles: PropTypes.array.isRequired,
    onUpdateCharacter: PropTypes.func.isRequired,
};

export default CharacterSort;
