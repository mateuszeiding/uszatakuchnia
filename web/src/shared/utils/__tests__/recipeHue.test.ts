import { describe, expect, it } from 'vitest';

import { recipeHue } from '../recipeHue';

describe('recipeHue', () => {
    it('is (id * 47) % 360', () => {
        expect(recipeHue(1)).toBe(47);
        expect(recipeHue(8)).toBe(16);
        expect(recipeHue(360)).toBe(0);
    });

    it('stays in 0–359', () => {
        for (const id of [0, 1, 7, 99, 1234, 99999]) {
            const h = recipeHue(id);
            expect(h).toBeGreaterThanOrEqual(0);
            expect(h).toBeLessThan(360);
        }
    });

    it('is the same for the same recipe', () => {
        expect(recipeHue(42)).toBe(recipeHue(42));
    });
});
