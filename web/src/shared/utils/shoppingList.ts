// Lista zakupów: sumowanie składników przepisu.
// Logika przeniesiona 1:1 z prototypu .design/sk-main.jsx (mergeParts, buildShopping, sumText, optText):
// 1) ten sam składnik (ingredientId) z różnych grup → jedna pozycja;
// 2) ta sama jednostka → suma (min i max osobno);
// 3) różne jednostki → przeliczenie przez gramy na jednostkę o największym udziale, jeśli oba mają przelicznik → „~”;
// 4) bez przelicznika → zostaje jako „+ 1 łyżeczka”; „do smaku” → tylko gdy nic innego.
// Formatowanie liczb wspólne ze stroną przepisu (quantity.ts: ułamki, tercje); w sklepie sztuki w górę.

import type { RecipeIngredientDto } from '@/data/dtos/recipe/RecipeIngredientDto';
import type { IUnit } from '@/data/dtos/taxonomy/DictionariesDto';
import { COUNT, fmtNum, parseQty, SPOON, unitLabel } from '@/shared/utils/quantity';

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

export type ShopSum =
    | { taste: true; rest: [] }
    | { taste?: false; main: IBucket; approx: boolean; rest: IBucket[] };

export interface IShopItem {
    id: number;
    name: string;
    type: { code: string; label: string };
    /** Tylko opcjonalne wystąpienia → sekcja „Opcjonalne”. */
    optional: boolean;
    sum: ShopSum | null;
    /** Opcjonalne wystąpienia składnika, który jest też wymagany („+ 20 g (opcjonalnie)”). */
    optSum: ShopSum | null;
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

    const g = (b: IBucket) => (b.gpu ? (((b.min + b.max) / 2) * (b.gpu.min + b.gpu.max)) / 2 : -1);
    buckets.sort((a, b) => g(b) - g(a) || b.max - a.max);
    const main: IBucket = { ...buckets[0]! };
    let approx = false;
    const rest: IBucket[] = [];
    for (const b of buckets.slice(1)) {
        if (b.gpu && main.gpu) {
            main.min += (b.min * b.gpu.min) / main.gpu.max;
            main.max += (b.max * b.gpu.max) / main.gpu.min;
            approx = true;
        } else rest.push(b);
    }
    // sumy przybliżone: ≥ 100 co 10. Odstępstwo od prototypu: łyżki / szklanki nie do całości
    // (fmtNum zaokrągli do ćwiartek / tercji) — inaczej 1¼ szklanki wychodziło „~1 szklanka”.
    if (approx && !isCount(main.unit)) {
        const m = (main.min + main.max) / 2;
        const r = SPOON.has(main.unit) ? m : m >= 100 ? Math.round(m / 10) * 10 : Math.round(m);
        main.min = main.max = r;
    }
    return { main, approx, rest };
}

/** Składniki przepisu → pozycje listy zakupów (kolejność pierwszego wystąpienia). */
export function buildShopping(
    lines: RecipeIngredientDto[],
    units: IUnit[],
    ratio: number
): IShopItem[] {
    const by = new Map<number, { ing: RecipeIngredientDto; parts: IShopPart[] }>();
    for (const x of lines) {
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
            optSum: req.length && opt.length ? mergeParts(opt) : null,
            parts,
            groups: [...new Set(parts.map((p) => p.group).filter((x): x is string => !!x))],
        });
    });
    return items;
}

export function sumText(s: ShopSum | null): string {
    if (!s) return '';
    if (s.taste) return 'do smaku';
    const m = s.main;
    let t = `${s.approx ? '~' : ''}${join(fmtQ(m, m.unit, true), uLabel(m, isCount(m.unit) ? Math.ceil(m.max - 0.001) : m.max))}`;
    for (const r of s.rest)
        t += ` + ${join(fmtQ(r, r.unit, true), uLabel(r, isCount(r.unit) ? Math.ceil(r.max - 0.001) : r.max))}`;
    return t;
}

/** Dopisek dla składnika wymaganego, który bywa też opcjonalny. */
export function optText(it: IShopItem): string {
    if (!it.optSum) return '';
    if (it.optSum.taste) {
        const groups = [
            ...new Set(
                it.parts.filter((p) => p.optional && p.group).map((p) => p.group!.toLowerCase())
            ),
        ];
        return groups.length ? `też ${groups.join(', ')} — do smaku` : 'też do smaku';
    }
    return `+ ${sumText(it.optSum)}`;
}

/** Części, z których składa się suma (opcjonalne tylko dla pozycji tylko-opcjonalnych). */
export const countedParts = (it: IShopItem) => it.parts.filter((p) => it.optional || !p.optional);

export interface IBreakdownRow {
    group: string;
    text: string;
    optional: boolean;
    note: string | null;
}

/** Rozwinięcie pozycji: „grupa · ilość (≈ przeliczenie) · notatka”. */
export function breakdown(it: IShopItem): IBreakdownRow[] {
    const mu = it.sum && !it.sum.taste ? it.sum.main.unit : null;
    return countedParts(it).map((p) => {
        const t = noQty(p) || !p.qty ? 'do smaku' : join(fmtQ(p.qty, p.unit), uLabel(p, p.qty.max));
        const conv =
            p.gpu && mu && p.unit !== 'g' && p.qty && mu !== p.unit
                ? ` ≈ ${Math.round((p.qty.min * p.gpu.min + p.qty.max * p.gpu.max) / 2)} ${mu === 'ml' ? 'ml' : 'g'}`
                : '';
        return { group: p.group ?? '—', text: t + conv, optional: p.optional, note: p.note };
    });
}

// Dopełniacz l.mn. jednostki („łyżeczek”, „gramów”) do zdania o braku przelicznika.
const GEN_PL: Record<string, string> = {
    g: 'gramów',
    kg: 'kilogramów',
    ml: 'mililitrów',
    l: 'litrów',
};
const genPl = (b: IBucket) => GEN_PL[b.unit] ?? b.unitDef?.forms?.[2] ?? b.unit;

/** „Łyżeczek nie da się dodać do gramów — ten składnik nie ma przelicznika.”; null, gdy wszystko zsumowane. */
export function noConversionNote(it: IShopItem): string | null {
    const s = it.sum;
    if (!s || s.taste || !s.rest.length) return null;
    const what = s.rest.map(genPl).join(' ani ');
    return `${what.charAt(0).toUpperCase()}${what.slice(1)} nie da się dodać do ${genPl(s.main)} — ten składnik nie ma przelicznika.`;
}
