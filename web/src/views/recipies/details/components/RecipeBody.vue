<script lang="ts" setup>
import { computed, ref } from 'vue';

import { publicLink } from '../recipeView';

import CookMode from './CookMode.vue';
import EquipmentCard from './EquipmentCard.vue';
import IconLegend from './IconLegend.vue';
import IngredientsPanel from './IngredientsPanel.vue';
import NutritionCard from './NutritionCard.vue';
import RecipeHero from './RecipeHero.vue';
import StepsList from './StepsList.vue';

import { useCookingProgress } from '@/composables/useCookingProgress';
import type { RecipeDto } from '@/data/dtos/recipe/RecipeDto';
import type { IDictionaries } from '@/data/dtos/taxonomy/DictionariesDto';
import { activeSubs, swapSummary, variantNutrition } from '@/shared/utils/substitutes';

const props = defineProps<{ recipe: RecipeDto; dicts: IDictionaries | undefined }>();

const { progress, toggle, setStep, setServings, setSub, clearSubs, hasProgress, reset } =
    useCookingProgress(props.recipe.id, props.recipe.servings || 1);
const cooking = ref(false);

// Ikony obecne w przepisie → karta „Ikony przy składnikach”.
const hasSwap = computed(() => props.recipe.ingredients.some((i) => i.substitutes.length));
const hasLink = computed(() => props.recipe.ingredients.some((i) => publicLink(i)));

// Wybrany wariant (zamienniki): makro i zdanie w karcie.
const subsActive = computed(
    () => activeSubs(props.recipe.ingredients, progress.value.subs).size > 0
);
const nutrition = computed(() =>
    variantNutrition(props.recipe.ingredients, progress.value.subs, props.recipe.nutrition)
);
const swapNote = computed(() =>
    subsActive.value ? swapSummary(props.recipe.ingredients, progress.value.subs) : ''
);
</script>

<template>
    <RecipeHero
        :recipe="recipe"
        :nutrition="nutrition"
        :servings="progress.servings"
    />
    <!-- Desktop: po lewej karta ikon (odjeżdża) i panel składników (przyklejony); po prawej Sprzęt → Kroki → Wartości odżywcze.
         Mobile: ikony → składniki → makro → sprzęt → kroki. -->
    <div class="recipe-grid">
        <div class="recipe-grid__side">
            <IconLegend
                v-if="hasSwap || hasLink"
                class="recipe-grid__legend"
                :swap="hasSwap"
                :link="hasLink"
            />
            <IngredientsPanel
                class="recipe-grid__ing"
                :recipe="recipe"
                :dicts="dicts"
                :servings="progress.servings"
                :checked="progress.ing"
                :subs="progress.subs"
                @toggle="toggle('ing', $event)"
                @servings="setServings"
                @sub="setSub"
                @clear-subs="clearSubs"
            />
        </div>
        <div class="recipe-grid__main">
            <EquipmentCard
                v-if="recipe.equipment.length"
                class="recipe-grid__equip"
                :items="recipe.equipment"
            />
            <StepsList
                class="recipe-grid__steps"
                :recipe="recipe"
                :done="progress.steps"
                :can-reset="hasProgress"
                @toggle="toggle('steps', $event)"
                @cook="cooking = true"
                @reset="reset"
            />
            <NutritionCard
                v-if="nutrition"
                class="recipe-grid__nutri"
                :nutrition="nutrition"
                :servings="progress.servings"
                :swap-note="swapNote"
                @restore="clearSubs"
            />
        </div>
    </div>
    <CookMode
        v-if="cooking && recipe.steps.length"
        :recipe="recipe"
        :dicts="dicts"
        :servings="progress.servings"
        :done="progress.steps"
        :checked-ing="progress.ing"
        :subs="progress.subs"
        @close="cooking = false"
        @set-step="setStep"
        @toggle-ing="toggle('ing', $event)"
        @sub="setSub"
    />
</template>

<style scoped>
.recipe-grid {
    display: grid;
    grid-template-columns: 340px minmax(0, 1fr);
    gap: 56px;
    margin-top: 56px;
    align-items: start;
}
/* lewa kolumna na całą wysokość siatki — panel składników przykleja się w niej, karta ikon odjeżdża */
.recipe-grid__side {
    align-self: stretch;
    display: flex;
    flex-direction: column;
    gap: 12px;
}
.recipe-grid__main {
    min-width: 0;
    display: flex;
    flex-direction: column;
    gap: 28px;
}

@media (max-width: 900px) {
    .recipe-grid {
        grid-template-columns: minmax(0, 1fr);
        gap: 40px;
        margin-top: 24px;
    }
    /* ikony → składniki → makro → sprzęt → kroki */
    .recipe-grid__side,
    .recipe-grid__main {
        display: contents;
    }
    .recipe-grid__legend {
        margin-bottom: -28px; /* karta ikon zaraz nad składnikami (12 px) */
    }
    .recipe-grid__nutri {
        order: 1;
        margin-top: -16px; /* karta zaraz po składnikach (jak w v1: odstęp 24px) */
    }
    .recipe-grid__equip {
        order: 2;
        margin-bottom: -16px; /* sprzęt czyta się razem z krokami */
    }
    .recipe-grid__steps {
        order: 3;
    }
}
</style>
