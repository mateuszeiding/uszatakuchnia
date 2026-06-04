<script lang="ts" setup>
import { computed, ref } from 'vue';

import { fetchIngredients } from '@/data/api/ingredients/fetch';
import type { IngredientDto } from '@/data/dtos/ingredients/IngredientDto';

const ingredients = ref<IngredientDto[]>([]);
const search = ref('');
fetchIngredients('list').then((v) => (ingredients.value = v));

const filtered = computed(() => {
    if (!search.value.trim()) return ingredients.value;
    const q = search.value.toLowerCase();
    return ingredients.value.filter((i) => i.name.toLowerCase().includes(q));
});
</script>

<template>
    <div class="container admin-ingredients-page">
        <div class="page-header">
            <div>
                <div class="kicker">// admin / składniki</div>
                <h1 class="page-title">Składniki</h1>
            </div>
        </div>

        <div class="toolbar">
            <input
                v-model="search"
                class="input search-input"
                placeholder="Szukaj składnika..."
            />
            <span class="count-label">{{ filtered.length }} składników</span>
        </div>

        <div class="ingredients-table">
            <div class="ingredients-table__head">
                <span>Nazwa</span>
                <span>Typ</span>
                <span>Alergen</span>
                <span></span>
            </div>
            <div
                v-for="ing in filtered"
                :key="ing.id"
                class="ingredients-table__row"
            >
                <span class="cell-name">
                    <img
                        v-if="ing.image?.url"
                        :src="ing.image.url"
                        :alt="ing.name"
                        class="ing-thumb"
                    />
                    <span
                        v-else
                        class="ing-thumb-placeholder"
                    />
                    {{ ing.name }}
                </span>
                <span class="cell-type">{{ ing.type?.name ?? '—' }}</span>
                <span>
                    <span
                        v-if="ing.isAllergen"
                        class="allergen-badge"
                    >
                        ⚠ alergen
                    </span>
                </span>
                <span class="cell-actions">
                    <button
                        class="action-btn"
                        title="Edytuj"
                        disabled
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
                    </button>
                </span>
            </div>
            <div
                v-if="filtered.length === 0"
                class="empty"
            >
                Brak składników
            </div>
        </div>
    </div>
</template>

<style scoped>
.admin-ingredients-page {
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
    gap: 16px;
    margin-bottom: 16px;
}
.search-input {
    max-width: 280px;
}
.count-label {
    font-family: var(--font-mono);
    font-size: 11px;
    color: var(--ink-muted);
    letter-spacing: 0.3px;
}
.ingredients-table {
    background: var(--bg-alt);
    border: 1px solid var(--rule);
    border-radius: var(--r-lg);
    overflow: hidden;
}
.ingredients-table__head {
    display: grid;
    grid-template-columns: 1fr 160px 100px 40px;
    gap: 14px;
    padding: 10px 20px;
    border-bottom: 1px solid var(--rule);
    font-family: var(--font-mono);
    font-size: 10px;
    color: var(--ink-muted);
    letter-spacing: 0.5px;
    text-transform: uppercase;
}
.ingredients-table__row {
    display: grid;
    grid-template-columns: 1fr 160px 100px 40px;
    gap: 14px;
    align-items: center;
    padding: 10px 20px;
    border-bottom: 1px solid var(--rule);
    transition: background 0.1s;
}
.ingredients-table__row:last-child {
    border-bottom: none;
}
.ingredients-table__row:hover {
    background: var(--bg);
}
.cell-name {
    display: flex;
    align-items: center;
    gap: 10px;
    font-size: 14px;
    font-weight: 600;
}
.ing-thumb {
    width: 28px;
    height: 28px;
    border-radius: var(--r-sm);
    object-fit: cover;
    flex-shrink: 0;
}
.ing-thumb-placeholder {
    width: 28px;
    height: 28px;
    border-radius: var(--r-sm);
    background: var(--rule);
    flex-shrink: 0;
}
.cell-type {
    font-size: 13px;
    color: var(--ink-muted);
}
.allergen-badge {
    font-family: var(--font-mono);
    font-size: 10px;
    font-weight: 600;
    padding: 2px 8px;
    border-radius: var(--r-sm);
    background: color-mix(in srgb, #c53728 12%, transparent);
    color: #c53728;
}
.cell-actions {
    display: flex;
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
    transition: all 0.12s;
}
.action-btn:not(:disabled):hover {
    border-color: var(--accent);
    color: var(--accent);
}
.action-btn:disabled {
    opacity: 0.4;
    cursor: not-allowed;
}
.empty {
    padding: 40px;
    text-align: center;
    font-family: var(--font-mono);
    font-size: 12px;
    color: var(--ink-muted);
}
</style>
