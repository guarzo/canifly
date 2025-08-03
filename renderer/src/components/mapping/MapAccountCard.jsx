// MapAccountCard.jsx
import React from 'react';
import PropTypes from 'prop-types';
import { Card, Typography, List, ListItem, ListItemText, IconButton, useTheme, Box } from '@mui/material';
import DeleteIcon from '@mui/icons-material/Delete';
import { formatDate } from '../../utils/formatter.jsx';

const AccountCard = ({ mapping, associations, handleUnassociate, handleDrop, mtimeToColor }) => {
    const theme = useTheme();
    const userId = mapping.userId;
    const userName = mapping.name || `Account ${userId}`;
    const associatedChars = associations.filter(assoc => assoc.userId === userId);

    const accountMtime = mapping.mtime || new Date().toISOString();
    const borderColor = mtimeToColor[accountMtime] || theme.palette.primary.main;

    return (
        <Card
            onDragOver={(e) => e.preventDefault()}
            onDrop={(e) => handleDrop(e, userId, userName)}
            sx={{
                marginBottom: 2,
                borderLeft: `4px solid ${borderColor}`,
                backgroundColor: theme.palette.background.paper,
                borderRadius: 2,
                paddingLeft: 2,
                cursor: 'pointer',
                boxShadow: 3, // Added shadow to the account card
                transition: 'background-color 0.2s ease-in-out, box-shadow 0.2s ease-in-out',
                '&:hover': {
                    backgroundColor: theme.palette.action.hover,
                    boxShadow: 4, // Increased shadow on hover
                },
            }}
        >
            <Typography variant="h6" color="text.primary" gutterBottom>
                {userName}
            </Typography>
            <Typography variant="body2" color="text.secondary" gutterBottom>
                {formatDate(accountMtime)}
            </Typography>
            <List>
                {associatedChars.map(assoc => (
                    <ListItem
                        key={`assoc-${assoc.charId}`}
                        secondaryAction={
                            <IconButton
                                edge="end"
                                aria-label={`Unassociate ${assoc.charName}`}
                                onClick={() => handleUnassociate(userId, assoc.charId, assoc.charName, userName)}
                                sx={{
                                    color: theme.palette.error.main, // Use error color for delete actions
                                }}
                            >
                                <DeleteIcon />
                            </IconButton>
                        }
                        sx={{
                            borderRadius: 1,
                            marginBottom: 1,
                            backgroundColor: theme.palette.background.default,
                            boxShadow: 1, // Subtle shadow around each associated character
                        }}
                    >
                        <ListItemText
                            primary={`${assoc.charName}`}
                            primaryTypographyProps={{ color: 'text.primary' }}
                        />
                    </ListItem>
                ))}
            </List>
        </Card>
    );
};

AccountCard.propTypes = {
    mapping: PropTypes.shape({
        userId: PropTypes.string.isRequired,
        name: PropTypes.string,
        mtime: PropTypes.string,
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
