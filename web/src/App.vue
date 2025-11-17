<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { fetchEnums } from './data/api/enums/fetch'
import UBadge from './components/UBadge.vue'
import UCard from './components/UCard.vue'
import UTopBar from './components/UTopBar.vue'
import { fetchIngredients } from './data/api/ingredients/fetch'
import { RouterView } from 'vue-router'
import UMenu from './components/UMenu.vue'

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
  <UTopBar />
  <UMenu />
  <main class="container pt-7 pb-6">
    <RouterView />
  </main>
</template>

<style scoped></style>
