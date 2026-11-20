<script lang="ts" setup>
import { computed, ref, useId } from 'vue';
import { useRouter } from 'vue-router';

import DiceIcon from '@/components/DiceIcon.vue';
import FlameIcon from '@/components/FlameIcon.vue';
import type { IChallenge } from '@/data/dtos/challenge/ChallengeDto';
import {
    DIFFICULTIES,
    levelName,
    losParam,
    rangeText,
    rollChallenge,
} from '@/shared/utils/challenges';

// "Losuj wyzwanie": difficulty chips 1–5 (multiple choice, none = any) + "Losuj".
// The pick skips done challenges; the result → the challenge page with ?los= ("Wylosowane" bar).
// Empty pool / only done ones → a message in the card, no navigation.
const props = defineProps<{ items: IChallenge[]; done: ReadonlySet<number>; doneCount: number }>();
const sel = defineModel<number[]>({ required: true });

const router = useRouter();
const id = useId();
const state = ref<'none' | 'alldone' | null>(null);
const range = computed(() => rangeText(sel.value));

function toggle(n: number) {
    sel.value = sel.value.includes(n)
        ? sel.value.filter((x) => x !== n)
        : [...sel.value, n].sort((a, b) => a - b);
    state.value = null;
}

function roll(includeDone = false) {
    const res = rollChallenge(props.items, sel.value, props.done, { includeDone });
    if (res.kind === 'ok') {
        router.push({ path: `/challenges/${res.id}`, query: { los: losParam(sel.value) } });
        return;
    }
    state.value = res.kind === 'none' ? 'none' : 'alldone';
}

function anyDifficulty() {
    sel.value = [];
    state.value = null;
}
</script>

<template>
    <section
        class="rnd"
        :aria-labelledby="`${id}-h`"
    >
        <h2
            :id="`${id}-h`"
            class="rnd__title"
        >
            <DiceIcon
                :size="17"
                class="rnd__dice"
            />
            Losuj wyzwanie
        </h2>
        <div
            role="group"
            :aria-labelledby="`${id}-t`"
        >
            <div
                :id="`${id}-t`"
                class="rnd__label"
            >
                <span>Trudność</span>
                <span class="rnd__range">{{ range }}</span>
            </div>
            <div class="rnd__chips">
                <button
                    v-for="n in DIFFICULTIES"
                    :key="n"
                    type="button"
                    class="rnd__chip"
                    :class="{ 'is-on': sel.includes(n) }"
                    :aria-pressed="sel.includes(n)"
                    :title="levelName(n)"
                    @click="toggle(n)"
                >
                    <FlameIcon
                        :size="12"
                        :filled="sel.includes(n)"
                        class="rnd__flame"
                    />
                    {{ n }}
                    <span class="sr-only">— {{ levelName(n) }}</span>
                </button>
            </div>
        </div>
        <button
            type="button"
            class="rnd__go"
            @click="roll()"
        >
            <DiceIcon />
            Losuj
        </button>
        <div aria-live="polite">
            <p
                v-if="!state"
                class="rnd__hint"
            >
                Pomija zrobione ({{ doneCount }}). Bez wybranej trudności — z wszystkich.
            </p>
            <div
                v-else
                class="rnd__msg"
            >
                <span v-if="state === 'none'">
                    <strong>Nie ma wyzwań o trudności {{ range }}.</strong>
                </span>
                <span v-else>
                    <strong>
                        Wszystkie
                        {{ sel.length ? `wyzwania o trudności ${range}` : 'wyzwania' }} są zrobione.
                    </strong>
                    Możesz wylosować coś jeszcze raz.
                </span>
                <button
                    type="button"
                    class="rnd__alt"
                    @click="state === 'none' ? anyDifficulty() : roll(true)"
                >
                    {{
                        state === 'none' ? 'Losuj z dowolnej trudności' : 'Losuj spośród zrobionych'
                    }}
                </button>
            </div>
        </div>
    </section>
</template>

<style scoped>
.rnd {
    background: var(--bg-alt);
    border: 1px solid var(--rule);
    border-radius: var(--r-xl);
    padding: 18px;
    display: flex;
    flex-direction: column;
    gap: 12px;
}
.rnd__title {
    margin: 0;
    font-size: 18px;
    font-weight: 600;
    letter-spacing: -0.3px;
    display: flex;
    align-items: center;
    gap: 8px;
}
.rnd__dice {
    color: var(--accent);
}
.rnd__label {
    display: flex;
    justify-content: space-between;
    margin-bottom: 8px;
    font-family: var(--font-mono);
    font-size: 12px;
    font-weight: 600;
    letter-spacing: 0.5px;
    text-transform: uppercase;
    color: var(--ink-muted);
}
.rnd__range {
    font-family: var(--font-sans);
    font-weight: 500;
    letter-spacing: 0;
    text-transform: none;
}
.rnd__chips {
    display: grid;
    grid-template-columns: repeat(5, minmax(0, 1fr));
    gap: 6px;
}
.rnd__chip {
    min-width: 0;
    height: 36px;
    padding: 0;
    border-radius: var(--r-lg);
    border: 1px solid var(--rule-strong);
    background: var(--bg-alt);
    color: var(--ink-dim);
    font-family: var(--font-mono);
    font-size: 14px;
    font-weight: 600;
    cursor: pointer;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    gap: 3px;
}
.rnd__chip:hover {
    border-color: var(--accent);
}
.rnd__flame {
    color: var(--ink-faint);
}
.rnd__chip.is-on {
    border-color: var(--accent-fill);
    background: var(--accent-soft);
    color: var(--accent);
}
.rnd__chip.is-on .rnd__flame {
    color: var(--c-tomato);
}
.rnd__go {
    height: 42px;
    border-radius: var(--r-lg);
    border: none;
    background: var(--accent-fill);
    color: var(--on-accent);
    font-family: var(--font-sans);
    font-size: 15px;
    font-weight: 600;
    cursor: pointer;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
}
.rnd__go:hover {
    box-shadow: var(--shadow-accent);
}
.rnd__hint {
    margin: 0;
    font-size: 12px;
    line-height: 1.45;
    color: var(--ink-muted);
}
.rnd__msg {
    background: var(--info-soft);
    border-radius: var(--r-lg);
    padding: 10px 12px;
    font-size: 13px;
    line-height: 1.45;
    color: var(--ink-dim);
    display: flex;
    flex-direction: column;
    gap: 8px;
}
.rnd__msg strong {
    font-weight: 600;
    color: var(--ink);
}
.rnd__alt {
    align-self: flex-start;
    min-height: 36px;
    padding: 0 12px;
    border-radius: var(--r-md);
    border: 1px solid var(--rule-strong);
    background: var(--bg-alt);
    color: var(--ink);
    font-family: var(--font-sans);
    font-size: 13px;
    font-weight: 600;
    cursor: pointer;
}
@media (max-width: 900px) {
    .rnd {
        padding: 16px;
    }
    .rnd__chip {
        height: 44px;
    }
    .rnd__go {
        height: 48px;
    }
}
</style>
