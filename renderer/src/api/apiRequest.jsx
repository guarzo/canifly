// src/api/apiRequest.js
import { toast } from 'react-toastify';
import { log, error as cerr } from '../utils/logger.jsx'
import { backEndURL} from "../Config.jsx";

export async function apiRequest(url, fetchOptions, {
    onSuccess,
    onError,
    successMessage,
    errorMessage
} = {}) {
    try {
        const response = await fetch(backEndURL+url, fetchOptions);

        let result;
        const contentType = response.headers.get('Content-Type');
        const isJSON = contentType && contentType.includes('application/json');
        if (isJSON) {
            result = await response.json();
            log(result)
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
            return result;
        } else {
            const errorMsg = result?.error || errorMessage || 'An unexpected error occurred.';
            if (response.status !== 401) {
                toast.error(errorMsg);
            }
            if (onError) {
                onError(errorMsg);
            }
        }
    } catch (error) {
        cerr('API request error:', error);
        toast.error(errorMessage || 'An error occurred during the request.');
        if (onError) {
            onError(error.message);
        }
    }
}
