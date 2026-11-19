<script lang="ts" setup>
import { computed, nextTick } from 'vue';

import { blankIngredient, type IChallengeIngRow } from '../challengeForm';

import UFieldError from '@/components/form/UFieldError.vue';
import { useListDnd } from '@/composables/useListDnd';

// Składniki wyzwania: wiersz = uchwyt (przeciąganie, ↑ ↓), tekst, „Obowiązkowy” (aria-pressed), ✕.
// Zwykły tekst bez bazy składników; puste wiersze są pomijane przy zapisie, usunięcie wiersza bez pytania.
defineProps<{ errors: Record<string, string> }>();
const rows = defineModel<IChallengeIngRow[]>({ required: true });

const dnd = useListDnd(computed(() => [{ id: 0, items: rows.value }]));

async function addRow() {
    const row = blankIngredient();
    rows.value.push(row);
    await nextTick();
    document.getElementById(`ching-${row.id}`)?.focus();
}

const removeRow = (r: IChallengeIngRow) => rows.value.splice(rows.value.indexOf(r), 1);
</script>

<template>
    <div class="ie">
        <div
            v-if="rows.length"
            class="ie__grid ie__head"
            aria-hidden="true"
        >
            <span />
            <span>Składnik</span>
            <span />
            <span />
        </div>
        <div
            v-for="r in rows"
            :key="r.id"
            v-bind="dnd.row(0, r.id)"
            class="ie__row"
            :class="{ 'is-dragged': dnd.isDragged(r.id) }"
        >
            <span
                v-if="dnd.lineAt(0, r.id)"
                class="drop-line"
                :class="`drop-line--${dnd.lineAt(0, r.id)}`"
                aria-hidden="true"
            />
            <div class="ie__grid">
                <button
                    type="button"
                    class="drag-handle"
                    aria-label="Zmień kolejność: przeciągnij albo użyj strzałek ↑ ↓"
                    title="Przeciągnij"
                    v-bind="dnd.handle(0, r.id)"
                >
                    <svg
                        width="10"
                        height="14"
                        viewBox="0 0 10 14"
                        fill="currentColor"
                        aria-hidden="true"
                    >
                        <circle
                            v-for="(c, ci) in [
                                [3, 3],
                                [7, 3],
                                [3, 7],
                                [7, 7],
                                [3, 11],
                                [7, 11],
                            ]"
                            :key="ci"
                            :cx="c[0]"
                            :cy="c[1]"
                            r="1.2"
                        />
                    </svg>
                </button>
                <input
                    :id="`ching-${r.id}`"
                    v-model="r.name"
                    class="input ie__name"
                    :class="{ 'is-invalid': errors[`ing-${r.id}`] }"
                    aria-label="Składnik"
                    placeholder="np. ryba białego mięsa"
                    maxlength="120"
                />
                <button
                    type="button"
                    class="ie__req"
                    :class="{ 'is-active': r.required }"
                    :aria-pressed="r.required"
                    @click="r.required = !r.required"
                >
                    {{ r.required ? '✓ ' : '' }}Obowiązkowy
                </button>
                <button
                    type="button"
                    class="row-icon-btn"
                    :aria-label="`Usuń: ${r.name || 'pusty wiersz'}`"
                    title="Usuń składnik"
                    @click="removeRow(r)"
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
            <UFieldError
                v-if="errors[`ing-${r.id}`]"
                class="ie__err"
            >
                {{ errors[`ing-${r.id}`] }}
            </UFieldError>
        </div>
        <div
            v-bind="dnd.groupEnd(0)"
            class="ie__end"
        >
            <span
                v-if="dnd.endLine(0)"
                class="drop-line drop-line--after"
                aria-hidden="true"
            />
        </div>
        <button
            type="button"
            class="add-link ie__add"
            @click="addRow"
        >
            + składnik
        </button>
    </div>
</template>

<style scoped>
.ie {
    display: flex;
    flex-direction: column;
}
.ie__grid {
    display: grid;
    grid-template-columns: 20px minmax(0, 1fr) 128px 30px;
    gap: 8px;
    align-items: center;
}
.ie__head {
    padding-bottom: 6px;
    font-family: var(--font-mono);
    font-size: 11px;
    font-weight: 600;
    letter-spacing: 0.5px;
    text-transform: uppercase;
    color: var(--ink-muted);
}
.ie__row {
    position: relative;
    padding: 5px 0;
}
.ie__row.is-dragged {
    opacity: 0.4;
}
.ie__name {
    padding: 9px 10px;
}
.ie__req {
    height: 100%;
    padding: 8px 10px;
    border-radius: var(--r-md);
    border: 1px solid var(--rule-strong);
    background: var(--bg);
    color: var(--ink-dim);
    font-family: var(--font-sans);
    font-size: 13px;
    text-align: left;
    white-space: nowrap;
    cursor: pointer;
}
.ie__req.is-active {
    color: var(--accent);
    border-color: var(--accent-fill);
    background: var(--accent-soft);
}
.ie__err {
    margin: 6px 0 0 28px;
}
.ie__end {
    position: relative;
    min-height: 2px;
}
.ie__add {
    align-self: flex-start;
    margin-top: 8px;
}
</style>
