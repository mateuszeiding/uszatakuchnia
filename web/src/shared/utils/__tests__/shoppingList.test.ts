import { describe, expect, it } from 'vitest';

import { buildShopping, forText, sumText } from '../shoppingList';

import { RecipeIngredientDto } from '@/data/dtos/recipe/RecipeIngredientDto';
import type { IUnit } from '@/data/dtos/taxonomy/DictionariesDto';

const units: IUnit[] = [
    { id: 1, name: 'g', noQty: false },
    { id: 3, name: 'ml', noQty: false },
    { id: 5, name: 'łyżka', forms: ['łyżka', 'łyżki', 'łyżek'], noQty: false },
    { id: 6, name: 'łyżeczka', forms: ['łyżeczka', 'łyżeczki', 'łyżeczek'], noQty: false },
    { id: 7, name: 'szklanka', forms: ['szklanka', 'szklanki', 'szklanek'], noQty: false },
    { id: 10, name: 'sztuka', forms: ['sztuka', 'sztuki', 'sztuk'], noQty: false },
    { id: 16, name: 'do smaku', noQty: true },
];
const U = Object.fromEntries(units.map((u) => [u.name, u.id]));

let so = 0;
const line = (
    section: string | null,
    id: number,
    name: string,
    quantity: string | null,
    unit: string,
    x: Partial<RecipeIngredientDto> = {}
) =>
    new RecipeIngredientDto({
        sortOrder: ++so,
        section,
        ingredientId: id,
        ingredientName: name,
        quantity,
        unitId: U[unit],
        unit,
        gramsPerUnit: unit === 'g' ? { min: 1, max: 1 } : null,
        ...x,
    });
const item = (lines: RecipeIngredientDto[], ratio = 1, id?: number) => {
    const items = buildShopping(lines, units, ratio);
    return id == null ? items[0]! : items.find((i) => i.id === id)!;
};

describe('buildShopping / sumText', () => {
    it('ta sama jednostka → suma', () => {
        const it = item([
            line('Ciasto', 1, 'mąka', '500', 'g'),
            line('Kruszonka', 1, 'mąka', '50', 'g'),
        ]);
        expect(sumText(it.sum)).toBe('550 g');
        expect(it.groups).toEqual(['Ciasto', 'Kruszonka']);
    });

    it('różne jednostki obok siebie, bez przeliczania — także z przelicznikiem (16.11)', () => {
        const it = item([
            line('Zaczyn', 1, 'mąka', '2', 'łyżka', { gramsPerUnit: { min: 8, max: 10 } }),
            line('Ciasto', 1, 'mąka', '500', 'g'),
            line('Kruszonka', 1, 'mąka', '50', 'g'),
        ]);
        expect(sumText(it.sum)).toBe('550 g + 2 łyżki');
    });

    it('bez przelicznika → „100 g + 1 łyżeczka”', () => {
        const it = item([
            line('Kruszonka', 13, 'cynamon', '100', 'g'),
            line('Ciasto', 13, 'cynamon', '1', 'łyżeczka'),
        ]);
        expect(sumText(it.sum)).toBe('100 g + 1 łyżeczka');
    });

    it('sztuki zaokrąglane w górę', () => {
        const eggs = [
            line('Ciasto', 3, 'jajko', '2', 'sztuka', { gramsPerUnit: { min: 50, max: 65 } }),
            line('Do smarowania', 3, 'jajko', '1', 'sztuka', {
                gramsPerUnit: { min: 50, max: 65 },
            }),
        ];
        expect(sumText(item(eggs).sum)).toBe('3 sztuki');
        expect(sumText(item(eggs, 0.75).sum)).toBe('3 sztuki'); // 2,25 → 3
        expect(sumText(item(eggs, 1.5).sum)).toBe('5 sztuk'); // 4,5 → 5
    });

    it('tylko-opcjonalne → pozycja opcjonalna z własną sumą', () => {
        const it = item([line('Ciasto', 9, 'rodzynki', '50', 'g', { optional: true })]);
        expect(it.optional).toBe(true);
        expect(sumText(it.sum)).toBe('50 g');
        expect(forText(it)).toBe('ciasto');
    });

    it('wymagane + opcjonalne → suma wymaganych, w „do:” grupa z dopiskiem', () => {
        const withQty = item([
            line('Ciasto', 4, 'masło', '80', 'g'),
            line('Kruszonka', 4, 'masło', '30', 'g'),
            line('Do podania', 4, 'masło', '20', 'g', { optional: true }),
        ]);
        expect(withQty.optional).toBe(false);
        expect(sumText(withQty.sum)).toBe('110 g');
        expect(forText(withQty)).toBe('ciasto, kruszonka · do podania (opcjonalnie)');

        const taste = item([
            line('Ciasto', 4, 'masło', '80', 'g'),
            line('Do podania', 4, 'masło', null, 'do smaku', { optional: true }),
        ]);
        expect(sumText(taste.sum)).toBe('80 g');
        expect(forText(taste)).toBe('ciasto · do podania (opcjonalnie)');
    });

    it('„do smaku” → tylko gdy nic innego', () => {
        expect(sumText(item([line(null, 7, 'sól', null, 'do smaku')]).sum)).toBe('do smaku');
        const mixed = item([
            line('Ciasto', 7, 'sól', '1', 'łyżeczka'),
            line('Sos', 7, 'sól', null, 'do smaku'),
        ]);
        expect(sumText(mixed.sum)).toBe('1 łyżeczka');
    });

    it('porcje skalują sumy, ułamki jak na stronie przepisu', () => {
        const lines = [
            line('Ciasto', 7, 'sól', '1/2', 'łyżeczka'),
            line('Zaczyn', 2, 'mleko', '100', 'ml'),
        ];
        expect(sumText(item(lines, 1, 7).sum)).toBe('0,5 łyżeczki');
        expect(sumText(item(lines, 2, 7).sum)).toBe('1 łyżeczka');
        expect(sumText(item(lines, 1.5, 2).sum)).toBe('150 ml');
    });

    it('jedna pozycja na ingredientId, typ składnika jako dział', () => {
        const items = buildShopping(
            [
                line('A', 1, 'mąka', '100', 'g', {
                    ingredientType: { code: 'GRAIN', label: 'zboża i makarony' },
                }),
                line('B', 2, 'mleko', '100', 'ml'),
                line('C', 1, 'mąka', '50', 'g'),
            ],
            units,
            1
        );
        expect(items.map((i) => i.name)).toEqual(['mąka', 'mleko']);
        expect(items[0]!.type.code).toBe('GRAIN');
        expect(items[1]!.type.code).toBe('OTHER');
    });
});

describe('forText — „do: …”', () => {
    it('grupy bez ilości, małymi literami, bez powtórzeń', () => {
        const it = item([
            line('Zaczyn', 1, 'mąka', '2', 'łyżka'),
            line('Ciasto', 1, 'mąka', '500', 'g'),
            line('Kruszonka', 1, 'mąka', '50', 'g'),
            line('Ciasto', 1, 'mąka', '20', 'g'),
        ]);
        expect(forText(it)).toBe('zaczyn, ciasto, kruszonka');
    });

    it('grupa wymagana i opcjonalna naraz → bez dopisku', () => {
        const it = item([
            line('Ciasto', 4, 'masło', '80', 'g'),
            line('Ciasto', 4, 'masło', '20', 'g', { optional: true }),
        ]);
        expect(forText(it)).toBe('ciasto');
    });

    it('przepis bez grup → pusto', () => {
        expect(forText(item([line(null, 7, 'sól', '1', 'łyżeczka')]))).toBe('');
    });
});

describe('kolejność jednostek', () => {
    it('najpierw jednostka o największym udziale w gramach, bez „~”', () => {
        const it = buildShopping(
            [
                line('Zaczyn', 2, 'mleko', '100', 'ml', { gramsPerUnit: { min: 1, max: 1 } }),
                line('Ciasto', 2, 'mleko', '3/4', 'szklanka', {
                    gramsPerUnit: { min: 240, max: 250 },
                }),
            ],
            units,
            1
        )[0]!;
        expect(sumText(it.sum)).toBe('0,75 szklanki + 100 ml');
    });
});

describe('skipShopping', () => {
    it('składnik oznaczony „nie dodawaj do listy” (woda) jest pomijany', () => {
        const items = buildShopping(
            [
                line('Ciasto', 1, 'mąka', '500', 'g'),
                line('Ciasto', 30, 'woda', '240', 'ml', { skipShopping: true }),
                line('Zaczyn', 30, 'woda', '50', 'ml', { skipShopping: true }),
            ],
            units,
            1
        );
        expect(items.map((i) => i.name)).toEqual(['mąka']);
    });
});

describe('notes — doprecyzowanie (przegląd 18.11, #3)', () => {
    it('jedno wystąpienie → jego doprecyzowanie', () => {
        expect(item([line('Mięso', 2, 'wołowina', '450', 'g', { note: 'mielona' })]).notes).toEqual(
            ['mielona']
        );
    });
    it('różne doprecyzowania tego samego składnika → wszystkie, bez powtórzeń i pustych', () => {
        const it = item([
            line('Ciasto', 3, 'jajka', '1', 'sztuka', { note: 'duże' }),
            line('Dodatkowo', 3, 'jajka', '1', 'sztuka', { note: 'duże, do smarowania' }),
            line('Inne', 3, 'jajka', '1', 'sztuka', { note: ' duże ' }),
            line('Inne', 3, 'jajka', '1', 'sztuka', { note: '' }),
        ]);
        expect(it.notes).toEqual(['duże', 'duże, do smarowania']);
    });
    it('bez doprecyzowania → pusta lista', () => {
        expect(item([line(null, 4, 'sól', null, 'do smaku')]).notes).toEqual([]);
    });
});

describe('jednostki z panelu i łyżki (przegląd 18.11, #12)', () => {
    const extra: IUnit[] = [
        ...units,
        { id: 17, name: 'łodyga', forms: ['łodyga', 'łodygi', 'łodyg'], noQty: false },
    ];
    it('łodyga w sklepie w górę, jak sztuka', () => {
        const items = buildShopping(
            [
                new RecipeIngredientDto({
                    ingredientId: 9,
                    ingredientName: 'seler',
                    quantity: '2',
                    unitId: 17,
                    unit: 'łodyga',
                }),
            ],
            extra,
            8 / 6
        );
        expect(sumText(items[0]!.sum)).toBe('3 łodygi');
    });
    it('łyżki zsumowane do połówek, odmiana po zaokrągleniu', () => {
        expect(
            sumText(
                item(
                    [
                        line('Baza', 1, 'oliwa', '2', 'łyżka'),
                        line('Mięso', 1, 'oliwa', '1', 'łyżka'),
                    ],
                    8 / 6
                ).sum
            )
        ).toBe('4 łyżki');
        expect(sumText(item([line('Baza', 1, 'oliwa', '2', 'łyżka')], 2.4).sum)).toBe('5 łyżek');
        expect(sumText(item([line('Baza', 1, 'oliwa', '2', 'łyżka')], 8 / 6).sum)).toBe(
            '2,5 łyżki'
        );
    });
});
