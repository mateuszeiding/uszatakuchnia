<script lang="ts" setup>
withDefaults(defineProps<{ size?: number; filled?: boolean }>(), { size: 16, filled: false });
</script>

<template>
    <svg
        :width="size"
        :height="size"
        viewBox="0 0 16 16"
        aria-hidden="true"
        class="flame"
    >
        <path
            d="M8 15c-2.7 0-4.6-1.9-4.6-4.5 0-2 1.1-3.3 2.1-4.3.2 1.2.7 2 1.5 2.4C6.8 6 7.5 3.3 9.4 1.2c.3 2.4 1.6 3.6 2.6 5 .8 1.2 1.3 2.5 1.3 4C13.3 13 11.1 15 8 15z"
            :fill="filled ? 'currentColor' : 'none'"
            stroke="currentColor"
            stroke-width="1.2"
            stroke-linejoin="round"
        />
        <path
            v-if="filled"
            d="M8 13.6c-1.2 0-2-.9-2-2 0-1 .6-1.6 1.1-2.2.1.6.4 1 .8 1.2.2-1 .7-1.8 1.4-2.4.2 1 .8 1.6 1.1 2.3.2.4.3.8.3 1.2 0 1.1-1.2 1.9-2.7 1.9z"
            fill="rgba(255,255,255,0.55)"
        />
    </svg>
</template>

<style scoped>
.flame {
    flex-shrink: 0;
}
</style>
