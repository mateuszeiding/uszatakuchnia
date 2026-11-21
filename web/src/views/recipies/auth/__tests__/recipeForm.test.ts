import { describe, expect, it, vi } from 'vitest';

// RecipeDto imports the API client, and createAuth builds the Auth0 client at import
// time, which needs a browser (window.location).
vi.mock('@/auth/createAuth', () => ({ auth0: {}, accessToken: () => Promise.resolve(null) }));

import {
    blankEquip,
    blankStep,
    fromRecipe,
    moveTipToStep,
    newForm,
    toRequest,
    validate,
} from '../recipeForm';

import { RecipeDto } from '@/data/dtos/recipe/RecipeDto';

describe('wskazówki w formularzu', () => {
    it('fromRecipe maps hints', () => {
        const r = new RecipeDto({
            name: 'x',
            steps: [
                { stepNo: 1, section: null, text: 'a', prep: false, prepLead: null, hints: ['h'] },
            ],
        } as unknown as Partial<RecipeDto>);
        expect(fromRecipe(r).stepGroups[0]!.items[0]!.hints).toEqual(['h']);
    });

    it('toRequest drops blank hints', () => {
        const f = newForm();
        f.stepGroups[0]!.items[0] = { ...blankStep(), text: 'krok', hints: [' a ', '', '  '] };
        expect(toRequest(f, 'draft').steps[0]!.hints).toEqual(['a']);
    });

    it('moveTipToStep appends and clears tip', () => {
        const f = newForm();
        const s = { ...blankStep(), text: 'krok', hints: ['x'] };
        f.stepGroups[0]!.items[0] = s;
        f.tip = ' rada ';
        moveTipToStep(f, s.id);
        expect(f.stepGroups[0]!.items[0]!.hints).toEqual(['x', 'rada']);
        expect(f.tip).toBe('');
    });

    it('moveTipToStep splits a tip over 300 chars on sentence boundaries', () => {
        const f = newForm();
        const s = { ...blankStep(), text: 'krok' };
        f.stepGroups[0]!.items[0] = s;
        const a = `${'a'.repeat(199)}.`;
        const b = `${'b'.repeat(149)}!`;
        const c = `${'c'.repeat(49)}?`;
        f.tip = `${a} ${b} ${c}`;
        moveTipToStep(f, s.id);
        expect(s.hints).toEqual([a, `${b} ${c}`]);
        expect(f.tip).toBe('');
        expect(validate(f, undefined).errors.steps).toBeUndefined();
    });

    it('validate flags hint over 300 chars', () => {
        const f = newForm();
        f.stepGroups[0]!.items[0] = { ...blankStep(), text: 'a' };
        f.stepGroups[0]!.items.push({ ...blankStep(), text: 'b', hints: ['ż'.repeat(301)] });
        expect(validate(f, undefined).errors.steps).toBe(
            'Krok 2: wskazówka może mieć najwyżej 300 znaków.'
        );
        f.stepGroups[0]!.items[1]!.hints = ['ż'.repeat(300)];
        expect(validate(f, undefined).errors.steps).toBeUndefined();
    });
});

describe('sprzęt w formularzu', () => {
    it('fromRecipe → toRequest zachowuje kolejność, doprecyzowanie i opcjonalny', () => {
        const r = new RecipeDto({
            name: 'x',
            equipment: [
                { equipmentId: 3, name: 'garnek', spec: '3–4 l', optional: false },
                { equipmentId: 1, name: 'robot', spec: null, optional: true },
            ],
        } as unknown as Partial<RecipeDto>);
        const f = fromRecipe(r);
        expect(f.equipment.map((e) => e.name)).toEqual(['garnek', 'robot']);
        expect(toRequest(f, 'draft').equipment).toEqual([
            { equipmentId: 3, spec: '3–4 l', optional: false },
            { equipmentId: 1, spec: null, optional: true },
        ]);
    });

    it('pusty wiersz pomijany, wpisany bez wyboru = błąd wiersza', () => {
        const f = newForm();
        const empty = blankEquip();
        const typed = { ...blankEquip(), name: 'termometr' };
        f.equipment = [empty, typed];
        expect(toRequest(f, 'draft').equipment).toEqual([]);
        const v = validate(f, undefined);
        expect(v.rows[typed.id]).toBeTruthy();
        expect(v.rows[empty.id]).toBeUndefined();
        expect(v.errors.equipment).toBe('Popraw sprzęt (1).');
    });
});

describe('zamienniki w formularzu', () => {
    it('fromRecipe → toRequest: zamiennik i pół na pół, bez zamiennika równego głównemu', () => {
        const r = new RecipeDto({
            name: 'x',
            ingredients: [
                {
                    sortOrder: 1,
                    section: null,
                    ingredientId: 10,
                    ingredientName: 'Wołowina',
                    quantity: '500',
                    unitId: 1,
                    optional: false,
                    substitutes: [
                        { ingredientId: 11, name: 'Wieprzowina', mix: false },
                        { ingredientId: 11, name: 'Wieprzowina', mix: true },
                    ],
                },
            ],
        } as unknown as Partial<RecipeDto>);
        const f = fromRecipe(r);
        const it = f.ingGroups[0]!.items[0]!;
        expect(it.substitutes).toEqual([
            { ingredientId: 11, name: 'Wieprzowina', mix: false },
            { ingredientId: 11, name: 'Wieprzowina', mix: true },
        ]);
        it.substitutes.push({ ingredientId: 10, name: 'Wołowina', mix: false });
        expect(toRequest(f, 'draft').ingredients[0]!.substitutes).toEqual([
            { ingredientId: 11, mix: false },
            { ingredientId: 11, mix: true },
        ]);
    });
});
