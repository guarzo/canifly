import { render, screen, fireEvent } from '@testing-library/react';
import { describe, test, expect, vi } from 'vitest';
import '@testing-library/jest-dom';
import { ThemeProvider } from '@mui/material/styles';
import theme from '../../Theme';
import PairSelect from './PairSelect';

const opts = [
    { value: 'a', primary: 'Alpha', secondary: '111' },
    { value: 'b', primary: 'Bravo', secondary: '222' },
];

function wrap(ui) {
    return render(<ThemeProvider theme={theme}>{ui}</ThemeProvider>);
}

describe('PairSelect', () => {
    test('renders the StatusDot with the given state', () => {
        wrap(
            <PairSelect
                state="ready"
                leftValue="a"
                leftOptions={opts}
                onLeftChange={() => {}}
                leftLabel="Left"
                rightValue="b"
                rightOptions={opts}
                onRightChange={() => {}}
                rightLabel="Right"
            />,
        );
        expect(screen.getByLabelText('Ready')).toBeInTheDocument();
    });

    test('calls onLeftChange when the left select changes', () => {
        const onLeftChange = vi.fn();
        wrap(
            <PairSelect
                state="idle"
                leftValue=""
                leftOptions={opts}
                onLeftChange={onLeftChange}
                leftLabel="Left"
                rightValue=""
                rightOptions={opts}
                onRightChange={() => {}}
                rightLabel="Right"
            />,
        );
        fireEvent.mouseDown(screen.getByLabelText('Left'));
        fireEvent.click(screen.getByText('Alpha'));
        expect(onLeftChange).toHaveBeenCalledWith('a');
    });

    test('shows option primary text in closed state when a value is selected', () => {
        wrap(
            <PairSelect
                state="queued"
                leftValue="a"
                leftOptions={opts}
                onLeftChange={() => {}}
                leftLabel="Left"
                rightValue=""
                rightOptions={opts}
                onRightChange={() => {}}
                rightLabel="Right"
            />,
        );
        // The closed Select should display the primary label ("Alpha"), not
        // the bare value ("a"). This is the behavior fix for the existing
        // SyncProfileRow which currently shows ID-only.
        expect(screen.getByText('Alpha')).toBeInTheDocument();
    });
});
