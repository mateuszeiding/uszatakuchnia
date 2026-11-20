import { createAuth0 } from '@auth0/auth0-vue';
import { computed, watch } from 'vue';

// One Auth0 instance per app; the API client uses it too (outside setup()).
export const auth0 = createAuth0({
    domain: import.meta.env.VITE_AUTH0_DOMAIN,
    clientId: import.meta.env.VITE_AUTH0_CLIENT_ID,
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

/** Waits until Auth0 has checked the session after page load. */
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

/** Token for the API, or null (logged out / Auth0 session expired). */
export async function accessToken(): Promise<string | null> {
    await authReady();
    if (!auth0.isAuthenticated.value) return null;
    try {
        return await auth0.getAccessTokenSilently();
    } catch {
        return null;
    }
}
