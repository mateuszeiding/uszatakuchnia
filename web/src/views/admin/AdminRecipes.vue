<script lang="ts" setup>
import { computed, ref } from 'vue';

import { deleteRecipe, fetchRecipes } from '@/data/api/recipe/fetch';
import type { RecipeBaseDto } from '@/data/dtos/recipe/RecipeDto';
import { confirmDialog } from '@/utils/dialog';

const recipes = ref<RecipeBaseDto[]>([]);
const tab = ref<'all' | 'pub' | 'draft'>('all');
const search = ref('');
fetchRecipes('list', true).then((v) => (recipes.value = v));

const filtered = computed(() => {
    let list = recipes.value;
    if (tab.value === 'pub') list = list.filter((r) => r.status === 'published');
    if (tab.value === 'draft') list = list.filter((r) => r.status === 'draft');
    if (search.value.trim()) {
        const q = search.value.toLowerCase();
        list = list.filter((r) => r.name.toLowerCase().includes(q));
    }
    return list;
});

const counts = computed(() => ({
    all: recipes.value.length,
    pub: recipes.value.filter((r) => r.status === 'published').length,
    draft: recipes.value.filter((r) => r.status === 'draft').length,
}));

async function handleDelete(r: RecipeBaseDto) {
    if (!confirmDialog(`Usuń przepis „${r.name}"?`)) return;
    await deleteRecipe(r.id);
    recipes.value = recipes.value.filter((x) => x.id !== r.id);
}

function fmtTime(min: number | null) {
    if (!min) return '—';
    if (min < 60) return `${min} min`;
    const h = Math.floor(min / 60);
    const m = min % 60;
    return m ? `${h}h ${m}min` : `${h}h`;
}
</script>

<template>
    <div class="container admin-recipes-page">
        <div class="page-header">
            <div>
                <div class="kicker">// admin / przepisy</div>
                <h1 class="page-title">Przepisy</h1>
            </div>
            <RouterLink
                to="/recipes/new"
                class="btn btn--accent"
                style="display: inline-flex; align-items: center; gap: 6px; text-decoration: none"
            >
                <svg
                    width="13"
                    height="13"
                    viewBox="0 0 13 13"
                    fill="none"
                >
                    <path
                        d="M6.5 2v9M2 6.5h9"
                        stroke="currentColor"
                        stroke-width="1.6"
                        stroke-linecap="round"
                    />
                </svg>
                Nowy przepis
            </RouterLink>
        </div>

        <div class="toolbar">
            <div class="tabs">
                <button
                    v-for="t in [
                        { id: 'all', label: 'Wszystkie', n: counts.all },
                        { id: 'pub', label: 'Opublikowane', n: counts.pub },
                        { id: 'draft', label: 'Szkice', n: counts.draft },
                    ] as const"
                    :key="t.id"
                    class="tab"
                    :class="{ 'tab--active': tab === t.id }"
                    @click="tab = t.id"
                >
                    {{ t.label }}
                    <span class="tab__n">{{ t.n }}</span>
                </button>
            </div>
            <input
                v-model="search"
                class="input search-input"
                placeholder="Szukaj przepisu..."
            />
        </div>

        <div class="recipe-table">
            <div class="recipe-table__head">
                <span>#</span>
                <span>Nazwa</span>
                <span>Kategoria</span>
                <span>Czas</span>
                <span>Status</span>
                <span></span>
            </div>
            <div
                v-for="r in filtered"
                :key="r.id"
                class="recipe-table__row"
            >
                <span class="cell-id">#{{ String(r.id).padStart(3, '0') }}</span>
                <span class="cell-name">
                    {{ r.name }}
                    <span
                        v-if="r.needsPrep"
                        class="prep-tag"
                    >
                        prep+
                    </span>
                </span>
                <span>
                    <span
                        v-if="r.category"
                        class="badge"
                        :class="`cat-${r.category}`"
                    >
                        {{ r.category }}
                    </span>
                </span>
                <span class="cell-meta">{{ fmtTime(r.timeMinutes) }}</span>
                <span>
                    <span
                        class="status-pill"
                        :class="
                            r.status === 'published' ? 'status-pill--pub' : 'status-pill--draft'
                        "
                    >
                        {{ r.status === 'published' ? 'live' : 'szkic' }}
                    </span>
                </span>
                <span class="cell-actions">
                    <RouterLink
                        :to="{ name: 'recipe-edit', params: { id: r.id } }"
                        class="action-btn"
                        title="Edytuj"
                    >
                        <svg
                            width="14"
                            height="14"
                            viewBox="0 0 14 14"
                            fill="none"
                        >
                            <path
                                d="M2 10.5L10 2.5l1.5 1.5L3.5 12H2v-1.5z"
                                stroke="currentColor"
                                stroke-width="1.3"
                                stroke-linejoin="round"
                            />
                        </svg>
                    </RouterLink>
                    <button
                        class="action-btn action-btn--danger"
                        title="Usuń"
                        @click="handleDelete(r)"
                    >
                        <svg
                            width="14"
                            height="14"
                            viewBox="0 0 14 14"
                            fill="none"
                        >
                            <path
                                d="M3 4h8M5.5 4V2.5h3V4M4 4l0.5 8h5L10 4"
                                stroke="currentColor"
                                stroke-width="1.3"
                                stroke-linecap="round"
                                stroke-linejoin="round"
                            />
                        </svg>
                    </button>
                </span>
            </div>
            <div
                v-if="filtered.length === 0"
                class="empty"
            >
                Brak przepisów
            </div>
        </div>
    </div>
</template>

<style scoped>
.admin-recipes-page {
    padding-top: 48px;
    padding-bottom: 80px;
}
.page-header {
    display: flex;
    align-items: flex-end;
    justify-content: space-between;
    gap: 24px;
    padding-bottom: 28px;
}
.kicker {
    font-family: var(--font-mono);
    font-size: 11px;
    color: var(--accent);
    letter-spacing: 0.4px;
    margin-bottom: 12px;
}
.page-title {
    font-size: 56px;
    line-height: 0.95;
    letter-spacing: -2px;
    font-weight: 500;
    margin: 0;
}
.toolbar {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 16px;
    margin-bottom: 16px;
}
.tabs {
    display: flex;
    gap: 0;
    border-bottom: 1px solid var(--rule);
}
.tab {
    background: transparent;
    border: none;
    padding: 10px 16px;
    border-bottom: 2px solid transparent;
    color: var(--ink-muted);
    font-size: 13px;
    font-weight: 600;
    cursor: pointer;
    font-family: var(--font-sans);
    display: inline-flex;
    align-items: center;
    gap: 7px;
    margin-bottom: -1px;
}
.tab--active {
    border-bottom-color: var(--accent);
    color: var(--ink);
}
.tab__n {
    font-family: var(--font-mono);
    font-size: 10px;
    background: var(--bg-alt);
    color: var(--ink-muted);
    padding: 2px 6px;
    border-radius: var(--r-pill);
}
.tab--active .tab__n {
    background: var(--accent-soft);
    color: var(--accent);
}
.search-input {
    max-width: 280px;
}
.recipe-table {
    background: var(--bg-alt);
    border: 1px solid var(--rule);
    border-radius: var(--r-lg);
    overflow: hidden;
}
.recipe-table__head {
    display: grid;
    grid-template-columns: 48px 1fr 120px 80px 70px 72px;
    gap: 14px;
    padding: 10px 20px;
    border-bottom: 1px solid var(--rule);
    font-family: var(--font-mono);
    font-size: 10px;
    color: var(--ink-muted);
    letter-spacing: 0.5px;
    text-transform: uppercase;
}
.recipe-table__row {
    display: grid;
    grid-template-columns: 48px 1fr 120px 80px 70px 72px;
    gap: 14px;
    align-items: center;
    padding: 12px 20px;
    border-bottom: 1px solid var(--rule);
    transition: background 0.1s;
}
.recipe-table__row:last-child {
    border-bottom: none;
}
.recipe-table__row:hover {
    background: var(--bg);
}
.cell-id {
    font-family: var(--font-mono);
    font-size: 11px;
    color: var(--ink-faint);
}
.cell-name {
    font-size: 14px;
    font-weight: 600;
    display: flex;
    align-items: center;
    gap: 8px;
    overflow: hidden;
}
.prep-tag {
    font-family: var(--font-mono);
    font-size: 9px;
    color: #d9a441;
    font-weight: 600;
    flex-shrink: 0;
}
.cell-meta {
    font-family: var(--font-mono);
    font-size: 12px;
    color: var(--ink-muted);
}
.cell-actions {
    display: flex;
    gap: 6px;
    justify-content: flex-end;
}
.action-btn {
    width: 28px;
    height: 28px;
    border-radius: var(--r-md);
    background: var(--bg);
    border: 1px solid var(--rule);
    color: var(--ink-muted);
    cursor: pointer;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    text-decoration: none;
    transition: all 0.12s;
}
.action-btn:hover {
    border-color: var(--accent);
    color: var(--accent);
}
.action-btn--danger:hover {
    border-color: #c53728;
    color: #c53728;
}
.status-pill {
    font-family: var(--font-mono);
    font-size: 9px;
    font-weight: 700;
    letter-spacing: 0.5px;
    text-transform: uppercase;
    padding: 3px 7px;
    border-radius: var(--r-sm);
}
.status-pill--pub {
    background: color-mix(in srgb, #2f6b3f 14%, transparent);
    color: #2f6b3f;
}
.status-pill--draft {
    background: color-mix(in srgb, #d9a441 18%, transparent);
    color: #d9a441;
}
.empty {
    padding: 40px;
    text-align: center;
    font-family: var(--font-mono);
    font-size: 12px;
    color: var(--ink-muted);
}
</style>
