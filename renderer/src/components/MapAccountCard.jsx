// AccountCard.jsx

import React from 'react';
import { Card, Typography, List, ListItem, ListItemText, IconButton, useTheme } from '@mui/material';
import DeleteIcon from '@mui/icons-material/Delete';

const AccountCard = ({ mapping, associations, handleUnassociate, handleDrop, mtimeToColor }) => {
    const theme = useTheme();
    const userId = mapping.subDir;
    const userName = mapping.availableUserFiles.length > 0 ? mapping.availableUserFiles[0].name : `Account ${userId.replace('settings_', '')}`;
    const associatedChars = associations.filter(assoc => assoc.userId === userId);

    // Assuming the first user file contains the latest mtime for the account
    const accountMtime = mapping.availableUserFiles.length > 0 ? mapping.availableUserFiles[0].mtime : new Date().toISOString();
    const borderColor = mtimeToColor[accountMtime] || theme.palette.primary.main;

    return (
        <Card
            onDragOver={(e) => e.preventDefault()}
            onDrop={(e) => handleDrop(e, userId)}
            aria-label={`Drop character to associate with account ${userId}`}
            sx={{
                marginBottom: 2,
                borderLeft: `4px solid ${borderColor}`,
                backgroundColor: theme.palette.background.paper,
                borderRadius: 2,
                paddingLeft: 2,
                cursor: 'pointer',
                '&:hover': {
                    backgroundColor: theme.palette.action.hover,
                },
            }}
        >
            <Typography variant="h6" color="text.primary" gutterBottom>
                {userName}
            </Typography>
            <Typography variant="body2" color="text.secondary" gutterBottom>
                {new Date(accountMtime).toLocaleString()}
            </Typography>
            <List>
                {associatedChars.map(assoc => (
                    <ListItem
                        key={`assoc-${assoc.charId}`}
                        secondaryAction={
                            <IconButton
                                edge="end"
                                aria-label={`Unassociate ${assoc.charName}`}
                                onClick={() => handleUnassociate(userId, assoc.charId, assoc.charName)}
                            >
                                <DeleteIcon color="secondary" />
                            </IconButton>
                        }
                        sx={{
                            '&:hover': {
                                backgroundColor: theme.palette.action.selected,
                            },
                        }}
                    >
                        <ListItemText
                            primary={`${assoc.charName} (${assoc.charId})`}
                            secondary={new Date(assoc.mtime).toLocaleString()}
                            primaryTypographyProps={{ color: 'text.primary' }}
                            secondaryTypographyProps={{ color: 'text.secondary' }}
                        />
                    </ListItem>
                ))}
            </List>
        </Card>
    );
};

export default AccountCard;
