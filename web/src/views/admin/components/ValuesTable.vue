<script lang="ts" setup>
import { ref } from 'vue';

import DeleteValueModal from './DeleteValueModal.vue';

import UIcon from '@/components/UIcon.vue';
import { confirmAction } from '@/composables/useConfirm';
import { showToast } from '@/composables/useToast';
import {
    createTaxonomyValue,
    deleteTaxonomyValue,
    updateTaxonomyValue,
} from '@/data/api/taxonomy/fetch';
import type { ITaxonomyAxis, ITaxonomyValue } from '@/data/dtos/taxonomy/DictionariesDto';

// Wartości jednej osi: edycja w wierszu (nazwa, kolor), dodawanie, usuwanie z przeniesieniem.
const props = defineProps<{ axis: ITaxonomyAxis }>();
const emit = defineEmits<{ changed: [] }>();

const SWATCHES = [
    'coral',
    'tomato',
    'terracotta',
    'lime',
    'plum',
    'pink',
    'brown',
    'mustard',
    'teal',
    'blue',
    'forest',
    'sage',
    'wheat',
    'olive',
];

const editId = ref<number | null>(null);
const buf = ref<{ label: string; color?: string }>({ label: '' });
const newLabel = ref('');
const deleting = ref<ITaxonomyValue | null>(null);
const busy = ref(false);

const recipesNoun = (n: number) =>
    n === 1
        ? 'przepis'
        : n % 10 >= 2 && n % 10 <= 4 && (n % 100 < 10 || n % 100 >= 20)
          ? 'przepisy'
          : 'przepisów';

function startEdit(v: ITaxonomyValue) {
    editId.value = v.id;
    buf.value = { label: v.label, color: v.color };
}

function apiError(e: unknown): string {
    return (e as { response?: { status?: number } }).response?.status === 409
        ? 'Taka wartość już istnieje'
        : 'Nie udało się zapisać';
}

async function saveEdit(v: ITaxonomyValue) {
    if (!buf.value.label.trim() || busy.value) return;
    busy.value = true;
    try {
        await updateTaxonomyValue(v.id, { label: buf.value.label, color: buf.value.color ?? null });
        editId.value = null;
        showToast('Zapisano');
        emit('changed');
    } catch (e) {
        showToast(apiError(e));
    } finally {
        busy.value = false;
    }
}

async function add() {
    const label = newLabel.value.trim();
    if (!label || busy.value) return;
    busy.value = true;
    const used = props.axis.values.map((v) => v.color);
    const color = props.axis.hasColor
        ? (SWATCHES.find((c) => !used.includes(c)) ?? SWATCHES[0])
        : null;
    try {
        await createTaxonomyValue({ axisId: props.axis.id, label, color });
        newLabel.value = '';
        showToast(`Dodano „${label}”`);
        emit('changed');
    } catch (e) {
        showToast(apiError(e));
    } finally {
        busy.value = false;
    }
}

async function remove(v: ITaxonomyValue) {
    if (v.recipeCount) {
        deleting.value = v;
        return;
    }
    const ok = await confirmAction({
        title: `Usunąć „${v.label}”?`,
        body: 'Ta wartość nie jest używana w żadnym przepisie.',
        confirmLabel: 'Usuń',
        danger: true,
    });
    if (ok) await doDelete(v, null);
}

async function doDelete(v: ITaxonomyValue, moveTo: number | null) {
    deleting.value = null;
    try {
        const res = await deleteTaxonomyValue(v.id, moveTo ?? undefined);
        showToast(
            moveTo && res.moved
                ? `Usunięto. ${res.moved} ${recipesNoun(res.moved)} przeniesiono.`
                : 'Usunięto'
        );
        emit('changed');
    } catch {
        showToast('Nie udało się usunąć');
    }
}
</script>

<template>
    <table class="data-table">
        <thead>
            <tr>
                <th>Nazwa</th>
                <th class="num">Przepisy</th>
                <th class="col-actions"><span class="sr-only">Akcje</span></th>
            </tr>
        </thead>
        <tbody>
            <tr
                v-for="v in axis.values"
                :key="v.id"
                :class="{ 'is-editing': editId === v.id }"
            >
                <td>
                    <div
                        v-if="editId === v.id"
                        class="edit"
                    >
                        <input
                            v-model="buf.label"
                            class="input edit__input"
                            aria-label="Nazwa"
                            autofocus
                            @keydown.enter="saveEdit(v)"
                            @keydown.esc="editId = null"
                        />
                        <div
                            v-if="axis.hasColor"
                            class="swatches"
                            role="group"
                            aria-label="Kolor"
                        >
                            <button
                                v-for="c in SWATCHES"
                                :key="c"
                                type="button"
                                class="swatch"
                                :class="[`cat-${c}`, { 'is-active': buf.color === c }]"
                                :aria-label="`Kolor ${c}`"
                                :aria-pressed="buf.color === c"
                                @click="buf.color = c"
                            />
                        </div>
                    </div>
                    <span
                        v-else
                        class="value"
                    >
                        <span
                            v-if="axis.hasColor"
                            class="value__dot"
                            :class="v.color ? `cat-${v.color}` : null"
                            aria-hidden="true"
                        />
                        {{ v.label }}
                    </span>
                </td>
                <td class="num">{{ v.recipeCount ?? 0 }}</td>
                <td>
                    <div class="actions">
                        <template v-if="editId === v.id">
                            <button
                                type="button"
                                class="tbl-btn"
                                aria-label="Zapisz"
                                title="Zapisz"
                                @click="saveEdit(v)"
                            >
                                <UIcon name="check" />
                            </button>
                            <button
                                type="button"
                                class="tbl-btn"
                                aria-label="Anuluj"
                                title="Anuluj"
                                @click="editId = null"
                            >
                                <UIcon name="x" />
                            </button>
                        </template>
                        <template v-else>
                            <button
                                type="button"
                                class="tbl-btn"
                                :aria-label="`Edytuj ${v.label}`"
                                title="Edytuj"
                                @click="startEdit(v)"
                            >
                                <UIcon name="edit" />
                            </button>
                            <button
                                type="button"
                                class="tbl-btn tbl-btn--danger"
                                :aria-label="`Usuń ${v.label}`"
                                title="Usuń"
                                @click="remove(v)"
                            >
                                <UIcon name="trash" />
                            </button>
                        </template>
                    </div>
                </td>
            </tr>
        </tbody>
    </table>
    <div class="add-row">
        <input
            v-model="newLabel"
            class="input add-row__input"
            :aria-label="`Nowa wartość: ${axis.label}`"
            :placeholder="`Nowa wartość: ${axis.label.toLowerCase()}`"
            @keydown.enter="add"
        />
        <button
            type="button"
            class="btn btn--secondary"
            :disabled="!newLabel.trim() || busy"
            @click="add"
        >
            + Dodaj
        </button>
    </div>

    <DeleteValueModal
        v-if="deleting"
        :label="deleting.label"
        :count="deleting.recipeCount ?? 0"
        :count-noun="recipesNoun(deleting.recipeCount ?? 0)"
        :others="axis.values.filter((v) => v.id !== deleting?.id)"
        :required="axis.required"
        @confirm="doDelete(deleting!, $event)"
        @close="deleting = null"
    />
</template>

<style scoped>
.col-actions {
    width: 84px;
}
tr.is-editing td {
    background: var(--accent-soft);
}
.value {
    display: inline-flex;
    align-items: center;
    gap: 8px;
    font-weight: 600;
}
.value__dot {
    width: 10px;
    height: 10px;
    border-radius: var(--r-pill);
    background: var(--cat, var(--rule-strong));
}
.edit {
    display: flex;
    flex-direction: column;
    gap: 8px;
}
.edit__input {
    padding: 6px 10px;
}
.swatches {
    display: flex;
    gap: 6px;
    flex-wrap: wrap;
}
.swatch {
    width: 22px;
    height: 22px;
    padding: 0;
    border-radius: var(--r-pill);
    background: var(--cat);
    border: 2px solid transparent;
    cursor: pointer;
}
.swatch.is-active {
    border-color: var(--ink);
}
.add-row {
    display: flex;
    gap: 8px;
    margin-top: 10px;
}
.add-row__input {
    flex: 1;
    background: var(--bg-alt);
}
</style>
