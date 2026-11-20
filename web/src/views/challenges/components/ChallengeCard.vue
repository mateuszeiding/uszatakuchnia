<script lang="ts" setup>
import { computed } from 'vue';

import ChallengeDifficulty from './ChallengeDifficulty.vue';
import DoneToggle from './DoneToggle.vue';

import type { IChallenge } from '@/data/dtos/challenge/ChallengeDto';
import { paragraphs, timeLimitText } from '@/shared/utils/challenges';

// Karta wyzwania: <li> = link (cała karta) + rodzeństwo — przycisk „zrobione” w prawym górnym rogu
// (osobny cel dla myszy i Tab). Zrobione zostaje na miejscu: tło --bg, tytuł --ink-dim, „✓ zrobione”.
// Tytuł i opis mają zarezerwowane 2 linie — karty tej samej wysokości (szkielet ma ten sam układ).
const props = defineProps<{ challenge: IChallenge; done: boolean }>();
defineEmits<{ toggle: [] }>();

const lead = computed(() => paragraphs(props.challenge.description)[0] ?? '');
</script>

<template>
    <li
        class="ch-card"
        :class="{ 'is-done': done }"
    >
        <RouterLink
            :to="`/challenges/${challenge.id}`"
            class="ch-card__link"
        >
            <span class="ch-card__format">
                <span class="badge-attr ch-card__badge">{{ challenge.format }}</span>
            </span>
            <span class="ch-card__title">{{ challenge.title }}</span>
            <span class="ch-card__desc">{{ lead }}</span>
            <span class="ch-card__foot">
                <ChallengeDifficulty
                    :level="challenge.difficulty"
                    :size="15"
                    compact
                />
                <span class="ch-card__time">
                    <svg
                        width="13"
                        height="13"
                        viewBox="0 0 12 12"
                        fill="none"
                        aria-hidden="true"
                    >
                        <circle
                            cx="6"
                            cy="6"
                            r="4.8"
                            stroke="currentColor"
                            stroke-width="1.2"
                        />
                        <path
                            d="M6 3.5v2.7l1.7 1"
                            stroke="currentColor"
                            stroke-width="1.2"
                            stroke-linecap="round"
                        />
                    </svg>
                    <span class="sr-only">Limit czasu:</span>
                    {{ timeLimitText(challenge.timeLimitMin) }}
                </span>
                <!-- miejsce zarezerwowane zawsze: przełączenie nie zmienia wysokości karty -->
                <span
                    class="ch-card__done"
                    :class="{ 'is-hidden': !done }"
                >
                    ✓ zrobione
                </span>
            </span>
        </RouterLink>
        <span class="ch-card__toggle">
            <DoneToggle
                :on="done"
                :title="challenge.title"
                @toggle="$emit('toggle')"
            />
        </span>
    </li>
</template>

<style scoped>
.ch-card {
    position: relative;
    list-style: none;
}
.ch-card__link {
    display: flex;
    flex-direction: column;
    gap: 10px;
    height: 100%;
    padding: 16px 18px 14px;
    background: var(--bg-alt);
    border: 1px solid var(--rule);
    border-radius: var(--r-xl);
    color: var(--ink);
    text-decoration: none;
    transition: border-color var(--dur-fast) var(--ease);
}
@media (hover: hover) {
    .ch-card__link:hover {
        border-color: var(--rule-strong);
    }
    .ch-card__link:hover .ch-card__title {
        color: var(--accent);
    }
}
.ch-card__link:focus-visible {
    border-color: var(--rule-strong);
}
.ch-card__link:focus-visible .ch-card__title {
    color: var(--accent);
}
.is-done .ch-card__link {
    background: var(--bg);
}
/* miejsce na przycisk „zrobione” w rogu */
.ch-card__format {
    display: flex;
    padding-right: 40px;
    min-width: 0;
}
.ch-card__badge {
    max-width: 100%;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    color: var(--ink-dim);
}
.ch-card__title {
    font-size: 18px;
    font-weight: 600;
    letter-spacing: -0.4px;
    line-height: 1.25;
    min-height: calc(2 * 1.25em);
    transition: color var(--dur-fast) var(--ease);
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
}
.is-done .ch-card__title {
    color: var(--ink-dim);
}
.ch-card__desc {
    flex: 1;
    font-size: 14px;
    line-height: 1.5;
    min-height: calc(2 * 1.5em);
    color: var(--ink-muted);
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
}
.ch-card__foot {
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    gap: 12px;
    padding-top: 10px;
    border-top: 1px solid var(--rule);
    font-size: 13px;
    color: var(--ink-dim);
}
.ch-card__time {
    display: inline-flex;
    align-items: center;
    gap: 5px;
    font-family: var(--font-mono);
}
.ch-card__done {
    margin-left: auto;
    font-size: 12px;
    font-weight: 600;
    color: var(--success);
}
.ch-card__done.is-hidden {
    visibility: hidden;
}
.ch-card__toggle {
    position: absolute;
    top: 14px;
    right: 14px;
}
</style>
