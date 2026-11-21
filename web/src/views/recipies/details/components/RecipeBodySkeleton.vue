<script lang="ts" setup>
import IngredientsPanel from './IngredientsPanel.vue';
import RecipeHero from './RecipeHero.vue';
import StepsList from './StepsList.vue';

// The same components and grid as RecipeBody, in skeleton mode (recipe = null).
</script>

<template>
    <RecipeHero
        :recipe="null"
        :servings="0"
    />
    <div class="recipe-grid">
        <IngredientsPanel
            class="recipe-grid__ing"
            :recipe="null"
            :dicts="undefined"
            :servings="0"
            :checked="[]"
            :subs="{}"
        />
        <StepsList
            class="recipe-grid__steps"
            :recipe="null"
            :done="[]"
            :can-reset="false"
        />
    </div>
</template>

<!-- grid copied from RecipeBody.vue: keep them in sync -->
<style scoped>
.recipe-grid {
    display: grid;
    grid-template-columns: 340px minmax(0, 1fr);
    grid-template-rows: auto 1fr;
    gap: 28px 56px;
    margin-top: 56px;
    align-items: start;
}
.recipe-grid__ing {
    grid-column: 1;
    grid-row: 1 / span 2;
}
.recipe-grid__steps {
    grid-column: 2;
    grid-row: 1;
    min-width: 0;
}
.recipe-grid__nutri {
    grid-column: 2;
    grid-row: 2;
}

@media (max-width: 900px) {
    .recipe-grid {
        grid-template-columns: minmax(0, 1fr);
        grid-template-rows: none;
        gap: 40px;
        margin-top: 24px;
    }
    .recipe-grid__ing,
    .recipe-grid__steps,
    .recipe-grid__nutri {
        grid-column: auto;
        grid-row: auto;
    }
    .recipe-grid__nutri {
        order: 1;
        margin-top: -16px; /* the card right after the ingredients (24 px gap) */
    }
    .recipe-grid__steps {
        order: 2;
    }
}
</style>
