// src/components/landing/LoginButton.jsx
import { useState } from 'react';
import { toast } from 'react-toastify';
import AccountPromptModal from '../common/AccountPromptModal.jsx';
import eveSsoImage from '../../assets/images/eve-sso.jpg';
import { initiateLogin } from '../../api/apiService';
import { error as cError } from '../../utils/logger';
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
            console.log('Login response data:', data);
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
                className="inline-flex items-center py-3 px-6 rounded-md transition duration-300 hover:bg-teal-700 bg-teal-600 dark:bg-teal-500 dark:hover:bg-teal-600"
            >
                <img
                    src={eveSsoImage}
                    alt="Login with Eve SSO"
                    className="h-16 w-auto object-contain"
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
