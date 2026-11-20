<script lang="ts" setup>
import { computed } from 'vue';

import { DIFFICULTY } from '../listFilters';

import { cuisineOf, mainDish, type RecipeBaseDto } from '@/data/dtos/recipe/RecipeDto';
import { avg, fmt, nf } from '@/shared/utils/nutrition';
import { formatTime } from '@/shared/utils/time';

const props = defineProps<{ recipe: RecipeBaseDto; hits?: string[] | null }>();

const dish = computed(() => mainDish(props.recipe));
const cuisine = computed(() => cuisineOf(props.recipe)?.label);
const meta = computed(() =>
    [
        formatTime(props.recipe.timeMin, props.recipe.timeMax),
        props.recipe.difficulty ? DIFFICULTY[props.recipe.difficulty] : null,
    ]
        .filter(Boolean)
        .join(' · ')
);
const kcal = computed(() => {
    const f = fmt(props.recipe.kcal, 'kcal', 'recipe');
    return f ? `~${nf(avg(f))} kcal` : null;
});
</script>

<template>
    <RouterLink
        :to="{ name: 'recipe-details', params: { id: recipe.id } }"
        class="r-card"
    >
        <!-- bez zdjęcia — bez kadru (20.11: żadnego zastępczego obrazka) -->
        <div
            v-if="recipe.photoUrl"
            class="r-card__photo"
        >
            <img
                :src="recipe.photoUrl"
                alt=""
                loading="lazy"
            />
        </div>
        <div class="r-card__body">
            <svg
                width="12"
                height="12"
                viewBox="0 0 12 12"
                aria-hidden="true"
            >
                <circle
                    cx="6"
                    cy="6"
                    r="4.8"
                    stroke="currentColor"
                    stroke-width="1.3"
                />
                <path
                    d="M6 3.5v2.7l1.7 1"
                    stroke="currentColor"
                    stroke-width="1.3"
                    stroke-linecap="round"
                />
            </svg>
            <div
                v-if="dish || cuisine"
                class="r-card__flags"
            >
                <span
                    v-if="dish"
                    class="badge"
                    :class="dish.color ? `cat-${dish.color}` : null"
                >
                    {{ dish.label }}
                </span>
                <span
                    v-if="cuisine"
                    class="r-card__cuisine"
                >
                    {{ cuisine }}
                </span>
            </div>
            <h2 class="r-card__title">{{ recipe.name }}</h2>
            <div
                v-if="hits?.length"
                class="r-card__hits"
            >
                zawiera:
                <mark>{{ hits.join(', ') }}</mark>
            </div>
            <div
                v-if="meta || kcal"
                class="r-card__meta"
            >
                {{ meta && kcal ? `${meta} · ` : meta }}
                <span
                    v-if="kcal"
                    class="r-card__kcal"
                >
                    {{ kcal }}
                </span>
                <span
                    v-if="kcal"
                    class="sr-only"
                >
                    na porcję, orientacyjnie
                </span>
            </div>
            <span
                v-if="recipe.needsPrep"
                class="r-card__prep"
            >
                <svg
                    width="13"
                    height="13"
                    viewBox="0 0 12 12"
                    fill="none"
                    aria-hidden="true"
                >
                    <circle
                        cx="6"
                        cy="6"
                        r="4.8"
                        stroke="currentColor"
                        stroke-width="1.3"
                    />
                    <path
                        d="M6 3.5v2.7l1.7 1"
                        stroke="currentColor"
                        stroke-width="1.3"
                        stroke-linecap="round"
                    />
                </svg>
                wymaga przygotowania
            </span>
        </div>
    </RouterLink>
</template>

<style scoped>
.r-card {
    display: flex;
    flex-direction: column;
    height: 100%;
    background: var(--bg-alt);
    border: 1px solid var(--rule);
    border-radius: var(--r-lg);
    overflow: hidden;
    color: var(--ink);
    text-decoration: none;
    transition: border-color var(--dur-fast) var(--ease);
}
/* hover nie przesuwa karty — translateY rozmywał tekst przy skalowaniu ekranu Windows */
@media (hover: hover) {
    .r-card:hover {
        border-color: var(--rule-strong);
    }
    .r-card:hover .r-card__title {
        color: var(--accent);
    }
}
.r-card:focus-visible {
    border-color: var(--rule-strong);
}
.r-card:focus-visible .r-card__title {
    color: var(--accent);
}
.r-card__photo {
    /* kwadrat bez linii pod spodem (linia dochodzi) — ten sam wymiar w szkielecie */
    box-sizing: content-box;
    aspect-ratio: 1 / 1;
    position: relative;
    border-bottom: 1px solid var(--rule);
}
.r-card__photo img {
    display: block; /* inline img dokładał linię tekstu — zdjęcie było o 1 px wyższe niż 1:1 */
    width: 100%;
    height: 100%;
    object-fit: cover;
}
/* „wymaga przygotowania” w treści karty (dawniej na kadrze) */
.r-card__prep {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    font-size: 12px;
    font-weight: 600;
    color: var(--prep);
}
.r-card__body {
    padding: 16px;
    display: flex;
    flex-direction: column;
    gap: 10px;
    flex: 1;
}
.r-card__flags {
    display: flex;
    align-items: center;
    gap: 8px;
    flex-wrap: wrap;
    font-size: 12px;
}
.r-card__cuisine {
    color: var(--ink-muted);
}
.r-card__title {
    font-size: 18px;
    font-weight: 600;
    letter-spacing: -0.4px;
    line-height: 1.25;
    text-wrap: balance;
    transition: color var(--dur-fast) var(--ease);
    flex: 1;
    /* dwie linie zarezerwowane — karta ma tę samą wysokość co szkielet i sąsiednie karty */
    min-height: calc(2 * 1.25em);
}
.r-card__hits {
    font-size: 13px;
    color: var(--ink-dim);
}
.r-card__hits mark {
    background: var(--accent-soft);
    color: var(--ink);
    padding: 0 3px;
    border-radius: 3px;
}
.r-card__meta {
    font-size: 13px;
    line-height: 20px; /* stała — mono w kcal nie podnosi linii; tyle samo w szkielecie */
    color: var(--ink-muted);
}
.r-card__kcal {
    font-family: var(--font-mono);
}
</style>
