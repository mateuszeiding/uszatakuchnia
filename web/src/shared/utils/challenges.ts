// Wyzwania: poziomy trudności, zakres wybranych trudności, losowanie (pomija zrobione), akapity opisu.
import { formatTime } from './time';

/** Poziomy trudności 1–5 (indeks = poziom − 1). */
export const LEVELS = ['rozgrzewka', 'łatwe', 'średnie', 'trudne', 'mistrzowskie'] as const;
export const DIFFICULTIES = [1, 2, 3, 4, 5] as const;

export const levelName = (n: number) => LEVELS[n - 1] ?? '';

/** Limit czasu: „45 min”, „1 h 15 min”, „bez limitu”. */
export const timeLimitText = (min: number | null) => (min ? formatTime(min) : null) ?? 'bez limitu';

/** Wybrane trudności słownie: [] → „dowolna”, [2,3] → „2–3”, [1,4] → „1, 4”. */
export function rangeText(sel: readonly number[]): string {
    const s = [...new Set(sel)].sort((a, b) => a - b);
    if (!s.length) return 'dowolna';
    const continuous = s.every((v, i) => i === 0 || v === s[i - 1]! + 1);
    return s.length > 1 && continuous ? `${s[0]}–${s[s.length - 1]}` : s.join(', ');
}

/** Parametr ?los= (trudności po przecinku). null = brak parametru, [] = dowolna. */
export function parseLos(raw: unknown): number[] | null {
    const v = Array.isArray(raw) ? raw[0] : raw;
    if (typeof v !== 'string') return null;
    const out = v
        .split(',')
        .map((x) => Number(x.trim()))
        .filter((n) => Number.isInteger(n) && n >= 1 && n <= 5);
    return [...new Set(out)].sort((a, b) => a - b);
}

export const losParam = (sel: readonly number[]) =>
    [...new Set(sel)].sort((a, b) => a - b).join(',');

interface IRollable {
    id: number;
    difficulty: number;
}

export type RollResult =
    | { kind: 'ok'; id: number }
    /** nie ma wyzwań o wybranej trudności */
    | { kind: 'none' }
    /** wszystkie z puli zrobione — można losować spośród zrobionych */
    | { kind: 'alldone' }
    /** jedynym niezrobionym w puli jest to, które jest otwarte („Losuj jeszcze raz”) */
    | { kind: 'onlycurrent' };

/**
 * Losuje wyzwanie z wybranych trudności (wszystkie, gdy nic), z pominięciem zrobionych.
 * includeDone — „Losuj spośród zrobionych”: z całej puli. exclude — otwarte wyzwanie (nie losuj go znowu).
 */
export function rollChallenge(
    items: readonly IRollable[],
    sel: readonly number[],
    done: ReadonlySet<number>,
    opt: { includeDone?: boolean; exclude?: number; random?: () => number } = {}
): RollResult {
    const pool = items.filter((c) => !sel.length || sel.includes(c.difficulty));
    if (!pool.length) return { kind: 'none' };
    const candidates = pool.filter(
        (c) => c.id !== opt.exclude && (opt.includeDone || !done.has(c.id))
    );
    if (!candidates.length) {
        const currentFresh =
            opt.exclude != null &&
            !opt.includeDone &&
            pool.some((c) => c.id === opt.exclude && !done.has(c.id));
        return { kind: currentFresh ? 'onlycurrent' : 'alldone' };
    }
    const r = (opt.random ?? Math.random)();
    return { kind: 'ok', id: candidates[Math.floor(r * candidates.length)]!.id };
}

/** Opis → akapity (pusta linia rozdziela akapity). */
export const paragraphs = (text: string) =>
    text
        .split(/\n\s*\n/)
        .map((p) => p.trim())
        .filter(Boolean);

/** 1 wyzwanie, 2 wyzwania, 5 wyzwań. */
export function plural(n: number, one: string, few: string, many: string) {
    if (n === 1) return one;
    const d = n % 10;
    const t = n % 100;
    return d >= 2 && d <= 4 && (t < 10 || t >= 20) ? few : many;
}
