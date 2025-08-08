// src/pages/CharacterOverview.jsx

import React, { useState, useMemo } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
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
    Refresh as RefreshIcon,
} from '@mui/icons-material';
import { overviewInstructions } from '../utils/instructions.jsx';
import PageHeader from '../components/common/SubPageHeader.jsx';
import { useAppData } from '../hooks/useAppData';
import { useAsyncOperation } from '../hooks/useAsyncOperation';
import * as apiService from '../api/apiService';
import { logger } from '../utils/logger';

const CharacterOverview = ({ roles, skillConversions }) => {
    const { accounts, updateAccount, deleteAccount, refreshData, fetchAccounts } = useAppData();
    const { execute } = useAsyncOperation();
    
    const [view, setView] = useState('account');
    const [sortOrder, setSortOrder] = useState('asc');
    const [showHiddenAccounts, setShowHiddenAccounts] = useState(false);
    const [isRefreshing, setIsRefreshing] = useState(false);
    
    // Debug: Log the accounts data
    React.useEffect(() => {
        logger.debug('CharacterOverview - accounts:', accounts);
        if (accounts && accounts.length > 0) {
            logger.debug('First account:', accounts[0]);
            if (accounts[0].Characters && accounts[0].Characters.length > 0) {
                logger.debug('First character:', accounts[0].Characters[0]);
                logger.debug('Character skills response:', accounts[0].Characters[0].Character?.CharacterSkillsResponse);
                logger.debug('Total SP:', accounts[0].Characters[0].Character?.CharacterSkillsResponse?.total_sp);
            }
        }
    }, [accounts]);

    const toggleSortOrder = () => {
        setSortOrder((prev) => (prev === 'asc' ? 'desc' : 'asc'));
    };
    
    // Handler functions to replace prop handlers
    const handleToggleAccountStatus = async (accountId) => {
        const account = accounts.find(acc => acc.Id === accountId);
        if (account) {
            await execute(
                () => updateAccount(accountId, { Status: !account.Status }),
                { successMessage: 'Account status updated' }
            );
        }
    };
    
    const handleUpdateCharacter = async (characterId, updates) => {
        await execute(
            () => apiService.updateCharacter(characterId, updates),
            { successMessage: 'Character updated' }
        );
    };
    
    const handleUpdateAccountName = async (accountId, name) => {
        await execute(
            () => updateAccount(accountId, { Name: name }),
            { successMessage: 'Account name updated' }
        );
    };
    
    const handleRemoveCharacter = async (characterId) => {
        await execute(
            () => apiService.deleteCharacter(characterId),
            { successMessage: 'Character removed' }
        );
    };
    
    const handleRefreshAll = async () => {
        setIsRefreshing(true);
        try {
            let successCount = 0;
            logger.debug('Starting refresh for all characters...');
            logger.debug('Accounts:', accounts);
            
            for (const account of accounts) {
                for (const character of account.Characters) {
                    try {
                        logger.debug(`Refreshing character ${character.Character.CharacterName} (ID: ${character.Character.CharacterID})`);
                        const result = await apiService.refreshCharacter(character.Character.CharacterID);
                        logger.debug(`Refresh result for ${character.Character.CharacterName}:`, result);
                        successCount++;
                    } catch (error) {
                        logger.error(`Failed to refresh character ${character.Character.CharacterName}:`, error);
                    }
                }
            }
            
            logger.debug(`Refreshed ${successCount} characters successfully`);
            
            if (successCount > 0) {
                // Refresh the accounts data to show updated skills
                logger.debug('Refreshing account data...');
                // Force a fresh fetch by calling fetchAccounts directly
                await fetchAccounts();
                logger.debug('Account data refreshed');
            }
        } finally {
            setIsRefreshing(false);
        }
    };
    
    const handleRemoveAccount = async (accountId) => {
        await execute(
            () => deleteAccount(accountId),
            { successMessage: 'Account removed' }
        );
    };
    
    const handleToggleAccountVisibility = async (accountId) => {
        const account = accounts.find(acc => acc.Id === accountId);
        if (account) {
            await execute(
                () => updateAccount(accountId, { Visible: !account.Visible }),
                { successMessage: 'Account visibility updated' }
            );
        }
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

    // 2. Build "allCharacters" for role/location grouping
    //    Skip hidden accounts if showHiddenAccounts === false
    const allCharacters = useMemo(() => {
        let chars = [];
        (accounts || []).forEach((account) => {
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
            // Be sure character.Character exists
            const location = character.Character?.LocationName || 'Unknown Location';
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

        let keys = Object.keys(mapToDisplay);

        // If grouping by role, omit roles that have no members:
        if (view === 'role') {
            keys = keys.filter((role) => (mapToDisplay[role] || []).length > 0);
        }
        // (If you wanted to hide empty locations too, do the same check here.)

        // Then do the alphabetical (asc/desc) sort
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
        <div className="min-h-screen px-4 pb-10 pt-16">
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

                {/* Show/Hide hidden accounts toggle and Refresh button */}
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 2 }}>
                    <Tooltip title="Refresh all character data from ESI">
                        <IconButton
                            onClick={handleRefreshAll}
                            disabled={isRefreshing}
                            sx={{
                                color: '#14b8a6',
                                '&:hover': {
                                    backgroundColor: 'rgba(20, 184, 166, 0.1)',
                                },
                                '&.Mui-disabled': {
                                    color: '#6b7280',
                                },
                            }}
                        >
                            <RefreshIcon sx={{ 
                                animation: isRefreshing ? 'spin 1s linear infinite' : 'none',
                                '@keyframes spin': {
                                    '0%': { transform: 'rotate(0deg)' },
                                    '100%': { transform: 'rotate(360deg)' },
                                },
                            }} />
                        </IconButton>
                    </Tooltip>
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
                            {showHiddenAccounts ? <VisibilityIcon /> : <VisibilityOffIcon />}
                        </IconButton>
                    </Tooltip>
                </Box>
            </Box>

            {/* RENDER LOGIC */}
            <AnimatePresence mode="wait">
            {view === 'account' ? (
                /* ------------- ACCOUNT VIEW ------------- */
                <motion.div 
                    key="account-view"
                    className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6"
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    exit={{ opacity: 0, y: -20 }}
                    transition={{ duration: 0.3 }}
                >
                    {filteredAndSortedAccounts.length === 0 ? (
                        <Box textAlign="center" mt={4} gridColumn="1 / -1">
                            <Typography variant="body1" sx={{ color: '#99f6e4' }}>
                                No accounts found.
                            </Typography>
                        </Box>
                    ) : (
                        filteredAndSortedAccounts.map((account, index) => (
                            <motion.div
                                key={account.ID}
                                initial={{ opacity: 0, scale: 0.9 }}
                                animate={{ opacity: 1, scale: 1 }}
                                transition={{ 
                                    duration: 0.3,
                                    delay: index * 0.05,
                                    type: "spring",
                                    stiffness: 100
                                }}
                            >
                                <AccountCard
                                    account={account}
                                    onToggleAccountStatus={handleToggleAccountStatus}
                                    onUpdateAccountName={handleUpdateAccountName}
                                    onUpdateCharacter={handleUpdateCharacter}
                                    onRemoveCharacter={handleRemoveCharacter}
                                    onRemoveAccount={handleRemoveAccount}
                                    roles={roles}
                                    skillConversions={skillConversions}
                                    onToggleAccountVisibility={handleToggleAccountVisibility}
                                />
                            </motion.div>
                        ))
                    )}
                </motion.div>
            ) : (
                /* ------------- GROUP VIEW (ROLE or LOCATION) ------------- */
                <motion.div 
                    key="group-view"
                    className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6"
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    exit={{ opacity: 0, y: -20 }}
                    transition={{ duration: 0.3 }}
                >
                    {sortedGroups.length === 0 ? (
                        <Box textAlign="center" mt={4} gridColumn="1 / -1">
                            <Typography variant="body1" sx={{ color: '#99f6e4' }}>
                                No characters found.
                            </Typography>
                        </Box>
                    ) : (
                        sortedGroups.map((group, index) => (
                            <motion.div
                                key={group}
                                initial={{ opacity: 0, scale: 0.9 }}
                                animate={{ opacity: 1, scale: 1 }}
                                transition={{ 
                                    duration: 0.3,
                                    delay: index * 0.05,
                                    type: "spring",
                                    stiffness: 100
                                }}
                            >
                                <GroupCard
                                    groupName={group}
                                    characters={mapToDisplay[group] || []}
                                    onUpdateCharacter={handleUpdateCharacter}
                                    roles={roles}
                                    skillConversions={skillConversions}
                                />
                            </motion.div>
                        ))
                    )}
                </motion.div>
            )}
            </AnimatePresence>
        </div>
    );
};

export default CharacterOverview;
