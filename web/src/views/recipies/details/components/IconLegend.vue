<script lang="ts" setup>
import { ref } from 'vue';

import BookIcon from '@/components/BookIcon.vue';
import SwapIcon from '@/components/SwapIcon.vue';

// Objaśnienie ikon przy składnikach — osobna karta nad kartą składników. Tylko ikony obecne w przepisie
// (rodzic nie renderuje karty, gdy żadnej nie ma). „Ukryj” zapamiętane na urządzeniu (uk:legend:hidden).
defineProps<{ swap: boolean; link: boolean }>();

const KEY = 'uk:legend:hidden';
const read = () => {
    try {
        return localStorage.getItem(KEY) === '1';
    } catch {
        return false;
    }
};
const hidden = ref(read());

function hide() {
    hidden.value = true;
    try {
        localStorage.setItem(KEY, '1');
    } catch {
        /* tryb prywatny — ukryte do odświeżenia */
    }
}
</script>

<template>
    <aside
        v-if="!hidden"
        class="legend"
        aria-labelledby="h-legend"
    >
        <div class="legend__head">
            <h2
                id="h-legend"
                class="legend__title"
            >
                Ikony przy składnikach
            </h2>
            <button
                type="button"
                class="legend__hide"
                @click="hide"
            >
                Ukryj
            </button>
        </div>
        <ul class="legend__list">
            <li
                v-if="swap"
                class="legend__item"
            >
                <span
                    class="legend__icon"
                    aria-hidden="true"
                >
                    <SwapIcon />
                </span>
                <span class="legend__text">
                    <strong>Można wymienić</strong>
                    — kliknij, żeby wybrać zamiennik. Zmienią się wartości odżywcze, alergeny i
                    lista zakupów.
                </span>
            </li>
            <li
                v-if="link"
                class="legend__item"
            >
                <span
                    class="legend__icon"
                    aria-hidden="true"
                >
                    <BookIcon />
                </span>
                <span class="legend__text">
                    <strong>Przepis na ten składnik</strong>
                    — otwiera przepis, z którego go zrobisz.
                </span>
            </li>
        </ul>
    </aside>
</template>

<style scoped>
.legend {
    background: var(--bg-alt);
    border: 1px solid var(--rule);
    border-radius: var(--r-xl);
    padding: 16px 20px;
}
.legend__head {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 8px;
    margin-bottom: 12px;
}
.legend__title {
    font-family: var(--font-mono);
    font-size: 12px;
    font-weight: 600;
    letter-spacing: 0.5px;
    text-transform: uppercase;
    color: var(--ink);
}
.legend__hide {
    padding: 4px 0;
    border: none;
    background: none;
    color: var(--ink-muted);
    font-size: 13px;
    text-decoration: underline;
    text-underline-offset: 3px;
    cursor: pointer;
}
.legend__list {
    list-style: none;
    margin: 0;
    padding: 0;
    display: flex;
    flex-direction: column;
    gap: 10px;
}
.legend__item {
    display: grid;
    grid-template-columns: 28px minmax(0, 1fr);
    gap: 12px;
    align-items: start;
}
.legend__icon {
    width: 28px;
    height: 28px;
    border-radius: var(--r-md);
    background: var(--accent-soft);
    color: var(--accent);
    display: inline-flex;
    align-items: center;
    justify-content: center;
}
.legend__text {
    padding-top: 4px;
    font-size: 14px;
    line-height: 1.45;
    color: var(--ink-dim);
}
.legend__text strong {
    font-weight: 600;
    color: var(--ink);
}

@media (max-width: 900px) {
    .legend {
        padding: 14px 16px;
    }
    /* cel dotyku 44 px bez wzrostu nagłówka */
    .legend__hide {
        padding: 12px 0 12px 12px;
        margin: -12px 0;
    }
}
</style>
