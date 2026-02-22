<script lang="ts" setup>
import type { Palette } from '@/scss/colors.config';
import { IngredientType } from '@/shared/enums/IngredientType';

type Union = `${IngredientType}`;

defineProps<{
    val: Union;
    name?: string;
}>();

const invertTextColorFor = [
    IngredientType.VEGETABLE.valueOf(),
    IngredientType.FRUIT.valueOf(),
    IngredientType.SPICE.valueOf(),
    IngredientType.MEAT.valueOf(),
] as const;

const ingredientMap: Record<IngredientType, NestedKeyUnion<Palette>> = {
    FISH: 'blue-picton-blue',
    FRUIT: 'green-persian-green',
    HERB: 'green-mint',
    MEAT: 'red-crimson',
    OTHER: 'brown-sandy-brown',
    SPICE: 'red-chili-red',
    VEGETABLE: 'green-dartmouth-green',
};

const mapMerge = {
    ...ingredientMap,
};
</script>
<template>
    <span
        class="radius-2 border px-3 py-1 fw-medium fs-xs cursor-default h-fit"
        :class="[
            `bg-${mapMerge[val]}`,
            {
                ['tx-gray-ghost-white']: invertTextColorFor.includes(val),
            },
        ]"
    >
        {{ name ?? val }}
    </span>
</template>
