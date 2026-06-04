<script lang="ts" setup>
defineProps<{
    anyFilter: boolean;
}>();

const search = defineModel<string>('search', { required: true });
const activeCategory = defineModel<string | null>('activeCategory', { required: true });
const maxTime = defineModel<number>('maxTime', { required: true });
const activeDiff = defineModel<number>('activeDiff', { required: true });

const emit = defineEmits<{ clear: [] }>();

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

const timeOptions = [
    { value: 0,   label: 'dowolny' },
    { value: 15,  label: 'do 15 min' },
    { value: 30,  label: 'do 30 min' },
    { value: 60,  label: 'do 1 h' },
    { value: 120, label: 'do 2 h' },
];

const diffOptions = [
    { value: 0, label: 'dowolna',  dots: 0 },
    { value: 1, label: 'łatwe',    dots: 1 },
    { value: 2, label: 'średnie',  dots: 2 },
    { value: 3, label: 'trudne',   dots: 3 },
];

function toggleCategory(key: string) {
    activeCategory.value = activeCategory.value === key ? null : key;
}
</script>

<template>
    <aside class="sidebar">
        <!-- Search -->
        <div class="search-wrap">
            <svg class="search-icon" width="14" height="14" viewBox="0 0 14 14" fill="none">
                <circle cx="6" cy="6" r="3.5" stroke="currentColor" stroke-width="1.3"/>
                <path d="M8.5 8.5L12 12" stroke="currentColor" stroke-width="1.3" stroke-linecap="round"/>
            </svg>
            <label for="recipes-search" class="visually-hidden">Szukaj przepisu</label>
            <input
                id="recipes-search"
                v-model="search"
                class="input search-input"
                type="search"
                placeholder="Szukaj przepisu..."
            />
        </div>

        <!-- Clear -->
        <button v-if="anyFilter" class="clear-btn" @click="emit('clear')">
            ← wyczyść wszystkie filtry
        </button>

        <!-- Kategoria -->
        <div class="filter-group">
            <div class="filter-group__label">Kategoria</div>
            <div class="filter-rows">
                <button
                    v-for="cat in categories"
                    :key="cat.key"
                    class="filter-row"
                    :class="{ 'is-active': activeCategory === cat.key }"
                    @click="toggleCategory(cat.key)"
                >
                    <span class="filter-row__box" />
                    <span class="filter-dot" :class="`cat-${cat.key}`" />
                    {{ cat.label }}
                </button>
            </div>
        </div>

        <!-- Czas -->
        <div class="filter-group">
            <div class="filter-group__label">Czas</div>
            <div class="filter-rows">
                <button
                    v-for="opt in timeOptions"
                    :key="opt.value"
                    class="filter-row"
                    :class="{ 'is-active': maxTime === opt.value }"
                    @click="maxTime = opt.value"
                >
                    <span class="filter-row__box filter-row__box--radio" />
                    {{ opt.label }}
                </button>
            </div>
        </div>

        <!-- Trudność -->
        <div class="filter-group">
            <div class="filter-group__label">Trudność</div>
            <div class="filter-rows">
                <button
                    v-for="opt in diffOptions"
                    :key="opt.value"
                    class="filter-row"
                    :class="{ 'is-active': activeDiff === opt.value }"
                    @click="activeDiff = opt.value"
                >
                    <span class="filter-row__box filter-row__box--radio" />
                    <span v-if="opt.dots" class="diff-dots">
                        <span
                            v-for="i in 3"
                            :key="i"
                            class="diff-dot"
                            :class="{ 'is-on': i <= opt.dots }"
                        />
                    </span>
                    {{ opt.label }}
                </button>
            </div>
        </div>
    </aside>
</template>

<style scoped>
.sidebar {
    position: sticky;
    top: 90px;
    display: flex;
    flex-direction: column;
    gap: 24px;
}

/* Search */
.search-wrap { position: relative; }
.search-icon {
    position: absolute;
    left: 12px;
    top: 50%;
    transform: translateY(-50%);
    color: var(--ink-muted);
    pointer-events: none;
}
.search-input { padding-left: 36px; }

/* Clear */
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
.filter-rows {
    display: flex;
    flex-direction: column;
    gap: 4px;
}

/* Filter row */
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

/* Checkbox box */
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

/* Radio variant */
.filter-row__box--radio { border-radius: var(--r-pill); }
.filter-row.is-active .filter-row__box--radio {
    background: transparent;
    border-color: var(--accent);
}
.filter-row.is-active .filter-row__box--radio::after {
    content: '';
    width: 6px;
    height: 6px;
    border-radius: var(--r-pill);
    background: var(--accent);
    border: none;
    transform: none;
    margin: 0;
}

/* Category dot */
.filter-dot {
    width: 8px;
    height: 8px;
    border-radius: var(--r-pill);
    flex-shrink: 0;
    background: var(--cat, var(--accent));
}

/* Difficulty dots */
.diff-dots { display: inline-flex; gap: 3px; }
.diff-dot {
    width: 5px;
    height: 5px;
    border-radius: var(--r-pill);
    border: 1px solid var(--rule-strong);
    background: transparent;
}
.diff-dot.is-on { background: var(--accent); border-color: var(--accent); }

/* Accessibility */
.visually-hidden {
    position: absolute;
    width: 1px;
    height: 1px;
    overflow: hidden;
    clip: rect(0 0 0 0);
    white-space: nowrap;
}

@media (max-width: 900px) {
    .sidebar { position: static; }
}
</style>
