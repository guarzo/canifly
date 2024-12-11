import  { useState } from 'react';
import LoginButton from './LoginButton';
import { Container, Box } from '@mui/material';
import PropTypes from "prop-types";

const Landing = ({backEndURL}) => {
    const [isLoginModalOpen, setIsLoginModalOpen] = useState(false);

    const handleModalOpenChange = (isOpen) => {
        setIsLoginModalOpen(isOpen);
    };

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
            {/* Only apply animate-pulse if modal is not open */}
            <Box textAlign="center" className={!isLoginModalOpen ? 'animate-pulse' : ''}>
                <LoginButton
                    onModalOpenChange={handleModalOpenChange}
                    backEndURL={backEndURL}
                />
            </Box>
        </Container>
    );
};

Landing.propTypes = {
    backEndURL: PropTypes.string.isRequired
};

export default Landing;
