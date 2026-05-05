import { render, screen } from '@testing-library/react';
import { describe, test, expect } from 'vitest';
import '@testing-library/jest-dom';
import SwatchBridge from './SwatchBridge';

describe('SwatchBridge', () => {
    test('renders one Connector per matched pair', () => {
        const containerRef = { current: { getBoundingClientRect: () => ({ left: 0, top: 0 }) } };
        const pairs = [
            {
                key: 'a',
                color: 'oklch(0.7 0.1 200)',
                fromRect: { right: 100, top: 50, height: 20 },
                toRect: { left: 200, top: 50, height: 20 },
            },
            {
                key: 'b',
                color: 'oklch(0.7 0.1 150)',
                fromRect: { right: 100, top: 100, height: 20 },
                toRect: { left: 200, top: 100, height: 20 },
            },
        ];
        const { container } = render(
            <SwatchBridge containerRef={containerRef} pairs={pairs} width={300} height={200} />,
        );
        expect(container.querySelectorAll('svg').length).toBe(1);
        expect(container.querySelectorAll('path').length).toBe(2);
    });

    test('renders nothing when pairs is empty', () => {
        const containerRef = { current: { getBoundingClientRect: () => ({ left: 0, top: 0 }) } };
        const { container } = render(
            <SwatchBridge containerRef={containerRef} pairs={[]} width={300} height={200} />,
        );
        expect(container.querySelectorAll('path').length).toBe(0);
    });
});
