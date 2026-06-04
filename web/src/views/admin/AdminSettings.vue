<script lang="ts" setup>
import { ref } from 'vue';

// Lokalne ustawienia — w przyszłości z API
const activeTab = ref<'categories' | 'regions' | 'diets' | 'practical'>('categories');

const categories = ref([
    { key: 'ryby', label: 'ryby' },
    { key: 'mieso', label: 'mięso' },
    { key: 'wege', label: 'wege' },
    { key: 'wegan', label: 'wegan' },
    { key: 'makarony', label: 'makarony' },
    { key: 'ryz', label: 'ryż / kasze' },
    { key: 'zupy', label: 'zupy' },
    { key: 'salatki', label: 'sałatki' },
    { key: 'pieczywo', label: 'pieczywo' },
    { key: 'desery', label: 'desery' },
    { key: 'sniadania', label: 'śniadania' },
    { key: 'przekaski', label: 'przekąski' },
    { key: 'napoje', label: 'napoje' },
]);

const regions = ref([
    'polska',
    'włoska',
    'francuska',
    'azjatycka',
    'skandynawska',
    'śródziemnomorska',
    'bliski wschód',
    'meksykańska',
    'amerykańska',
]);
const diets = ref([
    'bez glutenu',
    'bez laktozy',
    'bez cukru',
    'wysokobiałkowe',
    'niskokaloryczne',
    'keto',
    'low fodmap',
]);
const practical = ref([
    '30 min',
    '1 patelnia',
    '1 garnek',
    'do 5 składników',
    'można mrozić',
    'na imprezę',
    'dla dzieci',
    'meal prep',
]);

const newRegion = ref('');
const newDiet = ref('');
const newPractical = ref('');

const deleteConfirm = ref<string | null>(null);

function addToList(list: typeof regions, input: typeof newRegion) {
    const v = input.value.trim();
    if (v && !list.value.includes(v)) {
        list.value.push(v);
        input.value = '';
    }
}
function removeFromList(list: typeof regions, item: string) {
    if (deleteConfirm.value === item) {
        list.value = list.value.filter((x) => x !== item);
        deleteConfirm.value = null;
    } else {
        deleteConfirm.value = item;
    }
}

document.addEventListener('click', () => {
    deleteConfirm.value = null;
});
</script>

<template>
    <div class="container admin-settings-page">
        <div class="page-header">
            <div class="kicker">// admin / ustawienia</div>
            <h1 class="page-title">Ustawienia</h1>
        </div>

        <div class="settings-body">
            <!-- Tabs -->
            <div class="settings-tabs">
                <button
                    v-for="t in [
                        { id: 'categories', label: 'Kategorie' },
                        { id: 'regions', label: 'Regiony' },
                        { id: 'diets', label: 'Diety' },
                        { id: 'practical', label: 'Praktyczne' },
                    ]"
                    :key="t.id"
                    class="settings-tab"
                    :class="{ 'settings-tab--active': activeTab === t.id }"
                    @click="activeTab = t.id as any"
                >
                    {{ t.label }}
                </button>
            </div>

            <!-- Categories -->
            <div
                v-if="activeTab === 'categories'"
                class="settings-panel"
            >
                <div class="settings-panel__title">Kategorie przepisów</div>
                <div class="settings-panel__hint">
                    Kolory kategorii są zdefiniowane w Design System i nie są edytowalne tutaj.
                </div>
                <div class="tags-list">
                    <span
                        v-for="cat in categories"
                        :key="cat.key"
                        class="tag-item"
                    >
                        <span
                            class="badge"
                            :class="`cat-${cat.key}`"
                        >
                            {{ cat.label }}
                        </span>
                    </span>
                </div>
            </div>

            <!-- Regions -->
            <div
                v-if="activeTab === 'regions'"
                class="settings-panel"
            >
                <div class="settings-panel__title">Regiony / kuchnie</div>
                <div class="tags-list">
                    <span
                        v-for="r in regions"
                        :key="r"
                        class="tag-item"
                        @click.stop
                    >
                        <template v-if="deleteConfirm === r">
                            <span class="tag-item__confirm">
                                Usuń?
                                <button
                                    class="tag-item__yes"
                                    @click.stop="removeFromList(regions, r)"
                                >
                                    Tak
                                </button>
                                <button
                                    class="tag-item__no"
                                    @click.stop="deleteConfirm = null"
                                >
                                    Nie
                                </button>
                            </span>
                        </template>
                        <template v-else>
                            {{ r }}
                            <button
                                class="tag-item__del"
                                @click.stop="removeFromList(regions, r)"
                            >
                                ×
                            </button>
                        </template>
                    </span>
                </div>
                <div class="add-row">
                    <input
                        v-model="newRegion"
                        class="input add-input"
                        placeholder="Nowy region..."
                        @keydown.enter="addToList(regions, newRegion)"
                    />
                    <button
                        class="btn btn--secondary"
                        @click="addToList(regions, newRegion)"
                    >
                        Dodaj
                    </button>
                </div>
            </div>

            <!-- Diets -->
            <div
                v-if="activeTab === 'diets'"
                class="settings-panel"
            >
                <div class="settings-panel__title">Tagi grupy Diety</div>
                <div class="settings-panel__hint">
                    Tagi diety mają silniejszy styl wizualny — informacja krytyczna dla alergików.
                </div>
                <div class="tags-list">
                    <span
                        v-for="d in diets"
                        :key="d"
                        class="tag-item tag-item--diet"
                        @click.stop
                    >
                        <template v-if="deleteConfirm === d">
                            <span class="tag-item__confirm">
                                Usuń?
                                <button
                                    class="tag-item__yes"
                                    @click.stop="removeFromList(diets, d)"
                                >
                                    Tak
                                </button>
                                <button
                                    class="tag-item__no"
                                    @click.stop="deleteConfirm = null"
                                >
                                    Nie
                                </button>
                            </span>
                        </template>
                        <template v-else>
                            {{ d }}
                            <button
                                class="tag-item__del"
                                @click.stop="removeFromList(diets, d)"
                            >
                                ×
                            </button>
                        </template>
                    </span>
                </div>
                <div class="add-row">
                    <input
                        v-model="newDiet"
                        class="input add-input"
                        placeholder="Nowy tag diety..."
                        @keydown.enter="addToList(diets, newDiet)"
                    />
                    <button
                        class="btn btn--secondary"
                        @click="addToList(diets, newDiet)"
                    >
                        Dodaj
                    </button>
                </div>
            </div>

            <!-- Practical -->
            <div
                v-if="activeTab === 'practical'"
                class="settings-panel"
            >
                <div class="settings-panel__title">Tagi grupy Praktyczne</div>
                <div class="settings-panel__hint">
                    Tagi praktyczne mają lżejszy styl wizualny — informacja o wygodzie.
                </div>
                <div class="tags-list">
                    <span
                        v-for="p in practical"
                        :key="p"
                        class="tag-item tag-item--practical"
                        @click.stop
                    >
                        <template v-if="deleteConfirm === p">
                            <span class="tag-item__confirm">
                                Usuń?
                                <button
                                    class="tag-item__yes"
                                    @click.stop="removeFromList(practical, p)"
                                >
                                    Tak
                                </button>
                                <button
                                    class="tag-item__no"
                                    @click.stop="deleteConfirm = null"
                                >
                                    Nie
                                </button>
                            </span>
                        </template>
                        <template v-else>
                            {{ p }}
                            <button
                                class="tag-item__del"
                                @click.stop="removeFromList(practical, p)"
                            >
                                ×
                            </button>
                        </template>
                    </span>
                </div>
                <div class="add-row">
                    <input
                        v-model="newPractical"
                        class="input add-input"
                        placeholder="Nowy tag praktyczny..."
                        @keydown.enter="addToList(practical, newPractical)"
                    />
                    <button
                        class="btn btn--secondary"
                        @click="addToList(practical, newPractical)"
                    >
                        Dodaj
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
.admin-settings-page {
    padding-top: 48px;
    padding-bottom: 80px;
}
.page-header {
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
.settings-body {
    background: var(--bg-alt);
    border: 1px solid var(--rule);
    border-radius: var(--r-lg);
    overflow: hidden;
}
.settings-tabs {
    display: flex;
    border-bottom: 1px solid var(--rule);
    padding: 0 20px;
}
.settings-tab {
    background: transparent;
    border: none;
    padding: 12px 16px;
    border-bottom: 2px solid transparent;
    color: var(--ink-muted);
    font-size: 13px;
    font-weight: 600;
    cursor: pointer;
    font-family: var(--font-sans);
    margin-bottom: -1px;
    transition: color 0.12s;
}
.settings-tab--active {
    border-bottom-color: var(--accent);
    color: var(--ink);
}
.settings-panel {
    padding: 24px 28px;
    display: flex;
    flex-direction: column;
    gap: 16px;
}
.settings-panel__title {
    font-size: 16px;
    font-weight: 600;
    letter-spacing: -0.3px;
}
.settings-panel__hint {
    font-family: var(--font-mono);
    font-size: 11px;
    color: var(--ink-muted);
    letter-spacing: 0.3px;
}
.tags-list {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
}
.tag-item {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    padding: 5px 10px 5px 12px;
    border-radius: var(--r-pill);
    border: 1.5px solid var(--rule-strong);
    font-size: 13px;
    font-weight: 500;
    color: var(--ink);
    background: transparent;
    position: relative;
}
.tag-item--diet {
    font-family: var(--font-sans);
    border-width: 1.5px;
}
.tag-item--practical {
    font-family: var(--font-mono);
    font-size: 12px;
    font-weight: 400;
    color: var(--ink-muted);
    border-width: 1px;
    border-color: var(--rule);
}
.tag-item__del {
    background: transparent;
    border: none;
    color: var(--ink-faint);
    cursor: pointer;
    font-size: 14px;
    line-height: 1;
    padding: 0;
    width: 16px;
    height: 16px;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    border-radius: var(--r-pill);
    transition: all 0.12s;
}
.tag-item__del:hover {
    background: color-mix(in srgb, #c53728 12%, transparent);
    color: #c53728;
}
.tag-item__confirm {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    font-size: 12px;
    font-family: var(--font-mono);
}
.tag-item__yes,
.tag-item__no {
    background: transparent;
    border: 1px solid var(--rule);
    border-radius: var(--r-sm);
    padding: 2px 8px;
    font-size: 11px;
    cursor: pointer;
    font-family: var(--font-mono);
}
.tag-item__yes {
    color: #c53728;
    border-color: #c53728;
}
.tag-item__no {
    color: var(--ink-muted);
}
.add-row {
    display: flex;
    gap: 10px;
    align-items: center;
    padding-top: 8px;
    border-top: 1px solid var(--rule);
}
.add-input {
    flex: 1;
    max-width: 320px;
}
</style>
