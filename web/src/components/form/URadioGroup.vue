<script lang="ts" setup>
import { useField } from 'vee-validate';
import { toRef, useId } from 'vue';

type RadioOption = { label: string; value: string | number };

const props = defineProps<{
    name: string;
    options: RadioOption[];
    legend?: string;
}>();

const uid = useId();
const legendId = `${uid}-legend`;
const optionId = (val: RadioOption['value']) => `${uid}-${val}`;

const { value, errorMessage, handleChange, handleBlur } = useField<string | number>(
    toRef(props, 'name')
);
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
                class="radio"
                :for="optionId(o.value)"
            >
                <input
                    type="radio"
                    :id="optionId(o.value)"
                    :name="name"
                    :value="o.value"
                    :checked="value === o.value"
                    @change="handleChange"
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
