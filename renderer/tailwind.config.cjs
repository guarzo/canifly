// tailwind.config.js

module.exports = {
    darkMode: 'class', // Enables class-based dark mode
    content: [
        "./index.html",
        "./src/**/*.{js,jsx,ts,tsx}",
    ],
    theme: {
        extend: {
            colors: {
                eve: {
                    // Primary blues
                    blue: {
                        50: '#e6f3ff',
                        100: '#b3daff',
                        200: '#80c1ff',
                        300: '#4da8ff',
                        400: '#1a8fff',
                        500: '#0076e6',
                        600: '#005db3',
                        700: '#004480',
                        800: '#002b4d',
                        900: '#00121a'
                    },
                    // Accent teals
                    teal: {
                        50: '#e6fffb',
                        100: '#b3fff0',
                        200: '#80ffe5',
                        300: '#4dffda',
                        400: '#1affcf',
                        500: '#00e6b5',
                        600: '#00b38c',
                        700: '#008063',
                        800: '#004d3a',
                        900: '#001a11'
                    },
                    // Warning/danger oranges
                    orange: {
                        50: '#fff5e6',
                        100: '#ffe0b3',
                        200: '#ffcc80',
                        300: '#ffb74d',
                        400: '#ffa21a',
                        500: '#e68900',
                        600: '#b36b00',
                        700: '#804d00',
                        800: '#4d2f00',
                        900: '#1a1000'
                    }
                }
            },
            fontFamily: {
                'display': ['Orbitron', 'sans-serif'],
                'body': ['Inter', 'sans-serif'],
                'mono': ['JetBrains Mono', 'monospace'],
            },
            animation: {
                'glow': 'glow 2s ease-in-out infinite',
                'float': 'float 6s ease-in-out infinite',
                'pulse-slow': 'pulse 4s ease-in-out infinite',
                'shimmer': 'shimmer 1.5s ease-in-out infinite',
                'slide-in': 'slideIn 0.3s ease-out',
                'fade-in': 'fadeIn 0.5s ease-out',
            },
            keyframes: {
                glow: {
                    '0%, 100%': { opacity: '1' },
                    '50%': { opacity: '0.5' }
                },
                float: {
                    '0%, 100%': { transform: 'translateY(0)' },
                    '50%': { transform: 'translateY(-10px)' }
                },
                shimmer: {
                    '0%': { backgroundPosition: '-200px 0' },
                    '100%': { backgroundPosition: 'calc(200px + 100%) 0' }
                },
                slideIn: {
                    '0%': { transform: 'translateX(-100%)', opacity: '0' },
                    '100%': { transform: 'translateX(0)', opacity: '1' }
                },
                fadeIn: {
                    '0%': { opacity: '0' },
                    '100%': { opacity: '1' }
                }
            },
            backdropBlur: {
                xs: '2px',
            },
            boxShadow: {
                'glow': '0 0 20px rgba(20, 184, 166, 0.5)',
                'glow-lg': '0 0 40px rgba(20, 184, 166, 0.5)',
                'glass': '0 8px 32px 0 rgba(0, 0, 0, 0.37)',
                'neumorphic': '20px 20px 60px #0a0f1a, -20px -20px 60px #1a2940',
                'neumorphic-inset': 'inset 20px 20px 60px #0a0f1a, inset -20px -20px 60px #1a2940',
            },
            backgroundImage: {
                'gradient-radial': 'radial-gradient(var(--tw-gradient-stops))',
                'gradient-radial-at-t': 'radial-gradient(at top, var(--tw-gradient-stops))',
                'gradient-radial-at-b': 'radial-gradient(at bottom, var(--tw-gradient-stops))',
                'gradient-radial-at-l': 'radial-gradient(at left, var(--tw-gradient-stops))',
                'gradient-radial-at-r': 'radial-gradient(at right, var(--tw-gradient-stops))',
            }
        },
    },
    plugins: [],
};
