<script lang="ts" setup>
import { computed, nextTick } from 'vue';

import UFieldError from '@/components/form/UFieldError.vue';
import { useListDnd } from '@/composables/useListDnd';
import {
    blankCaution,
    CAUTION_COUNTER_FROM,
    type ICautionRow,
    MAX_CAUTION_LEN,
} from '@/shared/utils/cautions';

// Notes editor for the recipe and challenge forms: rows with a drag handle (↑ ↓), an auto-growing field and ✕.
// Deleting a row doesn't ask: nothing is saved until "Zapisz". Empty rows are skipped on save.
const props = defineProps<{
    /** "przepisie" / "wyzwaniu": completes the hint "czego nie wolno zawalić w całym …". */
    label: string;
    /** Row errors keyed caution-<row id> (as in challengeForm). */
    errors?: Record<string, string>;
}>();
const rows = defineModel<ICautionRow[]>({ required: true });

const dnd = useListDnd(computed(() => [{ id: 0, items: rows.value }]));
const counterId = (r: ICautionRow) => `caution-count-${r.id}`;
const showCounter = (r: ICautionRow) => r.text.length >= CAUTION_COUNTER_FROM;
const error = (r: ICautionRow) => props.errors?.[`caution-${r.id}`];

async function addRow() {
    const row = blankCaution();
    rows.value.push(row);
    await nextTick();
    document.getElementById(`caution-${row.id}`)?.focus();
}

const removeRow = (r: ICautionRow) => rows.value.splice(rows.value.indexOf(r), 1);
</script>

<template>
    <div class="ce">
        <p class="ce__hint">
            <svg
                class="ce__hint-icon"
                width="14"
                height="14"
                viewBox="0 0 16 16"
                fill="none"
                aria-hidden="true"
            >
                <path
                    d="M8 2.2L14.4 13.3H1.6L8 2.2z"
                    stroke="currentColor"
                    stroke-width="1.4"
                    stroke-linejoin="round"
                />
                <path
                    d="M8 6.4v3.3M8 11.6v.1"
                    stroke="currentColor"
                    stroke-width="1.6"
                    stroke-linecap="round"
                />
            </svg>
            <span>
                <strong>Uwaga</strong>
                — czego nie wolno zawalić w całym {{ label }}: zakupy, bezpieczeństwo, czas,
                planowanie. Czyta się ją przed startem. Jedno zdanie, góra dwa.
                <br />
                <strong>Wskazówka</strong>
                — jak zrobić lepiej konkretny krok („jeśli…”, „możesz też…”). Dodaj ją przy kroku,
                nie tutaj.
            </span>
        </p>

        <div
            v-for="(r, i) in rows"
            :key="r.id"
            v-bind="dnd.row(0, r.id)"
            class="ce__row"
            :class="{ 'is-dragged': dnd.isDragged(r.id) }"
        >
            <span
                v-if="dnd.lineAt(0, r.id)"
                class="drop-line"
                :class="`drop-line--${dnd.lineAt(0, r.id)}`"
                aria-hidden="true"
            />
            <button
                type="button"
                class="drag-handle ce__handle"
                :aria-label="`Uwaga ${i + 1}: zmień kolejność — przeciągnij albo użyj strzałek ↑ ↓`"
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
            <span class="ce__field">
                <textarea
                    :id="`caution-${r.id}`"
                    v-model="r.text"
                    v-autosize
                    class="textarea ce__input"
                    :class="{ 'is-invalid': error(r) }"
                    rows="2"
                    :maxlength="MAX_CAUTION_LEN"
                    :aria-label="`Uwaga ${i + 1}`"
                    :aria-describedby="showCounter(r) ? counterId(r) : undefined"
                    :aria-invalid="!!error(r) || undefined"
                    placeholder="np. Mięso tylko świeże od rzeźnika, nie paczkowane."
                />
                <span
                    v-if="showCounter(r)"
                    :id="counterId(r)"
                    class="ce__count"
                    :class="{ 'is-max': r.text.length >= MAX_CAUTION_LEN }"
                >
                    {{ r.text.length }} / {{ MAX_CAUTION_LEN }}
                </span>
                <UFieldError v-if="error(r)">{{ error(r) }}</UFieldError>
            </span>
            <button
                type="button"
                class="row-icon-btn ce__remove"
                :aria-label="`Usuń uwagę ${i + 1}`"
                :title="`Usuń uwagę ${i + 1}`"
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
        <div
            v-bind="dnd.groupEnd(0)"
            class="ce__end"
        >
            <span
                v-if="dnd.endLine(0)"
                class="drop-line drop-line--after"
                aria-hidden="true"
            />
        </div>
        <button
            type="button"
            class="add-link ce__add"
            @click="addRow"
        >
            + uwaga
        </button>
    </div>
</template>

<style scoped>
.ce {
    display: flex;
    flex-direction: column;
}
.ce__hint {
    display: grid;
    grid-template-columns: 16px minmax(0, 1fr);
    gap: 10px;
    margin: 0 0 14px;
    padding: 10px 12px;
    background: var(--bg);
    border: 1px solid var(--rule);
    border-radius: var(--r-lg);
    font-size: 13px;
    line-height: 1.5;
    color: var(--ink-dim);
}
.ce__hint-icon {
    margin-top: 2px;
    color: var(--prep);
}
.ce__hint strong {
    font-weight: 600;
    color: var(--ink);
}
.ce__row {
    position: relative;
    display: grid;
    grid-template-columns: 20px minmax(0, 1fr) 30px;
    gap: 8px;
    /* the field grows downward; the handle and ✕ stay by the first line */
    align-items: start;
    padding: 5px 0;
}
.ce__row.is-dragged {
    opacity: 0.4;
}
.ce__handle {
    margin-top: 4px;
}
.ce__field {
    display: flex;
    flex-direction: column;
    gap: 4px;
    min-width: 0;
}
/* at least 2 lines (rows=2), then grows with the content */
.ce__input {
    min-height: 0;
    padding: 9px 10px;
    line-height: 1.45;
}
.ce__count {
    align-self: flex-end;
    font-family: var(--font-mono);
    font-size: 12px;
    color: var(--ink-muted);
}
.ce__count.is-max {
    color: var(--danger);
}
.ce__remove {
    margin-top: 6px;
}
.ce__end {
    position: relative;
    min-height: 2px;
}
.ce__add {
    margin-top: 8px;
}
</style>
