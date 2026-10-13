<script lang="ts" setup>
import { useQueryClient } from '@tanstack/vue-query';
import { computed, onMounted, ref } from 'vue';

import UIcon from '@/components/UIcon.vue';
import USpinner from '@/components/USpinner.vue';
import { confirmAction } from '@/composables/useConfirm';
import { showToast } from '@/composables/useToast';
import { deleteRecipe, fetchRecipesWithClient, invalidateRecipes } from '@/data/api/recipe/fetch';
import {
    mainDish,
    type RecipeBaseDto,
    type RecipeStatus,
    STATUS_LABEL,
} from '@/data/dtos/recipe/RecipeDto';

// Jedyna lista przepisów dla admina: zakładki statusów, podgląd / edycja / usuwanie.
const qc = useQueryClient();
const recipes = ref<RecipeBaseDto[]>([]);
const loading = ref(true);
const tab = ref<'all' | RecipeStatus>('all');
const search = ref('');

async function load() {
    recipes.value = await fetchRecipesWithClient(qc, 'list', true);
    loading.value = false;
}
onMounted(load);

const counts = computed(() => ({
    all: recipes.value.length,
    published: recipes.value.filter((r) => r.status === 'published').length,
    draft: recipes.value.filter((r) => r.status === 'draft').length,
}));

const TABS = [
    { id: 'all', label: 'Wszystkie' },
    { id: 'published', label: 'Opublikowane' },
    { id: 'draft', label: 'Szkice' },
] as const;

const filtered = computed(() => {
    const q = search.value.trim().toLowerCase();
    return recipes.value.filter(
        (r) => (tab.value === 'all' || r.status === tab.value) && r.name.toLowerCase().includes(q)
    );
});

async function remove(r: RecipeBaseDto) {
    const ok = await confirmAction({
        title: `Usunąć „${r.name}”?`,
        body:
            r.status === 'published'
                ? 'Przepis jest opublikowany i zniknie ze strony. Tej operacji nie da się cofnąć.'
                : 'Tej operacji nie da się cofnąć.',
        confirmLabel: 'Usuń przepis',
        danger: true,
    });
    if (!ok) return;
    await deleteRecipe(r.id);
    recipes.value = recipes.value.filter((x) => x.id !== r.id);
    invalidateRecipes(qc);
    showToast('Przepis usunięty');
}
</script>

<template>
    <h1 class="admin-title">Przepisy</h1>
    <div
        class="tabs"
        role="tablist"
    >
        <button
            v-for="t in TABS"
            :key="t.id"
            type="button"
            role="tab"
            class="tab"
            :class="{ 'is-active': tab === t.id }"
            :aria-selected="tab === t.id"
            @click="tab = t.id"
        >
            {{ t.label }}
            <span class="count">{{ counts[t.id] }}</span>
        </button>
    </div>
    <div class="admin-toolbar">
        <input
            v-model="search"
            type="search"
            class="input admin-search"
            aria-label="Szukaj przepisu"
            placeholder="Szukaj przepisu…"
        />
        <RouterLink
            to="/recipes/new"
            class="btn btn--accent"
        >
            + Nowy przepis
        </RouterLink>
    </div>

    <div
        v-if="loading"
        class="admin-loading"
    >
        <USpinner :size="22" />
    </div>
    <table
        v-else
        class="data-table"
    >
        <thead>
            <tr>
                <th class="col-photo"><span class="sr-only">Zdjęcie</span></th>
                <th>Tytuł</th>
                <th>Typ dania</th>
                <th>Status</th>
                <th class="col-actions"><span class="sr-only">Akcje</span></th>
            </tr>
        </thead>
        <tbody>
            <tr
                v-for="r in filtered"
                :key="r.id"
            >
                <td class="col-photo">
                    <img
                        v-if="r.photoUrl"
                        :src="r.photoUrl"
                        alt=""
                        class="thumb"
                    />
                    <div
                        v-else
                        class="photo thumb"
                    />
                </td>
                <td class="cell-title">
                    <RouterLink :to="{ name: 'recipe-edit', params: { id: r.id } }">
                        {{ r.name }}
                    </RouterLink>
                </td>
                <td class="cell-dish">
                    <span
                        v-if="mainDish(r)"
                        class="dish"
                        :class="mainDish(r)?.color ? `cat-${mainDish(r)?.color}` : null"
                    >
                        <span
                            class="dish__dot"
                            aria-hidden="true"
                        />
                        {{ mainDish(r)?.label }}
                    </span>
                </td>
                <td>
                    <span
                        class="status-tag"
                        :class="`status-tag--${r.status}`"
                    >
                        {{ STATUS_LABEL[r.status] }}
                    </span>
                </td>
                <td>
                    <div class="actions">
                        <RouterLink
                            v-if="r.status === 'published'"
                            :to="{ name: 'recipe-details', params: { id: r.id } }"
                            class="tbl-btn"
                            :aria-label="`Zobacz „${r.name}”`"
                            title="Zobacz na stronie"
                        >
                            <UIcon name="eye" />
                        </RouterLink>
                        <RouterLink
                            :to="{ name: 'recipe-edit', params: { id: r.id } }"
                            class="tbl-btn"
                            :aria-label="`Edytuj „${r.name}”`"
                            title="Edytuj"
                        >
                            <UIcon name="edit" />
                        </RouterLink>
                        <button
                            type="button"
                            class="tbl-btn tbl-btn--danger"
                            :aria-label="`Usuń „${r.name}”`"
                            title="Usuń"
                            @click="remove(r)"
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
                    Brak przepisów
                </td>
            </tr>
        </tbody>
    </table>
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
    margin: 16px 0;
}
.admin-search {
    flex: 1;
    height: 40px;
    background: var(--bg-alt);
}
.admin-loading {
    display: flex;
    justify-content: center;
    padding: 60px 0;
}
.col-photo {
    width: 54px;
    padding-right: 0;
}
.col-actions {
    width: 112px;
}
.thumb {
    width: 40px;
    height: 40px;
    border-radius: var(--r-md);
    object-fit: cover;
}
.cell-title a {
    color: var(--ink);
    font-weight: 600;
}
.cell-dish {
    white-space: nowrap;
}
.dish {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    font-size: 13px;
}
.dish__dot {
    width: 7px;
    height: 7px;
    border-radius: var(--r-pill);
    background: var(--cat, var(--accent));
}
</style>
