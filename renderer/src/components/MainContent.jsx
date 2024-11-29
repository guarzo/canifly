import React, { useMemo, useRef, useEffect } from "react";
import PropTypes from "prop-types";
import TabulatorTable from "./TabulatorTable";
import { toast } from "react-toastify";

const MainContent = ({ identities = [], skillPlans = {} }) => {
    const tabulatorRef = useRef(null);

// Helper function to calculate days from today
    const calculateDaysFromToday = (date) => {
        console.log(date)
        if (!date) return "";
        const targetDate = new Date(date);
        const currentDate = new Date();
        const diffTime = targetDate - currentDate;
        if (diffTime <= 0) return "0 days"; // If the date has passed or is today
        return `${Math.ceil(diffTime / (1000 * 60 * 60 * 24))} days`; // Calculate remaining days
    };


    // Format number with commas (for thousands)
    const formatNumberWithCommas = (num) => {
        return num.toLocaleString(); // Using toLocaleString to format numbers with commas
    };

    // Character Table Data
    const characterData = useMemo(() => {
        return identities.map((identity) => {
            const characterDetails = identity.Character || {};
            let TotalSP = characterDetails.CharacterSkillsResponse?.total_sp || 0; // Use total_sp directly
            TotalSP = formatNumberWithCommas(TotalSP); // Format with commas

            const children = Object.keys(skillPlans).map((planName) => {
                const qualified = characterDetails.QualifiedPlans?.[planName];
                const pending = characterDetails.PendingPlans?.[planName];
                const pendingFinishDate = characterDetails.PendingFinishDates?.[planName];
                const missingSkillsForPlan = characterDetails.MissingSkills?.[planName] || {};
                const missingCount = Object.keys(missingSkillsForPlan).length;

                // Log the status computation to debug
                console.log(`Plan: ${planName}, Qualified: ${qualified}, Pending: ${pending}, Missing Count: ${missingCount}`);

                let status = '';
                if (qualified) {
                    status = `<i class="fas fa-check-circle" style="color: green;"></i> Qualified`;
                } else if (pending) {
                    const daysRemaining = calculateDaysFromToday(pendingFinishDate);
                    status = `<i class="fas fa-clock" style="color: orange;"></i> Pending${daysRemaining ? ` (${daysRemaining})` : ""}`;
                } else if (missingCount > 0) {
                    status = `<i class="fas fa-exclamation-circle" style="color: red;"></i> ${missingCount} missing`;
                }

                // Log the computed status for the row
                console.log(`Computed Status for ${planName}: ${status}`);

                return {
                    CharacterName: planName,
                    TotalSP: status,
                };
            }).filter(Boolean);

            // Log the character data being processed
            console.log(`Character Data for ${identity.CharacterName}:`, { TotalSP, children });

            return {
                ...identity,
                TotalSP,
                _children: children,
            };
        });
    }, [identities, skillPlans]);

    // Skill Plan Table Data (with parent row actions)
    const skillPlanData = useMemo(() => {
        return Object.values(skillPlans).map((skillPlan) => {
            const qualifiedCharacters = skillPlan.QualifiedCharacters || [];
            const pendingCharacters = skillPlan.PendingCharacters || [];
            const missingCharacters = skillPlan.MissingCharacters || [];

            const qualifiedChildren = qualifiedCharacters.map((characterName) => ({
                planName: characterName,
                Status: `<i class="fas fa-check-circle" style="color: green;"></i>`,
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
                        Status: `<i class="fas fa-clock" style="color: orange;"></i>`,
                    };
                }

                console.log(`Adding pending character: ${characterName} with remaining days: ${daysRemaining}`);
                return {
                    planName: characterName,
                    Status: `<i class="fas fa-clock" style="color: orange;"></i> Pending ${daysRemaining ? `(${daysRemaining})` : ""}`, // Pending icon and days remaining
                };
            }).filter(Boolean);

            const missingChildren = missingCharacters.map((characterName) => ({
                planName: characterName,
                Status: `<i class="fas fa-exclamation-circle" style="color: red;"></i> Missing`, // Missing indicator
            }));

            // Adding copy/delete icons to the parent row
            return {
                id: skillPlan.Name,
                planName: skillPlan.Name,
                Status: `
          <i class="fas fa-clipboard clipboard-icon text-teal-400 hover:text-teal-200" title="Copy Skill Plan" style="cursor: pointer; margin-right: 10px;" onclick="window.copySkillPlan('${skillPlan.Name}')"></i>
          <i class="fas fa-trash-alt delete-icon text-red-500 hover:text-red-300" title="Delete Skill Plan" style="cursor: pointer;" onclick="window.deleteSkillPlan('${skillPlan.Name}')"></i>
        `,
                _children: [
                    ...qualifiedChildren,   // Add qualified characters
                    ...pendingChildren,     // Add pending characters
                    ...missingChildren,     // Add missing characters
                ],
            };
        });
    }, [skillPlans, identities]);

    // Function to handle delete or copy actions
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
                                formatter: "html", // Ensure this is rendered as HTML
                            },
                        ]}
                        options={{
                            dataTree: true,
                            dataTreeStartExpanded: false,
                        }}
                        ref={tabulatorRef}
                    />
                </div>

                {/* Skill Plan Table */}
                <div className="mb-8 w-full">
                    <h2 className="text-2xl font-bold mb-4 text-teal-200">Skill Plan Table</h2>
                    <TabulatorTable
                        id="skill-plan-table"
                        data={skillPlanData}
                        columns={[
                            {
                                title: "Skill Plan",
                                field: "planName",
                                formatter: "html", // Render HTML for copy/delete icons in parent row
                            },
                            {
                                title: "Status",
                                field: "Status",
                                formatter: "html", // Render HTML for the status of both parent and child rows
                            },
                        ]}
                        options={{
                            dataTree: true,
                            dataTreeStartExpanded: false,
                        }}
                        ref={tabulatorRef}
                    />
                </div>
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
