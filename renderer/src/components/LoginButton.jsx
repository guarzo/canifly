import { useState } from 'react';
import { toast } from 'react-toastify';
import AccountPromptModal from './AccountPromptModal';
import PropTypes from 'prop-types';

const LoginButton = ({ onModalOpenChange }) => {
    const [modalOpen, setModalOpen] = useState(false);

    const handleOpenModal = () => {
        setModalOpen(true);
        onModalOpenChange(true); // Inform parent that modal is now open
    };

    const handleCloseModal = () => {
        setModalOpen(false);
        onModalOpenChange(false); // Inform parent that modal is closed
    };

    const handleLoginSubmit = async (account) => {
        try {
            const response = await fetch('/api/login', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ account }),
                credentials: 'include',
            });
            if (response.ok) {
                const data = await response.json();
                if (data.redirectURL) {
                    window.location.href = data.redirectURL;
                } else {
                    toast.error("No redirect URL received from server.");
                }
            } else {
                const errorText = await response.text();
                toast.error(`Failed to initiate login: ${errorText}`);
            }
        } catch (error) {
            console.error('Error initiating login:', error);
            toast.error('An error occurred while initiating login.');
        } finally {
            setModalOpen(false);
            onModalOpenChange(false); // Modal closed after submit
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
                    src="/images/eve-sso.jpg"
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

LoginButton.propTypes = {
    onModalOpenChange: PropTypes.func.isRequired,
};

export default LoginButton;
