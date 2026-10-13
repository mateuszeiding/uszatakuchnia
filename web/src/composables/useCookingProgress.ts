import { computed, ref, watch } from 'vue';

export interface ICookingProgress {
    /** Odhaczone składniki (sortOrder). */
    ing: number[];
    /** Zrobione kroki (stepNo). */
    steps: number[];
    servings: number;
}

// Postęp gotowania zapamiętany lokalnie, osobno dla każdego przepisu (klucz uk:progress:<id>).
export function useCookingProgress(recipeId: number, defaultServings: number) {
    const key = `uk:progress:${recipeId}`;
    const blank = (): ICookingProgress => ({ ing: [], steps: [], servings: defaultServings });

    function load(): ICookingProgress {
        try {
            const raw = localStorage.getItem(key);
            return raw
                ? { ...blank(), ...(JSON.parse(raw) as Partial<ICookingProgress>) }
                : blank();
        } catch {
            return blank();
        }
    }

    const progress = ref<ICookingProgress>(load());

    watch(
        progress,
        (p) => {
            try {
                localStorage.setItem(key, JSON.stringify(p));
            } catch {
                /* tryb prywatny / brak miejsca — postęp działa tylko w tej sesji */
            }
        },
        { deep: true }
    );

    const toggle = (kind: 'ing' | 'steps', id: number) => {
        const list = progress.value[kind];
        progress.value[kind] = list.includes(id) ? list.filter((x) => x !== id) : [...list, id];
    };

    const setStep = (id: number, done: boolean) => {
        const s = new Set(progress.value.steps);
        if (done) s.add(id);
        else s.delete(id);
        progress.value.steps = [...s];
    };

    const setServings = (n: number) => (progress.value.servings = Math.min(20, Math.max(1, n)));

    const hasProgress = computed(
        () =>
            progress.value.ing.length > 0 ||
            progress.value.steps.length > 0 ||
            progress.value.servings !== defaultServings
    );

    const reset = () => (progress.value = blank());

    return { progress, toggle, setStep, setServings, hasProgress, reset };
}
