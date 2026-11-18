import { type QueryClient, useQueryClient } from '@tanstack/vue-query';

import { API } from '../API';

import type {
    IDictionaries,
    IEquipment,
    ITaxonomyValue,
    IUnit,
} from '@/data/dtos/taxonomy/DictionariesDto';

export const DICTIONARIES_KEY = ['taxonomy/list'];

export const fetchDictionariesWithClient = (qc: QueryClient) =>
    qc.fetchQuery({
        queryKey: DICTIONARIES_KEY,
        queryFn: () => API.Client.get<IDictionaries>('taxonomy/list'),
        staleTime: 5 * 60 * 1000,
    });

// Use only from setup() context
export const fetchDictionaries = () => fetchDictionariesWithClient(useQueryClient());

export const invalidateDictionaries = (qc: QueryClient) =>
    qc.invalidateQueries({ queryKey: DICTIONARIES_KEY });

export interface ITaxonomyValuePayload {
    axisId?: number;
    label: string;
    color?: string | null;
}

export const createTaxonomyValue = (payload: ITaxonomyValuePayload) =>
    API.Client.post<ITaxonomyValue>('taxonomy/values', payload);

export const updateTaxonomyValue = (id: number, payload: ITaxonomyValuePayload) =>
    API.Client.put<ITaxonomyValue>(`taxonomy/values/${id}`, payload);

/** moveTo: przepisy z usuwaną wartością dostają tę (wymagane dla osi obowiązkowej, gdy wartość używana). */
export const deleteTaxonomyValue = (id: number, moveTo?: number) =>
    API.Client.delete<{ deleted: number; moved: number }>(
        `taxonomy/values/${id}${moveTo ? `?moveTo=${moveTo}` : ''}`
    );

export interface IUnitPayload {
    name: string;
    forms: string[];
    noQty: boolean;
}

export const createUnit = (payload: IUnitPayload) =>
    API.Client.post<IUnit>('taxonomy/units', payload);

export const updateUnit = (id: number, payload: IUnitPayload) =>
    API.Client.put<IUnit>(`taxonomy/units/${id}`, payload);

export const deleteUnit = (id: number, moveTo?: number) =>
    API.Client.delete<{ deleted: number; moved: number }>(
        `taxonomy/units/${id}${moveTo ? `?moveTo=${moveTo}` : ''}`
    );

export const createEquipment = (name: string) =>
    API.Client.post<IEquipment>('taxonomy/equipment', { name });

export const updateEquipment = (id: number, name: string) =>
    API.Client.put<IEquipment>(`taxonomy/equipment/${id}`, { name });

/** Używany sprzęt: replaceWith (inny sprzęt, doprecyzowania zostają) albo removeFromRecipes. */
export const deleteEquipment = (
    id: number,
    opt: { replaceWith?: number; removeFromRecipes?: boolean } = {}
) =>
    API.Client.delete<{ deleted: number; recipes: number }>(
        `taxonomy/equipment/${id}${
            opt.replaceWith
                ? `?replaceWith=${opt.replaceWith}`
                : opt.removeFromRecipes
                  ? '?removeFromRecipes=true'
                  : ''
        }`
    );
