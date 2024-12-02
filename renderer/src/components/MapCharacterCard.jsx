// CharacterCard.jsx

import React from 'react';
import { Card, Typography, CardContent, useTheme } from '@mui/material';

const CharacterCard = ({ char, handleDragStart, mtimeToColor }) => {
    const theme = useTheme();

    const borderColor = mtimeToColor[char.mtime] || theme.palette.secondary.main;

    return (
        <Card
            draggable
            onDragStart={(e) => handleDragStart(e, char.charId)}
            aria-label={`Draggable character ${char.name}`}
            title={`Drag to associate with an account`}
            sx={{
                borderLeft: `4px solid ${borderColor}`,
                backgroundColor: theme.palette.background.paper,
                borderRadius: 2,
                cursor: 'grab',
                '&:hover': {
                    backgroundColor: theme.palette.action.hover,
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

export default CharacterCard;
