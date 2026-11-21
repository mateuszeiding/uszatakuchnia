import { createAuth0 } from '@auth0/auth0-vue';
import { computed, watch } from 'vue';

// One Auth0 instance per app; the API client uses it too (outside setup()).
export const auth0 = createAuth0({
    domain: import.meta.env.VITE_AUTH0_DOMAIN,
    clientId: import.meta.env.VITE_AUTH0_CLIENT_ID,
    // The default in-memory cache restores the session through a hidden iframe on every
    // new tab or reload, and browsers block Auth0's third-party cookies there. Storing
    // tokens in localStorage and renewing them with a rotating refresh token keeps the
    // session without the iframe. No iframe fallback: it fails the same way, so a missing
    // or revoked refresh token ends in a 401 and the re-login popup (reauthenticate.ts).
    cacheLocation: 'localstorage',
    useRefreshTokens: true,
    authorizationParams: {
        redirect_uri: globalThis.location.origin,
        // access token as a JWT for our API (the backend checks signature, audience and `sub`)
        audience: import.meta.env.VITE_AUTH0_AUDIENCE,
    },
});

export const createAuth = () => auth0;

/**
 * The owner's logged-in account. Only it sees the admin UI; the backend checks the same
 * (`sub` from the token = AUTH0_OWNER_SUB / VITE_AUTH0_OWNER_SUB).
 */
export const isOwner = computed(
    () =>
        auth0.isAuthenticated.value &&
        !!import.meta.env.VITE_AUTH0_OWNER_SUB &&
        auth0.user.value?.sub === import.meta.env.VITE_AUTH0_OWNER_SUB
);

export function authReady(): Promise<void> {
    if (!auth0.isLoading.value) return Promise.resolve();
    return new Promise((resolve) => {
        const stop = watch(auth0.isLoading, (loading) => {
            if (loading) return;
            stop();
            resolve();
        });
    });
}

export async function accessToken(): Promise<string | null> {
    await authReady();
    if (!auth0.isAuthenticated.value) return null;
    try {
        return await auth0.getAccessTokenSilently();
    } catch {
        return null;
    }
}
