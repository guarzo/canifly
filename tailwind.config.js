module.exports = {
    darkMode: 'class', // Enable class-based dark mode
    content: [
        './internal/embed/templates/**/*.tmpl',
        './internal/embed/static/**/*.js',
        // Add other paths if necessary
    ],
    safelist: [
        'text-green-500',
        'text-red-500',
    ],
    theme: {
        extend: {},
    },
    plugins: [],
};
