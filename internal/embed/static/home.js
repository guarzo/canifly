const tables = {}; // Initialize tables object

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

// Function to copy missing skills to clipboard
function copyToClipboard(text, icon, skillsCount) {
    navigator.clipboard.writeText(text).then(() => {
        icon.style.color = "green";
        setTimeout(() => icon.style.color = "", 1000);
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
            status: `<i class="fas fa-check-circle" style="color: green;"></i> Qualified`
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

        return {
            planName: characterName,
            status: `
                <div style="display: flex; align-items: center;">
                    <i class="fas fa-clock" style="color: orange;"></i>
                    <span style="margin-left: 5px;">${daysRemaining} days</span>
                </div>
            `
        };
    }).filter(Boolean); // Filter out any null values

    const children = [...qualifiedChildren, ...pendingChildren];
    console.log(`Skill Plan: ${skillPlan.Name}, Children:`, children);

    const clipboardIconParent = `<i class="fas fa-clipboard clipboard-icon" style="cursor: pointer; margin-left: 5px;" data-plan-name="${skillPlan.Name}"></i>`;

    return { planName: skillPlan.Name, status: clipboardIconParent, _children: children };
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
            return { CharacterName: planName, TotalSP: '<i class="fas fa-check-circle" style="color: green;"></i>' };
        } else if (pending) {
            const daysRemaining = calculateDaysFromToday(pendingFinishDate);
            return { CharacterName: planName, TotalSP: `<div style="display: flex; align-items: center;"><i class="fas fa-clock" style="color: orange;"></i> <span style="margin-left: 5px;">${daysRemaining} days</span></div>` };
        } else if (skillMissingCount > 0) {
            const clipboardIcon = `<i class="fas fa-clipboard clipboard-icon" style="cursor: pointer; margin-left: 5px;" data-plan-name="${planName}" data-missing-skills='${JSON.stringify(missingSkillsForPlan)}'></i>`;
            return { CharacterName: planName, TotalSP: `<div style="display: flex; align-items: center;"><i class="fas fa-exclamation-circle" style="color: red;"></i> <span style="margin-left: 5px;">${skillMissingCount} missing</span>${clipboardIcon}</div>` };
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
            const missingSkillsText = event.target.getAttribute("data-missing-skills") || "";
            const skillsCount = event.target.getAttribute("data-skills-count") || 0;

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
                formatter: function (cell) {
                    return cell.getRow().getTreeParent() ? cell.getValue() : `<strong>${cell.getValue()}</strong>`;
                }
            },
            {
                title: "Total Skill Points",
                field: "TotalSP",
                headerSort: true,
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
        dataTree: true,
        dataTreeStartExpanded: false,
        columns: [
            {
                title: "Skill Plan Name",
                field: "planName",
                headerSort: true,
                formatter: function (cell) {
                    return cell.getRow().getTreeParent() ? cell.getValue() : `<strong>${cell.getValue()}</strong>`;
                }
            },
            {
                title: "Status",
                field: "status",
                headerSort: false,
                formatter: function (cell) {
                    const row = cell.getRow();
                    if (row.getTreeParent()) {
                        cell.getElement().style.textAlign = "left";
                    }
                    return cell.getValue();
                }
            }
        ]
    });

    tables["skill-plan-table"] = skillPlanTable;

    setTimeout(() => {
        resizeTabulatorTable("character-table");
        resizeTabulatorTable("skill-plan-table");
    }, 100);
});
