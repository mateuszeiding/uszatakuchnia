<script lang="ts" setup>
import { computed } from 'vue';

import BusyLabel from '@/components/BusyLabel.vue';
import UIcon from '@/components/UIcon.vue';
import type { RecipeStatus } from '@/data/dtos/recipe/RecipeDto';

export type BusyAction = 'draft' | 'preview' | 'publish' | 'save' | 'unpublish' | null;

const props = defineProps<{
    status: RecipeStatus;
    isNew: boolean;
    dirty: boolean;
    lastSaved: Date | null;
    busy: BusyAction;
    /** Number of errors blocking publishing, shown after a save attempt. */
    errorCount: number;
    /** Id of the saved recipe, for the "Zobacz na stronie" link when published. */
    recipeId: number | null;
}>();
const emit = defineEmits<{
    cancel: [];
    'save-draft': [];
    preview: [];
    publish: [];
    save: [];
    unpublish: [];
    'show-errors': [];
}>();

const pub = computed(() => props.status === 'published');

const plural = (n: number) =>
    n === 1
        ? 'pole'
        : n % 10 >= 2 && n % 10 <= 4 && (n % 100 < 10 || n % 100 >= 20)
          ? 'pola'
          : 'pól';

const state = computed<{ kind: 'error' | 'dirty' | 'saved' | 'idle'; text: string }>(() => {
    if (props.errorCount)
        return {
            kind: 'error',
            text: `Uzupełnij ${props.errorCount} ${plural(props.errorCount)}, żeby ${pub.value ? 'zapisać' : 'opublikować'}.`,
        };
    if (props.dirty)
        return {
            kind: 'dirty',
            text: props.isNew ? 'Nowy przepis, jeszcze niezapisany' : 'Niezapisane zmiany',
        };
    if (props.lastSaved)
        return {
            kind: 'saved',
            text: `Zapisano o ${props.lastSaved.toLocaleTimeString('pl-PL', { hour: '2-digit', minute: '2-digit' })}`,
        };
    return { kind: 'idle', text: props.isNew ? 'Nowy przepis' : 'Brak zmian' };
});
</script>

<template>
    <div class="action-bar">
        <div class="action-bar__inner">
            <div
                class="action-bar__state"
                :class="`is-${state.kind}`"
            >
                <span
                    class="action-bar__dot"
                    aria-hidden="true"
                />
                <span role="status">
                    {{ state.text }}
                    <button
                        v-if="state.kind === 'error'"
                        type="button"
                        class="action-bar__show"
                        @click="emit('show-errors')"
                    >
                        Pokaż
                    </button>
                </span>
            </div>
            <div class="action-bar__btns">
                <button
                    type="button"
                    class="btn btn--ghost"
                    :disabled="!!busy"
                    @click="emit('cancel')"
                >
                    Anuluj
                </button>
                <template v-if="pub">
                    <!-- published: no preview of unsaved changes (saving publishes right away), just a link -->
                    <RouterLink
                        v-if="recipeId"
                        :to="{ name: 'recipe-details', params: { id: recipeId } }"
                        target="_blank"
                        rel="noopener"
                        class="action-bar__view"
                    >
                        Zobacz na stronie
                        <UIcon
                            name="external"
                            :size="12"
                        />
                        <span class="sr-only">(nowa karta)</span>
                    </RouterLink>
                    <button
                        type="button"
                        class="btn btn--secondary"
                        :disabled="!!busy"
                        @click="emit('unpublish')"
                    >
                        <BusyLabel
                            :busy="busy === 'unpublish'"
                            label="Cofnij publikację"
                        />
                    </button>
                    <button
                        type="button"
                        class="btn btn--accent"
                        :disabled="!!busy || !dirty"
                        @click="emit('save')"
                    >
                        <BusyLabel
                            :busy="busy === 'save'"
                            label="Zapisz zmiany"
                            busy-label="Zapisywanie…"
                            spinner-color="#fff"
                        />
                    </button>
                </template>
                <template v-else>
                    <button
                        type="button"
                        class="btn btn--secondary"
                        :disabled="!!busy"
                        title="Zapisuje szkic i otwiera stronę w nowej karcie"
                        @click="emit('preview')"
                    >
                        <!-- while saving, the icon gives way to the spinner but keeps its space (the buttons don't shift) -->
                        <UIcon
                            name="external"
                            :class="{ 'action-bar__icon-off': busy === 'preview' }"
                        />
                        <BusyLabel
                            :busy="busy === 'preview'"
                            label="Podgląd"
                            busy-label="Zapisywanie…"
                        />
                    </button>
                    <button
                        type="button"
                        class="btn btn--secondary"
                        :disabled="!!busy || (!dirty && !isNew)"
                        @click="emit('save-draft')"
                    >
                        <BusyLabel
                            :busy="busy === 'draft'"
                            label="Zapisz szkic"
                            busy-label="Zapisywanie…"
                        />
                    </button>
                    <button
                        type="button"
                        class="btn btn--accent"
                        :disabled="!!busy"
                        @click="emit('publish')"
                    >
                        <BusyLabel
                            :busy="busy === 'publish'"
                            label="Opublikuj"
                            busy-label="Publikowanie…"
                            spinner-color="#fff"
                        />
                    </button>
                </template>
            </div>
        </div>
    </div>
</template>

<style scoped>
.action-bar {
    position: fixed;
    bottom: 0;
    left: 0;
    right: 0;
    z-index: var(--z-sticky);
    background: color-mix(in srgb, var(--bg) 94%, transparent);
    backdrop-filter: blur(12px);
    border-top: 1px solid var(--rule);
}
.action-bar__inner {
    max-width: 1280px;
    margin: 0 auto;
    padding: 14px 56px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 16px;
}
.action-bar__state {
    display: flex;
    align-items: center;
    gap: 10px;
    font-size: 13px;
    color: var(--ink-dim);
}
.action-bar__dot {
    width: 8px;
    height: 8px;
    border-radius: var(--r-pill);
    flex-shrink: 0;
    background: var(--rule-strong);
}
.is-error {
    color: var(--danger);
}
.is-error .action-bar__dot {
    background: var(--danger);
}
.is-dirty .action-bar__dot {
    background: var(--c-mustard);
}
.is-saved .action-bar__dot {
    background: var(--success);
}
.action-bar__show {
    background: none;
    border: none;
    padding: 0;
    color: var(--danger);
    font-weight: 600;
    text-decoration: underline;
    cursor: pointer;
}
.action-bar__btns {
    display: flex;
    align-items: center;
    gap: 8px;
}
.action-bar__icon-off {
    visibility: hidden;
}
.action-bar__view {
    display: inline-flex;
    align-items: center;
    gap: 5px;
    padding: 0 6px;
    font-size: 13px;
    color: var(--ink-muted);
}
.action-bar__view:hover {
    color: var(--ink);
}
</style>
