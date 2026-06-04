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
    value.value = (e.target as HTMLInputElement).checked;
}
</script>

<template>
    <div class="field">
        <label class="switch" :for="uid">
            <input
                type="checkbox"
                role="switch"
                :id="uid"
                :name="name"
                :checked="checked"
                @change="onChange"
                @blur="handleBlur"
            />
            <span class="track" />
            <span v-if="label">{{ label }}</span>
        </label>
        <span v-if="errorMessage" class="field-error">{{ errorMessage }}</span>
    </div>
</template>
