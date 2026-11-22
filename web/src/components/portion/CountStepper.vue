<script lang="ts" setup>
// − value + for the number of pieces (the servings counter keeps its own ServingsStepper).
const props = withDefaults(
    defineProps<{ value: number; label: string; min?: number; max?: number; big?: boolean }>(),
    { min: 1, max: 20, big: false }
);
const emit = defineEmits<{ set: [n: number] }>();
const step = (d: number) => emit('set', Math.min(props.max, Math.max(props.min, props.value + d)));
</script>

<template>
    <span
        class="cstep"
        :class="{ 'cstep--big': big }"
        role="group"
        :aria-label="label"
    >
        <button
            type="button"
            class="cstep__btn"
            :aria-label="`Mniej: ${label.toLowerCase()}`"
            :disabled="value <= min"
            @click="step(-1)"
        >
            −
        </button>
        <output
            class="cstep__value"
            aria-live="polite"
        >
            {{ value }}
        </output>
        <button
            type="button"
            class="cstep__btn"
            :aria-label="`Więcej: ${label.toLowerCase()}`"
            :disabled="value >= max"
            @click="step(1)"
        >
            +
        </button>
    </span>
</template>

<style scoped>
.cstep {
    display: inline-flex;
    align-items: center;
    gap: 4px;
}
.cstep__btn {
    width: 36px;
    height: 36px;
    border-radius: var(--r-md);
    background: var(--bg-alt);
    border: 1px solid var(--rule);
    color: var(--ink);
    font-size: 18px;
    cursor: pointer;
}
.cstep--big .cstep__btn {
    width: 44px;
    height: 44px;
}
.cstep__btn:disabled {
    background: transparent;
    color: var(--ink-faint);
    cursor: default;
}
.cstep__value {
    min-width: 36px;
    text-align: center;
    font-size: 22px;
    font-weight: 600;
}
</style>
