<script lang="ts" setup>
// Proste ikony liniowe (dekoracyjne — przycisk niesie aria-label).
defineProps<{
    name: 'edit' | 'trash' | 'eye' | 'x' | 'check' | 'lock' | 'external';
    size?: number;
}>();
</script>

<template>
    <svg
        :width="size ?? 13"
        :height="size ?? 13"
        viewBox="0 0 13 13"
        aria-hidden="true"
        stroke="currentColor"
        stroke-width="1.3"
        stroke-linecap="round"
        stroke-linejoin="round"
    >
        <path
            v-if="name === 'edit'"
            d="M2 9.5L8.5 3l1.5 1.5L3.5 11H2V9.5z"
        />
        <path
            v-else-if="name === 'trash'"
            d="M2.5 3.5h8M5 3.5V2.5h3v1M4 3.5l.4 7h4.2l.4-7"
        />
        <template v-else-if="name === 'eye'">
            <path d="M1 6.5s2-3.7 5.5-3.7S12 6.5 12 6.5s-2 3.7-5.5 3.7S1 6.5 1 6.5z" />
            <circle
                cx="6.5"
                cy="6.5"
                r="1.7"
            />
        </template>
        <path
            v-else-if="name === 'x'"
            d="M3 3l7 7M10 3l-7 7"
        />
        <path
            v-else-if="name === 'check'"
            d="M2.5 6.8L5.2 9.5 10.5 3.8"
        />
        <template v-else-if="name === 'lock'">
            <rect
                x="2.5"
                y="5.8"
                width="8"
                height="5.7"
                rx="1.2"
            />
            <path d="M4.5 5.8V4.1a2 2 0 014 0v1.7" />
        </template>
        <!-- nowa karta -->
        <path
            v-else-if="name === 'external'"
            d="M5.5 2.5H3a.5.5 0 00-.5.5v7a.5.5 0 00.5.5h7a.5.5 0 00.5-.5V7.5M7.5 2.5h3v3M10.5 2.5L6 7"
        />
    </svg>
</template>
