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
    Box,
    Tooltip
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

import { calculateDaysFromToday } from "../../utils/formatter.jsx";

const SkillPlanRow = ({ row, conversions, contents, onCopySkillPlan, onDeleteSkillPlan }) => {
    const [open, setOpen] = React.useState(false);

    // Lookup typeID from conversions map
    const typeID = conversions[row.planName];
    const planIconUrl = typeID ? `https://images.evetech.net/types/${typeID}/icon` : null;

    return (
        <React.Fragment>
            <TableRow className="hover:bg-gray-700 transition-colors duration-200">
                <TableCell sx={{ width: '40px', paddingX: '0.5rem' }}>
                    {row.children.length > 0 && (
                        <Tooltip title={open ? "Collapse" : "Expand"} arrow>
                            <IconButton
                                size="small"
                                onClick={() => setOpen(!open)}
                                sx={{ color: '#99f6e4', '&:hover': { color: '#ffffff' } }}
                            >
                                {open ? <KeyboardArrowUp /> : <KeyboardArrowDown />}
                            </IconButton>
                        </Tooltip>
                    )}
                </TableCell>
                <TableCell className="text-teal-200 font-semibold whitespace-nowrap px-2 py-2">
                    {planIconUrl && (
                        <img
                            src={planIconUrl}
                            alt={`${row.planName} icon`}
                            style={{
                                width: '24px',
                                height: '24px',
                                borderRadius: '50%',
                                verticalAlign: 'middle',
                                display: 'inline-block',
                                marginRight: '0.5rem',
                            }}
                        />
                    )}
                    {row.planName}
                </TableCell>
                <TableCell className="whitespace-nowrap px-2 py-2">
                    <Tooltip title="Copy Skill Plan" arrow>
                        <IconButton
                            size="small"
                            onClick={() => onCopySkillPlan(row.planName, row.contents)}
                            sx={{ color: '#14b8a6', '&:hover': { color: '#ffffff' }, mr: 1 }}
                        >
                            <ContentCopy fontSize="small" />
                        </IconButton>
                    </Tooltip>
                    <Tooltip title="Delete Skill Plan" arrow>
                        <IconButton
                            size="small"
                            onClick={() => onDeleteSkillPlan(row.planName)}
                            sx={{ color: '#ef4444', '&:hover': { color: '#ffffff' } }}
                        >
                            <Delete fontSize="small" />
                        </IconButton>
                    </Tooltip>
                </TableCell>
            </TableRow>
            {row.children.length > 0 && (
                <TableRow>
                    <TableCell style={{ paddingBottom: 0, paddingTop: 0 }} colSpan={3}>
                        <Collapse in={open} timeout="auto" unmountOnExit>
                            <Box margin={1}>
                                <Table size="small" className="bg-gray-700 rounded-md overflow-hidden">
                                    <TableBody>
                                        {row.children.map((child) => (
                                            <TableRow
                                                key={child.id}
                                                className="hover:bg-gray-600 transition-colors duration-200"
                                            >
                                                <TableCell className="pl-8 text-gray-300 flex items-center border-b border-gray-600 py-2">
                                                    {child.statusIcon}
                                                    <span className="ml-2">↳ {child.characterName}</span>
                                                </TableCell>
                                                <TableCell className="text-gray-300 border-b border-gray-600 py-2">
                                                    {child.statusText}
                                                </TableCell>
                                                <TableCell className="border-b border-gray-600 py-2" />
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
    conversions: PropTypes.object.isRequired,
    onCopySkillPlan: PropTypes.func.isRequired,
    onDeleteSkillPlan: PropTypes.func.isRequired,
};

const SkillPlanTable = ({ skillPlans, characters, conversions, onCopySkillPlan, onDeleteSkillPlan }) => {
    const skillPlanData = useMemo(() => {
        return Object.values(skillPlans).map((skillPlan) => {
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
                    const character = characters.find(
                        (c) => c.Character?.CharacterName === characterName
                    );
                    const characterData = character?.Character || null;
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
                contents: skillPlan.Skills,
                children,
            };
        });
    }, [skillPlans, characters]);

    return (
        <div className="mb-8 w-full">
            <TableContainer className="rounded-md border border-gray-700 overflow-hidden">
                <Table>
                    <TableHead>
                        <TableRow className="bg-gradient-to-r from-gray-900 to-gray-800">
                            <TableCell sx={{ width: '40px', paddingX: '0.5rem' }} />
                            <TableCell className="text-teal-200 font-bold uppercase py-2 px-2 text-sm">
                                Skill Plan
                            </TableCell>
                            <TableCell className="text-teal-200 font-bold uppercase py-2 px-2 text-sm">
                                Actions
                            </TableCell>
                        </TableRow>
                    </TableHead>
                    <TableBody>
                        {skillPlanData.map((row) => (
                            <SkillPlanRow key={row.id} row={row} conversions={conversions} contents={row.contents} onCopySkillPlan={onCopySkillPlan} onDeleteSkillPlan={onDeleteSkillPlan} />
                        ))}
                    </TableBody>
                </Table>
            </TableContainer>
        </div>
    );
};

SkillPlanTable.propTypes = {
    skillPlans: PropTypes.object.isRequired,
    characters: PropTypes.array.isRequired,
    conversions: PropTypes.object.isRequired,
    onCopySkillPlan: PropTypes.func.isRequired,
    onDeleteSkillPlan: PropTypes.func.isRequired,
};

export default SkillPlanTable;
