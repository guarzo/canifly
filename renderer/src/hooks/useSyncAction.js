//
// Shared envelope for sync-page actions (sync, sync-all, choose dir, backup,
// reset-to-default). Owns the loading flag and the success/error toast +
// logger pattern that was duplicated 5 times in Sync.jsx.

import { useCallback, useState } from 'react';
import { toast } from 'react-toastify';
import { logger } from '../utils/logger';

export function useSyncAction() {
    const [isLoading, setIsLoading] = useState(false);

    const run = useCallback(async (operation, options = {}) => {
        const { successMessage, errorContext = 'sync action' } = options;
        try {
            setIsLoading(true);
            const result = await operation();
            if (result?.success) {
                if (successMessage) toast.success(successMessage);
                else if (result.message) toast.success(result.message);
            }
            return result;
        } catch (err) {
            logger.error(`${errorContext} failed`, err);
            const userMessage = err?.message || `${errorContext} failed`;
            toast.error(userMessage);
            throw err;
        } finally {
            setIsLoading(false);
        }
    }, []);

    return { run, isLoading };
}
