import '@testing-library/jest-dom';

// Mock window.scrollTo to avoid jsdom "Not implemented" errors
window.scrollTo = () => {};

// Provide a localStorage mock (jsdom's may be missing/non-functional in this env)
const createStorageMock = () => {
    let store = {};
    return {
        getItem: (key) => (key in store ? store[key] : null),
        setItem: (key, value) => { store[key] = String(value); },
        removeItem: (key) => { delete store[key]; },
        clear: () => { store = {}; },
        key: (i) => Object.keys(store)[i] ?? null,
        get length() { return Object.keys(store).length; },
    };
};
Object.defineProperty(window, 'localStorage', { value: createStorageMock(), configurable: true });
Object.defineProperty(window, 'sessionStorage', { value: createStorageMock(), configurable: true });
