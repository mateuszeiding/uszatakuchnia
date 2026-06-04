<script lang="ts" setup>
import { computed, ref } from 'vue';
import { useAuth0 } from '@auth0/auth0-vue';

import RecipeCard from './components/RecipeCard.vue';
import { fetchRecipes } from '@/data/api/recipe/fetch';
import type { RecipeBaseDto } from '@/data/dtos/recipe/RecipeDto';

const { isAuthenticated } = useAuth0();

const recipes = ref<RecipeBaseDto[]>([]);
const search = ref('');
const activeCategory = ref<string | null>(null);
const sort = ref<'newest' | 'fastest' | 'az'>('newest');

fetchRecipes('list').then((v) => (recipes.value = v));

const categories = [
    { key: 'ryby',      label: 'ryby' },
    { key: 'mieso',     label: 'mięso' },
    { key: 'wege',      label: 'wege' },
    { key: 'wegan',     label: 'wegan' },
    { key: 'makarony',  label: 'makarony' },
    { key: 'ryz',       label: 'ryż / kasze' },
    { key: 'zupy',      label: 'zupy' },
    { key: 'salatki',   label: 'sałatki' },
    { key: 'pieczywo',  label: 'pieczywo' },
    { key: 'desery',    label: 'desery' },
    { key: 'sniadania', label: 'śniadania' },
    { key: 'przekaski', label: 'przekąski' },
    { key: 'napoje',    label: 'napoje' },
];

const filtered = computed(() => {
    let list = recipes.value;
    if (search.value.trim()) {
        const q = search.value.toLowerCase();
        list = list.filter((r) => r.name.toLowerCase().includes(q));
    }
    if (activeCategory.value) {
        list = list.filter((r) => r.category === activeCategory.value);
    }
    if (sort.value === 'fastest') list = [...list].sort((a, b) => (a.timeMinutes ?? 999) - (b.timeMinutes ?? 999));
    if (sort.value === 'az')      list = [...list].sort((a, b) => a.name.localeCompare(b.name, 'pl'));
    return list;
});

const anyFilter = computed(() => !!search.value || !!activeCategory.value);

function toggleCategory(key: string) {
    activeCategory.value = activeCategory.value === key ? null : key;
}

function clearFilters() {
    search.value = '';
    activeCategory.value = null;
}

function pluralPL(n: number) {
    if (n === 1) return 'przepis';
    const m10 = n % 10, m100 = n % 100;
    if (m10 >= 2 && m10 <= 4 && (m100 < 10 || m100 >= 20)) return 'przepisy';
    return 'przepisów';
}
</script>

<template>
    <div class="recipes-layout container">
        <!-- Sidebar -->
        <aside class="recipes-sidebar">
            <div style="position: relative; margin-bottom: 24px;">
                <svg class="search-icon" width="14" height="14" viewBox="0 0 14 14" fill="none">
                    <circle cx="6" cy="6" r="3.5" stroke="currentColor" stroke-width="1.3"/>
                    <path d="M8.5 8.5L12 12" stroke="currentColor" stroke-width="1.3" stroke-linecap="round"/>
                </svg>
                <input
                    v-model="search"
                    class="input search-input"
                    type="search"
                    placeholder="Szukaj przepisu..."
                />
            </div>

            <button
                v-if="anyFilter"
                class="clear-btn"
                @click="clearFilters"
            >← wyczyść wszystkie filtry</button>

            <div class="filter-group">
                <div class="filter-group__label">Kategoria</div>
                <div style="display: flex; flex-direction: column; gap: 4px;">
                    <button
                        v-for="cat in categories"
                        :key="cat.key"
                        class="filter-row"
                        :class="{ 'is-active': activeCategory === cat.key }"
                        @click="toggleCategory(cat.key)"
                    >
                        <span class="filter-row__box" />
                        <span class="filter-row__dot" :class="`cat-${cat.key}`" />
                        {{ cat.label }}
                    </button>
                </div>
            </div>
        </aside>

        <!-- Main -->
        <main class="recipes-main">
            <div class="page-header">
                <div class="page-header__kicker">
                    // {{ isAuthenticated ? 'admin / lista przepisów' : 'przepisy' }}
                </div>
                <div class="page-header__row">
                    <h1 class="page-title">
                        {{ isAuthenticated ? 'Zarządzaj przepisami' : 'Przepisy' }}
                    </h1>
                    <div class="page-header__count">
                        {{ filtered.length }} {{ pluralPL(filtered.length) }}
                    </div>
                </div>
            </div>

            <!-- Results bar -->
            <div class="results-bar">
                <span class="results-bar__count">{{ filtered.length }} {{ pluralPL(filtered.length) }}</span>
                <label class="results-bar__sort">
                    <span class="results-bar__sort-label">Sortuj</span>
                    <select v-model="sort" class="select" style="width: auto; padding: 7px 28px 7px 10px; font-size: 13px;">
                        <option value="newest">najnowsze</option>
                        <option value="fastest">najszybsze</option>
                        <option value="az">A–Z</option>
                    </select>
                </label>
            </div>

            <div v-if="filtered.length === 0" class="empty-state">
                <span class="empty-state__icon">
                    <svg width="14" height="14" viewBox="0 0 14 14" fill="none">
                        <circle cx="6" cy="6" r="3.5" stroke="currentColor" stroke-width="1.3"/>
                        <path d="M8.5 8.5L12 12" stroke="currentColor" stroke-width="1.3" stroke-linecap="round"/>
                    </svg>
                </span>
                <p style="font-size: 18px; font-weight: 600; letter-spacing: -0.3px;">Nic nie pasuje do filtrów</p>
                <p style="font-size: 13px; color: var(--ink-muted); max-width: 320px; line-height: 1.5;">
                    Spróbuj usunąć część filtrów albo poszukać innego hasła.
                </p>
                <button class="btn" @click="clearFilters">Wyczyść filtry</button>
            </div>

            <div v-else class="recipes-grid">
                <RouterLink v-if="isAuthenticated" to="/recipes/new" class="card-add">
                    <span class="card-add__icon">
                        <svg width="22" height="22" viewBox="0 0 22 22" fill="none">
                            <path d="M11 4v14M4 11h14" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
                        </svg>
                    </span>
                    <div style="text-align: center;">
                        <div style="font-size: 18px; font-weight: 600; letter-spacing: -0.4px;">Dodaj nowy przepis</div>
                        <div style="font-size: 13px; color: var(--ink-muted); margin-top: 6px; line-height: 1.5; max-width: 220px;">
                            Wgraj zdjęcie, składniki i kroki.
                        </div>
                    </div>
                </RouterLink>

                <RecipeCard
                    v-for="recipe in filtered"
                    :key="recipe.id"
                    :recipe
                />
            </div>
        </main>
    </div>
</template>

<style scoped>
.recipes-layout {
    display: grid;
    grid-template-columns: 260px 1fr;
    gap: 40px;
    padding-top: 48px;
    padding-bottom: 80px;
    align-items: start;
}

.recipes-sidebar {
    position: sticky;
    top: 90px;
    display: flex;
    flex-direction: column;
    gap: 24px;
}

/* Search */
.search-icon {
    position: absolute;
    left: 12px;
    top: 50%;
    transform: translateY(-50%);
    color: var(--ink-muted);
    pointer-events: none;
}
.search-input {
    padding-left: 36px;
}

.clear-btn {
    background: transparent;
    border: none;
    color: var(--accent);
    font-family: var(--font-mono);
    font-size: 11px;
    font-weight: 600;
    letter-spacing: 0.3px;
    cursor: pointer;
    text-align: left;
    padding: 0;
    margin-top: -16px;
}

/* Filter group */
.filter-group {}
.filter-group__label {
    font-family: var(--font-mono);
    font-size: 10px;
    color: var(--ink-muted);
    letter-spacing: 0.6px;
    text-transform: uppercase;
    font-weight: 600;
    margin-bottom: 10px;
    padding-bottom: 8px;
    border-bottom: 1px solid var(--rule);
}

.filter-row {
    display: flex;
    align-items: center;
    gap: 8px;
    background: transparent;
    border: none;
    padding: 6px 10px;
    border-radius: var(--r-md);
    cursor: pointer;
    font-family: var(--font-sans);
    font-size: 13px;
    color: var(--ink-dim);
    font-weight: 500;
    text-align: left;
    width: 100%;
    transition: background 0.12s, color 0.12s;
}
.filter-row.is-active {
    background: var(--accent-soft);
    color: var(--accent);
    font-weight: 600;
}
.filter-row__box {
    width: 14px;
    height: 14px;
    flex-shrink: 0;
    border-radius: 3px;
    border: 1.5px solid var(--rule-strong);
    background: transparent;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    transition: all 0.12s;
}
.filter-row.is-active .filter-row__box {
    background: var(--accent);
    border-color: var(--accent);
}
.filter-row.is-active .filter-row__box::after {
    content: '';
    width: 4px;
    height: 8px;
    margin-top: -2px;
    border: solid #fff;
    border-width: 0 2px 2px 0;
    transform: rotate(45deg);
}

/* Category dot in filter */
.filter-row__dot {
    width: 8px;
    height: 8px;
    border-radius: var(--r-pill);
    flex-shrink: 0;
    background: var(--cat, var(--accent));
}

/* Page header */
.page-header {
    padding-bottom: 24px;
}
.page-header__kicker {
    font-family: var(--font-mono);
    font-size: 11px;
    color: var(--accent);
    letter-spacing: 0.4px;
    font-weight: 500;
    margin-bottom: 14px;
}
.page-header__row {
    display: flex;
    align-items: flex-end;
    justify-content: space-between;
    gap: 24px;
    flex-wrap: wrap;
}
.page-title {
    font-size: 72px;
    line-height: 0.95;
    letter-spacing: -2.8px;
    font-weight: 500;
    margin: 0;
}
.page-header__count {
    font-family: var(--font-mono);
    font-size: 12px;
    color: var(--ink-muted);
    letter-spacing: 0.4px;
    padding-bottom: 10px;
}

/* Results bar */
.results-bar {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 18px;
}
.results-bar__count {
    font-family: var(--font-mono);
    font-size: 11px;
    color: var(--ink-muted);
    letter-spacing: 0.3px;
}
.results-bar__sort {
    display: flex;
    align-items: center;
    gap: 8px;
}
.results-bar__sort-label {
    font-family: var(--font-mono);
    font-size: 10px;
    color: var(--ink-muted);
    letter-spacing: 0.5px;
    text-transform: uppercase;
    font-weight: 600;
}

/* Grid */
.recipes-grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 20px;
}

/* Add card icon */
.card-add {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 16px;
    min-height: 380px;
}
.card-add__icon {
    width: 56px;
    height: 56px;
    border-radius: var(--r-pill);
    background: var(--accent);
    color: #fff;
    display: inline-flex;
    align-items: center;
    justify-content: center;
}

/* Empty state */
.empty-state {
    background: var(--bg-alt);
    border: 1px dashed var(--rule-strong);
    border-radius: var(--r-lg);
    padding: 60px 32px;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 14px;
    text-align: center;
}
.empty-state__icon {
    width: 48px;
    height: 48px;
    border-radius: var(--r-pill);
    background: var(--accent-soft);
    display: inline-flex;
    align-items: center;
    justify-content: center;
    color: var(--accent);
}

@media (max-width: 900px) {
    .recipes-layout {
        grid-template-columns: 1fr;
    }
    .recipes-sidebar {
        position: static;
    }
    .recipes-grid {
        grid-template-columns: repeat(2, 1fr);
    }
    .page-title {
        font-size: 48px;
    }
}
</style>
