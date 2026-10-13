<script lang="ts" setup>
import { computed } from 'vue';

import { DIFFICULTY, fmtTime, leadText } from '../recipeView';

import type { RecipeDto } from '@/data/dtos/recipe/RecipeDto';

const props = defineProps<{ recipe: RecipeDto; servings: number }>();

const dish = computed(() => props.recipe.taxonomy.dish ?? []);
const cuisine = computed(() => props.recipe.taxonomy.cuisine?.[0]?.label);
const diets = computed(() => props.recipe.taxonomy.diet ?? []);
const practical = computed(() => props.recipe.taxonomy.practical ?? []);
const prepSteps = computed(() => props.recipe.steps.filter((s) => s.prep));

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
    <div class="hero">
        <div class="hero__photo">
            <img
                v-if="recipe.photoUrl"
                :src="recipe.photoUrl"
                :alt="`Zdjęcie: ${recipe.name}`"
            />
            <div
                v-else
                class="photo hero__placeholder"
                aria-hidden="true"
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
                <div v-if="recipe.timeMinutes">
                    <dt>Czas</dt>
                    <dd>{{ fmtTime(recipe.timeMinutes) }}</dd>
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
                <div v-if="recipe.kcalPerServing">
                    <dt>Kcal / porcja</dt>
                    <dd>{{ recipe.kcalPerServing }}</dd>
                </div>
            </dl>
            <div
                v-if="prepSteps.length"
                class="hero__prep"
            >
                <svg
                    width="14"
                    height="14"
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
                <div>
                    <strong>Wymaga przygotowania</strong>
                    <div
                        v-for="s in prepSteps"
                        :key="s.stepNo"
                        class="hero__prep-line"
                    >
                        <template v-if="leadText(s)">
                            {{ leadText(s)!.charAt(0).toUpperCase() + leadText(s)!.slice(1) }}:
                        </template>
                        <button
                            type="button"
                            class="hero__prep-link"
                            @click="scrollToStep(s.stepNo)"
                        >
                            krok {{ s.stepNo }}
                        </button>
                    </div>
                </div>
            </div>
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
.hero__photo img,
.hero__placeholder {
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
.hero__prep {
    margin-top: 20px;
    display: flex;
    align-items: flex-start;
    gap: 12px;
    padding: 12px 14px;
    background: var(--prep-soft);
    border: 1px solid var(--prep-line);
    border-radius: var(--r-lg);
    max-width: 580px;
    font-size: 14px;
    line-height: 1.5;
    color: var(--prep);
}
.hero__prep svg {
    margin-top: 3px;
    flex-shrink: 0;
}
.hero__prep strong {
    font-weight: 600;
}
.hero__prep-line {
    color: var(--ink-dim);
}
.hero__prep-link {
    background: none;
    border: none;
    padding: 0;
    color: var(--ink);
    text-decoration: underline;
    text-underline-offset: 3px;
    cursor: pointer;
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
}
</style>
