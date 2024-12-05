import React from 'react';
import PropTypes from 'prop-types';
import { Card, Typography, List, ListItem, ListItemText, IconButton, useTheme } from '@mui/material';
import DeleteIcon from '@mui/icons-material/Delete';

const AccountCard = ({ mapping, associations, handleUnassociate, handleDrop, mtimeToColor }) => {
    const theme = useTheme();
    const userId = mapping.subDir;
    const userName = mapping.availableUserFiles?.[0]?.name || `Account ${userId.replace('settings_', '')}`;
    const associatedChars = associations.filter(assoc => assoc.userId === userId);

    const accountMtime = mapping.availableUserFiles?.[0]?.mtime || new Date().toISOString();
    const borderColor = mtimeToColor[accountMtime] || theme.palette.primary.main;

    return (
        <Card
            onDragOver={(e) => e.preventDefault()}
            onDrop={(e) => handleDrop(e, userId)}
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
                    >
                        <ListItemText
                            primary={`${assoc.charName} (${assoc.charId})`}
                            secondary={
                                assoc.mtime
                                    ? new Date(assoc.mtime).toLocaleString()
                                    : 'Date not available'
                            }
                            primaryTypographyProps={{ color: 'text.primary' }}
                            secondaryTypographyProps={{ color: 'text.secondary' }}
                        />
                    </ListItem>
                ))}
            </List>
        </Card>
    );
};

AccountCard.propTypes = {
    mapping: PropTypes.shape({
        subDir: PropTypes.string.isRequired,
        availableUserFiles: PropTypes.arrayOf(
            PropTypes.shape({
                name: PropTypes.string,
                mtime: PropTypes.string,
            })
        ),
    }).isRequired,
    associations: PropTypes.arrayOf(
        PropTypes.shape({
            userId: PropTypes.string.isRequired,
            charId: PropTypes.string.isRequired,
            charName: PropTypes.string.isRequired,
            mtime: PropTypes.string,
        })
    ).isRequired,
    handleUnassociate: PropTypes.func.isRequired,
    handleDrop: PropTypes.func.isRequired,
    mtimeToColor: PropTypes.objectOf(PropTypes.string).isRequired,
};

export default AccountCard;
