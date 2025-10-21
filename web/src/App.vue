<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { fetchEnums } from './data/api/enums/fetch'
import UBadge from './components/UBadge.vue'
import UCard from './components/UCard.vue'
import UTopBar from './components/UTopBar.vue'
import { fetchIngredients } from './data/api/ingredients/fetch'

const enums = ref<PropsOf<typeof UBadge>['val'][]>()
const ingredients = ref<PropsOf<typeof UCard>['ingredient'][]>([])

onMounted(async () => {
  const aroma = fetchEnums('aroma')
  const ingredient = fetchEnums('ingredient')
  const taste = fetchEnums('taste')
  const ingr = fetchIngredients('list')

  const result = await Promise.all([aroma, ingredient, taste])

  ingredients.value = await ingr

  enums.value = result.flatMap((arr) => arr.map((v) => v.Code))
})
</script>

<template>
  <u-top-bar />
  <main class="container pt-7 pb-6">
    <div class="d-flex gap-5 flex-wrap pb-8">
      <u-badge v-for="val in enums" :key="val" :val />
    </div>

    <div class="row row-cols-3-xl row-cols-2-lg row-cols-1-md row-cols-1-sm row-gap-6">
      <div class="col" v-for="ingredient in ingredients" :key="ingredient.id">
        <u-card :ingredient />
      </div>
    </div>
  </main>
</template>

<style scoped></style>
