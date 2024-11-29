import React, { useMemo, useEffect } from "react";
import PropTypes from "prop-types";
import TabulatorTable from "./TabulatorTable";
import { toast } from "react-toastify";

const MainContent = ({
                         identities = [],
                         skillPlans = {},
                         matchingSkillPlans = {},
                         matchingCharacters = {},
                     }) => {
    console.log("MainContent Props:", { identities, skillPlans, matchingSkillPlans, matchingCharacters });

    const calculateDaysFromToday = (date) => {
        if (!date) return "";
        const diffTime = new Date(date) - new Date();
        return diffTime > 0 ? `${Math.ceil(diffTime / (1000 * 60 * 60 * 24))} days` : "0 days";
    };

    useEffect(() => {
        window.copySkillPlan = (planName) => {
            const planSkills = skillPlans[planName]?.Skills || {};
            const skillText = Object.entries(planSkills)
                .map(([skill, detail]) => `${skill}: Level ${detail.Level}`)
                .join("\n");
            if (skillText) {
                navigator.clipboard.writeText(skillText)
                    .then(() => toast.success(`Copied ${Object.keys(planSkills).length} skills from ${planName}.`, { autoClose: 1500 }))
                    .catch((err) => {
                        console.error("Copy to clipboard failed:", err);
                        toast.error("Failed to copy skill plan.", { autoClose: 1500 });
                    });
            } else {
                toast.warning("No skills available to copy.", { autoClose: 1500 });
            }
        };

        window.deleteSkillPlan = async (planName) => {
            try {
                const response = await fetch(`/api/delete-skill-plan?planName=${encodeURIComponent(planName)}`, {
                    method: "DELETE",
                });

                if (response.ok) {
                    toast.success(`Deleted skill plan: ${planName}`, { autoClose: 1500 });

                    // Dynamically remove the row from Tabulator
                    const tableElement = document.getElementById("skill-plan-table");
                    if (tableElement?.tabulator) {
                        const row = tableElement.tabulator.getRow(planName);
                        if (row) {
                            tableElement.tabulator.deleteRow(planName);
                            console.log(`Deleted row with planName '${planName}' from the table.`);
                        } else {
                            console.warn(`Row with planName '${planName}' not found in Tabulator.`);
                        }
                    } else {
                        console.warn("Tabulator instance not found for #skill-plan-table");
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

    const characterData = useMemo(() => {
        return identities.map((identity) => {
            const characterID = identity.ID;
            const characterDetails = matchingCharacters[characterID]?.Character || {};

            const children = Object.keys(skillPlans).map((planName) => {
                const qualified = characterDetails.QualifiedPlans?.[planName];
                const pending = characterDetails.PendingPlans?.[planName];
                const pendingFinishDate = characterDetails.PendingFinishDates?.[planName];
                const missingSkillsForPlan = characterDetails.MissingSkills?.[planName] || {};
                const missingCount = Object.keys(missingSkillsForPlan).length;

                if (qualified) {
                    return {
                        CharacterName: planName,
                        TotalSP: `<i class="fas fa-check-circle" style="color: green;"></i>`,
                    };
                } else if (pending) {
                    const daysRemaining = calculateDaysFromToday(pendingFinishDate);
                    return {
                        CharacterName: planName,
                        TotalSP: `<i class="fas fa-clock" style="color: orange;"></i> ${daysRemaining}`,
                    };
                } else if (missingCount > 0) {
                    return {
                        CharacterName: planName,
                        TotalSP: `<i class="fas fa-exclamation-circle" style="color: red;"></i> ${missingCount} missing`,
                    };
                }

                return null;
            }).filter(Boolean);

            return {
                ...identity,
                _children: children,
            };
        });
    }, [identities, matchingCharacters, skillPlans]);

    const skillPlanData = useMemo(() => {
        return Object.values(matchingSkillPlans).map((skillPlan) => {
            const qualifiedCharacters = skillPlan.QualifiedCharacters || [];
            const pendingCharacters = skillPlan.PendingCharacters || [];

            const qualifiedChildren = qualifiedCharacters.map((characterName) => ({
                planName: characterName,
                status: `<i class="fas fa-check-circle" style="color: green;"></i> Qualified`,
            }));

            const pendingChildren = pendingCharacters.map((characterName) => {
                const characterData = matchingCharacters[characterName]?.Character || null;
                if (!characterData) return null;

                const pendingFinishDate = characterData.PendingFinishDates?.[skillPlan.Name] || "";
                const daysRemaining = calculateDaysFromToday(pendingFinishDate);

                return {
                    planName: characterName,
                    status: `<i class="fas fa-clock" style="color: orange;"></i> Pending (${daysRemaining})`,
                };
            }).filter(Boolean);

            return {
                id: skillPlan.Name, // Ensure unique ID for Tabulator
                planName: skillPlan.Name,
                _children: [...qualifiedChildren, ...pendingChildren],
            };
        });
    }, [matchingSkillPlans, matchingCharacters]);

    return (
        <main className="container mx-auto px-4 py-8 bg-gradient-to-b from-gray-800 to-gray-700 mt-16">
            <div className="bg-gray-800 text-gray-100 p-6 rounded-md shadow-md">
                {/* Character Table */}
                <div className="mb-8 w-full">
                    <h2 className="text-2xl font-bold mb-4 text-teal-200">Character Table</h2>
                    <TabulatorTable
                        id="character-table"
                        data={characterData}
                        columns={[
                            {
                                title: "Character Name",
                                field: "CharacterName",
                                formatter: "html",
                            },
                            {
                                title: "Total Skill Points",
                                field: "TotalSP",
                                formatter: (cell) => {
                                    const value = cell.getValue();
                                    if (cell.getRow().getTreeParent() === false) {
                                        return Number(value).toLocaleString(); // Format parent row as number
                                    }
                                    return value; // Child rows retain HTML
                                },
                            },
                        ]}
                        options={{
                            dataTree: true,
                            dataTreeStartExpanded: false,
                        }}
                    />
                </div>

                {/* Skill Plan Table */}
                <div className="w-full">
                    <h2 className="text-2xl font-bold mb-4 text-teal-200">Skill Plan Table</h2>
                    <TabulatorTable
                        id="skill-plan-table"
                        data={skillPlanData}
                        columns={[
                            {
                                title: "Skill Plan Name",
                                field: "planName",
                                formatter: "html",
                            },
                            {
                                title: "Actions",
                                field: "actions",
                                formatter: (cell) => {
                                    const data = cell.getRow().getData();
                                    const planName = data.planName;

                                    if (!cell.getRow().getTreeParent()) {
                                        return `
                                            <i class="fas fa-clipboard clipboard-icon text-teal-400 hover:text-teal-200"                                                 title="Copy Skill Plan" 
                                                style="cursor: pointer; margin-right: 10px;" 
                                                onclick="window.copySkillPlan('${planName}')"></i>
                                            <i class="fas fa-trash-alt delete-icon text-red-500 hover:text-red-300"                                                 title="Delete Skill Plan" 
                                                style="cursor: pointer; color: red;" 
                                                onclick="window.deleteSkillPlan('${planName}')"></i>
                                        `;
                                    }

                                    return data.status || "";
                                },
                                headerSort: false,
                                width: 200,
                            },
                        ]}
                        options={{
                            dataTree: true,
                            dataTreeStartExpanded: false,
                            key: "id", // Explicitly set unique key
                        }}
                    />
                </div>
            </div>
        </main>
    );
};

MainContent.propTypes = {
    identities: PropTypes.array.isRequired,
    skillPlans: PropTypes.object.isRequired,
    matchingSkillPlans: PropTypes.object.isRequired,
    matchingCharacters: PropTypes.object.isRequired,
};

export default MainContent;
