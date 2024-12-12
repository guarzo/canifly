// src/crypto/createRequestHandler.js
import { apiRequest } from './apiRequest';
import { backEndURL, isDev } from '../Config';
import { log } from './logger';
import { toast } from 'react-toastify';

/**
 * createPostHandler:
 * A helper to reduce repetitive POST request code.
 *
 * @param {string} endpoint - API endpoint (relative to backEndURL)
 * @param {Function} setAppData - A setter for appData if needed, or null if not
 * @param {string} [successMessage] - Optional success toast message
 * @param {string} [errorMessage] - Optional error toast message
 * @param {Function} [onSuccess] - Function to run on success, receives response data
 * @returns {Function} - A function that takes a `data` object to POST.
 */
export function createPostHandler({ endpoint, setAppData, successMessage, errorMessage, onSuccess }) {
    return async (data = {}) => {
        await apiRequest(`${backEndURL}${endpoint}`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data),
            credentials: 'include'
        }, {
            successMessage,
            errorMessage,
            onSuccess: onSuccess ? (result) => onSuccess(result) : undefined
        });
    };
}

/**
 * createAddCharacterHandler:
 * Special handler for adding a character since it has unique logic (redirect).
 */
export function createAddCharacterHandler() {
    return async (account) => {
        log("handleAddCharacter called with account:", account);
        await apiRequest(`${backEndURL}/api/add-character`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ account }),
            credentials: 'include'
        }, {
            errorMessage: 'An error occurred while adding character.',
            onSuccess: (data) => {
                if (data.redirectURL) {
                    if (isDev) {
                        log("Dev mode: redirecting internally to:", data.redirectURL);
                        window.location.href = data.redirectURL;
                    } else {
                        log("Production mode: opening external browser to:", data.redirectURL);
                        window.electronAPI.openExternal(data.redirectURL);
                        toast.info("Please authenticate in your browser");
                    }
                } else {
                    toast.error("No redirect URL received from server.");
                }
            }
        });
    };
}
