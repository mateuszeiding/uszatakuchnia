<script lang="ts" setup>
import { computed } from 'vue';

import { bySection, scaledAmount } from '../recipeView';

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
const sections = computed(() => bySection(props.recipe.ingredients));

function amount(ing: RecipeDto['ingredients'][number]) {
    const s = scaledAmount(ing, props.dicts?.units ?? [], ratio.value);
    return [s.qty, s.unit].filter(Boolean).join(' ');
}

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
</script>

<template>
    <section
        class="ing-panel"
        aria-labelledby="h-ing"
    >
        <h2
            id="h-ing"
            class="ing-panel__title"
        >
            Składniki
        </h2>
        <div class="servings">
            <span
                id="lbl-serv"
                class="servings__label"
            >
                Porcje
            </span>
            <div
                class="servings__ctrl"
                role="group"
                aria-labelledby="lbl-serv"
            >
                <button
                    type="button"
                    class="servings__btn"
                    aria-label="Mniej porcji"
                    :disabled="servings <= 1"
                    @click="emit('servings', servings - 1)"
                >
                    −
                </button>
                <output
                    class="servings__value"
                    aria-live="polite"
                >
                    {{ servings }}
                </output>
                <button
                    type="button"
                    class="servings__btn"
                    aria-label="Więcej porcji"
                    :disabled="servings >= 20"
                    @click="emit('servings', servings + 1)"
                >
                    +
                </button>
            </div>
        </div>

        <div
            v-for="(sec, si) in sections"
            :key="si"
        >
            <h3
                v-if="sec.label"
                class="ing-panel__group"
            >
                {{ sec.label }}
            </h3>
            <ul class="ing-list">
                <li
                    v-for="ing in sec.items"
                    :key="ing.sortOrder"
                >
                    <button
                        type="button"
                        role="checkbox"
                        class="ing-item"
                        :class="{ 'is-done': checked.includes(ing.sortOrder) }"
                        :aria-checked="checked.includes(ing.sortOrder)"
                        @click="emit('toggle', ing.sortOrder)"
                    >
                        <span
                            class="ing-item__box"
                            aria-hidden="true"
                        >
                            <svg
                                v-if="checked.includes(ing.sortOrder)"
                                width="9"
                                height="9"
                                viewBox="0 0 10 10"
                            >
                                <path
                                    d="M2 5.2L4.2 7.4L8 3.4"
                                    stroke="#fff"
                                    stroke-width="1.8"
                                    stroke-linecap="round"
                                />
                            </svg>
                        </span>
                        <span class="ing-item__qty">
                            {{ amount(ing) }}
                            <span
                                v-if="ing.note"
                                class="ing-item__note"
                            >
                                {{ ing.note }}
                            </span>
                        </span>
                        <span class="ing-item__name">
                            {{ ing.ingredientName }}
                            <em
                                v-if="ing.optional"
                                class="ing-item__opt"
                            >
                                (opcjonalnie)
                            </em>
                        </span>
                    </button>
                </li>
            </ul>
        </div>

        <div
            v-if="allergens.length"
            class="allergens"
        >
            <h3 class="ing-panel__group ing-panel__group--plain">Alergeny</h3>
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
    </section>
</template>

<style scoped>
.ing-panel {
    position: sticky;
    top: 84px;
    background: var(--bg-alt);
    border: 1px solid var(--rule);
    border-radius: var(--r-xl);
    padding: 24px;
    display: flex;
    flex-direction: column;
    gap: 18px;
}
.ing-panel__title {
    font-size: 26px;
    font-weight: 500;
    letter-spacing: -0.8px;
}
.servings {
    background: var(--bg);
    border: 1px solid var(--rule);
    border-radius: var(--r-lg);
    padding: 8px 8px 8px 14px;
    display: flex;
    align-items: center;
    justify-content: space-between;
}
.servings__label {
    font-size: 14px;
    font-weight: 500;
    color: var(--ink-dim);
}
.servings__ctrl {
    display: flex;
    align-items: center;
    gap: 4px;
}
.servings__btn {
    width: 36px;
    height: 36px;
    border-radius: var(--r-md);
    background: var(--bg-alt);
    border: 1px solid var(--rule);
    color: var(--ink);
    font-size: 18px;
    cursor: pointer;
}
.servings__btn:disabled {
    background: transparent;
    color: var(--ink-faint);
    cursor: default;
}
.servings__value {
    min-width: 40px;
    text-align: center;
    font-size: 22px;
    font-weight: 600;
}
.ing-panel__group {
    font-family: var(--font-mono);
    font-size: 12px;
    color: var(--ink-muted);
    letter-spacing: 0.5px;
    text-transform: uppercase;
    font-weight: 600;
    padding-bottom: 8px;
    margin-bottom: 4px;
    border-bottom: 1px solid var(--rule);
}
.ing-panel__group--plain {
    border: none;
    padding: 0;
    margin-bottom: 8px;
}
.ing-list {
    list-style: none;
    padding: 0;
    margin: 0;
}
.ing-item {
    width: 100%;
    display: grid;
    grid-template-columns: 18px 92px minmax(0, 1fr);
    gap: 10px;
    align-items: baseline;
    background: transparent;
    border: none;
    border-bottom: 1px solid var(--rule);
    padding: 8px 0;
    cursor: pointer;
    text-align: left;
    color: var(--ink);
}
.ing-item__box {
    width: 16px;
    height: 16px;
    border-radius: var(--r-sm);
    border: 1.5px solid var(--rule-strong);
    display: inline-flex;
    align-items: center;
    justify-content: center;
    align-self: center;
}
.ing-item.is-done .ing-item__box {
    background: var(--accent);
    border-color: var(--accent);
}
.ing-item__qty {
    font-family: var(--font-mono);
    font-size: 13px;
    font-weight: 600;
}
.ing-item__note {
    display: block;
    font-size: 12px;
    font-weight: 400;
    color: var(--ink-muted);
}
.ing-item__name {
    font-size: 15px;
    line-height: 1.4;
    color: var(--ink-dim);
}
.ing-item__opt {
    color: var(--ink-muted);
}
.ing-item.is-done .ing-item__qty,
.ing-item.is-done .ing-item__name {
    color: var(--ink-muted);
    text-decoration: line-through;
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
        padding: 20px;
    }
    /* licznik porcji przyklejony pod paskiem nawigacji podczas przewijania listy */
    .servings {
        position: sticky;
        top: 64px;
        z-index: 2;
    }
    .ing-item {
        padding: 11px 0;
    }
}
</style>
