// Dashboard.jsx
import PropTypes from 'prop-types';
import { useState, useMemo } from 'react';
import AccountCard from './AccountCard.jsx';
import GroupCard from './GroupCard.jsx';
import {
    Typography,
    Box,
    Tooltip,
    ToggleButtonGroup,
    ToggleButton,
    IconButton
} from '@mui/material';
import {
    ArrowUpward,
    ArrowDownward,
    AccountBalance,
    AccountCircle,
    Place
} from '@mui/icons-material';

/**
 * A unified dashboard that can display characters grouped by account, role, or location.
 *
 * @param {Object} props
 * @param {Array} props.accounts - Array of account objects.
 * @param {Function} props.onToggleAccountStatus
 * @param {Function} props.onUpdateCharacter
 * @param {Function} props.onUpdateAccountName
 * @param {Function} props.onRemoveCharacter
 * @param {Function} props.onRemoveAccount
 * @param {Array} props.roles - List of roles
 */
const Dashboard = ({
                       accounts,
                       onToggleAccountStatus,
                       onUpdateCharacter,
                       onUpdateAccountName,
                       onRemoveCharacter,
                       onRemoveAccount,
                       roles,
                   }) => {
    // 'account', 'role', 'location'
    const [view, setView] = useState('account');
    const [sortOrder, setSortOrder] = useState('asc');

    const toggleSortOrder = () => {
        setSortOrder((prev) => (prev === 'asc' ? 'desc' : 'asc'));
    };

    const sortIconColor = sortOrder === 'asc' ? '#14b8a6' : '#f59e0b';
    const sortIcon =
        sortOrder === 'asc' ? (
            <ArrowUpward fontSize="small" sx={{ color: sortIconColor }} />
        ) : (
            <ArrowDownward fontSize="small" sx={{ color: sortIconColor }} />
        );

    // Prepare all characters
    const allCharacters = useMemo(() => {
        let chars = [];
        (accounts || []).forEach((account) => {
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
    }, [accounts]);

    // Group by role
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

    // Group by location
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

    // Sorting logic for accounts
    const sortedAccounts = useMemo(() => {
        if (!accounts) return [];
        const accountsCopy = [...accounts];
        accountsCopy.sort((a, b) => {
            const nameA = a.Name || 'Unknown Account';
            const nameB = b.Name || 'Unknown Account';
            return sortOrder === 'asc' ? nameA.localeCompare(nameB) : nameB.localeCompare(nameA);
        });
        return accountsCopy;
    }, [accounts, sortOrder]);

    const mapToDisplay = view === 'role' ? roleMap : locationMap;

    const sortedGroups = useMemo(() => {
        if (view === 'account') return [];
        const keys = Object.keys(mapToDisplay);
        keys.sort((a, b) => (sortOrder === 'asc' ? a.localeCompare(b) : b.localeCompare(a)));
        return keys;
    }, [view, mapToDisplay, sortOrder]);

    const handleViewChange = (event, newValue) => {
        if (newValue !== null) {
            setView(newValue);
        }
    };

    return (
        <div className="bg-gray-900 min-h-screen text-teal-200 px-4 pb-10 pt-16">
            {/* Top bar with view toggle and sort order */}
            <Box
                sx={{
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'space-between',
                    mb: 3,
                }}
            >
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
                        <ToggleButton value="account">
                            <Tooltip title="Account">
                                <AccountBalance fontSize="small" />
                            </Tooltip>
                        </ToggleButton>
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
                </Box>

                {/* Sort Order Control (now always visible) */}
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

            {view === 'account' ? (
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                    {sortedAccounts.length === 0 ? (
                        <Box textAlign="center" mt={4}>
                            <Typography variant="body1" sx={{ color: '#99f6e4' }}>
                                No accounts found.
                            </Typography>
                        </Box>
                    ) : (
                        sortedAccounts.map((account) => (
                            <AccountCard
                                key={account.ID}
                                account={account}
                                onToggleAccountStatus={onToggleAccountStatus}
                                onUpdateAccountName={onUpdateAccountName}
                                onUpdateCharacter={onUpdateCharacter}
                                onRemoveCharacter={onRemoveCharacter}
                                onRemoveAccount={onRemoveAccount}
                                roles={roles}
                            />
                        ))
                    )}
                </div>
            ) : (
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
                            />
                        ))
                    )}
                </div>
            )}
        </div>
    );
};

Dashboard.propTypes = {
    accounts: PropTypes.array.isRequired,
    onToggleAccountStatus: PropTypes.func.isRequired,
    onUpdateCharacter: PropTypes.func.isRequired,
    onUpdateAccountName: PropTypes.func.isRequired,
    onRemoveCharacter: PropTypes.func.isRequired,
    onRemoveAccount: PropTypes.func.isRequired,
    roles: PropTypes.array.isRequired,
};

export default Dashboard;
