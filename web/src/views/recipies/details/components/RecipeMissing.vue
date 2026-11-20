<script lang="ts" setup>
import { useQueryClient } from '@tanstack/vue-query';
import { onMounted, ref } from 'vue';

import { fetchRecipesWithClient } from '@/data/api/recipe/fetch';
import { mainDish, type RecipeBaseDto } from '@/data/dtos/recipe/RecipeDto';
import { formatTime } from '@/shared/utils/time';
import ErrorState from '@/views/errors/ErrorState.vue';

// API 404 (missing, deleted or a draft): "Nie ma takiego przepisu" + up to 3 newest published recipes.
const qc = useQueryClient();
const latest = ref<RecipeBaseDto[]>([]);
onMounted(async () => {
    try {
        latest.value = (await fetchRecipesWithClient(qc, 'list')).slice(0, 3);
    } catch {
        /* no list: just the "not found" page */
    }
});
const sub = (r: RecipeBaseDto) =>
    [mainDish(r)?.label, formatTime(r.timeMin, r.timeMax)].filter(Boolean).join(' · ');
</script>

<template>
    <ErrorState
        kicker="Przepis"
        title="Nie ma takiego przepisu"
    >
        Mógł zostać usunięty albo właśnie jest poprawiany. Odświeżanie nic nie zmieni — zajrzyj do
        listy przepisów.
        <template #actions>
            <RouterLink
                to="/recipes"
                class="btn btn--accent"
            >
                Wszystkie przepisy
            </RouterLink>
        </template>
        <template #more>
            <section
                v-if="latest.length"
                class="latest"
                aria-labelledby="h-latest"
            >
                <h2
                    id="h-latest"
                    class="latest__title"
                >
                    Najnowsze przepisy
                </h2>
                <ul class="latest__list">
                    <li
                        v-for="r in latest"
                        :key="r.id"
                    >
                        <RouterLink
                            :to="`/recipes/${r.id}`"
                            class="latest__row"
                        >
                            <span class="latest__text">
                                <span class="latest__name">{{ r.name }}</span>
                                <span
                                    v-if="sub(r)"
                                    class="latest__sub"
                                >
                                    {{ sub(r) }}
                                </span>
                            </span>
                            <span
                                class="latest__arrow"
                                aria-hidden="true"
                            >
                                →
                            </span>
                        </RouterLink>
                    </li>
                </ul>
            </section>
        </template>
    </ErrorState>
</template>

<style scoped>
.latest {
    margin-top: 56px;
}
.latest__title {
    margin: 0 0 10px;
    font-family: var(--font-mono);
    font-size: 12px;
    font-weight: 600;
    letter-spacing: 0.5px;
    text-transform: uppercase;
    color: var(--ink-muted);
}
.latest__list {
    list-style: none;
    margin: 0;
    padding: 0;
    display: flex;
    flex-direction: column;
    gap: 8px;
}
.latest__row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 12px;
    min-height: 56px;
    padding: 12px 14px;
    background: var(--bg-alt);
    border: 1px solid var(--rule);
    border-radius: var(--r-lg);
    color: var(--ink);
    text-decoration: none;
    transition: border-color var(--dur-fast) var(--ease);
}
.latest__text {
    display: flex;
    flex-direction: column;
    gap: 4px;
    min-width: 0;
}
.latest__name {
    font-size: 15px;
    font-weight: 600;
    transition: color var(--dur-fast) var(--ease);
}
.latest__sub {
    font-size: 13px;
    color: var(--ink-muted);
}
.latest__arrow {
    color: var(--ink-muted);
}
@media (hover: hover) {
    .latest__row:hover {
        border-color: var(--rule-strong);
    }
    .latest__row:hover .latest__name {
        color: var(--accent);
    }
}
.latest__row:focus-visible {
    border-color: var(--rule-strong);
}
.latest__row:focus-visible .latest__name {
    color: var(--accent);
}

@media (max-width: 900px) {
    .latest {
        margin-top: 40px;
    }
}
</style>
