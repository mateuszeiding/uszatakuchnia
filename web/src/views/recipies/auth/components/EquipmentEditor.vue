<script lang="ts" setup>
import { computed, nextTick } from 'vue';

import { blankEquip, type IEquipItem } from '../recipeForm';

import EquipSuggest from './EquipSuggest.vue';

import UFieldError from '@/components/form/UFieldError.vue';
import { useListDnd } from '@/composables/useListDnd';
import { showToast } from '@/composables/useToast';
import { createEquipment } from '@/data/api/taxonomy/fetch';
import type { IEquipment } from '@/data/dtos/taxonomy/DictionariesDto';

// ✕ deletes without confirmation (nothing is saved until "Zapisz"); an empty row is skipped.
const items = defineModel<IEquipItem[]>({ required: true });
const props = defineProps<{
    equipment: IEquipment[];
    /** Row errors (id → message); empty until a save was attempted. */
    rowErrors: Record<number, string>;
}>();
const emit = defineEmits<{ created: [e: IEquipment] }>();

// useListDnd works on groups; equipment is a single list.
const dnd = useListDnd(computed(() => [{ id: 0, items: items.value }]));

function onName(it: IEquipItem, v: string) {
    it.name = v;
    // editing the text unlinks the chosen equipment until one is picked from the list
    const match = props.equipment.find((e) => e.id === it.equipmentId);
    if (!match || match.name !== v) it.equipmentId = 0;
}

function pick(it: IEquipItem, e: IEquipment) {
    it.equipmentId = e.id;
    it.name = e.name;
}

async function addToDictionary(it: IEquipItem, name: string) {
    try {
        const e = await createEquipment(name);
        emit('created', { ...e, usedIn: 0 });
        pick(it, e);
        showToast(`„${e.name}” dodany do słownika`);
    } catch (err) {
        const status = (err as { response?: { status?: number } }).response?.status;
        const existing = props.equipment.find((e) => e.name.toLowerCase() === name.toLowerCase());
        if (status === 409 && existing) pick(it, existing);
        else showToast('Nie udało się dodać do słownika');
    }
}

async function addRow() {
    const row = blankEquip();
    items.value.push(row);
    await nextTick();
    document.getElementById(`eq-${row.id}`)?.focus();
}

function removeRow(it: IEquipItem) {
    items.value.splice(items.value.indexOf(it), 1);
}
</script>

<template>
    <div class="eq-editor">
        <div
            v-if="items.length"
            class="eq-grid eq-editor__head"
            aria-hidden="true"
        >
            <span />
            <span>Sprzęt</span>
            <span>Doprecyzowanie</span>
            <span />
            <span />
        </div>
        <div
            v-for="it in items"
            :key="it.id"
            v-bind="dnd.row(0, it.id)"
            class="eq-row"
            :class="{ 'is-dragged': dnd.isDragged(it.id) }"
        >
            <span
                v-if="dnd.lineAt(0, it.id)"
                class="drop-line"
                :class="`drop-line--${dnd.lineAt(0, it.id)}`"
                aria-hidden="true"
            />
            <div class="eq-grid">
                <button
                    type="button"
                    class="drag-handle"
                    aria-label="Zmień kolejność: przeciągnij albo użyj strzałek ↑ ↓"
                    title="Przeciągnij"
                    v-bind="dnd.handle(0, it.id)"
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
                <EquipSuggest
                    :id="`eq-${it.id}`"
                    :name="it.name"
                    :equipment="equipment"
                    :invalid="!!rowErrors[it.id]"
                    @update:name="onName(it, $event)"
                    @select="pick(it, $event)"
                    @add-new="addToDictionary(it, $event)"
                />
                <input
                    v-model="it.spec"
                    class="input eq-row__spec"
                    :aria-label="`Doprecyzowanie: ${it.name || 'sprzęt'}`"
                    placeholder="np. 3–4 l, grube dno"
                    maxlength="160"
                />
                <button
                    type="button"
                    class="eq-row__opt"
                    :class="{ 'is-active': it.optional }"
                    :aria-pressed="it.optional"
                    @click="it.optional = !it.optional"
                >
                    {{ it.optional ? '✓ ' : '' }}Opcjonalny
                </button>
                <button
                    type="button"
                    class="row-icon-btn"
                    :aria-label="`Usuń: ${it.name || 'pusty wiersz'}`"
                    title="Usuń sprzęt"
                    @click="removeRow(it)"
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
                v-if="rowErrors[it.id]"
                class="eq-row__err"
            >
                {{ rowErrors[it.id] }}
            </UFieldError>
        </div>
        <div
            v-bind="dnd.groupEnd(0)"
            class="eq-editor__end"
        >
            <span
                v-if="dnd.endLine(0)"
                class="drop-line drop-line--after"
                aria-hidden="true"
            />
        </div>
        <button
            type="button"
            class="add-link eq-editor__add"
            @click="addRow"
        >
            + sprzęt
        </button>
    </div>
</template>

<style scoped>
.eq-editor {
    display: flex;
    flex-direction: column;
}
.eq-grid {
    display: grid;
    grid-template-columns: 20px minmax(0, 1fr) minmax(0, 1.2fr) 118px 30px;
    gap: 8px;
    align-items: center;
}
.eq-editor__head {
    padding-bottom: 6px;
    font-family: var(--font-mono);
    font-size: 11px;
    font-weight: 600;
    letter-spacing: 0.5px;
    text-transform: uppercase;
    color: var(--ink-muted);
}
.eq-row {
    position: relative;
    padding: 6px 0;
}
.eq-row.is-dragged {
    opacity: 0.4;
}
.eq-row__spec {
    padding: 8px 10px;
    font-size: 13px;
}
.eq-row__opt {
    height: 100%;
    padding: 8px 10px;
    border-radius: var(--r-md);
    border: 1px solid var(--rule-strong);
    background: var(--bg);
    color: var(--ink-dim);
    font-size: 13px;
    text-align: left;
    white-space: nowrap;
    cursor: pointer;
}
.eq-row__opt.is-active {
    color: var(--accent);
    border-color: var(--accent);
    background: var(--accent-soft);
}
.eq-row__err {
    margin: 6px 0 0 28px;
}
.eq-editor__end {
    position: relative;
    min-height: 2px;
}
.eq-editor__add {
    align-self: flex-start;
    margin-top: 8px;
}
</style>
