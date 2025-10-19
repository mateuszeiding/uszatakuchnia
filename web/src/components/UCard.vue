<script lang="ts" setup>
defineProps<{
  name: string
}>()
</script>
<template>
  <div class="border radius-5 p-6">
    {{ name }}
  </div>
</template>
