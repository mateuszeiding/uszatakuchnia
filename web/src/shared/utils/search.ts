// Recipe list search: ignores Polish characters and accents, with simple word inflection (no dictionary).
// 1) both sides normalized: lowercase, no diacritics (ł → l, ù → u), so "bulki" finds "Bułki";
//    a query WITH Polish characters is compared literally (lowercase only): "mąka" doesn't match "makaron";
// 2) the whole query contained in the text → match (after normalization);
// 3) otherwise every query word must match some word of the text: be contained in it, or
//    (word ≥ 4 chars) share a prefix ≥ max(4, length − 2, 70% of length) with it:
//    "bulka" → "bułki", "tymianek" → "tymiankiem", but "cukier" ≠ "cukinia".

/** Lowercase without diacritics. */
export function normalize(s: string): string {
    return s.toLowerCase().replace(/ł/g, 'l').normalize('NFD').replace(/\p{M}/gu, '');
}

const words = (s: string) => s.split(/[^\p{L}\p{N}]+/u).filter(Boolean);

function commonPrefix(a: string, b: string): number {
    let i = 0;
    while (i < a.length && i < b.length && a[i] === b[i]) i++;
    return i;
}

function wordMatches(q: string, w: string): boolean {
    if (w.includes(q)) return true;
    if (q.length < 4) return false;
    return commonPrefix(q, w) >= Math.max(4, q.length - 2, Math.ceil(q.length * 0.7));
}

/** Query prepared once for the whole list. */
export interface ISearchQuery {
    text: string;
    words: string[];
    /** Query without diacritics → compare without them on both sides. */
    fold: boolean;
}

const lower = (s: string) => s.toLowerCase().normalize('NFC');

export function prepareQuery(raw: string): ISearchQuery | null {
    const plain = lower(raw).trim();
    const folded = normalize(plain);
    const fold = folded === plain;
    const text = fold ? folded : plain;
    return text ? { text, words: words(text), fold } : null;
}

/** Whether the text (recipe name, cuisine, ingredient) matches the query. */
export function matches(q: ISearchQuery, target: string): boolean {
    const t = q.fold ? normalize(target) : lower(target);
    if (t.includes(q.text)) return true;
    if (!q.words.length) return false;
    const tw = words(t);
    return q.words.every((x) => tw.some((w) => wordMatches(x, w)));
}
