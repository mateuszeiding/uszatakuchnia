import { computed, ref, watch } from 'vue';

/** The row's chosen substitute: ingredient + "half and half". */
export interface ISubChoice {
    ingredientId: number;
    mix: boolean;
}

export interface ICookingProgress {
    /** Checked ingredients (sortOrder). */
    ing: number[];
    /** Done steps (stepNo). */
    steps: number[];
    servings: number;
    /** Chosen substitutes: row sortOrder → substitute. substitutes.ts skips stale entries. */
    subs: Record<number, ISubChoice>;
}

/** Checks from before "Odznacz wszystko", for "Cofnij". */
export interface ICheckSnapshot {
    ing: number[];
    steps: number[];
}

// Cooking progress stored locally, per recipe (key uk:progress:<id>).
// The recipe page, cooking mode and the shopping list read the same state (servings, substitutes).
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
                /* private mode / no storage: progress lasts only for this session */
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

    /** null = the main ingredient ("w przepisie"). */
    const setSub = (sortOrder: number, choice: ISubChoice | null) => {
        const subs = { ...progress.value.subs };
        if (choice) subs[sortOrder] = choice;
        else delete subs[sortOrder];
        progress.value.subs = subs;
    };

    /** "Przywróć przepis": every row goes back to its main ingredient. */
    const clearSubs = () => (progress.value.subs = {});

    /** Something is checked (an ingredient or a step), so "Odznacz wszystko" is enabled. */
    const hasChecks = computed(
        () => progress.value.ing.length > 0 || progress.value.steps.length > 0
    );

    // "Odznacz wszystko" clears ONLY the checks; servings have a counter, substitutes have "Przywróć przepis".
    // Returns a copy for "Cofnij".
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
