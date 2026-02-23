<script lang="ts" setup>
import { ref } from 'vue';

import RecipeCard from './components/RecipeCard.vue';

import { fetchRecipes } from '@/data/api/recipe/fetch';
import { RecipeBaseDto } from '@/data/dtos/recipe/RecipeDto';

const recipes = ref<RecipeBaseDto[]>([]);

fetchRecipes('list').then((v) => (recipes.value = v));
</script>
<template>
    <div>Recipes</div>
    <RecipeCard
        v-for="recipe in recipes"
        :key="recipe.id"
        :recipe
    />
</template>
