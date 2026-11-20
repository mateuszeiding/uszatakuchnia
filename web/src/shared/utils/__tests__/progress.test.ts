import { describe, expect, it } from 'vitest';

import { uncheckedText } from '../progress';

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
