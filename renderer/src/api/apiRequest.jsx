// src/api/apiRequest.js
import { toast } from 'react-toastify';
import { log, error as cerr } from '../utils/logger.jsx';
import { backEndURL } from '../Config.jsx';

export async function apiRequest(
    url,
    fetchOptions,
    {
        onSuccess,
        onError,
        successMessage,
        errorMessage,
        disableErrorToast = false, // 1) Add a new parameter with default false
    } = {}
) {
    try {
        const response = await fetch(backEndURL + url, fetchOptions);

        let result;
        const contentType = response.headers.get('Content-Type');
        const isJSON = contentType && contentType.includes('application/json');
        if (isJSON) {
            result = await response.json();
            log(result);
        } else {
            result = await response.text();
        }

        if (response.ok) {
            // If the request succeeded:
            if (successMessage) {
                toast.success(successMessage);
            }
            if (onSuccess) {
                onSuccess(result);
            }
            return result;
        } else {
            // If the request failed (4xx or 5xx status codes):
            const errorMsg =
                result?.error || errorMessage || 'An unexpected error occurred.';

            // Handle 401 specifically
            if (response.status === 401) {
                // 2) Only show the toast if disableErrorToast is false
                if (!disableErrorToast) {
                    toast.error(errorMsg);
                }

                if (onError) {
                    onError(errorMsg);
                }
            } else {
                // For other error statuses, log but still show error
                console.log('API error response:', response.status, result);
                
                if (!disableErrorToast) {
                    toast.error(errorMsg);
                }
                
                if (onError) {
                    onError(errorMsg);
                }
            }

            // Return null or throw error, depending on your preferred flow
            return null;
        }
    } catch (error) {
        // 3) Handle network errors or exceptions
        cerr('API request error:', error);

        // Again, only show the toast if disableErrorToast is false
        if (!disableErrorToast) {
            toast.error(errorMessage || 'An error occurred during the request.');
        }

        if (onError) {
            onError(error.message);
        }

        return null;
    }
}
