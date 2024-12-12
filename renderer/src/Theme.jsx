// theme.js
import { createTheme } from '@mui/material/styles';

const theme = createTheme({
    palette: {
        primary: {
            main: '#14b8a6'
        },
        secondary: {
            main: '#ef4444'
        },
        info: {
            main: '#14b8a6', // Change this to teal as well
        },
        warning: { main:
                '#f59e0b'
        },
        background: {
            default: '#1f2937',
            paper: '#2d3748',
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
