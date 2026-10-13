<script lang="ts" setup>
import { computed } from 'vue';

import { bySection, leadText } from '../recipeView';

import type { RecipeDto } from '@/data/dtos/recipe/RecipeDto';

const props = defineProps<{
    recipe: RecipeDto;
    done: number[];
    canReset: boolean;
}>();
const emit = defineEmits<{ toggle: [stepNo: number]; cook: []; reset: [] }>();

const sections = computed(() => bySection(props.recipe.steps));
const total = computed(() => props.recipe.steps.length);
const doneCount = computed(
    () => props.recipe.steps.filter((s) => props.done.includes(s.stepNo)).length
);
</script>

<template>
    <section aria-labelledby="h-steps">
        <div class="steps-head">
            <div class="steps-head__row">
                <h2
                    id="h-steps"
                    class="steps-head__title"
                >
                    Kroki
                </h2>
                <div class="steps-head__actions">
                    <span class="steps-head__count">{{ doneCount }} z {{ total }}</span>
                    <button
                        v-if="canReset"
                        type="button"
                        class="btn btn--ghost btn--sm"
                        @click="emit('reset')"
                    >
                        Zacznij od nowa
                    </button>
                    <button
                        v-if="total"
                        type="button"
                        class="btn btn--accent"
                        @click="emit('cook')"
                    >
                        Tryb gotowania
                    </button>
                </div>
            </div>
            <div
                class="steps-progress"
                role="progressbar"
                aria-label="Postęp kroków"
                aria-valuemin="0"
                :aria-valuemax="total"
                :aria-valuenow="doneCount"
            >
                <div
                    class="steps-progress__fill"
                    :style="{ width: `${total ? (doneCount / total) * 100 : 0}%` }"
                />
            </div>
        </div>

        <div class="steps">
            <div
                v-for="(sec, si) in sections"
                :key="si"
            >
                <h3
                    v-if="sec.label"
                    class="steps__group"
                >
                    {{ sec.label }}
                </h3>
                <ol class="steps__list">
                    <li
                        v-for="s in sec.items"
                        :id="`step-${s.stepNo}`"
                        :key="s.stepNo"
                        class="step"
                        :class="{ 'is-done': done.includes(s.stepNo), 'is-prep': s.prep }"
                    >
                        <button
                            type="button"
                            role="checkbox"
                            class="step__check"
                            :aria-checked="done.includes(s.stepNo)"
                            :aria-label="`Krok ${s.stepNo}: oznacz jako zrobiony`"
                            @click="emit('toggle', s.stepNo)"
                        >
                            <svg
                                v-if="done.includes(s.stepNo)"
                                width="20"
                                height="20"
                                viewBox="0 0 22 22"
                                aria-hidden="true"
                            >
                                <path
                                    d="M5 11.5L9.5 16L17 7"
                                    stroke="#fff"
                                    stroke-width="2.4"
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                />
                            </svg>
                            <template v-else>{{ s.stepNo }}</template>
                        </button>
                        <div class="step__body">
                            <div
                                v-if="s.prep"
                                class="step__prep"
                            >
                                Przygotuj wcześniej{{ s.prepLead ? ` · ${leadText(s)}` : '' }}
                            </div>
                            <p class="step__text">{{ s.text }}</p>
                        </div>
                    </li>
                </ol>
            </div>
        </div>

        <aside
            v-if="recipe.tip"
            class="tip-box"
        >
            <div class="tip-box__head">Wskazówka</div>
            <p class="tip-box__text">{{ recipe.tip }}</p>
        </aside>
    </section>
</template>

<style scoped>
.steps-head {
    margin-bottom: 22px;
}
.steps-head__row {
    display: flex;
    align-items: flex-end;
    justify-content: space-between;
    gap: 16px;
    flex-wrap: wrap;
}
.steps-head__title {
    font-size: 32px;
    font-weight: 500;
    letter-spacing: -1px;
}
.steps-head__actions {
    display: flex;
    align-items: center;
    gap: 12px;
    flex-wrap: wrap;
}
.steps-head__count {
    font-size: 13px;
    color: var(--ink-muted);
}
.steps-progress {
    margin-top: 18px;
    height: 4px;
    background: var(--rule);
    border-radius: var(--r-pill);
    overflow: hidden;
}
.steps-progress__fill {
    height: 100%;
    background: var(--accent);
    transition: width 0.25s;
}
.steps {
    display: flex;
    flex-direction: column;
    gap: 28px;
}
.steps__group {
    font-family: var(--font-mono);
    font-size: 12px;
    color: var(--ink-muted);
    letter-spacing: 0.5px;
    text-transform: uppercase;
    font-weight: 600;
    margin-bottom: 12px;
}
.steps__list {
    list-style: none;
    padding: 0;
    margin: 0;
    display: flex;
    flex-direction: column;
    gap: 12px;
}
.step {
    background: var(--bg-alt);
    border: 1px solid var(--rule);
    border-radius: var(--r-xl);
    padding: 18px 20px;
    display: grid;
    grid-template-columns: 48px minmax(0, 1fr);
    gap: 18px;
    align-items: start;
}
.step.is-prep {
    border-color: var(--prep-line);
}
.step.is-done {
    background: var(--bg);
    border-color: var(--rule);
}
.step__check {
    width: 48px;
    height: 48px;
    border-radius: var(--r-xl);
    background: var(--bg);
    border: 1.5px solid var(--rule-strong);
    color: var(--accent);
    font-size: 20px;
    font-weight: 600;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
}
.step.is-done .step__check {
    background: var(--accent);
    border-color: var(--accent);
}
.step__body {
    padding-top: 4px;
}
.step__prep {
    font-size: 12px;
    font-weight: 600;
    color: var(--prep);
    margin-bottom: 6px;
}
.step__text {
    font-size: 17px;
    line-height: 1.6;
    text-wrap: pretty;
}
.step.is-done .step__text {
    color: var(--ink-muted);
    text-decoration: line-through;
}
.tip-box {
    margin-top: 28px;
    padding: 20px 22px;
    background: var(--accent-soft);
    border-radius: var(--r-xl);
}
.tip-box__head {
    font-family: var(--font-mono);
    font-size: 12px;
    color: var(--accent);
    letter-spacing: 0.5px;
    text-transform: uppercase;
    font-weight: 600;
    margin-bottom: 6px;
}
.tip-box__text {
    font-size: 16px;
    line-height: 1.6;
    text-wrap: pretty;
}

@media (max-width: 900px) {
    .steps-head__actions .btn--accent {
        padding: 12px 18px;
    }
    .step {
        padding: 16px;
        gap: 14px;
        grid-template-columns: 44px minmax(0, 1fr);
    }
    .step__check {
        width: 44px;
        height: 44px;
    }
}
</style>
