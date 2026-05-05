import { render } from '@testing-library/react';
import { describe, test, expect } from 'vitest';
import '@testing-library/jest-dom';
import EveTypeIcon from './EveTypeIcon';

describe('EveTypeIcon', () => {
    test('renders an img with the evetech URL when conversions has the name', () => {
        const { container } = render(<EveTypeIcon name="Hurricane" conversions={{ Hurricane: 24702 }} />);
        const img = container.querySelector('img');
        expect(img.getAttribute('src')).toBe('https://images.evetech.net/types/24702/icon');
    });

    test('renders a placeholder when conversions does not have the name', () => {
        const { container } = render(<EveTypeIcon name="Unknown" conversions={{}} />);
        expect(container.querySelector('img')).toBeNull();
        expect(container.querySelector('span[aria-hidden]')).not.toBeNull();
    });

    test('respects size prop', () => {
        const { container } = render(<EveTypeIcon name="Hurricane" conversions={{ Hurricane: 24702 }} size={5} />);
        const img = container.querySelector('img');
        expect(img.className).toContain('h-5');
        expect(img.className).toContain('w-5');
    });
});
