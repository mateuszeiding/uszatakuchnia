import { type QueryClient, useQueryClient } from '@tanstack/vue-query';

import { API, type EndpointsConfig } from '../API';

import type { IUpsertRecipeRequest, RecipeBaseDto, RecipeDto } from '@/data/dtos/recipe/RecipeDto';

type RequestMap = {
    list: RecipeBaseDto[];
    [key: number]: RecipeDto;
};

function doFetch<E extends keyof RequestMap>(qc: QueryClient, endpoint: E, admin = false) {
    const base: NestedKeyUnion<EndpointsConfig, '/'> = `recipes/${endpoint}`;
    const url = (admin && endpoint === 'list' ? `${base}?admin=true` : base) as NestedKeyUnion<
        EndpointsConfig,
        '/'
    >;
    return qc.fetchQuery({
        queryKey: [base, admin],
        queryFn: () => API.Client.get<RequestMap[E]>(url),
    });
}

export const fetchRecipes = <E extends keyof RequestMap>(endpoint: E, admin = false) => {
    const qc = useQueryClient();
    return doFetch(qc, endpoint, admin);
};

export const createRecipe = (body: IUpsertRecipeRequest) =>
    API.Client.post<{ id: number }>('recipes', body);

export const updateRecipe = (id: number, body: IUpsertRecipeRequest) =>
    API.Client.put<{ id: number }>(`recipes/${id}`, body);

export const deleteRecipe = (id: number) => API.Client.delete<{ deleted: number }>(`recipes/${id}`);
