<script lang="ts" setup>
import { computed, ref } from 'vue';

import ShoppingHeader from './ShoppingHeader.vue';
import ShoppingRow from './ShoppingRow.vue';

import { useCookingProgress } from '@/composables/useCookingProgress';
import { useShoppingChecks } from '@/composables/useShoppingChecks';
import type { RecipeDto } from '@/data/dtos/recipe/RecipeDto';
import type { IDictionaries } from '@/data/dtos/taxonomy/DictionariesDto';
import { pluralPL } from '@/shared/utils/quantity';
import { buildShopping, type IShopItem } from '@/shared/utils/shoppingList';

const props = defineProps<{
    recipe: RecipeDto;
    dicts: IDictionaries | undefined;
    urlServings: number | null;
}>();
const emit = defineEmits<{ servings: [n: number] }>();

// Porcje wspólne z przepisem (uk:progress:<id>.servings); adres ?porcje= ma pierwszeństwo (link z panelu / udostępniony).
const { progress, setServings } = useCookingProgress(props.recipe.id, props.recipe.servings || 1);
if (props.urlServings != null) setServings(props.urlServings);
else emit('servings', progress.value.servings);
const changeServings = (n: number) => {
    setServings(n);
    emit('servings', progress.value.servings);
};

const { checked, toggle, clear } = useShoppingChecks(props.recipe.id);
const open = ref<number[]>([]);
const flip = (id: number) =>
    (open.value = open.value.includes(id)
        ? open.value.filter((x) => x !== id)
        : [...open.value, id]);
const doneOpen = ref(false);

const ratio = computed(() =>
    props.recipe.servings ? progress.value.servings / props.recipe.servings : 1
);
const items = computed(() =>
    buildShopping(props.recipe.ingredients, props.dicts?.units ?? [], ratio.value)
);
const isOn = (it: IShopItem) => checked.value.includes(it.id);
const done = computed(() => items.value.filter(isOn));
const optTodo = computed(() => items.value.filter((x) => x.optional && !isOn(x)));

/** Działy sklepu = typy składników w kolejności słownika; nieznane typy na końcu. */
const aisles = computed(() => {
    const todo = items.value.filter((x) => !x.optional && !isOn(x));
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

// Stopka wydruku
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
        :done="done.length"
        :total="items.length"
        @servings="changeServings"
    />
    <div class="shop">
        <p class="shop__hint">
            Odhacz to, co masz w domu albo już w koszyku. Ilości zsumowane ze wszystkich grup
            przepisu; liczba przy pozycji — z ilu miejsc się składa.
        </p>

        <div class="shop__cols">
            <section
                v-for="a in aisles"
                :key="a.code"
                class="aisle"
                :aria-label="a.label"
            >
                <h2 class="aisle__title">{{ a.label }}</h2>
                <ul class="aisle__list">
                    <ShoppingRow
                        v-for="it in a.items"
                        :key="it.id"
                        :item="it"
                        :checked="false"
                        :open="open.includes(it.id)"
                        @toggle="toggle(it.id)"
                        @flip="flip(it.id)"
                    />
                </ul>
            </section>

            <section
                v-if="optTodo.length"
                class="aisle aisle--opt"
                aria-label="Opcjonalne"
            >
                <h2 class="aisle__title">
                    <span>Opcjonalne</span>
                    <span class="aisle__aside">jeśli chcesz</span>
                </h2>
                <ul class="aisle__list">
                    <ShoppingRow
                        v-for="it in optTodo"
                        :key="it.id"
                        :item="it"
                        :checked="false"
                        :open="open.includes(it.id)"
                        @toggle="toggle(it.id)"
                        @flip="flip(it.id)"
                    />
                </ul>
            </section>
        </div>

        <section
            v-if="done.length"
            class="done"
            aria-label="Odhaczone"
        >
            <h2 class="done__head">
                <button
                    type="button"
                    class="done__toggle"
                    :aria-expanded="doneOpen"
                    aria-controls="shop-done"
                    @click="doneOpen = !doneOpen"
                >
                    <span class="aisle__title aisle__title--plain">
                        Odhaczone · {{ done.length }}
                    </span>
                    <svg
                        viewBox="0 0 12 12"
                        fill="none"
                        aria-hidden="true"
                        :class="{ 'is-open': doneOpen }"
                    >
                        <path
                            d="M2.5 4.5L6 8l3.5-3.5"
                            stroke="currentColor"
                            stroke-width="1.6"
                            stroke-linecap="round"
                            stroke-linejoin="round"
                        />
                    </svg>
                </button>
            </h2>
            <ul
                v-show="doneOpen"
                id="shop-done"
                class="aisle__list"
            >
                <ShoppingRow
                    v-for="it in done"
                    :key="it.id"
                    :item="it"
                    checked
                    :open="false"
                    @toggle="toggle(it.id)"
                />
            </ul>
            <button
                type="button"
                class="done__clear"
                @click="clear"
            >
                Odznacz wszystko
            </button>
        </section>

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
.done__head {
    font: inherit;
}
.done__toggle {
    width: 100%;
    display: flex;
    justify-content: space-between;
    align-items: center;
    background: none;
    border: none;
    border-bottom: 1px solid var(--rule);
    padding: 0 0 6px;
    color: var(--ink-muted);
    cursor: pointer;
}
.aisle__title--plain {
    padding: 0;
    border: none;
    color: var(--ink-muted);
}
.done__toggle svg {
    width: 11px;
    height: 11px;
    transition: transform var(--dur);
}
.done__toggle svg.is-open {
    transform: rotate(180deg);
}
.done__clear {
    margin-top: 10px;
    background: none;
    border: none;
    padding: 0;
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
    .shop__hint {
        margin-top: 10px;
    }
    .shop__cols {
        gap: 20px;
    }
}

/* wydruk A4: dwie kolumny wg działów, tylko nieodhaczone, stopka z pominiętymi */
@include m.print {
    .shop {
        max-width: none;
        padding: 18px 0 0;
        gap: 16px;
    }
    .shop__hint,
    .done {
        display: none;
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
