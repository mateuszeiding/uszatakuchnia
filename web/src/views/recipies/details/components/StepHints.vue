<script lang="ts" setup>
// Hints below the step text: recipe page (md) and cooking mode (lg). No hints = nothing.
import BulbIcon from '@/components/BulbIcon.vue';

withDefaults(defineProps<{ hints: string[]; done?: boolean; size?: 'md' | 'lg' }>(), {
    done: false,
    size: 'md',
});
</script>

<template>
    <ul
        v-if="hints.length"
        class="hints"
        :class="[`hints--${size}`, { 'is-done': done }]"
        aria-label="Wskazówki"
    >
        <li
            v-for="(h, i) in hints"
            :key="i"
            class="hints__row"
        >
            <span class="hints__icon"><BulbIcon :size="size === 'lg' ? 18 : 15" /></span>
            <span class="hints__text">{{ h }}</span>
        </li>
    </ul>
</template>

<style scoped>
.hints {
    list-style: none;
    margin: 14px 0 0;
    padding: 12px 0 0;
    border-top: 1px solid var(--rule);
    display: flex;
    flex-direction: column;
    gap: 8px;
}
.hints__row {
    display: grid;
    grid-template-columns: 16px minmax(0, 1fr);
    gap: 10px;
    align-items: start;
}
.hints__icon {
    display: inline-flex;
    padding-top: 3px;
    color: var(--accent);
}
.hints__text {
    font-size: 15px;
    line-height: 1.5;
    color: var(--ink-dim);
    text-wrap: pretty;
}
.hints.is-done .hints__icon,
.hints.is-done .hints__text {
    color: var(--ink-muted);
}

.hints--lg {
    margin-top: 20px;
    padding-top: 16px;
    gap: 12px;
}
.hints--lg .hints__row {
    grid-template-columns: 20px minmax(0, 1fr);
    gap: 12px;
}
.hints--lg .hints__icon {
    padding-top: 5px;
}
.hints--lg .hints__text {
    font-size: clamp(18px, 2.2vw, 22px);
}
</style>
