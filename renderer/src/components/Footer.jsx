import React from 'react';

const Footer = () => {
    return (
        <footer className="bg-gradient-to-r from-gray-900 to-gray-800 text-teal-200 py-4 mt-auto shadow-inner border-t-4 border-teal-500">
            <div className="container mx-auto px-4 flex flex-col items-center justify-center">
                <img src="/images/footer-logo.jpg" alt="Logo" className="h-8 w-8 mb-2" />
                <span className="text-sm">&copy; {new Date().getFullYear()} Can I Fly? All rights reserved.</span>
            </div>
        </footer>
    );
};

export default Footer;
