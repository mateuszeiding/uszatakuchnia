import { type FetchQueryOptions, useQueryClient } from '@tanstack/vue-query';

import { API, type EndpointsConfig } from '../API';

import type { RecipeBaseDto, RecipeDto } from '@/data/dtos/recipe/RecipeDto';

type RequestMap = {
    list: RecipeBaseDto[];
    [key: number]: RecipeDto;
};

export const fetchRecipes = async <E extends keyof RequestMap>(endpoint: E) => {
    const qc = useQueryClient();
    const e: NestedKeyUnion<EndpointsConfig, '/'> = `recipes/${endpoint}`;
    const options: FetchQueryOptions<RequestMap[E]> = {
        queryKey: [e],
        queryFn: () => API.Client.get<RequestMap[E]>(e),
    };
    return await qc.fetchQuery(options);
};
