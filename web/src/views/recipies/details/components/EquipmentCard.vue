<script lang="ts" setup>
import { computed } from 'vue';

import type { IRecipeEquipment } from '@/data/dtos/recipe/RecipeDto';

// Sprzęt potrzebny do przepisu: lista do przejrzenia przed gotowaniem (bez odhaczania).
// Desktop od 4 pozycji dwie kolumny (siatka wierszami). Linie tylko MIĘDZY pozycjami —
// ostatnia pozycja w każdej kolumnie bez linii (index + kolumny >= liczba pozycji).
// compact = sekcja arkusza „Sprzęt i składniki” w trybie gotowania: nagłówek jak grupa składników, jedna kolumna.
const props = defineProps<{ items: IRecipeEquipment[]; compact?: boolean }>();

const twoCols = computed(() => !props.compact && props.items.length >= 4);
const headingId = computed(() => (props.compact ? 'h-equip-cook' : 'h-equip'));
</script>

<template>
    <section
        class="equip"
        :class="{ 'equip--card': !compact, 'equip--compact': compact }"
        :aria-labelledby="headingId"
    >
        <h2
            :id="headingId"
            class="equip__title"
        >
            Sprzęt
        </h2>
        <ul
            class="equip__list"
            :class="{ 'equip__list--two': twoCols }"
        >
            <li
                v-for="(x, i) in items"
                :key="`${x.equipmentId}-${i}`"
                class="equip__item"
                :class="{ 'is-last-1': i + 1 >= items.length, 'is-last-2': i + 2 >= items.length }"
            >
                <span class="equip__name">
                    {{ x.name }}
                    <em
                        v-if="x.optional"
                        class="equip__opt"
                    >
                        (opcjonalnie)
                    </em>
                </span>
                <span
                    v-if="x.spec"
                    class="equip__spec"
                >
                    {{ x.spec }}
                </span>
            </li>
        </ul>
    </section>
</template>

<style scoped>
.equip--card {
    background: var(--bg-alt);
    border: 1px solid var(--rule);
    border-radius: var(--r-xl);
    padding: 20px 24px;
}
.equip--card .equip__title {
    font-size: 22px;
    font-weight: 500;
    letter-spacing: -0.6px;
    margin-bottom: 8px;
}
/* tryb gotowania: nagłówek jak nagłówek grupy składników */
.equip--compact .equip__title {
    font-family: var(--font-mono);
    font-size: 12px;
    font-weight: 600;
    letter-spacing: 0.5px;
    text-transform: uppercase;
    color: var(--ink);
    margin-bottom: 4px;
    padding-bottom: 8px;
    border-bottom: 1px solid var(--rule-strong);
}
.equip__list {
    list-style: none;
    margin: 0;
    padding: 0;
    display: grid;
    grid-template-columns: minmax(0, 1fr);
    column-gap: 28px;
}
.equip__item {
    padding: 9px 0;
    border-bottom: 1px solid var(--rule);
}
.equip__item.is-last-1 {
    border-bottom: none;
}
.equip__list--two {
    grid-template-columns: repeat(2, minmax(0, 1fr));
}
.equip__list--two .equip__item.is-last-2 {
    border-bottom: none;
}
.equip__name {
    display: block;
    font-size: 15px;
    line-height: 1.4;
    color: var(--ink);
}
.equip__opt {
    color: var(--ink-muted);
    font-style: italic;
}
.equip__spec {
    display: block;
    margin-top: 1px;
    font-size: 13px;
    line-height: 1.4;
    color: var(--ink-muted);
}

@media (max-width: 900px) {
    .equip--card {
        padding: 18px 20px;
    }
    .equip--card .equip__title {
        font-size: 20px;
    }
    .equip__list--two {
        grid-template-columns: minmax(0, 1fr);
    }
    .equip__list--two .equip__item.is-last-2 {
        border-bottom: 1px solid var(--rule);
    }
    .equip__list--two .equip__item.is-last-1 {
        border-bottom: none;
    }
}
</style>
