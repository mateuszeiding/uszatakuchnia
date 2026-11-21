<script lang="ts" setup>
import { computed, nextTick, onBeforeUnmount, onMounted, provide, ref, watch } from 'vue';

import { useCookFit } from '../cookFit';
import { leadText, prepRows } from '../recipeView';

import CookBar from './CookBar.vue';
import CookBefore from './CookBefore.vue';
import CookHints from './CookHints.vue';
import CookIngredientsSheet from './CookIngredientsSheet.vue';
import CookNav from './CookNav.vue';

import StepText from '@/components/StepText.vue';
import type { ISubChoice } from '@/composables/useCookingProgress';
import { COOK_TERMS } from '@/composables/useTermExplain';
import { useWakeLock } from '@/composables/useWakeLock';
import type { RecipeDto } from '@/data/dtos/recipe/RecipeDto';
import type { IDictionaries } from '@/data/dtos/taxonomy/DictionariesDto';
import { resumeIndex } from '@/shared/utils/progress';
import type { IRefContext } from '@/shared/utils/stepRefs';
import type { SubsMap } from '@/shared/utils/substitutes';

const props = defineProps<{
    recipe: RecipeDto;
    refs: IRefContext;
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
const prep = computed(() => prepRows(steps.value));
const hasBefore = computed(() => props.recipe.cautions.length > 0 || prep.value.length > 0);

// -1 = "Zanim zaczniesz". Resuming (some steps checked) skips it and opens the step after the last checked one.
const resumeAt = resumeIndex(steps.value, props.done);
const i = ref(resumeAt >= 0 ? resumeAt : hasBefore.value ? -1 : 0);

const showIng = ref(false);
/** The reader opened collapsed hints once: they stay open until the mode closes. */
const hintsOpen = ref(false);
useWakeLock();
// terms in the step open the big sheet with "Zamknij"
provide(COOK_TERMS, true);
const root = ref<HTMLElement>();
const stage = ref<HTMLElement>();
const bar = ref<InstanceType<typeof CookBar>>();

const { layout, size, tight, measuring, more, fit, updateMore } = useCookFit(stage);

const step = computed(() => (i.value >= 0 ? steps.value[i.value]! : null));
const last = computed(() => i.value === steps.value.length - 1);
const collapsible = computed(() => !measuring.value && tight.value);

// a new step always starts at the top (otherwise, after scrolling a long step, the next one starts halfway down)
watch(i, async () => {
    stage.value?.scrollTo({ top: 0 });
    await nextTick();
    void fit();
});
// a different substitute or servings change the text of the same step, so fit it again
watch(
    () => props.refs,
    async () => {
        await nextTick();
        void fit();
    }
);
watch(hintsOpen, async () => {
    await nextTick();
    updateMore();
});

// the "Sprzęt i składniki" sheet: after closing, focus returns to "Składniki"
async function closeIng() {
    showIng.value = false;
    await nextTick();
    bar.value?.focusIngredients();
}

function next() {
    if (!step.value) {
        i.value = 0;
        return;
    }
    emit('set-step', step.value.stepNo, true);
    if (last.value) emit('close');
    else i.value++;
}
/** An optional step: on without checking it ("Pomiń"); on the last step, "Pomiń i zakończ". */
function skip() {
    if (last.value) emit('close');
    else i.value++;
}
function prev() {
    if (i.value === 0) {
        if (hasBefore.value) i.value = -1;
        return;
    }
    if (i.value < 0) return;
    i.value--;
    emit('set-step', step.value!.stepNo, false);
}
const goTo = (stepNo: number) => {
    const idx = steps.value.findIndex((s) => s.stepNo === stepNo);
    if (idx >= 0) i.value = idx;
};

function moreDown() {
    const el = stage.value;
    el?.scrollBy({ top: el.clientHeight * 0.7, behavior: 'smooth' });
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
    root.value?.querySelector<HTMLElement>('.cook-nav__next')?.focus();
    document.addEventListener('keydown', onKey);
    void fit();
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
            :class="`cook--${layout}`"
            role="dialog"
            aria-modal="true"
            aria-label="Tryb gotowania"
        >
            <CookBar
                ref="bar"
                :n="i + 1"
                :total="steps.length"
                @ingredients="showIng = true"
                @close="emit('close')"
            />
            <div class="cook__body">
                <div
                    ref="stage"
                    class="cook__stage"
                    @scroll="updateMore"
                >
                    <div
                        class="cook__content"
                        aria-live="polite"
                    >
                        <CookBefore
                            v-if="!step"
                            :name="recipe.name"
                            :steps="steps.length"
                            :cautions="recipe.cautions"
                            :prep="prep"
                            @step="goTo"
                        />
                        <template v-else>
                            <div
                                v-if="step.section"
                                class="cook__section"
                            >
                                {{ step.section }}
                            </div>
                            <div
                                v-if="step.optional || step.prep"
                                class="cook__labels"
                            >
                                <em
                                    v-if="step.optional"
                                    class="cook__opt"
                                >
                                    (opcjonalnie)
                                </em>
                                <span
                                    v-if="step.prep"
                                    class="cook__prep"
                                >
                                    Przygotuj wcześniej{{
                                        step.prepLead ? ` · ${leadText(step)}` : ''
                                    }}
                                </span>
                            </div>
                            <p
                                class="cook__text"
                                :class="size && `cook__text--${size}`"
                            >
                                <StepText
                                    :text="step.text"
                                    :ctx="refs"
                                />
                            </p>
                            <CookHints
                                :hints="step.hints"
                                :refs="refs"
                                :collapsible="collapsible"
                                :open="hintsOpen"
                                @toggle="hintsOpen = !hintsOpen"
                            />
                        </template>
                    </div>
                </div>
                <template v-if="more">
                    <div
                        class="cook__fade"
                        aria-hidden="true"
                    />
                    <button
                        type="button"
                        class="cook__more"
                        @click="moreDown"
                    >
                        Ciąg dalszy
                        <svg
                            width="12"
                            height="12"
                            viewBox="0 0 12 12"
                            fill="none"
                            aria-hidden="true"
                        >
                            <path
                                d="M3 4.5L6 7.5 9 4.5"
                                stroke="currentColor"
                                stroke-width="1.6"
                                stroke-linecap="round"
                                stroke-linejoin="round"
                            />
                        </svg>
                    </button>
                </template>
            </div>
            <CookNav
                :start="!step"
                :can-back="i > 0 || (i === 0 && hasBefore)"
                :last="last"
                :optional="!!step?.optional"
                @back="prev"
                @next="next"
                @skip="skip"
            />
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
    display: grid;
    grid-template-columns: minmax(0, 1fr);
    grid-template-rows: auto auto minmax(0, 1fr) auto;
    grid-template-areas: 'bar' 'progress' 'body' 'nav';
}
.cook__body {
    grid-area: body;
    position: relative;
    min-height: 0;
}
.cook__stage {
    height: 100%;
    overflow-y: auto;
    padding: 22px 20px 20px;
}
.cook__section {
    font-family: var(--font-mono);
    font-size: 13px;
    color: var(--ink-muted);
    letter-spacing: 0.5px;
    text-transform: uppercase;
    font-weight: 600;
    margin-bottom: 12px;
}
.cook__labels {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    gap: 4px 10px;
    margin-bottom: 10px;
}
.cook__opt {
    font-size: 17px;
    color: var(--ink-muted);
}
.cook__prep {
    font-size: 14px;
    font-weight: 600;
    color: var(--prep);
}
.cook__text {
    font-size: 26px;
    line-height: 1.38;
    font-weight: 500;
    letter-spacing: -0.3px;
    text-wrap: pretty;
}
.cook__text--23 {
    font-size: 23px;
}
.cook__text--20 {
    font-size: 20px;
}
.cook__text--21 {
    font-size: 21px;
    line-height: 1.4;
}
.cook__text :deep(.ing-ref) {
    padding: 0.08em 0.32em;
}
.cook--portrait :deep(.hints__text),
.cook--landscape :deep(.hints__text) {
    font-size: 17px;
    line-height: 1.45;
}
.cook__fade {
    position: absolute;
    left: 0;
    right: 0;
    bottom: 0;
    height: 72px;
    background: linear-gradient(to bottom, transparent, var(--bg));
    pointer-events: none;
}
.cook__more {
    position: absolute;
    left: 50%;
    bottom: 8px;
    transform: translateX(-50%);
    display: inline-flex;
    align-items: center;
    gap: 6px;
    height: 32px;
    padding: 0 12px;
    border-radius: var(--r-pill);
    background: var(--bg-alt);
    border: 1px solid var(--rule-strong);
    color: var(--ink);
    font-size: 13px;
    font-weight: 600;
    cursor: pointer;
}

/* desktop: as before — content centred, text scaled with the window; auto margins, so content taller
   than the screen scrolls from the top instead of escaping past the edge */
.cook--desktop .cook__stage {
    display: flex;
    justify-content: center;
    padding: 32px 24px;
    scrollbar-gutter: stable;
}
.cook--desktop .cook__content {
    max-width: 760px;
    width: 100%;
    margin-block: auto;
}
.cook--desktop .cook__text {
    font-size: clamp(24px, 3.4vw, 40px);
    line-height: 1.4;
    letter-spacing: -0.4px;
}
.cook--desktop .cook__section {
    margin-bottom: 14px;
}

/* landscape phone: content on the left, the buttons in a 150 px column on the right */
.cook--landscape {
    grid-template-columns: minmax(0, 1fr) 150px;
    grid-template-rows: auto auto minmax(0, 1fr);
    grid-template-areas: 'bar nav' 'progress nav' 'body nav';
}
.cook--landscape .cook__stage {
    padding: 16px 20px;
}
</style>
