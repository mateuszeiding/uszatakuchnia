// ─── Unit definitions ─────────────────────────────────────────
// Each unit has a base value in ml (volume) or g (weight).
// "count" units scale numerically without conversion.

type UnitBase = 'ml' | 'g' | 'count';

// eslint-disable-next-line @typescript-eslint/naming-convention
interface UnitDef {
    label: string;
    base: UnitBase;
    toBase: number; // how many base units = 1 of this unit
    imperial?: boolean;
}

export const UNIT_DEFS: Record<string, UnitDef> = {
    // ── Volume (ml base) ─────────────────────────────────────
    kropla: { label: 'kropla', base: 'ml', toBase: 0.05 },
    łyżeczka: { label: 'łyżeczka', base: 'ml', toBase: 5 },
    łyżki: { label: 'łyżki', base: 'ml', toBase: 5 }, // plural alias
    łyżka: { label: 'łyżka', base: 'ml', toBase: 15 },
    szklanka: { label: 'szklanka', base: 'ml', toBase: 250 },
    ml: { label: 'ml', base: 'ml', toBase: 1 },
    l: { label: 'l', base: 'ml', toBase: 1000 },
    // ── Weight (g base) ──────────────────────────────────────
    szczypta: { label: 'szczypta', base: 'g', toBase: 0.5 },
    g: { label: 'g', base: 'g', toBase: 1 },
    dag: { label: 'dag', base: 'g', toBase: 10 },
    kg: { label: 'kg', base: 'g', toBase: 1000 },
    // ── Imperial volume ───────────────────────────────────────
    tsp: { label: 'tsp', base: 'ml', toBase: 4.93, imperial: true },
    tbsp: { label: 'tbsp', base: 'ml', toBase: 14.79, imperial: true },
    'fl oz': { label: 'fl oz', base: 'ml', toBase: 29.57, imperial: true },
    cup: { label: 'cup', base: 'ml', toBase: 236.6, imperial: true },
    pt: { label: 'pt', base: 'ml', toBase: 473.2, imperial: true },
    // ── Imperial weight ───────────────────────────────────────
    oz: { label: 'oz', base: 'g', toBase: 28.35, imperial: true },
    lb: { label: 'lb', base: 'g', toBase: 453.6, imperial: true },
    // ── Count (no conversion) ────────────────────────────────
    ząbek: { label: 'ząbek', base: 'count', toBase: 1 },
    sztuka: { label: 'sztuka', base: 'count', toBase: 1 },
    sztuki: { label: 'sztuki', base: 'count', toBase: 1 },
    filety: { label: 'filety', base: 'count', toBase: 1 },
    pęczek: { label: 'pęczek', base: 'count', toBase: 1 },
    garść: { label: 'garść', base: 'count', toBase: 1 },
};

// ─── Upgrade chains ───────────────────────────────────────────
// Defines the ordered list of units to consider when displaying a scaled value.
// Key = original unit → chain of candidates from smallest to largest.

const CHAINS: Record<string, string[]> = {
    // Volume: drop-based Polish
    kropla: ['kropla', 'łyżeczka', 'łyżka', 'szklanka', 'l'],
    łyżeczka: ['kropla', 'łyżeczka', 'łyżka', 'szklanka', 'l'],
    łyżki: ['kropla', 'łyżeczka', 'łyżka', 'szklanka', 'l'],
    łyżka: ['łyżeczka', 'łyżka', 'szklanka', 'l'],
    szklanka: ['łyżka', 'szklanka', 'l'],
    ml: ['ml', 'l'],
    l: ['ml', 'l'],
    // Dry: szczypta-based Polish
    szczypta: ['szczypta', 'łyżeczka', 'łyżka', 'szklanka'],
    // Weight
    g: ['g', 'dag', 'kg'],
    dag: ['g', 'dag', 'kg'],
    kg: ['g', 'dag', 'kg'],
    // Imperial volume
    tsp: ['tsp', 'tbsp', 'cup', 'pt'],
    tbsp: ['tsp', 'tbsp', 'cup', 'pt'],
    'fl oz': ['fl oz', 'cup', 'pt'],
    cup: ['tbsp', 'cup', 'pt'],
    pt: ['cup', 'pt'],
    // Imperial weight
    oz: ['oz', 'lb'],
    lb: ['oz', 'lb'],
};

// ─── Helpers ──────────────────────────────────────────────────

function fmtNum(n: number): string {
    if (n === 0) return '0';
    if (Math.abs(n - Math.round(n)) < 0.05) return String(Math.round(n));
    if (n < 10) return n.toFixed(1).replace(/\.0$/, '');
    return String(Math.round(n));
}

// ─── Public API ───────────────────────────────────────────────

/**
 * Scale an amount by a ratio and return the best display string.
 * e.g. scaleAmount(1, 'łyżeczka', 3) → '1 łyżka'
 *      scaleAmount(500, 'g', 3)       → '1.5 kg'
 */
export function scaleAmount(amount: number | null, unit: string | null, ratio: number): string {
    if (amount == null || amount === 0) return unit ? `— ${unit}` : '—';
    if (!unit) return fmtNum(amount * ratio);

    const def = UNIT_DEFS[unit];
    if (!def) return `${fmtNum(amount * ratio)} ${unit}`;

    const scaledAmount = amount * ratio;

    // Count units: just scale the number
    if (def.base === 'count') return `${fmtNum(scaledAmount)} ${unit}`;

    // Convert to base (ml or g)
    const baseValue = scaledAmount * def.toBase;

    // Find best unit in chain
    const chain = CHAINS[unit];
    if (!chain) return `${fmtNum(scaledAmount)} ${unit}`;

    let bestUnit = unit;
    let bestAmount = scaledAmount;

    for (const candidate of chain) {
        const candidateDef = UNIT_DEFS[candidate];
        if (!candidateDef || candidateDef.base !== def.base) continue;
        const val = baseValue / candidateDef.toBase;
        // Use this unit if value >= 1 (clean display)
        if (val >= 1) {
            bestUnit = candidate;
            bestAmount = val;
        }
    }

    return `${fmtNum(bestAmount)} ${bestUnit}`;
}

/**
 * Format amount + unit for display without scaling.
 */
export function formatAmount(amount: number | null, unit: string | null): string {
    return scaleAmount(amount, unit, 1);
}

// ─── Imperial conversion ──────────────────────────────────────

// Best imperial unit chains ordered from smallest to largest
const IMPERIAL_VOLUME_CHAIN = ['tsp', 'tbsp', 'fl oz', 'cup', 'pt'] as const;
const IMPERIAL_WEIGHT_CHAIN = ['oz', 'lb'] as const;

// Metric units that are "spoon-like" — map to tsp/tbsp first
const METRIC_VOLUME_UNITS = new Set([
    'kropla',
    'łyżeczka',
    'łyżki',
    'łyżka',
    'szklanka',
    'ml',
    'l',
]);
const METRIC_WEIGHT_UNITS = new Set(['szczypta', 'g', 'dag', 'kg']);

function bestImperialUnit(
    baseValue: number,
    chain: readonly string[]
): { amount: number; unit: string } {
    let bestUnit = chain[0]!;
    let bestAmount = baseValue / (UNIT_DEFS[bestUnit]?.toBase ?? 1);

    for (const candidate of chain) {
        const def = UNIT_DEFS[candidate];
        if (!def) continue;
        const val = baseValue / def.toBase;
        if (val >= 1) {
            bestUnit = candidate;
            bestAmount = val;
        }
    }
    return { amount: bestAmount, unit: bestUnit };
}

/**
 * Scale an amount and convert to imperial display.
 * e.g. scaleAmountImperial(3, 'łyżka', 1) → '3 tbsp'
 *      scaleAmountImperial(500, 'g', 1)    → '1.1 lb'
 *      scaleAmountImperial(2, 'ząbek', 1)  → '2 ząbek'  (count — no conversion)
 */
export function scaleAmountImperial(
    amount: number | null,
    unit: string | null,
    ratio: number
): string {
    if (amount == null || amount === 0) return unit ? `— ${unit}` : '—';
    if (!unit) return fmtNum(amount * ratio);

    const def = UNIT_DEFS[unit];
    if (!def) return `${fmtNum(amount * ratio)} ${unit}`;

    const scaledAmount = amount * ratio;

    // Count units: no conversion
    if (def.base === 'count') return `${fmtNum(scaledAmount)} ${unit}`;

    // Already imperial — just scale normally
    if (def.imperial) return scaleAmount(amount, unit, ratio);

    // Convert metric → base → best imperial
    const baseValue = scaledAmount * def.toBase;

    if (def.base === 'ml' && METRIC_VOLUME_UNITS.has(unit)) {
        const { amount: a, unit: u } = bestImperialUnit(baseValue, IMPERIAL_VOLUME_CHAIN);
        return `${fmtNum(a)} ${u}`;
    }

    if (def.base === 'g' && METRIC_WEIGHT_UNITS.has(unit)) {
        const { amount: a, unit: u } = bestImperialUnit(baseValue, IMPERIAL_WEIGHT_CHAIN);
        return `${fmtNum(a)} ${u}`;
    }

    // Fallback
    return `${fmtNum(scaledAmount)} ${unit}`;
}

// ─── Unit list for select ─────────────────────────────────────

export const UNIT_OPTIONS: { value: string; label: string; group: string }[] = [
    // Polish volume
    { value: '', label: '— jednostka —', group: '' },
    { value: 'kropla', label: 'kropla', group: 'objętość' },
    { value: 'łyżeczka', label: 'łyżeczka', group: 'objętość' },
    { value: 'łyżka', label: 'łyżka', group: 'objętość' },
    { value: 'szklanka', label: 'szklanka', group: 'objętość' },
    { value: 'ml', label: 'ml', group: 'objętość' },
    { value: 'l', label: 'l', group: 'objętość' },
    // Polish dry / weight
    { value: 'szczypta', label: 'szczypta', group: 'waga' },
    { value: 'g', label: 'g', group: 'waga' },
    { value: 'dag', label: 'dag', group: 'waga' },
    { value: 'kg', label: 'kg', group: 'waga' },
    // Count
    { value: 'ząbek', label: 'ząbek', group: 'sztuki' },
    { value: 'sztuka', label: 'sztuka', group: 'sztuki' },
    { value: 'filety', label: 'filety', group: 'sztuki' },
    { value: 'pęczek', label: 'pęczek', group: 'sztuki' },
    { value: 'garść', label: 'garść', group: 'sztuki' },
    // Imperial volume
    { value: 'tsp', label: 'tsp', group: 'imperial' },
    { value: 'tbsp', label: 'tbsp', group: 'imperial' },
    { value: 'fl oz', label: 'fl oz', group: 'imperial' },
    { value: 'cup', label: 'cup', group: 'imperial' },
    { value: 'pt', label: 'pt', group: 'imperial' },
    { value: 'oz', label: 'oz', group: 'imperial' },
    { value: 'lb', label: 'lb', group: 'imperial' },
];

export const UNIT_GROUPS = ['objętość', 'waga', 'sztuki', 'imperial'] as const;
