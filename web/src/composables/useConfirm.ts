import { shallowRef } from 'vue';

export interface IConfirmOptions {
    title: string;
    body?: string;
    confirmLabel?: string;
    danger?: boolean;
}

interface IPendingConfirm extends IConfirmOptions {
    resolve: (ok: boolean) => void;
}

// Jeden modal potwierdzenia na aplikację (UConfirmHost w App.vue).
export const pendingConfirm = shallowRef<IPendingConfirm | null>(null);

/** Pokazuje modal potwierdzenia; zwraca true, jeśli użytkownik potwierdził. */
export function confirmAction(options: IConfirmOptions): Promise<boolean> {
    pendingConfirm.value?.resolve(false);
    return new Promise((resolve) => {
        pendingConfirm.value = { ...options, resolve };
    });
}

export function settleConfirm(ok: boolean) {
    const p = pendingConfirm.value;
    pendingConfirm.value = null;
    p?.resolve(ok);
}
