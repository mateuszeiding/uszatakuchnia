<script lang="ts" setup>
import { computed } from 'vue';

import { DISH_ICONS, type DishIconKey } from '@/shared/dishIcons';

// Decorative: the dish type is always named next to it (the badge), so the icon is aria-hidden.
const props = defineProps<{ icon: DishIconKey }>();

const shapes = computed(() => DISH_ICONS[props.icon]);
</script>

<template>
    <svg
        viewBox="0 0 64 64"
        fill="none"
        stroke="currentColor"
        stroke-width="2.4"
        stroke-linecap="round"
        stroke-linejoin="round"
        aria-hidden="true"
        focusable="false"
    >
        <template
            v-for="(s, i) in shapes"
            :key="i"
        >
            <path
                v-if="'d' in s"
                :d="s.d"
            />
            <circle
                v-else
                :cx="s.cx"
                :cy="s.cy"
                :r="s.r"
            />
        </template>
    </svg>
</template>
