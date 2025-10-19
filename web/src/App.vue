<script setup lang="ts">
import UBadge from './components/UBadge.vue'
import { TasteType } from './shared/enums/TasteType'
</script>

<template>
  <main class="container">
    <h1>You did it!</h1>
    <UBadge :val="TasteType.BITTER" />
    <p>
      Visit <a href="https://vuejs.org/" target="_blank" rel="noopener">vuejs.org</a> to read the
      documentation
    </p>
  </main>
</template>

<style scoped></style>
