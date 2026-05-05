import { render, screen } from '@testing-library/react';
import { describe, test, expect } from 'vitest';
import '@testing-library/jest-dom';
import Connector from './Connector';

describe('Connector', () => {
    test('renders an svg path between from and to', () => {
        const { container } = render(
            <Connector from={{ x: 0, y: 10 }} to={{ x: 100, y: 10 }} />,
        );
        const path = container.querySelector('path');
        expect(path).not.toBeNull();
        // Cubic bezier through midpoint horizontally — the path string
        // includes both endpoints. We just check it starts at (0,10) and
        // ends at (100,10) — exact control point math is implementation
        // detail.
        expect(path.getAttribute('d')).toMatch(/M\s*0[ ,]\s*10/);
        expect(path.getAttribute('d')).toMatch(/100[ ,]\s*10\s*$/);
    });

    test('uses provided color', () => {
        const { container } = render(
            <Connector from={{ x: 0, y: 0 }} to={{ x: 10, y: 0 }} color="oklch(0.7 0.1 200)" />,
        );
        const path = container.querySelector('path');
        expect(path.getAttribute('stroke')).toBe('oklch(0.7 0.1 200)');
    });

    test('renders an arrow head when withArrow=true', () => {
        const { container } = render(
            <Connector from={{ x: 0, y: 0 }} to={{ x: 10, y: 0 }} withArrow />,
        );
        // Arrow renders as a small <polygon> at the to point.
        expect(container.querySelector('polygon')).not.toBeNull();
    });

    test('exposes role=presentation by default and an aria-label when given', () => {
        const { container, rerender } = render(
            <Connector from={{ x: 0, y: 0 }} to={{ x: 10, y: 0 }} />,
        );
        expect(container.querySelector('svg').getAttribute('role')).toBe('presentation');
        rerender(
            <Connector from={{ x: 0, y: 0 }} to={{ x: 10, y: 0 }} ariaLabel="paired" />,
        );
        expect(screen.getByLabelText('paired')).toBeInTheDocument();
    });
});
