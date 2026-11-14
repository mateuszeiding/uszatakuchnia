import { ref, watch } from 'vue';

// Odhaczenia na liście zakupów („mam w domu / w koszyku”) — osobno od postępu gotowania.
// localStorage `uk:shopping:<recipeId>` = tablica ingredientId; zmiana porcji ich nie czyści.
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
                /* tryb prywatny / brak miejsca — odhaczenia działają tylko w tej sesji */
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
