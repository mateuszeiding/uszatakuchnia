/** Longest note (same as MaxCautionLen in the API). */
export const MAX_CAUTION_LEN = 220;
export const CAUTION_COUNTER_FROM = 160;

export interface ICautionRow {
    id: number;
    text: string;
}

let seq = 0;
export const blankCaution = (text = ''): ICautionRow => ({ id: ++seq, text });

export const cautionRows = (list: string[] | undefined): ICautionRow[] =>
    (list ?? []).map((t) => blankCaution(t));

export const cautionTexts = (rows: ICautionRow[]): string[] =>
    rows.map((r) => r.text.trim()).filter(Boolean);
