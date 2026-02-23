<script lang="ts" setup>
import { ref } from 'vue';

import { fetchRecipes } from '@/data/api/recipe/fetch';
import { RecipeDto } from '@/data/dtos/recipe/RecipeDto';

const recipe = ref<RecipeDto>();
const props = defineProps<{
    id: number;
}>();
fetchRecipes(props.id).then((v) => (recipe.value = v));
</script>
<template>
    <div>{{ recipe }}</div>
</template>
