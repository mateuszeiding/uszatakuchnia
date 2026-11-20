// A draft needs only a name; publishing needs everything (name, ingredient, step, required axes, difficulty, servings, time).

import type { IUpsertRecipeRequest, RecipeDto, RecipeStatus } from '@/data/dtos/recipe/RecipeDto';
import type { ILinkedRecipe } from '@/data/dtos/recipe/RecipeIngredientDto';
import type { IDictionaries, IUnit } from '@/data/dtos/taxonomy/DictionariesDto';
import { cautionRows, cautionTexts, type ICautionRow } from '@/shared/utils/cautions';
import { parseQty } from '@/shared/utils/quantity';
import { defaultTimeUnit, fromMinutes, type TimeUnit, toMinutes } from '@/shared/utils/time';

let seq = 0;
export const uid = () => ++seq;

/** Longest step hint (as in the API). */
export const MAX_HINT_LEN = 300;

export const PREP_LEADS = ['30 min', '1 h', '2 h', '4 h', '12 h', 'dzień wcześniej'];

export interface IIngItem {
    id: number;
    ingredientId: number;
    /** Text in the autosuggest field; ingredientId > 0 only after picking from the database. */
    name: string;
    qty: string;
    unitId: number | null;
    note: string;
    optional: boolean;
    linked: ILinkedRef | null;
    /** Substitutes ("zamiast" or "pół na pół") in the order added. */
    substitutes: ISubItem[];
}

export interface ISubItem {
    ingredientId: number;
    name: string;
    mix: boolean;
}

export type ILinkedRef = Pick<ILinkedRecipe, 'id' | 'title' | 'status'>;

export interface IStepItem {
    id: number;
    text: string;
    prep: boolean;
    prepLead: string;
    /** Hints below the step text; empty ones are skipped on save. */
    hints: string[];
}

export interface IEquipItem {
    id: number;
    /** 0 = not picked from the dictionary yet (row skipped on save). */
    equipmentId: number;
    name: string;
    spec: string;
    optional: boolean;
}

export interface IGroup<T> {
    id: number;
    label: string;
    items: T[];
}

export interface IRecipeForm {
    name: string;
    tagline: string;
    tip: string;
    photoUrl: string | null;
    taxonomy: Record<string, number[]>;
    difficulty: number | null;
    /** "od" / "do" time as text in the timeUnit unit (e.g. "2,5" h). */
    time: string;
    timeMax: string;
    timeUnit: TimeUnit;
    servings: string;
    ingGroups: IGroup<IIngItem>[];
    stepGroups: IGroup<IStepItem>[];
    equipment: IEquipItem[];
    /** "Don't mess this up" notes; empty rows are skipped on save. */
    cautions: ICautionRow[];
}

export const blankIng = (): IIngItem => ({
    id: uid(),
    ingredientId: 0,
    name: '',
    qty: '',
    unitId: null,
    note: '',
    optional: false,
    linked: null,
    substitutes: [],
});

export const blankEquip = (): IEquipItem => ({
    id: uid(),
    equipmentId: 0,
    name: '',
    spec: '',
    optional: false,
});

export const blankStep = (): IStepItem => ({
    id: uid(),
    text: '',
    prep: false,
    prepLead: '',
    hints: [],
});

export const newForm = (): IRecipeForm => ({
    name: '',
    tagline: '',
    tip: '',
    photoUrl: null,
    taxonomy: {},
    difficulty: null,
    time: '',
    timeMax: '',
    timeUnit: 'min',
    servings: '4',
    ingGroups: [{ id: uid(), label: '', items: [blankIng()] }],
    stepGroups: [{ id: uid(), label: '', items: [blankStep()] }],
    equipment: [],
    cautions: [],
});

/** Groups consecutive items with the same section (API order preserved). */
function groupBy<T, R>(
    items: T[],
    section: (x: T) => string | null,
    map: (x: T) => R
): IGroup<R>[] {
    const out: IGroup<R>[] = [];
    for (const it of items) {
        const label = section(it) ?? '';
        const last = out[out.length - 1];
        if (last && last.label === label) last.items.push(map(it));
        else out.push({ id: uid(), label, items: [map(it)] });
    }
    return out;
}

/** Time from the API (minutes) → form fields; when editing, < 60 min in minutes, otherwise hours. */
function timeFields(min: number, max: number | null) {
    const timeUnit = defaultTimeUnit(min);
    return {
        time: fromMinutes(min, timeUnit),
        timeMax: max ? fromMinutes(max, timeUnit) : '',
        timeUnit,
    };
}

export function fromRecipe(r: RecipeDto): IRecipeForm {
    const taxonomy: Record<string, number[]> = {};
    for (const [axis, values] of Object.entries(r.taxonomy)) {
        taxonomy[axis] = (values ?? []).map((v) => v.id);
    }
    const ingGroups = groupBy(
        r.ingredients,
        (i) => i.section,
        (i): IIngItem => ({
            id: uid(),
            ingredientId: i.ingredientId,
            name: i.ingredientName,
            qty: i.quantity ?? '',
            unitId: i.unitId,
            note: i.note ?? '',
            optional: i.optional,
            linked: i.linkedRecipe
                ? {
                      id: i.linkedRecipe.id,
                      title: i.linkedRecipe.title,
                      status: i.linkedRecipe.status,
                  }
                : null,
            substitutes: i.substitutes.map((s) => ({
                ingredientId: s.ingredientId,
                name: s.name,
                mix: s.mix,
            })),
        })
    );
    const stepGroups = groupBy(
        r.steps,
        (s) => s.section,
        (s): IStepItem => ({
            id: uid(),
            text: s.text,
            prep: s.prep,
            prepLead: s.prepLead ?? '',
            hints: [...(s.hints ?? [])],
        })
    );
    return {
        name: r.name,
        tagline: r.tagline ?? '',
        tip: r.tip ?? '',
        photoUrl: r.photoPath,
        taxonomy,
        difficulty: r.difficulty,
        ...timeFields(r.timeMin, r.timeMax),
        servings: r.servings ? String(r.servings) : '',
        ingGroups: ingGroups.length ? ingGroups : newForm().ingGroups,
        stepGroups: stepGroups.length ? stepGroups : newForm().stepGroups,
        equipment: r.equipment.map((e) => ({
            id: uid(),
            equipmentId: e.equipmentId,
            name: e.name,
            spec: e.spec ?? '',
            optional: e.optional,
        })),
        cautions: cautionRows(r.cautions),
    };
}

const toInt = (s: string) => {
    const n = parseInt(s, 10);
    return Number.isFinite(n) && n > 0 ? n : null;
};
const orNull = (s: string) => s.trim() || null;

export function toRequest(f: IRecipeForm, status: RecipeStatus): IUpsertRecipeRequest {
    let stepNo = 0;
    let sortOrder = 0;
    return {
        name: f.name.trim(),
        tagline: orNull(f.tagline),
        tip: orNull(f.tip),
        cautions: cautionTexts(f.cautions),
        servings: toInt(f.servings) ?? 0,
        timeMin: toMinutes(f.time, f.timeUnit),
        timeMax: f.timeMax.trim() ? toMinutes(f.timeMax, f.timeUnit) || null : null,
        difficulty: f.difficulty,
        photoUrl: f.photoUrl ?? '',
        status,
        taxonomyValueIds: Object.values(f.taxonomy).flat(),
        steps: f.stepGroups.flatMap((g) =>
            g.items
                .filter((s) => s.text.trim())
                .map((s) => ({
                    stepNo: ++stepNo,
                    section: orNull(g.label),
                    text: s.text.trim(),
                    prep: s.prep,
                    prepLead: s.prep ? orNull(s.prepLead) : null,
                    hints: s.hints.map((h) => h.trim()).filter(Boolean),
                }))
        ),
        ingredients: f.ingGroups.flatMap((g) =>
            g.items
                .filter((i) => i.ingredientId > 0)
                .map((i) => ({
                    sortOrder: ++sortOrder,
                    section: orNull(g.label),
                    ingredientId: i.ingredientId,
                    quantity: orNull(i.qty),
                    unitId: i.unitId,
                    note: orNull(i.note),
                    optional: i.optional,
                    linkedRecipeId: i.linked?.id ?? null,
                    // a substitute equal to the main ingredient (the row's ingredient was changed): skipped
                    substitutes: i.substitutes
                        .filter((s) => s.ingredientId !== i.ingredientId)
                        .map((s) => ({ ingredientId: s.ingredientId, mix: s.mix })),
                }))
        ),
        // an empty row / not picked from the dictionary: skipped
        equipment: f.equipment
            .filter((e) => e.equipmentId > 0)
            .map((e) => ({
                equipmentId: e.equipmentId,
                spec: orNull(e.spec),
                optional: e.optional,
            })),
    };
}

/** Moves the old recipe hint (tip) to the end of the step's hints and clears tip. Saved to the database only on "Zapisz". */
export function moveTipToStep(f: IRecipeForm, stepId: number): void {
    const tip = f.tip.trim();
    const step = f.stepGroups.flatMap((g) => g.items).find((s) => s.id === stepId);
    if (!tip || !step) return;
    step.hints.push(...splitHint(tip));
    f.tip = '';
}

/** Splits text longer than MAX_HINT_LEN at sentence ends (. ! ? …) into chunks ≤ MAX_HINT_LEN.
 *  A sentence longer than the limit stays whole; validation will flag it at the step. */
function splitHint(text: string): string[] {
    const len = (s: string) => [...s].length;
    if (len(text) <= MAX_HINT_LEN) return [text];
    const sentences = text.match(/[^.!?…]+[.!?…]+|[^.!?…]+$/g)?.map((s) => s.trim()) ?? [text];
    const out: string[] = [];
    for (const s of sentences.filter(Boolean)) {
        const last = out[out.length - 1];
        if (last && len(`${last} ${s}`) <= MAX_HINT_LEN) out[out.length - 1] = `${last} ${s}`;
        else out.push(s);
    }
    return out;
}

export type FormErrors = Partial<
    Record<
        | 'name'
        | 'ingredients'
        | 'steps'
        | 'difficulty'
        | 'time'
        | 'servings'
        | 'equipment'
        | `taxonomy.${string}`,
        string
    >
>;

export interface IValidation {
    /** Errors required for publishing. */
    errors: FormErrors;
    /** Ingredient and equipment rows with errors (row id → message). */
    rows: Record<number, string>;
}

export function validate(f: IRecipeForm, dicts: IDictionaries | undefined): IValidation {
    const errors: FormErrors = {};
    const rows: Record<number, string> = {};
    const unitById = new Map<number, IUnit>((dicts?.units ?? []).map((u) => [u.id, u]));

    if (!f.name.trim()) errors.name = 'Podaj nazwę przepisu.';

    for (const axis of dicts?.axes ?? []) {
        if (axis.required && !(f.taxonomy[axis.code] ?? []).length) {
            errors[`taxonomy.${axis.code}`] =
                `Wybierz co najmniej jedną wartość: ${axis.label.toLowerCase()}.`;
        }
    }
    if (!f.difficulty) errors.difficulty = 'Wybierz trudność.';
    if (!toInt(f.servings)) errors.servings = 'Podaj liczbę porcji.';
    const tMin = toMinutes(f.time, f.timeUnit);
    const tMax = f.timeMax.trim() ? toMinutes(f.timeMax, f.timeUnit) : null;
    if (!tMin) errors.time = 'Podaj czas (np. 40 albo 2,5).';
    else if (tMax !== null && tMax < tMin)
        errors.time = 'Górna granica musi być większa od dolnej.';

    const items = f.ingGroups.flatMap((g) => g.items);
    for (const it of items) {
        const filled = it.name.trim() || it.qty.trim();
        if (!filled) continue;
        if (!it.ingredientId) rows[it.id] = 'Wybierz składnik z listy albo dodaj go do bazy.';
        else if (!unitById.get(it.unitId ?? -1)?.noQty && !parseQty(it.qty).ok)
            rows[it.id] = 'Nieznany format. Wpisz np. 2, 1,5, ½, 1/4 albo 2–3.';
    }
    const rowErrCount = Object.keys(rows).length;
    if (!items.some((i) => i.ingredientId > 0))
        errors.ingredients = 'Dodaj co najmniej jeden składnik.';
    else if (rowErrCount) errors.ingredients = `Popraw składniki (${rowErrCount}).`;

    let equipErr = 0;
    for (const e of f.equipment) {
        if (e.name.trim() && !e.equipmentId) {
            rows[e.id] = 'Wybierz sprzęt z listy albo dodaj go do słownika.';
            equipErr++;
        }
    }
    if (equipErr) errors.equipment = `Popraw sprzęt (${equipErr}).`;

    const steps = f.stepGroups.flatMap((g) => g.items);
    if (!steps.some((s) => s.text.trim())) errors.steps = 'Dodaj co najmniej jeden krok.';
    // Number as in the editor (continuous row numbering).
    const longAt = steps.findIndex((s) => s.hints.some((h) => [...h.trim()].length > MAX_HINT_LEN));
    if (longAt >= 0)
        errors.steps = `Krok ${longAt + 1}: wskazówka może mieć najwyżej ${MAX_HINT_LEN} znaków.`;

    return { errors, rows };
}

/** Order of fields to scroll to on error (form element ids). */
export const ERROR_ORDER: [keyof FormErrors | string, string][] = [
    ['name', 'f-name'],
    ['ingredients', 'f-ingredients'],
    ['equipment', 'f-equipment'],
    ['steps', 'f-steps'],
    ['taxonomy.dish', 'f-taxonomy-dish'],
    ['difficulty', 'f-difficulty'],
    ['time', 'f-time'],
    ['servings', 'f-servings'],
];
