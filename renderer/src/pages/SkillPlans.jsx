// src/components/skillplan/SkillPlans.jsx

import React, { useEffect, useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import CharacterTable from '../components/skillplan/CharacterTable.jsx';
import SkillPlanTable from '../components/skillplan/SkillPlanTable.jsx';
import {Typography, ToggleButtonGroup, ToggleButton, Box} from '@mui/material';
import {
    People as PeopleIcon,
    ListAlt as SkillPlansIcon,
} from '@mui/icons-material';
import {skillPlanInstructions} from "../utils/instructions.jsx";
import PageHeader from "../components/common/SubPageHeader.jsx";
import { useAsyncOperation } from '../hooks/useAsyncOperation';
import apiService from '../api/apiService';

const SkillPlans = ({ characters, skillPlans, conversions }) => {
    const [view, setView] = useState('characters'); // 'characters' or 'plans'
    const { execute } = useAsyncOperation();
    
    // Debug: Log the props
    React.useEffect(() => {
        console.log('SkillPlans - characters:', characters);
        console.log('SkillPlans - skillPlans:', skillPlans);
        console.log('SkillPlans - conversions:', conversions);
    }, [characters, skillPlans, conversions]);
    
    const handleCopySkillPlan = async (planName, newPlanName) => {
        return execute(
            () => apiService.copySkillPlan(planName, newPlanName),
            { successMessage: 'Skill plan copied successfully' }
        );
    };
    
    const handleDeleteSkillPlan = async (planName) => {
        return execute(
            () => apiService.deleteSkillPlan(planName),
            { successMessage: 'Skill plan deleted successfully' }
        );
    };

    const handleViewChange = (event, newValue) => {
        if (newValue) {
            setView(newValue);
        }
    };

    return (
        <div className="min-h-screen px-4 pt-16 pb-10">
            <div className="max-w-7xl mx-auto">
                <PageHeader
                    title="Skill Plans"
                    instructions={skillPlanInstructions}
                    storageKey="showSkillPlanInstructions"
                />
                <Box className="flex items-center justify-between mb-6">
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
                        <Typography variant="body2" className="text-gray-400">
                            View:
                        </Typography>
                        <ToggleButtonGroup
                            value={view}
                            exclusive
                            onChange={handleViewChange}
                            sx={{
                                backgroundColor: 'rgba(255,255,255,0.05)',
                                borderRadius: '9999px',
                                padding: '2px',
                                '.MuiToggleButton-root': {
                                    textTransform: 'none',
                                    color: '#99f6e4',
                                    fontWeight: 'normal',
                                    border: 'none',
                                    borderRadius: '9999px',
                                    '&.Mui-selected': {
                                        backgroundColor: '#14b8a6 !important',
                                        color: '#ffffff !important',
                                        fontWeight: 'bold',
                                    },
                                    '&:hover': {
                                        backgroundColor: 'rgba(255,255,255,0.1)',
                                    },
                                    minWidth: '40px',
                                    minHeight: '40px',
                                },
                            }}
                        >
                            <ToggleButton value="characters" title="View Characters">
                                <PeopleIcon fontSize="small" />
                            </ToggleButton>
                            <ToggleButton value="plans" title="View Skill Plans">
                                <SkillPlansIcon fontSize="small" />
                            </ToggleButton>
                        </ToggleButtonGroup>
                    </Box>
                </Box>

                <AnimatePresence mode="wait">
                    <motion.div 
                        key={view}
                        initial={{ opacity: 0, y: 20 }}
                        animate={{ opacity: 1, y: 0 }}
                        exit={{ opacity: 0, y: -20 }}
                        transition={{ duration: 0.3 }}
                        className="space-y-8"
                    >
                        {view === 'characters' && (
                            <motion.div 
                                className="glass rounded-xl p-6"
                                initial={{ scale: 0.95 }}
                                animate={{ scale: 1 }}
                                transition={{ duration: 0.3 }}
                            >
                                <Typography
                                    variant="h5"
                                    gutterBottom
                                    className="text-gradient font-display mb-4"
                                >
                                    By Character
                                </Typography>
                                <CharacterTable characters={characters} skillPlans={skillPlans} conversions={conversions} />
                            </motion.div>
                        )}

                        {view === 'plans' && (
                            <motion.div 
                                className="glass rounded-xl p-6"
                                initial={{ scale: 0.95 }}
                                animate={{ scale: 1 }}
                                transition={{ duration: 0.3 }}
                            >
                                <Typography
                                    variant="h5"
                                    gutterBottom
                                    className="text-gradient font-display mb-4"
                                >
                                    By Skill Plan
                                </Typography>
                                <SkillPlanTable skillPlans={skillPlans} characters={characters} conversions={conversions}
                                                onCopySkillPlan={handleCopySkillPlan} onDeleteSkillPlan={handleDeleteSkillPlan} />
                            </motion.div>
                        )}
                    </motion.div>
                </AnimatePresence>
            </div>
        </div>
    );
};

export default SkillPlans;
