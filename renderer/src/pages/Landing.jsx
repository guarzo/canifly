import  { useState } from 'react';
import LoginButton from '../components/landing/LoginButton.jsx';
import { Container, Box } from '@mui/material';
import PropTypes from "prop-types";

const Landing = ({backEndURL, logInCallBack}) => {
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
                    logInCallBack={logInCallBack}
                />
            </Box>
        </Container>
    );
};

Landing.propTypes = {
    backEndURL: PropTypes.string.isRequired,
    logInCallBack: PropTypes.func.isRequired,
};

export default Landing;
