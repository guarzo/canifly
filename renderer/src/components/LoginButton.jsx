import { useState } from 'react';
import { toast } from 'react-toastify';
import AccountPromptModal from './AccountPromptModal';
import PropTypes from 'prop-types';
import eveSsoImage from '../assets/images/eve-sso.jpg';

const LoginButton = ({ onModalOpenChange, backEndURL }) => {
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
            const response = await fetch(`${backEndURL}/api/login`, {
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

LoginButton.propTypes = {
    onModalOpenChange: PropTypes.func.isRequired,
    backEndURL: PropTypes.string.isRequired
};

export default LoginButton;
