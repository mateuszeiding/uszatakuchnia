<script lang="ts" setup>
import { type IRangeInput, rowState } from '../ingredientForm';

import NumInput from './NumInput.vue';

import MacroValue from '@/components/MacroValue.vue';
import { type Macros, nf, type NutrientKey, NUTRIENTS, rnd } from '@/shared/utils/nutrition';

// Makro na 100 g: 5 wierszy × min / max. Pusty wiersz pokazuje wartości odziedziczone (kursywą).
const props = defineProps<{
    /** Wartości z podskładników, gdy składnik ma dzieci; null = bez dziedziczenia. */
    inherited: Macros | null;
    errors: Record<string, string>;
}>();
const values = defineModel<Record<NutrientKey, IRangeInput>>({ required: true });

const ghost = (key: NutrientKey, side: 'min' | 'max') => {
    const r = props.inherited?.[key];
    const v = side === 'min' ? r?.min : r?.max;
    return v == null ? null : nf(rnd(v, key, '100g'));
};
const isOwn = (key: NutrientKey) => !('empty' in rowState(values.value[key]));
const shown = (key: NutrientKey) => {
    const st = rowState(values.value[key]);
    if ('val' in st) return st.val;
    if ('empty' in st) return props.inherited?.[key] ?? null;
    return null;
};
const restore = (key: NutrientKey) => {
    values.value[key] = { min: '', max: '' };
};
</script>

<template>
    <table class="data-table admin-grid">
        <thead>
            <tr>
                <th
                    scope="col"
                    class="admin-grid__label"
                >
                    Wartość
                </th>
                <th
                    scope="col"
                    class="admin-grid__num"
                >
                    Min
                </th>
                <th
                    scope="col"
                    class="admin-grid__num"
                >
                    Max
                    <span class="admin-grid__opt">(opcj.)</span>
                </th>
                <th scope="col">Na stronie</th>
                <th
                    v-if="inherited"
                    scope="col"
                    class="admin-grid__act"
                >
                    <span class="sr-only">Akcje</span>
                </th>
            </tr>
        </thead>
        <tbody>
            <tr
                v-for="n in NUTRIENTS"
                :key="n.key"
            >
                <th
                    scope="row"
                    class="admin-grid__row"
                >
                    {{ n.label }}
                </th>
                <td>
                    <NumInput
                        :id="`nut-${n.key}-min`"
                        v-model="values[n.key].min"
                        :label="`${n.label}, min`"
                        :unit="n.unit"
                        :invalid="!!errors[`nutrition.${n.key}`]"
                        :describedby="errors[`nutrition.${n.key}`] ? `nut-${n.key}-err` : undefined"
                        :ghost="!isOwn(n.key) ? ghost(n.key, 'min') : null"
                    />
                </td>
                <td>
                    <NumInput
                        :id="`nut-${n.key}-max`"
                        v-model="values[n.key].max"
                        :label="`${n.label}, max`"
                        :unit="n.unit"
                        :invalid="!!errors[`nutrition.${n.key}`]"
                        :describedby="errors[`nutrition.${n.key}`] ? `nut-${n.key}-err` : undefined"
                        :ghost="!isOwn(n.key) ? ghost(n.key, 'max') : null"
                    />
                </td>
                <td>
                    <span
                        v-if="errors[`nutrition.${n.key}`]"
                        :id="`nut-${n.key}-err`"
                        class="admin-grid__err"
                        role="alert"
                    >
                        {{ errors[`nutrition.${n.key}`] }}
                    </span>
                    <span
                        v-else
                        class="admin-grid__preview"
                    >
                        <MacroValue
                            :value="shown(n.key)"
                            :nutrient="n.key"
                            mode="100g"
                            :muted="!isOwn(n.key)"
                        />
                        <template v-if="inherited">
                            <span
                                v-if="isOwn(n.key)"
                                class="badge-attr admin-grid__own"
                            >
                                własne
                            </span>
                            <span
                                v-else-if="inherited[n.key]"
                                class="badge-attr"
                            >
                                odziedziczone
                            </span>
                        </template>
                    </span>
                </td>
                <td v-if="inherited">
                    <button
                        v-if="isOwn(n.key)"
                        type="button"
                        class="tbl-btn"
                        :aria-label="`${n.label}: przywróć dziedziczenie`"
                        title="Przywróć dziedziczenie"
                        @click="restore(n.key)"
                    >
                        <svg
                            width="13"
                            height="13"
                            viewBox="0 0 14 14"
                            fill="none"
                            aria-hidden="true"
                        >
                            <path
                                d="M4.5 3L2 5.5 4.5 8"
                                stroke="currentColor"
                                stroke-width="1.4"
                                stroke-linecap="round"
                                stroke-linejoin="round"
                            />
                            <path
                                d="M2.5 5.5h6a3 3 0 010 6H6"
                                stroke="currentColor"
                                stroke-width="1.4"
                                stroke-linecap="round"
                            />
                        </svg>
                    </button>
                </td>
            </tr>
        </tbody>
    </table>
</template>
