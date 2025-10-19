<script lang="ts" setup>
import type { IngredientDto } from '@/data/dtos/ingredients/IngredientDto'
import UBadge from './UBadge.vue'

defineProps<{
  ingredient: IngredientDto
}>()
</script>
<template>
  <div class="border radius-5 pb-6 pt-4 w-100 h-100">
    <div class="px-5">
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
      <div class="image-wrapper pt-4">
        <img class="radius-4" :src="ingredient.image.urls.small" />
      </div>
      <div class="fs-xs">
        <span>Photo by </span>
        <a :href="ingredient.image.author.profileURL">{{
          ingredient.image.author.name ?? ingredient.image.author.username
        }}</a>
        <span> on <b>Unsplash</b></span>
      </div>
    </div>
    <div class="flex-wrap px-5 py-4">
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
<style lang="scss" scoped>
.image-wrapper {
  width: 100%;
  aspect-ratio: 1 / 1;
  overflow: hidden;
  position: relative;

  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    object-position: center center;
    display: block;
  }
}
</style>
