<script lang="ts" setup>
import { useQueryClient } from '@tanstack/vue-query';
import { computed, onMounted, ref, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';

import ActiveFilters from './components/ActiveFilters.vue';
import FilterSheet from './components/FilterSheet.vue';
import FiltersPanel from './components/FiltersPanel.vue';
import RecipeCard from './components/RecipeCard.vue';
import RecipeCardSkeleton from './components/RecipeCardSkeleton.vue';
import {
    activeCount,
    applyFilters,
    emptyFilters,
    type IListFilters,
    PAGE_SIZE,
    readQuery,
    recipesLabel,
    toQuery,
} from './listFilters';

import { fetchRecipesWithClient } from '@/data/api/recipe/fetch';
import { fetchDictionariesWithClient } from '@/data/api/taxonomy/fetch';
import type { RecipeBaseDto } from '@/data/dtos/recipe/RecipeDto';
import type { ITaxonomyAxis } from '@/data/dtos/taxonomy/DictionariesDto';

const route = useRoute();
const router = useRouter();
const qc = useQueryClient();

const recipes = ref<RecipeBaseDto[]>([]);
const axes = ref<ITaxonomyAxis[]>([]);
const loading = ref(true);
const filters = ref<IListFilters>(emptyFilters());
const shown = ref(PAGE_SIZE);
const sheetOpen = ref(false);

onMounted(async () => {
    const [list, dicts] = await Promise.all([
        fetchRecipesWithClient(qc, 'list'),
        fetchDictionariesWithClient(qc),
    ]);
    recipes.value = list;
    axes.value = dicts.axes;
    filters.value = readQuery(route.query, dicts.axes);
    loading.value = false;
});

// Stan filtrów → URL (replace: bez zaśmiecania historii), nowe filtry = od pierwszej strony.
watch(
    filters,
    (f) => {
        if (loading.value) return;
        router.replace({ query: toQuery(f) });
        shown.value = PAGE_SIZE;
    },
    { deep: true }
);

const results = computed(() => applyFilters(recipes.value, filters.value, axes.value));
const nActive = computed(() => activeCount(filters.value));
const clear = () => (filters.value = emptyFilters(filters.value.sort));
</script>

<template>
    <div class="container list-page">
        <h1 class="list-page__title">Przepisy</h1>
        <div class="list-layout">
            <aside
                class="list-aside"
                aria-label="Filtry"
            >
                <div class="search">
                    <svg
                        class="search__icon"
                        width="15"
                        height="15"
                        viewBox="0 0 14 14"
                        aria-hidden="true"
                    >
                        <circle
                            cx="6"
                            cy="6"
                            r="3.8"
                            stroke="currentColor"
                            stroke-width="1.3"
                        />
                        <path
                            d="M8.8 8.8L12 12"
                            stroke="currentColor"
                            stroke-width="1.3"
                            stroke-linecap="round"
                        />
                    </svg>
                    <input
                        v-model="filters.q"
                        type="search"
                        class="search__input"
                        aria-label="Szukaj przepisu lub składnika"
                        placeholder="Przepis albo składnik…"
                    />
                </div>
                <div class="list-aside__scroll">
                    <FiltersPanel
                        v-model="filters"
                        :axes="axes"
                    />
                </div>
            </aside>

            <div>
                <div class="mobile-bar">
                    <div class="search">
                        <svg
                            class="search__icon"
                            width="15"
                            height="15"
                            viewBox="0 0 14 14"
                            aria-hidden="true"
                        >
                            <circle
                                cx="6"
                                cy="6"
                                r="3.8"
                                stroke="currentColor"
                                stroke-width="1.3"
                            />
                            <path
                                d="M8.8 8.8L12 12"
                                stroke="currentColor"
                                stroke-width="1.3"
                                stroke-linecap="round"
                            />
                        </svg>
                        <input
                            v-model="filters.q"
                            type="search"
                            class="search__input"
                            aria-label="Szukaj przepisu lub składnika"
                            placeholder="Przepis albo składnik…"
                        />
                    </div>
                    <div class="mobile-bar__row">
                        <button
                            type="button"
                            class="mobile-bar__filters"
                            @click="sheetOpen = true"
                        >
                            Filtry
                            <span
                                v-if="nActive"
                                class="mobile-bar__count"
                            >
                                {{ nActive }}
                            </span>
                        </button>
                        <select
                            v-model="filters.sort"
                            class="select sort-select"
                            aria-label="Sortowanie"
                        >
                            <option value="newest">najnowsze</option>
                            <option value="fastest">najszybsze</option>
                            <option value="az">A–Z</option>
                        </select>
                    </div>
                </div>

                <ActiveFilters
                    v-model="filters"
                    :axes="axes"
                />

                <div class="results-bar">
                    <span role="status">
                        {{
                            loading
                                ? 'Wczytywanie…'
                                : `${results.length} ${recipesLabel(results.length)}`
                        }}
                    </span>
                    <label class="results-bar__sort">
                        Sortuj
                        <select
                            v-model="filters.sort"
                            class="select sort-select"
                        >
                            <option value="newest">najnowsze</option>
                            <option value="fastest">najszybsze</option>
                            <option value="az">A–Z</option>
                        </select>
                    </label>
                </div>

                <ul
                    v-if="loading"
                    class="grid"
                    aria-busy="true"
                >
                    <li
                        v-for="n in 6"
                        :key="n"
                    >
                        <RecipeCardSkeleton />
                    </li>
                </ul>
                <div
                    v-else-if="!results.length"
                    class="empty"
                >
                    <div class="empty__title">Nic nie pasuje</div>
                    <p class="empty__text">
                        Usuń część filtrów albo wpisz inne hasło. Szukamy też po składnikach.
                    </p>
                    <button
                        type="button"
                        class="btn btn--secondary"
                        @click="clear"
                    >
                        Wyczyść filtry
                    </button>
                </div>
                <template v-else>
                    <ul class="grid">
                        <li
                            v-for="r in results.slice(0, shown)"
                            :key="r.recipe.id"
                        >
                            <RecipeCard
                                :recipe="r.recipe"
                                :hits="r.hits"
                            />
                        </li>
                    </ul>
                    <div
                        v-if="shown < results.length"
                        class="more"
                    >
                        <button
                            type="button"
                            class="btn btn--secondary more__btn"
                            @click="shown += PAGE_SIZE"
                        >
                            Pokaż więcej
                        </button>
                        <span class="more__count">{{ shown }} z {{ results.length }}</span>
                    </div>
                </template>
            </div>
        </div>

        <FilterSheet
            v-if="sheetOpen"
            :count="results.length"
            :active-count="nActive"
            @close="sheetOpen = false"
            @clear="clear"
        >
            <FiltersPanel
                v-model="filters"
                :axes="axes"
            />
        </FilterSheet>
    </div>
</template>

<style scoped>
.list-page {
    padding-bottom: 80px;
}
.list-page__title {
    font-size: 72px;
    line-height: 0.95;
    letter-spacing: -2.8px;
    font-weight: 500;
    margin: 56px 0 32px;
}
.list-layout {
    display: grid;
    grid-template-columns: 260px minmax(0, 1fr);
    gap: 40px;
}
.list-aside {
    position: sticky;
    top: 84px;
    align-self: start;
    max-height: calc(100vh - 100px);
    display: flex;
    flex-direction: column;
    gap: 16px;
}
.list-aside__scroll {
    overflow-y: auto;
    padding: 0 8px 16px 0;
}
.search {
    position: relative;
}
.search__icon {
    position: absolute;
    left: 13px;
    top: 50%;
    transform: translateY(-50%);
    color: var(--ink-muted);
}
.search__input {
    width: 100%;
    height: 44px;
    background: var(--bg-alt);
    border: 1px solid var(--rule-strong);
    border-radius: var(--r-lg);
    padding: 0 12px 0 38px;
    font-size: 15px;
    color: var(--ink);
}
.search__input:focus {
    outline: none;
    border-color: var(--accent);
    box-shadow: var(--focus-ring);
}
.mobile-bar {
    display: none;
}
.results-bar {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin: 4px 0 18px;
    font-size: 14px;
    color: var(--ink-muted);
}
.results-bar__sort {
    display: flex;
    align-items: center;
    gap: 8px;
}
.sort-select {
    height: 38px;
    padding-top: 0;
    padding-bottom: 0;
    width: auto;
    background-color: var(--bg-alt);
    border-color: var(--rule-strong);
}
.grid {
    list-style: none;
    padding: 0;
    margin: 0;
    display: grid;
    grid-template-columns: repeat(3, minmax(0, 1fr));
    gap: 20px;
}
.empty {
    background: var(--bg-alt);
    border: 1px dashed var(--rule-strong);
    border-radius: var(--r-lg);
    padding: 56px 28px;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 12px;
    text-align: center;
}
.empty__title {
    font-size: 19px;
    font-weight: 600;
}
.empty__text {
    font-size: 14px;
    color: var(--ink-muted);
    max-width: 340px;
    line-height: 1.5;
}
.more {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 8px;
    margin-top: 32px;
}
.more__btn {
    padding: 12px 24px;
}
.more__count {
    font-size: 13px;
    color: var(--ink-muted);
}

@media (max-width: 1180px) {
    .grid {
        grid-template-columns: repeat(2, minmax(0, 1fr));
    }
}

@media (max-width: 900px) {
    .list-page__title {
        font-size: 44px;
        letter-spacing: -1.6px;
        margin: 28px 0 20px;
    }
    .list-layout {
        grid-template-columns: minmax(0, 1fr);
    }
    .list-aside,
    .results-bar__sort {
        display: none;
    }
    .mobile-bar {
        display: flex;
        flex-direction: column;
        gap: 10px;
        margin-bottom: 12px;
    }
    .mobile-bar__row {
        display: flex;
        gap: 8px;
    }
    .mobile-bar__filters {
        flex: 1;
        height: 44px;
        border-radius: var(--r-lg);
        border: 1px solid var(--rule-strong);
        background: var(--bg-alt);
        color: var(--ink);
        font-size: 15px;
        font-weight: 600;
        cursor: pointer;
        display: inline-flex;
        align-items: center;
        justify-content: center;
        gap: 8px;
    }
    .mobile-bar__count {
        background: var(--accent);
        color: var(--on-accent);
        font-size: 12px;
        border-radius: var(--r-pill);
        padding: 1px 8px;
    }
    .mobile-bar .sort-select {
        flex: 1;
        height: 44px;
    }
    .results-bar {
        margin: 4px 0 14px;
    }
    .grid {
        gap: 16px;
    }
}

@media (max-width: 560px) {
    .grid {
        grid-template-columns: minmax(0, 1fr);
    }
}
</style>
