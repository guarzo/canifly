// AccountPromptModal.jsx
import React, { useState } from 'react';
import PropTypes from 'prop-types';

const AccountPromptModal = ({ isOpen, onClose, onSubmit, title }) => {
    const [account, setAccount] = useState('');

    const handleSubmit = () => {
        if (!account) {
            // Optionally show a toast error
            return;
        }
        onSubmit(account);
        setAccount(''); // reset
    };

    if (!isOpen) {
        return null;
    }

    return (
        <div className="fixed inset-0 flex items-center justify-center bg-black bg-opacity-50 z-50">
            <div className="bg-gray-800 text-teal-200 p-6 rounded shadow-md w-80">
                <h2 className="mb-4 text-lg font-semibold">{title || 'Enter Account Name'}</h2>
                <input
                    type="text"
                    className="w-full px-3 py-2 mb-4 border border-gray-600 rounded-md bg-gray-700 text-teal-200"
                    placeholder="Account Name"
                    value={account}
                    onChange={(e) => setAccount(e.target.value)}
                />
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
};

export default AccountPromptModal;
