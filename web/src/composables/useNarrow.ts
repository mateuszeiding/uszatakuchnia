import { onBeforeUnmount, ref } from 'vue';

/** Phone / narrow window (≤ 900 px, the same breakpoint as the recipe page layout). */
export function useNarrow() {
    const mq = window.matchMedia('(max-width: 900px)');
    const narrow = ref(mq.matches);
    const on = (e: MediaQueryListEvent) => (narrow.value = e.matches);
    mq.addEventListener('change', on);
    onBeforeUnmount(() => mq.removeEventListener('change', on));
    return narrow;
}
