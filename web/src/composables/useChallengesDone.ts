import { computed, ref } from 'vue';

// Challenges marked "done" on this device: localStorage `uk:challenges:done` = list of ids.
// One state per app (the list, the challenge page and the random pick see the same), changes from another
// browser tab too. Entries of challenges that no longer exist are ignored in counters (the caller filters).
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
        /* private mode / no storage: works until the session ends */
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
