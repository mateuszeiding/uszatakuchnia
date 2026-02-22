<script lang="ts" setup>
import { useField } from 'vee-validate';
import { computed, toRef, useId } from 'vue';

const props = defineProps<{
    name: string;
    label?: string;
    rows?: number;
}>();

const uid = useId();

const { value, errorMessage, handleBlur, handleChange } = useField<string>(toRef(props, 'name'));

const modelValue = computed(() => value.value ?? '');

function onInput(e: Event) {
    const el = e.target as HTMLTextAreaElement;
    handleChange(el.value);
}
</script>

<template>
    <div class="form-field">
        <label
            v-if="props.label"
            class="form-label"
            :for="uid"
        >
            {{ props.label }}
        </label>

        <textarea
            class="form-input"
            :id="uid"
            :name="props.name"
            :rows="props.rows ?? 3"
            :value="modelValue"
            @input="onInput"
            @blur="handleBlur"
        />

        <div
            v-if="errorMessage"
            class="form-error"
        >
            {{ errorMessage }}
        </div>
    </div>
</template>
