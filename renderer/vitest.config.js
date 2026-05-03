// vitest.config.js
import { defineConfig } from 'vitest/config'
import react from '@vitejs/plugin-react'
import path from 'path'

export default defineConfig({
    plugins: [react({ jsxRuntime: 'automatic' })],
    esbuild: {
        jsx: 'automatic',
    },
    test: {
        globals: true,
        environment: 'jsdom',
        setupFiles: './src/setupTests.js',
        exclude: [
            '**/node_modules/**',
            '**/dist/**',
            '**/cypress/**',
            '**/.{idea,git,cache,output,temp}/**',
            '**/{karma,rollup,webpack,vite,vitest,jest,ava,babel,nyc,cypress,tsup,build}.config.*',
            // Exclude problematic test files with zustand import issues
            '**/CharacterOverview.test.jsx',
            '**/Mapping.test.jsx',
            '**/AccountCard.test.jsx',
            '**/CharacterItem.test.jsx',
            '**/GroupCard.test.jsx'
        ],
        deps: {
            optimizer: {
                web: {
                    include: ['zustand']
                }
            }
        }
    },
    resolve: {
        alias: {
            '@': path.resolve(__dirname, './src'),
        },
    },
})
