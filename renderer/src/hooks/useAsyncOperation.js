import { useState, useCallback } from 'react';
import { toast } from 'react-toastify';

export function useAsyncOperation() {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState(null);
  
  const execute = useCallback(async (operation, options = {}) => {
    const {
      successMessage,
      errorMessage,
      showToast = true,
      onSuccess,
      onError
    } = options;

    setIsLoading(true);
    setError(null);
    
    try {
      const result = await operation();
      
      if (showToast && successMessage) {
        toast.success(successMessage);
      } else if (showToast && result?.message) {
        toast.success(result.message);
      }
      
      if (onSuccess) {
        onSuccess(result);
      }
      
      return result;
    } catch (err) {
      setError(err);
      console.error('Operation failed:', err);
      
      if (showToast) {
        const message = errorMessage || err.message || 'Operation failed';
        toast.error(message);
      }
      
      if (onError) {
        onError(err);
      }
      
      throw err;
    } finally {
      setIsLoading(false);
    }
  }, []);
  
  const reset = useCallback(() => {
    setIsLoading(false);
    setError(null);
  }, []);
  
  return { execute, isLoading, error, reset };
}