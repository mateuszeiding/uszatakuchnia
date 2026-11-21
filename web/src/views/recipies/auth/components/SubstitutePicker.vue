<script lang="ts" setup>
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue';

import type { ISubItem } from '../recipeForm';

import type { IngredientDto } from '@/data/dtos/ingredients/IngredientDto';

// Adding a substitute: search the ingredient database (excluding the main ingredient and ones already added both ways);
// each result offers "zamiast" and "pół na pół". ↑ ↓ to pick, Enter = "zamiast" ("pół na pół" if "zamiast" is already added), Esc closes.
const props = defineProps<{
    id: string;
    ingredientName: string;
    mainId: number;
    current: ISubItem[];
    ingredients: IngredientDto[];
}>();
const emit = defineEmits<{ add: [s: ISubItem]; close: [] }>();

const q = ref('');
const hi = ref(0);
const root = ref<HTMLElement>();
const input = ref<HTMLInputElement>();
const listId = computed(() => `${props.id}-list`);

const has = (id: number, mix: boolean) =>
    props.current.some((s) => s.ingredientId === id && s.mix === mix);

const norm = (s: string) => s.toLowerCase().normalize('NFD').replace(/[̀-ͯ]/g, '').replace(/ł/g, 'l');

const results = computed(() => {
    const n = norm(q.value.trim());
    return props.ingredients
        .filter((i) => i.id !== props.mainId && !(has(i.id, false) && has(i.id, true)))
        .filter((i) => !n || norm(i.name).includes(n))
        .sort(
            (a, b) =>
                Number(!norm(a.name).startsWith(n)) - Number(!norm(b.name).startsWith(n)) ||
                a.name.localeCompare(b.name, 'pl')
        )
        .slice(0, 8);
});

watch(q, () => (hi.value = 0));

function add(i: IngredientDto | undefined, mix: boolean) {
    if (!i || has(i.id, mix)) return;
    emit('add', { ingredientId: i.id, name: i.name, mix });
}

async function scrollHi() {
    await nextTick();
    document.getElementById(`${listId.value}-${hi.value}`)?.scrollIntoView({ block: 'nearest' });
}

function onKey(e: KeyboardEvent) {
    const n = results.value.length;
    if (e.key === 'ArrowDown' && n) {
        e.preventDefault();
        hi.value = (hi.value + 1) % n;
        scrollHi();
    } else if (e.key === 'ArrowUp' && n) {
        e.preventDefault();
        hi.value = (hi.value - 1 + n) % n;
        scrollHi();
    } else if (e.key === 'Enter') {
        e.preventDefault();
        const i = results.value[hi.value];
        add(i, !!i && has(i.id, false));
    } else if (e.key === 'Escape') {
        e.preventDefault();
        e.stopPropagation();
        emit('close');
    }
}

const onDocDown = (e: MouseEvent) => {
    if (root.value && !root.value.contains(e.target as Node)) emit('close');
};
onMounted(() => {
    input.value?.focus();
    document.addEventListener('mousedown', onDocDown);
});
onBeforeUnmount(() => document.removeEventListener('mousedown', onDocDown));
</script>

<template>
    <div
        ref="root"
        class="spick"
    >
        <label
            :for="`${id}-q`"
            class="spick__label"
        >
            Zamiennik dla „{{ ingredientName || 'składnika' }}”
        </label>
        <input
            :id="`${id}-q`"
            ref="input"
            v-model="q"
            class="input spick__input"
            role="combobox"
            aria-autocomplete="list"
            aria-expanded="true"
            :aria-controls="listId"
            :aria-activedescendant="results.length ? `${listId}-${hi}` : undefined"
            placeholder="Szukaj składnika"
            autocomplete="off"
            @keydown="onKey"
        />
        <ul
            :id="listId"
            role="listbox"
            aria-label="Składniki"
            class="spick__list"
        >
            <li
                v-for="(i, k) in results"
                :id="`${listId}-${k}`"
                :key="i.id"
                role="option"
                :aria-selected="k === hi"
                class="spick__opt"
                :class="{ 'is-hi': k === hi }"
                @mouseenter="hi = k"
            >
                <span class="spick__name">{{ i.name }}</span>
                <span class="spick__btns">
                    <button
                        type="button"
                        class="spick__btn"
                        :disabled="has(i.id, false)"
                        :aria-label="`${i.name} zamiast „${ingredientName}”`"
                        @click="add(i, false)"
                    >
                        zamiast
                    </button>
                    <button
                        type="button"
                        class="spick__btn"
                        :disabled="has(i.id, true)"
                        :aria-label="`${i.name} pół na pół z „${ingredientName}”`"
                        @click="add(i, true)"
                    >
                        pół na pół
                    </button>
                </span>
            </li>
            <li
                v-if="!results.length"
                class="spick__empty"
            >
                Brak składnika „{{ q }}” w bazie.
            </li>
        </ul>
        <div class="spick__hint">
            Zamiennik ma tę samą ilość. „Pół na pół” = po połowie ilości z tym składnikiem.
        </div>
    </div>
</template>

<style scoped>
.spick {
    position: absolute;
    right: 38px;
    top: calc(100% - 6px);
    z-index: var(--z-dropdown);
    width: 360px;
    padding: 10px;
    display: flex;
    flex-direction: column;
    gap: 8px;
    background: var(--bg-alt);
    border: 1px solid var(--rule-strong);
    border-radius: var(--r-xl);
    box-shadow: var(--shadow-menu);
}
.spick__label {
    font-family: var(--font-mono);
    font-size: 12px;
    font-weight: 600;
    letter-spacing: 0.5px;
    text-transform: uppercase;
    color: var(--ink-muted);
}
.spick__input {
    width: 100%;
}
.spick__list {
    list-style: none;
    margin: 0;
    padding: 0;
    max-height: 240px;
    overflow: auto;
    scrollbar-gutter: stable;
}
.spick__opt {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 8px;
    padding: 6px 8px 6px 10px;
    border-radius: var(--r-md);
    font-size: 14px;
}
.spick__opt.is-hi {
    background: var(--accent-soft);
}
.spick__name {
    min-width: 0;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
}
.spick__btns {
    display: flex;
    gap: 6px;
    flex-shrink: 0;
}
.spick__btn {
    padding: 2px 9px;
    border: 1px solid var(--rule-strong);
    border-radius: var(--r-pill);
    background: var(--bg-alt);
    color: var(--ink);
    font-size: 12px;
    font-weight: 600;
    cursor: pointer;
}
.spick__btn:hover:not(:disabled) {
    border-color: var(--accent);
    color: var(--accent);
}
.spick__btn:disabled {
    opacity: 0.4;
    cursor: default;
}
.spick__empty {
    padding: 10px;
    font-size: 13px;
    color: var(--ink-muted);
}
.spick__hint {
    padding-top: 8px;
    border-top: 1px solid var(--rule);
    font-size: 12px;
    color: var(--ink-muted);
}
</style>
