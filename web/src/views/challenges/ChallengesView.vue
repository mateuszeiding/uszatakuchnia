<script lang="ts" setup>
import { useQueryClient } from '@tanstack/vue-query';
import { computed, onMounted, ref, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';

import ChallengeCard from './components/ChallengeCard.vue';
import ChallengeCardSkeleton from './components/ChallengeCardSkeleton.vue';
import ChallengeRandomizer from './components/ChallengeRandomizer.vue';

import FlameIcon from '@/components/FlameIcon.vue';
import { useChallengesDone } from '@/composables/useChallengesDone';
import { fetchChallengesWithClient } from '@/data/api/challenges/fetch';
import type { IChallenge } from '@/data/dtos/challenge/ChallengeDto';
import { losParam, parseLos, plural } from '@/shared/utils/challenges';

// Challenge list: desktop, a column with "Losuj" (sticky) + a 2-column grid;
// phone, "Losuj" above the list, 1 column. Newest first. While loading, the same layout with skeletons.
const route = useRoute();
const router = useRouter();
const qc = useQueryClient();
const { done, isDone, toggle } = useChallengesDone();

const items = ref<IChallenge[]>([]);
const loading = ref(true);
const failed = ref(false);
// "Zmień trudność" on the challenge page arrives with ?los=; the choice goes into the chips.
const sel = ref<number[]>(parseLos(route.query.los) ?? []);

onMounted(async () => {
    try {
        items.value = await fetchChallengesWithClient(qc);
    } catch {
        failed.value = true;
    } finally {
        loading.value = false;
    }
});

// difficulty choice in the URL (replace, so the history stays clean)
watch(sel, (s) => router.replace({ query: s.length ? { los: losParam(s) } : {} }));

// "done" entries of challenges that no longer exist don't count
const doneCount = computed(() => items.value.filter((c) => done.value.has(c.id)).length);

// Skeletons: as many cards as last time (max 10), so the grid has a similar height after loading.
const SK_KEY = 'uk:challenges-skeletons';
const skeletons = (() => {
    try {
        const n = Number(localStorage.getItem(SK_KEY));
        return Number.isInteger(n) && n >= 1 && n <= 10 ? n : 6;
    } catch {
        return 6;
    }
})();
watch(loading, (l) => {
    if (l || !items.value.length) return;
    try {
        localStorage.setItem(SK_KEY, String(Math.min(items.value.length, 10)));
    } catch {
        /* no localStorage: the default stays */
    }
});
</script>

<template>
    <div class="container ch-page">
        <h1 class="ch-page__title">Wyzwania</h1>
        <p class="ch-page__lead">
            Zadania kulinarne do treningu. Wylosuj jedno albo wybierz z listy.
        </p>

        <div
            v-if="!loading && !failed && !items.length"
            class="ch-empty"
        >
            <span class="ch-empty__icon"><FlameIcon :size="20" /></span>
            <div class="ch-empty__title">Jeszcze nie ma wyzwań</div>
            <p class="ch-empty__text">
                Pierwsze pojawią się wkrótce. Do tego czasu —
                <!-- inline-flex: no space before the period -->
                <span class="ch-empty__link">
                    <RouterLink to="/recipes">przepisy</RouterLink>
                    <span>.</span>
                </span>
            </p>
        </div>

        <div
            v-else
            class="ch-layout"
        >
            <aside
                class="ch-aside"
                :inert="loading"
            >
                <!-- the card is in place while loading, just disabled -->
                <ChallengeRandomizer
                    v-model="sel"
                    :items="items"
                    :done="done"
                    :done-count="doneCount"
                />
            </aside>

            <div>
                <div
                    class="ch-results"
                    role="status"
                >
                    <template v-if="loading">Wczytywanie…</template>
                    <template v-else-if="failed">Nie udało się wczytać wyzwań.</template>
                    <template v-else>
                        <span>
                            {{ items.length }}
                            {{ plural(items.length, 'wyzwanie', 'wyzwania', 'wyzwań') }}
                        </span>
                        <span v-if="doneCount">
                            · {{ doneCount }}
                            {{ plural(doneCount, 'zrobione', 'zrobione', 'zrobionych') }}
                        </span>
                    </template>
                </div>
                <ul
                    v-if="loading"
                    class="ch-grid"
                    aria-busy="true"
                >
                    <li
                        v-for="n in skeletons"
                        :key="n"
                    >
                        <ChallengeCardSkeleton />
                    </li>
                </ul>
                <ul
                    v-else
                    class="ch-grid"
                >
                    <ChallengeCard
                        v-for="c in items"
                        :key="c.id"
                        :challenge="c"
                        :done="isDone(c.id)"
                        @toggle="toggle(c.id)"
                    />
                </ul>
            </div>
        </div>
    </div>
</template>

<style scoped>
.ch-page {
    padding-bottom: 48px;
}
.ch-page__title {
    font-size: 56px;
    line-height: 0.98;
    letter-spacing: -2.2px;
    font-weight: 500;
    margin: 40px 0 0;
}
.ch-page__lead {
    margin: 10px 0 0;
    max-width: 560px;
    font-size: 17px;
    line-height: 1.5;
    color: var(--ink-dim);
}
.ch-layout {
    display: grid;
    grid-template-columns: 260px minmax(0, 1fr);
    gap: 40px;
    align-items: start;
    margin-top: 32px;
}
.ch-aside {
    position: sticky;
    top: 84px;
    display: flex;
    flex-direction: column;
    gap: 14px;
}
.ch-results {
    display: flex;
    align-items: center;
    gap: 12px;
    margin: 0 0 14px;
    font-size: 14px;
    color: var(--ink-muted);
}
.ch-grid {
    list-style: none;
    margin: 0;
    padding: 0;
    display: grid;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: 16px;
}
.ch-empty {
    margin-top: 28px;
    background: var(--bg-alt);
    border: 1px dashed var(--rule-strong);
    border-radius: var(--r-xl);
    padding: 48px 24px;
    text-align: center;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 10px;
}
.ch-empty__icon {
    width: 44px;
    height: 44px;
    border-radius: var(--r-pill);
    background: var(--accent-soft);
    color: var(--accent);
    display: inline-flex;
    align-items: center;
    justify-content: center;
}
.ch-empty__title {
    font-size: 18px;
    font-weight: 600;
}
.ch-empty__text {
    margin: 0;
    max-width: 340px;
    font-size: 14px;
    line-height: 1.5;
    color: var(--ink-muted);
}
.ch-empty__link {
    display: inline-flex;
}
.ch-empty__text a {
    color: var(--accent);
}

@media (max-width: 900px) {
    .ch-page {
        padding-bottom: 32px;
    }
    .ch-page__title {
        font-size: 40px;
        letter-spacing: -1.4px;
        margin-top: 24px;
    }
    .ch-page__lead {
        font-size: 15px;
    }
    .ch-layout {
        grid-template-columns: minmax(0, 1fr);
        gap: 16px;
        margin-top: 22px;
    }
    .ch-aside {
        position: static;
    }
    .ch-grid {
        grid-template-columns: minmax(0, 1fr);
        gap: 12px;
    }
}
</style>
