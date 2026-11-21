<script lang="ts" setup>
import { computed } from 'vue';

import { DIFFICULTY, prepRows } from '../recipeView';

import UsefulRecipes from './UsefulRecipes.vue';

import CautionsBlock from '@/components/CautionsBlock.vue';
import USpinner from '@/components/USpinner.vue';
import type { RecipeDto } from '@/data/dtos/recipe/RecipeDto';
import { avg, fmt, type IRecipeNutrition, nf } from '@/shared/utils/nutrition';
import { formatTime } from '@/shared/utils/time';

// recipe null = loading: skeletons in place of the content.
// Without a photo the hero is one column; on a phone the title is 38 px,
// so the ingredients start on the first screen.
const props = defineProps<{
    recipe: RecipeDto | null;
    servings: number;
    /** Macros of the chosen variant (substitutes); none = the recipe's macros. */
    nutrition?: IRecipeNutrition | null;
}>();

const dish = computed(() => props.recipe?.taxonomy.dish ?? []);
const cuisine = computed(() => props.recipe?.taxonomy.cuisine?.[0]?.label);
const diets = computed(() => props.recipe?.taxonomy.diet ?? []);
const practical = computed(() => props.recipe?.taxonomy.practical ?? []);
// "prepare ahead" steps: the first items of the Uwagi block
const prep = computed(() => prepRows(props.recipe?.steps ?? []));
const kcal = computed(() => {
    const n = props.nutrition !== undefined ? props.nutrition : props.recipe?.nutrition;
    const f = fmt(n?.perPortion.kcal, 'kcal', 'recipe');
    return f ? nf(avg(f)) : null;
});

function scrollToStep(no: number) {
    const el = document.getElementById(`step-${no}`);
    if (!el) return;
    window.scrollTo({
        top: el.getBoundingClientRect().top + window.scrollY - 100,
        behavior: 'smooth',
    });
    el.querySelector<HTMLElement>('button')?.focus({ preventScroll: true });
}
</script>

<template>
    <div
        v-if="!recipe"
        class="hero hero--text"
        aria-busy="true"
    >
        <div
            class="hero__loading"
            role="status"
        >
            <USpinner :size="20" />
            <span>ładowanie przepisu…</span>
        </div>
        <div
            class="hero__text"
            aria-hidden="true"
        >
            <div class="hero__flags">
                <span class="badge hero__sk-badge">
                    <span class="sk hero__sk-pill" />
                    &nbsp;
                </span>
            </div>
            <div class="hero__title"><span class="sk hero__sk-title">&nbsp;</span></div>
            <dl class="hero__meta">
                <div>
                    <dt>Czas</dt>
                    <dd><span class="sk hero__sk-dd">&nbsp;</span></dd>
                </div>
                <div>
                    <dt>Trudność</dt>
                    <dd><span class="sk hero__sk-dd">&nbsp;</span></dd>
                </div>
                <div>
                    <dt>Porcje</dt>
                    <dd><span class="sk hero__sk-dd hero__sk-dd--short">&nbsp;</span></dd>
                </div>
                <div>
                    <dt>Kcal / porcja</dt>
                    <dd><span class="sk hero__sk-dd">&nbsp;</span></dd>
                </div>
            </dl>
        </div>
    </div>
    <div
        v-else
        class="hero"
        :class="{ 'hero--text': !recipe.photoUrl }"
    >
        <div
            v-if="recipe.photoUrl"
            class="hero__photo"
        >
            <img
                :src="recipe.photoUrl"
                :alt="`Zdjęcie: ${recipe.name}`"
            />
        </div>
        <div class="hero__text">
            <div class="hero__flags">
                <span
                    v-for="d in dish"
                    :key="d.id"
                    class="badge"
                    :class="d.color ? `cat-${d.color}` : null"
                >
                    {{ d.label }}
                </span>
                <span
                    v-if="cuisine"
                    class="hero__cuisine"
                >
                    kuchnia {{ cuisine }}
                </span>
            </div>
            <h1 class="hero__title">{{ recipe.name }}</h1>
            <p
                v-if="recipe.tagline"
                class="hero__tagline"
            >
                {{ recipe.tagline }}
            </p>
            <dl class="hero__meta">
                <div v-if="recipe.timeMin">
                    <dt>Czas</dt>
                    <dd>{{ formatTime(recipe.timeMin, recipe.timeMax) }}</dd>
                </div>
                <div v-if="recipe.difficulty">
                    <dt>Trudność</dt>
                    <dd class="hero__diff">
                        <span
                            class="hero__pips"
                            aria-hidden="true"
                        >
                            <span
                                v-for="i in 3"
                                :key="i"
                                class="hero__pip"
                                :class="{ 'is-on': i <= recipe.difficulty }"
                            />
                        </span>
                        {{ DIFFICULTY[recipe.difficulty] }}
                    </dd>
                </div>
                <div>
                    <dt>Porcje</dt>
                    <dd>{{ servings }}</dd>
                </div>
                <div v-if="kcal">
                    <dt>Kcal / porcja</dt>
                    <dd>
                        <a
                            href="#makro"
                            class="hero__kcal"
                        >
                            ~{{ kcal }}
                        </a>
                        <span class="hero__kcal-note">orientacyjnie</span>
                    </dd>
                </div>
            </dl>
            <CautionsBlock
                class="hero__cautions"
                :items="recipe.cautions"
                :prep="prep"
                @step="scrollToStep"
            />
            <ul
                v-if="diets.length || practical.length"
                class="hero__tags"
                aria-label="Cechy przepisu"
            >
                <li
                    v-for="d in diets"
                    :key="d.id"
                    class="hero__tag hero__tag--diet"
                >
                    {{ d.label }}
                </li>
                <li
                    v-for="p in practical"
                    :key="p.id"
                    class="hero__tag"
                >
                    {{ p.label }}
                </li>
            </ul>
            <UsefulRecipes :ingredients="recipe.ingredients" />
        </div>
    </div>
</template>

<style scoped>
.hero {
    display: grid;
    grid-template-columns: minmax(0, 1fr) 440px;
    gap: 56px;
    align-items: end;
    padding-top: 28px;
}
.hero__photo {
    order: 2;
    aspect-ratio: 1 / 1;
    border-radius: var(--r-2xl);
    border: 1px solid var(--rule);
    overflow: hidden;
}
.hero__photo img {
    width: 100%;
    height: 100%;
    object-fit: cover;
}
.hero__flags {
    display: flex;
    align-items: center;
    gap: 8px;
    flex-wrap: wrap;
    margin-bottom: 24px;
}
.hero__cuisine {
    font-size: 13px;
    color: var(--ink-dim);
    font-weight: 500;
}
.hero__title {
    font-size: clamp(48px, 6vw, 80px);
    line-height: 0.98;
    letter-spacing: -2.8px;
    font-weight: 500;
    text-wrap: balance;
}
.hero__tagline {
    margin-top: 22px;
    font-size: 19px;
    line-height: 1.5;
    color: var(--ink-dim);
    max-width: 580px;
    text-wrap: pretty;
}
.hero__meta {
    margin: 32px 0 0;
    padding-top: 20px;
    border-top: 1px solid var(--rule);
    display: flex;
    gap: 36px;
    flex-wrap: wrap;
}
.hero__meta div {
    display: flex;
    flex-direction: column;
    gap: 4px;
}
.hero__meta dt {
    font-family: var(--font-mono);
    font-size: 12px;
    color: var(--ink-muted);
    letter-spacing: 0.5px;
    text-transform: uppercase;
    font-weight: 600;
}
.hero__meta dd {
    margin: 0;
    font-size: 17px;
    font-weight: 600;
    letter-spacing: -0.3px;
}
.hero__kcal-note {
    margin-left: 6px;
    font-size: 12px;
    font-weight: 500;
    letter-spacing: 0;
    color: var(--ink-muted);
}
.hero__kcal {
    color: inherit;
    font-family: var(--font-mono);
    text-decoration: underline;
    text-decoration-color: var(--rule-strong);
    text-underline-offset: 4px;
}
.hero__kcal:hover {
    text-decoration-color: currentColor;
}
.hero__diff {
    display: inline-flex;
    align-items: center;
    gap: 8px;
}
.hero__pips {
    display: inline-flex;
    gap: 4px;
}
.hero__pip {
    width: 8px;
    height: 8px;
    border-radius: var(--r-pill);
    border: 1.5px solid var(--rule-strong);
}
.hero__pip.is-on {
    background: var(--accent);
    border-color: var(--accent);
}
.hero__cautions {
    margin-top: 20px;
    max-width: 580px;
}
.hero__tags {
    list-style: none;
    padding: 0;
    margin: 20px 0 0;
    display: flex;
    flex-wrap: wrap;
    gap: 6px;
}
.hero__tag {
    padding: 4px 12px;
    border-radius: var(--r-pill);
    font-size: 13px;
    color: var(--ink-dim);
    border: 1px solid var(--rule);
}
.hero__tag--diet {
    color: var(--ink);
    font-weight: 600;
    border: 1.5px solid var(--rule-strong);
}

.hero--text {
    display: block;
    max-width: 760px;
}
.hero--text .hero__flags {
    margin-bottom: 20px;
}
.hero--text .hero__title {
    font-size: clamp(48px, 6vw, 72px);
    line-height: 1;
    letter-spacing: -2.6px;
}
.hero--text .hero__tagline {
    margin-top: 18px;
}
.hero--text .hero__meta {
    margin-top: 28px;
}

.hero__loading {
    display: flex;
    align-items: center;
    gap: 10px;
    margin-bottom: 16px;
    font-size: 14px;
    color: var(--ink-muted);
}
.hero__sk-badge {
    position: relative;
    width: 110px;
    color: transparent;
    background: transparent;
    border-color: transparent;
}
.hero__sk-badge::before {
    visibility: hidden;
}
.hero__sk-pill {
    position: absolute;
    inset: 0;
    border-radius: var(--r-pill);
}
.hero__sk-title {
    display: inline-block;
    width: 70%;
    border-radius: var(--r-lg);
}
.hero__sk-dd {
    display: inline-block;
    width: 64px;
}
.hero__sk-dd--short {
    width: 24px;
}

@media (max-width: 900px) {
    .hero {
        grid-template-columns: minmax(0, 1fr);
        gap: 20px;
        padding-top: 16px;
    }
    .hero__photo {
        order: 0;
    }
    .hero__flags {
        margin-bottom: 16px;
    }
    .hero__title {
        font-size: 40px;
        letter-spacing: -1.4px;
    }
    .hero__tagline {
        margin-top: 14px;
        font-size: 17px;
    }
    .hero__meta {
        margin-top: 24px;
        gap: 24px;
    }
    .hero--text {
        padding-top: 14px;
    }
    .hero--text .hero__flags {
        margin-bottom: 12px;
    }
    .hero--text .hero__title {
        font-size: 38px;
        line-height: 1;
        letter-spacing: -1.3px;
    }
    .hero--text .hero__tagline {
        margin-top: 12px;
        font-size: 16px;
    }
    .hero--text .hero__meta {
        margin-top: 18px;
        padding-top: 14px;
        gap: 14px 22px;
    }
    .hero--text .hero__cautions,
    .hero--text .hero__tags {
        margin-top: 14px;
    }
}
</style>
