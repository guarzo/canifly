// SkillPlanTable.jsx
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
    ContentCopy,
    Delete,
    CheckCircle,
    AccessTime,
    Error as ErrorIcon,
} from '@mui/icons-material';
import { calculateDaysFromToday } from './Utils';

const SkillPlanRow = ({ row }) => {
    const [open, setOpen] = React.useState(false);

    console.log('Rendering SkillPlanRow for:', row.planName, 'with children:', row.children);

    return (
        <React.Fragment>
            <TableRow className="bg-gray-800">
                <TableCell>
                    {row.children.length > 0 && (
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
                    {row.planName}
                </TableCell>
                <TableCell>
                    <IconButton
                        size="small"
                        onClick={() => window.copySkillPlan(row.planName)}
                        className="hover:text-teal-200 mr-2"
                        title="Copy Skill Plan"
                    >
                        <ContentCopy style={{ color: '#14b8a6' }} />
                    </IconButton>
                    <IconButton
                        size="small"
                        onClick={() => window.deleteSkillPlan(row.planName)}
                        className="hover:text-red-300"
                        title="Delete Skill Plan"
                    >
                        <Delete style={{ color: '#ef4444' }} />
                    </IconButton>
                </TableCell>
            </TableRow>
            {row.children.length > 0 && (
                <TableRow className="bg-gray-800">
                    <TableCell style={{ paddingBottom: 0, paddingTop: 0 }} colSpan={3}>
                        <Collapse in={open} timeout="auto" unmountOnExit>
                            <Box margin={1}>
                                <Table size="small" className="bg-gray-700">
                                    <TableBody>
                                        {row.children.map((child) => (
                                            <TableRow key={child.id}>
                                                <TableCell className="pl-8 text-gray-300 flex items-center">
                                                    {child.statusIcon}
                                                    <span className="ml-2">↳ {child.characterName}</span>
                                                </TableCell>
                                                <TableCell className="text-gray-300">
                                                    {child.statusText}
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

SkillPlanRow.propTypes = {
    row: PropTypes.object.isRequired,
};

const SkillPlanTable = ({ skillPlans, identities }) => {
    const skillPlanData = useMemo(() => {
        const data = Object.values(skillPlans).map((skillPlan) => {
            const qualifiedCharacters = skillPlan.QualifiedCharacters || [];
            const pendingCharacters = skillPlan.PendingCharacters || [];
            const missingCharacters = skillPlan.MissingCharacters || [];

            const children = [
                ...qualifiedCharacters.map((characterName) => ({
                    id: `${skillPlan.Name}-${characterName}`,
                    characterName,
                    statusIcon: <CheckCircle style={{ color: 'green' }} fontSize="small" />,
                    statusText: 'Qualified',
                })),
                ...pendingCharacters.map((characterName) => {
                    const identity = identities.find(
                        (identity) => identity.Character?.CharacterName === characterName
                    );
                    const characterData = identity?.Character || null;
                    const pendingFinishDate =
                        characterData?.PendingFinishDates?.[skillPlan.Name] || '';
                    const daysRemaining = calculateDaysFromToday(pendingFinishDate);
                    return {
                        id: `${skillPlan.Name}-${characterName}`,
                        characterName,
                        statusIcon: <AccessTime style={{ color: 'orange' }} fontSize="small" />,
                        statusText: `Pending ${daysRemaining ? `(${daysRemaining})` : ''}`,
                    };
                }),
                ...missingCharacters.map((characterName) => ({
                    id: `${skillPlan.Name}-${characterName}`,
                    characterName,
                    statusIcon: <ErrorIcon style={{ color: 'red' }} fontSize="small" />,
                    statusText: 'Missing',
                })),
            ];

            console.log(`Processed Skill Plan: ${skillPlan.Name}`, {
                qualifiedCharacters,
                pendingCharacters,
                missingCharacters,
                children,
            });

            return {
                id: skillPlan.Name,
                planName: skillPlan.Name,
                children,
            };
        });

        console.log('Final SkillPlanData:', data);
        return data;
    }, [skillPlans, identities]);

    return (
        <div className="mb-8 w-full">
            <Typography variant="h5" gutterBottom className="text-teal-200">
                Skill Plan Table
            </Typography>
            <TableContainer>
                <Table>
                    <TableHead>
                        <TableRow className="bg-gray-900">
                            <TableCell />
                            <TableCell className="text-teal-200">Skill Plan</TableCell>
                            <TableCell className="text-teal-200">Actions</TableCell>
                        </TableRow>
                    </TableHead>
                    <TableBody>
                        {skillPlanData.map((row) => (
                            <SkillPlanRow key={row.id} row={row} />
                        ))}
                    </TableBody>
                </Table>
            </TableContainer>
        </div>
    );
};

SkillPlanTable.propTypes = {
    skillPlans: PropTypes.object.isRequired,
    identities: PropTypes.array.isRequired,
};

export default SkillPlanTable;
