// src/components/skillplan/SkillPlans.jsx

import React, { useEffect, useState } from 'react';
import PropTypes from 'prop-types';
import CharacterTable from '../components/skillplan/CharacterTable.jsx';
import SkillPlanTable from '../components/skillplan/SkillPlanTable.jsx';
import {Typography, ToggleButtonGroup, ToggleButton, Box} from '@mui/material';
import {
    People as PeopleIcon,
    ListAlt as SkillPlansIcon,
} from '@mui/icons-material';
import {skillPlanInstructions} from "../utils/instructions.jsx";
import PageHeader from "../components/common/SubPageHeader.jsx";

const SkillPlans = ({ characters, skillPlans, conversions, onCopySkillPlan, onDeleteSkillPlan }) => {
    const [view, setView] = useState('characters'); // 'characters' or 'plans'

    const handleViewChange = (event, newValue) => {
        if (newValue) {
            setView(newValue);
        }
    };

    return (
        <div className="bg-gray-900 min-h-screen text-teal-200 px-4 pt-16 pb-10">
            <div className="max-w-7xl mx-auto">
                <PageHeader
                    title="Skill Plans"
                    instructions={skillPlanInstructions}
                    storageKey="showSkillPlanInstructions"
                />
                <Box className="flex items-center justify-between mb-4">
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
                        <Typography variant="body2" sx={{ color: '#99f6e4' }}>
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

                <div className="space-y-8">
                    {view === 'characters' && (
                        <div className="bg-gray-800 rounded-md p-4 shadow-md">
                            <Typography
                                variant="h5"
                                gutterBottom
                                sx={{ color: '#14b8a6', fontWeight: 'bold', marginBottom: '1rem' }}
                            >
                                By Character
                            </Typography>
                            <CharacterTable characters={characters} skillPlans={skillPlans} conversions={conversions} />
                        </div>
                    )}

                    {view === 'plans' && (
                        <div className="bg-gray-800 rounded-md p-4 shadow-md">
                            <Typography
                                variant="h5"
                                gutterBottom
                                sx={{ color: '#14b8a6', fontWeight: 'bold', marginBottom: '1rem' }}
                            >
                                By Skill Plan
                            </Typography>
                            <SkillPlanTable skillPlans={skillPlans} characters={characters} conversions={conversions}
                                            onCopySkillPlan={onCopySkillPlan} onDeleteSkillPlan={onDeleteSkillPlan} />
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
};

SkillPlans.propTypes = {
    characters: PropTypes.array.isRequired,
    skillPlans: PropTypes.object.isRequired,
    onCopySkillPlan: PropTypes.func.isRequired,
    onDeleteSkillPlan: PropTypes.func.isRequired,
    conversions: PropTypes.object.isRequired,
};

export default SkillPlans;
