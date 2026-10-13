// Stan formularza przepisu: typy, konwersja z/do API, walidacja.
// Szkic wymaga tylko nazwy; publikacja — kompletu (nazwa, składnik, krok, wymagane osie, trudność, porcje).

import type { IUpsertRecipeRequest, RecipeDto, RecipeStatus } from '@/data/dtos/recipe/RecipeDto';
import type { IDictionaries, IUnit } from '@/data/dtos/taxonomy/DictionariesDto';
import { parseQty } from '@/shared/utils/quantity';

let seq = 0;
export const uid = () => ++seq;

export const PREP_LEADS = ['30 min', '1 h', '2 h', '4 h', '12 h', 'dzień wcześniej'];

export interface IIngItem {
    id: number;
    ingredientId: number;
    /** Tekst w polu podpowiadania; ingredientId > 0 dopiero po wyborze z bazy. */
    name: string;
    qty: string;
    unitId: number | null;
    note: string;
    optional: boolean;
}

export interface IStepItem {
    id: number;
    text: string;
    prep: boolean;
    prepLead: string;
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
    /** Wybrane wartości per kod osi (dish, base, cuisine, diet, practical, …). */
    taxonomy: Record<string, number[]>;
    difficulty: number | null;
    time: string;
    servings: string;
    kcal: string;
    ingGroups: IGroup<IIngItem>[];
    stepGroups: IGroup<IStepItem>[];
}

export const blankIng = (): IIngItem => ({
    id: uid(),
    ingredientId: 0,
    name: '',
    qty: '',
    unitId: null,
    note: '',
    optional: false,
});

export const blankStep = (): IStepItem => ({ id: uid(), text: '', prep: false, prepLead: '' });

export const newForm = (): IRecipeForm => ({
    name: '',
    tagline: '',
    tip: '',
    photoUrl: null,
    taxonomy: {},
    difficulty: null,
    time: '',
    servings: '4',
    kcal: '',
    ingGroups: [{ id: uid(), label: '', items: [blankIng()] }],
    stepGroups: [{ id: uid(), label: '', items: [blankStep()] }],
});

/** Grupuje kolejne elementy o tej samej sekcji (kolejność z API zachowana). */
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
        })
    );
    const stepGroups = groupBy(
        r.steps,
        (s) => s.section,
        (s): IStepItem => ({ id: uid(), text: s.text, prep: s.prep, prepLead: s.prepLead ?? '' })
    );
    return {
        name: r.name,
        tagline: r.tagline ?? '',
        tip: r.tip ?? '',
        photoUrl: r.photoPath,
        taxonomy,
        difficulty: r.difficulty,
        time: r.timeMinutes ? String(r.timeMinutes) : '',
        servings: r.servings ? String(r.servings) : '',
        kcal: r.kcalPerServing ? String(r.kcalPerServing) : '',
        ingGroups: ingGroups.length ? ingGroups : newForm().ingGroups,
        stepGroups: stepGroups.length ? stepGroups : newForm().stepGroups,
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
        servings: toInt(f.servings) ?? 0,
        timeMinutes: toInt(f.time),
        difficulty: f.difficulty,
        kcalPerServing: toInt(f.kcal),
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
                }))
        ),
    };
}

export type FormErrors = Partial<
    Record<
        'name' | 'ingredients' | 'steps' | 'difficulty' | 'servings' | `taxonomy.${string}`,
        string
    >
>;

export interface IValidation {
    /** Błędy wymagane do publikacji. */
    errors: FormErrors;
    /** Wiersze składników z błędem (id wiersza → komunikat). */
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

    if (!f.stepGroups.some((g) => g.items.some((s) => s.text.trim())))
        errors.steps = 'Dodaj co najmniej jeden krok.';

    return { errors, rows };
}

/** Kolejność pól do przewinięcia przy błędzie (id elementów w formularzu). */
export const ERROR_ORDER: [keyof FormErrors | string, string][] = [
    ['name', 'f-name'],
    ['ingredients', 'f-ingredients'],
    ['steps', 'f-steps'],
    ['taxonomy.dish', 'f-taxonomy-dish'],
    ['difficulty', 'f-difficulty'],
    ['servings', 'f-servings'],
];
