// src/components/Footer.jsx

import React from 'react';

const Footer = () => {
    return (
        <footer className="bg-white dark:bg-gray-800 text-gray-800 dark:text-gray-100 py-4 mt-auto">
            <div className="container mx-auto px-4 flex flex-col md:flex-row items-center justify-between">
                <div className="flex items-center mb-4 md:mb-0">
                    <img src="/images/footer-logo.jpg" alt="Logo" className="h-8 w-8 mr-2" />
                    <span className="font-semibold">Can I Fly?</span>
                </div>
                <div>
                    <span>&copy; {new Date().getFullYear()} Can I Fly. All rights reserved.</span>
                </div>
            </div>
        </footer>
    );
};

export default Footer;
