<script lang="ts" setup>
import FlameIcon from '@/components/FlameIcon.vue';
import { DIFFICULTIES, levelName } from '@/shared/utils/challenges';

// Trudność wyzwania 1–5: płomienie (wypełnione --c-tomato, puste obrysem --rule-strong) + podpis słowny.
// Ikony są aria-hidden; czytnik dostaje „trudność 3 z 5”. Bez podpisu (tabela w panelu) — sam tekst dla czytnika.
withDefaults(defineProps<{ level: number; size?: number; label?: boolean; compact?: boolean }>(), {
    size: 16,
    label: true,
    compact: false,
});
</script>

<template>
    <span
        class="diff"
        :class="{ 'diff--compact': compact }"
    >
        <span
            class="diff__icons"
            aria-hidden="true"
        >
            <FlameIcon
                v-for="i in DIFFICULTIES"
                :key="i"
                :size="size"
                :filled="i <= level"
                :class="i <= level ? 'diff__on' : 'diff__off'"
            />
        </span>
        <span
            v-if="label"
            class="diff__label"
        >
            {{ levelName(level) }}
            <span class="sr-only">, trudność {{ level }} z 5</span>
        </span>
        <span
            v-else
            class="sr-only"
        >
            Trudność {{ level }} z 5: {{ levelName(level) }}
        </span>
    </span>
</template>

<style scoped>
.diff {
    display: inline-flex;
    align-items: center;
    gap: 8px;
    white-space: nowrap;
}
.diff--compact {
    gap: 6px;
}
.diff__icons {
    display: inline-flex;
    align-items: center;
    gap: 1px;
}
.diff__on {
    color: var(--c-tomato);
}
.diff__off {
    color: var(--rule-strong);
}
.diff__label {
    font-size: 14px;
    font-weight: 500;
    color: var(--ink-dim);
}
.diff--compact .diff__label {
    font-size: 13px;
}
</style>
