// Classification axes are filtered generically: param = axis code, values = comma-separated ids.
// Axis mode from the database (filterMode): 'and' = the recipe has all of them (diets), 'or' = any of them.

import type { LocationQuery, LocationQueryRaw } from 'vue-router';

import type { RecipeBaseDto } from '@/data/dtos/recipe/RecipeDto';
import type { ITaxonomyAxis } from '@/data/dtos/taxonomy/DictionariesDto';
import { matches, prepareQuery } from '@/shared/utils/search';

export type SortKey = 'newest' | 'fastest' | 'az';

export interface IListFilters {
    q: string;
    axes: Record<string, number[]>;
    maxTime: number;
    difficulty: number;
    noPrep: boolean;
    sort: SortKey;
}

/** Time: disjoint ranges by time_min. Key = upper bound, 121 = "ponad 2 h" ("121+" in the URL). */
export const TIMES: [number, string][] = [
    [30, 'do 30 min'],
    [60, '30–60 min'],
    [120, '1–2 h'],
    [121, 'ponad 2 h'],
];
export const OVER_2H = 121;

/** The recipe's time range (0 = no time given). Old links (czas=15/30/60/120) → the nearest range. */
export function timeBucket(min: number): number {
    if (!min || min < 0) return 0;
    if (min <= 30) return 30;
    if (min <= 60) return 60;
    if (min <= 120) return 120;
    return OVER_2H;
}
export const DIFFICULTY = ['', 'łatwe', 'średnie', 'trudne'];
export const PAGE_SIZE = 12;

const SORTS: SortKey[] = ['newest', 'fastest', 'az'];

export const emptyFilters = (sort: SortKey = 'newest'): IListFilters => ({
    q: '',
    axes: {},
    maxTime: 0,
    difficulty: 0,
    noPrep: false,
    sort,
});

const str = (v: LocationQuery[string] | undefined) => (Array.isArray(v) ? v[0] : v) ?? '';

export function readQuery(query: LocationQuery, axes: ITaxonomyAxis[]): IListFilters {
    const f = emptyFilters();
    f.q = str(query.q);
    for (const a of axes) {
        const ids = str(query[a.code])
            .split(',')
            .map(Number)
            .filter((n) => Number.isInteger(n) && n > 0);
        if (ids.length) f.axes[a.code] = ids;
    }
    const czas = str(query.czas);
    f.maxTime = czas === '121+' ? OVER_2H : timeBucket(Number(czas) || 0);
    f.difficulty = Number(str(query.trud)) || 0;
    f.noPrep = str(query.bezprzyg) === '1';
    const sort = str(query.sort) as SortKey;
    f.sort = SORTS.includes(sort) ? sort : 'newest';
    return f;
}

export function toQuery(f: IListFilters): LocationQueryRaw {
    const q: LocationQueryRaw = {};
    if (f.q.trim()) q.q = f.q;
    for (const [code, ids] of Object.entries(f.axes)) if (ids.length) q[code] = ids.join(',');
    if (f.maxTime) q.czas = f.maxTime === OVER_2H ? '121+' : String(f.maxTime);
    if (f.difficulty) q.trud = String(f.difficulty);
    if (f.noPrep) q.bezprzyg = '1';
    if (f.sort !== 'newest') q.sort = f.sort;
    return q;
}

export const activeCount = (f: IListFilters) =>
    Object.values(f.axes).reduce((n, ids) => n + ids.length, 0) +
    (f.maxTime ? 1 : 0) +
    (f.difficulty ? 1 : 0) +
    (f.noPrep ? 1 : 0);

export interface IResult {
    recipe: RecipeBaseDto;
    /** Ingredients matching the search (when the match isn't in the name). */
    hits: string[] | null;
}

export function applyFilters(
    list: RecipeBaseDto[],
    f: IListFilters,
    axes: ITaxonomyAxis[]
): IResult[] {
    // ignores Polish characters, simple word inflection, also matches substitutes (shared/utils/search.ts)
    const q = prepareQuery(f.q);
    let out: IResult[] = [];
    for (const r of list) {
        if (!q) {
            out.push({ recipe: r, hits: null });
            continue;
        }
        const cuisine = r.taxonomy.cuisine?.[0]?.label ?? '';
        if (matches(q, r.name) || (cuisine && matches(q, cuisine))) {
            out.push({ recipe: r, hits: null });
            continue;
        }
        const hits = [
            ...new Set(
                [...r.ingredientNames, ...(r.substituteNames ?? [])].filter((n) => matches(q, n))
            ),
        ];
        if (hits.length) out.push({ recipe: r, hits });
    }

    out = out.filter(({ recipe: r }) => passes(r, f, axes));

    const cmp: Record<SortKey, (a: IResult, b: IResult) => number> = {
        newest: (a, b) => b.recipe.id - a.recipe.id,
        fastest: (a, b) => (a.recipe.timeMin || 9999) - (b.recipe.timeMin || 9999),
        az: (a, b) => a.recipe.name.localeCompare(b.recipe.name, 'pl'),
    };
    return out.sort(cmp[f.sort]);
}

export const recipesLabel = (n: number) =>
    n === 1
        ? 'przepis'
        : n % 10 >= 2 && n % 10 <= 4 && (n % 100 < 10 || n % 100 >= 20)
          ? 'przepisy'
          : 'przepisów';

/** Filters without search: axes (diets AND, the rest OR), time range by time_min, difficulty, no prep. */
function passes(r: RecipeBaseDto, f: IListFilters, axes: ITaxonomyAxis[]): boolean {
    for (const a of axes) {
        const wanted = f.axes[a.code] ?? [];
        if (!wanted.length) continue;
        const has = new Set((r.taxonomy[a.code] ?? []).map((v) => v.id));
        const ok =
            a.filterMode === 'and'
                ? wanted.every((id) => has.has(id))
                : wanted.some((id) => has.has(id));
        if (!ok) return false;
    }
    // a recipe without a time (0) doesn't pass the time filter
    if (f.maxTime && timeBucket(r.timeMin) !== f.maxTime) return false;
    if (f.difficulty && r.difficulty !== f.difficulty) return false;
    if (f.noPrep && r.needsPrep) return false;
    return true;
}

// The panel shows only values that at least one published recipe has (on the whole list, so nothing jumps),
// with a counter: how many recipes there would be after adding the value to the current choice (with the search).
// A selected value: how many of the current results have it. 0 and not selected → dimmed, disabled, in its place.

export type FacetDim = { kind: 'axis'; code: string } | { kind: 'time' } | { kind: 'difficulty' };

function hasValue(r: RecipeBaseDto, dim: FacetDim, v: number): boolean {
    if (dim.kind === 'axis') return (r.taxonomy[dim.code] ?? []).some((x) => x.id === v);
    if (dim.kind === 'time') return timeBucket(r.timeMin) === v;
    return r.difficulty === v;
}

function isSelected(f: IListFilters, dim: FacetDim, v: number): boolean {
    if (dim.kind === 'axis') return (f.axes[dim.code] ?? []).includes(v);
    if (dim.kind === 'time') return f.maxTime === v;
    return f.difficulty === v;
}

function withValue(f: IListFilters, dim: FacetDim, v: number): IListFilters {
    if (dim.kind === 'axis')
        return { ...f, axes: { ...f.axes, [dim.code]: [...(f.axes[dim.code] ?? []), v] } };
    if (dim.kind === 'time') return { ...f, maxTime: v };
    return { ...f, difficulty: v };
}

export const facetPresent = (list: RecipeBaseDto[], dim: FacetDim, v: number) =>
    list.some((r) => hasValue(r, dim, v));

export function facetCount(
    list: RecipeBaseDto[],
    f: IListFilters,
    axes: ITaxonomyAxis[],
    dim: FacetDim,
    v: number
): number {
    if (isSelected(f, dim, v))
        return applyFilters(list, f, axes).filter(({ recipe }) => hasValue(recipe, dim, v)).length;
    return applyFilters(list, withValue(f, dim, v), axes).length;
}
