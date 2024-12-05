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
        background: {
            default: '#1f2937', // Tailwind 'gray-800'
            paper: '#1f2937',
        },
        text: {
            primary: '#d1d5db', // Tailwind 'gray-300'
            secondary: '#9ca3af', // Tailwind 'gray-400'
        },
    },
    typography: {
        fontFamily: 'Roboto, sans-serif',
    },
    components: {
        MuiCard: {
            styleOverrides: {
                root: {
                    // Ensure consistent card styling
                    transition: 'box-shadow 0.3s ease-in-out',
                },
            },
        },
        MuiListItem: {
            styleOverrides: {
                root: {
                    // Ensure consistent list item styling
                    borderRadius: 4,
                },
            },
        },
    },
});

export default theme;
