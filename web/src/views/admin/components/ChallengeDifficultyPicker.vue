<script lang="ts" setup>
import FlameIcon from '@/components/FlameIcon.vue';
import { DIFFICULTIES, levelName } from '@/shared/utils/challenges';

// Challenge difficulty in the form: a radiogroup of 5 buttons (flames + "3 · średnie").
// Keyboard as in a radio group: Tab lands on the selected one, arrows change the choice.
defineProps<{ labelledby: string; invalid?: boolean }>();
const level = defineModel<number | null>({ required: true });

function onKey(e: KeyboardEvent, n: number) {
    const step = { ArrowRight: 1, ArrowDown: 1, ArrowLeft: -1, ArrowUp: -1 }[e.key];
    if (!step) return;
    e.preventDefault();
    const next = ((n - 1 + step + 5) % 5) + 1;
    level.value = next;
    const group = (e.currentTarget as HTMLElement).parentElement;
    (group?.children[next - 1] as HTMLElement | undefined)?.focus();
}
</script>

<template>
    <div
        role="radiogroup"
        class="dp"
        :class="{ 'is-invalid': invalid }"
        :aria-labelledby="labelledby"
        :aria-invalid="invalid || undefined"
    >
        <button
            v-for="n in DIFFICULTIES"
            :key="n"
            type="button"
            role="radio"
            class="dp__opt"
            :class="{ 'is-on': level === n }"
            :aria-checked="level === n"
            :tabindex="level === n || (!level && n === 1) ? 0 : -1"
            @click="level = n"
            @keydown="onKey($event, n)"
        >
            <span
                class="dp__flames"
                aria-hidden="true"
            >
                <FlameIcon
                    v-for="i in DIFFICULTIES"
                    :key="i"
                    :size="12"
                    :filled="i <= n"
                    :class="i <= n ? 'dp__on' : 'dp__off'"
                />
            </span>
            <span class="dp__label">{{ n }} · {{ levelName(n) }}</span>
        </button>
    </div>
</template>

<style scoped>
.dp {
    display: flex;
    flex-wrap: wrap;
    gap: 6px;
}
.dp__opt {
    display: inline-flex;
    flex-direction: column;
    align-items: center;
    gap: 4px;
    min-width: 104px;
    padding: 8px 12px;
    border-radius: var(--r-lg);
    border: 1px solid var(--rule-strong);
    background: var(--bg);
    font-family: var(--font-sans);
    cursor: pointer;
}
.dp__opt:hover {
    border-color: var(--accent);
}
.dp__opt.is-on {
    border-color: var(--accent-fill);
    background: var(--accent-soft);
}
.dp.is-invalid .dp__opt {
    border-color: var(--danger);
}
.dp__flames {
    display: inline-flex;
    gap: 1px;
}
.dp__on {
    color: var(--c-tomato);
}
.dp__off {
    color: var(--rule-strong);
}
.dp__label {
    font-size: 13px;
    font-weight: 500;
    color: var(--ink-dim);
}
.dp__opt.is-on .dp__label {
    font-weight: 600;
    color: var(--accent);
}
</style>
