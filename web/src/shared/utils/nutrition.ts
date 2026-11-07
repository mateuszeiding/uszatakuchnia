// Makroskładniki: typy z API, zaokrąglenia i formatowanie zakresów („20 – 22 g (~21)”).
// Logika przeniesiona z prototypu (.design/mk-main.jsx: rnd, fmt, fmtText). Wyliczenia są po stronie API (services/nutrition.go).

export type NutrientKey = 'kcal' | 'protein' | 'fat' | 'carbs' | 'fiber';

export const NUTRIENTS: { key: NutrientKey; label: string; unit: string }[] = [
    { key: 'kcal', label: 'Energia', unit: 'kcal' },
    { key: 'protein', label: 'Białko', unit: 'g' },
    { key: 'fat', label: 'Tłuszcz', unit: 'g' },
    { key: 'carbs', label: 'Węglowodany', unit: 'g' },
    { key: 'fiber', label: 'Błonnik', unit: 'g' },
];

/** Wartość składnika (na 100 g) lub gramy jednostki; max null = jedna wartość. */
export type Range = { min: number; max: number | null };
/** Wartość wyliczona dla przepisu. */
export type SumRange = { min: number; max: number };

export type Macros = Partial<Record<NutrientKey, Range>>;

export interface IRecipeNutrition {
    perPortion: Partial<Record<NutrientKey, SumRange>>;
    missing: string[];
}

/** '100g' = dane składnika (kcal co 1, g co 0,1) · 'recipe' = wyliczone (kcal co 10 / co 5 poniżej 100, g co 1 / co 0,5 poniżej 10). */
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

/** Średnia (dla pojedynczej wartości sama wartość) — do „~520 kcal”. */
export const avg = (f: Formatted) => (f.single ? f.a : f.m);

export const fmtText = (f: Formatted | null, unit: string) =>
    !f ? '—' : f.single ? `${nf(f.a)} ${unit}` : `${nf(f.a)} – ${nf(f.b)} ${unit} (~${nf(f.m)})`;

/** Tekst dla czytnika ekranu: „od 20 do 22 g, średnio około 21”. */
export const spokenText = (f: Formatted | null, unit: string) =>
    !f
        ? 'brak danych'
        : f.single
          ? `${nf(f.a)} ${unit}`
          : `od ${nf(f.a)} do ${nf(f.b)} ${unit}, średnio około ${nf(f.m)}`;

/** Gramy jednostki: „90 – 100 g” albo „10 g” (bez średniej). */
export function fmtGrams(r: Range): string {
    const f = fmt(r, 'g', '100g');
    if (!f) return '—';
    return f.single ? `${nf(f.a)} g` : `${nf(f.a)} – ${nf(f.b)} g`;
}

export const scaleRange = (r: SumRange | undefined, n: number): SumRange | undefined =>
    r && { min: r.min * n, max: r.max * n };

/** Udział białka, tłuszczu i węglowodanów w energii (B×4, T×9, W×4 ze średnich), w %. */
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
