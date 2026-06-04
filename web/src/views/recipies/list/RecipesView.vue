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

fetchRecipes('list').then((v) => (recipes.value = v));

const categories = [
    { key: 'ryby',      label: 'Ryby' },
    { key: 'mieso',     label: 'Mięso' },
    { key: 'wege',      label: 'Wege' },
    { key: 'wegan',     label: 'Wegan' },
    { key: 'makarony',  label: 'Makarony' },
    { key: 'ryz',       label: 'Ryż / Kasze' },
    { key: 'zupy',      label: 'Zupy' },
    { key: 'salatki',   label: 'Sałatki' },
    { key: 'pieczywo',  label: 'Pieczywo' },
    { key: 'desery',    label: 'Desery' },
    { key: 'sniadania', label: 'Śniadania' },
    { key: 'przekaski', label: 'Przekąski' },
    { key: 'napoje',    label: 'Napoje' },
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
    return list;
});

function toggleCategory(key: string) {
    activeCategory.value = activeCategory.value === key ? null : key;
}

function clearFilters() {
    search.value = '';
    activeCategory.value = null;
}
</script>

<template>
    <div class="recipes-layout container">
        <!-- Sidebar -->
        <aside class="recipes-sidebar">
            <div class="field" style="margin-bottom: 24px;">
                <input
                    v-model="search"
                    class="input"
                    type="search"
                    placeholder="Szukaj przepisu..."
                />
            </div>

            <div class="section-head">
                <span class="n">01</span>
                <h2>Kategoria</h2>
                <span class="line" />
            </div>

            <div style="display: flex; flex-direction: column; gap: 4px; margin-bottom: 24px;">
                <button
                    v-for="cat in categories"
                    :key="cat.key"
                    class="cat-filter-btn"
                    :class="{ 'is-active': activeCategory === cat.key }"
                    @click="toggleCategory(cat.key)"
                >
                    <span class="cat-dot badge" :class="`cat-${cat.key}`" style="padding: 2px 6px; font-size: 10px;">
                        {{ cat.label }}
                    </span>
                </button>
            </div>

            <button
                v-if="activeCategory || search"
                class="btn btn--ghost btn--sm"
                style="width: 100%;"
                @click="clearFilters"
            >
                Wyczyść filtry
            </button>
        </aside>

        <!-- Main -->
        <main class="recipes-main">
            <div style="display: flex; align-items: center; justify-content: space-between; margin-bottom: 32px;">
                <div>
                    <h1 style="font-size: var(--text-h2); font-weight: 600; letter-spacing: -0.9px;">
                        {{ isAuthenticated ? 'Zarządzaj przepisami' : 'Przepisy' }}
                    </h1>
                    <p class="muted" style="margin-top: 4px; font-size: var(--text-small);">
                        {{ filtered.length }} {{ filtered.length === 1 ? 'przepis' : 'przepisów' }}
                    </p>
                </div>
                <RouterLink v-if="isAuthenticated" to="/recipes/new" class="btn btn--accent">
                    + Nowy przepis
                </RouterLink>
            </div>

            <div v-if="filtered.length === 0" class="empty">
                <div class="empty-icon">
                    <svg width="20" height="20" viewBox="0 0 20 20" fill="none" stroke="currentColor" stroke-width="1.5">
                        <path d="M10 3v14M3 10h14" stroke-linecap="round"/>
                    </svg>
                </div>
                <p class="t-h3">Brak przepisów</p>
                <p class="muted t-small">Spróbuj innych filtrów lub dodaj nowy przepis.</p>
            </div>

            <div v-else class="recipes-grid">
                <RouterLink v-if="isAuthenticated" to="/recipes/new" class="card-add">
                    <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M12 5v14M5 12h14" stroke-linecap="round"/></svg>
                    <span class="t-small" style="font-weight: 600;">Dodaj nowy przepis</span>
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
    grid-template-columns: 240px 1fr;
    gap: 48px;
    padding-top: 48px;
    padding-bottom: 80px;
    align-items: start;
}

.recipes-sidebar {
    position: sticky;
    top: 80px;
}

.recipes-grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 20px;
}

.cat-filter-btn {
    background: none;
    border: none;
    padding: 0;
    cursor: pointer;
    text-align: left;
    border-radius: var(--r-md);
    transition: opacity var(--dur-fast);
}

.cat-filter-btn:not(.is-active) {
    opacity: 0.65;
}

.cat-filter-btn:hover {
    opacity: 1;
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
}
</style>
