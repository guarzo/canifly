// CharacterTable.jsx
import React, { useMemo } from 'react';
import PropTypes from 'prop-types';
import {
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
    Collapse,
    IconButton,
    Typography,
    Box,
} from '@mui/material';
import {
    KeyboardArrowDown,
    KeyboardArrowUp,
    CheckCircle,
    AccessTime,
    Error as ErrorIcon,
} from '@mui/icons-material';
import { calculateDaysFromToday, formatNumberWithCommas } from './Utils';

const generatePlanStatus = (planName, characterDetails) => {
    const qualified = characterDetails.QualifiedPlans?.[planName];
    const pending = characterDetails.PendingPlans?.[planName];
    const pendingFinishDate = characterDetails.PendingFinishDates?.[planName];
    const missingSkillsForPlan = characterDetails.MissingSkills?.[planName] || {};
    const missingCount = Object.keys(missingSkillsForPlan).length;

    let status = {
        statusIcon: null,
        statusText: '',
    };

    if (qualified) {
        status = {
            statusIcon: <CheckCircle style={{ color: 'green' }} fontSize="small" />,
            statusText: 'Qualified',
        };
    } else if (pending) {
        const daysRemaining = calculateDaysFromToday(pendingFinishDate);
        status = {
            statusIcon: <AccessTime style={{ color: 'orange' }} fontSize="small" />,
            statusText: `Pending ${daysRemaining ? `(${daysRemaining})` : ''}`,
        };
    } else if (missingCount > 0) {
        status = {
            statusIcon: <ErrorIcon style={{ color: 'red' }} fontSize="small" />,
            statusText: `${missingCount} missing`,
        };
    }

    console.log(`Character: ${characterDetails.CharacterName}, Plan: ${planName}, Status:`, status);
    return status;
};

const CharacterRow = ({ row }) => {
    const [open, setOpen] = React.useState(false);

    console.log('Rendering CharacterRow for:', row.CharacterName, 'with plans:', row.plans);

    return (
        <React.Fragment>
            <TableRow className="bg-gray-800">
                <TableCell>
                    {row.plans.length > 0 && (
                        <IconButton
                            size="small"
                            onClick={() => setOpen(!open)}
                            className="text-teal-200"
                        >
                            {open ? <KeyboardArrowUp /> : <KeyboardArrowDown />}
                        </IconButton>
                    )}
                </TableCell>
                <TableCell className="text-teal-200 font-semibold">
                    {row.CharacterName}
                </TableCell>
                <TableCell>{row.TotalSP}</TableCell>
            </TableRow>
            {row.plans.length > 0 && (
                <TableRow className="bg-gray-800">
                    <TableCell style={{ paddingBottom: 0, paddingTop: 0 }} colSpan={3}>
                        <Collapse in={open} timeout="auto" unmountOnExit>
                            <Box margin={1}>
                                <Table size="small" className="bg-gray-700">
                                    <TableBody>
                                        {row.plans.map((plan) => (
                                            <TableRow key={plan.id}>
                                                <TableCell className="pl-8 text-gray-300 flex items-center">
                                                    {plan.statusIcon}
                                                    <span className="ml-2">↳ {plan.planName}</span>
                                                </TableCell>
                                                <TableCell className="text-gray-300">
                                                    {plan.statusText}
                                                </TableCell>
                                                <TableCell />
                                            </TableRow>
                                        ))}
                                    </TableBody>
                                </Table>
                            </Box>
                        </Collapse>
                    </TableCell>
                </TableRow>
            )}
        </React.Fragment>
    );
};

CharacterRow.propTypes = {
    row: PropTypes.object.isRequired,
};

const CharacterTable = ({ identities, skillPlans }) => {
    const characterData = useMemo(() => {
        const data = identities.map((identity) => {
            const characterDetails = identity.Character || {};
            const TotalSP = formatNumberWithCommas(
                characterDetails.CharacterSkillsResponse?.total_sp || 0
            );

            const plans = Object.keys(skillPlans).map((planName) => {
                const status = generatePlanStatus(planName, characterDetails);
                return {
                    id: `${characterDetails.CharacterID}-${planName}`,
                    planName,
                    statusIcon: status.statusIcon,
                    statusText: status.statusText,
                };
            });

            console.log(`Processed Character: ${characterDetails.CharacterName}`, {
                plans,
            });

            return {
                id: characterDetails.CharacterID,
                CharacterName: characterDetails.CharacterName,
                TotalSP,
                plans,
            };
        });

        console.log('Final CharacterData:', data);
        return data;
    }, [identities, skillPlans]);

    return (
        <div className="mb-8 w-full">
            <Typography variant="h5" gutterBottom className="text-teal-200">
                Character Table
            </Typography>
            <TableContainer>
                <Table>
                    <TableHead>
                        <TableRow className="bg-gray-900">
                            <TableCell />
                            <TableCell className="text-teal-200">Character Name</TableCell>
                            <TableCell className="text-teal-200">Total Skill Points</TableCell>
                        </TableRow>
                    </TableHead>
                    <TableBody>
                        {characterData.map((row) => (
                            <CharacterRow key={row.id} row={row} />
                        ))}
                    </TableBody>
                </Table>
            </TableContainer>
        </div>
    );
};

CharacterTable.propTypes = {
    identities: PropTypes.array.isRequired,
    skillPlans: PropTypes.object.isRequired,
};

export default CharacterTable;
