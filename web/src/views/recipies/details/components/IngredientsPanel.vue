<script lang="ts" setup>
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue';

import { ingredientGroups, scaledAmount } from '../recipeView';

import IngGroupHeader from './IngGroupHeader.vue';
import IngRow from './IngRow.vue';
import IngRowSkeleton from './IngRowSkeleton.vue';
import ServingsStepper from './ServingsStepper.vue';

import PortionControl from '@/components/portion/PortionControl.vue';
import PortionStickyBar from '@/components/portion/PortionStickyBar.vue';
import RescaledNote from '@/components/portion/RescaledNote.vue';
import SwapIcon from '@/components/SwapIcon.vue';
import type { ISubChoice } from '@/composables/useCookingProgress';
import { useNarrow } from '@/composables/useNarrow';
import type { RecipeDto } from '@/data/dtos/recipe/RecipeDto';
import type { IDictionaries } from '@/data/dtos/taxonomy/DictionariesDto';
import {
    defaultSize,
    effectiveMode,
    isForm,
    sameSize,
    type SizeChoice,
    sizeLabel,
    sizeToQuery,
} from '@/shared/utils/portion';
import {
    activeSubs,
    chosenSub,
    effectiveLines,
    type SubsMap,
    subsToParam,
} from '@/shared/utils/substitutes';

const props = defineProps<{
    /** null = loading. */
    recipe: RecipeDto | null;
    dicts: IDictionaries | undefined;
    servings: number;
    /** Pieces in a size / the reader's form (usePortion); null in servings mode. */
    size?: SizeChoice | null;
    /** portionScale for the chosen servings or size; missing = 1 (skeleton). */
    scale?: number;
    checked: number[];
    /** Chosen substitutes (uk:progress:<id>.subs). */
    subs: SubsMap;
}>();
const emit = defineEmits<{
    toggle: [sortOrder: number];
    servings: [n: number];
    /** null = back to the recipe's size. */
    size: [s: SizeChoice | null];
    sub: [sortOrder: number, choice: ISubChoice | null];
    'clear-subs': [];
}>();

const activeCount = computed(() => activeSubs(props.recipe?.ingredients ?? [], props.subs).size);
const allergenLabel = (code: string) =>
    props.dicts?.allergens.find((a) => a.code === code)?.label ?? code;

const ratio = computed(() => props.scale ?? 1);
const mode = computed(() => (props.recipe ? effectiveMode(props.recipe) : 'servings'));
const changed = computed(
    () =>
        !!props.recipe &&
        mode.value !== 'servings' &&
        !sameSize(props.size ?? null, defaultSize(props.recipe))
);
const groups = computed(() => ingredientGroups(props.recipe?.ingredients ?? []));
// The shopping list link carries servings (or the size) and the chosen variant (substitutes), so the recipient
// sees what the sender sees.
const shopLink = computed(() => {
    const zamienniki = subsToParam(props.recipe?.ingredients ?? [], props.subs);
    const portion =
        mode.value === 'servings'
            ? { porcje: String(props.servings) }
            : sizeToQuery(props.size ?? null);
    return {
        path: `/recipes/${props.recipe?.id ?? 0}/zakupy`,
        query: { ...portion, ...(zamienniki ? { zamienniki } : {}) },
    };
});

// Phone: once the size control scrolls under the nav bar, the short bar sticks over the list.
const narrow = useNarrow();
const ctl = ref<HTMLElement>();
const stuck = ref(false);
let io: IntersectionObserver | undefined;
watch(
    [ctl, narrow],
    ([el, n]) => {
        io?.disconnect();
        stuck.value = false;
        if (!el || !n) return;
        io = new IntersectionObserver(
            ([e]) => (stuck.value = !!e && !e.isIntersecting && e.boundingClientRect.top < 64),
            { rootMargin: '-64px 0px 0px 0px' }
        );
        io.observe(el);
    },
    { flush: 'post' }
);
onBeforeUnmount(() => io?.disconnect());

function amount(ing: RecipeDto['ingredients'][number], factor = 1) {
    return scaledAmount(ing, props.dicts?.units ?? [], ratio.value * factor);
}
const doneIn = (items: RecipeDto['ingredients']) =>
    items.filter((i) => props.checked.includes(i.sortOrder)).length;

/** Allergen → the ingredients it comes from (EU list order), per the chosen variant. An ingredient from several groups appears once;
 *  "(opcj.)" only when it's optional in every occurrence. */
const allergens = computed(() =>
    (props.dicts?.allergens ?? [])
        .map((a) => {
            const by = new Map<number, { name: string; optional: boolean }>();
            for (const i of effectiveLines(props.recipe?.ingredients ?? [], props.subs)) {
                if (!i.allergens.includes(a.code)) continue;
                const seen = by.get(i.ingredientId);
                if (seen) seen.optional &&= i.optional;
                else by.set(i.ingredientId, { name: i.ingredientName, optional: i.optional });
            }
            const sources = [...by.values()].map((s) => s.name + (s.optional ? ' (opcj.)' : ''));
            return { label: a.label, sources };
        })
        .filter((a) => a.sources.length)
);

// Desktop: the list scrolls inside the panel; the bottom gradient shows only when there's more below.
const list = ref<HTMLElement>();
const more = ref(false);
const measure = () => {
    const el = list.value;
    more.value = !!el && el.scrollHeight - el.scrollTop - el.clientHeight > 2;
};
let ro: ResizeObserver | undefined;
onMounted(() => {
    measure();
    ro = new ResizeObserver(measure);
    if (list.value) ro.observe(list.value);
});
onBeforeUnmount(() => ro?.disconnect());
</script>

<template>
    <section
        id="skladniki"
        class="ing-panel"
        :class="{ 'ing-panel--bar': mode !== 'servings' }"
        aria-labelledby="h-ing"
        :aria-busy="!recipe"
        :inert="!recipe"
    >
        <div class="ing-panel__head">
            <h2
                id="h-ing"
                class="ing-panel__title"
            >
                Składniki
                <span
                    v-if="groups.single"
                    class="ing-panel__single"
                >
                    · {{ groups.single.toLowerCase() }}
                </span>
            </h2>
            <RouterLink
                :to="shopLink"
                class="btn btn--secondary btn--sm ing-panel__shop ing-panel__shop--top"
            >
                <svg
                    viewBox="0 0 16 16"
                    fill="none"
                    aria-hidden="true"
                >
                    <path
                        d="M1.5 2h2l1.6 8.2a1 1 0 001 .8h6.3a1 1 0 001-.8L14.5 5H4.3"
                        stroke="currentColor"
                        stroke-width="1.4"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                    />
                    <circle
                        cx="6.5"
                        cy="13.5"
                        r="1"
                        fill="currentColor"
                    />
                    <circle
                        cx="12"
                        cy="13.5"
                        r="1"
                        fill="currentColor"
                    />
                </svg>
                Lista zakupów
            </RouterLink>
        </div>
        <ServingsStepper
            v-if="mode === 'servings'"
            class="ing-panel__servings"
            :value="servings"
            @set="emit('servings', $event)"
        />
        <div
            v-else-if="recipe"
            ref="ctl"
        >
            <PortionControl
                :recipe="recipe"
                :size="size ?? null"
                :big="narrow"
                @size="emit('size', $event)"
            />
        </div>
        <RescaledNote
            v-if="changed && recipe"
            :text="sizeLabel(recipe, size ?? null)"
            :scale="ratio"
            :form="isForm(size)"
            @restore="emit('size', null)"
        />
        <div
            v-if="narrow && recipe && mode !== 'servings'"
            class="ing-panel__bar"
            :class="{ 'is-stuck': stuck }"
            :aria-hidden="!stuck || undefined"
            :inert="!stuck || undefined"
        >
            <PortionStickyBar
                :recipe="recipe"
                :size="size ?? null"
                :scale="ratio"
                @size="emit('size', $event)"
            />
        </div>
        <div
            v-if="activeCount"
            class="ing-panel__variant"
        >
            <span class="ing-panel__variant-n">
                <SwapIcon :size="13" />
                Z zamiennikami: {{ activeCount }}
            </span>
            <button
                type="button"
                class="ing-panel__restore"
                @click="emit('clear-subs')"
            >
                Przywróć przepis
            </button>
        </div>

        <div
            class="ing-panel__scroll"
            :class="{ 'has-more': more }"
        >
            <div
                ref="list"
                class="ing-panel__list"
                @scroll.passive="measure"
            >
                <ul
                    v-if="!recipe"
                    class="ing-list"
                >
                    <IngRowSkeleton
                        v-for="n in 8"
                        :key="n"
                        :wide="n % 3 === 0"
                    />
                </ul>
                <div
                    v-for="(sec, si) in groups.sections"
                    :key="si"
                >
                    <IngGroupHeader
                        v-if="groups.grouped && sec.label"
                        class="ing-panel__ghead"
                        :label="sec.label"
                        :count="sec.items.length"
                        :done="doneIn(sec.items)"
                    />
                    <ul class="ing-list">
                        <IngRow
                            v-for="ing in sec.items"
                            :key="ing.sortOrder"
                            :ing="ing"
                            :amount="amount(ing)"
                            :half="ing.substitutes.length ? amount(ing, 0.5) : undefined"
                            :checked="checked.includes(ing.sortOrder)"
                            :sub="chosenSub(ing, subs)"
                            :allergen-label="allergenLabel"
                            :recipe-allergens="recipe?.allergens"
                            @toggle="emit('toggle', ing.sortOrder)"
                            @sub="emit('sub', ing.sortOrder, $event)"
                        />
                    </ul>
                </div>

                <div
                    v-if="allergens.length"
                    class="allergens"
                >
                    <h3 class="allergens__title">Alergeny</h3>
                    <ul class="allergens__list">
                        <li
                            v-for="a in allergens"
                            :key="a.label"
                        >
                            <strong>{{ a.label }}</strong>
                            <span class="allergens__src">· {{ a.sources.join(', ') }}</span>
                        </li>
                    </ul>
                </div>
            </div>
        </div>

        <RouterLink
            :to="shopLink"
            class="btn btn--secondary ing-panel__shop ing-panel__shop--bottom"
        >
            <svg
                viewBox="0 0 16 16"
                fill="none"
                aria-hidden="true"
            >
                <path
                    d="M1.5 2h2l1.6 8.2a1 1 0 001 .8h6.3a1 1 0 001-.8L14.5 5H4.3"
                    stroke="currentColor"
                    stroke-width="1.4"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                />
                <circle
                    cx="6.5"
                    cy="13.5"
                    r="1"
                    fill="currentColor"
                />
                <circle
                    cx="12"
                    cy="13.5"
                    r="1"
                    fill="currentColor"
                />
            </svg>
            Lista zakupów
        </RouterLink>
    </section>
</template>

<style scoped>
/* Desktop: sticky panel, no taller than the window; the servings counter stays put, only the list scrolls. */
.ing-panel {
    position: sticky;
    top: 84px;
    max-height: calc(100vh - 84px - 24px);
    background: var(--bg-alt);
    border: 1px solid var(--rule);
    border-radius: var(--r-xl);
    padding: 24px;
    display: flex;
    flex-direction: column;
    gap: 16px;
}
.ing-panel__head {
    display: flex;
    align-items: baseline;
    justify-content: space-between;
    gap: 12px;
    flex-wrap: wrap;
}
.ing-panel__title {
    font-size: 26px;
    font-weight: 500;
    letter-spacing: -0.8px;
}
.ing-panel__single {
    font-size: 15px;
    letter-spacing: 0;
    color: var(--ink-muted);
}
.ing-panel__variant {
    margin-top: -4px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 8px;
    font-size: 13px;
    color: var(--ink-dim);
}
.ing-panel__variant-n {
    display: inline-flex;
    align-items: center;
    gap: 6px;
}
.ing-panel__variant-n svg {
    color: var(--accent);
}
.ing-panel__restore {
    padding: 6px 0;
    border: none;
    background: none;
    color: var(--ink-dim);
    font-size: 13px;
    text-decoration: underline;
    text-underline-offset: 3px;
    cursor: pointer;
}
.ing-panel__shop {
    gap: 7px;
    white-space: nowrap;
}
.ing-panel__shop--bottom {
    display: none;
}
.ing-panel__scroll {
    position: relative;
    display: flex;
    flex: 1 1 auto;
    min-height: 0;
}
.ing-panel__list {
    flex: 1;
    min-height: 0;
    overflow-y: auto;
    scrollbar-gutter: stable;
    padding-right: 6px;
    display: flex;
    flex-direction: column;
    gap: 22px;
}
.ing-panel__scroll.has-more::after {
    content: '';
    position: absolute;
    left: 0;
    right: 6px;
    bottom: 0;
    height: 36px;
    background: linear-gradient(to bottom, transparent, var(--bg-alt));
    pointer-events: none;
}
.ing-panel__ghead {
    position: sticky;
    top: 0;
    z-index: 1;
    background: var(--bg-alt);
    padding-top: 8px;
}
.ing-list {
    list-style: none;
    padding: 0;
    margin: 0;
}
.allergens__title {
    font-family: var(--font-mono);
    font-size: 12px;
    color: var(--ink-muted);
    letter-spacing: 0.5px;
    text-transform: uppercase;
    font-weight: 600;
    margin-bottom: 8px;
}
.allergens__list {
    list-style: none;
    padding: 0;
    margin: 0;
    display: flex;
    flex-direction: column;
    gap: 4px;
    font-size: 14px;
    line-height: 1.45;
}
.allergens__list strong {
    font-weight: 600;
}
.allergens__src {
    margin-left: 4px;
    color: var(--ink-muted);
}

@media (max-width: 900px) {
    .ing-panel {
        position: static;
        max-height: none;
        padding: 20px;
        /* coming back from the shopping list (#skladniki): the panel header below the nav bar */
        scroll-margin-top: 72px;
    }
    .ing-panel__list {
        overflow: visible;
        padding-right: 0;
    }
    .ing-panel__scroll.has-more::after {
        content: none;
    }
    /* servings counter sticks below the nav bar, the group header below the counter (64 + 54 px) */
    .ing-panel__servings {
        position: sticky;
        top: 64px;
        z-index: 2;
    }
    .ing-panel__ghead {
        top: 118px;
    }
    /* sized / form: the short bar (8 + 44 + 8 + 1 px) sticks instead of the counter; it takes no room in the
       flow (negative margin = its height + the panel gap) and shows only once the full control is gone */
    .ing-panel__bar {
        position: sticky;
        top: 64px;
        z-index: 2;
        margin: 0 -20px -77px;
        visibility: hidden;
    }
    .ing-panel__bar.is-stuck {
        visibility: visible;
    }
    .ing-panel--bar .ing-panel__ghead {
        top: 125px;
    }
    .ing-panel__shop--top {
        display: none;
    }
    .ing-panel__shop--bottom {
        display: inline-flex;
        width: 100%;
        justify-content: center;
    }
}
</style>
