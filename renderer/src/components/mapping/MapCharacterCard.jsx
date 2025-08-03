// CharacterCard.jsx
import React from 'react';
import PropTypes from 'prop-types';
import { Card, Typography, CardContent, useTheme } from '@mui/material';
import { formatDate } from '../../utils/formatter.jsx';

const CharacterCard = ({ char, handleDragStart, mtimeToColor }) => {
    const theme = useTheme();

    const borderColor = mtimeToColor[char.mtime] || theme.palette.secondary.main;

    // Format the date without the year and use 24-hour time
    const formattedDate = formatDate(char.mtime);

    return (
        <Card
            draggable
            onDragStart={(e) => handleDragStart(e, char.charId)}
            sx={{
                borderLeft: `4px solid ${borderColor}`,
                backgroundColor: theme.palette.background.paper,
                borderRadius: 2,
                cursor: 'grab',
                boxShadow: 3, // Added shadow
                transition: 'transform 0.2s ease-in-out',
                '&:hover': {
                    backgroundColor: theme.palette.action.hover,
                    transform: 'scale(1.02)',
                },
            }}
        >
            <CardContent>
                <Typography variant="h6" color="text.primary">
                    {char.name}
                </Typography>
                {/* Removed ID display */}
                <Typography variant="caption" color="text.secondary">
                    {formattedDate}
                </Typography>
            </CardContent>
        </Card>
    );
};

CharacterCard.propTypes = {
    char: PropTypes.shape({
        name: PropTypes.string.isRequired,
        charId: PropTypes.string.isRequired,
        mtime: PropTypes.string.isRequired,
    }).isRequired,
    handleDragStart: PropTypes.func.isRequired,
    mtimeToColor: PropTypes.objectOf(PropTypes.string).isRequired,
};

export default CharacterCard;
