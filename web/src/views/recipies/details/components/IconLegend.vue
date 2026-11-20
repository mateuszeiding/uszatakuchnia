<script lang="ts" setup>
import { ref } from 'vue';

import BookIcon from '@/components/BookIcon.vue';
import SwapIcon from '@/components/SwapIcon.vue';

// Ingredient icon legend, a separate card above the ingredients card. Only icons present in the recipe
// (the parent doesn't render the card when there are none). Collapsible: the header is
// a "Zwiń / Pokaż" button; collapsed, it shows the title and small icons. State on the device: uk:legend:collapsed
// (expanded by default); the old uk:legend:hidden ("Ukryj") = collapsed, and that key is removed.
defineProps<{ swap: boolean; link: boolean }>();

const KEY = 'uk:legend:collapsed';
const OLD_KEY = 'uk:legend:hidden';
function read(): boolean {
    try {
        const old = localStorage.getItem(OLD_KEY);
        if (old != null) {
            localStorage.removeItem(OLD_KEY);
            if (old === '1' || old === 'true') {
                localStorage.setItem(KEY, 'true');
                return true;
            }
        }
        return localStorage.getItem(KEY) === 'true';
    } catch {
        return false;
    }
}
const collapsed = ref(read());

function toggle() {
    collapsed.value = !collapsed.value;
    try {
        localStorage.setItem(KEY, String(collapsed.value));
    } catch {
        /* private mode: the state lasts until reload */
    }
}
</script>

<template>
    <aside
        class="legend"
        aria-labelledby="h-legend"
    >
        <h2
            id="h-legend"
            class="legend__h"
        >
            <button
                type="button"
                class="legend__head"
                :aria-expanded="!collapsed"
                aria-controls="legend-list"
                @click="toggle"
            >
                <span class="legend__title">Ikony przy składnikach</span>
                <span
                    v-if="collapsed"
                    class="legend__mini"
                    aria-hidden="true"
                >
                    <SwapIcon
                        v-if="swap"
                        :size="13"
                    />
                    <BookIcon
                        v-if="link"
                        external
                        :size="12"
                    />
                </span>
                <span class="legend__toggle">
                    {{ collapsed ? 'Pokaż' : 'Zwiń' }}
                    <svg
                        width="12"
                        height="12"
                        viewBox="0 0 12 12"
                        fill="none"
                        aria-hidden="true"
                        class="legend__chev"
                        :class="{ 'is-up': !collapsed }"
                    >
                        <path
                            d="M3 4.5L6 7.5 9 4.5"
                            stroke="currentColor"
                            stroke-width="1.5"
                            stroke-linecap="round"
                            stroke-linejoin="round"
                        />
                    </svg>
                </span>
            </button>
        </h2>
        <ul
            id="legend-list"
            class="legend__list"
            :hidden="collapsed"
        >
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
                    <BookIcon
                        external
                        :size="14"
                    />
                </span>
                <span class="legend__text">
                    <strong>Przepis na ten składnik</strong>
                    — otwiera w nowej karcie przepis, z którego go zrobisz.
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
    padding: 0 20px;
}
.legend__h {
    margin: 0;
}
.legend__head {
    width: 100%;
    min-height: 44px;
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 0;
    border: none;
    background: none;
    color: var(--ink);
    font-family: var(--font-sans);
    text-align: left;
    cursor: pointer;
}
.legend__mini {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    color: var(--accent);
}
.legend__toggle {
    display: inline-flex;
    align-items: center;
    gap: 4px;
    font-size: 13px;
    color: var(--ink-muted);
}
.legend__chev.is-up {
    transform: rotate(180deg);
}
.legend__title {
    flex: 1;
    font-family: var(--font-mono);
    font-size: 12px;
    font-weight: 600;
    letter-spacing: 0.5px;
    text-transform: uppercase;
    color: var(--ink);
}
.legend__list {
    list-style: none;
    margin: 0;
    padding: 2px 0 16px;
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

.legend__list[hidden] {
    display: none;
}

@media (max-width: 900px) {
    .legend {
        padding: 0 16px;
    }
    .legend__head {
        min-height: 48px;
    }
    .legend__list {
        padding-bottom: 14px;
    }
}
</style>
