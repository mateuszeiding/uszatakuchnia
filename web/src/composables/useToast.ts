import { ref } from 'vue';

export interface IToastAction {
    label: string;
    run: () => void;
}
export interface IToast {
    message: string;
    /** Przycisk w toaście (np. „Cofnij”) — toast jest wtedy klikalny i nie znika pod kursorem / fokusem. */
    action?: IToastAction;
}

export const toast = ref<IToast | null>(null);
/** Sam tekst — dla starego kodu i czytnika ekranu. */
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

/** Krótki komunikat na dole ekranu (role="status"). Z akcją: domyślnie 8 s. */
export function showToast(
    message: string,
    opts: { action?: IToastAction; duration?: number } = {}
) {
    toast.value = { message, action: opts.action };
    toastMessage.value = message;
    schedule(opts.duration ?? (opts.action ? 8000 : 2600));
}

/** Kursor / fokus na toaście z akcją — wstrzymaj odliczanie (nie znika pod ręką). */
export function pauseToast() {
    if (!toast.value) return;
    clearTimeout(timer);
    remaining = Math.max(0, endsAt - Date.now());
}
export function resumeToast() {
    if (!toast.value) return;
    schedule(Math.max(remaining, 1500));
}
