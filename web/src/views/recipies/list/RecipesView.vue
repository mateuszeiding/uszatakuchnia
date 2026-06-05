<script lang="ts" setup>
import { useAuth0 } from '@auth0/auth0-vue';
import { computed, ref } from 'vue';

import FilterSidebar from './components/FilterSidebar.vue';
import RecipeCard from './components/RecipeCard.vue';

import { fetchRecipes } from '@/data/api/recipe/fetch';
import type { RecipeBaseDto } from '@/data/dtos/recipe/RecipeDto';

const { isAuthenticated } = useAuth0();

const recipes = ref<RecipeBaseDto[]>([]);
const search = ref('');
const activeCategory = ref<string | null>(null);
const maxTime = ref(0);
const activeDiff = ref(0);
const activeDiets = ref<string[]>([]);
const activePractical = ref<string[]>([]);
const catMode = ref<'OR' | 'AND'>('OR');
const dietMode = ref<'OR' | 'AND'>('OR');
const practicalMode = ref<'OR' | 'AND'>('OR');
const hidePrep = ref(false);
const sort = ref<'newest' | 'fastest' | 'az'>('newest');
const tab = ref<'all' | 'pub' | 'draft'>('all');

fetchRecipes('list', true).then((v) => (recipes.value = v));

const filtered = computed(() => {
    let list = recipes.value;
    if (!isAuthenticated.value) {
        list = list.filter((r) => r.status === 'published');
    } else if (tab.value === 'pub') {
        list = list.filter((r) => r.status === 'published');
    } else if (tab.value === 'draft') {
        list = list.filter((r) => r.status === 'draft');
    }
    if (search.value.trim()) {
        const q = search.value.toLowerCase();
        list = list.filter((r) => r.name.toLowerCase().includes(q));
    }
    if (activeCategory.value) {
        list = list.filter((r) => r.category === activeCategory.value);
    }
    if (maxTime.value) {
        list = list.filter((r) => r.timeMinutes != null && r.timeMinutes <= maxTime.value);
    }
    if (activeDiff.value) {
        list = list.filter((r) => r.difficulty === activeDiff.value);
    }
    if (hidePrep.value) {
        list = list.filter((r) => !r.needsPrep);
    }
    if (activeDiets.value.length) {
        list =
            dietMode.value === 'AND'
                ? list.filter((r) => activeDiets.value.every((d) => r.dietTags.includes(d)))
                : list.filter((r) => activeDiets.value.some((d) => r.dietTags.includes(d)));
    }
    if (activePractical.value.length) {
        list =
            practicalMode.value === 'AND'
                ? list.filter((r) =>
                      activePractical.value.every((p) => r.practicalTags.includes(p))
                  )
                : list.filter((r) =>
                      activePractical.value.some((p) => r.practicalTags.includes(p))
                  );
    }
    if (sort.value === 'fastest')
        list = [...list].sort((a, b) => (a.timeMinutes ?? 999) - (b.timeMinutes ?? 999));
    if (sort.value === 'az') list = [...list].sort((a, b) => a.name.localeCompare(b.name, 'pl'));
    return list;
});

const counts = computed(() => ({
    all: recipes.value.length,
    pub: recipes.value.filter((r) => r.status === 'published').length,
    draft: recipes.value.filter((r) => r.status === 'draft').length,
}));

const anyFilter = computed(
    () =>
        !!(
            search.value ||
            activeCategory.value ||
            maxTime.value ||
            activeDiff.value ||
            activeDiets.value.length ||
            activePractical.value.length ||
            hidePrep.value
        )
);

const activeFilterGroups = computed(() => {
    const groups: {
        label: string;
        chips: { key: string; label: string; onRemove: () => void }[];
    }[] = [];
    if (activeCategory.value) {
        const cat = activeCategory.value;
        groups.push({
            label: 'KATEGORIA',
            chips: [{ key: 'cat', label: cat, onRemove: () => (activeCategory.value = null) }],
        });
    }
    if (activeDiets.value.length) {
        groups.push({
            label: 'DIETA',
            chips: activeDiets.value.map((d) => ({
                key: `diet-${d}`,
                label: d,
                onRemove: () => {
                    activeDiets.value = activeDiets.value.filter((x) => x !== d);
                },
            })),
        });
    }
    if (activePractical.value.length) {
        groups.push({
            label: 'PRAKTYCZNE',
            chips: activePractical.value.map((p) => ({
                key: `prac-${p}`,
                label: p,
                onRemove: () => {
                    activePractical.value = activePractical.value.filter((x) => x !== p);
                },
            })),
        });
    }
    if (maxTime.value) {
        const t = maxTime.value;
        groups.push({
            label: 'CZAS',
            chips: [{ key: 'time', label: `do ${t} min`, onRemove: () => (maxTime.value = 0) }],
        });
    }
    if (activeDiff.value) {
        const labels: Record<number, string> = { 1: 'łatwe', 2: 'średnie', 3: 'trudne' };
        const d = activeDiff.value;
        groups.push({
            label: 'TRUDNOŚĆ',
            chips: [{ key: 'diff', label: labels[d] ?? String(d), onRemove: () => (activeDiff.value = 0) }],
        });
    }
    if (search.value) {
        const q = search.value;
        groups.push({
            label: 'SZUKAJ',
            chips: [{ key: 'q', label: `„${q}"`, onRemove: () => (search.value = '') }],
        });
    }
    return groups;
});

function clearFilters() {
    search.value = '';
    activeCategory.value = null;
    maxTime.value = 0;
    activeDiff.value = 0;
    activeDiets.value = [];
    activePractical.value = [];
    hidePrep.value = false;
}

function pluralPL(n: number) {
    if (n === 1) return 'przepis';
    const m10 = n % 10;
    const m100 = n % 100;
    if (m10 >= 2 && m10 <= 4 && (m100 < 10 || m100 >= 20)) return 'przepisy';
    return 'przepisów';
}

function onDeleted(id: number) {
    recipes.value = recipes.value.filter((r) => r.id !== id);
}
</script>

<template>
    <div class="recipes-layout container">
        <FilterSidebar
            v-model:search="search"
            v-model:activeCategory="activeCategory"
            v-model:maxTime="maxTime"
            v-model:activeDiff="activeDiff"
            v-model:activeDiets="activeDiets"
            v-model:activePractical="activePractical"
            v-model:catMode="catMode"
            v-model:dietMode="dietMode"
            v-model:practicalMode="practicalMode"
            v-model:hidePrep="hidePrep"
            :anyFilter="anyFilter"
            @clear="clearFilters"
        />

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

            <div
                v-if="isAuthenticated"
                class="admin-tabs"
            >
                <button
                    v-for="t in [
                        { id: 'all', label: 'Wszystkie', n: counts.all },
                        { id: 'pub', label: 'Opublikowane', n: counts.pub },
                        { id: 'draft', label: 'Szkice', n: counts.draft },
                    ] as const"
                    :key="t.id"
                    class="admin-tab"
                    :class="{ 'admin-tab--active': tab === t.id }"
                    @click="tab = t.id"
                >
                    {{ t.label }}
                    <span class="admin-tab__count">{{ t.n }}</span>
                </button>
            </div>

            <div
                class="active-filter-bar"
                :style="{ visibility: activeFilterGroups.length ? 'visible' : 'hidden' }"
            >
                <template
                    v-for="(group, gi) in activeFilterGroups"
                    :key="group.label"
                >
                    <span
                        v-if="gi > 0"
                        class="active-filter-bar__sep"
                    />
                    <span class="active-filter-bar__group">
                        <span class="active-filter-bar__group-label">{{ group.label }}</span>
                        <span
                            v-for="chip in group.chips"
                            :key="chip.key"
                            class="active-filter-chip"
                        >
                            {{ chip.label }}
                            <button
                                class="active-filter-chip__remove"
                                @click="chip.onRemove"
                            >
                                ×
                            </button>
                        </span>
                    </span>
                </template>
                <button
                    class="active-filter-bar__clear"
                    @click="clearFilters"
                >
                    Wyczyść
                </button>
            </div>

            <div class="results-bar">
                <span class="results-bar__count">
                    {{ filtered.length }} {{ pluralPL(filtered.length) }}
                </span>
                <label class="results-bar__sort">
                    <span class="results-bar__sort-label">Sortuj</span>
                    <select
                        v-model="sort"
                        class="select"
                        style="width: auto; padding: 7px 28px 7px 10px; font-size: 13px"
                    >
                        <option value="newest">najnowsze</option>
                        <option value="fastest">najszybsze</option>
                        <option value="az">A–Z</option>
                    </select>
                </label>
            </div>

            <div
                v-if="filtered.length === 0"
                class="empty-state"
            >
                <span class="empty-state__icon">
                    <svg
                        width="14"
                        height="14"
                        viewBox="0 0 14 14"
                        fill="none"
                    >
                        <circle
                            cx="6"
                            cy="6"
                            r="3.5"
                            stroke="currentColor"
                            stroke-width="1.3"
                        />
                        <path
                            d="M8.5 8.5L12 12"
                            stroke="currentColor"
                            stroke-width="1.3"
                            stroke-linecap="round"
                        />
                    </svg>
                </span>
                <p class="empty-state__title">Nic nie pasuje do filtrów</p>
                <p class="empty-state__desc">
                    Spróbuj usunąć część filtrów albo poszukać innego hasła.
                </p>
                <button
                    class="btn"
                    @click="clearFilters"
                >
                    Wyczyść filtry
                </button>
            </div>

            <div
                v-else
                class="recipes-grid"
            >
                <RouterLink
                    v-if="isAuthenticated"
                    to="/recipes/new"
                    class="card-add"
                >
                    <span class="card-add__icon">
                        <svg
                            width="22"
                            height="22"
                            viewBox="0 0 22 22"
                            fill="none"
                        >
                            <path
                                d="M11 4v14M4 11h14"
                                stroke="currentColor"
                                stroke-width="2"
                                stroke-linecap="round"
                            />
                        </svg>
                    </span>
                    <div>
                        <div class="card-add__title">Dodaj nowy przepis</div>
                        <div class="card-add__sub">Wgraj zdjęcie, składniki i kroki.</div>
                    </div>
                </RouterLink>

                <RecipeCard
                    v-for="recipe in filtered"
                    :key="recipe.id"
                    :recipe
                    @deleted="onDeleted"
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
.admin-tabs {
    display: flex;
    gap: 24px;
    border-bottom: 1px solid var(--rule);
}
.admin-tab {
    background: transparent;
    border: none;
    padding: 10px 0;
    border-bottom: 2px solid transparent;
    color: var(--ink-muted);
    font-size: 14px;
    font-weight: 600;
    cursor: pointer;
    display: inline-flex;
    align-items: center;
    gap: 8px;
    font-family: var(--font-sans);
    margin-bottom: -1px;
}
.admin-tab--active {
    border-bottom-color: var(--accent);
    color: var(--ink);
}
.admin-tab__count {
    background: var(--bg-alt);
    color: var(--ink-muted);
    font-family: var(--font-mono);
    font-size: 10px;
    padding: 2px 7px;
    border-radius: var(--r-pill);
    font-weight: 600;
}
.admin-tab--active .admin-tab__count {
    background: var(--accent-soft);
    color: var(--accent);
}
.active-filter-bar {
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    gap: 0;
    padding: 10px 0 14px;
    border-bottom: 1px solid var(--rule);
    margin-bottom: 8px;
    min-height: 49px;
}
.active-filter-bar__sep {
    width: 1px;
    height: 16px;
    background: var(--rule-strong);
    margin: 0 12px;
    flex-shrink: 0;
}
.active-filter-bar__group {
    display: inline-flex;
    align-items: center;
    gap: 5px;
    flex-shrink: 0;
}
.active-filter-bar__group-label {
    font-family: var(--font-mono);
    font-size: 9px;
    color: var(--ink-faint);
    letter-spacing: 0.5px;
    text-transform: uppercase;
    font-weight: 600;
    margin-right: 1px;
}
.active-filter-chip {
    display: inline-flex;
    align-items: center;
    gap: 3px;
    background: var(--accent-soft);
    border: 1px solid color-mix(in srgb, var(--accent) 28%, transparent);
    border-radius: var(--r-pill);
    padding: 3px 7px 3px 10px;
    font-size: 12px;
    font-weight: 500;
    color: var(--accent);
    white-space: nowrap;
}
.active-filter-chip__remove {
    width: 16px;
    height: 16px;
    border-radius: var(--r-pill);
    border: none;
    cursor: pointer;
    background: transparent;
    color: var(--accent);
    padding: 0;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    font-size: 12px;
    opacity: 0.7;
    line-height: 1;
}
.active-filter-bar__clear {
    margin-left: auto;
    background: transparent;
    border: none;
    color: var(--accent);
    font-family: var(--font-mono);
    font-size: 11px;
    font-weight: 600;
    letter-spacing: 0.3px;
    cursor: pointer;
    padding-left: 14px;
    flex-shrink: 0;
}
.results-bar {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 18px;
    margin-top: 16px;
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
.recipes-grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 20px;
}
.card-add {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 16px;
    min-height: 380px;
    text-align: center;
    text-decoration: none;
    color: inherit;
    border: 1.5px dashed var(--rule-strong);
    border-radius: var(--r-lg);
    background: transparent;
    transition:
        border-color 0.15s,
        background 0.15s;
}
.card-add:hover {
    border-color: var(--accent);
    background: var(--accent-soft);
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
.card-add__title {
    font-size: 18px;
    font-weight: 600;
    letter-spacing: -0.4px;
}
.card-add__sub {
    font-size: 13px;
    color: var(--ink-muted);
    margin-top: 6px;
    line-height: 1.5;
    max-width: 220px;
}
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
.empty-state__title {
    font-size: 18px;
    font-weight: 600;
    letter-spacing: -0.3px;
    margin: 0;
}
.empty-state__desc {
    font-size: 13px;
    color: var(--ink-muted);
    max-width: 320px;
    line-height: 1.5;
    margin: 0;
}
@media (max-width: 900px) {
    .recipes-layout {
        grid-template-columns: 1fr;
    }
    .recipes-grid {
        grid-template-columns: repeat(2, 1fr);
    }
    .page-title {
        font-size: 48px;
    }
}
</style>
