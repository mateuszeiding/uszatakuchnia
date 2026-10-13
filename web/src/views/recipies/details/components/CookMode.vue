<script lang="ts" setup>
import { computed, onBeforeUnmount, onMounted, ref } from 'vue';

import { leadText, scaledAmount, servingsLabel } from '../recipeView';

import UModal from '@/components/UModal.vue';
import { useWakeLock } from '@/composables/useWakeLock';
import type { RecipeDto } from '@/data/dtos/recipe/RecipeDto';
import type { IDictionaries } from '@/data/dtos/taxonomy/DictionariesDto';

// Tryb gotowania: jeden krok naraz, duże Wstecz / Dalej, ← → na klawiaturze, ekran nie gaśnie.
const props = defineProps<{
    recipe: RecipeDto;
    dicts: IDictionaries | undefined;
    servings: number;
    done: number[];
    checkedIng: number[];
}>();
const emit = defineEmits<{
    close: [];
    'set-step': [stepNo: number, done: boolean];
    'toggle-ing': [sortOrder: number];
}>();

const steps = computed(() => props.recipe.steps);
const firstOpen = steps.value.findIndex((s) => !props.done.includes(s.stepNo));
const i = ref(firstOpen === -1 ? 0 : firstOpen);
const showIng = ref(false);
const wake = useWakeLock();
const root = ref<HTMLElement>();

const step = computed(() => steps.value[i.value]!);
const last = computed(() => i.value === steps.value.length - 1);
const ratio = computed(() => (props.recipe.servings ? props.servings / props.recipe.servings : 1));

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
    if (showIng.value) return; // panel składników ma własny modal (Esc, fokus)
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

function amount(ing: RecipeDto['ingredients'][number]) {
    const s = scaledAmount(ing, props.dicts?.units ?? [], ratio.value);
    return [s.qty, s.unit].filter(Boolean).join(' ');
}
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
                        type="button"
                        class="cook__pill"
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
            <div class="cook__stage">
                <div
                    class="cook__content"
                    aria-live="polite"
                >
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
            <div class="cook__hint">
                {{
                    wake
                        ? 'Ekran nie zgaśnie, dopóki tryb jest otwarty'
                        : 'Strzałki ← → zmieniają krok'
                }}
            </div>
        </div>

        <UModal
            v-if="showIng"
            :title="`Składniki · ${servings} ${servingsLabel(servings)}`"
            :width="480"
            @close="showIng = false"
        >
            <ul class="cook-ing">
                <li
                    v-for="ing in recipe.ingredients"
                    :key="ing.sortOrder"
                >
                    <button
                        type="button"
                        role="checkbox"
                        class="cook-ing__item"
                        :class="{ 'is-done': checkedIng.includes(ing.sortOrder) }"
                        :aria-checked="checkedIng.includes(ing.sortOrder)"
                        @click="emit('toggle-ing', ing.sortOrder)"
                    >
                        <span class="cook-ing__qty">{{ amount(ing) }}</span>
                        <span>
                            {{ ing.ingredientName }}{{ ing.optional ? ' (opcjonalnie)' : '' }}
                        </span>
                    </button>
                </li>
            </ul>
        </UModal>
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
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 32px 24px;
}
.cook__content {
    max-width: 760px;
    width: 100%;
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
    background: var(--accent);
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
.cook-ing {
    list-style: none;
    padding: 0;
    margin: 0;
    max-height: 60vh;
    overflow-y: auto;
}
.cook-ing__item {
    width: 100%;
    display: grid;
    grid-template-columns: 100px minmax(0, 1fr);
    gap: 10px;
    padding: 12px 0;
    border: none;
    border-bottom: 1px solid var(--rule);
    background: transparent;
    text-align: left;
    cursor: pointer;
    color: var(--ink);
    font-size: 16px;
}
.cook-ing__item.is-done {
    color: var(--ink-muted);
    text-decoration: line-through;
}
.cook-ing__qty {
    font-family: var(--font-mono);
    font-weight: 600;
    font-size: 14px;
}
</style>
