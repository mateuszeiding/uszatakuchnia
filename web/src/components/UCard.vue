<script lang="ts" setup>
import type { IngredientDto } from '@/data/dtos/ingredients/IngredientDto'
import UBadge from './UBadge.vue'
import type { IType } from '@/data/dtos/ingredients/IType'

const props = defineProps<{
  ingredient: IngredientDto
}>()

const badges: { name: string; arr: IType<PropsOf<typeof UBadge>['val']>[] }[] = [
  {
    name: 'Aromat',
    arr: props.ingredient.aromas,
  },
  {
    name: 'Smak',
    arr: props.ingredient.tastes,
  },
]
</script>
<template>
  <div class="container-cp border radius-4 px-5 py-4 w-100 h-100 bg-brown-lavender-blush">
    <div class="row">
      <div class="col">
        <div class="d-flex justify-content-between align-items-center pb-2">
          <div class="fs-xl">
            {{ ingredient.name }}
          </div>
          <u-badge class="d-inline-flex" :val="ingredient.type.code" :name="ingredient.type.name" />
        </div>
      </div>
    </div>
    <div class="row">
      <div class="col-6">
        <div class="py-4">
          <div class="image-wrapper">
            <img class="radius-5" :src="ingredient.image.urls.small" />
          </div>
          <div class="fs-xs">
            <span>Photo by </span>
            <a :href="ingredient.image.author.profileURL">{{
              ingredient.image.author.name ?? ingredient.image.author.username
            }}</a>
            <span> on <b>Unsplash</b></span>
          </div>
        </div>
      </div>
      <div class="col-6">
        <template v-for="v in badges" :key="v.name">
          <div class="py-2">
            <div class="pb-2 text-end">{{ v.name }}</div>
            <div class="d-flex gap-3 justify-content-end flex-wrap">
              <u-badge v-for="val in v.arr" :key="val.code" :val="val.code" :name="val.name" />
            </div>
          </div>
        </template>
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
