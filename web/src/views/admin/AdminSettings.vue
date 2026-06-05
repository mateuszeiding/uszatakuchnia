<script lang="ts" setup>
import { ref } from 'vue';

type Tab = 'cats' | 'regions' | 'diets' | 'practical';

const activeTab = ref<Tab>('cats');

const PALETTE = [
    '#3F70A8',
    '#B85A2E',
    '#7A9050',
    '#2F6B3F',
    '#D9A441',
    '#B8946A',
    '#C53728',
    '#A8B83D',
    '#7A4A2E',
    '#D27393',
    '#E89556',
    '#8A4A6D',
    '#2F8A8A',
    '#5C6E33',
];

const cats = ref([
    { id: 'ryby', label: 'ryby', color: '#3F70A8', n: 0 },
    { id: 'mieso', label: 'mięso', color: '#B85A2E', n: 0 },
    { id: 'wege', label: 'wege', color: '#7A9050', n: 0 },
    { id: 'wegan', label: 'wegan', color: '#2F6B3F', n: 0 },
    { id: 'makarony', label: 'makarony', color: '#D9A441', n: 0 },
    { id: 'ryz', label: 'ryż / kasze', color: '#B8946A', n: 0 },
    { id: 'zupy', label: 'zupy', color: '#C53728', n: 0 },
    { id: 'salatki', label: 'sałatki', color: '#A8B83D', n: 0 },
    { id: 'pieczywo', label: 'pieczywo', color: '#7A4A2E', n: 0 },
    { id: 'desery', label: 'desery', color: '#D27393', n: 0 },
    { id: 'sniadania', label: 'śniadania', color: '#E89556', n: 0 },
    { id: 'przekaski', label: 'przekąski', color: '#8A4A6D', n: 0 },
    { id: 'napoje', label: 'napoje', color: '#2F8A8A', n: 0 },
]);
const regions = ref([
    { t: 'polska', n: 0 },
    { t: 'włoska', n: 0 },
    { t: 'francuska', n: 0 },
    { t: 'azjatycka', n: 0 },
    { t: 'skandynawska', n: 0 },
    { t: 'śródziemnomorska', n: 0 },
    { t: 'bliski wschód', n: 0 },
    { t: 'meksykańska', n: 0 },
    { t: 'amerykańska', n: 0 },
]);
const diets = ref([
    { t: 'bez glutenu', n: 0 },
    { t: 'bez laktozy', n: 0 },
    { t: 'bez cukru', n: 0 },
    { t: 'wysokobiałkowe', n: 0 },
    { t: 'niskokaloryczne', n: 0 },
    { t: 'keto', n: 0 },
    { t: 'low fodmap', n: 0 },
]);
const practical = ref([
    { t: '30 min', n: 0 },
    { t: '1 patelnia', n: 0 },
    { t: '1 garnek', n: 0 },
    { t: 'do 5 składników', n: 0 },
    { t: 'można mrozić', n: 0 },
    { t: 'na imprezę', n: 0 },
    { t: 'dla dzieci', n: 0 },
    { t: 'meal prep', n: 0 },
]);

const TABS = [
    { id: 'cats' as Tab, label: 'Kategorie', hint: 'jedna na przepis · nadaje kolor' },
    { id: 'regions' as Tab, label: 'Regiony', hint: 'opcjonalne · widoczne przy nazwie' },
    { id: 'diets' as Tab, label: 'Diety', hint: 'filtr krytyczny · alergie i styl życia' },
    { id: 'practical' as Tab, label: 'Praktyczne', hint: 'wygoda · nie krytyczne' },
];

const counts = () => ({
    cats: cats.value.length,
    regions: regions.value.length,
    diets: diets.value.length,
    practical: practical.value.length,
});

// Categories
const openPicker = ref<string | null>(null);
const deleteCatConfirm = ref<string | null>(null);
const addingCat = ref(false);
const newCatLabel = ref('');
const newCatColor = ref(PALETTE[0] ?? '#3F70A8');

function setCatColor(id: string, color: string) {
    const cat = cats.value.find((c) => c.id === id);
    if (cat) cat.color = color;
    openPicker.value = null;
}
function setCatLabel(id: string, label: string) {
    const cat = cats.value.find((c) => c.id === id);
    if (cat) cat.label = label;
}
function removeCat(id: string) {
    cats.value = cats.value.filter((c) => c.id !== id);
    deleteCatConfirm.value = null;
}
function addCat() {
    if (!newCatLabel.value.trim()) return;
    cats.value.push({
        id: 'c' + Date.now(),
        label: newCatLabel.value.trim(),
        color: newCatColor.value,
        n: 0,
    });
    newCatLabel.value = '';
    addingCat.value = false;
}

// Simple list (regions)
const deleteRegionConfirm = ref<string | null>(null);
const newRegion = ref('');
function addRegion() {
    if (!newRegion.value.trim()) return;
    regions.value.push({ t: newRegion.value.trim(), n: 0 });
    newRegion.value = '';
}
function removeRegion(t: string) {
    regions.value = regions.value.filter((r) => r.t !== t);
    deleteRegionConfirm.value = null;
}

// Tag lists (diets + practical)
const deleteDietConfirm = ref<string | null>(null);
const deletePracticalConfirm = ref<string | null>(null);
const addingDiet = ref(false);
const addingPractical = ref(false);
const newDiet = ref('');
const newPractical = ref('');

function addDiet() {
    if (!newDiet.value.trim()) return;
    diets.value.push({ t: newDiet.value.trim(), n: 0 });
    newDiet.value = '';
    addingDiet.value = false;
}
function removeDiet(t: string) {
    diets.value = diets.value.filter((d) => d.t !== t);
    deleteDietConfirm.value = null;
}
function addPractical() {
    if (!newPractical.value.trim()) return;
    practical.value.push({ t: newPractical.value.trim(), n: 0 });
    newPractical.value = '';
    addingPractical.value = false;
}
function removePractical(t: string) {
    practical.value = practical.value.filter((p) => p.t !== t);
    deletePracticalConfirm.value = null;
}
</script>

<template>
    <div class="container settings-page">
        <div class="settings-header">
            <nav class="settings-breadcrumb">
                <RouterLink
                    to="/admin"
                    class="bc-link"
                >
                    Panel
                </RouterLink>
                <span class="bc-sep">/</span>
                <span class="bc-current">Ustawienia</span>
            </nav>
            <div class="settings-header__row">
                <h1 class="settings-title">Słowniki</h1>
                <div class="settings-meta">
                    {{ cats.length }} kategorii ·
                    {{ regions.length + diets.length + practical.length }} tagów
                </div>
            </div>
        </div>

        <!-- Pill tab bar -->
        <div class="tab-pill-bar">
            <button
                v-for="tab in TABS"
                :key="tab.id"
                class="tab-pill"
                :class="{ 'tab-pill--active': activeTab === tab.id }"
                @click="activeTab = tab.id"
            >
                {{ tab.label }}
                <span class="tab-pill__count">{{ counts()[tab.id] }}</span>
            </button>
        </div>

        <!-- Active section -->
        <div class="settings-panel">
            <div class="settings-panel__head">
                <div class="settings-panel__title">
                    {{ TABS.find((t) => t.id === activeTab)?.label ?? '' }}
                </div>
                <div class="settings-panel__hint">
                    // {{ TABS.find((t) => t.id === activeTab)?.hint ?? '' }}
                </div>
            </div>

            <!-- Categories -->
            <template v-if="activeTab === 'cats'">
                <div class="cat-table-head">
                    <span></span>
                    <span>Nazwa</span>
                    <span class="right">Przepisy</span>
                    <span></span>
                </div>
                <div
                    v-for="cat in cats"
                    :key="cat.id"
                    class="cat-row"
                >
                    <div class="color-cell">
                        <button
                            class="color-swatch"
                            :style="{ background: cat.color }"
                            @click.stop="openPicker = openPicker === cat.id ? null : cat.id"
                        />
                        <div
                            v-if="openPicker === cat.id"
                            class="color-picker"
                            @click.stop
                        >
                            <div class="color-picker__label">Kolor</div>
                            <div class="color-picker__grid">
                                <button
                                    v-for="hex in PALETTE"
                                    :key="hex"
                                    class="color-option"
                                    :class="{ 'color-option--active': cat.color === hex }"
                                    :style="{ background: hex }"
                                    @click="setCatColor(cat.id, hex)"
                                />
                            </div>
                        </div>
                    </div>
                    <input
                        class="cat-label-input"
                        :value="cat.label"
                        @blur="setCatLabel(cat.id, ($event.target as HTMLInputElement).value)"
                    />
                    <span class="cat-count right">{{ cat.n || '—' }}</span>
                    <div class="cat-actions">
                        <template v-if="deleteCatConfirm === cat.id">
                            <button
                                class="btn-danger-confirm"
                                @click="removeCat(cat.id)"
                            >
                                Usuń
                            </button>
                            <button
                                class="btn-cancel-confirm"
                                @click="deleteCatConfirm = null"
                            >
                                Nie
                            </button>
                        </template>
                        <button
                            v-else
                            class="del-btn"
                            @click="deleteCatConfirm = cat.id"
                        >
                            ×
                        </button>
                    </div>
                </div>
                <div
                    v-if="addingCat"
                    class="cat-add-row"
                >
                    <input
                        v-model="newCatLabel"
                        class="input cat-add-input"
                        placeholder="Nazwa kategorii"
                        autofocus
                        @keydown.enter="addCat"
                        @keydown.escape="addingCat = false"
                    />
                    <div class="color-palette">
                        <span class="color-palette__label">Kolor</span>
                        <button
                            v-for="hex in PALETTE"
                            :key="hex"
                            class="color-option"
                            :class="{ 'color-option--active': newCatColor === hex }"
                            :style="{ background: hex }"
                            @click="newCatColor = hex"
                        />
                    </div>
                    <div class="cat-add-actions">
                        <button
                            class="btn btn--accent btn--sm"
                            :disabled="!newCatLabel.trim()"
                            @click="addCat"
                        >
                            Dodaj
                        </button>
                        <button
                            class="btn btn--ghost btn--sm"
                            @click="addingCat = false"
                        >
                            Anuluj
                        </button>
                    </div>
                </div>
                <button
                    v-else
                    class="add-link"
                    @click="addingCat = true"
                >
                    + Dodaj kategorię
                </button>
            </template>

            <!-- Regions -->
            <template v-if="activeTab === 'regions'">
                <div
                    v-for="r in regions"
                    :key="r.t"
                >
                    <div class="list-row">
                        <div class="list-row__info">
                            <span class="list-row__name">{{ r.t }}</span>
                            <span class="list-row__count">{{ r.n }}</span>
                        </div>
                        <button
                            class="del-btn"
                            @click.stop="deleteRegionConfirm = r.t"
                        >
                            ×
                        </button>
                    </div>
                    <div
                        v-if="deleteRegionConfirm === r.t"
                        class="delete-confirm"
                        @click.stop
                    >
                        <span class="delete-confirm__text">
                            Usunąć
                            <b>{{ r.t }}</b>
                            {{ r.n > 0 ? ` · ${r.n} przepisów` : '' }}?
                        </span>
                        <div class="delete-confirm__actions">
                            <button
                                class="btn btn--ghost btn--sm"
                                @click="deleteRegionConfirm = null"
                            >
                                Anuluj
                            </button>
                            <button
                                class="btn btn--danger btn--sm"
                                @click="removeRegion(r.t)"
                            >
                                Usuń
                            </button>
                        </div>
                    </div>
                </div>
                <div class="add-row">
                    <input
                        v-model="newRegion"
                        class="input"
                        placeholder="Dodaj region..."
                        @keydown.enter="addRegion"
                    />
                    <button
                        class="btn btn--accent btn--sm"
                        :disabled="!newRegion.trim()"
                        @click="addRegion"
                    >
                        Dodaj
                    </button>
                </div>
            </template>

            <!-- Diets -->
            <template v-if="activeTab === 'diets'">
                <div class="tag-cloud">
                    <span
                        v-for="d in diets"
                        :key="d.t"
                    >
                        <span
                            v-if="deleteDietConfirm === d.t"
                            class="tag-confirm"
                            @click.stop
                        >
                            Usunąć?
                            <button
                                class="tag-confirm__yes"
                                @click="removeDiet(d.t)"
                            >
                                Tak
                            </button>
                            <button
                                class="tag-confirm__no"
                                @click="deleteDietConfirm = null"
                            >
                                Nie
                            </button>
                        </span>
                        <span
                            v-else
                            class="tag-item"
                        >
                            {{ d.t }}
                            <span class="tag-item__count">{{ d.n }}</span>
                            <button
                                class="tag-item__del"
                                @click.stop="deleteDietConfirm = d.t"
                            >
                                ×
                            </button>
                        </span>
                    </span>
                    <span
                        v-if="addingDiet"
                        class="tag-adding"
                        @click.stop
                    >
                        <input
                            v-model="newDiet"
                            class="tag-adding__input"
                            placeholder="nowy tag"
                            autofocus
                            @keydown.enter="addDiet"
                            @keydown.escape="
                                addingDiet = false;
                                newDiet = '';
                            "
                        />
                        <button
                            class="tag-adding__btn"
                            @click="addDiet"
                        >
                            +
                        </button>
                    </span>
                    <button
                        v-else
                        class="tag-add-btn"
                        @click="addingDiet = true"
                    >
                        + tag
                    </button>
                </div>
            </template>

            <!-- Practical -->
            <template v-if="activeTab === 'practical'">
                <div class="tag-cloud">
                    <span
                        v-for="p in practical"
                        :key="p.t"
                    >
                        <span
                            v-if="deletePracticalConfirm === p.t"
                            class="tag-confirm"
                            @click.stop
                        >
                            Usunąć?
                            <button
                                class="tag-confirm__yes"
                                @click="removePractical(p.t)"
                            >
                                Tak
                            </button>
                            <button
                                class="tag-confirm__no"
                                @click="deletePracticalConfirm = null"
                            >
                                Nie
                            </button>
                        </span>
                        <span
                            v-else
                            class="tag-item"
                        >
                            {{ p.t }}
                            <span class="tag-item__count">{{ p.n }}</span>
                            <button
                                class="tag-item__del"
                                @click.stop="deletePracticalConfirm = p.t"
                            >
                                ×
                            </button>
                        </span>
                    </span>
                    <span
                        v-if="addingPractical"
                        class="tag-adding"
                        @click.stop
                    >
                        <input
                            v-model="newPractical"
                            class="tag-adding__input"
                            placeholder="nowy tag"
                            autofocus
                            @keydown.enter="addPractical"
                            @keydown.escape="
                                addingPractical = false;
                                newPractical = '';
                            "
                        />
                        <button
                            class="tag-adding__btn"
                            @click="addPractical"
                        >
                            +
                        </button>
                    </span>
                    <button
                        v-else
                        class="tag-add-btn"
                        @click="addingPractical = true"
                    >
                        + tag
                    </button>
                </div>
            </template>
        </div>
    </div>
</template>

<style scoped>
.settings-page {
    padding-top: 48px;
    padding-bottom: 80px;
}
.settings-breadcrumb {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 13px;
    font-family: var(--font-mono);
    letter-spacing: 0.3px;
    margin-bottom: 14px;
}
.bc-link {
    color: var(--ink-muted);
    text-decoration: none;
}
.bc-link:hover {
    color: var(--ink);
}
.bc-sep {
    color: var(--ink-faint);
}
.bc-current {
    color: var(--ink);
    font-weight: 600;
}
.settings-header__row {
    display: flex;
    align-items: flex-end;
    justify-content: space-between;
    gap: 24px;
    padding-bottom: 24px;
}
.settings-title {
    font-size: 56px;
    line-height: 0.95;
    letter-spacing: -2.2px;
    font-weight: 500;
    margin: 0;
}
.settings-meta {
    font-family: var(--font-mono);
    font-size: 11px;
    color: var(--ink-muted);
    letter-spacing: 0.3px;
    padding-bottom: 8px;
}

/* Pill tabs */
.tab-pill-bar {
    display: flex;
    gap: 2px;
    background: var(--bg-alt);
    border: 1px solid var(--rule);
    border-radius: var(--r-lg);
    padding: 4px;
    width: fit-content;
    margin-bottom: 28px;
}
.tab-pill {
    display: flex;
    align-items: center;
    gap: 8px;
    background: transparent;
    border: 1px solid transparent;
    border-radius: var(--r-md);
    padding: 8px 16px;
    font-size: 13px;
    font-weight: 500;
    cursor: pointer;
    font-family: var(--font-sans);
    color: var(--ink);
    transition: all 0.15s;
}
.tab-pill--active {
    background: var(--bg);
    border-color: var(--rule-strong);
    font-weight: 600;
    box-shadow: 0 1px 3px rgba(10, 10, 15, 0.08);
}
.tab-pill__count {
    font-family: var(--font-mono);
    font-size: 10px;
    color: var(--ink-faint);
    font-weight: 600;
}
.tab-pill--active .tab-pill__count {
    color: var(--accent);
}

/* Settings panel */
.settings-panel {
    background: var(--bg-alt);
    border: 1px solid var(--rule);
    border-radius: var(--r-xl);
    padding: 28px;
    max-width: 800px;
}
.settings-panel__head {
    margin-bottom: 20px;
    padding-bottom: 16px;
    border-bottom: 1px solid var(--rule);
}
.settings-panel__title {
    font-size: 18px;
    font-weight: 600;
    letter-spacing: -0.4px;
}
.settings-panel__hint {
    font-family: var(--font-mono);
    font-size: 11px;
    color: var(--ink-faint);
    letter-spacing: 0.3px;
    margin-top: 4px;
}

/* Categories table */
.cat-table-head {
    display: grid;
    grid-template-columns: 28px 1fr 80px 80px;
    gap: 12px;
    padding: 6px 0 8px;
    border-bottom: 1px solid var(--rule);
    margin-bottom: 2px;
    font-family: var(--font-mono);
    font-size: 9px;
    color: var(--ink-faint);
    letter-spacing: 0.6px;
    text-transform: uppercase;
    font-weight: 600;
}
.cat-row {
    display: grid;
    grid-template-columns: 28px 1fr 80px 80px;
    gap: 12px;
    align-items: center;
    padding: 8px 0;
    border-bottom: 1px solid var(--rule);
    position: relative;
}
.color-cell {
    position: relative;
}
.color-swatch {
    width: 22px;
    height: 22px;
    border-radius: 5px;
    border: 1px solid rgba(0, 0, 0, 0.15);
    cursor: pointer;
    padding: 0;
}
.color-picker {
    position: absolute;
    top: 26px;
    left: 0;
    z-index: 40;
    background: var(--bg-alt);
    border: 1px solid var(--rule-strong);
    border-radius: var(--r-lg);
    padding: 10px;
    box-shadow: 0 14px 40px rgba(10, 10, 15, 0.18);
    width: 180px;
}
.color-picker__label {
    font-family: var(--font-mono);
    font-size: 9px;
    color: var(--ink-faint);
    letter-spacing: 0.6px;
    text-transform: uppercase;
    margin-bottom: 8px;
}
.color-picker__grid {
    display: grid;
    grid-template-columns: repeat(7, 1fr);
    gap: 5px;
}
.color-option {
    width: 20px;
    height: 20px;
    border-radius: 5px;
    border: 1px solid rgba(0, 0, 0, 0.12);
    cursor: pointer;
    padding: 0;
}
.color-option--active {
    border: 2px solid var(--ink) !important;
}
.color-palette {
    display: flex;
    align-items: center;
    gap: 6px;
    flex-wrap: wrap;
    margin-top: 8px;
}
.color-palette__label {
    font-family: var(--font-mono);
    font-size: 9px;
    color: var(--ink-faint);
    letter-spacing: 0.6px;
    text-transform: uppercase;
    font-weight: 600;
}
.cat-label-input {
    background: transparent;
    border: none;
    outline: none;
    font-size: 13px;
    font-weight: 600;
    color: var(--ink);
    font-family: var(--font-sans);
    padding: 0;
    width: 100%;
}
.cat-count {
    font-family: var(--font-mono);
    font-size: 12px;
    color: var(--ink-faint);
}
.right {
    text-align: right;
}
.cat-actions {
    display: flex;
    gap: 4px;
    justify-content: flex-end;
}
.cat-add-row {
    padding: 12px 0;
    border-bottom: 1px solid var(--rule);
    display: flex;
    flex-direction: column;
    gap: 10px;
}
.cat-add-input {
    border-color: var(--accent) !important;
}
.cat-add-actions {
    display: flex;
    gap: 8px;
}
.add-link {
    display: flex;
    align-items: center;
    gap: 6px;
    background: transparent;
    border: none;
    color: var(--accent);
    font-family: var(--font-mono);
    font-size: 11px;
    font-weight: 600;
    padding: 12px 0;
    cursor: pointer;
    letter-spacing: 0.3px;
}

/* Buttons */
.del-btn {
    width: 24px;
    height: 24px;
    border: 1px solid var(--rule);
    background: transparent;
    border-radius: 5px;
    cursor: pointer;
    color: var(--ink-faint);
    display: inline-flex;
    align-items: center;
    justify-content: center;
    font-size: 14px;
}
.del-btn:hover {
    border-color: #c53728;
    color: #c53728;
}
.btn-danger-confirm {
    background: #c53728;
    color: #fff;
    border: none;
    padding: 3px 8px;
    border-radius: 4px;
    font-size: 11px;
    font-weight: 700;
    cursor: pointer;
    font-family: var(--font-sans);
}
.btn-cancel-confirm {
    background: transparent;
    border: 1px solid var(--rule-strong);
    color: var(--ink-dim);
    padding: 3px 7px;
    border-radius: 4px;
    font-size: 11px;
    font-weight: 600;
    cursor: pointer;
    font-family: var(--font-sans);
}
.btn--danger {
    background: #c53728;
    color: #fff;
    border-color: #c53728;
}

/* List rows (regions) */
.list-row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 9px 0;
    border-bottom: 1px solid var(--rule);
}
.list-row__info {
    display: flex;
    align-items: center;
    gap: 12px;
}
.list-row__name {
    font-size: 14px;
    font-weight: 500;
    color: var(--ink);
}
.list-row__count {
    font-family: var(--font-mono);
    font-size: 10px;
    color: var(--ink-faint);
}
.delete-confirm {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 12px;
    background: color-mix(in srgb, #c53728 7%, transparent);
    border: 1px solid color-mix(in srgb, #c53728 25%, transparent);
    border-radius: var(--r-md);
    padding: 9px 14px;
    margin: 3px 0 6px;
}
.delete-confirm__text {
    font-size: 13px;
    color: var(--ink-dim);
}
.delete-confirm__actions {
    display: flex;
    gap: 7px;
    flex-shrink: 0;
}
.add-row {
    display: flex;
    gap: 8px;
    margin-top: 16px;
}

/* Tag cloud */
.tag-cloud {
    display: flex;
    flex-wrap: wrap;
    gap: 7px;
    align-items: center;
}
.tag-item {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    background: transparent;
    border: 1px solid var(--rule-strong);
    border-radius: var(--r-pill);
    padding: 5px 8px 5px 13px;
    font-size: 13px;
    font-weight: 500;
    color: var(--ink);
}
.tag-item__count {
    font-family: var(--font-mono);
    font-size: 10px;
    color: var(--ink-faint);
    margin-right: 2px;
}
.tag-item__del {
    width: 18px;
    height: 18px;
    border-radius: var(--r-pill);
    border: none;
    background: var(--rule);
    cursor: pointer;
    color: var(--ink-muted);
    display: inline-flex;
    align-items: center;
    justify-content: center;
    padding: 0;
    font-size: 12px;
    line-height: 1;
}
.tag-item__del:hover {
    background: color-mix(in srgb, #c53728 16%, transparent);
    color: #c53728;
}
.tag-confirm {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    background: color-mix(in srgb, #c53728 9%, transparent);
    border: 1px solid color-mix(in srgb, #c53728 38%, transparent);
    border-radius: var(--r-pill);
    padding: 4px 7px 4px 12px;
    font-size: 12px;
    color: #c53728;
}
.tag-confirm__yes {
    background: #c53728;
    color: #fff;
    border: none;
    border-radius: var(--r-pill);
    padding: 2px 9px;
    font-size: 11px;
    font-weight: 700;
    cursor: pointer;
    font-family: var(--font-sans);
}
.tag-confirm__no {
    background: var(--rule);
    color: var(--ink);
    border: none;
    border-radius: var(--r-pill);
    padding: 2px 9px;
    font-size: 11px;
    font-weight: 600;
    cursor: pointer;
    font-family: var(--font-sans);
}
.tag-adding {
    display: inline-flex;
    align-items: center;
    gap: 5px;
    background: var(--bg);
    border: 1px solid var(--accent);
    border-radius: var(--r-pill);
    padding: 3px 5px 3px 12px;
}
.tag-adding__input {
    background: transparent;
    border: none;
    outline: none;
    font-size: 13px;
    font-family: var(--font-sans);
    color: var(--ink);
    width: 90px;
}
.tag-adding__btn {
    width: 20px;
    height: 20px;
    border-radius: var(--r-pill);
    background: var(--accent);
    color: #fff;
    border: none;
    cursor: pointer;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    font-size: 14px;
    line-height: 1;
}
.tag-add-btn {
    display: inline-flex;
    align-items: center;
    gap: 4px;
    background: transparent;
    border: 1px dashed var(--rule-strong);
    color: var(--ink-muted);
    border-radius: var(--r-pill);
    padding: 5px 13px 5px 11px;
    font-size: 13px;
    cursor: pointer;
    font-family: var(--font-sans);
}
</style>
