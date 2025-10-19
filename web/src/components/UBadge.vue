<script setup lang="ts">
import type { Palette } from '@/scss/colors.config'
import { AromaType } from '@/shared/enums/AromaType'
import type { IngredientType } from '@/shared/enums/IngredientType'
import type { TasteType } from '@/shared/enums/TasteType'

type Union = `${AromaType}` | `${IngredientType}` | `${TasteType}`
type NestedKeyUnion<T> = {
  [K in keyof T]: K extends string
    ? keyof T[K] extends string
      ? `${K}-${keyof T[K]}`
      : never
    : never
}[keyof T]

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
  <div class="d-inline radius-2 px-4 py-2 fw-medium" :class="`bg-${mapMerge[val]}`">{{ val }}</div>
</template>
<style lang="scss" scoped></style>
