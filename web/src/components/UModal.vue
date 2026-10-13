<script lang="ts" setup>
import { onBeforeUnmount, onMounted, ref, useId } from 'vue';

// Modal: role="dialog", pułapka fokusu, Esc zamyka, fokus wraca na element, który go otworzył.
// Element z atrybutem autofocus dostaje fokus na starcie (inaczej pierwsza kontrolka).
const props = withDefaults(
    defineProps<{
        title: string;
        eyebrow?: string;
        width?: number;
    }>(),
    { eyebrow: undefined, width: 440 }
);
const emit = defineEmits<{ close: [] }>();

const titleId = useId();
const root = ref<HTMLElement>();
let previous: HTMLElement | null = null;

const FOCUSABLE = 'button,[href],input,select,textarea,[tabindex]:not([tabindex="-1"])';

function focusables(): HTMLElement[] {
    if (!root.value) return [];
    return [...root.value.querySelectorAll<HTMLElement>(FOCUSABLE)].filter(
        (el) => !el.hasAttribute('disabled') && el.offsetParent !== null
    );
}

function onKey(e: KeyboardEvent) {
    if (e.key === 'Escape') {
        e.stopPropagation();
        emit('close');
        return;
    }
    if (e.key !== 'Tab') return;
    const list = focusables();
    if (!list.length) return;
    const first = list[0]!;
    const last = list[list.length - 1]!;
    if (e.shiftKey && document.activeElement === first) {
        e.preventDefault();
        last.focus();
    } else if (!e.shiftKey && document.activeElement === last) {
        e.preventDefault();
        first.focus();
    }
}

onMounted(() => {
    previous = document.activeElement as HTMLElement | null;
    const auto = root.value?.querySelector<HTMLElement>('[autofocus]');
    (auto ?? focusables()[0])?.focus();
    document.addEventListener('keydown', onKey);
});

onBeforeUnmount(() => {
    document.removeEventListener('keydown', onKey);
    previous?.focus?.();
});
</script>

<template>
    <Teleport to="body">
        <div
            class="modal-backdrop"
            @click="emit('close')"
        />
        <div
            ref="root"
            class="modal"
            role="dialog"
            aria-modal="true"
            :aria-labelledby="titleId"
            :style="{ '--modal-w': `${props.width}px` }"
        >
            <div>
                <div
                    v-if="eyebrow"
                    class="modal__eyebrow"
                >
                    {{ eyebrow }}
                </div>
                <h2
                    :id="titleId"
                    class="modal__title"
                >
                    {{ title }}
                </h2>
            </div>
            <slot />
            <div
                v-if="$slots.footer"
                class="modal__footer"
            >
                <slot name="footer" />
            </div>
        </div>
    </Teleport>
</template>
