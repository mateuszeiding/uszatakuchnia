<script lang="ts" setup>
// Ikona wskazówki: żarówka (dekoracyjna — tekst wskazówki stoi obok).
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
            d="M6 12.5h4M6.5 14.5h3M8 1.5a4.5 4.5 0 00-2.6 8.2c.4.3.6.7.6 1.1v.2h4v-.2c0-.4.2-.8.6-1.1A4.5 4.5 0 008 1.5z"
            stroke="currentColor"
            stroke-width="1.3"
            stroke-linejoin="round"
        />
    </svg>
</template>
