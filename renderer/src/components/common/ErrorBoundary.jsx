// src/components/ErrorBoundary.jsx

import React from 'react';
import { ErrorBoundary as ReactErrorBoundary } from 'react-error-boundary';
import { error as cerr } from '../../utils/logger.jsx';

function ErrorFallback({ error, resetErrorBoundary }) {
    return (
        <div className="flex flex-col items-center justify-center min-h-screen bg-surface-0 p-4">
            <h2 className="text-h2 text-status-error mb-4">
                Something went wrong
            </h2>
            <pre className="bg-surface-1 border border-rule-1 rounded-md p-4 font-mono text-meta text-ink-2 mb-4 max-w-[65ch] whitespace-pre-wrap">
                {error.message}
            </pre>
            <button
                onClick={resetErrorBoundary}
                className="h-8 px-3 rounded-md bg-accent text-accent-ink text-meta font-medium hover:bg-accent-strong"
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
