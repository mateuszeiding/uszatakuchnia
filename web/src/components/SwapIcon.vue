<script lang="ts" setup>
// Ikona zamiennika składnika: dwie strzałki ⇄.
withDefaults(defineProps<{ size?: number }>(), { size: 15 });
</script>

<template>
    <svg
        :width="size"
        :height="size"
        viewBox="0 0 16 16"
        fill="none"
        aria-hidden="true"
    >
        <path
            d="M3 5h9.5M10 2.5L12.5 5 10 7.5M13 11H3.5M6 8.5L3.5 11 6 13.5"
            stroke="currentColor"
            stroke-width="1.4"
            stroke-linecap="round"
            stroke-linejoin="round"
        />
    </svg>
</template>
