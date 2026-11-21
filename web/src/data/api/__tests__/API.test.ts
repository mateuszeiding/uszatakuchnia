import { AxiosError, type InternalAxiosRequestConfig } from 'axios';
import { beforeEach, describe, expect, it, vi } from 'vitest';

type OnRequest = (c: InternalAxiosRequestConfig) => Promise<InternalAxiosRequestConfig>;
type OnError = (e: AxiosError) => Promise<unknown>;

const h = vi.hoisted(() => ({
    onRequest: null as OnRequest | null,
    onError: null as OnError | null,
    request: vi.fn(),
    token: null as string | null,
    owner: { value: false },
    reauthenticate: vi.fn(),
}));

// Capture the interceptors API.ts registers, so they run without a network.
vi.mock('axios', async (orig) => {
    const real = await orig<typeof import('axios')>();
    const instance = {
        interceptors: {
            request: { use: (fn: OnRequest) => (h.onRequest = fn) },
            response: { use: (_ok: unknown, fn: OnError) => (h.onError = fn) },
        },
        request: h.request,
    };
    return { ...real, default: { ...real.default, create: () => instance } };
});
vi.mock('@/auth/createAuth', () => ({
    isOwner: h.owner,
    accessToken: () => Promise.resolve(h.token),
}));
vi.mock('@/auth/reauthenticate', () => ({ reauthenticate: h.reauthenticate }));
vi.mock('@/composables/useToast', () => ({ showToast: vi.fn() }));

await import('../API');

async function send(status: number) {
    const config = await h.onRequest!({ headers: {} } as InternalAxiosRequestConfig);
    const error = new AxiosError('fail', undefined, config, undefined, {
        status,
    } as AxiosError['response']);
    return { config, result: h.onError!(error) };
}

describe('API re-login recovery', () => {
    beforeEach(() => {
        h.request.mockReset().mockResolvedValue('retried');
        h.reauthenticate.mockReset().mockResolvedValue(true);
        h.token = null;
        h.owner.value = false;
    });

    it('owner without a token: 404 prompts re-login and retries', async () => {
        h.owner.value = true;
        const { config, result } = await send(404);
        await expect(result).resolves.toBe('retried');
        expect(h.reauthenticate).toHaveBeenCalledOnce();
        expect(h.request).toHaveBeenCalledWith(config);
    });

    it('anonymous visitor: 404 never prompts', async () => {
        const { result } = await send(404);
        await expect(result).rejects.toBeInstanceOf(AxiosError);
        expect(h.reauthenticate).not.toHaveBeenCalled();
    });

    it('owner with a token: 404 is a real 404', async () => {
        h.owner.value = true;
        h.token = 'jwt';
        const { config, result } = await send(404);
        expect(config.headers.Authorization).toBe('Bearer jwt');
        await expect(result).rejects.toBeInstanceOf(AxiosError);
        expect(h.reauthenticate).not.toHaveBeenCalled();
    });

    it('401 still prompts for everyone', async () => {
        const { result } = await send(401);
        await expect(result).resolves.toBe('retried');
        expect(h.reauthenticate).toHaveBeenCalledOnce();
    });

    it('a retried request does not prompt twice', async () => {
        h.owner.value = true;
        const config = await h.onRequest!({ headers: {} } as InternalAxiosRequestConfig);
        (config as { _reauth?: boolean })._reauth = true;
        const again = new AxiosError('fail', undefined, config, undefined, {
            status: 404,
        } as AxiosError['response']);
        await expect(h.onError!(again)).rejects.toBe(again);
        expect(h.reauthenticate).not.toHaveBeenCalled();
    });

    it('declined re-login rethrows the 404', async () => {
        h.owner.value = true;
        h.reauthenticate.mockResolvedValue(false);
        const { result } = await send(404);
        await expect(result).rejects.toBeInstanceOf(AxiosError);
        expect(h.request).not.toHaveBeenCalled();
    });
});
