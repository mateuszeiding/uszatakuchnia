<script lang="ts" setup>
import { useQueryClient } from '@tanstack/vue-query';
import { onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';

import IngAutosuggest from './components/IngAutosuggest.vue';

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
const categories = ref<string[]>([]);
const region = ref<string | null>(null);
const tip = ref('');
const photoUrl = ref<string | null>(null);
const photoPreview = ref<string | null>(null);

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

let _id = 0;
const uid = () => ++_id;

function splitAmountText(amountText: string | null, unitField: string | null): { qty: string; unit: string } {
    if (unitField) return { qty: amountText ?? '', unit: unitField };
    if (!amountText) return { qty: '', unit: '' };
    const UNIT_LIST = ['g', 'kg', 'ml', 'l', 'łyżka', 'łyżki', 'łyżeczka', 'łyżeczki', 'szczypta', 'ząbek', 'sztuka', 'sztuki', 'filety', 'pęczek'];
    const parts = amountText.trim().split(/\s+/);
    if (parts.length >= 2) {
        const lastWord = (parts[parts.length - 1] ?? '').toLowerCase();
        if (UNIT_LIST.includes(lastWord)) {
            return { qty: parts.slice(0, -1).join(' '), unit: lastWord };
        }
    }
    return { qty: amountText, unit: '' };
}

const UNITS = [
    '',
    'g',
    'kg',
    'ml',
    'l',
    'łyżka',
    'łyżki',
    'łyżeczka',
    'łyżeczki',
    'szczypta',
    'ząbek',
    'sztuka',
    'sztuki',
    'filety',
    'pęczek',
];

type IngItem = { id: number; ingredientId: number; qty: string; unit: string; note: string };
type IngGroup = { id: number; label: string; items: IngItem[] };
type StepItem = { id: number; text: string };
type StepGroup = { id: number; label: string; items: StepItem[] };

const ingGroups = ref<IngGroup[]>([
    {
        id: uid(),
        label: '',
        items: [{ id: uid(), ingredientId: 0, qty: '', unit: '', note: '' }],
    },
]);
const stepGroups = ref<StepGroup[]>([{ id: uid(), label: '', items: [{ id: uid(), text: '' }] }]);

const CAT_COLORS: Record<string, string> = {
    ryby: '#3F70A8',
    mieso: '#B85A2E',
    wege: '#7A9050',
    wegan: '#2F6B3F',
    makarony: '#D9A441',
    ryz: '#B8946A',
    zupy: '#C53728',
    salatki: '#A8B83D',
    pieczywo: '#7A4A2E',
    desery: '#D27393',
    sniadania: '#E89556',
    przekaski: '#8A4A6D',
    napoje: '#2F8A8A',
};

const CATEGORY_LIST = [
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
        categories.value = recipe.categories ?? [];
        region.value = recipe.region;
        needsPrep.value = recipe.needsPrep ?? false;
        status.value = (recipe.status as 'draft' | 'published') ?? 'draft';
        diets.value = new Set(recipe.dietTags ?? []);
        practical.value = new Set(recipe.practicalTags ?? []);
        photoUrl.value = recipe.photoUrl ?? null;
        photoPreview.value = recipe.photoUrl ?? null;

        if (recipe.steps.length) {
            const bySection = new Map<string, StepItem[]>();
            recipe.steps.forEach((s) => {
                const key = s.section ?? '';
                if (!bySection.has(key)) bySection.set(key, []);
                bySection.get(key)!.push({ id: uid(), text: s.text });
            });
            stepGroups.value = [...bySection.entries()].map(([label, items]) => ({
                id: uid(),
                label,
                items,
            }));
        }

        if (recipe.ingredients.length) {
            const bySection = new Map<string, IngItem[]>();
            recipe.ingredients.forEach((i) => {
                const key = i.section ?? '';
                if (!bySection.has(key)) bySection.set(key, []);
                const { qty, unit } = splitAmountText(i.amountText, i.unit);
                bySection.get(key)!.push({
                    id: uid(),
                    ingredientId: i.ingredientId,
                    qty,
                    unit,
                    note: i.note ?? '',
                });
            });
            ingGroups.value = [...bySection.entries()].map(([label, items]) => ({
                id: uid(),
                label,
                items,
            }));
        }
    }
});

// ─── Photo ────────────────────────────────────────────────────
function onPhotoFile(e: Event) {
    const file = (e.target as HTMLInputElement).files?.[0];
    if (!file) return;
    const reader = new FileReader();
    reader.onload = (ev) => {
        photoPreview.value = ev.target?.result as string;
    };
    reader.readAsDataURL(file);
}
function clearPhoto() {
    photoPreview.value = null;
    photoUrl.value = null;
}

// ─── Ingredient groups ────────────────────────────────────────
function addIngGroup() {
    ingGroups.value.push({
        id: uid(),
        label: '',
        items: [{ id: uid(), ingredientId: 0, qty: '', unit: '', note: '' }],
    });
}
function removeIngGroup(gid: number) {
    ingGroups.value = ingGroups.value.filter((g) => g.id !== gid);
}
function addIngItem(gid: number) {
    const g = ingGroups.value.find((x) => x.id === gid)!;
    g.items.push({ id: uid(), ingredientId: 0, qty: '', unit: '', note: '' });
}
function removeIngItem(gid: number, iid: number) {
    const g = ingGroups.value.find((x) => x.id === gid)!;
    g.items = g.items.filter((i) => i.id !== iid);
}

// ─── Step groups ──────────────────────────────────────────────
function addStepGroup() {
    stepGroups.value.push({ id: uid(), label: '', items: [{ id: uid(), text: '' }] });
}
function removeStepGroup(gid: number) {
    stepGroups.value = stepGroups.value.filter((g) => g.id !== gid);
}
function addStepItem(gid: number) {
    const g = stepGroups.value.find((x) => x.id === gid)!;
    g.items.push({ id: uid(), text: '' });
}
function removeStepItem(gid: number, iid: number) {
    const g = stepGroups.value.find((x) => x.id === gid)!;
    g.items = g.items.filter((i) => i.id !== iid);
}

function ingGroupPlaceholder(gi: number) {
    return ingGroups.value.length === 1
        ? 'Nazwa grupy (opcjonalnie)'
        : `Grupa ${gi + 1} - np. "Sos"`;
}
function stepGroupPlaceholder(gi: number) {
    return stepGroups.value.length === 1
        ? 'Nazwa etapu (opcjonalnie)'
        : `Etap ${gi + 1} - np. "Smażenie"`;
}

function globalStepIdx(gi: number, ii: number): number {
    let n = 0;
    for (let g = 0; g < stepGroups.value.length; g++) {
        const grp = stepGroups.value[g];
        if (!grp) continue;
        for (let i = 0; i < grp.items.length; i++) {
            n++;
            if (g === gi && i === ii) return n;
        }
    }
    return n;
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

    let stepNo = 0;
    const body: IUpsertRecipeRequest = {
        name: name.value.trim(),
        tagline: tagline.value || null,
        description: description.value || null,
        servings: servings.value,
        timeMinutes: timeMinutes.value,
        difficulty: difficulty.value,
        kcalPerServing: kcalPerServing.value,
        photoUrl: photoUrl.value,
        categories: categories.value,
        region: region.value,
        status: status.value,
        needsPrep: needsPrep.value,
        dietTags: [...diets.value],
        practicalTags: [...practical.value],
        steps: stepGroups.value.flatMap((g) =>
            g.items
                .filter((s) => s.text.trim())
                .map((s) => ({
                    stepNo: ++stepNo,
                    section: g.label || null,
                    text: s.text.trim(),
                }))
        ),
        ingredients: ingGroups.value.flatMap((g, gi) =>
            g.items
                .filter((i) => i.ingredientId > 0)
                .map((i, ii) => ({
                    sortOrder: gi * 100 + ii + 1,
                    section: g.label || null,
                    ingredientId: i.ingredientId,
                    amountText: i.qty || null,
                    unit: i.unit || null,
                    note: i.note || null,
                }))
        ),
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

                    <div class="basics-grid">
                        <!-- Photo zone -->
                        <label class="photo-zone">
                            <input
                                type="file"
                                accept="image/*"
                                class="photo-zone__input"
                                @change="onPhotoFile"
                            />
                            <img
                                v-if="photoPreview"
                                :src="photoPreview"
                                class="photo-zone__img"
                                alt=""
                            />
                            <template v-else>
                                <svg
                                    width="28"
                                    height="28"
                                    viewBox="0 0 28 28"
                                    fill="none"
                                >
                                    <rect
                                        x="3"
                                        y="5"
                                        width="22"
                                        height="18"
                                        rx="2"
                                        stroke="currentColor"
                                        stroke-width="1.4"
                                    />
                                    <circle
                                        cx="9"
                                        cy="11"
                                        r="2"
                                        stroke="currentColor"
                                        stroke-width="1.4"
                                    />
                                    <path
                                        d="M3 19l6-5 6 4 4-3 6 5"
                                        stroke="currentColor"
                                        stroke-width="1.4"
                                        stroke-linejoin="round"
                                    />
                                </svg>
                                <span class="photo-zone__label">Dodaj zdjęcie</span>
                                <span class="photo-zone__hint">1:1 · JPG / PNG / WEBP</span>
                            </template>
                            <button
                                v-if="photoPreview"
                                class="photo-zone__clear"
                                @click.prevent="clearPhoto"
                            >
                                ×
                            </button>
                        </label>

                        <!-- Name + tagline -->
                        <div class="basics-fields">
                            <div class="field">
                                <label
                                    for="recipe-name"
                                    class="field-label"
                                >
                                    Nazwa przepisu
                                    <span class="req">*</span>
                                </label>
                                <input
                                    id="recipe-name"
                                    v-model="name"
                                    class="input upsert-name-input"
                                    type="text"
                                    placeholder="np. Pieczony dorsz z masłem cytrynowym"
                                />
                            </div>

                            <div class="field">
                                <label
                                    for="recipe-tagline"
                                    class="field-label"
                                >
                                    Krótki opis
                                </label>
                                <div class="field-hint">
                                    // jedno-dwa zdania, pokazuje się pod tytułem
                                </div>
                                <textarea
                                    id="recipe-tagline"
                                    v-model="tagline"
                                    class="textarea"
                                    rows="3"
                                    placeholder="Maślany, lekko kwaśny, gotowy w 30 minut. Jedna patelnia."
                                />
                            </div>
                        </div>
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

                    <div class="groups-list">
                        <div
                            v-for="(g, gi) in ingGroups"
                            :key="g.id"
                            class="group-box"
                        >
                            <div class="group-box__head">
                                <input
                                    v-model="g.label"
                                    class="group-box__label-input"
                                    :placeholder="ingGroupPlaceholder(gi)"
                                />
                                <button
                                    v-if="ingGroups.length > 1"
                                    class="row-del-btn"
                                    @click="removeIngGroup(g.id)"
                                >
                                    ×
                                </button>
                            </div>

                            <div class="ing-rows">
                                <div
                                    v-for="item in g.items"
                                    :key="item.id"
                                    class="ing-row"
                                >
                                    <input
                                        v-model="item.qty"
                                        class="input input--mono ing-row__qty"
                                        type="text"
                                        placeholder="Ilość"
                                    />
                                    <select
                                        v-model="item.unit"
                                        class="select ing-row__unit"
                                    >
                                        <option
                                            v-for="u in UNITS"
                                            :key="u"
                                            :value="u"
                                        >
                                            {{ u || '— jed. —' }}
                                        </option>
                                    </select>
                                    <IngAutosuggest
                                        v-model="item.ingredientId"
                                        :ingredients="ingredients"
                                    />

                                    <button
                                        class="row-del-btn"
                                        :disabled="g.items.length <= 1"
                                        @click="removeIngItem(g.id, item.id)"
                                    >
                                        ×
                                    </button>
                                </div>
                            </div>

                            <button
                                class="ghost-add-btn"
                                @click="addIngItem(g.id)"
                            >
                                + składnik
                            </button>
                        </div>

                        <button
                            class="ghost-add-btn ghost-add-btn--dashed"
                            @click="addIngGroup"
                        >
                            + Dodaj grupę składników
                        </button>
                    </div>
                </section>

                <!-- 03 Kroki -->
                <section class="upsert-section">
                    <div class="upsert-section__head">
                        <span class="upsert-section__n">03</span>
                        <h2 class="upsert-section__title">Kroki</h2>
                    </div>
                    <div class="upsert-section__hint">// grupuj jeśli przepis ma kilka etapów</div>

                    <!-- View mode stub -->
                    <div class="view-mode-stub">
                        <span class="view-mode-stub__active">etapy</span>
                        <span class="view-mode-stub__inactive">
                            timing
                            <span class="view-mode-stub__soon">wkrótce</span>
                        </span>
                    </div>

                    <div class="groups-list">
                        <div
                            v-for="(g, gi) in stepGroups"
                            :key="g.id"
                            class="group-box"
                        >
                            <div class="group-box__head">
                                <input
                                    v-model="g.label"
                                    class="group-box__label-input"
                                    :placeholder="stepGroupPlaceholder(gi)"
                                />
                                <button
                                    v-if="stepGroups.length > 1"
                                    class="row-del-btn"
                                    @click="removeStepGroup(g.id)"
                                >
                                    ×
                                </button>
                            </div>

                            <div class="step-rows">
                                <div
                                    v-for="(item, ii) in g.items"
                                    :key="item.id"
                                    class="step-row"
                                >
                                    <span class="step-num-large">{{ globalStepIdx(gi, ii) }}</span>
                                    <textarea
                                        v-model="item.text"
                                        class="textarea"
                                        rows="2"
                                        :placeholder="
                                            globalStepIdx(gi, ii) === 1
                                                ? 'Rozgrzej piekarnik do 200°C...'
                                                : 'Następny krok...'
                                        "
                                    />
                                    <button
                                        class="row-del-btn"
                                        :disabled="g.items.length <= 1 && stepGroups.length <= 1"
                                        @click="removeStepItem(g.id, item.id)"
                                    >
                                        ×
                                    </button>
                                </div>
                            </div>

                            <button
                                class="ghost-add-btn"
                                style="margin-left: 54px"
                                @click="addStepItem(g.id)"
                            >
                                + krok
                            </button>
                        </div>

                        <button
                            class="ghost-add-btn ghost-add-btn--dashed"
                            @click="addStepGroup"
                        >
                            + Dodaj etap
                        </button>
                    </div>
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
                    <label
                        for="recipe-tip"
                        class="sr-only"
                    >
                        Wskazówka
                    </label>
                    <textarea
                        id="recipe-tip"
                        v-model="tip"
                        class="textarea"
                        rows="3"
                        placeholder="Jeśli filety są bardzo cienkie, skróć pieczenie do 10 minut..."
                    />
                </section>
            </div>

            <!-- RIGHT: sticky sidebar -->
            <aside class="upsert-sidebar">
                <!-- Status — first -->
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

                <!-- Klasyfikacja -->
                <div class="side-card">
                    <div class="side-card__title">Klasyfikacja</div>

                    <div class="field">
                        <label
                            for="cat-picker"
                            class="field-label"
                        >
                            Kategoria
                            <span class="req">*</span>
                        </label>
                        <div
                            id="cat-picker"
                            class="cat-picker"
                        >
                            <button
                                v-for="cat in CATEGORY_LIST"
                                :key="cat.key"
                                class="cat-pill"
                                :class="{ 'cat-pill--active': categories.includes(cat.key) }"
                                :style="`--cat-color: ${CAT_COLORS[cat.key]}`"
                                @click="
                                    categories = categories.includes(cat.key)
                                        ? categories.filter((c) => c !== cat.key)
                                        : [...categories, cat.key]
                                "
                            >
                                <span
                                    class="cat-pill__dot"
                                    :style="{ background: CAT_COLORS[cat.key] }"
                                />
                                {{ cat.label }}
                            </button>
                        </div>
                    </div>

                    <div class="field">
                        <label
                            for="recipe-region"
                            class="field-label"
                        >
                            Region
                        </label>
                        <div class="field-hint">
                            // opcjonalnie — pokazywany po kropce w odznace
                        </div>
                        <select
                            id="recipe-region"
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
                        <label
                            for="diff-picker"
                            class="field-label"
                        >
                            Trudność
                            <span class="req">*</span>
                        </label>
                        <div
                            id="diff-picker"
                            style="display: flex; gap: 6px"
                        >
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
                    <div class="two-col-grid">
                        <div class="field">
                            <label
                                for="recipe-time"
                                class="field-label"
                            >
                                Czas (min)
                                <span class="req">*</span>
                            </label>
                            <input
                                id="recipe-time"
                                v-model.number="timeMinutes"
                                class="input input--mono"
                                type="number"
                                min="1"
                                placeholder="30"
                            />
                        </div>
                        <div class="field">
                            <label
                                for="recipe-servings"
                                class="field-label"
                            >
                                Porcje
                                <span class="req">*</span>
                            </label>
                            <input
                                id="recipe-servings"
                                v-model.number="servings"
                                class="input input--mono"
                                type="number"
                                min="1"
                            />
                        </div>
                    </div>
                    <div class="field">
                        <label
                            for="recipe-kcal"
                            class="field-label"
                        >
                            Kalorie / porcja
                        </label>
                        <div class="field-hint">// opcjonalnie</div>
                        <input
                            id="recipe-kcal"
                            v-model.number="kcalPerServing"
                            class="input input--mono"
                            type="number"
                            min="1"
                            placeholder="np. 490"
                        />
                    </div>
                </div>

                <!-- Atrybuty -->
                <div class="side-card">
                    <div class="side-card__title">Atrybuty</div>

                    <div class="attr-group">
                        <span class="field-label">Diety</span>
                        <div class="chip-grid">
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

                    <div class="attr-group">
                        <span class="field-label">Praktyczne</span>
                        <div class="chip-grid">
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

                    <div class="attr-sep">
                        <label class="prep-toggle-label">
                            <input
                                v-model="needsPrep"
                                type="checkbox"
                                class="prep-toggle-input"
                            />
                            <span class="prep-toggle-track">
                                <span class="prep-toggle-thumb" />
                            </span>
                            <div>
                                <div class="prep-toggle-text">Wymaga przygotowania</div>
                                <div class="prep-toggle-sub">// np. marynata dzień wcześniej</div>
                            </div>
                        </label>
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
.sr-only {
    position: absolute;
    width: 1px;
    height: 1px;
    padding: 0;
    margin: -1px;
    overflow: hidden;
    clip: rect(0, 0, 0, 0);
    white-space: nowrap;
    border: 0;
}

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

/* Section 01 basics grid */
.basics-grid {
    display: grid;
    grid-template-columns: 220px 1fr;
    gap: 24px;
    align-items: start;
    margin-top: 22px;
}
.basics-fields {
    display: flex;
    flex-direction: column;
    gap: 16px;
}

/* Photo zone */
.photo-zone {
    aspect-ratio: 1 / 1;
    width: 100%;
    background: repeating-linear-gradient(45deg, var(--bg) 0 8px, var(--bg-alt) 8px 9px);
    border: 1.5px dashed var(--rule-strong);
    border-radius: var(--r-lg);
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    flex-direction: column;
    gap: 10px;
    color: var(--ink-muted);
    position: relative;
    overflow: hidden;
}
.photo-zone__input {
    display: none;
}
.photo-zone__img {
    position: absolute;
    inset: 0;
    width: 100%;
    height: 100%;
    object-fit: cover;
}
.photo-zone__label {
    font-size: 13px;
    font-weight: 500;
    color: var(--ink);
}
.photo-zone__hint {
    font-family: var(--font-mono);
    font-size: 10px;
    color: var(--ink-faint);
    letter-spacing: 0.3px;
}
.photo-zone__clear {
    position: absolute;
    top: 8px;
    right: 8px;
    background: rgba(10, 10, 15, 0.6);
    color: #fff;
    border: none;
    border-radius: var(--r-pill);
    width: 26px;
    height: 26px;
    cursor: pointer;
    font-size: 16px;
    line-height: 1;
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1;
}

/* Name input — larger */
.upsert-name-input {
    font-size: 18px;
    font-weight: 500;
    letter-spacing: -0.3px;
    padding: 12px 14px;
}

/* Groups */
.groups-list {
    display: flex;
    flex-direction: column;
    gap: 18px;
}
.group-box {
    background: var(--bg);
    border: 1px solid var(--rule);
    border-radius: var(--r-lg);
    padding: 16px;
}
.group-box__head {
    display: flex;
    align-items: center;
    gap: 10px;
    margin-bottom: 12px;
}
.group-box__label-input {
    flex: 1;
    background: transparent;
    border: none;
    outline: none;
    font-size: 13px;
    font-weight: 600;
    color: var(--ink);
    font-family: var(--font-mono);
    letter-spacing: 0.4px;
    text-transform: uppercase;
    padding: 0;
}
.group-box__label-input::placeholder {
    color: var(--ink-faint);
    font-weight: 500;
    text-transform: none;
}

/* Ingredient rows */
.ing-rows {
    display: flex;
    flex-direction: column;
    gap: 6px;
}
.ing-row {
    display: grid;
    grid-template-columns: 90px 120px 1fr 28px;
    gap: 8px;
    align-items: center;
}
.ing-row__unit {
    padding: 8px 10px;
    font-size: 13px;
}
.ing-row__qty {
    padding: 8px 10px;
}
.ing-row__name {
    padding: 8px 10px;
    font-size: 13px;
}
.ing-row__note {
    padding: 8px 10px;
    font-size: 12px;
}

/* Step rows */
.step-rows {
    display: flex;
    flex-direction: column;
    gap: 12px;
}
.step-row {
    display: grid;
    grid-template-columns: 40px 1fr 28px;
    gap: 14px;
    align-items: start;
}
.step-num-large {
    font-size: 24px;
    font-weight: 600;
    color: var(--accent);
    line-height: 1.4;
    font-family: var(--font-sans);
    padding-top: 4px;
}

/* View mode stub */
.view-mode-stub {
    display: inline-flex;
    align-items: center;
    gap: 4px;
    background: var(--bg);
    border: 1px solid var(--rule);
    border-radius: var(--r-pill);
    padding: 4px;
    font-family: var(--font-mono);
    font-size: 11px;
    margin-bottom: 18px;
}
.view-mode-stub__active {
    padding: 5px 12px;
    border-radius: var(--r-pill);
    background: var(--accent);
    color: #fff;
    font-weight: 600;
}
.view-mode-stub__inactive {
    padding: 5px 12px;
    border-radius: var(--r-pill);
    color: var(--ink-faint);
    font-weight: 500;
    display: inline-flex;
    align-items: center;
    gap: 6px;
    cursor: not-allowed;
}
.view-mode-stub__soon {
    font-size: 9px;
}

/* Common row controls */
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
    flex-shrink: 0;
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
.ghost-add-btn--dashed {
    border: 1px dashed var(--rule-strong);
    padding: 10px;
    border-radius: var(--r-md);
    font-family: var(--font-sans);
    font-size: 13px;
    font-weight: 500;
    width: 100%;
    margin-top: 0;
    text-align: center;
}

/* Sidebar */
.upsert-sidebar {
    position: sticky;
    top: 90px;
    display: flex;
    flex-direction: column;
    gap: 18px;
}
.side-card {
    background: var(--bg-alt);
    border: 1px solid var(--rule);
    border-radius: var(--r-lg);
    padding: 18px;
    display: flex;
    flex-direction: column;
    gap: 14px;
}
.side-card__title {
    font-size: 13px;
    font-weight: 600;
    padding-bottom: 10px;
    border-bottom: 1px solid var(--rule);
}

/* Two-col grid for time/servings */
.two-col-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 10px;
}

/* Category picker */
.cat-picker {
    display: flex;
    flex-wrap: wrap;
    gap: 6px;
    margin-top: 4px;
}
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

/* Attrs */
.attr-group {
    display: flex;
    flex-direction: column;
    gap: 8px;
}
.attr-sep {
    padding-top: 12px;
    border-top: 1px solid var(--rule);
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

/* Status toggle */
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
    line-height: 1.3;
}
.prep-toggle-sub {
    font-family: var(--font-mono);
    font-size: 10px;
    color: var(--ink-faint);
    letter-spacing: 0.3px;
    margin-top: 1px;
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
    .basics-grid {
        grid-template-columns: 1fr;
    }
    .ing-row {
        grid-template-columns: 70px 100px 1fr 28px;
    }
}
</style>
