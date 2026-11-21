<script lang="ts" setup>
import { computed, onMounted, ref } from 'vue';
import { useRoute } from 'vue-router';

import DeleteTermModal from './components/DeleteTermModal.vue';
import TermFormModal from './components/TermFormModal.vue';

import UIcon from '@/components/UIcon.vue';
import USpinner from '@/components/USpinner.vue';
import { confirmAction } from '@/composables/useConfirm';
import { showToast } from '@/composables/useToast';
import { deleteTerm, fetchGlossary, type IGlossaryTerm } from '@/data/api/glossary/fetch';
import { matches, prepareQuery } from '@/shared/utils/search';

// Glossary: terms marked in recipe steps and hints. Content that will grow into a Wiki, so it's its own module,
// not a classification axis. ?edit=<id> opens the term (link "Edytuj w słowniczku" from the recipe form).
const route = useRoute();
const terms = ref<IGlossaryTerm[]>([]);
const loading = ref(true);
const search = ref('');
/** null = new term; undefined = closed. */
const editing = ref<IGlossaryTerm | null | undefined>(undefined);
const removing = ref<IGlossaryTerm | null>(null);

const byName = (a: IGlossaryTerm, b: IGlossaryTerm) => a.name.localeCompare(b.name, 'pl');

onMounted(async () => {
    try {
        terms.value = await fetchGlossary();
    } catch {
        showToast('Nie udało się wczytać słowniczka');
    }
    loading.value = false;
    const id = Number(route.query.edit);
    if (id) editing.value = terms.value.find((t) => t.id === id);
});

const filtered = computed(() => {
    const q = prepareQuery(search.value);
    return q ? terms.value.filter((t) => matches(q, `${t.name} ${t.explanation}`)) : terms.value;
});

const usedText = (n = 0) => `${n} przep.`;

function saved(t: IGlossaryTerm) {
    const i = terms.value.findIndex((x) => x.id === t.id);
    if (i >= 0) terms.value.splice(i, 1, t);
    else terms.value.push(t);
    terms.value.sort(byName);
    editing.value = undefined;
    showToast(i >= 0 ? 'Pojęcie zapisane' : 'Pojęcie dodane');
}

async function remove(t: IGlossaryTerm) {
    if (t.usedIn) {
        removing.value = t;
        return;
    }
    const ok = await confirmAction({
        title: `Usunąć „${t.name}”?`,
        body: 'Pojęcie nie jest oznaczone w żadnym przepisie.',
        confirmLabel: 'Usuń pojęcie',
        danger: true,
    });
    if (ok) await doDelete(t, {});
}

async function doDelete(t: IGlossaryTerm, opts: { replaceWith?: number; removeMarks?: boolean }) {
    removing.value = null;
    try {
        await deleteTerm(t.id, opts);
        terms.value = terms.value.filter((x) => x.id !== t.id);
        // replaced: the other term is now used in these recipes too; the count comes from the server
        if (opts.replaceWith) terms.value = (await fetchGlossary()).sort(byName);
        showToast('Pojęcie usunięte');
    } catch {
        showToast('Nie udało się usunąć pojęcia');
    }
}
</script>

<template>
    <h1 class="admin-title">Słowniczek</h1>
    <div class="admin-toolbar">
        <input
            v-model="search"
            type="search"
            class="input admin-search"
            aria-label="Szukaj pojęcia"
            placeholder="Szukaj pojęcia…"
        />
        <button
            type="button"
            class="btn btn--accent"
            @click="editing = null"
        >
            + Nowe pojęcie
        </button>
    </div>

    <div class="tbl-wrap">
        <div
            v-if="loading"
            class="tbl-loading"
            role="status"
            aria-label="Wczytywanie słowniczka"
        >
            <USpinner :size="22" />
        </div>
        <table
            class="data-table gloss-table"
            :aria-busy="loading"
        >
            <thead>
                <tr>
                    <th class="col-name">Pojęcie</th>
                    <th>Wyjaśnienie</th>
                    <th class="num col-used">Użyte w</th>
                    <th class="col-actions"><span class="sr-only">Akcje</span></th>
                </tr>
            </thead>
            <tbody v-if="!loading">
                <tr
                    v-for="t in filtered"
                    :key="t.id"
                >
                    <td class="cell-name">{{ t.name }}</td>
                    <td class="cell-expl">
                        <span class="cell-expl__text">{{ t.explanation }}</span>
                    </td>
                    <td class="num cell-used">{{ usedText(t.usedIn) }}</td>
                    <td>
                        <div class="actions">
                            <button
                                type="button"
                                class="tbl-btn"
                                :aria-label="`Edytuj ${t.name}`"
                                title="Edytuj"
                                @click="editing = t"
                            >
                                <UIcon name="edit" />
                            </button>
                            <button
                                type="button"
                                class="tbl-btn tbl-btn--danger"
                                :aria-label="`Usuń ${t.name}`"
                                title="Usuń"
                                @click="remove(t)"
                            >
                                <UIcon name="trash" />
                            </button>
                        </div>
                    </td>
                </tr>
                <tr v-if="!filtered.length">
                    <td
                        colspan="4"
                        class="empty-row"
                    >
                        {{
                            terms.length
                                ? 'Nic nie pasuje do wyszukiwania'
                                : 'Słowniczek jest pusty'
                        }}
                    </td>
                </tr>
            </tbody>
        </table>
    </div>

    <TermFormModal
        v-if="editing !== undefined"
        :term="editing"
        @saved="saved"
        @close="editing = undefined"
    />
    <DeleteTermModal
        v-if="removing"
        :term="removing"
        :others="terms.filter((t) => t.id !== removing!.id)"
        @confirm="doDelete(removing!, $event)"
        @close="removing = null"
    />
</template>

<style scoped>
.admin-title {
    font-size: 28px;
    font-weight: 600;
    letter-spacing: -0.8px;
    margin-bottom: 22px;
}
.admin-toolbar {
    display: flex;
    gap: 10px;
    align-items: center;
    margin-bottom: 16px;
}
.admin-search {
    flex: 1;
    max-width: 320px;
    height: 40px;
    background: var(--bg-alt);
}
.admin-toolbar .btn {
    margin-left: auto;
}
.tbl-wrap {
    position: relative;
    min-height: 120px;
}
.tbl-loading {
    position: absolute;
    left: 0;
    right: 0;
    top: 56px;
    display: flex;
    justify-content: center;
    pointer-events: none;
}
.gloss-table {
    table-layout: fixed;
    width: 100%;
}
.col-name {
    width: 26%;
}
.col-used {
    width: 96px;
}
.col-actions {
    width: 84px;
}
.cell-name {
    font-weight: 600;
}
.cell-expl {
    color: var(--ink-dim);
}
/* one line, cut off: the full text is in the edit dialog */
.cell-expl__text {
    display: block;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}
.cell-used {
    color: var(--ink-muted);
}
</style>
