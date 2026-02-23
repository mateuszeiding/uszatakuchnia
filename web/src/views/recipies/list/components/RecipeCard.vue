<script lang="ts" setup>
import type { RecipeBaseDto } from '@/data/dtos/recipe/RecipeDto';

defineProps<{
    recipe: RecipeBaseDto;
}>();
</script>
<template>
    <RouterLink
        class="d-block p-4 p-relative border radius-4 pe-5 w-100 h-100 bg-gray-ghost-white"
        :to="{ name: 'recipe-details', params: { id: recipe.id } }"
    >
        <div class="fs-lg pb-2">
            {{ recipe.name }}
        </div>
    </RouterLink>
</template>
