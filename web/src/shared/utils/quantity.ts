// Ilości składników: parser (liczby, ułamki, zakresy), odmiana jednostek, skalowanie porcji.
// Logika przeniesiona z prototypu (design/UszataKuchnia/cr2-ui.jsx, rv2-data.jsx).
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
const COUNT = new Set(['sztuka', 'ząbek', 'listek', 'plaster', 'pęczek', 'filet', '']);
const SPOON = new Set(['łyżka', 'łyżeczka', 'szklanka', 'szczypta', 'garść']);
const QUARTERS: Record<string, string> = { '0.25': '¼', '0.5': '½', '0.75': '¾' };

function quarters(v: number): string {
    const q = Math.max(0.25, Math.round(v * 4) / 4);
    const whole = Math.floor(q);
    const frac = QUARTERS[(q - whole).toFixed(2).replace(/0$/, '')];
    return frac ? `${whole || ''}${frac}` : String(whole);
}

function fmtNum(v: number, unit: string): string {
    if (SPOON.has(unit) || (COUNT.has(unit) && v < 1)) return quarters(v);
    if (v < 10 && !Number.isInteger(v))
        return (Math.round(v * 10) / 10).toString().replace('.', ',');
    if (v >= 50) return String(Math.round(v / 5) * 5);
    return String(Math.round(v));
}

/**
 * Skaluje ilość o `ratio` i zwraca tekst ilości + odmienioną jednostkę.
 * Jednostki sztukowe zaokrąglane do całości albo zakresu (1,5 cebuli → „1–2 sztuki”).
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
    if (ratio === 1) return { qty: raw!.trim(), unit: unitLabel(unit, unitName, q.max) };

    let lo = q.min * ratio;
    let hi = q.max * ratio;
    if (COUNT.has(name) && lo >= 1) {
        const near = (x: number) => Math.abs(x - Math.round(x)) < 0.05;
        lo = near(lo) ? Math.round(lo) : Math.floor(lo);
        hi = near(hi) ? Math.round(hi) : Math.ceil(hi);
        hi = Math.max(lo, hi);
    }
    const a = fmtNum(lo, name);
    const b = fmtNum(hi, name);
    const shownHi = SPOON.has(name) ? Math.round(hi * 4) / 4 : hi;
    return { qty: a === b ? a : `${a}–${b}`, unit: unitLabel(unit, unitName, shownHi) };
}
