// src/components/landing/LoginButton.jsx
import { useState } from 'react';
import { toast } from 'react-toastify';
import AccountPromptModal from '../common/AccountPromptModal.jsx';
import eveSsoImage from '../../assets/images/eve-sso.jpg';
import { initiateLogin } from '../../api/apiService';
import { error as cError, logger } from '../../utils/logger';
import { isDev } from '../../Config';
import { useAuth } from '../../hooks/useAuth';

const LoginButton = ({ onModalOpenChange }) => {
    const { refreshAuth } = useAuth();
    const [modalOpen, setModalOpen] = useState(false);

    const handleOpenModal = () => {
        setModalOpen(true);
        onModalOpenChange(true);
    };

    const handleCloseModal = () => {
        setModalOpen(false);
        onModalOpenChange(false);
    };

    const handleLoginSubmit = async (account) => {
        try {
            const data = await initiateLogin(account);
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
                
                // Start polling for auth status
                const pollAuth = setInterval(async () => {
                    const authenticated = await refreshAuth();
                    if (authenticated) {
                        clearInterval(pollAuth);
                        sessionStorage.removeItem('oauth_state');
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
            onModalOpenChange(false);
        }
    };

    return (
        <>
            <button
                onClick={handleOpenModal}
                aria-label="Login with Eve SSO"
                className="group relative inline-flex items-center py-4 px-8 rounded-xl bg-gradient-to-r from-teal-500 to-blue-600 hover:from-teal-400 hover:to-blue-500 transition-all duration-300 transform hover:scale-105 hover:shadow-glow-lg"
            >
                <span className="absolute inset-0 rounded-xl bg-gradient-to-r from-teal-400/20 to-blue-500/20 blur-xl group-hover:blur-2xl transition-all duration-300" />
                <img
                    src={eveSsoImage}
                    alt="Login with Eve SSO"
                    className="relative h-16 w-auto object-contain filter brightness-110"
                />
            </button>
            <AccountPromptModal
                isOpen={modalOpen}
                onClose={handleCloseModal}
                onSubmit={handleLoginSubmit}
                title="Account Name"
            />
        </>
    );
};

export default LoginButton;
