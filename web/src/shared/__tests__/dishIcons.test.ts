import { describe, expect, it } from 'vitest';

import { DISH_ICONS, dishIconOf } from '../dishIcons';

describe('dishIconOf', () => {
    it('returns the icon key stored on the taxonomy value', () => {
        expect(dishIconOf({ icon: 'zupa' })).toBe('zupa');
    });

    it('returns null without a value, without an icon or with an unknown key', () => {
        expect(dishIconOf(null)).toBeNull();
        expect(dishIconOf({})).toBeNull();
        expect(dishIconOf({ icon: 'pizza' })).toBeNull();
        expect(dishIconOf({ icon: 'toString' })).toBeNull();
    });
});

describe('DISH_ICONS', () => {
    it('has the nine dish-type icons, each with at least one shape', () => {
        expect(Object.keys(DISH_ICONS)).toHaveLength(9);
        for (const shapes of Object.values(DISH_ICONS)) expect(shapes.length).toBeGreaterThan(0);
    });
});
