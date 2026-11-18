import { auth0 } from './createAuth';

import { confirmAction } from '@/composables/useConfirm';
import { showToast } from '@/composables/useToast';

let pending: Promise<boolean> | null = null;

/**
 * API odpowiedziało 401 (sesja wygasła w trakcie pracy). Pyta o ponowne logowanie
 * w okienku — strona i wpisany formularz zostają nietknięte. Zwraca true, gdy
 * logowanie się udało i żądanie można powtórzyć. Kilka równoległych 401 = jedno pytanie.
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
