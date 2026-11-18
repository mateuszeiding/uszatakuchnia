<script lang="ts" setup>
import { ref } from 'vue';

import CookMode from './CookMode.vue';
import EquipmentCard from './EquipmentCard.vue';
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
    <!-- Desktop: panel składników (przyklejony) po lewej, po prawej Sprzęt → Kroki → Wartości odżywcze.
         Mobile: składniki → makro → sprzęt → kroki. -->
    <div class="recipe-grid">
        <IngredientsPanel
            class="recipe-grid__ing"
            :recipe="recipe"
            :dicts="dicts"
            :servings="progress.servings"
            :checked="progress.ing"
            @toggle="toggle('ing', $event)"
            @servings="setServings"
        />
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
                v-if="recipe.nutrition"
                class="recipe-grid__nutri"
                :nutrition="recipe.nutrition"
                :servings="progress.servings"
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
        margin-top: 32px;
    }
    /* składniki → makro → sprzęt → kroki */
    .recipe-grid__main {
        display: contents;
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
