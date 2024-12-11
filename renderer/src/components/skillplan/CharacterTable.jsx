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
    Box,
} from '@mui/material';
import {
    KeyboardArrowDown,
    KeyboardArrowUp,
    CheckCircle,
    AccessTime,
    Error as ErrorIcon,
} from '@mui/icons-material';
import { calculateDaysFromToday, formatNumberWithCommas } from '../../utils/utils.jsx';

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
            statusText: `${missingCount} skills missing`,
        };
    }

    return status;
};

const CharacterRow = ({ row }) => {
    const [open, setOpen] = React.useState(false);

    return (
        <React.Fragment>
            <TableRow
                className="bg-gray-800 hover:bg-gray-700 transition-colors"
                sx={{ borderBottom: row.plans.length === 0 ? '1px solid rgba(255,255,255,0.1)' : 'none' }}
            >
                <TableCell sx={{ width: '40px' }}>
                    {row.plans.length > 0 && (
                        <IconButton
                            size="small"
                            onClick={() => setOpen(!open)}
                            sx={{ color: '#99f6e4', '&:hover': { color: '#ffffff' } }}
                        >
                            {open ? <KeyboardArrowUp /> : <KeyboardArrowDown />}
                        </IconButton>
                    )}
                </TableCell>
                <TableCell className="text-teal-200 font-semibold whitespace-nowrap">
                    {row.CharacterName}
                </TableCell>
                <TableCell className="whitespace-nowrap text-teal-100">
                    {row.TotalSP}
                </TableCell>
            </TableRow>
            {row.plans.length > 0 && (
                <TableRow className="bg-gray-800">
                    <TableCell style={{ paddingBottom: 0, paddingTop: 0 }} colSpan={3}>
                        <Collapse in={open} timeout="auto" unmountOnExit>
                            <Box margin={1}>
                                <Table size="small" className="bg-gray-700 rounded-md overflow-hidden">
                                    <TableBody>
                                        {row.plans.map((plan) => (
                                            <TableRow
                                                key={plan.id}
                                                className="hover:bg-gray-600 transition-colors"
                                            >
                                                <TableCell className="pl-8 text-gray-300 flex items-center border-b border-gray-600">
                                                    {plan.statusIcon}
                                                    <span className="ml-2">↳ {plan.planName}</span>
                                                </TableCell>
                                                <TableCell className="text-gray-300 border-b border-gray-600">
                                                    {plan.statusText}
                                                </TableCell>
                                                <TableCell className="border-b border-gray-600" />
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
        return identities.map((identity) => {
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

            return {
                id: characterDetails.CharacterID,
                CharacterName: characterDetails.CharacterName,
                TotalSP,
                plans,
            };
        });
    }, [identities, skillPlans]);

    return (
        <div className="mb-8 w-full">
            <TableContainer className="rounded-md border border-gray-700 overflow-hidden">
                <Table>
                    <TableHead>
                        <TableRow className="bg-gray-900">
                            <TableCell sx={{ width: '40px' }} />
                            <TableCell className="text-teal-200 font-bold">Character Name</TableCell>
                            <TableCell className="text-teal-200 font-bold">Total Skill Points</TableCell>
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
