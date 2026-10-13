<script lang="ts" setup>
import { computed } from 'vue';

import UFieldError from '@/components/form/UFieldError.vue';
import type { ITaxonomyAxis } from '@/data/dtos/taxonomy/DictionariesDto';

// Pole jednej osi klasyfikacji: wielokrotny wybór = chipy (aria-pressed), pojedynczy = select.
const props = defineProps<{
    axis: ITaxonomyAxis;
    hint?: string;
    error?: string;
    /** Wariant wyróżnionych chipów (np. diety). */
    strong?: boolean;
}>();
const selected = defineModel<number[]>({ required: true });

const fieldId = computed(() => `f-taxonomy-${props.axis.code}`);
const labelId = computed(() => `${fieldId.value}-label`);

function toggle(id: number) {
    selected.value = selected.value.includes(id)
        ? selected.value.filter((x) => x !== id)
        : [...selected.value, id];
}

function onSelect(e: Event) {
    const v = (e.target as HTMLSelectElement).value;
    selected.value = v ? [Number(v)] : [];
}
</script>

<template>
    <div
        v-if="axis.multi"
        :id="fieldId"
        class="field"
        role="group"
        :aria-labelledby="labelId"
    >
        <span
            :id="labelId"
            class="field-label"
        >
            {{ axis.label }}
            <span
                v-if="axis.required"
                class="req"
                aria-hidden="true"
            >
                *
            </span>
        </span>
        <div class="tax-chips">
            <button
                v-for="v in axis.values"
                :key="v.id"
                type="button"
                class="chip"
                :class="[
                    {
                        'is-active': selected.includes(v.id),
                        'chip--color': v.color,
                        'chip--strong': strong,
                    },
                    v.color ? `cat-${v.color}` : null,
                ]"
                :aria-pressed="selected.includes(v.id)"
                @click="toggle(v.id)"
            >
                <span
                    v-if="v.color"
                    class="chip__dot"
                    aria-hidden="true"
                />
                {{ v.label }}
            </button>
        </div>
        <UFieldError v-if="error">{{ error }}</UFieldError>
        <span
            v-else-if="hint"
            class="field-hint"
        >
            {{ hint }}
        </span>
    </div>
    <div
        v-else
        :id="fieldId"
        class="field"
    >
        <label
            :for="`${fieldId}-select`"
            class="field-label"
        >
            {{ axis.label }}
        </label>
        <select
            :id="`${fieldId}-select`"
            class="select"
            :value="selected[0] ?? ''"
            @change="onSelect"
        >
            <option value="">nie dotyczy</option>
            <option
                v-for="v in axis.values"
                :key="v.id"
                :value="v.id"
            >
                {{ v.label }}
            </option>
        </select>
        <UFieldError v-if="error">{{ error }}</UFieldError>
    </div>
</template>

<style scoped>
.tax-chips {
    display: flex;
    flex-wrap: wrap;
    gap: 6px;
}
</style>
