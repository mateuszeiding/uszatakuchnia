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
    it('połówki dziesiętnie, tercje i ćwiartki cyframi (bez znaków ½ ⅓ ¼)', () => {
        expect(show('½', U.lyzeczka)).toBe('0,5 łyżeczki');
        expect(show('1,5', U.szklanka)).toBe('1,5 szklanki');
        expect(show('3/4', U.szklanka)).toBe('3/4 szklanki');
        expect(show('1/3', U.szklanka)).toBe('1/3 szklanki');
        expect(show('2/3', U.szklanka)).toBe('2/3 szklanki');
    });

    it('„1/2” i „1 1/2” formatowane także przy bazowej liczbie porcji', () => {
        expect(show('1/2', U.lyzeczka)).toBe('0,5 łyżeczki');
        expect(show('1 1/2', U.szklanka)).toBe('1,5 szklanki');
        expect(show('1/2', U.sztuka)).toBe('0,5 sztuki');
    });

    it('część całkowita oddzielona spacją od ułamka', () => {
        expect(show('1 1/4', U.szklanka)).toBe('1 1/4 szklanki');
        expect(show('1¼', U.szklanka)).toBe('1 1/4 szklanki');
        expect(show('2 3/4', U.szklanka)).toBe('2 3/4 szklanki');
        expect(show('2½', U.szklanka)).toBe('2,5 szklanki');
    });

    it('zakresy z ułamkami', () => {
        expect(show('1/4–1/2', U.lyzeczka)).toBe('1/4–0,5 łyżeczki');
        expect(show('⅓–⅔', U.szklanka)).toBe('1/3–2/3 szklanki');
    });

    it('nigdy nie pokazuje znaków ½ ⅓ ¼ ⅔ ¾', () => {
        const units = [U.lyzeczka, U.szklanka, U.sztuka, U.l, U.g];
        for (const u of units)
            for (let k = 1; k <= 48; k++) expect(show('1', u, k / 12)).not.toMatch(/[½⅓¼⅔¾⅛⅜⅝⅞⅕]/);
    });

    it('tercje po przeskalowaniu', () => {
        expect(show('1', U.szklanka, 1 / 3)).toBe('1/3 szklanki');
        expect(show('1', U.szklanka, 4 / 3)).toBe('1 1/3 szklanki');
        expect(show('1', U.szklanka, 2 / 3)).toBe('2/3 szklanki');
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

describe('przegląd 18.11, #12', () => {
    const lodyga = unit('łodyga', ['łodyga', 'łodygi', 'łodyg']);
    const lyzka = unit('łyżka', ['łyżka', 'łyżki', 'łyżek']);
    const garsc = unit('garść', ['garść', 'garście', 'garści']);

    it('jednostka z odmianą spoza miarek (łodyga z panelu) zaokrąglana jak sztuka', () => {
        expect(show('2', lodyga, 8 / 6)).toBe('2–3 łodygi');
        expect(show('2', lodyga)).toBe('2 łodygi');
        expect(show('3', lodyga, 0.5)).toBe('1–2 łodygi');
        expect(show('1', lodyga, 0.5)).toBe('0,5 łodygi');
    });

    it('łyżki i łyżeczki od 1 w górę do połówek', () => {
        expect(show('2', lyzka, 8 / 6)).toBe('2,5 łyżki');
        expect(show('1', lyzka, 8 / 6)).toBe('1,5 łyżki');
        expect(show('1', U.lyzeczka, 1.75)).toBe('2 łyżeczki');
        expect(show('2', lyzka, 2.4)).toBe('5 łyżek');
        expect(show('1,5', lyzka)).toBe('1,5 łyżki');
    });

    it('poniżej 1 łyżki bez zmian: ćwiartki i tercje', () => {
        expect(show('1/4', U.lyzeczka)).toBe('1/4 łyżeczki');
        expect(show('1', lyzka, 2 / 3)).toBe('2/3 łyżki');
    });

    it('szklanka i garść dalej z tercjami; g bez zmian', () => {
        expect(show('2', U.szklanka, 4 / 3)).toBe('2 2/3 szklanki');
        expect(show('2', garsc, 4 / 3)).toBe('2 2/3 garście');
        expect(show('500', U.g, 8 / 6)).toBe('665 g');
    });
});
