<script lang="ts" setup>
import type { IngredientDto } from '@/data/dtos/ingredients/IngredientDto'
import UBadge from './UBadge.vue'
import type { IType } from '@/data/dtos/ingredients/IType'
import { ref } from 'vue'

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

const johnHackintosh = ref([0])
</script>
<template>
  <div class="p-relative" @mouseover="johnHackintosh = [0, 1]" @mouseleave="johnHackintosh = [0]">
    <div class="card container-cp border radius-4 pe-5 w-100 h-100 bg-gray-ghost-white">
      <div class="row">
        <div class="col-6 ps-0">
          <div class="image-wrapper p-relative">
            <img class="radius-top-s-4 radius-bot-s-4" :src="ingredient.image.urls.small" />
            <a
              class="fs-xs p-absolute right-0 bottom-0 bg-purple-magnolia py-1 px-3 radius-top-s-2 tx-violet-french-violet"
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
            <u-badge
              class="d-inline-flex"
              :val="ingredient.type.code"
              :name="ingredient.type.name"
            />
            <template v-for="v in badges" :key="v.name">
              <div class="pt-4">
                <!-- <div class="pb-2 text-end">{{ v.name }}</div> -->
                <div class="d-flex pt-4 border-top gap-3 justify-content- flex-wrap">
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
  </div>
</template>
<style lang="scss" scoped>
.card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: -1;
  scale: 100%;
  border-radius: 1.2rem;
  box-shadow: 0 4px 12px 0 rgba(31, 38, 135, 0.37);

  animation: growing 0.2s 1 normal;
}
@keyframes growing {
  from {
    scale: 95%;
  }
  to {
    scale: 100%;
  }
}

// .card::before {
//   content: '';
//   display: block;
//   position: absolute;
//   z-index: -1;
//   width: 100%;
//   height: 100%;
//   box-shadow: 10px 10px 10px 10px var(--color-purple-magnolia);
// }
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
