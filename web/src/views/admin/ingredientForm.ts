// Ingredient form in the admin panel: field state (text, because of comma/period), validation as in the API
// (services/ingredients.go) and conversion to the payload.

import type { IIngredientPayload } from '@/data/api/ingredients/fetch';
import type { IngredientDto } from '@/data/dtos/ingredients/IngredientDto';
import type { IUnit } from '@/data/dtos/taxonomy/DictionariesDto';
import { type Macros, type NutrientKey, NUTRIENTS, type Range } from '@/shared/utils/nutrition';

export interface IRangeInput {
    min: string;
    max: string;
}

export interface IConvRow extends IRangeInput {
    /** Row key for v-for (not sent). */
    key: number;
    unitId: number | null;
}

export interface IIngredientForm {
    name: string;
    typeCode: string;
    parentId: number | null;
    allergens: string[];
    skipShopping: boolean;
    nutrition: Record<NutrientKey, IRangeInput>;
    conversions: IConvRow[];
}

export type RowState = { err: string } | { empty: true } | { val: Range };

/** Units converted automatically or without a quantity: no conversion (like fixedGrams / NoQty in the API). */
export const convertibleUnits = (units: IUnit[]) =>
    units.filter((u) => !u.noQty && u.name !== 'g' && u.name !== 'kg');

/** '' → null; comma or period; a negative number or garbage → NaN. */
export function parseNum(s: string): number | null {
    const t = s.trim().replace(',', '.');
    if (!t) return null;
    if (!/^\d+(\.\d+)?$/.test(t)) return NaN;
    return Number(t);
}

export function rowState(r: IRangeInput): RowState {
    const a = parseNum(r.min);
    const b = parseNum(r.max);
    if (Number.isNaN(a) || Number.isNaN(b)) return { err: 'Wpisz liczbę ≥ 0' };
    if (a == null && b != null) return { err: 'Podaj najpierw min' };
    if (a != null && b != null && b < a) return { err: 'Max nie może być mniejsze niż min' };
    if (a == null) return { empty: true };
    return { val: { min: a, max: b } };
}

const str = (v: number | null | undefined) => (v == null ? '' : String(v).replace('.', ','));

let rowKey = 0;
export const newConvRow = (unitId: number | null = null): IConvRow => ({
    key: ++rowKey,
    unitId,
    min: '',
    max: '',
});

export function formFromIngredient(ing?: IngredientDto): IIngredientForm {
    const nutrition = Object.fromEntries(
        NUTRIENTS.map(({ key }) => {
            const r = ing?.nutrition[key];
            return [key, { min: str(r?.min), max: str(r?.max) }];
        })
    ) as Record<NutrientKey, IRangeInput>;
    return {
        name: ing?.name ?? '',
        typeCode: ing?.type.code ?? '',
        parentId: ing?.parentId ?? null,
        allergens: [...(ing?.allergens ?? [])],
        skipShopping: ing?.skipShopping ?? false,
        nutrition,
        conversions: (ing?.conversions ?? [])
            .filter((c) => !c.inherited)
            .map((c) => ({ ...newConvRow(c.unitId), min: str(c.gramsMin), max: str(c.gramsMax) })),
    };
}

/** Form errors: keys `name`, `typeCode`, `nutrition.<key>`, `conversions[<i>]`. */
export function validateForm(f: IIngredientForm): Record<string, string> {
    const errs: Record<string, string> = {};
    if (!f.name.trim()) errs.name = 'Podaj nazwę';
    if (!f.typeCode) errs.typeCode = 'Wybierz typ';
    for (const { key } of NUTRIENTS) {
        const st = rowState(f.nutrition[key]);
        if ('err' in st) errs[`nutrition.${key}`] = st.err;
    }
    const seen = new Set<number>();
    f.conversions.forEach((c, i) => {
        const st = rowState(c);
        if (c.unitId == null) errs[`conversions[${i}]`] = 'Wybierz jednostkę';
        else if (seen.has(c.unitId)) errs[`conversions[${i}]`] = 'Ta jednostka ma już przelicznik';
        else if ('err' in st) errs[`conversions[${i}]`] = st.err;
        else if ('empty' in st) errs[`conversions[${i}]`] = 'Uzupełnij min';
        if (c.unitId != null) seen.add(c.unitId);
    });
    return errs;
}

/** API errors (422) → the same keys as validateForm (`nutrition.kcal.max` → `nutrition.kcal`). */
export function mapServerErrors(
    fields: { field: string; message: string }[]
): Record<string, string> {
    const out: Record<string, string> = {};
    for (const { field, message } of fields) {
        const key = field.replace(/\.(min|max|gramsMin|gramsMax|unitId)$/, '');
        out[key] ??= message;
    }
    return out;
}

export function toPayload(f: IIngredientForm): IIngredientPayload {
    const nutrition: Macros = {};
    for (const { key } of NUTRIENTS) {
        const st = rowState(f.nutrition[key]);
        if ('val' in st) nutrition[key] = st.val;
    }
    return {
        name: f.name.trim(),
        typeCode: f.typeCode,
        allergens: f.allergens,
        parentId: f.parentId,
        skipShopping: f.skipShopping,
        nutrition,
        conversions: f.conversions.flatMap((c) => {
            const st = rowState(c);
            return c.unitId != null && 'val' in st
                ? [{ unitId: c.unitId, gramsMin: st.val.min, gramsMax: st.val.max }]
                : [];
        }),
    };
}

/** Id of the ingredient and all its descendants: none of them can become its parent. */
export function selfAndDescendants(id: number, all: IngredientDto[]): Set<number> {
    const out = new Set([id]);
    let grew = true;
    while (grew) {
        grew = false;
        for (const i of all) {
            if (i.parentId != null && out.has(i.parentId) && !out.has(i.id)) {
                out.add(i.id);
                grew = true;
            }
        }
    }
    return out;
}

/** Range from sub-ingredients (lowest min, highest max), like ResolveMacros in the API; null = no children. */
export function inheritedFromChildren(children: IngredientDto[]): Macros | null {
    if (!children.length) return null;
    const out: Macros = {};
    for (const { key } of NUTRIENTS) {
        const vals = children.map((c) => c.resolvedNutrition[key]).filter((v): v is Range => !!v);
        if (!vals.length) continue;
        const min = Math.min(...vals.map((v) => v.min));
        const max = Math.max(...vals.map((v) => v.max ?? v.min));
        out[key] = { min, max: max > min ? max : null };
    }
    return out;
}
