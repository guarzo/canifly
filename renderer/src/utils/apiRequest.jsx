// src/utils/apiRequest.js
import { toast } from 'react-toastify';

/**
 * A helper function to reduce repetitive fetch and toast logic.
 *
 * @param {string} url - The full endpoint URL (including backEndURL).
 * @param {object} fetchOptions - Options for fetch (method, headers, body, credentials).
 * @param {object} [config={}] - Additional config for handling success and error states.
 * @param {Function} [config.onSuccess] - Called with the parsed response JSON if request succeeds.
 * @param {Function} [config.onError] - Called if request fails (receives error message).
 * @param {string} [config.successMessage] - Toast message on success.
 * @param {string} [config.errorMessage] - Default toast message on error if none returned by server.
 *
 * @returns {Promise<void>}
 */
export async function apiRequest(url, fetchOptions, {
    onSuccess,
    onError,
    successMessage,
    errorMessage
} = {}) {
    try {
        const response = await fetch(url, fetchOptions);

        let result;
        const contentType = response.headers.get('Content-Type');
        const isJSON = contentType && contentType.includes('application/json');
        // Attempt to parse JSON if available
        if (isJSON) {
            result = await response.json();
        } else {
            result = await response.text();
        }

        if (response.ok) {
            if (successMessage) {
                toast.success(successMessage);
            }
            if (onSuccess) {
                onSuccess(result);
            }
        } else {
            const errorMsg = result?.error || errorMessage || 'An unexpected error occurred.';
            toast.error(errorMsg);
            if (onError) {
                onError(errorMsg);
            }
        }
    } catch (error) {
        console.error('API request error:', error);
        toast.error(errorMessage || 'An error occurred during the request.');
        if (onError) {
            onError(error.message);
        }
    }
}
