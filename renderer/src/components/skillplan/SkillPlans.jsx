import { useEffect, useState } from 'react';
import PropTypes from 'prop-types';
import CharacterTable from './CharacterTable.jsx';
import SkillPlanTable from './SkillPlanTable.jsx';
import { toast } from 'react-toastify';
import { Typography, ToggleButtonGroup, ToggleButton, Box } from '@mui/material';
import { People as PeopleIcon, ListAlt as SkillPlansIcon } from '@mui/icons-material';

const SkillPlans = ({ identities, skillPlans, setAppData, backEndURL }) => {
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
            try {
                const response = await fetch(
                    `${backEndURL}/api/delete-skill-plan?planName=${encodeURIComponent(planName)}`,
                    {
                        method: 'DELETE',
                        credentials: 'include',
                    }
                );

                if (response.ok) {
                    toast.success(`Deleted skill plan: ${planName}`, { autoClose: 1500 });
                    // Update the state to remove the deleted skill plan
                    setAppData((prevAppData) => {
                        const updatedSkillPlans = { ...prevAppData.SkillPlans };
                        delete updatedSkillPlans[planName];
                        return { ...prevAppData, SkillPlans: updatedSkillPlans };
                    });
                } else {
                    const errorMessage = await response.text();
                    toast.error(`Failed to delete skill plan: ${errorMessage}`, { autoClose: 1500 });
                }
            } catch (error) {
                console.error('Error deleting skill plan:', error);
                toast.error('An error occurred while deleting the skill plan.', { autoClose: 1500 });
            }
        };
    }, [skillPlans, setAppData]);

    const handleViewChange = (event, newValue) => {
        if (newValue) {
            setView(newValue);
        }
    };

    return (
        <div className="bg-gray-900 min-h-screen text-teal-200 px-4 pt-16 pb-10">
            <div className="max-w-7xl mx-auto">
                <Box className="flex items-center justify-between mb-4">
                    <Typography variant="h4" sx={{ color: '#14b8a6', fontWeight: 'bold' }}>
                        Skill Plans
                    </Typography>
                    <ToggleButtonGroup
                        value={view}
                        exclusive
                        onChange={handleViewChange}
                        sx={{
                            backgroundColor: 'rgba(255,255,255,0.05)',
                            borderRadius: '9999px',
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

                <div className="space-y-8">
                    {view === 'characters' && (
                        <div className="bg-gray-800 rounded-md p-4 shadow-md">
                            <Typography variant="h5" gutterBottom sx={{ color: '#14b8a6', fontWeight: 'bold', marginBottom: '1rem' }}>
                                By Character
                            </Typography>
                            <CharacterTable identities={identities} skillPlans={skillPlans} />
                        </div>
                    )}

                    {view === 'plans' && (
                        <div className="bg-gray-800 rounded-md p-4 shadow-md">
                            <Typography variant="h5" gutterBottom sx={{ color: '#14b8a6', fontWeight: 'bold', marginBottom: '1rem' }}>
                                By Skill Plan
                            </Typography>
                            <SkillPlanTable skillPlans={skillPlans} identities={identities} />
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
};

SkillPlans.propTypes = {
    identities: PropTypes.array.isRequired,
    skillPlans: PropTypes.object.isRequired,
    setAppData: PropTypes.func.isRequired,
    backEndURL: PropTypes.string.isRequired,
};

export default SkillPlans;
