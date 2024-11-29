// SkillPlanTable.jsx
import React, { useMemo } from 'react';
import PropTypes from 'prop-types';
import TabulatorTable from './TabulatorTable';
import { calculateDaysFromToday } from './Utils';

const SkillPlanTable = ({ skillPlans, identities, tabulatorRef }) => {
    const skillPlanData = useMemo(() => {
        return Object.values(skillPlans).map((skillPlan) => {
            const qualifiedCharacters = skillPlan.QualifiedCharacters || [];
            const pendingCharacters = skillPlan.PendingCharacters || [];
            const missingCharacters = skillPlan.MissingCharacters || [];

            const qualifiedChildren = qualifiedCharacters.map((characterName) => ({
                planName: characterName,
                Status: `
                    <i class="fas fa-check-circle" style="color: green;"></i>
                `,
            }));

            const pendingChildren = pendingCharacters.map((characterName) => {
                const characterData = identities.find((identity) => identity.Character?.CharacterName === characterName)?.Character || null;
                if (!characterData) {
                    console.log(`Pending character not found: ${characterName}`);
                    return null;
                }

                const pendingFinishDate = characterData.PendingFinishDates?.[skillPlan.Name] || "";
                const daysRemaining = calculateDaysFromToday(pendingFinishDate);

                if (daysRemaining === "") {
                    console.log(`Adding pending character: ${characterName} but no pending days`);
                    return {
                        planName: characterName,
                        Status: `
                            <i class="fas fa-clock" style="color: orange;"></i>
                        `,
                    };
                }

                console.log(`Adding pending character: ${characterName} with remaining days: ${daysRemaining}`);
                return {
                    planName: characterName,
                    Status: `
                        <i class="fas fa-clock" style="color: orange;"></i> Pending ${daysRemaining ? `(${daysRemaining})` : ""}
                    `,
                };
            }).filter(Boolean);

            const missingChildren = missingCharacters.map((characterName) => ({
                planName: characterName,
                Status: `
                    <i class="fas fa-exclamation-circle" style="color: red;"></i> Missing
                `,
            }));

            return {
                id: skillPlan.Name,
                planName: skillPlan.Name,
                Status: `
                    <i class="fas fa-clipboard clipboard-icon text-teal-400 hover:text-teal-200" title="Copy Skill Plan" style="cursor: pointer; margin-right: 10px;" onclick="window.copySkillPlan('${skillPlan.Name}')"></i>
                    <i class="fas fa-trash-alt delete-icon text-red-500 hover:text-red-300" title="Delete Skill Plan" style="cursor: pointer;" onclick="window.deleteSkillPlan('${skillPlan.Name}')"></i>
                `,
                _children: [
                    ...qualifiedChildren,
                    ...pendingChildren,
                    ...missingChildren,
                ],
            };
        });
    }, [skillPlans, identities]);

    return (
        <div className="mb-8 w-full">
            <h2 className="text-2xl font-bold mb-4 text-teal-200">Skill Plan Table</h2>
            <TabulatorTable
                id="skill-plan-table"
                data={skillPlanData}
                columns={[
                    {
                        title: "Skill Plan",
                        field: "planName",
                        formatter: "html", // Use HTML for formatting
                    },
                    {
                        title: "Status",
                        field: "Status",
                        formatter: "html",
                    },
                ]}
                options={{
                    dataTree: true,
                    dataTreeStartExpanded: false,
                }}
                ref={tabulatorRef}
            />
        </div>
    );
};

SkillPlanTable.propTypes = {
    skillPlans: PropTypes.object.isRequired,
    identities: PropTypes.array.isRequired,
    tabulatorRef: PropTypes.object,
};

export default SkillPlanTable;
