// src/config.js

// Determine if the app is running in development mode.
// This depends on your build tools. For Vite, `import.meta.env.DEV` is used.
// If using Create React App, you might use `process.env.NODE_ENV !== 'production'`.
const isDev = import.meta.env.DEV;

// Define the back-end URL based on the environment.
// In production, you might adjust this URL accordingly.
const backEndURL = isDev ? '' : 'http://localhost:8713';

export { isDev, backEndURL };
