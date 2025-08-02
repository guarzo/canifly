// src/components/mapping/Mapping.jsx
import React, { useState, useEffect } from 'react';
import { motion } from 'framer-motion';
import { toast } from 'react-toastify';
import { Grid, Box, Typography } from '@mui/material';
import AccountCard from '../components/mapping/MapAccountCard.jsx';
import CharacterCard from '../components/mapping/MapCharacterCard.jsx';
import { useConfirmDialog } from '../hooks/useConfirmDialog.jsx';
import { associateCharacter, unassociateCharacter } from '../api/apiService.jsx';
import { mappingInstructions } from './../utils/instructions';
import PageHeader from '../components/common/SubPageHeader.jsx';
import { useAppData } from '../hooks/useAppData';

function roundToMinute(mtime) {
    const date = new Date(mtime);
    date.setSeconds(0);
    return date.toISOString();
}

const Mapping = ({ associations: initialAssociations, subDirs }) => {
    const { refreshData } = useAppData();
    const [accounts, setAccounts] = useState([]);
    const [availableCharacters, setAvailableCharacters] = useState([]);
    const [associations, setAssociations] = useState(initialAssociations);
    const [mtimeToColor, setMtimeToColor] = useState({});
    const [showConfirmDialog, confirmDialog] = useConfirmDialog();

    useEffect(() => {
        if (subDirs.length === 0) return;

        const userMap = {};
        subDirs.forEach((mapping) => {
            mapping.availableUserFiles.forEach((userFile) => {
                const roundedMtime = roundToMinute(userFile.mtime);
                if (
                    !userMap[userFile.userId] ||
                    new Date(roundedMtime) > new Date(userMap[userFile.userId].mtime)
                ) {
                    userMap[userFile.userId] = { ...userFile, mtime: roundedMtime };
                }
            });
        });

        const uniqueAccounts = Object.values(userMap).sort(
            (a, b) => new Date(b.mtime) - new Date(a.mtime)
        );

        const charMap = {};
        subDirs.forEach((mapping) => {
            mapping.availableCharFiles.forEach((charFile) => {
                const roundedMtime = roundToMinute(charFile.mtime);
                const { charId } = charFile;
                if (!charMap[charId] || new Date(roundedMtime) > new Date(charMap[charId].mtime)) {
                    charMap[charId] = { ...charFile, mtime: roundedMtime, profile: mapping.profile };
                }
            });
        });

        const associatedCharIds = new Set(associations.map((a) => a.charId));
        const uniqueChars = Object.values(charMap)
            .filter((ch) => !associatedCharIds.has(ch.charId))
            .sort((a, b) => new Date(b.mtime) - new Date(a.mtime));

        setAccounts(uniqueAccounts);
        setAvailableCharacters(uniqueChars);
        assignColors(uniqueAccounts, uniqueChars);
    }, [subDirs, associations]);

    const assignColors = (uniqueAccounts, uniqueChars) => {
        const predefinedColors = ['#4caf50', '#f44336', '#ff9800', '#9c27b0', '#00bcd4', '#e91e63'];

        const accountMtimes = uniqueAccounts.map((a) => a.mtime);
        const charMtimes = uniqueChars.map((c) => c.mtime);
        const allMtimes = [...accountMtimes, ...charMtimes];

        const uniqueMtimes = Array.from(new Set(allMtimes)).sort(
            (a, b) => new Date(a) - new Date(b)
        );

        const colorMapping = uniqueMtimes.reduce((acc, mtime, index) => {
            acc[mtime] = predefinedColors[index % predefinedColors.length];
            return acc;
        }, {});

        setMtimeToColor(colorMapping);
    };

    console.log(mtimeToColor);

    const handleDragStart = (event, charId) => {
        event.dataTransfer.setData('text/plain', charId);
    };

    const handleDrop = async (event, userId, userName) => {
        event.preventDefault();
        const charId = event.dataTransfer.getData('text/plain');

        if (associations.some((assoc) => assoc.charId === charId)) {
            toast.error('That character is already associated!');
            return;
        }

        const char = availableCharacters.find((c) => c.charId === charId);
        const charName = char?.name;
        if (!char) {
            toast.error('Character not found.');
            return;
        }

        const confirmAssoc = await showConfirmDialog({
            title: 'Confirm Association',
            message: `Associate "${charName}" with account "${userName}"?`,
        });

        if (!confirmAssoc.isConfirmed) return;

        const result = await associateCharacter(userId, charId, userName, charName);
        if (result && result.success) {
            toast.success(result.message);

            setAvailableCharacters((prev) => prev.filter((ch) => ch.charId !== charId));

            const newAssoc = { userId, charId, charName, mtime: char.mtime };
            setAssociations((prev) => [...prev, newAssoc]);

            if (refreshData) {
                await refreshData();
            }
        }
    };

    const handleUnassociate = async (userId, charId, charName, userName) => {
        const confirmUnassoc = await showConfirmDialog({
            title: 'Confirm Unassociation',
            message: `Unassociate "${charName}" from account "${userName}"?`,
        });

        if (!confirmUnassoc.isConfirmed) return;

        const result = await unassociateCharacter(userId, charId, userName, charName);
        if (result && result.success) {
            toast.success(result.message);
            setAssociations((prev) => prev.filter((a) => a.charId !== charId || a.userId !== userId));
            if (refreshData) {
                await refreshData();
            }
        }
    };

    const noCharacters = availableCharacters.length === 0;
    
    const containerVariants = {
        hidden: { opacity: 0 },
        visible: {
            opacity: 1,
            transition: {
                staggerChildren: 0.1
            }
        }
    };

    const itemVariants = {
        hidden: { opacity: 0, y: 20 },
        visible: {
            opacity: 1,
            y: 0,
            transition: { duration: 0.5 }
        }
    };

    return (
        <div className="min-h-screen px-4 pb-10 pt-16">
            <PageHeader
                title="Map Character Files to User Files"
                instructions={mappingInstructions}
                storageKey="showMappingInstructions"
            />
            <Box className="max-w-7xl mx-auto">
                <motion.div
                    variants={containerVariants}
                    initial="hidden"
                    animate="visible"
                >
                    <Grid container spacing={4}>
                        <Grid item xs={12} md={noCharacters ? 12 : 6}>
                            {accounts.length === 0 ? (
                                <Box textAlign="center" className="text-gray-300">
                                    <Typography variant="body1" className="text-gray-400">
                                        No accounts found.
                                    </Typography>
                                </Box>
                            ) : (
                                <Grid container spacing={2}>
                                    {accounts.map((mapping, index) => (
                                        <Grid
                                            item
                                            xs={12}
                                            sm={noCharacters ? 6 : 12}
                                            key={`${mapping.userId}-${mapping.mtime}`}
                                        >
                                            <motion.div
                                                variants={itemVariants}
                                                initial="hidden"
                                                animate="visible"
                                                transition={{ delay: index * 0.05 }}
                                            >
                                                <AccountCard
                                            mapping={mapping}
                                            associations={associations}
                                            handleUnassociate={handleUnassociate}
                                            handleDrop={handleDrop}
                                            mtimeToColor={mtimeToColor}
                                        />
                                            </motion.div>
                                        </Grid>
                                    ))}
                                </Grid>
                            )}
                        </Grid>

                    {/* Only display the Characters column if there are actually unassociated characters */}
                    {!noCharacters && (
                        <Grid item xs={12} md={6}>
                            <Grid container spacing={2} data-testid="available-characters">
                                {availableCharacters.map((char, index) => (
                                    <Grid item xs={12} sm={6} key={`${char.charId}-${char.mtime}`}>
                                        <motion.div
                                            variants={itemVariants}
                                            initial="hidden"
                                            animate="visible"
                                            transition={{ delay: index * 0.05 }}
                                        >
                                            <CharacterCard
                                                char={char}
                                                handleDragStart={handleDragStart}
                                                mtimeToColor={mtimeToColor}
                                            />
                                        </motion.div>
                                    </Grid>
                                ))}
                            </Grid>
                        </Grid>
                    )}
                    </Grid>
                </motion.div>
            </Box>

            {confirmDialog}
        </div>
    );
};

export default Mapping;
