<script lang="ts" setup>
import { useQueryClient } from '@tanstack/vue-query';
import { ref, watch } from 'vue';

import RecipeBody from './components/RecipeBody.vue';

import USpinner from '@/components/USpinner.vue';
import { fetchRecipesWithClient } from '@/data/api/recipe/fetch';
import { fetchDictionariesWithClient } from '@/data/api/taxonomy/fetch';
import type { RecipeDto } from '@/data/dtos/recipe/RecipeDto';
import type { IDictionaries } from '@/data/dtos/taxonomy/DictionariesDto';

const props = defineProps<{ id: string | number }>();

const recipe = ref<RecipeDto>();
const dicts = ref<IDictionaries>();
const loading = ref(true);
const failed = ref(false);

const qc = useQueryClient();

async function load(id: number) {
    loading.value = true;
    failed.value = false;
    try {
        [recipe.value, dicts.value] = await Promise.all([
            fetchRecipesWithClient(qc, id),
            fetchDictionariesWithClient(qc),
        ]);
        document.title = `${recipe.value.name} · uszatakuchnia`;
    } catch {
        failed.value = true;
    } finally {
        loading.value = false;
    }
}

watch(
    () => Number(props.id),
    (id) => load(id),
    { immediate: true }
);
</script>

<template>
    <div class="container recipe-page">
        <RouterLink
            to="/recipes"
            class="recipe-page__back"
        >
            ← Wszystkie przepisy
        </RouterLink>

        <div
            v-if="loading"
            class="recipe-page__state"
        >
            <USpinner :size="28" />
            <span>ładowanie przepisu…</span>
        </div>
        <div
            v-else-if="failed || !recipe"
            class="recipe-page__state"
        >
            Nie udało się wczytać przepisu.
        </div>
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
.recipe-page__state {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 12px;
    padding: 80px 0;
    color: var(--ink-muted);
    font-size: 14px;
}
.recipe-page__footer {
    margin-top: 80px;
    padding: 40px 0;
    border-top: 1px solid var(--rule);
    font-size: 13px;
}
.recipe-page__footer a {
    color: var(--ink-muted);
}

@media (max-width: 900px) {
    .recipe-page__back {
        padding-top: 16px;
    }
    .recipe-page__footer {
        padding: 32px 0 40px;
    }
}
</style>
