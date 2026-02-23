<script setup lang="ts">
import { fetchEnums } from '@api/enums/fetch';
import { fetchIngredients } from '@api/ingredients/fetch';
import UBadge from '@cmp/UBadge.vue';
import UCard from '@cmp/UCard.vue';
import { onMounted, ref } from 'vue';

const enums = ref<PropsOf<typeof UBadge>['val'][]>();
const ingredients = ref<PropsOf<typeof UCard>['ingredient'][]>([]);

onMounted(async () => {
    const ingredient = fetchEnums('ingredient');
    const ingr = fetchIngredients('list');

    ingredients.value = await ingr;
    enums.value = (await ingredient).map((v) => v.code);
});
</script>
<template>
    <div class="d-flex gap-5 flex-wrap pb-8">
        <u-badge
            v-for="val in enums"
            :key="val"
            :val
        />
    </div>

    <div class="grid">
        <u-card
            v-for="ingredient in ingredients"
            :key="ingredient.id"
            :ingredient
        />
    </div>
</template>
<style lang="scss" scoped>
.grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: var(--space-6);
}
</style>
