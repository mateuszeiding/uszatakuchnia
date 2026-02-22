<script lang="ts" setup>
import { useField } from 'vee-validate';
import { computed, toRef, useId } from 'vue';

const props = defineProps<{
    name: string;
    label?: string;
}>();

const uid = useId();
const { value, errorMessage, handleBlur } = useField<boolean>(toRef(props, 'name'));

const checked = computed(() => Boolean(value.value));

function onChange(e: Event) {
    const el = e.target as HTMLInputElement;
    value.value = el.checked;
}
</script>

<template>
    <div class="form-field-h">
        <input
            class="form-input"
            type="checkbox"
            role="switch"
            :id="uid"
            :name="props.name"
            :checked="checked"
            @change="onChange"
            @blur="handleBlur"
        />

        <label
            v-if="props.label"
            class="form-label"
            :for="uid"
        >
            {{ props.label }}
        </label>

        <div
            v-if="errorMessage"
            class="form-error"
        >
            {{ errorMessage }}
        </div>
    </div>
</template>
