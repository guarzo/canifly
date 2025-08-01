// vite.config.js
import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';

export default defineConfig({
  plugins: [react({
    // Explicitly set the JSX runtime to automatic
    jsxRuntime: 'automatic',
  })],
  base: "./",
  server: {
    port: 3113, // Changed from 5173 to 3000
    strictPort: true, // Fail if port is already in use
    proxy: {
      '/api': {
        target: 'http://localhost:8713',
        changeOrigin: true,
        secure: false,
      },
      'v2': {
        target: 'https://login.eveonline.com',
        changeOrigin: true,
        secure: true,
      }
    },
  },
});
