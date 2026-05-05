// FirstRunDialog — first thing a new user sees. Calm tool tone:
// short copy, one primary action, no hero, no gradients.

import { useState } from 'react';
import PropTypes from 'prop-types';
import { Dialog, DialogContent } from '@mui/material';
import { saveEVECredentials } from '../../api/configApi';

const CALLBACK = 'http://localhost:42423/callback';

const openExternal = (url) => (e) => {
    e.preventDefault();
    if (window.electronAPI?.openExternal) window.electronAPI.openExternal(url);
    else window.open(url, '_blank', 'noopener,noreferrer');
};

const FirstRunDialog = ({ open, onComplete }) => {
    const [clientId, setClientId] = useState('');
    const [clientSecret, setClientSecret] = useState('');
    const [error, setError] = useState('');
    const [loading, setLoading] = useState(false);

    const canSubmit = clientId.trim() && clientSecret.trim() && !loading;

    const handleSubmit = async (e) => {
        e?.preventDefault();
        if (!canSubmit) {
            setError('Both Client ID and Client Secret are required.');
            return;
        }
        setLoading(true);
        setError('');
        try {
            await saveEVECredentials(clientId.trim(), clientSecret.trim());
            onComplete();
        } catch (err) {
            setError(err?.message || 'Failed to save EVE credentials.');
        } finally {
            setLoading(false);
        }
    };

    return (
        <Dialog
            open={open}
            maxWidth="sm"
            fullWidth
            disableEscapeKeyDown
            slotProps={{
                paper: {
                    sx: {
                        backgroundColor: 'var(--surface-1)',
                        border: '1px solid var(--rule-1)',
                        backgroundImage: 'none',
                    },
                },
            }}
        >
            <DialogContent sx={{ p: 0 }}>
                <form onSubmit={handleSubmit} className="px-6 py-6">
                    <h2 className="text-h2 text-ink-1">Connect EVE application</h2>
                    <p className="mt-1 text-meta text-ink-3">
                        CanIFly needs an EVE Online developer application to read character data.
                    </p>

                    <ol className="mt-5 max-w-[65ch] list-decimal space-y-1.5 pl-5 text-body text-ink-2">
                        <li>
                            Open{' '}
                            <a
                                href="https://developers.eveonline.com"
                                onClick={openExternal('https://developers.eveonline.com')}
                                className="text-accent hover:text-accent-strong"
                            >
                                developers.eveonline.com
                            </a>{' '}
                            and create a new application.
                        </li>
                        <li>
                            Set the callback URL to{' '}
                            <code className="font-mono text-meta text-ink-1 px-1 py-0.5 rounded-sm bg-surface-2 border border-rule-1">
                                {CALLBACK}
                            </code>
                            .
                        </li>
                        <li>Select the scopes for character, skill, and asset access.</li>
                        <li>Paste the Client ID and Secret below.</li>
                    </ol>

                    <div className="mt-6 space-y-4 max-w-[480px]">
                        <Field
                            label="Client ID"
                            value={clientId}
                            onChange={setClientId}
                            disabled={loading}
                            autoFocus
                        />
                        <Field
                            label="Client Secret"
                            type="password"
                            value={clientSecret}
                            onChange={setClientSecret}
                            disabled={loading}
                        />
                    </div>

                    {error ? (
                        <p
                            role="alert"
                            className="mt-4 text-meta text-status-error"
                        >
                            {error}
                        </p>
                    ) : null}

                    <div className="mt-6 flex items-center justify-between gap-4 pt-4 border-t border-rule-1">
                        <span className="text-micro text-ink-3 font-mono tabular">
                            Callback: {CALLBACK}
                        </span>
                        <button
                            type="submit"
                            disabled={!canSubmit}
                            className="h-9 px-4 rounded-md bg-accent text-accent-ink text-meta hover:bg-accent-strong disabled:opacity-50 disabled:cursor-not-allowed"
                        >
                            {loading ? 'Saving…' : 'Save and continue'}
                        </button>
                    </div>
                </form>
            </DialogContent>
        </Dialog>
    );
};

FirstRunDialog.propTypes = {
    open: PropTypes.bool.isRequired,
    onComplete: PropTypes.func.isRequired,
};

const Field = ({ label, value, onChange, type = 'text', disabled, autoFocus }) => (
    <label className="block">
        <span className="block text-meta text-ink-3 mb-1">{label}</span>
        <input
            type={type}
            value={value}
            onChange={(e) => onChange(e.target.value)}
            disabled={disabled}
            autoFocus={autoFocus}
            className="w-full h-9 px-2.5 rounded-md border border-rule-1 bg-surface-0 text-body text-ink-1 placeholder:text-ink-3 focus:border-accent focus:outline-hidden disabled:opacity-50"
        />
    </label>
);

Field.propTypes = {
    label: PropTypes.string.isRequired,
    value: PropTypes.string.isRequired,
    onChange: PropTypes.func.isRequired,
    type: PropTypes.string,
    disabled: PropTypes.bool,
    autoFocus: PropTypes.bool,
};

export default FirstRunDialog;
