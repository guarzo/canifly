import { render, screen } from '@testing-library/react';
import { describe, test, expect } from 'vitest';
import '@testing-library/jest-dom';
import ProgressRail from './ProgressRail';

describe('ProgressRail', () => {
    const segs = [
        { key: 'a', weight: 2, state: 'ready', label: 'A' },
        { key: 'b', weight: 1, state: 'training', label: 'B' },
        { key: 'c', weight: 1, state: 'idle', label: 'C' },
    ];

    test('renders one element per segment', () => {
        render(<ProgressRail segments={segs} ariaLabel="Plan readiness" />);
        const rail = screen.getByLabelText('Plan readiness');
        expect(rail.querySelectorAll('[data-segment]').length).toBe(3);
    });

    test('segment widths are proportional to weight (sums to 100%)', () => {
        const { container } = render(<ProgressRail segments={segs} ariaLabel="x" />);
        const widths = Array.from(container.querySelectorAll('[data-segment]'))
            .map((el) => parseFloat(el.style.width));
        expect(widths[0]).toBeCloseTo(50, 1);
        expect(widths[1]).toBeCloseTo(25, 1);
        expect(widths[2]).toBeCloseTo(25, 1);
    });

    test('renders an empty placeholder when segments is empty', () => {
        render(<ProgressRail segments={[]} ariaLabel="empty" />);
        expect(screen.getByLabelText('empty')).toBeInTheDocument();
    });
});
