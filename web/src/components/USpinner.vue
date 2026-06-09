<script setup lang="ts">
withDefaults(
    defineProps<{
        size?: number;
        color?: string;
    }>(),
    {
        size: 16,
        color: undefined,
    }
);
</script>

<template>
    <span
        class="u-spinner"
        :style="{ '--sz': `${size}px`, '--color': color }"
    />
</template>

<style scoped>
.u-spinner {
    --sz: 16px;
    --color: var(--accent);
    --track: color-mix(in srgb, var(--color) 20%, transparent);

    display: inline-block;
    flex-shrink: 0;
    width: var(--sz);
    height: var(--sz);
    border-radius: 50%;
    border: max(1.5px, calc(var(--sz) / 10)) solid var(--track);
    border-top-color: var(--color);
    animation: spin 0.65s linear infinite;
}

@keyframes spin {
    to {
        transform: rotate(360deg);
    }
}
</style>
