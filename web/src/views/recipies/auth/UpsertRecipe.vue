<script lang="ts" setup>
import { useQueryClient } from '@tanstack/vue-query';
import { computed, nextTick, onBeforeUnmount, onMounted, ref } from 'vue';
import { onBeforeRouteLeave, useRouter } from 'vue-router';

import FormActionBar, { type BusyAction } from './components/FormActionBar.vue';
import FormSection from './components/FormSection.vue';
import IngredientsEditor from './components/IngredientsEditor.vue';
import PhotoField from './components/PhotoField.vue';
import StepsEditor from './components/StepsEditor.vue';
import TaxonomyField from './components/TaxonomyField.vue';
import {
    ERROR_ORDER,
    type FormErrors,
    fromRecipe,
    type IIngItem,
    type IRecipeForm,
    newForm,
    PREP_LEADS,
    toRequest,
    validate,
} from './recipeForm';

import UFieldError from '@/components/form/UFieldError.vue';
import IngredientModal from '@/components/IngredientModal.vue';
import USpinner from '@/components/USpinner.vue';
import { confirmAction } from '@/composables/useConfirm';
import { showToast } from '@/composables/useToast';
import { fetchIngredientsWithClient, invalidateIngredients } from '@/data/api/ingredients/fetch';
import {
    createRecipe,
    deleteRecipe,
    fetchRecipesWithClient,
    invalidateRecipes,
    updateRecipe,
} from '@/data/api/recipe/fetch';
import { fetchDictionaries } from '@/data/api/taxonomy/fetch';
import type { IngredientDto } from '@/data/dtos/ingredients/IngredientDto';
import { type IFieldError, type RecipeStatus, STATUS_LABEL } from '@/data/dtos/recipe/RecipeDto';
import type { IDictionaries } from '@/data/dtos/taxonomy/DictionariesDto';

const props = defineProps<{ id?: string | number }>();
const router = useRouter();
const qc = useQueryClient();

const PANEL = '/admin/recipes';
const SIDEBAR_AXES_LAST = ['diet', 'practical'];

const loading = ref(true);
const recipeId = ref<number | null>(props.id ? Number(props.id) : null);
const isNew = computed(() => recipeId.value == null);
const form = ref<IRecipeForm>(newForm());
const saved = ref(JSON.stringify(form.value));
const status = ref<RecipeStatus>('draft');
const lastSaved = ref<Date | null>(null);
const busy = ref<BusyAction>(null);
const showErrors = ref<'draft' | 'full' | null>(null);
const serverError = ref('');
const dicts = ref<IDictionaries>();
const ingredients = ref<IngredientDto[]>([]);
const addIng = ref<{ name: string; item: IIngItem } | null>(null);
let leaving = false;

const dirty = computed(() => JSON.stringify(form.value) !== saved.value);

onMounted(async () => {
    const [d, ings] = await Promise.all([fetchDictionaries(), fetchIngredientsWithClient(qc)]);
    dicts.value = d;
    ingredients.value = ings;
    if (recipeId.value) {
        const r = await fetchRecipesWithClient(qc, recipeId.value);
        form.value = fromRecipe(r);
        status.value = r.status;
    }
    saved.value = JSON.stringify(form.value);
    loading.value = false;
});

// ─── Walidacja ────────────────────────────────────────────────
const validation = computed(() => validate(form.value, dicts.value));
const errors = computed<FormErrors>(() => {
    if (showErrors.value === 'full') return validation.value.errors;
    if (showErrors.value === 'draft' && validation.value.errors.name)
        return { name: validation.value.errors.name };
    return {};
});
const rowErrors = computed(() => (showErrors.value ? validation.value.rows : {}));
const errorCount = computed(() =>
    showErrors.value === 'full' ? Object.keys(validation.value.errors).length : 0
);

async function scrollToFirstError() {
    await nextTick();
    const hit = ERROR_ORDER.find(([k]) => errors.value[k as keyof FormErrors]);
    const el = hit && document.getElementById(hit[1]);
    if (!el) return;
    window.scrollTo({
        top: el.getBoundingClientRect().top + window.scrollY - 110,
        behavior: 'smooth',
    });
    setTimeout(
        () =>
            el
                .querySelector<HTMLElement>('input,textarea,select,button')
                ?.focus({ preventScroll: true }),
        350
    );
}

// ─── Zapis ────────────────────────────────────────────────────
async function submit(kind: Exclude<BusyAction, null>, nextStatus: RecipeStatus, toast: string) {
    busy.value = kind;
    serverError.value = '';
    try {
        const body = toRequest(form.value, nextStatus);
        const res = recipeId.value
            ? await updateRecipe(recipeId.value, body)
            : await createRecipe(body);
        const created = recipeId.value == null;
        recipeId.value = res.id;
        status.value = nextStatus;
        saved.value = JSON.stringify(form.value);
        lastSaved.value = new Date();
        showErrors.value = null;
        invalidateRecipes(qc);
        showToast(toast);
        if (created) router.replace({ name: 'recipe-edit', params: { id: res.id } });
    } catch (e) {
        const data = (e as { response?: { data?: { fields?: IFieldError[] } } }).response?.data;
        serverError.value = data?.fields?.length
            ? data.fields.map((f) => f.message).join(' ')
            : 'Nie udało się zapisać. Spróbuj ponownie.';
    } finally {
        busy.value = null;
    }
}

function saveDraft() {
    if (validation.value.errors.name) {
        showErrors.value = 'draft';
        scrollToFirstError();
        return;
    }
    submit('draft', 'draft', 'Szkic zapisany');
}

function requireComplete(): boolean {
    if (!Object.keys(validation.value.errors).length) return true;
    showErrors.value = 'full';
    scrollToFirstError();
    return false;
}

function publish() {
    if (requireComplete()) submit('publish', 'published', 'Przepis opublikowany');
}

function saveChanges() {
    if (requireComplete()) submit('save', 'published', 'Zmiany zapisane');
}

async function unpublish() {
    const ok = await confirmAction({
        title: 'Cofnąć publikację?',
        body: 'Przepis zniknie z publicznej listy i wróci do szkiców. Niezapisane zmiany też zostaną zapisane.',
        confirmLabel: 'Cofnij publikację',
    });
    if (ok) submit('unpublish', 'draft', 'Przepis przeniesiony do szkiców');
}

async function cancel() {
    if (dirty.value) {
        const ok = await confirmAction({
            title: 'Odrzucić niezapisane zmiany?',
            body: 'Zmiany od ostatniego zapisu przepadną.',
            confirmLabel: 'Odrzuć i wyjdź',
            danger: true,
        });
        if (!ok) return;
    }
    leaving = true;
    router.push(PANEL);
}

async function removeRecipe() {
    if (!recipeId.value) return;
    const ok = await confirmAction({
        title: `Usunąć „${form.value.name || 'ten przepis'}”?`,
        body: 'Przepis zniknie z listy i z panelu. Tej operacji nie da się cofnąć.',
        confirmLabel: 'Usuń przepis',
        danger: true,
    });
    if (!ok) return;
    await deleteRecipe(recipeId.value);
    invalidateRecipes(qc);
    showToast('Przepis usunięty');
    leaving = true;
    router.push(PANEL);
}

// ─── Ostrzeżenie przy wyjściu ─────────────────────────────────
const onBeforeUnload = (e: BeforeUnloadEvent) => {
    if (!dirty.value) return;
    e.preventDefault();
    e.returnValue = '';
};
onMounted(() => window.addEventListener('beforeunload', onBeforeUnload));
onBeforeUnmount(() => window.removeEventListener('beforeunload', onBeforeUnload));

onBeforeRouteLeave((to) => {
    if (leaving || !dirty.value) return true;
    // przejście /recipes/new → /recipes/:id/edit po pierwszym zapisie
    if (to.name === 'recipe-edit' && Number(to.params.id) === recipeId.value) return true;
    return confirmAction({
        title: 'Wyjść bez zapisywania?',
        body: 'Zmiany od ostatniego zapisu przepadną.',
        confirmLabel: 'Wyjdź',
        danger: true,
    });
});

// ─── Składnik spoza bazy ──────────────────────────────────────
function onIngredientCreated(ing: IngredientDto) {
    ingredients.value = [...ingredients.value, ing];
    if (addIng.value) {
        addIng.value.item.ingredientId = ing.id;
        addIng.value.item.name = ing.name;
    }
    invalidateIngredients(qc);
    showToast(`„${ing.name}” dodany do bazy`);
    addIng.value = null;
}

// ─── Sidebar ──────────────────────────────────────────────────
const classificationAxes = computed(() =>
    (dicts.value?.axes ?? []).filter((a) => !SIDEBAR_AXES_LAST.includes(a.code))
);
const dietAxis = computed(() => dicts.value?.axes.find((a) => a.code === 'diet'));
const practicalAxis = computed(() => dicts.value?.axes.find((a) => a.code === 'practical'));

const axisModel = (code: string) => form.value.taxonomy[code] ?? [];
function setAxis(code: string, ids: number[]) {
    form.value.taxonomy = { ...form.value.taxonomy, [code]: ids };
}

const prepSummary = computed(() => {
    const prepSteps = form.value.stepGroups.flatMap((g) => g.items).filter((s) => s.prep);
    const maxLead = prepSteps.reduce((m, s) => Math.max(m, PREP_LEADS.indexOf(s.prepLead)), -1);
    return { count: prepSteps.length, lead: maxLead >= 0 ? PREP_LEADS[maxLead] : null };
});
</script>

<template>
    <div class="container upsert-page">
        <header class="upsert-head">
            <nav
                class="breadcrumb"
                aria-label="Ścieżka"
            >
                <RouterLink to="/admin">Panel</RouterLink>
                <span
                    class="sep"
                    aria-hidden="true"
                >
                    /
                </span>
                <RouterLink :to="PANEL">Przepisy</RouterLink>
                <span
                    class="sep"
                    aria-hidden="true"
                >
                    /
                </span>
                <span class="current">{{ isNew ? 'Nowy' : 'Edycja' }}</span>
            </nav>
            <div class="upsert-head__row">
                <div class="upsert-head__title">
                    <h1 class="upsert-title">{{ isNew ? 'Nowy przepis' : 'Edytuj przepis' }}</h1>
                    <span
                        class="status-pill"
                        :class="`status-pill--${status}`"
                    >
                        <span
                            class="status-pill__dot"
                            aria-hidden="true"
                        />
                        {{ STATUS_LABEL[status] }}
                    </span>
                </div>
                <button
                    v-if="!isNew"
                    type="button"
                    class="btn btn--danger"
                    @click="removeRecipe"
                >
                    Usuń przepis
                </button>
            </div>
        </header>

        <div
            v-if="loading"
            class="upsert-loading"
        >
            <USpinner :size="22" />
        </div>

        <div
            v-else
            class="upsert-body"
        >
            <div class="upsert-main">
                <FormSection
                    n="01"
                    title="Podstawy"
                >
                    <div class="basics">
                        <PhotoField v-model="form.photoUrl" />
                        <div class="basics__fields">
                            <div
                                id="f-name"
                                class="field"
                            >
                                <label
                                    for="in-name"
                                    class="field-label"
                                >
                                    Nazwa przepisu
                                    <span class="req">*</span>
                                </label>
                                <input
                                    id="in-name"
                                    v-model="form.name"
                                    class="input basics__name"
                                    :class="{ 'is-invalid': errors.name }"
                                    :aria-invalid="!!errors.name || undefined"
                                    placeholder="np. Pieczony dorsz z masłem cytrynowym"
                                />
                                <UFieldError v-if="errors.name">{{ errors.name }}</UFieldError>
                            </div>
                            <div class="field">
                                <label
                                    for="in-tagline"
                                    class="field-label"
                                >
                                    Krótki opis
                                </label>
                                <textarea
                                    id="in-tagline"
                                    v-model="form.tagline"
                                    class="textarea"
                                    rows="3"
                                    placeholder="Maślany, lekko kwaśny, gotowy w 30 minut."
                                />
                                <span class="field-hint">Jedno, dwa zdania pod tytułem.</span>
                            </div>
                        </div>
                    </div>
                </FormSection>

                <FormSection
                    id="f-ingredients"
                    n="02"
                    title="Składniki"
                    hint="Grupuj, jeśli przepis ma części (np. ciasto i krem). Kolejność zmienisz, przeciągając za uchwyt."
                    :error="errors.ingredients"
                >
                    <IngredientsEditor
                        v-model:groups="form.ingGroups"
                        :ingredients="ingredients"
                        :units="dicts?.units ?? []"
                        :row-errors="rowErrors"
                        @add-new="(name, item) => (addIng = { name, item })"
                    />
                </FormSection>

                <FormSection
                    id="f-steps"
                    n="03"
                    title="Kroki"
                    hint="Jeśli krok trzeba zrobić dużo wcześniej (marynata, namaczanie, wyrastanie), zaznacz „Przygotuj wcześniej”. Przepis dostanie wtedy oznaczenie."
                    :error="errors.steps"
                >
                    <StepsEditor v-model:groups="form.stepGroups" />
                </FormSection>

                <FormSection
                    n="04"
                    title="Wskazówka"
                    optional
                >
                    <textarea
                        v-model="form.tip"
                        class="textarea"
                        rows="3"
                        aria-label="Wskazówka"
                        placeholder="Jeśli filety są bardzo cienkie, skróć pieczenie do 10 minut…"
                    />
                </FormSection>
            </div>

            <aside class="upsert-side">
                <div class="side-card">
                    <h2 class="side-card__title">Klasyfikacja</h2>
                    <TaxonomyField
                        v-for="a in classificationAxes"
                        :key="a.id"
                        :axis="a"
                        :hint="
                            a.code === 'base'
                                ? 'Główne składniki, po których się szuka.'
                                : undefined
                        "
                        :error="errors[`taxonomy.${a.code}`]"
                        :model-value="axisModel(a.code)"
                        @update:model-value="setAxis(a.code, $event)"
                    />
                    <div
                        id="f-difficulty"
                        class="field"
                        role="group"
                        aria-labelledby="f-difficulty-label"
                    >
                        <span
                            id="f-difficulty-label"
                            class="field-label"
                        >
                            Trudność
                            <span
                                class="req"
                                aria-hidden="true"
                            >
                                *
                            </span>
                        </span>
                        <div class="diff-picker">
                            <button
                                v-for="(label, i) in ['łatwe', 'średnie', 'trudne']"
                                :key="label"
                                type="button"
                                class="diff-btn"
                                :class="{
                                    'is-active': form.difficulty === i + 1,
                                    'is-invalid': errors.difficulty,
                                }"
                                :aria-pressed="form.difficulty === i + 1"
                                @click="form.difficulty = i + 1"
                            >
                                <span
                                    class="diff-btn__dots"
                                    aria-hidden="true"
                                >
                                    <span
                                        v-for="k in 3"
                                        :key="k"
                                        class="diff-btn__dot"
                                        :class="{ 'is-on': k <= i + 1 }"
                                    />
                                </span>
                                {{ label }}
                            </button>
                        </div>
                        <UFieldError v-if="errors.difficulty">{{ errors.difficulty }}</UFieldError>
                    </div>
                </div>

                <div class="side-card">
                    <h2 class="side-card__title">Czas i porcje</h2>
                    <div class="side-card__two">
                        <div class="field">
                            <label
                                for="in-time"
                                class="field-label"
                            >
                                Czas (min)
                            </label>
                            <input
                                id="in-time"
                                v-model="form.time"
                                class="input input--mono"
                                type="number"
                                min="1"
                                placeholder="40"
                            />
                        </div>
                        <div
                            id="f-servings"
                            class="field"
                        >
                            <label
                                for="in-servings"
                                class="field-label"
                            >
                                Porcje
                                <span class="req">*</span>
                            </label>
                            <input
                                id="in-servings"
                                v-model="form.servings"
                                class="input input--mono"
                                :class="{ 'is-invalid': errors.servings }"
                                :aria-invalid="!!errors.servings || undefined"
                                type="number"
                                min="1"
                            />
                        </div>
                    </div>
                    <UFieldError v-if="errors.servings">{{ errors.servings }}</UFieldError>
                    <div class="field">
                        <label
                            for="in-kcal"
                            class="field-label"
                        >
                            Kcal / porcja
                        </label>
                        <input
                            id="in-kcal"
                            v-model="form.kcal"
                            class="input input--mono"
                            type="number"
                            min="1"
                            placeholder="opcjonalnie"
                        />
                    </div>
                </div>

                <div class="side-card">
                    <h2 class="side-card__title">Atrybuty</h2>
                    <TaxonomyField
                        v-if="dietAxis"
                        :axis="dietAxis"
                        strong
                        :model-value="axisModel('diet')"
                        @update:model-value="setAxis('diet', $event)"
                    />
                    <TaxonomyField
                        v-if="practicalAxis"
                        :axis="practicalAxis"
                        :model-value="axisModel('practical')"
                        @update:model-value="setAxis('practical', $event)"
                    />
                    <div class="prep-summary">
                        <template v-if="prepSummary.count">
                            <span class="prep-badge">
                                wymaga przygotowania{{
                                    prepSummary.lead ? ` · ${prepSummary.lead}` : ''
                                }}
                            </span>
                            <span class="field-hint">
                                Wynika z
                                {{
                                    prepSummary.count === 1
                                        ? '1 kroku'
                                        : `${prepSummary.count} kroków`
                                }}
                                oznaczonych „Przygotuj wcześniej”.
                            </span>
                        </template>
                        <span
                            v-else
                            class="field-hint"
                        >
                            Oznaczenie „wymaga przygotowania” pojawi się, gdy zaznaczysz taki krok.
                        </span>
                    </div>
                </div>

                <UFieldError v-if="serverError">{{ serverError }}</UFieldError>
            </aside>
        </div>

        <FormActionBar
            :status="status"
            :is-new="isNew"
            :dirty="dirty"
            :last-saved="lastSaved"
            :busy="busy"
            :error-count="errorCount"
            @cancel="cancel"
            @save-draft="saveDraft"
            @publish="publish"
            @save="saveChanges"
            @unpublish="unpublish"
            @show-errors="scrollToFirstError"
        />

        <IngredientModal
            v-if="addIng && dicts"
            :dicts="dicts"
            :initial-name="addIng.name"
            create-label="Dodaj i wstaw"
            @saved="onIngredientCreated"
            @close="addIng = null"
        />
    </div>
</template>

<style scoped>
.upsert-page {
    padding-top: 44px;
    padding-bottom: 120px;
}
.upsert-head .breadcrumb {
    margin-bottom: 14px;
}
.upsert-head__row {
    display: flex;
    align-items: flex-end;
    justify-content: space-between;
    gap: 24px;
}
.upsert-head__title {
    display: flex;
    align-items: center;
    gap: 16px;
}
.upsert-title {
    font-size: 52px;
    line-height: 1;
    letter-spacing: -2px;
    font-weight: 500;
}
.status-pill {
    display: inline-flex;
    align-items: center;
    gap: 7px;
    margin-top: 8px;
    padding: 5px 12px;
    border-radius: var(--r-pill);
    font-size: 13px;
    font-weight: 600;
}
.status-pill__dot {
    width: 7px;
    height: 7px;
    border-radius: var(--r-pill);
}
.status-pill--draft {
    background: var(--prep-soft);
    color: var(--prep);
}
.status-pill--draft .status-pill__dot {
    background: var(--c-mustard);
}
.status-pill--published {
    background: var(--success-soft);
    color: var(--success);
}
.status-pill--published .status-pill__dot {
    background: var(--success);
}
.upsert-loading {
    display: flex;
    justify-content: center;
    padding: 80px 0;
}
.upsert-body {
    display: grid;
    grid-template-columns: minmax(0, 1fr) 360px;
    gap: 40px;
    margin-top: 32px;
}
.upsert-main {
    display: flex;
    flex-direction: column;
    gap: 28px;
}
.basics {
    display: grid;
    grid-template-columns: 220px minmax(0, 1fr);
    gap: 24px;
    align-items: start;
}
.basics__fields {
    display: flex;
    flex-direction: column;
    gap: 18px;
}
.basics__name {
    font-size: 22px;
    font-weight: 500;
    letter-spacing: -0.4px;
    padding: 12px 14px;
}
.upsert-side {
    position: sticky;
    top: 84px;
    align-self: start;
    display: flex;
    flex-direction: column;
    gap: 16px;
}
.side-card {
    background: var(--bg-alt);
    border: 1px solid var(--rule);
    border-radius: var(--r-lg);
    padding: 18px;
    display: flex;
    flex-direction: column;
    gap: 16px;
}
.side-card__title {
    font-size: 14px;
    font-weight: 600;
    letter-spacing: 0;
    padding-bottom: 10px;
    border-bottom: 1px solid var(--rule);
}
.side-card__two {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 10px;
}
.diff-picker {
    display: flex;
    gap: 6px;
}
.diff-btn {
    flex: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 6px;
    padding: 10px 6px;
    border-radius: var(--r-md);
    border: 1px solid var(--rule-strong);
    background: transparent;
    color: var(--ink);
    font-size: 13px;
    font-weight: 500;
    cursor: pointer;
}
.diff-btn.is-invalid {
    border-color: var(--danger);
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
    border: 1px solid var(--rule-strong);
}
.diff-btn__dot.is-on {
    background: var(--accent);
    border-color: var(--accent);
}
.prep-summary {
    display: flex;
    flex-direction: column;
    gap: 4px;
    padding-top: 14px;
    border-top: 1px solid var(--rule);
}
.prep-badge {
    align-self: flex-start;
    font-size: 12px;
    font-weight: 600;
    color: var(--prep);
    background: var(--prep-soft);
    border: 1px solid var(--prep-line);
    padding: 4px 10px;
    border-radius: var(--r-pill);
}
</style>
