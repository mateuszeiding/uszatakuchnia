import type { QueryClient } from '@tanstack/vue-query';

import { API, type EndpointsConfig } from '../API';

import type {
    IChallenge,
    IChallengeFormat,
    IUpsertChallengeRequest,
} from '@/data/dtos/challenge/ChallengeDto';

export const CHALLENGES_KEY = ['challenges/list'];
const STALE = 60 * 1000;

/** Wszystkie wyzwania, najnowsze na górze (losowanie i liczniki liczone w przeglądarce). */
export const fetchChallengesWithClient = (qc: QueryClient) =>
    qc.fetchQuery({
        queryKey: CHALLENGES_KEY,
        queryFn: () => API.Client.get<IChallenge[]>('challenges/list'),
        staleTime: STALE,
    });

/** Wyzwanie po id. Z listy w pamięci, jeśli już jest (bez szkieletu po kliknięciu karty). */
export const fetchChallengeWithClient = (qc: QueryClient, id: number) =>
    qc.fetchQuery({
        queryKey: ['challenges/', id],
        queryFn: () =>
            API.Client.get<IChallenge>(`challenges/${id}` as NestedKeyUnion<EndpointsConfig, '/'>),
        initialData: () => qc.getQueryData<IChallenge[]>(CHALLENGES_KEY)?.find((c) => c.id === id),
        initialDataUpdatedAt: () => qc.getQueryState(CHALLENGES_KEY)?.dataUpdatedAt,
        staleTime: STALE,
    });

/** Po zapisie / usunięciu: świeża lista i szczegóły. */
export const invalidateChallenges = (qc: QueryClient) =>
    qc.invalidateQueries({ predicate: (q) => String(q.queryKey[0]).startsWith('challenges/') });

export const createChallenge = (body: IUpsertChallengeRequest) =>
    API.Client.post<IChallenge>('challenges', body);

export const updateChallenge = (id: number, body: IUpsertChallengeRequest) =>
    API.Client.put<IChallenge>(`challenges/${id}`, body);

export const deleteChallenge = (id: number) =>
    API.Client.delete<{ deleted: number }>(`challenges/${id}`);

// ─── Formaty (Panel → Klasyfikacja → Formaty wyzwań; lista jest w słownikach) ───

export const createChallengeFormat = (name: string) =>
    API.Client.post<IChallengeFormat>('taxonomy/challenge-formats', { name });

export const updateChallengeFormat = (id: number, name: string) =>
    API.Client.put<IChallengeFormat>(`taxonomy/challenge-formats/${id}`, { name });

/** Używany format: moveTo = format, do którego przechodzą wyzwania (format jest wymagany). */
export const deleteChallengeFormat = (id: number, moveTo?: number) =>
    API.Client.delete<{ deleted: number; moved: number }>(
        `taxonomy/challenge-formats/${id}${moveTo ? `?moveTo=${moveTo}` : ''}`
    );
