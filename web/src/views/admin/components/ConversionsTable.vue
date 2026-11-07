<script lang="ts" setup>
import { computed } from 'vue';

import { type IConvRow, newConvRow, rowState } from '../ingredientForm';

import NumInput from './NumInput.vue';

import type { IConversion } from '@/data/dtos/ingredients/IngredientDto';
import type { IUnit } from '@/data/dtos/taxonomy/DictionariesDto';
import { fmtGrams } from '@/shared/utils/nutrition';

// Ile gramów waży 1 jednostka składnika. g / kg liczone automatycznie, „do smaku” bez przelicznika.
const props = defineProps<{
    units: IUnit[];
    /** Przeliczniki odziedziczone od składnika nadrzędnego (tylko do odczytu). */
    inherited: IConversion[];
    name: string;
    errors: Record<string, string>;
}>();
const rows = defineModel<IConvRow[]>({ required: true });

const unitName = (id: number | null) => props.units.find((u) => u.id === id)?.name ?? '';
const preview = (r: IConvRow) => {
    const st = rowState(r);
    return 'val' in st ? `1 ${unitName(r.unitId)} ≈ ${fmtGrams(st.val)}` : null;
};
const shadowed = computed(() => new Set(rows.value.map((r) => r.unitId)));
const visibleInherited = computed(() =>
    props.inherited.filter((c) => !shadowed.value.has(c.unitId))
);

function add() {
    const free = props.units.find((u) => !shadowed.value.has(u.id));
    rows.value = [...rows.value, newConvRow(free?.id ?? null)];
}
function remove(key: number) {
    rows.value = rows.value.filter((r) => r.key !== key);
}
</script>

<template>
    <div class="conv">
        <table class="data-table admin-grid">
            <thead>
                <tr>
                    <th
                        scope="col"
                        class="admin-grid__label"
                    >
                        Jednostka
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
                    <th scope="col">Podgląd</th>
                    <th
                        scope="col"
                        class="admin-grid__act"
                    >
                        <span class="sr-only">Usuń</span>
                    </th>
                </tr>
            </thead>
            <tbody>
                <tr
                    v-for="(r, i) in rows"
                    :key="r.key"
                >
                    <td>
                        <label
                            class="sr-only"
                            :for="`conv-${r.key}-unit`"
                        >
                            Jednostka
                        </label>
                        <select
                            :id="`conv-${r.key}-unit`"
                            v-model="r.unitId"
                            class="select"
                        >
                            <option
                                :value="null"
                                disabled
                            >
                                wybierz
                            </option>
                            <option
                                v-for="u in units"
                                :key="u.id"
                                :value="u.id"
                            >
                                1 {{ u.name }}
                            </option>
                        </select>
                    </td>
                    <td>
                        <NumInput
                            :id="`conv-${r.key}-min`"
                            v-model="r.min"
                            :label="`1 ${unitName(r.unitId)}, min`"
                            unit="g"
                            :invalid="!!errors[`conversions[${i}]`]"
                            :describedby="
                                errors[`conversions[${i}]`] ? `conv-${r.key}-err` : undefined
                            "
                        />
                    </td>
                    <td>
                        <NumInput
                            :id="`conv-${r.key}-max`"
                            v-model="r.max"
                            :label="`1 ${unitName(r.unitId)}, max`"
                            unit="g"
                            :invalid="!!errors[`conversions[${i}]`]"
                            :describedby="
                                errors[`conversions[${i}]`] ? `conv-${r.key}-err` : undefined
                            "
                        />
                    </td>
                    <td>
                        <span
                            v-if="errors[`conversions[${i}]`]"
                            :id="`conv-${r.key}-err`"
                            class="admin-grid__err"
                            role="alert"
                        >
                            {{ errors[`conversions[${i}]`] }}
                        </span>
                        <span
                            v-else-if="preview(r)"
                            class="admin-grid__preview conv__preview"
                        >
                            {{ preview(r) }}
                        </span>
                        <span
                            v-else
                            class="admin-grid__preview conv__muted"
                        >
                            uzupełnij min
                        </span>
                    </td>
                    <td>
                        <button
                            type="button"
                            class="tbl-btn tbl-btn--danger"
                            :aria-label="`Usuń przelicznik 1 ${unitName(r.unitId)}`"
                            title="Usuń"
                            @click="remove(r.key)"
                        >
                            <svg
                                width="11"
                                height="11"
                                viewBox="0 0 11 11"
                                fill="none"
                                aria-hidden="true"
                            >
                                <path
                                    d="M2 2l7 7M9 2L2 9"
                                    stroke="currentColor"
                                    stroke-width="1.5"
                                    stroke-linecap="round"
                                />
                            </svg>
                        </button>
                    </td>
                </tr>
                <tr
                    v-for="c in visibleInherited"
                    :key="`inh-${c.unitId}`"
                    class="conv__inherited"
                >
                    <td>1 {{ c.unitName }}</td>
                    <td
                        colspan="2"
                        class="conv__muted"
                    >
                        od składnika nadrzędnego
                    </td>
                    <td>
                        <span class="admin-grid__preview conv__preview">
                            1 {{ c.unitName }} ≈
                            {{ fmtGrams({ min: c.gramsMin, max: c.gramsMax }) }}
                            <span class="badge-attr">odziedziczone</span>
                        </span>
                    </td>
                    <td />
                </tr>
                <tr v-if="!rows.length && !visibleInherited.length">
                    <td
                        colspan="5"
                        class="empty-row"
                    >
                        Brak przeliczników — „{{ name || 'składnik' }}” liczy się w makro tylko w
                        gramach i mililitrach.
                    </td>
                </tr>
            </tbody>
        </table>
        <div>
            <button
                type="button"
                class="btn btn--secondary btn--sm"
                @click="add"
            >
                + Dodaj przelicznik
            </button>
        </div>
    </div>
</template>

<style scoped>
.conv {
    display: flex;
    flex-direction: column;
    gap: 12px;
}
.conv__preview {
    font-family: var(--font-mono);
}
.conv__muted {
    color: var(--ink-muted);
    font-size: 13px;
}
.conv__inherited td {
    padding-top: 14px;
    font-size: 14px;
}
</style>
