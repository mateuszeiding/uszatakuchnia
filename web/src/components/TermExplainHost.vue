<script lang="ts" setup>
import { computed, nextTick, onBeforeUnmount, ref, useId, watch } from 'vue';
import { useRoute } from 'vue-router';

import TermExplainCard from '@/components/TermExplainCard.vue';
import { closeTerm, openTerm } from '@/composables/useTermExplain';

// The one open term explanation (App.vue). Popover by the word on desktop, a bottom sheet on a phone and in
// cooking mode. Esc / ✕ / "Zamknij" / the backdrop close it and return focus to the word; on desktop a click
// elsewhere closes it too. Keys are caught in the capture phase, so Esc doesn't also close cooking mode.
const uid = useId();
const box = ref<HTMLElement>();
const pos = ref<{ left: number; top: number } | null>(null);
const cur = computed(() => openTerm.value);
const sheet = computed(() => !!cur.value && cur.value.mode !== 'popover');

const GAP = 8;
function place() {
    const o = cur.value;
    if (!o || o.mode !== 'popover' || !box.value) return;
    // a word that wraps has several boxes: anchor to the first line
    const r = o.anchor.getClientRects()[0] ?? o.anchor.getBoundingClientRect();
    const w = box.value.offsetWidth;
    const h = box.value.offsetHeight;
    const left = Math.min(Math.max(GAP, r.left), window.innerWidth - w - GAP);
    const below = r.bottom + GAP;
    const top =
        below + h > window.innerHeight - GAP && r.top - h - GAP > GAP ? r.top - h - GAP : below;
    pos.value = { left, top };
}

function trap(e: KeyboardEvent) {
    const list = [...(box.value?.querySelectorAll<HTMLElement>('button') ?? [])];
    if (!list.length) return;
    const first = list[0]!;
    const last = list[list.length - 1]!;
    if (!box.value?.contains(document.activeElement)) {
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

function onKey(e: KeyboardEvent) {
    if (e.key === 'Escape') {
        e.stopPropagation();
        e.preventDefault();
        closeTerm(true);
    } else if (sheet.value && e.key === 'Tab') {
        e.stopPropagation();
        trap(e);
    } else if (sheet.value && (e.key === 'ArrowLeft' || e.key === 'ArrowRight')) {
        // cooking mode turns steps with the arrows; the step stays while the sheet is open
        e.stopPropagation();
    }
}

function onDown(e: PointerEvent) {
    const t = e.target as Node;
    if (box.value?.contains(t)) return;
    // the word itself toggles (TermMark)
    if (cur.value?.anchor.contains(t)) return;
    if (!sheet.value) closeTerm();
}

function listen(on: boolean) {
    const m = on ? 'addEventListener' : 'removeEventListener';
    document[m]('keydown', onKey as EventListener, true);
    document[m]('pointerdown', onDown as EventListener, true);
    window[m]('resize', place);
    window[m]('scroll', place, true);
}

watch(
    cur,
    async (o, before) => {
        if (!!o !== !!before) listen(!!o);
        pos.value = null;
        if (!o) return;
        await nextTick();
        place();
        // a sheet takes the focus (it's modal); the popover leaves it on the word
        if (o.mode !== 'popover') box.value?.querySelector<HTMLElement>('button')?.focus();
    },
    { flush: 'post' }
);
// another page: the word is gone
const route = useRoute();
watch(
    () => route.fullPath,
    () => closeTerm()
);
onBeforeUnmount(() => {
    if (cur.value) listen(false);
});
</script>

<template>
    <Teleport to="body">
        <template v-if="cur">
            <div
                v-if="sheet"
                class="explain-scrim"
                :class="{ 'explain-scrim--cook': cur.mode === 'cook' }"
                aria-hidden="true"
                @click="closeTerm(true)"
            />
            <div
                ref="box"
                class="explain-box"
                :class="sheet ? `explain-box--sheet explain-box--${cur.mode}` : 'explain-box--pop'"
                :style="!sheet && pos ? { left: `${pos.left}px`, top: `${pos.top}px` } : undefined"
                role="dialog"
                :aria-modal="sheet || undefined"
                :aria-labelledby="`${uid}-h`"
            >
                <TermExplainCard
                    :name="cur.term.name"
                    :explanation="cur.term.explanation"
                    :mode="cur.mode"
                    :heading-id="`${uid}-h`"
                    @close="closeTerm(true)"
                />
            </div>
        </template>
    </Teleport>
</template>

<style scoped>
.explain-box {
    z-index: calc(var(--z-modal) + 10);
    background: var(--bg-alt);
    border: 1px solid var(--rule-strong);
}
.explain-box--pop {
    position: fixed;
    /* off-screen until measured, so it doesn't flash at the top left */
    left: -9999px;
    top: 0;
    width: 340px;
    max-width: calc(100vw - 16px);
    border-radius: var(--r-xl);
    padding: 14px 16px;
    box-shadow: 0 12px 32px rgba(10, 10, 15, 0.18);
}
.explain-box--sheet {
    position: fixed;
    left: 0;
    right: 0;
    bottom: 0;
    border-radius: 16px 16px 0 0;
    border-bottom: none;
    padding: 8px 20px calc(20px + env(safe-area-inset-bottom));
    box-shadow: 0 -12px 40px rgba(10, 10, 15, 0.25);
}
.explain-box--cook {
    padding-bottom: calc(16px + env(safe-area-inset-bottom));
}
.explain-scrim {
    position: fixed;
    inset: 0;
    z-index: calc(var(--z-modal) + 9);
    background: rgba(10, 10, 15, 0.45);
}
.explain-scrim--cook {
    background: rgba(10, 10, 15, 0.5);
}
@media (min-width: 901px) {
    /* cooking mode on desktop: the sheet keeps the content width */
    .explain-box--cook {
        left: 50%;
        right: auto;
        width: min(640px, 100vw);
        transform: translateX(-50%);
    }
}
</style>
