<script lang="ts" setup>
import { computed } from 'vue';

import QtyText from '@/components/QtyText.vue';
import { forText, type IShopItem, sumText } from '@/shared/utils/shoppingList';

// Pozycja listy zakupów: odhaczenie, nazwa + „do: …” (grupy przepisu), suma.
// Odhaczona zostaje na swoim miejscu — przekreślona i przygaszona (16.11).
const props = defineProps<{ item: IShopItem; checked: boolean }>();
const emit = defineEmits<{ toggle: [] }>();

const forGroups = computed(() => forText(props.item));
</script>

<template>
    <li
        class="srow"
        :class="{ 'is-done': checked }"
    >
        <button
            type="button"
            role="checkbox"
            class="srow__check"
            :aria-checked="checked"
            @click="emit('toggle')"
        >
            <span
                class="srow__box"
                aria-hidden="true"
            >
                <svg
                    v-if="checked"
                    viewBox="0 0 10 10"
                >
                    <path
                        d="M2 5.2L4.2 7.4L8 3.4"
                        stroke="currentColor"
                        stroke-width="1.8"
                        stroke-linecap="round"
                    />
                </svg>
            </span>
            <span class="srow__text">
                <span class="srow__name">{{ item.name }}</span>
                <span
                    v-if="forGroups"
                    class="srow__for"
                >
                    do: {{ forGroups }}
                </span>
            </span>
            <span class="srow__sum"><QtyText :text="sumText(item.sum)" /></span>
        </button>
    </li>
</template>

<style lang="scss" scoped>
@use '@/scss/mixins/media' as m;

.srow {
    border-bottom: 1px solid var(--rule);
}
.srow__check {
    width: 100%;
    display: grid;
    grid-template-columns: 18px minmax(0, 1fr) auto;
    gap: 12px;
    /* pole i suma przy pierwszej linii nazwy (pod nią „do: …"); pojedyncza linia wyśrodkowana w min-height */
    align-items: start;
    align-content: center;
    min-height: 44px;
    padding: 8px 0;
    background: transparent;
    border: none;
    text-align: left;
    color: var(--ink);
    cursor: pointer;
}
.srow__box {
    width: 16px;
    height: 16px;
    border-radius: var(--r-sm);
    border: 1.5px solid var(--rule-strong);
    display: inline-flex;
    align-items: center;
    justify-content: center;
    color: var(--on-accent);
    margin-top: 2px;
}
.srow__box svg {
    width: 9px;
    height: 9px;
}
.is-done .srow__box {
    background: var(--accent);
    border-color: var(--accent);
}
.srow__text {
    min-width: 0;
}
.srow__name {
    display: block;
    font-size: 15px;
    line-height: 1.35;
}
.srow__for {
    display: block;
    margin-top: 1px;
    font-size: 13px;
    color: var(--ink-muted);
}
.srow__sum {
    max-width: 220px;
    font-family: var(--font-mono);
    font-size: 14px;
    font-weight: 600;
    line-height: 20px;
    text-align: right;
}
.is-done .srow__name,
.is-done .srow__sum {
    color: var(--ink-muted);
    text-decoration: line-through;
}

/* telefon w sklepie: większe wiersze i pole odhaczenia */
@media (max-width: 900px) {
    .srow__check {
        grid-template-columns: 24px minmax(0, 1fr) auto;
        min-height: 56px;
        padding: 10px 0;
    }
    .srow__box {
        width: 22px;
        height: 22px;
        border-radius: var(--r-md);
        margin-top: 0.5px;
    }
    .srow__box svg {
        width: 12px;
        height: 12px;
    }
    .srow__name {
        font-size: 17px;
    }
    .srow__sum {
        max-width: 150px;
        font-size: 15px;
        line-height: 23px;
    }
}

/* wydruk: pusty kwadrat do odhaczenia długopisem, „do: …” drobnym drukiem */
@include m.print {
    .srow {
        break-inside: avoid;
    }
    .srow.is-done {
        display: none;
    }
    .srow__check {
        grid-template-columns: 11px minmax(0, 1fr) auto;
        min-height: 0;
        padding: 5px 0;
        gap: 8px;
        align-items: baseline;
    }
    .srow__box {
        width: 11px;
        height: 11px;
        border: 1pt solid #000;
        border-radius: 2px;
        margin-top: 0;
    }
    .srow__name {
        font-size: 12.5px;
    }
    .srow__for {
        font-size: 11px;
        color: #444;
    }
    .srow__sum {
        max-width: none;
        font-size: 12px;
    }
}
</style>
