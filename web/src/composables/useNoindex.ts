import { onBeforeUnmount, onMounted } from 'vue';

/** <meta name="robots" content="noindex"> na czas życia widoku (strony błędów, 20.11). */
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
