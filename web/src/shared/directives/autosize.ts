// v-autosize: textarea rośnie (i maleje) razem z treścią. Minimum daje `rows` / min-height z DS.
import type { Directive } from 'vue';

type El = HTMLTextAreaElement & { _autosize?: { ro: ResizeObserver; onInput: () => void } };

function fit(el: HTMLTextAreaElement) {
    if (!el.isConnected) return;
    // „auto” na chwilę zwija pole — zachowujemy przewinięcie strony, żeby nie skakała.
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
        // Zmiana szerokości (okno, panel) zmienia łamanie linii.
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
    // Wartość ustawiona z kodu (wczytanie przepisu, przeniesienie wskazówki).
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
