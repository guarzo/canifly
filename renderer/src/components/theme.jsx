// theme.js
import { createTheme } from '@mui/material/styles';

const theme = createTheme({
    palette: {
        primary: {
            main: '#14b8a6', // Tailwind 'teal-500'
        },
        secondary: {
            main: '#ef4444', // Tailwind 'red-500'
        },
        info: {
            main: '#3b82f6', // Tailwind 'blue-500'
        },
        warning: {
            main: '#f59e0b', // Tailwind 'yellow-500'
        },
        background: {
            default: '#1f2937', // Tailwind 'gray-800'
            paper: '#2d3748', // Slightly lighter for contrast
        },
        text: {
            primary: '#d1d5db', // Tailwind 'gray-300'
            secondary: '#9ca3af', // Tailwind 'gray-400'
        },
    },
    typography: {
        fontFamily: 'Roboto, sans-serif',
        h6: {
            fontWeight: 600,
        },
    },
    components: {
        MuiCard: {
            styleOverrides: {
                root: {
                    transition: 'box-shadow 0.3s ease-in-out, transform 0.3s ease-in-out',
                },
            },
        },
        MuiListItem: {
            styleOverrides: {
                root: {
                    borderRadius: 4,
                },
            },
        },
        MuiButton: {
            styleOverrides: {
                root: {
                    textTransform: 'none', // Prevent uppercase transformation
                },
            },
        },
        MuiSelect: {
            styleOverrides: {
                select: {
                    backgroundColor: 'background.paper',
                    borderRadius: 1,
                },
            },
        },
    },
});

export default theme;
