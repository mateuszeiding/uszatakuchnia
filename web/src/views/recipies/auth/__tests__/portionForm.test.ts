import { describe, expect, it, vi } from 'vitest';

vi.mock('@/auth/createAuth', () => ({
    auth0: {},
    isOwner: { value: false },
    accessToken: () => Promise.resolve(null),
}));

import {
    blankSize,
    newPortion,
    pieceUnits,
    portionError,
    portionFromRecipe,
    portionToRequest,
} from '../portionForm';
import { newForm, validate } from '../recipeForm';

import { RecipeDto } from '@/data/dtos/recipe/RecipeDto';

const placek = {
    id: 20,
    name: 'placek',
    forms: ['placek', 'placki', 'placków'] as [string, string, string],
    noQty: false,
};

describe('formularz: czym jest porcja', () => {
    it('zapis i odczyt sztuk w rozmiarach', () => {
        const r = new RecipeDto({
            name: 'pizza',
            servings: 2,
            portion: {
                mode: 'sized',
                pieceUnit: placek,
                sizes: [
                    { diameterCm: 25, doughG: 190, isBase: false },
                    { diameterCm: 30, doughG: 255, isBase: true },
                ],
                form: null,
            },
        } as unknown as Partial<RecipeDto>);
        const p = portionFromRecipe(r);
        expect(p.mode).toBe('sized');
        expect(p.pieceUnitId).toBe(20);
        expect(portionToRequest(p)).toEqual({
            mode: 'sized',
            pieceUnitId: 20,
            sizes: [
                { diameterCm: 25, doughG: 190, isBase: false },
                { diameterCm: 30, doughG: 255, isBase: true },
            ],
            form: null,
        });
    });

    it('forma: wysyła tylko wymiary wybranego kształtu, pola innych trybów zostają', () => {
        const p = {
            ...newPortion(),
            mode: 'form' as const,
            shape: 'rect' as const,
            d: '24',
            a: '25',
            b: '35',
            pieces: '20',
            sizes: [{ ...blankSize(true), d: '30', g: '255' }],
        };
        expect(portionToRequest(p)).toEqual({
            mode: 'form',
            pieceUnitId: null,
            sizes: [{ diameterCm: 30, doughG: 255, isBase: true }],
            form: { shape: 'rect', dCm: null, aCm: 25, bCm: 35, pieces: 20 },
        });
    });

    it('stary przepis bez porcji = porcje', () => {
        const r = new RecipeDto({ name: 'x', servings: 4 } as Partial<RecipeDto>);
        expect(portionFromRecipe(r)).toEqual(newPortion());
        expect(portionToRequest(newPortion())).toEqual({
            mode: 'servings',
            pieceUnitId: null,
            sizes: [],
            form: null,
        });
    });

    it('walidacja przy publikacji', () => {
        const sized = { ...newPortion(), mode: 'sized' as const, pieceUnitId: 20 };
        expect(portionError(sized, true)).toMatch(/co najmniej jeden rozmiar/);
        expect(
            portionError({ ...sized, sizes: [{ ...blankSize(false), d: '30', g: '255' }] }, true)
        ).toMatch(/w przepisie/);
        expect(
            portionError({ ...sized, sizes: [{ ...blankSize(true), d: '65', g: '255' }] }, true)
        ).toMatch(/od 10 do 60/);
        expect(
            portionError({ ...sized, sizes: [{ ...blankSize(true), d: '30', g: '255' }] }, false)
        ).toMatch(/ile sztuk/);
        expect(
            portionError({ ...sized, sizes: [{ ...blankSize(true), d: '30', g: '255' }] }, true)
        ).toBeNull();
        expect(portionError({ ...newPortion(), mode: 'form' }, false)).toMatch(/średnicę/);
        expect(portionError({ ...newPortion(), mode: 'form', d: '24' }, false)).toBeNull();
        expect(
            portionError({ ...newPortion(), mode: 'form', shape: 'rect', a: '25' }, false)
        ).toMatch(/długość i szerokość/);
    });

    it('forma nie wymaga liczby porcji; porcje nadal tak', () => {
        const f = newForm();
        f.servings = '';
        expect(validate(f, undefined).errors.servings).toBeDefined();
        f.portion = { ...newPortion(), mode: 'form', d: '24' };
        expect(validate(f, undefined).errors.servings).toBeUndefined();
        expect(validate(f, undefined).errors.portion).toBeUndefined();
    });

    it('nazwa sztuki tylko z jednostek z odmianą', () => {
        const g = { id: 1, name: 'g', noQty: false };
        const smak = { id: 2, name: 'do smaku', noQty: true };
        expect(pieceUnits([g, smak, placek]).map((u) => u.name)).toEqual(['placek']);
    });
});
