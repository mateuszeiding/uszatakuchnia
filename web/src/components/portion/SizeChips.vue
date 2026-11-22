<script lang="ts" setup>
import { ref } from 'vue';

import type { IPieceSize } from '@/shared/utils/portion';

// Only the author's sizes (owner's decision 22.11: no custom diameter, dough never rolls out exactly).
const props = defineProps<{ sizes: IPieceSize[]; value: number; label: string; big?: boolean }>();
const emit = defineEmits<{ pick: [d: number] }>();

const root = ref<HTMLElement>();
// arrow keys move the choice like native radios (one tab stop for the group)
function onKey(e: KeyboardEvent, i: number) {
    const d = { ArrowRight: 1, ArrowDown: 1, ArrowLeft: -1, ArrowUp: -1 }[e.key];
    if (!d) return;
    e.preventDefault();
    const n = props.sizes.length;
    const next = props.sizes[(i + d + n) % n]!;
    emit('pick', next.diameterCm);
    root.value?.querySelectorAll<HTMLElement>('[role="radio"]')[(i + d + n) % n]?.focus();
}
</script>

<template>
    <div
        ref="root"
        class="schips"
        role="radiogroup"
        :aria-label="label"
    >
        <button
            v-for="(s, i) in sizes"
            :key="s.diameterCm"
            type="button"
            role="radio"
            class="schips__chip"
            :class="{ 'is-on': s.diameterCm === value, 'schips__chip--big': big }"
            :aria-checked="s.diameterCm === value"
            :tabindex="s.diameterCm === value ? 0 : -1"
            @click="emit('pick', s.diameterCm)"
            @keydown="onKey($event, i)"
        >
            {{ s.diameterCm }} cm
        </button>
    </div>
</template>

<style scoped>
.schips {
    display: flex;
    flex-wrap: wrap;
    gap: 6px;
}
.schips__chip {
    min-height: 34px;
    padding: 0 12px;
    border-radius: var(--r-pill);
    border: 1px solid var(--rule-strong);
    background: var(--bg-alt);
    color: var(--ink);
    font-size: 14px;
    font-weight: 500;
    cursor: pointer;
}
.schips__chip--big {
    min-height: 44px;
}
.schips__chip:hover {
    border-color: var(--ink-muted);
}
.schips__chip.is-on {
    border-color: var(--accent);
    background: var(--accent-soft);
    color: var(--accent);
    font-weight: 600;
}
</style>
