<script lang="ts" setup>
import { computed } from 'vue';

import {
    DIFFICULTY,
    facetCount,
    type FacetDim,
    facetPresent,
    type IListFilters,
    TIMES,
} from '../listFilters';

import type { RecipeBaseDto } from '@/data/dtos/recipe/RecipeDto';
import type { ITaxonomyAxis } from '@/data/dtos/taxonomy/DictionariesDto';

// Facets: only values that at least one published recipe has; a group without values
// disappears; each value has a counter of how many recipes there would be after adding it. 0 and not selected → disabled,
// dimmed, in its place. "Bez przygotowania wcześniej" only when some recipe requires prep.
const props = defineProps<{
    axes: ITaxonomyAxis[];
    /** Published recipes (the whole list, unfiltered): the source of facets and counters. */
    recipes: RecipeBaseDto[];
    loading?: boolean;
}>();

/** Skeleton chip rows: dish type, time, difficulty, diets, two axes (facets are usually 1–2 rows). */
const SK_ROWS = [2, 1, 1, 1, 1, 1];
const f = defineModel<IListFilters>({ required: true });

interface IOption {
    id: number;
    label: string;
    color?: string;
    count: number;
    on: boolean;
}
interface IGroup {
    key: string;
    label: string;
    dim: FacetDim;
    hint?: string;
    kind: 'color' | 'strong' | 'plain';
    options: IOption[];
}

function isOn(dim: FacetDim, id: number) {
    if (dim.kind === 'axis') return (f.value.axes[dim.code] ?? []).includes(id);
    if (dim.kind === 'time') return f.value.maxTime === id;
    return f.value.difficulty === id;
}

function options(dim: FacetDim, values: { id: number; label: string; color?: string }[]) {
    const out: IOption[] = [];
    for (const v of values) {
        const on = isOn(dim, v.id);
        // a value selected from a link stays visible even without recipes, so it can be unselected
        if (!on && !facetPresent(props.recipes, dim, v.id)) continue;
        out.push({ ...v, on, count: facetCount(props.recipes, f.value, props.axes, dim, v.id) });
    }
    return out;
}

function axisGroup(a: ITaxonomyAxis, kind: IGroup['kind'], hint?: string): IGroup {
    const dim: FacetDim = { kind: 'axis', code: a.code };
    return { key: a.code, label: a.label, dim, kind, hint, options: options(dim, a.values) };
}

// Order: dish type, time, difficulty, diets, then the remaining axes.
const groups = computed<IGroup[]>(() => {
    const dish = props.axes.find((a) => a.code === 'dish');
    const diet = props.axes.find((a) => a.code === 'diet');
    const g: IGroup[] = [];
    if (dish) g.push(axisGroup(dish, 'color'));
    const time: FacetDim = { kind: 'time' };
    g.push({
        key: 'time',
        label: 'Czas',
        dim: time,
        kind: 'plain',
        options: options(
            time,
            TIMES.map(([id, label]) => ({ id, label }))
        ),
    });
    const diff: FacetDim = { kind: 'difficulty' };
    g.push({
        key: 'difficulty',
        label: 'Trudność',
        dim: diff,
        kind: 'plain',
        options: options(
            diff,
            [1, 2, 3].map((id) => ({ id, label: DIFFICULTY[id]! }))
        ),
    });
    if (diet) g.push(axisGroup(diet, 'strong', 'Przepis spełnia wszystkie zaznaczone.'));
    for (const a of props.axes)
        if (a.code !== 'dish' && a.code !== 'diet') g.push(axisGroup(a, 'plain'));
    return g.filter((x) => x.options.length);
});

const showPrep = computed(() => f.value.noPrep || props.recipes.some((r) => r.needsPrep));

function toggle(dim: FacetDim, id: number) {
    if (dim.kind === 'time') {
        f.value = { ...f.value, maxTime: f.value.maxTime === id ? 0 : id };
        return;
    }
    if (dim.kind === 'difficulty') {
        f.value = { ...f.value, difficulty: f.value.difficulty === id ? 0 : id };
        return;
    }
    const cur = f.value.axes[dim.code] ?? [];
    f.value = {
        ...f.value,
        axes: {
            ...f.value.axes,
            [dim.code]: cur.includes(id) ? cur.filter((x) => x !== id) : [...cur, id],
        },
    };
}
</script>

<template>
    <div
        class="filters"
        :inert="loading"
    >
        <template v-if="loading">
            <div
                v-for="(n, i) in SK_ROWS"
                :key="i"
                class="filters__group"
                aria-hidden="true"
            >
                <div class="filters__legend"><span class="sk fsk__legend" /></div>
                <div class="filters__chips">
                    <div
                        v-for="r in n"
                        :key="r"
                        class="fsk__row"
                    >
                        <span
                            v-for="c in r % 2 ? 3 : 2"
                            :key="c"
                            class="chip fsk__chip"
                        >
                            <span class="sk fsk__fill" />
                            &nbsp;
                        </span>
                    </div>
                </div>
            </div>
        </template>
        <template v-else>
            <label
                v-if="showPrep"
                class="filters__switch"
            >
                <span>Bez przygotowania wcześniej</span>
                <input
                    v-model="f.noPrep"
                    type="checkbox"
                    role="switch"
                    class="sr-only"
                />
                <span
                    class="filters__track"
                    aria-hidden="true"
                />
            </label>

            <fieldset
                v-for="g in groups"
                :key="g.key"
                class="filters__group"
            >
                <legend class="filters__legend">{{ g.label }}</legend>
                <p
                    v-if="g.hint"
                    class="filters__hint"
                >
                    {{ g.hint }}
                </p>
                <div class="filters__chips">
                    <button
                        v-for="o in g.options"
                        :key="o.id"
                        type="button"
                        class="chip"
                        :class="[
                            {
                                'is-active': o.on,
                                'chip--color': g.kind === 'color',
                                'chip--strong': g.kind === 'strong',
                            },
                            o.color ? `cat-${o.color}` : null,
                        ]"
                        :aria-pressed="o.on"
                        :disabled="!o.on && !o.count"
                        @click="toggle(g.dim, o.id)"
                    >
                        <span
                            v-if="g.kind === 'color'"
                            class="chip__dot"
                            aria-hidden="true"
                        />
                        {{ o.label }}
                        <span class="filters__count">
                            <span class="sr-only">,</span>
                            {{ o.count }}
                        </span>
                    </button>
                </div>
            </fieldset>
        </template>
    </div>
</template>

<style scoped>
.filters__switch {
    position: relative;
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 12px;
    padding: 4px 0 18px;
    margin-bottom: 18px;
    border-bottom: 1px solid var(--rule);
    font-size: 14px;
    font-weight: 500;
    line-height: 1.35;
    cursor: pointer;
}
.filters__track {
    width: 36px;
    height: 20px;
    border-radius: var(--r-pill);
    flex-shrink: 0;
    background: var(--rule-strong);
    position: relative;
    transition: background 0.15s;
}
.filters__track::after {
    content: '';
    position: absolute;
    top: 2px;
    left: 2px;
    width: 16px;
    height: 16px;
    border-radius: var(--r-pill);
    background: #fff;
    transition: left 0.15s;
}
input:checked + .filters__track {
    background: var(--accent);
}
input:checked + .filters__track::after {
    left: 18px;
}
input:focus-visible + .filters__track {
    outline: 2px solid var(--accent);
    outline-offset: 2px;
}
.filters__group {
    border: none;
    padding: 0;
    margin: 0 0 22px;
    min-width: 0;
}
.filters__legend {
    font-family: var(--font-mono);
    font-size: 12px;
    color: var(--ink-muted);
    letter-spacing: 0.5px;
    text-transform: uppercase;
    font-weight: 600;
    padding: 0;
    margin-bottom: 10px;
}
.filters__hint {
    font-size: 12px;
    color: var(--ink-muted);
    margin: -4px 0 10px;
}
.filters__chips {
    display: flex;
    flex-wrap: wrap;
    gap: 6px;
}
.filters__count {
    font-family: var(--font-mono);
    font-size: 12px;
    font-weight: 500;
    color: var(--ink-muted);
}
.chip.is-active .filters__count {
    color: inherit;
}
.chip:disabled {
    opacity: 0.55;
    color: var(--ink-faint);
    cursor: default;
}
.chip:disabled:hover {
    border-color: var(--rule-strong);
}
.chip:disabled .filters__count {
    color: var(--ink-faint);
}
/* skeleton: the real .chip / .filters__legend (same heights), transparent content */
.fsk__row {
    flex-basis: 100%;
    display: flex;
    gap: 6px;
}
.fsk__chip {
    position: relative;
    width: 72px;
    color: transparent;
    background: transparent;
    border-color: transparent;
}
.fsk__fill {
    position: absolute;
    inset: 0;
    border-radius: var(--r-pill);
}
.fsk__chip:nth-child(2) {
    width: 96px;
}
.fsk__chip:nth-child(3) {
    width: 58px;
}
.fsk__row:nth-child(even) .fsk__chip:first-child {
    width: 110px;
}
.fsk__legend {
    display: inline-block;
    width: 88px;
    height: 12px;
    vertical-align: middle;
}
</style>
