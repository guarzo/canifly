import React from 'react';
import { Container, Grid } from '@mui/material';
import { motion } from 'framer-motion';
import PageHeader from '../components/common/SubPageHeader';
import FuzzworksSettings from '../components/settings/FuzzworksSettings';

const Settings = () => {
  const containerVariants = {
    hidden: { opacity: 0 },
    visible: {
      opacity: 1,
      transition: {
        staggerChildren: 0.1
      }
    }
  };

  const itemVariants = {
    hidden: { opacity: 0, y: 20 },
    visible: {
      opacity: 1,
      y: 0,
      transition: { duration: 0.5 }
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900 px-4 pb-10 pt-16">
      <PageHeader title="Settings" />
      <Container maxWidth="lg" sx={{ mt: 3 }}>
        <motion.div
          variants={containerVariants}
          initial="hidden"
          animate="visible"
        >
          <Grid container spacing={3}>
            <Grid item xs={12} md={6}>
              <motion.div variants={itemVariants}>
                <FuzzworksSettings />
              </motion.div>
            </Grid>
            {/* Additional settings cards can be added here */}
          </Grid>
        </motion.div>
      </Container>
    </div>
  );
};

export default Settings;