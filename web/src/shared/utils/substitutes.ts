// Zamienniki składników: wybrany wariant przepisu (decyzje 18.11, prototyp Claude Design zm-main.jsx).
// Wybór trzymany w uk:progress:<id>.subs (sortOrder → { ingredientId, mix }). Nieaktualny wpis
// (wiersz albo zamiennik zniknął z przepisu) jest pomijany — przepis wygląda wtedy jak bez wyboru.

import type { ISubChoice } from '@/composables/useCookingProgress';
import type {
    IRecipeSubstitute,
    RecipeIngredientDto,
} from '@/data/dtos/recipe/RecipeIngredientDto';
import type { IRecipeNutrition, NutrientKey, SumRange } from '@/shared/utils/nutrition';

export type SubsMap = Record<number, ISubChoice>;

/** Wybrany zamiennik wiersza albo null (składnik główny / nieaktualny wpis). */
export function chosenSub(ing: RecipeIngredientDto, subs: SubsMap): IRecipeSubstitute | null {
    const c = subs[ing.sortOrder];
    if (!c) return null;
    return (
        ing.substitutes.find((s) => s.ingredientId === c.ingredientId && s.mix === c.mix) ?? null
    );
}

/** Aktualne wybory (sortOrder → zamiennik), bez nieaktualnych wpisów. */
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

/** Nazwa opcji: składnik główny, zamiennik albo mieszanka „Wołowina + wieprzowina mielona”. */
export function optionName(ing: RecipeIngredientDto, sub: IRecipeSubstitute | null): string {
    if (!sub) return ing.ingredientName;
    return sub.mix ? `${ing.ingredientName} + ${lowerFirst(sub.name)}` : sub.name;
}

/** Alergeny opcji (mieszanka = główny + zamiennik). */
export function optionAllergens(ing: RecipeIngredientDto, sub: IRecipeSubstitute | null): string[] {
    if (!sub) return ing.allergens;
    return sub.mix ? [...new Set([...ing.allergens, ...sub.allergens])] : sub.allergens;
}

/** Pozycja w wybranym wariancie; factor 0,5 = połowa ilości (mieszanka). */
export type EffectiveLine = RecipeIngredientDto & { factor?: number };

/** Pozycje przepisu w wybranym wariancie: zamiennik zamiast głównego, mieszanka = dwie pozycje po połowie. */
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

/** Makro na porcję wybranego wariantu: suma wkładów pozycji (jak services.Compute). Bez wyboru = makro z API. */
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

/** „wieprzowina mielona zamiast wołowiny” — do zdania w karcie makro. */
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
