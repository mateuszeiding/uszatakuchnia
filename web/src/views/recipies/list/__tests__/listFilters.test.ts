import { describe, expect, it } from 'vitest';

import {
    applyFilters,
    emptyFilters,
    facetCount,
    facetPresent,
    readQuery,
    timeBucket,
    toQuery,
} from '../listFilters';

import type { RecipeBaseDto } from '@/data/dtos/recipe/RecipeDto';
import type { ITaxonomyAxis } from '@/data/dtos/taxonomy/DictionariesDto';

// Facets and time ranges
const axis = (id: number, code: string, filterMode: 'and' | 'or'): ITaxonomyAxis =>
    ({
        id,
        code,
        label: code,
        multi: true,
        required: false,
        hasColor: false,
        filterMode,
        values: [],
    }) as unknown as ITaxonomyAxis;
const axes = [axis(1, 'dish', 'or'), axis(2, 'diet', 'and')];
const v = (id: number) => ({ id, label: `v${id}` });
const recipe = (
    id: number,
    timeMin: number,
    difficulty: number,
    dish: number[],
    diet: number[] = []
) =>
    ({
        id,
        name: `r${id}`,
        timeMin,
        difficulty,
        taxonomy: { dish: dish.map(v), diet: diet.map(v) },
        ingredientNames: [],
        substituteNames: [],
        needsPrep: false,
    }) as unknown as RecipeBaseDto;
// ragù: main course (3), 3 h, easy; bułki: bread (7), 3 h, medium, vegetarian (32)
const list = [recipe(6, 180, 1, [3]), recipe(5, 180, 2, [7], [32])];

describe('timeBucket i czas w URL', () => {
    it('przedziały rozłączne wg time_min', () => {
        expect([0, 15, 30, 31, 60, 61, 120, 121, 300].map(timeBucket)).toEqual([
            0, 30, 30, 60, 60, 120, 120, 121, 121,
        ]);
    });
    it('stare linki czas=15/30/60/120 → najbliższy przedział; „121+” = ponad 2 h', () => {
        const q = (czas: string) => readQuery({ czas }, axes).maxTime;
        expect([q('15'), q('30'), q('60'), q('120'), q('121+')]).toEqual([30, 30, 60, 120, 121]);
        expect(toQuery({ ...emptyFilters(), maxTime: 121 }).czas).toBe('121+');
    });
    it('filtr „ponad 2 h” bierze przepisy 3 h', () => {
        expect(applyFilters(list, { ...emptyFilters(), maxTime: 121 }, axes)).toHaveLength(2);
        expect(applyFilters(list, { ...emptyFilters(), maxTime: 120 }, axes)).toHaveLength(0);
    });
});

describe('fasety', () => {
    const dish = { kind: 'axis', code: 'dish' } as const;
    const diet = { kind: 'axis', code: 'diet' } as const;
    it('obecne tylko wartości, które ma któryś przepis', () => {
        expect(facetPresent(list, dish, 3)).toBe(true);
        expect(facetPresent(list, dish, 4)).toBe(false);
        expect(facetPresent(list, { kind: 'time' }, 121)).toBe(true);
        expect(facetPresent(list, { kind: 'time' }, 30)).toBe(false);
    });
    it('bez wyboru: licznik = przepisy z wartością', () => {
        const f = emptyFilters();
        expect(facetCount(list, f, axes, dish, 3)).toBe(1);
        expect(facetCount(list, f, axes, { kind: 'time' }, 121)).toBe(2);
    });
    it('po wyborze „danie główne”: oś OR liczy po dołożeniu, AND i inne wymiary zawężają', () => {
        const f = { ...emptyFilters(), axes: { dish: [3] } };
        expect(facetCount(list, f, axes, dish, 3)).toBe(1); // selected: how many results have it
        expect(facetCount(list, f, axes, dish, 7)).toBe(2); // OR: ragù + bułki
        expect(facetCount(list, f, axes, diet, 32)).toBe(0); // ragù isn't vegetarian → disabled
        expect(facetCount(list, f, axes, { kind: 'difficulty' }, 2)).toBe(0);
        expect(facetCount(list, f, axes, { kind: 'difficulty' }, 1)).toBe(1);
        expect(facetCount(list, f, axes, { kind: 'time' }, 121)).toBe(1);
    });
});
