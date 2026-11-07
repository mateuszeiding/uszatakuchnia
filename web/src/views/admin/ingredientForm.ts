// Formularz składnika w panelu: stan pól (tekst, bo przecinek/kropka), walidacja jak w API
// (services/ingredients.go) i zamiana na payload.

import type { IIngredientPayload } from '@/data/api/ingredients/fetch';
import type { IngredientDto } from '@/data/dtos/ingredients/IngredientDto';
import type { IUnit } from '@/data/dtos/taxonomy/DictionariesDto';
import { type Macros, type NutrientKey, NUTRIENTS, type Range } from '@/shared/utils/nutrition';

export interface IRangeInput {
    min: string;
    max: string;
}

export interface IConvRow extends IRangeInput {
    /** Klucz wiersza w v-for (nie wysyłany). */
    key: number;
    unitId: number | null;
}

export interface IIngredientForm {
    name: string;
    typeCode: string;
    parentId: number | null;
    allergens: string[];
    nutrition: Record<NutrientKey, IRangeInput>;
    conversions: IConvRow[];
}

export type RowState = { err: string } | { empty: true } | { val: Range };

/** Jednostki przeliczane automatycznie albo bez ilości — bez przelicznika (jak fixedGrams / NoQty w API). */
export const convertibleUnits = (units: IUnit[]) =>
    units.filter((u) => !u.noQty && u.name !== 'g' && u.name !== 'kg');

/** '' → null; przecinek lub kropka; liczba ujemna lub śmieci → NaN. */
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
        nutrition,
        conversions: (ing?.conversions ?? [])
            .filter((c) => !c.inherited)
            .map((c) => ({ ...newConvRow(c.unitId), min: str(c.gramsMin), max: str(c.gramsMax) })),
    };
}

/** Błędy formularza: klucze `name`, `typeCode`, `nutrition.<key>`, `conversions[<i>]`. */
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

/** Błędy z API (422) → te same klucze co validateForm (`nutrition.kcal.max` → `nutrition.kcal`). */
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
        nutrition,
        conversions: f.conversions.flatMap((c) => {
            const st = rowState(c);
            return c.unitId != null && 'val' in st
                ? [{ unitId: c.unitId, gramsMin: st.val.min, gramsMax: st.val.max }]
                : [];
        }),
    };
}

/** Id składnika i wszystkich jego potomków — nie mogą zostać jego rodzicem. */
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

/** Zakres z podskładników (najniższe min, najwyższe max) — jak ResolveMacros w API; null = brak dzieci. */
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
