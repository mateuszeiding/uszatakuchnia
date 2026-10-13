<script lang="ts" setup>
import { ref } from 'vue';

import DeleteValueModal from './DeleteValueModal.vue';

import UIcon from '@/components/UIcon.vue';
import { confirmAction } from '@/composables/useConfirm';
import { showToast } from '@/composables/useToast';
import { createUnit, deleteUnit, updateUnit } from '@/data/api/taxonomy/fetch';
import type { IUnit } from '@/data/dtos/taxonomy/DictionariesDto';

// Jednostki: nazwa + formy odmiany „1, 2–4, 5+” (puste dla skrótów g, ml).
const props = defineProps<{ units: IUnit[] }>();
const emit = defineEmits<{ changed: [] }>();

const editId = ref<number | null>(null);
const buf = ref({ name: '', forms: '' });
const newName = ref('');
const deleting = ref<IUnit | null>(null);
const busy = ref(false);

const parseForms = (s: string) =>
    s
        .split(',')
        .map((x) => x.trim())
        .filter(Boolean);

const usesNoun = (n: number) => (n === 1 ? 'składnik w przepisach' : 'składników w przepisach');

function startEdit(u: IUnit) {
    editId.value = u.id;
    buf.value = { name: u.name, forms: u.forms?.join(', ') ?? '' };
}

async function saveEdit(u: IUnit) {
    const forms = parseForms(buf.value.forms);
    if (forms.length && forms.length !== 3) {
        showToast('Podaj trzy formy: 1, 2–4, 5+ (np. ząbek, ząbki, ząbków)');
        return;
    }
    if (!buf.value.name.trim() || busy.value) return;
    busy.value = true;
    try {
        await updateUnit(u.id, { name: buf.value.name, forms, noQty: u.noQty });
        editId.value = null;
        showToast('Zapisano');
        emit('changed');
    } catch (e) {
        showToast(
            (e as { response?: { status?: number } }).response?.status === 409
                ? 'Taka jednostka już istnieje'
                : 'Nie udało się zapisać'
        );
    } finally {
        busy.value = false;
    }
}

async function add() {
    const name = newName.value.trim();
    if (!name || busy.value) return;
    busy.value = true;
    try {
        await createUnit({ name, forms: [], noQty: false });
        newName.value = '';
        showToast(`Dodano „${name}”`);
        emit('changed');
    } catch {
        showToast('Nie udało się dodać (może już istnieje)');
    } finally {
        busy.value = false;
    }
}

async function remove(u: IUnit) {
    if (u.usedIn) {
        deleting.value = u;
        return;
    }
    const ok = await confirmAction({
        title: `Usunąć „${u.name}”?`,
        body: 'Ta jednostka nie jest używana w żadnym przepisie.',
        confirmLabel: 'Usuń',
        danger: true,
    });
    if (ok) await doDelete(u, null);
}

async function doDelete(u: IUnit, moveTo: number | null) {
    deleting.value = null;
    try {
        await deleteUnit(u.id, moveTo ?? undefined);
        showToast('Usunięto');
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
                <th>Odmiana</th>
                <th class="num">Użycia</th>
                <th class="col-actions"><span class="sr-only">Akcje</span></th>
            </tr>
        </thead>
        <tbody>
            <tr
                v-for="u in units"
                :key="u.id"
                :class="{ 'is-editing': editId === u.id }"
            >
                <td>
                    <input
                        v-if="editId === u.id"
                        v-model="buf.name"
                        class="input edit-input"
                        aria-label="Nazwa"
                        autofocus
                        @keydown.enter="saveEdit(u)"
                        @keydown.esc="editId = null"
                    />
                    <strong v-else>{{ u.name }}</strong>
                </td>
                <td>
                    <input
                        v-if="editId === u.id"
                        v-model="buf.forms"
                        class="input edit-input"
                        aria-label="Formy odmiany"
                        placeholder="ząbek, ząbki, ząbków"
                        @keydown.enter="saveEdit(u)"
                        @keydown.esc="editId = null"
                    />
                    <span
                        v-else
                        class="forms"
                        :class="{ 'is-empty': !u.forms }"
                    >
                        {{ u.noQty ? 'bez ilości' : u.forms ? u.forms.join(' · ') : 'bez odmiany' }}
                    </span>
                </td>
                <td class="num">{{ u.usedIn ?? 0 }}</td>
                <td>
                    <div class="actions">
                        <template v-if="editId === u.id">
                            <button
                                type="button"
                                class="tbl-btn"
                                aria-label="Zapisz"
                                title="Zapisz"
                                @click="saveEdit(u)"
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
                                :aria-label="`Edytuj ${u.name}`"
                                title="Edytuj"
                                @click="startEdit(u)"
                            >
                                <UIcon name="edit" />
                            </button>
                            <button
                                type="button"
                                class="tbl-btn tbl-btn--danger"
                                :aria-label="`Usuń ${u.name}`"
                                title="Usuń"
                                @click="remove(u)"
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
            v-model="newName"
            class="input add-row__input"
            aria-label="Nowa jednostka"
            placeholder="Nowa jednostka"
            @keydown.enter="add"
        />
        <button
            type="button"
            class="btn btn--secondary"
            :disabled="!newName.trim() || busy"
            @click="add"
        >
            + Dodaj
        </button>
    </div>

    <DeleteValueModal
        v-if="deleting"
        :label="deleting.name"
        :count="deleting.usedIn ?? 0"
        :count-noun="usesNoun(deleting.usedIn ?? 0)"
        :others="
            props.units
                .filter((u) => u.id !== deleting?.id)
                .map((u) => ({ id: u.id, label: u.name }))
        "
        :required="false"
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
.edit-input {
    padding: 6px 10px;
    font-size: 13px;
}
.forms {
    font-size: 13px;
    color: var(--ink-dim);
}
.forms.is-empty {
    color: var(--ink-muted);
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
