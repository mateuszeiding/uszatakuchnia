import { onBeforeUnmount, onMounted } from 'vue';

export function useNoindex() {
    let el: HTMLMetaElement | null = null;
    onMounted(() => {
        el = document.createElement('meta');
        el.name = 'robots';
        el.content = 'noindex';
        document.head.appendChild(el);
    });
    onBeforeUnmount(() => el?.remove());
}
