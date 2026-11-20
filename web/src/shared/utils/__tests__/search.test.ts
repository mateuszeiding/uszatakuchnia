import { describe, expect, it } from 'vitest';

import { matches, normalize, prepareQuery } from '../search';

const m = (q: string, t: string) => matches(prepareQuery(q)!, t);

describe('normalize', () => {
    it('małe litery, bez diakrytyków, ł → l', () => {
        expect(normalize('Bułki ŻÓŁĆ Ragù')).toBe('bulki zolc ragu');
    });
});

describe('matches (przegląd 18.11, #8)', () => {
    it('bez polskich znaków i akcentów', () => {
        expect(m('bulki', 'Bułki do burgerów')).toBe(true);
        expect(m('ragu', 'Ragù alla bolognese')).toBe(true);
        expect(m('wolowina', 'Wołowina')).toBe(true);
        expect(m('Bułki', 'Bułki do burgerów')).toBe(true);
    });
    it('formy słowa po wspólnym początku', () => {
        expect(m('bulka', 'Bułki do burgerów')).toBe(true);
        expect(m('tymianek', 'Pieczony łosoś z cytryną i tymiankiem')).toBe(true);
        expect(m('pomidory', 'Passata pomidorowa')).toBe(true);
    });
    it('bez głupich trafień', () => {
        expect(m('cukier', 'Cukinia')).toBe(false);
        expect(m('naci', 'Seler naciowy')).toBe(true); // substring match
        expect(m('bul', 'Bułka')).toBe(true);
        expect(m('bulx', 'Bułka')).toBe(false);
        expect(m('kurczak', 'Ragù alla bolognese')).toBe(false);
    });
    it('kilka słów: każde musi pasować', () => {
        expect(m('ragu bolognese', 'Ragù alla bolognese')).toBe(true);
        expect(m('bolognese ragu', 'Ragù alla bolognese')).toBe(true);
        expect(m('ragu kurczak', 'Ragù alla bolognese')).toBe(false);
    });
    it('zapytanie z polskimi znakami porównywane dosłownie', () => {
        expect(m('mąka', 'Makaron spaghetti')).toBe(false);
        expect(m('maka', 'Makaron spaghetti')).toBe(true);
        expect(m('mąka', 'Mąka pszenna typ 550')).toBe(true);
        expect(m('Bułka', 'Bułki do burgerów')).toBe(true);
    });
    it('puste zapytanie → brak zapytania', () => {
        expect(prepareQuery('   ')).toBeNull();
    });
});
