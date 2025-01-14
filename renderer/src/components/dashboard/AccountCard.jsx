import React, { useState } from 'react';
import PropTypes from 'prop-types';
import {
    IconButton,
    Menu,
    MenuItem,
    Tooltip,
} from '@mui/material';
import {
    MoreVert as MoreVertIcon,
    Visibility as VisibilityIcon,
    VisibilityOff as VisibilityOffIcon,
} from '@mui/icons-material';
import CharacterItem from './CharacterItem.jsx';

const AccountCard = ({
                         account,
                         onToggleAccountStatus,
                         onToggleAccountVisibility,
                         onUpdateAccountName,
                         onUpdateCharacter,
                         onRemoveCharacter,
                         onRemoveAccount,
                         roles,
                         skillConversions,
                     }) => {
    const [isEditingName, setIsEditingName] = useState(false);
    const [accountName, setAccountName] = useState(account.Name);

    const [anchorEl, setAnchorEl] = useState(null);
    const menuOpen = Boolean(anchorEl);

    const handleNameChange = (e) => setAccountName(e.target.value);

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

    const handleMenuClick = (event) => setAnchorEl(event.currentTarget);
    const handleMenuClose = () => setAnchorEl(null);

    const handleRemoveAccountClick = () => {
        handleMenuClose();
        onRemoveAccount(account.Name);
    };

    // IMPORTANT:
    // Make sure we pass the account.ID to the onToggleAccountVisibility callback.
    // If we just passed onToggleAccountVisibility, it wouldn't have the account ID.
    const handleVisibilityClick = () => {
        onToggleAccountVisibility(account.ID);
    };

    return (
        <div className="p-4 rounded-md shadow-md bg-gray-800 text-teal-200 max-w-sm">
            {/* Account Header */}
            <div className="flex justify-between items-center mb-4">
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
                    <span
                        className="text-sm font-bold cursor-pointer"
                        onClick={startEditingName}
                    >
            {account.Name}
          </span>
                )}

                <div className="flex items-center space-x-2">
                    <Tooltip title="Toggle Account Status (Alpha/Omega)">
                        <button
                            onClick={() => onToggleAccountStatus(account.ID)}
                            className="text-xl font-bold text-white"
                        >
                            {account.Status === 'Alpha' ? 'α' : 'Ω'}
                        </button>
                    </Tooltip>

                    {/* NEW: Icon-based visibility toggle */}
                    <Tooltip
                        title={account.Visible ? 'Hide Account' : 'Show Account'}
                    >
                        <IconButton
                            onClick={handleVisibilityClick}
                            // You can style color differently depending on whether it's visible or not
                            sx={{ color: account.Visible ? '#10b981' : '#6b7280' }}
                        >
                            {account.Visible ? <VisibilityIcon /> : <VisibilityOffIcon />}
                        </IconButton>
                    </Tooltip>

                    {/* Existing More Options menu */}
                    <IconButton
                        onClick={handleMenuClick}
                        size="small"
                        sx={{ color: '#9ca3af' }}
                        aria-label="more options"
                    >
                        <MoreVertIcon fontSize="inherit" />
                    </IconButton>

                    <Menu
                        anchorEl={anchorEl}
                        open={menuOpen}
                        onClose={handleMenuClose}
                        PaperProps={{
                            style: {
                                backgroundColor: '#1f2937',
                                color: '#14b8a6',
                            },
                        }}
                        anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
                        transformOrigin={{ vertical: 'top', horizontal: 'right' }}
                    >
                        <MenuItem onClick={handleRemoveAccountClick}>
                            Remove Account
                        </MenuItem>
                    </Menu>
                </div>
            </div>

            {/* Characters in Account */}
            <div className="space-y-2">
                {account.Characters.map((character) => (
                    <CharacterItem
                        key={character.Character.CharacterID}
                        character={character}
                        onUpdateCharacter={onUpdateCharacter}
                        onRemoveCharacter={onRemoveCharacter}
                        roles={roles}
                        skillConversions={skillConversions}
                    />
                ))}
            </div>
        </div>
    );
};

AccountCard.propTypes = {
    account: PropTypes.object.isRequired,
    onToggleAccountStatus: PropTypes.func.isRequired,
    onToggleAccountVisibility: PropTypes.func.isRequired,
    onUpdateAccountName: PropTypes.func.isRequired,
    onUpdateCharacter: PropTypes.func.isRequired,
    onRemoveCharacter: PropTypes.func.isRequired,
    onRemoveAccount: PropTypes.func.isRequired,
    roles: PropTypes.array.isRequired,
    skillConversions: PropTypes.object.isRequired,
};

export default AccountCard;
