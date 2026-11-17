<script lang="ts" setup>
import { computed } from 'vue';

import { DIFFICULTY, leadText } from '../recipeView';

import UsefulRecipes from './UsefulRecipes.vue';

import USpinner from '@/components/USpinner.vue';
import type { RecipeDto } from '@/data/dtos/recipe/RecipeDto';
import { avg, fmt, nf } from '@/shared/utils/nutrition';
import { formatTime } from '@/shared/utils/time';

// recipe null = ładowanie: ten sam .hero (zdjęcie 1:1, metadane) ze szkieletami w miejscu treści.
const props = defineProps<{ recipe: RecipeDto | null; servings: number }>();

const dish = computed(() => props.recipe?.taxonomy.dish ?? []);
const cuisine = computed(() => props.recipe?.taxonomy.cuisine?.[0]?.label);
const diets = computed(() => props.recipe?.taxonomy.diet ?? []);
const practical = computed(() => props.recipe?.taxonomy.practical ?? []);
const prepSteps = computed(() => props.recipe?.steps.filter((s) => s.prep) ?? []);
const kcal = computed(() => {
    const f = fmt(props.recipe?.nutrition?.perPortion.kcal, 'kcal', 'recipe');
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
        class="hero"
        aria-busy="true"
    >
        <div class="hero__photo hero__photo--sk">
            <div class="sk hero__sk-fill" />
            <div
                class="hero__loading"
                role="status"
            >
                <USpinner :size="28" />
                <span>ładowanie przepisu…</span>
            </div>
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
    >
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
.hero__photo img,
.hero__placeholder,
.hero__sk-fill {
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

/* szkielet: te same klasy co treść, wypełnienia .sk w miejscu tekstu */
.hero__photo--sk {
    position: relative;
}
.hero__sk-fill {
    border-radius: 0;
}
.hero__loading {
    position: absolute;
    inset: 0;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 12px;
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
}
</style>
