import { onBeforeUnmount, onMounted } from 'vue';

/** <meta name="robots" content="noindex"> while the view is alive (error pages). */
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
