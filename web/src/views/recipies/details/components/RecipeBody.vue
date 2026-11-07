<script lang="ts" setup>
import { ref } from 'vue';

import CookMode from './CookMode.vue';
import IngredientsPanel from './IngredientsPanel.vue';
import NutritionCard from './NutritionCard.vue';
import RecipeHero from './RecipeHero.vue';
import StepsList from './StepsList.vue';

import { useCookingProgress } from '@/composables/useCookingProgress';
import type { RecipeDto } from '@/data/dtos/recipe/RecipeDto';
import type { IDictionaries } from '@/data/dtos/taxonomy/DictionariesDto';

const props = defineProps<{ recipe: RecipeDto; dicts: IDictionaries | undefined }>();

const { progress, toggle, setStep, setServings, hasProgress, reset } = useCookingProgress(
    props.recipe.id,
    props.recipe.servings || 1
);
const cooking = ref(false);
</script>

<template>
    <RecipeHero
        :recipe="recipe"
        :servings="progress.servings"
    />
    <div class="recipe-grid">
        <div
            class="recipe-side"
            :class="{ 'recipe-side--nutrition': recipe.nutrition }"
        >
            <IngredientsPanel
                :recipe="recipe"
                :dicts="dicts"
                :servings="progress.servings"
                :checked="progress.ing"
                @toggle="toggle('ing', $event)"
                @servings="setServings"
            />
            <NutritionCard
                v-if="recipe.nutrition"
                :nutrition="recipe.nutrition"
                :servings="progress.servings"
            />
        </div>
        <StepsList
            :recipe="recipe"
            :done="progress.steps"
            :can-reset="hasProgress"
            @toggle="toggle('steps', $event)"
            @cook="cooking = true"
            @reset="reset"
        />
    </div>
    <CookMode
        v-if="cooking && recipe.steps.length"
        :recipe="recipe"
        :dicts="dicts"
        :servings="progress.servings"
        :done="progress.steps"
        :checked-ing="progress.ing"
        @close="cooking = false"
        @set-step="setStep"
        @toggle-ing="toggle('ing', $event)"
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
.recipe-side {
    display: flex;
    flex-direction: column;
    gap: 24px;
    min-width: 0;
}
/* panel składników jest przyklejony; z kartą makro pod spodem zasłaniałby ją przy przewijaniu */
.recipe-side--nutrition :deep(.ing-panel) {
    position: static;
}

@media (max-width: 900px) {
    .recipe-grid {
        grid-template-columns: minmax(0, 1fr);
        gap: 40px;
        margin-top: 32px;
    }
}
</style>
