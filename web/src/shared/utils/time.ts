// Recipe time: always stored in minutes (from–to), entered in min or h (with a decimal comma).

export type TimeUnit = 'min' | 'h';

/** 90 → "1 h 30 min", 40 → "40 min", 120 → "2 h". */
function formatOne(min: number): string {
    if (min < 60) return `${min} min`;
    const h = Math.floor(min / 60);
    const m = min % 60;
    return m ? `${h} h ${m} min` : `${h} h`;
}

/** Minutes as hours with a decimal comma: 150 → "2,5", 75 → "1,25". */
const hours = (min: number) => (min / 60).toLocaleString('pl-PL', { maximumFractionDigits: 2 });

/**
 * Recipe time for display:
 * - no max → "40 min", "1 h 30 min"
 * - both < 60 → "20–30 min"
 * - both ≥ 60 and divisible by 15 → "2,5–4,5 h"
 * - otherwise → "45 min – 1 h 30 min"
 * Returns null when no time is given (0).
 */
export function formatTime(min: number, max?: number | null): string | null {
    if (!min) return null;
    if (!max || max <= min) return formatOne(min);
    if (max < 60) return `${min}–${max} min`;
    if (min >= 60 && min % 15 === 0 && max % 15 === 0) return `${hours(min)}–${hours(max)} h`;
    return `${formatOne(min)} – ${formatOne(max)}`;
}

/** Field text (e.g. "2,5") in the given unit → minutes (rounded). 0 = empty / invalid. */
export function toMinutes(value: string, unit: TimeUnit): number {
    const n = parseFloat(value.trim().replace(',', '.'));
    if (!Number.isFinite(n) || n <= 0) return 0;
    return Math.round(unit === 'h' ? n * 60 : n);
}

/** Minutes → field text in the given unit: 150 + 'h' → "2,5", 40 + 'min' → "40". */
export function fromMinutes(min: number, unit: TimeUnit): string {
    if (!min) return '';
    const v = unit === 'h' ? Math.round((min / 60) * 100) / 100 : min;
    return String(v).replace('.', ',');
}

/** When editing: < 60 min → minutes, otherwise hours. */
export const defaultTimeUnit = (min: number): TimeUnit => (min >= 60 ? 'h' : 'min');
