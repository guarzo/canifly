import { useState, useEffect, useRef } from 'react';
import PropTypes from 'prop-types';
import { motion, AnimatePresence } from 'framer-motion';
import { IconButton, Select, MenuItem, TextField, Typography } from '@mui/material';
import { Check as CheckIcon, Close as CloseIcon } from '@mui/icons-material';
import GlassCard from '../ui/GlassCard';
import FuturisticButton from '../ui/FuturisticButton';

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

    const hasExistingAccounts = existingAccounts && existingAccounts.length > 0;

    return (
        <AnimatePresence>
            {isOpen && (
                <motion.div 
                    className="fixed inset-0 flex items-center justify-center bg-black/70 backdrop-blur-sm z-50"
                    initial={{ opacity: 0 }}
                    animate={{ opacity: 1 }}
                    exit={{ opacity: 0 }}
                    onClick={onClose}
                >
                    <motion.div
                        initial={{ scale: 0.9, opacity: 0 }}
                        animate={{ scale: 1, opacity: 1 }}
                        exit={{ scale: 0.9, opacity: 0 }}
                        transition={{ type: "spring", damping: 25, stiffness: 300 }}
                        onClick={(e) => e.stopPropagation()}
                    >
                        <GlassCard className="w-96 p-6">
                            <div className="flex justify-between items-center mb-4">
                                <Typography variant="h6" className="text-gradient font-display">
                                    {title || 'Enter Account Name'}
                                </Typography>
                                <IconButton 
                                    onClick={onClose}
                                    className="hover:bg-red-500/20"
                                    size="small"
                                >
                                    <CloseIcon className="text-gray-400" />
                                </IconButton>
                            </div>

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

                            <div className="flex justify-end space-x-3 mt-6">
                                <FuturisticButton
                                    onClick={onClose}
                                    variant="ghost"
                                    size="md"
                                >
                                    Cancel
                                </FuturisticButton>
                                <FuturisticButton
                                    onClick={handleSubmit}
                                    variant="primary"
                                    size="md"
                                >
                                    Submit
                                </FuturisticButton>
                            </div>
                        </GlassCard>
                    </motion.div>
                </motion.div>
            )}
        </AnimatePresence>
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
