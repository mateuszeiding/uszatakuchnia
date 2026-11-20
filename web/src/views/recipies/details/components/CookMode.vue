<script lang="ts" setup>
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue';

import { leadText } from '../recipeView';

import CookIngredientsSheet from './CookIngredientsSheet.vue';
import StepHints from './StepHints.vue';

import CautionsBlock from '@/components/CautionsBlock.vue';
import type { ISubChoice } from '@/composables/useCookingProgress';
import { useWakeLock } from '@/composables/useWakeLock';
import type { RecipeDto } from '@/data/dtos/recipe/RecipeDto';
import type { IDictionaries } from '@/data/dtos/taxonomy/DictionariesDto';
import type { SubsMap } from '@/shared/utils/substitutes';

// Cooking mode: one step at a time, big Wstecz / Dalej, ← → on the keyboard, the screen stays on.
const props = defineProps<{
    recipe: RecipeDto;
    dicts: IDictionaries | undefined;
    servings: number;
    done: number[];
    checkedIng: number[];
    /** Chosen substitutes, shared with the recipe page. */
    subs: SubsMap;
}>();
const emit = defineEmits<{
    close: [];
    'set-step': [stepNo: number, done: boolean];
    'toggle-ing': [sortOrder: number];
    sub: [sortOrder: number, choice: ISubChoice | null];
    servings: [n: number];
    'clear-subs': [];
}>();

const steps = computed(() => props.recipe.steps);
const firstOpen = steps.value.findIndex((s) => !props.done.includes(s.stepNo));
const i = ref(firstOpen === -1 ? 0 : firstOpen);
const showIng = ref(false);
const wake = useWakeLock();
const root = ref<HTMLElement>();
const stage = ref<HTMLElement>();
// a new step always starts at the top (otherwise, after scrolling a long step, the next one starts halfway down)
watch(i, () => stage.value?.scrollTo({ top: 0 }));

const step = computed(() => steps.value[i.value]!);
const last = computed(() => i.value === steps.value.length - 1);
const ingBtn = ref<HTMLButtonElement>();

// the "Sprzęt i składniki" sheet: after closing, focus returns to "Składniki"
async function closeIng() {
    showIng.value = false;
    await nextTick();
    ingBtn.value?.focus();
}

function next() {
    emit('set-step', step.value.stepNo, true);
    if (last.value) emit('close');
    else i.value++;
}
function prev() {
    if (i.value === 0) return;
    i.value--;
    emit('set-step', step.value.stepNo, false);
}

function trapTab(e: KeyboardEvent) {
    const list = [...(root.value?.querySelectorAll<HTMLElement>('button:not([disabled])') ?? [])];
    if (!list.length) return;
    const first = list[0]!;
    const lastEl = list[list.length - 1]!;
    if (e.shiftKey && document.activeElement === first) {
        e.preventDefault();
        lastEl.focus();
    } else if (!e.shiftKey && document.activeElement === lastEl) {
        e.preventDefault();
        first.focus();
    }
}

function onKey(e: KeyboardEvent) {
    if (showIng.value) return; // the ingredients sheet has its own handling (Esc, focus)
    if (e.key === 'Tab') trapTab(e);
    else if (e.key === 'Escape') emit('close');
    else if (e.key === 'ArrowRight') next();
    else if (e.key === 'ArrowLeft') prev();
}

let prevOverflow = '';
let prevFocus: HTMLElement | null = null;
onMounted(() => {
    prevOverflow = document.body.style.overflow;
    document.body.style.overflow = 'hidden';
    prevFocus = document.activeElement as HTMLElement | null;
    root.value?.querySelector<HTMLElement>('.cook__next')?.focus();
    document.addEventListener('keydown', onKey);
});
onBeforeUnmount(() => {
    document.body.style.overflow = prevOverflow;
    document.removeEventListener('keydown', onKey);
    prevFocus?.focus?.();
});
</script>

<template>
    <Teleport to="body">
        <div
            ref="root"
            class="cook"
            role="dialog"
            aria-modal="true"
            aria-label="Tryb gotowania"
        >
            <div class="cook__bar">
                <div class="cook__titles">
                    <div class="cook__recipe">{{ recipe.name }}</div>
                    <div class="cook__count">Krok {{ i + 1 }} z {{ steps.length }}</div>
                </div>
                <div class="cook__tools">
                    <button
                        ref="ingBtn"
                        type="button"
                        class="cook__pill"
                        aria-haspopup="dialog"
                        @click="showIng = true"
                    >
                        Składniki
                    </button>
                    <button
                        type="button"
                        class="cook__round"
                        aria-label="Zamknij tryb gotowania"
                        @click="emit('close')"
                    >
                        <svg
                            width="11"
                            height="11"
                            viewBox="0 0 11 11"
                            aria-hidden="true"
                        >
                            <path
                                d="M2 2l7 7M9 2L2 9"
                                stroke="currentColor"
                                stroke-width="1.5"
                                stroke-linecap="round"
                            />
                        </svg>
                    </button>
                </div>
            </div>
            <div
                class="cook__progress"
                aria-hidden="true"
            >
                <div
                    class="cook__progress-fill"
                    :style="{ width: `${((i + 1) / steps.length) * 100}%` }"
                />
            </div>
            <div
                ref="stage"
                class="cook__stage"
            >
                <div
                    class="cook__content"
                    aria-live="polite"
                >
                    <!-- notes only on the step 1 screen, above the step text -->
                    <CautionsBlock
                        v-if="i === 0"
                        class="cook__cautions"
                        :items="recipe.cautions"
                        variant="large"
                    />
                    <div
                        v-if="step.section"
                        class="cook__section"
                    >
                        {{ step.section }}
                    </div>
                    <div
                        v-if="step.prep"
                        class="cook__prep"
                    >
                        Przygotuj wcześniej{{ step.prepLead ? ` · ${leadText(step)}` : '' }}
                    </div>
                    <p class="cook__text">{{ step.text }}</p>
                    <StepHints
                        :hints="step.hints"
                        size="lg"
                    />
                </div>
            </div>
            <div class="cook__nav">
                <button
                    type="button"
                    class="cook__btn cook__btn--back"
                    :disabled="i === 0"
                    @click="prev"
                >
                    ← Wstecz
                </button>
                <button
                    type="button"
                    class="cook__btn cook__next"
                    @click="next"
                >
                    {{ last ? 'Gotowe' : 'Dalej →' }}
                </button>
            </div>
            <div
                v-if="wake"
                class="cook__hint"
            >
                Ekran nie zgaśnie, dopóki tryb jest otwarty
            </div>
            <div
                v-else
                class="cook__hint cook__hint--keys"
            >
                Strzałki ← → zmieniają krok
            </div>
        </div>

        <CookIngredientsSheet
            v-if="showIng"
            :recipe="recipe"
            :dicts="dicts"
            :servings="servings"
            :checked-ing="checkedIng"
            :subs="subs"
            :step-no="i + 1"
            @close="closeIng"
            @toggle-ing="emit('toggle-ing', $event)"
            @sub="(so, c) => emit('sub', so, c)"
            @servings="emit('servings', $event)"
            @clear-subs="emit('clear-subs')"
        />
    </Teleport>
</template>

<style scoped>
.cook {
    position: fixed;
    inset: 0;
    z-index: var(--z-overlay);
    background: var(--bg);
    color: var(--ink);
    display: flex;
    flex-direction: column;
}
.cook__bar {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 12px;
    padding: 14px 20px;
    border-bottom: 1px solid var(--rule);
}
.cook__titles {
    min-width: 0;
}
.cook__recipe {
    font-size: 13px;
    color: var(--ink-muted);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}
.cook__count {
    font-size: 16px;
    font-weight: 600;
}
.cook__tools {
    display: flex;
    align-items: center;
    gap: 8px;
}
.cook__pill {
    height: 44px;
    padding: 0 16px;
    border-radius: var(--r-pill);
    border: 1px solid var(--rule-strong);
    background: var(--bg-alt);
    color: var(--ink);
    font-size: 15px;
    font-weight: 600;
    cursor: pointer;
}
.cook__round {
    width: 44px;
    height: 44px;
    border-radius: var(--r-pill);
    border: 1px solid var(--rule-strong);
    background: var(--bg-alt);
    color: var(--ink);
    cursor: pointer;
    display: inline-flex;
    align-items: center;
    justify-content: center;
}
.cook__progress {
    height: 4px;
    background: var(--rule);
}
.cook__progress-fill {
    height: 100%;
    background: var(--accent);
    transition: width 0.25s;
}
.cook__stage {
    flex: 1;
    overflow-y: auto;
    scrollbar-gutter: stable;
    display: flex;
    justify-content: center;
    padding: 32px 24px;
}
/* vertically centered with auto margins: when the content (e.g. notes on step 1) is taller
   than the screen, it scrolls from the top instead of escaping past the edge */
.cook__content {
    max-width: 760px;
    width: 100%;
    margin-block: auto;
}
.cook__cautions {
    margin-bottom: 22px;
}
.cook__section {
    font-family: var(--font-mono);
    font-size: 13px;
    color: var(--ink-muted);
    letter-spacing: 0.5px;
    text-transform: uppercase;
    font-weight: 600;
    margin-bottom: 14px;
}
.cook__prep {
    font-size: 14px;
    font-weight: 600;
    color: var(--prep);
    margin-bottom: 12px;
}
.cook__text {
    font-size: clamp(24px, 3.4vw, 40px);
    line-height: 1.4;
    font-weight: 500;
    letter-spacing: -0.4px;
    text-wrap: pretty;
}
.cook__nav {
    display: flex;
    gap: 12px;
    padding: 16px 20px calc(16px + env(safe-area-inset-bottom));
    border-top: 1px solid var(--rule);
    max-width: 800px;
    width: 100%;
    margin: 0 auto;
}
.cook__btn {
    flex: 2;
    min-height: 64px;
    border-radius: var(--r-2xl);
    font-size: 18px;
    font-weight: 600;
    cursor: pointer;
    background: var(--accent-fill);
    color: var(--on-accent);
    border: none;
}
.cook__btn--back {
    flex: 1;
    background: var(--bg-alt);
    color: var(--ink);
    border: 1px solid var(--rule-strong);
}
.cook__btn--back:disabled {
    opacity: 0.4;
    cursor: default;
}
.cook__hint {
    text-align: center;
    font-size: 12px;
    color: var(--ink-muted);
    padding-bottom: 10px;
}
/* the arrow-key hint only where there's a keyboard and a pointer (noise on a phone) */
.cook__hint--keys {
    display: none;
}
@media (hover: hover) and (pointer: fine) {
    .cook__hint--keys {
        display: block;
    }
}
</style>
