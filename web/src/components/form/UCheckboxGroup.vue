<script setup lang="ts">
import { useField } from 'vee-validate';
import { toRef, useId } from 'vue';

type CheckboxOption = {
    label: string;
    value: string | number;
};

const props = defineProps<{
    name: string;
    options: CheckboxOption[];
    legend?: string;
}>();

const uid = useId();
const legendId = `${uid}-legend`;
const optionId = (val: CheckboxOption['value']) => `${uid}-${val}`;

const { value, errorMessage, handleChange, handleBlur } = useField<Array<string | number>>(
    toRef(props, 'name')
);
</script>

<template>
    <fieldset
        class="form-group"
        :aria-labelledby="legend ? legendId : undefined"
    >
        <legend
            v-if="legend"
            :id="legendId"
        >
            {{ legend }}
        </legend>

        <div
            v-for="o in options"
            :key="o.value"
            class="form-field-h"
        >
            <input
                type="checkbox"
                class="form-input"
                :id="optionId(o.value)"
                :name="props.name"
                :value="o.value"
                :checked="value?.includes(o.value)"
                @change="handleChange"
                @blur="handleBlur"
            />

            <label
                class="form-label"
                :for="optionId(o.value)"
            >
                {{ o.label }}
            </label>
        </div>

        <div v-if="errorMessage">
            {{ errorMessage }}
        </div>
    </fieldset>
</template>
