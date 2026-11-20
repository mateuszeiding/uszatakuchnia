// Uwagi w formularzach przepisu i wyzwania: wiersze edytora ⇄ lista tekstów z API.

/** Najdłuższa uwaga (jak MaxCautionLen w API). */
export const MAX_CAUTION_LEN = 220;
/** Od tylu znaków edytor pokazuje licznik. */
export const CAUTION_COUNTER_FROM = 160;

export interface ICautionRow {
    id: number;
    text: string;
}

let seq = 0;
export const blankCaution = (text = ''): ICautionRow => ({ id: ++seq, text });

export const cautionRows = (list: string[] | undefined): ICautionRow[] =>
    (list ?? []).map((t) => blankCaution(t));

/** Do zapisu: przycięte, bez pustych wierszy. */
export const cautionTexts = (rows: ICautionRow[]): string[] =>
    rows.map((r) => r.text.trim()).filter(Boolean);
