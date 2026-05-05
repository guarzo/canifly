import { render, screen } from '@testing-library/react';
import { describe, test, expect } from 'vitest';
import '@testing-library/jest-dom';
import { ThemeProvider } from '@mui/material/styles';
import theme from '../../Theme';
import PlanTimeline from './PlanTimeline';

const wrap = (ui) => render(<ThemeProvider theme={theme}>{ui}</ThemeProvider>);

const characters = [
    { Character: { CharacterID: 1, CharacterName: 'Alpha', QualifiedPlans: { 'Hurricane': true } } },
    { Character: { CharacterID: 2, CharacterName: 'Bravo', PendingPlans: { 'Hurricane': true }, SkillQueue: [], MissingSkills: {} } },
    { Character: { CharacterID: 3, CharacterName: 'Charlie', MissingSkills: { 'Hurricane': { Gunnery: 5 } } } },
];

const skillPlans = {
    'Hurricane': {
        Name: 'Hurricane',
        QualifiedCharacters: ['Alpha'],
        PendingCharacters: ['Bravo'],
        MissingCharacters: ['Charlie'],
    },
};

describe('PlanTimeline', () => {
    test('renders one rail per plan', () => {
        wrap(<PlanTimeline characters={characters} skillPlans={skillPlans} conversions={{}} filter="" />);
        expect(screen.getByLabelText('Hurricane readiness')).toBeInTheDocument();
    });

    test('omits plans that do not match the filter', () => {
        wrap(<PlanTimeline characters={characters} skillPlans={skillPlans} conversions={{}} filter="Vexor" />);
        expect(screen.queryByLabelText('Hurricane readiness')).toBeNull();
        expect(screen.getByText(/no skill plans match/i)).toBeInTheDocument();
    });

    test('renders the plan name', () => {
        wrap(<PlanTimeline characters={characters} skillPlans={skillPlans} conversions={{}} filter="" />);
        expect(screen.getByText('Hurricane')).toBeInTheDocument();
    });

    test('shows "no plans" message when skillPlans is empty', () => {
        wrap(<PlanTimeline characters={characters} skillPlans={{}} conversions={{}} filter="" />);
        expect(screen.getByText(/no skill plans yet/i)).toBeInTheDocument();
    });
});
