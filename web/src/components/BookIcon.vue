<script lang="ts" setup>
// Ikona linku do przepisu: książka; external = „otwiera się w nowej karcie” (tryb gotowania).
withDefaults(defineProps<{ size?: number; external?: boolean }>(), { size: 15 });
</script>

<template>
    <svg
        v-if="external"
        :width="size"
        :height="size"
        viewBox="0 0 12 12"
        fill="none"
        aria-hidden="true"
    >
        <path
            d="M5 2H2.5A.5.5 0 002 2.5v7a.5.5 0 00.5.5h7a.5.5 0 00.5-.5V7M7 2h3v3M10 2L5.5 6.5"
            stroke="currentColor"
            stroke-width="1.3"
            stroke-linecap="round"
            stroke-linejoin="round"
        />
    </svg>
    <svg
        v-else
        :width="size"
        :height="size"
        viewBox="0 0 14 14"
        fill="none"
        aria-hidden="true"
    >
        <path
            d="M2 2.5h3.5A1.5 1.5 0 017 4v8a1.2 1.2 0 00-1.2-1.2H2v-8.3zM12 2.5H8.5A1.5 1.5 0 007 4v8a1.2 1.2 0 011.2-1.2H12v-8.3z"
            stroke="currentColor"
            stroke-width="1.2"
            stroke-linejoin="round"
        />
    </svg>
</template>
