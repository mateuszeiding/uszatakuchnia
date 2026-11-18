import axios, { type AxiosError, type InternalAxiosRequestConfig } from 'axios';

import { accessToken } from '@/auth/createAuth';
import { reauthenticate } from '@/auth/reauthenticate';
import { showToast } from '@/composables/useToast';

type RetriableConfig = InternalAxiosRequestConfig & { _reauth?: boolean };

const http = axios.create({
    baseURL: `${import.meta.env.DEV ? 'http://localhost:3000' : ''}/api`,
});

// Każde żądanie zalogowanego idzie z tokenem Auth0 (zapis i szkice wymagają właściciela).
http.interceptors.request.use(async (config) => {
    const token = await accessToken();
    if (token) config.headers.Authorization = `Bearer ${token}`;
    return config;
});

// 401 = sesja wygasła: zaloguj w okienku i powtórz żądanie (formularz zostaje na stronie).
http.interceptors.response.use(undefined, async (error: AxiosError) => {
    const config = error.config as RetriableConfig | undefined;
    const status = error.response?.status;
    if (status === 401 && config && !config._reauth) {
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

    /** Wysyła surowe bajty (np. obraz) z raportowaniem postępu 0–100. */
    public async upload<T>(endpoint: string, file: Blob, onProgress?: (pct: number) => void) {
        const response = await http.post<T>(endpoint, file, {
            headers: { 'Content-Type': file.type },
            onUploadProgress: (e) =>
                e.total && onProgress?.(Math.round((e.loaded / e.total) * 100)),
        });
        return response.data;
    }

    /** Adres zasobu z API (np. zdjęcia /api/photos/1) — lokalnie API działa na innym porcie. */
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
};
