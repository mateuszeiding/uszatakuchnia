<script lang="ts" setup>
import { useForm } from 'vee-validate'

const form = useForm()
</script>
<template>
  <div>Form {{ form.values }}</div>
</template>
