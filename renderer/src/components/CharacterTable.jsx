// CharacterTable.jsx
import React, { useMemo, useRef } from 'react';
import PropTypes from 'prop-types';
import TabulatorTable from './TabulatorTable';
import { calculateDaysFromToday, formatNumberWithCommas } from './Utils';

// CharacterTable.jsx
const generatePlanStatus = (planName, characterDetails) => {
    const qualified = characterDetails.QualifiedPlans?.[planName];
    const pending = characterDetails.PendingPlans?.[planName];
    const pendingFinishDate = characterDetails.PendingFinishDates?.[planName];
    const missingSkillsForPlan = characterDetails.MissingSkills?.[planName] || {};
    const missingCount = Object.keys(missingSkillsForPlan).length;

    if (qualified) {
        return `
            <i class="fas fa-check-circle text-green-500"></i> Qualified
        `;
    } else if (pending) {
        const daysRemaining = calculateDaysFromToday(pendingFinishDate);
        return `
            <i class="fas fa-clock text-orange-500"></i> Pending ${daysRemaining ? `(${daysRemaining})` : ""}
        `;
    } else if (missingCount > 0) {
        return `
            <i class="fas fa-exclamation-circle text-red-500"></i> ${missingCount} missing
        `;
    }
    return '';
};


const CharacterTable = ({ identities, skillPlans, tabulatorRef }) => {
    const characterData = useMemo(() => {
        return identities.map((identity) => {
            const characterDetails = identity.Character || {};
            let TotalSP = characterDetails.CharacterSkillsResponse?.total_sp || 0;
            TotalSP = formatNumberWithCommas(TotalSP); // Format with commas

            const children = Object.keys(skillPlans)
                .map((planName) => ({
                    CharacterName: planName,
                    TotalSP: generatePlanStatus(planName, characterDetails),
                }))
                .filter(Boolean);

            return {
                ...identity,
                TotalSP,
                _children: children,
            };
        });
    }, [identities, skillPlans]);

    return (
        <div className="mb-8 w-full">
            <h2 className="text-2xl font-bold mb-4 text-teal-200">Character Table</h2>
            <TabulatorTable
                id="character-table"
                data={characterData}
                columns={[
                    {
                        title: "Character Name",
                        field: "CharacterName",
                        formatter: "html", // Use HTML for formatting
                    },
                    {
                        title: "Total Skill Points",
                        field: "TotalSP",
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

CharacterTable.propTypes = {
    identities: PropTypes.array.isRequired,
    skillPlans: PropTypes.object.isRequired,
    tabulatorRef: PropTypes.object,
};

export default CharacterTable;
