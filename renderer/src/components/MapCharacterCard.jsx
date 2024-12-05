import React from 'react';
import PropTypes from 'prop-types';
import { Card, Typography, CardContent, useTheme } from '@mui/material';

const CharacterCard = ({ char, handleDragStart, mtimeToColor }) => {
    const theme = useTheme();

    const borderColor = mtimeToColor[char.mtime] || theme.palette.secondary.main;

    return (
        <Card
            draggable
            onDragStart={(e) => handleDragStart(e, char.charId)}
            sx={{
                borderLeft: `4px solid ${borderColor}`,
                backgroundColor: theme.palette.background.paper,
                borderRadius: 2,
                cursor: 'grab',
                '&:hover': {
                    backgroundColor: theme.palette.action.hover,
                    transform: 'scale(1.02)',
                    transition: 'transform 0.2s ease-in-out',
                },
            }}
        >
            <CardContent>
                <Typography variant="h6" color="text.primary">
                    {char.name}
                </Typography>
                <Typography variant="body2" color="text.secondary">
                    ID: {char.charId}
                </Typography>
                <Typography variant="caption" color="text.secondary">
                    {new Date(char.mtime).toLocaleString()}
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
