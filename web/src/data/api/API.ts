import axios from 'axios';

export class API {
    static #instance: API | null = null;

    private readonly __baseUrl = `${import.meta.env.DEV ? 'http://localhost:3000' : ''}/api`;

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
        const response = await axios.get<T>(endpoint, {
            baseURL: this.__baseUrl,
            params,
        });
        return response.data;
    }

    public async post<T>(endpoint: string, body: unknown) {
        const response = await axios.post<T>(endpoint, body, { baseURL: this.__baseUrl });
        return response.data;
    }

    /** Wysyła surowe bajty (np. obraz) z raportowaniem postępu 0–100. */
    public async upload<T>(endpoint: string, file: Blob, onProgress?: (pct: number) => void) {
        const response = await axios.post<T>(endpoint, file, {
            baseURL: this.__baseUrl,
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
        const response = await axios.put<T>(endpoint, body, { baseURL: this.__baseUrl });
        return response.data;
    }

    public async delete<T>(endpoint: string) {
        const response = await axios.delete<T>(endpoint, { baseURL: this.__baseUrl });
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
