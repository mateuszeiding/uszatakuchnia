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

    public async put<T>(endpoint: string, body: unknown) {
        const response = await axios.put<T>(endpoint, body, { baseURL: this.__baseUrl });
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
};
