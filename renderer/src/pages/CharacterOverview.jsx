import PropTypes from 'prop-types';
import React, { useState, useMemo } from 'react';
import AccountCard from '../components/dashboard/AccountCard.jsx';
import GroupCard from '../components/dashboard/GroupCard.jsx';
import {
    Typography,
    Box,
    Tooltip,
    ToggleButtonGroup,
    ToggleButton,
    IconButton,
} from '@mui/material';
import {
    ArrowUpward,
    ArrowDownward,
    AccountBalance,
    AccountCircle,
    Place,
    Visibility as VisibilityIcon,
    VisibilityOff as VisibilityOffIcon,
} from '@mui/icons-material';
import { overviewInstructions } from '../utils/instructions.jsx';
import PageHeader from '../components/common/SubPageHeader.jsx';

const CharacterOverview = ({
                               accounts,
                               onToggleAccountStatus,
                               onUpdateCharacter,
                               onUpdateAccountName,
                               onRemoveCharacter,
                               onRemoveAccount,
                               roles,
                               skillConversions,
                               onToggleAccountVisibility,
                           }) => {
    const [view, setView] = useState('account');
    const [sortOrder, setSortOrder] = useState('asc');
    const [showHiddenAccounts, setShowHiddenAccounts] = useState(false);

    const toggleSortOrder = () => {
        setSortOrder((prev) => (prev === 'asc' ? 'desc' : 'asc'));
    };

    // 1. Filter + sort the accounts for "account" view
    const filteredAndSortedAccounts = useMemo(() => {
        if (!accounts) return [];

        const accountsCopy = [...accounts];
        // Sort accounts by Name
        accountsCopy.sort((a, b) => {
            const nameA = a.Name || 'Unknown Account';
            const nameB = b.Name || 'Unknown Account';
            return sortOrder === 'asc'
                ? nameA.localeCompare(nameB)
                : nameB.localeCompare(nameA);
        });

        // Hide them if showHiddenAccounts is false
        if (!showHiddenAccounts) {
            return accountsCopy.filter((acct) => acct.Visible !== false);
        }

        return accountsCopy;
    }, [accounts, sortOrder, showHiddenAccounts]);

    // 2. For role/location grouping, build "allCharacters"
    //    but skip hidden accounts if showHiddenAccounts === false
    const allCharacters = useMemo(() => {
        let chars = [];
        (accounts || []).forEach((account) => {
            // If we are NOT showing hidden accounts, skip characters
            // from an account that is hidden (Visible === false).
            if (!showHiddenAccounts && account.Visible === false) {
                return;
            }

            const accountName = account.Name || 'Unknown Account';
            chars = chars.concat(
                (account.Characters || []).map((char) => ({
                    ...char,
                    accountName,
                    Role: char.Role || '',
                    MCT: typeof char.MCT === 'boolean' ? char.MCT : false,
                }))
            );
        });
        return chars;
    }, [accounts, showHiddenAccounts]);

    // 3. Build the role map
    const roleMap = useMemo(() => {
        const map = { Unassigned: [] };
        roles.forEach((r) => {
            map[r] = [];
        });
        allCharacters.forEach((character) => {
            const charRole = character.Role || 'Unassigned';
            if (!map[charRole]) {
                map[charRole] = [];
            }
            map[charRole].push(character);
        });
        return map;
    }, [allCharacters, roles]);

    // 4. Build the location map
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

    // Decide which map to display if user selects 'role' or 'location'
    const mapToDisplay = view === 'role' ? roleMap : locationMap;

    // 5. Sort the group keys (role/location) for the GroupCard display
    const sortedGroups = useMemo(() => {
        if (view === 'account') return [];
        const keys = Object.keys(mapToDisplay);
        keys.sort((a, b) =>
            sortOrder === 'asc' ? a.localeCompare(b) : b.localeCompare(a)
        );
        return keys;
    }, [view, mapToDisplay, sortOrder]);

    // Switch the "view" among account/role/location
    const handleViewChange = (event, newValue) => {
        if (newValue !== null) {
            setView(newValue);
        }
    };

    // Determine icon color and direction
    const sortIconColor = sortOrder === 'asc' ? '#14b8a6' : '#f59e0b';
    const sortIcon =
        sortOrder === 'asc' ? (
            <ArrowUpward fontSize="small" sx={{ color: sortIconColor }} />
        ) : (
            <ArrowDownward fontSize="small" sx={{ color: sortIconColor }} />
        );

    return (
        <div className="bg-gray-900 min-h-screen text-teal-200 px-4 pb-10 pt-16">
            <PageHeader
                title="Character Overview"
                instructions={overviewInstructions}
                storageKey="showDashboardInstructions"
            />

            <Box
                sx={{
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'space-between',
                    mb: 3,
                }}
            >
                {/* Group By: Account, Role, Location */}
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
                    <Typography variant="body2" sx={{ color: '#99f6e4' }}>
                        Group by:
                    </Typography>
                    <ToggleButtonGroup
                        value={view}
                        exclusive
                        onChange={handleViewChange}
                        sx={{
                            backgroundColor: 'rgba(255,255,255,0.05)',
                            borderRadius: '9999px',
                            padding: '2px',
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
                        <ToggleButton value="account" aria-label="Account">
                            <Tooltip title="Account">
                                <AccountBalance fontSize="small" />
                            </Tooltip>
                        </ToggleButton>
                        <ToggleButton value="role" aria-label="Role">
                            <Tooltip title="Role">
                                <AccountCircle fontSize="small" />
                            </Tooltip>
                        </ToggleButton>
                        <ToggleButton value="location" aria-label="Location">
                            <Tooltip title="Location">
                                <Place fontSize="small" />
                            </Tooltip>
                        </ToggleButton>
                    </ToggleButtonGroup>
                </Box>

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
                    <Typography variant="body2" sx={{ color: '#99f6e4' }}>
                        Sort:
                    </Typography>
                    <IconButton
                        onClick={toggleSortOrder}
                        aria-label="Sort"
                        sx={{
                            '&:hover': {
                                backgroundColor: 'rgba(255,255,255,0.1)',
                            },
                            padding: '4px',
                        }}
                        size="small"
                    >
                        {sortIcon}
                    </IconButton>
                </Box>

                {/* Show/Hide hidden accounts toggle (only in 'account' view) */}
                {view === 'account' && (
                    <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
                        <Tooltip
                            title={
                                showHiddenAccounts
                                    ? 'Hide hidden accounts'
                                    : 'Show hidden accounts'
                            }
                        >
                            <IconButton
                                onClick={() => setShowHiddenAccounts(!showHiddenAccounts)}
                                sx={{
                                    color: showHiddenAccounts ? '#10b981' : '#6b7280',
                                }}
                            >
                                {showHiddenAccounts ? (
                                    <VisibilityIcon />
                                ) : (
                                    <VisibilityOffIcon />
                                )}
                            </IconButton>
                        </Tooltip>
                    </Box>
                )}
            </Box>

            {/* RENDER LOGIC */}
            {view === 'account' ? (
                /* ------------- ACCOUNT VIEW ------------- */
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                    {filteredAndSortedAccounts.length === 0 ? (
                        <Box textAlign="center" mt={4}>
                            <Typography variant="body1" sx={{ color: '#99f6e4' }}>
                                No accounts found.
                            </Typography>
                        </Box>
                    ) : (
                        filteredAndSortedAccounts.map((account) => (
                            <AccountCard
                                key={account.ID}
                                account={account}
                                onToggleAccountStatus={onToggleAccountStatus}
                                onUpdateAccountName={onUpdateAccountName}
                                onUpdateCharacter={onUpdateCharacter}
                                onRemoveCharacter={onRemoveCharacter}
                                onRemoveAccount={onRemoveAccount}
                                roles={roles}
                                skillConversions={skillConversions}
                                onToggleAccountVisibility={onToggleAccountVisibility}
                            />
                        ))
                    )}
                </div>
            ) : (
                /* ------------- GROUP VIEW (ROLE or LOCATION) ------------- */
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                    {sortedGroups.length === 0 ? (
                        <Box textAlign="center" mt={4}>
                            <Typography variant="body1" sx={{ color: '#99f6e4' }}>
                                No characters found.
                            </Typography>
                        </Box>
                    ) : (
                        sortedGroups.map((group) => (
                            <GroupCard
                                key={group}
                                groupName={group}
                                characters={mapToDisplay[group] || []}
                                onUpdateCharacter={onUpdateCharacter}
                                roles={roles}
                                skillConversions={skillConversions}
                            />
                        ))
                    )}
                </div>
            )}
        </div>
    );
};

CharacterOverview.propTypes = {
    accounts: PropTypes.array.isRequired,
    onToggleAccountStatus: PropTypes.func.isRequired,
    onToggleAccountVisibility: PropTypes.func.isRequired,
    onUpdateCharacter: PropTypes.func.isRequired,
    onUpdateAccountName: PropTypes.func.isRequired,
    onRemoveCharacter: PropTypes.func.isRequired,
    onRemoveAccount: PropTypes.func.isRequired,
    roles: PropTypes.array.isRequired,
    skillConversions: PropTypes.object.isRequired,
};

export default CharacterOverview;
