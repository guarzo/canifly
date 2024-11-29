import React from 'react';
import LoginButton from './LoginButton';

const Landing = () => {
    return (
        <div className="flex flex-col min-h-screen bg-gradient-to-b from-gray-800 to-gray-700 text-teal-200">
            <header className="bg-gradient-to-r from-gray-900 to-gray-800 text-teal-200 p-5 shadow-md">
                <div className="container mx-auto">
                    <h1 className="text-3xl font-bold tracking-wide">Can I Fly?</h1>
                </div>
            </header>
            <main className="flex-grow flex items-center justify-center">
                <div className="text-center">
                    <h2 className="text-2xl mb-6">Welcome to Our Application</h2>
                    <LoginButton />
                </div>
            </main>
            <footer className="bg-gradient-to-r from-gray-900 to-gray-800 text-teal-200 p-5 shadow-md">
                <div className="text-center">
                    <img src="/static/img/logo.png" alt="Zoolanders Logo" className="mx-auto h-10" />
                </div>
            </footer>
        </div>
    );
};

export default Landing;
