<script lang="ts" setup>
import { computed, useId } from 'vue';

export interface ICautionPrep {
    lead: string | null;
    step: number;
}

const props = withDefaults(
    defineProps<{
        items: string[];
        prep?: ICautionPrep[];
        variant?: 'default' | 'compact' | 'large' | 'print';
    }>(),
    { prep: () => [], variant: 'default' }
);
const emit = defineEmits<{ step: [no: number] }>();

const headingId = useId();
const isPrint = computed(() => props.variant === 'print');
// the printout (shopping list) shows only the notes
const prepRows = computed(() => (isPrint.value ? [] : props.prep));
const visible = computed(() => props.items.length > 0 || prepRows.value.length > 0);
</script>

<template>
    <section
        v-if="visible"
        class="cautions"
        :class="`cautions--${variant}`"
        :aria-labelledby="headingId"
    >
        <h2
            :id="headingId"
            class="cautions__title"
        >
            <svg
                v-if="!isPrint"
                class="cautions__warn"
                viewBox="0 0 16 16"
                fill="none"
                aria-hidden="true"
            >
                <path
                    d="M8 2.2L14.4 13.3H1.6L8 2.2z"
                    stroke="currentColor"
                    stroke-width="1.4"
                    stroke-linejoin="round"
                />
                <path
                    d="M8 6.4v3.3M8 11.6v.1"
                    stroke="currentColor"
                    stroke-width="1.6"
                    stroke-linecap="round"
                />
            </svg>
            Uwagi
        </h2>
        <ul class="cautions__list">
            <li
                v-for="p in prepRows"
                :key="`prep-${p.step}`"
                class="cautions__item"
            >
                <span class="cautions__clock">
                    <svg
                        width="13"
                        height="13"
                        viewBox="0 0 12 12"
                        fill="none"
                        aria-hidden="true"
                    >
                        <circle
                            cx="6"
                            cy="6"
                            r="4.8"
                            stroke="currentColor"
                            stroke-width="1.3"
                        />
                        <path
                            d="M6 3.5v2.7l1.7 1"
                            stroke="currentColor"
                            stroke-width="1.3"
                            stroke-linecap="round"
                        />
                    </svg>
                </span>
                <!-- the period after the link comes from CSS (::after), so no space before it -->
                <span class="cautions__text cautions__text--prep">
                    Przygotuj {{ p.lead ?? 'wcześniej' }}:
                    <button
                        type="button"
                        class="cautions__step"
                        @click="emit('step', p.step)"
                    >
                        krok {{ p.step }}
                    </button>
                </span>
            </li>
            <li
                v-for="(c, i) in items"
                :key="i"
                class="cautions__item"
            >
                <span
                    class="cautions__dot"
                    aria-hidden="true"
                />
                <span class="cautions__text">{{ c }}</span>
            </li>
        </ul>
    </section>
</template>

<style scoped>
.cautions {
    padding: 14px 16px;
    background: var(--prep-soft);
    border: 1px solid var(--prep-line);
    border-radius: var(--r-xl);
}
.cautions__title {
    margin: 0 0 8px;
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 14px;
    font-weight: 600;
    letter-spacing: 0;
    line-height: 1.3;
    color: var(--prep);
}
.cautions__warn {
    width: 15px;
    height: 15px;
    flex-shrink: 0;
}
.cautions__list {
    list-style: none;
    margin: 0;
    padding: 0;
    display: flex;
    flex-direction: column;
    gap: 8px;
}
.cautions__item {
    display: grid;
    grid-template-columns: 16px minmax(0, 1fr);
    gap: 10px;
    align-items: start;
    font-size: 15px;
    line-height: 1.45;
    font-weight: 500;
    color: var(--ink);
}
.cautions__dot {
    width: 6px;
    height: 6px;
    margin: 8px 0 0 5px;
    border-radius: 1px;
    background: var(--prep);
}
.cautions__clock {
    display: inline-flex;
    padding-top: 3px;
    color: var(--prep);
}
.cautions__text {
    text-wrap: pretty;
    overflow-wrap: anywhere;
}
.cautions__text--prep::after {
    content: '.';
}
.cautions__step {
    background: none;
    border: none;
    padding: 0;
    font: inherit;
    color: var(--ink);
    text-decoration: underline;
    text-underline-offset: 3px;
    cursor: pointer;
}

.cautions--compact {
    padding: 12px 14px;
}
.cautions--compact .cautions__list {
    gap: 6px;
}
.cautions--compact .cautions__item {
    font-size: 14px;
}

.cautions--large {
    padding: 16px 18px;
}
.cautions--large .cautions__title {
    font-size: 15px;
}
.cautions--large .cautions__warn {
    width: 17px;
    height: 17px;
}
.cautions--large .cautions__list {
    gap: 10px;
}
.cautions--large .cautions__item {
    font-size: clamp(16px, 2vw, 19px);
}
.cautions--large .cautions__dot {
    margin-top: 10px;
}
.cautions--large .cautions__clock {
    padding-top: 5px;
}

.cautions--print {
    padding: 8px 12px;
    background: none;
    border: 1pt solid #000;
    border-radius: 0;
    color: #000;
    break-inside: avoid;
}
.cautions--print .cautions__title {
    margin: 0 0 4px;
    font-family: var(--font-mono);
    font-size: 10px;
    letter-spacing: 0.6px;
    text-transform: uppercase;
    color: #000;
}
.cautions--print .cautions__list {
    display: block;
    list-style: disc;
    padding-left: 14px;
}
.cautions--print .cautions__item {
    display: list-item;
    font-size: 11.5px;
    font-weight: 400;
    color: #000;
}
.cautions--print .cautions__dot {
    display: none;
}
</style>
