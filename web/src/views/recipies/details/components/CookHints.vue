<script lang="ts" setup>
import { computed, useId } from 'vue';

import StepHints from './StepHints.vue';

import BulbIcon from '@/components/BulbIcon.vue';
import { pluralPL } from '@/shared/utils/quantity';

// Hints in cooking mode. When they don't fit next to the step text, they collapse into one button;
// the parent keeps "open" for the rest of the session, so the reader taps it once.
const props = defineProps<{ hints: string[]; collapsible: boolean; open: boolean }>();
const emit = defineEmits<{ toggle: [] }>();

const id = useId();
const label = computed(
    () =>
        `${props.hints.length} ${pluralPL(props.hints.length, 'wskazówka', 'wskazówki', 'wskazówek')}`
);
</script>

<template>
    <template v-if="hints.length">
        <button
            v-if="collapsible"
            type="button"
            class="cook-hints__toggle"
            :aria-expanded="open"
            :aria-controls="id"
            @click="emit('toggle')"
        >
            <span class="cook-hints__bulb"><BulbIcon :size="16" /></span>
            {{ label }}
            <svg
                class="cook-hints__chev"
                width="12"
                height="12"
                viewBox="0 0 12 12"
                fill="none"
                aria-hidden="true"
            >
                <path
                    d="M3 4.5L6 7.5 9 4.5"
                    stroke="currentColor"
                    stroke-width="1.6"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                />
            </svg>
        </button>
        <div
            v-if="!collapsible || open"
            :id="id"
        >
            <StepHints
                :hints="hints"
                size="lg"
            />
        </div>
    </template>
</template>

<style scoped>
.cook-hints__toggle {
    margin-top: 16px;
    min-height: 44px;
    display: inline-flex;
    align-items: center;
    gap: 8px;
    padding: 0 14px;
    border-radius: var(--r-pill);
    border: 1px solid var(--rule-strong);
    background: var(--bg-alt);
    color: var(--ink);
    font-size: 15px;
    font-weight: 600;
    cursor: pointer;
}
.cook-hints__bulb {
    display: inline-flex;
    color: var(--accent);
}
.cook-hints__toggle[aria-expanded='true'] .cook-hints__chev {
    transform: rotate(180deg);
}
</style>
