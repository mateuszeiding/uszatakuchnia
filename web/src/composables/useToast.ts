import { ref } from 'vue';

export const toastMessage = ref<string | null>(null);
let timer: ReturnType<typeof setTimeout> | undefined;

/** Krótki komunikat na dole ekranu (role="status"). */
export function showToast(message: string) {
    toastMessage.value = message;
    clearTimeout(timer);
    timer = setTimeout(() => (toastMessage.value = null), 2600);
}
