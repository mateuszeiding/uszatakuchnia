import { onBeforeUnmount, ref } from 'vue';

/** Telefon / wąskie okno (≤ 900 px — ten sam próg co układ strony przepisu). */
export function useNarrow() {
    const mq = window.matchMedia('(max-width: 900px)');
    const narrow = ref(mq.matches);
    const on = (e: MediaQueryListEvent) => (narrow.value = e.matches);
    mq.addEventListener('change', on);
    onBeforeUnmount(() => mq.removeEventListener('change', on));
    return narrow;
}
