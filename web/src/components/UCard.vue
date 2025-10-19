<script lang="ts" setup>
import type { IngredientDto } from '@/data/dtos/ingredients/ingredientDto'
import UBadge from './UBadge.vue'

defineProps<{
  ingredient: IngredientDto
}>()
</script>
<template>
  <div class="border radius-5 pb-6 pt-4">
    <div class="pb-4 px-5">
      <span class="fs-lg">
        {{ ingredient.name }}
      </span>
    </div>
    <div class="border-top px-5 py-4">
      <span>Smak</span>
      <div class="d-flex flex-wrap gap-3">
        <u-badge
          v-for="val in ingredient.tastes"
          :key="val.code"
          :val="val.code"
          :name="val.name"
        />
      </div>
    </div>
    <div class="border-top flex-wrap px-5 py-4">
      <span>Aromat</span>
      <div class="d-flex gap-3">
        <u-badge
          v-for="val in ingredient.aromas"
          :key="val.code"
          :val="val.code"
          :name="val.name"
        />
      </div>
    </div>
  </div>
</template>
