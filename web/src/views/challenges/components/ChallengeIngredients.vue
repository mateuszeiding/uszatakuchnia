<script lang="ts" setup>
import { computed, useId } from 'vue';

import type { IChallengeIngredient } from '@/data/dtos/challenge/ChallengeDto';

// The challenge "Składniki" card: "Musisz użyć" and "Do dyspozycji" / "Możesz też użyć".
// Headers like a recipe ingredient group header. From 8 items (desktop), 2 columns; lines only between items.
const props = defineProps<{ ingredients: IChallengeIngredient[]; narrow: boolean }>();

const id = useId();
const required = computed(() => props.ingredients.filter((i) => i.required));
const available = computed(() => props.ingredients.filter((i) => !i.required));
const sections = computed(() =>
    [
        { key: 'req', title: 'Musisz użyć', items: required.value, cols: 1 },
        {
            key: 'avail',
            title: required.value.length ? 'Możesz też użyć' : 'Do dyspozycji',
            items: available.value,
            cols: !props.narrow && available.value.length >= 8 ? 2 : 1,
        },
    ].filter((s) => s.items.length)
);
</script>

<template>
    <section
        class="ch-ing"
        :aria-labelledby="id"
    >
        <h2
            :id="id"
            class="ch-ing__title"
        >
            Składniki
        </h2>
        <div
            v-for="s in sections"
            :key="s.key"
        >
            <h3 class="ch-ing__head">{{ s.title }}</h3>
            <ul
                class="ch-ing__list"
                :class="{ 'is-two': s.cols === 2 }"
            >
                <li
                    v-for="(x, i) in s.items"
                    :key="i"
                    class="ch-ing__item"
                    :class="{ 'is-last-row': i + s.cols >= s.items.length }"
                >
                    <span
                        class="ch-ing__dot"
                        :class="{ 'is-req': x.required }"
                        aria-hidden="true"
                    />
                    {{ x.name }}
                </li>
            </ul>
        </div>
    </section>
</template>

<style scoped>
.ch-ing {
    background: var(--bg-alt);
    border: 1px solid var(--rule);
    border-radius: var(--r-xl);
    padding: 20px 24px;
    display: flex;
    flex-direction: column;
    gap: 16px;
}
.ch-ing__title {
    margin: 0;
    font-size: 22px;
    font-weight: 500;
    letter-spacing: -0.6px;
}
/* like a recipe ingredient group header */
.ch-ing__head {
    margin: 0 0 2px;
    padding-bottom: 8px;
    border-bottom: 1px solid var(--rule-strong);
    font-family: var(--font-mono);
    font-size: 12px;
    font-weight: 600;
    letter-spacing: 0.5px;
    text-transform: uppercase;
    color: var(--ink);
}
.ch-ing__list {
    list-style: none;
    margin: 0;
    padding: 0;
    display: grid;
    grid-template-columns: minmax(0, 1fr);
    column-gap: 20px;
}
.ch-ing__list.is-two {
    grid-template-columns: repeat(2, minmax(0, 1fr));
}
.ch-ing__item {
    display: grid;
    grid-template-columns: 10px minmax(0, 1fr);
    gap: 10px;
    align-items: baseline;
    padding: 8px 0;
    border-bottom: 1px solid var(--rule);
    font-size: 15px;
    line-height: 1.4;
    color: var(--ink-dim);
}
.ch-ing__item.is-last-row {
    border-bottom: none;
}
.ch-ing__dot {
    width: 6px;
    height: 6px;
    border-radius: var(--r-pill);
    background: var(--rule-strong);
    transform: translateY(-1px);
}
.ch-ing__dot.is-req {
    background: var(--c-tomato);
}
@media (max-width: 900px) {
    .ch-ing {
        padding: 18px 20px;
    }
}
</style>
