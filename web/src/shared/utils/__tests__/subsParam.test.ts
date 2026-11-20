import { describe, expect, it } from 'vitest';

import { parseSubsParam, subsToParam } from '../substitutes';

import {
    type IRecipeSubstitute,
    RecipeIngredientDto,
} from '@/data/dtos/recipe/RecipeIngredientDto';

// ?zamienniki= w linku do listy zakupów (przegląd 18.11, #13)
const sub = (id: number, mix = false): IRecipeSubstitute => ({
    ingredientId: id,
    name: `s${id}`,
    mix,
    allergens: [],
    ingredientType: { code: 'MEAT', label: 'mięso i drób' },
    skipShopping: false,
    gramsPerUnit: null,
    hasNutrition: true,
    nutrition: null,
});
const row = (sortOrder: number, subs: IRecipeSubstitute[] = []) =>
    new RecipeIngredientDto({ sortOrder, ingredientId: 100 + sortOrder, substitutes: subs });
const items = [row(1), row(6, [sub(21), sub(21, true)]), row(8, [sub(40), sub(41)])];

describe('subsToParam', () => {
    it('przepis bez zamienników → bez parametru', () => {
        expect(subsToParam([row(1), row(2)], {})).toBeUndefined();
    });
    it('nic nie wybrano → „brak”', () => {
        expect(subsToParam(items, {})).toBe('brak');
    });
    it('wybór i pół na pół; nieaktualne wpisy pominięte', () => {
        expect(
            subsToParam(items, {
                6: { ingredientId: 21, mix: true },
                8: { ingredientId: 40, mix: false },
                9: { ingredientId: 1, mix: false },
            })
        ).toBe('6:21m,8:40');
    });
});

describe('parseSubsParam', () => {
    it('brak parametru → null (zostaje pamięć urządzenia)', () => {
        expect(parseSubsParam(undefined, items)).toBeNull();
        expect(parseSubsParam('', items)).toBeNull();
    });
    it('„brak” → przepis bez zamienników', () => {
        expect(parseSubsParam('brak', items)).toEqual({});
    });
    it('round trip', () => {
        const subs = { 6: { ingredientId: 21, mix: true }, 8: { ingredientId: 41, mix: false } };
        expect(parseSubsParam(subsToParam(items, subs), items)).toEqual(subs);
    });
    it('pozycje spoza przepisu pomijane; same śmieci → null', () => {
        expect(parseSubsParam('6:21,8:99,7:1,xx', items)).toEqual({
            6: { ingredientId: 21, mix: false },
        });
        expect(parseSubsParam('6:99,abc', items)).toBeNull();
    });
});
