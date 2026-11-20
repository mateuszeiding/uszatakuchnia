import { computed, ref, watch } from 'vue';

/** Wybrany zamiennik wiersza: składnik + „pół na pół”. */
export interface ISubChoice {
    ingredientId: number;
    mix: boolean;
}

export interface ICookingProgress {
    /** Odhaczone składniki (sortOrder). */
    ing: number[];
    /** Zrobione kroki (stepNo). */
    steps: number[];
    servings: number;
    /** Wybrane zamienniki: sortOrder wiersza → zamiennik. Nieaktualne wpisy pomija substitutes.ts. */
    subs: Record<number, ISubChoice>;
}

/** Odhaczenia sprzed „Odznacz wszystko” — do „Cofnij”. */
export interface ICheckSnapshot {
    ing: number[];
    steps: number[];
}

// Postęp gotowania zapamiętany lokalnie, osobno dla każdego przepisu (klucz uk:progress:<id>).
// Strona przepisu, tryb gotowania i lista zakupów czytają ten sam stan (porcje, zamienniki).
export function useCookingProgress(recipeId: number, defaultServings: number) {
    const key = `uk:progress:${recipeId}`;
    const blank = (): ICookingProgress => ({
        ing: [],
        steps: [],
        servings: defaultServings,
        subs: {},
    });

    function load(): ICookingProgress {
        try {
            const raw = localStorage.getItem(key);
            if (!raw) return blank();
            const p = { ...blank(), ...(JSON.parse(raw) as Partial<ICookingProgress>) };
            if (!p.subs || typeof p.subs !== 'object') p.subs = {};
            return p;
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

    /** null = składnik główny („w przepisie”). */
    const setSub = (sortOrder: number, choice: ISubChoice | null) => {
        const subs = { ...progress.value.subs };
        if (choice) subs[sortOrder] = choice;
        else delete subs[sortOrder];
        progress.value.subs = subs;
    };

    /** „Przywróć przepis” — wszystkie wiersze wracają do składników głównych. */
    const clearSubs = () => (progress.value.subs = {});

    /** Coś odhaczone (składnik albo krok) — „Odznacz wszystko” aktywne. */
    const hasChecks = computed(
        () => progress.value.ing.length > 0 || progress.value.steps.length > 0
    );

    // „Odznacz wszystko” (20.11) czyści TYLKO odhaczenia; porcje mają licznik, zamienniki „Przywróć przepis”.
    // Zwraca kopię do „Cofnij”.
    const uncheckAll = (): ICheckSnapshot => {
        const snap = { ing: [...progress.value.ing], steps: [...progress.value.steps] };
        progress.value = { ...progress.value, ing: [], steps: [] };
        return snap;
    };
    const restoreChecks = (snap: ICheckSnapshot) =>
        (progress.value = { ...progress.value, ing: [...snap.ing], steps: [...snap.steps] });

    return {
        progress,
        toggle,
        setStep,
        setServings,
        setSub,
        clearSubs,
        hasChecks,
        uncheckAll,
        restoreChecks,
    };
}
