import axios, { type AxiosError, type InternalAxiosRequestConfig } from 'axios';

import { accessToken, isOwner } from '@/auth/createAuth';
import { reauthenticate } from '@/auth/reauthenticate';
import { showToast } from '@/composables/useToast';

type RetriableConfig = InternalAxiosRequestConfig & { _reauth?: boolean; _ownerNoToken?: boolean };

const http = axios.create({
    baseURL: `${import.meta.env.DEV ? 'http://localhost:3000' : ''}/api`,
});

// Every request from a logged-in user carries the Auth0 token (writes and drafts require the owner).
http.interceptors.request.use(async (config) => {
    const token = await accessToken();
    if (token) config.headers.Authorization = `Bearer ${token}`;
    // The owner still looks logged in (id token cached) but the refresh token is gone, so the
    // request goes out anonymously. A draft read then gets 404 rather than 401.
    else if (isOwner.value) (config as RetriableConfig)._ownerNoToken = true;
    return config;
});

// 401 = session expired: log in via a popup and retry the request (the form stays on the page).
// A 404 sent without the owner's token may be a draft, so it gets the same recovery.
http.interceptors.response.use(undefined, async (error: AxiosError) => {
    const config = error.config as RetriableConfig | undefined;
    const status = error.response?.status;
    const sessionLost = status === 401 || (status === 404 && config?._ownerNoToken);
    if (sessionLost && config && !config._reauth) {
        config._reauth = true;
        if (await reauthenticate()) return http.request(config);
    } else if (status === 403) {
        showToast('To konto nie ma uprawnień do zmian.');
    }
    throw error;
});

export class API {
    static #instance: API | null = null;

    private constructor() {}

    public static get Client(): API {
        if (API.#instance === null) {
            API.#instance = new API();
        }

        return API.#instance;
    }

    public async get<T>(
        endpoint: NestedKeyUnion<EndpointsConfig, '/'>,
        params?: Record<string, string>
    ) {
        const response = await http.get<T>(endpoint, { params });
        return response.data;
    }

    public async post<T>(endpoint: string, body: unknown) {
        const response = await http.post<T>(endpoint, body);
        return response.data;
    }

    public async upload<T>(endpoint: string, file: Blob, onProgress?: (pct: number) => void) {
        const response = await http.post<T>(endpoint, file, {
            headers: { 'Content-Type': file.type },
            onUploadProgress: (e) =>
                e.total && onProgress?.(Math.round((e.loaded / e.total) * 100)),
        });
        return response.data;
    }

    /** API resource URL (e.g. photos at /api/photos/1); locally the API runs on another port. */
    public static resolveUrl(url: string | null | undefined): string | null {
        if (!url) return null;
        return import.meta.env.DEV && url.startsWith('/api/') ? `http://localhost:3000${url}` : url;
    }

    public async put<T>(endpoint: string, body: unknown) {
        const response = await http.put<T>(endpoint, body);
        return response.data;
    }

    public async delete<T>(endpoint: string) {
        const response = await http.delete<T>(endpoint);
        return response.data;
    }
}

export type EndpointsConfig = {
    enums: {
        aroma: string;
        ingredient: string;
        taste: string;
    };
    ingredients: {
        list: string;
    };
    recipes: {
        list: string;
        [key: number]: string;
    };
    taxonomy: {
        list: string;
    };
    challenges: {
        list: string;
        [key: number]: string;
    };
};
