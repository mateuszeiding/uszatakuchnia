<script lang="ts" setup>
import { nextTick, onMounted, ref } from 'vue';

import BookIcon from '@/components/BookIcon.vue';
import type { ITerm } from '@/shared/utils/stepRefs';

// Clicking a marked word in the field: change the term, edit it in the glossary (new tab), or remove the mark.
const props = defineProps<{ term: ITerm | null; word: string }>();
const emit = defineEmits<{ change: []; remove: []; close: [] }>();

const root = ref<HTMLElement>();
const editHref = props.term ? `/admin/glossary?edit=${props.term.id}` : '';
onMounted(async () => {
    await nextTick();
    root.value?.querySelector<HTMLElement>('button, a')?.focus();
});
</script>

<template>
    <div
        ref="root"
        class="mark-menu"
        role="dialog"
        :aria-label="`Pojęcie: ${term?.name ?? word}`"
        @keydown.esc.stop.prevent="emit('close')"
    >
        <div class="mark-menu__title">
            <span class="mark-menu__icon"><BookIcon :size="13" /></span>
            {{ term?.name ?? word }}
        </div>
        <p
            v-if="term"
            class="mark-menu__text"
        >
            {{ term.explanation }}
        </p>
        <p
            v-else
            class="mark-menu__text mark-menu__text--danger"
        >
            Tego pojęcia nie ma już w słowniczku. Wybierz inne albo usuń oznaczenie — do tego czasu
            nie da się opublikować.
        </p>
        <div class="mark-menu__actions">
            <button
                type="button"
                class="mark-menu__btn"
                @click="emit('change')"
            >
                {{ term ? 'Zmień pojęcie' : 'Wybierz pojęcie' }}
            </button>
            <a
                v-if="term"
                class="mark-menu__btn"
                :href="editHref"
                target="_blank"
                rel="noopener"
            >
                Edytuj w słowniczku ↗
            </a>
            <button
                type="button"
                class="mark-menu__btn mark-menu__btn--danger"
                @click="emit('remove')"
            >
                Usuń oznaczenie
            </button>
        </div>
    </div>
</template>

<style scoped>
.mark-menu {
    width: 340px;
    max-width: calc(100vw - 32px);
    background: var(--bg-alt);
    border: 1px solid var(--rule-strong);
    border-radius: var(--r-xl);
    box-shadow: var(--shadow-menu);
    padding: 12px;
    font-size: 13px;
    line-height: 1.45;
}
.mark-menu__title {
    display: flex;
    align-items: center;
    gap: 6px;
    font-weight: 600;
    margin-bottom: 4px;
}
.mark-menu__icon {
    display: inline-flex;
    color: var(--accent);
}
.mark-menu__text {
    margin: 0;
    color: var(--ink-dim);
    display: -webkit-box;
    -webkit-line-clamp: 2;
    line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
}
.mark-menu__text--danger {
    color: var(--danger);
    display: block;
}
.mark-menu__actions {
    display: flex;
    flex-wrap: wrap;
    gap: 8px 12px;
    margin-top: 10px;
}
.mark-menu__btn {
    background: none;
    border: none;
    padding: 0;
    font: inherit;
    font-weight: 600;
    color: var(--accent);
    text-decoration: none;
    cursor: pointer;
}
.mark-menu__btn--danger {
    color: var(--danger);
    margin-left: auto;
}
</style>
