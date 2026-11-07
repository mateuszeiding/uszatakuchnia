<script lang="ts" setup>
import { computed } from 'vue';

import type { IngredientDto } from '@/data/dtos/ingredients/IngredientDto';
import { type NutrientKey, NUTRIENTS } from '@/shared/utils/nutrition';

// Komunikat nad tabelą makro — zależy od stanu: podskładniki × liczba własnych wartości.
const props = defineProps<{
    children: IngredientDto[];
    /** Wartości z własnymi danymi (wypełniony wiersz). */
    own: NutrientKey[];
    name: string;
    /** Baza linków do podskładników, np. /admin/ingredients. */
    linkBase: string;
}>();

const hasChildren = computed(() => props.children.length > 0);
const ownCount = computed(() => props.own.length);
const lacking = computed(() =>
    NUTRIENTS.filter((n) => !props.own.includes(n.key))
        .map((n) => n.label.toLowerCase())
        .join(', ')
);
const state = computed(() => {
    const all = NUTRIENTS.length;
    if (hasChildren.value) {
        if (ownCount.value === 0) return 'inherit';
        return ownCount.value < all ? 'mixed' : 'own';
    }
    return ownCount.value > 0 && ownCount.value < all ? 'partial' : null;
});
</script>

<template>
    <div
        v-if="state === 'inherit' || state === 'mixed'"
        class="note note--info"
    >
        <svg
            class="note__icon"
            width="13"
            height="13"
            viewBox="0 0 14 14"
            fill="none"
            aria-hidden="true"
        >
            <rect
                x="1.5"
                y="1.5"
                width="4"
                height="3"
                rx="1"
                stroke="currentColor"
                stroke-width="1.2"
            />
            <rect
                x="8.5"
                y="5.5"
                width="4"
                height="3"
                rx="1"
                stroke="currentColor"
                stroke-width="1.2"
            />
            <rect
                x="8.5"
                y="10"
                width="4"
                height="3"
                rx="1"
                stroke="currentColor"
                stroke-width="1.2"
            />
            <path
                d="M3.5 4.5v7h5M3.5 7h5"
                stroke="currentColor"
                stroke-width="1.2"
            />
        </svg>
        <div class="note__text">
            <template v-if="state === 'inherit'">
                <strong>Brak własnych danych — zakres wyliczony z podskładników</strong>
                (najniższe min, najwyższe max):
                <span
                    v-for="c in children"
                    :key="c.id"
                    class="note__child"
                >
                    <RouterLink :to="`${linkBase}/${c.id}`">{{ c.name }}</RouterLink>
                </span>
                <span class="note__sub">
                    Wpisz wartość w wierszu, żeby go nadpisać. Szeroki zakres sugeruje, że w
                    przepisie lepiej wybrać konkretny podskładnik.
                </span>
            </template>
            <template v-else>
                <strong>Część wartości z podskładników:</strong>
                {{ lacking }}. Pozostałe są własne. Podskładniki:
                <span
                    v-for="c in children"
                    :key="c.id"
                    class="note__child"
                >
                    <RouterLink :to="`${linkBase}/${c.id}`">{{ c.name }}</RouterLink>
                </span>
            </template>
        </div>
    </div>

    <p
        v-else-if="state === 'own'"
        class="note note--plain"
    >
        <svg
            class="note__icon"
            width="13"
            height="13"
            viewBox="0 0 14 14"
            fill="none"
            aria-hidden="true"
        >
            <rect
                x="1.5"
                y="1.5"
                width="4"
                height="3"
                rx="1"
                stroke="currentColor"
                stroke-width="1.2"
            />
            <rect
                x="8.5"
                y="5.5"
                width="4"
                height="3"
                rx="1"
                stroke="currentColor"
                stroke-width="1.2"
            />
            <rect
                x="8.5"
                y="10"
                width="4"
                height="3"
                rx="1"
                stroke="currentColor"
                stroke-width="1.2"
            />
            <path
                d="M3.5 4.5v7h5M3.5 7h5"
                stroke="currentColor"
                stroke-width="1.2"
            />
        </svg>
        <span class="note__text">
            Wszystkie wartości są własne — zakres podskładników
            <span class="note__paren">
                <span
                    v-for="c in children"
                    :key="c.id"
                    class="note__child"
                >
                    <RouterLink :to="`${linkBase}/${c.id}`">{{ c.name }}</RouterLink>
                </span>
            </span>
            nie jest używany. Wyczyść wiersz, żeby wrócić do dziedziczenia.
        </span>
    </p>

    <div
        v-else-if="state === 'partial'"
        class="note note--warn"
    >
        <svg
            class="note__icon"
            width="13"
            height="13"
            viewBox="0 0 14 14"
            fill="none"
            aria-hidden="true"
        >
            <circle
                cx="7"
                cy="7"
                r="5.8"
                stroke="currentColor"
                stroke-width="1.3"
            />
            <path
                d="M7 6.3v3.6M7 4.2v.2"
                stroke="currentColor"
                stroke-width="1.5"
                stroke-linecap="round"
            />
        </svg>
        <div class="note__text">
            <strong>Brakuje: {{ lacking }}.</strong>
            Składnik i tak wliczy się do przepisów, ale te wartości będą tam niepełne (dopisek „bez:
            {{ name.trim() || 'nazwa' }}”).
        </div>
    </div>
</template>

<style scoped>
.note {
    display: flex;
    gap: 12px;
    align-items: flex-start;
    margin: 0 0 16px;
    font-size: 13px;
    line-height: 1.5;
}
.note--info,
.note--warn {
    padding: 12px 14px;
    border-radius: var(--r-lg);
    border: 1px solid transparent;
}
.note--info {
    background: var(--info-soft);
    color: var(--ink-dim);
}
.note--info .note__icon {
    color: var(--info);
}
.note--warn {
    background: var(--prep-soft);
    border-color: var(--prep-line);
    color: var(--ink);
}
.note--warn .note__icon {
    color: var(--prep);
}
.note--plain {
    gap: 10px;
    color: var(--ink-dim);
}
.note--plain .note__icon {
    color: var(--ink-muted);
}
.note__icon {
    flex-shrink: 0;
    margin-top: 3px;
}
.note__text strong {
    font-weight: 600;
    color: var(--ink);
}
.note__text a {
    color: var(--ink);
    text-underline-offset: 3px;
}
/* przecinki między podskładnikami i kropka na końcu — bez spacji z formatowania szablonu */
.note__child:not(:last-child)::after {
    content: ', ';
}
.note--info .note__child:last-child::after {
    content: '.';
}
.note__paren::before {
    content: '(';
}
.note__paren::after {
    content: ')';
}
.note__sub {
    display: block;
    margin-top: 4px;
}
</style>
