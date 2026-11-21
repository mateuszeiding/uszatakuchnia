<script lang="ts" setup>
import { nextTick, onMounted, ref } from 'vue';

import { chipLabel, type IRefRow } from '../stepRefsForm';

import SwapIcon from '@/components/SwapIcon.vue';
import { partLabel, PARTS, type RefPart } from '@/shared/utils/stepRefs';

// Clicking a ref in the step field: change the row, its part, or remove it.
defineProps<{ row: IRefRow | null; part: RefPart | null }>();
const emit = defineEmits<{
    change: [];
    part: [p: RefPart | null];
    remove: [];
    close: [];
}>();

const parts = ref(false);
const root = ref<HTMLElement>();
onMounted(async () => {
    await nextTick();
    root.value?.querySelector<HTMLElement>('button')?.focus();
});
</script>

<template>
    <div
        ref="root"
        class="ref-menu"
        role="dialog"
        :aria-label="row ? chipLabel(row, part) : 'Brak składnika'"
        @keydown.esc.stop.prevent="emit('close')"
    >
        <div class="ref-menu__title">{{ row ? chipLabel(row, part) : 'Brak składnika' }}</div>
        <p
            v-if="!row"
            class="ref-menu__note ref-menu__note--danger"
        >
            Ten wiersz usunięto z sekcji Składniki. Wybierz inny składnik albo usuń odwołanie.
        </p>
        <p
            v-else-if="row.hasSubs"
            class="ref-menu__note"
        >
            <span class="ref-menu__swap"><SwapIcon :size="13" /></span>
            Ten składnik ma zamienniki. Słowo przed odwołaniem się nie zmieni — napisz ogólnie, np.
            „mięso”, a nazwę zamiennika pokaże znacznik.
        </p>
        <div class="ref-menu__actions">
            <button
                type="button"
                class="ref-menu__btn"
                @click="emit('change')"
            >
                {{ row ? 'Zmień składnik' : 'Wybierz składnik' }}
            </button>
            <button
                v-if="row"
                type="button"
                class="ref-menu__btn"
                :aria-expanded="parts"
                @click="parts = !parts"
            >
                Część: {{ partLabel(part) }} ⌄
            </button>
            <button
                type="button"
                class="ref-menu__btn ref-menu__btn--danger"
                @click="emit('remove')"
            >
                Usuń
            </button>
        </div>
        <div
            v-if="parts"
            class="ref-menu__radios"
            role="radiogroup"
            aria-label="Część"
        >
            <button
                v-for="p in PARTS"
                :key="p.label"
                type="button"
                role="radio"
                class="ref-menu__radio"
                :aria-checked="p.part === part"
                @click="emit('part', p.part)"
            >
                {{ p.label }}
            </button>
        </div>
    </div>
</template>

<style scoped>
.ref-menu {
    position: absolute;
    z-index: 20;
    width: 360px;
    max-width: calc(100vw - 32px);
    background: var(--bg-alt);
    border: 1px solid var(--rule-strong);
    border-radius: var(--r-lg);
    box-shadow: var(--shadow-menu);
    padding: 12px;
    font-size: 13px;
    line-height: 1.5;
}
.ref-menu__title {
    font-weight: 600;
    margin-bottom: 4px;
}
.ref-menu__note {
    display: flex;
    gap: 6px;
    margin: 0;
    color: var(--ink-dim);
}
.ref-menu__note--danger {
    color: var(--danger);
}
.ref-menu__swap {
    display: inline-flex;
    color: var(--accent);
    padding-top: 2px;
}
.ref-menu__actions {
    display: flex;
    gap: 12px;
    margin-top: 10px;
}
.ref-menu__btn {
    background: none;
    border: none;
    padding: 0;
    font: inherit;
    font-weight: 600;
    color: var(--accent);
    cursor: pointer;
}
.ref-menu__btn--danger {
    color: var(--danger);
    margin-left: auto;
}
.ref-menu__radios {
    display: flex;
    gap: 4px;
    flex-wrap: wrap;
    margin-top: 10px;
}
.ref-menu__radio {
    padding: 5px 10px;
    border-radius: var(--r-pill);
    font-size: 13px;
    font-weight: 600;
    border: 1px solid var(--rule-strong);
    background: var(--bg-alt);
    color: var(--ink);
    cursor: pointer;
}
.ref-menu__radio[aria-checked='true'] {
    border-color: var(--accent);
    background: var(--accent-soft);
    color: var(--accent);
}
</style>
