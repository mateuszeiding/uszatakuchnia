import type { QueryClient } from '@tanstack/vue-query';

import { API, type EndpointsConfig } from '../API';

import type {
    IChallenge,
    IChallengeFormat,
    IUpsertChallengeRequest,
} from '@/data/dtos/challenge/ChallengeDto';

export const CHALLENGES_KEY = ['challenges/list'];
const STALE = 60 * 1000;

/** All challenges, newest first (the random pick and counters are computed in the browser). */
export const fetchChallengesWithClient = (qc: QueryClient) =>
    qc.fetchQuery({
        queryKey: CHALLENGES_KEY,
        queryFn: () => API.Client.get<IChallenge[]>('challenges/list'),
        staleTime: STALE,
    });

/** Challenge by id. From the in-memory list if it's there (no skeleton after clicking a card). */
export const fetchChallengeWithClient = (qc: QueryClient, id: number) =>
    qc.fetchQuery({
        queryKey: ['challenges/', id],
        queryFn: () =>
            API.Client.get<IChallenge>(`challenges/${id}` as NestedKeyUnion<EndpointsConfig, '/'>),
        initialData: () => qc.getQueryData<IChallenge[]>(CHALLENGES_KEY)?.find((c) => c.id === id),
        initialDataUpdatedAt: () => qc.getQueryState(CHALLENGES_KEY)?.dataUpdatedAt,
        staleTime: STALE,
    });

/** After save / delete: fresh list and details. */
export const invalidateChallenges = (qc: QueryClient) =>
    qc.invalidateQueries({ predicate: (q) => String(q.queryKey[0]).startsWith('challenges/') });

export const createChallenge = (body: IUpsertChallengeRequest) =>
    API.Client.post<IChallenge>('challenges', body);

export const updateChallenge = (id: number, body: IUpsertChallengeRequest) =>
    API.Client.put<IChallenge>(`challenges/${id}`, body);

export const deleteChallenge = (id: number) =>
    API.Client.delete<{ deleted: number }>(`challenges/${id}`);

// ─── Formats (admin → Classification → Challenge formats; the list comes with the dictionaries) ───

export const createChallengeFormat = (name: string) =>
    API.Client.post<IChallengeFormat>('taxonomy/challenge-formats', { name });

export const updateChallengeFormat = (id: number, name: string) =>
    API.Client.put<IChallengeFormat>(`taxonomy/challenge-formats/${id}`, { name });

/** Format in use: moveTo = the format its challenges move to (a format is required). */
export const deleteChallengeFormat = (id: number, moveTo?: number) =>
    API.Client.delete<{ deleted: number; moved: number }>(
        `taxonomy/challenge-formats/${id}${moveTo ? `?moveTo=${moveTo}` : ''}`
    );
