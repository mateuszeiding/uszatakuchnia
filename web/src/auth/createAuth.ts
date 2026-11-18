import { createAuth0 } from '@auth0/auth0-vue';
import { computed, watch } from 'vue';

// Jedna instancja Auth0 na aplikację — używa jej też klient API (poza setup()).
export const auth0 = createAuth0({
    domain: import.meta.env.VITE_AUTH0_DOMAIN,
    clientId: import.meta.env.VITE_AUTH0_CLIENT_ID,
    authorizationParams: {
        redirect_uri: globalThis.location.origin,
        // token dostępu jako JWT dla naszego API (backend sprawdza podpis, audience i `sub`)
        audience: import.meta.env.VITE_AUTH0_AUDIENCE,
    },
});

export const createAuth = () => auth0;

/**
 * Zalogowane konto właściciela. Tylko ono widzi UI admina; backend sprawdza to samo
 * (`sub` z tokenu = AUTH0_OWNER_SUB / VITE_AUTH0_OWNER_SUB).
 */
export const isOwner = computed(
    () =>
        auth0.isAuthenticated.value &&
        !!import.meta.env.VITE_AUTH0_OWNER_SUB &&
        auth0.user.value?.sub === import.meta.env.VITE_AUTH0_OWNER_SUB
);

/** Czeka, aż Auth0 sprawdzi sesję po załadowaniu strony. */
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

/** Token dla API albo null (niezalogowany / sesja Auth0 wygasła). */
export async function accessToken(): Promise<string | null> {
    await authReady();
    if (!auth0.isAuthenticated.value) return null;
    try {
        return await auth0.getAccessTokenSilently();
    } catch {
        return null;
    }
}
