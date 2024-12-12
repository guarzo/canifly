// src/crypto/apiRequest.js
import { toast } from 'react-toastify';

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
