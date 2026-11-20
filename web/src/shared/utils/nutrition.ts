// Macronutrients: API types, rounding and range formatting ("20 – 22 g (~21)").
// The calculations happen in the API (services/nutrition.go).

export type NutrientKey = 'kcal' | 'protein' | 'fat' | 'carbs' | 'fiber';

export const NUTRIENTS: { key: NutrientKey; label: string; unit: string }[] = [
    { key: 'kcal', label: 'Energia', unit: 'kcal' },
    { key: 'protein', label: 'Białko', unit: 'g' },
    { key: 'fat', label: 'Tłuszcz', unit: 'g' },
    { key: 'carbs', label: 'Węglowodany', unit: 'g' },
    { key: 'fiber', label: 'Błonnik', unit: 'g' },
];

/** Ingredient value (per 100 g) or grams per unit; max null = single value. */
export type Range = { min: number; max: number | null };
/** Value computed for a recipe. */
export type SumRange = { min: number; max: number };

export type Macros = Partial<Record<NutrientKey, Range>>;

export interface IRecipeNutrition {
    perPortion: Partial<Record<NutrientKey, SumRange>>;
    /** Skipped entirely (no macros or no conversion). */
    missing: string[];
    /** Counted, but missing a value → "bez: tofu naturalne". Absent in the old cache. */
    incomplete?: Partial<Record<NutrientKey, string[]>>;
}

/** '100g' = ingredient data (kcal step 1, g step 0.1) · 'recipe' = computed (kcal step 10 / 5 below 100, g step 1 / 0.5 below 10). */
export type RoundMode = '100g' | 'recipe';

export type Formatted =
    | { single: true; a: number }
    | { single: false; a: number; b: number; m: number };

export const nf = (v: number) => v.toLocaleString('pl-PL', { maximumFractionDigits: 1 });

export function rnd(v: number, key: NutrientKey | 'g', mode: RoundMode): number {
    if (mode === '100g') return key === 'kcal' ? Math.round(v) : Math.round(v * 10) / 10;
    if (key === 'kcal') return v >= 100 ? Math.round(v / 10) * 10 : Math.round(v / 5) * 5;
    return v >= 10 ? Math.round(v) : Math.round(v * 2) / 2;
}

export function fmt(
    val: Range | SumRange | null | undefined,
    key: NutrientKey | 'g',
    mode: RoundMode
): Formatted | null {
    if (!val || val.min == null) return null;
    const a = rnd(val.min, key, mode);
    if (val.max == null || val.max <= val.min) return { single: true, a };
    const b = rnd(val.max, key, mode);
    if (a === b) return { single: true, a };
    return { single: false, a, b, m: rnd((val.min + val.max) / 2, key, mode) };
}

/** Average (for a single value, the value itself), for "~520 kcal". */
export const avg = (f: Formatted) => (f.single ? f.a : f.m);

export const fmtText = (f: Formatted | null, unit: string) =>
    !f ? '—' : f.single ? `${nf(f.a)} ${unit}` : `${nf(f.a)} – ${nf(f.b)} ${unit} (~${nf(f.m)})`;

/** Screen reader text: "od 20 do 22 g, średnio około 21". */
export const spokenText = (f: Formatted | null, unit: string) =>
    !f
        ? 'brak danych'
        : f.single
          ? `${nf(f.a)} ${unit}`
          : `od ${nf(f.a)} do ${nf(f.b)} ${unit}, średnio około ${nf(f.m)}`;

/** Grams per unit: "90 – 100 g" or "10 g" (no average). */
export function fmtGrams(r: Range): string {
    const f = fmt(r, 'g', '100g');
    if (!f) return '—';
    return f.single ? `${nf(f.a)} g` : `${nf(f.a)} – ${nf(f.b)} g`;
}

export const scaleRange = (r: SumRange | undefined, n: number): SumRange | undefined =>
    r && { min: r.min * n, max: r.max * n };

/** Share of protein, fat and carbs in energy (P×4, F×9, C×4 from the averages), in %. */
export function energySplit(
    pp: IRecipeNutrition['perPortion']
): { protein: number; fat: number; carbs: number } | null {
    const mid = (r?: SumRange) => (r ? (r.min + r.max) / 2 : 0);
    const e = { protein: mid(pp.protein) * 4, fat: mid(pp.fat) * 9, carbs: mid(pp.carbs) * 4 };
    const sum = e.protein + e.fat + e.carbs;
    if (sum <= 0) return null;
    return {
        protein: Math.round((e.protein / sum) * 100),
        fat: Math.round((e.fat / sum) * 100),
        carbs: Math.round((e.carbs / sum) * 100),
    };
}
