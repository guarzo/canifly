import React from 'react';
import { motion } from 'framer-motion';
import { Typography } from '@mui/material';
import helloImg from '../../assets/images/hello.png';

const LoadingScreen = ({ message = "Loading..." }) => {
    return (
        <div className="flex flex-col items-center justify-center min-h-screen relative overflow-hidden">
            {/* Animated background */}
            <div className="absolute inset-0">
                <div className="absolute inset-0 bg-gradient-radial-at-t from-eve-blue-900/20 via-gray-900 to-gray-900" />
                <motion.div
                    className="absolute inset-0"
                    animate={{
                        background: [
                            'radial-gradient(circle at 20% 50%, rgba(20, 184, 166, 0.1) 0%, transparent 50%)',
                            'radial-gradient(circle at 80% 50%, rgba(0, 118, 230, 0.1) 0%, transparent 50%)',
                            'radial-gradient(circle at 20% 50%, rgba(20, 184, 166, 0.1) 0%, transparent 50%)',
                        ],
                    }}
                    transition={{ duration: 10, repeat: Infinity }}
                />
            </div>

            {/* Content */}
            <motion.div
                className="relative z-10 flex flex-col items-center space-y-8"
                initial={{ opacity: 0, scale: 0.9 }}
                animate={{ opacity: 1, scale: 1 }}
                transition={{ duration: 0.5 }}
            >
                {/* Logo with glow effect */}
                <motion.div
                    className="relative"
                    animate={{ y: [0, -10, 0] }}
                    transition={{ duration: 2, repeat: Infinity }}
                >
                    <div className="absolute inset-0 blur-2xl bg-teal-500/30 rounded-full" />
                    <img
                        src={helloImg}
                        alt="Loading"
                        className="w-32 h-auto object-contain relative z-10"
                    />
                </motion.div>

                {/* Loading text */}
                <div className="text-center space-y-4">
                    <Typography variant="h4" className="text-gradient font-display">
                        CanIFly
                    </Typography>
                    
                    <div className="flex items-center space-x-2">
                        <Typography variant="body1" className="text-gray-400">
                            {message}
                        </Typography>
                        
                        {/* Animated dots */}
                        <div className="flex space-x-1">
                            {[0, 1, 2].map((i) => (
                                <motion.span
                                    key={i}
                                    className="w-2 h-2 bg-teal-400 rounded-full"
                                    animate={{ 
                                        opacity: [0.3, 1, 0.3],
                                        scale: [0.8, 1.2, 0.8]
                                    }}
                                    transition={{
                                        duration: 1.5,
                                        repeat: Infinity,
                                        delay: i * 0.2,
                                    }}
                                />
                            ))}
                        </div>
                    </div>
                </div>

                {/* Progress bar */}
                <div className="w-64 h-1 bg-gray-800 rounded-full overflow-hidden">
                    <motion.div
                        className="h-full bg-gradient-to-r from-teal-500 to-blue-600"
                        animate={{ x: ['-100%', '100%'] }}
                        transition={{ duration: 2, repeat: Infinity, ease: 'linear' }}
                        style={{ width: '100%' }}
                    />
                </div>
            </motion.div>
        </div>
    );
};

export default LoadingScreen;