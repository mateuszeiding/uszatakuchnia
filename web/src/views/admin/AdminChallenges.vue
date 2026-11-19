<script lang="ts" setup>
import { useQueryClient } from '@tanstack/vue-query';
import { computed, onMounted, ref } from 'vue';

import UIcon from '@/components/UIcon.vue';
import USpinner from '@/components/USpinner.vue';
import { confirmAction } from '@/composables/useConfirm';
import { showToast } from '@/composables/useToast';
import {
    deleteChallenge,
    fetchChallengesWithClient,
    invalidateChallenges,
} from '@/data/api/challenges/fetch';
import { invalidateDictionaries } from '@/data/api/taxonomy/fetch';
import type { IChallenge } from '@/data/dtos/challenge/ChallengeDto';
import { timeLimitText } from '@/shared/utils/challenges';
import ChallengeDifficulty from '@/views/challenges/components/ChallengeDifficulty.vue';

// Panel → Wyzwania: wyszukiwanie, tabela (tytuł, format, trudność, czas), edycja i usuwanie z potwierdzeniem.
const qc = useQueryClient();
const items = ref<IChallenge[]>([]);
const loading = ref(true);
const search = ref('');

onMounted(async () => {
    try {
        items.value = await fetchChallengesWithClient(qc);
    } catch {
        showToast('Nie udało się wczytać wyzwań');
    } finally {
        loading.value = false;
    }
});

const filtered = computed(() => {
    const q = search.value.trim().toLowerCase();
    return items.value.filter((c) => `${c.title} ${c.format}`.toLowerCase().includes(q));
});

async function remove(c: IChallenge) {
    const ok = await confirmAction({
        title: `Usunąć „${c.title}”?`,
        body: 'Wyzwanie zniknie ze strony. Tej operacji nie da się cofnąć.',
        confirmLabel: 'Usuń wyzwanie',
        danger: true,
    });
    if (!ok) return;
    try {
        await deleteChallenge(c.id);
        items.value = items.value.filter((x) => x.id !== c.id);
        invalidateChallenges(qc);
        // liczba użyć formatu w Klasyfikacji
        invalidateDictionaries(qc);
        showToast('Wyzwanie usunięte');
    } catch {
        showToast('Nie udało się usunąć wyzwania');
    }
}
</script>

<template>
    <h1 class="admin-title">Wyzwania</h1>
    <div class="admin-toolbar">
        <input
            v-model="search"
            type="search"
            class="input admin-search"
            aria-label="Szukaj wyzwania"
            placeholder="Szukaj wyzwania…"
        />
        <RouterLink
            to="/admin/challenges/new"
            class="btn btn--accent"
        >
            + Nowe wyzwanie
        </RouterLink>
    </div>

    <div class="tbl-wrap">
        <div
            v-if="loading"
            class="tbl-loading"
            role="status"
            aria-label="Wczytywanie wyzwań"
        >
            <USpinner :size="22" />
        </div>
        <table
            class="data-table"
            :aria-busy="loading"
        >
            <thead>
                <tr>
                    <th>Tytuł</th>
                    <th>Format</th>
                    <th>Trudność</th>
                    <th>Czas</th>
                    <th class="col-actions"><span class="sr-only">Akcje</span></th>
                </tr>
            </thead>
            <tbody
                v-if="loading"
                aria-hidden="true"
            >
                <tr
                    v-for="n in 8"
                    :key="n"
                >
                    <td><span class="sk sk-cell sk-cell--wide" /></td>
                    <td class="cell-dim"><span class="sk sk-cell" /></td>
                    <td><span class="sk sk-cell sk-cell--diff" /></td>
                    <td class="cell-time"><span class="sk sk-cell sk-cell--short" /></td>
                    <td>
                        <div class="actions">
                            <span class="tbl-btn sk-box"><span class="sk sk-fill" /></span>
                            <span class="tbl-btn sk-box"><span class="sk sk-fill" /></span>
                        </div>
                    </td>
                </tr>
            </tbody>
            <tbody v-else>
                <tr
                    v-for="c in filtered"
                    :key="c.id"
                >
                    <td class="cell-title">{{ c.title }}</td>
                    <td class="cell-dim">{{ c.format }}</td>
                    <td>
                        <ChallengeDifficulty
                            :level="c.difficulty"
                            :size="14"
                            :label="false"
                        />
                    </td>
                    <td class="cell-time">{{ timeLimitText(c.timeLimitMin) }}</td>
                    <td>
                        <div class="actions">
                            <RouterLink
                                :to="`/admin/challenges/${c.id}`"
                                class="tbl-btn"
                                :aria-label="`Edytuj „${c.title}”`"
                                title="Edytuj"
                            >
                                <UIcon name="edit" />
                            </RouterLink>
                            <button
                                type="button"
                                class="tbl-btn tbl-btn--danger"
                                :aria-label="`Usuń „${c.title}”`"
                                title="Usuń"
                                @click="remove(c)"
                            >
                                <UIcon name="trash" />
                            </button>
                        </div>
                    </td>
                </tr>
                <tr v-if="!filtered.length">
                    <td
                        colspan="5"
                        class="empty-row"
                    >
                        {{ items.length ? 'Nic nie pasuje' : 'Jeszcze nie ma wyzwań' }}
                    </td>
                </tr>
            </tbody>
        </table>
    </div>
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
    height: 40px;
    background: var(--bg-alt);
}
/* ładowanie: tabela od razu, wiersze-szkielety w tych samych komórkach, spinner nad tabelą */
.tbl-wrap {
    position: relative;
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
.sk-cell {
    display: inline-block;
    width: 120px;
    height: 12px;
    vertical-align: middle;
}
.sk-cell--wide {
    width: 220px;
}
.sk-cell--diff {
    width: 74px;
}
.sk-cell--short {
    width: 48px;
}
.sk-box {
    position: relative;
    border-color: transparent;
    background: transparent;
}
.sk-fill {
    position: absolute;
    inset: 0;
    border-radius: inherit;
}
.col-actions {
    width: 84px;
}
.cell-title {
    font-weight: 600;
}
.cell-dim {
    font-size: 13px;
    color: var(--ink-dim);
}
.cell-time {
    font-family: var(--font-mono);
    font-size: 13px;
    color: var(--ink-dim);
    white-space: nowrap;
}
</style>
