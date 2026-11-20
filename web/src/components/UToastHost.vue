<script lang="ts" setup>
import { hideToast, pauseToast, resumeToast, toast, toastMessage } from '@/composables/useToast';

function runAction() {
    const a = toast.value?.action;
    hideToast();
    a?.run();
}
</script>

<template>
    <div
        class="sr-only"
        role="status"
    >
        {{ toastMessage }}
    </div>
    <!-- z akcją („Cofnij”): pełny pasek na dole, klikalny; nie znika pod kursorem ani z fokusem -->
    <div
        v-if="toast?.action"
        class="toast-bar"
        @mouseenter="pauseToast"
        @mouseleave="resumeToast"
        @focusin="pauseToast"
        @focusout="resumeToast"
    >
        <span
            class="toast-bar__text"
            aria-hidden="true"
        >
            {{ toast.message }}
        </span>
        <button
            type="button"
            class="toast-bar__action"
            @click="runAction"
        >
            {{ toast.action.label }}
        </button>
    </div>
    <div
        v-else-if="toastMessage"
        class="toast"
        aria-hidden="true"
    >
        {{ toastMessage }}
    </div>
</template>

<style scoped>
/* prototyp pz-main.jsx (Toast): tło odwrotne do motywu, 12 px od krawędzi, nad safe-area */
.toast-bar {
    position: fixed;
    left: 12px;
    right: 12px;
    bottom: calc(16px + env(safe-area-inset-bottom));
    z-index: var(--z-toast);
    max-width: 440px;
    margin: 0 auto;
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 6px 6px 6px 16px;
    border-radius: var(--r-xl);
    background: var(--ink);
    color: var(--bg);
    box-shadow: 0 12px 32px rgba(10, 10, 15, 0.28);
    font-size: 14px;
    line-height: 1.35;
}
.toast-bar__text {
    flex: 1;
}
.toast-bar__action {
    min-height: 44px;
    padding: 0 14px;
    border: none;
    border-radius: var(--r-lg);
    background: transparent;
    color: var(--toast-action);
    font-family: var(--font-sans);
    font-size: 15px;
    font-weight: 700;
    cursor: pointer;
}
.toast-bar__action:focus-visible {
    outline: 2px solid var(--toast-action);
    outline-offset: -2px;
}
@media (min-width: 901px) {
    .toast-bar {
        bottom: 24px;
    }
}
</style>
