import { type FetchQueryOptions, useQueryClient } from '@tanstack/vue-query';

import { API } from '../API';

import type { IngredientDto } from '@/data/dtos/ingredients/IngredientDto';

type RequestMap = {
    list: IngredientDto[];
};

export const fetchIngredients = async <E extends keyof RequestMap>(endpoint: E) => {
    const qc = useQueryClient();
    const options: FetchQueryOptions<RequestMap[E]> = {
        queryKey: [endpoint],
        queryFn: () => API.Client.get<RequestMap[E]>(`ingredients/${endpoint}`),
    };
    return await qc.fetchQuery(options);
};
