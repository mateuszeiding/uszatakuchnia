<script lang="ts" setup>
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue';

import type { ILinkedRef } from '../recipeForm';

import { mainDish, type RecipeBaseDto } from '@/data/dtos/recipe/RecipeDto';

// Search by title (also without Polish characters), excluding the current recipe; drafts last with a label.
const props = defineProps<{
    id: string;
    ingredientName: string;
    recipes: RecipeBaseDto[];
    /** The edited recipe: it can't link to itself. */
    currentId: number | null;
}>();
const emit = defineEmits<{ select: [r: ILinkedRef]; close: [] }>();

const q = ref('');
const hi = ref(0);
const root = ref<HTMLElement>();
const input = ref<HTMLInputElement>();
const listId = computed(() => `${props.id}-list`);

const norm = (s: string) => s.toLowerCase().normalize('NFD').replace(/[̀-ͯ]/g, '').replace(/ł/g, 'l');

const results = computed(() => {
    const n = norm(q.value.trim());
    return props.recipes
        .filter((r) => r.id !== props.currentId && norm(r.name).includes(n))
        .sort((a, b) => Number(a.status === 'draft') - Number(b.status === 'draft'));
});

watch(q, () => (hi.value = 0));

function choose(r: RecipeBaseDto | undefined) {
    if (!r) return;
    emit('select', { id: r.id, title: r.name, status: r.status });
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
        choose(results.value[hi.value]);
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
        class="lpick"
    >
        <label
            :for="`${id}-q`"
            class="lpick__label"
        >
            Przepis dla „{{ ingredientName || 'składnika' }}”
        </label>
        <div class="lpick__search">
            <svg
                width="14"
                height="14"
                viewBox="0 0 14 14"
                fill="none"
                aria-hidden="true"
            >
                <circle
                    cx="6"
                    cy="6"
                    r="4.2"
                    stroke="currentColor"
                    stroke-width="1.4"
                />
                <path
                    d="M9.2 9.2L12.5 12.5"
                    stroke="currentColor"
                    stroke-width="1.4"
                    stroke-linecap="round"
                />
            </svg>
            <input
                :id="`${id}-q`"
                ref="input"
                v-model="q"
                class="input lpick__input"
                role="combobox"
                aria-autocomplete="list"
                aria-expanded="true"
                :aria-controls="listId"
                :aria-activedescendant="results.length ? `${listId}-${hi}` : undefined"
                placeholder="Szukaj po nazwie przepisu"
                autocomplete="off"
                @keydown="onKey"
            />
        </div>
        <ul
            :id="listId"
            role="listbox"
            aria-label="Przepisy"
            class="lpick__list"
        >
            <li
                v-for="(r, i) in results"
                :id="`${listId}-${i}`"
                :key="r.id"
                role="option"
                :aria-selected="i === hi"
                class="lpick__opt"
                :class="{ 'is-hi': i === hi }"
                @mouseenter="hi = i"
                @mousedown.prevent="choose(r)"
            >
                <span
                    class="lpick__dot"
                    :class="mainDish(r)?.color ? `cat-${mainDish(r)!.color}` : null"
                    aria-hidden="true"
                />
                <span class="lpick__name">{{ r.name }}</span>
                <span
                    v-if="r.status === 'draft'"
                    class="lpick__draft"
                >
                    szkic
                </span>
                <span
                    v-else
                    class="lpick__dish"
                >
                    {{ mainDish(r)?.label }}
                </span>
            </li>
            <li
                v-if="!results.length"
                class="lpick__empty"
            >
                Brak przepisu „{{ q }}”.
            </li>
        </ul>
        <div class="lpick__hint">
            ↑ ↓ wybór · Enter zatwierdza · Esc zamyka. Bieżącego przepisu nie ma na liście.
        </div>
    </div>
</template>

<style scoped>
.lpick {
    position: absolute;
    right: 38px;
    top: calc(100% - 6px);
    z-index: var(--z-dropdown);
    width: 400px;
    padding: 10px;
    display: flex;
    flex-direction: column;
    gap: 8px;
    background: var(--bg-alt);
    border: 1px solid var(--rule-strong);
    border-radius: var(--r-xl);
    box-shadow: var(--shadow-menu);
}
.lpick__label {
    font-family: var(--font-mono);
    font-size: 12px;
    font-weight: 600;
    letter-spacing: 0.5px;
    text-transform: uppercase;
    color: var(--ink-muted);
}
.lpick__search {
    position: relative;
}
.lpick__search svg {
    position: absolute;
    left: 10px;
    top: 50%;
    transform: translateY(-50%);
    color: var(--ink-muted);
}
.lpick__input {
    width: 100%;
    padding-left: 32px;
}
.lpick__list {
    list-style: none;
    margin: 0;
    padding: 0;
    max-height: 220px;
    overflow: auto;
    scrollbar-gutter: stable;
}
.lpick__opt {
    display: grid;
    grid-template-columns: 8px minmax(0, 1fr) auto;
    gap: 10px;
    align-items: center;
    padding: 8px 10px;
    border-radius: var(--r-md);
    cursor: pointer;
}
.lpick__opt.is-hi {
    background: var(--accent-soft);
}
.lpick__dot {
    width: 7px;
    height: 7px;
    border-radius: var(--r-pill);
    background: var(--cat, var(--ink-muted));
}
.lpick__name {
    font-size: 14px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
}
.lpick__dish {
    font-size: 12px;
    color: var(--ink-muted);
}
.lpick__draft {
    padding: 1px 8px;
    border-radius: var(--r-pill);
    font-size: 12px;
    font-weight: 600;
    background: var(--prep-soft);
    color: var(--prep);
}
.lpick__empty {
    padding: 10px;
    font-size: 13px;
    color: var(--ink-muted);
}
.lpick__hint {
    padding-top: 8px;
    border-top: 1px solid var(--rule);
    font-size: 12px;
    color: var(--ink-muted);
}
</style>
