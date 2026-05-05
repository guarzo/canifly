import { render, screen, fireEvent } from '@testing-library/react';
import { describe, test, expect, vi } from 'vitest';
import '@testing-library/jest-dom';
import { ThemeProvider } from '@mui/material/styles';
import theme from '../../Theme';
import PlanList from './PlanList';

const wrap = (ui) => render(<ThemeProvider theme={theme}>{ui}</ThemeProvider>);

const skillPlans = {
    'Hurricane': {
        Name: 'Hurricane',
        Skills: { 'Gunnery': 5 },
        QualifiedCharacters: [],
        PendingCharacters: [],
        MissingCharacters: [],
    },
};

describe('PlanList copy action', () => {
    test('onCopy is called with only the source plan name (string)', () => {
        const onCopy = vi.fn();
        wrap(
            <PlanList
                skillPlans={skillPlans}
                characters={[]}
                conversions={{}}
                filter=""
                onCopy={onCopy}
                onDelete={() => {}}
            />,
        );
        fireEvent.click(screen.getByRole('button', { name: 'Copy Hurricane' }));
        expect(onCopy).toHaveBeenCalledTimes(1);
        const args = onCopy.mock.calls[0];
        expect(args.length).toBe(1);
        expect(args[0]).toBe('Hurricane');
    });
});
