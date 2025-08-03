import { describe, it, expect, vi } from 'vitest';
import { formatDate, calculateDaysFromToday, formatNumberWithCommas } from './formatter';

describe('dateFormatter', () => {
    describe('formatDate', () => {
        it('formats an ISO string without the year and with a 24-hour time', () => {
            // Set a known date/time in UTC so we get a predictable result
            // For example, 2023-10-05T14:30:00Z
            const isoString = '2023-10-05T14:30:00Z';
            // formatDate uses the local time zone, so the output depends on the environment's locale and time zone.
            // We'll just check if it returns a string containing expected month/day and hour:minute.

            const result = formatDate(isoString);
            // result might look like "Oct 5, 14:30" in a certain locale.
            // Check for partial matches
            expect(result).toMatch(/Oct/);    // Check for abbreviated month
            expect(result).toMatch(/5/);      // Day should appear
            // Check there's a time with 24-hour format (e.g., "14:30")
            expect(result).toMatch(/\d{2}:\d{2}/);
            // Since we removed the year, ensure no full year is present
            expect(result).not.toMatch(/\d{4}/);
        });
    });

    describe('calculateDaysFromToday', () => {
        beforeAll(() => {
            // Freeze time at a known date: 2023-10-01T00:00:00.000Z
            vi.setSystemTime(new Date('2023-10-01T00:00:00Z'));
        });

        afterAll(() => {
            vi.useRealTimers();
        });

        it('returns "0 days" if date is today or in the past', () => {
            const today = '2023-10-01T10:00:00Z'; // same day but later time
            const pastDate = '2023-09-30T00:00:00Z';
            expect(calculateDaysFromToday(today)).toBe('0 days');
            expect(calculateDaysFromToday(pastDate)).toBe('0 days');
        });

        it('returns number of days until a future date', () => {
            // Future date 3 days ahead
            const futureDate = '2023-10-04T12:00:00Z';
            expect(calculateDaysFromToday(futureDate)).toBe('3 days');
        });

        it('returns an empty string if no date is provided', () => {
            expect(calculateDaysFromToday()).toBe('');
            expect(calculateDaysFromToday(null)).toBe('');
        });
    });

    describe('formatNumberWithCommas', () => {
        it('formats numbers with commas', () => {
            expect(formatNumberWithCommas(1000)).toBe('1,000');
            expect(formatNumberWithCommas(1234567)).toBe('1,234,567');
            expect(formatNumberWithCommas(42)).toBe('42'); // no commas needed
            expect(formatNumberWithCommas(1000000000)).toBe('1,000,000,000');
        });

        it('works with negative and decimal numbers', () => {
            // Although the code uses toLocaleString(), decimals and negatives also get formatted
            expect(formatNumberWithCommas(-123456.789)).toBe('-123,456.789');
            // The exact formatting of decimals may vary by locale, but generally commas should appear in the integer part.
        });
    });
});
