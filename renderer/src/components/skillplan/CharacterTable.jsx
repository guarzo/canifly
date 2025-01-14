// src/components/CharacterTable.jsx
import React, { useMemo, useState } from 'react';
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
    KeyboardArrowUp
} from '@mui/icons-material';
import {
    CheckCircle,
    AccessTime,
    Error as ErrorIcon,
    ArrowUpward,
    ArrowDownward
} from '@mui/icons-material';
import { calculateDaysFromToday, formatNumberWithCommas } from '../../utils/formatter.jsx';
import CharacterDetailModal from '../common/CharacterDetailModal.jsx';

const generatePlanStatus = (planName, characterDetails) => {
    const qualified = characterDetails.QualifiedPlans?.[planName];
    const pending = characterDetails.PendingPlans?.[planName];
    const pendingFinishDate = characterDetails.PendingFinishDates?.[planName];
    const missingSkillsForPlan = characterDetails.MissingSkills?.[planName] || {};
    const missingCount = Object.keys(missingSkillsForPlan).length;

    let status = {
        statusIcon: null,
        statusText: ''
    };

    if (qualified) {
        status = {
            statusIcon: <CheckCircle style={{ color: 'green' }} fontSize="small" />,
            statusText: 'Qualified'
        };
    } else if (pending) {
        const daysRemaining = calculateDaysFromToday(pendingFinishDate);
        status = {
            statusIcon: <AccessTime style={{ color: 'orange' }} fontSize="small" />,
            statusText: `Pending ${daysRemaining ? `(${daysRemaining})` : ''}`
        };
    } else if (missingCount > 0) {
        status = {
            statusIcon: <ErrorIcon style={{ color: 'red' }} fontSize="small" />,
            statusText: `${missingCount} skills missing`
        };
    }

    return status;
};

const CharacterRow = ({ row, conversions }) => {
    const [open, setOpen] = useState(false);
    const [detailOpen, setDetailOpen] = useState(false);

    return (
        <React.Fragment>
            <CharacterDetailModal
                open={detailOpen}
                onClose={() => setDetailOpen(false)}
                character={row.fullCharacter}
                skillConversions={conversions}
            />
            <TableRow
                className="hover:bg-gray-700 transition-colors duration-200"
                sx={{
                    borderBottom:
                        row.plans.length === 0 ? '1px solid rgba(255,255,255,0.1)' : 'none'
                }}
            >
                <TableCell sx={{ width: '40px', paddingX: '0.5rem' }}>
                    {row.plans.length > 0 && (
                        <Tooltip title={open ? 'Collapse' : 'Expand'} arrow>
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
                    <img
                        src={`https://images.evetech.net/characters/${row.id}/portrait?size=32`}
                        alt={`${row.CharacterName}'s portrait`}
                        style={{
                            width: '24px',
                            height: '24px',
                            borderRadius: '50%',
                            verticalAlign: 'middle',
                            display: 'inline-block',
                            marginRight: '0.5rem'
                        }}
                    />
                    <span
                        style={{ verticalAlign: 'middle' }}
                        className="font-semibold text-sm text-teal-200 cursor-pointer underline"
                        onClick={() => setDetailOpen(true)}
                    >
                        {row.CharacterName}
                    </span>
                </TableCell>
                <TableCell className="whitespace-nowrap text-teal-100 px-2 py-2">
                    {row.TotalSP}
                </TableCell>
            </TableRow>
            {row.plans.length > 0 && (
                <TableRow>
                    <TableCell style={{ paddingBottom: 0, paddingTop: 0 }} colSpan={3}>
                        <Collapse in={open} timeout="auto" unmountOnExit>
                            <Box margin={1}>
                                <Table
                                    size="small"
                                    className="bg-gray-700 rounded-md overflow-hidden"
                                >
                                    <TableBody>
                                        {row.plans.map((plan) => (
                                            <TableRow
                                                key={plan.id}
                                                className="hover:bg-gray-600 transition-colors duration-200"
                                            >
                                                <TableCell className="pl-8 text-gray-300 flex items-center border-b border-gray-600 py-2">
                                                    {plan.statusIcon}
                                                    <span className="ml-2">
                                                        ↳ {plan.planName}
                                                    </span>
                                                </TableCell>
                                                <TableCell className="text-gray-300 border-b border-gray-600 py-2">
                                                    {plan.statusText}
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

CharacterRow.propTypes = {
    row: PropTypes.object.isRequired,
    conversions: PropTypes.object.isRequired
};

const CharacterTable = ({ characters, skillPlans, conversions }) => {
    // sortBy can be "name" or "sp"
    const [sortBy, setSortBy] = useState('sp');
    const [sortDirection, setSortDirection] = useState('asc');

    const handleSort = (column) => {
        // If we're already sorting by that column, flip direction
        // Otherwise, switch to that column with ascending direction
        if (sortBy === column) {
            setSortDirection((prev) => (prev === 'asc' ? 'desc' : 'asc'));
        } else {
            setSortBy(column);
            setSortDirection('asc');
        }
    };

    // Prepare the raw data
    const characterData = useMemo(() => {
        return characters.map((character) => {
            const characterDetails = character.Character || {};
            const totalSP = characterDetails.CharacterSkillsResponse?.total_sp || 0;
            const TotalSPFormatted = formatNumberWithCommas(totalSP);

            const plans = Object.keys(skillPlans).map((planName) => {
                const status = generatePlanStatus(planName, characterDetails);
                return {
                    id: `${characterDetails.CharacterID}-${planName}`,
                    planName,
                    statusIcon: status.statusIcon,
                    statusText: status.statusText
                };
            });

            return {
                id: characterDetails.CharacterID,
                CharacterName: characterDetails.CharacterName || '',
                rawSP: totalSP,
                TotalSP: TotalSPFormatted,
                plans,
                fullCharacter: character
            };
        });
    }, [characters, skillPlans]);

    // Sort the data based on sortBy & sortDirection
    const sortedData = useMemo(() => {
        const sorted = [...characterData];

        sorted.sort((a, b) => {
            if (sortBy === 'name') {
                // Compare by CharacterName (string)
                const nameA = a.CharacterName.toLowerCase();
                const nameB = b.CharacterName.toLowerCase();
                return nameA.localeCompare(nameB);
            } else {
                // sortBy === 'sp'
                // Compare by rawSP (numeric)
                return a.rawSP - b.rawSP;
            }
        });

        if (sortDirection === 'desc') {
            sorted.reverse();
        }

        return sorted;
    }, [characterData, sortBy, sortDirection]);

    return (
        <div className="mb-8 w-full">
            <TableContainer className="rounded-md border border-gray-700 overflow-hidden">
                <Table>
                    <TableHead>
                        <TableRow className="bg-gradient-to-r from-gray-900 to-gray-800">
                            {/* Expand/Collapse Column */}
                            <TableCell sx={{ width: '40px', paddingX: '0.5rem' }} />

                            {/* Character Name Column */}
                            <TableCell
                                onClick={() => handleSort('name')}
                                className="text-teal-200 font-bold uppercase py-2 px-2 text-sm cursor-pointer select-none"
                            >
                                Character Name
                                {sortBy === 'name' && (
                                    sortDirection === 'asc' ? (
                                        <ArrowUpward
                                            fontSize="small"
                                            style={{ marginLeft: 4, verticalAlign: 'middle' }}
                                        />
                                    ) : (
                                        <ArrowDownward
                                            fontSize="small"
                                            style={{ marginLeft: 4, verticalAlign: 'middle' }}
                                        />
                                    )
                                )}
                            </TableCell>

                            {/* Skill Points Column */}
                            <TableCell
                                onClick={() => handleSort('sp')}
                                className="text-teal-200 font-bold uppercase py-2 px-2 text-sm cursor-pointer select-none"
                            >
                                Total Skill Points
                                {sortBy === 'sp' && (
                                    sortDirection === 'asc' ? (
                                        <ArrowUpward
                                            fontSize="small"
                                            style={{ marginLeft: 4, verticalAlign: 'middle' }}
                                        />
                                    ) : (
                                        <ArrowDownward
                                            fontSize="small"
                                            style={{ marginLeft: 4, verticalAlign: 'middle' }}
                                        />
                                    )
                                )}
                            </TableCell>
                        </TableRow>
                    </TableHead>
                    <TableBody>
                        {sortedData.map((row) => (
                            <CharacterRow key={row.id} row={row} conversions={conversions} />
                        ))}
                    </TableBody>
                </Table>
            </TableContainer>
        </div>
    );
};

CharacterTable.propTypes = {
    characters: PropTypes.array.isRequired,
    skillPlans: PropTypes.object.isRequired,
    conversions: PropTypes.object.isRequired
};

export default CharacterTable;
