<script lang="ts" setup>
import { computed, nextTick, onMounted, ref, useId, watch } from 'vue';

import { fold } from '../stepRefsForm';

import type { ITerm } from '@/shared/utils/stepRefs';

// Marking a selected word as a glossary term: search (pre-filled with the word), the list, and "+ Nowe pojęcie"
// at the end. ↑ ↓ Enter Esc; the search field keeps focus.
const props = defineProps<{ terms: ITerm[]; word: string; current?: number | null }>();
const emit = defineEmits<{ pick: [id: number]; create: [name: string]; close: [] }>();

const uid = useId();
const q = ref(props.word);
const input = ref<HTMLInputElement>();
const listEl = ref<HTMLElement>();

const matches = computed(() => {
    const f = fold(q.value.trim());
    const all = [...props.terms].sort((a, b) => a.name.localeCompare(b.name, 'pl'));
    if (!f) return all;
    // the word is usually inflected ("deglasuj" → "deglasowanie"): also match on its first letters
    const stem = f.slice(0, Math.max(4, Math.ceil(f.length * 0.6)));
    const hit = (t: ITerm) => fold(t.name).includes(f) || fold(t.name).startsWith(stem);
    return all.filter(hit);
});
/** Options: the terms, then "+ Nowe pojęcie". */
const count = computed(() => matches.value.length + 1);
const active = ref(0);
watch(matches, () => (active.value = 0));

async function scrollActive() {
    await nextTick();
    listEl.value?.querySelector('[aria-selected="true"]')?.scrollIntoView({ block: 'nearest' });
}

function choose(i: number) {
    const t = matches.value[i];
    if (t) emit('pick', t.id);
    else emit('create', q.value.trim() || props.word);
}

function onKey(e: KeyboardEvent) {
    if (e.key === 'ArrowDown' || e.key === 'ArrowUp') {
        e.preventDefault();
        active.value = (active.value + (e.key === 'ArrowDown' ? 1 : count.value - 1)) % count.value;
        void scrollActive();
    } else if (e.key === 'Enter') {
        e.preventDefault();
        choose(active.value);
    } else if (e.key === 'Escape') {
        e.preventDefault();
        e.stopPropagation();
        emit('close');
    }
}

onMounted(() => {
    input.value?.focus();
    input.value?.select();
});
</script>

<template>
    <div
        class="term-pick"
        role="dialog"
        aria-label="Oznacz pojęcie"
    >
        <label
            :for="`${uid}-q`"
            class="sr-only"
        >
            Szukaj pojęcia
        </label>
        <input
            :id="`${uid}-q`"
            ref="input"
            v-model="q"
            class="input term-pick__q"
            role="combobox"
            aria-autocomplete="list"
            aria-expanded="true"
            :aria-controls="`${uid}-list`"
            :aria-activedescendant="`${uid}-o${active}`"
            autocomplete="off"
            placeholder="Szukaj pojęcia…"
            @keydown="onKey"
        />
        <ul
            :id="`${uid}-list`"
            ref="listEl"
            class="term-pick__list"
            role="listbox"
            aria-label="Pojęcia"
        >
            <li
                v-for="(t, i) in matches"
                :id="`${uid}-o${i}`"
                :key="t.id"
                role="option"
                class="term-pick__opt"
                :class="{ 'is-active': i === active }"
                :aria-selected="i === active"
                @mousedown.prevent
                @click="choose(i)"
                @mousemove="active = i"
            >
                <span class="term-pick__name">
                    {{ t.name }}
                    <span
                        v-if="t.id === current"
                        class="term-pick__current"
                    >
                        · teraz
                    </span>
                </span>
                <span class="term-pick__text">{{ t.explanation }}</span>
            </li>
            <li
                :id="`${uid}-o${matches.length}`"
                role="option"
                class="term-pick__opt term-pick__new"
                :class="{ 'is-active': active === matches.length }"
                :aria-selected="active === matches.length"
                @mousedown.prevent
                @click="choose(matches.length)"
                @mousemove="active = matches.length"
            >
                + Nowe pojęcie „{{ q.trim() || word }}”…
            </li>
        </ul>
        <p class="term-pick__foot">
            Słowo w kroku zostaje takie, jak je napisałeś. ↑ ↓ · Enter · Esc
        </p>
    </div>
</template>

<style scoped>
.term-pick {
    width: 360px;
    max-width: calc(100vw - 32px);
    background: var(--bg-alt);
    border: 1px solid var(--rule-strong);
    border-radius: var(--r-xl);
    box-shadow: var(--shadow-menu);
    padding: 10px;
}
.term-pick__q {
    width: 100%;
    font-size: 14px;
}
.term-pick__list {
    list-style: none;
    margin: 8px 0 0;
    padding: 0;
    max-height: 260px;
    overflow-y: auto;
}
.term-pick__opt {
    display: flex;
    flex-direction: column;
    gap: 2px;
    padding: 8px 10px;
    border-radius: var(--r-md);
    cursor: pointer;
}
.term-pick__opt.is-active {
    background: var(--accent-soft);
}
.term-pick__name {
    font-size: 14px;
    font-weight: 600;
}
.term-pick__current {
    font-weight: 500;
    color: var(--ink-muted);
}
.term-pick__text {
    font-size: 12.5px;
    color: var(--ink-muted);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}
.term-pick__new {
    margin-top: 6px;
    border-top: 1px solid var(--rule);
    border-radius: 0 0 var(--r-md) var(--r-md);
    font-size: 14px;
    font-weight: 600;
    color: var(--accent);
}
.term-pick__foot {
    margin: 6px 0 0;
    padding-top: 8px;
    border-top: 1px solid var(--rule);
    font-size: 12px;
    color: var(--ink-muted);
}
</style>
