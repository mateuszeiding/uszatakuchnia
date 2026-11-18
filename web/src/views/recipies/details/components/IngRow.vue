<script lang="ts" setup>
import { computed, nextTick, ref } from 'vue';

import { publicLink } from '../recipeView';

import SubstituteChooser from './SubstituteChooser.vue';

import BookIcon from '@/components/BookIcon.vue';
import QtyText from '@/components/QtyText.vue';
import SwapIcon from '@/components/SwapIcon.vue';
import type { ISubChoice } from '@/composables/useCookingProgress';
import { useNarrow } from '@/composables/useNarrow';
import type {
    IRecipeSubstitute,
    RecipeIngredientDto,
} from '@/data/dtos/recipe/RecipeIngredientDto';
import { optionName } from '@/shared/utils/substitutes';

// Wiersz składnika z odhaczaniem (panel przepisu i tryb gotowania).
// Przekreślenie po odhaczeniu tylko na ilości i nazwie — notatka zostaje czytelna.
// Ikony (⇄ zamienniki, książka — link do przepisu) to RODZEŃSTWO przycisku (nie w nim): klik w ikonę nie odhacza,
// Tab trafia osobno. Wybrany zamiennik: nazwa wiersza = zamiennik, pod nią „zamiast: …” / „pół na pół”.
const props = defineProps<{
    ing: RecipeIngredientDto;
    amount: { qty: string; unit: string };
    /** Połowa ilości — opcja „pół na pół” w wyborze zamiennika. */
    half?: { qty: string; unit: string };
    checked: boolean;
    /** Wybrany zamiennik (null = składnik główny). */
    sub?: IRecipeSubstitute | null;
    allergenLabel?: (code: string) => string;
    /** Większe wiersze (tryb gotowania). */
    roomy?: boolean;
    /** Link otwiera nową kartę (tryb gotowania — gotowanie zostaje). */
    newTab?: boolean;
}>();
const emit = defineEmits<{ toggle: []; sub: [choice: ISubChoice | null] }>();

const link = computed(() => publicLink(props.ing));
const hasSubs = computed(() => props.ing.substitutes.length > 0);
const swapped = computed(() => props.sub ?? null);
const name = computed(() => optionName(props.ing, swapped.value));

const narrow = useNarrow();
const open = ref(false);
const li = ref<HTMLElement>();
const swapBtn = ref<HTMLButtonElement>();

const swapLabel = computed(
    () =>
        `${props.ing.ingredientName}: zamienniki (${props.ing.substitutes.length})` +
        (swapped.value ? `, wybrany: ${name.value}` : '')
);

async function close() {
    open.value = false;
    await nextTick();
    swapBtn.value?.focus();
}
</script>

<template>
    <li
        ref="li"
        class="ing-li"
        :class="{ 'has-icons': link || hasSubs }"
    >
        <button
            type="button"
            role="checkbox"
            class="ing-item"
            :class="{ 'is-done': checked, 'ing-item--roomy': roomy }"
            :aria-checked="checked"
            @click="emit('toggle')"
        >
            <span
                class="ing-item__box"
                aria-hidden="true"
            >
                <svg
                    v-if="checked"
                    width="9"
                    height="9"
                    viewBox="0 0 10 10"
                >
                    <path
                        d="M2 5.2L4.2 7.4L8 3.4"
                        stroke="#fff"
                        stroke-width="1.8"
                        stroke-linecap="round"
                    />
                </svg>
            </span>
            <span class="ing-item__qty">
                <span class="ing-item__strike">
                    <QtyText :text="amount.qty" />
                    {{ amount.unit }}
                </span>
                <span
                    v-if="ing.note"
                    class="ing-item__note"
                >
                    {{ ing.note }}
                </span>
            </span>
            <span class="ing-item__name">
                <span class="ing-item__strike">{{ name }}</span>
                <em
                    v-if="ing.optional"
                    class="ing-item__opt"
                >
                    (opcjonalnie)
                </em>
                <span
                    v-if="swapped"
                    class="ing-item__instead"
                >
                    {{
                        swapped.mix ? 'pół na pół' : `zamiast: ${ing.ingredientName.toLowerCase()}`
                    }}
                </span>
            </span>
        </button>
        <span
            v-if="link || hasSubs"
            class="ing-icons"
        >
            <button
                v-if="hasSubs"
                ref="swapBtn"
                type="button"
                class="ing-icon ing-swap"
                :class="{ 'ing-icon--roomy': roomy, 'is-active': swapped || open }"
                aria-haspopup="dialog"
                :aria-expanded="open"
                :aria-label="swapLabel"
                :title="`Zamienniki: ${ing.substitutes.length}`"
                @click="open ? close() : (open = true)"
            >
                <span class="ing-icon__sq">
                    <SwapIcon />
                </span>
            </button>
            <RouterLink
                v-if="link"
                :to="`/recipes/${link.id}`"
                class="ing-icon ing-link"
                :class="{ 'ing-icon--roomy': roomy }"
                :title="`Przepis: ${link.title}`"
                :aria-label="`Przepis: ${link.title}${newTab ? ' (otwiera się w nowej karcie)' : ''}`"
                :target="newTab ? '_blank' : undefined"
                :rel="newTab ? 'noopener' : undefined"
            >
                <span class="ing-icon__sq">
                    <BookIcon
                        :external="newTab"
                        :size="newTab ? 14 : 15"
                    />
                </span>
            </RouterLink>
        </span>
        <SubstituteChooser
            v-if="open"
            :ing="ing"
            :selected="swapped"
            :amount="amount"
            :half="half ?? amount"
            :allergen-label="allergenLabel ?? ((c: string) => c)"
            :anchor="li ?? null"
            :sheet="narrow"
            @select="emit('sub', $event)"
            @close="close"
        />
    </li>
</template>

<style scoped>
.ing-li {
    border-bottom: 1px solid var(--rule);
}
.ing-li.has-icons {
    display: grid;
    grid-template-columns: minmax(0, 1fr) auto;
    column-gap: 8px;
    align-items: start;
}
.ing-item {
    width: 100%;
    display: grid;
    grid-template-columns: 18px 92px minmax(0, 1fr);
    gap: 10px;
    align-items: baseline;
    background: transparent;
    border: none;
    padding: 8px 0;
    cursor: pointer;
    text-align: left;
    color: var(--ink);
}
.ing-item--roomy {
    padding: 11px 0;
}
.ing-item__box {
    width: 16px;
    height: 16px;
    border-radius: var(--r-sm);
    border: 1.5px solid var(--rule-strong);
    display: inline-flex;
    align-items: center;
    justify-content: center;
    /* przy pierwszej linii tekstu (nazwa 15 px × 1,4 = 21 px), nie na środku wiersza z notatką */
    align-self: start;
    margin-top: 2.5px;
}
.ing-item.is-done .ing-item__box {
    background: var(--accent);
    border-color: var(--accent);
}
.ing-item__qty {
    font-family: var(--font-mono);
    font-size: 13px;
    font-weight: 600;
}
.ing-item__note {
    display: block;
    font-family: var(--font-sans);
    font-size: 12px;
    font-weight: 400;
    color: var(--ink-muted);
}
.ing-item__name {
    font-size: 15px;
    line-height: 1.4;
    color: var(--ink-dim);
}
.ing-item__opt {
    color: var(--ink-muted);
}
.ing-item__instead {
    display: block;
    font-size: 12px;
    color: var(--ink-muted);
}
.ing-item.is-done .ing-item__qty,
.ing-item.is-done .ing-item__name {
    color: var(--ink-muted);
}
.ing-item.is-done .ing-item__strike {
    text-decoration: line-through;
}

/* gołe ikony: wyrównane do pierwszej linii nazwy, pole kliknięcia 28 px (telefon 44 px, bez wzrostu wiersza);
   tło tylko na wewnętrznym kwadracie 28 px */
.ing-icons {
    display: inline-flex;
}
.ing-icon {
    width: 28px;
    height: 28px;
    /* bez ujemnego marginesu z prawej: panel na desktopie przewija się w środku (poziomy pasek) */
    margin: 4px 0 0;
    padding: 0;
    border: none;
    background: transparent;
    color: var(--accent);
    display: inline-flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
}
.ing-icon--roomy {
    margin-top: 7px;
}
.ing-icon__sq {
    width: 28px;
    height: 28px;
    border-radius: var(--r-md);
    display: inline-flex;
    align-items: center;
    justify-content: center;
}
.ing-icon:hover .ing-icon__sq,
.ing-icon.is-active .ing-icon__sq {
    background: var(--accent-soft);
}
.ing-icon:focus-visible {
    outline: none;
}
.ing-icon:focus-visible .ing-icon__sq {
    background: var(--accent-soft);
    outline: 2px solid var(--accent);
    outline-offset: 0;
}

@media (max-width: 900px) {
    .ing-item {
        padding: 11px 0;
    }
    .ing-icon {
        width: 44px;
        height: 44px;
        margin: 0 0 -4px;
    }
    .ing-icons {
        margin-right: -12px;
    }
    /* tryb gotowania: lista w modalu się przewija — bez wystawania w prawo */
    .ing-icon--roomy {
        margin-top: 0;
    }
    .ing-li:has(.ing-icon--roomy) .ing-icons {
        margin-right: 0;
    }
}
</style>
