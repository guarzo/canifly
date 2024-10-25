document.addEventListener('DOMContentLoaded', function () {
    console.log("TabulatorIdentities:", TabulatorIdentities);
    console.log("TabulatorSkillPlans:", TabulatorSkillPlans);
    console.log("MatchingSkillPlans:", MatchingSkillPlans);
    console.log("MatchingCharacters:", MatchingCharacters);

    // Function to adjust table height dynamically (if necessary)
    function adjustTableHeight(tableId) {
        const tableContainer = document.getElementById(tableId);
        const rowElements = Array.from(tableContainer.querySelectorAll(".tabulator-row"));
        const visibleRows = rowElements.filter(row => row.style.display !== "none").length;

        const expandedContainers = Array.from(tableContainer.querySelectorAll(".expand-container"))
            .filter(container => container.style.display === "block");
        const expandedRows = expandedContainers.reduce((count, container) => count + container.childElementCount, 0);

        const expandedCorrection = expandedRows * 14; // Adjusted correction factor
        let visibleCorrection = visibleRows * 5;

        if (visibleRows === 1) {
            visibleCorrection = -5;
        }

        if (visibleRows === 2) {
            visibleCorrection = -3;
        }

        const totalHeight = 40 + (visibleRows + expandedRows) * 36 + 18 - expandedCorrection - visibleCorrection;
        tableContainer.style.height = `${totalHeight}px`;
    }

    adjustTableHeight("character-table");
    adjustTableHeight("skill-plan-table");

    function calculateDaysFromToday(date) {
        if (!date) return "";  // No display for N/A or null dates
        const finishDate = new Date(date);
        const currentDate = new Date();
        const diffTime = finishDate - currentDate;
        return diffTime > 0 ? `${Math.ceil(diffTime / (1000 * 60 * 60 * 24))} days` : "Today";
    }

    /**
     * Copies text to the clipboard and provides visual feedback.
     * @param {string} text - The text to copy.
     * @param {HTMLElement} icon - The icon element to provide feedback.
     */
    function copyToClipboard(text, icon) {
        console.log("Attempting to copy to clipboard:", text); // Debugging log
        navigator.clipboard.writeText(text).then(() => {
            icon.classList.add("active");
            showToast("Copied to clipboard!", "success");
            setTimeout(() => icon.classList.remove("active"), 600); // Remove active class after some time
        }).catch((err) => {
            console.error("Clipboard write failed:", err);
            icon.classList.add("error");
            showToast("Failed to copy.", "error");
            setTimeout(() => icon.classList.remove("error"), 600); // Remove error class after some time
        });
    }

    /**
     * Formats missing skills into a readable string.
     * @param {Object} missingSkills - An object containing missing skills and their levels.
     * @returns {string} - Formatted string of missing skills.
     */
    function formatMissingSkills(missingSkills) {
        return Object.entries(missingSkills).map(([skill, level]) => `${skill}: Level ${level}`).join('\n');
    }

    /**
     * Displays a toast notification.
     * @param {string} message - The message to display.
     * @param {string} type - The type of the toast ("success" or "error").
     */
    function showToast(message, type = 'success') {
        const toastContainer = document.getElementById('toast-container');
        const toast = document.createElement('div');
        toast.classList.add('toast', type);
        toast.textContent = message;
        toastContainer.appendChild(toast);
        setTimeout(() => {
            toast.remove();
        }, 3000);
    }

    /**
     * Copies the full skill plan to the clipboard.
     * @param {string} planName - The name of the skill plan.
     * @param {HTMLElement} icon - The clipboard icon element for feedback.
     */
    function copySkillPlanToClipboard(planName, icon) {
        const planSkills = TabulatorSkillPlans[planName]?.Skills || {};
        const skillText = Object.entries(planSkills).map(([skill, detail]) => `${skill}: Level ${detail.Level}`).join('\n');
        console.log(`Copying skill plan for ${planName}:`, skillText);  // Debugging log
        if (skillText.length > 0) {
            copyToClipboard(skillText, icon);
        } else {
            console.error(`No skills found to copy for ${planName}`);
            showToast("No skills available to copy.", "error");
        }
    }

    /**
     * Copies the missing skills to the clipboard.
     * @param {string} planName - The name of the skill plan.
     * @param {Object} missingSkills - An object containing missing skills.
     * @param {HTMLElement} icon - The clipboard icon element for feedback.
     */
    function copyMissingSkillsToClipboard(planName, missingSkills, icon) {
        const missingSkillsText = formatMissingSkills(missingSkills);
        console.log(`Copying missing skills for ${planName}:`, missingSkillsText); // Debugging log
        if (missingSkillsText.length > 0) {
            copyToClipboard(missingSkillsText, icon);
        } else {
            console.error(`No missing skills found to copy for ${planName}`);
            showToast("No missing skills available to copy.", "error");
        }
    }

    /**
     * Initializes the Character Table using Tabulator.
     */
    const processedData = TabulatorIdentities.map(identity => ({
        CharacterName: identity.CharacterName,
        TotalSP: isNaN(Number(identity.TotalSP)) ? 0 : Number(identity.TotalSP),
        ID: identity.ID // Assuming ID is used for mapping
    }));

    const characterTable = new Tabulator("#character-table", {
        data: processedData,
        layout: "fitColumns",
        columns: [
            { 
                title: "Name", 
                field: "CharacterName", 
                headerSort: true 
            },
            { 
                title: "Total Skill Points", 
                field: "TotalSP", 
                headerSort: true, 
                sorter: "number",
                formatter: "number",
                formatterParams: {thousands: ","}
            }
        ],
        rowClick: function (e, row) {
            if (e.target.classList.contains('clipboard-icon')) {
                return;  // Prevent row click when clipboard icon is clicked
            }
            const rowData = row.getData();
            const characterID = rowData.ID;
            const characterData = MatchingCharacters[characterID]?.Character;

            if (!characterData) {
                showToast("Character data not found.", "error");
                return;
            }

            let expandContainer = row.getElement().querySelector(".expand-container");
            if (!expandContainer) {
                expandContainer = document.createElement("div");
                expandContainer.classList.add("expand-container");
                expandContainer.style.display = "block";
                row.getElement().appendChild(expandContainer);
            } else {
                expandContainer.style.display = expandContainer.style.display === 'none' ? 'block' : 'none';
            }

            const qualifiedPlans = characterData.QualifiedPlans || {};
            const pendingPlans = characterData.PendingPlans || {};
            const pendingFinishDates = characterData.PendingFinishDates || {};
            const missingSkills = characterData.MissingSkills || {};

            expandContainer.innerHTML = Object.keys(TabulatorSkillPlans).map(planName => {
                const isQualified = !!qualifiedPlans[planName];
                const isPending = !!pendingPlans[planName];
                const missingSkillsForPlan = missingSkills[planName] || {};
                const skillMissingCount = Object.keys(missingSkillsForPlan).length;

                let displayContent = `<strong>${planName}</strong>`;
                let clipboardIcon = '';

                if (isQualified) {
                    clipboardIcon = `<i class="fas fa-clipboard clipboard-icon" data-copy-type="skill-plan" data-plan-name="${planName}" title="Copy Skill Plan"></i>`;
                    return `
                        <div class="expand-row custom-expand-row">
                            <i class="fas fa-check-circle" style="color: var(--secondary-color);"></i>
                            <div>${displayContent}</div>
                            ${clipboardIcon}
                        </div>`;
                } else if (isPending) {
                    const daysRemaining = calculateDaysFromToday(pendingFinishDates[planName]);
                    clipboardIcon = `<i class="fas fa-clipboard clipboard-icon" data-copy-type="missing-skills" data-plan-name="${planName}" data-missing-skills='${JSON.stringify(missingSkillsForPlan)}' title="Copy Missing Skills"></i>`;
                    return `
                        <div class="expand-row custom-expand-row">
                            <i class="fas fa-exclamation-circle" style="color: var(--primary-color);"></i>
                            <div>${displayContent} - ${daysRemaining}</div>
                            ${clipboardIcon}
                        </div>`;
                }
                return '';
            }).join('');

            adjustTableHeight("character-table");
        }
    });

    /**
     * Event listener for clipboard icon clicks in expanded rows and skill plan table.
     */
    document.addEventListener('click', function (event) {
        if (event.target.classList.contains('clipboard-icon')) {
            event.stopPropagation();  // Prevent row click when clipboard icon is clicked
            const copyType = event.target.getAttribute("data-copy-type");
            const planName = event.target.getAttribute("data-plan-name");
            const missingSkillsAttr = event.target.getAttribute("data-missing-skills");
            let missingSkills = {};

            if (copyType === "skill-plan") {
                copySkillPlanToClipboard(planName, event.target);
            } else if (copyType === "missing-skills") {
                try {
                    missingSkills = JSON.parse(missingSkillsAttr);
                } catch (error) {
                    console.error("Invalid JSON in data-missing-skills:", error);
                    showToast("Failed to copy missing skills.", "error");
                    return;
                }
                copyMissingSkillsToClipboard(planName, missingSkills, event.target);
            }
        }
    });

    /**
     * Initializes the Skill Plan Table using Tabulator.
     */
    const skillPlanData = Object.keys(TabulatorSkillPlans).map(planName => ({
        planName
    }));

    const skillPlanTable = new Tabulator("#skill-plan-table", {
        data: skillPlanData,
        layout: "fitColumns",
        columns: [
            { 
                title: "Skill Plan Name", 
                field: "planName", 
                headerSort: true 
            },
            {
                title: "Copy",
                field: "copy", // Added field for clarity
                formatter: function(cell, formatterParams, onRendered){
                    const planName = cell.getRow().getData().planName;
                    const missingSkills = TabulatorSkillPlans[planName]?.MissingSkills || {};
                    return `
                        <i class="fas fa-clipboard clipboard-icon" 
                           style="cursor: pointer;"
                           data-copy-type="skill-plan"
                           data-plan-name="${planName}" 
                           title="Copy Skill Plan"></i>`;
                },
                headerSort: false,
                hozAlign: "center",
                width: 80, // Fixed width to prevent overflow
                tooltip: "Copy Skill Plan" // Added tooltip for better UX
            }
        ],
        rowClick: function (e, row) {
            if (e.target.classList.contains('clipboard-icon')) {
                return; // Prevent row click when clipboard icon is clicked
            }
            const expandContainer = row.getElement().querySelector(".expand-container");
            if (expandContainer) {
                expandContainer.style.display = expandContainer.style.display === 'none' ? 'block' : 'none';
                adjustTableHeight("skill-plan-table");
            }
        },
        rowFormatter: function (row) {
            const rowData = row.getData();
            const planName = rowData.planName;
            const qualifiedCharacters = MatchingSkillPlans[planName]?.QualifiedCharacters || [];
            const pendingCharacters = MatchingSkillPlans[planName]?.PendingCharacters || [];

            let expandContainer = row.getElement().querySelector(".expand-container");
            if (expandContainer) expandContainer.remove();

            expandContainer = document.createElement("div");
            expandContainer.classList.add("expand-container");
            expandContainer.style.display = "none";
            row.getElement().appendChild(expandContainer);
            expandContainer.innerHTML = [
                ...qualifiedCharacters.map(characterName => `
                    <div class="expand-row custom-expand-row"> <!-- Add custom-expand-row class -->
                        <i class="fas fa-check-circle" style="color: var(--secondary-color);"></i>
                        <div><strong>${characterName}</strong></div>
                        <i class="fas fa-clipboard clipboard-icon" 
                           data-copy-type="skill-plan"
                           data-plan-name="${planName}" 
                           title="Copy Skill Plan"></i>
                    </div>`),
                ...pendingCharacters.map(characterName => `
                    <div class="expand-row custom-expand-row"> <!-- Add custom-expand-row class -->
                        <i class="fas fa-exclamation-circle" style="color: var(--primary-color);"></i>
                        <div><strong>${characterName}</strong></div>
                        <i class="fas fa-clipboard clipboard-icon" 
                           data-copy-type="missing-skills"
                           data-plan-name="${planName}" 
                           data-missing-skills='${JSON.stringify({})}' 
                           title="Copy Missing Skills"></i>
                    </div>`)
            ].join('');
        }
    });

    /**
     * Theme Toggle Functionality
     */
    const themeToggleButton = document.getElementById('theme-toggle-button');
    const themeIcon = document.getElementById('theme-icon');

    // Check for saved theme preference
    if (localStorage.getItem('theme') === 'light') {
        document.body.classList.add('light-theme');
        themeIcon.classList.remove('fa-moon');
        themeIcon.classList.add('fa-sun');
    } else {
        themeIcon.classList.remove('fa-sun');
        themeIcon.classList.add('fa-moon');
    }

    themeToggleButton.addEventListener('click', function() {
        document.body.classList.toggle('light-theme');
        // Toggle icon between moon and sun
        if (document.body.classList.contains('light-theme')) {
            themeIcon.classList.remove('fa-moon');
            themeIcon.classList.add('fa-sun');
            localStorage.setItem('theme', 'light');
        } else {
            themeIcon.classList.remove('fa-sun');
            themeIcon.classList.add('fa-moon');
            localStorage.setItem('theme', 'dark');
        }
    });

    /**
     * Modal Logic for Adding New Skill Plan
     */
    const addSkillPlanModal = document.getElementById('add-skill-plan-modal');
    const closeModalButton = document.getElementById('close-modal');
    const saveSkillPlanButton = document.getElementById('save-skill-plan-button');
    const cancelSkillPlanButton = document.getElementById('cancel-skill-plan-button'); // Reference to Cancel button
    const addSkillPlanButton = document.getElementById('add-skill-plan-button'); // Reference to Add Skill Plan button

    // Event listener for opening the Add Skill Plan modal
    addSkillPlanButton.addEventListener('click', function () {
        addSkillPlanModal.style.display = 'block';
    });

    // Event listener for closing the modal via the close icon
    closeModalButton.addEventListener('click', function () {
        addSkillPlanModal.style.display = 'none';
    });

    // Event listener for closing the modal via the Cancel button
    cancelSkillPlanButton.addEventListener('click', function () {
        addSkillPlanModal.style.display = 'none';
    });

    // Event listener for saving the skill plan
    saveSkillPlanButton.addEventListener('click', async () => {
        const skillPlanName = document.getElementById("skill-plan-name").value.trim();
        const skillPlanContents = document.getElementById("skill-plan-contents").value.trim();

        if (skillPlanName && skillPlanContents) {
            showToast("Saving Skill Plan...", "success"); // Optional: Indicate saving
            try {
                const response = await fetch('/save-skill-plan', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ name: skillPlanName, contents: skillPlanContents })
                });

                if (response.ok) {
                    showToast("Skill Plan Saved!", "success");
                    addSkillPlanModal.style.display = "none";
                    window.location.reload();
                } else {
                    showToast("Error saving skill plan.", "error");
                }
            } catch (error) {
                console.error("Error:", error);
                showToast("Failed to save skill plan.", "error");
            }
        } else {
            showToast("Please fill in both fields.", "error");
        }
    });

    // Toggle Charts Logic (New)
    const toggleChartButton = document.getElementById('toggle-chart-button');
    const toggleChartIcon = document.getElementById('toggle-chart-icon');
    
    let showingCharacterTable = true; // Flag to track which table is currently shown
    
    toggleChartButton.addEventListener('click', () => {
        if (showingCharacterTable) {
            // Hide Character Table and Show Skill Plan Table
            document.getElementById('character-table').style.display = 'none';
            document.getElementById('skill-plan-table').style.display = 'block';
    
            skillPlanTable.redraw(); // Ensure the table redraws itself
            adjustTableHeight("skill-plan-table"); // Optional, if needed
    
            showingCharacterTable = false;
            // Update Toggle Icon
            toggleChartIcon.classList.remove('fa-chart-bar');
            toggleChartIcon.classList.add('fa-user');
            toggleChartButton.title = "Show Character Chart";
            toggleChartButton.setAttribute('aria-label', "Show Character Chart");
        } else {
            // Hide Skill Plan Table and Show Character Table
            document.getElementById('character-table').style.display = 'block';
            document.getElementById('skill-plan-table').style.display = 'none';
    
            characterTable.redraw(); // Ensure the table redraws itself
            adjustTableHeight("character-table"); // Optional, if needed
    
            showingCharacterTable = true;
            // Update Toggle Icon
            toggleChartIcon.classList.remove('fa-user');
            toggleChartIcon.classList.add('fa-chart-bar');
            toggleChartButton.title = "Show Skill Plan Chart";
            toggleChartButton.setAttribute('aria-label', "Show Skill Plan Chart");
        }
    });
    

    // Close modal when clicking outside the modal content
    window.onclick = function(event) {
        if (event.target === addSkillPlanModal) {
            addSkillPlanModal.style.display = 'none';
        }
    };

    /**
     * Event listener for Add Person Icon
     */
    const addPersonButton = document.getElementById('add-person-button');
    if (addPersonButton) {
        addPersonButton.addEventListener('click', function(event) {
            // Since it's an <a> tag, default behavior will handle navigation
            // Optionally, you can add additional functionality here
        });
    }

    adjustTableHeight("character-table");
    adjustTableHeight("skill-plan-table");
});
