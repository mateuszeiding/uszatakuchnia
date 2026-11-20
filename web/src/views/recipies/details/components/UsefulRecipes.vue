<script lang="ts" setup>
import { computed } from 'vue';

import { usefulRecipes } from '../recipeView';

import type { RecipeIngredientDto } from '@/data/dtos/recipe/RecipeIngredientDto';

// "Przydatne przepisy": published recipes that the ingredients link to.
// Card = one link: 1:1 thumbnail, title (max 2 lines), arrow. No "do: …" (it repeated the title;
// the row icon tells which ingredient).
// Desktop: a 3-column grid; ≤ 900 px: a list separated by lines. No links → no block.
// A card opens a new tab.
const props = defineProps<{ ingredients: RecipeIngredientDto[] }>();
const items = computed(() => usefulRecipes(props.ingredients));
</script>

<template>
    <section
        v-if="items.length"
        class="useful"
        aria-labelledby="h-useful"
    >
        <h2
            id="h-useful"
            class="useful__title"
        >
            Przydatne przepisy
        </h2>
        <ul class="useful__list">
            <li
                v-for="r in items"
                :key="r.id"
            >
                <RouterLink
                    :to="`/recipes/${r.id}`"
                    class="useful__card"
                    target="_blank"
                    rel="noopener"
                    :aria-label="`${r.title} (otwiera się w nowej karcie)`"
                >
                    <img
                        v-if="r.image"
                        :src="r.image"
                        alt=""
                        class="useful__thumb"
                    />
                    <span
                        v-else
                        class="photo useful__thumb"
                        aria-hidden="true"
                    />
                    <span class="useful__text">
                        <span class="useful__name">{{ r.title }}</span>
                    </span>
                    <svg
                        class="useful__arrow"
                        width="12"
                        height="12"
                        viewBox="0 0 12 12"
                        fill="none"
                        aria-hidden="true"
                    >
                        <path
                            d="M2 6h7.5M6.5 3l3 3-3 3"
                            stroke="currentColor"
                            stroke-width="1.5"
                            stroke-linecap="round"
                            stroke-linejoin="round"
                        />
                    </svg>
                </RouterLink>
            </li>
        </ul>
    </section>
</template>

<style scoped>
.useful {
    margin-top: 28px;
}
.useful__title {
    margin: 0 0 10px;
    font-family: var(--font-mono);
    font-size: 12px;
    font-weight: 600;
    letter-spacing: 0.5px;
    text-transform: uppercase;
    color: var(--ink-muted);
}
.useful__list {
    list-style: none;
    padding: 0;
    margin: 0;
    display: grid;
    grid-template-columns: repeat(3, minmax(0, 1fr));
    gap: 10px;
}
.useful__card {
    height: 100%;
    display: grid;
    grid-template-columns: 48px minmax(0, 1fr) 14px;
    gap: 12px;
    align-items: center;
    min-height: 56px;
    padding: 10px 12px;
    background: var(--bg-alt);
    border: 1px solid var(--rule);
    border-radius: var(--r-lg);
    color: var(--ink);
    text-decoration: none;
    transition: border-color var(--dur-fast) var(--ease);
}
@media (hover: hover) {
    .useful__card:hover {
        border-color: var(--rule-strong);
    }
    .useful__card:hover .useful__name {
        color: var(--accent);
    }
}
.useful__card:focus-visible {
    border-color: var(--rule-strong);
}
.useful__card:focus-visible .useful__name {
    color: var(--accent);
}
.useful__thumb {
    width: 48px;
    height: 48px;
    border-radius: var(--r-md);
    border: 1px solid var(--rule);
    object-fit: cover;
}
.useful__text {
    min-width: 0;
}
.useful__name {
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
    font-size: 14px;
    font-weight: 600;
    line-height: 1.3;
    transition: color var(--dur-fast) var(--ease);
}
.useful__arrow {
    color: var(--ink-muted);
}

@media (max-width: 900px) {
    .useful {
        margin-top: 24px;
    }
    .useful__list {
        display: block;
        border-top: 1px solid var(--rule);
    }
    .useful__card {
        grid-template-columns: 40px minmax(0, 1fr) 14px;
        padding: 8px 0;
        background: transparent;
        border: none;
        border-bottom: 1px solid var(--rule);
        border-radius: 0;
    }
    .useful__thumb {
        width: 40px;
        height: 40px;
    }
}
</style>
