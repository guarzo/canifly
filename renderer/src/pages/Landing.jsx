import { useState, useEffect } from 'react';
import { motion } from 'framer-motion';
import LoginButton from '../components/landing/LoginButton.jsx';
import { Container, Box, Typography } from '@mui/material';
import GlassCard from '../components/ui/GlassCard.jsx';
import ParticleBackground from '../components/effects/ParticleBackground.jsx';

const Landing = () => {
    const [isLoginModalOpen, setIsLoginModalOpen] = useState(false);
    const [mousePosition, setMousePosition] = useState({ x: 0, y: 0 });

    useEffect(() => {
        const handleMouseMove = (e) => {
            setMousePosition({
                x: (e.clientX / window.innerWidth) * 100,
                y: (e.clientY / window.innerHeight) * 100,
            });
        };
        window.addEventListener('mousemove', handleMouseMove);
        return () => window.removeEventListener('mousemove', handleMouseMove);
    }, []);

    const handleModalOpenChange = (isOpen) => {
        setIsLoginModalOpen(isOpen);
    };

    return (
        <div className="relative min-h-screen overflow-hidden">
            {/* Particle Background */}
            <ParticleBackground />
            
            {/* Animated background gradient */}
            <div 
                className="absolute inset-0 opacity-30"
                style={{
                    background: `radial-gradient(circle at ${mousePosition.x}% ${mousePosition.y}%, rgba(20, 184, 166, 0.15), transparent 50%)`,
                }}
            />

            <Container
                maxWidth="lg"
                sx={{
                    minHeight: '100vh',
                    display: 'flex',
                    justifyContent: 'center',
                    alignItems: 'center',
                    position: 'relative',
                    zIndex: 1,
                }}
            >
                <motion.div
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ duration: 0.8 }}
                >
                    <GlassCard className="max-w-2xl mx-auto p-12 text-center">
                        <motion.div
                            initial={{ scale: 0.9, opacity: 0 }}
                            animate={{ scale: 1, opacity: 1 }}
                            transition={{ delay: 0.2, duration: 0.5 }}
                        >
                            <Typography 
                                variant="h1" 
                                className="text-gradient font-display text-6xl md:text-8xl mb-4"
                                sx={{ fontWeight: 900 }}
                            >
                                CanIFly
                            </Typography>
                        </motion.div>
                        
                        <motion.div
                            initial={{ opacity: 0 }}
                            animate={{ opacity: 1 }}
                            transition={{ delay: 0.4, duration: 0.5 }}
                        >
                            <Typography 
                                variant="h5" 
                                className="text-gray-400 mb-8 font-light"
                            >
                                EVE Online Character & Skill Management
                            </Typography>
                        </motion.div>

                        <motion.div
                            initial={{ opacity: 0, y: 10 }}
                            animate={{ opacity: 1, y: 0 }}
                            transition={{ delay: 0.6, duration: 0.5 }}
                            className="space-y-4"
                        >
                            <Typography 
                                variant="body1" 
                                className="text-gray-300 mb-8 max-w-md mx-auto"
                            >
                                Track your characters, manage skill plans, and synchronize your EVE settings across all your accounts.
                            </Typography>
                            
                            <Box className={!isLoginModalOpen ? 'animate-float' : ''}>
                                <LoginButton
                                    onModalOpenChange={handleModalOpenChange}
                                />
                            </Box>
                        </motion.div>

                        {/* Feature highlights */}
                        <motion.div
                            initial={{ opacity: 0 }}
                            animate={{ opacity: 1 }}
                            transition={{ delay: 0.8, duration: 0.5 }}
                            className="grid grid-cols-3 gap-4 mt-12"
                        >
                            {[
                                { icon: '🚀', text: 'Multi-Account' },
                                { icon: '📊', text: 'Skill Planning' },
                                { icon: '🔄', text: 'Settings Sync' },
                            ].map((feature, index) => (
                                <motion.div
                                    key={index}
                                    whileHover={{ scale: 1.05 }}
                                    className="text-center"
                                >
                                    <div className="text-3xl mb-2">{feature.icon}</div>
                                    <Typography variant="body2" className="text-gray-400">
                                        {feature.text}
                                    </Typography>
                                </motion.div>
                            ))}
                        </motion.div>
                    </GlassCard>
                </motion.div>
            </Container>
        </div>
    );
};

export default Landing;
