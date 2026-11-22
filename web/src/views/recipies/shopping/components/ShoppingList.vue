<script lang="ts" setup>
import { computed } from 'vue';

import ShoppingHeader from './ShoppingHeader.vue';
import ShoppingRow from './ShoppingRow.vue';

import CautionsBlock from '@/components/CautionsBlock.vue';
import { useCookingProgress } from '@/composables/useCookingProgress';
import { usePortion } from '@/composables/usePortion';
import { useShoppingChecks } from '@/composables/useShoppingChecks';
import type { RecipeDto } from '@/data/dtos/recipe/RecipeDto';
import type { IDictionaries } from '@/data/dtos/taxonomy/DictionariesDto';
import { parseSizeQuery, type SizeChoice, sizeToQuery } from '@/shared/utils/portion';
import { pluralPL } from '@/shared/utils/quantity';
import { buildShopping, type IShopItem } from '@/shared/utils/shoppingList';
import { effectiveLines, parseSubsParam, subsToParam } from '@/shared/utils/substitutes';

const props = defineProps<{
    recipe: RecipeDto;
    dicts: IDictionaries | undefined;
    urlServings: number | null;
    /** The whole query: ?sztuki=&rozmiar= / ?forma= are parsed here (the recipe decides which apply). */
    query: Record<string, unknown>;
    /** Raw ?zamienniki= from the URL (the link sender's variant). */
    urlSubs: string | null;
}>();
const emit = defineEmits<{
    servings: [n: number];
    size: [s: SizeChoice | null];
    /** Fills the URL in with the stored state when the link didn't have it (so the URL can be passed on). */
    sync: [query: Record<string, string>];
}>();

// Servings and substitutes are shared with the recipe (uk:progress:<id>); the URL takes precedence (a shared link)
// and, like servings, gets saved on the device, so "‹ Przepis" shows the same variant.
const { progress, setServings, setSize } = useCookingProgress(
    props.recipe.id,
    props.recipe.servings || 1
);
const portion = usePortion(() => props.recipe, progress);
const missing: Record<string, string> = {};
if (portion.mode.value === 'servings') {
    if (props.urlServings != null) setServings(props.urlServings);
    else missing.porcje = String(progress.value.servings);
} else {
    const fromUrl = parseSizeQuery(props.recipe, props.query);
    if (fromUrl) setSize(fromUrl);
    // also when the link's size fell back (e.g. a diameter the author removed): the URL shows what's on screen
    const shown = sizeToQuery(portion.size.value);
    for (const [k, v] of Object.entries(shown)) if (props.query[k] !== v) missing[k] = v;
}
const urlSubs = parseSubsParam(props.urlSubs, props.recipe.ingredients);
if (urlSubs) progress.value.subs = urlSubs;
else {
    const param = subsToParam(props.recipe.ingredients, progress.value.subs);
    if (param) missing.zamienniki = param;
}
if (Object.keys(missing).length) emit('sync', missing);
const changeServings = (n: number) => {
    setServings(n);
    emit('servings', progress.value.servings);
};
const changeSize = (s: SizeChoice | null) => {
    setSize(s);
    emit('size', portion.size.value);
};
const headerPortion = computed(() => ({
    recipe: props.recipe,
    mode: portion.mode.value,
    size: portion.size.value,
    scale: portion.scale.value,
    changed: portion.changed.value,
    label: portion.label.value,
}));

const { checked, toggle, clear } = useShoppingChecks(props.recipe.id);

const ratio = computed(() => portion.scale.value);
// Chosen recipe variant (substitutes from uk:progress:<id>.subs): the substitute instead of the main one, "pół na pół" = 2 items.
const items = computed(() =>
    buildShopping(
        effectiveLines(props.recipe.ingredients, progress.value.subs),
        props.dicts?.units ?? [],
        ratio.value
    )
);
const isOn = (it: IShopItem) => checked.value.includes(it.id);
const done = computed(() => items.value.filter(isOn));
const optional = computed(() => items.value.filter((x) => x.optional));

/** Store aisles = ingredient types in dictionary order; unknown types last.
 *  Checked items stay in their aisle; the printout skips them. */
const aisles = computed(() => {
    const todo = items.value.filter((x) => !x.optional);
    const order = (props.dicts?.ingredientTypes ?? []).map((t) => t.code as string);
    for (const it of todo) if (!order.includes(it.type.code)) order.push(it.type.code);
    return order
        .map((code) => ({
            code,
            label: todo.find((x) => x.type.code === code)?.type.label ?? '',
            items: todo.filter((x) => x.type.code === code),
        }))
        .filter((a) => a.items.length);
});

const skipped = computed(() => {
    const n = done.value.length;
    if (!n) return '';
    const what = pluralPL(n, 'odhaczoną pozycję', 'odhaczone pozycje', 'odhaczonych pozycji');
    return `Pominięto ${n} ${what} (mam w domu): ${done.value.map((x) => x.name).join(', ')}.`;
});
const recipeUrl = computed(() => `${location.host}/recipes/${props.recipe.id}`);
const today = new Date().toLocaleDateString('pl-PL');
</script>

<template>
    <ShoppingHeader
        :recipe-id="recipe.id"
        :recipe-name="recipe.name"
        :servings="progress.servings"
        :portion="headerPortion"
        :done="done.length"
        :total="items.length"
        @servings="changeServings"
        @size="changeSize"
    />
    <div class="shop">
        <CautionsBlock
            class="shop__cautions"
            :items="recipe.cautions"
            variant="compact"
        />
        <p class="shop__hint">
            Odhacz to, co masz w domu albo już w koszyku. Ilości zsumowane ze wszystkich grup
            przepisu.
        </p>
        <CautionsBlock
            class="shop__cautions-print"
            :items="recipe.cautions"
            variant="print"
        />

        <div class="shop__cols">
            <section
                v-for="a in aisles"
                :key="a.code"
                class="aisle"
                :class="{ 'is-all-done': a.items.every(isOn) }"
                :aria-label="a.label"
            >
                <h2 class="aisle__title">{{ a.label }}</h2>
                <ul class="aisle__list">
                    <ShoppingRow
                        v-for="it in a.items"
                        :key="it.id"
                        :item="it"
                        :checked="isOn(it)"
                        @toggle="toggle(it.id)"
                    />
                </ul>
            </section>

            <section
                v-if="optional.length"
                class="aisle aisle--opt"
                :class="{ 'is-all-done': optional.every(isOn) }"
                aria-label="Opcjonalne"
            >
                <h2 class="aisle__title">
                    <span>Opcjonalne</span>
                    <span class="aisle__aside">jeśli chcesz</span>
                </h2>
                <ul class="aisle__list">
                    <ShoppingRow
                        v-for="it in optional"
                        :key="it.id"
                        :item="it"
                        :checked="isOn(it)"
                        @toggle="toggle(it.id)"
                    />
                </ul>
            </section>
        </div>

        <div
            v-if="done.length"
            class="done"
        >
            <button
                type="button"
                class="done__clear"
                @click="clear"
            >
                Odznacz wszystko
            </button>
        </div>

        <footer class="print-foot">
            <span>{{ skipped }}</span>
            <span>{{ recipeUrl }} · {{ today }}</span>
        </footer>
    </div>
</template>

<style lang="scss" scoped>
@use '@/scss/mixins/media' as m;

.shop {
    max-width: 640px;
    margin: 0 auto;
    padding: 8px 0 40px;
    display: flex;
    flex-direction: column;
    gap: 24px;
}
.shop__cautions {
    margin-top: 14px;
}
/* the notes take the space under the header, so the hint has no gap of its own */
.shop__cautions + .shop__hint {
    margin-top: 0;
}
.shop__cautions-print {
    display: none;
}
.shop__hint {
    margin-top: 14px;
    font-size: 13px;
    line-height: 1.5;
    color: var(--ink-muted);
}
.shop__cols {
    display: flex;
    flex-direction: column;
    gap: 24px;
}
.aisle__title {
    display: flex;
    justify-content: space-between;
    padding-bottom: 6px;
    border-bottom: 1px solid var(--rule-strong);
    font-family: var(--font-mono);
    font-size: 12px;
    font-weight: 600;
    letter-spacing: 0.5px;
    text-transform: uppercase;
    color: var(--ink);
}
.aisle--opt .aisle__title {
    color: var(--ink-muted);
    border-bottom-style: dashed;
}
.aisle__aside {
    font-family: var(--font-sans);
    font-weight: 500;
    letter-spacing: 0;
    text-transform: none;
}
.aisle__list {
    list-style: none;
    padding: 0;
    margin: 0;
}
.done__clear {
    background: none;
    border: none;
    padding: 6px 0;
    color: var(--ink-dim);
    font-size: 13px;
    text-decoration: underline;
    text-underline-offset: 3px;
    cursor: pointer;
}
.print-foot {
    display: none;
}

@media (max-width: 900px) {
    .shop {
        padding: 4px 16px 28px;
        gap: 20px;
    }
    .shop__hint,
    .shop__cautions {
        margin-top: 10px;
    }
    .shop__cautions + .shop__hint {
        margin-top: 0;
    }
    .shop__cols {
        gap: 20px;
    }
}

/* A4 print: two columns by aisle, only unchecked items, a footer with the skipped ones */
@include m.print {
    .shop {
        max-width: none;
        padding: 18px 0 0;
        gap: 16px;
    }
    .shop__hint,
    .shop__cautions,
    .done {
        display: none;
    }
    .shop__cautions-print {
        display: block;
    }
    .shop__cols {
        display: block;
        column-count: 2;
        column-gap: 32px;
    }
    .aisle {
        break-inside: avoid;
        margin-bottom: 16px;
    }
    .aisle.is-all-done {
        display: none;
    }
    .aisle__title {
        padding-bottom: 2px;
        border: none;
        font-size: 10px;
        letter-spacing: 0.6px;
    }
    .aisle--opt .aisle__title {
        color: #000;
        border: none;
    }
    .print-foot {
        display: flex;
        justify-content: space-between;
        gap: 16px;
        padding-top: 6px;
        border-top: 0.5pt solid #999;
        font-size: 10px;
        color: #444;
    }
}
</style>
