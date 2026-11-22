<script lang="ts" setup>
import { nextTick, onBeforeUnmount, onMounted, ref, useId } from 'vue';

// Sheet from the bottom (phone) or a popover under its button (desktop); Esc, ✕ and the backdrop close it.
// Teleported, so the sticky ingredients panel doesn't clip it. Same behaviour as SubstituteChooser.
const props = defineProps<{ title: string; sheet: boolean; anchor: HTMLElement | null }>();
const emit = defineEmits<{ close: [] }>();

const uid = useId();
const root = ref<HTMLElement>();
const pos = ref<{ left: number; top: number } | null>(null);

function place() {
    if (props.sheet || !props.anchor || !root.value) return;
    const r = props.anchor.getBoundingClientRect();
    const w = root.value.offsetWidth;
    const h = root.value.offsetHeight;
    const left = Math.min(Math.max(8, r.left), window.innerWidth - w - 8);
    let top = r.bottom + 6;
    if (top + h > window.innerHeight - 8) top = Math.max(8, r.top - h - 6);
    pos.value = { left, top };
}

const focusables = () => [
    ...(root.value?.querySelectorAll<HTMLElement>(
        'button:not([disabled]):not([tabindex="-1"]), input:not([disabled])'
    ) ?? []),
];

function trap(e: KeyboardEvent) {
    const list = focusables();
    if (!list.length) return;
    const first = list[0]!;
    const last = list[list.length - 1]!;
    if (!root.value?.contains(document.activeElement)) {
        e.preventDefault();
        first.focus();
    } else if (e.shiftKey && document.activeElement === first) {
        e.preventDefault();
        last.focus();
    } else if (!e.shiftKey && document.activeElement === last) {
        e.preventDefault();
        first.focus();
    }
}

// capture phase: Esc doesn't also close the cooking mode modal underneath
function onKey(e: KeyboardEvent) {
    if (e.key === 'Escape') {
        e.stopPropagation();
        e.preventDefault();
        emit('close');
    } else if (e.key === 'Tab') {
        e.stopPropagation();
        trap(e);
    }
}

function onDown(e: PointerEvent) {
    const t = e.target as Node;
    if (root.value?.contains(t)) return;
    // the opening button toggles the dialog itself
    if (props.anchor?.contains(t)) return;
    emit('close');
}

onMounted(async () => {
    document.addEventListener('keydown', onKey, true);
    document.addEventListener('pointerdown', onDown, true);
    window.addEventListener('resize', place);
    window.addEventListener('scroll', place, true);
    await nextTick();
    place();
    (
        root.value?.querySelector<HTMLElement>('[aria-checked="true"], [aria-selected="true"]') ??
        focusables()[0]
    )?.focus();
});
onBeforeUnmount(() => {
    document.removeEventListener('keydown', onKey, true);
    document.removeEventListener('pointerdown', onDown, true);
    window.removeEventListener('resize', place);
    window.removeEventListener('scroll', place, true);
    props.anchor?.focus();
});
</script>

<template>
    <Teleport to="body">
        <div
            v-if="sheet"
            class="pdlg-scrim"
            aria-hidden="true"
        />
        <div
            ref="root"
            class="pdlg"
            :class="sheet ? 'pdlg--sheet' : 'pdlg--pop'"
            :style="!sheet && pos ? { left: `${pos.left}px`, top: `${pos.top}px` } : undefined"
            role="dialog"
            :aria-modal="sheet || undefined"
            :aria-labelledby="`${uid}-h`"
        >
            <div
                v-if="sheet"
                class="pdlg__grip"
                aria-hidden="true"
            />
            <div class="pdlg__head">
                <h3
                    :id="`${uid}-h`"
                    class="pdlg__title"
                >
                    {{ title }}
                </h3>
                <button
                    v-if="sheet"
                    type="button"
                    class="pdlg__close"
                    aria-label="Zamknij"
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
            <slot />
        </div>
    </Teleport>
</template>

<style scoped>
.pdlg {
    z-index: calc(var(--z-modal) + 10);
    background: var(--bg-alt);
    border: 1px solid var(--rule-strong);
    color: var(--ink);
    display: flex;
    flex-direction: column;
    gap: 14px;
}
.pdlg--pop {
    position: fixed;
    width: 380px;
    max-width: calc(100vw - 16px);
    border-radius: var(--r-xl);
    padding: 16px;
    box-shadow: 0 12px 32px rgba(10, 10, 15, 0.18);
}
.pdlg--sheet {
    position: fixed;
    left: 0;
    right: 0;
    bottom: 0;
    max-height: 85vh;
    overflow-y: auto;
    border-radius: 16px 16px 0 0;
    border-bottom: none;
    padding: 8px 20px calc(20px + env(safe-area-inset-bottom));
    box-shadow: 0 -12px 40px rgba(10, 10, 15, 0.25);
}
.pdlg-scrim {
    position: fixed;
    inset: 0;
    z-index: calc(var(--z-modal) + 9);
    background: rgba(10, 10, 15, 0.45);
}
.pdlg__grip {
    width: 36px;
    height: 4px;
    margin: 0 auto;
    border-radius: var(--r-pill);
    background: var(--rule-strong);
}
.pdlg__head {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 8px;
}
.pdlg__title {
    font-size: 16px;
    font-weight: 600;
}
.pdlg--sheet .pdlg__title {
    font-size: 18px;
}
.pdlg__close {
    width: 44px;
    height: 44px;
    margin: -10px -12px -10px 0;
    border: none;
    background: none;
    color: var(--ink-muted);
    display: inline-flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    flex-shrink: 0;
}
</style>
