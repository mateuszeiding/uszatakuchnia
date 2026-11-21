import { auth0 } from './createAuth';

import { confirmAction } from '@/composables/useConfirm';
import { showToast } from '@/composables/useToast';

let pending: Promise<boolean> | null = null;

/**
 * The API answered 401 (the session expired mid-work). Asks to log in again in a popup;
 * the page and the filled-in form stay untouched. Resolves true when login succeeded
 * and the request can be retried. Several parallel 401s = one prompt.
 */
export function reauthenticate(): Promise<boolean> {
    pending ??= ask().finally(() => (pending = null));
    return pending;
}

async function ask(): Promise<boolean> {
    const ok = await confirmAction({
        title: 'Sesja wygasła',
        body: 'Zaloguj się ponownie, żeby zapisać. Wpisana treść zostaje na stronie.',
        note: 'Logowanie otworzy się w nowym okienku.',
        confirmLabel: 'Zaloguj się',
    });
    if (!ok) return false;
    try {
        await auth0.loginWithPopup();
        return auth0.isAuthenticated.value;
    } catch {
        showToast('Nie udało się zalogować. Spróbuj ponownie.');
        return false;
    }
}
