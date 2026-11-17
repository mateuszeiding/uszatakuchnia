<script lang="ts" setup>
import { computed } from 'vue';

import { publicLink } from '../recipeView';

import BookIcon from '@/components/BookIcon.vue';
import QtyText from '@/components/QtyText.vue';
import type { RecipeIngredientDto } from '@/data/dtos/recipe/RecipeIngredientDto';

// Wiersz składnika z odhaczaniem (panel przepisu i tryb gotowania).
// Przekreślenie po odhaczeniu tylko na ilości i nazwie — notatka zostaje czytelna.
// Link do przepisu to RODZEŃSTWO przycisku (nie w nim): klik w link nie odhacza, Tab trafia osobno.
const props = defineProps<{
    ing: RecipeIngredientDto;
    amount: { qty: string; unit: string };
    checked: boolean;
    /** Większe wiersze (tryb gotowania). */
    roomy?: boolean;
    /** Link otwiera nową kartę (tryb gotowania — gotowanie zostaje). */
    newTab?: boolean;
}>();
const link = computed(() => publicLink(props.ing));
const emit = defineEmits<{ toggle: [] }>();
</script>

<template>
    <li
        class="ing-li"
        :class="{ 'has-link': link }"
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
                <span class="ing-item__strike">{{ ing.ingredientName }}</span>
                <em
                    v-if="ing.optional"
                    class="ing-item__opt"
                >
                    (opcjonalnie)
                </em>
            </span>
        </button>
        <RouterLink
            v-if="link"
            :to="`/recipes/${link.id}`"
            class="ing-link"
            :class="{ 'ing-link--roomy': roomy }"
            :title="`Przepis: ${link.title}`"
            :aria-label="`Przepis: ${link.title}${newTab ? ' (otwiera się w nowej karcie)' : ''}`"
            :target="newTab ? '_blank' : undefined"
            :rel="newTab ? 'noopener' : undefined"
        >
            <BookIcon
                :external="newTab"
                :size="newTab ? 14 : 15"
            />
        </RouterLink>
    </li>
</template>

<style scoped>
.ing-li {
    border-bottom: 1px solid var(--rule);
}
.ing-li.has-link {
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
    align-self: center;
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
.ing-item.is-done .ing-item__qty,
.ing-item.is-done .ing-item__name {
    color: var(--ink-muted);
}
.ing-item.is-done .ing-item__strike {
    text-decoration: line-through;
}

/* goła ikona: wyrównana do pierwszej linii nazwy, pole kliknięcia 28 px (telefon 44 px, bez wzrostu wiersza) */
.ing-link {
    width: 28px;
    height: 28px;
    /* bez ujemnego marginesu z prawej: panel na desktopie przewija się w środku (poziomy pasek) */
    margin: 4px 0 0;
    border-radius: var(--r-md);
    color: var(--accent);
    display: inline-flex;
    align-items: center;
    justify-content: center;
}
.ing-link--roomy {
    margin-top: 7px;
}
.ing-link:hover,
.ing-link:focus-visible {
    background: var(--accent-soft);
}

@media (max-width: 900px) {
    .ing-item {
        padding: 11px 0;
    }
    .ing-link {
        width: 44px;
        height: 44px;
        margin: 0 -12px -4px 0;
    }
    /* tryb gotowania: lista w modalu się przewija — bez wystawania w prawo */
    .ing-link--roomy {
        margin-right: 0;
    }
}
</style>
