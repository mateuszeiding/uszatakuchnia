// Filtry listy przepisów: stan ↔ query w URL (do udostępniania linkiem), dopasowanie przepisów.
// Osie klasyfikacji filtrowane generycznie: param = kod osi, wartości = id po przecinku.
// Tryb osi z bazy (filterMode): 'and' = przepis ma wszystkie (diety), 'or' = dowolną.

import type { LocationQuery, LocationQueryRaw } from 'vue-router';

import type { RecipeBaseDto } from '@/data/dtos/recipe/RecipeDto';
import type { ITaxonomyAxis } from '@/data/dtos/taxonomy/DictionariesDto';

export type SortKey = 'newest' | 'fastest' | 'az';

export interface IListFilters {
    q: string;
    axes: Record<string, number[]>;
    maxTime: number;
    difficulty: number;
    noPrep: boolean;
    sort: SortKey;
}

export const TIMES: [number, string][] = [
    [15, 'do 15 min'],
    [30, 'do 30 min'],
    [60, 'do 1 h'],
    [120, 'do 2 h'],
];
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
    f.maxTime = Number(str(query.czas)) || 0;
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
    if (f.maxTime) q.czas = String(f.maxTime);
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
    /** Składniki pasujące do wyszukiwania (gdy trafienie nie jest w nazwie). */
    hits: string[] | null;
}

export function applyFilters(
    list: RecipeBaseDto[],
    f: IListFilters,
    axes: ITaxonomyAxis[]
): IResult[] {
    const q = f.q.trim().toLowerCase();
    let out: IResult[] = [];
    for (const r of list) {
        if (!q) {
            out.push({ recipe: r, hits: null });
            continue;
        }
        const cuisine = r.taxonomy.cuisine?.[0]?.label.toLowerCase() ?? '';
        if (r.name.toLowerCase().includes(q) || cuisine.includes(q)) {
            out.push({ recipe: r, hits: null });
            continue;
        }
        const hits = r.ingredientNames.filter((n) => n.toLowerCase().includes(q));
        if (hits.length) out.push({ recipe: r, hits });
    }

    out = out.filter(({ recipe: r }) => {
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
        // po czasie minimalnym (czas „od”); przepis bez czasu (0) nie przechodzi filtra
        if (f.maxTime && !(r.timeMin > 0 && r.timeMin <= f.maxTime)) return false;
        if (f.difficulty && r.difficulty !== f.difficulty) return false;
        if (f.noPrep && r.needsPrep) return false;
        return true;
    });

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
