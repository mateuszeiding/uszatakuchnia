<script lang="ts" setup>
import { useQueryClient } from '@tanstack/vue-query';
import { computed, ref, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';

import ShoppingList from './components/ShoppingList.vue';
import ShoppingListSkeleton from './components/ShoppingListSkeleton.vue';

import { useDraftBar } from '@/composables/useDraftBar';
import { fetchRecipesWithClient } from '@/data/api/recipe/fetch';
import { fetchDictionariesWithClient } from '@/data/api/taxonomy/fetch';
import type { RecipeDto } from '@/data/dtos/recipe/RecipeDto';
import type { IDictionaries } from '@/data/dtos/taxonomy/DictionariesDto';
import { type SizeChoice, sizeToQuery } from '@/shared/utils/portion';
import { setTitle } from '@/shared/utils/title';

const props = defineProps<{ id: string | number }>();

const recipe = ref<RecipeDto>();
const dicts = ref<IDictionaries>();
const loading = ref(true);
const failed = ref(false);

const qc = useQueryClient();
const route = useRoute();
const router = useRouter();
// the bar only above the loaded list of the current recipe; it doesn't linger from the previous one while loading or on error
useDraftBar(
    'list',
    computed(() =>
        !loading.value && !failed.value && recipe.value?.id === Number(props.id)
            ? recipe.value
            : undefined
    )
);

/** Servings from the URL (?porcje=): a link from the panel or a shared one; null = take them from the recipe progress. */
const urlServings = computed(() => {
    const n = Number(route.query.porcje);
    return Number.isInteger(n) && n >= 1 && n <= 20 ? n : null;
});

async function load(id: number) {
    loading.value = true;
    failed.value = false;
    try {
        [recipe.value, dicts.value] = await Promise.all([
            fetchRecipesWithClient(qc, id),
            fetchDictionariesWithClient(qc),
        ]);
        setTitle(`Lista zakupów: ${recipe.value.name}`);
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

/** Substitutes from the URL (?zamienniki=): the link sender's variant; ShoppingList parses it (it knows the recipe). */
const urlSubs = computed(() => {
    const v = route.query.zamienniki;
    return (Array.isArray(v) ? v[0] : v) ?? null;
});

const onServings = (n: number) => router.replace({ query: { ...route.query, porcje: String(n) } });
// the size replaces its own params only (a form link drops ?sztuki=, and the other way round)
const onSize = (s: SizeChoice | null) => {
    const { sztuki: _n, rozmiar: _d, forma: _f, ...rest } = route.query;
    router.replace({ query: { ...rest, ...sizeToQuery(s) } });
};
const onSync = (q: Record<string, string>) => router.replace({ query: { ...route.query, ...q } });
</script>

<template>
    <div class="shop-page">
        <ShoppingListSkeleton
            v-if="loading"
            :recipe-id="Number(id)"
        />
        <div
            v-else-if="failed || !recipe"
            class="shop-page__state"
        >
            Nie udało się wczytać przepisu.
            <RouterLink :to="`/recipes/${id}`">← Przepis</RouterLink>
        </div>
        <ShoppingList
            v-else
            :key="recipe.id"
            :recipe="recipe"
            :dicts="dicts"
            :url-servings="urlServings"
            :url-subs="urlSubs"
            :query="route.query"
            @servings="onServings"
            @size="onSize"
            @sync="onSync"
        />
    </div>
</template>

<style lang="scss" scoped>
@use '@/scss/mixins/media' as m;

.shop-page__state {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 12px;
    padding: 80px 16px;
    color: var(--ink-muted);
    font-size: 14px;
}

@include m.print {
    :global(.topbar) {
        display: none;
    }
    :global(body) {
        background: #fff;
    }
    .shop-page {
        --bg: #fff;
        --bg-alt: #fff;
        --ink: #000;
        --ink-dim: #000;
        --ink-muted: #444;
        --rule: #bbb;
        --rule-strong: #000;
        --accent: #000;
        color: #000;
        background: #fff;
    }
}
</style>
