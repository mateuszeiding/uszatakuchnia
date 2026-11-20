<script lang="ts" setup>
// Number field for admin tables: unit inside the field, an inherited value shown as a ghost.
defineProps<{
    id: string;
    /** Screen reader label (e.g. "Białko, min"); visually it's the row and column headers. */
    label: string;
    unit: string;
    invalid?: boolean;
    describedby?: string;
    /** Hint shown when the field is empty (the inherited value). */
    ghost?: string | null;
}>();
const model = defineModel<string>({ required: true });
</script>

<template>
    <label
        class="num-input"
        :for="id"
    >
        <span class="sr-only">{{ label }}</span>
        <input
            :id="id"
            v-model="model"
            class="input input--mono num-input__field"
            :class="{ 'is-invalid': invalid }"
            inputmode="decimal"
            autocomplete="off"
            :aria-invalid="invalid || undefined"
            :aria-describedby="describedby"
        />
        <span
            v-if="ghost && !model"
            class="num-input__ghost"
            aria-hidden="true"
        >
            {{ ghost }}
        </span>
        <span
            class="num-input__unit"
            aria-hidden="true"
        >
            {{ unit }}
        </span>
    </label>
</template>

<style scoped>
.num-input {
    position: relative;
    display: block;
}
.num-input__field {
    width: 100%;
    padding-right: 44px;
    text-align: right;
    font-variant-numeric: tabular-nums;
}
.num-input__ghost,
.num-input__unit {
    position: absolute;
    top: 50%;
    transform: translateY(-50%);
    font-family: var(--font-mono);
    color: var(--ink-muted);
    pointer-events: none;
}
.num-input__ghost {
    right: 44px;
    font-size: 14px;
    font-style: italic;
}
.num-input__unit {
    right: 10px;
    font-size: 12px;
}
</style>
