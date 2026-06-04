<script lang="ts" setup>
import { computed, ref } from 'vue';
import { useAuth0 } from '@auth0/auth0-vue';
import { useRouter } from 'vue-router';

import { fetchRecipes } from '@/data/api/recipe/fetch';
import type { RecipeDto } from '@/data/dtos/recipe/RecipeDto';
import type { RecipeIngredientDto } from '@/data/dtos/recipe/RecipeIngredientDto';

const props = defineProps<{ id: number }>();
const { isAuthenticated } = useAuth0();
const router = useRouter();

const recipe = ref<RecipeDto>();
const servings = ref(1);
const checkedSteps = ref<Set<number>>(new Set());
const checkedIngs = ref<Set<string>>(new Set());

fetchRecipes(props.id).then((v) => {
    recipe.value = v;
    servings.value = v.servings;
});

const diffLabel: Record<number, string> = { 1: 'Łatwe', 2: 'Średnie', 3: 'Trudne' };

function fmtTime(min: number | null) {
    if (!min) return null;
    if (min < 60) return `${min} min`;
    const h = Math.floor(min / 60);
    const m = min % 60;
    return m ? `${h}h ${m}min` : `${h}h`;
}

function scaleAmount(ing: RecipeIngredientDto) {
    if (!recipe.value || !ing.amount) return ing.amountText ?? '—';
    const ratio = servings.value / recipe.value.servings;
    const scaled = ing.amount * ratio;
    const fmt = Math.abs(scaled - Math.round(scaled)) < 0.05
        ? String(Math.round(scaled))
        : scaled < 10 ? scaled.toFixed(1) : String(Math.round(scaled));
    return `${fmt}${ing.unit ? ' ' + ing.unit : ''}`;
}

const stepsBySection = computed(() => {
    if (!recipe.value) return [];
    const sections: { name: string | null; steps: typeof recipe.value.steps }[] = [];
    for (const step of recipe.value.steps) {
        const last = sections[sections.length - 1];
        if (last && last.name === step.section) last.steps.push(step);
        else sections.push({ name: step.section, steps: [step] });
    }
    return sections;
});

const ingredientsBySection = computed(() => {
    if (!recipe.value) return [];
    const sections: { name: string | null; items: typeof recipe.value.ingredients }[] = [];
    for (const ing of recipe.value.ingredients) {
        const last = sections[sections.length - 1];
        if (last && last.name === ing.section) last.items.push(ing);
        else sections.push({ name: ing.section, items: [ing] });
    }
    return sections;
});

function toggleStep(no: number) {
    if (checkedSteps.value.has(no)) checkedSteps.value.delete(no);
    else checkedSteps.value.add(no);
}

function toggleIng(key: string) {
    if (checkedIngs.value.has(key)) checkedIngs.value.delete(key);
    else checkedIngs.value.add(key);
}

const totalSteps = computed(() => recipe.value?.steps.length ?? 0);
const doneSteps = computed(() => checkedSteps.value.size);
const progress = computed(() => totalSteps.value ? (doneSteps.value / totalSteps.value) * 100 : 0);
const allDone = computed(() => recipe.value ? doneSteps.value === totalSteps.value : false);
</script>

<template>
    <div v-if="!recipe" style="padding: 80px; text-align: center;" class="muted t-small">Ładowanie…</div>

    <article v-else class="container recipe-page">
        <!-- Breadcrumb -->
        <nav class="breadcrumb">
            <RouterLink to="/recipes" class="breadcrumb__link">przepisy</RouterLink>
            <span class="breadcrumb__sep">/</span>
            <span class="breadcrumb__link" v-if="recipe.category">{{ recipe.category }}</span>
            <span class="breadcrumb__sep" v-if="recipe.category">/</span>
            <span class="breadcrumb__current">#{{ String(recipe.id).padStart(3, '0') }}</span>
        </nav>

        <!-- Hero: text left, photo right -->
        <header class="recipe-hero">
            <div class="recipe-hero__text">
                <div style="display: flex; align-items: center; gap: 10px; flex-wrap: wrap; margin-bottom: 24px;">
                    <span v-if="recipe.category" class="badge" :class="`cat-${recipe.category}`">
                        {{ recipe.category }}<template v-if="recipe.region"> · {{ recipe.region }}</template>
                    </span>
                </div>

                <h1 class="recipe-hero__title">{{ recipe.name }}</h1>

                <p v-if="recipe.tagline" class="recipe-hero__tagline">{{ recipe.tagline }}</p>

                <div class="recipe-meta-strip">
                    <div v-if="recipe.timeMinutes" class="meta-item">
                        <span class="meta-item__label">Czas</span>
                        <span class="meta-item__value">{{ fmtTime(recipe.timeMinutes) }}</span>
                    </div>
                    <div v-if="recipe.difficulty" class="meta-item">
                        <span class="meta-item__label">Trudność</span>
                        <span class="meta-item__value">
                            <span class="diff-dots">
                                <span
                                    v-for="i in 3"
                                    :key="i"
                                    class="diff-dot"
                                    :class="{ 'is-on': i <= recipe.difficulty }"
                                />
                            </span>
                            {{ diffLabel[recipe.difficulty] }}
                        </span>
                    </div>
                    <div class="meta-item">
                        <span class="meta-item__label">Porcje</span>
                        <span class="meta-item__value">{{ recipe.servings }}</span>
                    </div>
                    <div v-if="recipe.kcalPerServing" class="meta-item">
                        <span class="meta-item__label">Na porcję</span>
                        <span class="meta-item__value">{{ recipe.kcalPerServing }} kcal</span>
                    </div>
                </div>

                <RouterLink
                    v-if="isAuthenticated"
                    :to="{ name: 'recipe-edit', params: { id: recipe.id } }"
                    class="btn btn--secondary btn--sm"
                    style="margin-top: 20px;"
                >Edytuj przepis</RouterLink>
            </div>

            <div class="recipe-hero__photo">
                <div class="photo photo--hero" :style="`--hue: 220`">
                    <img
                        v-if="recipe.photoUrl"
                        :src="recipe.photoUrl"
                        :alt="recipe.name"
                        style="position:absolute;inset:0;width:100%;height:100%;object-fit:cover;"
                    />
                    <div class="photo--hero__label">// foto · 4:5</div>
                </div>
            </div>
        </header>

        <!-- Body: ingredients left, steps right -->
        <div class="recipe-body">
            <!-- Ingredients sidebar -->
            <aside class="recipe-ingredients">
                <div class="card" style="padding: 24px; display: flex; flex-direction: column; gap: 18px;">
                    <div>
                        <div class="ingredients-kicker">// składniki</div>
                        <h2 class="ingredients-title">Czego potrzebujesz</h2>
                    </div>

                    <!-- Serving stepper -->
                    <div class="servings-control">
                        <div style="display: flex; flex-direction: column; gap: 2px;">
                            <span class="servings-control__label">Porcje</span>
                            <span class="servings-control__hint">ilości skalują się automatycznie</span>
                        </div>
                        <div style="display: flex; align-items: center; gap: 4px;">
                            <button class="stepper-btn" :disabled="servings <= 1" @click="servings = Math.max(1, servings - 1)">−</button>
                            <span class="servings-control__value">{{ servings }}</span>
                            <button class="stepper-btn" :disabled="servings >= 20" @click="servings++">+</button>
                        </div>
                    </div>

                    <!-- Ingredient groups -->
                    <div style="display: flex; flex-direction: column; gap: 22px;">
                        <template v-for="section in ingredientsBySection" :key="section.name">
                            <div>
                                <div v-if="section.name" class="ing-section-label">{{ section.name }}</div>
                                <ul style="list-style: none; padding: 0; margin: 0; display: flex; flex-direction: column; gap: 2px;">
                                    <li v-for="(ing, ii) in section.items" :key="ing.sortOrder">
                                        <button
                                            class="ing-row"
                                            :class="{ 'is-checked': checkedIngs.has(`${section.name}-${ii}`) }"
                                            @click="toggleIng(`${section.name}-${ii}`)"
                                        >
                                            <span class="ing-row__check">
                                                <svg v-if="checkedIngs.has(`${section.name}-${ii}`)" width="8" height="8" viewBox="0 0 10 10" fill="none">
                                                    <path d="M2 5.2L4.2 7.4L8 3.4" stroke="#fff" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"/>
                                                </svg>
                                            </span>
                                            <span class="ing-row__qty mono">{{ scaleAmount(ing) }}</span>
                                            <span class="ing-row__name">
                                                {{ ing.ingredientName }}
                                                <span v-if="ing.note" class="ing-row__note">{{ ing.note }}</span>
                                            </span>
                                        </button>
                                    </li>
                                </ul>
                            </div>
                        </template>
                    </div>
                </div>
            </aside>

            <!-- Steps main -->
            <main class="recipe-steps">
                <div class="steps-header">
                    <div>
                        <div class="steps-kicker">// kroki</div>
                        <h2 class="steps-title">Jak to zrobić</h2>
                    </div>
                    <span class="steps-progress-label">{{ doneSteps }}/{{ totalSteps }} zrobione</span>
                </div>
                <!-- Progress bar -->
                <div class="steps-progress-bar">
                    <div class="steps-progress-bar__fill" :style="`width: ${progress}%`" />
                </div>

                <template v-for="section in stepsBySection" :key="section.name">
                    <p v-if="section.name" class="step-section-label">{{ section.name }}</p>
                    <ol style="list-style: none; padding: 0; margin: 0 0 20px; display: flex; flex-direction: column; gap: 14px;">
                        <li
                            v-for="step in section.steps"
                            :key="step.stepNo"
                            class="step-card"
                            :class="{ 'is-done': checkedSteps.has(step.stepNo) }"
                        >
                            <button
                                class="step-num-btn"
                                :class="{ 'is-done': checkedSteps.has(step.stepNo) }"
                                @click="toggleStep(step.stepNo)"
                            >
                                <svg v-if="checkedSteps.has(step.stepNo)" width="22" height="22" viewBox="0 0 22 22" fill="none">
                                    <path d="M5 11.5L9.5 16L17 7" stroke="#fff" stroke-width="2.4" stroke-linecap="round" stroke-linejoin="round"/>
                                </svg>
                                <span v-else>{{ step.stepNo }}</span>
                            </button>
                            <p class="step-text">{{ step.text }}</p>
                        </li>
                    </ol>
                </template>

                <div v-if="allDone" class="done-block done-block--done">
                    <div>
                        <div class="done-block__kicker">// gotowe ✓</div>
                        <div class="done-block__msg">Smacznego! Daj znać jak wyszło.</div>
                    </div>
                    <div style="display: flex; gap: 8px;">
                        <button class="btn btn--secondary">★ Oceń przepis</button>
                        <button class="btn">Dodaj komentarz</button>
                    </div>
                </div>
                <div v-else class="done-block">
                    <div>
                        <div class="done-block__kicker done-block__kicker--muted">// po przepisie</div>
                        <div class="done-block__msg">Smacznego!</div>
                    </div>
                </div>
            </main>
        </div>
    </article>
</template>

<style scoped>
.recipe-page {
    padding-bottom: 80px;
}

/* Breadcrumb */
.breadcrumb {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 12px;
    color: var(--ink-muted);
    padding-top: 28px;
    padding-bottom: 28px;
    font-family: var(--font-mono);
    letter-spacing: 0.3px;
}
.breadcrumb__link {
    color: var(--ink-muted);
    text-decoration: none;
}
.breadcrumb__link:hover { color: var(--ink); }
.breadcrumb__sep { color: var(--ink-faint); }
.breadcrumb__current { color: var(--ink); font-weight: 600; }

/* Hero */
.recipe-hero {
    display: grid;
    grid-template-columns: 1fr 460px;
    gap: 56px;
    align-items: end;
    padding-bottom: 8px;
}
.recipe-hero__title {
    font-size: clamp(48px, 6.5vw, 84px);
    line-height: 0.95;
    letter-spacing: -3px;
    font-weight: 500;
    margin: 0 0 22px;
    text-wrap: balance;
}
.recipe-hero__tagline {
    font-size: 19px;
    line-height: 1.5;
    color: var(--ink-dim);
    max-width: 580px;
    text-wrap: pretty;
    font-weight: 400;
    margin: 0 0 36px;
}

/* Photo */
.photo--hero {
    aspect-ratio: 4 / 5 !important;
    border-radius: var(--r-2xl);
    background: radial-gradient(
        circle at 30% 20%,
        hsl(var(--hue, 220) 35% 88%) 0%,
        hsl(calc(var(--hue, 220) + 20) 50% 68%) 55%,
        hsl(calc(var(--hue, 220) + 40) 45% 60%) 100%
    );
}
.photo--hero__label {
    position: absolute;
    bottom: 14px;
    left: 14px;
    font-family: var(--font-mono);
    font-size: 10px;
    color: rgba(10, 10, 15, 0.45);
    letter-spacing: 0.5px;
    text-transform: uppercase;
    background: rgba(255,255,255,0.55);
    padding: 4px 8px;
    border-radius: var(--r-sm);
    backdrop-filter: blur(8px);
}

/* Meta strip */
.recipe-meta-strip {
    display: flex;
    gap: 32px;
    flex-wrap: wrap;
    padding-top: 24px;
    border-top: 1px solid var(--rule);
}
.meta-item {
    display: flex;
    flex-direction: column;
    gap: 4px;
    min-width: 80px;
}
.meta-item__label {
    font-family: var(--font-mono);
    font-size: 10px;
    color: var(--ink-muted);
    letter-spacing: 0.6px;
    text-transform: uppercase;
    font-weight: 600;
}
.meta-item__value {
    font-size: 17px;
    font-weight: 600;
    letter-spacing: -0.3px;
    color: var(--ink);
    display: flex;
    align-items: center;
    gap: 8px;
}

/* Diff dots */
.diff-dots { display: inline-flex; gap: 4px; }
.diff-dot {
    width: 8px; height: 8px; border-radius: var(--r-pill);
    border: 1.5px solid var(--rule-strong);
    background: transparent;
}
.diff-dot.is-on { background: var(--accent); border-color: var(--accent); }

/* Body layout */
.recipe-body {
    display: grid;
    grid-template-columns: 340px 1fr;
    gap: 56px;
    align-items: start;
    margin-top: 56px;
}

/* Ingredients */
.recipe-ingredients { position: sticky; top: 90px; }
.ingredients-kicker {
    font-family: var(--font-mono);
    font-size: 10px;
    color: var(--accent);
    letter-spacing: 0.6px;
    text-transform: uppercase;
    font-weight: 600;
    margin-bottom: 8px;
}
.ingredients-title {
    font-size: 26px;
    font-weight: 500;
    letter-spacing: -0.8px;
    margin: 0;
}

/* Servings */
.servings-control {
    background: var(--bg);
    border: 1px solid var(--rule);
    border-radius: var(--r-lg);
    padding: 10px 12px;
    display: flex;
    align-items: center;
    justify-content: space-between;
}
.servings-control__label {
    font-family: var(--font-mono);
    font-size: 10px;
    color: var(--ink-muted);
    letter-spacing: 0.6px;
    text-transform: uppercase;
    font-weight: 600;
}
.servings-control__hint {
    font-family: var(--font-mono);
    font-size: 10px;
    color: var(--ink-faint);
    letter-spacing: 0.3px;
}
.servings-control__value {
    min-width: 36px;
    text-align: center;
    font-size: 22px;
    font-weight: 600;
    letter-spacing: -0.4px;
}
.stepper-btn {
    width: 30px; height: 30px;
    border-radius: var(--r-md);
    background: var(--bg-alt);
    border: 1px solid var(--rule);
    color: var(--ink);
    font-size: 18px;
    line-height: 1;
    cursor: pointer;
    font-weight: 500;
}
.stepper-btn:disabled { color: var(--ink-faint); background: transparent; cursor: not-allowed; }

/* Ingredient section label */
.ing-section-label {
    font-family: var(--font-mono);
    font-size: 10px;
    color: var(--ink-muted);
    letter-spacing: 0.6px;
    text-transform: uppercase;
    font-weight: 600;
    padding-bottom: 8px;
    margin-bottom: 8px;
    border-bottom: 1px solid var(--rule);
}

/* Ingredient row */
.ing-row {
    width: 100%;
    display: grid;
    grid-template-columns: 16px 78px 1fr;
    gap: 10px;
    align-items: baseline;
    background: transparent;
    border: none;
    padding: 7px 0;
    cursor: pointer;
    text-align: left;
    border-bottom: 1px solid var(--rule);
}
.ing-row__check {
    width: 14px; height: 14px;
    margin-top: 4px;
    border-radius: var(--r-sm);
    border: 1.5px solid var(--rule-strong);
    background: transparent;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
    transition: all 0.12s;
}
.ing-row.is-checked .ing-row__check {
    background: var(--accent);
    border-color: var(--accent);
}
.ing-row__qty {
    font-family: var(--font-mono);
    font-size: 13px;
    font-weight: 600;
    color: var(--ink);
    letter-spacing: 0.2px;
    white-space: nowrap;
}
.ing-row.is-checked .ing-row__qty,
.ing-row.is-checked .ing-row__name {
    color: var(--ink-faint);
    text-decoration: line-through;
}
.ing-row__name {
    font-size: 14px;
    color: var(--ink-dim);
    line-height: 1.4;
}
.ing-row__note {
    font-family: var(--font-mono);
    font-size: 10px;
    color: var(--ink-faint);
    margin-left: 6px;
    letter-spacing: 0.3px;
}

/* Steps */
.steps-header {
    display: flex;
    align-items: flex-end;
    justify-content: space-between;
    gap: 16px;
    flex-wrap: wrap;
    margin-bottom: 18px;
}
.steps-kicker {
    font-family: var(--font-mono);
    font-size: 10px;
    color: var(--accent);
    letter-spacing: 0.6px;
    text-transform: uppercase;
    font-weight: 600;
    margin-bottom: 8px;
}
.steps-title {
    font-size: 32px;
    font-weight: 500;
    letter-spacing: -1px;
    margin: 0;
}
.steps-progress-label {
    font-family: var(--font-mono);
    font-size: 11px;
    color: var(--ink-muted);
    letter-spacing: 0.3px;
    padding-bottom: 4px;
}
.steps-progress-bar {
    height: 3px;
    background: var(--rule);
    border-radius: var(--r-pill);
    overflow: hidden;
    margin-bottom: 22px;
}
.steps-progress-bar__fill {
    height: 100%;
    background: var(--accent);
    transition: width 0.25s ease;
}

/* Step card */
.step-card {
    background: var(--bg-alt);
    border: 1px solid var(--rule);
    border-radius: var(--r-xl);
    padding: 22px 24px;
    display: grid;
    grid-template-columns: 64px 1fr;
    gap: 22px;
    align-items: start;
    transition: all 0.2s;
}
.step-card.is-done {
    opacity: 0.45;
    background: var(--bg);
}

.step-num-btn {
    width: 56px; height: 56px;
    border-radius: 12px;
    background: var(--bg);
    border: 1.5px solid var(--rule-strong);
    color: var(--accent);
    font-size: 24px;
    font-weight: 600;
    letter-spacing: -0.5px;
    cursor: pointer;
    font-family: var(--font-sans);
    padding: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.15s;
}
.step-num-btn:hover { background: var(--accent-soft); }
.step-num-btn.is-done {
    background: var(--accent);
    border-color: var(--accent);
    color: #fff;
}

.step-text {
    font-size: 17px;
    line-height: 1.55;
    color: var(--ink);
    text-wrap: pretty;
    padding-top: 8px;
    margin: 0;
}
.step-card.is-done .step-text {
    color: var(--ink-muted);
    text-decoration: line-through;
    text-decoration-color: var(--ink-faint);
}

/* Step section label */
.step-section-label {
    font-family: var(--font-mono);
    font-size: 10px;
    color: var(--accent);
    letter-spacing: 0.6px;
    text-transform: uppercase;
    font-weight: 600;
    margin: 0 0 12px;
}

/* Done block */
.done-block {
    margin-top: 32px;
    padding: 32px 28px;
    border: 1px dashed var(--rule-strong);
    border-radius: var(--r-xl);
    background: var(--bg-alt);
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 24px;
    flex-wrap: wrap;
    transition: all 0.3s;
}
.done-block--done {
    background: color-mix(in srgb, var(--c-sage) 12%, transparent);
    border-style: solid;
    border-color: var(--c-sage);
}
.done-block__kicker {
    font-family: var(--font-mono);
    font-size: 10px;
    color: var(--c-forest);
    letter-spacing: 0.6px;
    text-transform: uppercase;
    font-weight: 600;
    margin-bottom: 6px;
}
.done-block__kicker--muted { color: var(--ink-muted); }
.done-block__msg {
    font-size: 22px;
    font-weight: 500;
    letter-spacing: -0.6px;
}

@media (max-width: 1100px) {
    .recipe-hero { grid-template-columns: 1fr 340px; gap: 40px; }
    .recipe-body  { grid-template-columns: 300px 1fr; gap: 40px; }
}

@media (max-width: 900px) {
    .recipe-hero { grid-template-columns: 1fr; }
    .recipe-hero__photo { order: -1; }
    .photo--hero { aspect-ratio: 16/9 !important; }
    .recipe-body { grid-template-columns: 1fr; }
    .recipe-ingredients { position: static; }
}
</style>
