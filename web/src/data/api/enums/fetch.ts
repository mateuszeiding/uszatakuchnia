import { type FetchQueryOptions, useQueryClient } from '@tanstack/vue-query';

import { API } from '../API';

import type { IngredientType } from '@/shared/enums/IngredientType';

type RequestMap = {
    ingredient: IngredientType;
};

export const fetchEnums = async <E extends keyof RequestMap>(endpoint: E) => {
    const qc = useQueryClient();
    const options: FetchQueryOptions<RequestMap[E][]> = {
        queryKey: [endpoint],
        queryFn: () => API.Client.get<RequestMap[E][]>(`enums/${endpoint}`),
    };
    return await qc.fetchQuery(options);
};
