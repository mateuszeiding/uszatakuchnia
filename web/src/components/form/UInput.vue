<script lang="ts" setup>
import { computed, useId } from 'vue'

type Base = {
  name: string
  label: string
}

type Rest = Base & {
  type: 'text' | 'radio'
}

type Checkbox = Base & {
  type: 'checkbox'
  switch?: boolean
}

const props = defineProps<Rest | Checkbox>()

const uid = useId()

const switchRole = computed(() =>
  props.type === 'checkbox' && props.switch ? 'switch' : undefined,
)
console.log(props, switchRole.value)
</script>
<template>
  <div class="form-field">
    <label class="form-label" v-if="label" :for="uid">{{ label }}</label>
    <input class="form-input" :id="uid" :name :type :role="switchRole" />
  </div>
</template>
