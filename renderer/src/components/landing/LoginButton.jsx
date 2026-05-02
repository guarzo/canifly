// src/components/landing/LoginButton.jsx
import { useState } from 'react';
import PropTypes from 'prop-types';
import { toast } from 'react-toastify';
import AccountPromptModal from '../common/AccountPromptModal.jsx';
import { initiateLogin } from '../../api/apiService';
import { error as cError, logger } from '../../utils/logger';
import { isDev } from '../../Config';
import { useAuth } from '../../hooks/useAuth';

const LoginButton = ({ onModalOpenChange }) => {
    const { refreshAuth } = useAuth();
    const [modalOpen, setModalOpen] = useState(false);
    const [rememberMe, setRememberMe] = useState(true);

    const handleOpenModal = () => {
        setModalOpen(true);
        onModalOpenChange?.(true);
    };

    const handleCloseModal = () => {
        setModalOpen(false);
        onModalOpenChange?.(false);
    };

    const handleLoginSubmit = async (account) => {
        try {
            const data = await initiateLogin(account, rememberMe);
            logger.debug('Login response data:', data);
            // Data should have {redirectURL, state} on success
            if (data && data.redirectURL && data.state) {
                // Store the state for the callback
                sessionStorage.setItem('oauth_state', data.state);
                
                if (isDev) {
                    // In development, just redirect within Electron's internal browser
                    window.location.href = data.redirectURL;
                } else {
                    // In production, open system browser
                    window.electronAPI.openExternal(data.redirectURL);
                    toast.info("Please complete the login in your browser");
                }
                
                // Start polling for auth completion using finalize-login
                let attempts = 0;
                const maxAttempts = 30; // 30 attempts * 2 seconds = 1 minute
                const pollAuth = setInterval(async () => {
                    attempts++;
                    
                    if (attempts > maxAttempts) {
                        clearInterval(pollAuth);
                        sessionStorage.removeItem('oauth_state');
                        toast.warning('Login timeout. Please try again.');
                        return;
                    }
                    
                    try {
                        // Import finalizelogin from apiService
                        const { finalizelogin } = await import('../../api/apiService');
                        logger.debug(`Polling attempt ${attempts} - calling finalize-login with state:`, data.state);
                        const result = await finalizelogin(data.state, rememberMe);
                        logger.debug(`Finalize-login response:`, result);
                        
                        if (result && result.success) {
                            logger.info('Login finalized successfully!');
                            clearInterval(pollAuth);
                            sessionStorage.removeItem('oauth_state');
                            
                            // Store the token if provided (workaround for file:// protocol)
                            if (result.token) {
                                localStorage.setItem('session_token', result.token);
                                logger.debug('Stored session token');
                            }
                            
                            // Now refresh auth to update the UI
                            await refreshAuth();
                            toast.success('Login successful!');
                        } else if (result && result.pending) {
                            logger.debug('OAuth callback pending, continuing to poll...');
                        } else if (result && result.error) {
                            logger.error('Finalize-login error:', result.error, result.message);
                            if (result.error === 'state_not_found') {
                                // State not found - might need to restart login
                                clearInterval(pollAuth);
                                toast.error('Login session expired. Please try again.');
                            }
                        }
                        // If not success, continue polling
                    } catch (error) {
                        // Log errors but continue polling
                        logger.error('Polling error:', error);
                    }
                }, 2000);
            } else {
                toast.error("No redirect URL received from server.");
            }
        } catch (error) {
            cError('Error initiating login:', error);
            // Toast is handled by apiRequest if needed
        } finally {
            setModalOpen(false);
            onModalOpenChange?.(false);
        }
    };

    return (
        <>
            <button
                onClick={handleOpenModal}
                aria-label="Log in with EVE SSO"
                className="inline-flex items-center justify-center h-10 px-4 rounded-md bg-accent text-accent-ink text-meta font-medium hover:bg-accent-strong transition-colors duration-fast ease-out-quart"
            >
                Log in with EVE SSO
            </button>
            <AccountPromptModal
                isOpen={modalOpen}
                onClose={handleCloseModal}
                onSubmit={handleLoginSubmit}
                title="Account Name"
                showRememberMe={true}
                rememberMe={rememberMe}
                onRememberMeChange={setRememberMe}
            />
        </>
    );
};

LoginButton.propTypes = {
    onModalOpenChange: PropTypes.func,
};

export default LoginButton;
