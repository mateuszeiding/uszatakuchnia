<script lang="ts" setup>
import type { IngredientDto } from '@/data/dtos/ingredients/IngredientDto'
import UBadge from './UBadge.vue'

defineProps<{
  ingredient: IngredientDto
}>()
</script>
<template>
  <div class="border radius-5 pb-6 pt-4 w-100 h-100">
    <div class="pb-4 px-5">
      <div class="fs-xl pb-2">
        {{ ingredient.name }}
      </div>
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
      <div class="pb-2">Aromat</div>
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
