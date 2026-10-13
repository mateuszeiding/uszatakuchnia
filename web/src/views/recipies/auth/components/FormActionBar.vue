<script lang="ts" setup>
import { computed } from 'vue';

import USpinner from '@/components/USpinner.vue';
import type { RecipeStatus } from '@/data/dtos/recipe/RecipeDto';

export type BusyAction = 'draft' | 'publish' | 'save' | 'unpublish' | null;

const props = defineProps<{
    status: RecipeStatus;
    isNew: boolean;
    dirty: boolean;
    lastSaved: Date | null;
    busy: BusyAction;
    /** Liczba błędów blokujących publikację — pokazywana po próbie zapisu. */
    errorCount: number;
}>();
const emit = defineEmits<{
    cancel: [];
    'save-draft': [];
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
                    <button
                        type="button"
                        class="btn btn--secondary"
                        :disabled="!!busy"
                        @click="emit('unpublish')"
                    >
                        <USpinner
                            v-if="busy === 'unpublish'"
                            :size="14"
                        />
                        Cofnij publikację
                    </button>
                    <button
                        type="button"
                        class="btn btn--accent"
                        :disabled="!!busy || !dirty"
                        @click="emit('save')"
                    >
                        <USpinner
                            v-if="busy === 'save'"
                            :size="14"
                            color="#fff"
                        />
                        {{ busy === 'save' ? 'Zapisywanie…' : 'Zapisz zmiany' }}
                    </button>
                </template>
                <template v-else>
                    <button
                        type="button"
                        class="btn btn--secondary"
                        :disabled="!!busy || (!dirty && !isNew)"
                        @click="emit('save-draft')"
                    >
                        <USpinner
                            v-if="busy === 'draft'"
                            :size="14"
                        />
                        {{ busy === 'draft' ? 'Zapisywanie…' : 'Zapisz szkic' }}
                    </button>
                    <button
                        type="button"
                        class="btn btn--accent"
                        :disabled="!!busy"
                        @click="emit('publish')"
                    >
                        <USpinner
                            v-if="busy === 'publish'"
                            :size="14"
                            color="#fff"
                        />
                        {{ busy === 'publish' ? 'Publikowanie…' : 'Opublikuj' }}
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
    gap: 8px;
}
</style>
