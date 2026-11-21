// Ingredient substitutes: the chosen recipe variant.
// The choice lives in uk:progress:<id>.subs (sortOrder → { ingredientId, mix }). A stale entry
// (the row or the substitute is gone from the recipe) is skipped; the recipe then looks as if nothing was chosen.

import type { ISubChoice } from '@/composables/useCookingProgress';
import type {
    IRecipeSubstitute,
    RecipeIngredientDto,
} from '@/data/dtos/recipe/RecipeIngredientDto';
import type { IRecipeNutrition, NutrientKey, SumRange } from '@/shared/utils/nutrition';

export type SubsMap = Record<number, ISubChoice>;

/** The row's chosen substitute, or null (main ingredient / stale entry). */
export function chosenSub(ing: RecipeIngredientDto, subs: SubsMap): IRecipeSubstitute | null {
    const c = subs[ing.sortOrder];
    if (!c) return null;
    return (
        ing.substitutes.find((s) => s.ingredientId === c.ingredientId && s.mix === c.mix) ?? null
    );
}

/** Current choices (sortOrder → substitute), without stale entries. */
export function activeSubs(
    items: RecipeIngredientDto[],
    subs: SubsMap
): Map<number, IRecipeSubstitute> {
    const out = new Map<number, IRecipeSubstitute>();
    for (const ing of items) {
        const s = chosenSub(ing, subs);
        if (s) out.set(ing.sortOrder, s);
    }
    return out;
}

const lowerFirst = (s: string) => s.charAt(0).toLowerCase() + s.slice(1);

/** Option name: the main ingredient, a substitute or a mix "Wołowina + wieprzowina mielona". */
export function optionName(ing: RecipeIngredientDto, sub: IRecipeSubstitute | null): string {
    if (!sub) return ing.ingredientName;
    return sub.mix ? `${ing.ingredientName} + ${lowerFirst(sub.name)}` : sub.name;
}

/** Option allergens (mix = main + substitute). */
export function optionAllergens(ing: RecipeIngredientDto, sub: IRecipeSubstitute | null): string[] {
    if (!sub) return ing.allergens;
    return sub.mix ? [...new Set([...ing.allergens, ...sub.allergens])] : sub.allergens;
}

/** Row in the chosen variant; factor 0.5 = half the quantity (mix). */
export type EffectiveLine = RecipeIngredientDto & { factor?: number };

/** Recipe rows in the chosen variant: the substitute instead of the main one, a mix = two rows at half each. */
export function effectiveLines(items: RecipeIngredientDto[], subs: SubsMap): EffectiveLine[] {
    const out: EffectiveLine[] = [];
    for (const ing of items) {
        const sub = chosenSub(ing, subs);
        if (!sub) {
            out.push(ing);
            continue;
        }
        const alt: EffectiveLine = {
            ...ing,
            ingredientId: sub.ingredientId,
            ingredientName: sub.name,
            allergens: sub.allergens,
            ingredientType: sub.ingredientType,
            skipShopping: sub.skipShopping,
            gramsPerUnit: sub.gramsPerUnit,
            linkedRecipe: null,
            substitutes: [],
        };
        if (sub.mix) out.push({ ...ing, factor: 0.5 }, { ...alt, factor: 0.5 });
        else out.push(alt);
    }
    return out;
}

/** Per-serving macros of the chosen variant: sum of the row shares (like services.Compute). No choice = macros from the API. */
export function variantNutrition(
    items: RecipeIngredientDto[],
    subs: SubsMap,
    base: IRecipeNutrition | null
): IRecipeNutrition | null {
    const active = activeSubs(items, subs);
    if (!active.size) return base;

    const per: Partial<Record<NutrientKey, SumRange>> = {};
    const missing: string[] = [];
    const incomplete: Partial<Record<NutrientKey, string[]>> = {};
    let contributed = false;
    for (const ing of items) {
        const sub = active.get(ing.sortOrder);
        const n = sub ? sub.nutrition : ing.nutrition;
        if (!n) continue;
        for (const name of n.missing) if (!missing.includes(name)) missing.push(name);
        for (const [k, names] of Object.entries(n.incomplete ?? {}) as [NutrientKey, string[]][]) {
            const list = (incomplete[k] ??= []);
            for (const name of names) if (!list.includes(name)) list.push(name);
        }
        for (const [k, v] of Object.entries(n.perPortion) as [NutrientKey, SumRange][]) {
            contributed = true;
            const s = (per[k] ??= { min: 0, max: 0 });
            s.min += v.min;
            s.max += v.max;
        }
    }
    return contributed ? { perPortion: per, missing, incomplete } : null;
}

/** "wieprzowina mielona zamiast wołowiny", for the sentence in the macro card. */
export function swapSummary(items: RecipeIngredientDto[], subs: SubsMap): string {
    const parts: string[] = [];
    for (const [sortOrder, sub] of activeSubs(items, subs)) {
        const ing = items.find((i) => i.sortOrder === sortOrder)!;
        parts.push(
            sub.mix
                ? `${lowerFirst(ing.ingredientName)} pół na pół z: ${lowerFirst(sub.name)}`
                : `${lowerFirst(sub.name)} zamiast: ${lowerFirst(ing.ingredientName)}`
        );
    }
    return parts.join(', ');
}

// ─── Substitutes in the shopping list URL (?zamienniki=) ───────────────────────
// Format: "<row sortOrder>:<substitute ingredientId>[m]", comma-separated, "m" = half and half,
// e.g. ?porcje=6&zamienniki=6:21,8:40m. "brak" = recipe without substitutes (the recipient sees the sender's variant
// even if they chose something themselves). A recipe without substitutes gets no parameter.

export const SUBS_NONE = 'brak';

/** Parameter value for the chosen variant; undefined = the recipe has no substitutes. */
export function subsToParam(items: RecipeIngredientDto[], subs: SubsMap): string | undefined {
    if (!items.some((i) => i.substitutes.length)) return undefined;
    const parts = [...activeSubs(items, subs)].map(
        ([sortOrder, s]) => `${sortOrder}:${s.ingredientId}${s.mix ? 'm' : ''}`
    );
    return parts.length ? parts.join(',') : SUBS_NONE;
}

/** Variant from the URL; null = no parameter or nothing usable (the device's memory stays).
 *  Rows missing from the recipe are skipped. */
export function parseSubsParam(
    raw: string | null | undefined,
    items: RecipeIngredientDto[]
): SubsMap | null {
    const s = (raw ?? '').trim();
    if (!s) return null;
    if (s === SUBS_NONE) return {};
    const out: SubsMap = {};
    for (const part of s.split(',')) {
        const m = part.trim().match(/^(\d+):(\d+)(m?)$/);
        if (!m) continue;
        const sortOrder = Number(m[1]);
        const choice = { ingredientId: Number(m[2]), mix: m[3] === 'm' };
        const ing = items.find((i) => i.sortOrder === sortOrder);
        if (ing && chosenSub(ing, { [sortOrder]: choice })) out[sortOrder] = choice;
    }
    return Object.keys(out).length ? out : null;
}
