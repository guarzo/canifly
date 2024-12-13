import { useState, useEffect, useRef } from 'react';
import PropTypes from 'prop-types';
import { IconButton, Select, MenuItem, TextField } from '@mui/material';
import { Check as CheckIcon } from '@mui/icons-material';

const AccountPromptModal = ({ isOpen, onClose, onSubmit, title, existingAccounts }) => {
    const [account, setAccount] = useState('');
    const [isAddingAccount, setIsAddingAccount] = useState(false);
    const [newAccount, setNewAccount] = useState('');

    const selectRef = useRef(null);
    const textFieldRef = useRef(null);

    useEffect(() => {
        if (!isOpen) {
            setAccount('');
            setIsAddingAccount(false);
            setNewAccount('');
        } else {
            // Modal just opened
            // Focus the appropriate element
            if (existingAccounts && existingAccounts.length > 0 && !isAddingAccount) {
                // Focus the select component
                if (selectRef.current) {
                    selectRef.current.focus();
                }
            } else {
                // Focus the text field for adding a new account
                if (textFieldRef.current) {
                    textFieldRef.current.querySelector('input')?.focus();
                }
            }
        }
    }, [isOpen, existingAccounts, isAddingAccount]);

    const handleAccountChange = (event) => {
        const selectedValue = event.target.value;
        if (selectedValue === 'add_new_account') {
            setIsAddingAccount(true);
            setAccount('');
            // Focus the text field when switching to add new account
            setTimeout(() => {
                if (textFieldRef.current) {
                    textFieldRef.current.querySelector('input')?.focus();
                }
            }, 0);
        } else {
            setIsAddingAccount(false);
            setAccount(selectedValue);
        }
    };

    const handleAddAccount = () => {
        if (newAccount.trim() !== '') {
            const trimmed = newAccount.trim();
            setAccount(trimmed);
            setIsAddingAccount(false);
            setNewAccount('');
        }
    };

    const handleSubmit = () => {
        const finalAccount = isAddingAccount ? newAccount.trim() : account;
        if (!finalAccount) return;
        onSubmit(finalAccount);
        setAccount('');
        setIsAddingAccount(false);
        setNewAccount('');
    };

    if (!isOpen) return null;

    const hasExistingAccounts = existingAccounts && existingAccounts.length > 0;

    return (
        <div className="fixed inset-0 flex items-center justify-center bg-black bg-opacity-50 z-50">
            <div className="bg-gray-800 text-teal-200 p-6 rounded shadow-md w-80">
                <h2 className="mb-4 text-lg font-semibold">{title || 'Enter Account Name'}</h2>

                {!isAddingAccount && hasExistingAccounts && (
                    <div className="mb-4">
                        <Select
                            value={account || ''}
                            onChange={handleAccountChange}
                            displayEmpty
                            fullWidth
                            inputRef={selectRef}
                            sx={{
                                backgroundColor: 'background.paper',
                                borderRadius: 1,
                                '& .MuiSelect-select': {
                                    padding: '10px 14px',
                                },
                            }}
                        >
                            <MenuItem value="" disabled>
                                Select Account
                            </MenuItem>
                            {existingAccounts.map((acc) => (
                                <MenuItem key={acc} value={acc}>
                                    {acc}
                                </MenuItem>
                            ))}
                            <MenuItem value="add_new_account">Add New Account</MenuItem>
                        </Select>
                    </div>
                )}

                {(!hasExistingAccounts || isAddingAccount) && (
                    <div className="mb-4 flex items-center space-x-2" ref={textFieldRef}>
                        <TextField
                            size="small"
                            value={isAddingAccount ? newAccount : account}
                            onChange={(e) =>
                                isAddingAccount ? setNewAccount(e.target.value) : setAccount(e.target.value)
                            }
                            placeholder="Enter account name"
                            fullWidth
                            sx={{
                                '& .MuiInputBase-root': {
                                    padding: '2px',
                                },
                            }}
                        />
                        {isAddingAccount && (
                            <IconButton onClick={handleAddAccount} size="small" color="primary">
                                <CheckIcon fontSize="small" />
                            </IconButton>
                        )}
                    </div>
                )}

                <div className="flex justify-end space-x-3">
                    <button
                        onClick={onClose}
                        className="py-2 px-4 bg-gray-600 text-white rounded hover:bg-gray-700"
                    >
                        Cancel
                    </button>
                    <button
                        onClick={handleSubmit}
                        className="py-2 px-4 bg-teal-600 text-white rounded hover:bg-teal-700"
                    >
                        Submit
                    </button>
                </div>
            </div>
        </div>
    );
};

AccountPromptModal.propTypes = {
    isOpen: PropTypes.bool.isRequired,
    onClose: PropTypes.func.isRequired,
    onSubmit: PropTypes.func.isRequired,
    title: PropTypes.string,
    existingAccounts: PropTypes.arrayOf(PropTypes.string),
};

export default AccountPromptModal;
