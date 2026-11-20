import { describe, expect, it } from 'vitest';

import { cautionRows, cautionTexts, MAX_CAUTION_LEN } from '../cautions';

import {
    formFromChallenge,
    mapServerErrors,
    snapshot,
    toPayload,
    validateForm,
} from '@/views/admin/challengeForm';

describe('uwagi — wiersze edytora', () => {
    it('lista z API → wiersze z unikalnym id → lista do zapisu bez pustych', () => {
        const rows = cautionRows(['a', 'b']);
        expect(rows.map((r) => r.text)).toEqual(['a', 'b']);
        expect(new Set(rows.map((r) => r.id)).size).toBe(2);
        rows.push({ id: 999, text: '   ' });
        rows[0]!.text = '  a  ';
        expect(cautionTexts(rows)).toEqual(['a', 'b']);
        expect(cautionRows(undefined)).toEqual([]);
    });
});

describe('uwagi w formularzu wyzwania', () => {
    const base = {
        id: 1,
        title: 'T',
        formatId: 1,
        format: 'na temat',
        description: 'D',
        timeLimitMin: null,
        difficulty: 2,
        source: null,
        ingredients: [],
        createdAt: '',
        updatedAt: '',
    };

    it('payload: kolejność zachowana, puste pominięte; pusty wiersz nie zmienia stanu „niezapisane”', () => {
        const f = formFromChallenge({ ...base, cautions: ['x', 'y'] });
        const before = snapshot(f);
        f.cautions.push({ id: 1000, text: '' });
        expect(snapshot(f)).toBe(before);
        f.cautions.reverse();
        expect(toPayload(f).cautions).toEqual(['y', 'x']);
    });

    it('walidacja: ponad 220 znaków → błąd na wierszu', () => {
        const f = formFromChallenge({ ...base, cautions: ['ok', 'x'.repeat(MAX_CAUTION_LEN + 1)] });
        expect(validateForm(f)).toEqual({
            [`caution-${f.cautions[1]!.id}`]: 'Uwaga może mieć najwyżej 220 znaków',
        });
    });

    it('błąd z API cautions[i] trafia do i-tego niepustego wiersza', () => {
        const f = formFromChallenge({ ...base, cautions: ['a', 'b'] });
        f.cautions.splice(1, 0, { id: 2000, text: '' });
        const e = mapServerErrors([{ field: 'cautions[1]', message: 'za długa' }], f);
        expect(e).toEqual({ [`caution-${f.cautions[2]!.id}`]: 'za długa' });
    });
});
