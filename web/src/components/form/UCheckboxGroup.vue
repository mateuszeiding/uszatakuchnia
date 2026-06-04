<script setup lang="ts">
import { useField } from 'vee-validate';
import { toRef, useId } from 'vue';

type CheckboxOption = { label: string; value: string | number };

const props = defineProps<{
    name: string;
    options: CheckboxOption[];
    legend?: string;
}>();

const uid = useId();
const legendId = `${uid}-legend`;
const optionId = (val: CheckboxOption['value']) => `${uid}-${val}`;

const { value, errorMessage, handleBlur } = useField<Array<string | number>>(toRef(props, 'name'));

const toggle = (optionValue: string | number, e: Event) => {
    const checked = (e.target as HTMLInputElement).checked;
    const current = value.value ?? [];
    value.value = checked
        ? current.includes(optionValue)
            ? current
            : [...current, optionValue]
        : current.filter((v) => v !== optionValue);
};
</script>

<template>
    <fieldset
        class="field"
        style="border: none; padding: 0; margin: 0"
        :aria-labelledby="legend ? legendId : undefined"
    >
        <legend
            v-if="legend"
            :id="legendId"
            class="field-label"
        >
            {{ legend }}
        </legend>
        <div style="display: flex; flex-direction: column; gap: 10px; margin-top: 6px">
            <label
                v-for="o in options"
                :key="o.value"
                class="check"
                :for="optionId(o.value)"
            >
                <input
                    type="checkbox"
                    :id="optionId(o.value)"
                    :name="name"
                    :value="o.value"
                    :checked="value?.includes(o.value)"
                    @change="(e) => toggle(o.value, e)"
                    @blur="handleBlur"
                />
                <span class="box" />
                {{ o.label }}
            </label>
        </div>
        <span
            v-if="errorMessage"
            class="field-error"
        >
            {{ errorMessage }}
        </span>
    </fieldset>
</template>
