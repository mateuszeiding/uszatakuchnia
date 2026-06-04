<script lang="ts" setup>
import { computed, ref } from 'vue';
import { useAuth0 } from '@auth0/auth0-vue';

import { fetchRecipes } from '@/data/api/recipe/fetch';
import type { RecipeDto } from '@/data/dtos/recipe/RecipeDto';
import type { RecipeIngredientDto } from '@/data/dtos/recipe/RecipeIngredientDto';

const props = defineProps<{ id: number }>();
const { isAuthenticated } = useAuth0();

const recipe = ref<RecipeDto>();
const servings = ref(1);
const checkedSteps = ref<Set<number>>(new Set());

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
    const fmt = scaled % 1 === 0 ? scaled.toString() : scaled.toFixed(1);
    return `${fmt}${ing.unit ? ' ' + ing.unit : ''}`;
}

const stepsBySection = computed(() => {
    if (!recipe.value) return [];
    const sections: { name: string | null; steps: typeof recipe.value.steps }[] = [];
    for (const step of recipe.value.steps) {
        const last = sections[sections.length - 1];
        if (last && last.name === step.section) {
            last.steps.push(step);
        } else {
            sections.push({ name: step.section, steps: [step] });
        }
    }
    return sections;
});

const ingredientsBySection = computed(() => {
    if (!recipe.value) return [];
    const sections: { name: string | null; items: typeof recipe.value.ingredients }[] = [];
    for (const ing of recipe.value.ingredients) {
        const last = sections[sections.length - 1];
        if (last && last.name === ing.section) {
            last.items.push(ing);
        } else {
            sections.push({ name: ing.section, items: [ing] });
        }
    }
    return sections;
});

function toggleStep(no: number) {
    if (checkedSteps.value.has(no)) checkedSteps.value.delete(no);
    else checkedSteps.value.add(no);
}

const allDone = computed(() =>
    recipe.value ? checkedSteps.value.size === recipe.value.steps.length : false
);
</script>

<template>
    <div v-if="!recipe" style="padding: 80px; text-align: center;" class="muted t-small">Ładowanie…</div>

    <article v-else class="container recipe-page">
        <!-- Header -->
        <header class="recipe-header">
            <div class="recipe-header__photo">
                <div class="photo" style="border-radius: var(--r-lg); --hue: 220;">
                    <img
                        v-if="recipe.photoUrl"
                        :src="recipe.photoUrl"
                        :alt="recipe.name"
                        style="position:absolute;inset:0;width:100%;height:100%;object-fit:cover;"
                    />
                </div>
            </div>

            <div class="recipe-header__meta">
                <div style="display: flex; gap: 8px; align-items: center; margin-bottom: 12px;">
                    <span v-if="recipe.category" class="badge" :class="`cat-${recipe.category}`">
                        {{ recipe.category }}<template v-if="recipe.region"> · {{ recipe.region }}</template>
                    </span>
                </div>

                <h1 style="font-size: var(--text-h1); font-weight: 500; letter-spacing: -2.2px; line-height: 0.95; margin-bottom: 12px;">
                    {{ recipe.name }}
                </h1>

                <p v-if="recipe.tagline" class="muted" style="font-size: var(--text-body); margin-bottom: 20px;">
                    {{ recipe.tagline }}
                </p>

                <div class="recipe-stats">
                    <div v-if="recipe.timeMinutes" class="recipe-stats__item">
                        <span class="t-meta muted">Czas</span>
                        <span style="font-weight: 600;">{{ fmtTime(recipe.timeMinutes) }}</span>
                    </div>
                    <div v-if="recipe.difficulty" class="recipe-stats__item">
                        <span class="t-meta muted">Trudność</span>
                        <span style="font-weight: 600;">{{ diffLabel[recipe.difficulty] }}</span>
                    </div>
                    <div v-if="recipe.kcalPerServing" class="recipe-stats__item">
                        <span class="t-meta muted">Kcal / porcja</span>
                        <span style="font-weight: 600;">{{ recipe.kcalPerServing }}</span>
                    </div>
                </div>

                <RouterLink
                    v-if="isAuthenticated"
                    :to="{ name: 'recipe-edit', params: { id: recipe.id } }"
                    class="btn btn--secondary btn--sm"
                    style="margin-top: 20px;"
                >
                    Edytuj przepis
                </RouterLink>
            </div>
        </header>

        <!-- Body: steps left, ingredients right -->
        <div class="recipe-body">
            <!-- Steps -->
            <section class="recipe-steps">
                <div class="section-head">
                    <span class="n">01</span><h2>Przygotowanie</h2><span class="line" />
                </div>

                <template v-for="section in stepsBySection" :key="section.name">
                    <p v-if="section.name" class="t-meta" style="color: var(--accent); margin-bottom: 12px;">
                        {{ section.name }}
                    </p>
                    <ol style="list-style: none; padding: 0; margin: 0 0 24px;">
                        <li
                            v-for="step in section.steps"
                            :key="step.stepNo"
                            class="recipe-step"
                            :class="{ 'is-done': checkedSteps.has(step.stepNo) }"
                            @click="toggleStep(step.stepNo)"
                        >
                            <span class="recipe-step__num">{{ step.stepNo }}</span>
                            <p class="recipe-step__text">{{ step.text }}</p>
                        </li>
                    </ol>
                </template>

                <div v-if="allDone" style="padding: 24px; text-align: center; border: 1px solid var(--rule); border-radius: var(--r-lg); color: var(--success);">
                    <p style="font-weight: 600;">Smacznego! 🎉</p>
                </div>
            </section>

            <!-- Ingredients sidebar -->
            <aside class="recipe-sidebar">
                <div class="surface" style="padding: 24px; position: sticky; top: 80px;">
                    <div class="section-head" style="margin-bottom: 16px;">
                        <span class="n">02</span><h2>Składniki</h2><span class="line" />
                    </div>

                    <!-- Serving stepper -->
                    <div style="display: flex; align-items: center; gap: 12px; margin-bottom: 20px; padding-bottom: 16px; border-bottom: 1px solid var(--rule);">
                        <span class="t-meta muted">Porcje</span>
                        <div style="display: flex; align-items: center; gap: 8px; margin-left: auto;">
                            <button class="icon-btn" @click="servings = Math.max(1, servings - 1)">−</button>
                            <span style="font-weight: 600; min-width: 24px; text-align: center;">{{ servings }}</span>
                            <button class="icon-btn" @click="servings++">+</button>
                        </div>
                    </div>

                    <template v-for="section in ingredientsBySection" :key="section.name">
                        <p v-if="section.name" class="t-meta" style="color: var(--accent); margin-bottom: 8px;">
                            {{ section.name }}
                        </p>
                        <ul style="list-style: none; padding: 0; margin: 0 0 16px;">
                            <li
                                v-for="ing in section.items"
                                :key="ing.sortOrder"
                                style="display: flex; justify-content: space-between; align-items: baseline; gap: 8px; padding: 6px 0; border-bottom: 1px solid var(--rule);"
                            >
                                <span style="font-size: var(--text-small);">{{ ing.ingredientName }}</span>
                                <span class="mono muted" style="font-size: 12px; flex-shrink: 0;">
                                    {{ scaleAmount(ing) }}
                                </span>
                            </li>
                        </ul>
                    </template>
                </div>
            </aside>
        </div>
    </article>
</template>

<style scoped>
.recipe-page {
    padding-top: 48px;
    padding-bottom: 80px;
}

.recipe-header {
    display: grid;
    grid-template-columns: 360px 1fr;
    gap: 48px;
    align-items: start;
    margin-bottom: 56px;
}

.recipe-stats {
    display: flex;
    gap: 0;
    border: 1px solid var(--rule);
    border-radius: var(--r-lg);
    overflow: hidden;
}

.recipe-stats__item {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 4px;
    padding: 12px 16px;
    border-right: 1px solid var(--rule);
}

.recipe-stats__item:last-child {
    border-right: none;
}

.recipe-body {
    display: grid;
    grid-template-columns: 1fr 320px;
    gap: 48px;
    align-items: start;
}

.recipe-step {
    display: flex;
    gap: 16px;
    padding: 16px 0;
    border-bottom: 1px solid var(--rule);
    cursor: pointer;
    transition: opacity var(--dur-fast);
}

.recipe-step.is-done {
    opacity: 0.4;
}

.recipe-step.is-done .recipe-step__num {
    background: var(--rule);
    color: var(--ink-faint);
}

.recipe-step__num {
    width: 28px;
    height: 28px;
    border-radius: var(--r-pill);
    background: var(--accent);
    color: var(--on-accent);
    display: flex;
    align-items: center;
    justify-content: center;
    font-family: var(--font-mono);
    font-size: 12px;
    font-weight: 700;
    flex-shrink: 0;
    margin-top: 2px;
    transition: all var(--dur-fast);
}

.recipe-step__text {
    font-size: var(--text-body);
    line-height: 1.6;
}

@media (max-width: 900px) {
    .recipe-header {
        grid-template-columns: 1fr;
    }

    .recipe-body {
        grid-template-columns: 1fr;
    }

    .recipe-sidebar aside {
        position: static !important;
    }
}
</style>
