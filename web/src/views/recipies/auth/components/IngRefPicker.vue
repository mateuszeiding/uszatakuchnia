<script lang="ts" setup>
import { computed, nextTick, ref, watch } from 'vue';

import { type IRefRow } from '../stepRefsForm';

import SwapIcon from '@/components/SwapIcon.vue';
import { PARTS, type RefPart } from '@/shared/utils/stepRefs';

// The list behind "@": rows from the Składniki section grouped like there, then the part of the row.
// The step field keeps the focus (the caret stays after "@…"), so the keyboard is handled by the editor;
// here mousedown is prevented for the same reason.
const props = defineProps<{
    listId: string;
    /** Rows already filtered by the typed text, in form order. */
    rows: IRefRow[];
    query: string;
    active: number;
    part: RefPart | null;
}>();
const emit = defineEmits<{
    pick: [key: string];
    'update:active': [i: number];
    'update:part': [p: RefPart | null];
}>();

const groups = computed(() => {
    const out: { label: string | null; items: { row: IRefRow; i: number }[] }[] = [];
    props.rows.forEach((row, i) => {
        const last = out[out.length - 1];
        if (last && last.label === row.group) last.items.push({ row, i });
        else out.push({ label: row.group, items: [{ row, i }] });
    });
    return out;
});
const activeRow = computed(() => props.rows[props.active]);
const optId = (i: number) => `${props.listId}-o${i}`;

// focus stays in the step field, so the browser doesn't scroll the list to the active option by itself
const list = ref<HTMLElement>();
watch(
    () => [props.active, props.rows],
    async () => {
        await nextTick();
        list.value
            ?.querySelector<HTMLElement>(`#${CSS.escape(optId(props.active))}`)
            ?.scrollIntoView({ block: 'nearest' });
    }
);
defineExpose({ optId });
</script>

<template>
    <div
        class="ref-pop"
        role="dialog"
        aria-label="Wstaw ilość składnika"
        @mousedown.prevent
    >
        <div class="ref-pop__label">Składnik{{ query ? ` · „${query}”` : '' }}</div>
        <ul
            v-if="rows.length"
            :id="listId"
            ref="list"
            class="ref-pop__list"
            role="listbox"
            aria-label="Składniki"
        >
            <template
                v-for="(g, gi) in groups"
                :key="gi"
            >
                <li
                    v-if="g.label"
                    class="ref-pop__group"
                    role="presentation"
                >
                    {{ g.label }}
                </li>
                <li
                    v-for="o in g.items"
                    :id="optId(o.i)"
                    :key="o.row.key"
                    class="ref-pop__opt"
                    :class="{ 'is-active': o.i === active }"
                    role="option"
                    :aria-selected="o.i === active"
                    @mouseenter="emit('update:active', o.i)"
                    @click="emit('pick', o.row.key)"
                >
                    <span class="ref-pop__name">{{ o.row.name }}</span>
                    <span class="ref-pop__qty">{{ o.row.qty || 'bez ilości' }}</span>
                </li>
            </template>
        </ul>
        <p
            v-else
            class="ref-pop__empty"
        >
            Brak składnika „{{ query }}” w sekcji Składniki.
        </p>
        <p
            v-if="activeRow?.hasSubs"
            class="ref-pop__subs"
        >
            <span class="ref-pop__swap"><SwapIcon :size="13" /></span>
            Ten składnik ma zamienniki. Słowo przed odwołaniem się nie zmieni — napisz ogólnie, np.
            „mięso”, a nazwę zamiennika pokaże znacznik.
        </p>
        <div class="ref-pop__parts">
            <div class="ref-pop__sub">Ile z tego wiersza</div>
            <div
                class="ref-pop__radios"
                role="radiogroup"
                aria-label="Część"
            >
                <button
                    v-for="p in PARTS"
                    :key="p.label"
                    type="button"
                    role="radio"
                    class="ref-pop__radio"
                    :aria-checked="p.part === part"
                    tabindex="-1"
                    @click="emit('update:part', p.part)"
                >
                    {{ p.label }}
                </button>
            </div>
        </div>
        <div class="ref-pop__keys">↑ ↓ wybór · ← → część · Enter wstawia · Esc zamyka</div>
    </div>
</template>

<style scoped>
.ref-pop {
    position: absolute;
    z-index: 20;
    width: 380px;
    max-width: calc(100vw - 32px);
    background: var(--bg-alt);
    border: 1px solid var(--rule-strong);
    border-radius: var(--r-lg);
    box-shadow: var(--shadow-menu);
    padding: 10px;
}
.ref-pop__label,
.ref-pop__group {
    font-family: var(--font-mono);
    font-size: 11px;
    font-weight: 600;
    letter-spacing: 0.5px;
    text-transform: uppercase;
    color: var(--ink-muted);
    margin: 2px 4px 8px;
}
.ref-pop__group {
    margin: 0;
    padding: 8px 8px 4px;
    color: var(--ink-faint);
}
.ref-pop__list {
    list-style: none;
    margin: 0;
    padding: 0;
    max-height: 260px;
    overflow-y: auto;
}
.ref-pop__opt {
    display: flex;
    justify-content: space-between;
    gap: 8px;
    padding: 7px 8px;
    border-radius: var(--r-md);
    font-size: 14px;
    cursor: pointer;
}
.ref-pop__opt.is-active {
    background: var(--accent-soft);
}
.ref-pop__qty {
    font-family: var(--font-mono);
    font-size: 12.5px;
    color: var(--ink-muted);
    white-space: nowrap;
}
.ref-pop__empty {
    font-size: 13px;
    color: var(--ink-muted);
    margin: 4px 8px 8px;
}
.ref-pop__subs {
    display: flex;
    gap: 6px;
    margin: 8px 4px 0;
    font-size: 13px;
    line-height: 1.5;
    color: var(--ink-dim);
}
.ref-pop__swap {
    display: inline-flex;
    color: var(--accent);
    padding-top: 2px;
}
.ref-pop__parts {
    border-top: 1px solid var(--rule);
    margin-top: 8px;
    padding-top: 10px;
}
.ref-pop__sub {
    font-size: 12px;
    color: var(--ink-muted);
    margin: 0 4px 6px;
}
.ref-pop__radios {
    display: flex;
    gap: 4px;
    flex-wrap: wrap;
}
.ref-pop__radio {
    padding: 5px 10px;
    border-radius: var(--r-pill);
    font-size: 13px;
    font-weight: 600;
    border: 1px solid var(--rule-strong);
    background: var(--bg-alt);
    color: var(--ink);
    cursor: pointer;
}
.ref-pop__radio[aria-checked='true'] {
    border-color: var(--accent);
    background: var(--accent-soft);
    color: var(--accent);
}
.ref-pop__keys {
    font-size: 12px;
    color: var(--ink-muted);
    border-top: 1px solid var(--rule);
    margin-top: 10px;
    padding-top: 8px;
}
</style>
