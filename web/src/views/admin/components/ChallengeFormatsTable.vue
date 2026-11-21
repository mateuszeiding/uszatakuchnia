<script lang="ts" setup>
import { ref } from 'vue';

import DeleteFormatModal from './DeleteFormatModal.vue';

import UIcon from '@/components/UIcon.vue';
import { confirmAction } from '@/composables/useConfirm';
import { showToast } from '@/composables/useToast';
import {
    createChallengeFormat,
    deleteChallengeFormat,
    updateChallengeFormat,
} from '@/data/api/challenges/fetch';
import type { IChallengeFormat } from '@/data/dtos/challenge/ChallengeDto';
import { plural } from '@/shared/utils/challenges';

const props = defineProps<{ items: IChallengeFormat[] }>();
const emit = defineEmits<{ changed: [] }>();

const newName = ref('');
const editId = ref<number | null>(null);
const buf = ref('');
const deleting = ref<IChallengeFormat | null>(null);
const busy = ref(false);

const usedText = (n: number) =>
    n ? `${n} ${plural(n, 'wyzwaniu', 'wyzwaniach', 'wyzwaniach')}` : '—';

function apiError(e: unknown): string {
    return (e as { response?: { status?: number } }).response?.status === 409
        ? 'Taki format już jest'
        : 'Nie udało się zapisać';
}

async function add() {
    const name = newName.value.trim();
    if (!name || busy.value) return;
    busy.value = true;
    try {
        await createChallengeFormat(name);
        newName.value = '';
        showToast(`Dodano „${name}”`);
        emit('changed');
    } catch (e) {
        showToast(apiError(e));
    } finally {
        busy.value = false;
    }
}

function startEdit(x: IChallengeFormat) {
    editId.value = x.id;
    buf.value = x.name;
}

async function saveEdit(x: IChallengeFormat) {
    if (!buf.value.trim() || busy.value) return;
    busy.value = true;
    try {
        await updateChallengeFormat(x.id, buf.value);
        editId.value = null;
        showToast('Zapisano');
        emit('changed');
    } catch (e) {
        showToast(apiError(e));
    } finally {
        busy.value = false;
    }
}

async function remove(x: IChallengeFormat) {
    if (x.usedIn) {
        deleting.value = x;
        return;
    }
    const ok = await confirmAction({
        title: `Usunąć „${x.name}”?`,
        body: 'Ten format nie jest używany w żadnym wyzwaniu.',
        confirmLabel: 'Usuń',
        danger: true,
    });
    if (ok) await doDelete(x);
}

async function doDelete(x: IChallengeFormat, moveTo?: number) {
    deleting.value = null;
    try {
        const res = await deleteChallengeFormat(x.id, moveTo);
        const target = props.items.find((i) => i.id === moveTo);
        showToast(
            res.moved && target
                ? `Usunięto. ${res.moved} ${plural(res.moved, 'wyzwanie ma', 'wyzwania mają', 'wyzwań ma')} teraz „${target.name}”.`
                : 'Usunięto'
        );
        emit('changed');
    } catch {
        showToast('Nie udało się usunąć');
    }
}
</script>

<template>
    <div class="cf-add">
        <input
            v-model="newName"
            class="input cf-add__input"
            aria-label="Nowy format"
            placeholder="Nowy format, np. „w 30 minut”"
            maxlength="64"
            @keydown.enter="add"
        />
        <button
            type="button"
            class="btn btn--accent"
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
                v-for="x in items"
                :key="x.id"
                :class="{ 'is-editing': editId === x.id }"
            >
                <td>
                    <input
                        v-if="editId === x.id"
                        v-model="buf"
                        class="input cf-edit"
                        :aria-label="`Nazwa: ${x.name}`"
                        maxlength="64"
                        autofocus
                        @keydown.enter="saveEdit(x)"
                        @keydown.esc="editId = null"
                    />
                    <span v-else>{{ x.name }}</span>
                </td>
                <td
                    class="col-used cf-used"
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
            <tr v-if="!items.length">
                <td
                    colspan="3"
                    class="cf-empty"
                >
                    Nie ma formatów. Dodaj pierwszy powyżej — format jest wymagany w wyzwaniu.
                </td>
            </tr>
        </tbody>
    </table>

    <DeleteFormatModal
        v-if="deleting"
        :item="deleting"
        :others="items.filter((x) => x.id !== deleting?.id)"
        @confirm="doDelete(deleting!, $event)"
        @close="deleting = null"
    />
</template>

<style scoped>
.cf-add {
    display: flex;
    gap: 10px;
    margin-bottom: 14px;
}
.cf-add__input {
    flex: 1;
    background: var(--bg-alt);
}
.col-used {
    width: 150px;
}
.col-actions {
    width: 90px;
}
tr.is-editing td {
    background: var(--accent-soft);
}
.cf-edit {
    padding: 6px 10px;
}
.cf-used {
    font-family: var(--font-mono);
    font-size: 13px;
    color: var(--ink-dim);
}
.cf-used.is-unused {
    color: var(--ink-muted);
}
.cf-empty {
    color: var(--ink-muted);
    font-size: 13px;
}
</style>
