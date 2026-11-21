<script lang="ts" setup>
import { computed, nextTick, onBeforeUnmount, onMounted, ref } from 'vue';

import { ingredientGroups, publicLink, scaledAmount } from '../recipeView';

import EquipmentCard from './EquipmentCard.vue';
import IngGroupHeader from './IngGroupHeader.vue';
import IngRow from './IngRow.vue';
import ServingsStepper from './ServingsStepper.vue';

import SwapIcon from '@/components/SwapIcon.vue';
import type { ISubChoice } from '@/composables/useCookingProgress';
import { useNarrow } from '@/composables/useNarrow';
import type { RecipeDto } from '@/data/dtos/recipe/RecipeDto';
import type { IDictionaries } from '@/data/dtos/taxonomy/DictionariesDto';
import { activeSubs, chosenSub, type SubsMap } from '@/shared/utils/substitutes';

const props = defineProps<{
    recipe: RecipeDto;
    dicts: IDictionaries | undefined;
    servings: number;
    checkedIng: number[];
    subs: SubsMap;
    stepNo: number;
}>();
const emit = defineEmits<{
    close: [];
    'toggle-ing': [sortOrder: number];
    sub: [sortOrder: number, choice: ISubChoice | null];
    servings: [n: number];
    'clear-subs': [];
}>();

const title = computed(() => (props.recipe.equipment.length ? 'Sprzęt i składniki' : 'Składniki'));
const groups = computed(() => ingredientGroups(props.recipe.ingredients));
const ratio = computed(() => (props.recipe.servings ? props.servings / props.recipe.servings : 1));
const hasLinks = computed(() => props.recipe.ingredients.some((x) => publicLink(x)));
const subsCount = computed(() => activeSubs(props.recipe.ingredients, props.subs).size);
const allergenLabel = (code: string) =>
    props.dicts?.allergens.find((a) => a.code === code)?.label ?? code;
const amount = (ing: RecipeDto['ingredients'][number], factor = 1) =>
    scaledAmount(ing, props.dicts?.units ?? [], ratio.value * factor);

const narrow = useNarrow();
const root = ref<HTMLElement>();
const heading = ref<HTMLElement>();

const drag = ref(0);
let startY: number | null = null;
function onDragStart(e: PointerEvent) {
    if (!narrow.value) return;
    startY = e.clientY;
    try {
        (e.currentTarget as HTMLElement).setPointerCapture(e.pointerId);
    } catch {
        /* no pointer capture: moves keep coming while the finger is over the handle */
    }
}
function onDragMove(e: PointerEvent) {
    if (startY == null) return;
    drag.value = Math.max(0, e.clientY - startY);
}
function onDragEnd() {
    if (startY == null) return;
    startY = null;
    if (drag.value > 80) emit('close');
    drag.value = 0;
}

function trap(e: KeyboardEvent) {
    const list = [
        ...(root.value?.querySelectorAll<HTMLElement>(
            'button:not([disabled]), a[href], [tabindex="-1"]'
        ) ?? []),
    ].filter((el) => el.tabIndex >= 0 || el === heading.value);
    if (!list.length) return;
    const first = list[0]!;
    const last = list[list.length - 1]!;
    if (!root.value?.contains(document.activeElement)) {
        e.preventDefault();
        first.focus();
    } else if (e.shiftKey && document.activeElement === first) {
        e.preventDefault();
        last.focus();
    } else if (!e.shiftKey && document.activeElement === last) {
        e.preventDefault();
        first.focus();
    }
}
function onKey(e: KeyboardEvent) {
    // the substitute chooser has its own key handling (capture phase, stopPropagation)
    if (e.key === 'Escape') {
        e.preventDefault();
        emit('close');
    } else if (e.key === 'Tab') trap(e);
}

onMounted(async () => {
    document.addEventListener('keydown', onKey);
    await nextTick();
    heading.value?.focus();
});
onBeforeUnmount(() => document.removeEventListener('keydown', onKey));
</script>

<template>
    <Teleport to="body">
        <div
            class="cis-scrim"
            aria-hidden="true"
            @click="emit('close')"
        />
        <div
            ref="root"
            class="cis"
            role="dialog"
            aria-modal="true"
            aria-labelledby="cis-title"
            :style="drag ? { transform: `translateY(${drag}px)`, transition: 'none' } : undefined"
        >
            <div
                class="cis__top"
                @pointerdown="onDragStart"
                @pointermove="onDragMove"
                @pointerup="onDragEnd"
                @pointercancel="onDragEnd"
            >
                <div
                    class="cis__grip"
                    aria-hidden="true"
                />
                <div class="cis__head">
                    <h2
                        id="cis-title"
                        ref="heading"
                        class="cis__title"
                        tabindex="-1"
                    >
                        {{ title }}
                    </h2>
                    <button
                        type="button"
                        class="cis__close"
                        :aria-label="`Zamknij ${title.toLowerCase()}`"
                        @pointerdown.stop
                        @click="emit('close')"
                    >
                        <svg
                            width="12"
                            height="12"
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

            <div class="cis__body">
                <ServingsStepper
                    :value="servings"
                    big
                    @set="emit('servings', $event)"
                />
                <div
                    v-if="subsCount"
                    class="cis__variant"
                >
                    <span class="cis__variant-n">
                        <SwapIcon :size="13" />
                        Z zamiennikami: {{ subsCount }}
                    </span>
                    <button
                        type="button"
                        class="cis__restore"
                        @click="emit('clear-subs')"
                    >
                        Przywróć przepis
                    </button>
                </div>
                <EquipmentCard
                    v-if="recipe.equipment.length"
                    :items="recipe.equipment"
                    compact
                />
                <div
                    v-for="(sec, si) in groups.sections"
                    :key="si"
                >
                    <IngGroupHeader
                        v-if="groups.grouped && sec.label"
                        :label="sec.label"
                        :count="sec.items.length"
                        :done="sec.items.filter((x) => checkedIng.includes(x.sortOrder)).length"
                    />
                    <ul class="cis__list">
                        <IngRow
                            v-for="ing in sec.items"
                            :key="ing.sortOrder"
                            :ing="ing"
                            :amount="amount(ing)"
                            :half="ing.substitutes.length ? amount(ing, 0.5) : undefined"
                            :checked="checkedIng.includes(ing.sortOrder)"
                            :sub="chosenSub(ing, subs)"
                            :allergen-label="allergenLabel"
                            :recipe-allergens="recipe.allergens"
                            roomy
                            @toggle="emit('toggle-ing', ing.sortOrder)"
                            @sub="emit('sub', ing.sortOrder, $event)"
                        />
                    </ul>
                </div>
                <p
                    v-if="hasLinks"
                    class="cis__note"
                >
                    Linki do przepisów otwierają się w nowej karcie — gotowanie zostaje tutaj.
                </p>
            </div>

            <div class="cis__foot">
                <button
                    type="button"
                    class="btn btn--secondary cis__back"
                    @click="emit('close')"
                >
                    Wróć do kroku {{ stepNo }}
                </button>
            </div>
        </div>
    </Teleport>
</template>

<style scoped>
.cis-scrim {
    position: fixed;
    inset: 0;
    z-index: calc(var(--z-modal) - 1);
    background: rgba(10, 10, 15, 0.5);
}
.cis {
    position: fixed;
    top: 0;
    right: 0;
    bottom: 0;
    width: 420px;
    max-width: 100vw;
    z-index: var(--z-modal);
    display: flex;
    flex-direction: column;
    background: var(--bg-alt);
    color: var(--ink);
    border-left: 1px solid var(--rule-strong);
    box-shadow: -12px 0 40px rgba(10, 10, 15, 0.25);
    transition: transform 0.2s var(--ease);
}
.cis__top {
    touch-action: none;
}
.cis__grip {
    display: none;
}
.cis__head {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 12px 8px 8px 20px;
}
.cis__title {
    margin: 0;
    font-size: 19px;
    font-weight: 600;
}
.cis__title:focus {
    outline: none;
}
.cis__close {
    width: 44px;
    height: 44px;
    flex-shrink: 0;
    border-radius: var(--r-pill);
    border: 1px solid var(--rule-strong);
    background: var(--bg);
    color: var(--ink);
    display: inline-flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
}
.cis__body {
    flex: 1;
    min-height: 0;
    overflow-y: auto;
    overscroll-behavior: contain;
    padding: 0 20px 16px;
    display: flex;
    flex-direction: column;
    gap: 14px;
}
.cis__variant {
    margin-top: -4px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 8px;
    font-size: 13px;
    color: var(--ink-dim);
}
.cis__variant-n {
    display: inline-flex;
    align-items: center;
    gap: 6px;
}
.cis__variant-n svg {
    color: var(--accent);
}
.cis__restore {
    padding: 12px 0;
    margin: -12px 0;
    border: none;
    background: none;
    color: var(--ink-dim);
    font-size: 13px;
    text-decoration: underline;
    text-underline-offset: 3px;
    cursor: pointer;
}
.cis__list {
    list-style: none;
    padding: 0;
    margin: 0;
}
.cis__note {
    margin: 0;
    font-size: 12px;
    line-height: 1.5;
    color: var(--ink-muted);
}
.cis__foot {
    padding: 12px 20px calc(16px + env(safe-area-inset-bottom));
    border-top: 1px solid var(--rule);
    background: var(--bg-alt);
}
.cis__back {
    width: 100%;
    min-height: 56px;
    font-size: 17px;
    border-radius: var(--r-2xl);
}

/* phone: a bottom sheet (modeled on the substitute sheet) */
@media (max-width: 900px) {
    .cis {
        top: 56px;
        left: 0;
        width: auto;
        border-left: none;
        border-top: 1px solid var(--rule-strong);
        border-radius: 16px 16px 0 0;
        box-shadow: 0 -12px 40px rgba(10, 10, 15, 0.25);
    }
    .cis__grip {
        display: block;
        width: 36px;
        height: 4px;
        margin: 8px auto 0;
        border-radius: var(--r-pill);
        background: var(--rule-strong);
    }
    .cis__head {
        padding: 4px 8px 8px 20px;
    }
}
</style>
