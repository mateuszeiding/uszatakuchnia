<script lang="ts" setup>
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue';

import type { IngredientDto } from '@/data/dtos/ingredients/IngredientDto';

const props = defineProps<{
    modelValue: number; // ingredientId
    ingredients: IngredientDto[];
}>();

const emit = defineEmits<{
    'update:modelValue': [id: number];
}>();

const query = ref('');
const open = ref(false);
const root = ref<HTMLElement | null>(null);

// Sync display name when modelValue changes externally (e.g. on load)
watch(
    () => props.modelValue,
    (id) => {
        if (id > 0) {
            const found = props.ingredients.find((i) => i.id === id);
            if (found) query.value = found.name;
        } else {
            query.value = '';
        }
    },
    { immediate: true }
);

const filtered = computed(() => {
    const q = query.value.trim().toLowerCase();
    if (!q) return [];
    return props.ingredients.filter((i) => i.name.toLowerCase().includes(q));
});

const grouped = computed(() => {
    const map = new Map<string, IngredientDto[]>();
    for (const ing of filtered.value) {
        const key = ing.type?.name || 'inne';
        if (!map.has(key)) map.set(key, []);
        map.get(key)!.push(ing);
    }
    return [...map.entries()].map(([type, items]) => ({ type, items }));
});

function onInput(e: Event) {
    query.value = (e.target as HTMLInputElement).value;
    open.value = true;
    // Clear selection if user edits
    emit('update:modelValue', 0);
}

function select(ing: IngredientDto) {
    query.value = ing.name;
    emit('update:modelValue', ing.id);
    open.value = false;
}

function onFocus() {
    if (query.value.trim()) open.value = true;
}

function onClickOutside(e: MouseEvent) {
    if (root.value && !root.value.contains(e.target as Node)) {
        open.value = false;
    }
}

onMounted(() => document.addEventListener('mousedown', onClickOutside));
onBeforeUnmount(() => document.removeEventListener('mousedown', onClickOutside));
</script>

<template>
    <div
        ref="root"
        class="ing-suggest"
    >
        <input
            :value="query"
            class="input ing-suggest__input"
            type="text"
            placeholder="np. masło"
            autocomplete="off"
            @input="onInput"
            @focus="onFocus"
        />

        <div
            v-if="open && (grouped.length > 0 || query.trim().length > 0)"
            class="ing-suggest__dropdown"
        >
            <template v-if="grouped.length > 0">
                <div
                    v-for="group in grouped"
                    :key="group.type"
                >
                    <div class="ing-suggest__group-label">{{ group.type }}</div>
                    <button
                        v-for="ing in group.items"
                        :key="ing.id"
                        class="ing-suggest__item"
                        :class="{ 'ing-suggest__item--active': modelValue === ing.id }"
                        type="button"
                        @mousedown.prevent="select(ing)"
                    >
                        {{ ing.name }}
                    </button>
                </div>
            </template>

            <div
                v-if="query.trim().length > 0"
                class="ing-suggest__add"
            >
                + Nie ma w bazie? Dodaj „{{ query }}"
            </div>
        </div>
    </div>
</template>

<style scoped>
.ing-suggest {
    position: relative;
    flex: 1;
}

.ing-suggest__input {
    width: 100%;
    padding: 8px 10px;
    font-size: 13px;
}

.ing-suggest__dropdown {
    position: absolute;
    top: calc(100% + 3px);
    left: 0;
    right: 0;
    z-index: 200;
    background: var(--bg-alt);
    border: 1px solid var(--rule-strong);
    border-radius: var(--r-lg);
    box-shadow: 0 6px 20px rgba(10, 10, 15, 0.12);
    max-height: 260px;
    overflow-y: auto;
}

.ing-suggest__group-label {
    padding: 6px 12px 4px;
    font-family: var(--font-mono);
    font-size: 9px;
    letter-spacing: 0.7px;
    text-transform: uppercase;
    font-weight: 600;
    color: var(--ink-faint);
    position: sticky;
    top: 0;
    background: var(--bg-alt);
    border-bottom: 1px solid var(--rule);
}

.ing-suggest__item {
    display: block;
    width: 100%;
    padding: 8px 14px;
    font-size: 13px;
    text-align: left;
    cursor: pointer;
    color: var(--ink);
    background: transparent;
    border: none;
    font-family: var(--font-sans);
    font-weight: 400;
}
.ing-suggest__item:hover,
.ing-suggest__item--active {
    background: var(--accent-soft);
}
.ing-suggest__item--active {
    font-weight: 600;
}

.ing-suggest__add {
    padding: 8px 14px;
    font-size: 12px;
    color: var(--accent);
    font-family: var(--font-mono);
    font-weight: 600;
    border-top: 1px solid var(--rule);
    cursor: default;
}
</style>
