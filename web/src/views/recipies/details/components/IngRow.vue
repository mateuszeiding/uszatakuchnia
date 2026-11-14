<script lang="ts" setup>
import QtyText from '@/components/QtyText.vue';
import type { RecipeIngredientDto } from '@/data/dtos/recipe/RecipeIngredientDto';

// Wiersz składnika z odhaczaniem (panel przepisu i tryb gotowania).
// Przekreślenie po odhaczeniu tylko na ilości i nazwie — notatka i „też: …” zostają czytelne.
defineProps<{
    ing: RecipeIngredientDto;
    amount: { qty: string; unit: string };
    checked: boolean;
    /** Inne grupy z tym samym składnikiem („też: ciasto, kruszonka”). */
    also: string[];
    /** Większe wiersze (tryb gotowania). */
    roomy?: boolean;
}>();
const emit = defineEmits<{ toggle: [] }>();
</script>

<template>
    <li>
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
                <span
                    v-if="also.length"
                    class="ing-item__also"
                >
                    też: {{ also.join(', ') }}
                </span>
            </span>
        </button>
    </li>
</template>

<style scoped>
.ing-item {
    width: 100%;
    display: grid;
    grid-template-columns: 18px 92px minmax(0, 1fr);
    gap: 10px;
    align-items: baseline;
    background: transparent;
    border: none;
    border-bottom: 1px solid var(--rule);
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
.ing-item__note,
.ing-item__also {
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

@media (max-width: 900px) {
    .ing-item {
        padding: 11px 0;
    }
}
</style>
