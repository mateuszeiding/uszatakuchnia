import { computed, nextTick, onBeforeUnmount, onMounted, type Ref, ref } from 'vue';

export type CookLayout = 'desktop' | 'portrait' | 'landscape';

// Fixed steps, not clamp(): two steps of similar length should not "breathe" between sizes.
const PORTRAIT_SIZES = [26, 23, 20] as const;
const LANDSCAPE_SIZE = 21;

function useMedia(query: string) {
    const mq = window.matchMedia(query);
    const on = ref(mq.matches);
    const set = (e: MediaQueryListEvent) => (on.value = e.matches);
    mq.addEventListener('change', set);
    onBeforeUnmount(() => mq.removeEventListener('change', set));
    return on;
}

/**
 * Fits a cooking step into the phone screen without scrolling, in order: text 26 → 23 → 20 px with hints
 * expanded; still too long → hints collapsed into one button; still too long (or hints opened) → the stage
 * scrolls and `more` shows "Ciąg dalszy". Landscape: 21 px, hints always collapsed. Desktop: untouched.
 */
export function useCookFit(stage: Ref<HTMLElement | undefined>) {
    const narrow = useMedia('(max-width: 900px)');
    const short = useMedia('(orientation: landscape) and (max-height: 500px)');
    const layout = computed<CookLayout>(() =>
        short.value ? 'landscape' : narrow.value ? 'portrait' : 'desktop'
    );

    const size = ref<number | null>(null);
    /** Hints don't fit next to the text and go into a button (until the reader opens them). */
    const tight = ref(false);
    /** While measuring, hints are expanded regardless of the button. */
    const measuring = ref(false);
    const more = ref(false);

    const overflows = () => {
        const el = stage.value;
        return !!el && el.scrollHeight - el.clientHeight > 1;
    };

    function updateMore() {
        const el = stage.value;
        more.value =
            layout.value !== 'desktop' &&
            !!el &&
            el.scrollTop + el.clientHeight < el.scrollHeight - 2;
    }

    let run = 0;
    async function fit() {
        const id = ++run;
        measuring.value = true;
        if (layout.value === 'desktop') {
            size.value = null;
            tight.value = false;
        } else if (layout.value === 'landscape') {
            size.value = LANDSCAPE_SIZE;
            tight.value = true;
        } else {
            tight.value = false;
            for (const s of PORTRAIT_SIZES) {
                size.value = s;
                await nextTick();
                if (id !== run) return;
                if (!overflows()) break;
            }
            tight.value = overflows();
        }
        measuring.value = false;
        await nextTick();
        if (id === run) updateMore();
    }

    let ro: ResizeObserver | null = null;
    let frame = 0;
    onMounted(() => {
        ro = new ResizeObserver(() => {
            cancelAnimationFrame(frame);
            frame = requestAnimationFrame(() => void fit());
        });
        if (stage.value) ro.observe(stage.value);
        // the first measurement may run before Geist is loaded (different line breaks)
        document.fonts?.ready.then(() => void fit());
    });
    onBeforeUnmount(() => {
        ro?.disconnect();
        cancelAnimationFrame(frame);
    });

    return { layout, size, tight, measuring, more, fit, updateMore };
}
