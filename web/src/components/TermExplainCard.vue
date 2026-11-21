<script lang="ts" setup>
import BookIcon from '@/components/BookIcon.vue';
import type { ExplainMode } from '@/composables/useTermExplain';

// The explanation of a term: name + text. The host places it (popover / sheet); the panel shows it as a preview.
withDefaults(
    defineProps<{
        name: string;
        explanation: string;
        mode?: ExplainMode;
        headingId?: string;
        /** Panel preview: no close button. */
        preview?: boolean;
    }>(),
    { mode: 'popover', headingId: undefined, preview: false }
);
const emit = defineEmits<{ close: [] }>();
</script>

<template>
    <div
        class="explain"
        :class="`explain--${mode}`"
    >
        <div
            v-if="mode !== 'popover'"
            class="explain__grip"
            aria-hidden="true"
        />
        <div class="explain__head">
            <h3
                :id="headingId"
                class="explain__name"
            >
                <span class="explain__icon"><BookIcon :size="mode === 'cook' ? 18 : 15" /></span>
                {{ name }}
            </h3>
            <button
                v-if="mode !== 'cook' && !preview"
                type="button"
                class="explain__close"
                aria-label="Zamknij wyjaśnienie"
                @click="emit('close')"
            >
                <svg
                    width="11"
                    height="11"
                    viewBox="0 0 11 11"
                    aria-hidden="true"
                >
                    <path
                        d="M2 2l7 7M9 2L2 9"
                        stroke="currentColor"
                        stroke-width="1.5"
                        stroke-linecap="round"
                    />
                </svg>
            </button>
        </div>
        <!-- a long text scrolls here; the name and ✕ stay. Below the text: room for "Więcej w Wiki →" (later). -->
        <div class="explain__body">
            <p class="explain__text">{{ explanation }}</p>
        </div>
        <button
            v-if="mode === 'cook'"
            type="button"
            class="explain__done"
            @click="emit('close')"
        >
            Zamknij
        </button>
    </div>
</template>

<style scoped>
.explain {
    color: var(--ink);
}
.explain__grip {
    width: 36px;
    height: 4px;
    margin: 0 auto 10px;
    border-radius: var(--r-pill);
    background: var(--rule-strong);
}
.explain__head {
    display: flex;
    align-items: flex-start;
    gap: 8px;
}
.explain__name {
    flex: 1;
    margin: 0;
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 16px;
    font-weight: 600;
    letter-spacing: -0.3px;
    line-height: 1.3;
}
.explain__icon {
    display: inline-flex;
    color: var(--accent);
}
.explain__close {
    width: 32px;
    height: 32px;
    margin: -6px -8px 0 0;
    border: none;
    background: none;
    color: var(--ink-muted);
    display: inline-flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    flex-shrink: 0;
}
.explain__body {
    max-height: 200px;
    overflow-y: auto;
    margin-top: 8px;
}
.explain__text {
    margin: 0;
    font-size: 15px;
    line-height: 1.5;
    color: var(--ink-dim);
    text-wrap: pretty;
}

.explain--sheet .explain__name {
    font-size: 18px;
}
.explain--sheet .explain__close {
    width: 44px;
    height: 44px;
    margin: -10px -10px 0 0;
}
.explain--sheet .explain__body {
    max-height: 280px;
}
.explain--sheet .explain__text {
    font-size: 16px;
}

.explain--cook .explain__name {
    font-size: 22px;
}
.explain--cook .explain__body {
    max-height: 260px;
}
.explain--cook .explain__text {
    font-size: 20px;
}
.explain__done {
    margin-top: 16px;
    width: 100%;
    min-height: 64px;
    border-radius: 14px;
    border: 1px solid var(--rule-strong);
    background: var(--bg);
    color: var(--ink);
    font-size: 19px;
    font-weight: 600;
    cursor: pointer;
}
</style>
