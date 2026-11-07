<script lang="ts" setup>
import { useQueryClient } from '@tanstack/vue-query';
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue';
import { onBeforeRouteLeave, onBeforeRouteUpdate, useRouter } from 'vue-router';

import ConversionsTable from './components/ConversionsTable.vue';
import IngredientBasics from './components/IngredientBasics.vue';
import NutritionTable from './components/NutritionTable.vue';
import {
    convertibleUnits,
    formFromIngredient,
    type IIngredientForm,
    inheritedFromChildren,
    mapServerErrors,
    rowState,
    selfAndDescendants,
    toPayload,
    validateForm,
} from './ingredientForm';

import USpinner from '@/components/USpinner.vue';
import { confirmAction } from '@/composables/useConfirm';
import { showToast } from '@/composables/useToast';
import {
    createIngredient,
    fetchIngredientsWithClient,
    invalidateIngredients,
    updateIngredient,
} from '@/data/api/ingredients/fetch';
import { fetchDictionariesWithClient } from '@/data/api/taxonomy/fetch';
import type { IngredientDto } from '@/data/dtos/ingredients/IngredientDto';
import type { IDictionaries } from '@/data/dtos/taxonomy/DictionariesDto';
import { NUTRIENTS } from '@/shared/utils/nutrition';

// Edycja składnika jako osobna strona: podstawowe, makro na 100 g (z dziedziczeniem), przeliczniki jednostek.
const props = defineProps<{ id?: string }>();
const LIST = '/admin/ingredients';

const qc = useQueryClient();
const router = useRouter();

const all = ref<IngredientDto[]>([]);
const dicts = ref<IDictionaries>();
const loading = ref(true);
const notFound = ref(false);
const form = ref<IIngredientForm>(formFromIngredient());
const saved = ref('');
const errors = ref<Record<string, string>>({});
const serverError = ref('');
const saving = ref(false);
let leaving = false;

const ingredientId = computed(() => (props.id ? Number(props.id) : null));
const current = computed(() => all.value.find((i) => i.id === ingredientId.value));

// Przejście między składnikami (okruszki, baner) to ta sama trasa — wczytaj ponownie przy zmianie id.
async function load() {
    loading.value = true;
    [all.value, dicts.value] = await Promise.all([
        fetchIngredientsWithClient(qc),
        fetchDictionariesWithClient(qc),
    ]);
    notFound.value = ingredientId.value != null && !current.value;
    form.value = formFromIngredient(current.value);
    saved.value = JSON.stringify(form.value);
    errors.value = {};
    serverError.value = '';
    loading.value = false;
}
watch(ingredientId, load, { immediate: true });

const dirty = computed(() => !loading.value && JSON.stringify(form.value) !== saved.value);

const children = computed(() =>
    ingredientId.value == null ? [] : all.value.filter((i) => i.parentId === ingredientId.value)
);
const parents = computed(() => {
    const blocked =
        ingredientId.value == null
            ? new Set<number>()
            : selfAndDescendants(ingredientId.value, all.value);
    return all.value
        .filter((i) => !blocked.has(i.id))
        .sort((a, b) => a.name.localeCompare(b.name, 'pl'));
});
const parent = computed(() => all.value.find((i) => i.id === form.value.parentId));
const inherited = computed(() => inheritedFromChildren(children.value));
const ownCount = computed(
    () => NUTRIENTS.filter(({ key }) => !('empty' in rowState(form.value.nutrition[key]))).length
);
const units = computed(() => convertibleUnits(dicts.value?.units ?? []));
const inheritedConversions = computed(() =>
    (parent.value?.conversions ?? []).map((c) => ({ ...c, inherited: true }))
);

async function save() {
    errors.value = validateForm(form.value);
    serverError.value = '';
    if (Object.keys(errors.value).length) {
        serverError.value = 'Popraw zaznaczone pola.';
        return;
    }
    saving.value = true;
    try {
        const payload = toPayload(form.value);
        if (ingredientId.value != null) await updateIngredient(ingredientId.value, payload);
        else await createIngredient(payload);
        await invalidateIngredients(qc);
        showToast(ingredientId.value != null ? 'Składnik zapisany' : 'Składnik dodany');
        leaving = true;
        router.push(LIST);
    } catch (e) {
        const res = (
            e as {
                response?: {
                    status?: number;
                    data?: { fields?: { field: string; message: string }[] };
                };
            }
        ).response;
        if (res?.status === 422 && res.data?.fields) {
            errors.value = mapServerErrors(res.data.fields);
            serverError.value = 'Popraw zaznaczone pola.';
        } else {
            serverError.value =
                res?.status === 409
                    ? 'Taki składnik już jest w bazie.'
                    : 'Nie udało się zapisać składnika.';
        }
    } finally {
        saving.value = false;
    }
}

function cancel() {
    router.push(LIST);
}

// ─── Ostrzeżenie przy wyjściu ─────────────────────────────────
const onBeforeUnload = (e: BeforeUnloadEvent) => {
    if (!dirty.value) return;
    e.preventDefault();
    e.returnValue = '';
};
onMounted(() => window.addEventListener('beforeunload', onBeforeUnload));
onBeforeUnmount(() => window.removeEventListener('beforeunload', onBeforeUnload));

const guard = () => {
    if (leaving || !dirty.value) return true;
    return confirmAction({
        title: 'Wyjść bez zapisywania?',
        body: 'Zmiany w składniku przepadną.',
        confirmLabel: 'Wyjdź',
        danger: true,
    });
};
onBeforeRouteLeave(guard);
onBeforeRouteUpdate(guard);
</script>

<template>
    <nav
        class="crumbs"
        aria-label="Ścieżka"
    >
        <RouterLink :to="LIST">Składniki</RouterLink>
        <span aria-hidden="true">/</span>
        <template v-if="parent">
            <RouterLink :to="`${LIST}/${parent.id}`">{{ parent.name }}</RouterLink>
            <span aria-hidden="true">/</span>
        </template>
        <span class="crumbs__current">{{ current?.name ?? 'Nowy składnik' }}</span>
    </nav>
    <h1 class="admin-title">{{ ingredientId != null ? 'Edytuj składnik' : 'Nowy składnik' }}</h1>

    <div
        v-if="loading"
        class="admin-loading"
    >
        <USpinner :size="22" />
    </div>
    <p
        v-else-if="notFound"
        class="edit__missing"
    >
        Nie ma takiego składnika.
        <RouterLink :to="LIST">Wróć do listy</RouterLink>
    </p>
    <form
        v-else-if="dicts"
        class="edit"
        novalidate
        @submit.prevent="save"
    >
        <section class="edit__section">
            <div class="edit__head">
                <span class="edit__no">01</span>
                <h2 class="edit__title">Podstawowe</h2>
            </div>
            <IngredientBasics
                v-model="form"
                :dicts="dicts"
                :parents="parents"
                :child-count="children.length"
                :errors="errors"
            />
        </section>

        <section class="edit__section">
            <div class="edit__head">
                <span class="edit__no">02</span>
                <h2 class="edit__title">Wartości odżywcze na 100 g</h2>
                <span
                    v-if="inherited"
                    class="badge-attr edit__aside"
                >
                    {{ ownCount ? `${ownCount} z 5 własne` : 'dziedziczy' }}
                </span>
            </div>
            <p class="edit__hint">
                {{
                    inherited
                        ? 'Wartości kursywą w polach są odziedziczone. Wpisz własne, by nadpisać wiersz; max puste = jedna wartość.'
                        : 'Wpisz zakres, jeśli marki się różnią. Max puste = jedna wartość. Liczby z przecinkiem lub kropką.'
                }}
            </p>
            <div
                v-if="inherited"
                class="edit__banner"
            >
                <strong>Brak własnych danych — zakres wyliczony z podskładników</strong>
                (najniższe min, najwyższe max):
                <template
                    v-for="(c, i) in children"
                    :key="c.id"
                >
                    <RouterLink :to="`${LIST}/${c.id}`">{{ c.name }}</RouterLink>
                    {{ i < children.length - 1 ? ', ' : '.' }}
                </template>
                <span class="edit__banner-note">
                    Wpisz wartość w wierszu, żeby go nadpisać. Szeroki zakres sugeruje, że w
                    przepisie lepiej wybrać konkretny podskładnik.
                </span>
            </div>
            <NutritionTable
                v-model="form.nutrition"
                :inherited="inherited"
                :errors="errors"
            />
        </section>

        <section class="edit__section">
            <div class="edit__head">
                <span class="edit__no">03</span>
                <h2 class="edit__title">Przeliczniki jednostek</h2>
            </div>
            <p class="edit__hint">
                Ile gramów waży 1 jednostka tego składnika. Dzięki temu „2 szklanki” w przepisie
                wliczają się do makro. Gramy i kilogramy liczymy automatycznie, mililitry jako 1 ml
                = 1 g, jeśli nie podasz przelicznika.
            </p>
            <ConversionsTable
                v-model="form.conversions"
                :units="units"
                :inherited="inheritedConversions"
                :name="form.name"
                :errors="errors"
            />
        </section>

        <div class="edit__bar">
            <span
                v-if="serverError"
                class="edit__bar-error"
                role="alert"
            >
                {{ serverError }}
            </span>
            <button
                type="button"
                class="btn btn--ghost"
                @click="cancel"
            >
                Anuluj
            </button>
            <button
                type="submit"
                class="btn btn--accent"
                :class="{ 'btn--saving': saving }"
                :disabled="saving"
            >
                {{ saving ? 'Zapisywanie…' : 'Zapisz składnik' }}
            </button>
        </div>
    </form>
</template>

<style scoped>
.crumbs {
    display: flex;
    gap: 8px;
    align-items: center;
    font-size: 13px;
    color: var(--ink-muted);
}
.crumbs a {
    color: var(--ink-muted);
}
.crumbs__current {
    color: var(--ink);
    font-weight: 600;
}
.admin-title {
    font-size: 28px;
    font-weight: 600;
    letter-spacing: -0.8px;
    margin: 10px 0 22px;
}
.admin-loading {
    display: flex;
    justify-content: center;
    padding: 60px 0;
}
.edit {
    display: flex;
    flex-direction: column;
    gap: 16px;
    max-width: 860px;
}
.edit__section {
    background: var(--bg-alt);
    border: 1px solid var(--rule);
    border-radius: var(--r-lg);
    padding: 28px;
}
.edit__head {
    display: flex;
    align-items: baseline;
    gap: 12px;
    flex-wrap: wrap;
    margin-bottom: 20px;
}
.edit__head:has(+ .edit__hint) {
    margin-bottom: 6px;
}
.edit__no {
    font-family: var(--font-mono);
    font-size: 12px;
    font-weight: 600;
    color: var(--accent);
}
.edit__title {
    margin: 0;
    font-size: 20px;
    font-weight: 600;
    letter-spacing: -0.5px;
}
.edit__aside {
    margin-left: auto;
}
.edit__hint {
    margin: 0 0 20px;
    font-size: 13px;
    line-height: 1.45;
    color: var(--ink-muted);
}
.edit__banner {
    margin-bottom: 16px;
    padding: 12px 14px;
    background: var(--info-soft);
    border-radius: var(--r-lg);
    font-size: 13px;
    line-height: 1.5;
    color: var(--ink-dim);
}
.edit__banner strong {
    color: var(--ink);
    font-weight: 600;
}
.edit__banner a {
    color: var(--ink);
    text-underline-offset: 3px;
}
.edit__banner-note {
    display: block;
    margin-top: 4px;
}
.edit__missing {
    color: var(--ink-muted);
}
.edit__bar {
    position: sticky;
    bottom: 0;
    display: flex;
    justify-content: flex-end;
    align-items: center;
    gap: 10px;
    padding: 12px 0;
    background: var(--bg);
    border-top: 1px solid var(--rule);
}
.edit__bar-error {
    margin-right: auto;
    font-size: 13px;
    font-weight: 500;
    color: var(--danger);
}
</style>
