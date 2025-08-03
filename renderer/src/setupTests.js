import '@testing-library/jest-dom';

// Mock window.scrollTo to avoid jsdom "Not implemented" errors
window.scrollTo = () => {};
