import { describe, expect, it } from 'vitest';

import { resumeIndex, stepProgress, uncheckedText } from '../progress';

describe('uncheckedText — toast „Odznacz wszystko”', () => {
    it('odmiana i obie części', () => {
        expect(uncheckedText(3, 5)).toBe('Odznaczono 3 kroki i 5 składników.');
        expect(uncheckedText(1, 1)).toBe('Odznaczono 1 krok i 1 składnik.');
        expect(uncheckedText(12, 22)).toBe('Odznaczono 12 kroków i 22 składniki.');
    });
    it('gdy jedna liczba to 0 — tylko druga część', () => {
        expect(uncheckedText(0, 2)).toBe('Odznaczono 2 składniki.');
        expect(uncheckedText(5, 0)).toBe('Odznaczono 5 kroków.');
    });
});

describe('optional steps', () => {
    const steps = [1, 2, 3, 4, 5, 6, 7].map((n) => ({
        stepNo: n,
        optional: [3, 6, 7].includes(n),
    }));
    it('"N z M" counts mandatory steps only', () => {
        expect(stepProgress(steps, [1, 2, 3])).toEqual({ done: 2, total: 4 });
        // done without the optional ones
        expect(stepProgress(steps, [1, 2, 4, 5])).toEqual({ done: 4, total: 4 });
    });
    it('resumes after the last checked step, so a skipped optional one is not offered again', () => {
        expect(resumeIndex(steps, [])).toBe(-1);
        expect(resumeIndex(steps, [1, 2])).toBe(2); // optional step 3 is next, not skipped yet
        expect(resumeIndex(steps, [1, 2, 4])).toBe(4); // 3 skipped → step 5
        expect(resumeIndex(steps, [1, 2, 4, 5, 7])).toBe(6); // past the end → first open mandatory (none) → last
        expect(resumeIndex(steps, [7])).toBe(0); // only the last one checked → first open mandatory
    });
});
