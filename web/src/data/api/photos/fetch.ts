import { API } from '../API';

export const uploadPhoto = (file: Blob, onProgress?: (pct: number) => void) =>
    API.Client.upload<{ id: number; url: string }>('photos', file, onProgress);
