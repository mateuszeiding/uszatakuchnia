// Ingredient refs in the step editor: rows to pick from, the label in the field and the cooking-mode limits.

import type { IRecipeForm } from './recipeForm';

import { RecipeIngredientDto } from '@/data/dtos/recipe/RecipeIngredientDto';
import type { IUnit } from '@/data/dtos/taxonomy/DictionariesDto';
import { scaleQuantity } from '@/shared/utils/quantity';
import {
    type IRefContext,
    partLabel,
    type RefPart,
    stepRefs,
    usedFraction,
} from '@/shared/utils/stepRefs';

/** Cooking mode fits a step of about this many characters on a phone (refs counted as shown). */
export const STEP_LIMIT = 280;
/** The counter under the step field appears from here. */
export const COUNT_FROM = 220;
export const HINT_LIMIT = 140;
export const HINTS_LIMIT = 3;

export interface IRefRow {
    key: string;
    name: string;
    group: string | null;
    /** Base quantity, e.g. "350 g"; '' = none. */
    qty: string;
    hasSubs: boolean;
}

/** Rows that can be referred to: picked from the database (others aren't saved). */
export function refRows(f: IRecipeForm, units: IUnit[]): IRefRow[] {
    return f.ingGroups.flatMap((g) =>
        g.items
            .filter((i) => i.ingredientId > 0)
            .map((i) => {
                const unit = units.find((u) => u.id === i.unitId);
                const a = scaleQuantity(i.qty || null, unit?.name ?? null, unit, 1);
                return {
                    key: i.key,
                    name: i.name,
                    group: g.label.trim() || null,
                    qty: a.qty ? `${a.qty}${a.unit ? ` ${a.unit}` : ''}` : '',
                    hasSubs: i.substitutes.length > 0,
                };
            })
    );
}

/** Context as on the recipe page at the default servings, without substitutes (for the length counter). */
export function formRefContext(f: IRecipeForm, units: IUnit[]): IRefContext {
    return {
        ingredients: f.ingGroups.flatMap((g) =>
            g.items
                .filter((i) => i.ingredientId > 0)
                .map(
                    (i) =>
                        new RecipeIngredientDto({
                            key: i.key,
                            section: g.label.trim() || null,
                            ingredientId: i.ingredientId,
                            ingredientName: i.name,
                            quantity: i.qty || null,
                            unitId: i.unitId,
                            unit: units.find((u) => u.id === i.unitId)?.name ?? null,
                            optional: i.optional,
                        })
                )
        ),
        steps: f.stepGroups.flatMap((g) => g.items.map((s) => s.text)),
        units,
        ratio: 1,
        subs: {},
    };
}

/** Rows whose fractions add up to more than the whole, or to the whole while "reszta" still refers to them. */
export function overusedKeys(f: IRecipeForm): Set<string> {
    const texts = f.stepGroups.flatMap((g) => g.items.map((s) => s.text));
    const refs = texts.flatMap((t) => stepRefs(t));
    const out = new Set<string>();
    for (const r of refs) {
        const used = usedFraction(r.key, texts);
        const rest = refs.some((x) => x.key === r.key && x.part === 'rest');
        if (used > 1 + 1e-9 || (rest && used >= 1 - 1e-9)) out.add(r.key);
    }
    return out;
}

/** Label of a ref in the field: "@mąka pszenna · Zaczyn", with a part "@1/2 · masło · Ciasto". No quantity (it depends on servings). */
export function chipLabel(row: IRefRow, part: RefPart | null): string {
    return `@${part ? `${partLabel(part)} · ` : ''}${row.name}${row.group ? ` · ${row.group}` : ''}`;
}

/** For filtering without Polish letters: "maka" finds "mąka". */
export const fold = (s: string) =>
    s.normalize('NFD').replace(/\p{M}/gu, '').replace(/ł/g, 'l').replace(/Ł/g, 'L').toLowerCase();
