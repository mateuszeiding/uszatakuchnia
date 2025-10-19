<script setup lang="ts">
import type { Palette } from '@/scss/colors.config'
import { AromaType } from '@/shared/enums/AromaType'
import type { IngredientType } from '@/shared/enums/IngredientType'
import type { TasteType } from '@/shared/enums/TasteType'

type Union = `${AromaType}` | `${IngredientType}` | `${TasteType}`

defineProps<{
  val: Union
}>()

const aromaMap: Record<AromaType, NestedKeyUnion<Palette>> = {
  CREAMY: 'brown-lavender-blush',
  EARTHY: 'brown-wenge',
  FRUITY: 'green-persian-green',
  HERBAL: 'green-mint',
  OTHER: 'brown-sandy-brown',
  SMOKY: 'purple-dark-purple',
  SPICY: 'red-chili-red',
}

const tasteMap: Record<TasteType, NestedKeyUnion<Palette>> = {
  BITTER: 'yellow-maize',
  SALTY: 'gray-ghost-white',
  SOUR: 'orange-ut-orange',
  SWEET: 'violet-plum',
  UMAMI: 'gray-outer-space',
  UNKNOWN: 'brown-sandy-brown',
}

const ingredientMap: Record<IngredientType, NestedKeyUnion<Palette>> = {
  FISH: 'blue-picton-blue',
  FRUIT: 'green-persian-green',
  HERB: 'green-mint',
  MEAT: 'red-crimson',
  OTHER: 'brown-sandy-brown',
  SPICE: 'red-chili-red',
  VEGETABLE: 'green-dartmouth-green',
}

const mapMerge = {
  ...aromaMap,
  ...tasteMap,
  ...ingredientMap,
}
</script>
<template>
  <span class="radius-2 border px-4 py-2 fw-medium" :class="`bg-${mapMerge[val]}`">
    {{ val }}
  </span>
</template>
<style lang="scss" scoped></style>
