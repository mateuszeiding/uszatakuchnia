<script lang="ts" setup>
import { ref } from 'vue';

import { fetchRecipes } from '@/data/api/recipe/fetch';
import { RecipeDto } from '@/data/dtos/recipe/RecipeDto';

const recipes = ref<RecipeDto[]>([]);

fetchRecipes('list').then((v) => (recipes.value = v));
</script>
<template>
    <div>Recipes</div>
    <div>{{ recipes }}</div>
</template>
