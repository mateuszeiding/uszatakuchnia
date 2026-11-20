import { shallowRef } from 'vue';

export interface IConfirmOptions {
    title: string;
    body?: string;
    /** List below the text (e.g. recipes linking to the one being deleted). */
    items?: string[];
    /** Footnote in smaller text ("Tej operacji nie da się cofnąć."). */
    note?: string;
    confirmLabel?: string;
    danger?: boolean;
}

interface IPendingConfirm extends IConfirmOptions {
    resolve: (ok: boolean) => void;
}

// One confirmation modal per app (UConfirmHost in App.vue).
export const pendingConfirm = shallowRef<IPendingConfirm | null>(null);

/** Shows the confirmation modal; resolves true if the user confirmed. */
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
