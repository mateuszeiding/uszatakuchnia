<script lang="ts" setup>
import { onBeforeUnmount, onMounted, ref } from 'vue';

import { recipesLabel } from '../listFilters';

// Filtry na telefonie: bottom sheet z pułapką fokusu, Esc zamyka.
defineProps<{ count: number; activeCount: number }>();
const emit = defineEmits<{ close: []; clear: [] }>();

const root = ref<HTMLElement>();
let prevOverflow = '';
let prevFocus: HTMLElement | null = null;

function onKey(e: KeyboardEvent) {
    if (e.key === 'Escape') {
        emit('close');
        return;
    }
    if (e.key !== 'Tab' || !root.value) return;
    const list = [
        ...root.value.querySelectorAll<HTMLElement>(
            'button:not([disabled]),input,select,[tabindex="0"]'
        ),
    ];
    const first = list[0];
    const last = list[list.length - 1];
    if (!first || !last) return;
    if (e.shiftKey && document.activeElement === first) {
        e.preventDefault();
        last.focus();
    } else if (!e.shiftKey && document.activeElement === last) {
        e.preventDefault();
        first.focus();
    }
}

onMounted(() => {
    prevOverflow = document.body.style.overflow;
    document.body.style.overflow = 'hidden';
    prevFocus = document.activeElement as HTMLElement | null;
    root.value?.querySelector<HTMLElement>('button')?.focus();
    document.addEventListener('keydown', onKey);
});
onBeforeUnmount(() => {
    document.body.style.overflow = prevOverflow;
    document.removeEventListener('keydown', onKey);
    prevFocus?.focus?.();
});
</script>

<template>
    <Teleport to="body">
        <div
            class="sheet-backdrop"
            @click="emit('close')"
        />
        <div
            ref="root"
            class="sheet"
            role="dialog"
            aria-modal="true"
            aria-labelledby="sheet-title"
        >
            <div class="sheet__head">
                <h2
                    id="sheet-title"
                    class="sheet__title"
                >
                    Filtry
                </h2>
                <button
                    type="button"
                    class="sheet__close"
                    aria-label="Zamknij filtry"
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
            <div class="sheet__body">
                <slot />
            </div>
            <div class="sheet__foot">
                <button
                    type="button"
                    class="btn btn--ghost sheet__btn"
                    :disabled="!activeCount"
                    @click="emit('clear')"
                >
                    Wyczyść
                </button>
                <button
                    type="button"
                    class="btn btn--accent sheet__btn sheet__btn--main"
                    @click="emit('close')"
                >
                    Pokaż {{ count }} {{ recipesLabel(count) }}
                </button>
            </div>
        </div>
    </Teleport>
</template>

<style scoped>
.sheet-backdrop {
    position: fixed;
    inset: 0;
    z-index: var(--z-overlay);
    background: rgba(10, 10, 15, 0.45);
}
.sheet {
    position: fixed;
    left: 0;
    right: 0;
    bottom: 0;
    z-index: var(--z-modal);
    max-height: 88vh;
    background: var(--bg-alt);
    border-radius: 16px 16px 0 0;
    display: flex;
    flex-direction: column;
    box-shadow: 0 -12px 40px rgba(10, 10, 15, 0.2);
}
.sheet__head {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 14px 20px;
    border-bottom: 1px solid var(--rule);
}
.sheet__title {
    font-size: 18px;
    font-weight: 600;
    letter-spacing: 0;
}
.sheet__close {
    width: 44px;
    height: 44px;
    border-radius: var(--r-pill);
    border: 1px solid var(--rule-strong);
    background: var(--bg);
    color: var(--ink);
    cursor: pointer;
    display: inline-flex;
    align-items: center;
    justify-content: center;
}
.sheet__body {
    overflow-y: auto;
    padding: 18px 20px 4px;
}
.sheet__foot {
    display: flex;
    gap: 10px;
    padding: 12px 20px calc(12px + env(safe-area-inset-bottom));
    border-top: 1px solid var(--rule);
}
.sheet__btn {
    height: 48px;
}
.sheet__btn--main {
    flex: 1;
    font-size: 15px;
}
</style>
