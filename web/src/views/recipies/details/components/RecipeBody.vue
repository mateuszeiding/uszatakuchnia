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
import { usePortion } from '@/composables/usePortion';
import { showToast } from '@/composables/useToast';
import type { RecipeDto } from '@/data/dtos/recipe/RecipeDto';
import type { IDictionaries } from '@/data/dtos/taxonomy/DictionariesDto';
import { uncheckedText } from '@/shared/utils/progress';
import type { IRefContext } from '@/shared/utils/stepRefs';
import { activeSubs, swapSummary, variantNutrition } from '@/shared/utils/substitutes';

const props = defineProps<{ recipe: RecipeDto; dicts: IDictionaries | undefined }>();

const {
    progress,
    toggle,
    setStep,
    setServings,
    setSize,
    setSub,
    clearSubs,
    hasChecks,
    uncheckAll,
    restoreChecks,
} = useCookingProgress(props.recipe.id, props.recipe.servings || 1);
// One scale for the list, step quantities, cooking mode and macros: servings or the chosen size.
const { size, scale } = usePortion(() => props.recipe, progress);

// "Odznacz wszystko": no confirmation, with "Cofnij" in the toast for 8 s.
function onUncheckAll() {
    const snap = uncheckAll();
    showToast(uncheckedText(snap.steps.length, snap.ing.length), {
        action: { label: 'Cofnij', run: () => restoreChecks(snap) },
    });
}
const cooking = ref(false);

const hasSwap = computed(() => props.recipe.ingredients.some((i) => i.substitutes.length));
const hasLink = computed(() => props.recipe.ingredients.some((i) => publicLink(i)));

const subsActive = computed(
    () => activeSubs(props.recipe.ingredients, progress.value.subs).size > 0
);
const nutrition = computed(() =>
    variantNutrition(props.recipe.ingredients, progress.value.subs, props.recipe.nutrition)
);
// Quantities in step text follow the same servings and substitutes as the ingredient list.
const refCtx = computed<IRefContext>(() => ({
    ingredients: props.recipe.ingredients,
    steps: props.recipe.steps.map((s) => s.text),
    units: props.dicts?.units ?? [],
    ratio: scale.value,
    subs: progress.value.subs,
    terms: props.recipe.terms,
}));
const swapNote = computed(() =>
    subsActive.value ? swapSummary(props.recipe.ingredients, progress.value.subs) : ''
);
</script>

<template>
    <RecipeHero
        :recipe="recipe"
        :nutrition="nutrition"
        :servings="progress.servings"
        :size="size"
    />
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
                :size="size"
                :scale="scale"
                :checked="progress.ing"
                :subs="progress.subs"
                @toggle="toggle('ing', $event)"
                @servings="setServings"
                @size="setSize"
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
                :refs="refCtx"
                :done="progress.steps"
                :can-reset="hasChecks"
                @toggle="toggle('steps', $event)"
                @cook="cooking = true"
                @reset="onUncheckAll"
            />
            <NutritionCard
                v-if="nutrition"
                class="recipe-grid__nutri"
                :nutrition="nutrition"
                :servings="progress.servings"
                :recipe="recipe"
                :size="size"
                :swap-note="swapNote"
                @restore="clearSubs"
            />
        </div>
    </div>
    <CookMode
        v-if="cooking && recipe.steps.length"
        :recipe="recipe"
        :refs="refCtx"
        :dicts="dicts"
        :servings="progress.servings"
        :portion-size="size"
        :scale="scale"
        :done="progress.steps"
        :checked-ing="progress.ing"
        :subs="progress.subs"
        @close="cooking = false"
        @set-step="setStep"
        @toggle-ing="toggle('ing', $event)"
        @sub="setSub"
        @servings="setServings"
        @size="setSize"
        @clear-subs="clearSubs"
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
/* the left column spans the full grid height: the ingredients panel sticks inside it, the icons card scrolls away */
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
    /* icons → ingredients → macros → equipment → steps */
    .recipe-grid__side,
    .recipe-grid__main {
        display: contents;
    }
    .recipe-grid__legend {
        margin-bottom: -28px; /* the icons card right above the ingredients (12 px) */
    }
    .recipe-grid__nutri {
        order: 1;
        margin-top: -16px; /* the card right after the ingredients (24 px gap) */
    }
    .recipe-grid__equip {
        order: 2;
        margin-bottom: -16px; /* equipment reads together with the steps */
    }
    .recipe-grid__steps {
        order: 3;
    }
}
</style>
