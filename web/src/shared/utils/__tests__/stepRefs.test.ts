import { describe, expect, it } from 'vitest';

import {
    type IRefContext,
    parseStep,
    refText,
    resolveRef,
    usedFraction,
    viewSegments,
    visibleLength,
} from '../stepRefs';

import {
    type IRecipeSubstitute,
    RecipeIngredientDto,
} from '@/data/dtos/recipe/RecipeIngredientDto';
import type { IUnit } from '@/data/dtos/taxonomy/DictionariesDto';

const units: IUnit[] = [
    { id: 1, name: 'g', noQty: false },
    { id: 2, name: 'sztuka', forms: ['sztuka', 'sztuki', 'sztuk'], noQty: false },
    { id: 3, name: 'ząbek', forms: ['ząbek', 'ząbki', 'ząbków'], noQty: false },
    { id: 4, name: 'łyżka', forms: ['łyżka', 'łyżki', 'łyżek'], noQty: false },
    { id: 5, name: 'do smaku', noQty: true },
];
let so = 0;
const row = (
    key: string,
    ingredientId: number,
    name: string,
    quantity: string | null,
    unitId: number,
    x: Partial<RecipeIngredientDto> = {}
) =>
    new RecipeIngredientDto({
        sortOrder: ++so,
        key,
        ingredientId,
        ingredientName: name,
        quantity,
        unitId,
        unit: units.find((u) => u.id === unitId)!.name,
        ...x,
    });
const sub = (id: number, name: string, mix = false): IRecipeSubstitute => ({
    ingredientId: id,
    name,
    mix,
    allergens: [],
    ingredientType: { code: 'MEAT', label: 'mięso' },
    skipShopping: false,
    gramsPerUnit: null,
    hasNutrition: true,
    nutrition: null,
});

const ingredients = [
    row('mz', 10, 'Mąka pszenna', '150', 1, { section: 'Zaczyn' }),
    row('mc', 10, 'Mąka pszenna', '350', 1, { section: 'Ciasto' }),
    row('ms', 11, 'Masło', '50', 1, { section: 'Ciasto' }),
    row('jc', 12, 'Jajko', '1', 2, { section: 'Ciasto' }),
    row('se', 13, 'Sezam', '1', 4, { section: 'Do posmarowania', optional: true }),
    row('ma', 14, 'Marchew', '1–2', 2, { section: 'Ragù' }),
    row('cz', 15, 'Czosnek', '2', 3, { section: 'Ragù' }),
    row('sl', 16, 'Sól', null, 5, { section: 'Ragù' }),
    row('wo', 17, 'Wołowina mielona', '500', 1, {
        section: 'Ragù',
        substitutes: [sub(18, 'Wieprzowina mielona'), sub(18, 'Wieprzowina mielona', true)],
    }),
];
const steps = ['Dodaj połowę masła {{ing:ms|1/2}}.', 'Dodaj resztę masła {{ing:ms|rest}}.'];
const ctx = (x: Partial<IRefContext> = {}): IRefContext => ({
    ingredients,
    steps,
    units,
    ratio: 1,
    subs: {},
    ...x,
});
const show = (text: string, c = ctx()) =>
    parseStep(text)
        .map((s) => (typeof s === 'string' ? s : s.kind === 'ing' ? refText(resolveRef(s, c)) : ''))
        .join('');

describe('parseStep', () => {
    it('splits text and refs; plain text stays one segment', () => {
        expect(parseStep('Zwykły krok.')).toEqual(['Zwykły krok.']);
        expect(parseStep('Wsyp mąkę {{ing:mz}}, potem {{ing:ms|1/2}}.')).toEqual([
            'Wsyp mąkę ',
            { kind: 'ing', key: 'mz', part: null },
            ', potem ',
            { kind: 'ing', key: 'ms', part: '1/2' },
            '.',
        ]);
    });
    it('keeps unknown kinds and bad parts apart (rendered as nothing)', () => {
        expect(parseStep('{{term:x}}{{ing:a|5/8}}')).toEqual([
            { kind: 'other', raw: '{{term:x}}' },
            { kind: 'other', raw: '{{ing:a|5/8}}' },
        ]);
    });
});

describe('resolveRef', () => {
    it('adds the group only for an ingredient that is in more than one row', () => {
        expect(show('Wsyp mąkę {{ing:mz}} i mąkę {{ing:mc}}.')).toBe(
            'Wsyp mąkę  150 g · zaczyn i mąkę  350 g · ciasto.'
        );
        expect(show('Jajko {{ing:jc}}')).toBe('Jajko  1 sztuka');
    });
    it('scales with servings like the ingredient list (inflection, ranges, rounding)', () => {
        expect(show('{{ing:cz}}', ctx({ ratio: 1.5 }))).toBe(' 3 ząbki');
        expect(show('{{ing:ma}}')).toBe(' 1–2 sztuki');
        expect(show('{{ing:mc}}', ctx({ ratio: 2 }))).toBe(' 700 g · ciasto');
    });
    it('½ and "reszta" = whole minus the fractions used in other steps', () => {
        expect(usedFraction('ms', steps)).toBe(0.5);
        expect(show(steps[0]!)).toBe('Dodaj połowę masła  25 g.');
        expect(show(steps[1]!)).toBe('Dodaj resztę masła  25 g.');
    });
    it('optional row and a row without a quantity', () => {
        expect(show('Posyp sezamem {{ing:se}}.')).toBe('Posyp sezamem  1 łyżka (opcjonalnie).');
        expect(resolveRef({ kind: 'ing', key: 'sl', part: null }, ctx())).toBeNull();
        expect(resolveRef({ kind: 'ing', key: 'gone', part: null }, ctx())).toBeNull();
    });
    it('substitute and half-and-half from the chosen variant', () => {
        const pick = (mix: boolean) => ctx({ subs: { 9: { ingredientId: 18, mix } } });
        expect(show('{{ing:wo}}', pick(false))).toBe(' ⇄ wieprzowina mielona · 500 g');
        expect(show('{{ing:wo}}', pick(true))).toBe(' ⇄ wołowina + wieprzowina · 250 + 250 g');
    });
    it('title: name, group, part of the whole', () => {
        expect(resolveRef({ kind: 'ing', key: 'ms', part: '1/2' }, ctx())?.title).toBe(
            'Masło · Ciasto · ½ z 50 g'
        );
    });
});

describe('visibleLength', () => {
    it('counts refs as the reader sees them', () => {
        // "Dodaj mąkę 350 g · ciasto." = 26
        expect(visibleLength('Dodaj mąkę {{ing:mc}}.', ctx())).toBe(26);
        expect(visibleLength('Dopraw solą {{ing:sl}}.', ctx())).toBe('Dopraw solą.'.length);
    });
});

describe('viewSegments', () => {
    const text = (t: string) =>
        viewSegments(t, ctx())
            .map((s) => (typeof s === 'string' ? s : `[${s.qty}]`))
            .join('');
    it('a ref without a quantity takes the space before it when punctuation follows', () => {
        expect(text('Dopraw solą {{ing:sl}}.')).toBe('Dopraw solą.');
        expect(text('Dopraw solą {{ing:sl}} i pieprzem.')).toBe('Dopraw solą  i pieprzem.');
        expect(text('Dodaj solą {{ing:sl}}')).toBe('Dodaj solą');
        expect(text('Dodaj masło {{ing:ms}}.')).toBe('Dodaj masło [50 g].');
    });
});
