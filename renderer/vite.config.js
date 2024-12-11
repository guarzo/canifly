// https://vite.dev/config/
import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';

export default defineConfig({
  plugins: [react()],
  base: "./",
  server: {
    port: 5173,
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
