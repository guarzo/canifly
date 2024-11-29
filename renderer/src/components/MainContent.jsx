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
        const finishDate = new Date(date);
        const currentDate = new Date();
        const diffTime = finishDate - currentDate;
        return diffTime > 0 ? `${Math.ceil(diffTime / (1000 * 60 * 60 * 24))}` : "0";
    };

    const transformCharacterData = (character) => {
        const children = Object.keys(skillPlans).map((planName) => {
            const missingSkillsCount = Object.keys(character.MissingSkills?.[planName] || {}).length;
            if (character.QualifiedPlans?.[planName]) {
                return { CharacterName: planName, TotalSP: "✔️" };
            } else if (character.PendingPlans?.[planName]) {
                const daysRemaining = calculateDaysFromToday(character.PendingFinishDates?.[planName]);
                return { CharacterName: planName, TotalSP: `${daysRemaining} days` };
            } else if (missingSkillsCount > 0) {
                return { CharacterName: planName, TotalSP: `${missingSkillsCount} missing` };
            }
            return null;
        }).filter(Boolean);

        console.log(`Transformed character data: ${character.CharacterName}`, children);
        return { ...character, _children: children };
    };

    const transformSkillPlanData = (skillPlan) => {
        const children = (skillPlan.QualifiedCharacters || []).map((name) => ({
            planName: name,
            isQualified: true,
        }));
        console.log(`Transformed skill plan data: ${skillPlan.Name}`, children);
        return { planName: skillPlan.Name, _children: children };
    };

    const characterColumns = useMemo(() => [
        { title: "Name", field: "CharacterName", sorter: "string", headerFilter: "input" },
        { title: "Total Skill Points", field: "TotalSP" },
    ], []);

    const skillPlanColumns = useMemo(() => [
        { title: "Skill Plan Name", field: "planName", sorter: "string", headerFilter: "input" },
        {
            title: "Actions",
            formatter: (cell) => {
                const { planName } = cell.getRow().getData();
                return `
                    <button onclick="window.handleCopySkillPlan('${planName}')" class="copy-skill-plan p-2 bg-blue-500 text-white rounded hover:bg-blue-600">
                        <i class="fas fa-copy"></i>
                    </button>
                    <button onclick="window.handleDeleteSkillPlan('${planName}')" class="delete-skill-plan p-2 bg-red-500 text-white rounded hover:bg-red-600">
                        <i class="fas fa-trash-alt"></i>
                    </button>
                `;
            },
        },
    ], []);

    const characterData = useMemo(() => identities.map(transformCharacterData), [identities, skillPlans]);
    const skillPlanData = useMemo(() => Object.values(skillPlans).map(transformSkillPlanData), [skillPlans]);

    useEffect(() => {
        window.handleCopySkillPlan = async (planName) => {
            try {
                const skillPlan = skillPlans[planName];
                if (!skillPlan) {
                    toast.error(`Skill plan "${planName}" not found.`);
                    return;
                }
                const planDetails = Object.entries(skillPlan.Skills || {})
                    .map(([skill, detail]) => `${skill}: Level ${detail.Level}`)
                    .join("\n");
                await navigator.clipboard.writeText(planDetails);
                toast.success(`Copied skill plan "${planName}" to clipboard.`);
            } catch {
                toast.error("Failed to copy skill plan.");
            }
        };

        window.handleDeleteSkillPlan = async (planName) => {
            try {
                const response = await fetch(`/api/delete-skill-plan?planName=${encodeURIComponent(planName)}`, {
                    method: "DELETE",
                });
                if (response.ok) {
                    toast.success(`Skill plan "${planName}" deleted.`);
                    const tableEvent = new CustomEvent("updateSkillPlanTable", { detail: { planName } });
                    document.dispatchEvent(tableEvent);
                } else {
                    toast.error("Failed to delete skill plan.");
                }
            } catch (error) {
                console.error("Error deleting skill plan:", error);
                toast.error("An error occurred.");
            }
        };

        return () => {
            delete window.handleCopySkillPlan;
            delete window.handleDeleteSkillPlan;
        };
    }, [skillPlans]);

    return (
        <main className="container mx-auto px-4 py-8 bg-gray-50 dark:bg-gray-800 mt-16">
            <div className="bg-white dark:bg-gray-700 text-gray-800 dark:text-gray-100 p-6 rounded-md shadow-md">
                <div className="w-full mb-8">
                    <h2 className="text-2xl font-bold mb-4">Character Table</h2>
                    <TabulatorTable
                        data={characterData}
                        columns={characterColumns}
                        options={{
                            dataTree: true,
                            dataTreeStartExpanded: false,
                            layout: "fitColumns",
                        }}
                    />
                </div>
                <div className="w-full">
                    <h2 className="text-2xl font-bold mb-4">Skill Plan Table</h2>
                    <TabulatorTable
                        data={skillPlanData}
                        columns={skillPlanColumns}
                        options={{
                            dataTree: true,
                            layout: "fitColumns",
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
