// LoginButton.jsx
import React, { useState } from 'react';
import { toast } from 'react-toastify';
import AccountPromptModal from './AccountPromptModal';

const LoginButton = () => {
    const [modalOpen, setModalOpen] = useState(false);

    const handleOpenModal = () => {
        setModalOpen(true);
    };

    const handleCloseModal = () => {
        setModalOpen(false);
    };

    const handleLoginSubmit = async (account) => {
        try {
            const response = await fetch('/api/login', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ account }),
                credentials: 'include',
            })
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
        }
    };

    return (
        <>
            <button
                onClick={handleOpenModal}
                aria-label="Login with Eve SSO"
                className="inline-flex items-center bg-teal-600 text-white py-3 px-6 rounded-md text-lg font-semibold transition duration-300 hover:bg-teal-700 dark:bg-teal-500 dark:hover:bg-teal-600"
            >
                <i className="fas fa-sign-in-alt mr-2"></i> Login with Eve SSO
            </button>
            <AccountPromptModal
                isOpen={modalOpen}
                onClose={handleCloseModal}
                onSubmit={handleLoginSubmit}
                title="Login - Enter Account Name"
            />
        </>
    );
};

export default LoginButton;
