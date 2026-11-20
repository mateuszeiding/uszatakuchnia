<script lang="ts" setup>
import USpinner from './USpinner.vue';

// Button label with a busy state: both variants ("Zapisz zmiany" / spinner + "Zapisywanie…")
// share one grid cell and the hidden one holds the space, so the button keeps its width
// and nothing next to it shifts.
withDefaults(
    defineProps<{
        busy: boolean;
        label: string;
        /** Text while busy; defaults to label (only the spinner is added). */
        busyLabel?: string;
        spinnerColor?: string;
    }>(),
    { busyLabel: undefined, spinnerColor: undefined }
);
</script>

<template>
    <span class="busy-label">
        <span
            class="busy-label__part"
            :class="{ 'is-hidden': busy }"
            :aria-hidden="busy"
        >
            {{ label }}
        </span>
        <span
            class="busy-label__part"
            :class="{ 'is-hidden': !busy }"
            :aria-hidden="!busy"
        >
            <USpinner
                :size="14"
                :color="spinnerColor"
            />
            {{ busyLabel ?? label }}
        </span>
    </span>
</template>

<style scoped>
.busy-label {
    display: inline-grid;
}
.busy-label__part {
    grid-area: 1 / 1;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
    white-space: nowrap;
}
.busy-label__part.is-hidden {
    visibility: hidden;
}
</style>
