import { describe, expect, it } from 'vitest';

import { scaleQuantity } from '../quantity';

import type { IUnit } from '@/data/dtos/taxonomy/DictionariesDto';

const unit = (name: string, forms?: [string, string, string]): IUnit => ({
    id: 0,
    name,
    forms,
    noQty: false,
});
const U = {
    g: unit('g'),
    ml: unit('ml'),
    l: unit('l'),
    lyzeczka: unit('łyżeczka', ['łyżeczka', 'łyżeczki', 'łyżeczek']),
    szklanka: unit('szklanka', ['szklanka', 'szklanki', 'szklanek']),
    sztuka: unit('sztuka', ['sztuka', 'sztuki', 'sztuk']),
    zabek: unit('ząbek', ['ząbek', 'ząbki', 'ząbków']),
};
const show = (raw: string, u: IUnit, ratio = 1) => {
    const s = scaleQuantity(raw, u.name, u, ratio);
    return `${s.qty} ${s.unit}`;
};

describe('scaleQuantity — ułamki', () => {
    it('znaki ułamków dla łyżek i szklanek', () => {
        expect(show('½', U.lyzeczka)).toBe('½ łyżeczki');
        expect(show('1,5', U.szklanka)).toBe('1½ szklanki');
        expect(show('3/4', U.szklanka)).toBe('¾ szklanki');
        expect(show('1/3', U.szklanka)).toBe('⅓ szklanki');
        expect(show('2/3', U.szklanka)).toBe('⅔ szklanki');
    });

    it('„1/2” i „1 1/2” formatowane także przy bazowej liczbie porcji', () => {
        expect(show('1/2', U.lyzeczka)).toBe('½ łyżeczki');
        expect(show('1 1/2', U.szklanka)).toBe('1½ szklanki');
        expect(show('1/2', U.sztuka)).toBe('½ sztuki');
    });

    it('tercje po przeskalowaniu', () => {
        expect(show('1', U.szklanka, 1 / 3)).toBe('⅓ szklanki');
        expect(show('1', U.szklanka, 4 / 3)).toBe('1⅓ szklanki');
        expect(show('1', U.szklanka, 2 / 3)).toBe('⅔ szklanki');
    });

    it('gramy, mililitry i litry bez ułamków', () => {
        expect(show('0,5', U.l)).toBe('0,5 l');
        expect(show('1/2', U.l)).toBe('0,5 l');
        expect(show('150', U.g)).toBe('150 g');
        expect(show('100', U.ml, 1 / 3)).toBe('33 ml');
        expect(show('1', U.l, 1 / 3)).toBe('0,3 l');
    });
});

describe('scaleQuantity — bez zmian', () => {
    it('odmiana jednostek', () => {
        expect(show('1', U.zabek)).toBe('1 ząbek');
        expect(show('3', U.zabek)).toBe('3 ząbki');
        expect(show('5', U.zabek)).toBe('5 ząbków');
        expect(show('2', U.lyzeczka)).toBe('2 łyżeczki');
    });

    it('sztuki zaokrąglane do całości albo zakresu', () => {
        expect(show('1', U.sztuka, 1.5)).toBe('1–2 sztuki');
        expect(show('2', U.sztuka, 2)).toBe('4 sztuki');
        expect(show('3', U.sztuka, 0.5)).toBe('1–2 sztuki');
    });

    it('zakresy „2–3”', () => {
        expect(show('2–3', U.zabek)).toBe('2–3 ząbki');
        expect(show('2-3', U.lyzeczka, 2)).toBe('4–6 łyżeczek');
        expect(show('100–150', U.g, 2)).toBe('200–300 g');
    });

    it('brak ilości i tekst nieparsowalny', () => {
        expect(scaleQuantity(null, 'do smaku', undefined, 2)).toEqual({
            qty: '',
            unit: 'do smaku',
        });
        expect(scaleQuantity('trochę', 'g', U.g, 2)).toEqual({ qty: 'trochę', unit: 'g' });
    });
});
