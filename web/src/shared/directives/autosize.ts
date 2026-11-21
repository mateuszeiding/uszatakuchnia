// v-autosize: the textarea grows (and shrinks) with its content. The minimum comes from `rows` / min-height in the DS.
import type { Directive } from 'vue';

type El = HTMLTextAreaElement & { _autosize?: { ro: ResizeObserver; onInput: () => void } };

function fit(el: HTMLTextAreaElement) {
    if (!el.isConnected) return;
    // Setting "auto" briefly collapses the field, so we keep the page scroll position to avoid a jump.
    const y = window.scrollY;
    el.style.height = 'auto';
    const border = el.offsetHeight - el.clientHeight;
    el.style.height = `${el.scrollHeight + border}px`;
    if (window.scrollY !== y) window.scrollTo({ top: y });
}

export const vAutosize: Directive<El> = {
    mounted(el) {
        el.classList.add('textarea--auto');
        const onInput = () => fit(el);
        let width = el.clientWidth;
        // A width change (window, panel) changes line wrapping.
        const ro = new ResizeObserver(() => {
            if (el.clientWidth === width) return;
            width = el.clientWidth;
            fit(el);
        });
        ro.observe(el);
        el.addEventListener('input', onInput);
        el._autosize = { ro, onInput };
        fit(el);
    },
    // Value set from code (recipe load, moving a hint).
    updated(el) {
        fit(el);
    },
    unmounted(el) {
        el._autosize?.ro.disconnect();
        if (el._autosize) el.removeEventListener('input', el._autosize.onInput);
        delete el._autosize;
    },
};

declare module 'vue' {
    // eslint-disable-next-line @typescript-eslint/naming-convention -- nazwa narzucona przez Vue
    interface GlobalDirectives {
        vAutosize: typeof vAutosize;
    }
}
