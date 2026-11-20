import { ref, watch } from 'vue';

// Shopping list checks ("mam w domu / w koszyku"), separate from cooking progress.
// localStorage `uk:shopping:<recipeId>` = array of ingredientId; changing servings doesn't clear them.
export function useShoppingChecks(recipeId: number) {
    const key = `uk:shopping:${recipeId}`;

    function load(): number[] {
        try {
            const raw = JSON.parse(localStorage.getItem(key) ?? '[]') as unknown;
            return Array.isArray(raw) ? raw.filter((x): x is number => typeof x === 'number') : [];
        } catch {
            return [];
        }
    }

    const checked = ref<number[]>(load());
    watch(
        checked,
        (c) => {
            try {
                localStorage.setItem(key, JSON.stringify(c));
            } catch {
                /* private mode / no storage: checks last only for this session */
            }
        },
        { deep: true }
    );

    const toggle = (id: number) => {
        checked.value = checked.value.includes(id)
            ? checked.value.filter((x) => x !== id)
            : [...checked.value, id];
    };
    const clear = () => (checked.value = []);

    return { checked, toggle, clear };
}
