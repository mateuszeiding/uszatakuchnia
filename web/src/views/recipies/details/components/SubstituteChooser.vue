<script lang="ts" setup>
import { computed, nextTick, onBeforeUnmount, onMounted, ref, useId } from 'vue';

import QtyText from '@/components/QtyText.vue';
import type { ISubChoice } from '@/composables/useCookingProgress';
import type {
    IRecipeSubstitute,
    RecipeIngredientDto,
} from '@/data/dtos/recipe/RecipeIngredientDto';
import { optionAllergens, optionName } from '@/shared/utils/substitutes';

// Teleported: the ingredient list scrolls inside the panel and must not clip the popover.
const props = defineProps<{
    ing: RecipeIngredientDto;
    selected: IRecipeSubstitute | null;
    amount: { qty: string; unit: string };
    half: { qty: string; unit: string };
    allergenLabel: (code: string) => string;
    /** Whole-recipe allergens: "zawiera: X" only when the option adds an allergen the recipe doesn't have. */
    recipeAllergens?: string[];
    anchor: HTMLElement | null;
    sheet: boolean;
}>();
const emit = defineEmits<{ select: [choice: ISubChoice | null]; close: [] }>();

const uid = useId();
const root = ref<HTMLElement>();
const pos = ref<{ left: number; top: number } | null>(null);

const options = computed(() => [null, ...props.ing.substitutes]);
const isOn = (o: IRecipeSubstitute | null) =>
    o === null
        ? !props.selected
        : !!props.selected &&
          props.selected.ingredientId === o.ingredientId &&
          props.selected.mix === o.mix;

/** Allergens the option adds to the recipe (neither the main ingredient nor the rest of the recipe has them). */
const added = (o: IRecipeSubstitute | null) =>
    optionAllergens(props.ing, o)
        .filter((a) => !props.ing.allergens.includes(a) && !props.recipeAllergens?.includes(a))
        .map(props.allergenLabel);

const note = (o: IRecipeSubstitute | null) =>
    o === null ? 'w przepisie' : o.mix ? 'po połowie ilości' : '';

function place() {
    if (props.sheet || !props.anchor || !root.value) return;
    const r = props.anchor.getBoundingClientRect();
    const w = root.value.offsetWidth;
    const h = root.value.offsetHeight;
    const gap = 4;
    let left = r.right + gap;
    if (left + w > window.innerWidth - 8) left = Math.max(8, r.left - w - gap);
    const top = Math.min(Math.max(8, r.top), window.innerHeight - h - 8);
    pos.value = { left, top };
}

function trap(e: KeyboardEvent) {
    const list = [...(root.value?.querySelectorAll<HTMLElement>('button, input:checked') ?? [])];
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

// Capture phase: Esc doesn't also close the modal underneath (cooking mode).
function onKey(e: KeyboardEvent) {
    if (e.key === 'Escape') {
        e.stopPropagation();
        e.preventDefault();
        emit('close');
    } else if (e.key === 'Tab') {
        e.stopPropagation();
        trap(e);
    }
}

function onDown(e: PointerEvent) {
    const t = e.target as Node;
    if (root.value?.contains(t)) return;
    // clicking the icon is handled by the row itself (toggles)
    if (props.anchor?.querySelector('.ing-swap')?.contains(t)) return;
    emit('close');
}

onMounted(async () => {
    document.addEventListener('keydown', onKey, true);
    document.addEventListener('pointerdown', onDown, true);
    window.addEventListener('resize', place);
    window.addEventListener('scroll', place, true);
    await nextTick();
    place();
    root.value?.querySelector<HTMLElement>('input:checked')?.focus();
});
onBeforeUnmount(() => {
    document.removeEventListener('keydown', onKey, true);
    document.removeEventListener('pointerdown', onDown, true);
    window.removeEventListener('resize', place);
    window.removeEventListener('scroll', place, true);
});
</script>

<template>
    <Teleport to="body">
        <div
            v-if="sheet"
            class="subs-scrim"
            aria-hidden="true"
        />
        <div
            ref="root"
            class="subs"
            :class="sheet ? 'subs--sheet' : 'subs--pop'"
            :style="!sheet && pos ? { left: `${pos.left}px`, top: `${pos.top}px` } : undefined"
            role="dialog"
            :aria-modal="sheet || undefined"
            :aria-labelledby="`${uid}-h`"
        >
            <div
                v-if="sheet"
                class="subs__grip"
                aria-hidden="true"
            />
            <div class="subs__head">
                <h3
                    :id="`${uid}-h`"
                    class="subs__title"
                >
                    Czym zastąpić: {{ ing.ingredientName.toLowerCase() }}
                </h3>
                <button
                    v-if="sheet"
                    type="button"
                    class="subs__close"
                    aria-label="Zamknij"
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
            <div
                role="radiogroup"
                class="subs__list"
                :aria-labelledby="`${uid}-h`"
            >
                <label
                    v-for="(o, i) in options"
                    :key="o ? `${o.ingredientId}-${o.mix}` : 'main'"
                    class="subs__opt"
                    :class="{ 'is-on': isOn(o) }"
                >
                    <input
                        type="radio"
                        :name="`${uid}-r`"
                        :checked="isOn(o)"
                        class="subs__radio"
                        @change="
                            emit('select', o ? { ingredientId: o.ingredientId, mix: o.mix } : null)
                        "
                    />
                    <span class="subs__body">
                        <span class="subs__name">{{ optionName(ing, o) }}</span>
                        <span
                            v-if="note(o) || added(o).length || (o && !o.mix && !o.hasNutrition)"
                            class="subs__desc"
                        >
                            {{ note(o) }}
                            <span
                                v-if="added(o).length"
                                class="subs__alg"
                            >
                                {{ note(o) ? ' · ' : '' }}zawiera: {{ added(o).join(', ') }}
                            </span>
                            <template v-if="o && !o.mix && !o.hasNutrition">
                                {{ added(o).length ? ' · ' : '' }}bez danych o wartościach
                                odżywczych
                            </template>
                        </span>
                    </span>
                    <span class="subs__qty">
                        <template v-if="o?.mix">
                            <QtyText :text="half.qty" />
                            +
                            <QtyText :text="half.qty" />
                            {{ half.unit }}
                        </template>
                        <template v-else>
                            <QtyText :text="amount.qty" />
                            {{ amount.unit }}
                        </template>
                    </span>
                    <span
                        v-if="i === 0"
                        class="sr-only"
                    >
                        (składnik z przepisu)
                    </span>
                </label>
            </div>
            <p class="subs__foot">
                Wybór zmienia wartości odżywcze, alergeny i listę zakupów. Zapamiętany na tym
                urządzeniu.
            </p>
        </div>
    </Teleport>
</template>

<style scoped>
.subs {
    z-index: calc(var(--z-modal) + 10);
    background: var(--bg-alt);
    border: 1px solid var(--rule-strong);
    color: var(--ink);
}
.subs--pop {
    position: fixed;
    width: 340px;
    max-width: calc(100vw - 16px);
    border-radius: var(--r-xl);
    padding: 12px;
    box-shadow: 0 12px 32px rgba(10, 10, 15, 0.18);
}
.subs--sheet {
    position: fixed;
    left: 0;
    right: 0;
    bottom: 0;
    max-height: 85vh;
    overflow-y: auto;
    border-radius: 16px 16px 0 0;
    border-bottom: none;
    padding: 8px 16px calc(20px + env(safe-area-inset-bottom));
    box-shadow: 0 -12px 40px rgba(10, 10, 15, 0.25);
}
.subs-scrim {
    position: fixed;
    inset: 0;
    z-index: calc(var(--z-modal) + 9);
    background: rgba(10, 10, 15, 0.45);
}
.subs__grip {
    width: 36px;
    height: 4px;
    margin: 4px auto 12px;
    border-radius: var(--r-pill);
    background: var(--rule-strong);
}
.subs__head {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 8px;
    margin-bottom: 6px;
}
.subs__title {
    font-size: 15px;
    font-weight: 600;
}
.subs--sheet .subs__title {
    font-size: 17px;
}
.subs__close {
    width: 44px;
    height: 44px;
    margin: -8px -10px -8px 0;
    border: none;
    background: none;
    color: var(--ink-muted);
    display: inline-flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    flex-shrink: 0;
}
.subs__list {
    display: flex;
    flex-direction: column;
}
.subs__opt {
    display: grid;
    grid-template-columns: 20px minmax(0, 1fr) auto;
    gap: 10px;
    align-items: start;
    padding: 9px 4px;
    cursor: pointer;
}
.subs__opt + .subs__opt {
    border-top: 1px solid var(--rule);
}
.subs--sheet .subs__opt {
    min-height: 52px;
    padding: 12px 4px;
}
.subs__radio {
    width: 16px;
    height: 16px;
    margin: 3px 0 0;
    accent-color: var(--accent-fill);
}
.subs__body {
    min-width: 0;
}
.subs__name {
    display: block;
    font-size: 14px;
    font-weight: 500;
}
.subs--sheet .subs__name {
    font-size: 16px;
}
.subs__opt.is-on .subs__name {
    font-weight: 600;
}
.subs__desc {
    display: block;
    margin-top: 1px;
    font-size: 12px;
    color: var(--ink-muted);
}
.subs__alg {
    color: var(--prep);
    font-weight: 600;
}
.subs__qty {
    padding-top: 1px;
    font-family: var(--font-mono);
    font-size: 13px;
    font-weight: 600;
    color: var(--ink-dim);
    white-space: nowrap;
}
.subs__foot {
    margin: 8px 0 0;
    padding-top: 8px;
    border-top: 1px solid var(--rule);
    font-size: 12px;
    line-height: 1.45;
    color: var(--ink-muted);
}
</style>
