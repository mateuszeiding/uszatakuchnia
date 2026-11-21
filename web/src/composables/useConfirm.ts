import { shallowRef } from 'vue';

export interface IConfirmOptions {
    title: string;
    body?: string;
    items?: string[];
    note?: string;
    confirmLabel?: string;
    danger?: boolean;
}

interface IPendingConfirm extends IConfirmOptions {
    resolve: (ok: boolean) => void;
}

// One confirmation modal per app (UConfirmHost in App.vue).
export const pendingConfirm = shallowRef<IPendingConfirm | null>(null);

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
