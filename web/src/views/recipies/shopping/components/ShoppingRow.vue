<script lang="ts" setup>
import { computed } from 'vue';

import ShoppingBreakdown from './ShoppingBreakdown.vue';

import QtyText from '@/components/QtyText.vue';
import { countedParts, type IShopItem, optText, sumText } from '@/shared/utils/shoppingList';

const props = defineProps<{ item: IShopItem; checked: boolean; open: boolean }>();
const emit = defineEmits<{ toggle: []; flip: [] }>();

const parts = computed(() => countedParts(props.item).length);
const multi = computed(() => parts.value > 1);
const bdId = computed(() => `bd-${props.item.id}`);
</script>

<template>
    <li
        class="srow"
        :class="{ 'is-done': checked }"
    >
        <div class="srow__main">
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
                <span class="srow__name">{{ item.name }}</span>
                <span class="srow__sum"><QtyText :text="sumText(item.sum)" /></span>
            </button>
            <button
                v-if="multi && !checked"
                type="button"
                class="srow__parts"
                :class="{ 'is-open': open }"
                :aria-expanded="open"
                :aria-controls="bdId"
                :aria-label="`${item.name}: z czego się składa`"
                @click="emit('flip')"
            >
                {{ parts }}
                <svg
                    viewBox="0 0 12 12"
                    fill="none"
                    aria-hidden="true"
                >
                    <path
                        d="M2.5 4.5L6 8l3.5-3.5"
                        stroke="currentColor"
                        stroke-width="1.6"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                    />
                </svg>
            </button>
            <span
                v-else
                class="srow__spacer"
            />
        </div>
        <div
            v-if="item.optSum && !checked"
            class="srow__extra srow__opt"
        >
            <QtyText :text="optText(item)" />
            (opcjonalnie)
        </div>
        <div
            v-if="multi && !checked"
            :id="bdId"
            class="srow__extra srow__bd"
            :class="{ 'is-open': open }"
        >
            <ShoppingBreakdown :item="item" />
        </div>
    </li>
</template>

<style lang="scss" scoped>
@use '@/scss/mixins/media' as m;

.srow {
    --box: 18px;
    --parts-w: 30px;
    border-bottom: 1px solid var(--rule);
}
.srow__main {
    display: grid;
    grid-template-columns: minmax(0, 1fr) auto;
    align-items: start;
    gap: 8px;
}
.srow__check {
    display: grid;
    grid-template-columns: var(--box) minmax(0, 1fr) auto;
    gap: 12px;
    align-items: center;
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
}
.srow__box svg {
    width: 9px;
    height: 9px;
}
.is-done .srow__box {
    background: var(--accent);
    border-color: var(--accent);
}
.srow__name {
    font-size: 15px;
    line-height: 1.35;
}
.srow__sum {
    max-width: 220px;
    font-family: var(--font-mono);
    font-size: 14px;
    font-weight: 600;
    text-align: right;
}
.is-done .srow__name,
.is-done .srow__sum {
    color: var(--ink-muted);
    text-decoration: line-through;
}
.srow__parts {
    margin-top: 10px;
    width: var(--parts-w);
    height: 26px;
    border-radius: var(--r-md);
    border: 1px solid var(--rule);
    background: var(--bg-alt);
    color: var(--ink-muted);
    display: inline-flex;
    align-items: center;
    justify-content: center;
    gap: 3px;
    font-family: var(--font-mono);
    font-size: 11px;
    font-weight: 600;
    cursor: pointer;
}
.srow__parts svg {
    width: 11px;
    height: 11px;
    transition: transform var(--dur);
}
.srow__parts.is-open {
    background: var(--accent-soft);
    color: var(--accent);
}
.srow__parts.is-open svg {
    transform: rotate(180deg);
}
.srow__spacer {
    width: var(--parts-w);
}
/* dopiski pod nazwą (wyrównane do nazwy, za polem odhaczenia) */
.srow__extra {
    margin: 0 0 8px calc(var(--box) + 12px);
}
.srow__opt {
    margin-top: -4px;
    font-size: 13px;
    color: var(--ink-muted);
}
.srow__bd {
    display: none;
}
.srow__bd.is-open {
    display: block;
}

/* telefon w sklepie: większe wiersze i pole odhaczenia */
@media (max-width: 900px) {
    .srow {
        --box: 24px;
        --parts-w: 36px;
    }
    .srow__check {
        min-height: 56px;
        padding: 10px 0;
    }
    .srow__box {
        width: 22px;
        height: 22px;
        border-radius: var(--r-md);
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
    }
    .srow__parts {
        margin-top: 14px;
        height: 32px;
    }
    .srow__opt {
        margin-top: -6px;
        margin-bottom: 10px;
    }
    .srow__bd {
        margin-bottom: 12px;
    }
}

/* wydruk: pusty kwadrat do odhaczenia długopisem, składowe zawsze widoczne */
@include m.print {
    .srow {
        --box: 11px;
        break-inside: avoid;
    }
    .srow__check {
        min-height: 0;
        padding: 5px 0 2px;
        gap: 8px;
        align-items: baseline;
    }
    .srow__box {
        width: 11px;
        height: 11px;
        border: 1pt solid #000;
        border-radius: 2px;
    }
    .srow__name {
        font-size: 12.5px;
    }
    .srow__sum {
        max-width: none;
        font-size: 12px;
    }
    .srow__parts,
    .srow__spacer {
        display: none;
    }
    .srow__main {
        grid-template-columns: minmax(0, 1fr);
    }
    .srow__extra {
        margin: 0 0 4px 19px;
        font-size: 11px;
    }
    .srow__bd {
        display: block;
    }
}
</style>
