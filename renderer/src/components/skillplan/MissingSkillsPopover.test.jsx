import { render, screen, fireEvent } from '@testing-library/react';
import { describe, test, expect } from 'vitest';
import '@testing-library/jest-dom';
import { ThemeProvider } from '@mui/material/styles';
import theme from '../../Theme';
import MissingSkillsPopover from './MissingSkillsPopover';

const wrap = (ui) => render(<ThemeProvider theme={theme}>{ui}</ThemeProvider>);

describe('MissingSkillsPopover', () => {
    test('opens on trigger click and lists each missing skill at level', () => {
        wrap(
            <MissingSkillsPopover
                planName="Hurricane"
                missingSkills={{ Gunnery: 5, 'Motion Prediction': 3 }}
                trigger={<button type="button">Show 2 missing</button>}
            />,
        );
        fireEvent.click(screen.getByRole('button', { name: /show 2 missing/i }));
        expect(screen.getByText('Gunnery')).toBeInTheDocument();
        expect(screen.getByText('Motion Prediction')).toBeInTheDocument();
        // Levels rendered as small caps roman numerals or numerics — assert
        // the numeric appears either way.
        expect(screen.getAllByText(/^[1-5]$/i).length).toBeGreaterThan(0);
    });

    test('renders nothing in the popover when missingSkills is empty', () => {
        wrap(
            <MissingSkillsPopover
                planName="Hurricane"
                missingSkills={{}}
                trigger={<button type="button">Open</button>}
            />,
        );
        fireEvent.click(screen.getByRole('button', { name: 'Open' }));
        expect(screen.getByText(/no missing skills/i)).toBeInTheDocument();
    });
});
