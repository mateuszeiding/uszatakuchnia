import { API, type EndpointsConfig } from '../API';

import type { ITerm } from '@/shared/utils/stepRefs';

export interface IGlossaryTerm extends ITerm {
    /** Recipes (drafts too) that mark the term; only for the owner. */
    usedIn?: number;
}

export interface ITermUsage {
    id: number;
    name: string;
    status: 'draft' | 'published';
}

export interface IUpsertTerm {
    name: string;
    explanation: string;
}

const path = (p: string) => p as NestedKeyUnion<EndpointsConfig, '/'>;

/** The whole glossary A–Z (small: one request, search in the browser). */
export const fetchGlossary = () => API.Client.get<IGlossaryTerm[]>(path('glossary'));

export const fetchTermUsage = (id: number) =>
    API.Client.get<ITermUsage[]>(path(`glossary/${id}/usage`));

export const createTerm = (body: IUpsertTerm) => API.Client.post<IGlossaryTerm>('glossary', body);

export const updateTerm = (id: number, body: IUpsertTerm) =>
    API.Client.put<IGlossaryTerm>(`glossary/${id}`, body);

/** A term in use: replaceWith (marks point at another term) or removeMarks (words stay as plain text). */
export const deleteTerm = (
    id: number,
    opts: { replaceWith?: number; removeMarks?: boolean } = {}
) =>
    API.Client.delete<{ deleted: number; recipes: number }>(
        `glossary/${id}${
            opts.replaceWith
                ? `?replaceWith=${opts.replaceWith}`
                : opts.removeMarks
                  ? '?removeMarks=true'
                  : ''
        }`
    );
