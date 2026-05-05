import { render, screen } from '@testing-library/react';
import { describe, test, expect } from 'vitest';
import '@testing-library/jest-dom';
import MatrixShell from './MatrixShell';

describe('MatrixShell', () => {
    test('renders the frozen header slot, column headers, and rows', () => {
        render(
            <MatrixShell
                ariaLabel="My matrix"
                frozenHeader={<div data-testid="frozen-h">Char</div>}
                columnHeaders={['A', 'B'].map((p) => (
                    <div key={p} data-testid={`ch-${p}`}>{p}</div>
                ))}
                rows={[
                    {
                        key: 'r1',
                        frozenCell: <div data-testid="r1-f">Row1</div>,
                        cells: ['A', 'B'].map((p) => (
                            <div key={p} data-testid={`r1-${p}`}>{p}</div>
                        )),
                    },
                ]}
            />,
        );
        expect(screen.getByLabelText('My matrix')).toBeInTheDocument();
        expect(screen.getByTestId('frozen-h')).toBeInTheDocument();
        expect(screen.getByTestId('ch-A')).toBeInTheDocument();
        expect(screen.getByTestId('r1-f')).toBeInTheDocument();
        expect(screen.getByTestId('r1-B')).toBeInTheDocument();
    });
});
