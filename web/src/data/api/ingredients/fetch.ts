import { type QueryClient, useQueryClient } from '@tanstack/vue-query';

import { API } from '../API';

import type { IngredientDto } from '@/data/dtos/ingredients/IngredientDto';

const LIST_KEY = ['ingredients/list'];

export const fetchIngredientsWithClient = (qc: QueryClient) =>
    qc.fetchQuery({
        queryKey: LIST_KEY,
        queryFn: () => API.Client.get<IngredientDto[]>('ingredients/list'),
    });

// Use only from setup() context
export const fetchIngredients = (_endpoint: 'list' = 'list') =>
    fetchIngredientsWithClient(useQueryClient());

export const invalidateIngredients = (qc: QueryClient) =>
    qc.invalidateQueries({ queryKey: LIST_KEY });

export interface IIngredientPayload {
    name: string;
    typeCode: string;
    allergens: string[];
}

export const createIngredient = (payload: IIngredientPayload): Promise<IngredientDto> =>
    API.Client.post<IngredientDto>('ingredients', payload);

export const updateIngredient = (id: number, payload: IIngredientPayload): Promise<IngredientDto> =>
    API.Client.put<IngredientDto>(`ingredients/${id}`, payload);

/** 409 = składnik używany w przepisach (body: { usedIn }). */
export const deleteIngredient = (id: number) =>
    API.Client.delete<{ deleted: number }>(`ingredients/${id}`);
