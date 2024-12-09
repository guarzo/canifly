// Landing.jsx

import React from 'react';
import LoginButton from './LoginButton';
import { Container, Box } from '@mui/material';

const Landing = () => {
    return (
        <Container
            maxWidth="md"
            sx={{
                minHeight: '100vh',
                display: 'flex',
                justifyContent: 'center',
                alignItems: 'center'
            }}
        >
            <Box textAlign="center" className="animate-pulse">
                <LoginButton />
            </Box>
        </Container>
    );
};

export default Landing;
