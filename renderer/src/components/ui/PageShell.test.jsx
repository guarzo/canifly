import { render, screen } from '@testing-library/react';
import { describe, test, expect } from 'vitest';
import '@testing-library/jest-dom';
import PageShell from './PageShell';

describe('PageShell', () => {
    test('renders the title and meta', () => {
        render(
            <PageShell title="My Page" meta="3 things">
                <div>body</div>
            </PageShell>,
        );
        expect(screen.getByRole('heading', { name: 'My Page' })).toBeInTheDocument();
        expect(screen.getByText('3 things')).toBeInTheDocument();
    });

    test('renders header actions slot', () => {
        render(
            <PageShell title="P" actions={<button type="button">Do</button>}>
                <div>body</div>
            </PageShell>,
        );
        expect(screen.getByRole('button', { name: 'Do' })).toBeInTheDocument();
    });

    test('renders children in the body', () => {
        render(
            <PageShell title="P">
                <div data-testid="body" />
            </PageShell>,
        );
        expect(screen.getByTestId('body')).toBeInTheDocument();
    });
});
