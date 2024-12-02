// SkillPlans.jsx
import React, { useEffect } from 'react';
import PropTypes from 'prop-types';
import CharacterTable from './CharacterTable';
import SkillPlanTable from './SkillPlanTable';
import { toast } from 'react-toastify';
import { Container, Paper } from '@mui/material';

const SkillPlans = ({ identities, skillPlans, setHomeData }) => {
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
                    `/api/delete-skill-plan?planName=${encodeURIComponent(planName)}`,
                    {
                        method: 'DELETE',
                        credentials: 'include',
                    }
                );

                if (response.ok) {
                    toast.success(`Deleted skill plan: ${planName}`, { autoClose: 1500 });
                    // Update the state to remove the deleted skill plan
                    setHomeData((prevHomeData) => {
                        const updatedSkillPlans = { ...prevHomeData.SkillPlans };
                        delete updatedSkillPlans[planName];
                        return { ...prevHomeData, SkillPlans: updatedSkillPlans };
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
    }, [skillPlans, setHomeData]);

    return (
        <Container maxWidth="lg" className="mt-20">
            <Paper elevation={3} className="p-6 bg-gray-800 text-gray-100">
                {/* Character Table */}
                <CharacterTable identities={identities} skillPlans={skillPlans} />

                {/* Skill Plan Table */}
                <SkillPlanTable skillPlans={skillPlans} identities={identities} />
            </Paper>
        </Container>
    );
};

SkillPlans.propTypes = {
    identities: PropTypes.array.isRequired,
    skillPlans: PropTypes.object.isRequired,
    setHomeData: PropTypes.func.isRequired,
};

export default SkillPlans;
