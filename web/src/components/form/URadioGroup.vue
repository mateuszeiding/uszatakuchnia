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
        class="form-group"
        :aria-labelledby="legend ? legendId : undefined"
    >
        <legend
            v-if="legend"
            class="form-legend"
            :id="legendId"
        >
            {{ legend }}
        </legend>

        <div
            v-for="o in options"
            :key="o.value"
            class="form-field"
        >
            <input
                type="radio"
                class="form-input"
                :id="optionId(o.value)"
                :name="props.name"
                :value="o.value"
                :checked="value === o.value"
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

        <div
            v-if="errorMessage"
            class="form-error"
        >
            {{ errorMessage }}
        </div>
    </fieldset>
</template>
