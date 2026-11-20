<script lang="ts" setup>
import { useQueryClient } from '@tanstack/vue-query';
import { computed, onMounted, ref } from 'vue';

import UIcon from '@/components/UIcon.vue';
import UModal from '@/components/UModal.vue';
import USpinner from '@/components/USpinner.vue';
import { confirmAction } from '@/composables/useConfirm';
import { showToast } from '@/composables/useToast';
import {
    deleteIngredient,
    fetchIngredientsWithClient,
    invalidateIngredients,
} from '@/data/api/ingredients/fetch';
import { fetchDictionariesWithClient } from '@/data/api/taxonomy/fetch';
import type { IngredientDto } from '@/data/dtos/ingredients/IngredientDto';
import type { IDictionaries } from '@/data/dtos/taxonomy/DictionariesDto';

// Ingredient database: type, EU allergens, usage count. An ingredient in use can't be deleted.
const qc = useQueryClient();
const ingredients = ref<IngredientDto[]>([]);
const dicts = ref<IDictionaries>();
const loading = ref(true);
const search = ref('');
const onlyAllergens = ref(false);
const blocked = ref<IngredientDto | null>(null);

onMounted(async () => {
    [ingredients.value, dicts.value] = await Promise.all([
        fetchIngredientsWithClient(qc),
        fetchDictionariesWithClient(qc),
    ]);
    loading.value = false;
});

const allergenLabel = computed(() =>
    Object.fromEntries((dicts.value?.allergens ?? []).map((a) => [a.code, a.label]))
);

const filtered = computed(() => {
    const q = search.value.trim().toLowerCase();
    return ingredients.value.filter(
        (i) =>
            `${i.name} ${i.type?.name ?? ''}`.toLowerCase().includes(q) &&
            (!onlyAllergens.value || i.allergens.length)
    );
});

const plural = (n: number) => (n === 1 ? 'przepisie' : 'przepisach');

async function remove(ing: IngredientDto) {
    if (ing.usedIn > 0) {
        blocked.value = ing;
        return;
    }
    const ok = await confirmAction({
        title: `Usunąć „${ing.name}”?`,
        body: 'Składnik nie jest używany w żadnym przepisie.',
        confirmLabel: 'Usuń',
        danger: true,
    });
    if (!ok) return;
    try {
        await deleteIngredient(ing.id);
        ingredients.value = ingredients.value.filter((i) => i.id !== ing.id);
        invalidateIngredients(qc);
        showToast('Składnik usunięty');
    } catch {
        showToast('Nie udało się usunąć składnika');
    }
}
</script>

<template>
    <h1 class="admin-title">Składniki</h1>
    <div class="admin-toolbar">
        <input
            v-model="search"
            type="search"
            class="input admin-search"
            aria-label="Szukaj składnika"
            placeholder="Szukaj składnika…"
        />
        <button
            type="button"
            class="chip"
            :class="{ 'is-active': onlyAllergens }"
            :aria-pressed="onlyAllergens"
            @click="onlyAllergens = !onlyAllergens"
        >
            z alergenami
        </button>
        <RouterLink
            to="/admin/ingredients/new"
            class="btn btn--accent"
        >
            + Nowy składnik
        </RouterLink>
    </div>

    <div class="tbl-wrap">
        <div
            v-if="loading"
            class="tbl-loading"
            role="status"
            aria-label="Wczytywanie składników"
        >
            <USpinner :size="22" />
        </div>
        <table
            class="data-table"
            :aria-busy="loading"
        >
            <thead>
                <tr>
                    <th>Nazwa</th>
                    <th>Typ</th>
                    <th>Alergeny</th>
                    <th class="num">Użycia</th>
                    <th class="col-actions"><span class="sr-only">Akcje</span></th>
                </tr>
            </thead>
            <tbody
                v-if="loading"
                aria-hidden="true"
            >
                <tr
                    v-for="n in 12"
                    :key="n"
                >
                    <td class="cell-name"><span class="sk sk-cell" /></td>
                    <td class="cell-type"><span class="sk sk-cell" /></td>
                    <td><span class="muted">—</span></td>
                    <td class="num"><span class="sk sk-cell sk-cell--short" /></td>
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
                    v-for="ing in filtered"
                    :key="ing.id"
                >
                    <td class="cell-name">{{ ing.name }}</td>
                    <td class="cell-type">{{ ing.type?.name }}</td>
                    <td>
                        <div
                            v-if="ing.allergens.length"
                            class="allergen-tags"
                        >
                            <span
                                v-for="a in ing.allergens"
                                :key="a"
                                class="allergen-tag"
                            >
                                {{ allergenLabel[a] ?? a }}
                            </span>
                        </div>
                        <span
                            v-else
                            class="muted"
                        >
                            —
                        </span>
                    </td>
                    <td class="num">{{ ing.usedIn }}</td>
                    <td>
                        <div class="actions">
                            <RouterLink
                                :to="`/admin/ingredients/${ing.id}`"
                                class="tbl-btn"
                                :aria-label="`Edytuj ${ing.name}`"
                                title="Edytuj"
                            >
                                <UIcon name="edit" />
                            </RouterLink>
                            <button
                                type="button"
                                class="tbl-btn tbl-btn--danger"
                                :aria-label="`Usuń ${ing.name}`"
                                title="Usuń"
                                @click="remove(ing)"
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
                        Brak składników
                    </td>
                </tr>
            </tbody>
        </table>
    </div>

    <UModal
        v-if="blocked"
        :title="`Nie można usunąć „${blocked.name}”`"
        @close="blocked = null"
    >
        <p class="modal__text">
            Składnik jest używany w {{ blocked.usedIn }} {{ plural(blocked.usedIn) }}. Najpierw usuń
            go z tych przepisów albo zmień jego nazwę, jeśli to duplikat.
        </p>
        <template #footer>
            <button
                type="button"
                class="btn btn--accent"
                autofocus
                @click="blocked = null"
            >
                Rozumiem
            </button>
        </template>
    </UModal>
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
/* loading: the table right away, skeleton rows on the real cell classes (same height),
   spinner above the table, outside the layout */
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
.sk-cell--short {
    width: 28px;
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
.cell-name {
    font-weight: 600;
}
.cell-type,
.muted {
    font-size: 13px;
    color: var(--ink-muted);
}
.allergen-tags {
    display: flex;
    flex-wrap: wrap;
    gap: 4px;
}
.allergen-tag {
    font-size: 12px;
    font-weight: 600;
    color: var(--danger);
    background: var(--danger-soft);
    padding: 2px 8px;
    border-radius: var(--r-pill);
}
</style>
