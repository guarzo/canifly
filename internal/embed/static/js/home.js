// home.js

// Initialize tables object
const tables = {};

// Configure Toastr options as a constant
const TOASTR_OPTIONS = {
    closeButton: true,
    debug: false,
    newestOnTop: true,
    progressBar: true,
    positionClass: "toast-top-right",
    preventDuplicates: false,
    showDuration: 300,
    hideDuration: 1000,
    timeOut: 1500,
    extendedTimeOut: 1000,
    showEasing: "swing",
    hideEasing: "linear",
    showMethod: "fadeIn",
    hideMethod: "fadeOut"
};

toastr.options = TOASTR_OPTIONS;

// Helper function to calculate days remaining
function calculateDaysFromToday(date) {
    if (!date) return "";  // No display for N/A or null dates
    const finishDate = new Date(date);
    const currentDate = new Date();
    const diffTime = finishDate - currentDate;
    return diffTime > 0 ? `${Math.ceil(diffTime / (1000 * 60 * 60 * 24))}` : "0";
}

// Function to copy text to clipboard
function copyToClipboard(text, icon, skillsCount) {
    navigator.clipboard.writeText(text).then(() => {
        // Add visual indicator
        icon.classList.add('active');
        setTimeout(() => icon.classList.remove('active'), 1000);
        toastr.success(`${skillsCount} skills copied!`);
    }).catch((err) => {
        console.error("Failed to copy to clipboard: ", err);
        toastr.error("Failed to copy skills to clipboard.");
    });
}

// Function to transform skill plan data into hierarchical format
function transformSkillPlanData(skillPlan) {
    const qualifiedCharacters = skillPlan.QualifiedCharacters || [];
    const pendingCharacters = skillPlan.PendingCharacters || [];

    console.log("Processing Skill Plan:", skillPlan);

    function getCharacterData(characterName) {
        const characterEntry = Object.values(MatchingCharacters).find(
            (entry) => entry.Character.CharacterName === characterName
        );
        return characterEntry?.Character || null;
    }

    const qualifiedChildren = qualifiedCharacters.map(characterName => {
        return {
            planName: characterName,
            missingSkillsText: '',
            skillsCount: 0,
            characterName: characterName,
            daysRemaining: '',
            isQualified: true
        };
    });

    const pendingChildren = pendingCharacters.map(characterName => {
        const characterData = getCharacterData(characterName);
        if (!characterData) {
            console.warn(`No data found for pending character: ${characterName}`);
            return null; // Skip if no character data is available
        }

        const pendingFinishDate = characterData.PendingFinishDates?.[skillPlan.Name] || "";
        const daysRemaining = calculateDaysFromToday(pendingFinishDate);

        const missingSkillsForPlan = characterData.MissingSkills?.[skillPlan.Name] || {};
        const skillsCount = Object.keys(missingSkillsForPlan).length;

        return {
            planName: characterName,
            missingSkills: missingSkillsForPlan, // Store missing skills object
            skillsCount: skillsCount,
            characterName: characterName,
            daysRemaining: daysRemaining,
            isQualified: false
        };
    }).filter(Boolean); // Filter out any null values

    const children = [...qualifiedChildren, ...pendingChildren];
    console.log(`Skill Plan: ${skillPlan.Name}, Children:`, children);

    return {
        planName: skillPlan.Name,
        _children: children
    };
}

// Function to transform character data into hierarchical format
function transformCharacterData(character) {
    const children = Object.keys(TabulatorSkillPlans).map(planName => {
        const qualified = character.QualifiedPlans?.[planName];
        const pending = character.PendingPlans?.[planName];
        const pendingFinishDate = character.PendingFinishDates?.[planName];
        const missingSkillsForPlan = character.MissingSkills?.[planName] || {};
        const skillMissingCount = Object.keys(missingSkillsForPlan).length;

        if (qualified) {
            return {
                CharacterName: planName,
                TotalSP: '<i class="fas fa-check-circle" style="color: green;"></i>'
            };
        } else if (pending) {
            const daysRemaining = calculateDaysFromToday(pendingFinishDate);
            const missingSkillsJson = JSON.stringify(missingSkillsForPlan);

            // Add clipboard icon for copying missing skills
            const clipboardIcon = `<i class="fas fa-clipboard clipboard-icon" 
                title="Copy Missing Skills" style="cursor: pointer; margin-left: 5px;" 
                data-plan-name="${planName}" data-character-name="${character.CharacterName}" 
                data-missing-skills='${missingSkillsJson}' data-skills-count="${skillMissingCount}"></i>`;

            return {
                CharacterName: planName,
                TotalSP: `<div style="display: flex; align-items: center;">
                            <i class="fas fa-clock" style="color: orange;"></i>
                            <span style="margin-left: 5px;">${daysRemaining} days</span>
                            ${clipboardIcon}
                          </div>`
            };
        } else if (skillMissingCount > 0) {
            // For characters missing skills (not pending or qualified)
            const missingSkillsJson = JSON.stringify(missingSkillsForPlan);

            // Add clipboard icon for copying missing skills
            const clipboardIcon = `<i class="fas fa-clipboard clipboard-icon" 
                title="Copy Missing Skills" style="cursor: pointer; margin-left: 5px;" 
                data-plan-name="${planName}" data-character-name="${character.CharacterName}" 
                data-missing-skills='${missingSkillsJson}' data-skills-count="${skillMissingCount}"></i>`;

            return {
                CharacterName: planName,
                TotalSP: `<div style="display: flex; align-items: center;">
                            <i class="fas fa-exclamation-circle" style="color: red;"></i>
                            <span style="margin-left: 5px;">${skillMissingCount} missing</span>
                            ${clipboardIcon}
                          </div>`
            };
        }
        return null;
    }).filter(Boolean);

    return { ...character, _children: children };
}


function copySkillPlanToClipboard(planName, icon) {
    const planSkills = TabulatorSkillPlans[planName]?.Skills || {};
    const skillText = Object.entries(planSkills).map(([skill, detail]) => `${skill}: Level ${detail.Level}`).join('\n');
    console.log(`Copying skill plan for ${planName}:`, skillText);
    if (skillText.length > 0) {
        copyToClipboard(skillText, icon, Object.keys(planSkills).length);
    } else {
        console.error(`No skills found to copy for ${planName}`);
        toastr.error("No skills available to copy.");
    }
}

// Define the clipboard click handler
function clipboardClickHandler(event) {
    if (event.target.classList.contains('clipboard-icon')) {
        event.stopPropagation();

        const planName = event.target.getAttribute("data-plan-name");
        const characterName = event.target.getAttribute("data-character-name");

        if (characterName) {
            const missingSkills = JSON.parse(event.target.getAttribute("data-missing-skills")) || {};
            const skillsCount = Object.keys(missingSkills).length;

            console.log('Missing Skills:', missingSkills);

            const missingSkillsText = Object.entries(missingSkills)
                .map(([skillName, skillData]) => {
                    console.log(`Skill Data for ${skillName}:`, skillData);
                    // Adjust access based on actual data structure
                    const requiredLevel = skillData.Level || skillData || 'N/A';
                    return `${skillName}: Level ${requiredLevel}`;
                })
                .join('\n');

            console.log(`Copying missing skills for ${characterName} in ${planName}:`, missingSkillsText);
            copyToClipboard(missingSkillsText, event.target, skillsCount);
        } else {
            copySkillPlanToClipboard(planName, event.target);
        }
    }
}


// Attach the clipboard click handler to the document
document.addEventListener('click', clipboardClickHandler);

document.addEventListener('DOMContentLoaded', function () {
    const addSkillPlanButton = document.getElementById('add-skill-plan-button');
    const addSkillPlanModal = document.getElementById('add-skill-plan-modal');
    const closeModalButton = document.getElementById('close-modal');
    const saveSkillPlanButton = document.getElementById('save-skill-plan-button');

    // Open modal on button click
    addSkillPlanButton.addEventListener('click', function () {
        addSkillPlanModal.style.display = 'block';
    });

    // Close modal on clicking the close button
    closeModalButton.addEventListener('click', function () {
        addSkillPlanModal.style.display = 'none';
    });

    // Save the skill plan on button click
    saveSkillPlanButton.addEventListener('click', async () => {
        const skillPlanName = document.getElementById("skill-plan-name").value.trim();
        const skillPlanContents = document.getElementById("skill-plan-contents").value.trim();

        if (skillPlanName && skillPlanContents) {
            try {
                const response = await fetch('/save-skill-plan', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ name: skillPlanName, contents: skillPlanContents })
                });

                if (response.ok) {
                    toastr.success("Skill Plan Saved!");
                    addSkillPlanModal.style.display = "none";
                    window.location.reload();
                } else {
                    toastr.error("Error saving skill plan.");
                }
            } catch (error) {
                console.error("Error:", error);
                toastr.error("Failed to save skill plan.");
            }
        } else {
            toastr.warning("Please fill in both fields.");
        }
    });

    // Close modal if clicking outside of it
    window.onclick = function(event) {
        if (event.target === addSkillPlanModal) {
            addSkillPlanModal.style.display = "none";
        }
    };
});

document.addEventListener('DOMContentLoaded', function () {
    console.log("TabulatorIdentities:", TabulatorIdentities);
    console.log("TabulatorSkillPlans:", TabulatorSkillPlans);
    console.log("MatchingSkillPlans:", MatchingSkillPlans);
    console.log("MatchingCharacters:", MatchingCharacters);

    function resizeTabulatorTable(tableId) {
        const tableInstance = tables[tableId];
        const tableElement = document.getElementById(tableId);
        if (tableInstance && tableElement && tableElement.offsetParent !== null) {
            tableInstance.redraw(true);
        } else {
            console.warn(`Table with ID '${tableId}' is not visible or not found.`);
        }
    }

    const characterData = TabulatorIdentities.map(identity => {
        const characterID = identity.ID;
        const characterDetails = MatchingCharacters[characterID]?.Character || {};
        return transformCharacterData({ ...identity, ...characterDetails });
    });

    const skillPlanData = Object.values(MatchingSkillPlans).map(skillPlan => {
        const transformedData = transformSkillPlanData(skillPlan);
        console.log("Transformed Skill Plan Data:", transformedData);
        return transformedData;
    });

    console.log("Final Skill Plan Data Structure:", skillPlanData);

    const characterTable = new Tabulator("#character-table", {
        data: characterData,
        layout: "fitColumns",
        dataTree: true,
        dataTreeStartExpanded: false,
        columns: [
            {
                title: "Name",
                field: "CharacterName",
                headerSort: true,
                headerFilter: "input",
                sorter: "string",
                formatter: function (cell) {
                    return cell.getRow().getTreeParent() ? cell.getValue() : `<strong>${cell.getValue()}</strong>`;
                }
            },
            {
                title: "Total Skill Points",
                field: "TotalSP",
                headerSort: true,
                sorter: "number",
                formatter: function (cell) {
                    const value = cell.getValue();
                    const row = cell.getRow();
                    if (!row.getTreeParent()) {
                        return value ? parseInt(value).toLocaleString() : "N/A";
                    } else {
                        return value || "";
                    }
                }
            }
        ]
    });

    tables["character-table"] = characterTable;

    const skillPlanTable = new Tabulator("#skill-plan-table", {
        data: skillPlanData,
        layout: "fitColumns",
        index: "planName",
        dataTree: true,
        dataTreeStartExpanded: false,
        columns: [
            {
                title: "Skill Plan Name",
                field: "planName",
                headerSort: true,
                headerFilter: "input",
                sorter: "string",
                formatter: function (cell) {
                    return cell.getRow().getTreeParent() ? cell.getValue() : `<strong>${cell.getValue()}</strong>`;
                }
            },
            {
                title: "Actions",
                formatter: function(cell, formatterParams, onRendered) {
                    const row = cell.getRow();
                    const data = row.getData();
                    const planName = data.planName || data.name;

                    if (!row.getTreeParent()) {
                        // Parent row actions
                        return `
                            <i class="fas fa-clipboard clipboard-icon" title="Copy Skill Plan" style="cursor: pointer; margin-right: 10px;" data-plan-name="${planName}"></i>
                            <i class="fas fa-trash-alt delete-icon" title="Delete Skill Plan" style="cursor: pointer; color: red;" data-plan-name="${planName}"></i>
                        `;
                    } else {
                        // Child row actions
                        const characterName = data.characterName || data.planName || data.CharacterName;
                        const missingSkills = data.missingSkills || {};
                        const skillsCount = data.skillsCount || 0;

                        if (data.isQualified) {
                            // Qualified character
                            return `
                                <div style="display: flex; align-items: center;">
                                    <i class="fas fa-check-circle" style="color: green;"></i>
                                    <span style="margin-left: 5px;">Qualified</span>
                                </div>
                            `;
                        } else if (skillsCount > 0) {
                            // Pending character with missing skills
                            const daysRemaining = data.daysRemaining || '';
                            // Inside the formatter function for the "Actions" column
                            const missingSkillsJson = JSON.stringify(missingSkills);
                            return `
                                <div style="display: flex; align-items: center;">
                                    <i class="fas fa-clock" style="color: orange;"></i>
                                    <span style="margin-left: 5px;">${daysRemaining} days</span>
                                    <i class="fas fa-clipboard clipboard-icon" title="Copy Missing Skills" style="cursor: pointer; margin-left: 10px;" data-plan-name="${planName}" data-character-name="${characterName}" data-missing-skills='${missingSkillsJson}' data-skills-count="${skillsCount}"></i>
                                </div>
                            `;

                        } else {
                            // No missing skills and not qualified
                            return '';
                        }
                    }
                },
                width: 200,
                headerSort: false,
                formatterParams: {
                    sanitize: false, // Allow HTML content
                },
            }
        ]
    });

    tables["skill-plan-table"] = skillPlanTable;

    // Resize tables after a short delay
    setTimeout(() => {
        resizeTabulatorTable("character-table");
        resizeTabulatorTable("skill-plan-table");
    }, 100);
});

// Implement focus trapping in the modal
function trapFocus(element) {
    const focusableElements = element.querySelectorAll('input, textarea, button, [tabindex]:not([tabindex="-1"])');
    const firstElement = focusableElements[0];
    const lastElement = focusableElements[focusableElements.length - 1];

    element.addEventListener('keydown', (e) => {
        if (e.key === 'Tab') {
            if (e.shiftKey) { // Shift + Tab
                if (document.activeElement === firstElement) {
                    e.preventDefault();
                    lastElement.focus();
                }
            } else { // Tab
                if (document.activeElement === lastElement) {
                    e.preventDefault();
                    firstElement.focus();
                }
            }
        }
    });
}

document.addEventListener('click', function(event) {
    if (event.target && event.target.matches('.delete-icon')) {
        const planName = event.target.getAttribute('data-plan-name');
        if (planName && confirm(`Are you sure you want to delete the skill plan "${planName}"?`)) {
            deleteSkillPlan(planName);
        }
    }
});

function deleteSkillPlan(planName) {
    fetch(`/delete-skill-plan?planName=${encodeURIComponent(planName)}`, {
        method: 'DELETE',
    })
        .then(response => {
            if (response.ok) {
                // Remove the plan from the table
                tables["skill-plan-table"].deleteRow(planName);
                toastr.success(`Skill plan "${planName}" has been deleted.`);
            } else {
                response.text().then(text => {
                    toastr.error(`Failed to delete skill plan: ${text}`);
                });
            }
        })
        .catch(error => {
            console.error('Error deleting skill plan:', error);
            toastr.error('An error occurred while deleting the skill plan.');
        });
}

// Theme Toggle Setup
document.addEventListener('DOMContentLoaded', () => {
    const htmlElement = document.documentElement;
    const themeIcon = document.getElementById('theme-icon');
    const themeToggle = document.getElementById('theme-toggle');
    const savedTheme = localStorage.getItem('theme');

    if (savedTheme === 'dark') {
        htmlElement.classList.add('dark');
        themeIcon.classList.remove('fa-sun');
        themeIcon.classList.add('fa-moon');
    } else {
        htmlElement.classList.remove('dark');
        themeIcon.classList.remove('fa-moon');
        themeIcon.classList.add('fa-sun');
    }

    // On theme toggle
    themeToggle.addEventListener('click', () => {
        htmlElement.classList.toggle('dark');

        // Save preference
        if (htmlElement.classList.contains('dark')) {
            localStorage.setItem('theme', 'dark');
            themeIcon.classList.remove('fa-sun');
            themeIcon.classList.add('fa-moon');
        } else {
            localStorage.setItem('theme', 'light');
            themeIcon.classList.remove('fa-moon');
            themeIcon.classList.add('fa-sun');
        }
    });
});
