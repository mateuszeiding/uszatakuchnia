<script lang="ts" setup>
import { computed } from 'vue';

import { fmtScale, scaleWarning } from '@/shared/utils/portion';

// "Przeliczone: tortownica Ø 26 cm ×1,17 · Przywróć przepis" (only when it differs from the recipe), like the
// substitutes line. A form also gets the baking time note and, under 0,5 / over 2, a warning (we still scale).
// Pieces get no warning: fewer or more pizzas don't make the dough thinner (owner's decision 22.11).
const props = defineProps<{ text: string; scale: number; form: boolean }>();
const emit = defineEmits<{ restore: [] }>();

const warn = computed(() => (props.form ? scaleWarning(props.scale) : null));
</script>

<template>
    <div
        class="rsc"
        role="status"
    >
        <div class="rsc__main">
            <span class="rsc__what">
                <svg
                    class="rsc__icon"
                    width="13"
                    height="13"
                    viewBox="0 0 16 16"
                    fill="none"
                    aria-hidden="true"
                >
                    <rect
                        x="2"
                        y="2"
                        width="8"
                        height="8"
                        rx="1.5"
                        stroke="currentColor"
                        stroke-width="1.3"
                    />
                    <path
                        d="M8 14h6V8M10.5 11.5L14 8"
                        stroke="currentColor"
                        stroke-width="1.3"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                    />
                </svg>
                <span>
                    Przeliczone: {{ text }}
                    <span class="rsc__scale">{{ fmtScale(scale) }}</span>
                </span>
            </span>
            <button
                type="button"
                class="rsc__restore"
                @click="emit('restore')"
            >
                Przywróć przepis
            </button>
        </div>
        <p
            v-if="form"
            class="rsc__note"
        >
            <svg
                class="rsc__info"
                width="13"
                height="13"
                viewBox="0 0 16 16"
                fill="none"
                aria-hidden="true"
            >
                <circle
                    cx="8"
                    cy="8"
                    r="6.3"
                    stroke="currentColor"
                    stroke-width="1.3"
                />
                <path
                    d="M8 7.2v4M8 4.9v.1"
                    stroke="currentColor"
                    stroke-width="1.6"
                    stroke-linecap="round"
                />
            </svg>
            Czas pieczenia może się zmienić — sprawdź patyczkiem kilka minut wcześniej.
        </p>
        <p
            v-if="warn"
            class="rsc__note rsc__note--warn"
        >
            <svg
                class="rsc__info"
                width="13"
                height="13"
                viewBox="0 0 16 16"
                fill="none"
                aria-hidden="true"
            >
                <circle
                    cx="8"
                    cy="8"
                    r="6.3"
                    stroke="currentColor"
                    stroke-width="1.3"
                />
                <path
                    d="M8 7.2v4M8 4.9v.1"
                    stroke="currentColor"
                    stroke-width="1.6"
                    stroke-linecap="round"
                />
            </svg>
            {{ warn }}
        </p>
    </div>
</template>

<style scoped>
.rsc {
    margin-top: -4px;
    display: flex;
    flex-direction: column;
    gap: 4px;
}
.rsc__main {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 8px;
    font-size: 13px;
    color: var(--ink-dim);
}
.rsc__what {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    min-width: 0;
}
.rsc__icon {
    flex-shrink: 0;
    color: var(--accent);
}
.rsc__scale {
    font-family: var(--font-mono);
    color: var(--ink-muted);
    white-space: nowrap;
}
.rsc__restore {
    padding: 6px 0;
    border: none;
    background: none;
    color: var(--ink-dim);
    font-size: 13px;
    text-decoration: underline;
    text-underline-offset: 3px;
    white-space: nowrap;
    cursor: pointer;
}
.rsc__note {
    margin: 0;
    display: flex;
    gap: 6px;
    font-size: 12.5px;
    line-height: 1.4;
    color: var(--ink-muted);
}
.rsc__info {
    flex-shrink: 0;
    margin-top: 1px;
}
.rsc__note--warn {
    color: var(--prep);
    font-weight: 600;
}
</style>
