<script lang="ts" setup>
import { computed } from 'vue';

import { DIFFICULTY } from '../listFilters';

import { cuisineOf, mainDish, type RecipeBaseDto } from '@/data/dtos/recipe/RecipeDto';
import { avg, fmt, nf } from '@/shared/utils/nutrition';
import { recipeHue } from '@/shared/utils/recipeHue';
import { formatTime } from '@/shared/utils/time';

const props = defineProps<{ recipe: RecipeBaseDto; hits?: string[] | null }>();

// read by the placeholder frame through v-bind() in <style>
const hue = computed(() => recipeHue(props.recipe.id));

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
        <!-- the 1:1 frame is always there (even rows in a mixed grid); without a photo:
             a gradient in the recipe's hue + a plate. Decorative — "needs prep" is repeated in the body for screen readers. -->
        <div
            class="r-card__photo"
            aria-hidden="true"
        >
            <img
                v-if="recipe.photoUrl"
                :src="recipe.photoUrl"
                alt=""
                loading="lazy"
            />
            <span
                v-else
                class="r-card__plate"
            />
            <span
                v-if="recipe.needsPrep"
                class="r-card__prep"
            >
                <svg
                    width="12"
                    height="12"
                    viewBox="0 0 12 12"
                    fill="none"
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
                class="sr-only"
            >
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
/* hover doesn't move the card: translateY blurred the text with Windows display scaling */
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
    /* placeholder colors in the recipe's hue (h = id * 47 % 360); dark overrides below */
    --h: v-bind(hue);
    --thumb-from: hsl(var(--h) 35% 90%);
    --thumb-to: hsl(calc(var(--h) + 20) 45% 74%);
    --plate: rgba(255, 255, 255, 0.8);
    --plate-shadow: 0 6px 18px rgba(0, 0, 0, 0.08);
    --plate-dot: hsl(var(--h) 60% 62%);
    --prep-pill: rgba(255, 248, 230, 0.95);

    /* a square without an extra line below; the same size in the skeleton */
    box-sizing: content-box;
    aspect-ratio: 1 / 1;
    position: relative;
    overflow: hidden;
    background: linear-gradient(135deg, var(--thumb-from), var(--thumb-to));
    border-bottom: 1px solid var(--rule);
}
/* dark: dimmed (lightness 19–36 %), no bright spot; the same two selectors as the theme tokens.
   No :global() here: it would drop ".r-card__photo" and set the variables on <html>. */
[data-theme='dark'] .r-card__photo {
    --thumb-from: hsl(var(--h) 22% 19%);
    --thumb-to: hsl(calc(var(--h) + 20) 28% 27%);
    --plate: rgba(248, 248, 255, 0.07);
    --plate-shadow: none;
    --plate-dot: hsl(var(--h) 30% 36%);
    --prep-pill: rgba(42, 34, 19, 0.95);
}
@media (prefers-color-scheme: dark) {
    :root:not([data-theme='light']):not([data-theme='dark']) .r-card__photo {
        --thumb-from: hsl(var(--h) 22% 19%);
        --thumb-to: hsl(calc(var(--h) + 20) 28% 27%);
        --plate: rgba(248, 248, 255, 0.07);
        --plate-shadow: none;
        --plate-dot: hsl(var(--h) 30% 36%);
        --prep-pill: rgba(42, 34, 19, 0.95);
    }
}
.r-card__photo img {
    display: block; /* an inline img added a text line, so the photo was 1 px taller than 1:1 */
    width: 100%;
    height: 100%;
    object-fit: cover;
}
/* the plate: a 52 % circle with a 28 % circle in the middle */
.r-card__plate {
    position: absolute;
    inset: 24%;
    border-radius: 50%;
    background: var(--plate);
    box-shadow: var(--plate-shadow);
}
.r-card__plate::after {
    content: '';
    position: absolute;
    inset: calc(100% * 12 / 52); /* (52 % - 28 %) / 2 of the frame, relative to the plate */
    border-radius: 50%;
    background: var(--plate-dot);
}
.r-card__prep {
    position: absolute;
    left: 10px;
    bottom: 10px;
    display: inline-flex;
    align-items: center;
    gap: 5px;
    padding: 4px 10px;
    border-radius: var(--r-pill);
    background: var(--prep-pill);
    color: var(--prep);
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
    transition: color var(--dur-fast) var(--ease);
    flex: 1;
    /* two lines reserved: the card has the same height as the skeleton and its neighbors */
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
    line-height: 20px; /* fixed: the mono font in kcal doesn't raise the line; the skeleton uses the same */
    color: var(--ink-muted);
}
.r-card__kcal {
    font-family: var(--font-mono);
}
</style>
