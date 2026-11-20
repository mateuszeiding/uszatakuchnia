<script lang="ts" setup>
import { computed, ref } from 'vue';

import DeleteEquipmentModal from './DeleteEquipmentModal.vue';

import UIcon from '@/components/UIcon.vue';
import { confirmAction } from '@/composables/useConfirm';
import { showToast } from '@/composables/useToast';
import { createEquipment, deleteEquipment, updateEquipment } from '@/data/api/taxonomy/fetch';
import type { IEquipment } from '@/data/dtos/taxonomy/DictionariesDto';

// Equipment dictionary: add, rename inline, delete (in use → replace or remove from recipes).
const props = defineProps<{ items: IEquipment[] }>();
const emit = defineEmits<{ changed: [] }>();

const sorted = computed(() => [...props.items].sort((a, b) => a.name.localeCompare(b.name, 'pl')));
const newName = ref('');
const editId = ref<number | null>(null);
const buf = ref('');
const deleting = ref<IEquipment | null>(null);
const busy = ref(false);

const usedText = (n: number) => (n ? `${n} ${n === 1 ? 'przepisie' : 'przepisach'}` : '—');
const recipesNoun = (n: number) =>
    n === 1
        ? 'przepis'
        : n % 10 >= 2 && n % 10 <= 4 && (n % 100 < 10 || n % 100 >= 20)
          ? 'przepisy'
          : 'przepisów';

function apiError(e: unknown): string {
    return (e as { response?: { status?: number } }).response?.status === 409
        ? 'Taki sprzęt już jest w słowniku'
        : 'Nie udało się zapisać';
}

async function add() {
    const name = newName.value.trim();
    if (!name || busy.value) return;
    busy.value = true;
    try {
        await createEquipment(name);
        newName.value = '';
        showToast(`Dodano „${name}”`);
        emit('changed');
    } catch (e) {
        showToast(apiError(e));
    } finally {
        busy.value = false;
    }
}

function startEdit(x: IEquipment) {
    editId.value = x.id;
    buf.value = x.name;
}

async function saveEdit(x: IEquipment) {
    if (!buf.value.trim() || busy.value) return;
    busy.value = true;
    try {
        await updateEquipment(x.id, buf.value);
        editId.value = null;
        showToast('Zapisano');
        emit('changed');
    } catch (e) {
        showToast(apiError(e));
    } finally {
        busy.value = false;
    }
}

async function remove(x: IEquipment) {
    if (x.usedIn) {
        deleting.value = x;
        return;
    }
    const ok = await confirmAction({
        title: `Usunąć „${x.name}”?`,
        body: 'Ten sprzęt nie jest używany w żadnym przepisie.',
        confirmLabel: 'Usuń',
        danger: true,
    });
    if (ok) await doDelete(x, {});
}

async function doDelete(x: IEquipment, opt: { replaceWith?: number; removeFromRecipes?: boolean }) {
    deleting.value = null;
    try {
        const res = await deleteEquipment(x.id, opt);
        const target = props.items.find((i) => i.id === opt.replaceWith);
        showToast(
            res.recipes && target
                ? `Usunięto. ${res.recipes} ${recipesNoun(res.recipes)} ma teraz „${target.name}”.`
                : res.recipes
                  ? `Usunięto z ${res.recipes} ${res.recipes === 1 ? 'przepisu' : 'przepisów'}.`
                  : 'Usunięto'
        );
        emit('changed');
    } catch {
        showToast('Nie udało się usunąć');
    }
}
</script>

<template>
    <div class="eq-add">
        <input
            v-model="newName"
            class="input eq-add__input"
            aria-label="Nowy sprzęt"
            placeholder="Nowy sprzęt, np. termometr kuchenny"
            maxlength="120"
            @keydown.enter="add"
        />
        <button
            type="button"
            class="btn btn--primary"
            :disabled="!newName.trim() || busy"
            @click="add"
        >
            Dodaj
        </button>
    </div>
    <table class="data-table">
        <thead>
            <tr>
                <th>Nazwa</th>
                <th class="col-used">Użyty w</th>
                <th class="col-actions"><span class="sr-only">Akcje</span></th>
            </tr>
        </thead>
        <tbody>
            <tr
                v-for="x in sorted"
                :key="x.id"
                :class="{ 'is-editing': editId === x.id }"
            >
                <td>
                    <input
                        v-if="editId === x.id"
                        v-model="buf"
                        class="input eq-edit"
                        :aria-label="`Nazwa: ${x.name}`"
                        maxlength="120"
                        autofocus
                        @keydown.enter="saveEdit(x)"
                        @keydown.esc="editId = null"
                    />
                    <span
                        v-else
                        class="eq-name"
                    >
                        {{ x.name }}
                    </span>
                </td>
                <td
                    class="col-used eq-used"
                    :class="{ 'is-unused': !x.usedIn }"
                >
                    {{ usedText(x.usedIn) }}
                </td>
                <td>
                    <div class="actions">
                        <template v-if="editId === x.id">
                            <button
                                type="button"
                                class="tbl-btn"
                                aria-label="Zapisz"
                                title="Zapisz"
                                @click="saveEdit(x)"
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
                                :aria-label="`Zmień nazwę: ${x.name}`"
                                title="Zmień nazwę"
                                @click="startEdit(x)"
                            >
                                <UIcon name="edit" />
                            </button>
                            <button
                                type="button"
                                class="tbl-btn tbl-btn--danger"
                                :aria-label="`Usuń: ${x.name}`"
                                title="Usuń"
                                @click="remove(x)"
                            >
                                <UIcon name="trash" />
                            </button>
                        </template>
                    </div>
                </td>
            </tr>
            <tr v-if="!sorted.length">
                <td
                    colspan="3"
                    class="eq-empty"
                >
                    Słownik jest pusty. Dodaj pierwszy sprzęt powyżej albo z formularza przepisu.
                </td>
            </tr>
        </tbody>
    </table>

    <DeleteEquipmentModal
        v-if="deleting"
        :item="deleting"
        :others="sorted.filter((x) => x.id !== deleting?.id)"
        @confirm="doDelete(deleting!, $event)"
        @close="deleting = null"
    />
</template>

<style scoped>
.eq-add {
    display: flex;
    gap: 10px;
    margin-bottom: 14px;
}
.eq-add__input {
    flex: 1;
    background: var(--bg-alt);
}
.col-used {
    width: 130px;
}
.col-actions {
    width: 90px;
}
tr.is-editing td {
    background: var(--accent-soft);
}
.eq-name {
    line-height: 1.4;
}
.eq-edit {
    padding: 6px 10px;
}
.eq-used {
    font-family: var(--font-mono);
    font-size: 13px;
    color: var(--ink-dim);
}
.eq-used.is-unused {
    color: var(--ink-muted);
}
.eq-empty {
    color: var(--ink-muted);
    font-size: 13px;
}
</style>
