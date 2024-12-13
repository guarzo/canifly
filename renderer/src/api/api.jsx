// src/api/api.js
import { toast } from 'react-toastify';

/**
 * A generic fetch function to handle various endpoints and update state.
 *
 * @param {Object} params
 * @param {string} params.backEndURL - Base URL for the backend.
 * @param {string} params.endpoint - The endpoint to fetch (e.g. "/api/app-data").
 * @param {boolean} [params.loggedOut=false] - Whether the user is logged out.
 * @param {boolean} [params.isAuthenticated=false] - Whether the user is authenticated.
 * @param {Function} [params.setIsLoading] - Setter for isLoading state.
 * @param {Function} [params.setIsRefreshing] - Setter for isRefreshing state.
 * @param {Function} [params.setIsAuthenticated] - Setter for isAuthenticated state.
 * @param {Function} [params.setAppData] - Setter for appData state.
 * @param {Object} [options={}] - Additional options.
 * @param {boolean} [options.setLoading=false] - Whether to set loading state during fetch.
 * @param {boolean} [options.setRefreshing=false] - Whether to set refreshing state during fetch.
 * @param {boolean} [options.returnSuccess=false] - Return boolean indicating success/failure.
 *
 * @returns {Promise<boolean|undefined>} If returnSuccess is true, returns boolean indicating success/failure; otherwise, undefined.
 */
export async function fetchAppEndpoint({
                                           backEndURL,
                                           endpoint,
                                           loggedOut,
                                           isAuthenticated,
                                           setIsLoading,
                                           setIsRefreshing,
                                           setIsAuthenticated,
                                           setAppData
                                       }, {
                                           setLoading = false,
                                           setRefreshing = false,
                                           returnSuccess = false
                                       } = {}) {
    console.log(`fetchAppEndpoint called with endpoint=${endpoint}, setLoading=${setLoading}, setRefreshing=${setRefreshing}, returnSuccess=${returnSuccess}`);
    console.log("State before fetch:", {
        isAuthenticated,
        loggedOut,
        appDataExists: typeof setAppData === 'function' // just to indicate we have a setter
    });

    if (loggedOut) {
        console.log("fetchAppEndpoint: loggedOut is true, returning early");
        return returnSuccess ? false : undefined;
    }

    if (setLoading && setIsLoading) {
        setIsLoading(true);
    }
    if (setRefreshing && setIsRefreshing) {
        setIsRefreshing(true);
    }

    try {
        const response = await fetch(`${backEndURL}${endpoint}`, { credentials: 'include' });

        if (response.status === 401) {
            console.log("fetchAppEndpoint: got 401 Unauthorized");
            if (setIsAuthenticated) setIsAuthenticated(false);
            if (setAppData) setAppData(null);
            toast.warning('Please log in to access your data.');
            return returnSuccess ? false : undefined;
        }

        if (response.ok) {
            const data = await response.json();
            console.log("fetchAppEndpoint: response OK");
            if (setIsAuthenticated) setIsAuthenticated(true);
            if (setAppData) setAppData(data);
            return returnSuccess ? true : undefined;
        } else {
            const errorData = await response.json();
            console.log("fetchAppEndpoint: response not OK, error:", errorData.error);
            toast.error(errorData.error || 'An unexpected error occurred.');
            return returnSuccess ? false : undefined;
        }
    } catch (error) {
        console.error('Error fetching data:', error);
        if (isAuthenticated) {
            toast.error('Failed to load data. Please try again later.');
        }
        return returnSuccess ? false : undefined;
    } finally {
        if (setLoading && setIsLoading) {
            setIsLoading(false);
        }
        if (setRefreshing && setIsRefreshing) {
            setIsRefreshing(false);
        }
    }
}
