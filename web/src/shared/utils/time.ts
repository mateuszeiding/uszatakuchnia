// Czas przepisu: zapis zawsze w minutach (od–do), wpisywanie w min albo h (z przecinkiem).

export type TimeUnit = 'min' | 'h';

/** 90 → „1 h 30 min”, 40 → „40 min”, 120 → „2 h”. */
function formatOne(min: number): string {
    if (min < 60) return `${min} min`;
    const h = Math.floor(min / 60);
    const m = min % 60;
    return m ? `${h} h ${m} min` : `${h} h`;
}

/** Minuty jako godziny z przecinkiem: 150 → „2,5”, 75 → „1,25”. */
const hours = (min: number) => (min / 60).toLocaleString('pl-PL', { maximumFractionDigits: 2 });

/**
 * Czas przepisu do wyświetlenia:
 * - bez max → „40 min”, „1 h 30 min”
 * - oba < 60 → „20–30 min”
 * - oba ≥ 60 i podzielne przez 15 → „2,5–4,5 h”
 * - inaczej → „45 min – 1 h 30 min”
 * Zwraca null, gdy czasu nie podano (0).
 */
export function formatTime(min: number, max?: number | null): string | null {
    if (!min) return null;
    if (!max || max <= min) return formatOne(min);
    if (max < 60) return `${min}–${max} min`;
    if (min >= 60 && min % 15 === 0 && max % 15 === 0) return `${hours(min)}–${hours(max)} h`;
    return `${formatOne(min)} – ${formatOne(max)}`;
}

/** Tekst z pola (np. „2,5”) w danej jednostce → minuty (zaokrąglone). 0 = puste / niepoprawne. */
export function toMinutes(value: string, unit: TimeUnit): number {
    const n = parseFloat(value.trim().replace(',', '.'));
    if (!Number.isFinite(n) || n <= 0) return 0;
    return Math.round(unit === 'h' ? n * 60 : n);
}

/** Minuty → tekst pola w danej jednostce: 150 + 'h' → „2,5”, 40 + 'min' → „40”. */
export function fromMinutes(min: number, unit: TimeUnit): string {
    if (!min) return '';
    const v = unit === 'h' ? Math.round((min / 60) * 100) / 100 : min;
    return String(v).replace('.', ',');
}

/** Przy edycji: < 60 min → minuty, inaczej godziny. */
export const defaultTimeUnit = (min: number): TimeUnit => (min >= 60 ? 'h' : 'min');
