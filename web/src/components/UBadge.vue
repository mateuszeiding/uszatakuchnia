<script lang="ts" setup>
import { IngredientType } from '@/shared/enums/IngredientType';

defineProps<{
    val: `${IngredientType}`;
    name?: string;
}>();

const catClass: Record<IngredientType, string> = {
    [IngredientType.FISH]: 'cat-ryby-ing',
    [IngredientType.MEAT]: 'cat-mieso-ing',
    [IngredientType.VEGETABLE]: 'cat-warzywo',
    [IngredientType.FRUIT]: 'cat-owoc',
    [IngredientType.DAIRY]: 'cat-nabial',
    [IngredientType.GRAIN]: 'cat-zboze',
    [IngredientType.LEGUME]: 'cat-straczki',
    [IngredientType.SPICE]: 'cat-przyprawa',
    [IngredientType.FAT]: 'cat-tluszcz',
    [IngredientType.OTHER]: 'cat-inne',
};
</script>

<template>
    <span
        class="badge"
        :class="catClass[val as IngredientType]"
    >
        {{ name ?? val }}
    </span>
</template>
