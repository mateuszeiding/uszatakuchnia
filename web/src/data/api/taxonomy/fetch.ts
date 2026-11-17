import { type QueryClient, useQueryClient } from '@tanstack/vue-query';

import { API } from '../API';

import type { IDictionaries, ITaxonomyValue, IUnit } from '@/data/dtos/taxonomy/DictionariesDto';

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
