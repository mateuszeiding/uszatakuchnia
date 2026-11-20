// Ilości składników: parser (liczby, ułamki, zakresy), odmiana jednostek, skalowanie porcji.
// Logika przeniesiona z prototypu (Claude Design: cr2-ui.jsx, rv2-data.jsx).
// Ten sam parser jest po stronie API: services/quantity.go.

import type { IUnit } from '@/data/dtos/taxonomy/DictionariesDto';

export type Qty =
    | { ok: true; empty: true }
    | { ok: true; empty: false; min: number; max: number }
    | { ok: false };

const FR: Record<string, number> = { '½': 0.5, '¼': 0.25, '¾': 0.75, '⅓': 1 / 3, '⅔': 2 / 3 };

function num(raw: string): number | null {
    const s = raw.trim();
    if (!s) return null;
    if (FR[s] != null) return FR[s];
    let m: RegExpMatchArray | null;
    if ((m = s.match(/^(\d+)\s*([½¼¾⅓⅔])$/))) return +m[1]! + FR[m[2]!]!;
    if ((m = s.match(/^(\d+)\s+(\d+)\/(\d+)$/)) && +m[3]!) return +m[1]! + +m[2]! / +m[3]!;
    if ((m = s.match(/^(\d+)\/(\d+)$/)) && +m[2]!) return +m[1]! / +m[2]!;
    if (/^\d+([.,]\d+)?$/.test(s)) return parseFloat(s.replace(',', '.'));
    return null;
}

export function parseQty(raw: string | null | undefined): Qty {
    const s = (raw ?? '').trim();
    if (!s) return { ok: true, empty: true };
    const p = s.split(/\s*[-–—]\s*/);
    if (p.length === 2) {
        const a = num(p[0]!);
        const b = num(p[1]!);
        return a != null && b != null && b > a
            ? { ok: true, empty: false, min: a, max: b }
            : { ok: false };
    }
    const v = num(s);
    return v != null && v > 0 ? { ok: true, empty: false, min: v, max: v } : { ok: false };
}

export function pluralPL(n: number, one: string, few: string, many: string): string {
    if (n === 1) return one;
    const m10 = n % 10;
    const m100 = n % 100;
    return m10 >= 2 && m10 <= 4 && (m100 < 10 || m100 >= 20) ? few : many;
}

/** Odmiana jednostki dla ilości (1 ząbek / 3 ząbki / 5 ząbków; ułamki → forma 2–4). */
export function unitLabel(
    unit: IUnit | undefined,
    unitName: string | null,
    amount: number | null
): string {
    if (!unit) return unitName ?? '';
    if (!unit.forms || amount == null) return unit.name;
    const [one, few, many] = unit.forms;
    if (!Number.isInteger(amount)) return few;
    return pluralPL(amount, one, few, many);
}

// ─── Skalowanie ───────────────────────────────────────────────
export const COUNT = new Set(['sztuka', 'ząbek', 'listek', 'plaster', 'pęczek', 'filet', '']);
export const SPOON = new Set(['łyżka', 'łyżeczka', 'szklanka', 'szczypta', 'garść']);
/** Łyżki i łyżeczki od 1 w górę: do połówek (2⅔ → 2½), poniżej 1 jak inne miarki (¼, ⅓, ½, ¾). */
const HALVES = new Set(['łyżka', 'łyżeczka']);

/**
 * Jednostka sztukowa (zaokrąglana do całości albo zakresu, w sklepie w górę): znana sztuka (COUNT)
 * albo każda jednostka ze słownika z odmianą, która nie jest miarką (SPOON) — np. „łodyga” dodana
 * w panelu. g / ml / l nie mają odmiany, więc nigdy nie są sztukami.
 */
export function isCountUnit(name: string, unit?: IUnit): boolean {
    if (COUNT.has(name)) return true;
    return !!unit?.forms?.length && !SPOON.has(name);
}
const QUARTERS: Record<string, string> = { '0.25': '¼', '0.5': '½', '0.75': '¾' };
const THIRDS: Record<string, string> = { '0.33': '⅓', '0.67': '⅔' };

function quarters(v: number): string {
    const q = Math.max(0.25, Math.round(v * 4) / 4);
    const whole = Math.floor(q);
    const frac = QUARTERS[(q - whole).toFixed(2).replace(/0$/, '')];
    return frac ? `${whole || ''}${frac}` : String(whole);
}

/** Łyżki / szklanki / sztuki: tercje (⅓ ⅔), a poza nimi ćwiartki. */
function fraction(v: number): string {
    const whole = Math.floor(v);
    const third = THIRDS[(v - whole).toFixed(2)];
    return third ? `${whole || ''}${third}` : quarters(v);
}

/** Do połówek: 2⅔ → 2½, 1⅓ → 1½, 1¾ → 2. */
function halves(v: number): string {
    const h = Math.round(v * 2) / 2;
    const whole = Math.floor(h);
    return h > whole ? `${whole || ''}½` : String(whole);
}

/** Ilość do odmiany jednostki — po tym samym zaokrągleniu co liczba na ekranie. */
export function shownAmount(v: number, unit: string): number {
    if (HALVES.has(unit) && v >= 1) return Math.round(v * 2) / 2;
    if (SPOON.has(unit)) return Math.round(v * 4) / 4;
    return v;
}

/** Liczba do wyświetlenia: ułamki tylko dla łyżek/szklanek i sztuk < 1; g / ml / l zawsze dziesiętnie.
 *  `count` — jednostka sztukowa (isCountUnit); domyślnie znane sztuki z COUNT. */
export function fmtNum(v: number, unit: string, count = COUNT.has(unit)): string {
    if (HALVES.has(unit) && v >= 1) return halves(v);
    if (SPOON.has(unit) || (count && v < 1)) return fraction(v);
    if (v < 10 && !Number.isInteger(v))
        return (Math.round(v * 10) / 10).toString().replace('.', ',');
    if (v >= 50) return String(Math.round(v / 5) * 5);
    return String(Math.round(v));
}

/**
 * Skaluje ilość o `ratio` i zwraca tekst ilości + odmienioną jednostkę.
 * Jednostki sztukowe zaokrąglane do całości albo zakresu (1,5 cebuli → „1–2 sztuki”).
 * Zawsze formatuje — także przy bazowej liczbie porcji („1/2” → „½”, „1 1/2” → „1½”).
 */
export function scaleQuantity(
    raw: string | null,
    unitName: string | null,
    unit: IUnit | undefined,
    ratio: number
): { qty: string; unit: string } {
    const name = unit?.name ?? unitName ?? '';
    const q = parseQty(raw);
    if (!q.ok) return { qty: raw ?? '', unit: name };
    if (q.empty) return { qty: '', unit: name };
    let lo = q.min * ratio;
    let hi = q.max * ratio;
    const count = isCountUnit(name, unit);
    if (count && lo >= 1) {
        const near = (x: number) => Math.abs(x - Math.round(x)) < 0.05;
        lo = near(lo) ? Math.round(lo) : Math.floor(lo);
        hi = near(hi) ? Math.round(hi) : Math.ceil(hi);
        hi = Math.max(lo, hi);
    }
    const a = fmtNum(lo, name, count);
    const b = fmtNum(hi, name, count);
    const shownHi = shownAmount(hi, name);
    return { qty: a === b ? a : `${a}–${b}`, unit: unitLabel(unit, unitName, shownHi) };
}
