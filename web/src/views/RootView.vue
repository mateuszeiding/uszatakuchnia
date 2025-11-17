<script lang="ts" setup>
import UInput from '@/components/form/UInput.vue'
</script>
<template>
  <div>
    <UInput name="test" label="test" type="text" />
    <UInput name="test" label="test" type="radio" />
    <UInput name="test" label="test" type="checkbox" />
    <UInput name="test" label="test" type="checkbox" switch />
  </div>
</template>
