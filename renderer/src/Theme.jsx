// theme.js
import { createTheme } from '@mui/material/styles';

const theme = createTheme({
    palette: {
        primary: {
            main: '#14b8a6',
            light: '#4dffda',
            dark: '#008063',
            contrastText: '#ffffff',
        },
        secondary: {
            main: '#0076e6',
            light: '#4da8ff',
            dark: '#004480',
            contrastText: '#ffffff',
        },
        error: {
            main: '#ef4444',
            light: '#fca5a5',
            dark: '#991b1b',
        },
        warning: { 
            main: '#f59e0b',
            light: '#fbbf24',
            dark: '#92400e',
        },
        info: {
            main: '#0891b2',
            light: '#22d3ee',
            dark: '#0c4a6e',
        },
        success: {
            main: '#10b981',
            light: '#34d399',
            dark: '#047857',
        },
        background: {
            default: 'rgba(17, 24, 39, 0.95)', // Semi-transparent for glass effect
            paper: 'rgba(31, 41, 59, 0.8)', // Glass morphism background
        },
        text: {
            primary: '#e5e7eb', // Tailwind 'gray-200'
            secondary: '#9ca3af', // Tailwind 'gray-400'
        },
        divider: 'rgba(148, 163, 184, 0.2)', // Semi-transparent divider
    },
    typography: {
        fontFamily: '"Inter", sans-serif',
        h1: {
            fontFamily: '"Orbitron", sans-serif',
            fontWeight: 800,
            letterSpacing: '-0.02em',
        },
        h2: {
            fontFamily: '"Orbitron", sans-serif',
            fontWeight: 700,
            letterSpacing: '-0.01em',
        },
        h3: {
            fontFamily: '"Orbitron", sans-serif',
            fontWeight: 600,
        },
        h4: {
            fontFamily: '"Orbitron", sans-serif',
            fontWeight: 600,
        },
        h5: {
            fontFamily: '"Orbitron", sans-serif',
            fontWeight: 500,
        },
        h6: {
            fontFamily: '"Orbitron", sans-serif',
            fontWeight: 500,
        },
        button: {
            fontWeight: 500,
            letterSpacing: '0.02em',
        },
        code: {
            fontFamily: '"JetBrains Mono", monospace',
        },
    },
    shape: {
        borderRadius: 12,
    },
    components: {
        MuiCssBaseline: {
            styleOverrides: {
                body: {
                    scrollbarColor: '#14b8a6 #1f2937',
                    '&::-webkit-scrollbar, & *::-webkit-scrollbar': {
                        width: 8,
                    },
                    '&::-webkit-scrollbar-thumb, & *::-webkit-scrollbar-thumb': {
                        borderRadius: 8,
                        backgroundColor: '#14b8a6',
                        minHeight: 24,
                    },
                    '&::-webkit-scrollbar-track, & *::-webkit-scrollbar-track': {
                        borderRadius: 8,
                        backgroundColor: '#1f2937',
                    },
                },
            },
        },
        MuiCard: {
            styleOverrides: {
                root: {
                    backgroundImage: 'none',
                    backgroundColor: 'rgba(31, 41, 59, 0.6)',
                    backdropFilter: 'blur(10px)',
                    border: '1px solid rgba(148, 163, 184, 0.1)',
                    transition: 'all 0.3s cubic-bezier(0.4, 0, 0.2, 1)',
                    '&:hover': {
                        transform: 'translateY(-4px)',
                        boxShadow: '0 20px 40px rgba(0, 0, 0, 0.3)',
                        borderColor: 'rgba(20, 184, 166, 0.3)',
                    },
                },
            },
        },
        MuiPaper: {
            styleOverrides: {
                root: {
                    backgroundImage: 'none',
                    backgroundColor: 'rgba(31, 41, 59, 0.8)',
                    backdropFilter: 'blur(10px)',
                },
            },
        },
        MuiButton: {
            styleOverrides: {
                root: {
                    textTransform: 'none',
                    fontWeight: 500,
                    borderRadius: 8,
                    padding: '8px 16px',
                    transition: 'all 0.3s ease',
                    '&:hover': {
                        transform: 'translateY(-2px)',
                        boxShadow: '0 8px 16px rgba(0, 0, 0, 0.2)',
                    },
                },
                contained: {
                    background: 'linear-gradient(135deg, #14b8a6 0%, #0891b2 100%)',
                    '&:hover': {
                        background: 'linear-gradient(135deg, #10b981 0%, #0ea5e9 100%)',
                    },
                },
                outlined: {
                    borderWidth: 2,
                    '&:hover': {
                        borderWidth: 2,
                    },
                },
            },
        },
        MuiIconButton: {
            styleOverrides: {
                root: {
                    transition: 'all 0.3s ease',
                    '&:hover': {
                        transform: 'scale(1.1)',
                        backgroundColor: 'rgba(20, 184, 166, 0.1)',
                    },
                },
            },
        },
        MuiTextField: {
            styleOverrides: {
                root: {
                    '& .MuiOutlinedInput-root': {
                        backgroundColor: 'rgba(17, 24, 39, 0.5)',
                        backdropFilter: 'blur(10px)',
                        '&:hover fieldset': {
                            borderColor: '#14b8a6',
                        },
                        '&.Mui-focused fieldset': {
                            borderColor: '#14b8a6',
                            borderWidth: 2,
                        },
                    },
                },
            },
        },
        MuiSelect: {
            styleOverrides: {
                root: {
                    backgroundColor: 'rgba(17, 24, 39, 0.5)',
                    backdropFilter: 'blur(10px)',
                    '&:hover .MuiOutlinedInput-notchedOutline': {
                        borderColor: '#14b8a6',
                    },
                    '&.Mui-focused .MuiOutlinedInput-notchedOutline': {
                        borderColor: '#14b8a6',
                        borderWidth: 2,
                    },
                },
            },
        },
        MuiListItem: {
            styleOverrides: {
                root: {
                    borderRadius: 8,
                    marginBottom: 4,
                    transition: 'all 0.3s ease',
                    '&:hover': {
                        backgroundColor: 'rgba(20, 184, 166, 0.1)',
                    },
                },
            },
        },
        MuiTooltip: {
            styleOverrides: {
                tooltip: {
                    backgroundColor: 'rgba(17, 24, 39, 0.95)',
                    backdropFilter: 'blur(10px)',
                    border: '1px solid rgba(148, 163, 184, 0.2)',
                    fontSize: '0.875rem',
                },
            },
        },
        MuiDivider: {
            styleOverrides: {
                root: {
                    borderColor: 'rgba(148, 163, 184, 0.2)',
                },
            },
        },
    },
});

export default theme;
