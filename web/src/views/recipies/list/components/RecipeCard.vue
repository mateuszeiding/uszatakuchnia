<script lang="ts" setup>
import { computed } from 'vue';

import { DIFFICULTY } from '../listFilters';

import { cuisineOf, mainDish, type RecipeBaseDto } from '@/data/dtos/recipe/RecipeDto';
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
</script>

<template>
    <RouterLink
        :to="{ name: 'recipe-details', params: { id: recipe.id } }"
        class="r-card"
    >
        <div class="r-card__photo">
            <img
                v-if="recipe.photoUrl"
                :src="recipe.photoUrl"
                alt=""
                loading="lazy"
            />
            <div
                v-else
                class="photo r-card__placeholder"
            />
            <span
                v-if="recipe.needsPrep"
                class="r-card__prep"
            >
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
                wymaga przygotowania
            </span>
        </div>
        <div class="r-card__body">
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
                v-if="meta"
                class="r-card__meta"
            >
                {{ meta }}
            </div>
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
    transition:
        transform 0.15s,
        box-shadow 0.15s;
}
.r-card:hover {
    text-decoration: none;
    transform: translateY(-2px);
    box-shadow: var(--shadow-card);
}
.r-card__photo {
    aspect-ratio: 1 / 1;
    position: relative;
    border-bottom: 1px solid var(--rule);
}
.r-card__photo img,
.r-card__placeholder {
    width: 100%;
    height: 100%;
    object-fit: cover;
}
.r-card__prep {
    position: absolute;
    left: 10px;
    bottom: 10px;
    display: inline-flex;
    align-items: center;
    gap: 5px;
    background: rgba(255, 248, 230, 0.95);
    color: #5c420a;
    padding: 4px 10px;
    border-radius: var(--r-pill);
    font-size: 12px;
    font-weight: 600;
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
    flex: 1;
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
    color: var(--ink-muted);
}
</style>
