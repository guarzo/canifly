// Landing.jsx

import React from 'react';
import LoginButton from './LoginButton';
import { Typography, Container, Box } from '@mui/material';

const Landing = () => {
    return (
        <Container maxWidth="md" style={{ marginTop: '80px' }}>
            <Box textAlign="center" py={5}>
                <Typography variant="h3" gutterBottom>
                    Welcome to Can I Fly?
                </Typography>
                <Typography variant="h6" gutterBottom>
                    Please log in to continue.
                </Typography>
                <LoginButton />
            </Box>
        </Container>
    );
};

export default Landing;
