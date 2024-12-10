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
    ContentCopy,
    Delete,
    CheckCircle,
    AccessTime,
    Error as ErrorIcon,
} from '@mui/icons-material';
import { calculateDaysFromToday } from './Utils';

const SkillPlanRow = ({ row }) => {
    const [open, setOpen] = React.useState(false);

    return (
        <React.Fragment>
            <TableRow
                className="bg-gray-800 hover:bg-gray-700 transition-colors"
                sx={{ borderBottom: row.children.length === 0 ? '1px solid rgba(255,255,255,0.1)' : 'none' }}
            >
                <TableCell sx={{ width: '40px' }}>
                    {row.children.length > 0 && (
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
                    {row.planName}
                </TableCell>
                <TableCell className="whitespace-nowrap">
                    <IconButton
                        size="small"
                        onClick={() => window.copySkillPlan(row.planName)}
                        sx={{ color: '#14b8a6', '&:hover': { color: '#ffffff' }, mr: 1 }}
                        title="Copy Skill Plan"
                    >
                        <ContentCopy fontSize="small" />
                    </IconButton>
                    <IconButton
                        size="small"
                        onClick={() => window.deleteSkillPlan(row.planName)}
                        sx={{ color: '#ef4444', '&:hover': { color: '#ffffff' } }}
                        title="Delete Skill Plan"
                    >
                        <Delete fontSize="small" />
                    </IconButton>
                </TableCell>
            </TableRow>
            {row.children.length > 0 && (
                <TableRow className="bg-gray-800">
                    <TableCell style={{ paddingBottom: 0, paddingTop: 0 }} colSpan={3}>
                        <Collapse in={open} timeout="auto" unmountOnExit>
                            <Box margin={1}>
                                <Table size="small" className="bg-gray-700 rounded-md overflow-hidden">
                                    <TableBody>
                                        {row.children.map((child) => (
                                            <TableRow
                                                key={child.id}
                                                className="hover:bg-gray-600 transition-colors"
                                            >
                                                <TableCell className="pl-8 text-gray-300 flex items-center border-b border-gray-600">
                                                    {child.statusIcon}
                                                    <span className="ml-2">↳ {child.characterName}</span>
                                                </TableCell>
                                                <TableCell className="text-gray-300 border-b border-gray-600">
                                                    {child.statusText}
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

            return {
                id: skillPlan.Name,
                planName: skillPlan.Name,
                children,
            };
        });

        return data;
    }, [skillPlans, identities]);

    return (
        <div className="mb-8 w-full">
            <TableContainer className="rounded-md border border-gray-700 overflow-hidden">
                <Table>
                    <TableHead>
                        <TableRow className="bg-gray-900">
                            <TableCell sx={{ width: '40px' }} />
                            <TableCell className="text-teal-200 font-bold">Skill Plan</TableCell>
                            <TableCell className="text-teal-200 font-bold">Actions</TableCell>
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
