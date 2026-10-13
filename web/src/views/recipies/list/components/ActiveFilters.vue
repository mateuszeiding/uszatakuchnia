<script lang="ts" setup>
import { computed } from 'vue';

import { DIFFICULTY, emptyFilters, type IListFilters, TIMES } from '../listFilters';

import type { ITaxonomyAxis } from '@/data/dtos/taxonomy/DictionariesDto';

// Pasek aktywnych filtrów: każdy chip usuwa jeden filtr, „Wyczyść wszystko” — wszystkie.
const props = defineProps<{ axes: ITaxonomyAxis[] }>();
const f = defineModel<IListFilters>({ required: true });

const chips = computed(() => {
    const out: { label: string; remove: () => void }[] = [];
    const v = f.value;
    if (v.q.trim())
        out.push({ label: `„${v.q.trim()}”`, remove: () => (f.value = { ...v, q: '' }) });
    if (v.noPrep)
        out.push({ label: 'bez przygotowania', remove: () => (f.value = { ...v, noPrep: false }) });
    if (v.maxTime)
        out.push({
            label: TIMES.find((t) => t[0] === v.maxTime)?.[1] ?? `do ${v.maxTime} min`,
            remove: () => (f.value = { ...v, maxTime: 0 }),
        });
    if (v.difficulty)
        out.push({
            label: DIFFICULTY[v.difficulty] ?? '',
            remove: () => (f.value = { ...v, difficulty: 0 }),
        });
    for (const a of props.axes) {
        for (const id of v.axes[a.code] ?? []) {
            const label = a.values.find((x) => x.id === id)?.label;
            if (!label) continue;
            out.push({
                label,
                remove: () =>
                    (f.value = {
                        ...v,
                        axes: {
                            ...v.axes,
                            [a.code]: (v.axes[a.code] ?? []).filter((x) => x !== id),
                        },
                    }),
            });
        }
    }
    return out;
});

const clear = () => (f.value = emptyFilters(f.value.sort));
</script>

<template>
    <div
        v-if="chips.length"
        class="active"
    >
        <span
            v-for="c in chips"
            :key="c.label"
            class="active__chip"
        >
            {{ c.label }}
            <button
                type="button"
                class="active__x"
                :aria-label="`Usuń filtr: ${c.label}`"
                @click="c.remove"
            >
                <svg
                    width="11"
                    height="11"
                    viewBox="0 0 11 11"
                    aria-hidden="true"
                >
                    <path
                        d="M2 2l7 7M9 2L2 9"
                        stroke="currentColor"
                        stroke-width="1.5"
                        stroke-linecap="round"
                    />
                </svg>
            </button>
        </span>
        <button
            type="button"
            class="active__clear"
            @click="clear"
        >
            Wyczyść wszystko
        </button>
    </div>
    <div
        v-else
        class="active__rule"
    />
</template>

<style scoped>
.active {
    display: flex;
    align-items: center;
    gap: 6px;
    flex-wrap: wrap;
    padding-bottom: 14px;
    margin-bottom: 14px;
    border-bottom: 1px solid var(--rule);
}
.active__rule {
    height: 1px;
    background: var(--rule);
    margin-bottom: 14px;
}
.active__chip {
    display: inline-flex;
    align-items: center;
    gap: 2px;
    background: var(--accent-soft);
    color: var(--accent);
    border-radius: var(--r-pill);
    padding: 2px 2px 2px 12px;
    font-size: 13px;
    font-weight: 500;
    white-space: nowrap;
    flex-shrink: 0;
}
.active__x {
    width: 28px;
    height: 28px;
    border: none;
    background: transparent;
    color: var(--accent);
    cursor: pointer;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    border-radius: var(--r-pill);
}
.active__clear {
    margin-left: auto;
    background: none;
    border: none;
    color: var(--accent);
    font-size: 14px;
    font-weight: 600;
    cursor: pointer;
    padding: 6px 4px;
    white-space: nowrap;
    flex-shrink: 0;
}

@media (max-width: 900px) {
    .active {
        flex-wrap: nowrap;
        overflow-x: auto;
        border-bottom: none;
        margin-bottom: 4px;
    }
    .active__rule {
        display: none;
    }
    .active__clear {
        margin-left: 4px;
    }
}
</style>
