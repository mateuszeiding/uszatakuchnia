<script lang="ts" setup>
import { computed, onBeforeUnmount, onMounted, ref } from 'vue';

import { alsoIn, ingredientGroups, scaledAmount } from '../recipeView';

import IngGroupHeader from './IngGroupHeader.vue';
import IngRow from './IngRow.vue';
import ServingsStepper from './ServingsStepper.vue';

import type { RecipeDto } from '@/data/dtos/recipe/RecipeDto';
import type { IDictionaries } from '@/data/dtos/taxonomy/DictionariesDto';

const props = defineProps<{
    recipe: RecipeDto;
    dicts: IDictionaries | undefined;
    servings: number;
    checked: number[];
}>();
const emit = defineEmits<{ toggle: [sortOrder: number]; servings: [n: number] }>();

const ratio = computed(() => (props.recipe.servings ? props.servings / props.recipe.servings : 1));
const groups = computed(() => ingredientGroups(props.recipe.ingredients));
const shopLink = computed(() => ({
    path: `/recipes/${props.recipe.id}/zakupy`,
    query: { porcje: String(props.servings) },
}));

function amount(ing: RecipeDto['ingredients'][number]) {
    return scaledAmount(ing, props.dicts?.units ?? [], ratio.value);
}
const doneIn = (items: RecipeDto['ingredients']) =>
    items.filter((i) => props.checked.includes(i.sortOrder)).length;

/** Alergen → składniki, z których pochodzi (kolejność listy UE). */
const allergens = computed(() =>
    (props.dicts?.allergens ?? [])
        .map((a) => ({
            label: a.label,
            sources: props.recipe.ingredients
                .filter((i) => i.allergens.includes(a.code))
                .map((i) => i.ingredientName + (i.optional ? ' (opcj.)' : '')),
        }))
        .filter((a) => a.sources.length)
);

// Desktop: lista przewija się w panelu — gradient na dole tylko, gdy coś jest jeszcze niżej.
const list = ref<HTMLElement>();
const more = ref(false);
const measure = () => {
    const el = list.value;
    more.value = !!el && el.scrollHeight - el.scrollTop - el.clientHeight > 2;
};
let ro: ResizeObserver | undefined;
onMounted(() => {
    measure();
    ro = new ResizeObserver(measure);
    if (list.value) ro.observe(list.value);
});
onBeforeUnmount(() => ro?.disconnect());
</script>

<template>
    <section
        class="ing-panel"
        aria-labelledby="h-ing"
    >
        <div class="ing-panel__head">
            <h2
                id="h-ing"
                class="ing-panel__title"
            >
                Składniki
                <span
                    v-if="groups.single"
                    class="ing-panel__single"
                >
                    · {{ groups.single.toLowerCase() }}
                </span>
            </h2>
            <RouterLink
                :to="shopLink"
                class="btn btn--secondary btn--sm ing-panel__shop ing-panel__shop--top"
            >
                <svg
                    viewBox="0 0 16 16"
                    fill="none"
                    aria-hidden="true"
                >
                    <path
                        d="M1.5 2h2l1.6 8.2a1 1 0 001 .8h6.3a1 1 0 001-.8L14.5 5H4.3"
                        stroke="currentColor"
                        stroke-width="1.4"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                    />
                    <circle
                        cx="6.5"
                        cy="13.5"
                        r="1"
                        fill="currentColor"
                    />
                    <circle
                        cx="12"
                        cy="13.5"
                        r="1"
                        fill="currentColor"
                    />
                </svg>
                Lista zakupów
            </RouterLink>
        </div>
        <ServingsStepper
            class="ing-panel__servings"
            :value="servings"
            @set="emit('servings', $event)"
        />

        <div
            class="ing-panel__scroll"
            :class="{ 'has-more': more }"
        >
            <div
                ref="list"
                class="ing-panel__list"
                @scroll.passive="measure"
            >
                <div
                    v-for="(sec, si) in groups.sections"
                    :key="si"
                >
                    <IngGroupHeader
                        v-if="groups.grouped && sec.label"
                        class="ing-panel__ghead"
                        :label="sec.label"
                        :count="sec.items.length"
                        :done="doneIn(sec.items)"
                    />
                    <ul class="ing-list">
                        <IngRow
                            v-for="ing in sec.items"
                            :key="ing.sortOrder"
                            :ing="ing"
                            :amount="amount(ing)"
                            :checked="checked.includes(ing.sortOrder)"
                            :also="alsoIn(ing, recipe.ingredients)"
                            @toggle="emit('toggle', ing.sortOrder)"
                        />
                    </ul>
                </div>

                <div
                    v-if="allergens.length"
                    class="allergens"
                >
                    <h3 class="allergens__title">Alergeny</h3>
                    <ul class="allergens__list">
                        <li
                            v-for="a in allergens"
                            :key="a.label"
                        >
                            <strong>{{ a.label }}</strong>
                            <span class="allergens__src">· {{ a.sources.join(', ') }}</span>
                        </li>
                    </ul>
                </div>
            </div>
        </div>

        <RouterLink
            :to="shopLink"
            class="btn btn--secondary ing-panel__shop ing-panel__shop--bottom"
        >
            <svg
                viewBox="0 0 16 16"
                fill="none"
                aria-hidden="true"
            >
                <path
                    d="M1.5 2h2l1.6 8.2a1 1 0 001 .8h6.3a1 1 0 001-.8L14.5 5H4.3"
                    stroke="currentColor"
                    stroke-width="1.4"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                />
                <circle
                    cx="6.5"
                    cy="13.5"
                    r="1"
                    fill="currentColor"
                />
                <circle
                    cx="12"
                    cy="13.5"
                    r="1"
                    fill="currentColor"
                />
            </svg>
            Lista zakupów
        </RouterLink>
    </section>
</template>

<style scoped>
/* Desktop: panel przyklejony, nie wyższy niż okno; licznik porcji stoi, przewija się sama lista. */
.ing-panel {
    position: sticky;
    top: 84px;
    max-height: calc(100vh - 84px - 24px);
    background: var(--bg-alt);
    border: 1px solid var(--rule);
    border-radius: var(--r-xl);
    padding: 24px;
    display: flex;
    flex-direction: column;
    gap: 16px;
}
.ing-panel__head {
    display: flex;
    align-items: baseline;
    justify-content: space-between;
    gap: 12px;
    flex-wrap: wrap;
}
.ing-panel__title {
    font-size: 26px;
    font-weight: 500;
    letter-spacing: -0.8px;
}
.ing-panel__single {
    font-size: 15px;
    letter-spacing: 0;
    color: var(--ink-muted);
}
.ing-panel__shop {
    gap: 7px;
    white-space: nowrap;
}
.ing-panel__shop--bottom {
    display: none;
}
.ing-panel__scroll {
    position: relative;
    display: flex;
    flex: 1 1 auto;
    min-height: 0;
}
.ing-panel__list {
    flex: 1;
    min-height: 0;
    overflow-y: auto;
    padding-right: 6px;
    display: flex;
    flex-direction: column;
    gap: 22px;
}
.ing-panel__scroll.has-more::after {
    content: '';
    position: absolute;
    left: 0;
    right: 6px;
    bottom: 0;
    height: 36px;
    background: linear-gradient(to bottom, transparent, var(--bg-alt));
    pointer-events: none;
}
/* nagłówek grupy przykleja się u góry przewijanej listy */
.ing-panel__ghead {
    position: sticky;
    top: 0;
    z-index: 1;
    background: var(--bg-alt);
    padding-top: 8px;
}
.ing-list {
    list-style: none;
    padding: 0;
    margin: 0;
}
.allergens__title {
    font-family: var(--font-mono);
    font-size: 12px;
    color: var(--ink-muted);
    letter-spacing: 0.5px;
    text-transform: uppercase;
    font-weight: 600;
    margin-bottom: 8px;
}
.allergens__list {
    list-style: none;
    padding: 0;
    margin: 0;
    display: flex;
    flex-direction: column;
    gap: 4px;
    font-size: 14px;
    line-height: 1.45;
}
.allergens__list strong {
    font-weight: 600;
}
.allergens__src {
    margin-left: 4px;
    color: var(--ink-muted);
}

@media (max-width: 900px) {
    .ing-panel {
        position: static;
        max-height: none;
        padding: 20px;
    }
    .ing-panel__list {
        overflow: visible;
        padding-right: 0;
    }
    .ing-panel__scroll.has-more::after {
        content: none;
    }
    /* licznik porcji przyklejony pod paskiem nawigacji, nagłówek grupy pod licznikiem (64 + 54 px) */
    .ing-panel__servings {
        position: sticky;
        top: 64px;
        z-index: 2;
    }
    .ing-panel__ghead {
        top: 118px;
    }
    .ing-panel__shop--top {
        display: none;
    }
    .ing-panel__shop--bottom {
        display: inline-flex;
        width: 100%;
        justify-content: center;
    }
}
</style>
