<script lang="ts" setup>
import { useField } from 'vee-validate';
import { computed, toRef, useId } from 'vue';

const props = defineProps<{
    name: string;
    label?: string;
    hint?: string;
    required?: boolean;
    rows?: number;
    placeholder?: string;
}>();

const uid = useId();
const { value, errorMessage, handleBlur, handleChange } = useField<string>(toRef(props, 'name'));
const modelValue = computed(() => value.value ?? '');

function onInput(e: Event) {
    handleChange((e.target as HTMLTextAreaElement).value);
}
</script>

<template>
    <div class="field">
        <label v-if="label" class="field-label" :for="uid">
            {{ label }}<span v-if="required" class="req">*</span>
        </label>
        <textarea
            class="textarea"
            :class="{ 'is-invalid': errorMessage }"
            :id="uid"
            :name="name"
            :rows="rows ?? 3"
            :placeholder="placeholder"
            :value="modelValue"
            @input="onInput"
            @blur="handleBlur"
        />
        <span v-if="hint && !errorMessage" class="field-hint">{{ hint }}</span>
        <span v-if="errorMessage" class="field-error">{{ errorMessage }}</span>
    </div>
</template>
