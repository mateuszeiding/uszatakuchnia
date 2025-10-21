<script lang="ts" setup>
import type { IngredientDto } from '@/data/dtos/ingredients/IngredientDto'
import UBadge from './UBadge.vue'
import type { IType } from '@/data/dtos/ingredients/IType'

const props = defineProps<{
  ingredient: IngredientDto
}>()

const badges: { name: string; arr: IType<PropsOf<typeof UBadge>['val']>[] }[] = [
  // {
  //   name: 'Aromat',
  //   arr: props.ingredient.aromas,
  // },
  {
    name: 'Smak',
    arr: props.ingredient.tastes,
  },
]
</script>
<template>
  <div class="container-cp border radius-4 pe-5 pb-4 w-100 h-100 bg-purple-magnolia">
    <div class="row pb-4">
      <div class="col-6 ps-0">
        <div class="image-wrapper p-relative">
          <img class="radius-top-s-4 radius-bot-e-4" :src="ingredient.image.urls.small" />
          <a
            class="fs-xs p-absolute left-0 bottom-0 bg-purple-magnolia py-1 px-3 radius-top-e-2 tx-violet-french-violet"
            :href="ingredient.image.author.profileURL"
            >{{ ingredient.image.author.name ?? ingredient.image.author.username }}</a
          >
        </div>
      </div>
      <div class="col-6 pt-4">
        <div>
          <div class="fs-xl pb-2">
            {{ ingredient.name }}
          </div>
          <u-badge class="d-inline-flex" :val="ingredient.type.code" :name="ingredient.type.name" />
          <template v-for="v in badges" :key="v.name">
            <div class="pt-4">
              <!-- <div class="pb-2 text-end">{{ v.name }}</div> -->
              <div class="d-flex pt-3 border-top gap-3 justify-content- flex-wrap">
                <u-badge v-for="val in v.arr" :key="val.code" :val="val.code" :name="val.name" />
              </div>
            </div>
          </template>
        </div>
      </div>
    </div>

    <div class="row">
      <div class="col"></div>
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
