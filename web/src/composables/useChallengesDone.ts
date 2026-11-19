import { computed, ref } from 'vue';

// „Zrobione” wyzwania na tym urządzeniu: localStorage `uk:challenges:done` = lista id.
// Jeden stan na aplikację (lista, strona wyzwania i losowanie widzą to samo), zmiany z innej karty
// przeglądarki też. Wpisy wyzwań, których już nie ma, są ignorowane w licznikach (filtruje wywołujący).
const KEY = 'uk:challenges:done';

function load(): number[] {
    try {
        const raw = JSON.parse(localStorage.getItem(KEY) ?? '[]') as unknown;
        return Array.isArray(raw) ? raw.filter((x): x is number => Number.isInteger(x)) : [];
    } catch {
        return [];
    }
}

const ids = ref<number[]>(load());

function save() {
    try {
        localStorage.setItem(KEY, JSON.stringify(ids.value));
    } catch {
        /* tryb prywatny / brak miejsca — działa do końca sesji */
    }
}

if (typeof window !== 'undefined') {
    window.addEventListener('storage', (e) => {
        if (e.key === KEY) ids.value = load();
    });
}

export function useChallengesDone() {
    const done = computed(() => new Set(ids.value));
    const isDone = (id: number) => done.value.has(id);
    const toggle = (id: number) => {
        ids.value = isDone(id) ? ids.value.filter((x) => x !== id) : [...ids.value, id];
        save();
    };
    return { done, isDone, toggle };
}
