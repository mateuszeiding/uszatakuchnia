<script lang="ts" setup>
import { computed, onBeforeUnmount, ref } from 'vue';

import MacroValue from '@/components/MacroValue.vue';
import {
    avg,
    energySplit,
    fmt,
    type IRecipeNutrition,
    nf,
    NUTRIENTS,
    scaleRange,
} from '@/shared/utils/nutrition';

// Makro na porcję (wyliczone w API). Licznik porcji zmienia tylko „Całe danie”.
const props = defineProps<{ nutrition: IRecipeNutrition; servings: number }>();

const pp = computed(() => props.nutrition.perPortion);
const incomplete = computed(() => props.nutrition.incomplete ?? {});
const hasIncomplete = computed(() => Object.values(incomplete.value).some((v) => v?.length));
const kcal = computed(() => fmt(pp.value.kcal, 'kcal', 'recipe'));
const kcalAvg = computed(() => (kcal.value ? nf(avg(kcal.value)) : null));
const split = computed(() => energySplit(pp.value));
const splitParts = computed(() =>
    split.value
        ? [
              { key: 'protein', label: 'białko', pct: split.value.protein },
              { key: 'fat', label: 'tłuszcz', pct: split.value.fat },
              { key: 'carbs', label: 'węglowodany', pct: split.value.carbs },
          ]
        : []
);

const portionsLabel = computed(() => {
    const n = props.servings;
    if (n === 1) return 'porcja';
    const d = n % 10;
    const t = n % 100;
    return d >= 2 && d <= 4 && (t < 10 || t >= 20) ? 'porcje' : 'porcji';
});

// Mobile: sekcja zwijana, domyślnie zwinięta.
const mq = window.matchMedia('(max-width: 900px)');
const mobile = ref(mq.matches);
const open = ref(!mq.matches);
const onMq = (e: MediaQueryListEvent) => {
    mobile.value = e.matches;
    open.value = !e.matches;
};
mq.addEventListener('change', onMq);
onBeforeUnmount(() => mq.removeEventListener('change', onMq));
</script>

<template>
    <section
        id="makro"
        class="nutri"
        aria-labelledby="makro-h"
    >
        <h2
            id="makro-h"
            class="nutri__h"
        >
            <button
                v-if="mobile"
                type="button"
                class="nutri__toggle"
                :aria-expanded="open"
                aria-controls="makro-b"
                @click="open = !open"
            >
                <span class="nutri__title">
                    Wartości odżywcze
                    <span class="badge-attr nutri__badge">orientacyjnie</span>
                </span>
                <span class="nutri__toggle-side">
                    <span
                        v-if="!open && kcalAvg"
                        class="nutri__peek"
                    >
                        <span class="nutri__peek-num">~{{ kcalAvg }}</span>
                        kcal/porcja
                    </span>
                    <svg
                        class="nutri__chev"
                        :class="{ 'is-open': open }"
                        width="14"
                        height="14"
                        viewBox="0 0 14 14"
                        fill="none"
                        aria-hidden="true"
                    >
                        <path
                            d="M3.5 5.25 7 8.75l3.5-3.5"
                            stroke="currentColor"
                            stroke-width="1.5"
                            stroke-linecap="round"
                            stroke-linejoin="round"
                        />
                    </svg>
                </span>
            </button>
            <span
                v-else
                class="nutri__title"
            >
                Wartości odżywcze
                <span class="badge-attr nutri__badge">orientacyjnie</span>
            </span>
        </h2>

        <div
            v-show="open"
            id="makro-b"
            class="nutri__body"
        >
            <div class="nutri__cols">
                <div class="nutri__col">
                    <div v-if="kcal">
                        <div class="nutri__label">Na porcję</div>
                        <div class="nutri__kcal">
                            <span class="sr-only">około {{ kcalAvg }}</span>
                            <span
                                class="nutri__kcal-num"
                                aria-hidden="true"
                            >
                                {{ `~${kcalAvg}` }}
                            </span>
                            <span class="nutri__kcal-unit">kcal</span>
                        </div>
                        <div
                            v-if="!kcal.single"
                            class="nutri__range"
                        >
                            zakres {{ nf(kcal.a) }} – {{ nf(kcal.b) }} kcal
                        </div>
                        <span
                            v-if="incomplete.kcal?.length"
                            class="nutri__without"
                        >
                            bez: {{ incomplete.kcal.join(', ') }}
                        </span>
                    </div>

                    <dl class="nutri__rows">
                        <template
                            v-for="n in NUTRIENTS.slice(1)"
                            :key="n.key"
                        >
                            <dt>{{ n.label }}</dt>
                            <dd>
                                <MacroValue
                                    :value="pp[n.key]"
                                    :nutrient="n.key"
                                    mode="recipe"
                                />
                                <span
                                    v-if="incomplete[n.key]?.length"
                                    class="nutri__without"
                                >
                                    bez: {{ incomplete[n.key]!.join(', ') }}
                                </span>
                            </dd>
                        </template>
                    </dl>
                </div>

                <div class="nutri__col">
                    <figure
                        v-if="split"
                        class="nutri__split"
                    >
                        <figcaption class="nutri__label">Skąd energia · średnio</figcaption>
                        <div
                            class="nutri__bar"
                            aria-hidden="true"
                        >
                            <span
                                v-for="p in splitParts"
                                :key="p.key"
                                :class="`nutri__seg nutri__seg--${p.key}`"
                                :style="{ width: `${p.pct}%` }"
                            />
                        </div>
                        <ul class="nutri__legend">
                            <li
                                v-for="p in splitParts"
                                :key="p.key"
                            >
                                <span
                                    :class="`nutri__dot nutri__seg--${p.key}`"
                                    aria-hidden="true"
                                />
                                {{ p.label }}
                                <span class="nutri__pct">~{{ p.pct }}%</span>
                            </li>
                        </ul>
                    </figure>

                    <div class="nutri__dish">
                        <div class="nutri__dish-head">
                            <span class="nutri__label">Całe danie</span>
                            <span class="nutri__dish-count">
                                {{ servings }} {{ portionsLabel }}
                            </span>
                        </div>
                        <dl class="nutri__dish-rows">
                            <template
                                v-for="n in NUTRIENTS"
                                :key="n.key"
                            >
                                <dt>{{ n.label }}</dt>
                                <dd>
                                    <MacroValue
                                        :value="scaleRange(pp[n.key], servings)"
                                        :nutrient="n.key"
                                        mode="recipe"
                                        size="sm"
                                        muted
                                    />
                                </dd>
                            </template>
                        </dl>
                    </div>
                </div>
            </div>

            <p
                v-if="nutrition.missing.length || hasIncomplete"
                class="nutri__missing"
            >
                <template v-if="nutrition.missing.length">
                    <strong>Nie wliczono:</strong>
                    {{ nutrition.missing.join(', ') }}
                    <span class="nutri__muted">— brak danych lub przelicznika jednostki.</span>
                </template>
                <span
                    v-if="hasIncomplete"
                    class="nutri__muted nutri__inc-note"
                    :class="{ 'nutri__inc-note--below': nutrition.missing.length }"
                >
                    „bez: …” przy wartości — składnik jest wliczony, ale nie ma tej jednej wartości,
                    więc ona może być zaniżona.
                </span>
            </p>
            <p class="nutri__note">
                Wartości orientacyjne, wyliczone ze składników. Produkty różnych marek się różnią —
                stąd zakres. Liczba porcji zmienia tylko sumę dla całego dania.
            </p>
        </div>
    </section>
</template>

<style scoped>
.nutri {
    background: var(--bg-alt);
    border: 1px solid var(--rule);
    border-radius: var(--r-xl);
    padding: 24px;
    scroll-margin-top: 84px;
    container-type: inline-size;
}
.nutri__h {
    margin: 0;
    font: inherit;
}
.nutri__title {
    display: flex;
    align-items: center;
    gap: 10px;
    flex-wrap: wrap;
    font-size: 22px;
    font-weight: 500;
    letter-spacing: -0.6px;
    color: var(--ink);
}
.nutri__badge {
    color: var(--info);
    border-color: color-mix(in srgb, var(--info) 40%, transparent);
    letter-spacing: 0;
}
.nutri__toggle {
    width: 100%;
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 12px;
    background: none;
    border: 0;
    padding: 0;
    color: var(--ink);
    text-align: left;
    cursor: pointer;
    font: inherit;
}
.nutri__toggle-side {
    display: inline-flex;
    align-items: center;
    gap: 10px;
    color: var(--ink-muted);
}
.nutri__peek {
    white-space: nowrap;
    font-size: 12px;
}
.nutri__peek-num {
    font-family: var(--font-mono);
    font-size: 15px;
    font-weight: 600;
    color: var(--ink);
}
.nutri__chev {
    flex-shrink: 0;
    transition: transform 0.15s;
}
.nutri__chev.is-open {
    transform: rotate(180deg);
}
.nutri__body {
    margin-top: 18px;
    display: flex;
    flex-direction: column;
    gap: 18px;
}
.nutri__label {
    font-family: var(--font-mono);
    font-size: 12px;
    font-weight: 600;
    letter-spacing: 0.5px;
    text-transform: uppercase;
    color: var(--ink-muted);
    margin-bottom: 6px;
}
.nutri__kcal {
    display: flex;
    align-items: baseline;
    gap: 8px;
    flex-wrap: wrap;
}
.nutri__kcal-num {
    font-family: var(--font-mono);
    font-size: 36px;
    font-weight: 600;
    letter-spacing: -1.2px;
    line-height: 1;
}
.nutri__kcal-unit {
    font-size: 15px;
    font-weight: 500;
    color: var(--ink-dim);
}
.nutri__range {
    font-family: var(--font-mono);
    font-size: 13px;
    color: var(--ink-muted);
    margin-top: 6px;
}
.nutri__rows {
    margin: 0;
    display: grid;
    grid-template-columns: minmax(0, 1fr) auto;
    border-top: 1px solid var(--rule);
}
.nutri__rows dt,
.nutri__rows dd {
    padding: 9px 0;
    border-bottom: 1px solid var(--rule);
}
.nutri__rows dt {
    font-size: 14px;
    color: var(--ink-dim);
}
.nutri__rows dd {
    margin: 0;
    text-align: right;
}
.nutri__without {
    display: block;
    margin-top: 2px;
    font-size: 12px;
    color: var(--ink-muted);
}
.nutri__inc-note {
    display: block;
}
.nutri__inc-note--below {
    margin-top: 4px;
}
.nutri__cols {
    display: grid;
    grid-template-columns: minmax(0, 1fr);
    gap: 18px;
    align-items: start;
}
/* desktop (prawa kolumna, pod krokami): kcal + wiersze | rozkład energii + całe danie — gdy jest miejsce */
@container (min-width: 600px) {
    .nutri__cols {
        grid-template-columns: minmax(0, 1fr) minmax(0, 1fr);
        gap: 32px;
    }
}
.nutri__col {
    display: flex;
    flex-direction: column;
    gap: 18px;
    min-width: 0;
}
.nutri__split {
    margin: 0;
}
.nutri__split .nutri__label {
    margin-bottom: 8px;
}
.nutri__bar {
    display: flex;
    gap: 2px;
    height: 6px;
    border-radius: var(--r-pill);
    overflow: hidden;
}
.nutri__seg--protein {
    background: var(--accent);
}
.nutri__seg--fat {
    background: var(--c-mustard);
}
.nutri__seg--carbs {
    background: var(--c-teal);
}
.nutri__legend {
    list-style: none;
    padding: 0;
    margin: 8px 0 0;
    display: flex;
    flex-wrap: wrap;
    gap: 4px 14px;
    font-size: 12px;
    color: var(--ink-dim);
}
.nutri__legend li {
    display: inline-flex;
    align-items: center;
    gap: 6px;
}
.nutri__dot {
    width: 7px;
    height: 7px;
    border-radius: 2px;
}
.nutri__pct {
    font-family: var(--font-mono);
    color: var(--ink);
}
.nutri__dish {
    background: var(--bg);
    border: 1px solid var(--rule);
    border-radius: var(--r-lg);
    padding: 12px 14px;
}
.nutri__dish-head {
    display: flex;
    justify-content: space-between;
    align-items: baseline;
    gap: 8px;
    margin-bottom: 8px;
}
.nutri__dish-head .nutri__label {
    margin-bottom: 0;
}
.nutri__dish-count {
    font-size: 12px;
    color: var(--ink-muted);
}
.nutri__dish-rows {
    margin: 0;
    display: grid;
    grid-template-columns: minmax(0, 1fr) auto;
    row-gap: 4px;
    font-size: 12px;
}
.nutri__dish-rows dt {
    color: var(--ink-muted);
}
.nutri__dish-rows dd {
    margin: 0;
    text-align: right;
}
.nutri__missing {
    margin: 0;
    font-size: 13px;
    line-height: 1.45;
    color: var(--ink-dim);
    padding-left: 10px;
    border-left: 2px solid var(--info);
}
.nutri__missing strong {
    font-weight: 600;
    color: var(--ink);
}
.nutri__muted {
    color: var(--ink-muted);
}
.nutri__note {
    margin: 0;
    font-size: 12px;
    line-height: 1.5;
    color: var(--ink-muted);
}

@media (max-width: 900px) {
    .nutri {
        padding: 20px;
    }
    .nutri__title {
        font-size: 20px;
    }
    .nutri__kcal-num {
        font-size: 32px;
    }
}
</style>
