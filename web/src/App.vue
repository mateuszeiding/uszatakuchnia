<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { RouterView } from 'vue-router';

import UBadge from './components/UBadge.vue';
import UCard from './components/UCard.vue';
import UMenu from './components/UMenu.vue';
import UTopBar from './components/UTopBar.vue';
import { fetchEnums } from './data/api/enums/fetch';
import { fetchIngredients } from './data/api/ingredients/fetch';

const enums = ref<PropsOf<typeof UBadge>['val'][]>();
const ingredients = ref<PropsOf<typeof UCard>['ingredient'][]>([]);

onMounted(async () => {
    const ingredient = fetchEnums('ingredient');
    const ingr = fetchIngredients('list');

    ingredients.value = await ingr;

    enums.value = await ingredient;
});
</script>

<template>
    <UTopBar />
    <UMenu />
    <main class="container pt-7 pb-6">
        <RouterView />
    </main>
</template>
