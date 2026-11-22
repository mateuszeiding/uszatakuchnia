import { describe, expect, it } from 'vitest';

import {
    defaultSize,
    effectiveMode,
    fmtScale,
    type IPortionRecipe,
    parseSizeQuery,
    pieceCount,
    pieceMacroFactor,
    portionScale,
    resolveSize,
    sameSize,
    scaleWarning,
    sizeLabel,
    sizeToQuery,
} from '../portion';
import { scaleQuantity } from '../quantity';

const placek = {
    id: 20,
    name: 'placek',
    forms: ['placek', 'placki', 'placków'] as [string, string, string],
    noQty: false,
};

const pizza: IPortionRecipe = {
    servings: 2,
    portion: {
        mode: 'sized',
        pieceUnit: placek,
        sizes: [
            { diameterCm: 25, doughG: 190, isBase: false },
            { diameterCm: 30, doughG: 255, isBase: true },
            { diameterCm: 35, doughG: 340, isBase: false },
            { diameterCm: 40, doughG: 420, isBase: false },
        ],
        form: null,
    },
};

const szarlotka: IPortionRecipe = {
    servings: 0,
    portion: {
        mode: 'form',
        pieceUnit: null,
        sizes: [],
        form: { shape: 'round', dCm: 24, aCm: null, bCm: null, pieces: 12 },
    },
};

const blacha: IPortionRecipe = {
    servings: 0,
    portion: {
        mode: 'form',
        pieceUnit: null,
        sizes: [],
        form: { shape: 'rect', dCm: null, aCm: 25, bCm: 35, pieces: null },
    },
};

const round2 = (x: number) => Math.round(x * 100) / 100;

describe('portionScale', () => {
    it('porcje: jak dotąd', () => {
        expect(portionScale({ servings: 4 }, 6, null)).toBe(1.5);
        expect(portionScale({ servings: 4, portion: null }, 2, null)).toBe(0.5);
        expect(portionScale({ servings: 0 }, 3, null)).toBe(1);
    });

    it('pizza: sztuki × masa na sztukę / masa w przepisie', () => {
        expect(portionScale(pizza, 0, { n: 2, d: 30 })).toBe(1);
        expect(round2(portionScale(pizza, 0, { n: 3, d: 25 }))).toBe(1.12); // 570 / 510
        expect(round2(portionScale(pizza, 0, { n: 1, d: 40 }))).toBe(0.82); // 420 / 510
        expect(round2(portionScale(pizza, 0, { n: 4, d: 25 }))).toBe(1.49); // 760 / 510
    });

    it('forma: stosunek pól dna, okrągła ↔ prostokątna', () => {
        expect(fmtScale(portionScale(szarlotka, 0, { shape: 'round', d: 26 }))).toBe('×1,17');
        expect(fmtScale(portionScale(szarlotka, 0, { shape: 'round', d: 18 }))).toBe('×0,56');
        expect(fmtScale(portionScale(blacha, 0, { shape: 'round', d: 26 }))).toBe('×0,61');
        expect(fmtScale(portionScale(szarlotka, 0, { shape: 'rect', a: 25, b: 35 }))).toBe('×1,93');
    });

    it('szkic bez rozmiaru bazowego zachowuje się jak porcje', () => {
        const draft: IPortionRecipe = { servings: 2, portion: { ...pizza.portion!, sizes: [] } };
        expect(effectiveMode(draft)).toBe('servings');
        expect(portionScale(draft, 4, null)).toBe(2);
    });
});

describe('resolveSize', () => {
    it('rozmiar spoza listy autora → rozmiar z przepisu, liczba sztuk zostaje', () => {
        expect(resolveSize(pizza, { n: 3, d: 45 })).toEqual({ n: 3, d: 30 });
        expect(resolveSize(pizza, { n: 99, d: 25 })).toEqual({ n: 2, d: 25 });
        expect(resolveSize(pizza, undefined)).toEqual({ n: 2, d: 30 });
        expect(resolveSize(pizza, { shape: 'round', d: 26 })).toEqual({ n: 2, d: 30 });
    });

    it('forma spoza 10–60 cm albo innego rodzaju → forma z przepisu', () => {
        expect(resolveSize(szarlotka, { shape: 'round', d: 9 })).toEqual({ shape: 'round', d: 24 });
        expect(resolveSize(szarlotka, { shape: 'rect', a: 30, b: 61 })).toEqual({
            shape: 'round',
            d: 24,
        });
        expect(resolveSize(szarlotka, { n: 2, d: 30 })).toEqual({ shape: 'round', d: 24 });
        expect(resolveSize(szarlotka, { shape: 'rect', a: 30, b: 40 })).toEqual({
            shape: 'rect',
            a: 30,
            b: 40,
        });
    });

    it('porcje: brak rozmiaru', () => {
        expect(resolveSize({ servings: 4 }, { n: 2, d: 30 })).toBeNull();
    });
});

describe('URL', () => {
    it('?sztuki=&rozmiar= i ?forma=', () => {
        expect(sizeToQuery({ n: 3, d: 25 })).toEqual({ sztuki: '3', rozmiar: '25' });
        expect(sizeToQuery({ shape: 'round', d: 26 })).toEqual({ forma: 'o26' });
        expect(sizeToQuery({ shape: 'rect', a: 30, b: 40 })).toEqual({ forma: 'p30x40' });
        expect(sizeToQuery(null)).toEqual({});
    });

    it('parsowanie z powrotem i rozmiary spoza listy', () => {
        expect(parseSizeQuery(pizza, { sztuki: '3', rozmiar: '25' })).toEqual({ n: 3, d: 25 });
        expect(parseSizeQuery(pizza, { sztuki: '3', rozmiar: '45' })).toEqual({ n: 3, d: 30 });
        expect(parseSizeQuery(pizza, { sztuki: '4' })).toEqual({ n: 4, d: 30 });
        expect(parseSizeQuery(pizza, { porcje: '4' })).toBeNull();
        expect(parseSizeQuery(szarlotka, { forma: 'o26' })).toEqual({ shape: 'round', d: 26 });
        expect(parseSizeQuery(szarlotka, { forma: 'p30x40' })).toEqual({
            shape: 'rect',
            a: 30,
            b: 40,
        });
        expect(parseSizeQuery(szarlotka, { forma: 'o99' })).toEqual({ shape: 'round', d: 24 });
        expect(parseSizeQuery(szarlotka, { forma: 'x' })).toBeNull();
        expect(parseSizeQuery({ servings: 4 }, { forma: 'o26' })).toBeNull();
    });
});

describe('opisy i makro', () => {
    it('nazwy rozmiarów', () => {
        expect(sizeLabel(pizza, { n: 3, d: 25 })).toBe('3 placki 25 cm');
        expect(sizeLabel(pizza, { n: 5, d: 25 })).toBe('5 placków 25 cm');
        expect(sizeLabel(szarlotka, { shape: 'round', d: 26 })).toBe('tortownica Ø 26 cm');
        expect(sizeLabel(szarlotka, { shape: 'rect', a: 25, b: 35 })).toBe('blacha 25 × 35 cm');
    });

    it('ostrzeżenia poniżej 0,5 i powyżej 2', () => {
        expect(scaleWarning(0.49)).toMatch(/Mniej niż połowa/);
        expect(scaleWarning(2.01)).toMatch(/Ponad 2×/);
        expect(scaleWarning(1)).toBeNull();
    });

    it('makro na sztukę i liczba sztuk', () => {
        expect(round2(pieceMacroFactor(pizza, { n: 3, d: 40 }))).toBe(1.65); // 420 / 255
        expect(pieceCount(pizza, { n: 3, d: 40 })).toBe(3);
        expect(fmtScale(pieceMacroFactor(szarlotka, { shape: 'round', d: 26 }))).toBe('×1,17');
        expect(pieceCount(szarlotka, defaultSize(szarlotka))).toBe(12);
        expect(pieceCount(blacha, defaultSize(blacha))).toBe(1);
    });

    it('sameSize: prostokąt obrócony to ta sama forma', () => {
        expect(sameSize({ shape: 'rect', a: 25, b: 35 }, { shape: 'rect', a: 35, b: 25 })).toBe(
            true
        );
        expect(sameSize({ n: 2, d: 30 }, { n: 2, d: 25 })).toBe(false);
    });

    it('sztuki przy ×1,17 jak dziś: 2 jajka → 2–3 sztuki', () => {
        const s = portionScale(szarlotka, 0, { shape: 'round', d: 26 });
        const jajko = {
            id: 10,
            name: 'sztuka',
            forms: ['sztuka', 'sztuki', 'sztuk'] as [string, string, string],
            noQty: false,
        };
        const out = scaleQuantity('2', 'sztuka', jajko, s);
        expect(`${out.qty} ${out.unit}`).toBe('2–3 sztuki');
    });
});
