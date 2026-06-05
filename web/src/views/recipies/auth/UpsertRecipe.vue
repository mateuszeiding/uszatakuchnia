<script lang="ts" setup>
import { useQueryClient } from '@tanstack/vue-query';
import { onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';

import { fetchIngredients } from '@/data/api/ingredients/fetch';
import { createRecipe, fetchRecipesWithClient, updateRecipe } from '@/data/api/recipe/fetch';
import type { IngredientDto } from '@/data/dtos/ingredients/IngredientDto';
import type { IUpsertRecipeRequest } from '@/data/dtos/recipe/RecipeDto';

const props = defineProps<{ id?: number }>();
const router = useRouter();
const qc = useQueryClient();

const saving = ref(false);
const ingredients = ref<IngredientDto[]>([]);

const name = ref('');
const tagline = ref('');
const description = ref('');
const servings = ref(2);
const timeMinutes = ref<number | null>(null);
const difficulty = ref<number | null>(null);
const kcalPerServing = ref<number | null>(null);
const category = ref<string | null>(null);
const region = ref<string | null>(null);
const tip = ref('');

const diets = ref<Set<string>>(new Set());
const practical = ref<Set<string>>(new Set());
const needsPrep = ref(false);
const status = ref<'draft' | 'published'>('draft');

const DIETS = [
    'bez glutenu',
    'bez laktozy',
    'bez cukru',
    'wysokobiałkowe',
    'niskokaloryczne',
    'keto',
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
const REGIONS = [
    '—',
    'polska',
    'włoska',
    'francuska',
    'azjatycka',
    'skandynawska',
    'śródziemnomorska',
    'bliski wschód',
    'meksykańska',
    'amerykańska',
];

const steps = ref<{ stepNo: number; section: string; text: string }[]>([
    { stepNo: 1, section: '', text: '' },
]);

const recipeIngredients = ref<
    {
        sortOrder: number;
        section: string;
        ingredientId: number;
        amountText: string;
        note: string;
    }[]
>([{ sortOrder: 1, section: '', ingredientId: 0, amountText: '', note: '' }]);

const CAT_COLORS: Record<string, string> = {
    ryby: '#3F70A8', mieso: '#B85A2E', wege: '#7A9050', wegan: '#2F6B3F',
    makarony: '#D9A441', ryz: '#B8946A', zupy: '#C53728', salatki: '#A8B83D',
    pieczywo: '#7A4A2E', desery: '#D27393', sniadania: '#E89556', przekaski: '#8A4A6D', napoje: '#2F8A8A',
};

const categories = [
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

onMounted(async () => {
    ingredients.value = await fetchIngredients('list');

    if (props.id) {
        const recipe = await fetchRecipesWithClient(qc, props.id);
        name.value = recipe.name;
        tagline.value = recipe.tagline ?? '';
        description.value = recipe.description ?? '';
        servings.value = recipe.servings;
        timeMinutes.value = recipe.timeMinutes;
        difficulty.value = recipe.difficulty;
        kcalPerServing.value = recipe.kcalPerServing ?? null;
        category.value = recipe.category;
        region.value = recipe.region;
        needsPrep.value = recipe.needsPrep ?? false;
        status.value = (recipe.status as 'draft' | 'published') ?? 'draft';
        diets.value = new Set(recipe.dietTags ?? []);
        practical.value = new Set(recipe.practicalTags ?? []);

        if (recipe.steps.length) {
            steps.value = recipe.steps.map((s) => ({
                stepNo: s.stepNo,
                section: s.section ?? '',
                text: s.text,
            }));
        }

        if (recipe.ingredients.length) {
            recipeIngredients.value = recipe.ingredients.map((i, idx) => ({
                sortOrder: i.sortOrder ?? idx + 1,
                section: i.section ?? '',
                ingredientId: i.ingredientId,
                amountText: i.amountText ?? '',
                note: i.note ?? '',
            }));
        }
    }
});

function addStep() {
    steps.value.push({ stepNo: steps.value.length + 1, section: '', text: '' });
}

function removeStep(idx: number) {
    steps.value.splice(idx, 1);
    steps.value.forEach((s, i) => (s.stepNo = i + 1));
}

function addIngredient() {
    recipeIngredients.value.push({
        sortOrder: recipeIngredients.value.length + 1,
        section: '',
        ingredientId: 0,
        amountText: '',
        note: '',
    });
}

function removeIngredient(idx: number) {
    recipeIngredients.value.splice(idx, 1);
    recipeIngredients.value.forEach((i, x) => (i.sortOrder = x + 1));
}

function toggleSet(set: Set<string>, v: string) {
    const n = new Set(set);
    if (n.has(v)) n.delete(v);
    else n.add(v);
    return n;
}

async function submit() {
    if (!name.value.trim() || servings.value < 1) return;
    saving.value = true;

    const body: IUpsertRecipeRequest = {
        name: name.value.trim(),
        tagline: tagline.value || null,
        description: description.value || null,
        servings: servings.value,
        timeMinutes: timeMinutes.value,
        difficulty: difficulty.value,
        kcalPerServing: kcalPerServing.value,
        category: category.value,
        region: region.value,
        status: status.value,
        needsPrep: needsPrep.value,
        dietTags: [...diets.value],
        practicalTags: [...practical.value],
        steps: steps.value
            .filter((s) => s.text.trim())
            .map((s) => ({
                stepNo: s.stepNo,
                section: s.section || null,
                text: s.text.trim(),
            })),
        ingredients: recipeIngredients.value
            .filter((i) => i.ingredientId > 0)
            .map((i) => ({
                sortOrder: i.sortOrder,
                section: i.section || null,
                ingredientId: i.ingredientId,
                amountText: i.amountText || null,
                note: i.note || null,
            })),
    };

    try {
        const result = props.id ? await updateRecipe(props.id, body) : await createRecipe(body);
        router.push({ name: 'recipe-details', params: { id: result.id } });
    } finally {
        saving.value = false;
    }
}
</script>

<template>
    <div class="container upsert-page">
        <!-- Page header -->
        <div class="upsert-header">
            <nav class="upsert-breadcrumb">
                <RouterLink
                    to="/recipes"
                    class="upsert-breadcrumb__link"
                >
                    Przepisy
                </RouterLink>
                <span class="upsert-breadcrumb__sep">/</span>
                <span class="upsert-breadcrumb__current">{{ id ? 'Edytuj' : 'Nowy' }}</span>
            </nav>
            <div class="upsert-header__row">
                <h1 class="upsert-title">{{ id ? 'Edytuj przepis' : 'Nowy przepis' }}</h1>
                <div class="upsert-status">
                    <span
                        class="upsert-status__dot"
                        :class="
                            status === 'published'
                                ? 'upsert-status__dot--pub'
                                : 'upsert-status__dot--draft'
                        "
                    />
                    {{
                        status === 'published'
                            ? 'opublikowany · widoczny publicznie'
                            : 'szkic · zmiany zapisują się auto'
                    }}
                </div>
            </div>
        </div>

        <div class="upsert-body">
            <!-- LEFT: main content -->
            <div class="upsert-main">
                <!-- 01 Podstawy -->
                <section class="upsert-section">
                    <div class="upsert-section__head">
                        <span class="upsert-section__n">01</span>
                        <h2 class="upsert-section__title">Podstawy</h2>
                    </div>

                    <div class="field">
                        <label class="field-label">
                            Nazwa przepisu
                            <span class="req">*</span>
                        </label>
                        <input
                            v-model="name"
                            class="input upsert-name-input"
                            type="text"
                            placeholder="np. Pieczony dorsz z masłem cytrynowym"
                        />
                    </div>

                    <div
                        class="field"
                        style="margin-top: 16px"
                    >
                        <label class="field-label">Krótki opis</label>
                        <div class="field-hint">// jedno-dwa zdania, pokazuje się pod tytułem</div>
                        <textarea
                            v-model="tagline"
                            class="textarea"
                            rows="3"
                            placeholder="Maślany, lekko kwaśny, gotowy w 30 minut. Jedna patelnia."
                        />
                    </div>
                </section>

                <!-- 02 Składniki -->
                <section class="upsert-section">
                    <div class="upsert-section__head">
                        <span class="upsert-section__n">02</span>
                        <h2 class="upsert-section__title">Składniki</h2>
                    </div>
                    <div class="upsert-section__hint">
                        // grupuj jeśli to ma sens (np. „na rybę / sos") — w prostych przepisach
                        jedna grupa wystarczy
                    </div>

                    <div
                        v-for="(ing, idx) in recipeIngredients"
                        :key="idx"
                        class="ingredient-row"
                    >
                        <select
                            v-model="ing.ingredientId"
                            class="select"
                            style="flex: 2"
                        >
                            <option
                                :value="0"
                                disabled
                            >
                                Wybierz składnik
                            </option>
                            <option
                                v-for="i in ingredients"
                                :key="i.id"
                                :value="i.id"
                            >
                                {{ i.name }}
                            </option>
                        </select>
                        <input
                            v-model="ing.amountText"
                            class="input"
                            type="text"
                            placeholder="Ilość (np. 2 ząbki)"
                            style="flex: 1"
                        />
                        <input
                            v-model="ing.section"
                            class="input input--mono"
                            type="text"
                            placeholder="Sekcja"
                            style="flex: 1"
                        />
                        <button
                            class="row-del-btn"
                            @click="removeIngredient(idx)"
                            :disabled="recipeIngredients.length <= 1"
                        >
                            ×
                        </button>
                    </div>

                    <button
                        class="ghost-add-btn"
                        @click="addIngredient"
                    >
                        + składnik
                    </button>
                </section>

                <!-- 03 Kroki -->
                <section class="upsert-section">
                    <div class="upsert-section__head">
                        <span class="upsert-section__n">03</span>
                        <h2 class="upsert-section__title">Kroki</h2>
                    </div>
                    <div class="upsert-section__hint">// grupuj jeśli przepis ma kilka etapów</div>

                    <div
                        v-for="(step, idx) in steps"
                        :key="idx"
                        class="step-row"
                    >
                        <span class="step-num-large">{{ step.stepNo }}</span>
                        <div style="flex: 1; display: flex; flex-direction: column; gap: 8px">
                            <input
                                v-model="step.section"
                                class="input input--mono"
                                type="text"
                                placeholder="Etap (opcja, np. Sos)"
                            />
                            <textarea
                                v-model="step.text"
                                class="textarea"
                                rows="2"
                                :placeholder="`Krok ${step.stepNo}…`"
                            />
                        </div>
                        <button
                            class="row-del-btn"
                            @click="removeStep(idx)"
                            :disabled="steps.length <= 1"
                        >
                            ×
                        </button>
                    </div>

                    <button
                        class="ghost-add-btn"
                        style="margin-left: 54px"
                        @click="addStep"
                    >
                        + krok
                    </button>
                </section>

                <!-- 04 Wskazówka -->
                <section class="upsert-section">
                    <div class="upsert-section__head">
                        <span class="upsert-section__n">04</span>
                        <h2 class="upsert-section__title">Wskazówka</h2>
                        <span class="upsert-section__optional">opcjonalne</span>
                    </div>
                    <div class="upsert-section__hint">
                        // opcjonalna porada na koniec przepisu (np. czego unikać, czym podać)
                    </div>
                    <textarea
                        v-model="tip"
                        class="textarea"
                        rows="3"
                        placeholder="Jeśli filety są bardzo cienkie, skróć pieczenie do 10 minut..."
                    />
                </section>
            </div>

            <!-- RIGHT: sticky sidebar -->
            <aside class="upsert-sidebar">
                <!-- Klasyfikacja -->
                <div class="side-card">
                    <div class="side-card__title">Klasyfikacja</div>

                    <div class="field">
                        <label class="field-label">
                            Kategoria
                            <span class="req">*</span>
                        </label>
                        <div class="cat-picker">
                            <button
                                v-for="cat in categories"
                                :key="cat.key"
                                class="cat-pill"
                                :class="{ 'cat-pill--active': category === cat.key }"
                                :style="category === cat.key ? `--cat-color: ${CAT_COLORS[cat.key]}` : `--cat-color: ${CAT_COLORS[cat.key]}`"
                                @click="category = category === cat.key ? null : cat.key"
                            >
                                <span class="cat-pill__dot" :style="{ background: CAT_COLORS[cat.key] }" />
                                {{ cat.label }}
                            </button>
                        </div>
                    </div>

                    <div class="field">
                        <label class="field-label">Region</label>
                        <div class="field-hint">
                            // opcjonalnie — pokazywany po kropce w odznace
                        </div>
                        <select
                            v-model="region"
                            class="select"
                        >
                            <option
                                v-for="r in REGIONS"
                                :key="r"
                                :value="r === '—' ? null : r"
                            >
                                {{ r }}
                            </option>
                        </select>
                    </div>

                    <div class="field">
                        <label class="field-label">
                            Trudność
                            <span class="req">*</span>
                        </label>
                        <div style="display: flex; gap: 6px">
                            <button
                                v-for="d in [1, 2, 3]"
                                :key="d"
                                class="diff-btn"
                                :class="{ 'is-active': difficulty === d }"
                                @click="difficulty = difficulty === d ? null : d"
                            >
                                <span class="diff-btn__dots">
                                    <span
                                        v-for="i in 3"
                                        :key="i"
                                        class="diff-btn__dot"
                                        :class="{ 'is-on': i <= d }"
                                    />
                                </span>
                                {{ { 1: 'Łatwe', 2: 'Średnie', 3: 'Trudne' }[d] }}
                            </button>
                        </div>
                    </div>
                </div>

                <!-- Czas i porcje -->
                <div class="side-card">
                    <div class="side-card__title">Czas i porcje</div>
                    <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 10px">
                        <div class="field">
                            <label class="field-label">
                                Czas (min)
                                <span class="req">*</span>
                            </label>
                            <input
                                v-model.number="timeMinutes"
                                class="input input--mono"
                                type="number"
                                min="1"
                                placeholder="30"
                            />
                        </div>
                        <div class="field">
                            <label class="field-label">
                                Porcje
                                <span class="req">*</span>
                            </label>
                            <input
                                v-model.number="servings"
                                class="input input--mono"
                                type="number"
                                min="1"
                            />
                        </div>
                    </div>
                    <div class="field">
                        <label class="field-label">Kalorie / porcja</label>
                        <div class="field-hint">// opcjonalnie</div>
                        <input
                            v-model.number="kcalPerServing"
                            class="input input--mono"
                            type="number"
                            min="1"
                            placeholder="np. 490"
                        />
                    </div>
                </div>

                <!-- Status -->
                <div class="side-card">
                    <div class="side-card__title">Status</div>
                    <div class="status-pill-toggle">
                        <button
                            v-for="opt in [
                                { id: 'draft', label: 'Szkic', dot: 'var(--c-mustard)' },
                                { id: 'published', label: 'Opublikowany', dot: 'var(--c-sage)' },
                            ]"
                            :key="opt.id"
                            class="status-pill-btn"
                            :class="{ 'status-pill-btn--active': status === opt.id }"
                            @click="status = opt.id as 'draft' | 'published'"
                        >
                            <span
                                class="status-pill-dot"
                                :style="{ background: opt.dot }"
                            />
                            {{ opt.label }}
                        </button>
                    </div>
                    <div class="status-hint">
                        {{
                            status === 'published'
                                ? '// widoczny publicznie'
                                : '// tylko dla administratora'
                        }}
                    </div>
                </div>

                <!-- Atrybuty -->
                <div class="side-card">
                    <div class="side-card__title">Atrybuty</div>
                    <div class="field">
                        <label class="prep-toggle-label">
                            <input
                                v-model="needsPrep"
                                type="checkbox"
                                class="prep-toggle-input"
                            />
                            <span class="prep-toggle-track">
                                <span class="prep-toggle-thumb" />
                            </span>
                            <span class="prep-toggle-text">Wymaga przygotowania</span>
                        </label>
                        <div class="field-hint">
                            // marynata, namaczanie, ciasto dzień wcześniej
                        </div>
                    </div>
                    <div class="field">
                        <label
                            for="diets-label"
                            class="field-label"
                            style="margin-bottom: 6px"
                        >
                            Diety
                        </label>
                        <div
                            class="chip-grid"
                            id="diets-label"
                        >
                            <button
                                v-for="d in DIETS"
                                :key="d"
                                class="chip"
                                :class="{ 'is-active': diets.has(d) }"
                                @click="diets = toggleSet(diets, d)"
                            >
                                {{ d }}
                            </button>
                        </div>
                    </div>
                    <div class="field">
                        <label
                            for="practical-label"
                            class="field-label"
                            style="margin-bottom: 6px"
                        >
                            Praktyczne
                        </label>
                        <div
                            class="chip-grid"
                            id="practical-label"
                        >
                            <button
                                v-for="p in PRACTICAL"
                                :key="p"
                                class="chip chip--mono"
                                :class="{ 'is-active': practical.has(p) }"
                                @click="practical = toggleSet(practical, p)"
                            >
                                {{ p }}
                            </button>
                        </div>
                    </div>
                </div>
            </aside>
        </div>
    </div>

    <!-- Fixed action bar -->
    <div class="action-bar">
        <div class="container action-bar__inner">
            <span class="action-bar__hint">
                // tryb edycji · zmiany zapisują się automatycznie jako szkic
            </span>
            <div style="display: flex; gap: 8px">
                <RouterLink
                    to="/recipes"
                    class="btn btn--ghost"
                >
                    Anuluj
                </RouterLink>
                <button class="btn btn--secondary">Podgląd</button>
                <button
                    class="btn"
                    :class="status === 'published' ? 'btn--pub' : 'btn--accent'"
                    :disabled="saving || !name.trim()"
                    @click="submit"
                >
                    {{
                        saving
                            ? 'Zapisywanie…'
                            : status === 'published'
                              ? 'Opublikowany · Zapisz'
                              : 'Opublikuj przepis'
                    }}
                </button>
            </div>
        </div>
    </div>
</template>

<style scoped>
.upsert-page {
    padding-top: 48px;
    padding-bottom: 120px;
}

/* Header */
.upsert-breadcrumb {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 13px;
    color: var(--ink-muted);
    margin-bottom: 14px;
}
.upsert-breadcrumb__link {
    color: var(--ink-muted);
    text-decoration: none;
}
.upsert-breadcrumb__link:hover {
    color: var(--ink);
}
.upsert-breadcrumb__sep {
    color: var(--ink-faint);
}
.upsert-breadcrumb__current {
    color: var(--ink);
    font-weight: 500;
}

.upsert-header__row {
    display: flex;
    align-items: flex-end;
    justify-content: space-between;
    gap: 24px;
    margin-bottom: 32px;
}
.upsert-title {
    font-size: 56px;
    line-height: 0.95;
    letter-spacing: -2.2px;
    font-weight: 500;
    margin: 0;
}
.upsert-status {
    display: flex;
    align-items: center;
    gap: 6px;
    font-family: var(--font-mono);
    font-size: 11px;
    color: var(--ink-muted);
    letter-spacing: 0.3px;
    padding-bottom: 10px;
}
.upsert-status__dot {
    width: 6px;
    height: 6px;
    border-radius: var(--r-pill);
    background: var(--c-sage);
}
.upsert-status__dot--draft {
    background: var(--c-mustard);
}
.upsert-status__dot--pub {
    background: var(--c-sage);
}

/* Body layout */
.upsert-body {
    display: grid;
    grid-template-columns: 1fr 360px;
    gap: 40px;
    align-items: start;
}

.upsert-main {
    display: flex;
    flex-direction: column;
    gap: 32px;
}

/* Section cards */
.upsert-section {
    background: var(--bg-alt);
    border: 1px solid var(--rule);
    border-radius: var(--r-lg);
    padding: 28px;
}
.upsert-section__head {
    display: flex;
    align-items: baseline;
    gap: 14px;
    margin-bottom: 6px;
}
.upsert-section__n {
    font-family: var(--font-mono);
    font-size: 11px;
    color: var(--accent);
    font-weight: 600;
}
.upsert-section__title {
    font-size: 20px;
    font-weight: 600;
    letter-spacing: -0.5px;
    margin: 0;
}
.upsert-section__optional {
    font-family: var(--font-mono);
    font-size: 10px;
    color: var(--ink-faint);
    letter-spacing: 0.5px;
    text-transform: uppercase;
}
.upsert-section__hint {
    font-family: var(--font-mono);
    font-size: 11px;
    color: var(--ink-muted);
    letter-spacing: 0.3px;
    margin-bottom: 22px;
}

/* Name input — larger */
.upsert-name-input {
    font-size: 18px;
    font-weight: 500;
    letter-spacing: -0.3px;
    padding: 12px 14px;
}

/* Ingredients */
.ingredient-row {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 8px;
}

/* Steps */
.step-row {
    display: grid;
    grid-template-columns: 40px 1fr 28px;
    gap: 14px;
    align-items: start;
    margin-bottom: 12px;
}
.step-num-large {
    font-size: 24px;
    font-weight: 600;
    color: var(--accent);
    line-height: 1.4;
    font-family: var(--font-sans);
    padding-top: 4px;
}

/* Common */
.row-del-btn {
    width: 28px;
    height: 28px;
    border: 1px solid var(--rule);
    background: var(--bg-alt);
    border-radius: var(--r-md);
    font-size: 18px;
    cursor: pointer;
    line-height: 1;
    color: var(--ink-muted);
    font-family: var(--font-mono);
    padding: 0;
}
.row-del-btn:disabled {
    opacity: 0.3;
    cursor: not-allowed;
}

.ghost-add-btn {
    background: transparent;
    border: none;
    color: var(--accent);
    font-family: var(--font-mono);
    font-size: 12px;
    font-weight: 600;
    letter-spacing: 0.3px;
    padding: 8px 0;
    cursor: pointer;
    margin-top: 8px;
    display: block;
}

/* Sidebar */
.upsert-sidebar {
    position: sticky;
    top: 90px;
    display: flex;
    flex-direction: column;
    gap: 20px;
}
.side-card {
    background: var(--bg-alt);
    border: 1px solid var(--rule);
    border-radius: var(--r-lg);
    padding: 20px;
    display: flex;
    flex-direction: column;
    gap: 14px;
}
.side-card__title {
    font-size: 13px;
    font-weight: 600;
    padding-bottom: 12px;
    border-bottom: 1px solid var(--rule);
}

/* Category picker */
.cat-picker { display: flex; flex-wrap: wrap; gap: 6px; margin-top: 4px; }
.cat-pill {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    background: transparent;
    border: 1px solid var(--rule-strong);
    color: var(--ink);
    padding: 4px 10px 4px 8px;
    border-radius: var(--r-pill);
    font-size: 11px;
    font-weight: 500;
    cursor: pointer;
    font-family: var(--font-sans);
    transition: all 0.12s;
}
.cat-pill--active {
    background: color-mix(in srgb, var(--cat-color) 18%, transparent);
    border-color: var(--cat-color);
    font-weight: 600;
}
.cat-pill__dot {
    width: 6px;
    height: 6px;
    border-radius: var(--r-pill);
    flex-shrink: 0;
}

/* Difficulty */
.diff-btn {
    flex: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 6px;
    background: transparent;
    border: 1px solid var(--rule-strong);
    color: var(--ink);
    padding: 10px 6px;
    border-radius: var(--r-md);
    font-size: 12px;
    font-weight: 500;
    cursor: pointer;
    font-family: var(--font-sans);
    transition: all 0.12s;
}
.diff-btn.is-active {
    background: var(--accent-soft);
    border-color: var(--accent);
    font-weight: 600;
}
.diff-btn__dots {
    display: flex;
    gap: 3px;
}
.diff-btn__dot {
    width: 6px;
    height: 6px;
    border-radius: var(--r-pill);
    background: var(--rule);
    border: 1px solid var(--rule-strong);
}
.diff-btn__dot.is-on {
    background: var(--accent);
}

/* Chips */
.chip-grid {
    display: flex;
    flex-wrap: wrap;
    gap: 6px;
}

/* Action bar */
.action-bar {
    position: fixed;
    bottom: 0;
    left: 0;
    right: 0;
    z-index: 40;
    background: rgba(248, 248, 255, 0.92);
    backdrop-filter: blur(12px);
    border-top: 1px solid var(--rule);
}
[data-theme='dark'] .action-bar {
    background: rgba(10, 10, 15, 0.92);
}

.action-bar__inner {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 16px;
    padding-top: 14px;
    padding-bottom: 14px;
}
.action-bar__hint {
    font-family: var(--font-mono);
    font-size: 11px;
    color: var(--ink-muted);
    letter-spacing: 0.3px;
}

@media (max-width: 900px) {
    .upsert-body {
        grid-template-columns: 1fr;
    }
    .upsert-sidebar {
        position: static;
    }
    .upsert-title {
        font-size: 40px;
    }
    .action-bar__hint {
        display: none;
    }
}

/* Status toggle */
/* Status pill toggle */
.status-pill-toggle {
    display: flex;
    background: var(--bg);
    border: 1px solid var(--rule);
    border-radius: var(--r-lg);
    padding: 3px;
    gap: 3px;
}
.status-pill-btn {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 7px;
    background: transparent;
    border: 1px solid transparent;
    border-radius: var(--r-md);
    padding: 9px 6px;
    font-size: 13px;
    font-weight: 500;
    cursor: pointer;
    font-family: var(--font-sans);
    color: var(--ink);
    transition: all 0.15s;
}
.status-pill-btn--active {
    background: var(--bg-alt);
    border-color: var(--rule-strong);
    font-weight: 600;
    box-shadow: 0 1px 3px rgba(10, 10, 15, 0.08);
}
.status-pill-dot {
    width: 7px;
    height: 7px;
    border-radius: var(--r-pill);
    flex-shrink: 0;
}
.status-hint {
    font-family: var(--font-mono);
    font-size: 10px;
    color: var(--ink-faint);
    letter-spacing: 0.3px;
}

/* btn--pub (sage green) */
.btn--pub {
    background: var(--c-sage);
    color: #fff;
    border: none;
}

/* NeedsPrep toggle */
.prep-toggle-label {
    display: flex;
    align-items: center;
    gap: 10px;
    cursor: pointer;
    user-select: none;
}
.prep-toggle-input {
    position: absolute;
    opacity: 0;
    width: 0;
    height: 0;
}
.prep-toggle-track {
    width: 32px;
    height: 18px;
    border-radius: var(--r-pill);
    flex-shrink: 0;
    background: var(--rule-strong);
    position: relative;
    transition: background 0.18s;
    display: inline-block;
}
.prep-toggle-input:checked ~ .prep-toggle-track {
    background: var(--c-mustard);
}
.prep-toggle-thumb {
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
.prep-toggle-input:checked ~ .prep-toggle-track .prep-toggle-thumb {
    left: 16px;
}
.prep-toggle-text {
    font-size: 13px;
    font-weight: 500;
    color: var(--ink);
}
</style>
