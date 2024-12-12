// src/utils/logger.jsx
import { isDev } from '../Config';

export function log(...args) {
    if (isDev) console.log(...args);
}

export function warn(...args) {
    if (isDev) console.warn(...args);
}

export function error(...args) {
    if (isDev) console.error(...args);
}

export function trace(...args) {
    if (isDev) console.trace(...args);
}
