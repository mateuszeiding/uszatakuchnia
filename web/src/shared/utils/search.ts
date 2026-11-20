// Wyszukiwarka listy przepisów: bez polskich znaków i akcentów, z prostą odmianą słów (bez słownika).
// 1) obie strony znormalizowane: małe litery, bez diakrytyków (ł → l, ù → u) — „bulki” znajduje „Bułki”;
//    zapytanie Z polskimi znakami porównujemy dosłownie (tylko małe litery): „mąka” nie łapie „makaronu”;
// 2) całe zapytanie zawarte w tekście → trafienie (jak dotąd, tylko po normalizacji);
// 3) inaczej każde słowo zapytania musi pasować do jakiegoś słowa tekstu: zawiera się w nim albo
//    (słowo ≥ 4 znaki) ma z nim wspólny początek ≥ max(4, długość − 2, 70% długości) —
//    „bulka” → „bułki”, „tymianek” → „tymiankiem”, ale „cukier” ≠ „cukinia”.

/** Małe litery bez diakrytyków. */
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

/** Zapytanie przygotowane raz na całą listę. */
export interface ISearchQuery {
    text: string;
    words: string[];
    /** Zapytanie bez znaków diakrytycznych → porównujemy bez nich po obu stronach. */
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

/** Czy tekst (nazwa przepisu, kuchnia, składnik) pasuje do zapytania. */
export function matches(q: ISearchQuery, target: string): boolean {
    const t = q.fold ? normalize(target) : lower(target);
    if (t.includes(q.text)) return true;
    if (!q.words.length) return false;
    const tw = words(t);
    return q.words.every((x) => tw.some((w) => wordMatches(x, w)));
}
