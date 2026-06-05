<script lang="ts" setup>
const props = defineProps<{
    search: string;
    activeCategory: string | null;
    maxTime: number;
    activeDiff: number;
    activeDiets: string[];
    activePractical: string[];
    catMode: 'OR' | 'AND';
    dietMode: 'OR' | 'AND';
    practicalMode: 'OR' | 'AND';
    hidePrep: boolean;
    anyFilter: boolean;
}>();

const emit = defineEmits<{
    'update:search': [v: string];
    'update:activeCategory': [v: string | null];
    'update:maxTime': [v: number];
    'update:activeDiff': [v: number];
    'update:activeDiets': [v: string[]];
    'update:activePractical': [v: string[]];
    'update:catMode': [v: 'OR' | 'AND'];
    'update:dietMode': [v: 'OR' | 'AND'];
    'update:practicalMode': [v: 'OR' | 'AND'];
    'update:hidePrep': [v: boolean];
    clear: [];
}>();

const CATEGORIES = [
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
];

const TIME_OPTIONS = [
    { v: 0, l: 'dowolny' },
    { v: 15, l: 'do 15 min' },
    { v: 30, l: 'do 30 min' },
    { v: 60, l: 'do 1 h' },
    { v: 120, l: 'do 2 h' },
];

const DIFF_OPTIONS = [
    { v: 0, l: 'dowolna' },
    { v: 1, l: 'łatwe' },
    { v: 2, l: 'średnie' },
    { v: 3, l: 'trudne' },
];

const DIETS = [
    'bez glutenu',
    'bez laktozy',
    'bez cukru',
    'wysokobiałkowe',
    'niskokaloryczne',
    'keto',
    'low fodmap',
];
const PRACTICAL = [
    '30 min',
    '1 patelnia',
    '1 garnek',
    'do 5 składników',
    'można mrozić',
    'na imprezę',
    'dla dzieci',
    'meal prep',
];

function toggleCat(key: string) {
    emit('update:activeCategory', props.activeCategory === key ? null : key);
}

function toggleDiet(tag: string) {
    const next = props.activeDiets.includes(tag)
        ? props.activeDiets.filter((d) => d !== tag)
        : [...props.activeDiets, tag];
    emit('update:activeDiets', next);
}

function togglePractical(tag: string) {
    const next = props.activePractical.includes(tag)
        ? props.activePractical.filter((p) => p !== tag)
        : [...props.activePractical, tag];
    emit('update:activePractical', next);
}
</script>

<template>
    <aside class="filter-sidebar">
        <!-- Search — zawsze widoczne -->
        <div class="filter-search">
            <svg
                class="filter-search__icon"
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
            <input
                class="input filter-search__input"
                :value="search"
                placeholder="Szukaj przepisu..."
                @input="emit('update:search', ($event.target as HTMLInputElement).value)"
            />
        </div>

        <!-- Prep switch — zawsze widoczne -->
        <div class="filter-prep-switch">
            <label class="prep-switch">
                <input
                    type="checkbox"
                    class="prep-switch__input"
                    :checked="hidePrep"
                    @change="emit('update:hidePrep', ($event.target as HTMLInputElement).checked)"
                />
                <span class="prep-switch__track">
                    <span class="prep-switch__thumb" />
                </span>
                <span class="prep-switch__label">
                    Ukryj wymagające
                    <br />
                    przygotowania
                </span>
            </label>
        </div>

        <!-- Scrollowalne filtry -->
        <div class="filter-scroll">
            <!-- Kategoria -->
            <div class="filter-group">
                <div class="filter-group__header">
                    <span class="filter-group__label">Kategoria</span>
                    <button
                        class="or-and-toggle"
                        :class="{ 'or-and-toggle--disabled': !activeCategory }"
                        :disabled="!activeCategory"
                        @click="emit('update:catMode', catMode === 'OR' ? 'AND' : 'OR')"
                    >
                        <span :class="{ 'or-and-toggle__opt--active': catMode === 'OR' }">LUB</span>
                        <span :class="{ 'or-and-toggle__opt--active': catMode === 'AND' }">I</span>
                    </button>
                </div>
                <div class="filter-rows">
                    <button
                        v-for="cat in CATEGORIES"
                        :key="cat.key"
                        class="filter-row"
                        :class="{ 'filter-row--active': activeCategory === cat.key }"
                        @click="toggleCat(cat.key)"
                    >
                        <span
                            class="badge"
                            :class="`cat-${cat.key}`"
                        >
                            {{ cat.label }}
                        </span>
                    </button>
                </div>
            </div>

            <!-- Czas -->
            <div class="filter-group">
                <div class="filter-group__header">
                    <span class="filter-group__label">Czas</span>
                </div>
                <div class="filter-rows">
                    <button
                        v-for="opt in TIME_OPTIONS"
                        :key="opt.v"
                        class="filter-row"
                        :class="{ 'filter-row--active': maxTime === opt.v }"
                        @click="emit('update:maxTime', opt.v)"
                    >
                        {{ opt.l }}
                    </button>
                </div>
            </div>

            <!-- Trudność -->
            <div class="filter-group">
                <div class="filter-group__header">
                    <span class="filter-group__label">Trudność</span>
                </div>
                <div class="filter-rows">
                    <button
                        v-for="opt in DIFF_OPTIONS"
                        :key="opt.v"
                        class="filter-row"
                        :class="{ 'filter-row--active': activeDiff === opt.v }"
                        @click="emit('update:activeDiff', opt.v)"
                    >
                        <span
                            v-if="opt.v > 0"
                            class="diff-dots"
                        >
                            <span
                                v-for="i in 3"
                                :key="i"
                                class="diff-dot"
                                :class="{ 'diff-dot--on': i <= opt.v }"
                            />
                        </span>
                        {{ opt.l }}
                    </button>
                </div>
            </div>

            <!-- Diety — silna hierarchia (health-critical) -->
            <div class="filter-group">
                <div class="filter-group__header">
                    <span class="filter-group__label">Diety</span>
                    <button
                        class="or-and-toggle"
                        :class="{ 'or-and-toggle--disabled': activeDiets.length < 2 }"
                        :disabled="activeDiets.length < 2"
                        @click="emit('update:dietMode', dietMode === 'OR' ? 'AND' : 'OR')"
                    >
                        <span :class="{ 'or-and-toggle__opt--active': dietMode === 'OR' }">
                            LUB
                        </span>
                        <span :class="{ 'or-and-toggle__opt--active': dietMode === 'AND' }">I</span>
                    </button>
                </div>
                <div class="filter-chips">
                    <button
                        v-for="tag in DIETS"
                        :key="tag"
                        class="chip-diet"
                        :class="{ 'chip-diet--active': activeDiets.includes(tag) }"
                        @click="toggleDiet(tag)"
                    >
                        {{ tag }}
                    </button>
                </div>
            </div>

            <!-- Praktyczne — lżejsza hierarchia -->
            <div class="filter-group">
                <div class="filter-group__header">
                    <span class="filter-group__label">Praktyczne</span>
                    <button
                        class="or-and-toggle"
                        :class="{ 'or-and-toggle--disabled': activePractical.length < 2 }"
                        :disabled="activePractical.length < 2"
                        @click="emit('update:practicalMode', practicalMode === 'OR' ? 'AND' : 'OR')"
                    >
                        <span :class="{ 'or-and-toggle__opt--active': practicalMode === 'OR' }">
                            LUB
                        </span>
                        <span :class="{ 'or-and-toggle__opt--active': practicalMode === 'AND' }">
                            I
                        </span>
                    </button>
                </div>
                <div class="filter-chips">
                    <button
                        v-for="tag in PRACTICAL"
                        :key="tag"
                        class="chip-practical"
                        :class="{ 'chip-practical--active': activePractical.includes(tag) }"
                        @click="togglePractical(tag)"
                    >
                        {{ tag }}
                    </button>
                </div>
            </div>
        </div>

        <!-- Wyczyść — zawsze widoczne na dole -->
        <div class="filter-footer">
            <button
                class="btn filter-clear"
                :class="{ 'filter-clear--disabled': !anyFilter }"
                :disabled="!anyFilter"
                @click="emit('clear')"
            >
                <svg
                    width="11"
                    height="11"
                    viewBox="0 0 11 11"
                    fill="none"
                >
                    <path
                        d="M1.5 1.5l8 8M9.5 1.5l-8 8"
                        stroke="currentColor"
                        stroke-width="1.5"
                        stroke-linecap="round"
                    />
                </svg>
                Wyczyść filtry
            </button>
        </div>
    </aside>
</template>

<style scoped>
.filter-sidebar {
    position: sticky;
    top: 90px;
    align-self: start;
    max-height: calc(100vh - 106px);
    display: flex;
    flex-direction: column;
    gap: 0;
}

.filter-search {
    position: relative;
    flex-shrink: 0;
    padding-bottom: 14px;
}
.filter-search__icon {
    position: absolute;
    left: 12px;
    top: 11px;
    color: var(--ink-muted);
}
.filter-search__input {
    padding-left: 36px;
    width: 100%;
}

.filter-prep-switch {
    flex-shrink: 0;
    padding-bottom: 14px;
    border-bottom: 1px solid var(--rule);
}
.prep-switch {
    display: flex;
    align-items: center;
    gap: 10px;
    cursor: pointer;
    user-select: none;
}
.prep-switch__input {
    position: absolute;
    opacity: 0;
    width: 0;
    height: 0;
}
.prep-switch__track {
    width: 32px;
    height: 18px;
    border-radius: var(--r-pill);
    flex-shrink: 0;
    background: var(--rule-strong);
    position: relative;
    transition: background 0.18s;
    display: inline-block;
}
.prep-switch__input:checked ~ .prep-switch__track {
    background: var(--accent);
}
.prep-switch__thumb {
    position: absolute;
    top: 2px;
    left: 2px;
    width: 14px;
    height: 14px;
    border-radius: var(--r-pill);
    background: #fff;
    transition: left 0.18s;
    box-shadow: 0 1px 2px rgba(10, 10, 15, 0.18);
}
.prep-switch__input:checked ~ .prep-switch__track .prep-switch__thumb {
    left: 16px;
}
.prep-switch__label {
    font-size: 13px;
    font-weight: 500;
    line-height: 1.4;
    color: var(--ink);
}

.filter-scroll {
    flex: 1;
    overflow-y: auto;
    overflow-x: hidden;
    display: flex;
    flex-direction: column;
    gap: 20px;
    scrollbar-width: thin;
    scrollbar-color: var(--rule-strong) transparent;
    padding-right: 2px;
    padding-top: 16px;
}

.filter-group__header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 8px;
    padding-bottom: 7px;
    border-bottom: 1px solid var(--rule);
}
.filter-group__label {
    font-family: var(--font-mono);
    font-size: 10px;
    color: var(--ink-muted);
    letter-spacing: 0.6px;
    text-transform: uppercase;
    font-weight: 600;
}

.or-and-toggle {
    display: flex;
    background: var(--bg);
    border: 1px solid var(--rule-strong);
    border-radius: var(--r-pill);
    padding: 2px;
    flex-shrink: 0;
    cursor: pointer;
    transition: opacity 0.15s;
}
.or-and-toggle--disabled {
    opacity: 0.38;
    cursor: default;
    border-color: var(--rule);
}
.or-and-toggle span {
    border-radius: var(--r-pill);
    padding: 2px 9px;
    font-size: 9px;
    font-family: var(--font-mono);
    font-weight: 700;
    letter-spacing: 0.5px;
    color: var(--ink-faint);
    background: transparent;
    line-height: 1.5;
}
.or-and-toggle__opt--active {
    background: var(--ink) !important;
    color: var(--bg) !important;
}

.filter-rows {
    display: flex;
    flex-direction: column;
    gap: 3px;
}
.filter-row {
    display: flex;
    align-items: center;
    gap: 8px;
    background: transparent;
    border: none;
    padding: 5px 10px;
    border-radius: var(--r-md);
    cursor: pointer;
    font-family: var(--font-sans);
    font-size: 13px;
    color: var(--ink-dim);
    font-weight: 500;
    text-align: left;
    width: 100%;
    transition:
        background 0.12s,
        color 0.12s;
}
.filter-row--active {
    background: var(--accent-soft);
    color: var(--accent);
    font-weight: 600;
}
.filter-row__check {
    width: 14px;
    height: 14px;
    flex-shrink: 0;
    border-radius: var(--r-sm);
    border: 1.5px solid var(--rule-strong);
    background: transparent;
    display: inline-flex;
    align-items: center;
    justify-content: center;
}
.filter-row__check--radio {
    border-radius: var(--r-pill);
}
.filter-row__radio-dot {
    width: 6px;
    height: 6px;
    border-radius: var(--r-pill);
    background: var(--accent);
}
.filter-row--active .filter-row__check--radio {
    border-color: var(--accent);
}

.diff-dots {
    display: inline-flex;
    gap: 3px;
}
.diff-dot {
    width: 5px;
    height: 5px;
    border-radius: var(--r-pill);
    border: 1px solid var(--rule-strong);
    background: transparent;
}
.diff-dot--on {
    background: var(--accent);
    border-color: var(--accent);
}

.filter-chips {
    display: flex;
    flex-wrap: wrap;
    gap: 6px;
}

/* Diet chips — strong visual weight (health-critical) */
.chip-diet {
    display: inline-flex;
    align-items: center;
    padding: 4px 11px;
    border-radius: var(--r-pill);
    font-size: 12px;
    font-weight: 500;
    cursor: pointer;
    font-family: var(--font-sans);
    background: transparent;
    color: var(--ink);
    border: 1.5px solid var(--rule-strong);
    transition: all 0.12s;
}
.chip-diet--active {
    background: var(--accent);
    color: #fff;
    border-color: var(--accent);
    font-weight: 600;
}

/* Practical chips — lighter visual weight */
.chip-practical {
    display: inline-flex;
    align-items: center;
    padding: 3px 9px;
    border-radius: var(--r-pill);
    font-size: 11px;
    font-weight: 400;
    cursor: pointer;
    font-family: var(--font-mono);
    letter-spacing: 0.2px;
    background: transparent;
    color: var(--ink-muted);
    border: 1px solid var(--rule);
    transition: all 0.12s;
}
.chip-practical--active {
    background: var(--ink);
    color: var(--bg);
    border-color: var(--ink);
    font-weight: 600;
}

.filter-footer {
    flex-shrink: 0;
    padding-top: 12px;
    border-top: 1px solid var(--rule);
}
.filter-clear {
    width: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 6px;
    background: transparent;
    border: 1px solid var(--rule-strong);
    color: var(--ink);
    font-weight: 600;
    transition: all 0.15s;
}
.filter-clear--disabled {
    opacity: 0.45;
    border-color: var(--rule);
    color: var(--ink-faint);
    cursor: default;
}
</style>
