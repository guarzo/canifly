// MainContent.jsx
import React, { useRef, useEffect } from 'react';
import PropTypes from 'prop-types';
import CharacterTable from './CharacterTable';
import SkillPlanTable from './SkillPlanTable';
import { toast } from 'react-toastify';

const MainContent = ({ identities, skillPlans }) => {
    const tabulatorRef = useRef(null);

    // Function to handle delete or copy actions
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
                toast.warning(`No skills available to copy in the plan: ${planName}.`, { autoClose: 1500 });
                return;
            }

            const skillText = Object.entries(planSkills)
                .map(([skill, detail]) => `${skill} ${detail.Level}`)
                .join("\n");

            navigator.clipboard.writeText(skillText)
                .then(() => {
                    toast.success(`Copied ${Object.keys(planSkills).length} skills from ${planName}.`, { autoClose: 1500 });
                })
                .catch((err) => {
                    console.error("Copy to clipboard failed:", err);
                    toast.error("Failed to copy skill plan.", { autoClose: 1500 });
                });
        };

        window.deleteSkillPlan = async (planName) => {
            try {
                if (!tabulatorRef.current) {
                    console.error("Tabulator instance is not ready yet.");
                    return;
                }

                const response = await fetch(`/api/delete-skill-plan?planName=${encodeURIComponent(planName)}`, {
                    method: "DELETE",
                });

                if (response.ok) {
                    toast.success(`Deleted skill plan: ${planName}`, { autoClose: 1500 });

                    const row = tabulatorRef.current.getRow(planName);
                    if (row) {
                        row.delete();
                        console.log(`Deleted row with planName '${planName}' from the table.`);
                    } else {
                        console.warn(`Row with planName '${planName}' not found in Tabulator.`);
                    }
                } else {
                    const errorMessage = await response.text();
                    toast.error(`Failed to delete skill plan: ${errorMessage}`, { autoClose: 1500 });
                }
            } catch (error) {
                console.error("Error deleting skill plan:", error);
                toast.error("An error occurred while deleting the skill plan.", { autoClose: 1500 });
            }
        };
    }, [skillPlans]);

    return (
        <main className="container mx-auto px-4 py-8 bg-gradient-to-b from-gray-800 to-gray-700 mt-16">
            <div className="bg-gray-800 text-gray-100 p-6 rounded-md shadow-md">
                {/* Character Table */}
                <CharacterTable
                    identities={identities}
                    skillPlans={skillPlans}
                    tabulatorRef={tabulatorRef}
                />

                {/* Skill Plan Table */}
                <SkillPlanTable
                    skillPlans={skillPlans}
                    identities={identities}
                    tabulatorRef={tabulatorRef}
                />
            </div>
        </main>
    );
};

MainContent.propTypes = {
    identities: PropTypes.arrayOf(
        PropTypes.shape({
            CharacterID: PropTypes.number.isRequired,
            CharacterName: PropTypes.string.isRequired,
            Character: PropTypes.shape({
                QualifiedPlans: PropTypes.object,
                PendingPlans: PropTypes.object,
                PendingFinishDates: PropTypes.object,
                MissingSkills: PropTypes.object,
            }),
        })
    ),
    skillPlans: PropTypes.objectOf(
        PropTypes.shape({
            Name: PropTypes.string.isRequired,
            QualifiedCharacters: PropTypes.arrayOf(PropTypes.string),
            PendingCharacters: PropTypes.arrayOf(PropTypes.string),
            Skills: PropTypes.object,
        })
    ),
};

export default MainContent;
