import { render, screen, fireEvent } from '@testing-library/react';
import { describe, test, expect, vi } from 'vitest';
import { useRef } from 'react';
import '@testing-library/jest-dom';
import FilterBar from './FilterBar';

function Harness({ value, onChange, ...rest }) {
    const ref = useRef(null);
    return <FilterBar inputRef={ref} value={value} onChange={onChange} {...rest} />;
}

describe('FilterBar', () => {
    test('renders the placeholder', () => {
        render(<Harness value="" onChange={() => {}} placeholder="Filter things" />);
        expect(screen.getByPlaceholderText('Filter things')).toBeInTheDocument();
    });

    test('shows kbd hint when empty', () => {
        render(<Harness value="" onChange={() => {}} placeholder="P" />);
        expect(screen.getByText('/')).toBeInTheDocument();
    });

    test('shows clear button when filled and calls onChange("") on click', () => {
        const onChange = vi.fn();
        render(<Harness value="abc" onChange={onChange} placeholder="P" />);
        fireEvent.click(screen.getByRole('button', { name: /clear/i }));
        expect(onChange).toHaveBeenCalledWith('');
    });

    test('renders a count slot when provided', () => {
        render(<Harness value="abc" onChange={() => {}} placeholder="P" count="3 matches" />);
        expect(screen.getByText('3 matches')).toBeInTheDocument();
    });

    test('calls onChange with new value on input', () => {
        const onChange = vi.fn();
        render(<Harness value="" onChange={onChange} placeholder="P" />);
        fireEvent.change(screen.getByPlaceholderText('P'), { target: { value: 'q' } });
        expect(onChange).toHaveBeenCalledWith('q');
    });
});
