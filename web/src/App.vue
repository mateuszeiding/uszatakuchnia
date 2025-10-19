<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { fetchEnums } from './data/api/enums/fetch'
import UBadge from './components/UBadge.vue'
import UTopBar from './components/UTopBar.vue'

const enums = ref<PropsOf<typeof UBadge>['val'][]>()

onMounted(async () => {
  const aroma = fetchEnums('aroma')
  const ingredient = fetchEnums('ingredient')
  const taste = fetchEnums('taste')

  const result = await Promise.all([aroma, ingredient, taste])

  enums.value = result.flatMap((arr) => arr.map((v) => v.Code))
})
</script>

<template>
  <UTopBar />
  <main class="container pt-7">
    <div class="d-flex gap-5 flex-wrap">
      <u-badge v-for="val in enums" :key="val" :val />
    </div>
  </main>
</template>

<style scoped></style>
