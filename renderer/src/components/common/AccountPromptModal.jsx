import { useState, useEffect, useRef } from 'react';
import {
    Dialog,
    Select,
    MenuItem,
    TextField,
    IconButton,
    Checkbox,
    FormControlLabel,
} from '@mui/material';
import { Close as CloseIcon, Check as CheckIcon } from '@mui/icons-material';

const AccountPromptModal = ({
    isOpen,
    onClose,
    onSubmit,
    title,
    existingAccounts,
    showRememberMe,
    rememberMe,
    onRememberMeChange,
}) => {
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
            return;
        }
        if (existingAccounts && existingAccounts.length > 0 && !isAddingAccount) {
            selectRef.current?.focus();
        } else {
            textFieldRef.current?.querySelector('input')?.focus();
        }
    }, [isOpen, existingAccounts, isAddingAccount]);

    const handleAccountChange = (e) => {
        const v = e.target.value;
        if (v === 'add_new_account') {
            setIsAddingAccount(true);
            setAccount('');
            setTimeout(() => textFieldRef.current?.querySelector('input')?.focus(), 0);
        } else {
            setIsAddingAccount(false);
            setAccount(v);
        }
    };

    const handleAddAccount = () => {
        const trimmed = newAccount.trim();
        if (!trimmed) return;
        setAccount(trimmed);
        setIsAddingAccount(false);
        setNewAccount('');
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
        <Dialog
            open={!!isOpen}
            onClose={onClose}
            maxWidth="xs"
            fullWidth
            slotProps={{
                paper: {
                    sx: {
                        bgcolor: 'var(--surface-1)',
                        backgroundImage: 'none',
                        color: 'var(--ink-1)',
                        border: '1px solid var(--rule-1)',
                        borderRadius: '10px',
                        boxShadow: '0 6px 20px oklch(0 0 0 / 0.4)',
                    },
                },
            }}
        >
            <header className="flex items-center justify-between px-5 py-3 border-b border-rule-1">
                <h2 className="text-h3 text-ink-1 truncate">{title || 'Add account'}</h2>
                <IconButton
                    aria-label="Close"
                    onClick={onClose}
                    size="small"
                    sx={{ color: 'var(--ink-3)', '&:hover': { color: 'var(--ink-1)', bgcolor: 'var(--surface-2)' } }}
                >
                    <CloseIcon fontSize="small" />
                </IconButton>
            </header>

            <div className="px-5 py-5 space-y-4">
                {!isAddingAccount && hasExistingAccounts ? (
                    <div>
                        <label className="block text-meta text-ink-3 mb-1.5">Account</label>
                        <Select
                            value={account || ''}
                            onChange={handleAccountChange}
                            displayEmpty
                            fullWidth
                            size="small"
                            inputRef={selectRef}
                        >
                            <MenuItem value="" disabled>Select account</MenuItem>
                            {existingAccounts.map((acc) => (
                                <MenuItem key={acc} value={acc}>{acc}</MenuItem>
                            ))}
                            <MenuItem value="add_new_account">+ Add new account…</MenuItem>
                        </Select>
                    </div>
                ) : null}

                {(!hasExistingAccounts || isAddingAccount) ? (
                    <div ref={textFieldRef}>
                        <label className="block text-meta text-ink-3 mb-1.5">New account name</label>
                        <div className="flex items-center gap-2">
                            <TextField
                                size="small"
                                fullWidth
                                value={isAddingAccount ? newAccount : account}
                                onChange={(e) =>
                                    isAddingAccount ? setNewAccount(e.target.value) : setAccount(e.target.value)
                                }
                                onKeyDown={(e) => {
                                    if (e.key === 'Enter') {
                                        if (isAddingAccount) handleAddAccount();
                                        else handleSubmit();
                                    }
                                }}
                                placeholder="e.g. Main"
                            />
                            {isAddingAccount ? (
                                <IconButton
                                    onClick={handleAddAccount}
                                    aria-label="Confirm account name"
                                    size="small"
                                    sx={{ color: 'var(--accent)', '&:hover': { bgcolor: 'var(--surface-2)' } }}
                                >
                                    <CheckIcon fontSize="small" />
                                </IconButton>
                            ) : null}
                        </div>
                    </div>
                ) : null}

                {showRememberMe ? (
                    <FormControlLabel
                        control={
                            <Checkbox
                                checked={!!rememberMe}
                                onChange={(e) => onRememberMeChange && onRememberMeChange(e.target.checked)}
                                size="small"
                                sx={{
                                    color: 'var(--ink-3)',
                                    '&.Mui-checked': { color: 'var(--accent)' },
                                }}
                            />
                        }
                        label={<span className="text-meta text-ink-2">Keep me logged in for 30 days</span>}
                    />
                ) : null}
            </div>

            <footer className="flex justify-end gap-2 px-5 py-3 border-t border-rule-1">
                <button
                    type="button"
                    onClick={onClose}
                    className="h-8 px-3 rounded-md border border-rule-1 text-ink-2 text-meta hover:bg-surface-2 hover:text-ink-1 transition-colors duration-fast ease-out-quart"
                >
                    Cancel
                </button>
                <button
                    type="button"
                    onClick={handleSubmit}
                    className="h-8 px-3 rounded-md bg-accent text-accent-ink text-meta hover:bg-accent-strong transition-colors duration-fast ease-out-quart"
                >
                    Submit
                </button>
            </footer>
        </Dialog>
    );
};

export default AccountPromptModal;
