// src/utils/logger.jsx
import { isDev } from '../Config';

// Logger functions that respect development/production modes
export function log(...args) {
    if (isDev) console.log(...args);
}

export function warn(...args) {
    console.warn(...args); // Warnings shown in all environments
}

export function error(...args) {
    console.error(...args); // Errors shown in all environments
}

export function trace(...args) {
    if (isDev) console.trace(...args);
}

// Main logger object for consistency
export const logger = {
    debug: log,
    info: log,
    warn,
    error,
    trace
};

