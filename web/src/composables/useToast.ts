import { ref } from 'vue';

export interface IToastAction {
    label: string;
    run: () => void;
}
export interface IToast {
    message: string;
    /** Toast button (e.g. "Cofnij"): the toast becomes clickable and doesn't disappear under the cursor / focus. */
    action?: IToastAction;
}

export const toast = ref<IToast | null>(null);
/** Plain text: for older code and screen readers. */
export const toastMessage = ref<string | null>(null);

let timer: ReturnType<typeof setTimeout> | undefined;
let endsAt = 0;
let remaining = 0;

function schedule(ms: number) {
    clearTimeout(timer);
    remaining = ms;
    endsAt = Date.now() + ms;
    timer = setTimeout(hideToast, ms);
}

export function hideToast() {
    clearTimeout(timer);
    toast.value = null;
    toastMessage.value = null;
}

/** Short message at the bottom of the screen (role="status"). With an action: 8 s by default. */
export function showToast(
    message: string,
    opts: { action?: IToastAction; duration?: number } = {}
) {
    toast.value = { message, action: opts.action };
    toastMessage.value = message;
    schedule(opts.duration ?? (opts.action ? 8000 : 2600));
}

/** Cursor / focus on a toast with an action: pause the countdown (it doesn't vanish under your hand). */
export function pauseToast() {
    if (!toast.value) return;
    clearTimeout(timer);
    remaining = Math.max(0, endsAt - Date.now());
}
export function resumeToast() {
    if (!toast.value) return;
    schedule(Math.max(remaining, 1500));
}
