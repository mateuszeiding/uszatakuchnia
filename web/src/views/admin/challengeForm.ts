// Formularz wyzwania (panel): model, walidacja przed wysłaniem, payload, błędy z serwera.
import type { IChallenge, IUpsertChallengeRequest } from '@/data/dtos/challenge/ChallengeDto';

export interface IChallengeIngRow {
    id: number;
    name: string;
    required: boolean;
}

export interface IChallengeForm {
    title: string;
    formatId: number | '';
    /** Minuty jako tekst z pola; puste = bez limitu. */
    timeLimit: string;
    difficulty: number | null;
    source: string;
    description: string;
    ingredients: IChallengeIngRow[];
}

let seq = 0;
export const blankIngredient = (): IChallengeIngRow => ({ id: ++seq, name: '', required: false });

export function formFromChallenge(c?: IChallenge): IChallengeForm {
    return {
        title: c?.title ?? '',
        formatId: c?.formatId ?? '',
        timeLimit: c?.timeLimitMin ? String(c.timeLimitMin) : '',
        difficulty: c?.difficulty ?? null,
        source: c?.source ?? '',
        description: c?.description ?? '',
        ingredients: (c?.ingredients ?? []).map((i) => ({ ...blankIngredient(), ...i })),
    };
}

const filledRows = (f: IChallengeForm) => f.ingredients.filter((r) => r.name.trim());

export function toPayload(f: IChallengeForm): IUpsertChallengeRequest {
    const t = f.timeLimit.trim();
    return {
        title: f.title.trim(),
        formatId: f.formatId === '' ? 0 : f.formatId,
        description: f.description.trim(),
        timeLimitMin: t ? Number(t) : null,
        difficulty: f.difficulty ?? 0,
        source: f.source.trim() || null,
        // puste wiersze pomijane
        ingredients: filledRows(f).map((r) => ({ name: r.name.trim(), required: r.required })),
    };
}

/** Porównanie „czy są niezapisane zmiany” — bez id wierszy i pustych wierszy. */
export const snapshot = (f: IChallengeForm) => JSON.stringify(toPayload(f));

/** Błędy pól; klucze: title, formatId, difficulty, description, timeLimitMin, source, ing-<id wiersza>. */
export function validateForm(f: IChallengeForm): Record<string, string> {
    const e: Record<string, string> = {};
    if (!f.title.trim()) e.title = 'Wpisz tytuł';
    else if (f.title.trim().length > 255) e.title = 'Tytuł może mieć najwyżej 255 znaków';
    if (f.formatId === '') e.formatId = 'Wybierz format';
    if (!f.difficulty) e.difficulty = 'Wybierz trudność';
    if (!f.description.trim()) e.description = 'Opisz zadanie';
    const t = f.timeLimit.trim();
    if (t && !(/^\d+$/.test(t) && Number(t) > 0))
        e.timeLimitMin = 'Wpisz liczbę minut większą od zera albo zostaw puste';
    if (f.source.trim().length > 255) e.source = 'Źródło może mieć najwyżej 255 znaków';
    for (const r of filledRows(f))
        if (r.name.trim().length > 120) e[`ing-${r.id}`] = 'Składnik może mieć najwyżej 120 znaków';
    return e;
}

/** Błędy z API (422) → klucze formularza; ingredients[i] liczy tylko niepuste wiersze. */
export function mapServerErrors(
    fields: { field: string; message: string }[],
    f: IChallengeForm
): Record<string, string> {
    const rows = filledRows(f);
    const e: Record<string, string> = {};
    for (const { field, message } of fields) {
        const m = /^ingredients\[(\d+)\]/.exec(field);
        const row = m ? rows[Number(m[1])] : undefined;
        e[row ? `ing-${row.id}` : field] = message;
    }
    return e;
}
