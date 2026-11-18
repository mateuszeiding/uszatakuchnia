import { describe, expect, it } from 'vitest';

import { blankStep, fromRecipe, moveTipToStep, newForm, toRequest, validate } from '../recipeForm';

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
