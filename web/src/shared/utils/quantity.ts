// The same parser lives in the API: services/quantity.go.

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

/** Unit inflection for a quantity (1 ząbek / 3 ząbki / 5 ząbków; fractions → the 2–4 form). */
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

export const COUNT = new Set(['sztuka', 'ząbek', 'listek', 'plaster', 'pęczek', 'filet', '']);
export const SPOON = new Set(['łyżka', 'łyżeczka', 'szklanka', 'szczypta', 'garść']);
/** Tablespoons and teaspoons from 1 up: to halves (2 2/3 → 2,5); below 1 like other measures (0,25, 1/3, 0,5, 0,75). */
const HALVES = new Set(['łyżka', 'łyżeczka']);

/**
 * Count unit (rounded to a whole number or a range, up in the store): a known count unit (COUNT)
 * or any inflected dictionary unit that isn't a measure (SPOON), e.g. "łodyga" added
 * in the admin panel. g / ml / l don't inflect, so they're never count units.
 */
export function isCountUnit(name: string, unit?: IUnit): boolean {
    if (COUNT.has(name)) return true;
    return !!unit?.forms?.length && !SPOON.has(name);
}
// Display never uses the glyphs ½ ¼ ¾ ⅓ ⅔ (hard to read). Halves and quarters are decimal with a Polish comma
// ('0,25', '0,5', '1,75'); only thirds stay fractions, with the whole part apart ('1/3', '1 1/3') — 0,33 can't be measured.
const DECIMAL: Record<string, string> = { '0.25': '25', '0.50': '5', '0.75': '75' };
const THIRDS: Record<string, string> = { '0.33': '1/3', '0.67': '2/3' };

/** Whole part and a third → '1/3' / '1 1/3'. */
function mixed(whole: number, frac: string): string {
    return whole ? `${whole} ${frac}` : frac;
}

/** Whole part and decimal digits of a quarter or half → '0,25' / '1,5' / '2,75'. */
const decimal = (whole: number, digits: string) => `${whole},${digits}`;

function quarters(v: number): string {
    const q = Math.max(0.25, Math.round(v * 4) / 4);
    const whole = Math.floor(q);
    const rest = DECIMAL[(q - whole).toFixed(2)];
    return rest ? decimal(whole, rest) : String(whole);
}

/** Spoons / cups / pieces: thirds (1/3, 2/3), otherwise quarters. */
function fraction(v: number): string {
    const whole = Math.floor(v);
    const third = THIRDS[(v - whole).toFixed(2)];
    return third ? mixed(whole, third) : quarters(v);
}

/** To halves: 2 2/3 → 2,5, 1 1/3 → 1,5, 1,75 → 2. */
function halves(v: number): string {
    const h = Math.round(v * 2) / 2;
    const whole = Math.floor(h);
    return h > whole ? decimal(whole, '5') : String(whole);
}

/** Quantity for unit inflection, after the same rounding as the displayed number. */
export function shownAmount(v: number, unit: string): number {
    if (HALVES.has(unit) && v >= 1) return Math.round(v * 2) / 2;
    if (SPOON.has(unit)) return Math.max(0.25, Math.round(v * 4) / 4);
    return v;
}

/** Number for display: fractions only for spoons/cups and pieces < 1; g / ml / l always decimal.
 *  `count`: a count unit (isCountUnit); defaults to the known COUNT units. */
export function fmtNum(v: number, unit: string, count = COUNT.has(unit)): string {
    if (HALVES.has(unit) && v >= 1) return halves(v);
    if (SPOON.has(unit) || (count && v < 1)) return fraction(v);
    if (v < 10 && !Number.isInteger(v))
        return (Math.round(v * 10) / 10).toString().replace('.', ',');
    if (v >= 50) return String(Math.round(v / 5) * 5);
    return String(Math.round(v));
}

/**
 * Scales a quantity by `ratio` and returns the quantity text + the inflected unit.
 * Count units are rounded to a whole number or a range (1.5 onions → "1–2 sztuki").
 * Always formats, also at the base servings ("½" → "0,5", "1¼" → "1,25").
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
