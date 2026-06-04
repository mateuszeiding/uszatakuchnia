<script lang="ts" setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';

import { fetchRecipes, createRecipe, updateRecipe } from '@/data/api/recipe/fetch';
import { fetchIngredients } from '@/data/api/ingredients/fetch';
import type { UpsertRecipeRequest } from '@/data/dtos/recipe/RecipeDto';
import type { IngredientDto } from '@/data/dtos/ingredients/IngredientDto';

const props = defineProps<{ id?: number }>();
const router = useRouter();

const saving = ref(false);
const ingredients = ref<IngredientDto[]>([]);

// Form state
const name = ref('');
const tagline = ref('');
const description = ref('');
const servings = ref(2);
const timeMinutes = ref<number | null>(null);
const difficulty = ref<number | null>(null);
const kcalPerServing = ref<number | null>(null);
const category = ref<string | null>(null);
const region = ref<string | null>(null);

const steps = ref<{ stepNo: number; section: string; text: string }[]>([
    { stepNo: 1, section: '', text: '' },
]);

const recipeIngredients = ref<{
    sortOrder: number;
    section: string;
    ingredientId: number;
    amountText: string;
    note: string;
}[]>([
    { sortOrder: 1, section: '', ingredientId: 0, amountText: '', note: '' },
]);

const categories = [
    { key: 'ryby', label: 'Ryby' },
    { key: 'mieso', label: 'Mięso' },
    { key: 'wege', label: 'Wege' },
    { key: 'wegan', label: 'Wegan' },
    { key: 'makarony', label: 'Makarony' },
    { key: 'ryz', label: 'Ryż / Kasze' },
    { key: 'zupy', label: 'Zupy' },
    { key: 'salatki', label: 'Sałatki' },
    { key: 'pieczywo', label: 'Pieczywo' },
    { key: 'desery', label: 'Desery' },
    { key: 'sniadania', label: 'Śniadania' },
    { key: 'przekaski', label: 'Przekąski' },
    { key: 'napoje', label: 'Napoje' },
];

onMounted(async () => {
    ingredients.value = await fetchIngredients('list');

    if (props.id) {
        const recipe = await fetchRecipes(props.id);
        name.value = recipe.name;
        tagline.value = recipe.tagline ?? '';
        description.value = recipe.description ?? '';
        servings.value = recipe.servings;
        timeMinutes.value = recipe.timeMinutes;
        difficulty.value = recipe.difficulty;
        kcalPerServing.value = recipe.kcalPerServing ?? null;
        category.value = recipe.category;
        region.value = recipe.region;

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

async function submit() {
    if (!name.value.trim() || servings.value < 1) return;
    saving.value = true;

    const body: UpsertRecipeRequest = {
        name: name.value.trim(),
        tagline: tagline.value || null,
        description: description.value || null,
        servings: servings.value,
        timeMinutes: timeMinutes.value,
        difficulty: difficulty.value,
        kcalPerServing: kcalPerServing.value,
        category: category.value,
        region: region.value,
        steps: steps.value.filter((s) => s.text.trim()).map((s) => ({
            stepNo: s.stepNo,
            section: s.section || null,
            text: s.text.trim(),
        })),
        ingredients: recipeIngredients.value.filter((i) => i.ingredientId > 0).map((i) => ({
            sortOrder: i.sortOrder,
            section: i.section || null,
            ingredientId: i.ingredientId,
            amountText: i.amountText || null,
            note: i.note || null,
        })),
    };

    try {
        const result = props.id
            ? await updateRecipe(props.id, body)
            : await createRecipe(body);
        router.push({ name: 'recipe-details', params: { id: result.id } });
    } finally {
        saving.value = false;
    }
}
</script>

<template>
    <div class="container upsert-layout">
        <div class="upsert-header">
            <h1 class="t-h2">{{ id ? 'Edytuj przepis' : 'Nowy przepis' }}</h1>
            <RouterLink to="/recipes" class="btn btn--ghost btn--sm">Anuluj</RouterLink>
        </div>

        <div class="upsert-body">
            <!-- LEFT: main content -->
            <div class="upsert-main">
                <!-- Podstawy -->
                <section class="upsert-section">
                    <div class="section-head"><span class="n">01</span><h2>Podstawy</h2><span class="line" /></div>

                    <div class="field" style="margin-bottom: 16px;">
                        <label class="field-label">Nazwa przepisu <span class="req">*</span></label>
                        <input v-model="name" class="input" type="text" placeholder="np. Soljanka" />
                    </div>

                    <div class="field" style="margin-bottom: 16px;">
                        <label class="field-label">Tagline</label>
                        <input v-model="tagline" class="input" type="text" placeholder="Krótki opis dania" />
                    </div>

                    <div class="field">
                        <label class="field-label">Opis</label>
                        <textarea v-model="description" class="textarea" rows="3" placeholder="Opcjonalny opis przepisu…" />
                    </div>
                </section>

                <!-- Składniki -->
                <section class="upsert-section">
                    <div class="section-head"><span class="n">02</span><h2>Składniki</h2><span class="line" /></div>

                    <div
                        v-for="(ing, idx) in recipeIngredients"
                        :key="idx"
                        class="ingredient-row"
                    >
                        <select v-model="ing.ingredientId" class="select" style="flex: 2;">
                            <option :value="0" disabled>Wybierz składnik</option>
                            <option v-for="i in ingredients" :key="i.id" :value="i.id">{{ i.name }}</option>
                        </select>
                        <input v-model="ing.amountText" class="input" type="text" placeholder="Ilość (np. 2 ząbki)" style="flex: 1;" />
                        <input v-model="ing.section" class="input" type="text" placeholder="Sekcja (opcja)" style="flex: 1;" />
                        <button class="icon-btn" @click="removeIngredient(idx)" :disabled="recipeIngredients.length <= 1">×</button>
                    </div>

                    <button class="btn btn--ghost btn--sm" style="margin-top: 8px;" @click="addIngredient">
                        + Dodaj składnik
                    </button>
                </section>

                <!-- Kroki -->
                <section class="upsert-section">
                    <div class="section-head"><span class="n">03</span><h2>Kroki</h2><span class="line" /></div>

                    <div
                        v-for="(step, idx) in steps"
                        :key="idx"
                        class="step-row"
                    >
                        <span class="step-num">{{ step.stepNo }}</span>
                        <div style="flex: 1; display: flex; flex-direction: column; gap: 8px;">
                            <input v-model="step.section" class="input input--mono" type="text" placeholder="Etap (opcja, np. Sos)" />
                            <textarea v-model="step.text" class="textarea" rows="2" :placeholder="`Krok ${step.stepNo}…`" />
                        </div>
                        <button class="icon-btn" @click="removeStep(idx)" :disabled="steps.length <= 1">×</button>
                    </div>

                    <button class="btn btn--ghost btn--sm" style="margin-top: 8px;" @click="addStep">
                        + Dodaj krok
                    </button>
                </section>
            </div>

            <!-- RIGHT: metadata sidebar -->
            <aside class="upsert-sidebar">
                <div class="surface" style="padding: 24px; position: sticky; top: 80px; display: flex; flex-direction: column; gap: 20px;">
                    <!-- Kategoria -->
                    <div class="field">
                        <label class="field-label">Kategoria</label>
                        <div style="display: flex; flex-wrap: wrap; gap: 6px; margin-top: 4px;">
                            <button
                                v-for="cat in categories"
                                :key="cat.key"
                                class="badge"
                                :class="[`cat-${cat.key}`, { 'is-active': category === cat.key }]"
                                style="cursor: pointer; border: 2px solid transparent;"
                                :style="category === cat.key ? 'border-color: var(--accent);' : ''"
                                @click="category = category === cat.key ? null : cat.key"
                            >
                                {{ cat.label }}
                            </button>
                        </div>
                    </div>

                    <!-- Region -->
                    <div class="field">
                        <label class="field-label">Region kuchni</label>
                        <input v-model="region" class="input" type="text" placeholder="np. polska, włoska, azjatycka" />
                    </div>

                    <!-- Trudność -->
                    <div class="field">
                        <label class="field-label">Trudność</label>
                        <div style="display: flex; gap: 8px;">
                            <button
                                v-for="d in [1,2,3]"
                                :key="d"
                                class="btn btn--sm"
                                :class="difficulty === d ? 'btn--accent' : 'btn--secondary'"
                                @click="difficulty = difficulty === d ? null : d"
                            >
                                {{ { 1: 'Łatwe', 2: 'Średnie', 3: 'Trudne' }[d] }}
                            </button>
                        </div>
                    </div>

                    <!-- Czas + porcje + kcal -->
                    <div class="field">
                        <label class="field-label">Czas (minuty)</label>
                        <input v-model.number="timeMinutes" class="input" type="number" min="1" placeholder="np. 45" />
                    </div>

                    <div class="field">
                        <label class="field-label">Porcje <span class="req">*</span></label>
                        <input v-model.number="servings" class="input" type="number" min="1" />
                    </div>

                    <div class="field">
                        <label class="field-label">Kcal / porcja</label>
                        <input v-model.number="kcalPerServing" class="input" type="number" min="1" placeholder="opcjonalnie" />
                    </div>

                    <button class="btn btn--accent" :disabled="saving || !name.trim()" @click="submit">
                        {{ saving ? 'Zapisywanie…' : (id ? 'Zapisz zmiany' : 'Opublikuj przepis') }}
                    </button>
                </div>
            </aside>
        </div>
    </div>
</template>

<style scoped>
.upsert-layout {
    padding-top: 48px;
    padding-bottom: 80px;
}

.upsert-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 40px;
}

.upsert-body {
    display: grid;
    grid-template-columns: 1fr 300px;
    gap: 40px;
    align-items: start;
}

.upsert-section {
    margin-bottom: 40px;
}

.ingredient-row {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 8px;
}

.step-row {
    display: flex;
    gap: 14px;
    align-items: flex-start;
    margin-bottom: 16px;
}

.step-num {
    width: 28px;
    height: 28px;
    border-radius: var(--r-pill);
    background: var(--accent);
    color: var(--on-accent);
    display: flex;
    align-items: center;
    justify-content: center;
    font-family: var(--font-mono);
    font-size: 12px;
    font-weight: 700;
    flex-shrink: 0;
    margin-top: 10px;
}

@media (max-width: 900px) {
    .upsert-body {
        grid-template-columns: 1fr;
    }
}
</style>
