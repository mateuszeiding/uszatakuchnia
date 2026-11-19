<script lang="ts" setup>
// Kość do gry — losowanie wyzwania.
withDefaults(defineProps<{ size?: number }>(), { size: 16 });
</script>

<template>
    <svg
        :width="size"
        :height="size"
        viewBox="0 0 16 16"
        fill="none"
        aria-hidden="true"
    >
        <rect
            x="2"
            y="2"
            width="12"
            height="12"
            rx="3"
            stroke="currentColor"
            stroke-width="1.4"
        />
        <circle
            cx="5.5"
            cy="5.5"
            r="1.1"
            fill="currentColor"
        />
        <circle
            cx="10.5"
            cy="10.5"
            r="1.1"
            fill="currentColor"
        />
        <circle
            cx="8"
            cy="8"
            r="1.1"
            fill="currentColor"
        />
    </svg>
</template>
