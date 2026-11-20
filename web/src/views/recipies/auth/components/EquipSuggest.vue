<script lang="ts" setup>
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue';

import type { IEquipment } from '@/data/dtos/taxonomy/DictionariesDto';

// Equipment combobox: suggestions from the dictionary, ↑ ↓ to pick, Enter confirms, Esc closes.
// The last option is "+ Dodaj „…” do słownika" (when there's no exact match).
const props = defineProps<{
    id: string;
    name: string;
    equipment: IEquipment[];
    invalid?: boolean;
}>();
const emit = defineEmits<{
    'update:name': [v: string];
    select: [e: IEquipment];
    'add-new': [name: string];
}>();

const open = ref(false);
const hi = ref(0);
const root = ref<HTMLElement>();
const listId = computed(() => `${props.id}-list`);

type Opt = { kind: 'match'; e: IEquipment } | { kind: 'add' };

const options = computed<Opt[]>(() => {
    const q = props.name.trim().toLowerCase();
    if (!q) return [];
    const matches = props.equipment
        .filter((e) => e.name.toLowerCase().includes(q))
        .sort(
            (a, b) =>
                Number(!a.name.toLowerCase().startsWith(q)) -
                    Number(!b.name.toLowerCase().startsWith(q)) ||
                a.name.localeCompare(b.name, 'pl')
        )
        .slice(0, 8);
    const exact = matches.some((m) => m.name.toLowerCase() === q);
    return [
        ...matches.map((e): Opt => ({ kind: 'match', e })),
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
        emit('update:name', o.e.name);
        emit('select', o.e);
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
        class="eq-suggest"
    >
        <input
            :id="id"
            :value="name"
            class="input eq-suggest__input"
            :class="{ 'is-invalid': invalid }"
            role="combobox"
            aria-label="Sprzęt"
            aria-autocomplete="list"
            :aria-expanded="expanded"
            :aria-controls="listId"
            :aria-activedescendant="expanded ? `${listId}-${hi}` : undefined"
            :aria-invalid="invalid || undefined"
            placeholder="np. robot planetarny"
            autocomplete="off"
            @input="onInput"
            @focus="name.trim() && (open = true)"
            @keydown="onKey"
        />
        <ul
            v-if="expanded"
            :id="listId"
            class="eq-suggest__list"
            role="listbox"
        >
            <li
                v-for="(o, i) in options"
                :id="`${listId}-${i}`"
                :key="o.kind === 'match' ? o.e.id : 'add'"
                role="option"
                class="eq-suggest__opt"
                :class="{ 'is-hi': i === hi, 'is-add': o.kind === 'add' }"
                :aria-selected="i === hi"
                @mousedown.prevent="choose(o)"
                @mouseenter="hi = i"
            >
                <template v-if="o.kind === 'add'">+ Dodaj „{{ name.trim() }}” do słownika</template>
                <template v-else>{{ o.e.name }}</template>
            </li>
        </ul>
    </div>
</template>

<style scoped>
.eq-suggest {
    position: relative;
    min-width: 0;
}
.eq-suggest__input {
    padding: 8px 10px;
}
.eq-suggest__list {
    position: absolute;
    top: calc(100% + 4px);
    left: 0;
    z-index: var(--z-dropdown);
    width: max(100%, 300px);
    margin: 0;
    padding: 6px;
    list-style: none;
    background: var(--bg-alt);
    border: 1px solid var(--rule-strong);
    border-radius: var(--r-lg);
    box-shadow: var(--shadow-menu);
}
.eq-suggest__opt {
    padding: 8px 10px;
    border-radius: var(--r-md);
    font-size: 14px;
    cursor: pointer;
}
.eq-suggest__opt.is-hi {
    background: var(--accent-soft);
}
.eq-suggest__opt.is-add {
    color: var(--accent);
    font-weight: 600;
}
.eq-suggest__opt.is-add:not(:first-child) {
    margin-top: 4px;
    border-top: 1px solid var(--rule);
    border-radius: 0 0 var(--r-md) var(--r-md);
}
</style>
