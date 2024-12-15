// src/components/ErrorBoundary.jsx

import React from 'react';
import { error as cerr } from '../../utils/logger.jsx'

class ErrorBoundary extends React.Component {
    constructor(props) {
        super(props);
        this.state = { hasError: false };
    }

    static getDerivedStateFromError(error) {
        // Update state to show fallback UI
        return { hasError: true };
    }

    componentDidCatch(error, info) {
        // Log error details
        cerr("ErrorBoundary caught an error:", error, info);
    }

    render() {
        if (this.state.hasError) {
            return <h2>Something went wrong. Please try again later.</h2>;
        }

        return this.props.children;
    }
}

export default ErrorBoundary;
