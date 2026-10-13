<script lang="ts" setup>
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue';

import type { IngredientDto } from '@/data/dtos/ingredients/IngredientDto';

// Combobox składnika: ↑ ↓ wybór, Enter zatwierdza, Esc zamyka. Ostatnia opcja „Dodaj do bazy”.
const props = defineProps<{
    id: string;
    name: string;
    ingredients: IngredientDto[];
    invalid?: boolean;
}>();

const emit = defineEmits<{
    'update:name': [v: string];
    select: [ing: IngredientDto];
    'add-new': [name: string];
}>();

const open = ref(false);
const hi = ref(0);
const root = ref<HTMLElement>();
const listId = computed(() => `${props.id}-list`);

type Opt = { kind: 'match'; ing: IngredientDto } | { kind: 'add' };

const options = computed<Opt[]>(() => {
    const q = props.name.trim().toLowerCase();
    if (!q) return [];
    const matches = props.ingredients
        .filter((i) => i.name.toLowerCase().includes(q))
        .sort(
            (a, b) =>
                Number(!a.name.toLowerCase().startsWith(q)) -
                Number(!b.name.toLowerCase().startsWith(q))
        )
        .slice(0, 8);
    const exact = matches.some((m) => m.name.toLowerCase() === q);
    return [
        ...matches.map((ing): Opt => ({ kind: 'match', ing })),
        ...(exact ? [] : [{ kind: 'add' } as Opt]),
    ];
});

const expanded = computed(() => open.value && options.value.length > 0);

watch(
    () => props.name,
    () => (hi.value = 0)
);

function choose(o: Opt | undefined) {
    if (!o) return;
    if (o.kind === 'add') emit('add-new', props.name.trim());
    else {
        emit('update:name', o.ing.name);
        emit('select', o.ing);
    }
    open.value = false;
}

function onKey(e: KeyboardEvent) {
    const n = options.value.length;
    if (!expanded.value) {
        if (e.key === 'ArrowDown' && n) {
            open.value = true;
            e.preventDefault();
        }
        return;
    }
    if (e.key === 'ArrowDown') {
        e.preventDefault();
        hi.value = (hi.value + 1) % n;
    } else if (e.key === 'ArrowUp') {
        e.preventDefault();
        hi.value = (hi.value - 1 + n) % n;
    } else if (e.key === 'Enter') {
        e.preventDefault();
        choose(options.value[hi.value]);
    } else if (e.key === 'Escape') {
        e.stopPropagation();
        open.value = false;
    }
}

function onInput(e: Event) {
    emit('update:name', (e.target as HTMLInputElement).value);
    open.value = true;
}

const onDocDown = (e: MouseEvent) => {
    if (root.value && !root.value.contains(e.target as Node)) open.value = false;
};
onMounted(() => document.addEventListener('mousedown', onDocDown));
onBeforeUnmount(() => document.removeEventListener('mousedown', onDocDown));
</script>

<template>
    <div
        ref="root"
        class="ing-suggest"
    >
        <input
            :id="id"
            :value="name"
            class="input ing-suggest__input"
            :class="{ 'is-invalid': invalid }"
            role="combobox"
            aria-label="Składnik"
            aria-autocomplete="list"
            :aria-expanded="expanded"
            :aria-controls="listId"
            :aria-activedescendant="expanded ? `${listId}-${hi}` : undefined"
            :aria-invalid="invalid || undefined"
            placeholder="np. masło"
            autocomplete="off"
            @input="onInput"
            @focus="name.trim() && (open = true)"
            @keydown="onKey"
        />
        <ul
            v-if="expanded"
            :id="listId"
            class="ing-suggest__list"
            role="listbox"
        >
            <li
                v-for="(o, i) in options"
                :id="`${listId}-${i}`"
                :key="o.kind === 'match' ? o.ing.id : 'add'"
                role="option"
                class="ing-suggest__opt"
                :class="{ 'is-hi': i === hi, 'is-add': o.kind === 'add' }"
                :aria-selected="i === hi"
                @mousedown.prevent="choose(o)"
                @mouseenter="hi = i"
            >
                <template v-if="o.kind === 'add'">+ Dodaj „{{ name.trim() }}” do bazy</template>
                <template v-else>
                    <span>{{ o.ing.name }}</span>
                    <span class="ing-suggest__type">{{ o.ing.type?.name }}</span>
                </template>
            </li>
        </ul>
    </div>
</template>

<style scoped>
.ing-suggest {
    position: relative;
    min-width: 0;
}
.ing-suggest__input {
    padding: 8px 10px;
}
.ing-suggest__list {
    position: absolute;
    top: 100%;
    left: 0;
    right: 0;
    z-index: var(--z-dropdown);
    margin: 4px 0 0;
    padding: 0;
    list-style: none;
    background: var(--bg-alt);
    border: 1px solid var(--rule-strong);
    border-radius: var(--r-lg);
    box-shadow: var(--shadow-menu);
    overflow: hidden;
}
.ing-suggest__opt {
    display: flex;
    justify-content: space-between;
    gap: 10px;
    padding: 8px 12px;
    font-size: 13px;
    cursor: pointer;
}
.ing-suggest__opt.is-hi {
    background: var(--accent-soft);
}
.ing-suggest__opt.is-add {
    color: var(--accent);
    font-weight: 600;
}
.ing-suggest__opt.is-add:not(:first-child) {
    border-top: 1px solid var(--rule);
}
.ing-suggest__type {
    font-size: 12px;
    color: var(--ink-muted);
}
</style>
