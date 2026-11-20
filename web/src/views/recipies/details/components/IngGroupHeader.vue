<script lang="ts" setup>
import { computed } from 'vue';

import { pluralPL } from '@/shared/utils/quantity';

// Ingredient group header with a counter on the right ("8", after checking off "2/8").
// The whole group checked → "✓" in the accent color and a dimmed header.
const props = defineProps<{ label: string; count: number; done: number }>();

const all = computed(() => props.count > 0 && props.done === props.count);
const srCount = computed(() =>
    props.done
        ? `${props.done} z ${props.count} odhaczonych`
        : `${props.count} ${pluralPL(props.count, 'składnik', 'składniki', 'składników')}`
);
</script>

<template>
    <h3
        class="ing-group"
        :class="{ 'is-done': all }"
    >
        <span class="ing-group__label">
            <span
                v-if="all"
                class="ing-group__tick"
                aria-hidden="true"
            >
                ✓
            </span>
            {{ label }}
        </span>
        <span class="ing-group__count">
            <span aria-hidden="true">{{ done ? `${done}/${count}` : count }}</span>
            <span class="sr-only">{{ srCount }}</span>
        </span>
    </h3>
</template>

<style scoped>
.ing-group {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 8px;
    padding-bottom: 8px;
    border-bottom: 1px solid var(--rule-strong);
    font-family: var(--font-mono);
    font-size: 12px;
    font-weight: 600;
}
.ing-group__label {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    color: var(--ink);
    letter-spacing: 0.5px;
    text-transform: uppercase;
}
.ing-group.is-done .ing-group__label {
    color: var(--ink-muted);
}
.ing-group__tick {
    color: var(--accent);
}
.ing-group__count {
    font-weight: 400;
    color: var(--ink-muted);
}
</style>
