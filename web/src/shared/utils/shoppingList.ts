// Shopping list: adding up a recipe's ingredients.
// 1) the same ingredient (ingredientId) from different groups → one item;
// 2) the same unit → a sum (min and max separately);
// 3) different units → side by side: "550 g + 2 łyżki" (no conversion);
// 4) "do smaku" → only when there's nothing else;
// 5) ingredients with skipShopping (e.g. water) are skipped.
// Number formatting is shared with the recipe page (quantity.ts: fractions, thirds); in the store, pieces round up.

import type { RecipeIngredientDto } from '@/data/dtos/recipe/RecipeIngredientDto';
import type { IUnit } from '@/data/dtos/taxonomy/DictionariesDto';
import { fmtNum, isCountUnit, parseQty, shownAmount, unitLabel } from '@/shared/utils/quantity';

export interface IRange {
    min: number;
    max: number;
}

/** One occurrence of an ingredient in the recipe (scaled to the chosen servings). */
export interface IShopPart {
    group: string | null;
    unit: string;
    unitDef: IUnit | undefined;
    qty: IRange | null;
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
    /** Only optional occurrences → the "Opcjonalne" section. */
    optional: boolean;
    /** Sum of the required occurrences (for optional-only items: of the optional ones). */
    sum: ShopSum | null;
    parts: IShopPart[];
    groups: string[];
    /** Details ("mielona", "wędzony") from all occurrences, deduplicated, in recipe order. */
    notes: string[];
}

const GPU: Record<string, number> = { g: 1, kg: 1000 };

const isCount = (unit: string, def?: IUnit) => isCountUnit(unit, def);
const noQty = (p: { unitDef: IUnit | undefined }) => !!p.unitDef?.noQty;

/** Number for the list: pieces are rounded up in the store (you can't buy 2,5 eggs). */
function fmtN(v: number, unit: string, shop: boolean, def?: IUnit): string {
    const count = isCount(unit, def);
    if (shop && count) return String(Math.ceil(v - 0.001));
    return fmtNum(v, unit, count);
}
export function fmtQ(q: IRange | null, unit: string, shop = false, def?: IUnit): string {
    if (!q) return '';
    const a = fmtN(q.min, unit, shop, def);
    const b = fmtN(q.max, unit, shop, def);
    return a === b ? a : `${a}–${b}`;
}
const uLabel = (b: { unit: string; unitDef: IUnit | undefined }, n: number) =>
    unitLabel(b.unitDef, b.unit || null, n);
const join = (qty: string, unit: string) => [qty, unit].filter(Boolean).join(' ');

/** Recipe row → part of a shopping list item (quantity × ratio). */
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

    // Order: biggest share in grams first ("550 g + 2 łyżki"), no conversion.
    const g = (b: IBucket) => (b.gpu ? (((b.min + b.max) / 2) * (b.gpu.min + b.gpu.max)) / 2 : -1);
    buckets.sort((a, b) => g(b) - g(a) || b.max - a.max);
    const [main, ...rest] = buckets;
    return { main: main!, rest };
}

/** Recipe ingredients → shopping list items (in order of first occurrence). */
export function buildShopping(
    lines: (RecipeIngredientDto & { factor?: number })[],
    units: IUnit[],
    ratio: number
): IShopItem[] {
    const by = new Map<number, { ing: RecipeIngredientDto; parts: IShopPart[] }>();
    for (const x of lines) {
        if (x.skipShopping) continue;
        const k = x.ingredientId;
        if (!by.has(k)) by.set(k, { ing: x, parts: [] });
        // factor 0.5 = "half and half" with a substitute (substitutes.effectiveLines)
        by.get(k)!.parts.push(toPart(x, units, ratio * (x.factor ?? 1)));
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
            notes: [...new Set(parts.map((p) => p.note?.trim()).filter((x): x is string => !!x))],
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
                fmtQ(b, b.unit, true, b.unitDef),
                uLabel(
                    b,
                    isCount(b.unit, b.unitDef)
                        ? Math.ceil(b.max - 0.001)
                        : shownAmount(b.max, b.unit)
                )
            )
        )
        .join(' + ');
}

/** "What for": the recipe groups the ingredient appears in, without quantities.
 *  Groups where it's only optional → "(opcjonalnie)"; a recipe without groups → ''. */
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
