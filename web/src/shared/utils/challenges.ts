import { formatTime } from './time';

/** Difficulty levels 1–5 (index = level − 1). */
export const LEVELS = ['rozgrzewka', 'łatwe', 'średnie', 'trudne', 'mistrzowskie'] as const;
export const DIFFICULTIES = [1, 2, 3, 4, 5] as const;

export const levelName = (n: number) => LEVELS[n - 1] ?? '';

export const timeLimitText = (min: number | null) => (min ? formatTime(min) : null) ?? 'bez limitu';

/** Selected difficulties in words: [] → "dowolna", [2,3] → "2–3", [1,4] → "1, 4". */
export function rangeText(sel: readonly number[]): string {
    const s = [...new Set(sel)].sort((a, b) => a - b);
    if (!s.length) return 'dowolna';
    const continuous = s.every((v, i) => i === 0 || v === s[i - 1]! + 1);
    return s.length > 1 && continuous ? `${s[0]}–${s[s.length - 1]}` : s.join(', ');
}

/** The ?los= parameter (comma-separated difficulties). null = no parameter, [] = any. */
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
    /** no challenges of the selected difficulty */
    | { kind: 'none' }
    /** everything in the pool is done: picking among the done ones is possible */
    | { kind: 'alldone' }
    /** the only undone one in the pool is the one already open ("Losuj jeszcze raz") */
    | { kind: 'onlycurrent' };

/**
 * Picks a challenge from the selected difficulties (all of them when none), skipping done ones.
 * includeDone: "Losuj spośród zrobionych", picks from the whole pool. exclude: the open challenge (don't pick it again).
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
