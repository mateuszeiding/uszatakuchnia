import { type FetchQueryOptions, useQueryClient } from '@tanstack/vue-query';

import { API } from '../API';

import type { IType } from '@/data/dtos/ingredients/IType';
import type { IngredientType } from '@/shared/enums/IngredientType';

type RequestMap = {
    ingredient: IType<IngredientType>;
};

export const fetchEnums = async <E extends keyof RequestMap>(endpoint: E) => {
    const qc = useQueryClient();
    const options: FetchQueryOptions<RequestMap[E][]> = {
        queryKey: [endpoint],
        queryFn: () => API.Client.get<RequestMap[E][]>(`enums/${endpoint}`),
    };
    return await qc.fetchQuery(options);
};
