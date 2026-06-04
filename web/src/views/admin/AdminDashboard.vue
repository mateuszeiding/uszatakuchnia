<script lang="ts" setup>
import { computed, ref } from 'vue';

import { fetchRecipes } from '@/data/api/recipe/fetch';
import type { RecipeBaseDto } from '@/data/dtos/recipe/RecipeDto';

const recipes = ref<RecipeBaseDto[]>([]);
fetchRecipes('list', true).then((v) => (recipes.value = v));

const stats = computed(() => ({
    total: recipes.value.length,
    published: recipes.value.filter((r) => r.status === 'published').length,
    drafts: recipes.value.filter((r) => r.status === 'draft').length,
}));

const recent = computed(() => [...recipes.value].reverse().slice(0, 6));
const drafts = computed(() => recipes.value.filter((r) => r.status === 'draft').slice(0, 4));

const NAV = [
    { id: 'recipes', label: 'Przepisy', to: '/admin/recipes', active: true },
    { id: 'ingredients', label: 'Składniki', to: '/admin/ingredients', active: true },
    { id: 'settings', label: 'Ustawienia', to: '/admin/settings', active: true },
    { id: 'challenges', label: 'Wyzwania', to: null, active: false, soon: true },
    { id: 'wiki', label: 'Wiki', to: null, active: false, soon: true },
];
</script>

<template>
    <div class="admin-layout container">
        <div class="admin-header">
            <div>
                <div class="admin-kicker">// admin / panel</div>
                <h1 class="admin-title">Panel admina</h1>
            </div>
            <RouterLink
                to="/recipes/new"
                class="btn btn--accent admin-new-btn"
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

        <div class="admin-body">
            <!-- Sidebar -->
            <aside class="admin-nav">
                <div class="admin-nav__label">Moduły</div>
                <template
                    v-for="m in NAV"
                    :key="m.id"
                >
                    <RouterLink
                        v-if="m.active && m.to"
                        :to="m.to"
                        class="admin-nav__item"
                        active-class="admin-nav__item--active"
                    >
                        <span class="admin-nav__item-label">{{ m.label }}</span>
                    </RouterLink>
                    <div
                        v-else
                        class="admin-nav__item admin-nav__item--soon"
                    >
                        <span class="admin-nav__item-label">{{ m.label }}</span>
                        <span class="admin-nav__soon-badge">wkrótce</span>
                    </div>
                </template>
            </aside>

            <!-- Main content -->
            <main class="admin-main">
                <!-- Stats -->
                <div class="stats-grid">
                    <div class="stat-card">
                        <div class="stat-card__n">{{ stats.total }}</div>
                        <div class="stat-card__label">przepisy</div>
                    </div>
                    <div class="stat-card">
                        <div class="stat-card__n">{{ stats.published }}</div>
                        <div class="stat-card__label">opublikowane</div>
                        <div class="stat-card__sub">
                            {{
                                stats.total ? Math.round((stats.published / stats.total) * 100) : 0
                            }}% bazy
                        </div>
                    </div>
                    <div class="stat-card stat-card--accent">
                        <div class="stat-card__n">{{ stats.drafts }}</div>
                        <div class="stat-card__label">szkice</div>
                        <div class="stat-card__sub">do dokończenia</div>
                    </div>
                </div>

                <!-- Recently edited -->
                <section class="admin-section">
                    <div class="admin-section__head">
                        <span class="admin-section__title">Ostatnio dodane</span>
                        <RouterLink
                            to="/admin/recipes"
                            class="admin-section__link"
                        >
                            wszystkie →
                        </RouterLink>
                    </div>
                    <div class="recipe-table">
                        <div
                            v-for="r in recent"
                            :key="r.id"
                            class="recipe-table__row"
                        >
                            <span class="recipe-table__id">
                                #{{ String(r.id).padStart(3, '0') }}
                            </span>
                            <span class="recipe-table__name">{{ r.name }}</span>
                            <span
                                class="status-pill"
                                :class="
                                    r.status === 'published'
                                        ? 'status-pill--pub'
                                        : 'status-pill--draft'
                                "
                            >
                                {{ r.status === 'published' ? 'live' : 'szkic' }}
                            </span>
                            <RouterLink
                                :to="{ name: 'recipe-edit', params: { id: r.id } }"
                                class="recipe-table__edit"
                            >
                                edytuj →
                            </RouterLink>
                        </div>
                    </div>
                </section>

                <!-- Drafts -->
                <section
                    v-if="drafts.length"
                    class="admin-section"
                >
                    <div class="admin-section__head">
                        <span class="admin-section__title">Szkice do dokończenia</span>
                        <span class="admin-section__count">{{ stats.drafts }}</span>
                    </div>
                    <div class="drafts-grid">
                        <RouterLink
                            v-for="d in drafts"
                            :key="d.id"
                            :to="{ name: 'recipe-edit', params: { id: d.id } }"
                            class="draft-card card"
                        >
                            <span
                                v-if="d.category"
                                class="badge"
                                :class="`cat-${d.category}`"
                            >
                                {{ d.category }}
                            </span>
                            <div class="draft-card__title">{{ d.name }}</div>
                            <div class="draft-card__action">dokończ →</div>
                        </RouterLink>
                    </div>
                </section>
            </main>
        </div>
    </div>
</template>

<style scoped>
.admin-layout {
    padding-top: 48px;
    padding-bottom: 80px;
}
.admin-header {
    display: flex;
    align-items: flex-end;
    justify-content: space-between;
    gap: 24px;
    padding-bottom: 32px;
}
.admin-kicker {
    font-family: var(--font-mono);
    font-size: 11px;
    color: var(--accent);
    letter-spacing: 0.4px;
    font-weight: 500;
    margin-bottom: 12px;
}
.admin-title {
    font-size: 64px;
    line-height: 0.95;
    letter-spacing: -2.6px;
    font-weight: 500;
    margin: 0;
}
.admin-new-btn {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    text-decoration: none;
}

.admin-body {
    display: grid;
    grid-template-columns: 220px 1fr;
    gap: 40px;
    align-items: start;
}

/* Nav */
.admin-nav {
    position: sticky;
    top: 90px;
    display: flex;
    flex-direction: column;
    gap: 4px;
}
.admin-nav__label {
    font-family: var(--font-mono);
    font-size: 9px;
    color: var(--ink-faint);
    letter-spacing: 0.6px;
    text-transform: uppercase;
    font-weight: 600;
    margin-bottom: 8px;
    padding-bottom: 8px;
    border-bottom: 1px solid var(--rule);
}
.admin-nav__item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 9px 12px;
    border-radius: var(--r-md);
    cursor: pointer;
    font-size: 13px;
    font-weight: 500;
    color: var(--ink);
    text-decoration: none;
    transition: background 0.12s;
    border: 1px solid transparent;
}
.admin-nav__item:hover {
    background: var(--bg-alt);
}
.admin-nav__item--active {
    background: var(--accent-soft) !important;
    border-color: color-mix(in srgb, var(--accent) 28%, transparent);
    color: var(--accent);
    font-weight: 600;
}
.admin-nav__item--soon {
    opacity: 0.4;
    cursor: not-allowed;
}
.admin-nav__soon-badge {
    font-family: var(--font-mono);
    font-size: 8px;
    color: var(--ink-faint);
    letter-spacing: 0.5px;
    text-transform: uppercase;
    font-weight: 600;
    border: 1px solid var(--rule);
    padding: 2px 6px;
    border-radius: var(--r-pill);
}

/* Stats */
.stats-grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 12px;
    margin-bottom: 32px;
}
.stat-card {
    background: var(--bg-alt);
    border: 1px solid var(--rule);
    border-radius: var(--r-lg);
    padding: 18px;
    position: relative;
    overflow: hidden;
}
.stat-card--accent::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 2px;
    background: var(--accent);
}
.stat-card__n {
    font-family: var(--font-mono);
    font-size: 38px;
    font-weight: 500;
    letter-spacing: -1.2px;
    line-height: 1;
}
.stat-card__label {
    font-size: 12px;
    font-weight: 600;
    margin-top: 10px;
}
.stat-card__sub {
    font-family: var(--font-mono);
    font-size: 10px;
    color: var(--ink-muted);
    margin-top: 2px;
    letter-spacing: 0.3px;
}

/* Sections */
.admin-main {
    display: flex;
    flex-direction: column;
    gap: 32px;
}
.admin-section {
    background: var(--bg-alt);
    border: 1px solid var(--rule);
    border-radius: var(--r-lg);
    overflow: hidden;
}
.admin-section__head {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 14px 20px;
    border-bottom: 1px solid var(--rule);
}
.admin-section__title {
    font-size: 14px;
    font-weight: 600;
}
.admin-section__link {
    font-family: var(--font-mono);
    font-size: 11px;
    color: var(--accent);
    text-decoration: none;
    font-weight: 600;
    letter-spacing: 0.3px;
}
.admin-section__count {
    font-family: var(--font-mono);
    font-size: 10px;
    color: var(--accent);
    background: var(--accent-soft);
    padding: 2px 7px;
    border-radius: var(--r-pill);
    font-weight: 600;
}

/* Recipe table */
.recipe-table__row {
    display: grid;
    grid-template-columns: 48px 1fr 60px 80px;
    gap: 14px;
    align-items: center;
    padding: 11px 20px;
    border-bottom: 1px solid var(--rule);
}
.recipe-table__row:last-child {
    border-bottom: none;
}
.recipe-table__id {
    font-family: var(--font-mono);
    font-size: 11px;
    color: var(--ink-faint);
}
.recipe-table__name {
    font-size: 14px;
    font-weight: 600;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
}
.recipe-table__edit {
    font-family: var(--font-mono);
    font-size: 11px;
    color: var(--accent);
    text-decoration: none;
    font-weight: 600;
    text-align: right;
}

/* Status pill */
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

/* Drafts grid */
.drafts-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 14px;
    padding: 16px 20px;
}
.draft-card {
    padding: 14px;
    text-decoration: none;
    color: inherit;
    display: flex;
    flex-direction: column;
    gap: 10px;
}
.draft-card__title {
    font-size: 13px;
    font-weight: 600;
    letter-spacing: -0.2px;
    line-height: 1.3;
}
.draft-card__action {
    font-family: var(--font-mono);
    font-size: 10px;
    color: var(--accent);
    font-weight: 600;
    margin-top: auto;
}

@media (max-width: 900px) {
    .admin-body {
        grid-template-columns: 1fr;
    }
    .admin-title {
        font-size: 40px;
    }
}
</style>
