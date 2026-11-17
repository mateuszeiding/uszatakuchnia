<script lang="ts" setup>
import { computed } from 'vue';

import { DIFFICULTY, type IListFilters, TIMES } from '../listFilters';

import type { ITaxonomyAxis } from '@/data/dtos/taxonomy/DictionariesDto';

// Panel filtrów (sidebar na desktopie, zawartość bottom sheetu na mobile).
// loading: słowniki jeszcze się wczytują — sekcje statyczne (przełącznik, czas, trudność) są od razu prawdziwe,
// sekcje osi to szkielety w tych samych miejscach, z tymi samymi klasami (legenda, chipy) i liczbą rzędów
// chipów jak po załadowaniu → panel ma od początku docelową wysokość.
const props = defineProps<{ axes: ITaxonomyAxis[]; loading?: boolean }>();

/** Rzędy chipów w panelu 242 px (słowniki z 17.11: 9 / 9 / 11 / 11 / 7 wartości). */
const SK_ROWS = { dish: 5, diet: 5, rest: [4, 5, 4] };
const f = defineModel<IListFilters>({ required: true });

// Kolejność jak w prototypie: typ dania, czas, trudność, diety, potem reszta osi.
const dish = computed(() => props.axes.find((a) => a.code === 'dish'));
const diet = computed(() => props.axes.find((a) => a.code === 'diet'));
const rest = computed(() => props.axes.filter((a) => a.code !== 'dish' && a.code !== 'diet'));

const isOn = (code: string, id: number) => (f.value.axes[code] ?? []).includes(id);

function toggle(code: string, id: number) {
    const cur = f.value.axes[code] ?? [];
    f.value = {
        ...f.value,
        axes: {
            ...f.value.axes,
            [code]: cur.includes(id) ? cur.filter((x) => x !== id) : [...cur, id],
        },
    };
}
</script>

<template>
    <div
        class="filters"
        :inert="loading"
    >
        <label class="filters__switch">
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

        <div
            v-if="loading"
            class="filters__group"
            aria-hidden="true"
        >
            <div class="filters__legend"><span class="sk fsk__legend" /></div>
            <div class="filters__chips">
                <div
                    v-for="r in SK_ROWS.dish"
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
        <fieldset
            v-if="dish"
            class="filters__group"
        >
            <legend class="filters__legend">{{ dish.label }}</legend>
            <div class="filters__chips">
                <button
                    v-for="v in dish.values"
                    :key="v.id"
                    type="button"
                    class="chip chip--color"
                    :class="[
                        { 'is-active': isOn('dish', v.id) },
                        v.color ? `cat-${v.color}` : null,
                    ]"
                    :aria-pressed="isOn('dish', v.id)"
                    @click="toggle('dish', v.id)"
                >
                    <span
                        class="chip__dot"
                        aria-hidden="true"
                    />
                    {{ v.label }}
                </button>
            </div>
        </fieldset>

        <fieldset class="filters__group">
            <legend class="filters__legend">Czas</legend>
            <div class="filters__chips">
                <button
                    v-for="[v, l] in TIMES"
                    :key="v"
                    type="button"
                    class="chip"
                    :class="{ 'is-active': f.maxTime === v }"
                    :aria-pressed="f.maxTime === v"
                    @click="f = { ...f, maxTime: f.maxTime === v ? 0 : v }"
                >
                    {{ l }}
                </button>
            </div>
        </fieldset>

        <fieldset class="filters__group">
            <legend class="filters__legend">Trudność</legend>
            <div class="filters__chips">
                <button
                    v-for="v in [1, 2, 3]"
                    :key="v"
                    type="button"
                    class="chip"
                    :class="{ 'is-active': f.difficulty === v }"
                    :aria-pressed="f.difficulty === v"
                    @click="f = { ...f, difficulty: f.difficulty === v ? 0 : v }"
                >
                    {{ DIFFICULTY[v] }}
                </button>
            </div>
        </fieldset>

        <div
            v-if="loading"
            class="filters__group"
            aria-hidden="true"
        >
            <div class="filters__legend"><span class="sk fsk__legend" /></div>
            <p class="filters__hint">Przepis spełnia wszystkie zaznaczone.</p>
            <div class="filters__chips">
                <div
                    v-for="r in SK_ROWS.diet"
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
        <fieldset
            v-if="diet"
            class="filters__group"
        >
            <legend class="filters__legend">{{ diet.label }}</legend>
            <p class="filters__hint">Przepis spełnia wszystkie zaznaczone.</p>
            <div class="filters__chips">
                <button
                    v-for="v in diet.values"
                    :key="v.id"
                    type="button"
                    class="chip chip--strong"
                    :class="{ 'is-active': isOn('diet', v.id) }"
                    :aria-pressed="isOn('diet', v.id)"
                    @click="toggle('diet', v.id)"
                >
                    {{ v.label }}
                </button>
            </div>
        </fieldset>

        <div
            v-for="(n, i) in loading ? SK_ROWS.rest : []"
            :key="`sk-${i}`"
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
        <fieldset
            v-for="a in rest"
            :key="a.id"
            class="filters__group"
        >
            <legend class="filters__legend">{{ a.label }}</legend>
            <div class="filters__chips">
                <button
                    v-for="v in a.values"
                    :key="v.id"
                    type="button"
                    class="chip"
                    :class="{ 'is-active': isOn(a.code, v.id) }"
                    :aria-pressed="isOn(a.code, v.id)"
                    @click="toggle(a.code, v.id)"
                >
                    {{ v.label }}
                </button>
            </div>
        </fieldset>
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
/* szkielet: prawdziwe .chip / .filters__legend (te same wysokości), treść przezroczysta */
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
