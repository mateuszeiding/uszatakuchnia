<script lang="ts" setup>
import USpinner from './USpinner.vue';

// Etykieta przycisku z trybem zajętości: oba warianty („Zapisz zmiany” / spinner + „Zapisywanie…”)
// leżą w tej samej komórce siatki, niewidoczny trzyma miejsce — przycisk ma zawsze tę samą szerokość
// i nic obok się nie przesuwa.
withDefaults(
    defineProps<{
        busy: boolean;
        label: string;
        /** Tekst w trakcie; domyślnie ten sam co label (sam spinner dochodzi). */
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
