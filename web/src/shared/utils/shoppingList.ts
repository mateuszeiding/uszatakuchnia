// Lista zakupów: sumowanie składników przepisu.
// Logika przeniesiona 1:1 z prototypu Claude Design sk-main.jsx (mergeParts, buildShopping, sumText, forText):
// 1) ten sam składnik (ingredientId) z różnych grup → jedna pozycja;
// 2) ta sama jednostka → suma (min i max osobno);
// 3) różne jednostki → obok siebie: „550 g + 2 łyżki” (bez przeliczania, 16.11);
// 4) „do smaku” → tylko gdy nic innego;
// 5) składniki z skipShopping (np. woda) pomijane.
// Formatowanie liczb wspólne ze stroną przepisu (quantity.ts: ułamki, tercje); w sklepie sztuki w górę.

import type { RecipeIngredientDto } from '@/data/dtos/recipe/RecipeIngredientDto';
import type { IUnit } from '@/data/dtos/taxonomy/DictionariesDto';
import { COUNT, fmtNum, parseQty, unitLabel } from '@/shared/utils/quantity';

export interface IRange {
    min: number;
    max: number;
}

/** Jedno wystąpienie składnika w przepisie (przeskalowane do wybranej liczby porcji). */
export interface IShopPart {
    group: string | null;
    unit: string;
    unitDef: IUnit | undefined;
    qty: IRange | null;
    /** Gramy na 1 jednostkę; null = nie da się przeliczyć. */
    gpu: IRange | null;
    optional: boolean;
    note: string | null;
}

export interface IBucket extends IRange {
    unit: string;
    unitDef: IUnit | undefined;
    gpu: IRange | null;
}

export type ShopSum = { taste: true; rest: [] } | { taste?: false; main: IBucket; rest: IBucket[] };

export interface IShopItem {
    id: number;
    name: string;
    type: { code: string; label: string };
    /** Tylko opcjonalne wystąpienia → sekcja „Opcjonalne”. */
    optional: boolean;
    /** Suma wymaganych wystąpień (dla pozycji tylko-opcjonalnych — opcjonalnych). */
    sum: ShopSum | null;
    parts: IShopPart[];
    groups: string[];
}

const GPU: Record<string, number> = { g: 1, kg: 1000 };

const isCount = (unit: string) => COUNT.has(unit);
const noQty = (p: { unitDef: IUnit | undefined }) => !!p.unitDef?.noQty;

/** Liczba do listy: w sklepie sztuki zaokrąglone w górę (nie kupisz 2½ jajka). */
function fmtN(v: number, unit: string, shop: boolean): string {
    if (shop && isCount(unit)) return String(Math.ceil(v - 0.001));
    return fmtNum(v, unit);
}
export function fmtQ(q: IRange | null, unit: string, shop = false): string {
    if (!q) return '';
    const a = fmtN(q.min, unit, shop);
    const b = fmtN(q.max, unit, shop);
    return a === b ? a : `${a}–${b}`;
}
const uLabel = (b: { unit: string; unitDef: IUnit | undefined }, n: number) =>
    unitLabel(b.unitDef, b.unit || null, n);
const join = (qty: string, unit: string) => [qty, unit].filter(Boolean).join(' ');

/** Pozycja przepisu → część listy zakupów (ilość × ratio). */
export function toPart(ing: RecipeIngredientDto, units: IUnit[], ratio: number): IShopPart {
    const unitDef = units.find((u) => u.id === ing.unitId);
    const q = parseQty(ing.quantity);
    return {
        group: ing.section,
        unit: unitDef?.name ?? ing.unit ?? '',
        unitDef,
        qty: q.ok && !q.empty ? { min: q.min * ratio, max: q.max * ratio } : null,
        gpu: ing.gramsPerUnit,
        optional: ing.optional,
        note: ing.note,
    };
}

export function mergeParts(parts: IShopPart[]): ShopSum | null {
    const byU: Record<string, IBucket> = {};
    const tasteOnly: IShopPart[] = [];
    for (const p of parts) {
        if (noQty(p) || !p.qty) {
            tasteOnly.push(p);
            continue;
        }
        const fixed = GPU[p.unit];
        const b = (byU[p.unit] ??= {
            unit: p.unit,
            unitDef: p.unitDef,
            min: 0,
            max: 0,
            gpu: p.gpu ?? (fixed ? { min: fixed, max: fixed } : null),
        });
        b.min += p.qty.min;
        b.max += p.qty.max;
    }
    const buckets = Object.values(byU);
    if (!buckets.length) return tasteOnly.length ? { taste: true, rest: [] } : null;

    // Kolejność: największy udział w gramach na początku („550 g + 2 łyżki”), bez przeliczania (16.11).
    const g = (b: IBucket) => (b.gpu ? (((b.min + b.max) / 2) * (b.gpu.min + b.gpu.max)) / 2 : -1);
    buckets.sort((a, b) => g(b) - g(a) || b.max - a.max);
    const [main, ...rest] = buckets;
    return { main: main!, rest };
}

/** Składniki przepisu → pozycje listy zakupów (kolejność pierwszego wystąpienia). */
export function buildShopping(
    lines: RecipeIngredientDto[],
    units: IUnit[],
    ratio: number
): IShopItem[] {
    const by = new Map<number, { ing: RecipeIngredientDto; parts: IShopPart[] }>();
    for (const x of lines) {
        if (x.skipShopping) continue; // np. woda — nie kupujemy
        const k = x.ingredientId;
        if (!by.has(k)) by.set(k, { ing: x, parts: [] });
        by.get(k)!.parts.push(toPart(x, units, ratio));
    }
    const items: IShopItem[] = [];
    by.forEach(({ ing, parts }) => {
        const req = parts.filter((p) => !p.optional);
        const opt = parts.filter((p) => p.optional);
        const base = req.length ? req : opt;
        items.push({
            id: ing.ingredientId,
            name: ing.ingredientName,
            type: ing.ingredientType ?? { code: 'OTHER', label: 'inne' },
            optional: !req.length,
            sum: mergeParts(base),
            parts,
            groups: [...new Set(parts.map((p) => p.group).filter((x): x is string => !!x))],
        });
    });
    return items;
}

export function sumText(s: ShopSum | null): string {
    if (!s) return '';
    if (s.taste) return 'do smaku';
    return [s.main, ...s.rest]
        .map((b) =>
            join(
                fmtQ(b, b.unit, true),
                uLabel(b, isCount(b.unit) ? Math.ceil(b.max - 0.001) : b.max)
            )
        )
        .join(' + ');
}

/** „Do czego”: grupy przepisu, w których składnik występuje, bez ilości (16.11).
 *  Grupy, w których jest tylko opcjonalnie → „(opcjonalnie)”; przepis bez grup → ''. */
export function forText(it: IShopItem): string {
    const groupsOf = (optional: boolean) => [
        ...new Set(
            it.parts
                .filter((p) => p.optional === optional && p.group)
                .map((p) => p.group!.toLowerCase())
        ),
    ];
    const req = groupsOf(false);
    const opt = groupsOf(true).filter((g) => !req.includes(g));
    return [
        req.join(', '),
        opt.length && !it.optional ? `${opt.join(', ')} (opcjonalnie)` : '',
        it.optional ? opt.join(', ') : '',
    ]
        .filter(Boolean)
        .join(' · ');
}
