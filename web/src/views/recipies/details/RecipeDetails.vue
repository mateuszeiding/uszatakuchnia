<script lang="ts" setup>
import { useQueryClient } from '@tanstack/vue-query';
import { computed, ref, watch } from 'vue';

import RecipeBody from './components/RecipeBody.vue';
import RecipeBodySkeleton from './components/RecipeBodySkeleton.vue';
import RecipeLoadError from './components/RecipeLoadError.vue';
import RecipeMissing from './components/RecipeMissing.vue';

import { useDraftBar } from '@/composables/useDraftBar';
import { fetchRecipesWithClient } from '@/data/api/recipe/fetch';
import { fetchDictionariesWithClient } from '@/data/api/taxonomy/fetch';
import type { RecipeDto } from '@/data/dtos/recipe/RecipeDto';
import type { IDictionaries } from '@/data/dtos/taxonomy/DictionariesDto';
import { setTitle } from '@/shared/utils/title';

const props = defineProps<{ id: string | number }>();

const recipe = ref<RecipeDto>();
const dicts = ref<IDictionaries>();
const loading = ref(true);
/** missing = API 404 (doesn't exist, deleted, a draft for a non-owner); error = network / 5xx / timeout. */
const failed = ref<'missing' | 'error' | null>(null);
const retrying = ref(false);

const qc = useQueryClient();
// the bar only above the loaded current recipe (not the previous one while loading, nor on error)
useDraftBar(
    'recipe',
    computed(() =>
        !loading.value && !failed.value && recipe.value?.id === Number(props.id)
            ? recipe.value
            : undefined
    )
);

const isNotFound = (e: unknown) =>
    (e as { response?: { status?: number } } | null)?.response?.status === 404;

async function fetchAll(id: number) {
    try {
        [recipe.value, dicts.value] = await Promise.all([
            fetchRecipesWithClient(qc, id),
            fetchDictionariesWithClient(qc),
        ]);
        failed.value = null;
        setTitle(recipe.value.name);
    } catch (e) {
        recipe.value = undefined;
        failed.value = isNotFound(e) ? 'missing' : 'error';
    }
}

async function load(id: number) {
    loading.value = true;
    failed.value = null;
    await fetchAll(id);
    loading.value = false;
}

/** "Spróbuj ponownie": in place; the error page stays until the response ("Wczytuję…"). */
async function retry() {
    retrying.value = true;
    await fetchAll(Number(props.id));
    retrying.value = false;
}

watch(
    () => Number(props.id),
    (id) => load(id),
    { immediate: true }
);
</script>

<template>
    <RecipeMissing v-if="failed === 'missing'" />
    <RecipeLoadError
        v-else-if="failed === 'error'"
        :retrying="retrying"
        @retry="retry"
    />
    <div
        v-else
        class="container recipe-page"
    >
        <RouterLink
            to="/recipes"
            class="recipe-page__back"
        >
            ← Wszystkie przepisy
        </RouterLink>

        <RecipeBodySkeleton v-if="loading || !recipe" />
        <RecipeBody
            v-else
            :key="recipe.id"
            :recipe="recipe"
            :dicts="dicts"
        />

        <footer class="recipe-page__footer">
            <RouterLink to="/recipes">← Wróć do listy przepisów</RouterLink>
        </footer>
    </div>
</template>

<style scoped>
.recipe-page__back {
    display: inline-block;
    padding-top: 28px;
    font-size: 13px;
    color: var(--ink-muted);
}
/* link below the recipe; the border and "Zaloguj" are in the page footer (UFooter) */
.recipe-page__footer {
    margin-top: 56px;
    font-size: 13px;
}
.recipe-page__footer a {
    color: var(--ink-muted);
}

@media (max-width: 900px) {
    .recipe-page__back {
        padding-top: 12px;
    }
    .recipe-page__footer {
        margin-top: 40px;
    }
}
</style>
