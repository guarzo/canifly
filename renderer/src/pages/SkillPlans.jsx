// src/components/skillplan/SkillPlans.jsx

import React, { useEffect, useState } from 'react';
import PropTypes from 'prop-types';
import CharacterTable from '../components/skillplan/CharacterTable.jsx';
import SkillPlanTable from '../components/skillplan/SkillPlanTable.jsx';
import { toast } from 'react-toastify';
import {Typography, ToggleButtonGroup, ToggleButton, Box} from '@mui/material';
import {
    People as PeopleIcon,
    ListAlt as SkillPlansIcon,
} from '@mui/icons-material';
import { deleteSkillPlan as deleteSkillPlanApi } from '../api/apiService.jsx';
import {skillPlanInstructions} from "../utils/instructions.jsx";
import PageHeader from "../components/common/SubPageHeader.jsx";

const SkillPlans = ({ characters, skillPlans, setAppData, backEndURL, conversions }) => {
    const [view, setView] = useState('characters'); // 'characters' or 'plans'

    useEffect(() => {
        window.copySkillPlan = (planName) => {
            const plan = skillPlans[planName];
            if (!plan) {
                console.error(`Skill plan not found: ${planName}`);
                toast.warning(`Skill plan not found: ${planName}`, { autoClose: 1500 });
                return;
            }

            const planSkills = plan.Skills || {};

            if (Object.keys(planSkills).length === 0) {
                console.warn(`No skills available to copy in the plan: ${planName}`);
                toast.warning(`No skills available to copy in the plan: ${planName}.`, {
                    autoClose: 1500,
                });
                return;
            }

            const skillText = Object.entries(planSkills)
                .map(([skill, detail]) => `${skill} ${detail.Level}`)
                .join('\n');

            navigator.clipboard
                .writeText(skillText)
                .then(() => {
                    toast.success(`Copied ${Object.keys(planSkills).length} skills from ${planName}.`, {
                        autoClose: 1500,
                    });
                })
                .catch((err) => {
                    console.error('Copy to clipboard failed:', err);
                    toast.error('Failed to copy skill plan.', { autoClose: 1500 });
                });
        };

        window.deleteSkillPlan = async (planName) => {
            const result = await deleteSkillPlanApi(planName, backEndURL);
            if (result && result.success) {
                toast.success(`Deleted skill plan: ${planName}`, { autoClose: 1500 });
                // Update the state to remove the deleted skill plan
                setAppData((prevAppData) => {
                    const updatedSkillPlans = { ...prevAppData.EveData.SkillPlans };
                    delete updatedSkillPlans[planName];

                    return {
                        ...prevAppData,
                        EveData: {
                            ...prevAppData.EveData,
                            SkillPlans: updatedSkillPlans
                        }
                    };
                });
            }
        };
    }, [skillPlans, setAppData, backEndURL]);

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
                            <SkillPlanTable skillPlans={skillPlans} characters={characters} conversions={conversions} />
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
    setAppData: PropTypes.func.isRequired,
    backEndURL: PropTypes.string.isRequired,
    conversions: PropTypes.object.isRequired,
};

export default SkillPlans;
