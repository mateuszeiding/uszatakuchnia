import { describe, expect, it } from 'vitest';

import { skipAdmin } from '../analytics';

const view = (url: string) => ({ type: 'pageview' as const, url });

describe('skipAdmin', () => {
    it('drops admin panel events', () => {
        expect(skipAdmin(view('https://example.com/admin'))).toBeNull();
        expect(skipAdmin(view('https://example.com/admin/recipes?x=1'))).toBeNull();
    });

    it('keeps public events untouched, query included', () => {
        const event = view('https://example.com/recipes/12?porcje=4&zamienniki=3');
        expect(skipAdmin(event)).toBe(event);
        expect(skipAdmin(view('https://example.com/recipes?q=admin'))).not.toBeNull();
    });
});
