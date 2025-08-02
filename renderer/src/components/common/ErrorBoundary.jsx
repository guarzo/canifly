// src/components/ErrorBoundary.jsx

import React from 'react';
import { ErrorBoundary as ReactErrorBoundary } from 'react-error-boundary';
import { error as cerr } from '../../utils/logger.jsx';

function ErrorFallback({ error, resetErrorBoundary }) {
    return (
        <div className="flex flex-col items-center justify-center min-h-screen p-4">
            <h2 className="text-2xl font-bold text-red-600 mb-4">
                Something went wrong
            </h2>
            <pre className="bg-gray-100 p-4 rounded text-sm mb-4">
                {error.message}
            </pre>
            <button
                onClick={resetErrorBoundary}
                className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
            >
                Try again
            </button>
        </div>
    );
}

function ErrorBoundary({ children }) {
    return (
        <ReactErrorBoundary
            FallbackComponent={ErrorFallback}
            onError={(error, errorInfo) => {
                cerr('ErrorBoundary caught an error:', error, errorInfo);
            }}
            onReset={() => {
                // Reset any state if needed
            }}
        >
            {children}
        </ReactErrorBoundary>
    );
}

export default ErrorBoundary;
