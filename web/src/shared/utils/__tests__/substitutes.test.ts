import { describe, expect, it } from 'vitest';

import { buildShopping, sumText } from '../shoppingList';
import {
    activeSubs,
    effectiveLines,
    optionAllergens,
    optionName,
    swapSummary,
    variantNutrition,
} from '../substitutes';

import {
    type IRecipeSubstitute,
    RecipeIngredientDto,
} from '@/data/dtos/recipe/RecipeIngredientDto';
import type { IUnit } from '@/data/dtos/taxonomy/DictionariesDto';
import type { IRecipeNutrition } from '@/shared/utils/nutrition';

const units: IUnit[] = [{ id: 1, name: 'g', noQty: false }];
const kcal = (min: number, max: number, missing: string[] = []): IRecipeNutrition => ({
    perPortion: min || max ? { kcal: { min, max } } : {},
    missing,
    incomplete: {},
});
const sub = (id: number, name: string, x: Partial<IRecipeSubstitute> = {}): IRecipeSubstitute => ({
    ingredientId: id,
    name,
    mix: false,
    allergens: [],
    ingredientType: { code: 'MEAT', label: 'mięso i drób' },
    skipShopping: false,
    gramsPerUnit: { min: 1, max: 1 },
    hasNutrition: true,
    nutrition: null,
    ...x,
});

const beef = new RecipeIngredientDto({
    sortOrder: 1,
    ingredientId: 10,
    ingredientName: 'Wołowina',
    quantity: '500',
    unitId: 1,
    unit: 'g',
    gramsPerUnit: { min: 1, max: 1 },
    nutrition: kcal(200, 300),
    substitutes: [
        sub(11, 'Wieprzowina mielona', { nutrition: kcal(300, 350) }),
        sub(11, 'Wieprzowina mielona', { mix: true, nutrition: kcal(250, 325) }),
        sub(12, 'Indyk mielony', { hasNutrition: false, nutrition: kcal(0, 0, ['Indyk mielony']) }),
    ],
});
const pasta = new RecipeIngredientDto({
    sortOrder: 2,
    ingredientId: 20,
    ingredientName: 'Makaron z soczewicy',
    quantity: '400',
    unitId: 1,
    unit: 'g',
    gramsPerUnit: { min: 1, max: 1 },
    nutrition: kcal(100, 110),
    substitutes: [
        sub(21, 'Makaron orkiszowy', { allergens: ['gluten'], nutrition: kcal(120, 130) }),
    ],
});
const items = [beef, pasta];
const base = kcal(300, 410);

describe('zamienniki — wybrany wariant', () => {
    it('bez wyboru: makro z API, pozycje bez zmian', () => {
        expect(variantNutrition(items, {}, base)).toBe(base);
        expect(effectiveLines(items, {})).toEqual(items);
    });

    it('nieaktualny wpis (wiersz albo zamiennik zniknął) jest pomijany', () => {
        const subs = { 1: { ingredientId: 99, mix: false }, 7: { ingredientId: 11, mix: false } };
        expect(activeSubs(items, subs).size).toBe(0);
        expect(variantNutrition(items, subs, base)).toBe(base);
    });

    it('makro = suma wkładów wybranego wariantu', () => {
        const n = variantNutrition(items, { 1: { ingredientId: 11, mix: false } }, base)!;
        expect(n.perPortion.kcal).toEqual({ min: 400, max: 460 });
    });

    it('zamiennik bez makro trafia do „Nie wliczono”', () => {
        const n = variantNutrition(items, { 1: { ingredientId: 12, mix: false } }, base)!;
        expect(n.perPortion.kcal).toEqual({ min: 100, max: 110 });
        expect(n.missing).toEqual(['Indyk mielony']);
    });

    it('lista zakupów: mieszanka = dwie pozycje po połowie', () => {
        const lines = effectiveLines([beef], { 1: { ingredientId: 11, mix: true } });
        const shop = buildShopping(lines, units, 1);
        expect(shop.map((s) => [s.name, sumText(s.sum)])).toEqual([
            ['Wołowina', '250 g'],
            ['Wieprzowina mielona', '250 g'],
        ]);
    });

    it('lista zakupów: zamiennik zamiast głównego, skalowanie porcji działa', () => {
        const shop = buildShopping(
            effectiveLines(items, { 2: { ingredientId: 21, mix: false } }),
            units,
            2
        );
        expect(shop.map((s) => [s.name, sumText(s.sum)])).toEqual([
            ['Wołowina', '1000 g'],
            ['Makaron orkiszowy', '800 g'],
        ]);
    });

    it('nazwa, alergeny i opis wyboru', () => {
        const mix = beef.substitutes[1]!;
        expect(optionName(beef, mix)).toBe('Wołowina + wieprzowina mielona');
        expect(optionAllergens(pasta, pasta.substitutes[0]!)).toEqual(['gluten']);
        expect(swapSummary(items, { 2: { ingredientId: 21, mix: false } })).toBe(
            'makaron orkiszowy zamiast: makaron z soczewicy'
        );
    });
});
