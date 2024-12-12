// AccountCard.test.jsx
import React from 'react';
import { render, screen } from '@testing-library/react'
import AccountCard from './AccountCard'

test('renders account name and status', () => {
    const mockAccount = { Name: 'TestAccount', Status: 'Omega', Characters: [] };
    render(<AccountCard account={mockAccount} />);

    expect(screen.getByText('TestAccount')).toBeInTheDocument();
    expect(screen.getByText('Ω')).toBeInTheDocument();
});
