<script lang="ts" setup>
import { computed } from 'vue';

import DiceIcon from '@/components/DiceIcon.vue';
import { losParam, rangeText } from '@/shared/utils/challenges';

// Pasek „Wylosowane” nad wyzwaniem (strona otwarta z losowania, ?los=). „Zmień trudność” wraca na listę
// z tym wyborem w „Losuj”; „Losuj jeszcze raz” losuje z tej samej puli. Gdy się nie da — komunikat pod paskiem.
export type RerollState = 'none' | 'alldone' | 'onlycurrent' | 'single' | null;

const props = defineProps<{ sel: number[]; state: RerollState }>();
defineEmits<{ reroll: [includeDone: boolean]; anyDifficulty: [] }>();

const range = computed(() => rangeText(props.sel));
const scope = computed(() =>
    props.sel.length ? `wyzwania o trudności ${range.value}` : 'wyzwania'
);
</script>

<template>
    <div class="rolled">
        <div class="rolled__bar">
            <span class="rolled__what">
                <DiceIcon class="rolled__dice" />
                Wylosowane
                <span class="rolled__range">· trudność {{ range }}</span>
            </span>
            <RouterLink
                :to="{ path: '/challenges', query: sel.length ? { los: losParam(sel) } : {} }"
                class="rolled__change"
            >
                Zmień trudność
            </RouterLink>
            <button
                type="button"
                class="rolled__again"
                @click="$emit('reroll', false)"
            >
                <DiceIcon :size="14" />
                Losuj jeszcze raz
            </button>
        </div>
        <div aria-live="polite">
            <div
                v-if="state"
                class="rolled__msg"
            >
                <template v-if="state === 'none'">
                    <strong>Nie ma wyzwań o trudności {{ range }}.</strong>
                    <button
                        type="button"
                        class="rolled__alt"
                        @click="$emit('anyDifficulty')"
                    >
                        Losuj z dowolnej trudności
                    </button>
                </template>
                <template v-else-if="state === 'single'">
                    <strong>
                        To jedyne {{ sel.length ? `wyzwanie o trudności ${range}` : 'wyzwanie' }}.
                    </strong>
                </template>
                <template v-else>
                    <span>
                        <strong v-if="state === 'alldone'">
                            Wszystkie {{ scope }} są zrobione.
                        </strong>
                        <strong v-else>Pozostałe {{ scope }} są zrobione.</strong>
                        Możesz wylosować coś jeszcze raz.
                    </span>
                    <button
                        type="button"
                        class="rolled__alt"
                        @click="$emit('reroll', true)"
                    >
                        Losuj spośród zrobionych
                    </button>
                </template>
            </div>
        </div>
    </div>
</template>

<style scoped>
.rolled {
    display: flex;
    flex-direction: column;
    gap: 8px;
    margin-bottom: 24px;
}
.rolled__bar {
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    gap: 12px;
    padding: 10px 14px;
    background: var(--accent-soft);
    border-radius: var(--r-xl);
}
.rolled__what {
    flex: 1;
    min-width: 180px;
    display: inline-flex;
    align-items: center;
    gap: 8px;
    white-space: nowrap;
    font-size: 14px;
    font-weight: 600;
    color: var(--ink);
}
.rolled__dice {
    color: var(--accent);
}
.rolled__range {
    font-weight: 400;
    color: var(--ink-dim);
}
.rolled__change {
    font-size: 13px;
    color: var(--ink-dim);
    display: inline-flex;
    align-items: center;
}
.rolled__again {
    min-height: 34px;
    padding: 0 12px;
    border-radius: var(--r-md);
    border: none;
    background: var(--accent-fill);
    color: var(--on-accent);
    font-family: var(--font-sans);
    font-size: 13px;
    font-weight: 600;
    display: inline-flex;
    align-items: center;
    gap: 6px;
    cursor: pointer;
}
.rolled__msg {
    background: var(--info-soft);
    border-radius: var(--r-lg);
    padding: 10px 12px;
    font-size: 13px;
    line-height: 1.45;
    color: var(--ink-dim);
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
}
.rolled__msg strong {
    font-weight: 600;
    color: var(--ink);
}
.rolled__alt {
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
    .rolled {
        margin-bottom: 18px;
    }
    .rolled__bar {
        padding: 10px 12px;
    }
    .rolled__change,
    .rolled__again {
        min-height: 44px;
    }
}
</style>
