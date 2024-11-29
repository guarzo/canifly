// src/components/LoginButton.jsx

import React from 'react';

/**
 * LoginButton Component
 *
 * Renders a button that redirects the user to the /login endpoint with a unique query parameter.
 */
const LoginButton = () => {
    // Generate a unique query parameter to prevent caching
    // const loginUrl = `/login?ts=${Date.now()}`;
    const loginUrl = `http://localhost:8713/login?ts=${Date.now()}`

    return (
        <a
            href={loginUrl}
            aria-label="Login with Eve SSO"
            className="inline-flex items-center bg-teal-600 text-white py-3 px-6 rounded-md text-lg font-semibold transition duration-300 hover:bg-teal-700 dark:bg-teal-500 dark:hover:bg-teal-600"
        >
            <i className="fas fa-sign-in-alt mr-2"></i> Login with Eve SSO
        </a>
    );
};

export default LoginButton;
