import { type QueryClient, useQueryClient } from '@tanstack/vue-query';

import { API, type EndpointsConfig } from '../API';

import { type IUpsertRecipeRequest, RecipeBaseDto, RecipeDto } from '@/data/dtos/recipe/RecipeDto';

type RequestMap = {
    list: RecipeBaseDto[];
    [key: number]: RecipeDto;
};

export function fetchRecipesWithClient<E extends keyof RequestMap>(
    qc: QueryClient,
    endpoint: E,
    admin = false
): Promise<RequestMap[E]> {
    const base: NestedKeyUnion<EndpointsConfig, '/'> = `recipes/${endpoint}`;
    const url = (admin && endpoint === 'list' ? `${base}?admin=true` : base) as NestedKeyUnion<
        EndpointsConfig,
        '/'
    >;
    return qc.fetchQuery({
        queryKey: [base, admin],
        // Lista: 60 s świeża — powrót z przepisu bez ponownego pobierania (unieważnienie po zapisie / usunięciu).
        staleTime: endpoint === 'list' ? LIST_STALE_MS : 0,
        // API zwraca zwykłe obiekty — budujemy klasy DTO (photoUrl, domyślne wartości).
        queryFn: async () => {
            if (endpoint === 'list') {
                const raw = await API.Client.get<Partial<RecipeBaseDto>[]>(url);
                return raw.map((r) => new RecipeBaseDto(r)) as RequestMap[E];
            }
            const raw = await API.Client.get<Partial<RecipeDto>>(url);
            return new RecipeDto(raw) as RequestMap[E];
        },
    });
}

const LIST_STALE_MS = 60_000;

/** Lista z cache do pokazania od razu, bez szkieletów (odświeżenie w tle); undefined = brak danych
 *  albo cache unieważniony po zapisie / usunięciu — wtedy nie pokazujemy nieaktualnej listy. */
export function cachedRecipeList(qc: QueryClient, admin = false) {
    const key = ['recipes/list', admin];
    if (qc.getQueryState(key)?.isInvalidated) return undefined;
    return qc.getQueryData<RecipeBaseDto[]>(key);
}

// Use only from setup() context
export const fetchRecipes = <E extends keyof RequestMap>(endpoint: E, admin = false) =>
    fetchRecipesWithClient(useQueryClient(), endpoint, admin);

/** Po zapisie/usunięciu: wymuś świeże dane listy i szczegółów. */
export const invalidateRecipes = (qc: QueryClient) =>
    qc.invalidateQueries({ predicate: (q) => String(q.queryKey[0]).startsWith('recipes/') });

export const createRecipe = (body: IUpsertRecipeRequest) =>
    API.Client.post<{ id: number }>('recipes', body);

export const updateRecipe = (id: number, body: IUpsertRecipeRequest) =>
    API.Client.put<{ id: number }>(`recipes/${id}`, body);

export const deleteRecipe = (id: number) => API.Client.delete<{ deleted: number }>(`recipes/${id}`);
