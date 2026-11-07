<script lang="ts" setup>
import { computed } from 'vue';

import {
    fmt,
    nf,
    type NutrientKey,
    type Range,
    type RoundMode,
    spokenText,
    type SumRange,
} from '@/shared/utils/nutrition';

// „20 – 22 g (~21)” — liczby w Geist Mono, średnia przygaszona; czytnik słyszy pełne zdanie.
const props = withDefaults(
    defineProps<{
        value: Range | SumRange | null | undefined;
        nutrient: NutrientKey;
        mode: RoundMode;
        unit?: string;
        size?: 'sm' | 'md';
        muted?: boolean;
    }>(),
    { unit: undefined, size: 'md', muted: false }
);

const f = computed(() => fmt(props.value, props.nutrient, props.mode));
const unitLabel = computed(() => props.unit ?? (props.nutrient === 'kcal' ? 'kcal' : 'g'));
</script>

<template>
    <span
        class="macro-val"
        :class="[`macro-val--${size}`, { 'is-muted': muted }]"
    >
        <template v-if="f">
            <span class="sr-only">{{ spokenText(f, unitLabel) }}</span>
            <span
                v-if="f.single"
                aria-hidden="true"
            >
                {{ nf(f.a) }} {{ unitLabel }}
            </span>
            <span
                v-else
                aria-hidden="true"
            >
                {{ nf(f.a) }} – {{ nf(f.b) }} {{ unitLabel }}
                <span class="macro-val__avg">(~{{ nf(f.m) }})</span>
            </span>
        </template>
        <span
            v-else
            class="macro-val__empty"
        >
            —
        </span>
    </span>
</template>

<style scoped>
.macro-val {
    font-family: var(--font-mono);
    font-weight: 500;
    font-variant-numeric: tabular-nums;
    white-space: nowrap;
    color: var(--ink);
}
.macro-val--sm {
    font-size: 12px;
}
.macro-val--md {
    font-size: 13px;
}
.macro-val.is-muted {
    font-weight: 400;
    color: var(--ink-dim);
}
.macro-val__avg,
.macro-val__empty {
    font-weight: 400;
    color: var(--ink-muted);
}
</style>
