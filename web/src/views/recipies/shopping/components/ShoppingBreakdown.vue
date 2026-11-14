<script lang="ts" setup>
import { computed } from 'vue';

import QtyText from '@/components/QtyText.vue';
import { breakdown, type IShopItem, noConversionNote } from '@/shared/utils/shoppingList';

// Składowe pozycji: „grupa · ilość (≈ przeliczenie) · notatka”. Na wydruku zawsze widoczne, w jednej linii.
const props = defineProps<{ item: IShopItem }>();

const rows = computed(() => breakdown(props.item));
const note = computed(() => noConversionNote(props.item));
</script>

<template>
    <ul class="bd">
        <li
            v-for="(r, i) in rows"
            :key="i"
            class="bd__row"
        >
            <span class="bd__group">{{ r.group.toLowerCase() }}</span>
            {{ ' ' }}
            <span class="bd__qty"><QtyText :text="r.text" /></span>
            <em
                v-if="r.optional"
                class="bd__muted"
            >
                (opcj.)
            </em>
            <span
                v-if="r.note"
                class="bd__muted"
            >
                · {{ r.note }}
            </span>
        </li>
        <li
            v-if="note"
            class="bd__why"
        >
            {{ note }}
        </li>
    </ul>
</template>

<style lang="scss" scoped>
@use '@/scss/mixins/media' as m;

.bd {
    list-style: none;
    margin: 6px 0 2px;
    padding: 0 0 0 12px;
    border-left: 2px solid var(--rule);
    display: flex;
    flex-direction: column;
    gap: 3px;
    font-size: 13px;
    color: var(--ink-dim);
}
.bd__group,
.bd__muted {
    color: var(--ink-muted);
}
.bd__qty {
    font-family: var(--font-mono);
    font-size: 12.5px;
}
.bd__why {
    margin-top: 2px;
    font-size: 12px;
    color: var(--ink-muted);
}

@include m.print {
    .bd {
        margin: 2px 0 0;
        padding: 0;
        border: none;
        flex-flow: row wrap;
        gap: 0 10px;
        font-size: 11px;
    }
    .bd__qty {
        font-size: 10.5px;
    }
    .bd__why {
        display: none;
    }
}
</style>
