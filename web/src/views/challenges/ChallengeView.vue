<script lang="ts" setup>
import { useQueryClient } from '@tanstack/vue-query';
import { computed, ref, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';

import ChallengeDifficulty from './components/ChallengeDifficulty.vue';
import ChallengeIngredients from './components/ChallengeIngredients.vue';
import ChallengeRolledBar, { type RerollState } from './components/ChallengeRolledBar.vue';
import ChallengeViewSkeleton from './components/ChallengeViewSkeleton.vue';

import CautionsBlock from '@/components/CautionsBlock.vue';
import { useChallengesDone } from '@/composables/useChallengesDone';
import { useNarrow } from '@/composables/useNarrow';
import { fetchChallengesWithClient, fetchChallengeWithClient } from '@/data/api/challenges/fetch';
import type { IChallenge } from '@/data/dtos/challenge/ChallengeDto';
import {
    losParam,
    paragraphs,
    parseLos,
    rollChallenge,
    timeLimitText,
} from '@/shared/utils/challenges';
import { setTitle } from '@/shared/utils/title';

// Challenge page. Coming from a random pick (?los=), the "Wylosowane" bar replaces "← Wszystkie wyzwania".
// Desktop: content + the sticky "Składniki" card on the right; without ingredients, a single column.
// The "Uwagi" block sits below the metadata (before "Oznacz jako zrobione" and the description).
// Phone: format, title, metadata, notes, description, ingredients, full-width "Oznacz jako zrobione", source.
const props = defineProps<{ id: string }>();

const route = useRoute();
const router = useRouter();
const qc = useQueryClient();
const narrow = useNarrow();
const { done, isDone, toggle } = useChallengesDone();

const challenge = ref<IChallenge>();
const loading = ref(true);
const failed = ref(false);
const reroll = ref<RerollState>(null);

const los = computed(() => parseLos(route.query.los));
const hasIngredients = computed(() => !!challenge.value?.ingredients.length);
const isOn = computed(() => !!challenge.value && isDone(challenge.value.id));

async function load(id: number) {
    reroll.value = null;
    failed.value = false;
    loading.value = !qc.getQueryData(['challenges/', id]);
    try {
        challenge.value = await fetchChallengeWithClient(qc, id);
        setTitle(challenge.value.title);
    } catch (e) {
        // unknown challenge → list
        if ((e as { response?: { status?: number } }).response?.status === 404) {
            router.replace('/challenges');
            return;
        }
        failed.value = true;
    } finally {
        loading.value = false;
    }
}
watch(
    () => Number(props.id),
    (id) => (Number.isInteger(id) && id > 0 ? load(id) : router.replace('/challenges')),
    { immediate: true }
);

// "Losuj jeszcze raz": the same pool, without done ones and the open one; the result replaces the history entry.
async function rollAgain(includeDone: boolean, sel = los.value ?? []) {
    const current = challenge.value?.id;
    let list: IChallenge[];
    try {
        list = await fetchChallengesWithClient(qc);
    } catch {
        failed.value = true;
        return;
    }
    const res = rollChallenge(list, sel, done.value, { includeDone, exclude: current });
    if (res.kind === 'ok') {
        router.replace({ path: `/challenges/${res.id}`, query: { los: losParam(sel) } });
        return;
    }
    const poolSize = list.filter((c) => !sel.length || sel.includes(c.difficulty)).length;
    reroll.value = res.kind === 'onlycurrent' && poolSize === 1 ? 'single' : res.kind;
}
</script>

<template>
    <div class="container ch-detail">
        <ChallengeRolledBar
            v-if="los"
            :sel="los"
            :state="reroll"
            @reroll="rollAgain"
            @any-difficulty="rollAgain(false, [])"
        />
        <RouterLink
            v-else
            to="/challenges"
            class="ch-detail__back"
        >
            ← Wszystkie wyzwania
        </RouterLink>

        <ChallengeViewSkeleton
            v-if="loading"
            :narrow="narrow"
        />
        <p
            v-else-if="failed || !challenge"
            class="ch-detail__state"
        >
            Nie udało się wczytać wyzwania.
        </p>
        <div
            v-else
            class="ch-detail__grid"
            :class="{ 'is-single': !hasIngredients }"
        >
            <div class="ch-detail__main">
                <span class="badge-attr ch-detail__format">{{ challenge.format }}</span>
                <h1 class="ch-detail__title">{{ challenge.title }}</h1>
                <dl class="ch-detail__meta">
                    <div class="ch-detail__meta-item">
                        <dt class="ch-detail__dt">Trudność</dt>
                        <dd>
                            <ChallengeDifficulty
                                :level="challenge.difficulty"
                                :size="18"
                            />
                        </dd>
                    </div>
                    <div class="ch-detail__meta-item">
                        <dt class="ch-detail__dt">Limit czasu</dt>
                        <dd class="ch-detail__time">{{ timeLimitText(challenge.timeLimitMin) }}</dd>
                    </div>
                </dl>
                <CautionsBlock :items="challenge.cautions" />
                <div v-if="!narrow">
                    <button
                        type="button"
                        class="ch-done"
                        :class="{ 'is-on': isOn }"
                        :aria-pressed="isOn"
                        @click="toggle(challenge.id)"
                    >
                        {{ isOn ? '✓ Zrobione' : 'Oznacz jako zrobione' }}
                    </button>
                </div>
                <section
                    class="ch-detail__desc"
                    aria-label="Zadanie"
                >
                    <p
                        v-for="(p, i) in paragraphs(challenge.description)"
                        :key="i"
                    >
                        {{ p }}
                    </p>
                </section>
                <ChallengeIngredients
                    v-if="narrow && hasIngredients"
                    :ingredients="challenge.ingredients"
                    :narrow="narrow"
                />
                <button
                    v-if="narrow"
                    type="button"
                    class="ch-done"
                    :class="{ 'is-on': isOn }"
                    :aria-pressed="isOn"
                    @click="toggle(challenge.id)"
                >
                    {{ isOn ? '✓ Zrobione' : 'Oznacz jako zrobione' }}
                </button>
                <p
                    v-if="challenge.source"
                    class="ch-detail__source"
                >
                    Źródło: {{ challenge.source }}
                </p>
            </div>
            <ChallengeIngredients
                v-if="!narrow && hasIngredients"
                class="ch-detail__side"
                :ingredients="challenge.ingredients"
                :narrow="narrow"
            />
        </div>
    </div>
</template>

<style scoped>
.ch-detail {
    padding-top: 24px;
    padding-bottom: 48px;
}
.ch-detail__back {
    display: inline-block;
    margin-bottom: 20px;
    font-size: 13px;
    color: var(--ink-muted);
}
.ch-detail__state {
    color: var(--ink-muted);
}
.ch-detail__grid {
    display: grid;
    grid-template-columns: minmax(0, 1fr) 340px;
    gap: 56px;
    align-items: start;
}
.ch-detail__grid.is-single {
    grid-template-columns: minmax(0, 1fr);
    max-width: 760px;
}
.ch-detail__main {
    display: flex;
    flex-direction: column;
    gap: 22px;
    min-width: 0;
}
.ch-detail__format {
    align-self: flex-start;
    max-width: 100%;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    color: var(--ink-dim);
}
.ch-detail__title {
    margin: 0;
    font-size: 48px;
    line-height: 1.02;
    letter-spacing: -1.8px;
    font-weight: 500;
    text-wrap: balance;
    overflow-wrap: anywhere;
}
.ch-detail__meta {
    margin: 0;
    padding-top: 16px;
    border-top: 1px solid var(--rule);
    display: flex;
    flex-wrap: wrap;
    gap: 36px;
}
.ch-detail__meta-item {
    display: flex;
    flex-direction: column;
    gap: 6px;
}
.ch-detail__meta-item dd {
    margin: 0;
}
.ch-detail__dt {
    font-family: var(--font-mono);
    font-size: 12px;
    font-weight: 600;
    letter-spacing: 0.5px;
    text-transform: uppercase;
    color: var(--ink-muted);
}
.ch-detail__time {
    font-family: var(--font-mono);
    font-size: 16px;
    font-weight: 600;
}
.ch-done {
    min-height: 40px;
    padding: 0 16px;
    border-radius: var(--r-lg);
    border: 1px solid var(--rule-strong);
    background: var(--bg-alt);
    color: var(--ink);
    font-family: var(--font-sans);
    font-size: 14px;
    font-weight: 600;
    cursor: pointer;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
}
.ch-done:hover {
    border-color: var(--accent);
}
.ch-done.is-on {
    background: var(--accent-fill);
    border-color: var(--accent-fill);
    color: var(--on-accent);
}
.ch-detail__desc {
    display: flex;
    flex-direction: column;
    gap: 14px;
}
.ch-detail__desc p {
    margin: 0;
    max-width: 680px;
    font-size: 18px;
    line-height: 1.6;
    color: var(--ink);
    text-wrap: pretty;
    white-space: pre-line;
}
.ch-detail__source {
    margin: 0;
    padding-top: 14px;
    border-top: 1px solid var(--rule);
    font-size: 13px;
    color: var(--ink-muted);
}
.ch-detail__side {
    position: sticky;
    top: 84px;
}

@media (max-width: 900px) {
    .ch-detail {
        padding-top: 16px;
        padding-bottom: 32px;
    }
    .ch-detail__back {
        margin-bottom: 14px;
    }
    .ch-detail__grid {
        grid-template-columns: minmax(0, 1fr);
        gap: 24px;
    }
    .ch-detail__main {
        gap: 18px;
    }
    .ch-detail__title {
        font-size: 32px;
        letter-spacing: -1px;
    }
    .ch-detail__meta {
        gap: 20px;
    }
    .ch-detail__desc p {
        font-size: 17px;
    }
    .ch-done {
        width: 100%;
        min-height: 48px;
    }
}
</style>
